package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"os"
)

func VerifyAppName(appName string) bool {
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
