package project

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		json        string
		wantID      string
		wantErrText string
		wantNotExist bool
	}{
		{
			name: "valid project",
			json: `{
				"id": "test-proj",
				"name": "Test Project",
				"platform": "xiaohongshu",
				"audience": "food lovers",
				"persona": "casual reviewer"
			}`,
			wantID: "test-proj",
		},
		{
			name:        "missing required field id",
			json:        `{"name": "No ID", "platform": "xiaohongshu"}`,
			wantErrText: "missing required field: id",
		},
		{
			name:        "missing required field name",
			json:        `{"id": "proj-1", "platform": "xiaohongshu"}`,
			wantErrText: "missing required field: name",
		},
		{
			name:        "missing required field platform",
			json:        `{"id": "proj-1", "name": "Proj"}`,
			wantErrText: "missing required field: platform",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "project.json")
			if err := os.WriteFile(path, []byte(tc.json), 0o644); err != nil {
				t.Fatal(err)
			}

			p, err := Load(path)
			if tc.wantErrText != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.wantErrText)
				}
				if !strings.Contains(err.Error(), tc.wantErrText) {
					t.Fatalf("expected error %q to contain %q", err.Error(), tc.wantErrText)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if p.ID != tc.wantID {
				t.Errorf("ID = %q, want %q", p.ID, tc.wantID)
			}
		})
	}
}

func TestLoad_fileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/project.json")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected os.ErrNotExist, got: %v", err)
	}
}
