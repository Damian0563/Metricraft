package db

import (
	"backend/types"
	"context"
	"encoding/json"
	"sync"
	"time"
)

type MetricOrchestrator interface {
	Initialize(context.Context) error
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

func (m *CustomMetric) Initialize(ctx context.Context) error {
	return dbAddCustomMetric(ctx, m)
}

func (m *CustomMetric) PrepareData(ctx context.Context, loc *time.Location) (map[string]string, error) {
	return make(map[string]string), nil
}

func (m *CustomMetric) Delete(ctx context.Context) error {
	return dbDeleteCustomMetric(ctx, m)
}

func (m *CustomMetric) Edit(ctx context.Context, updated CustomMetric) error {
	return dbEditCustomMetric(ctx, m, updated)
}

func dbEditCustomMetric(ctx context.Context, original *CustomMetric, updated CustomMetric) error {
	conn, err := GetLogsPool()
	if err != nil {
		return err
	}
	original.LastUpdate = ""
	var formattedOriginal []byte
	if formattedOriginal, err = json.Marshal(original); err != nil {
		return err
	}
	updated.LastUpdate = ""
	var formattedUpdated []byte
	if formattedUpdated, err = json.Marshal(updated); err != nil {
		return err
	}
	now := time.Now().UTC()
	_, err = conn.Exec(ctx, "UPDATE custom_metrics SET metric=$1, date=$2 WHERE metric=$3", string(formattedUpdated), now, string(formattedOriginal))
	return err
}

func dbDeleteCustomMetric(ctx context.Context, metric *CustomMetric) error {
	conn, err := GetLogsPool()
	if err != nil {
		return err
	}
	metric.LastUpdate = ""
	var formattedMetric []byte
	if formattedMetric, err = json.Marshal(metric); err != nil {
		return err
	}
	_, err = conn.Exec(ctx, "DELETE FROM custom_metrics WHERE metric=$1", string(formattedMetric))
	return err
}

func dbAddCustomMetric(ctx context.Context, metric *CustomMetric) error {
	conn, err := GetLogsPool()
	if err != nil {
		return err
	}
	metric.LastUpdate = ""
	formattedMetric, err := json.Marshal(metric)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	_, err = conn.Exec(ctx, "INSERT INTO custom_metrics (date, metric) VALUES ($1,$2)", now, string(formattedMetric))
	return err
}

func ListMetrics(ctx context.Context) ([]CustomMetric, error) {
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
		metric.LastUpdate = date.UTC().Format(time.RFC3339)
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

func GetCustomMetricsData(ctx context.Context, result *[]types.MetricData, loc *time.Location) error {
	conn, err := GetLogsPool()
	if err != nil {
		return err
	}
	rows, err := conn.Query(ctx, "SELECT metric FROM custom_metrics")
	if err != nil {
		return err
	}
	defer rows.Close()
	var wg sync.WaitGroup
	for rows.Next() {
		var rawMetric string
		if err := rows.Scan(&rawMetric); err != nil {
			return err
		}
		var metric CustomMetric
		if err := json.Unmarshal([]byte(rawMetric), &metric); err != nil {
			return err
		}
		var mu sync.Mutex
		go func(metric *CustomMetric, metrics *[]types.MetricData) {
			defer wg.Done()
			data, err := metric.PrepareData(ctx, loc)
			if err != nil {
				return
			}
			mu.Lock()
			*result = append(*result, types.MetricData{Name: metric.Name, Metrics: data, Timeframe: metric.Timeframe, CustomMetrics: true})
			mu.Unlock()

		}(&metric, result)
		wg.Add(1)
	}
	wg.Wait()
	return nil
}
