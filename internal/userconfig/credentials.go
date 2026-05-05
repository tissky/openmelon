package userconfig

// credentials.go — ~/.openmelon/credentials.json (0600 perms).
//
// Lives separate from config.json so we can write it with restrictive
// permissions and exclude it from any sync / backup the user might be
// doing on their config dir. Same pattern as Codex CLI's auth.json and
// Claude Code's credentials store.
//
// Schema is intentionally simple: a string→string map of provider name
// (anthropic / openai / openrouter) to API key. Future OAuth tokens
// will get their own typed shape; for now we just store API keys.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Credentials is the on-disk shape of ~/.openmelon/credentials.json.
type Credentials struct {
	// APIKeys maps provider slug ("anthropic" / "openai" / "openrouter")
	// to the user's API key for that provider. Missing entry = no key
	// stored for that provider; the runtime falls back to the matching
	// env var.
	APIKeys map[string]string `json:"api_keys,omitempty"`
}

// LoadCredentials reads ~/.openmelon/credentials.json. Missing file →
// empty Credentials (not an error — first run is fine).
func LoadCredentials() (*Credentials, error) {
	home, err := Home()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(home, "credentials.json")
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Credentials{APIKeys: map[string]string{}}, nil
		}
		return nil, fmt.Errorf("userconfig: read %s: %w", path, err)
	}
	var c Credentials
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("userconfig: parse %s: %w", path, err)
	}
	if c.APIKeys == nil {
		c.APIKeys = map[string]string{}
	}
	return &c, nil
}

// SaveCredentials writes ~/.openmelon/credentials.json with mode 0600.
func SaveCredentials(c *Credentials) error {
	home, err := EnsureHome()
	if err != nil {
		return err
	}
	if c.APIKeys == nil {
		c.APIKeys = map[string]string{}
	}
	path := filepath.Join(home, "credentials.json")
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("userconfig: marshal credentials: %w", err)
	}
	// Atomic write with restrictive perms.
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".cred-")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if err := os.Chmod(tmpPath, 0o600); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(append(b, '\n')); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}

// SetAPIKey is a convenience wrapper for the common "set one key" path.
func SetAPIKey(provider, key string) error {
	creds, err := LoadCredentials()
	if err != nil {
		return err
	}
	creds.APIKeys[provider] = key
	return SaveCredentials(creds)
}

// GetAPIKey returns the stored key for provider, or "".
func GetAPIKey(provider string) string {
	creds, err := LoadCredentials()
	if err != nil {
		return ""
	}
	return creds.APIKeys[provider]
}
