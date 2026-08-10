# Review — realtime removal, custom-metric SQL, workers/overwatch fixes

Scope: unstaged diff on `master` (backend settings/WS cleanup, frontend realtime removal, worker metric SQL, overwatch aggregation, workers error handling, SQL test fixtures).

## Cross-cutting (realtime removal incomplete)

- `metricraft/app/composables/helpers.ts`:L9: 🔴 bug: returns `wsshost` but `config` type dropped it — `nuxi typecheck` fails (TS2353).
- `metricraft/app/calls/settings.ts`:L1-10: 🔵 nit: dead `toggleRealtime` still POSTs removed `/settings/realtime`.
- `metricraft/app/composables/types/views.ts`:L6: 🔵 nit: `dashboardInitPayload.settings.realtime` still typed; backend no longer sends it.
- `metricraft/app/calls/dashboard.ts`:L17: 🔵 nit: error fallback still sets `realtime: false`.

## backend/api/general_views.go

- (removed `ToggleRealtime`) — clean; no findings.

## backend/cmd/main.go

- (removed route + `ws` env) — clean; WS server still starts via `StartWebSocketServer` elsewhere.

## backend/db/settings.go

- `L62`: 🔵 nit: existing DBs still have `realtime` column; no migration/drop. Harmless orphan column.

## backend/types/special_types.go

- (removed `Settings.Realtime`) — clean.

## metricraft/app/components/Dashboard.vue

- (WS/realtime props removed) — clean.

## metricraft/app/components/Graph.vue

- `L136`: 🔵 nit: debug `console.log(data)` in `populateChart` when `props.custom`; remove before ship.

## metricraft/app/components/Settings.vue

- `L12-29`: 🔴 bug: "Real-time updates" label + tooltip remain but toggle/checkbox removed — dead, misleading UI block.

## metricraft/app/pages/dashboard.vue

- (realtime state removed) — clean given backend change.

## metricraft/app/pages/overwatch.vue

- `L360-364`: ✅ fixed prior bug: `allowedAggregations` intersects value-type + chart-type allowlists.
- `L365-369`: 🟡 risk: watch has no `{ immediate: true }`; invalid agg combo on first paint not reset until `valueType`/`chartType` changes.
- `L368`: 🟡 risk: empty `allowedAggregations` sets `aggregation` to `undefined` (`allowed[0]!`); guard with fallback e.g. `'count'`.

## metricraft/app/pages/workers.vue

- `L182-184,L193-195`: 🟡 risk: sequential `if (fetchError)` then `if (teamUsersError)` — second clobbers first dual-failure message.
- `L215`: 🔵 nit: `workersList.value = workers` inside `watch(workersList)` is noop; drop line.
- `L223-224,L227-228`: 🔴 bug: catch returns `{ data: undefined }` but filter keeps fulfilled rows; `WorkerUptimeGraph` gets `undefined` data. Filter `result.value.data != null` or omit failed workers.
- `L229-232`: 🟡 risk: uptime fetch errors overwrite workers/team-users init errors; append or prefer first non-empty.

## metricraft/nuxt.config.ts

- (removed `wsshost` from runtime config) — clean; pair with `helpers.ts` fix.

## metricraft/app/ws/visitors.ts

- (deleted) — clean.

## worker/db/db.go

- `L25`: 🟡 risk: `CREATE TABLE IF NOT EXISTS` drops `realtime` only on fresh DB; existing installs keep old schema — document or migrate.

## worker/db/metrics.go

- `L542-549,L614-621`: 🟡 risk: empty selector leaves `valueExpr` `""` — inner SELECT becomes ` AS value` / ` AS predicate`; invalid SQL for no-selector metrics.
- `L626-640`: ✅ fixes prior cumulative pie bug: inner query now defines `predicate` + `value`, `Scan` uses `sql.NullString`.
- `L572,L643,L651,L658`: 🔵 nit: `fmt.Errorf(...)` drops root `err`; use `%w`.

## others/custom_metrics_man_tests/mailing_body_selector.sql

- (jsonb casts + `@?` operator) — aligns with worker jsonb path; no findings.

## Noise / out of scope

- `backend/cmd/tmp/build-errors.log`, `worker/cmd/tmp/build-errors.log`: 🔵 nit: tmp build logs in diff; exclude from commit.
- `?? others/custom_metrics_man_tests/checkout_status_pie.sql`: untracked; not reviewed.
