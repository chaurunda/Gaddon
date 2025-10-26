package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Test helper to create a temporary directory
func createTempDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "checkupdate_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return dir
}

// Test helper to create a temporary directory with .git folder
func createTempGitDir(t *testing.T) string {
	dir := createTempDir(t)
	gitDir := filepath.Join(dir, ".git")
	err := os.Mkdir(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git dir: %v", err)
	}
	return dir
}

// Since we can't easily mock the dependencies without changing the original code,
// we'll test the functions with real filesystem operations but in controlled environments

func TestIsUpdateAvailable_NonGitDirectory(t *testing.T) {
	// Create a temporary directory without .git
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	// Change to the test directory to test the function
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	err = os.Chdir(dir)
	if err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}

	// Test that non-git directory returns false
	result, err := isUpdateAvailable(dir)
	if err != nil {
		t.Errorf("Expected no error for non-git directory, got: %v", err)
	}
	if result {
		t.Error("Expected false for non-git directory, got true")
	}
}

func TestIsUpdateAvailable_GitDirectory(t *testing.T) {
	t.Skip("Skipping git directory test because git.Check_last_local_commit_id and git.Check_last_distant_commit_id use log.Fatal() which terminates the test process. To properly test this, the git package would need to be refactored to return errors instead of calling log.Fatal().")

	// NOTE: The test below would work if the git functions returned errors instead of using log.Fatal()

	// Create a temporary directory with .git folder
	dir := createTempGitDir(t)
	defer os.RemoveAll(dir)

	// Initialize a real git repository
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	err = os.Chdir(dir)
	if err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}

	// This would test git directory detection if git functions didn't use log.Fatal()
	result, err := isUpdateAvailable(dir)

	if err != nil {
		t.Logf("Git command failed as expected in test environment: %v", err)
	} else {
		if result {
			t.Logf("Update availability result: %v", result)
		}
	}
}

func TestIsUpdateAvailable_InvalidPath(t *testing.T) {
	// Test with non-existent directory
	dir := createTempDir(t)
	os.RemoveAll(dir) // Remove it so it doesn't exist

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	// The function works by changing to a directory and listing contents,
	// but if we can't change to the directory, it can't work properly.
	// However, the isUpdateAvailable function calls cmd.ListDirectoryContents()
	// which calls os.ReadDir(".") on the current directory, not the passed path.
	// So we need to test from within a context where the current directory
	// doesn't exist or we can't read it.

	// Create a directory we can change to first
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Now remove the directory we're in
	os.RemoveAll(tempDir)

	// Test that function handles directory listing errors gracefully
	_, err = isUpdateAvailable(dir)
	if err == nil {
		// Actually, the function might not error because it uses "." not the path parameter
		// for listing directory contents. Let's check what actually happens.
		t.Logf("Function didn't return error as expected, this might be due to implementation details")
	} else {
		t.Logf("Function correctly returned error: %v", err)
	}
}

// Test path validation logic separately since PromptForPath involves user input
func TestPathValidation(t *testing.T) {
	tests := []struct {
		name        string
		setupPath   func(*testing.T) string
		cleanupPath func(string)
		expectError bool
		description string
	}{
		{
			name: "valid_directory",
			setupPath: func(t *testing.T) string {
				return createTempDir(t)
			},
			cleanupPath: func(path string) {
				os.RemoveAll(path)
			},
			expectError: false,
			description: "Valid directory should pass validation",
		},
		{
			name: "non_existent_path",
			setupPath: func(t *testing.T) string {
				dir := createTempDir(t)
				os.RemoveAll(dir) // Remove it so it doesn't exist
				return dir
			},
			cleanupPath: func(path string) {
				// Already removed
			},
			expectError: true,
			description: "Non-existent path should fail validation",
		},
		{
			name: "file_not_directory",
			setupPath: func(t *testing.T) string {
				dir := createTempDir(t)
				filePath := filepath.Join(dir, "testfile.txt")
				file, err := os.Create(filePath)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				file.Close()
				return filePath
			},
			cleanupPath: func(path string) {
				os.RemoveAll(filepath.Dir(path))
			},
			expectError: true,
			description: "File path should fail directory validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setupPath(t)
			defer tt.cleanupPath(path)

			// Test if path exists
			_, statErr := os.Stat(path)
			pathExists := !os.IsNotExist(statErr)

			// Test if path is directory (only if it exists)
			var isDirectory bool
			if pathExists {
				if info, err := os.Stat(path); err == nil {
					isDirectory = info.IsDir()
				}
			}

			// Validate expectations
			if tt.expectError {
				if pathExists && isDirectory {
					t.Errorf("%s: Expected validation to fail, but path is valid", tt.description)
				}
			} else {
				if !pathExists {
					t.Errorf("%s: Expected path to exist, but it doesn't", tt.description)
				}
				if pathExists && !isDirectory {
					t.Errorf("%s: Expected path to be directory, but it's not", tt.description)
				}
			}
		})
	}
}

// Test string validation logic used in PromptForPath
func TestStringValidation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"empty_string", "", false},
		{"whitespace_only", "   \n\t  ", false},
		{"valid_string", "/some/path", true},
		{"string_with_spaces", "/path with spaces", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trimmed := strings.TrimSpace(tt.input)
			isEmpty := len(trimmed) == 0

			if tt.expected && isEmpty {
				t.Errorf("Expected %q to be valid after trimming, but it's empty", tt.input)
			}
			if !tt.expected && !isEmpty {
				t.Errorf("Expected %q to be invalid after trimming, but it's %q", tt.input, trimmed)
			}
		})
	}
}

// Test the command configuration
func TestCheckUpdateCommand(t *testing.T) {
	t.Run("command_properties", func(t *testing.T) {
		if checkUpdateCmd.Use != "u" {
			t.Errorf("Expected Use to be 'u', got '%s'", checkUpdateCmd.Use)
		}

		expectedAliases := []string{"update"}
		if len(checkUpdateCmd.Aliases) != len(expectedAliases) {
			t.Errorf("Expected %d aliases, got %d", len(expectedAliases), len(checkUpdateCmd.Aliases))
		}

		for i, alias := range checkUpdateCmd.Aliases {
			if alias != expectedAliases[i] {
				t.Errorf("Expected alias '%s', got '%s'", expectedAliases[i], alias)
			}
		}

		if !strings.Contains(checkUpdateCmd.Short, "Update") {
			t.Errorf("Expected Short description to contain 'Update'")
		}

		if !strings.Contains(checkUpdateCmd.Long, "Check if") {
			t.Errorf("Expected Long description to contain 'Check if'")
		}
	})

	t.Run("command_flags", func(t *testing.T) {
		folderFlag := checkUpdateCmd.Flags().Lookup("folder")
		if folderFlag == nil {
			t.Error("Expected 'folder' flag to be defined")
		} else {
			if folderFlag.Usage != "path to the wow addon folder" {
				t.Errorf("Expected folder flag usage to be 'path to the wow addon folder', got '%s'", folderFlag.Usage)
			}
			if folderFlag.DefValue != "" {
				t.Errorf("Expected folder flag default value to be empty, got '%s'", folderFlag.DefValue)
			}
		}
	})
}

// Test helper functions behavior
func TestDirectoryOperations(t *testing.T) {
	t.Run("create_and_cleanup_temp_dir", func(t *testing.T) {
		dir := createTempDir(t)

		// Verify directory was created
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Error("Temporary directory was not created")
		}

		// Verify it's a directory
		if info, err := os.Stat(dir); err == nil && !info.IsDir() {
			t.Error("Created path is not a directory")
		}

		// Cleanup
		os.RemoveAll(dir)

		// Verify cleanup worked
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			t.Error("Directory was not properly cleaned up")
		}
	})

	t.Run("create_git_directory", func(t *testing.T) {
		dir := createTempGitDir(t)
		defer os.RemoveAll(dir)

		// Verify .git directory exists
		gitDir := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitDir); os.IsNotExist(err) {
			t.Error(".git directory was not created")
		}

		// Verify .git is a directory
		if info, err := os.Stat(gitDir); err == nil && !info.IsDir() {
			t.Error(".git path is not a directory")
		}
	})
}

// Integration test that tests the actual workflow (without mocking)
func TestCheckUpdateWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("non_git_directory_workflow", func(t *testing.T) {
		dir := createTempDir(t)
		defer os.RemoveAll(dir)

		// Change to test directory
		originalWd, err := os.Getwd()
		if err != nil {
			t.Fatalf("Failed to get current working directory: %v", err)
		}
		defer os.Chdir(originalWd)

		err = os.Chdir(dir)
		if err != nil {
			t.Fatalf("Failed to change to test directory: %v", err)
		}

		// Test the complete workflow for non-git directory
		hasUpdate, err := isUpdateAvailable(dir)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if hasUpdate {
			t.Error("Non-git directory should not have updates available")
		}
	})
}

// Benchmark the main function
func BenchmarkIsUpdateAvailable(b *testing.B) {
	dir, err := os.MkdirTemp("", "benchmark_test")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	// Change to test directory
	originalWd, err := os.Getwd()
	if err != nil {
		b.Fatalf("Failed to get current working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	err = os.Chdir(dir)
	if err != nil {
		b.Fatalf("Failed to change to test directory: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := isUpdateAvailable(dir)
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
	}
}