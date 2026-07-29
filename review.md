# Caveman review — current git diff

Hot-hours feature wired end-to-end on backend/proto; frontend chart missing; a few landmines in seed JSON and debug prints.

## Findings

`backend/api/metric_views.go:L122`: 🔵 nit: debug `fmt.Println(response, err)` left in Hot-hours case. Delete before merge.

`backend/api/metric_views.go:L127`: 🔵 nit: existing `fmt.Println(err)` on every metric error — same cleanup if you care about prod logs.

`metricraft/app/components/Graph.vue:L139`: 🔴 bug: no `"Hot hours"` branch in `populateChart`. Metric renders empty canvas despite API data. Add `createHotHours` (or reuse unique-visitors chart).

`metricraft/app/composables/charts/charts.ts`: 🔴 bug: no `createHotHours` export in diff. Required companion to `Graph.vue` wiring.

`worker/db/db.go:L37`: 🟡 risk: `"Hot hours"` only in seed INSERT — existing DBs never get the metric in `settings.enabled`. Add migration or merge-on-read for upgrades.

`worker/db/metrics.go:L507-L508`: 🟡 risk: `end`/`startAdjusted` use ad-hoc truncation, not `alignStart`/`rangeEnd` like other metrics. Window may drift from selected timeframe at day boundaries.

`worker/db/metrics.go:L506`: 🔵 nit: `loadLocation(timezone)` while `tz := validTimezone(timezone)` sits unused for loc. Use `loadLocation(tz)` for consistency.

`worker/db/metrics.go:L511-L513`: 🔵 nit: label built twice (init loop + scan loop). Fine, but could index a `labels [24]string` slice once.

`backend/cmd/tmp/main`: 🔵 nit: compiled binary in diff. Drop from commit; build in CI/local only.

`worker/cmd/tmp/main`: 🔵 nit: same — don't commit binary artifacts.

`metricraft/app/pages/invite.vue:L144-L181`: 🟡 risk: `v-else` empty state moved outside `ClientOnly`. SSR/hydration may flash "No users..." before client list mounts.

`metricraft/app/pages/invite.vue:L42`: 🟡 risk: `:key="email"` breaks if duplicate invite emails allowed. Use stable id or `email+index` if dupes possible.

`metricraft/app/pages/workers.vue:L137-L150`: 🔵 nit: `motion.label` + nested checkbox is valid a11y pattern. Good if intentional; confirm no outer `<label>` ancestor.

`proto/service.proto` + generated `*.pb.go`: ✅ RPC surface looks consistent; no issues beyond generated noise.

`worker/rpc/rpc_metrics.go:L53-L55`: ✅ thin pass-through to `db.GetHotHours` — fine.

`worker/db/metrics.go:L517`: ✅ `EXTRACT(HOUR FROM (date AT TIME ZONE $4))` correctly buckets by user tz hour-of-day.

`metricraft/app/components/Settings.vue:L148`: ✅ metric registered in UI settings list.
