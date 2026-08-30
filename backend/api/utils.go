package api

import (
	"backend/types"
	pb "metricraft/proto/metricraft/proto"
	"strconv"
	"strings"
	"time"
)

var mapping = map[string]string{
	"Last 12 hours": "0.5d",
	"Last 24 hours": "1d",
	"Last 7 days":   "7d",
	"Last 30 days":  "30d",
	"Last 90 days":  "90d",
	"Last 180 days": "180d",
	"Last 365 days": "365d",
	"This week":     "7t",
	"This month":    "30t",
	"This year":     "365t",
}

func standardizeTimeframe(inputTimeframe string) string {
	if inputTimeframe == "" {
		return "7d"
	}
	if val, ok := mapping[inputTimeframe]; ok {
		return val
	}
	for _, code := range mapping {
		if code == inputTimeframe {
			return inputTimeframe
		}
	}
	return "7d"
}

func customMetricDataFromProto(resp *pb.CustomMetricData) []types.MetricDataItems {
	if resp == nil || resp.Metrics == nil {
		return nil
	}
	mappings := make([]types.MetricDataItems, 0, len((*resp.Metrics).Metrics))
	for _, mapping := range (*resp.Metrics).Metrics {
		mappings = append(mappings, types.MetricDataItems{Grouping: mapping.Grouping, Value: float64(mapping.Value)})
	}
	return mappings
}

func convertTimeframe(inputTimeframe string, timezone string) (time.Time, types.ResolutionDays) {
	inputTimeframe = standardizeTimeframe(inputTimeframe)
	numFloat := float32(7)
	if len(inputTimeframe) > 1 {
		unit := inputTimeframe[len(inputTimeframe)-1]
		if unit == 'd' || unit == 't' {
			if num64, err := strconv.ParseFloat(inputTimeframe[:len(inputTimeframe)-1], 32); err == nil {
				numFloat = float32(num64)
			}
		}
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
		midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		switch numFloat {
		case 7:
			weekday := int(now.Weekday()) - 1
			if weekday < 0 {
				weekday = 6
			}
			start = midnight.AddDate(0, 0, -weekday)
		case 30:
			start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		default: //365
			start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, loc)
		}
	} else {
		start = now.Add(-time.Duration(int(numFloat)) * 24 * time.Hour)
	}
	return start.UTC(), timeResolution[numFloat]
}
