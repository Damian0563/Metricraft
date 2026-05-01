package enter

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v4"
	"metricraft/worker/external"
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
	_, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS settings (realtime BOOL, enabled TEXT)")
	if err != nil {
		errChannel <- err
		return
	}
	conn.Exec(context.Background(), "INSERT INTO settings (realtime, enabled) VALUES (true, '{}')")
	if _, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS logs (date TIMESTAMP,responseTime INTEGER, url TEXT NOT NULL, \"user\" TEXT NOT NULL,country TEXT NOT NULL ,payload TEXT, headers TEXT NOT NULL, method TEXT NOT NULL, status INTEGER NOT NULL)"); err != nil {
		errChannel <- err
		return
	}
	if _, err = conn.Exec(context.Background(), "CREATE INDEX IF NOT EXISTS idx_date ON logs (date)"); err != nil {
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
	headers, _ := json.Marshal(payload.Headers)
	body, _ := json.Marshal(payload.Body)
	var realip string = payload.Headers["X-Real-IP"]
	var country string = "Unknown"
	if realip == "" {
		realip = "Unknown"
	} else {
		country, err = external.GetCountryOrigin(realip)
		if err != nil {
			country = "Unknown"
		}
	}
	_, err = conn.Exec(ctx, "INSERT INTO logs (date, responseTime, url, \"user\",country, payload, headers, method, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", payload.Time, payload.Metrics.Duration, payload.Url, realip, country, string(body), string(headers), payload.Method, payload.Metrics.StatusCode)
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
