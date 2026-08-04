package types

type MetricData struct {
	Name          string                              `json:"name"`
	Metrics       []MetricAggregatorResolutionMapping `json:"metrics"`
	Timeframe     string                              `json:"timeframe"`
	CustomMetrics bool                                `json:"customMetrics"`
}

type Rule struct {
	Rule    string   `json:"rule"`
	Matches []string `json:"matches"`
	Mode    string   `json:"mode"`
}

type Worker struct {
	Url          string            `json:"url"`
	PollInterval int               `json:"pollInterval"`
	Headers      map[string]string `json:"headers,omitempty"`
}

type EnabledMetric struct {
	Enabled   bool   `json:"enabled"`
	Timeframe string `json:"timeframe"`
}

type Metric struct {
	Name      string `json:"name"`
	Enabled   bool   `json:"enabled"`
	Timeframe string `json:"timeframe"`
}

type MetricAggregatorResolutionMapping struct {
	Data []MetricAggregatorData `json:"data"`
}

type MetricAggregatorData struct {
	Timerange string  `json:"timerange"`
	Data      float64 `json:"value"`
}
