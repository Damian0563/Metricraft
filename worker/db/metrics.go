package db

import (
	"context"
	"fmt"
	pb "metricraft/proto/metricraft/proto"
	"strconv"
	"time"
)

func GetTrafficCongestion(ctx context.Context, rules []*pb.Rule, startDate time.Time, resolution int32, timezone string) (*pb.Congestion, error) {
	conn, err := getLogsPool()
	if err != nil {
		return nil, err
	}
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
	blacklisted := blacklistedRules(rules)
	grouping := groupingRules(rules)
	res, err := conn.Query(ctx, fmt.Sprintf(`
		SELECT bucket, grouped_url, COUNT(*)
		FROM (
			SELECT
				FLOOR(
					EXTRACT(EPOCH FROM
						date_trunc($6, (date AT TIME ZONE $4) AT TIME ZONE $5)
						- date_trunc($6, ($1 AT TIME ZONE $4) AT TIME ZONE $5)
					) / $3
				)::int AS bucket,
				%s AS grouped_url
			FROM logs
			WHERE date > $1 AND date <= $2
				AND %s
		) filtered
		GROUP BY bucket, grouped_url
		ORDER BY COUNT(*) DESC
	`, groupedUrlSQL("url", "$8"), blacklistFilterSQL("url", "$7")),
		alignedStart, endDate, increment.Seconds(), storageTimezone, tz, truncPeriod, blacklisted, grouping)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	for res.Next() {
		var url string
		var bucket int
		var count int
		if err := res.Scan(&bucket, &url, &count); err != nil {
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

func GetGeographicalTraffic(ctx context.Context, rules []*pb.Rule, startDate time.Time, timezone string) (*pb.Distribution, error) {
	conn, err := getLogsPool()
	if err != nil {
		return nil, err
	}
	loc := loadLocation(timezone)
	endDate := time.Now().In(loc).UTC()
	distribution := make(map[string]int32)
	blacklisted := blacklistedRules(rules)
	res, err := conn.Query(ctx, fmt.Sprintf(`
		SELECT country, COUNT(*)
		FROM logs
		WHERE date > $1 AND date <= $2 AND %s
		GROUP BY country
		ORDER BY COUNT(*) DESC`, blacklistFilterSQL("url", "$3")),
		startDate, endDate, blacklisted)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	for res.Next() {
		var country string
		var count int
		err = res.Scan(&country, &count)
		if err != nil {
			return nil, err
		}
		distribution[country] += int32(count)
	}
	return &pb.Distribution{Distribution: &pb.StringInt32Map{Values: distribution}}, nil
}

func GetP95Latency(ctx context.Context, rules []*pb.Rule, startDate time.Time, timezone string) (*pb.Distribution, error) {
	conn, err := getLogsPool()
	if err != nil {
		return nil, err
	}
	loc := loadLocation(timezone)
	endDate := time.Now().In(loc).UTC()
	distribution := make(map[string]int32)
	blacklisted := blacklistedRules(rules)
	grouping := groupingRules(rules)
	res, err := conn.Query(ctx, fmt.Sprintf(`
		SELECT grouped_url,
			percentile_cont(0.95) WITHIN GROUP (ORDER BY responsetime) AS percentile
		FROM (
			SELECT responsetime, %s AS grouped_url
			FROM logs
			WHERE date > $1 AND date <= $2 AND status BETWEEN 200 AND 299
				AND %s
		) filtered
		GROUP BY grouped_url
		ORDER BY percentile DESC`, groupedUrlSQL("url", "$4"), blacklistFilterSQL("url", "$3")),
		startDate, endDate, blacklisted, grouping)
	if err != nil {
		return nil, err
	}
	defer res.Close()
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

func GetUptimeScore(ctx context.Context, rules []*pb.Rule, startDate time.Time, timezone string) (*pb.FloatDistribution, error) {
	conn, err := getLogsPool()
	if err != nil {
		return nil, err
	}
	loc := loadLocation(timezone)
	endDate := time.Now().In(loc).UTC()
	distribution := make(map[string]float32)
	blacklisted := blacklistedRules(rules)
	grouping := groupingRules(rules)
	res, err := conn.Query(ctx, fmt.Sprintf(`
		SELECT grouped_url,
			100.0 * COUNT(*) FILTER (WHERE status BETWEEN 200 AND 299) / NULLIF(COUNT(*), 0) AS availability
		FROM (
			SELECT status, %s AS grouped_url
			FROM logs
			WHERE date > $1 AND date <= $2
				AND %s
		) filtered
		GROUP BY grouped_url
		ORDER BY availability DESC`, groupedUrlSQL("url", "$4"), blacklistFilterSQL("url", "$3")),
		startDate, endDate, blacklisted, grouping)
	if err != nil {
		return nil, err
	}
	defer res.Close()
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

func GetThroughput(ctx context.Context, rules []*pb.Rule, start time.Time, resolution int32, timezone string) (*pb.Throughput, error) {
	conn, err := getLogsPool()
	if err != nil {
		return nil, err
	}
	tz := validTimezone(timezone)
	loc := loadLocation(tz)
	alignedStart := alignStart(start, loc, resolution)
	var increment time.Duration
	var truncPeriod string
	if resolution == 0 {
		increment = time.Minute * 60
		truncPeriod = "hour"
	} else {
		increment = time.Hour * 24 * time.Duration(resolution)
		truncPeriod = "day"
	}
	endDate := rangeEnd(loc, increment, resolution)
	var throughput []*pb.ThroughputEntry
	for cursor := alignedStart; cursor.Before(endDate); cursor = cursor.Add(increment) {
		throughput = append(throughput, &pb.ThroughputEntry{
			Timerange: timerangeLabel(cursor.In(loc), increment, resolution, "throughput"),
			Value:     0,
		})
	}
	blacklisted := blacklistedRules(rules)
	res, err := conn.Query(ctx, fmt.Sprintf(`
		SELECT
			FLOOR(
				EXTRACT(EPOCH FROM
					date_trunc($6, (date AT TIME ZONE $4) AT TIME ZONE $5)
					- date_trunc($6, ($1 AT TIME ZONE $4) AT TIME ZONE $5)
				) / $3
			)::int AS bucket,
			COUNT(*)
		FROM logs
		WHERE date > $1 AND date <= $2 AND %s
		GROUP BY bucket
	`, blacklistFilterSQL("url", "$7")),
		alignedStart, endDate, increment.Seconds(), storageTimezone, tz, truncPeriod, blacklisted)
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
			throughput[bucket].Value += int32(count)
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
	err = conn.QueryRow(ctx, fmt.Sprintf(`
		SELECT COUNT(DISTINCT "user")
		FROM logs
		WHERE date > $1 AND date <= $2 AND %s`, blacklistFilterSQL("url", "$3")),
		alignedStart, endDate, blacklisted).Scan(&uniqUsers)
	if err != nil {
		return nil, err
	}
	return &pb.Throughput{Values: throughput, ComputedThroughput: computedThroughput, UniqUsers: uniqUsers}, nil
}

func GetGeographicalPerformance(ctx context.Context, rules []*pb.Rule, startDate time.Time, timezone string) (*pb.FloatDistribution, error) {
	conn, err := getLogsPool()
	if err != nil {
		return nil, err
	}
	loc := loadLocation(timezone)
	endDate := time.Now().In(loc).UTC()
	alignedStart := alignStart(startDate, loc, 0)
	blacklisted := blacklistedRules(rules)
	res, err := conn.Query(ctx, fmt.Sprintf(`
		SELECT country,
			percentile_cont(0.5) WITHIN GROUP (ORDER BY responsetime) AS median
		FROM (
			SELECT country, responsetime
			FROM logs
			WHERE date >= $1 AND date < $2
				AND %s
		) filtered
		GROUP BY country
		ORDER BY median DESC`, blacklistFilterSQL("url", "$3")),
		alignedStart, endDate, blacklisted)
	if err != nil {
		return nil, err
	}
	defer res.Close()
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

func GetStatusCodeDistribution(ctx context.Context, rules []*pb.Rule, startDate time.Time, resolution int32, timezone string) (*pb.Distribution, error) {
	conn, err := getLogsPool()
	if err != nil {
		return nil, err
	}
	end := time.Now().In(loadLocation(timezone)).UTC()
	start := alignStart(startDate, loadLocation(timezone), resolution)
	blacklisted := blacklistedRules(rules)
	res, err := conn.Query(ctx, fmt.Sprintf(`
		SELECT status, COUNT(*)
		FROM logs
		WHERE date >= $1 AND date < $2 AND %s
		GROUP BY status
		ORDER BY COUNT(*) DESC`, blacklistFilterSQL("url", "$3")),
		start, end, blacklisted)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	result := make(map[string]int32)
	for res.Next() {
		var status int
		var count int
		err = res.Scan(&status, &count)
		if err != nil {
			return nil, err
		}
		result[strconv.Itoa(status)] += int32(count)
	}
	return &pb.Distribution{Distribution: &pb.StringInt32Map{Values: result}}, nil
}

func GetRouteCongestion(ctx context.Context, rules []*pb.Rule, start time.Time, timezone string) (*pb.Distribution, error) {
	conn, err := getLogsPool()
	if err != nil {
		return nil, err
	}
	loc := loadLocation(timezone)
	alignedStart := alignStart(start, loc, 0)
	end := rangeEnd(loc, time.Hour, 0)
	blacklisted := blacklistedRules(rules)
	grouping := groupingRules(rules)
	res, err := conn.Query(ctx, fmt.Sprintf(`
		SELECT grouped_url, COUNT(*)
		FROM (
			SELECT %s AS grouped_url
			FROM logs
			WHERE date >= $1 AND date < $2 AND status BETWEEN 200 AND 299
				AND %s
		) filtered
		GROUP BY grouped_url
		ORDER BY COUNT(*) DESC`, groupedUrlSQL("url", "$4"), blacklistFilterSQL("url", "$3")),
		alignedStart, end, blacklisted, grouping)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	result := make(map[string]int32)
	for res.Next() {
		var url string
		var count int
		err = res.Scan(&url, &count)
		if err != nil {
			return nil, err
		}
		result[url] += int32(count)
	}
	return &pb.Distribution{Distribution: &pb.StringInt32Map{Values: result}}, nil
}

func GetHttpMethodDistribution(ctx context.Context, rules []*pb.Rule, startDate time.Time, resolution int32, timezone string) (*pb.Congestion, error) {
	conn, err := getLogsPool()
	if err != nil {
		return nil, err
	}
	loc := loadLocation(timezone)
	alignedStart := alignStart(startDate, loc, resolution)
	end := rangeEnd(loc, time.Hour, resolution)
	var increment time.Duration
	var truncPeriod string
	tz := validTimezone(timezone)
	if resolution == 0 {
		increment = time.Hour
		truncPeriod = "hour"
	} else {
		increment = time.Hour * 24 * time.Duration(resolution)
		truncPeriod = "day"
	}
	var congestion []*pb.CongestionEntry
	for cursor := alignedStart; cursor.Before(end); cursor = cursor.Add(increment) {
		congestion = append(congestion, &pb.CongestionEntry{
			Timerange: timerangeLabel(cursor.In(loc), increment, resolution, "http method"),
			Pairing:   &pb.StringInt32Map{Values: map[string]int32{}},
		})
	}
	blacklisted := blacklistedRules(rules)
	res, err := conn.Query(ctx, fmt.Sprintf(`
		SELECT bucket, method, COUNT(*)
		FROM (
			SELECT
				FLOOR(
					EXTRACT(EPOCH FROM
						date_trunc($6, (date AT TIME ZONE $4) AT TIME ZONE $5)
						- date_trunc($6, ($1 AT TIME ZONE $4) AT TIME ZONE $5)
					) / $3
				)::int AS bucket,
				method
			FROM logs
			WHERE date >= $1 AND date < $2 AND %s
		) filtered
		GROUP BY bucket, method`, blacklistFilterSQL("url", "$7")),
		alignedStart, end, increment.Seconds(), storageTimezone, tz, truncPeriod, blacklisted)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	for res.Next() {
		var bucket int
		var method string
		var count int
		if err = res.Scan(&bucket, &method, &count); err != nil {
			return nil, err
		}
		if bucket >= 0 && bucket < len(congestion) {
			congestion[bucket].Pairing.Values[method] += int32(count)
		}
	}
	return &pb.Congestion{Values: congestion}, nil
}

func GetUniqueVisitors(ctx context.Context, rules []*pb.Rule, start time.Time, resolution int32, timezone string) (*pb.SimpleRepeatedDistribution, error) {
	conn, err := getLogsPool()
	if err != nil {
		return nil, err
	}
	loc := loadLocation(timezone)
	alignedStart := alignStart(start, loc, resolution)
	end := rangeEnd(loc, time.Hour, resolution)
	var increment time.Duration
	var truncPeriod string
	tz := validTimezone(timezone)
	if resolution == 0 {
		increment = time.Hour
		truncPeriod = "hour"
	} else {
		increment = time.Hour * 24 * time.Duration(resolution)
		truncPeriod = "day"
	}
	var dist []*pb.StringInt32Map
	for cursor := alignedStart; cursor.Before(end); cursor = cursor.Add(increment) {
		timerange := timerangeLabel(cursor.In(loc), increment, resolution, "unique visitors")
		curr := map[string]int32{timerange: 0}
		dist = append(dist, &pb.StringInt32Map{Values: curr})
	}
	blacklisted := blacklistedRules(rules)
	res, err := conn.Query(ctx, fmt.Sprintf(`
		SELECT
			FLOOR(
				EXTRACT(EPOCH FROM
					date_trunc($6, (date AT TIME ZONE $4) AT TIME ZONE $5)
					- date_trunc($6, ($1 AT TIME ZONE $4) AT TIME ZONE $5)
				) / $3
			)::int AS bucket,
			COUNT(DISTINCT "user") AS user
		FROM logs
		WHERE date >= $1 AND date < $2 AND %s
		GROUP BY bucket
	`, blacklistFilterSQL("url", "$7")),
		alignedStart, end, increment.Seconds(), storageTimezone, tz, truncPeriod, blacklisted)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	for res.Next() {
		var bucket int
		var user int
		if err = res.Scan(&bucket, &user); err != nil {
			return nil, err
		}
		if bucket >= 0 && bucket < len(dist) {
			curr := dist[bucket].Values
			for k := range curr {
				curr[k] = int32(user)
			}
			dist[bucket].Values = curr
		}
	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	return &pb.SimpleRepeatedDistribution{Distribution: dist}, nil
}
