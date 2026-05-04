package artifacts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStableID_deterministic(t *testing.T) {
	id1 := StableID("proj", "workflow", "stage")
	id2 := StableID("proj", "workflow", "stage")
	if id1 != id2 {
		t.Errorf("StableID is not deterministic: %q != %q", id1, id2)
	}
}

func TestStableID_differentParts(t *testing.T) {
	a := StableID("a", "b", "c")
	b := StableID("a", "b", "d")
	if a == b {
		t.Error("StableID should differ for different parts")
	}
}

func TestStableID_length(t *testing.T) {
	id := StableID("x")
	if len(id) != 16 {
		t.Errorf("StableID length = %d, want 16", len(id))
	}
}

func TestWrite_createsFiles(t *testing.T) {
	dir := t.TempDir()
	a := &Artifact{
		ID:         "abcd1234",
		Type:       TypeImagePrompt,
		Content:    "some prompt content",
		Provenance: `{"artifact_id":"abcd1234"}`,
	}

	if err := Write(dir, a); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	contentPath := filepath.Join(dir, "abcd1234.image.prompt.txt")
	data, err := os.ReadFile(contentPath)
	if err != nil {
		t.Fatalf("content file not found: %v", err)
	}
	if string(data) != "some prompt content" {
		t.Errorf("content = %q, want %q", string(data), "some prompt content")
	}

	provPath := filepath.Join(dir, "abcd1234.provenance.json")
	pData, err := os.ReadFile(provPath)
	if err != nil {
		t.Fatalf("provenance file not found: %v", err)
	}
	if !strings.Contains(string(pData), "abcd1234") {
		t.Errorf("provenance file missing artifact_id, got: %q", string(pData))
	}
}

func TestWrite_noProvenanceField(t *testing.T) {
	dir := t.TempDir()
	a := &Artifact{
		ID:      "ef567890",
		Type:    TypeCopyDraft,
		Content: "copy draft content",
	}

	if err := Write(dir, a); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	provPath := filepath.Join(dir, "ef567890.provenance.json")
	if _, err := os.Stat(provPath); !os.IsNotExist(err) {
		t.Error("provenance file should NOT be written when Provenance is empty")
	}
}
