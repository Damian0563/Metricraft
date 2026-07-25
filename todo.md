# Rules — remaining work

## Statistical metrics (rules not applied)

Rules are passed into these handlers but intentionally ignored until aggregation is implemented correctly.

| Metric | Function | File |
|--------|----------|------|
| P95 Latency | `GetP95Latency` | `worker/db/metrics.go` |
| Uptime Score | `GetUptimeScore` | `worker/db/metrics.go` |
| Geographic performance | `GetGeographicalPerformance` | `worker/db/metrics.go` |

### Why

Post-query `checkUrl` + assignment (`=`) does not produce valid merged statistics:

- **P95 latency** — percentiles cannot be merged by taking the last value per grouped key; grouping collapses multiple URLs into one label but each URL already has its own SQL-computed percentile.
- **Uptime score** — availability is a ratio per URL; grouped keys need combined success/total counts, not overwritten percentages.
- **Geographic performance** — median is computed per `(country, url)` in SQL; filtering rows in Go and assigning `distribution[country] = …` leaves one arbitrary URL’s median per country.

Blacklisting-only (skip excluded URLs, keep raw URL keys) is straightforward. **Grouping** requires re-aggregation from raw rows or a SQL-side approach.

### Possible fixes

1. **Blacklisting only** — filter log rows before computing stats (SQL `WHERE` or scan raw rows in Go, then `percentile_cont` / ratio logic on the filtered set).
2. **Grouping** — aggregate raw `responsetime` / status counts per grouped URL key, then compute percentile or uptime from that combined set.
3. **SQL-side** — express rule prefix match and grouping in queries (harder to maintain, better at scale).
