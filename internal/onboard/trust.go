package onboard

// trust.go — first wizard. Prompts the user to confirm trust on the
// current working directory before any agent loop runs.

import (
	"fmt"
	"os"
	"strings"

	"github.com/eight-acres-lab/openmelon/internal/userconfig"
)

// EnsureTrust runs the trust prompt if cwd is not in the trusted list.
// Returns:
//   trusted=true  → proceed
//   trusted=false → user declined; caller should exit cleanly
//
// Adds cwd to the trusted_dirs list when the user confirms.
func EnsureTrust(cwd string) (trusted bool, err error) {
	cfg, err := userconfig.LoadConfig()
	if err != nil {
		return false, err
	}
	if cfg.IsTrusted(cwd) {
		return true, nil
	}

	header := headerStyle.Render(fmt.Sprintf("> You are in %s", pathStyle.Render(cwd)))
	body := bodyStyle.Render(strings.Join([]string{
		"",
		"Do you trust the contents of this directory?",
		"openmelon will read project files, registered characters and references,",
		"and may invoke tools the agent decides to call. Only continue if you trust",
		"the contents here.",
	}, "\n"))

	res, err := RunList(listOpts{
		Title: header + "\n" + body,
		Items: []ListItem{
			{Title: "Yes, continue"},
			{Title: "No, quit"},
		},
		Help: "↑/↓ to choose · 1/2 shortcut · enter to continue · ctrl+c to quit",
	})
	if err != nil {
		return false, err
	}
	if res.Cancelled || res.Index != 0 {
		fmt.Fprintln(os.Stderr, "openmelon: not trusted; exiting")
		return false, nil
	}

	cfg.AddTrusted(cwd)
	if err := userconfig.SaveConfig(cfg); err != nil {
		return false, fmt.Errorf("trust: save config: %w", err)
	}
	return true, nil
}
