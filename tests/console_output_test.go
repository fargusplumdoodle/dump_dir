package tests

import (
	"github.com/fargusplumdoodle/dump_dir/tests/e2e"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestConsoleOutput(t *testing.T) {
	t.Run("test file content markers", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./test1.txt": "content 1",
				"./test2.txt": "content 2",
			}).
			WithArgs(".")

		result := env.Run()

		// Check if the console output contains the total number of files
		if !strings.Contains(result.Output, "Total files found: 2") {
			t.Error("Expected console output to show correct file count")
		}

		// Verify parsed files section exists
		if !strings.Contains(result.Output, "🔍 Parsed files:") {
			t.Error("Expected console output to show parsed files section")
		}

		// Check if both files are listed in the output
		if !strings.Contains(result.Output, "./test1.txt") || !strings.Contains(result.Output, "./test2.txt") {
			t.Error("Expected console output to list all processed files")
		}
	})

	t.Run("test line counting", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./test.txt": "line1\nline2\nline3\n",
			}).
			WithArgs(".")

		result := env.Run()

		// Extract line count from output
		re := regexp.MustCompile(`Total lines across all parsed files: (\d+)`)
		matches := re.FindStringSubmatch(result.Output)
		if len(matches) < 2 {
			t.Fatal("Could not find line count in output")
		}

		lineCount, _ := strconv.Atoi(matches[1])
		if lineCount != 3 {
			t.Errorf("Expected 3 lines, got %d", lineCount)
		}
	})

	t.Run("test token estimation accuracy", func(t *testing.T) {
		// Create a file with known token count (approximately)
		content := strings.Repeat("test word ", 100) // Should be around 200 tokens

		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./test.txt": content,
			}).
			WithArgs(".")

		result := env.Run()

		re := regexp.MustCompile(`Estimated tokens: (\d+(?:\.\d+)?)k`)
		matches := re.FindStringSubmatch(result.Output)
		if len(matches) < 2 {
			t.Fatal("Could not find token count in output")
		}

		tokenStr := matches[1]
		baseCount, _ := strconv.ParseFloat(tokenStr, 64)
		tokenCount := baseCount * 1000

		if tokenCount != 200 {
			t.Errorf("Token count %f is outside expected range (150-250)", tokenCount)
		}
	})

	t.Run("test different line endings", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./unix.txt":    "line1\nline2\nline3",
				"./windows.txt": "line1\r\nline2\r\nline3",
			}).
			WithArgs(".")

		result := env.Run()

		// Both files should report 3 lines regardless of line ending
		re := regexp.MustCompile(`Total lines across all parsed files: (\d+)`)
		matches := re.FindStringSubmatch(result.Output)
		if len(matches) < 2 {
			t.Fatal("Could not find line count in output")
		}

		lineCount, _ := strconv.Atoi(matches[1])
		if lineCount != 6 { // 3 lines per file
			t.Errorf("Expected 6 total lines, got %d", lineCount)
		}
	})

	t.Run("test skipped files output formatting", func(t *testing.T) {
		// Create a file that's too large (over 500KB)
		largeContent := strings.Repeat("a", 501*1024)

		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./normal.txt":   "normal content",
				"./toolarge.txt": largeContent,
				"./empty.txt":    "",
			}).
			WithArgs(".")

		result := env.Run()

		// Check for appropriate section headers
		if !strings.Contains(result.Output, "🔍 Parsed files:") {
			t.Error("Missing parsed files section")
		}
		if !strings.Contains(result.Output, "🪨 Skipped large files:") {
			t.Error("Missing skipped large files section")
		}

		// Verify file categorization
		if !strings.Contains(result.Output, "./normal.txt") {
			t.Error("Normal file should be listed in parsed files")
		}
		if !strings.Contains(result.Output, "./toolarge.txt") {
			t.Error("Large file should be listed in skipped files")
		}
		if !strings.Contains(result.Output, "./empty.txt") {
			t.Error("Empty file should be listed in parsed files")
		}
	})
}
