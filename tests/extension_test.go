package tests

import (
	"github.com/fargusplumdoodle/dump_dir/tests/e2e"
	"testing"
)

func TestExtensionFiltering(t *testing.T) {
	t.Run("multiple extensions", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithWorkingDir("/test/project").
			WithFiles(map[string]string{
				"/test/project/main.go":    "package main\n\nfunc main() {}\n",
				"/test/project/script.py":  "def main():\n    pass\n",
				"/test/project/index.js":   "console.log('hello');\n",
				"/test/project/README.md":  "# Test Project\n",
				"/test/project/styles.css": "body { margin: 0; }\n",
			}).
			WithArgs("-e go,py,js /test/project")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("/test/project/main.go").
			AssertFileInOutput("/test/project/script.py").
			AssertFileInOutput("/test/project/index.js").
			AssertFileCount(3)
	})

	t.Run("nonexistent extension", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithWorkingDir("/test/project").
			WithFiles(map[string]string{
				"/test/project/main.go":  "package main\n\nfunc main() {}\n",
				"/test/project/test.txt": "some text\n",
			}).
			WithArgs("-e xyz /test/project")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileCount(0)
	})

	t.Run("mixed valid and invalid extensions", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithWorkingDir("/test/project").
			WithFiles(map[string]string{
				"/test/project/main.go":  "package main\n\nfunc main() {}\n",
				"/test/project/test.txt": "some text\n",
				"/test/project/util.py":  "def util():\n    pass\n",
			}).
			WithArgs("-e go,xyz,py /test/project")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("/test/project/main.go").
			AssertFileInOutput("/test/project/util.py").
			AssertFileCount(2)
	})

	t.Run("no extension specified", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithWorkingDir("/test/project").
			WithFiles(map[string]string{
				"/test/project/main.go":  "package main\n\nfunc main() {}\n",
				"/test/project/test.txt": "some text\n",
				"/test/project/README":   "# Test Project\n",
			}).
			WithArgs("/test/project")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("/test/project/main.go").
			AssertFileInOutput("/test/project/test.txt").
			AssertFileInOutput("/test/project/README").
			AssertFileCount(3)
	})
}
