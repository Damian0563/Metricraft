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
With docker-compose:

```yaml
services:
  metricraft:
    image: your-username/metricraft:latest
    ports:
      - "8080:8000"
		environment:
			-DEST_PORT=3000
```

<strong>Make sure the port is allowed by CORS policy in your backend.</strong>

## Building and Pushing

To build and push the image to Docker Hub:

```bash
# Build with environment variables
docker build \
  --build-arg SECRET=your-secret \
  --build-arg DATABASE_USERS=your-supabase-url \
  --build-arg DATABASE_LOGS=your-postgres-url \
  -t your-username/metricraft:latest \
  ./metricraft

# Push to Docker Hub
docker push your-username/metricraft:latest
```

## .Env Configuration

Each service in the stack reads its configuration from a local `.env` file (loaded via [`godotenv`](https://github.com/joho/godotenv) for the Go services and Vite/Nuxt for the frontend). Create one file per service at the paths shown below. All `.env` files are git-ignored by default.

### File locations

| Service | File path |
|---------|-----------|
| API Server | `backend/.env` |
| Worker Proxy | `worker/.env` |
| Frontend (Nuxt) | `metricraft/.env` |

### Shared variables

These variables must be **identical** across the services that use them, otherwise authentication and inter-service calls will fail.

| Variable | Used by | Description |
|----------|---------|-------------|
| `SECRET` | backend, worker, metricraft | Shared bearer token used for service-to-service authorization (sent as the `Authorization` header). Use a long random string. |
| `MODE` | backend, worker, metricraft | Deployment mode. One of `local`, `docker`, or `prod`. Controls hostnames used for inter-service communication (e.g. `localhost` vs. `metricraft-backend-1`). |

### `backend/.env`

| Variable | Required | Description |
|----------|----------|-------------|
| `SECRET` | yes | Shared service token (see above). |
| `MODE` | yes | `local` \| `docker` \| `prod`. |
| `DATABASE_USERS` | yes | PostgreSQL connection string for the Supabase user database, e.g. `postgresql://postgres.<project>:<password>@<host>:5432/postgres`. |
| `DATABASE_LOGS` | yes | PostgreSQL connection string for the metrics/logs database, e.g. `postgresql://postgres:password@localhost:5432/postgres?sslmode=disable`. |
| `ALLOWED_ORIGINS` | yes | Comma-separated list of origins permitted by CORS (e.g. `http://localhost:8000,https://app.example.com`). |
| `API_RESEND` | optional | [Resend](https://resend.com) API key used by the serverless mailing integration to send verification emails. Required only if email delivery is enabled. |

Example:

```dotenv
SECRET=replace-with-a-long-random-string
MODE=local
DATABASE_USERS=postgresql://postgres.<project>:<password>@aws-1-eu-west-3.pooler.supabase.com:5432/postgres
DATABASE_LOGS=postgresql://postgres:password@localhost:5432/postgres?sslmode=disable
ALLOWED_ORIGINS=http://localhost:8000
GOOGLE_APP_PASSWORD=re_xxxxxxxxxxxxxxxxxxxxxxxx
```

### `worker/.env`

| Variable | Required | Description |
|----------|----------|-------------|
| `SECRET` | yes | Shared service token (must match `backend/.env`). |
| `MODE` | yes | `local` \| `docker` \| `prod`. |
| `DATABASE_LOGS` | yes | PostgreSQL connection string for writing captured request/response metrics. Must point to the same database as the backend. |
| `DEST_PORT` | optional | Port the worker proxy forwards captured traffic to (your upstream application). Defaults to the port present in the request `Host` header when unset. |

Example:

```dotenv
SECRET=replace-with-a-long-random-string
MODE=local
DATABASE_LOGS=postgresql://postgres:password@localhost:5432/postgres?sslmode=disable
DEST_PORT=3000
```

### `metricraft/.env`

| Variable | Required | Description |
|----------|----------|-------------|
| `SECRET` | yes | Shared service token (must match `backend/.env`); exposed to the client through Nuxt's `runtimeConfig.public`. |
| `MODE` | yes | `local` \| `docker` \| `prod`. Selects the backend host the frontend talks to (`http://localhost:8080`, `http://metricraft-backend-1:8080`, or the production URL). |
| `PORT` | optional | Port the Nuxt dev server binds to. Defaults to `8000`. |

Example:

```dotenv
PORT=8000
SECRET=replace-with-a-long-random-string
MODE=local
```

### Notes & best practices

- **Never commit `.env` files.** They are included in `.gitignore` (`**.env`); rotate any secret that is accidentally pushed.
- When using Docker Compose, the values referenced as `${SECRET}`, `${DATABASE_USERS}`, and `${DATABASE_LOGS}` in `docker-compose.dev.yml` are read from a project-level `.env` file at the repository root or from the shell environment.
- The `MODE` value must be consistent across all three services; mixing `local` and `docker` will cause hostname resolution errors.
- For production, prefer injecting variables through your orchestrator's secret store (Kubernetes secrets, Docker secrets, etc.) rather than baking them into images via `--build-arg`.

## Useful Commands

```bash
# Compile proto files
protoc -I=./proto --go_out=proto proto/service.proto
```
```

```bash
# Run dev container
docker-compose up -d
```
## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.
