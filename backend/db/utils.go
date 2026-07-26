package db

import (
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
	"time"
)

var usersPool *pgxpool.Pool

func SetUsersPool(pool *pgxpool.Pool) {
	usersPool = pool
}
func getUsersPool() (*pgxpool.Pool, error) {
	if usersPool == nil {
		return nil, errors.New("users database pool is not initialized")
	}
	return usersPool, nil
}

var logsPool *pgxpool.Pool

func SetLogsPool(pool *pgxpool.Pool) {
	logsPool = pool
}

func GetLogsPool() (*pgxpool.Pool, error) {
	if logsPool == nil {
		return nil, errors.New("logs database pool is not initialized")
	}
	return logsPool, nil
}

func urlMatchesRulePrefix(url, rule string) bool {
	if !strings.HasPrefix(url, rule) {
		return false
	}
	if len(url) == len(rule) {
		return true
	}
	return url[len(rule)] == '/'
}

var ErrInvalidTimezone = errors.New("invalid timezone")

func getLocation(tz string) (*time.Location, error) {
	tz = strings.TrimSpace(tz)
	if tz == "" {
		tz = "UTC"
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, errors.Join(ErrInvalidTimezone, err)
	}
	return loc, nil
}
