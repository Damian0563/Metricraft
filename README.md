# Metricraft (

<div align="center">
  <img src="metricraft/public/logo.svg" alt="Metricraft" width="400" />
	
</div>

An analytics platform for log observability, focused on visual dashboards and reporting capabilities.
[https://hub.docker.com/repository/docker/damianek952/metricraft/general]
<div align="center">
  <img src="images/dashboard.png" alt="Metricraft dashboard" width="900" />
</div>

<table>
  <tr>
    <td width="50%"><img src="images/overwatch.png" alt="Overwatch custom metrics" /><p align="center"><b>Overwatch</b>: custom metrics from your API</p></td>
    <td width="50%"><img src="images/custom-layout.png" alt="Custom dashboard layout" /><p align="center"><b>Custom layouts</b>: drag and drop your dashboard</p></td>
  </tr>
  <tr>
    <td width="50%"><img src="images/workers.png" alt="Workers and uptime" /><p align="center"><b>Workers</b>: uptime and downtime alerts</p></td>
    <td width="50%"><img src="images/rules-grouping.png" alt="Grouping rules" /><p align="center"><b>Rules</b>: group and blacklist routes</p></td>
  </tr>
</table>

<details>
<summary>More screenshots</summary>

| Landing page | Settings |
|:---:|:---:|
| <img src="images/landing.png" alt="Landing page" /> | <img src="images/settings.png" alt="Settings" /> |
| **Sign in** | **Documentation** |
| <img src="images/sign-in.png" alt="Sign in" /> | <img src="images/documentation.png" alt="Documentation" /> |

</details>

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

## Self-Hosting

Metricraft ships as a single Docker image that bundles the frontend, API, worker proxy, PostgreSQL and Redis. Everything, including user accounts, can stay on your own server.

### Requirements

- A Linux server with Docker Engine and the Compose plugin (`docker compose version`)
- A Gmail account with an [app password](https://support.google.com/accounts/answer/185833). Sign-up sends a verification code by email.
- Optional: a domain name, if you want HTTPS

### 1. Create the configuration

```bash
mkdir metricraft && cd metricraft
```

Create `.env` and fill in your values:

```dotenv
# Name of the app you are monitoring
APPNAME=my-app
# Port your app listens on; the worker proxy forwards traffic here
DEST_PORT=3000
# API URL as the browser reaches it (port 8080, or your API domain)
NUXT_PUBLIC_HTTPHOST=http://localhost:8080
# Generate with: openssl rand -hex 32
SECRET=
# Users database. The bundled PostgreSQL works; any PostgreSQL (e.g. Supabase) does too
DATABASE_USERS=postgresql://postgres:password@127.0.0.1:5432/postgres?sslmode=disable
# Gmail account that sends verification, invite and alert emails
GOOGLE_MAIL_ADDRESS=you@gmail.com
GOOGLE_APP_PASSWORD=
```

Create `compose.yaml`:

```yaml
services:
  metricraft:
    image: damianek952/metricraft:latest   # pin a release tag in production
    restart: unless-stopped
    stop_grace_period: 30s
    env_file: .env
    ports:
      - "127.0.0.1:8000:8000"   # dashboard
      - "127.0.0.1:8080:8080"   # API
      - "127.0.0.1:8081:8081"   # worker proxy
    volumes:
      - metricraft-db:/var/lib/postgresql/data

volumes:
  metricraft-db:
```

The ports only listen on `127.0.0.1`, so they are reachable from the server itself or through a reverse proxy (see [step 4](#4-expose-it-on-a-domain)). To expose them directly, remove the `127.0.0.1:` prefix.

### 2. Start Metricraft

```bash
docker compose up -d
docker compose logs -f metricraft
```

### 3. Create the users tables

Metricraft creates its log tables automatically, but you create the account tables once in the database `DATABASE_USERS` points to. For the bundled PostgreSQL:

```bash
docker compose exec -T metricraft psql -U postgres -d postgres <<'SQL'
CREATE TABLE IF NOT EXISTS public.users (
  created_at    timestamptz NOT NULL DEFAULT now(),
  app_name      text,
  mail          text PRIMARY KEY,
  secret        text NOT NULL,
  uuid          uuid,
  allowed_users jsonb DEFAULT '[]'::jsonb,
  pending_users jsonb DEFAULT '[]'::jsonb,
  owner         boolean DEFAULT false
);
CREATE TABLE IF NOT EXISTS public.workers (
  id         bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  created_at timestamptz NOT NULL DEFAULT now(),
  app_name   text,
  workers    json DEFAULT '[]'::json
);
SQL
docker compose restart metricraft
```

If you use Supabase or another external database, run the same SQL there. Column details are in [`others/users.md`](others/users.md) and [`others/workers.md`](others/workers.md).

Open http://localhost:8000 and sign up.

### 4. Expose it on a domain

Put a reverse proxy in front of Metricraft for HTTPS. For example, a `Caddyfile` for [Caddy](https://caddyserver.com), which issues certificates automatically:

```caddyfile
metrics.example.com {
	reverse_proxy 127.0.0.1:8000
}

api.metrics.example.com {
	reverse_proxy 127.0.0.1:8080
}

app.example.com {
	reverse_proxy 127.0.0.1:8081
}
```

Then set `NUXT_PUBLIC_HTTPHOST=https://api.metrics.example.com` in `.env` and run `docker compose up -d`. No rebuild is needed.

### 5. Route your app's traffic through the worker proxy

The worker proxy on `:8081` records each request and forwards it to `http://<request host>:<DEST_PORT>`. In the example above, a request to `app.example.com` is forwarded to `http://app.example.com:3000`, so:

- your app must be reachable on `DEST_PORT` at that hostname from inside the container
- the proxy must pass the original `Host` header without a port (Caddy does by default; for nginx use `proxy_set_header Host $host;`)

### Configuration reference

| Variable | Required | Description |
|----------|----------|-------------|
| `APPNAME` | yes | Name of the monitored app. Metrics and accounts are grouped under it. |
| `SECRET` | yes | Token shared by the frontend and API. It is visible to the browser, so treat it as an API key, not a password. |
| `DATABASE_USERS` | yes | PostgreSQL connection string for the users database. |
| `NUXT_PUBLIC_HTTPHOST` | yes | Public API URL the browser uses. Must not be an internal Docker hostname. |
| `GOOGLE_MAIL_ADDRESS` | yes | Gmail address that sends verification, invite, recovery and alert emails. |
| `GOOGLE_APP_PASSWORD` | yes | App password for `GOOGLE_MAIL_ADDRESS`. |
| `DEST_PORT` | no | Port the worker proxy forwards to. Default: `3000`. |

### Updating

```bash
docker compose pull
docker compose up -d
```

Logs and metrics live in the `metricraft-db` volume and survive updates. The image bundles PostgreSQL 16, so read the release notes before upgrading across a PostgreSQL major version.

### Backup and restore

Restore into a fresh `metricraft-db` volume.

```bash
docker compose exec -T metricraft pg_dump -U postgres postgres > metricraft-backup.sql
docker compose exec -T metricraft psql -U postgres -d postgres < metricraft-backup.sql
```

## Local Development

1. Start PostgreSQL, Redis and pgAdmin: `docker compose up -d` (uses the repo's `docker-compose.yml`).
2. Create `backend/.env` and `worker/.env` with `MODE=local`, `SECRET`, `APPNAME`, `DATABASE_USERS`, `GOOGLE_MAIL_ADDRESS`, `GOOGLE_APP_PASSWORD` and `DATABASE_LOGS=postgresql://postgres:password@localhost:5432/postgres?sslmode=disable`.
3. Create `metricraft/.env` with the same `SECRET` and `NUXT_PUBLIC_HTTPHOST=http://localhost:8080`.
4. Run each service in its own terminal:

```bash
cd backend && go run ./cmd
cd worker && go run ./cmd
cd metricraft && npm install && npm run dev
```

Regenerate the gRPC code after changing `proto/service.proto`:

```bash
protoc -I=./proto --go_out=proto --go-grpc_out=proto proto/service.proto
```

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.
