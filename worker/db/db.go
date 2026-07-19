package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
	"time"
	"worker/external"
	"worker/types"
)

func timerangeLabel(rangeStart time.Time, increment time.Duration, resolution int32, mode string) string {
	if resolution == 0 {
		if mode == "congestion" {
			return fmt.Sprintf("%02d:00", rangeStart.Hour())
		} else {
			return fmt.Sprintf("%02d:%02d", rangeStart.Hour(), rangeStart.Minute())
		}
	} else if resolution != 1 {
		rangeEnd := rangeStart.Add(increment)
		return fmt.Sprintf("%v-%v", rangeStart.Format("02.01"), rangeEnd.Add(-time.Hour*24).Format("02.01"))
	}
	return rangeStart.Format("02.01")
}

func InitDB(ctx context.Context, errChannel chan error) {
	time.Sleep(15 * time.Second)
	appName := os.Getenv("APPNAME")
	if appName == "" {
		errChannel <- fmt.Errorf("APPNAME must be set")
		return
	}
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_LOGS"))
	if err != nil {
		errChannel <- err
		return
	}
	_, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS settings (realtime BOOL, enabled TEXT, retention INTEGER, appName TEXT)")
	if err != nil {
		errChannel <- err
		return
	}
	var count int
	err = conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM settings").Scan(&count)
	if err != nil {
		errChannel <- err
		return
	}
	if count == 0 {
		_, err = conn.Exec(context.Background(), "INSERT INTO settings (realtime, enabled, retention, appname) VALUES (true, '{\"Geographical traffic\":{\"enabled\":true,\"timeframe\":\"7d\"},\"P95 Latency\":{\"enabled\":true,\"timeframe\":\"7d\"},\"Traffic congestion trends\":{\"enabled\":false,\"timeframe\":\"7d\"},\"Uptime Score\":{\"enabled\":true,\"timeframe\":\"7d\"},\"Geographic performance\":{\"enabled\":false,\"timeframe\":\"7d\"},\"Status code distribution\":{\"enabled\":false,\"timeframe\":\"7d\"},\"Median response time\":{\"enabled\":true,\"timeframe\":\"7d\"},\"Throughput\":{\"enabled\":true,\"timeframe\":\"7d\"},\"Route congestion\":{\"enabled\":false, \"timeframe\":\"7d\"},\"Unique visitors\":{\"enabled\":false, \"timeframe\":\"7d\"},\"HTTP method distribution\":{\"enabled\":false, \"timeframe\":\"7d\"}}', 30, $1)", appName)
		if err != nil {
			errChannel <- err
			return
		}
	}
	if _, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS logs (date TIMESTAMP,responseTime INTEGER, url TEXT NOT NULL, \"user\" TEXT NOT NULL,country TEXT NOT NULL ,payload TEXT, headers TEXT NOT NULL, method TEXT NOT NULL, status INTEGER NOT NULL)"); err != nil {
		errChannel <- err
		return
	}
	if _, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS worker_logs (date TIMESTAMP, url TEXT, up INTEGER)"); err != nil {
		errChannel <- err
		return
	}
	if _, err = conn.Exec(context.Background(), "CREATE INDEX IF NOT EXISTS idx_date ON logs (date)"); err != nil {
		errChannel <- err
		return
	}
	if _, err = conn.Exec(context.Background(), "CREATE INDEX IF NOT EXISTS idx_url ON logs (url)"); err != nil {
		errChannel <- err
		return
	}
	errChannel <- conn.Close(context.Background())
}

func Insert(payload types.Payload) error {
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
	_, err = conn.Exec(ctx, "INSERT INTO logs (date, responseTime, url, \"user\",country, payload, headers, method, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", payload.Time.UTC(), payload.Metrics.Duration, payload.Url, realip, country, string(body), string(headers), payload.Method, payload.Metrics.StatusCode)
	return err
}

func GetAppname(ctx context.Context) (string, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return "", err
	}
	defer conn.Close(ctx)
	var appName string
	err = conn.QueryRow(ctx, "SELECT appname FROM settings").Scan(&appName)
	if err != nil {
		return "", err
	}
	return appName, nil
}
