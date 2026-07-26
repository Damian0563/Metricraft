# Code review — custom metrics / overwatch / types split

## Backend

- `backend/db/custom_metric.go:L33-42`: 🔴 bug: `PrepareData`/`Delete`/`Edit` omit `tz string` — doesn't satisfy `MetricOrchestrator`. Add `tz` param to match interface or drop the interface.
- `backend/db/custom_metric.go:L33-42`: 🟡 risk: `PrepareData`/`Delete`/`Edit` are no-op stubs. Implement or remove from interface until needed.
- `backend/db/custom_metric.go:L59`: 🟡 risk: PK is full JSON `metric` blob — resubmit same config → duplicate-key error; rename-only change → duplicate row. Use `name` (or UUID) as PK, store JSON in a column.
- `backend/db/custom_metric.go:L59`: 🟡 risk: no duplicate-name check before INSERT. `SELECT`/upsert on `name` or return 409.
- `backend/api/overwatch_views.go:L22-31`: 🟡 risk: no validation on decoded `CustomMetric` (empty name/path/selector). Validate before INSERT.
- `backend/api/overwatch_views.go:L29-31`: 🟡 risk: DB errors swallowed — always 500, no log, no body. Log `err`; return JSON `{ err: ... }` like other handlers.
- `backend/api/overwatch_views.go`: 🟡 risk: no GET/list/delete/edit routes — frontend can't reload saved metrics. Add list endpoint or document as WIP.
- `worker/db/db.go:L51` + backend: 🟡 risk: `custom_metrics` table only created in worker `InitDB` — backend INSERT fails if worker hasn't started. Mirror CREATE in worker path backend uses, or migrate on backend boot.

## Frontend

- `metricraft/app/pages/overwatch.vue:L392-394`: 🔴 bug: `resetForm()` in `finally` runs on failure — user loses filled form. Only reset on success; set `loading = false` in `finally`.
- `metricraft/app/pages/overwatch.vue:L219-226`: 🟡 risk: "Configured metrics" panel shows count only — no list UI. Render `metrics` or remove section until list API exists.
- `metricraft/app/pages/overwatch.vue:L388`: 🟡 risk: `metrics` pushed locally only — refresh loses state, no fetch on mount. Load from backend when list endpoint exists.
- `metricraft/app/pages/overwatch.vue:L387`: 🔵 nit: success text stored in `errorMessage` ref. Use separate `statusMessage` or rename ref.
- `metricraft/app/composables/types/additional.ts:L8` vs `overwatch.vue:L245`: 🔵 nit: `ChartType` includes `'geographic'` but UI options are `line|bar|pie`. Align type with UI or add option.
- `metricraft/nuxt.config.ts:L35`: 🔵 nit: `motion-v/nuxt` module added but unused in templates. Remove dep or use it.

## Repo / hygiene

- `backend/cmd/tmp/main`, `worker/cmd/tmp/main`, `**/build-errors.log`: 🔵 nit: build artifacts in diff — add to `.gitignore`, don't commit.
- `backend/types/special_types.go` / `types.go`: 🔵 nit: `Settings`/`DashboardInitPayload` moved to `special_types.go` — verify all Go imports still compile (they do today).

## Still open (from prior todo — rules/metrics)

Rules are passed into these handlers but intentionally ignored until aggregation is implemented correctly.

| Metric | Function | File |
|--------|----------|------|
| P95 Latency | `GetP95Latency` | `worker/db/metrics.go` |
| Uptime Score | `GetUptimeScore` | `worker/db/metrics.go` |
| Geographic performance | `GetGeographicalPerformance` | `worker/db/metrics.go` |

Post-query `checkUrl` + assignment does not produce valid merged statistics for grouped keys — needs re-aggregation from raw rows or SQL-side grouping.
