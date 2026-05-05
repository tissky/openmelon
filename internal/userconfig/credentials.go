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

// GetAPIKey returns the stored global key for provider, or "".
//
// Use ResolveAPIKey(workdir, provider) when you want the
// project-overrides-global semantics.
func GetAPIKey(provider string) string {
	creds, err := LoadCredentials()
	if err != nil {
		return ""
	}
	return creds.APIKeys[provider]
}

// projectCredentialsPath returns <workdir>/.openmelon/credentials.json.
// We hard-code the dir name to avoid an import cycle with projectx —
// projectx.DirName is the same string.
func projectCredentialsPath(workdir string) string {
	return filepath.Join(workdir, ".openmelon", "credentials.json")
}

// LoadProjectCredentials reads <workdir>/.openmelon/credentials.json.
// Missing file → empty Credentials (not an error). Used by the
// resolver below; callers can also use it directly to inspect.
func LoadProjectCredentials(workdir string) (*Credentials, error) {
	path := projectCredentialsPath(workdir)
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

// SaveProjectCredentials writes the per-project credentials file with
// mode 0600. Atomic via temp + rename like SaveCredentials.
func SaveProjectCredentials(workdir string, c *Credentials) error {
	if c.APIKeys == nil {
		c.APIKeys = map[string]string{}
	}
	path := projectCredentialsPath(workdir)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("userconfig: mkdir %s: %w", dir, err)
	}
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("userconfig: marshal project credentials: %w", err)
	}
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

// SetProjectAPIKey is the project-scoped equivalent of SetAPIKey.
func SetProjectAPIKey(workdir, provider, key string) error {
	creds, err := LoadProjectCredentials(workdir)
	if err != nil {
		return err
	}
	creds.APIKeys[provider] = key
	return SaveProjectCredentials(workdir, creds)
}

// UnsetProjectAPIKey removes a provider's key from the project
// credentials. No-op if the key wasn't set.
func UnsetProjectAPIKey(workdir, provider string) error {
	creds, err := LoadProjectCredentials(workdir)
	if err != nil {
		return err
	}
	if _, ok := creds.APIKeys[provider]; !ok {
		return nil
	}
	delete(creds.APIKeys, provider)
	return SaveProjectCredentials(workdir, creds)
}

// KeySource describes where a resolved API key came from.
type KeySource string

const (
	// SourceProject means the key came from <workdir>/.openmelon/credentials.json.
	SourceProject KeySource = "project"
	// SourceGlobal means it came from ~/.openmelon/credentials.json.
	SourceGlobal KeySource = "global"
	// SourceNone means no key was found at either scope. Caller may
	// still fall back to env vars (llm.New does this internally).
	SourceNone KeySource = "none"
)

// ResolveAPIKey returns the API key for provider, with project-overrides-
// global semantics. Pass the project workdir; an empty workdir skips the
// project lookup. Returns ("", SourceNone) when neither scope has a key.
//
// Env vars are NOT consulted here — llm.New / imagegen.New do that as
// a final fallback when the passed-in key is empty.
func ResolveAPIKey(workdir, provider string) (string, KeySource) {
	if workdir != "" {
		if proj, err := LoadProjectCredentials(workdir); err == nil {
			if k := proj.APIKeys[provider]; k != "" {
				return k, SourceProject
			}
		}
	}
	if global, err := LoadCredentials(); err == nil {
		if k := global.APIKeys[provider]; k != "" {
			return k, SourceGlobal
		}
	}
	return "", SourceNone
}
