package tests

import (
	"github.com/fargusplumdoodle/dump_dir/tests/e2e"
	"testing"
)

func TestGlobOptions(t *testing.T) {
	t.Run("Include files matching single glob pattern", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./main.go":       "package main\nfunc main() {}\n",
				"./helper.py":     "def helper():\n    pass\n",
				"./README.md":     "# Project README\n",
				"./utils/util.js": "console.log('util');\n",
			}).
			WithArgs("-g *.go .")

		result := env.Run()
		result.PrintOutput()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./main.go").
			AssertFileNotInOutput("./helper.py").
			AssertFileNotInOutput("./README.md").
			AssertFileNotInOutput("./utils/util.js").
			AssertFileCount(1)
	})

	t.Run("Include files matching multiple glob patterns", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./main.go":        "package main\nfunc main() {}\n",
				"./helper.py":      "def helper():\n    pass\n",
				"./README.md":      "# Project README\n",
				"./utils/util.js":  "console.log('util');\n",
				"./scripts/run.sh": "#!/bin/bash\necho 'Run'\n",
			}).
			WithArgs("-g *.go -g *.py .")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./main.go").
			AssertFileInOutput("./helper.py").
			AssertFileNotInOutput("./README.md").
			AssertFileNotInOutput("./utils/util.js").
			AssertFileNotInOutput("./scripts/run.sh").
			AssertFileCount(2)
	})

	t.Run("Glob pattern with no matching files", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./main.go":       "package main\nfunc main() {}\n",
				"./helper.py":     "def helper():\n    pass\n",
				"./README.md":     "# Project README\n",
				"./utils/util.js": "console.log('util');\n",
			}).
			WithArgs("-g *.rb .")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileCount(0)
	})

	t.Run("Combine glob with skip to ignore matched files", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./main.go":           "package main\nfunc main() {}\n",
				"./helper.go":         "package main\nfunc helper() {}\n",
				"./utils/util.go":     "package utils\nfunc Util() {}\n",
				"./utils/helper.go":   "package utils\nfunc Helper() {}\n",
				"./README.md":         "# Project README\n",
				"./scripts/deploy.go": "package scripts\nfunc Deploy() {}\n",
			}).
			WithArgs("-g *.go --skip ./utils .")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./main.go").
			AssertFileInOutput("./helper.go").
			AssertFileInOutput("./scripts/deploy.go").
			AssertFileNotInOutput("./utils/util.go").
			AssertFileNotInOutput("./utils/helper.go").
			AssertFileCount(3)
	})

	t.Run("Glob pattern with nested directories and exclusions", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./src/main.go":            "package main\nfunc main() {}\n",
				"./src/utils/util.go":      "package utils\nfunc Util() {}\n",
				"./src/utils/test_util.go": "package utils\nfunc TestUtil() {}\n",
				"./src/helpers/helper.go":  "package helpers\nfunc Helper() {}\n",
				"./tests/test_main.go":     "package tests\nfunc TestMain() {}\n",
				"./docs/guide.md":          "# Guide\n",
				"./scripts/deploy.go":      "package scripts\nfunc Deploy() {}\n",
			}).
			WithArgs("-g *.go --skip ./src/utils .")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./src/main.go").
			AssertFileInOutput("./src/helpers/helper.go").
			AssertFileInOutput("./scripts/deploy.go").
			AssertFileInOutput("./tests/test_main.go").
			AssertFileNotInOutput("./src/utils/util.go").
			AssertFileNotInOutput("./src/utils/test_util.go").
			AssertFileCount(4)
	})
}
