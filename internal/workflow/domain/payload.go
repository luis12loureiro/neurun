package domain

type Payload interface {
	Type() TaskType
}

type LogPayload struct {
	Message string
}

func (l *LogPayload) Type() TaskType {
	return TaskTypeLog
}

type HTTPPayload struct {
	URL     string
	Method  string
	Body    string
	Headers map[string]string
	//Auth    *HTTPAuth
}

func (h *HTTPPayload) Type() TaskType {
	return TaskTypeHTTP
}
