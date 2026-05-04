package skillplus

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// validCompiledSkillJSON is a minimal valid response from the Python compiler.
const validCompiledSkillJSON = `{
	"target": "openmelon",
	"package": {"id": "food-street-realism", "version": "1.0.0"},
	"compiled_prompt": "test compiled prompt",
	"runtime_vars": {"realism_level": "high"},
	"model_profile": "gpt-image-family",
	"evaluation": {"checklist": ["check sharpness", "check realism"]}
}`

func TestCompiler_pythonNotFound(t *testing.T) {
	c := &Compiler{
		CompilerPath: "/fake",
		PythonCmd:    "/nonexistent/python99",
	}
	req := &CompileRequest{
		PackagePath:  "/some/package",
		Target:       "openmelon",
		ModelProfile: "gpt-image-family",
	}
	_, err := c.Compile(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for missing python, got nil")
	}
	if !strings.Contains(err.Error(), "not found in PATH") {
		t.Errorf("expected 'not found in PATH' in error, got: %v", err)
	}
}

func TestCompiler_successMockExec(t *testing.T) {
	// Create a fake python3 script that outputs valid CompiledSkill JSON.
	tmpDir := t.TempDir()
	fakePython := filepath.Join(tmpDir, "fake_python3")
	script := "#!/bin/sh\ncat << 'ENDJSON'\n" + validCompiledSkillJSON + "\nENDJSON\n"
	if err := os.WriteFile(fakePython, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	c := &Compiler{
		CompilerPath: tmpDir,
		PythonCmd:    fakePython,
	}
	req := &CompileRequest{
		PackagePath:  "/some/food.skillplus",
		Target:       "openmelon",
		ModelProfile: "gpt-image-family",
		Locale:       "zh-CN",
		Vars:         map[string]string{"realism_level": "high"},
	}

	got, err := c.Compile(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.PackageID != "food-street-realism" {
		t.Errorf("PackageID = %q, want %q", got.PackageID, "food-street-realism")
	}
	if got.Prompt != "test compiled prompt" {
		t.Errorf("Prompt = %q, want %q", got.Prompt, "test compiled prompt")
	}
	if len(got.Evaluation) != 2 {
		t.Errorf("Evaluation len = %d, want 2", len(got.Evaluation))
	}
}
