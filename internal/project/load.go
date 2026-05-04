package project

import (
	"encoding/json"
	"fmt"
	"os"
)

// Load reads a project JSON file from path, validates required fields, and returns a Project.
// Returns a wrapped os.ErrNotExist error if the file does not exist.
// Returns a descriptive error if required fields are missing.
func Load(path string) (*Project, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("load project %q: %w", path, err)
	}

	var p Project
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("parse project %q: %w", path, err)
	}

	if err := validate(&p); err != nil {
		return nil, fmt.Errorf("invalid project %q: %w", path, err)
	}

	return &p, nil
}

func validate(p *Project) error {
	if p.ID == "" {
		return fmt.Errorf("missing required field: id")
	}
	if p.Name == "" {
		return fmt.Errorf("missing required field: name")
	}
	if p.Platform == "" {
		return fmt.Errorf("missing required field: platform")
	}
	return nil
}
