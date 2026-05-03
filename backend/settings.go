package main

import (
	"context"
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
	err = conn.QueryRow(ctx, "SELECT realtime,enabled,retention FROM settings WHERE TRUE").Scan(&settings.Realtime, &settings.Enabled, &settings.Retention)
	if err != nil {
		return Settings{}, err
	}
	return settings, nil
}
