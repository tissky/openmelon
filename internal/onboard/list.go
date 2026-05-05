// Package onboard runs first-launch and per-launch wizards: trust
// confirmation, API-key setup, and project initialization. Each wizard
// is its own short-lived bubbletea Program; onboard.Ensure stitches
// them together.
//
// We deliberately don't reuse the main TUI's Model — these screens are
// modal, single-purpose, and easier to read as small focused programs.
package onboard

// list.go — reusable arrow-key list selector. Used by the trust prompt
// and the auth-provider picker. Modeled on the screens in Codex's
// onboarding (numbered list, ">" arrow on the active row, optional
// dim subtitle under each option).

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ListItem is one selectable row.
type ListItem struct {
	Title    string // headline; rendered in the accent color when selected
	Subtitle string // optional dim subtitle under the title
}

// ListResult is what RunList returns. Index is -1 if the user
// cancelled (Ctrl+C / Esc / "q").
type ListResult struct {
	Index     int
	Cancelled bool
}

// listOpts configures a list selector.
type listOpts struct {
	Title    string // shown as the header above the list
	Help     string // shown below the list (e.g. "Press enter to continue")
	Items    []ListItem
	Initial  int // initial selection
}

// listModel is the bubbletea Model behind RunList.
type listModel struct {
	opts     listOpts
	cursor   int
	chosen   int
	cancel   bool
	finished bool
}

func (m *listModel) Init() tea.Cmd { return nil }

func (m *listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(k, key.NewBinding(key.WithKeys("ctrl+c", "esc", "q"))):
			m.cancel = true
			m.finished = true
			return m, finishCancelled()
		case key.Matches(k, key.NewBinding(key.WithKeys("up", "k"))):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(k, key.NewBinding(key.WithKeys("down", "j"))):
			if m.cursor < len(m.opts.Items)-1 {
				m.cursor++
			}
		case key.Matches(k, key.NewBinding(key.WithKeys("enter"))):
			m.chosen = m.cursor
			m.finished = true
			return m, finishWith(m.cursor)
		}
		// Number-key shortcut: 1..9 picks that item.
		if len(k.String()) == 1 && k.String()[0] >= '1' && k.String()[0] <= '9' {
			n := int(k.String()[0]-'1') // 0-indexed
			if n < len(m.opts.Items) {
				m.cursor = n
				m.chosen = n
				m.finished = true
				return m, finishWith(n)
			}
		}
	}
	return m, nil
}

func (m *listModel) View() string {
	var b strings.Builder
	if m.opts.Title != "" {
		b.WriteString(m.opts.Title)
		b.WriteString("\n\n")
	}
	for i, it := range m.opts.Items {
		arrow := "  "
		if i == m.cursor {
			arrow = arrowStyle.Render("> ")
		}
		num := fmt.Sprintf("%d.", i+1)
		title := it.Title
		if i == m.cursor {
			num = activeNumStyle.Render(num)
			title = activeTitleStyle.Render(title)
		} else {
			num = numStyle.Render(num)
			title = titleStyle.Render(title)
		}
		fmt.Fprintf(&b, "%s%s %s\n", arrow, num, title)
		if it.Subtitle != "" {
			fmt.Fprintf(&b, "     %s\n", subtitleStyle.Render(it.Subtitle))
		}
		if i < len(m.opts.Items)-1 {
			b.WriteString("\n")
		}
	}
	if m.opts.Help != "" {
		b.WriteString("\n")
		b.WriteString(helpStyle.Render(m.opts.Help))
		b.WriteString("\n")
	}
	return b.String()
}

// RunList runs a list selector until the user picks or cancels.
//
// Standalone use: wraps the model in singleShotRunner so wizardDoneMsg
// → tea.Quit. Used by `openmelon project set-key`'s provider picker.
// The orchestrator (Run) does NOT use this — it hosts listModel
// directly so it can transition to the next state instead of quitting.
func RunList(opts listOpts) (ListResult, error) {
	if opts.Initial < 0 || opts.Initial >= len(opts.Items) {
		opts.Initial = 0
	}
	m := &listModel{opts: opts, cursor: opts.Initial, chosen: -1}
	runner := &singleShotRunner{inner: m}
	if _, err := tea.NewProgram(runner, tea.WithAltScreen()).Run(); err != nil {
		return ListResult{}, err
	}
	if runner.cancelled {
		return ListResult{Index: -1, Cancelled: true}, nil
	}
	if idx, ok := runner.payload.(int); ok {
		return ListResult{Index: idx}, nil
	}
	return ListResult{Index: m.chosen}, nil
}

// --- styles (shared with the rest of onboard) ---

var (
	accentColor   = lipgloss.Color("4")
	mutedColor    = lipgloss.Color("8")
	activeColor   = lipgloss.Color("6") // cyan, like Codex's selection

	arrowStyle       = lipgloss.NewStyle().Foreground(accentColor).Bold(true)
	activeNumStyle   = lipgloss.NewStyle().Foreground(activeColor).Bold(true)
	activeTitleStyle = lipgloss.NewStyle().Foreground(activeColor).Bold(true)
	numStyle         = lipgloss.NewStyle()
	titleStyle       = lipgloss.NewStyle()
	subtitleStyle    = lipgloss.NewStyle().Foreground(mutedColor)
	helpStyle        = lipgloss.NewStyle().Foreground(mutedColor)
	headerStyle      = lipgloss.NewStyle().Bold(true)
	pathStyle        = lipgloss.NewStyle().Foreground(activeColor)
	bodyStyle        = lipgloss.NewStyle()
)
