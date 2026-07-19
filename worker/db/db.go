package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "metricraft/proto/metricraft/proto"
	"os"
	"strconv"
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

func GetTrafficCongestion(ctx context.Context, startDate time.Time, resolution int32, timezone string) (*pb.Congestion, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)
	tz := validTimezone(timezone)
	loc := loadLocation(tz)
	alignedStart := alignStart(startDate, loc, resolution)
	var increment time.Duration
	if resolution == 0 {
		increment = time.Hour
	} else {
		increment = time.Hour * 24 * time.Duration(resolution)
	}
	endDate := rangeEnd(loc, increment, resolution)
	congestion := make([]*pb.CongestionEntry, 0)
	for cursor := alignedStart; cursor.Before(endDate); cursor = cursor.Add(increment) {
		congestion = append(congestion, &pb.CongestionEntry{
			Timerange: timerangeLabel(cursor.In(loc), increment, resolution, "congestion"),
			Pairing:   &pb.StringInt32Map{Values: map[string]int32{}},
		})
	}
	var truncPeriod string
	if resolution != 0 {
		truncPeriod = "day"
	} else {
		truncPeriod = "hour"
	}
	res, err := conn.Query(ctx, `
		SELECT url,
			FLOOR(
				EXTRACT(EPOCH FROM
					date_trunc($6, (date AT TIME ZONE $4) AT TIME ZONE $5)
					- date_trunc($6, ($1 AT TIME ZONE $4) AT TIME ZONE $5)
				) / $3
			)::int AS bucket,
			COUNT(*)
		FROM logs
		WHERE date > $1 AND date <= $2
		GROUP BY url, bucket
		ORDER BY COUNT(*) DESC
	`, alignedStart, endDate, increment.Seconds(), storageTimezone, tz, truncPeriod)
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

func InsertWorkerLog(ctx context.Context, url string, success bool) error {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	status := 0
	if success {
		status = 1
	}
	_, err = conn.Exec(ctx, "INSERT INTO worker_logs (date, url, up) VALUES ($1, $2, $3)", time.Now(), url, status)
	return err
}

func GetGeographicalTraffic(ctx context.Context, startDate time.Time, timezone string) (*pb.Distribution, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)
	loc := loadLocation(timezone)
	endDate := time.Now().In(loc).UTC()
	distribution := make(map[string]int32)
	res, err := conn.Query(ctx, "SELECT country,COUNT(*) FROM logs WHERE date > $1 AND date <= $2 GROUP BY country ORDER BY COUNT(*) DESC", startDate, endDate)
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

func GetP95Latency(ctx context.Context, startDate time.Time, timezone string) (*pb.Distribution, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)
	loc := loadLocation(timezone)
	endDate := time.Now().In(loc).UTC()
	distribution := make(map[string]int32)
	res, err := conn.Query(ctx, "SELECT url,percentile_cont(0.95) WITHIN GROUP (ORDER BY responsetime) AS percentile FROM logs WHERE date > $1 AND date <= $2 AND status BETWEEN 200 AND 299 GROUP BY url ORDER BY percentile DESC", startDate, endDate)
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

func GetUptimeScore(ctx context.Context, startDate time.Time, timezone string) (*pb.FloatDistribution, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)
	loc := loadLocation(timezone)
	endDate := time.Now().In(loc).UTC()
	distribution := make(map[string]float32)
	res, err := conn.Query(ctx, "SELECT url, 100.0 * COUNT(*) FILTER (WHERE status BETWEEN 200 AND 299) / NULLIF(COUNT(*), 0) AS availability FROM logs WHERE date > $1 AND date <= $2 GROUP BY url ORDER BY availability DESC", startDate, endDate)
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

func GetThroughput(ctx context.Context, start time.Time, resolution int32, timezone string) (*pb.Throughput, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)
	tz := validTimezone(timezone)
	loc := loadLocation(tz)
	alignedStart := alignStart(start, loc, resolution)
	var increment time.Duration
	if resolution == 0 {
		increment = time.Minute * 60
	} else {
		increment = time.Hour * 24 * time.Duration(resolution)
	}
	endDate := rangeEnd(loc, increment, resolution)
	var throughput []*pb.ThroughputEntry
	for cursor := alignedStart; cursor.Before(endDate); cursor = cursor.Add(increment) {
		throughput = append(throughput, &pb.ThroughputEntry{
			Timerange: timerangeLabel(cursor.In(loc), increment, resolution, "throughput"),
			Value:     0,
		})
	}
	var truncPeriod string
	if resolution != 0 {
		truncPeriod = "day"
	} else {
		truncPeriod = "hour"
	}
	res, err := conn.Query(ctx, `
		SELECT
			FLOOR(
				EXTRACT(EPOCH FROM
					date_trunc($6, (date AT TIME ZONE $4) AT TIME ZONE $5)
					- date_trunc($6, ($1 AT TIME ZONE $4) AT TIME ZONE $5)
				) / $3
			)::int AS bucket,
			COUNT(*)
		FROM logs
		WHERE date > $1 AND date <= $2
		GROUP BY bucket
	`, alignedStart, endDate, increment.Seconds(), storageTimezone, tz, truncPeriod)
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
	if seconds := endDate.Sub(alignedStart).Seconds(); seconds > 0 {
		computedThroughput = float32(total) / float32(seconds)
	}
	var uniqUsers int32
	if err = conn.QueryRow(ctx, "SELECT COUNT(DISTINCT(\"user\")) FROM logs WHERE date > $1 AND date <= $2", alignedStart, endDate).Scan(&uniqUsers); err != nil {
		return nil, err
	}
	return &pb.Throughput{Values: throughput, ComputedThroughput: computedThroughput, UniqUsers: uniqUsers}, nil
}

func GetWorkerUptime(ctx context.Context, url string, timezone string, pollIntervalSetting int32) (*pb.WorkerUptime, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)
	loc := loadLocation(timezone)
	uptime := make([]*pb.WorkerUptimeEntry, 0)
	endDate := time.Now().In(loc).UTC().Add(-time.Hour * 24 * 31)
	res, err := conn.Query(ctx, "SELECT date, up FROM worker_logs WHERE url = $1 AND date > $2 ORDER BY date ASC", url, endDate)
	if err != nil {
		return nil, err
	}
	pollInterval := time.Minute * time.Duration(pollIntervalSetting)
	var lastPoll time.Time
	for res.Next() {
		var date time.Time
		var status int
		if err := res.Scan(&date, &status); err != nil {
			return nil, err
		}
		if !lastPoll.IsZero() {
			for date.After(lastPoll.Add(pollInterval)) {
				lastPoll = lastPoll.Add(pollInterval)
				uptime = append(uptime, &pb.WorkerUptimeEntry{
					Status: -1,
					Stamp:  timestamppb.New(lastPoll),
				})
			}
		}
		uptime = append(uptime, &pb.WorkerUptimeEntry{
			Status: int32(status),
			Stamp:  timestamppb.New(date),
		})
		lastPoll = date
	}
	now := time.Now().In(loc).UTC()
	for !lastPoll.IsZero() && now.After(lastPoll.Add(pollInterval)) {
		lastPoll = lastPoll.Add(pollInterval)
		uptime = append(uptime, &pb.WorkerUptimeEntry{
			Status: -1,
			Stamp:  timestamppb.New(lastPoll),
		})
	}
	//delete old logs, do not throw error if occured
	_, _ = conn.Exec(ctx, "DELETE FROM worker_logs WHERE url = $1 AND date < $2", url, endDate)
	return &pb.WorkerUptime{Entries: uptime}, nil
}

func GetGeographicalPerformance(ctx context.Context, startDate time.Time, timezone string) (*pb.FloatDistribution, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)
	loc := loadLocation(timezone)
	endDate := time.Now().In(loc).UTC()
	alignedStart := alignStart(startDate, loc, 0)
	res, err := conn.Query(ctx, "SELECT country, percentile_cont(0.5) WITHIN GROUP (ORDER BY responsetime) AS median FROM logs WHERE date >= $1 AND date < $2 AND status BETWEEN 200 AND 299 GROUP BY country ORDER BY median DESC", alignedStart, endDate)
	if err != nil {
		return nil, err
	}
	distribution := make(map[string]float32)
	for res.Next() {
		var country string
		var percentile float64
		err = res.Scan(&country, &percentile)
		if err != nil {
			return nil, err
		}
		fortmatted := fmt.Sprintf("%.2f", percentile)
		finalPercentile, err := strconv.ParseFloat(fortmatted, 32)
		if err != nil {
			return nil, err
		}
		distribution[country] = float32(finalPercentile)
	}
	return &pb.FloatDistribution{Distribution: &pb.StringFloat32Map{Values: distribution}}, nil
}

func GetStatusCodeDistribution(ctx context.Context, startDate time.Time, resolution int32, timezone string) (*pb.Distribution, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)
	end := time.Now().In(loadLocation(timezone)).UTC()
	start := alignStart(startDate, loadLocation(timezone), resolution)
	res, err := conn.Query(ctx, "SELECT status,COUNT(*) FROM logs WHERE date >= $1 AND date < $2 GROUP BY status ORDER BY COUNT(*) DESC", start, end)
	if err != nil {
		return nil, err
	}
	result := make(map[string]int32)
	for res.Next() {
		var status int
		var count int
		err = res.Scan(&status, &count)
		if err != nil {
			return nil, err
		}
		result[strconv.Itoa(status)] = int32(count)
	}
	return &pb.Distribution{Distribution: &pb.StringInt32Map{Values: result}}, nil
}

func GetRouteCongestion(ctx context.Context, start time.Time, timezone string) (*pb.Distribution, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)
	loc := loadLocation(timezone)
	alignedStart := alignStart(start, loc, 0)
	end := rangeEnd(loc, time.Hour, 0)
	res, err := conn.Query(ctx, "SELECT url,COUNT(*) FROM logs WHERE date >= $1 AND date < $2 AND status BETWEEN 200 AND 299 GROUP BY url ORDER BY COUNT(*) DESC", alignedStart, end)
	result := make(map[string]int32)
	for res.Next() {
		var url string
		var count int
		err = res.Scan(&url, &count)
		if err != nil {
			return nil, err
		}
		result[url] = int32(count)
	}
	return &pb.Distribution{Distribution: &pb.StringInt32Map{Values: result}}, nil
}

func DeleteWorkerlogs(ctx context.Context, url string) error {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "DELETE FROM worker_logs WHERE url=$1", url)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	err = tx.Commit(ctx)
	return err
}
