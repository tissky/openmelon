package generation

// Request describes a generation request for a model or tool adapter.
type Request struct {
	ArtifactType string
	Prompt       string
	Model        string
	Params       map[string]string
}
