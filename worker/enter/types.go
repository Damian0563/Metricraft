package enter

type Request interface {
}

type Payload struct {
	Headers map[string]any
	Body    map[string]any
	Method  string
}

type ResponseMetrics struct {
	StatusCode int
	Duration   time.Duration
	Method     string
	URL        string
}
