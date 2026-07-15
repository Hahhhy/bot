package image

import (
	"os"
	"testing"
)

func TestGenerateTom(t *testing.T) {
	path, err := GenerateTom("D O")
	if err != nil {
		t.Fatalf("GenerateTom failed: %v", err)
	}
	if path == "" {
		t.Fatal("expected path, got empty string")
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not found: %s", path)
	}
}
