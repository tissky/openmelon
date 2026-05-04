package generation

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ShellProvider implements Provider by running a shell command.
// The request Prompt is written to the command's stdin; stdout is returned as content.
type ShellProvider struct {
	// Command is the executable and any fixed arguments (space-separated).
	Command string
	// Model is recorded in the Trace for provenance purposes.
	Model string
	// Env contains additional environment variables merged with os.Environ().
	Env map[string]string
}

// Generate runs the configured Command, writes req.Prompt to stdin,
// and returns stdout as the artifact content.
func (p *ShellProvider) Generate(ctx context.Context, req *Request) (string, *Trace, error) {
	start := time.Now()

	parts := strings.Fields(p.Command)
	if len(parts) == 0 {
		return "", nil, &ProviderError{Code: "invalid_command", Message: "ShellProvider.Command is empty"}
	}

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	cmd.Stdin = strings.NewReader(req.Prompt)

	env := os.Environ()
	for k, v := range p.Env {
		env = append(env, k+"="+v)
	}
	cmd.Env = env

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return "", nil, &ProviderError{
				Code:    "timeout",
				Message: fmt.Sprintf("command timed out: %s", p.Command),
				Wrapped: ctx.Err(),
			}
		}
		detail := strings.TrimSpace(stderr.String())
		if detail == "" {
			detail = err.Error()
		}
		return "", nil, &ProviderError{
			Code:    "non_zero_exit",
			Message: fmt.Sprintf("command failed: %s", detail),
			Wrapped: err,
		}
	}

	content := strings.TrimRight(stdout.String(), "\n")
	trace := &Trace{
		ProviderType: "shell",
		Model:        p.Model,
		Command:      p.Command,
		DurationSec:  time.Since(start).Seconds(),
	}
	return content, trace, nil
}
