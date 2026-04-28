package provenance

// Record describes how an artifact was produced.
type Record struct {
	ProjectID        string
	WorkflowID       string
	Stage            string
	SkillPackage     string
	CompiledTarget   string
	Model            string
	Prompt           string
	GenerationParams map[string]string
	EvaluationResult string
}
