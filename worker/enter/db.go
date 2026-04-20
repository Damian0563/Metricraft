package enter

import (
	"context"
	"github.com/jackc/pgx/v4"
	"os"
	"time"
)

func InitDB(ctx context.Context, errChannel chan error) {
	time.Sleep(15 * time.Second)
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_LOGS"))
	if err != nil {
		errChannel <- err
		return
	}
	_, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS settings (id UUID PRIMARY KEY, realtime BOOL, enabled BOOL)")
	if err != nil {
		errChannel <- err
		return
	}
	conn.Exec(context.Background(), "INSERT INTO settings (realtime, enabled) VALUES (true, false)")
	_, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS logs (id UUID PRIMARY KEY,date TIMESTAMP,responseTime TEXT NOT NULL, url TEXT NOT NULL, \"user\" TEXT NOT NULL, payload TEXT, headers TEXT NOT NULL, method TEXT NOT NULL, status INTEGER NOT NULL)")
	if err != nil {
		errChannel <- err
		return
	}
	errChannel <- conn.Close(context.Background())
}

func Insert(payload Payload) error {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	_, err = conn.Exec(ctx, "INSERT INTO logs (date, responseTime, url, \"user\", payload, headers, method, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", payload.Metrics.Duration, payload.Metrics.Duration.String(), payload.Url, payload.Headers["User-Agent"], payload.Body, payload.Headers, payload.Method, payload.Metrics.StatusCode)
	return err
}

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

func GetRealtime() (bool, error) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return false, err
	}
	defer conn.Close(ctx)
	var realtime bool
	err = conn.QueryRow(ctx, "SELECT realtime FROM settings WHERE TRUE").Scan(&realtime)
	if err != nil {
		return false, err
	}
	return realtime, nil
}
