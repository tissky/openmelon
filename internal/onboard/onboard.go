package onboard

// onboard.go — public entry point.
//
// Ensure() runs the full onboarding flow as a single bubbletea Program
// in alt-screen mode, transitioning smoothly between the trust, auth,
// and project-init screens. See orchestrator.go for the state machine.
//
// EnsureAuth and EnsureProject (formerly separate functions called from
// the old multi-Program Ensure) are kept as standalone helpers used by
// the `openmelon setup` and `openmelon project set-key` subcommands —
// those still want to drop into one specific wizard rather than the
// full onboarding flow.

// Ensure runs trust → auth → project-init in one bubbletea Program.
// Steps are skipped silently when their precondition is already met
// (trusted dir, configured API key, existing project).
//
// Returns ASAP if the user declines a step.
func Ensure() (Result, error) {
	return Run()
}
