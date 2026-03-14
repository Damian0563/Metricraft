package enter

import "time"

type Request interface {
}

type Payload struct {
	Headers map[string]any
	Url     string
	Body    map[string]any
	Method  string
	Metrics ResponseMetrics
}

type ResponseMetrics struct {
	StatusCode int
	Duration   time.Duration
}
