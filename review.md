# Review — aggregation/chart filters, workers SSR load, cumulative metric SQL

Scope: `metricraft/app/calls/dashboard.ts`, `metricraft/app/pages/overwatch.vue`, `metricraft/app/pages/workers.vue`, `worker/db/metrics.go`, `worker/db/utils.go`.

## metricraft/app/calls/dashboard.ts

- `L32-33`: 🔵 nit: removed `console.error(e)` in `getUrls`; silent `[]` on failure loses dev signal.

## metricraft/app/pages/overwatch.vue

- `L106`: 🟡 risk: filtered `<option>` list can be empty while `v-model="aggregation"` still holds invalid combo (e.g. pie + `p50`).
- `L357`: 🔵 nit: typo `aggregaationTypesPerChartType` (triple `a`); rename before it spreads.
- `L363-367`: 🔴 bug: `watch(valueType)` only checks value-type allowlist, not chart-type; boolean on pie keeps `p50`.
- `L369-373`: 🟡 risk: `watch(chartType)` resets agg but not coordinated with value-type watch order; intersect both allowlists in one place.

## metricraft/app/pages/workers.vue

- `L5`: 🟡 risk: dropped `loading` from `<Spinner>`; no full-page load indicator while uptimes fetch in background.
- `L191`: 🔴 bug: `watch(teamUsersError)` sets `errorMessage` to `''` when error clears — wipes `'Failed to load workers.'` from `fetchError`.
- `L191`: 🟡 risk: team-users error overwrites workers error; lost prior dual-failure combined message.
- `L209-218`: 🟡 risk: failed `getWorkerUptime` rows silently dropped via `allSettled`; no partial-failure notice.
- `L182,L209`: 🔵 nit: `workersList` mirrors `existingWorkers` then diverges on local CRUD; consider single source or refresh `existingWorkers` after mutations.

## worker/db/metrics.go

- `L628-630`: 🔴 bug: inner SELECT has trailing comma (`%s AS value,`) — invalid SQL syntax.
- `L625,L629`: 🔴 bug: outer SELECT references `predicate` but inner never defines it; query fails at runtime for pie charts.
- `L623-636`: 🔴 bug: cumulative query needs inner `%s AS predicate, %s AS value` (like buckets fn), not one column.
- `L593,L639,L647,L654`: 🟡 risk: `fmt.Errorf(...)` drops root `err`. Use `%w`.
- `L653-654`: 🔵 nit: Scan vs `res.Err()` share identical message; tag stage in text.

## worker/db/utils.go

- `L85-86`: 🟡 risk: removed empty-`selectorParam` guard; `jsonb_path_query_first(%s, ::jsonpath)` breaks if selector omitted server-side.
- `L87-88`: 🟡 risk: url/query branch always uses `substring(...)`; old empty-selector path returned `url`/`NULL`.
