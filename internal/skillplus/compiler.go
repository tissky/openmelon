package skillplus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CompileRequest holds the parameters for a Skill-Plus compilation.
type CompileRequest struct {
	PackagePath  string
	Target       string
	ModelProfile string
	Locale       string
	Vars         map[string]string
}

// Compiler invokes the Skill-Plus Python reference compiler.
type Compiler struct {
	// CompilerPath is the directory added to PYTHONPATH (where skillplus_compile module lives).
	CompilerPath string
	// PythonCmd overrides the Python executable used (default: "python3"). Useful in tests.
	PythonCmd string
}

// rawOutput is the intermediate struct matching the Python compiler JSON output.
type rawOutput struct {
	Target  string `json:"target"`
	Package struct {
		ID      string `json:"id"`
		Version string `json:"version"`
	} `json:"package"`
	CompiledPrompt string            `json:"compiled_prompt"`
	RuntimeVars    map[string]string `json:"runtime_vars"`
	ModelProfile   string            `json:"model_profile"`
	Evaluation     struct {
		Checklist []string `json:"checklist"`
	} `json:"evaluation"`
}

// Compile invokes the Skill-Plus compiler subprocess and returns a CompiledSkill.
func (c *Compiler) Compile(ctx context.Context, req *CompileRequest) (*CompiledSkill, error) {
	pythonCmd := c.pythonCmd()

	if _, err := exec.LookPath(pythonCmd); err != nil {
		return nil, fmt.Errorf(
			"%q not found in PATH — install Python 3.9+ and ensure it is in PATH: %w",
			pythonCmd, err,
		)
	}

	args := []string{"-m", "skillplus_compile", req.PackagePath,
		"--target", req.Target,
		"--model-profile", req.ModelProfile,
	}
	if req.Locale != "" {
		args = append(args, "--locale", req.Locale)
	}
	for k, v := range req.Vars {
		args = append(args, "--var", k+"="+v)
	}

	cmd := exec.CommandContext(ctx, pythonCmd, args...)
	cmd.Env = append(os.Environ(), "PYTHONPATH="+filepath.Clean(c.CompilerPath))

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		detail := strings.TrimSpace(stderr.String())
		if detail == "" {
			detail = err.Error()
		}
		return nil, fmt.Errorf("skillplus compile failed for %q: %s", req.PackagePath, detail)
	}

	var raw rawOutput
	if err := json.Unmarshal(stdout.Bytes(), &raw); err != nil {
		return nil, fmt.Errorf("skillplus: failed to parse compiler output: %w", err)
	}

	return &CompiledSkill{
		PackageID:      raw.Package.ID,
		PackageVersion: raw.Package.Version,
		Target:         raw.Target,
		ModelProfile:   raw.ModelProfile,
		RuntimeVars:    raw.RuntimeVars,
		Prompt:         raw.CompiledPrompt,
		Evaluation:     raw.Evaluation.Checklist,
	}, nil
}

func (c *Compiler) pythonCmd() string {
	if c.PythonCmd != "" {
		return c.PythonCmd
	}
	return "python3"
}
