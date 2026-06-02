# Metricraft

<div align="center">
  <img src="metricraft/public/logo.svg" alt="Metricraft" width="400" />
</div>

An analytics platform for log observability, focused on visual dashboards and reporting capabilities.

## Who Is This For?

Metriccraft is designed for **small to medium teams** that want to self-host their analytics infrastructure. Whether you're a startup tracking user behavior, a DevOps team monitoring service health, or an engineering org needing visibility into HTTP traffic — Metriccraft gives you full control without relying on third-party SaaS platforms.

Key benefits:
- **Self-hosted**: No data leaves your infrastructure
- **Privacy-first**: Your logs and metrics stay on your servers
- **Scalable for teams**: Built for collaborative analysis across small to medium engineering teams
- **Customizable**: Extend with serverless integrations and gRPC communication between services

## Features

- **Log Observability**: Monitor and track application logs in real-time
- **Visual Dashboards**: Interactive charts and visualizations for data analysis
- **Real-time Metrics**: Live HTTP request/response tracking with performance insights
- **User Authentication**: Secure account management for team collaboration
- **Serverless Mailing Integration**: Send reports and alerts via email using serverless functions
- **gRPC Backend-Worker Communication**: High-performance gRPC communication between backend and worker proxy for efficient metric streaming

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Metricraft Stack                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   ┌──────────────┐        ┌──────────────┐        ┌───────────┐ │
│   │              │        │              │        │           │ │
│   │    Nuxt 4    │◄──────►│   Go API     │◄───gRPC│  Redis    │ │
│   │  (Frontend)  │  HTTP  │   Server     │  Auth  │  (Cache)  │ │
│   │              │        │   :8080      │        │  :6379    │ │
│   └──────────────┘        └──────┬───────┘        └───────────┘ │
│                                  │                              │
│                            WebSocket                            │
│                                  │                              │
│   ┌──────────────┐        ┌──────▼───────┐                     │
│   │              │        │              │                     │
│   │   PostgreSQL │◄───────│  Go Worker   │◄─── User Traffic    │
│   │  (Metrics)   │        │   Proxy      │                     │
│   │              │        │   (gRPC)     │                     │
│   └──────────────┘        └──────────────┘                     │
│                                                                  │
│   ┌──────────────┐                                               │
│   │              │                                               │
│   │  Serverless  │◄── Email Reports & Alerts                      │
│   │  Mail Func   │                                               │
│   └──────────────┘                                               │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Components

| Component | Technology | Description |
|-----------|------------|-------------|
| Frontend | Nuxt 4 + Vue 3 | Server-side rendered web application |
| API Server | Go | REST API and WebSocket server for real-time updates |
| Worker Proxy | Go | Reverse proxy that captures HTTP metrics, communicates with backend via gRPC |
| Serverless Mail | Go Functions | Serverless email service for reports and alerts |
| Metrics Store | PostgreSQL | Database for log storage and analytics |
| Session Cache | Redis | Fast token validation and session management |
| User Database | Supabase | User accounts and authentication |

### Data Flow

1. **Worker Proxy** intercepts incoming HTTP traffic and captures:
   - Request headers and body
   - Response status codes
   - Request duration/latency

2. **Metrics Streaming** via WebSocket to the API server

3. **PostgreSQL** for efficient analytical queries on log data

4. **Real-time Dashboard** updates through Nuxt frontend

## Tech Stack

| Category | Technology |
|----------|------------|
| Frontend Framework | Nuxt 4 |
| UI Framework | Vue 3 |
| Backend Language | Go |
| Metrics Database | PostgreSQL |
| Session Cache | Redis |
| User Database | Supabase (external) |
| Containerization | Docker Compose |

## Getting Started

Metricraft is designed to sit **behind a reverse proxy**. The browser loads the UI and calls the API/WebSocket on the same public hostname; the proxy routes API paths to the backend and everything else to the frontend. Captured app traffic can use a separate hostname or stay on the internal Docker network (`worker:8081`).

### All-in-one image (recommended)

Build the bundled image from the repo root, then run it with a reverse proxy — **do not publish backend or worker ports directly to the host**:

```yaml
services:
  metricraft:
    image: your-username/metricraft:latest
    environment:
      DEST_PORT: "3000"            # port of YOUR upstream app
    expose:
      - "8000"                     # frontend (internal)
      - "8080"                     # backend API + WebSocket (internal)
      - "8081"                     # worker proxy (internal)
    volumes:
      - metricraft-db:/var/lib/postgresql/data

  reverse-proxy:
    image: caddy:2-alpine
    ports:
      - "80:80"                    # UI + API (path-routed)
      - "8081:8081"                # worker proxy ingress (optional)
    volumes:
      - ./docker/Caddyfile.standalone:/etc/caddy/Caddyfile
    depends_on:
      - metricraft

volumes:
  metricraft-db:
```

Set the `DOMAIN` build arg to the public URL the reverse proxy serves (e.g. `http://localhost` when using the Caddyfile above). See [Environment Configuration](#environment-configuration).

## Building and Pushing

### Build arguments

All configuration that is fixed at image-build time is passed via `--build-arg`. Each Dockerfile also declares default values for every `ENV` variable so images run without extra runtime configuration (override placeholders before production use).
| Build arg | Services | Required | Description |
|-----------|----------|----------|-------------|
| `SECRET` | backend, worker, frontend | yes | Shared bearer token for service-to-service authorization; must be identical across all three. |
| `APPNAME` | backend, worker, frontend | yes | Identifier of the application/tenant the deployment serves. |
| `DATABASE_USERS` | backend, all-in-one | yes | Connection string for the Supabase user/auth database. Use a least-privilege role and append `?sslmode=require`. |
| `DOMAIN` | frontend, all-in-one | yes | Public base URL of the backend as reached from the end user's browser through the reverse proxy, e.g. `https://metrics.example.com`. |

### All-in-one image

The repository root contains a single multi-stage `Dockerfile` that bundles **everything** — PostgreSQL (logs/metrics), Redis, the Go backend, the Go worker proxy, and the Nuxt frontend — into one image, managed by `supervisord`. Users do **not** need to provide their own database or cache: they build the image once, push it to their registry, and run a single container alongside their application with reverse-proxy in place.

Internal services communicate over `localhost`. `DATABASE_LOGS`, `GOOGLE_APP_PASSWORD`, and `DEST_PORT` are set in the Dockerfile; override `DEST_PORT` at runtime for your upstream app.

Build it from the repo root (note the trailing `.`):

```bash
docker build \
  --build-arg SECRET=your-secret \
  --build-arg APPNAME=your-app-name \
  --build-arg DATABASE_USERS=your-supabase-url \
  --build-arg DOMAIN=https://metrics.example.com \
  -t your-username/metricraft:latest \
  .

docker push your-username/metricraft:latest
```

Then run it behind a reverse proxy (see [Getting Started](#getting-started)). The container exposes `8000` (frontend) and `8080` (backend), on the Docker network only — publish **one** public port through the proxy on `:80` (UI + API):

```yaml
services:
  metricraft:
    image: your-username/metricraft:latest
    environment:
      DEST_PORT: "3000"
    expose:
      - "8000"
      - "8080"
    volumes:
      - metricraft-db:/var/lib/postgresql/data

  reverse-proxy:
    image: caddy:2-alpine
    ports:
      - "80:80"
    volumes:
      - ./docker/Caddyfile.standalone:/etc/caddy/Caddyfile
    depends_on:
      - metricraft

volumes:
  metricraft-db:
```

The PostgreSQL data directory is a volume (`/var/lib/postgresql/data`) so the logs database survives container recreation.

### Reverse proxy routing

The included Caddyfiles route by path on a single public hostname:

| Path pattern | Upstream | Purpose |
|--------------|----------|---------|
| `/sign`, `/dashboard/*`, `/settings/*`, `/verify/*`, `/ws/*` | backend `:8080` | REST API and WebSocket |
| everything else | frontend `:8000` | Nuxt UI |

Because the browser calls the API on the **same origin** as the UI, set `DOMAIN` to that public URL (e.g. `https://metrics.example.com`), not an internal Docker hostname.

For multi-service development (`docker-compose.dev.yml`), use `docker/Caddyfile` — it targets the `backend` and `metricraft` service names instead.

## Environment Configuration

Configuration is supplied in one of three ways, depending on how you run the stack:

| Setup | Where variables live |
|-------|---------------------|
| **Local development** | Per-service `.env` files (see below), loaded via [`godotenv`](https://github.com/joho/godotenv) for Go services and at Nuxt build/dev time for the frontend. |
| **All-in-one image** (root `Dockerfile`) | Build args baked into the image; override `DEST_PORT` at runtime. Served via `docker/Caddyfile.standalone`. See [deployment.md](deployment.md). |

All `.env` files are git-ignored by default (`**.env` in `.gitignore`).

### Deployment modes (`MODE`)
Set `MODE=local` in `backend/.env` and `worker/.env` for host-based development, for production leave this unset.
### File locations (local development)

| Service | File path |
|---------|-----------|
| API Server | `backend/.env` |
| Worker Proxy | `worker/.env` |
| Frontend (Nuxt) | `metricraft/.env` |

### Shared variables

These values must be **identical** wherever they appear, or authentication and inter-service calls will fail.

| Variable | Used by | Description |
|----------|---------|-------------|
| `SECRET` | backend, worker, frontend | Shared bearer token for service-to-service authorization (sent as the `Authorization` header). Use a long random string. |
### `backend/.env`

| Variable | Required | Description |
|----------|----------|-------------|
| `SECRET` | yes | Shared service token (see above). |
| `MODE` | yes | `local` for host development; Docker images set `standalone` automatically. |
| `DATABASE_USERS` | yes | PostgreSQL connection string for the Supabase user database, e.g. `postgresql://postgres.<project>:<password>@<host>:5432/postgres?sslmode=require`. |
| `DATABASE_LOGS` | yes | PostgreSQL connection string for the metrics/logs database, e.g. `postgresql://postgres:password@localhost:5432/postgres?sslmode=disable`. |
| `GOOGLE_APP_PASSWORD` | optional | SMTP/app password used to send verification emails. Required only if email delivery is enabled. |

Example:

```dotenv
SECRET=replace-with-a-long-random-string
MODE=local
DATABASE_USERS=postgresql://postgres.<project>:<password>@aws-1-eu-west-3.pooler.supabase.com:5432/postgres?sslmode=require
DATABASE_LOGS=postgresql://postgres:password@localhost:5432/postgres?sslmode=disable
GOOGLE_APP_PASSWORD=your-smtp-app-password
```

### `worker/.env`

| Variable | Required | Description |
|----------|----------|-------------|
| `SECRET` | yes | Shared service token (must match `backend/.env`). |
| `APPNAME` | yes | Application identifier (must match the other services). Used when bootstrapping the logs database. |
| `MODE` | yes | `local` for host development; Docker images set `standalone` automatically. |
| `DATABASE_LOGS` | yes | PostgreSQL connection string for writing captured request/response metrics. Must point to the same database as the backend. |
| `DEST_PORT` | optional | Port the worker proxy forwards captured traffic to (your upstream application). Defaults to the port present in the request `Host` header when unset. |

Example:

```dotenv
SECRET=replace-with-a-long-random-string
APPNAME=my-app
MODE=local
DATABASE_LOGS=postgresql://postgres:password@localhost:5432/postgres?sslmode=disable
DEST_PORT=3000
```

### `metricraft/.env`

The frontend reads `SECRET` and `DOMAIN` at build/dev time (`nuxt.config.ts`). `DOMAIN` is the public backend base URL the **browser** uses to reach the API and WebSocket.

| Variable | Required | Description |
|----------|----------|-------------|
| `SECRET` | yes | Shared service token (must match `backend/.env`); exposed to the client through Nuxt's `runtimeConfig.public`. |
| `DOMAIN` | yes | Public backend URL as seen by the browser **through the reverse proxy**, e.g. `http://localhost` or `https://metrics.example.com`. Same origin as the UI when using path-based routing. |
| `PORT` | optional | Port the Nuxt dev server binds to. Defaults to `8000`. |

Example:

```dotenv
SECRET=replace-with-a-long-random-string
DOMAIN=http://localhost
PORT=8000
```

### Docker Compose dev (root `.env`)

When running `docker-compose -f docker-compose.dev.yml up`, create a `.env` file at the repository root. Compose passes values as **build args** and starts a Caddy reverse proxy — application services are not published to the host directly.

| Variable | Required | Used by |
|----------|----------|---------|
| `SECRET` | yes | backend, worker, frontend |
| `APPNAME` | yes | backend, worker, frontend |
| `DATABASE_USERS` | yes | backend |
| `DOMAIN` | yes | frontend — public URL served by the reverse proxy, e.g. `http://localhost` |
| `GOOGLE_APP_PASSWORD` | no | backend |
| `DEST_PORT` | no | worker — defaults to `3000` |

Example root `.env`:

```dotenv
SECRET=replace-with-a-long-random-string
APPNAME=my-app
DOMAIN=http://localhost
DATABASE_USERS=postgresql://postgres.<project>:<password>@<host>:5432/postgres?sslmode=require
GOOGLE_APP_PASSWORD=your-smtp-app-password
DEST_PORT=3000
```

Only the reverse proxy publishes host ports (`80` for UI + API, `8081` for worker ingress). Route captured traffic through `http://localhost:8081` or, from another container on the same network, `http://worker:8081`.

### All-in-one image runtime

When running the bundled image, build-time configuration is already baked in. Override at runtime if needed:

| Variable | Required | Description |
|----------|----------|-------------|
| `DEST_PORT` | yes | Port of your upstream application that the worker proxy forwards captured traffic to. |

### Notes & best practices

- **Never commit `.env` files.** Rotate any secret that is accidentally pushed.
- **`DOMAIN` must match the reverse proxy's public URL** (e.g. `https://metrics.example.com`), not internal Docker hostnames or backend port numbers.
- Put the reverse proxy in front of Metricraft; use `expose` (not `ports`) on Metricraft services and publish only the proxy's `:80`/`:443`.
- Worker ingress can stay on the Docker network (`worker:8081`) when your upstream app runs in the same compose stack; publish `:8081` on the proxy only when traffic enters from outside Docker.
- `DATABASE_USERS` is passed as a build arg and ends up in the backend image layer. Push backend/all-in-one images to a **private** registry, or override the variable in your orchestrator at runtime.
- For local development, run PostgreSQL and Redis yourself (`docker-compose.yml` starts only those services) and point `DATABASE_LOGS` at your local Postgres instance.
- For production, prefer injecting secrets through your orchestrator's secret store rather than committing them to `.env` files.

## Useful Commands

```bash
# Compile proto files
protoc -I=./proto --go_out=proto proto/service.proto
```

```bash
# Run dev stack (Postgres + Redis only)
docker-compose up -d

# Run full stack behind reverse proxy
docker-compose -f docker-compose.dev.yml up -d
```
## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.
