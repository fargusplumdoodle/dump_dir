package tests

import (
	"fmt"
	"github.com/fargusplumdoodle/dump_dir/tests/e2e"
	"testing"
)

func TestGitignoreIntegration(t *testing.T) {
	t.Run("respects basic gitignore patterns", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				".gitignore":                "*.log\nnode_modules/\nbuild/\n",
				"main.go":                   "package main\nfunc main() {}\n",
				"app.log":                   "some logs\n",
				"node_modules/package.json": "{\n  \"name\": \"test\"\n}\n",
				"build/output.js":           "console.log('built');\n",
				"src/feature.go":            "package src\nfunc feature() {}\n",
			}).
			WithArgs(".")

		result := env.Run()
		fmt.Println(result.Output)

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./main.go").
			AssertFileInOutput("./.gitignore").
			AssertFileInOutput("./src/feature.go").
			AssertFileCount(3)
	})

	t.Run("includes ignored files with --include-ignored flag", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				".gitignore":                "*.log\nnode_modules/\nbuild/\n",
				"main.go":                   "package main\nfunc main() {}\n",
				"app.log":                   "some logs\n",
				"node_modules/package.json": "{\n  \"name\": \"test\"\n}\n",
				"build/output.js":           "console.log('built');\n",
				"src/feature.go":            "package src\nfunc feature() {}\n",
			}).
			WithArgs(". --include-ignored")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./.gitignore").
			AssertFileInOutput("./main.go").
			AssertFileInOutput("./app.log").
			AssertFileInOutput("./node_modules/package.json").
			AssertFileInOutput("./build/output.js").
			AssertFileInOutput("./src/feature.go").
			AssertFileCount(6)
	})
}
