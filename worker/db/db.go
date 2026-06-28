package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	pb "metricraft/proto/metricraft/proto"
	"os"
	"time"
	"worker/external"
	"worker/types"
)

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
		_, err = conn.Exec(context.Background(), "INSERT INTO settings (realtime, enabled, retention, appname) VALUES (true, '{\"Geographical traffic\":{\"enabled\":true,\"timeframe\":\"7d\"},\"P95 Latency\":{\"enabled\":true,\"timeframe\":\"7d\"},\"Traffic congestion trends\":{\"enabled\":false,\"timeframe\":\"7d\"},\"Uptime Score\":{\"enabled\":true,\"timeframe\":\"7d\"},\"Geographic performance\":{\"enabled\":false,\"timeframe\":\"7d\"},\"Status code distribution\":{\"enabled\":false,\"timeframe\":\"7d\"},\"Median response time\":{\"enabled\":true,\"timeframe\":\"7d\"},\"Throughput\":{\"enabled\":true,\"timeframe\":\"7d\"}}', 30, $1)", appName)
		if err != nil {
			errChannel <- err
			return
		}
	}
	if _, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS logs (date TIMESTAMP,responseTime INTEGER, url TEXT NOT NULL, \"user\" TEXT NOT NULL,country TEXT NOT NULL ,payload TEXT, headers TEXT NOT NULL, method TEXT NOT NULL, status INTEGER NOT NULL)"); err != nil {
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
	_, err = conn.Exec(ctx, "INSERT INTO logs (date, responseTime, url, \"user\",country, payload, headers, method, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", payload.Time, payload.Metrics.Duration, payload.Url, realip, country, string(body), string(headers), payload.Method, payload.Metrics.StatusCode)
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
	return &pb.Congestion{Values: congestion}, nil
}

func GetGeographicalTraffic(ctx context.Context, startDate time.Time) (*pb.Distribution, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)
	endDate := time.Now()
	distribution := make(map[string]int32)
	res, err := conn.Query(ctx, "SELECT country,COUNT(*) FROM logs WHERE date BETWEEN $1 AND $2 GROUP BY country ORDER BY COUNT(*) DESC", startDate, endDate)
	if err != nil {
		return nil, err
	}
	for res.Next() {
		var country string
		var count int
		err = res.Scan(&country, &count)
		if err != nil {
			return nil, err
		}
		distribution[country] = int32(count)
	}
	return &pb.Distribution{Distribution: &pb.StringInt32Map{Values: distribution}}, nil
}

func GetP95Latency(ctx context.Context, startDate time.Time) (*pb.Distribution, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)
	endDate := time.Now()
	distribution := make(map[string]int32)
	res, err := conn.Query(ctx, "SELECT url,percentile_cont(0.95) WITHIN GROUP (ORDER BY responsetime) AS percentile FROM logs WHERE date BETWEEN $1 AND $2 AND status BETWEEN 200 AND 299 GROUP BY url ORDER BY percentile DESC", startDate, endDate)
	if err != nil {
		return nil, err
	}
	for res.Next() {
		var url string
		var percentile float64
		err = res.Scan(&url, &percentile)
		if err != nil {
			return nil, err
		}
		distribution[url] = int32(percentile)
	}
	return &pb.Distribution{Distribution: &pb.StringInt32Map{Values: distribution}}, nil
}

func GetUptimeScore(ctx context.Context, startDate time.Time) (*pb.FloatDistribution, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)
	endDate := time.Now()
	distribution := make(map[string]float32)
	res, err := conn.Query(ctx, "SELECT url, 100.0 * COUNT(*) FILTER (WHERE status BETWEEN 200 AND 299) / NULLIF(COUNT(*), 0) AS availability FROM logs WHERE date BETWEEN $1 AND $2 GROUP BY url ORDER BY availability DESC", startDate, endDate)
	if err != nil {
		return nil, err
	}
	for res.Next() {
		var url string
		var availability float64
		err = res.Scan(&url, &availability)
		if err != nil {
			return nil, err
		}
		distribution[url] = float32(availability)
	}
	return &pb.FloatDistribution{Distribution: &pb.StringFloat32Map{Values: distribution}}, nil
}
