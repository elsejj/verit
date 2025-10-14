package utils

import (
	"os"
	"path"
	"regexp"
	"testing"
)

func TestSed(t *testing.T) {
	tmpDir := t.TempDir()

	filePath := path.Join(tmpDir, "test.txt")

	originalContent := "Hello, World!\nThis is a test file.\nVersion: 1.0.0\n"
	err := os.WriteFile(filePath, []byte(originalContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	re := regexp.MustCompile(`Version:\s*(\d\.\d\.\d)`)
	newVersion := "2.0.0"

	err = Sed(filePath, re, newVersion)
	if err != nil {
		t.Fatalf("Sed failed: %v", err)
	}

	updatedContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read updated file: %v", err)
	}

	expectedContent := "Hello, World!\nThis is a test file.\nVersion: 2.0.0\n"
	if string(updatedContent) != expectedContent {
		t.Errorf("Sed did not update content as expected.\nGot:\n%s\nExpected:\n%s", string(updatedContent), expectedContent)
	}

}
