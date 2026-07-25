package db

import (
	"backend/types"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
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
