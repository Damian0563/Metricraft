package main

import (
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
	_, err = conn.Exec(ctx, "UPDATE settings SET realtime = $1 WHERE TRUE", enabled)
	return err
}

func ChangeLogsRetention(retention int) error {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	_, err = conn.Exec(ctx, "UPDATE settings SET retention = $1 WHERE TRUE", retention)
	return err
}

func GetSettings() (Settings, error) {
	var settings Settings
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return Settings{}, err
	}
	defer conn.Close(ctx)
	var enabled string
	err = conn.QueryRow(ctx, "SELECT realtime,enabled,retention FROM settings WHERE TRUE").Scan(&settings.Realtime, &enabled, &settings.Retention)
	if err != nil {
		return Settings{}, err
	}
	settings.Enabled = make(map[string]bool)
	err = json.Unmarshal([]byte(enabled), &settings.Enabled)
	if err != nil {
		return Settings{}, err
	}
	return settings, nil
}

func ChangeMetrics(metrics []Metric) error {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	mapMetrics := make(map[string]bool)
	for _, metric := range metrics {
		mapMetrics[metric.Name] = metric.Enabled
	}
	var stringifiedMetrics []byte
	stringifiedMetrics, err = json.Marshal(mapMetrics)
	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx, "UPDATE settings SET enabled = $1 WHERE TRUE", string(stringifiedMetrics))
	return err
}
