package onboard

// initproject.go — wizard that creates a new project in cwd when no
// project is found there.

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/eight-acres-lab/openmelon/internal/projectx"
	"github.com/eight-acres-lab/openmelon/internal/userconfig"
)

// EnsureProject runs the project-init wizard if cwd is not already in
// (or under) an openmelon project. Returns the project's workdir.
func EnsureProject(cwd string) (workdir string, err error) {
	wd, err := projectx.Discover(cwd)
	if err != nil {
		return "", err
	}
	if wd != "" {
		return wd, nil
	}

	// Step 1: confirm.
	res, err := RunList(listOpts{
		Title: headerStyle.Render(fmt.Sprintf("> No openmelon project found in %s", cwd)) + "\n\n" +
			bodyStyle.Render("Create one here? It just adds a `.openmelon/` directory with a project.json,\nplus subdirs for characters, references, materials, and sessions."),
		Items: []ListItem{
			{Title: "Yes, initialize a new project here"},
			{Title: "No, quit"},
		},
		Help: "enter to continue · esc to cancel",
	})
	if err != nil {
		return "", err
	}
	if res.Cancelled || res.Index != 0 {
		return "", nil
	}

	// Step 2: project id (slug).
	defaultID := slugify(filepath.Base(cwd))
	id, err := runProjectField("Project id (kebab-case, this is the registry key)", defaultID)
	if err != nil {
		return "", err
	}
	if id == "" {
		id = defaultID
	}
	if err := projectx.ValidateID(id); err != nil {
		return "", fmt.Errorf("init: %w", err)
	}

	// Step 3: name.
	name, err := runProjectField("Project name (free text shown in the UI)", id)
	if err != nil {
		return "", err
	}
	if name == "" {
		name = id
	}

	// Step 4: description (optional).
	description, err := runProjectField("One-line description (optional)", "")
	if err != nil {
		return "", err
	}

	p, err := projectx.Init(cwd, id, name)
	if err != nil {
		return "", err
	}
	if description != "" {
		p.Description = description
		if err := projectx.Save(cwd, p); err != nil {
			return "", err
		}
	}
	if err := userconfig.Register(id, name, cwd); err != nil {
		return "", fmt.Errorf("init: register: %w", err)
	}
	if err := userconfig.SetCurrent(id); err != nil {
		return "", fmt.Errorf("init: set current: %w", err)
	}
	return cwd, nil
}

// --- text-input helper ---

type projectFieldModel struct {
	prompt string
	input  textinput.Model
	done   bool
	cancel bool
}

func (m *projectFieldModel) Init() tea.Cmd { return textinput.Blink }

func (m *projectFieldModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(k, key.NewBinding(key.WithKeys("ctrl+c", "esc"))):
			m.cancel = true
			return m, finishCancelled()
		case key.Matches(k, key.NewBinding(key.WithKeys("enter"))):
			m.done = true
			return m, finishWith(strings.TrimSpace(m.input.Value()))
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *projectFieldModel) View() string {
	var b strings.Builder
	b.WriteString(headerStyle.Render(m.prompt))
	b.WriteString("\n\n")
	b.WriteString(m.input.View())
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("enter to continue · esc to cancel"))
	b.WriteString("\n")
	return b.String()
}

func runProjectField(prompt, def string) (string, error) {
	ti := textinput.New()
	ti.Placeholder = def
	if def != "" {
		ti.SetValue(def)
	}
	ti.CharLimit = 200
	ti.Width = 60
	ti.Focus()
	m := &projectFieldModel{prompt: prompt, input: ti}
	runner := &singleShotRunner{inner: m}
	if _, err := tea.NewProgram(runner, tea.WithAltScreen()).Run(); err != nil {
		return "", err
	}
	if runner.cancelled {
		return "", fmt.Errorf("cancelled")
	}
	if s, ok := runner.payload.(string); ok {
		return s, nil
	}
	return strings.TrimSpace(m.input.Value()), nil
}

// slugify reduces a directory basename to a kebab-case project id.
// Mirrors cmd/openmelon/cmd_init.go's slugFromBase but lives here so
// we don't reach into cmd/.
func slugify(base string) string {
	base = strings.ToLower(base)
	var b strings.Builder
	prevHy := false
	for _, r := range base {
		switch {
		case r >= 'a' && r <= 'z' || r >= '0' && r <= '9':
			b.WriteRune(r)
			prevHy = false
		case r == ' ' || r == '_' || r == '-' || r == '.':
			if !prevHy && b.Len() > 0 {
				b.WriteByte('-')
				prevHy = true
			}
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" || (out[0] < 'a' || out[0] > 'z') {
		out = "project-" + out
		out = strings.TrimRight(out, "-")
	}
	if len(out) < 2 {
		out = "project"
	}
	return out
}
