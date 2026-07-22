# ---------- Stage 1: build Go binaries (backend + worker) ----------
FROM golang:1.26-alpine AS gobuild
WORKDIR /src
ENV CGO_ENABLED=0
COPY go.work go.work.sum ./
COPY proto ./proto
COPY backend ./backend
COPY worker ./worker
RUN go build -o /out/backend ./backend/cmd \
 && go build -o /out/worker ./worker/cmd

# ---------- Stage 2: build the Nuxt frontend ----------
FROM node:20-alpine AS webbuild
WORKDIR /web
COPY metricraft/package*.json ./
RUN npm ci
COPY metricraft ./
ARG SECRET
ARG APPNAME=metricraft
ARG NUXT_PUBLIC_HTTPHOST=http://localhost
ARG NUXT_PUBLIC_WSSHOST=ws://localhost
ENV SECRET=$SECRET
ENV APPNAME=$APPNAME
ENV NUXT_PUBLIC_HTTPHOST=$NUXT_PUBLIC_HTTPHOST
ENV NUXT_PUBLIC_WSSHOST=$NUXT_PUBLIC_WSSHOST
ENV MODE=standalone
RUN npm run build

# ---------- Stage 3: single runtime image (postgres + redis + node + go) ----------
FROM postgres:16-alpine

RUN apk add --no-cache redis nodejs supervisor

ARG SECRET
ARG DATABASE_USERS
ARG GOOGLE_APP_PASSWORD
ARG APPNAME=metricraft
ARG NUXT_PUBLIC_HTTPHOST=http://localhost
ARG NUXT_PUBLIC_WSSHOST=ws://localhost
ARG DEST_PORT=3000


ENV SECRET=$SECRET
ENV APPNAME=$APPNAME
ENV DATABASE_USERS=$DATABASE_USERS
ENV NUXT_PUBLIC_HTTPHOST=$NUXT_PUBLIC_HTTPHOST
ENV NUXT_PUBLIC_WSSHOST=$NUXT_PUBLIC_WSSHOST
ENV GOOGLE_APP_PASSWORD=$GOOGLE_APP_PASSWORD
ENV MODE=standalone
ENV DATABASE_LOGS=postgresql://postgres:password@127.0.0.1:5432/postgres?sslmode=disable
ENV DEST_PORT=3000
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=password
ENV POSTGRES_DB=postgres
ENV NITRO_PORT=8000
ENV PORT=8000
ENV HOST=0.0.0.0
ENV NITRO_HOST=0.0.0.0

COPY --from=gobuild /out/backend /usr/local/bin/backend
COPY --from=gobuild /out/worker /usr/local/bin/worker
COPY --from=webbuild /web/.output /app/web/.output
COPY docker/supervisord.conf /etc/supervisord.conf

WORKDIR /app

EXPOSE 8000 8080 8081

VOLUME ["/var/lib/postgresql/data"]

CMD ["supervisord", "-c", "/etc/supervisord.conf"]
