package tests

import (
	"github.com/fargusplumdoodle/dump_dir/tests/e2e"
	"strings"
	"testing"
)

func TestClipboardOutput(t *testing.T) {
	t.Run("normal file content", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./test.txt": "Hello\nWorld\n",
			}).
			WithArgs(".")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileCount(1).
			AssertWholeFileContent("./test.txt", "Hello\nWorld\n")
	})

	t.Run("empty file content", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./empty.txt": "",
			}).
			WithArgs(".")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileCount(1).
			AssertEmptyFile("./empty.txt")
	})

	t.Run("large file content", func(t *testing.T) {
		// Create a file larger than 500KB (default limit)
		largeContent := strings.Repeat("a", 501*1024)

		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./large.txt": largeContent,
			}).
			WithArgs(".")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileCount(1).
			AssertFileTooLarge("./large.txt", 501*1024)
	})

	t.Run("binary file content", func(t *testing.T) {
		// Create a binary-like file with null bytes
		binaryContent := string([]byte{0x00, 0x01, 0x02, 0x03})

		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./binary.bin": binaryContent,
			}).
			WithArgs(".")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileCount(1).
			AssertBinaryFile("./binary.bin")
	})

	t.Run("mixed file types", func(t *testing.T) {
		largeContent := strings.Repeat("a", 501*1024)
		binaryContent := string([]byte{0x00, 0x01, 0x02, 0x03})

		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./normal.txt": "Hello\nWorld\n",
				"./empty.txt":  "",
				"./large.txt":  largeContent,
				"./binary.bin": binaryContent,
			}).
			WithArgs(".")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileCount(4).
			AssertWholeFileContent("./normal.txt", "Hello\nWorld\n").
			AssertEmptyFile("./empty.txt").
			AssertFileTooLarge("./large.txt", 501*1024).
			AssertBinaryFile("./binary.bin")
	})
}
