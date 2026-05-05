package onboard

// projectkey.go — interactive wizard for `openmelon project set-key`.
//
// Same look and feel as the auth wizard but scoped to one project: pick
// provider, paste key (masked), save into <workdir>/.openmelon/credentials.json.

import (
	"fmt"
	"os"

	"github.com/eight-acres-lab/openmelon/internal/userconfig"
)

// RunProjectKeyWizard runs the per-project key wizard. If providerHint
// is non-empty, the provider picker is skipped and the wizard goes
// straight to the key input for that provider.
//
// Returns (provider, ok, err). ok=false means the user cancelled — the
// caller should treat that as a clean no-op.
func RunProjectKeyWizard(workdir, providerHint string) (provider string, ok bool, err error) {
	chosen := -1
	if providerHint != "" {
		for i, p := range providerOptions {
			if p.slug == providerHint {
				chosen = i
				break
			}
		}
		if chosen < 0 {
			return "", false, fmt.Errorf("unknown provider %q (supported: openrouter, openai, anthropic)", providerHint)
		}
	} else {
		header := headerStyle.Render(fmt.Sprintf("Set a project-scoped key for %s", workdir)) + "\n\n" +
			bodyStyle.Render("Pick the provider this key is for. The key will be stored at\n<project>/.openmelon/credentials.json (mode 0600) and override the global key when you run openmelon in this project.")
		items := make([]ListItem, len(providerOptions))
		for i, p := range providerOptions {
			items[i] = ListItem{Title: p.title, Subtitle: p.subtitle}
		}
		res, err := RunList(listOpts{
			Title: header,
			Items: items,
			Help:  "↑/↓ to choose · 1/2/3 shortcut · enter to continue · ctrl+c to cancel",
		})
		if err != nil {
			return "", false, err
		}
		if res.Cancelled {
			return "", false, nil
		}
		chosen = res.Index
	}

	p := providerOptions[chosen]
	apiKey, err := runKeyInput(p)
	if err != nil {
		return "", false, err
	}
	if apiKey == "" {
		return "", false, nil
	}
	if err := userconfig.SetProjectAPIKey(workdir, p.slug, apiKey); err != nil {
		return "", false, fmt.Errorf("save: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Project key saved for %s.\n", p.slug)
	return p.slug, true, nil
}
