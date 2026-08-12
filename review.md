# Caveman review

Scope: unstaged diff on `master` — proto custom-metric nesting, `accumulate` flag (pie), realtime UI/API cleanup, overwatch/workers fixes.

## backend/api/utils.go

- L38-39: 🔴 bug: `(*resp.Metrics).Metrics` panics if `resp.Metrics` nil. Guard `resp.Metrics == nil` before dereference.

## proto/service.proto

- L126-132: 🟡 risk: `customMetricData.metrics` moved from repeated→nested message; breaking wire change — redeploy worker + backend together.

## worker/db/metrics.go

## metricraft/app/components/Graph.vue

- L184-191: 🔴 bug: custom branch only `console.log`; no chart factory — custom metrics render blank.
- L185: 🟡 risk: annotates `props.data` as `MetricData` but `GraphGrid` passes `entry.metrics` (array).
- L186-190: 🔵 nit: debug `console.log`s; remove before ship.


- L226: 🔴 bug: `result.value.data !== null` still passes `undefined` from catch. Use `result.value.data != null`.

