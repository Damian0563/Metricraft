package db

import (
	"backend/types"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
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

func ConvertTimeframe(inputTimeframe string, timezone string) (time.Time, types.ResolutionDays) {
	if inputTimeframe == "" {
		inputTimeframe = "7d"
	}
	numFloat := float32(7)
	timeframe := strings.ReplaceAll(inputTimeframe, "d", "")
	timeframe = strings.ReplaceAll(timeframe, "t", "")
	if num64, err := strconv.ParseFloat(timeframe, 32); err == nil {
		numFloat = float32(num64)
	}
	loc := time.UTC
	if timezone != "" {
		if loaded, loadErr := time.LoadLocation(timezone); loadErr == nil {
			loc = loaded
		}
	}
	timeResolution := map[float32]types.ResolutionDays{
		0.5: {Days: -1},
		1:   {Days: 0},
		7:   {Days: 1},
		14:  {Days: 2},
		30:  {Days: 3},
		90:  {Days: 7},
		180: {Days: 14},
		365: {Days: 30},
	}
	now := time.Now().In(loc)
	var start time.Time
	if strings.Contains(inputTimeframe, "d") {
		switch numFloat {
		case 1:
			start = now.Add(-23 * time.Hour)
		case 0.5:
			start = now.Add(-11 * time.Hour)
		default:
			start = now.Add(-time.Duration(int(numFloat)) * 24 * time.Hour)
		}
	} else if strings.Contains(inputTimeframe, "t") {
		switch numFloat {
		case 7:
			weekday := int(now.Weekday()) - 1
			if weekday < 0 {
				weekday = 6
			}
			start = now.Add(time.Duration(weekday) * -23 * time.Hour)
		case 30:
			monthday := now.Day()
			start = now.Add(time.Duration(monthday) * -23 * time.Hour)
		default: //365
			yearday := now.YearDay()
			start = now.Add(time.Duration(yearday) * -23 * time.Hour)
		}
	} else {
		start = now.Add(-time.Duration(int(numFloat)) * 24 * time.Hour)
	}
	return start.UTC(), timeResolution[numFloat]
}
