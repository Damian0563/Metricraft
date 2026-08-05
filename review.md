# Review — custom metric flat proto shape + workers client fetch

Scope: `backend/api/metric_views.go`, `backend/api/utils.go`, `backend/types/metric_types.go`, `metricraft/app/pages/workers.vue`, `proto/service.proto`, `worker/db/metrics.go`, `worker/rpc/rpc_metrics.go`.

## backend/api/metric_views.go

- `L155-157`: 🔴 bug: on gRPC error writes raw `err.Error()` to `w` but no `return`; goroutine keeps going and final handler also `w.Write`s JSON — corrupt response body.
- `L155-157`: 🔴 bug: concurrent goroutines call `w.Write` on same `ResponseWriter`; not safe — collect errors under `mu`, set status once after `wg.Wait`.
- `L160-165`: 🟡 risk: on error `resp` may be nil; still appends metric with empty `Metrics` instead of skipping or surfacing failure.

## backend/types/metric_types.go

- `L3-8`: 🟡 risk: JSON shape changed from nested `[{data:[{timerange,value}]}]` to flat `[{grouping,value}]`; update any chart consumer before deploy.

## metricraft/app/pages/workers.vue

- `L207-228`: 🟡 risk: dropped `useAsyncData` SSR prefetch for workers/team users; page renders empty until `onMounted` — expect loading flash/regression vs prior behavior.
- `L226-228`: 🟡 risk: `loadWorkers` and `loadTeamUsers` fire in parallel with no loading state; partial failure can show workers without recipients (or vice versa) with no combined error UX.




