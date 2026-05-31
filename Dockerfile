# syntax=docker/dockerfile:1

# ---------- Stage 1: build Go binaries (backend + worker) ----------
FROM golang:1.26-alpine AS gobuild
WORKDIR /src
ENV CGO_ENABLED=0
# Workspace ties backend, worker and proto together
COPY go.work go.work.sum ./
COPY proto ./proto
COPY backend ./backend
COPY worker ./worker
RUN go build -o /out/backend ./backend \
 && go build -o /out/worker ./worker

# ---------- Stage 2: build the Nuxt frontend ----------
FROM node:20-alpine AS webbuild
WORKDIR /web
COPY metricraft/package*.json ./
RUN npm ci
COPY metricraft ./
# Frontend bakes these at build time (read via import.meta.env)
ARG SECRET
ARG APPNAME
ARG DOMAIN
ENV SECRET=$SECRET
ENV APPNAME=$APPNAME
ENV DOMAIN=$DOMAIN
ENV MODE=standalone
RUN npm run build

# ---------- Stage 3: single runtime image (postgres + redis + node + go) ----------
FROM postgres:16-alpine

# Redis, Node.js runtime and the process supervisor
RUN apk add --no-cache redis nodejs supervisor

# Build-time configuration baked into the image
ARG SECRET
ARG APPNAME
ARG ALLOWED_ORIGINS
ARG DATABASE_USERS
ARG DOMAIN

ENV SECRET=$SECRET
ENV APPNAME=$APPNAME
ENV ALLOWED_ORIGINS=$ALLOWED_ORIGINS
ENV DATABASE_USERS=$DATABASE_USERS
ENV DOMAIN=$DOMAIN

# Everything runs in one container and talks over localhost
ENV MODE=standalone
# Internal logs/metrics database (bundled postgres on localhost)
ENV DATABASE_LOGS=postgresql://postgres:password@127.0.0.1:5432/postgres?sslmode=disable
# Bootstrap values for the bundled postgres instance
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=password
ENV POSTGRES_DB=postgres

# Artifacts
COPY --from=gobuild /out/backend /usr/local/bin/backend
COPY --from=gobuild /out/worker /usr/local/bin/worker
COPY --from=webbuild /web/.output /app/web/.output
COPY docker/supervisord.conf /etc/supervisord.conf

WORKDIR /app

# Frontend (8000), backend API/WS (8080), worker proxy (8081)
EXPOSE 8000 8080 8081

# Persist the logs/metrics database
VOLUME ["/var/lib/postgresql/data"]

CMD ["supervisord", "-c", "/etc/supervisord.conf"]
