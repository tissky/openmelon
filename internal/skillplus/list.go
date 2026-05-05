package skillplus

// list.go — wrapper around `skillplus list --json` for the TUI's
// /skill picker. Shelling out (vs. reading catalog.json + ~/.skillplus/
// directly) keeps openmelon decoupled from skillplus's on-disk layout
// — the CLI is the contract.

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
)

// SkillInfo is one entry from `skillplus list --json`. Same shape the
// skillplus catalog uses; openmelon only cares about the few fields
// below for the picker UI.
type SkillInfo struct {
	ID          string   `json:"id"`
	Version     string   `json:"version"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Path        string   `json:"path"`
	Source      string   `json:"source"` // "local" or "bundled"
}

// ListSkills runs `skillplus list --json` and parses the result.
// Returns an empty slice (not an error) when the skillplus CLI isn't
// installed — the picker just shows "no skills found" in that case.
func ListSkills(ctx context.Context) ([]SkillInfo, error) {
	if _, err := exec.LookPath("skillplus"); err != nil {
		return nil, nil
	}
	cmd := exec.CommandContext(ctx, "skillplus", "list", "--json")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("skillplus list: %w", err)
	}
	var skills []SkillInfo
	if err := json.Unmarshal(out, &skills); err != nil {
		return nil, fmt.Errorf("skillplus list: parse: %w", err)
	}
	return skills, nil
}
