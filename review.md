# Review — custom metric float aggregation + grouping filter restore

Scope: `.gitignore`, `backend/api/metric_views.go`, `backend/types/metric_types.go`, `worker/db/metrics.go`, `worker/db/utils.go`. Fixes prior float truncation and NULL-scan issues; grouping filter restored.


## backend/api/metric_views.go

- `L189`: ✅ `float64(point.Value)` — fixes prior int truncation.

## worker/db/metrics.go

- `L576-581`: ✅ `sql.NullFloat64` — fixes NULL aggregate scan panic.
- `L581`: 🟡 risk: ignores `count.Valid`; NULL bucket silently becomes `0`. Skip assign when `!count.Valid` or use `COALESCE` in SQL.
- `L538-540`: ✅ grouping param only bound when `applyRules && len(grouping) > 0`.

## worker/db/utils.go

- `L162-166`: ✅ grouping EXISTS filter restored for `applyRules` metrics.
- `L163-166`: 🔵 nit: unnest/EXISTS block mirrors `blacklistFilterSQL`; extract if a third copy appears.

## .gitignore

- `L12`: ✅ `**/cmd/tmp/*.log` — addresses prior build-artifact nit.
