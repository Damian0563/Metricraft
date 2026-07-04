package db

import (
	"time"
)

const storageTimezone = "UTC"

func validTimezone(timezone string) string {
	if timezone == "" {
		return "UTC"
	}
	if _, err := time.LoadLocation(timezone); err != nil {
		return "UTC"
	}
	return timezone
}

func loadLocation(timezone string) *time.Location {
	loc, err := time.LoadLocation(validTimezone(timezone))
	if err != nil {
		return time.UTC
	}
	return loc
}

func alignStart(start time.Time, loc *time.Location, resolution int32) time.Time {
	local := start.In(loc)
	if resolution == 0 {
		y, m, d := local.Date()
		return time.Date(y, m, d, local.Hour(), 0, 0, 0, loc).UTC()
	}
	y, m, d := local.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, loc).UTC()
}

func rangeEnd(loc *time.Location, increment time.Duration, resolution int32) time.Time {
	end := time.Now().In(loc)
	if resolution == 0 || resolution == 1 {
		end = end.Add(increment).Truncate(increment)
	}
	return end.UTC()
}
