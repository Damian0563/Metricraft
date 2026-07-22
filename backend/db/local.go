package db

import (
	"context"
)

func VerifyAppName(appName string) bool {
	ctx := context.Background()
	conn, err := GetLogsPool()
	if err != nil {
		return false
	}
	var exists bool
	err = conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM settings WHERE appname = $1)", appName).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}
