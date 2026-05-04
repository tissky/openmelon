package workflow

import (
	"encoding/json"
	"fmt"
	"os"
)

// WorkflowDefinition describes a multi-stage content production workflow loaded from project.json.
type WorkflowDefinition struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Vertical string            `json:"vertical"`
	Stages   []StageDefinition `json:"stages"`
}

// StageDefinition describes a single stage within a WorkflowDefinition.
type StageDefinition struct {
	Stage            Stage             `json:"stage"`
	SkillPlusPackage string            `json:"skillplus_package"`
	CompileTarget    string            `json:"compile_target"`
	ModelProfile     string            `json:"model_profile"`
	ArtifactType     string            `json:"artifact_type"`
	Locale           string            `json:"locale,omitempty"`
	Vars             map[string]string `json:"vars,omitempty"`
}

// projectWorkflowsFile is a minimal struct used only for extracting the "workflows" key.
type projectWorkflowsFile struct {
	Workflows map[string]*WorkflowDefinition `json:"workflows"`
}

// LoadWorkflows reads the project JSON file at path and returns a map of WorkflowDefinitions.
// Each WorkflowDefinition has its ID field populated from the map key if not already set.
func LoadWorkflows(path string) (map[string]*WorkflowDefinition, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("workflow.LoadWorkflows: read %q: %w", path, err)
	}

	var pf projectWorkflowsFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return nil, fmt.Errorf("workflow.LoadWorkflows: parse %q: %w", path, err)
	}

	if len(pf.Workflows) == 0 {
		return nil, fmt.Errorf("workflow.LoadWorkflows: no workflows found in %q", path)
	}

	for key, wf := range pf.Workflows {
		if wf.ID == "" {
			wf.ID = key
		}
	}

	return pf.Workflows, nil
}
