package onboard

// onboard.go — orchestrator. Ensure() runs all needed wizards in order:
//
//	1. Trust    (every launch on untrusted dirs)
//	2. Auth     (only when no API keys configured)
//	3. Project  (only when cwd is not in/under an openmelon project)
//
// Each step is skipped silently when its precondition is already met.
// Returns the project workdir for the caller to load.

import (
	"fmt"
	"os"
)

// Result describes what Ensure produced. Workdir is the project root
// the caller should hand to repl/tui. Quit signals the user explicitly
// declined a step (e.g. "No, quit" on the trust prompt) — caller
// should exit cleanly with no error.
type Result struct {
	Workdir string
	Quit    bool
}

// Ensure runs the full onboarding. Returns ASAP if the user declines
// any step.
func Ensure() (Result, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return Result{}, fmt.Errorf("onboard: cwd: %w", err)
	}

	trusted, err := EnsureTrust(cwd)
	if err != nil {
		return Result{}, err
	}
	if !trusted {
		return Result{Quit: true}, nil
	}

	configured, err := EnsureAuth()
	if err != nil {
		return Result{}, err
	}
	if !configured {
		return Result{Quit: true}, nil
	}

	wd, err := EnsureProject(cwd)
	if err != nil {
		return Result{}, err
	}
	if wd == "" {
		return Result{Quit: true}, nil
	}

	return Result{Workdir: wd}, nil
}
