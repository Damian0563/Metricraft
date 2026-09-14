package db

import (
	"backend/types"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"math"
	"slices"
)

func VerifyAppName(appName string) bool {
	ctx := context.Background()
	conn, err := GetLogsPool()
	if err != nil {
		return false
	}
	var exists bool
	err = conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM settings WHERE appname = $1)", appName).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func GetLogCapacity() (float64, error) {
	conn, err := GetLogsPool()
	if err != nil {
		return 0, err
	}
	var sizeBytes int64
	err = conn.QueryRow(context.Background(), "SELECT pg_total_relation_size('logs')").Scan(&sizeBytes)
	if err != nil {
		return 0, err
	}
	return float64(sizeBytes) / 1024 / 1024, nil
}

func DeleteLogs(ctx context.Context, toDeleteMbs float64) (float64, error) {
	conn, err := GetLogsPool()
	if err != nil {
		return 0, err
	}
	var sizeBytes, rowCount int64
	if err = conn.QueryRow(ctx, "SELECT pg_total_relation_size('logs'), (SELECT COUNT(*) FROM logs)").Scan(&sizeBytes, &rowCount); err != nil {
		return 0, err
	}
	if rowCount == 0 {
		return float64(sizeBytes) / 1024 / 1024, nil
	}
	avgRowBytes := float64(sizeBytes) / float64(rowCount)
	rowsToDelete := int64(math.Ceil(toDeleteMbs * 1024 * 1024 / avgRowBytes))
	rowsToDelete = int64(math.Min(float64(rowCount), float64(rowsToDelete)))
	if _, err = conn.Exec(ctx, "DELETE FROM logs WHERE ctid IN (SELECT ctid FROM logs ORDER BY date ASC NULLS FIRST LIMIT $1)", rowsToDelete); err != nil {
		return 0, err
	}
	if _, err = conn.Exec(ctx, "VACUUM FULL logs"); err != nil {
		return 0, err
	}
	return GetLogCapacity()
}

func AddRule(ctx context.Context, rule types.Rule) error {
	conn, err := GetLogsPool()
	if err != nil {
		return err
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var exists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM rules WHERE rule = $1)", rule.Rule).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		_, err = tx.Exec(ctx, "DELETE FROM rules WHERE rule = $1", rule.Rule)
		if err != nil {
			return err
		}
	}
	if _, err = tx.Exec(ctx, "INSERT INTO rules (rule, matches, mode) VALUES ($1, $2, $3)", rule.Rule, rule.Matches, rule.Mode); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func reComputeMatches(conn *pgxpool.Pool, ctx context.Context, rules []types.Rule) ([]types.Rule, error) {
	res, err := conn.Query(ctx, "SELECT DISTINCT url from logs")
	if err != nil {
		return nil, err
	}
	defer res.Close()
	for res.Next() {
		var url string
		if err := res.Scan(&url); err != nil {
			return nil, err
		}
		for i, rule := range rules {
			if urlMatchesRulePrefix(url, rule.Rule) && !slices.Contains(rule.Matches, url) {
				rules[i].Matches = append(rules[i].Matches, url)
			}
		}
	}
	return rules, nil
}

func GetRules(ctx context.Context) ([]types.Rule, error) {
	conn, err := GetLogsPool()
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query(ctx, "SELECT rule, matches, mode FROM rules")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rules := []types.Rule{}
	for rows.Next() {
		var rule types.Rule
		if err := rows.Scan(&rule.Rule, &rule.Matches, &rule.Mode); err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	if rules, err = reComputeMatches(conn, ctx, rules); err != nil {
		return nil, err
	}
	return rules, nil
}

func DeleteRule(ctx context.Context, rule types.Rule) error {
	conn, err := GetLogsPool()
	if err != nil {
		return err
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, "DELETE FROM rules WHERE rule = $1", rule.Rule)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
