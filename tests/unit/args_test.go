package unit

import (
	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"
	"os"
	"testing"

	. "github.com/fargusplumdoodle/dump_dir/src"
)

func TestParseArgs(t *testing.T) {
	// Set up a test filesystem with some sample files and directories
	fs := afero.NewMemMapFs()

	// Create test directory structure
	testDirs := []string{
		"project/src",
		"project/src/vendor",
		"project/src/generated",
		"project/tests",
		"project/dist",
	}

	for _, dir := range testDirs {
		fs.MkdirAll(dir, 0755)
	}

	// Create some test files
	testFiles := map[string]string{
		"project/src/main.go":        "package main",
		"project/src/util.go":        "package main",
		"project/tests/main_test.go": "package tests",
		"project/config.json":        "{}",
	}

	for path, content := range testFiles {
		afero.WriteFile(fs, path, []byte(content), 0644)
	}

	// Mock os.Stat to use our memory filesystem
	originalStat := OsStat
	OsStat = func(name string) (os.FileInfo, error) {
		return fs.Stat(name)
	}
	defer func() { OsStat = originalStat }()

	tests := []struct {
		name           string
		args           []string
		expectedConfig *Config
		expectedError  error
	}{
		{
			name: "Default dump_dir action",
			args: []string{"project/src"},
			expectedConfig: BuildConfig(
				WithAction("dump_dir"),
				WithDirectories("./project/src"),
			),
		},
		{
			name: "Help action",
			args: []string{"--help"},
			expectedConfig: BuildConfig(
				WithAction("help"),
			),
		},
		{
			name: "Version action",
			args: []string{"--version"},
			expectedConfig: BuildConfig(
				WithAction("version"),
			),
		},
		{
			name: "Multiple extensions and directories",
			args: []string{"project/src", "project/tests", "-e", "go,json"},
			expectedConfig: BuildConfig(
				WithAction("dump_dir"),
				WithExtensions("go", "json"),
				WithDirectories("./project/src", "./project/tests"),
			),
		},
		{
			name: "Skip directories",
			args: []string{"project/src", "-s", "project/src/vendor", "--skip", "project/src/generated"},
			expectedConfig: BuildConfig(
				WithAction("dump_dir"),
				WithDirectories("./project/src"),
				WithSkipDirs("./project/src/vendor", "./project/src/generated"),
			),
		},
		{
			name: "Include ignored files",
			args: []string{"project/src", "--include-ignored"},
			expectedConfig: BuildConfig(
				WithAction("dump_dir"),
				WithDirectories("./project/src"),
				WithIncludeIgnored(true),
			),
		},
		{
			name: "Complex case",
			args: []string{
				"-e", "go,json",
				"project/src",
				"project/tests",
				"-s", "project/src/vendor",
				"--include-ignored",
				"project/config.json",
			},
			expectedConfig: BuildConfig(
				WithAction("dump_dir"),
				WithExtensions("go", "json"),
				WithDirectories("./project/src", "./project/tests"),
				WithSkipDirs("./project/src/vendor"),
				WithIncludeIgnored(true),
				WithSpecificFiles("./project/config.json"),
			),
		},
		{
			name: "Max filesize in Bytes",
			args: []string{"project/src", "--max-filesize", "1000B"},
			expectedConfig: BuildConfig(
				WithAction("dump_dir"),
				WithDirectories("./project/src"),
				WithMaxFileSize(1000),
			),
		},
		{
			name: "Max filesize in Kilobytes",
			args: []string{"project/src", "-m", "500KB"},
			expectedConfig: BuildConfig(
				WithAction("dump_dir"),
				WithDirectories("./project/src"),
				WithMaxFileSize(500*1024),
			),
		},
		{
			name: "Max filesize in Megabytes",
			args: []string{"project/src", "--max-filesize", "2MB"},
			expectedConfig: BuildConfig(
				WithAction("dump_dir"),
				WithDirectories("./project/src"),
				WithMaxFileSize(2*1024*1024),
			),
		},
		{
			name:           "Invalid max filesize",
			args:           []string{"project/src", "--max-filesize", "invalid"},
			expectedConfig: nil,
			expectedError:  ErrInvalidMaxFileSize{Value: "invalid"},
		},
		{
			name:           "Missing max filesize value",
			args:           []string{"project/src", "--max-filesize"},
			expectedConfig: nil,
			expectedError:  ErrInvalidMaxFileSize{Value: ""},
		},
		{
			name: "Single glob pattern",
			args: []string{".", "-g", "*.go"},
			expectedConfig: BuildConfig(
				WithAction("dump_dir"),
				WithDirectories("."),
				WithGlobPatterns("*.go"),
			),
		},
		{
			name: "Multiple glob patterns",
			args: []string{".", "-g", "*.go", "-g", "*.md"},
			expectedConfig: BuildConfig(
				WithAction("dump_dir"),
				WithDirectories("."),
				WithGlobPatterns("*.go", "*.md"),
			),
		},
		{
			name: "Glob with other arguments",
			args: []string{".", "-g", "*.go", "-s", "node_modules"},
			expectedConfig: BuildConfig(
				WithAction("dump_dir"),
				WithDirectories("."),
				WithGlobPatterns("*.go"),
				WithSkipDirs("node_modules"),
			),
		},
		{
			name: "Using --glob instead of -g",
			args: []string{".", "--glob", "*.go"},
			expectedConfig: BuildConfig(
				WithAction("dump_dir"),
				WithDirectories("."),
				WithGlobPatterns("*.go"),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := ParseArgs(tt.args)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error %v, but got nil", tt.expectedError)
				} else if err != tt.expectedError {
					t.Errorf("Expected error %v, but got %v", tt.expectedError, err)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.expectedConfig != nil {
				if diff := cmp.Diff(*tt.expectedConfig, config); diff != "" {
					t.Errorf("ParseArgs() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}

	// Keep the filesystem-aware test case separate as it uses a different file structure
	t.Run("Filesystem-aware specific files and directories", func(t *testing.T) {
		// Create a new filesystem for this specific test
		fs := afero.NewMemMapFs()
		fs.MkdirAll("project/src", 0755)
		fs.MkdirAll("project/tests", 0755)
		afero.WriteFile(fs, "project/src/main.go", []byte("package main"), 0644)
		afero.WriteFile(fs, "project/config.json", []byte("{}"), 0644)

		// Update os.Stat for this test
		OsStat = func(name string) (os.FileInfo, error) {
			return fs.Stat(name)
		}

		args := []string{"-e", "go,json", "project/src/main.go", "project/config.json", "project/tests"}
		config, err := ParseArgs(args)

		expectedConfig := BuildConfig(
			WithAction("dump_dir"),
			WithExtensions("go", "json"),
			WithSpecificFiles("./project/src/main.go", "./project/config.json"),
			WithDirectories("./project/tests"),
		)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if diff := cmp.Diff(expectedConfig, &config); diff != "" {
			t.Errorf("ParseArgs() mismatch (-want +got):\n%s", diff)
		}
	})
}
