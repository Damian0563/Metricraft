package db

import (
	"fmt"
	"strings"

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

func resolveFieldName(customMetricSrc string) string {
	switch customMetricSrc {
	case "body":
		return "payload"
	case "header":
		return "headers"
	case "query":
		return "url"
	}
	return ""
}

func urlPrefixMatchSQL(ruleCol string) string {
	return fmt.Sprintf(`(url = %s OR (starts_with(url, %s) AND length(url) > length(%s) AND substring(url, length(%s) + 1, 1) = '/'))`,
		ruleCol, ruleCol, ruleCol, ruleCol)
}

func blacklistFilterSQL(param string) string {
	return fmt.Sprintf(`NOT EXISTS (
		SELECT 1 FROM unnest(%s::text[]) AS bl(rule)
		WHERE %s
	)`, param, urlPrefixMatchSQL("bl.rule"))
}

func groupedUrlSQL(groupingParam string) string {
	return fmt.Sprintf(`COALESCE(
		(
			SELECT g.rule
			FROM unnest(%s::text[]) WITH ORDINALITY AS g(rule, ord)
			WHERE %s
			ORDER BY ord
			LIMIT 1
		),
		url
	)`, groupingParam, urlPrefixMatchSQL("g.rule"))
}

func appendQueryParam(args []any, value any) ([]any, string) {
	args = append(args, value)
	return args, fmt.Sprintf("$%d", len(args))
}

func customMetricInnerWhere(pathParam, methodParam, groupingParam string, applyRules bool) string {
	if !applyRules {
		return fmt.Sprintf("url = %s AND method = %s", pathParam, methodParam)
	}
	parts := []string{
		urlPrefixMatchSQL(pathParam),
		fmt.Sprintf("method = %s", methodParam),
	}
	if groupingParam != "" {
		parts = append(parts, fmt.Sprintf(`EXISTS (
			SELECT 1 FROM unnest(%s::text[]) AS g(rule)
			WHERE %s
		)`, groupingParam, urlPrefixMatchSQL("g.rule")))
	}
	return strings.Join(parts, " AND ")
}
