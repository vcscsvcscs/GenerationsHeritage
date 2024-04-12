package utilities

import (
	"os"
	"testing"
)

func TestFileExists(t *testing.T) {
	// Create a temporary file for testing
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	// Test with an existing file
	exists := FileExists(file.Name())
	if !exists {
		t.Errorf("FileExists(%q) = false, want true", file.Name())
	}

	// Test with a non-existing file
	exists = FileExists("non_existing_file.txt")
	if exists {
		t.Errorf("FileExists(%q) = true, want false", "non_existing_file.txt")
	}
}
