package tests

import (
	"github.com/fargusplumdoodle/dump_dir/tests/e2e"
	"path/filepath"
	"testing"
)

func TestPathHandling(t *testing.T) {
	t.Run("absolute paths", func(t *testing.T) {
		// Get absolute path for test directory
		absPath, _ := filepath.Abs(".")
		testPath := filepath.Join(absPath, "testdir")

		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				filepath.Join(testPath, "main.go"):   "package main\nfunc main() {}\n",
				filepath.Join(testPath, "helper.go"): "package main\nfunc helper() {}\n",
				"./outside.go":                       "package main\nfunc outside() {}\n",
			}).
			WithArgs(testPath)

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput(filepath.Join(testPath, "main.go")).
			AssertFileInOutput(filepath.Join(testPath, "helper.go")).
			AssertFileCount(2)
	})

	t.Run("relative paths", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./src/main.go":        "package main\nfunc main() {}\n",
				"./src/lib/helper.go":  "package lib\nfunc helper() {}\n",
				"./src/util/util.go":   "package util\nfunc util() {}\n",
				"./outside/outside.go": "package outside\nfunc outside() {}\n",
			}).
			WithArgs("./src")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./src/main.go").
			AssertFileInOutput("./src/lib/helper.go").
			AssertFileInOutput("./src/util/util.go").
			AssertFileCount(3)
	})

	t.Run("multiple input directories", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./src/main.go":     "package main\nfunc main() {}\n",
				"./lib/helper.go":   "package lib\nfunc helper() {}\n",
				"./test/test.go":    "package test\nfunc test() {}\n",
				"./docs/README.md":  "# Documentation\n",
				"./build/output.go": "package build\nfunc output() {}\n",
			}).
			WithArgs("./src ./lib ./test")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./src/main.go").
			AssertFileInOutput("./lib/helper.go").
			AssertFileInOutput("./test/test.go").
			AssertFileNotInOutput("./docs/README.md").
			AssertFileNotInOutput("./build/output.go").
			AssertFileCount(3)
	})

	t.Run("duplicate paths", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./src/main.go":   "package main\nfunc main() {}\n",
				"./src/helper.go": "package main\nfunc helper() {}\n",
			}).
			WithArgs("./src ./src")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./src/main.go").
			AssertFileInOutput("./src/helper.go").
			AssertFileCount(2) // Should only count files once even if path is duplicated
	})

	t.Run("nonexistent paths", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./src/main.go":   "package main\nfunc main() {}\n",
				"./lib/helper.go": "package lib\nfunc helper() {}\n",
			}).
			WithArgs("./src ./nonexistent ./lib")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./src/main.go").
			AssertFileInOutput("./lib/helper.go").
			AssertFileCount(2) // Should only count files from existing directories
	})
}
