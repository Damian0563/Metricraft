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

func resolveAggregationType(aggregation string, inspectField string) string {
	switch aggregation {
	case "sum":
		return fmt.Sprintf("SUM (%s)", inspectField)
	case "avg":
		return fmt.Sprintf("AVG (%s)", inspectField)
	case "min":
		return fmt.Sprintf("MIN (%s)", inspectField)
	case "max":
		return fmt.Sprintf("MAX (%s)", inspectField)
	case "p95":
		return fmt.Sprintf("percentile_cont(0.95) WITHIN GROUP (ORDER BY %s)", inspectField)
	case "p50":
		return fmt.Sprintf("percentile_cont(0.5) WITHIN GROUP (ORDER BY %s)", inspectField)
	case "unique":
		return fmt.Sprintf("COUNT(DISTINCT %s)", inspectField)
	default:
		return "COUNT (*)"
	}
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

// jsonPathSelector turns a dot-notation selector such as items[0].price into a
// postgres jsonpath expression ($."items"[0]."price"). Keys are always quoted so
// that names containing dashes or digits stay valid.
var replacer = strings.NewReplacer(`\`, `\\`, `"`, `\"`)

func jsonPathSelector(selector string) string {
	var path strings.Builder
	path.WriteString("$")
	sequence := strings.Split(selector, ".")
	for _, segment := range sequence {
		if segment == "" {
			continue
		}
		key := segment
		indexes := ""
		if open := strings.Index(segment, "["); open >= 0 {
			key = segment[:open]
			indexes = segment[open:]
		}
		if key != "" {
			path.WriteString(`."`)
			path.WriteString(replacer.Replace(key))
			path.WriteString(`"`)
		}
		path.WriteString(indexes)
	}
	return path.String()
}

// customMetricLogicMatch builds the predicate deciding whether a log row carries
// the value a custom metric is tracking: a jsonpath into the stored JSON of the
// request body or headers, or a query string parameter of the url.
func customMetricLogicMatch(inspectedField, selectorParam string) string {
	if selectorParam == "" {
		return "TRUE"
	}
	switch inspectedField {
	case "payload", "headers":
		return fmt.Sprintf("%s @? %s::jsonpath", inspectedField, selectorParam)
	default:
		return fmt.Sprintf("(strpos(url, '?' || %s || '=') > 0 OR strpos(url, '&' || %s || '=') > 0)",
			selectorParam, selectorParam)
	}
}
