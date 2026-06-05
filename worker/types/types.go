package types

import "time"

type CountryOrigin struct {
	Country string `json:"country"`
}

type Payload struct {
	Headers map[string]string
	Url     string
	Body    map[string]any
	Method  string
	Metrics ResponseMetrics
	Time    time.Time
}

type ResponseMetrics struct {
	StatusCode int
	Duration   int64
}
