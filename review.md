# Review — custom metric float aggregation + uptime rounding

Scope: `worker/db/metrics.go`, `worker/db/utils.go`, `proto/service.proto`, `metricraft/app/composables/charts/charts.ts`, deleted `others/migrate_logs_jsonb.sql`. Aggregation wiring fix looks correct; downstream API/types still truncate floats.

## worker/db/utils.go

- `L157-165`: 🟡 risk: grouping-rule EXISTS filter dropped from custom metrics; `applyRules` metrics now span URLs outside rule groups. Restore filter or document intentional removal.

## worker/db/metrics.go

- `L571-572`: 🟡 risk: `Scan` into plain `float32` errors on SQL NULL (`SUM`/`AVG`/`percentile_cont` over all-null values). Use `sql.NullFloat64` or wrap aggregate in `COALESCE(..., 0)`.

## proto/service.proto

- `L123`: 🟡 risk: `value` wire type `int32`→`float` breaks mixed-version clients during rollout. Deploy worker, backend, and frontend together.

## backend/api/metric_views.go

- `L189`: 🔴 bug: `int(point.Value)` truncates fractional `avg`/`p50`/`p95` after proto change. Emit `float64` (or round explicitly per `valueType`).

## backend/types/metric_types.go

- `L39`: 🔴 bug: `MetricAggregatorData.Data` is `int`; JSON API cannot represent fractional custom-metric values. Change to `float64`.


- `L1`: 🔵 nit: temp build artifact; gitignore `**/cmd/tmp/*.log` and drop from commit.
