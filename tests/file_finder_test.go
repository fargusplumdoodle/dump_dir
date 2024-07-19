package tests

import (
	. "github.com/fargusplumdoodle/dump_dir/src"
	"testing"
)

// TestFileFinder runs tests for the FileFinder
func TestFileFinder(t *testing.T) {
	tests := []struct {
		name            string
		fileSystem      map[string]string
		config          Config
		expectedFiles   []string
		unexpectedFiles []string
	}{
		{
			name: "Find Go files in src directory",
			fileSystem: map[string]string{
				"src/file1.go":        "package main",
				"src/file2.go":        "package main",
				"src/subdir/file3.go": "package subdir",
				"src/file4.txt":       "text file",
			},
			config: BuildConfig(),
			expectedFiles: []string{
				"src/file1.go",
				"src/file2.go",
				"src/subdir/file3.go",
			},
			unexpectedFiles: []string{
				"src/file4.txt",
			},
		},
		{
			name: "Find specific file types",
			fileSystem: map[string]string{
				"src/file1.go":  "package main",
				"src/file2.js":  "console.log('Hello');",
				"src/file3.py":  "print('Hello')",
				"src/file4.txt": "text file",
			},
			config: BuildConfig(WithExtensions("js", "py")),
			expectedFiles: []string{
				"src/file2.js",
				"src/file3.py",
			},
			unexpectedFiles: []string{
				"src/file1.go",
				"src/file4.txt",
			},
		},
		{
			name: "Skip specified directories",
			fileSystem: map[string]string{
				"src/file1.go":        "package main",
				"src/subdir/file2.go": "package subdir",
				"src/skipme/file3.go": "package skipme",
			},
			config: BuildConfig(WithSkipDirs("src/skipme")),
			expectedFiles: []string{
				"src/file1.go",
				"src/subdir/file2.go",
			},
			unexpectedFiles: []string{
				"src/skipme/file3.go",
			},
		},
		{
			name: "Include specific files",
			fileSystem: map[string]string{
				"src/file1.go":  "package main",
				"src/file2.js":  "console.log('Hello');",
				"root/file3.go": "package root",
			},
			config: BuildConfig(WithSpecificFiles("root/file3.go")),
			expectedFiles: []string{
				"src/file1.go",
				"root/file3.go",
			},
			unexpectedFiles: []string{
				"src/file2.js",
			},
		},
		{
			name: "Leading ./ is optional",
			fileSystem: map[string]string{
				"src/file1.go":   "package main",
				"tests/file2.go": "console.log('Hello');",
				"root/file3.go":  "package root",
			},
			config: BuildConfig(WithDirectories("src", "tests", "./root")),
			expectedFiles: []string{
				"src/file1.go",
				"tests/file2.go",
				"root/file3.go",
			},
			unexpectedFiles: []string{},
		},
		{
			name: "Can ignore contents of subdirectories",
			fileSystem: map[string]string{
				"src/ignore/notfound.go":  "package main",
				"src/file1.go":            "console.log('Hello');",
				"src/dir/file2.go":        "package root",
				"src/dir/ignore/file2.go": "package root",
			},
			config: BuildConfig(WithSkipDirs("src/ignore", "src/dir/ignore/")),
			expectedFiles: []string{
				"src/file1.go",
				"src/dir/file2.go",
			},
			unexpectedFiles: []string{
				"src/ignore/notfound.go",
				"src/dir/ignore/file2.go",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := setupTestFileSystem(tt.fileSystem)
			fileFinder := NewFileFinder(tt.config, fs)

			foundFiles := fileFinder.DiscoverFiles()

			assertFilesFound(t, foundFiles, tt.expectedFiles, tt.unexpectedFiles)
		})
	}
}
