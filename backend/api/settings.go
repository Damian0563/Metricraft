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
	settings.Enabled = make(map[string]bool)
	err = json.Unmarshal([]byte(enabled), &settings.Enabled)
	if err != nil {
		return types.Settings{}, err
	}
	return settings, nil
}

func verifyAppName(appName string) bool {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return false
	}
	defer conn.Close(ctx)
	var exists bool
	err = conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM settings WHERE appname = $1)", appName).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func GetUrls() ([]string, error) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)
	res, err := conn.Query(ctx, "SELECT DISTINCT url FROM logs WHERE TRUE ORDER BY url")
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
