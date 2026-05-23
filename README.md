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

## .Env configuration



## Useful Commands

```bash
# Compile proto files
protoc -I=./proto --go_out=proto proto/service.proto
```
```
# Run dev container
docker-compose up -d
```
## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.
