# Review — custom metric error envelope + workers parallel load

Scope: `backend/api/metric_views.go`, `backend/types/metric_types.go`, `metricraft/app/calls/dashboard.ts`, `metricraft/app/composables/types/additional.ts`, `metricraft/app/composables/types/metrics.ts`, `metricraft/app/components/GraphGrid.vue`, `metricraft/app/pages/overwatch.vue`, `metricraft/app/pages/workers.vue`, `worker/db/metrics.go`.

Fixes prior concurrent `w.Write` / corrupt-body bugs in `CustomMetricFetch`.

## backend/api/metric_views.go

- `L158`: 🟡 risk: error string has no metric name. Prefix with `metric.Name` before append.
- `L160-165`: 🟡 risk: on gRPC error still appends metric with empty `Metrics` via `customMetricDataFromProto(resp)`. Skip failed metrics or include per-metric failure flag.
- `L177`: 🔵 nit: partial failures return HTTP 200 with `errors` array; fine if intentional — document contract.

## backend/types/metric_types.go

- `L10-13`: 🟡 risk: response shape changed from bare `MetricData[]` to `{metrics,errors}`. Confirm no other API consumers before deploy.

## metricraft/app/calls/dashboard.ts

- `L59-60`: 🟡 risk: `errorMessage.value = response.errors.join('')` overwrites standard-metric errors set earlier in same `fetchAllMetrics` pass. Append or accumulate.
- `L63`: 🔵 nit: removed `console.error`; loses dev diagnostics on hard failure.

## metricraft/app/composables/types/additional.ts

- `L23`: 🔵 nit: extra blank line before `CustomMetricResponse`.
- `L29-34`: 🔵 nit: `MetricData` moved here from `metrics.ts`; name collides conceptually with `CustomMetric` — consider `FetchedMetricData` or keep in `metrics.ts`.

## metricraft/app/pages/overwatch.vue

- `L329`: ❓ q: `as CustomMetric[]` default-only change — intentional drive-by unrelated to error envelope?

## metricraft/app/pages/workers.vue

- `L208-214`: 🟡 risk: one `getWorkerUptime` failure rejects entire `fetchWorkers` via `Promise.all`. Use `allSettled` per worker.
- `L217-244`: 🟡 risk: `loading.value = false` not in `finally`; spinner sticks if handler throws.
- `L247-248`: 🟡 risk: page data still loads only in `onMounted`; initial render empty until client mount (SSR flash).

## worker/db/metrics.go

- `L514,L571,L578,L585`: 🟡 risk: `fmt.Errorf(...)` drops root cause. Use `%w` with underlying `err`.
- `L578,L585`: 🔵 nit: Scan error and `res.Err()` share identical message. Differentiate stage in text.
- `L590-593`: 🔵 nit: `GetCustomMetricDataCummulative` still returns bare `err`; wrap consistently with buckets fn.
