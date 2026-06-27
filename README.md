# Metricraft

<div align="center">
  <img src="metricraft/public/logo.svg" alt="Metricraft" width="400" />
</div>

An analytics platform for log observability, focused on visual dashboards and reporting capabilities.

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
│                         Metricraft Stack                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│                            ┌───────────┐                        │
│                            │  Supabase │                        │
│                            │  (Users)  │                        │
│                            └───────────┘                        │
│                                 ◄──►                            │
│   ┌──────────────┐        ┌──────────────┐        ┌───────────┐ │
│   │              │        │              │        │           │ │
│   │    Nuxt 4    │◄──────►│    Go API    │◄──────►│   Redis   │ │
│   │  (Frontend)  │ HTTP   │    Server    │        │   :6379   │ │
│   │              │        │    :8080     │        │           │ │
│   └──────────────┘        └───────┬──────┘        └───────────┘ │
│                                 ◄──►                            │
│                          gRPC / WebSocket                       │
│   ┌──────────────┐        ┌──────────────┐                      │
│   │              │        │              │                      │
│   │  PostgreSQL  │◄───────│  Go Worker   │◄─── User Traffic     │
│   │  (Metrics)   │        │    Proxy     │                      │
│   │              │        │              │                      │
│   └──────────────┘        └──────────────┘                      │
│                                                                 │
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

Metricraft is designed to run from the prebuilt all-in-one Docker image. The image bundles PostgreSQL for logs/metrics, Redis, the Go backend, the Go worker proxy, and the Nuxt frontend under `supervisord`, so deployment only needs a Docker Compose file, the app-specific runtime variables, and a volume for PostgreSQL data.

### Docker Compose Deployment

Create a `.env` file next to your Compose file:

```dotenv
APPNAME=my-app
METRICRAFT_PUBLIC_URL=http://localhost:8080
METRICRAFT_WS_URL=ws://localhost:8080
DEST_PORT=3000
```

Then run the prebuilt image:

```yaml
services:
  metricraft:
    image: damianek952/metricraft:latest
    environment:
      APPNAME: ${APPNAME}
      DEST_PORT: ${DEST_PORT}
      NUXT_PUBLIC_HTTPHOST: ${METRICRAFT_PUBLIC_URL}
      NUXT_PUBLIC_WSSHOST: ${METRICRAFT_WS_URL}
    ports:
      - "8000:8000"
      - "8080:8080"
      - "8081:8081"
    volumes:
      - metricraft-db:/var/lib/postgresql/data

volumes:
  metricraft-db:
```

The PostgreSQL data directory is mounted at `/var/lib/postgresql/data`, so captured logs and metrics survive container recreation. The UI is served on `:8000`, the API and WebSocket server on `:8080`, and the worker proxy ingress on `:8081`.

## Environment Configuration

| Variable | Required | Description |
|----------|----------|-------------|
| `APPNAME` | yes | Application identifier used when initializing and grouping captured metrics. |
| `DEST_PORT` | yes | Port of your upstream application that the worker proxy forwards captured traffic to. |
| `NUXT_PUBLIC_HTTPHOST` | yes | Public HTTP(S) origin the browser uses for API calls, e.g. `https://metrics.example.com`. |
| `NUXT_PUBLIC_WSSHOST` | yes | Public WebSocket origin the browser uses, e.g. `wss://metrics.example.com`. |

All `.env` files are git-ignored by default (`**.env` in `.gitignore`).

### Deployment modes (`MODE`)
Set `MODE=local` in `backend/.env` and `worker/.env` for host-based development. Leave it unset when running the all-in-one image because the image sets `MODE=standalone`.

### File locations (local development)

| Service | File path |
|---------|-----------|
| API Server | `backend/.env` |
| Worker Proxy | `worker/.env` |
| Frontend (Nuxt) | `metricraft/.env` |

### `backend/.env`

| Variable | Required | Description |
|----------|----------|-------------|
| `SECRET` | yes | Shared service token (see above). |
| `MODE` | yes | `local` for host development; Docker images set `standalone` automatically. |
| `DATABASE_USERS` | yes | PostgreSQL connection string for the Supabase user database, e.g. `postgresql://postgres.<project>:<password>@<host>:5432/postgres`. See [`users.md`](users.md) for the required table schema. |
| `DATABASE_LOGS` | yes | PostgreSQL connection string for the metrics/logs database, e.g. `postgresql://postgres:password@localhost:5432/postgres?sslmode=disable`. |
| `GOOGLE_APP_PASSWORD` | optional | SMTP/app password used to send verification emails. Required only if email delivery is enabled. |

Example:

```dotenv
SECRET=replace-with-a-long-random-string
MODE=local
DATABASE_USERS=postgresql://postgres.supabase_pooler_creds.pooler.supabase.com:5432/postgres
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

The frontend reads `SECRET`, `NUXT_PUBLIC_HTTPHOST`, and `NUXT_PUBLIC_WSSHOST` at build/dev time (`nuxt.config.ts`). The public host values are the API and WebSocket origins the **browser** uses through the reverse proxy.

| Variable | Required | Description |
|----------|----------|-------------|
| `SECRET` | yes | Shared service token (must match `backend/.env`); exposed to the client through Nuxt's `runtimeConfig.public`. |
| `NUXT_PUBLIC_HTTPHOST` | yes | Public HTTP(S) backend URL as seen by the browser **through the reverse proxy**, e.g. `http://localhost` or `https://metrics.example.com`. |
| `NUXT_PUBLIC_WSSHOST` | optional | Public WebSocket backend URL as seen by the browser, e.g. `ws://localhost` or `wss://metrics.example.com`. Defaults by replacing the HTTP scheme when unset. |
| `PORT` | optional | Port the Nuxt dev server binds to. Defaults to `8000`. |

Example:

```dotenv
SECRET=replace-with-a-long-random-string
NUXT_PUBLIC_HTTPHOST=http://localhost
NUXT_PUBLIC_WSSHOST=ws://localhost
PORT=8000
```

### Notes & best practices

- **Never commit `.env` files.** Rotate any secret that is accidentally pushed.
- **`NUXT_PUBLIC_HTTPHOST` and `NUXT_PUBLIC_WSSHOST` must match the public API/WebSocket URL the browser can reach**, not an internal Docker hostname.
- Worker ingress can stay on the Docker network (`metricraft:8081`) when your upstream app runs in the same compose stack; publish `:8081` only when traffic enters from outside Docker.
- For local development, run PostgreSQL and Redis yourself (`docker-compose.yml` starts only those services) and point `DATABASE_LOGS` at your local Postgres instance.
- For production, prefer injecting secrets through your orchestrator's secret store rather than committing them to `.env` files.

## Useful Commands

```bash
# Compile proto files
protoc -I=./proto --go_out=proto proto/service.proto
protoc -I=./proto --go_out=proto --go-grpc_out=proto proto/service.proto
```
```
# Run dev container
docker-compose up -d

```
## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.
