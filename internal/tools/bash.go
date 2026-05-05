package tools

// bash.go — the bash tool. Lets the agent run shell commands inside
// the project workdir, with mandatory user approval before execution.
//
// Why approval is mandatory:
//   - The model is good at "let me check if this image opened
//     correctly" type investigations, but bad at "wait, that command
//     would also delete my downloads folder."
//   - The whole point of having bash is interactive triage; if we
//     run anything the model emits, we lose the human-in-the-loop
//     property that makes bash safe.
//   - There's no useful sandbox we can apply that's both cheap and
//     correct. Approval is.
//
// When env.Approve is nil (which happens in headless / CI / piped
// stdin contexts where there's no UI to render the modal), bash
// returns an error the model can read instead of running the command.

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

func bashTool(env *Env) Tool {
	return Tool{
		Spec: Spec{
			Name: "bash",
			Description: "Run a shell command inside the project workdir and return its combined stdout/stderr. " +
				"Use sparingly — for inspecting files (file, ls, du), checking output (open, identify), " +
				"or quick text edits the user expects. Each call requires explicit user confirmation.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"command": {
						"type": "string",
						"description": "The shell command to run. Will be passed to /bin/sh -c."
					},
					"description": {
						"type": "string",
						"description": "One-line plain-English explanation of why you're running this. Shown to the user in the approval prompt."
					},
					"timeout_seconds": {
						"type": "number",
						"description": "Kill the command after this many seconds. Default 30, max 300."
					}
				},
				"required": ["command", "description"]
			}`),
		},
		Handler: func(ctx context.Context, raw json.RawMessage) (any, error) {
			var args struct {
				Command        string  `json:"command"`
				Description    string  `json:"description"`
				TimeoutSeconds float64 `json:"timeout_seconds"`
			}
			if err := json.Unmarshal(raw, &args); err != nil {
				return nil, fmt.Errorf("invalid args: %w", err)
			}
			if args.Command == "" {
				return map[string]any{"error": "command is required"}, nil
			}
			timeout := time.Duration(args.TimeoutSeconds) * time.Second
			if timeout <= 0 {
				timeout = 30 * time.Second
			}
			if timeout > 5*time.Minute {
				timeout = 5 * time.Minute
			}

			if env.Approve == nil {
				return map[string]any{
					"error": "bash is unavailable: no approval gate is wired (running headless?)",
				}, nil
			}
			ok := env.Approve(ApprovalRequest{
				Tool:        "bash",
				Command:     args.Command,
				Description: args.Description,
			})
			if !ok {
				return map[string]any{"error": "user denied execution"}, nil
			}

			execCtx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			cmd := exec.CommandContext(execCtx, "/bin/sh", "-c", args.Command)
			cmd.Dir = env.Workdir
			out, err := cmd.CombinedOutput()
			res := map[string]any{
				"stdout":    string(out),
				"exit_code": cmd.ProcessState.ExitCode(),
			}
			if err != nil {
				if execCtx.Err() == context.DeadlineExceeded {
					res["error"] = fmt.Sprintf("timed out after %s", timeout)
				} else {
					// exec.ExitError is normal — the model wants the
					// exit code + output regardless.
					if _, isExit := err.(*exec.ExitError); !isExit {
						res["error"] = err.Error()
					}
				}
			}
			return res, nil
		},
	}
}
