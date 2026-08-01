# Review — custom metric query + timeframe extraction
Solid direction: dropping the `urlCol` param from the SQL helpers removes a lot of duplication, and moving `ConvertTimeframe` out of `db` into `api` is the right layer.

## worker/db/metrics.go

- `L509-581`: 🔴 bug: `metric.Source`, `Selector`, `Aggregation`, `ValueType` are ignored — query is always `COUNT(*)`. `resolveFieldName` was added for exactly this and is never called. A "sum of `body.amount`" metric now silently returns request counts instead of the old `not implemented` error. Either wire aggregation/selector in, or keep returning an error for `aggregation != "count"`.
- `L540-542,557`: 🔴 bug: `grouped_url` is selected but never used in `SELECT`/`GROUP BY`. The correlated subquery in `groupedUrlSQL` runs per row for nothing. Drop `grouped_url` from the inner select, or group by it if that was the intent.
- `L122-126 (utils.go), via L559`: 🔴 bug: with `applyRules` + grouping rules, the inner `EXISTS(... g.rule)` *excludes* rows not matched by a grouping rule. Everywhere else grouping only relabels (`COALESCE(..., url)`). Custom metrics will undercount. Remove the `EXISTS` clause.
- `L116`: 🟡 risk: `!applyRules` uses `url = $8`, exact match. Logged urls carry query strings (`resolveFieldName` maps `query`→`url`), so `/api/x?p=1` never matches `/api/x`. Strip query before compare or match on prefix.

