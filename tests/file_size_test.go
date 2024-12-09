package tests

import (
	"github.com/fargusplumdoodle/dump_dir/tests/e2e"
	"strings"
	"testing"
)

func TestFileSizeHandling(t *testing.T) {
	t.Run("default size limit (500KB)", func(t *testing.T) {
		// Create a string just at the 500KB limit
		content := strings.Repeat("a", 500*1024)
		largeContent := strings.Repeat("b", 501*1024) // Just over limit

		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./at-limit.txt":        content,
				"./just-over-limit.txt": largeContent,
				"./small.txt":           "small file content",
			}).
			WithArgs(".")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./at-limit.txt").
			AssertFileInOutput("./just-over-limit.txt").
			AssertFileInOutput("./small.txt").
			AssertFileTooLarge("./just-over-limit.txt", 501*1024).
			AssertFileCount(3) // All files should be counted, including large ones
	})

	t.Run("custom max file size with different units", func(t *testing.T) {
		smallContent := strings.Repeat("a", 100)     // 100B
		mediumContent := strings.Repeat("b", 2*1024) // 2KB
		largeContent := strings.Repeat("c", 3*1024)  // 3KB

		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./small.txt":  smallContent,
				"./medium.txt": mediumContent,
				"./large.txt":  largeContent,
			}).
			WithArgs(". --max-filesize 2KB")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./small.txt").
			AssertFileInOutput("./medium.txt").
			AssertFileInOutput("./large.txt").
			AssertFileTooLarge("./large.txt", 3*1024).
			AssertFileCount(3) // All files should be counted
	})

	t.Run("mixed files above and below size limit", func(t *testing.T) {
		withinLimit := strings.Repeat("a", 200*1024)    // 200KB
		exceedsLimit := strings.Repeat("b", 600*1024)   // 600KB
		wayTooLarge := strings.Repeat("c", 2*1024*1024) // 2MB

		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./within-limit.txt":  withinLimit,
				"./exceeds-limit.txt": exceedsLimit,
				"./way-too-large.txt": wayTooLarge,
				"./tiny.txt":          "tiny file",
			}).
			WithArgs(".")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./within-limit.txt").
			AssertFileInOutput("./exceeds-limit.txt").
			AssertFileInOutput("./way-too-large.txt").
			AssertFileInOutput("./tiny.txt").
			AssertFileTooLarge("./exceeds-limit.txt", 600*1024).
			AssertFileTooLarge("./way-too-large.txt", 2*1024*1024).
			AssertFileCount(4) // All files should be counted
	})

	t.Run("empty files", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./empty1.txt": "",
				"./normal.txt": "some content",
			}).
			WithArgs(".")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./empty1.txt").
			AssertFileInOutput("./normal.txt").
			AssertEmptyFile("./empty1.txt").
			AssertFileCount(2)
	})

	t.Run("size limit with MB units", func(t *testing.T) {
		content := strings.Repeat("a", 1024*1024)        // 1MB
		largeContent := strings.Repeat("b", 3*1024*1024) // 3MB

		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./onemb.txt":   content,
				"./threemb.txt": largeContent,
			}).
			WithArgs(". --max-filesize 2MB")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./onemb.txt").
			AssertFileInOutput("./threemb.txt").
			AssertFileTooLarge("./threemb.txt", 3*1024*1024).
			AssertFileCount(2) // Both files should be counted
	})
}
