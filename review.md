# Review — JSONB migration + custom metric aggregation

Scope: `.gitignore`, `worker/db/db.go`, `worker/db/metrics.go`, `worker/db/utils.go`, deleted `**/cmd/tmp/main` binaries. Prior grouped-url/blacklist and TEXT→jsonb cast issues are addressed; aggregation wiring is not.

## worker/db/metrics.go

- `L544`: 🔴 bug: `resolveAggregationType(metric.Aggregation, selectorParam)` passes the bound placeholder (e.g. `$8` = jsonpath string), not a value expression. `SUM($8)` aggregates the path literal, not field data. Build `(jsonb_path_query_first(<col>, $n::jsonpath)#>>'{}')::numeric` (or query-string extract for `source: query`) and pass that into `resolveAggregationType`.
- `L575-580`: 🔴 bug: `Scan` into `int` / assign `int32` — `avg`, `p50`, `p95` return float/numeric. Use `float64` (or scan `sql.NullFloat64`) and respect `value_type` before truncating.


