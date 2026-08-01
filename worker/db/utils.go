package db

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	pb "metricraft/proto/metricraft/proto"
	"time"
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

var logsPool *pgxpool.Pool

func SetLogsPool(pool *pgxpool.Pool) {
	logsPool = pool
}

func getLogsPool() (*pgxpool.Pool, error) {
	if logsPool == nil {
		return nil, fmt.Errorf("logs database pool is not initialized")
	}
	return logsPool, nil
}

func blacklistedRules(rules []*pb.Rule) []string {
	var blacklisted []string
	for _, rule := range rules {
		if rule.Mode == "blacklisting" {
			blacklisted = append(blacklisted, rule.Rule)
		}
	}
	return blacklisted
}

func groupingRules(rules []*pb.Rule) []string {
	var grouping []string
	for _, rule := range rules {
		if rule.Mode == "grouping" {
			grouping = append(grouping, rule.Rule)
		}
	}
	return grouping
}

func getIncrementTruncPeriod(resolution int32) (time.Duration, string) {
	var increment time.Duration
	var truncPeriod string
	if resolution == 0 {
		increment = time.Hour
		truncPeriod = "hour"
	} else {
		increment = time.Hour * 24 * time.Duration(resolution)
		truncPeriod = "day"
	}
	return increment, truncPeriod
}

func urlPrefixMatchSQL(urlCol, ruleCol string) string {
	return fmt.Sprintf(`(%s = %s OR (starts_with(%s, %s) AND length(%s) > length(%s) AND substring(%s, length(%s) + 1, 1) = '/'))`,
		urlCol, ruleCol, urlCol, ruleCol, urlCol, ruleCol, urlCol, ruleCol)
}

func blacklistFilterSQL(urlCol, param string) string {
	return fmt.Sprintf(`NOT EXISTS (
		SELECT 1 FROM unnest(%s::text[]) AS bl(rule)
		WHERE %s
	)`, param, urlPrefixMatchSQL(urlCol, "bl.rule"))
}

func groupedUrlSQL(urlCol, groupingParam string) string {
	return fmt.Sprintf(`COALESCE(
		(
			SELECT g.rule
			FROM unnest(%s::text[]) WITH ORDINALITY AS g(rule, ord)
			WHERE %s
			ORDER BY ord
			LIMIT 1
		),
		%s
	)`, groupingParam, urlPrefixMatchSQL(urlCol, "g.rule"), urlCol)
}
