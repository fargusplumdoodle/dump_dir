package tests

import (
	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"
	"os"
	"testing"

	. "github.com/fargusplumdoodle/dump_dir/src"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedConfig Config
	}{
		{
			name: "Default dump_dir action",
			args: []string{"go", "./src"},
			expectedConfig: BuildConfig(
				WithAction("dump_dir"),
				WithExtensions("go"),
				WithDirectories("./src"),
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
			args: []string{"go,js,py", "./src", "./tests"},
			expectedConfig: BuildConfig(
				WithAction("dump_dir"),
				WithExtensions("go", "js", "py"),
				WithDirectories("./src", "./tests"),
			),
		},
		{
			name: "Skip directories",
			args: []string{"go", "./src", "-s", "./src/vendor", "-s", "./src/generated"},
			expectedConfig: BuildConfig(
				WithAction("dump_dir"),
				WithExtensions("go"),
				WithDirectories("./src"),
				WithSkipDirs("./src/vendor", "./src/generated"),
			),
		},
		{
			name: "Include ignored files",
			args: []string{"go", "./src", "--include-ignored"},
			expectedConfig: BuildConfig(
				WithAction("dump_dir"),
				WithExtensions("go"),
				WithDirectories("./src"),
				WithIncludeIgnored(true),
			),
		},
		{
			name: "Complex case",
			args: []string{"go,js", "./src", "./tests", "-s", "./src/vendor", "--include-ignored", "./config.go"},
			expectedConfig: BuildConfig(
				WithAction("dump_dir"),
				WithExtensions("go", "js"),
				WithDirectories("./src", "./tests"),
				WithSkipDirs("./src/vendor"),
				WithIncludeIgnored(true),
				WithSpecificFiles("./config.go"),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ParseArgs(tt.args)
			if diff := cmp.Diff(tt.expectedConfig, config); diff != "" {
				t.Errorf("ParseArgs() mismatch (-want +got):\n%s", diff)
			}
		})
	}

	t.Run("Filesystem-aware specific files and directories", func(t *testing.T) {
		// Set up a mock file system
		fs := afero.NewMemMapFs()
		afero.WriteFile(fs, "/project/src/main.go", []byte("package main"), 0644)
		afero.WriteFile(fs, "/project/config.json", []byte("{}"), 0644)
		fs.MkdirAll("/project/tests", 0755)

		// Store the original os.Stat function
		originalStat := OsStat
		// Replace os.Stat with our mock version
		OsStat = func(name string) (os.FileInfo, error) {
			return fs.Stat(name)
		}
		// Restore the original os.Stat function after the test
		defer func() { OsStat = originalStat }()

		args := []string{"go,json", "/project/src/main.go", "/project/config.json", "/project/tests"}
		config := ParseArgs(args)

		expectedConfig := BuildConfig(
			WithAction("dump_dir"),
			WithExtensions("go", "json"),
			WithSpecificFiles("/project/src/main.go", "/project/config.json"),
			WithDirectories("/project/tests"),
		)

		if diff := cmp.Diff(expectedConfig, config); diff != "" {
			t.Errorf("ParseArgs() mismatch (-want +got):\n%s", diff)
		}
	})
}
