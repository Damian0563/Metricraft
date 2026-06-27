package api

import (
	"backend/types"
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"os"
)

func ChangeRealtime(enabled bool) error {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, "UPDATE settings SET realtime = $1 WHERE TRUE", enabled)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func ChangeLogsRetention(retention int) error {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
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
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return types.Settings{}, err
	}
	defer conn.Close(ctx)
	var enabled string
	err = conn.QueryRow(ctx, "SELECT realtime,enabled,retention FROM settings WHERE TRUE").Scan(&settings.Realtime, &enabled, &settings.Retention)
	if err != nil {
		return types.Settings{}, err
	}
	allSettings := make(map[string]types.EnabledMetric)
	err = json.Unmarshal([]byte(enabled), &allSettings)
	if err != nil {
		return types.Settings{}, err
	}
	settings.Enabled = make(map[string]types.EnabledMetric)
	for name, metric := range allSettings {
		if metric.Enabled {
			settings.Enabled[name] = metric
		}
	}
	return settings, nil
}

func persistTimeframeSelection(persistChan chan error, metric string, timeframe string) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		persistChan <- err
		return
	}
	defer conn.Close(ctx)
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
	mapMetrics := make(map[string]types.EnabledMetric)
	if err = json.Unmarshal([]byte(existing), &mapMetrics); err != nil {
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
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)
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
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
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
	mapMetrics := make(map[string]types.EnabledMetric)
	if err = json.Unmarshal([]byte(existing), &mapMetrics); err != nil {
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
