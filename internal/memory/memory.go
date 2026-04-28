package memory

// Memory stores durable project facts and feedback.
type Memory struct {
	Facts    map[string]string
	Feedback []string
}
