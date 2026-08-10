package db

import (
	"backend/types"
	"context"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

var trailingCommaRE = regexp.MustCompile(`,(\s*[}\]])`)

func unmarshalEnabledMetrics(raw string) (map[string]types.EnabledMetric, bool, error) {
	raw = strings.TrimSpace(raw)
	metrics := make(map[string]types.EnabledMetric)
	if err := json.Unmarshal([]byte(raw), &metrics); err == nil {
		return metrics, false, nil
	}
	cleaned := trailingCommaRE.ReplaceAllString(raw, "$1")
	if err := json.Unmarshal([]byte(cleaned), &metrics); err != nil {
		return nil, false, err
	}
	return metrics, true, nil
}

func repairEnabledMetrics(ctx context.Context, conn *pgxpool.Pool, metrics map[string]types.EnabledMetric) {
	fixed, err := json.Marshal(metrics)
	if err != nil {
		return
	}
	_, _ = conn.Exec(ctx, "UPDATE settings SET enabled = $1 WHERE TRUE", string(fixed))
}

func ChangeLogsRetention(retention int) error {
	ctx := context.Background()
	conn, err := GetLogsPool()
	if err != nil {
		return err
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, "UPDATE settings SET retention = $1 WHERE TRUE", retention)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func GetSettings() (types.Settings, error) {
	var settings types.Settings
	ctx := context.Background()
	conn, err := GetLogsPool()
	if err != nil {
		return types.Settings{}, err
	}
	var enabled string
	err = conn.QueryRow(ctx, "SELECT enabled,retention FROM settings WHERE TRUE").Scan(&enabled, &settings.Retention)
	if err != nil {
		return types.Settings{}, err
	}
	allSettings, repaired, err := unmarshalEnabledMetrics(enabled)
	if err != nil {
		return types.Settings{}, err
	}
	if repaired {
		repairEnabledMetrics(ctx, conn, allSettings)
	}
	settings.Enabled = make(map[string]types.EnabledMetric)
	for name, metric := range allSettings {
		if metric.Enabled {
			settings.Enabled[name] = metric
		}
	}
	return settings, nil
}

func PersistTimeframeSelection(persistChan chan error, metric string, timeframe string) {
	ctx := context.Background()
	conn, err := GetLogsPool()
	if err != nil {
		persistChan <- err
		return
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		persistChan <- err
		return
	}
	defer tx.Rollback(ctx)
	var existing string
	err = tx.QueryRow(ctx, "SELECT enabled FROM settings").Scan(&existing)
	if err != nil {
		persistChan <- err
		return
	}
	mapMetrics, _, err := unmarshalEnabledMetrics(existing)
	if err != nil {
		persistChan <- err
		return
	}
	current := mapMetrics[metric]
	current.Timeframe = timeframe
	mapMetrics[metric] = current
	var stringifiedMetrics []byte
	stringifiedMetrics, err = json.Marshal(mapMetrics)
	if err != nil {
		persistChan <- err
		return
	}
	_, err = tx.Exec(ctx, "UPDATE settings SET enabled = $1", string(stringifiedMetrics))
	if err != nil {
		persistChan <- err
		return
	}
	persistChan <- tx.Commit(ctx)
}

func GetUrls() ([]string, error) {
	ctx := context.Background()
	conn, err := GetLogsPool()
	if err != nil {
		return nil, err
	}
	res, err := conn.Query(ctx, "SELECT DISTINCT url FROM logs ORDER BY url")
	if err != nil {
		return nil, err
	}
	urls := make([]string, 0)
	for res.Next() {
		var url string
		err = res.Scan(&url)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func ChangeMetrics(metrics []types.Metric) error {
	ctx := context.Background()
	conn, err := GetLogsPool()
	if err != nil {
		return err
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var existing string
	err = tx.QueryRow(ctx, "SELECT enabled FROM settings").Scan(&existing)
	if err != nil {
		return err
	}
	mapMetrics, _, err := unmarshalEnabledMetrics(existing)
	if err != nil {
		return err
	}
	for _, metric := range metrics {
		current := mapMetrics[metric.Name]
		current.Enabled = metric.Enabled
		if metric.Timeframe != "" {
			current.Timeframe = metric.Timeframe
		}
		mapMetrics[metric.Name] = current
	}
	var stringifiedMetrics []byte
	stringifiedMetrics, err = json.Marshal(mapMetrics)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "UPDATE settings SET enabled = $1 WHERE TRUE", string(stringifiedMetrics))
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
