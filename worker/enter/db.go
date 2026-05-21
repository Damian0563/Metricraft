package enter

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	pb "metricraft/proto/metricraft/proto"
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
	_, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS settings (realtime BOOL, enabled TEXT, retention INTEGER)")
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
		conn.Exec(context.Background(), "INSERT INTO settings (realtime, enabled, retention) VALUES (true, '{\"Geographical traffic\":true,\"P95 Latency\":true,\"Traffic congestion trends\":false,\"Uptime Score\":true,\"Geographic performance\":false,\"Status code distribution\":false,\"Median response time\":true,\"Throughput\":true}', 30)")
	}
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

func GetTrafficCongestion(ctx context.Context, startDate time.Time, resolution int32) (*pb.Congestion, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)
	endDate := time.Now()
	originalStartDate := startDate
	increment := time.Hour * 24 * time.Duration(resolution)
	congestion := make([]*pb.CongestionEntry, 0)
	index := make(map[string]int)
	for startDate.Before(endDate) {
		rangeStart := startDate
		rangeEnd := startDate.Add(increment)
		var timerange string
		if resolution != 1 {
			timerange = fmt.Sprintf("%v-%v", rangeStart.Format("02/01"), rangeEnd.Add(-time.Hour*24).Format("02/01"))
		} else {
			timerange = rangeStart.Format("02/01")
		}
		index[timerange] = len(congestion)
		congestion = append(congestion, &pb.CongestionEntry{
			Timerange: timerange,
			Pairing:   &pb.StringInt32Map{Values: map[string]int32{}},
		})
		startDate = rangeEnd
	}
	res, err := conn.Query(ctx, "SELECT url,COUNT(*),date FROM logs WHERE date BETWEEN $1 AND $2 GROUP BY url,date ORDER BY COUNT(*) DESC", originalStartDate, endDate)
	if err != nil {
		return nil, err
	}
	for res.Next() {
		var url string
		var count int
		var date time.Time
		err = res.Scan(&url, &count, &date)
		if err != nil {
			return nil, err
		}
		var key string
		if resolution != 1 {
			elapsed := date.Sub(originalStartDate)
			intervals := int(elapsed / increment)
			rangeStart := originalStartDate.Add(time.Duration(intervals) * increment)
			rangeEnd := rangeStart.Add(increment)
			key = fmt.Sprintf("%v-%v", rangeStart.Format("02/01"), rangeEnd.Add(-time.Hour*24).Format("02/01"))
		} else {
			key = date.Format("02/01")
		}
		if i, ok := index[key]; ok {
			congestion[i].Pairing.Values[url] = int32(count)
		}
	}
	fmt.Println(congestion)
	return &pb.Congestion{Values: congestion}, nil
}
