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

func timerangeLabel(rangeStart time.Time, increment time.Duration, resolution int32) string {
	if resolution == 0 {
		return fmt.Sprintf("%s %d:00", rangeStart.Format("02.01"), rangeStart.Hour())
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
	now := time.Now()
	endDate := now
	if resolution == 0 {
		endDate = now.Truncate(time.Hour).Add(2 * time.Hour)
	}
	var increment time.Duration
	if resolution == 0 {
		increment = time.Hour
	} else {
		increment = time.Hour * 24 * time.Duration(resolution)
	}
	congestion := make([]*pb.CongestionEntry, 0)
	for cursor := startDate; cursor.Before(endDate); cursor = cursor.Add(increment) {
		congestion = append(congestion, &pb.CongestionEntry{
			Timerange: timerangeLabel(cursor, increment, resolution),
			Pairing:   &pb.StringInt32Map{Values: map[string]int32{}},
		})
	}
	res, err := conn.Query(ctx, `
		SELECT url, FLOOR(EXTRACT(EPOCH FROM (date - $1)) / $3)::int, COUNT(*)
		FROM logs
		WHERE date >= $1 AND date < $2
		GROUP BY url, 2
		ORDER BY COUNT(*) DESC
	`, startDate, endDate, increment.Seconds())
	if err != nil {
		return nil, err
	}
	defer res.Close()
	for res.Next() {
		var url string
		var bucket int
		var count int
		if err := res.Scan(&url, &bucket, &count); err != nil {
			return nil, err
		}
		if bucket >= 0 && bucket < len(congestion) {
			congestion[bucket].Pairing.Values[url] += int32(count)
		}
	}
	if err := res.Err(); err != nil {
		return nil, err
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

func GetThroughput(ctx context.Context, start time.Time, resolution int32) (*pb.Throughput, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)
	endDate := time.Now()
	if resolution == 0 {
		endDate = endDate.Truncate(time.Hour).Add(time.Hour)
	}
	var increment time.Duration
	if resolution == 0 {
		increment = time.Hour
	} else {
		increment = time.Hour * 24 * time.Duration(resolution)
	}
	throughput := make([]*pb.ThroughputEntry, 0)
	for cursor := start; cursor.Before(endDate); cursor = cursor.Add(increment) {
		throughput = append(throughput, &pb.ThroughputEntry{
			Timerange: timerangeLabel(cursor, increment, resolution),
			Value:     0,
		})
	}
	res, err := conn.Query(ctx, `
		SELECT FLOOR(EXTRACT(EPOCH FROM (date - $1)) / $3)::int, COUNT(*)
		FROM logs
		WHERE date >= $1 AND date < $2
		GROUP BY 1
	`, start, endDate, increment.Seconds())
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var total int32
	for res.Next() {
		var bucket int
		var count int
		if err := res.Scan(&bucket, &count); err != nil {
			return nil, err
		}
		if bucket >= 0 && bucket < len(throughput) {
			throughput[bucket].Value = int32(count)
		}
		total += int32(count)
	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	var computedThroughput float32
	if seconds := endDate.Sub(start).Seconds(); seconds > 0 {
		computedThroughput = float32(total) / float32(seconds)
	}
	return &pb.Throughput{Values: throughput, ComputedThroughput: computedThroughput}, nil
}
