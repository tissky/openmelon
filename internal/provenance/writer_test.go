package provenance

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestAppendRecord_createAndAppend(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "provenance.jsonl")

	rec1 := &Record{
		ArtifactID: "id1",
		ProjectID:  "proj",
		Stage:      "stage1",
		Timestamp:  "2024-01-01T00:00:00Z",
	}
	rec2 := &Record{
		ArtifactID: "id2",
		ProjectID:  "proj",
		Stage:      "stage2",
		Timestamp:  "2024-01-01T00:01:00Z",
	}

	if err := AppendRecord(path, rec1); err != nil {
		t.Fatalf("AppendRecord(rec1) error: %v", err)
	}
	if err := AppendRecord(path, rec2); err != nil {
		t.Fatalf("AppendRecord(rec2) error: %v", err)
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var records []Record
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var r Record
		if err := json.Unmarshal([]byte(line), &r); err != nil {
			t.Fatalf("failed to unmarshal line %q: %v", line, err)
		}
		records = append(records, r)
	}

	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if records[0].ArtifactID != "id1" {
		t.Errorf("records[0].ArtifactID = %q, want %q", records[0].ArtifactID, "id1")
	}
	if records[1].ArtifactID != "id2" {
		t.Errorf("records[1].ArtifactID = %q, want %q", records[1].ArtifactID, "id2")
	}
}

func TestAppendRecord_createsDirIfMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "deep", "provenance.jsonl")
	rec := &Record{ArtifactID: "x"}

	if err := AppendRecord(path, rec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}
