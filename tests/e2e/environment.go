package e2e

import (
	"bytes"
	"github.com/fargusplumdoodle/dump_dir/src"
	"github.com/spf13/afero"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Environment encapsulates all the mocked dependencies and utilities
// needed for end-to-end testing
type Environment struct {
	t          *testing.T
	fs         afero.Fs
	stdout     *bytes.Buffer
	stderr     *bytes.Buffer
	clipboard  *MockClipboard
	originalWd string
	currentWd  string
	args       []string
}

// NewEnvironment creates a new test environment with mocked dependencies
func NewEnvironment(t *testing.T) *Environment {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	fs := afero.NewMemMapFs()

	// Mock the OsStat function to use our memory filesystem
	src.OsStat = func(name string) (os.FileInfo, error) {
		return fs.Stat(name)
	}

	return &Environment{
		t:         t,
		fs:        fs,
		stdout:    stdout,
		stderr:    stderr,
		clipboard: NewMockClipboard(),
	}
}

// WithArgs sets the command line arguments for the test
func (e *Environment) WithArgs(args string) *Environment {
	e.args = strings.Fields(args)
	return e
}

// WithFiles creates files in the virtual filesystem
func (e *Environment) WithFiles(files map[string]string) *Environment {
	for path, content := range files {
		// Clean the path to ensure consistent formatting
		cleanPath := filepath.Clean(path)

		// Create directory if it doesn't exist
		dir := filepath.Dir(cleanPath)
		if dir != "" {
			err := e.fs.MkdirAll(dir, 0755)
			if err != nil {
				e.t.Fatalf("Failed to create directory %s: %v", dir, err)
			}
		}

		// Write the file
		err := afero.WriteFile(e.fs, cleanPath, []byte(content), 0644)
		if err != nil {
			e.t.Fatalf("Failed to write file %s: %v", cleanPath, err)
		}
	}
	return e
}

// WithWorkingDir sets up a working directory for the test
func (e *Environment) WithWorkingDir(path string) *Environment {
	cleanPath := filepath.Clean(path)
	err := e.fs.MkdirAll(cleanPath, 0755)
	if err != nil {
		e.t.Fatalf("Failed to create working directory %s: %v", cleanPath, err)
	}
	e.currentWd = cleanPath
	return e
}

// Run executes the dump_dir command in the test environment
func (e *Environment) Run() *Result {
	// Save original stdout/stderr
	originalStdout := os.Stdout
	originalStderr := os.Stderr

	// Create pipe for capturing Output
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	// Run the command
	err := src.Run(e.args, src.RunConfig{
		Fs:        e.fs,
		Clipboard: e.clipboard,
		Version:   "test",
		Commit:    "test",
		Date:      "test",
	})

	// Restore original stdout/stderr
	w.Close()
	os.Stdout = originalStdout
	os.Stderr = originalStderr

	// Capture Output
	var buf bytes.Buffer
	io.Copy(&buf, r)

	return &Result{
		env:       e,
		Output:    buf.String(),
		Clipboard: e.clipboard.Content,
		Err:       err,
	}
}
