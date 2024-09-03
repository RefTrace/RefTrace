package parser

import (
	"path/filepath"
	"testing"
)

func TestParseSpeed(t *testing.T) {
	dir := filepath.Join(getTestDataDir(), "nf-core")
	lines, err := processDirectory(dir)
	if err != nil {
		t.Fatalf("Error processing directory: %v", err)
	}
	t.Logf("Total file count: %d", lines)
}
