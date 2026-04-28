package project

// Project stores durable creative context for a content production workflow.
type Project struct {
	ID          string
	Name        string
	Platform    string
	Audience    string
	Persona     string
	Memory      map[string]string
	Constraints []string
}
