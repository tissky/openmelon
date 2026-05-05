package userconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCredentialsMissingReturnsEmpty(t *testing.T) {
	withTmpHome(t)
	c, err := LoadCredentials()
	if err != nil {
		t.Fatalf("LoadCredentials: %v", err)
	}
	if len(c.APIKeys) != 0 {
		t.Errorf("expected empty map, got %v", c.APIKeys)
	}
}

func TestSaveCredentialsRoundtrip(t *testing.T) {
	withTmpHome(t)
	in := &Credentials{APIKeys: map[string]string{
		"openrouter": "sk-or-test",
		"openai":     "sk-test",
	}}
	if err := SaveCredentials(in); err != nil {
		t.Fatalf("SaveCredentials: %v", err)
	}
	out, err := LoadCredentials()
	if err != nil {
		t.Fatalf("LoadCredentials: %v", err)
	}
	if out.APIKeys["openrouter"] != "sk-or-test" || out.APIKeys["openai"] != "sk-test" {
		t.Errorf("roundtrip mismatch: %+v", out.APIKeys)
	}
}

func TestSaveCredentialsWritesMode0600(t *testing.T) {
	home := withTmpHome(t)
	if err := SaveCredentials(&Credentials{APIKeys: map[string]string{"openrouter": "k"}}); err != nil {
		t.Fatalf("SaveCredentials: %v", err)
	}
	st, err := os.Stat(filepath.Join(home, "credentials.json"))
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	mode := st.Mode().Perm()
	if mode != 0o600 {
		t.Errorf("expected mode 0600, got %o", mode)
	}
}

func TestSetAndGetAPIKey(t *testing.T) {
	withTmpHome(t)
	if err := SetAPIKey("openrouter", "sk-or-1"); err != nil {
		t.Fatalf("SetAPIKey: %v", err)
	}
	if got := GetAPIKey("openrouter"); got != "sk-or-1" {
		t.Errorf("GetAPIKey: got %q want sk-or-1", got)
	}
	if got := GetAPIKey("ghost"); got != "" {
		t.Errorf("missing key should return empty, got %q", got)
	}
}

func TestProjectCredentialsRoundtripWithMode0600(t *testing.T) {
	withTmpHome(t)
	wd := t.TempDir()
	if err := SetProjectAPIKey(wd, "openai", "sk-proj"); err != nil {
		t.Fatalf("SetProjectAPIKey: %v", err)
	}
	st, err := os.Stat(filepath.Join(wd, ".openmelon", "credentials.json"))
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if mode := st.Mode().Perm(); mode != 0o600 {
		t.Errorf("expected mode 0600, got %o", mode)
	}
	out, err := LoadProjectCredentials(wd)
	if err != nil {
		t.Fatalf("LoadProjectCredentials: %v", err)
	}
	if out.APIKeys["openai"] != "sk-proj" {
		t.Errorf("project key roundtrip: %v", out.APIKeys)
	}
}

func TestUnsetProjectAPIKey(t *testing.T) {
	withTmpHome(t)
	wd := t.TempDir()
	_ = SetProjectAPIKey(wd, "openai", "sk-proj")
	_ = SetProjectAPIKey(wd, "openrouter", "sk-or")
	if err := UnsetProjectAPIKey(wd, "openai"); err != nil {
		t.Fatalf("UnsetProjectAPIKey: %v", err)
	}
	out, _ := LoadProjectCredentials(wd)
	if _, ok := out.APIKeys["openai"]; ok {
		t.Errorf("openai not unset: %v", out.APIKeys)
	}
	if out.APIKeys["openrouter"] != "sk-or" {
		t.Errorf("openrouter clobbered: %v", out.APIKeys)
	}
	// Unsetting an absent key is a no-op (no error).
	if err := UnsetProjectAPIKey(wd, "ghost"); err != nil {
		t.Errorf("unset of absent key should be no-op, got %v", err)
	}
}

func TestResolveAPIKey(t *testing.T) {
	withTmpHome(t)
	wd := t.TempDir()

	// Both scopes empty.
	k, src := ResolveAPIKey(wd, "openrouter")
	if k != "" || src != SourceNone {
		t.Errorf("empty: got (%q, %s), want ('', none)", k, src)
	}

	// Global only.
	_ = SetAPIKey("openrouter", "sk-global")
	k, src = ResolveAPIKey(wd, "openrouter")
	if k != "sk-global" || src != SourceGlobal {
		t.Errorf("global only: got (%q, %s)", k, src)
	}

	// Both — project wins.
	_ = SetProjectAPIKey(wd, "openrouter", "sk-proj")
	k, src = ResolveAPIKey(wd, "openrouter")
	if k != "sk-proj" || src != SourceProject {
		t.Errorf("both: got (%q, %s), want (sk-proj, project)", k, src)
	}

	// Different provider — falls back to global only.
	_ = SetAPIKey("openai", "sk-openai-global")
	k, src = ResolveAPIKey(wd, "openai")
	if k != "sk-openai-global" || src != SourceGlobal {
		t.Errorf("global fallback: got (%q, %s)", k, src)
	}

	// Empty workdir → only checks global.
	k, src = ResolveAPIKey("", "openrouter")
	if k != "sk-global" || src != SourceGlobal {
		t.Errorf("no workdir: got (%q, %s)", k, src)
	}
}

func TestIsTrustedExactAndSubdir(t *testing.T) {
	c := &Config{TrustedDirs: []string{"/work/ai-talks"}}
	if !c.IsTrusted("/work/ai-talks") {
		t.Error("exact match should be trusted")
	}
	if !c.IsTrusted("/work/ai-talks/subdir") {
		t.Error("subdir should be trusted")
	}
	if !c.IsTrusted("/work/ai-talks/deeply/nested") {
		t.Error("deep subdir should be trusted")
	}
	if c.IsTrusted("/work/ai-talks-other") {
		t.Error("prefix-only should NOT be trusted")
	}
	if c.IsTrusted("/elsewhere") {
		t.Error("unrelated path should NOT be trusted")
	}
}

func TestAddTrustedIsIdempotent(t *testing.T) {
	c := &Config{}
	if !c.AddTrusted("/work/ai-talks") {
		t.Error("first AddTrusted should return true")
	}
	if c.AddTrusted("/work/ai-talks") {
		t.Error("second AddTrusted should return false")
	}
	if len(c.TrustedDirs) != 1 {
		t.Errorf("expected 1 entry, got %d", len(c.TrustedDirs))
	}
}
