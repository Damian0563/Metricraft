package db

import (
	"context"
	"encoding/json"
	"time"
)

type MetricOrchestrator interface {
	Initialize(context.Context, string) error
	PrepareData(context.Context, string) (string, error) //stringified data and error
	Delete(context.Context, string) error
	Edit(context.Context, string) error
}

type CustomMetric struct {
	Name        string `json:"name"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	Source      string `json:"source"`
	Selector    string `json:"selector"`
	Aggregation string `json:"aggregation"`
	Timeframe   string `json:"timeframe"`
	ValueType   string `json:"valueType"`
	ApplyRules  bool   `json:"applyRules"`
	ChartType   string `json:"chartType"`
	LastUpdate  string `json:"lastUpdate"`
}

func (m *CustomMetric) Initialize(ctx context.Context, tz string) error {
	return dbAddCustomMetric(ctx, m, tz)
}

func (m *CustomMetric) PrepareData(ctx context.Context) (string, error) {
	return "", nil
}

func (m *CustomMetric) Delete(ctx context.Context) error {
	return nil
}

func (m *CustomMetric) Edit(ctx context.Context) error {
	return nil
}

func dbAddCustomMetric(ctx context.Context, metric *CustomMetric, tz string) error {
	conn, err := GetLogsPool()
	if err != nil {
		return err
	}
	loc, err := getLocation(tz)
	if err != nil {
		return err
	}
	metric.LastUpdate = ""
	formattedMetric, err := json.Marshal(metric)
	if err != nil {
		return err
	}
	now := time.Now().In(loc)
	_, err = conn.Exec(ctx, "INSERT INTO custom_metrics (date, metric) VALUES ($1,$2)", now, string(formattedMetric))
	return err
}

func ListMetrics(ctx context.Context, tz string) ([]CustomMetric, error) {
	conn, err := GetLogsPool()
	if err != nil {
		return nil, err
	}
	var metrics []CustomMetric
	rows, err := conn.Query(ctx, "SELECT date, metric FROM custom_metrics WHERE TRUE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	loc, err := getLocation(tz)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var metric CustomMetric
		var date time.Time
		var rawMetric string
		if err := rows.Scan(&date, &rawMetric); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(rawMetric), &metric); err != nil {
			return nil, err
		}
		metric.LastUpdate = date.In(loc).Format("Jan 2, 2006 at 3:04pm")
		metrics = append(metrics, metric)
	}
	return metrics, nil
}
