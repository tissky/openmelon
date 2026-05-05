package main

// cmd_setup.go — `openmelon setup` re-runs the auth wizard so users can
// rotate keys or switch providers without editing files.

import (
	"fmt"
	"os"

	"github.com/eight-acres-lab/openmelon/internal/onboard"
	"github.com/eight-acres-lab/openmelon/internal/userconfig"
)

func runSetup(_ []string) error {
	// Force the wizard to run by clearing existing credentials in
	// memory before calling EnsureAuth — but persist nothing until the
	// user finishes the flow successfully.
	creds, err := userconfig.LoadCredentials()
	if err != nil {
		return err
	}
	had := len(creds.APIKeys) > 0
	if had {
		// Stash and clear so EnsureAuth's "already configured?" check
		// fires the wizard. We restore on user cancel.
		empty := &userconfig.Credentials{APIKeys: map[string]string{}}
		if err := userconfig.SaveCredentials(empty); err != nil {
			return err
		}
	}
	configured, err := onboard.EnsureAuth()
	if err != nil {
		// Restore on error.
		if had {
			_ = userconfig.SaveCredentials(creds)
		}
		return err
	}
	if !configured {
		// User cancelled — restore.
		if had {
			_ = userconfig.SaveCredentials(creds)
			fmt.Fprintln(os.Stderr, "setup cancelled; previous credentials kept")
		}
		return nil
	}
	fmt.Fprintln(os.Stderr, "setup complete.")
	return nil
}
