package tests

import (
	"testing"
)

func TestDumpDir(t *testing.T) {
	t.Run("basic file dump", func(t *testing.T) {
		env := NewTestEnvironment(t).
			WithWorkingDir("/test/project").
			WithFiles(map[string]string{
				"/test/project/main.go": "package main\n\nfunc main() {}\n",
				"/test/project/util.go": "package main\n\nfunc helper() {}\n",
			}).
			WithArgs("-e go /test/project")

		env.Run().
			AssertOutputContains("✅ File contents have been copied to clipboard").
			AssertClipboardContains("START FILE: /test/project/main.go").
			AssertClipboardContains("START FILE: /test/project/util.go")
	})

	t.Run("respects gitignore", func(t *testing.T) {
		env := NewTestEnvironment(t).
			WithWorkingDir("/test/project").
			WithFiles(map[string]string{
				"/test/project/.gitignore":        "*.generated.go\n",
				"/test/project/main.go":           "package main\n",
				"/test/project/auto.generated.go": "package main\n",
			}).
			WithArgs("-e go /test/project")

		env.Run().
			AssertOutputContains("✅ File contents have been copied to clipboard").
			AssertClipboardContains("START FILE: /test/project/main.go").
	})

	t.Run("handles binary files", func(t *testing.T) {
		env := NewTestEnvironment(t).
			WithWorkingDir("/test/project").
			WithFiles(map[string]string{
				"/test/project/main.go":    "package main\n",
				"/test/project/binary.exe": "\x00\x01\x02\x03",
			}).
			WithArgs("dump_dir /test/project")

		env.Run().
			AssertOutputContains("💽 Skipped binary files:").
			AssertOutputContains("binary.exe")
	})
}
