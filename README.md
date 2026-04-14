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
| Metrics Store | PostgreSQL | Columnar database for high-volume log storage |
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
| User Database | Supabase |

## Getting Started

### Configuration

Create `.env` files in both `backend/` and `worker/` directories:

**backend/.env**
```
PORT=8080
SECRET=<your-secret-key>
DATABASE_URL=postgres://user:pass@localhost:5432/metricraft
```

**worker/.env**
```
PORT=<incoming-port>
DEST_PORT=<destination-port>
SECRET=<your-secret-key>
```

### Running

1. Start PostgreSQL, Redis, and ClickHouse
2. Start the backend server:
   ```bash
   cd backend && go run main.go
   ```
3. Start the worker proxy:
   ```bash
   cd worker && go run main.go
   ```
4. Start the frontend:
   ```bash
   cd metricraft && npm install && npm run dev
   ```

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.
