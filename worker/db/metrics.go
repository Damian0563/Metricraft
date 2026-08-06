package db

import (
	"context"
	"database/sql"
	"fmt"
	pb "metricraft/proto/metricraft/proto"
	"strconv"
	"strings"
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
	increment, truncPeriod := getIncrementTruncPeriod(resolution)
	endDate := rangeEnd(loc, increment, resolution)
	congestion := make([]*pb.CongestionEntry, 0)
	for cursor := alignedStart; cursor.Before(endDate); cursor = cursor.Add(increment) {
		congestion = append(congestion, &pb.CongestionEntry{
			Timerange: timerangeLabel(cursor.In(loc), increment, resolution, "congestion"),
			Pairing:   &pb.StringInt32Map{Values: map[string]int32{}},
		})
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
	`, groupedUrlSQL("$8"), blacklistFilterSQL("$7")),
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
		ORDER BY COUNT(*) DESC`, blacklistFilterSQL("$3")),
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
		ORDER BY percentile DESC`, groupedUrlSQL("$4"), blacklistFilterSQL("$3")),
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
		ORDER BY availability DESC`, groupedUrlSQL("$4"), blacklistFilterSQL("$3")),
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
	increment, truncPeriod := getIncrementTruncPeriod(resolution)
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
	`, blacklistFilterSQL("$7")),
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
		WHERE date > $1 AND date <= $2 AND %s`, blacklistFilterSQL("$3")),
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
		ORDER BY median DESC`, blacklistFilterSQL("$3")),
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
		ORDER BY COUNT(*) DESC`, blacklistFilterSQL("$3")),
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
		ORDER BY COUNT(*) DESC`, groupedUrlSQL("$4"), blacklistFilterSQL("$3")),
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
	tz := validTimezone(timezone)
	increment, truncPeriod := getIncrementTruncPeriod(resolution)
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
		GROUP BY bucket, method`, blacklistFilterSQL("$7")),
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
	tz := validTimezone(timezone)
	increment, truncPeriod := getIncrementTruncPeriod(resolution)
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
	`, blacklistFilterSQL("$7")),
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

func GetHotHours(ctx context.Context, rules []*pb.Rule, start time.Time, timezone string) (*pb.SimpleRepeatedDistribution, error) {
	conn, err := getLogsPool()
	if err != nil {
		return nil, err
	}
	tz := validTimezone(timezone)
	loc := loadLocation(timezone)
	end := time.Now().Add(time.Hour).Truncate(time.Hour).In(loc)
	startAdjusted := start.In(loc).Truncate(time.Hour)
	blacklisted := blacklistedRules(rules)
	dist := make([]*pb.StringInt32Map, 24)
	for hour := 0; hour < 24; hour++ {
		label := fmt.Sprintf("%02d:00", hour)
		dist[hour] = &pb.StringInt32Map{Values: map[string]int32{label: 0}}
	}
	res, err := conn.Query(ctx, fmt.Sprintf(`
		SELECT
			EXTRACT(HOUR FROM (date AT TIME ZONE $4))::int AS bucket,
			COUNT(*) AS count
		FROM logs
		WHERE date >= $1 AND date < $2 AND %s
		GROUP BY bucket
	`, blacklistFilterSQL("$3")), startAdjusted, end, blacklisted, tz)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	for res.Next() {
		var bucket int
		var count int
		if err = res.Scan(&bucket, &count); err != nil {
			return nil, err
		}
		if bucket >= 0 && bucket < len(dist) {
			label := fmt.Sprintf("%02d:00", bucket)
			dist[bucket].Values[label] = int32(count)
		}
	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	return &pb.SimpleRepeatedDistribution{Distribution: dist}, nil
}

func GetCustomMetricDataBuckets(ctx context.Context, metric *pb.CustomMetric, rules []*pb.Rule, start time.Time, resolution int32, timezone string) (*pb.CustomMetricData, error) {
	conn, err := getLogsPool()
	if err != nil {
		return nil, fmt.Errorf("failed to get established connection to database while fetching %s", metric.Name)
	}
	tz := validTimezone(timezone)
	loc := loadLocation(tz)
	alignedStart := alignStart(start, loc, resolution)
	increment, truncPeriod := getIncrementTruncPeriod(resolution)
	endDate := rangeEnd(loc, increment, resolution)
	buckets := make([]*pb.CustomMetricDataPoint, 0)
	for cursor := alignedStart; cursor.Before(endDate); cursor = cursor.Add(increment) {
		buckets = append(buckets, &pb.CustomMetricDataPoint{
			Grouping: timerangeLabel(cursor.In(loc), increment, resolution, "custom"),
			Value:    0,
		})
	}
	var blacklisted []string
	grouping := groupingRules(rules)
	applyRules := metric.ApplyRules
	if applyRules {
		blacklisted = blacklistedRules(rules)
	}
	args := []any{alignedStart, endDate, increment.Seconds(), storageTimezone, tz, truncPeriod, blacklisted}
	var pathParam, methodParam, groupingParam, selectorParam string
	args, pathParam = appendQueryParam(args, metric.Path)
	args, methodParam = appendQueryParam(args, metric.Method)
	if applyRules && len(grouping) > 0 {
		args, groupingParam = appendQueryParam(args, grouping)
	}
	inspectedField := resolveFieldName(metric.Source)
	if selector := strings.TrimSpace(metric.Selector); selector != "" {
		if inspectedField == "payload" || inspectedField == "headers" {
			selector = jsonPathSelector(selector)
		}
		args, selectorParam = appendQueryParam(args, selector)
	}
	valueExpr := customMetricExpr(inspectedField, selectorParam)
	aggregationExpr := resolveAggregationType(metric.Aggregation, metric.ValueType)
	query := fmt.Sprintf(`
		SELECT
			FLOOR(
				EXTRACT(EPOCH FROM
					date_trunc($6, (date AT TIME ZONE $4) AT TIME ZONE $5)
					- date_trunc($6, ($1 AT TIME ZONE $4) AT TIME ZONE $5)
				) / $3
			)::int AS bucket,
			%s AS count
		FROM (
			SELECT
				date,
				%s AS value
			FROM logs
			WHERE %s AND %s AND %s
		) filtered
		WHERE date >= $1 AND date < $2
		GROUP BY bucket
	`, aggregationExpr, valueExpr, customMetricInnerWhere(pathParam, methodParam, groupingParam, applyRules), customMetricLogicMatch(inspectedField, selectorParam), blacklistFilterSQL("$7"))
	res, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query custom metric data: %s", metric.Name)
	}
	defer res.Close()
	for res.Next() {
		var bucket int
		var count sql.NullFloat64
		if err = res.Scan(&bucket, &count); err != nil {
			return nil, fmt.Errorf("failed to scan custom metric data: %s", metric.Name)
		}
		if bucket >= 0 && bucket < len(buckets) && count.Valid {
			buckets[bucket].Value = float32(count.Float64)
		}
	}
	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan custom metric data: %s", metric.Name)
	}
	return &pb.CustomMetricData{Metrics: buckets}, nil
}

func GetCustomMetricDataCummulative(ctx context.Context, metric *pb.CustomMetric, rules []*pb.Rule, start time.Time, timezone string) (*pb.CustomMetricData, error) {
	conn, err := getLogsPool()
	if err != nil {
		return nil, fmt.Errorf("failed to get established connection to database while fetching %s", metric.Name)
	}
	tz := validTimezone(timezone)
	loc := loadLocation(tz)
	alignedStart := start.In(loc)
	alignedEnd := time.Now().In(loc)
	var blacklisted []string
	grouping := groupingRules(rules)
	applyRules := metric.ApplyRules
	if applyRules {
		blacklisted = blacklistedRules(rules)
	}
	args := []any{alignedStart, alignedEnd}
	var pathParam, methodParam, groupingParam, selectorParam string
	args, pathParam = appendQueryParam(args, metric.Path)
	args, methodParam = appendQueryParam(args, metric.Method)
	if applyRules && len(grouping) > 0 {
		args, groupingParam = appendQueryParam(args, grouping)
	}
	inspectedField := resolveFieldName(metric.Source)
	if selector := strings.TrimSpace(metric.Selector); selector != "" {
		if inspectedField == "payload" || inspectedField == "headers" {
			selector = jsonPathSelector(selector)
		}
		args, selectorParam = appendQueryParam(args, selector)
	}
	args = append(args, blacklisted)
	blacklistParam := fmt.Sprintf("$%d", len(args))
	valueExpr := customMetricExpr(inspectedField, selectorParam)
	aggregationExpr := resolveAggregationType(metric.Aggregation, metric.ValueType)
	query := fmt.Sprintf(`
		SELECT
			predicate,
			%s AS value
		FROM (
			SELECT
				%s AS value,
			FROM logs
			WHERE %s AND %s AND %s
		) filtered
		WHERE date >= $1 AND date < $2
		GROUP BY predicate
		ORDER BY value DESC
	`, aggregationExpr, valueExpr, customMetricInnerWhere(pathParam, methodParam, groupingParam, applyRules), customMetricLogicMatch(inspectedField, selectorParam), blacklistFilterSQL(blacklistParam))
	res, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query custom metric data: %s", metric.Name)
	}
	defer res.Close()
	result := make([]*pb.CustomMetricDataPoint, 0)
	for res.Next() {
		var predicate string
		var value sql.NullFloat64
		if err = res.Scan(&predicate, &value); err != nil {
			return nil, fmt.Errorf("failed to scan custom metric data: %s", metric.Name)
		}
		if value.Valid {
			result = append(result, &pb.CustomMetricDataPoint{Grouping: predicate, Value: float32(value.Float64)})
		}
	}
	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan custom metric data: %s", metric.Name)
	}
	return &pb.CustomMetricData{Metrics: result}, nil
}
