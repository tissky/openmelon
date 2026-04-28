package review

// Result stores review output for an artifact.
type Result struct {
	ArtifactID string
	Passed     bool
	Notes      []string
	Labels     map[string]string
}
