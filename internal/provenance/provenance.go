package provenance

// Record describes how an artifact was produced.
type Record struct {
	ArtifactID       string            `json:"artifact_id"`
	ProjectID        string            `json:"project_id"`
	WorkflowID       string            `json:"workflow_id"`
	Stage            string            `json:"stage"`
	SkillPackage     string            `json:"skill_package"`
	CompiledTarget   string            `json:"compiled_target"`
	Model            string            `json:"model"`
	PromptHash       string            `json:"prompt_hash"`
	GenerationParams map[string]string `json:"generation_params,omitempty"`
	EvaluationResult string            `json:"evaluation_result,omitempty"`
	Timestamp        string            `json:"timestamp"`
}
