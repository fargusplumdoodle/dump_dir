package tests

import (
	"github.com/fargusplumdoodle/dump_dir/tests/e2e"
	"testing"
)

func TestDumpDir(t *testing.T) {
	t.Run("basic file dump", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithWorkingDir("/test/project").
			WithFiles(map[string]string{
				"/test/project/main.go": "package main\n\nfunc main() {}\n",
				"/test/project/util.go": "package main\n\nfunc helper() {}\n",
			}).
			WithArgs("-e go /test/project")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("/test/project/main.go").
			AssertFileInOutput("/test/project/util.go").
			AssertFileCount(2).
			AssertLineCount(6).
			AssertTokenCount(26)
	})

}
