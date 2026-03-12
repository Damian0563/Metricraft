package enter

type Request interface {
}

type Payload struct {
	Headers map[string]any
	Body    map[string]any
	Method  string
}
