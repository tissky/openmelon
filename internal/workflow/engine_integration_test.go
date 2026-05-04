//go:build integration

package workflow_test

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Jackyffight/openmelon/internal/generation"
	"github.com/Jackyffight/openmelon/internal/project"
	"github.com/Jackyffight/openmelon/internal/provenance"
	"github.com/Jackyffight/openmelon/internal/skillplus"
	"github.com/Jackyffight/openmelon/internal/workflow"
)

// TestEngine_runEndToEnd runs the engine with:
//   - a fake Python3 compiler script that returns a valid CompiledSkill JSON
//   - a ShellProvider using "cat" (echoes stdin to stdout as the generated content)
//   - project.json from examples/food-exploration/project.json
//
// It verifies that artifact files and provenance.jsonl are created correctly.
func TestEngine_runEndToEnd(t *testing.T) {
	// --- Setup: fake Python compiler ---
	tmpDir := t.TempDir()
	fakePython := filepath.Join(tmpDir, "fake_python3")
	fakeCompilerOutput := `{
		"target": "openmelon",
		"package": {"id": "food-street-realism", "version": "1.0.0"},
		"compiled_prompt": "A vivid street food photo",
		"runtime_vars": {"realism_level": "high"},
		"model_profile": "image_generator",
		"evaluation": {"checklist": ["check focus", "check lighting"]}
	}`
	script := "#!/bin/sh\ncat << 'ENDJSON'\n" + fakeCompilerOutput + "\nENDJSON\n"
	if err := os.WriteFile(fakePython, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	// --- Setup: artifact and provenance dirs ---
	artifactDir := filepath.Join(tmpDir, "artifacts")
	provPath := filepath.Join(tmpDir, "provenance.jsonl")

	// --- Load project ---
	// Use the canonical example project relative to workspace root.
	projectPath := "../../../examples/food-exploration/project.json"
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		t.Skip("example project.json not found, skipping integration test")
	}
	proj, err := project.Load(projectPath)
	if err != nil {
		t.Fatalf("project.Load: %v", err)
	}

	// --- Load workflows ---
	workflows, err := workflow.LoadWorkflows(projectPath)
	if err != nil {
		t.Fatalf("workflow.LoadWorkflows: %v", err)
	}
	var wfDef *workflow.WorkflowDefinition
	for _, wf := range workflows {
		wfDef = wf
		break
	}
	if wfDef == nil {
		t.Fatal("no workflow found")
	}

	// --- Build request ---
	compiler := &skillplus.Compiler{
		CompilerPath: tmpDir,
		PythonCmd:    fakePython,
	}
	provider := &generation.ShellProvider{Command: "cat"} // echoes stdin (prompt) back as output

	req := &workflow.RunRequest{
		Project:        proj,
		WorkflowDef:    wfDef,
		Intent:         "show a real late-night noodle shop vibe",
		ArtifactDir:    artifactDir,
		CompilerPath:   tmpDir,
		ProvenancePath: provPath,
		Compiler:       compiler,
		Provider:       provider,
		Generate:       true,
	}

	// --- Run ---
	engine := &workflow.Engine{}
	results, err := engine.Run(context.Background(), req)
	if err != nil {
		t.Fatalf("engine.Run: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least one stage result")
	}

	// --- Verify artifact files ---
	for _, r := range results {
		if r.Artifact == nil {
			t.Errorf("stage %q: artifact is nil", r.Stage)
			continue
		}
		// Content file should exist
		entries, err := os.ReadDir(artifactDir)
		if err != nil {
			t.Fatalf("read artifact dir: %v", err)
		}
		found := false
		for _, e := range entries {
			if len(e.Name()) > 0 && e.Name()[:16] == r.Artifact.ID[:min(16, len(r.Artifact.ID))] {
				found = true
			}
		}
		if !found {
			t.Errorf("stage %q: no artifact file found for id %s", r.Stage, r.Artifact.ID)
		}

		// Provenance file for artifact should exist
		provFile := filepath.Join(artifactDir, r.Artifact.ID+".provenance.json")
		if _, err := os.Stat(provFile); err != nil {
			t.Errorf("stage %q: provenance file missing: %v", r.Stage, err)
		}
	}

	// --- Verify provenance.jsonl ---
	f, err := os.Open(provPath)
	if err != nil {
		t.Fatalf("provenance.jsonl not created: %v", err)
	}
	defer f.Close()

	var records []provenance.Record
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var rec provenance.Record
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			t.Fatalf("unmarshal provenance line: %v", err)
		}
		records = append(records, rec)
	}

	if len(records) != len(results) {
		t.Errorf("provenance records = %d, want %d", len(records), len(results))
	}
	for _, rec := range records {
		if rec.ArtifactID == "" {
			t.Error("provenance record missing artifact_id")
		}
		if rec.ProjectID != proj.ID {
			t.Errorf("provenance project_id = %q, want %q", rec.ProjectID, proj.ID)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
