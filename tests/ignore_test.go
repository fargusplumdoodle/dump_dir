package tests

import (
	. "github.com/fargusplumdoodle/dump_dir/src"
	"github.com/spf13/afero"
	"os/exec"
	"testing"
)

// TestIgnoreFunctionality runs tests for the ignore functionality
func TestIgnoreFunctionality(t *testing.T) {
	tests := []struct {
		name            string
		files           []string
		globalGitignore string
		localGitignore  string
		config          *Config
		expectedFiles   []string
		unexpectedFiles []string
	}{
		{
			name: "Ignore .git directory",
			files: []string{
				"src/file1.go",
				".git/config",
				".git/objects/abc",
			},
			config: BuildConfig(
				WithExtensions("go"),
				WithDirectories("src"),
			),
			expectedFiles: []string{
				"src/file1.go",
			},
			unexpectedFiles: []string{
				".git/config",
				".git/objects/abc",
			},
		},
		{
			name: "Respect global gitignore",
			files: []string{
				"src/file1.go",
				"src/file2.log",
				"build/output.txt",
			},
			globalGitignore: "*.log\nbuild/",
			config: BuildConfig(
				WithDirectories("src"),
				WithExtensions("go"),
			),
			expectedFiles: []string{
				"src/file1.go",
			},
			unexpectedFiles: []string{
				"src/file2.log",
				"build/output.txt",
			},
		},
		{
			name: "Respect local .gitignore",
			files: []string{
				"src/file1.go",
				"src/file2.tmp",
				"dist/bundle.js",
			},
			localGitignore: "*.tmp\ndist/",
			config: BuildConfig(
				WithDirectories("src"),
				WithExtensions("go"),
			),
			expectedFiles: []string{
				"src/file1.go",
			},
			unexpectedFiles: []string{
				"src/file2.tmp",
				"dist/bundle.js",
			},
		},
		{
			name: "Include ignored files when specified",
			files: []string{
				"src/file1.go",
				"src/file2.tmp",
				"dist/bundle.js",
			},
			localGitignore: "*.tmp\ndist/",
			config: BuildConfig(
				WithExtensions("go", "tmp", "js"),
				WithDirectories("."),
				WithIncludeIgnored(true),
			),
			expectedFiles: []string{
				"src/file1.go",
				"src/file2.tmp",
				"dist/bundle.js",
			},
			unexpectedFiles: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := setupIgnoreTestEnvironment(tt.files, tt.globalGitignore, tt.localGitignore)
			fileFinder := NewFileFinder(*tt.config, fs)

			foundFiles := fileFinder.DiscoverFiles()

			assertFilesFound(t, foundFiles, tt.expectedFiles, tt.unexpectedFiles)
		})
	}
}

// Helper functions

func setupIgnoreTestEnvironment(files []string, globalGitignore, localGitignore string) afero.Fs {
	fs := afero.NewMemMapFs()

	// Create test files
	for _, file := range files {
		afero.WriteFile(fs, file, []byte("content"), 0644)
	}

	// Set up global gitignore
	if globalGitignore != "" {
		afero.WriteFile(fs, "/home/user/.gitignore_global", []byte(globalGitignore), 0644)
		// Mock the global gitignore path
		oldExec := ExecCommand
		ExecCommand = func(name string, arg ...string) *exec.Cmd {
			return exec.Command("echo", "/home/user/.gitignore_global")
		}
		defer func() { ExecCommand = oldExec }()
	}

	// Set up local .gitignore
	if localGitignore != "" {
		afero.WriteFile(fs, ".gitignore", []byte(localGitignore), 0644)
	}

	return fs
}

func assertFilesFound(t *testing.T, foundFiles, expectedFiles, unexpectedFiles []string) {
	for _, expectedFile := range expectedFiles {
		if !contains(foundFiles, expectedFile) {
			t.Errorf("Expected file not found: %s", expectedFile)
		}
	}

	for _, unexpectedFile := range unexpectedFiles {
		if contains(foundFiles, unexpectedFile) {
			t.Errorf("Unexpected file found: %s", unexpectedFile)
		}
	}
}
