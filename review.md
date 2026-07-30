# Caveman review — current git diff

Hot-hours chart wired; timeframe tokens extended (`0.5d`, `7t`/`30t`/`365t`); details modal animated. `0.5d` backend parse and details DOM move are the main landmines.

## Findings

`backend/api/metric_views.go:L84`: 🟡 risk: `timeResolution[numFloat]` misses bad tokens (zero value). Default resolution when key absent.


