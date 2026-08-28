package types

type CustomMetricDefinition struct {
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
	LastUpdate  string `json:"lastUpdate,omitempty"`
}

type MetricData struct {
	Name          string                   `json:"name"`
	Metrics       []MetricDataItems        `json:"metrics"`
	Timeframe     string                   `json:"timeframe"`
	CustomMetrics bool                     `json:"customMetrics"`
	Accumulate    bool                     `json:"accumulate"`
	Definition    *CustomMetricDefinition  `json:"definition,omitempty"`
}

type MetricDataResponse struct {
	Metrics []MetricData `json:"metrics"`
	Errors  []string     `json:"errors"`
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

type MetricDataItems struct {
	Grouping string  `json:"grouping"`
	Value    float64 `json:"value"`
}
