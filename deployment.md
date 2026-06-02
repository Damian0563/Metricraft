# Deployment — Docker Self-Hosting

Metricraft ships as a single, self-contained Docker image. One container bundles the
entire stack — PostgreSQL, Redis, the Go API backend, the Go worker proxy, and the Nuxt
frontend — supervised by `supervisord`. Self-hosters build the image once with their own
build arguments, push it to a registry, and run it next to their application. They do
**not** need to provide their own database or cache.

## Architecture

```
┌──────────────────────────────────────────────────────────────────────┐
│                    metricraft (single container)                       │
│                         MODE=standalone                                │
│                                                                        │
│   supervisord (PID 1)                                                  │
│   │                                                                    │
│   ├─ postgresql        :5432   ← logs/metrics DB (volume-backed)       │
│   ├─ redis             :6379   ← session / token cache                 │
│   ├─ backend           :8080   ← REST API + WebSocket                  │
│   │     ├── gRPC client ───────────► worker :50051                     │
│   │     ├── redis  ───────────────► 127.0.0.1:6379                     │
│   │     └── DATABASE_LOGS ────────► 127.0.0.1:5432                      │
│   ├─ worker            :8081   ← reverse proxy + gRPC server :50051     │
│   │     ├── DATABASE_LOGS ────────► 127.0.0.1:5432                      │
│   │     └── WebSocket ────────────► 127.0.0.1:8080/ws/workers           │
│   └─ frontend          :8000   ← Nuxt server (Nitro)                   │
│                                                                        │
│   All inter-process traffic stays on localhost.                        │
└──────────────────────────────────────────────────────────────────────┘
        │ 8000 (UI)        │ 8080 (API/WS)        │ 8081 (proxy ingress)
        ▼                  ▼                      ▼
     browser            browser              your app's traffic
```

External dependencies that remain **outside** the container:

- **Supabase / user-auth database** (`DATABASE_USERS`) — your external Postgres for accounts.
- **The end user's browser** — reaches the frontend (`:8000`) and the backend (`:8080`) via `DOMAIN`.
- **Your upstream application** — the worker proxy forwards captured traffic to it (`DEST_PORT`).

## Process supervision

`supervisord` (see `docker/supervisord.conf`) is PID 1 and manages five long-running
programs with `autorestart=true`. Start priorities ensure data services come up first:

| Priority | Program | Command | Listens |
|----------|---------|---------|---------|
| 10 | `postgres` | `docker-entrypoint.sh postgres` | `5432` |
| 10 | `redis` | `redis-server --save "" --appendonly no` | `6379` |
| 20 | `backend` | `/usr/local/bin/backend` | `8080` |
| 20 | `worker` | `/usr/local/bin/worker` | `8081`, gRPC `50051` |
| 30 | `frontend` | `node /app/web/.output/server/index.mjs` | `8000` |

PostgreSQL is initialized on first boot by the official `postgres` entrypoint (using
`POSTGRES_USER`/`POSTGRES_PASSWORD`/`POSTGRES_DB`). The worker's `InitDB` waits ~15s, then
creates the `settings` and `logs` tables if absent — no manual migration step is required.

## Internal networking (`MODE=standalone`)

The image sets `MODE=standalone`, which wires every service to `localhost`:

- backend → redis: `127.0.0.1:6379`
- backend → worker (gRPC): `127.0.0.1:50051`
- backend & worker → logs DB: `127.0.0.1:5432`
- worker → backend (WebSocket): `127.0.0.1:8080/ws/workers`

Only `local` (host development) and `standalone` (this image) modes exist in the codebase.

## Image build (multi-stage)

The root `Dockerfile` has three stages:

1. **`gobuild`** (`golang:1.26-alpine`) — builds `backend` and `worker` as static binaries
   (`CGO_ENABLED=0`) using the Go workspace (`go.work` over `backend`, `worker`, `proto`).
2. **`webbuild`** (`node:20-alpine`) — `npm ci` + `npm run build` for the Nuxt frontend.
   `SECRET`, `APPNAME`, `DOMAIN`, `MODE` are baked at this stage.
3. **runtime** (`postgres:16-alpine` + `redis`, `nodejs`, `supervisor`) — copies the
   binaries and the frontend `.output`, bakes configuration, declares the data volume, and
   launches `supervisord`.

### Build arguments

| Build arg | Required | Description |
|-----------|----------|-------------|
| `SECRET` | yes | Shared bearer token for service-to-service authorization. |
| `APPNAME` | yes | Identifier of the application/tenant the deployment serves. |
| `DATABASE_USERS` | yes | Connection string for your external Supabase user/auth database. Use a least-privilege role and `?sslmode=require`. |
| `DOMAIN` | yes | Public base URL of the backend as reached from the end user's browser, e.g. `http://your-host:8080` or `https://metrics.example.com` (behind a reverse proxy). |

`DATABASE_LOGS` is **not** a build arg — it is fixed inside the image to the bundled
PostgreSQL (`postgresql://postgres:password@127.0.0.1:5432/postgres?sslmode=disable`).

### Build & push

```bash
docker build \
  --build-arg SECRET=your-secret \
  --build-arg APPNAME=your-app-name \
  --build-arg DATABASE_USERS=your-supabase-url \
  --build-arg DOMAIN=http://your-host:8080 \
  -t your-username/metricraft:latest \
  .

docker push your-username/metricraft:latest
```

## Runtime configuration

| Variable | Set by | Description |
|----------|--------|-------------|
| `DEST_PORT` | self-hoster (runtime) | Port of your upstream application the worker proxy forwards captured traffic to. |

| Exposed port | Purpose |
|--------------|---------|
| `8000` | Frontend (UI) |
| `8080` | Backend REST API + WebSocket |
| `8081` | Worker reverse proxy — route the traffic you want to measure here |

| Volume | Purpose |
|--------|---------|
| `/var/lib/postgresql/data` | Persists the logs/metrics database across container recreation |

## Running alongside your application

```yaml
services:
  metricraft:
    image: your-username/metricraft:latest
    environment:
      DEST_PORT: "3000"            # port of YOUR upstream app
    ports:
      - "80:8000"                  # UI
      - "8080:8080"                # API + WebSocket
      - "8081:8081"                # proxy ingress
    volumes:
      - metricraft-db:/var/lib/postgresql/data

  # your own application service(s) ...

volumes:
  metricraft-db:
```

The browser loads the UI from `:8000` and calls the backend at `DOMAIN` (`:8080`). Route
the user traffic you want to observe through the worker proxy on `:8081`, which forwards it
to `DEST_PORT` and streams metrics to the backend.

## Data flow

1. User traffic hits the **worker proxy** (`:8081`), which forwards it to `DEST_PORT` and
   captures method, status, latency, headers, and geo data.
2. The worker streams each captured request to the **backend** over a WebSocket
   (`/ws/workers`) and persists it to the **logs DB**.
3. The **backend** pushes live updates to the **frontend** over a WebSocket
   (`/ws/visitors`) and serves dashboard queries from the logs DB.
4. Account/auth operations use the external **Supabase** database (`DATABASE_USERS`).

## Notes & limitations

- **Single-container trade-off:** bundling five processes simplifies distribution but means
  the whole stack shares one lifecycle; there is no per-service scaling or rolling restart.
- **Secrets in the image:** `SECRET`, `DATABASE_USERS`, and the baked frontend `SECRET` live
  in image layers. Build per-deployment and push only to a **private** registry, or front
  the database with a least-privilege role and network restrictions.
- **Reverse proxy / TLS:** for HTTPS, terminate TLS in front of the container and forward to
  `:8000` (UI) and `:8080` (API/WS, including WebSocket upgrades); set `DOMAIN` to the
  public HTTPS origin.
- **Backups:** back up the `/var/lib/postgresql/data` volume to retain historical metrics.
