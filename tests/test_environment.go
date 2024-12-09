package tests

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

// TestEnvironment encapsulates all the mocked dependencies and utilities
// needed for end-to-end testing
type TestEnvironment struct {
	t          *testing.T
	fs         afero.Fs
	stdout     *bytes.Buffer
	stderr     *bytes.Buffer
	clipboard  *MockClipboard
	originalWd string
	currentWd  string
	args       []string
}

// NewTestEnvironment creates a new test environment with mocked dependencies
func NewTestEnvironment(t *testing.T) *TestEnvironment {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	fs := afero.NewMemMapFs()

	// Mock the OsStat function to use our memory filesystem
	src.OsStat = func(name string) (os.FileInfo, error) {
		return fs.Stat(name)
	}

	return &TestEnvironment{
		t:         t,
		fs:        fs,
		stdout:    stdout,
		stderr:    stderr,
		clipboard: NewMockClipboard(),
	}
}

// WithArgs sets the command line arguments for the test
func (e *TestEnvironment) WithArgs(args string) *TestEnvironment {
	e.args = strings.Fields(args)
	return e
}

// WithFiles creates files in the virtual filesystem
func (e *TestEnvironment) WithFiles(files map[string]string) *TestEnvironment {
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
func (e *TestEnvironment) WithWorkingDir(path string) *TestEnvironment {
	cleanPath := filepath.Clean(path)
	err := e.fs.MkdirAll(cleanPath, 0755)
	if err != nil {
		e.t.Fatalf("Failed to create working directory %s: %v", cleanPath, err)
	}
	e.currentWd = cleanPath
	return e
}

// Run executes the dump_dir command in the test environment
func (e *TestEnvironment) Run() *TestResult {
	// Save original stdout/stderr
	originalStdout := os.Stdout
	originalStderr := os.Stderr

	// Create pipe for capturing output
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

	// Capture output
	var buf bytes.Buffer
	io.Copy(&buf, r)

	return &TestResult{
		env:       e,
		output:    buf.String(),
		clipboard: e.clipboard.Content,
		err:       err,
	}
}

// TestResult contains the results of running the command
type TestResult struct {
	env       *TestEnvironment
	output    string
	clipboard string
	err       error
}

// AssertOutputContains checks if the command output contains expected content
func (r *TestResult) AssertOutputContains(expected string) *TestResult {
	if !strings.Contains(r.output, expected) {
		r.env.t.Errorf("Expected output to contain %q, got: %q", expected, r.output)
	}
	return r
}

// AssertClipboardContains checks if the clipboard contains expected content
func (r *TestResult) AssertClipboardContains(expected string) *TestResult {
	if !strings.Contains(r.clipboard, expected) {
		r.env.t.Errorf("Expected clipboard to contain %q, got: %q", expected, r.clipboard)
	}
	return r
}

// AssertNoError checks if the command completed without error
func (r *TestResult) AssertNoError() *TestResult {
	if r.err != nil {
		r.env.t.Errorf("Expected no error, got: %v", r.err)
	}
	return r
}

// AssertError checks if the command completed with an error
func (r *TestResult) AssertError() *TestResult {
	if r.err == nil {
		r.env.t.Error("Expected an error, got none")
	}
	return r
}
