# syntax=docker/dockerfile:1

# ---------- Stage 1: build Go binaries (backend + worker) ----------
# Build stages run on the runner's native arch and cross-compile, so multi-arch builds skip emulation
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS gobuild
ARG TARGETOS TARGETARCH
WORKDIR /src
ENV CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH
# Copy only module manifests first so dependency downloads are cached across source changes
COPY go.work go.work.sum ./
COPY proto/go.mod proto/go.sum ./proto/
COPY backend/go.mod backend/go.sum ./backend/
COPY worker/go.mod worker/go.sum ./worker/
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY proto ./proto
COPY backend ./backend
COPY worker ./worker
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o /out/backend ./backend/cmd \
 && go build -trimpath -ldflags="-s -w" -o /out/worker ./worker/cmd

# ---------- Stage 2: build the Nuxt frontend ----------
# Nuxt output is plain JS, so it is identical for every target platform
FROM --platform=$BUILDPLATFORM node:22-alpine AS webbuild
WORKDIR /web
COPY metricraft/package*.json ./
RUN --mount=type=cache,target=/root/.npm npm ci
COPY metricraft ./
# No secrets at build time: runtimeConfig is overridden from NUXT_* env vars when the server starts
ENV MODE=standalone
RUN npm run build

# ---------- Stage 3: single runtime image (postgres + redis + node + go) ----------
FROM postgres:16-alpine

RUN apk add --no-cache redis nodejs supervisor \
 && adduser -D -H -s /sbin/nologin metricraft

# Runtime configuration. Secrets (SECRET, DATABASE_USERS, GOOGLE_APP_PASSWORD)
# must be supplied with `docker run -e` / compose `environment:`, never baked in.
ENV MODE=standalone \
    APPNAME=metricraft \
    DEST_PORT=3000 \
    NUXT_PUBLIC_HTTPHOST=http://localhost:8080 \
    POSTGRES_USER=postgres \
    POSTGRES_PASSWORD=password \
    POSTGRES_DB=postgres \
    PGDATA=/var/lib/postgresql/data \
    DATABASE_LOGS=postgresql://postgres:password@127.0.0.1:5432/postgres?sslmode=disable

COPY --from=gobuild /out/backend /out/worker /usr/local/bin/
COPY --from=webbuild /web/.output /app/web/.output
COPY docker/supervisord.conf /etc/supervisord.conf

WORKDIR /app

EXPOSE 8000 8080 8081

VOLUME ["/var/lib/postgresql/data"]

HEALTHCHECK --interval=30s --timeout=5s --start-period=60s --retries=3 \
  CMD pg_isready -q -h 127.0.0.1 -U postgres && wget -qO /dev/null http://127.0.0.1:8000/ || exit 1

CMD ["supervisord", "-c", "/etc/supervisord.conf"]
