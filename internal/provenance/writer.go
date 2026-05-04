package provenance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// AppendRecord serializes rec to JSON and appends it as a single line to the file at path.
// The directory is created automatically if it does not exist.
// The file is fsynced after writing to guard against data loss on crash.
func AppendRecord(path string, rec *Record) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("provenance.AppendRecord: mkdir: %w", err)
	}

	data, err := json.Marshal(rec)
	if err != nil {
		return fmt.Errorf("provenance.AppendRecord: marshal: %w", err)
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return fmt.Errorf("provenance.AppendRecord: open %q: %w", path, err)
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("provenance.AppendRecord: write: %w", err)
	}
	if _, err := f.WriteString("\n"); err != nil {
		return fmt.Errorf("provenance.AppendRecord: write newline: %w", err)
	}
	if err := f.Sync(); err != nil {
		return fmt.Errorf("provenance.AppendRecord: sync: %w", err)
	}
	return nil
}
