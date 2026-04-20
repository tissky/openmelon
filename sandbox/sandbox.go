// Copyright 2026 Point Eight AI Pte. Ltd.
// Licensed under the Apache License, Version 2.0

// Package sandbox provides gVisor-based sandboxed isolation for Skill execution.
//
// Python Skills run inside gVisor containers with:
//   - Network egress allowlist
//   - Memory hard cap (default 512MB)
//   - Per-skill timeout enforcement
//   - Read-only filesystem (except /tmp, cleared after execution)
//   - Static analysis rejection of eval/exec/os.system/subprocess
package sandbox

import "context"

// Runner executes a Skill in a sandboxed environment.
type Runner struct {
	// TODO: gVisor runtime configuration
}

// Run executes the given skill in a sandbox and returns the raw output.
func (r *Runner) Run(ctx context.Context, skillSlug string, input []byte) ([]byte, error) {
	// TODO: Implement sandboxed execution
	return nil, nil
}
