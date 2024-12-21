package tests

import (
	"github.com/fargusplumdoodle/dump_dir/tests/e2e"
	"testing"
)

func TestConfigurationFile(t *testing.T) {
	t.Run("valid configuration file", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				".dump_dir.yml": `---
include:
  - ./src
  - ./README.md
ignore:
  - ./vendor
  - ./dist`,
				"./src/main.go":   "package main\nfunc main() {}\n",
				"./src/util.go":   "package main\nfunc util() {}\n",
				"./vendor/lib.go": "package vendor\nfunc lib() {}\n",
				"./dist/app.go":   "package dist\nfunc app() {}\n",
				"./README.md":     "# Test Project\n",
			}).
			WithArgs(".")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./.dump_dir.yml").
			AssertFileInOutput("./src/main.go").
			AssertFileInOutput("./src/util.go").
			AssertFileInOutput("./README.md").
			AssertFileCount(4)
	})

	t.Run("invalid YAML syntax in config", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				".dump_dir.yml": `---
include:
  - ./src
  [ invalid yaml
ignore:
  - ./vendor`,
				"./src/main.go":   "package main\nfunc main() {}\n",
				"./vendor/lib.go": "package vendor\nfunc lib() {}\n",
			}).
			WithArgs(".")

		result := env.Run()

		if result.Err == nil {
			t.Errorf("Expected an error!")
		}
	})

	t.Run("config file directory inclusion rules", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				".dump_dir.yml": `---
include:
  - ./docs
  - ./src/core`,
				"./docs/guide.md":             "# User Guide\n",
				"./docs/api.md":               "# API Reference\n",
				"./src/core/main.go":          "package core\nfunc main() {}\n",
				"./src/plugins/plugin.go":     "package plugins\nfunc plugin() {}\n",
				"./src/core/internal/util.go": "package internal\nfunc util() {}\n",
			}).
			WithArgs(".")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./.dump_dir.yml").
			AssertFileInOutput("./docs/guide.md").
			AssertFileInOutput("./docs/api.md").
			AssertFileInOutput("./src/core/main.go").
			AssertFileInOutput("./src/plugins/plugin.go").
			AssertFileInOutput("./src/core/internal/util.go").
			AssertFileCount(6)
	})

	t.Run("config file directory exclusion rules", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				".dump_dir.yml": `---
ignore:
  - ./test
  - ./build`,
				"./src/main.go":          "package main\nfunc main() {}\n",
				"./test/test.go":         "package test\nfunc test() {}\n",
				"./build/output.go":      "package build\nfunc output() {}\n",
				"./src/internal/util.go": "package internal\nfunc util() {}\n",
				"./src/public/public.go": "package public\nfunc public() {}\n",
			}).
			WithArgs(".")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./.dump_dir.yml").
			AssertFileInOutput("./src/main.go").
			AssertFileInOutput("./src/public/public.go").
			AssertFileInOutput("./src/internal/util.go").
			AssertFileCount(4)
	})

	t.Run("CLI skip flag overrides config include", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				".dump_dir.yml": `---
include:
  - ./src
  - ./docs`,
				"./src/main.go":   "package main\nfunc main() {}\n",
				"./docs/guide.md": "# User Guide\n",
			}).
			WithArgs(". --skip ./docs") // Skip flag should override config include

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./.dump_dir.yml").
			AssertFileInOutput("./src/main.go").
			AssertFileCount(2)
	})

	t.Run("extension flag overrides config file", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				".dump_dir.yml": `---
include:
  - ./src
  - ./docs`,
				"./src/main.go":   "package main\nfunc main() {}\n",
				"./src/util.py":   "def util():\n    pass\n",
				"./docs/guide.md": "# User Guide\n",
			}).
			WithArgs(". -e go") // Extension flag should filter regardless of config includes

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./src/main.go").
			AssertFileCount(1) // Only includes .go files
	})

	t.Run("ignore config file when --no-config is set", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				".dump_dir.yml": `---
ignore:
  - ./vendor
  - ./src
  - ./README.md`,
				"./src/main.go":   "package main\nfunc main() {}\n",
				"./src/util.go":   "package main\nfunc util() {}\n",
				"./vendor/lib.go": "package vendor\nfunc lib() {}\n",
				"./README.md":     "# Test Project\n",
				"./extra.txt":     "Extra file content",
			}).
			WithArgs(". --no-config") // Use the new flag

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./.dump_dir.yml").
			AssertFileInOutput("./src/main.go").
			AssertFileInOutput("./src/util.go").
			AssertFileInOutput("./vendor/lib.go").
			AssertFileInOutput("./README.md").
			AssertFileInOutput("./extra.txt").
			AssertFileCount(6)
	})
}
