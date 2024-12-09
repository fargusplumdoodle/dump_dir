package tests

import (
	"fmt"
	"github.com/fargusplumdoodle/dump_dir/tests/e2e"
	"testing"
)

func TestSkipDirectory(t *testing.T) {
	t.Run("skip single directory", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./main.go":          "package main\nfunc main() {}\n",
				"./lib/helper.go":    "package lib\nfunc helper() {}\n",
				"./build/output.txt": "build output\n",
				"./src/feature.go":   "package src\nfunc feature() {}\n",
			}).
			WithArgs("-e go --skip ./lib .")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./main.go").
			AssertFileInOutput("./src/feature.go").
			AssertFileCount(2)
	})

	t.Run("skip multiple directories", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./main.go":          "package main\nfunc main() {}\n",
				"./lib/helper.go":    "package lib\nfunc helper() {}\n",
				"./build/output.txt": "build output\n",
				"./src/feature.go":   "package src\nfunc feature() {}\n",
				"./test/test.go":     "package test\nfunc test() {}\n",
			}).
			WithArgs("-e go --skip ./lib -s ./test .")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./main.go").
			AssertFileInOutput("./src/feature.go").
			AssertFileCount(2)
	})

	t.Run("skip nested directories", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./main.go":                 "package main\nfunc main() {}\n",
				"./lib/helper.go":           "package lib\nfunc helper() {}\n",
				"./lib/internal/util.go":    "package internal\nfunc util() {}\n",
				"./lib/external/wrapper.go": "package external\nfunc wrapper() {}\n",
			}).
			WithArgs("-e go --skip ./lib .")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./main.go").
			AssertFileCount(1)
	})

	t.Run("skip files", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./main.go":                 "package main\nfunc main() {}\n",
				"./lib/helper.go":           "package lib\nfunc helper() {}\n",
				"./lib/internal/util.go":    "package internal\nfunc util() {}\n",
				"./lib/external/wrapper.go": "package external\nfunc wrapper() {}\n",
			}).
			WithArgs("-s lib/helper.go .")

		result := env.Run()

		fmt.Println(result.Output)

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./main.go").
			AssertFileInOutput("./lib/internal/util.go").
			AssertFileInOutput("./lib/external/wrapper.go").
			AssertFileCount(3)
	})

	t.Run("Normalize skip file paths", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./src/main.go":       "package main\nfunc main() {}\n",
				"./src/lib/helper.go": "package lib\nfunc helper() {}\n",
				"./src/lib/util.go":   "package internal\nfunc util() {}\n",
				"./src/ay/ay.go":      "package internal\nfunc util() {}\n",
			}).
			WithArgs(". -s src/lib/helper.go -s ./src/lib/util.go --skip src/ay")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./src/main.go").
			AssertFileCount(1)
	})

	t.Run("skip directory with special characters", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./main.go":            "package main\nfunc main() {}\n",
				"./lib-v1.2/helper.go": "package lib\nfunc helper() {}\n",
				"./lib@temp/util.go":   "package temp\nfunc util() {}\n",
				"./src/feature.go":     "package src\nfunc feature() {}\n",
			}).
			WithArgs("-e go --skip ./lib-v1.2 --skip ./lib@temp .")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./main.go").
			AssertFileInOutput("./src/feature.go").
			AssertFileCount(2)
	})

	t.Run("skip nonexistent directory", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithFiles(map[string]string{
				"./main.go":        "package main\nfunc main() {}\n",
				"./lib/helper.go":  "package lib\nfunc helper() {}\n",
				"./src/feature.go": "package src\nfunc feature() {}\n",
			}).
			WithArgs("-e go --skip ./nonexistent .")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("./main.go").
			AssertFileInOutput("./lib/helper.go").
			AssertFileInOutput("./src/feature.go").
			AssertFileCount(3)
	})
}
