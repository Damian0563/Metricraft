# Caveman review — custom metrics dashboard + styling

Custom metrics fetched into `GraphGrid`, styled via `Graph` `custom` prop. Wiring is mostly there; timeframe refresh and chart rendering are the gaps.

## Findings

`GraphGrid.vue:L66-74`: 🔴 bug: `handleTimeframeChange` always calls `fetchMetric` (built-in `/dashboard/fetch`). Custom cards hit wrong endpoint; timeframe change won't work. Branch on `entry.customMetrics` and refetch via custom API (or re-call `fetchCustomMetrics` + merge).

`Graph.vue:L144-191`: 🔴 bug: `populateChart` has no custom-metric branch. Custom cards render styled shell + empty canvas. Add renderer keyed on chart type (or generic line/bar/pie from Overwatch config).

`GraphGrid.vue:L6`: 🟡 risk: `v-for` `:key="entry.name"` collides if custom metric name matches a built-in. Key on `${entry.customMetrics}-${entry.name}` or stable id.
