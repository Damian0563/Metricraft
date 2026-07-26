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
	formattedMetric, err := json.Marshal(metric)
	if err != nil {
		return err
	}
	now := time.Now().In(loc)
	_, err = conn.Exec(ctx, "INSERT INTO custom_metrics (date, metric) VALUES ($1,$2)", now, string(formattedMetric))
	return err
}
