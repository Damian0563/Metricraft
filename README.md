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

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.
