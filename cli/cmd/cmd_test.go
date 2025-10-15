package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListDirectoryContents(t *testing.T) {
	// Test in current directory
	entries, err := ListDirectoryContents()
	if err != nil {
		t.Fatalf("ListDirectoryContents() error = %v", err)
	}

	// Should return at least some entries (unless directory is completely empty)
	// We can't guarantee specific entries, but we can test the structure
	for _, entry := range entries {
		if entry.name == "" {
			t.Error("Entry should have a non-empty name")
		}
		// entry.isDir is a boolean, so no specific validation needed beyond type
	}
}

func TestListDirectoryContentsInTempDir(t *testing.T) {
	// Create a temporary directory with known contents
	tempDir, err := os.MkdirTemp("", "test_cmd")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create some test files and directories
	testFile := filepath.Join(tempDir, "test.txt")
	testDir := filepath.Join(tempDir, "testdir")
	gitDir := filepath.Join(tempDir, ".git")

	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if err := os.Mkdir(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Test ListDirectoryContents
	entries, err := ListDirectoryContents()
	if err != nil {
		t.Fatalf("ListDirectoryContents() error = %v", err)
	}

	// Should have exactly 3 entries
	if len(entries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(entries))
	}

	// Check that we have the expected entries
	entryMap := make(map[string]bool)
	for _, entry := range entries {
		entryMap[entry.name] = entry.isDir
	}

	// Verify file
	if isDir, exists := entryMap["test.txt"]; !exists {
		t.Error("Expected test.txt to be in entries")
	} else if isDir {
		t.Error("test.txt should not be marked as directory")
	}

	// Verify directory
	if isDir, exists := entryMap["testdir"]; !exists {
		t.Error("Expected testdir to be in entries")
	} else if !isDir {
		t.Error("testdir should be marked as directory")
	}

	// Verify .git directory
	if isDir, exists := entryMap[".git"]; !exists {
		t.Error("Expected .git to be in entries")
	} else if !isDir {
		t.Error(".git should be marked as directory")
	}
}

func TestEntryIsGitDir(t *testing.T) {
	tests := []struct {
		name     string
		entry    Entry
		expected bool
	}{
		{
			name:     "git directory",
			entry:    Entry{name: ".git", isDir: true},
			expected: true,
		},
		{
			name:     "git file (edge case)",
			entry:    Entry{name: ".git", isDir: false},
			expected: true,
		},
		{
			name:     "regular directory",
			entry:    Entry{name: "src", isDir: true},
			expected: false,
		},
		{
			name:     "regular file",
			entry:    Entry{name: "main.go", isDir: false},
			expected: false,
		},
		{
			name:     "gitignore file",
			entry:    Entry{name: ".gitignore", isDir: false},
			expected: false,
		},
		{
			name:     "git with suffix",
			entry:    Entry{name: ".git-backup", isDir: true},
			expected: false,
		},
		{
			name:     "empty name",
			entry:    Entry{name: "", isDir: false},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.entry.IsGitDir()
			if result != tt.expected {
				t.Errorf("Entry.IsGitDir() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestEntryStruct(t *testing.T) {
	// Test Entry struct creation and field access
	entry := Entry{name: "test.go", isDir: false}

	if entry.name != "test.go" {
		t.Errorf("Expected name 'test.go', got '%s'", entry.name)
	}

	if entry.isDir != false {
		t.Errorf("Expected isDir false, got %v", entry.isDir)
	}

	dirEntry := Entry{name: "src", isDir: true}

	if dirEntry.name != "src" {
		t.Errorf("Expected name 'src', got '%s'", dirEntry.name)
	}

	if dirEntry.isDir != true {
		t.Errorf("Expected isDir true, got %v", dirEntry.isDir)
	}
}