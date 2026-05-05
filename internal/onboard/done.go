package onboard

// done.go — wizard-completion signaling.
//
// Sub-models (list, key input, model input, project field) used to call
// tea.Quit when the user pressed Enter / Ctrl+C. That kills the
// bubbletea Program — fine when each wizard runs in its own Program,
// but useless once we want one Program to host every wizard in
// sequence (no inter-wizard screen flash).
//
// The new contract: sub-models emit wizardDoneMsg via a tea.Cmd. The
// orchestrator (orchestrator.go) intercepts these and transitions to
// the next state. Standalone callers (runProjectKeyWizard,
// `openmelon setup`) wrap the sub-model in singleShotRunner, which
// translates wizardDoneMsg to tea.Quit and captures the result.

import tea "github.com/charmbracelet/bubbletea"

// wizardDoneMsg signals one wizard step finished. Payload is whatever
// the step produced (an int for list selections, a string for text
// inputs, etc.). Cancelled=true means the user bailed (Ctrl+C / Esc).
type wizardDoneMsg struct {
	Cancelled bool
	Payload   any
}

// finishWith returns a Cmd that emits wizardDoneMsg with payload.
func finishWith(payload any) tea.Cmd {
	return func() tea.Msg { return wizardDoneMsg{Payload: payload} }
}

// finishCancelled returns a Cmd that emits a cancellation
// wizardDoneMsg.
func finishCancelled() tea.Cmd {
	return func() tea.Msg { return wizardDoneMsg{Cancelled: true} }
}

// singleShotRunner hosts one sub-model and translates wizardDoneMsg
// into tea.Quit, capturing the payload + cancellation flag for the
// caller to read after Run() returns.
//
// Used by RunList / runKeyInput / runModelInput / runProjectField so
// they can keep their old "blocking call returns a value" contract
// while the orchestrator uses the same sub-models internally.
type singleShotRunner struct {
	inner     tea.Model
	payload   any
	cancelled bool
	finished  bool
}

func (s *singleShotRunner) Init() tea.Cmd { return s.inner.Init() }

func (s *singleShotRunner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if d, ok := msg.(wizardDoneMsg); ok {
		s.payload = d.Payload
		s.cancelled = d.Cancelled
		s.finished = true
		return s, tea.Quit
	}
	var cmd tea.Cmd
	s.inner, cmd = s.inner.Update(msg)
	return s, cmd
}

func (s *singleShotRunner) View() string { return s.inner.View() }
