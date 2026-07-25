package types

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
