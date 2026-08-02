# Review — custom metric selector filtering

Scope: `worker/db/metrics.go`, `worker/db/utils.go`, `backend/api/metric_views.go` (uncommitted diff).
Solid direction: parameterizing the selector instead of interpolating it keeps this injection-free, and reusing `appendQueryParam` fits the existing style.

## worker/db/metrics.go

- `L561`: 🔴 bug: `%s AS url` aliases the *grouped* url over the raw one, so the outer `blacklistFilterSQL("$7")` (`L565`) now matches blacklist rules against the grouping rule instead of the real url. Every other query in this file (`L41-48`, `L117-123`, `L334-340`) blacklists on raw `url`. Restore `url, %s AS grouped_url`, or drop the grouped select entirely — nothing downstream reads it.

## worker/db/utils.go

- `L167`: 🔴 bug: `payload` is `TEXT` (`db.go:43`), and `::jsonb` aborts the whole query on the first row that isn't valid JSON. Postgres gives no ordering guarantee that the url/method predicates run before the cast, so one bad row kills unrelated metrics. Migrate the column to `jsonb`, or gate on a validity check (`payload ~ '^\s*[{\[]'`) inside a `CASE`.
- `L167`: 🟡 risk: full-table `::jsonb` cast per row, unindexable. With only `idx_url` on `logs` this scans everything in the date range. Add a GIN index on the jsonb column once migrated.
- `L144-147`: 🟡 risk: `indexes` is copied verbatim into the jsonpath, so `items[abc]` produces `$."items"[abc]` and Postgres returns a `42601` syntax error to the user at query time. Validate the bracket contents as digits (or `n to m`) and reject the selector up front.
- `L150`: 🔵 nit: `strings.NewReplacer` allocated per segment inside the loop. Hoist to a package-level `var`.
- `L169-170`: 🔵 nit: `strpos(url, '?'||$n||'=')` misses valueless params (`?debug`), url-encoded keys, and is case-sensitive. If that's acceptable, say so in the comment; otherwise split the query string and compare keys.

## Repo hygiene

- 🟡 risk: `backend/cmd/tmp/main` (29 MB) and `worker/cmd/tmp/main` (24 MB) are tracked and change on every build. Add `**/cmd/tmp/main` to `.gitignore` and `git rm --cached` them.

## Questions

- ❓ q: `customMetricLogicMatch` returns `TRUE` for an empty selector, so a `source: body` metric with no selector counts every matching request. Intended, or should an empty selector be a validation error?
