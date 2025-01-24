package tests

import (
	"github.com/fargusplumdoodle/dump_dir/tests/e2e"
	"testing"
)

func TestDumpDir(t *testing.T) {
	t.Run("basic file dump", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"src/main.go": "package main\n\nfunc main() {}\n",
				"src/util.go": "package main\n\nfunc helper() {}\n",
			}).
			WithArgs("-e go src")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./src/main.go").
			AssertFileInOutput("./src/util.go").
			AssertFileCount(2).
			AssertLineCount(6).
			AssertTokenCount(26)
	})

}
