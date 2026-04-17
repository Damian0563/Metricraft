# Metricraft

<div align="center">
  <img src="metricraft/public/logo.svg" alt="Metricraft" width="400" />
</div>

An analytics platform for log observability, focused on visual dashboards and reporting capabilities.

## Features

- **Log Observability**: Monitor and track application logs in real-time
- **Visual Dashboards**: Interactive charts and visualizations for data analysis
- **Real-time Metrics**: Live HTTP request/response tracking with performance insights
- **User Authentication**: Secure account management for team collaboration

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Metricraft Stack                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   ┌──────────────┐        ┌──────────────┐        ┌───────────┐ │
│   │              │        │              │        │           │ │
│   │    Nuxt 4    │◄──────►│   Go API     │◄──────►│  Redis    │ │
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
│   │              │        └──────────────┘                     │
│   └──────────────┘                                              │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Components

| Component | Technology | Description |
|-----------|------------|-------------|
| Frontend | Nuxt 4 + Vue 3 | Server-side rendered web application |
| API Server | Go | REST API and WebSocket server for real-time updates |
| Worker Proxy | Go | Reverse proxy that captures HTTP metrics |
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

### Configuration

Create a root `.env` file:

```
SECRET=<your-secret-key>
DATABASE_USERS=<supabase-connection-string>
DATABASE_LOGS=<postgres-connection-string>
```

- `DATABASE_USERS`: Supabase PostgreSQL connection for user accounts
- `DATABASE_LOGS`: Local PostgreSQL connection for log storage

### Running

1. Configure the `.env` file with your secrets
2. Start all services with Docker:
   ```bash
   docker compose up --build
   ```

### Port Binding

Only the frontend is exposed externally. The internal services (backend, worker, PostgreSQL, Redis) are hidden within the Docker network.

To bind the frontend to your desired host port, modify the port mapping in docker-compose.yml:

```yaml
# Change "80:8000" to your desired port
ports:
  - "8080:8000"  # Host port : Container port
```

Then access at http://localhost:8080 (or your chosen port)

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

## Running (Users)

Users only need to bind the port:

```bash
docker run -p 8080:8000 your-username/metricraft:latest
```

Or with docker-compose:

```yaml
services:
  metricraft:
    image: your-username/metricraft:latest
    ports:
      - "8080:8000"
```

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.
