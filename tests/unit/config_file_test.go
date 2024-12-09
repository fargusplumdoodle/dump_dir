package unit

import (
	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"
	"os"
	"testing"

	. "github.com/fargusplumdoodle/dump_dir/src"
)

func TestConfigLoader(t *testing.T) {
	tests := []struct {
		name           string
		configContent  string
		baseConfig     Config
		expectedConfig Config
		expectError    bool
	}{
		{
			name: "No config file present",
			baseConfig: *BuildConfig(
				WithDirectories("./src"),
			),
			expectedConfig: *BuildConfig(
				WithDirectories("./src"),
			),
		},
		{
			name:          "Empty config file",
			configContent: ``,
			baseConfig: *BuildConfig(
				WithDirectories("./src"),
			),
			expectedConfig: *BuildConfig(
				WithDirectories("./src"),
			),
		},
		{
			name: "Config with ignore paths",
			configContent: `
ignore:
  - ./src/subdir
`,
			baseConfig: *BuildConfig(
				WithDirectories("./src"),
			),
			expectedConfig: *BuildConfig(
				WithDirectories("./src"),
				WithSkipDirs("./src/subdir"),
			),
		},
		{
			name: "Config with include paths",
			configContent: `
include:
  - ./src/main.go
  - ./src/subdir
`,
			baseConfig: *BuildConfig(),
			expectedConfig: *BuildConfig(
				WithSpecificFiles("./src/main.go"),
				WithDirectories("./src/subdir"),
			),
		},
		{
			name: "Config with both include and ignore",
			configContent: `
include:
  - ./src/main.go
  - ./src
ignore:
  - ./src/subdir
`,
			baseConfig: *BuildConfig(),
			expectedConfig: *BuildConfig(
				WithSpecificFiles("./src/main.go"),
				WithDirectories("./src"),
				WithSkipDirs("./src/subdir"),
			),
		},
		{
			name: "Invalid YAML",
			configContent: `
include: [
  - invalid
  yaml: content
`,
			baseConfig:  *BuildConfig(),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up test filesystem with all necessary paths
			fs := afero.NewMemMapFs()

			// Create the base directory structure and files
			fs.MkdirAll("src/subdir", 0755)
			afero.WriteFile(fs, "src/main.go", []byte("package main"), 0644)
			afero.WriteFile(fs, "src/subdir/file.go", []byte("package subdir"), 0644)

			// Mock the OsStat function to use our memory filesystem
			originalStat := OsStat
			OsStat = func(name string) (os.FileInfo, error) {
				return fs.Stat(name)
			}
			defer func() { OsStat = originalStat }()

			// Add config file if content provided
			if tt.configContent != "" {
				err := afero.WriteFile(fs, ".dump_dir.yml", []byte(tt.configContent), 0644)
				if err != nil {
					t.Fatalf("Failed to write config file: %v", err)
				}
			}

			// Create config loader and load config
			loader := NewConfigLoader(fs)
			config, err := loader.LoadAndMergeConfig(tt.baseConfig)

			// Check error expectations
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Compare configs
			if diff := cmp.Diff(tt.expectedConfig, config); diff != "" {
				t.Errorf("Config mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
