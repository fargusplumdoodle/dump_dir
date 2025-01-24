package tests

import (
	"github.com/fargusplumdoodle/dump_dir/tests/e2e"
	"testing"
)

func TestExtensionFiltering(t *testing.T) {
	t.Run("multiple extensions", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./main.go":    "package main\n\nfunc main() {}\n",
				"./script.py":  "def main():\n    pass\n",
				"./index.js":   "console.log('hello');\n",
				"./README.md":  "# Test Project\n",
				"./styles.css": "body { margin: 0; }\n",
			}).
			WithArgs("-e go,py,js .")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./main.go").
			AssertFileInOutput("./script.py").
			AssertFileInOutput("./index.js").
			AssertFileCount(3)
	})

	t.Run("nonexistent extension", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./main.go":  "package main\n\nfunc main() {}\n",
				"./test.txt": "some text\n",
			}).
			WithArgs("-e xyz .")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileCount(0)
	})

	t.Run("mixed valid and invalid extensions", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./main.go":  "package main\n\nfunc main() {}\n",
				"./test.txt": "some text\n",
				"./util.py":  "def util():\n    pass\n",
			}).
			WithArgs("-e go,xyz,py .")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./main.go").
			AssertFileInOutput("./util.py").
			AssertFileCount(2)
	})

	t.Run("no extension specified", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./main.go":  "package main\n\nfunc main() {}\n",
				"./test.txt": "some text\n",
				"./README":   "# Test Project\n",
			}).
			WithArgs(".")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./main.go").
			AssertFileInOutput("./test.txt").
			AssertFileInOutput("./README").
			AssertFileCount(3)
	})
}
