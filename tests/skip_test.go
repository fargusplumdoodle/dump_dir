package tests

import (
	"fmt"
	"github.com/fargusplumdoodle/dump_dir/tests/e2e"
	"testing"
)

func TestSkipDirectory(t *testing.T) {
	t.Run("skip single directory", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithWorkingDir("/test/project").
			WithFiles(map[string]string{
				"/test/project/main.go":          "package main\nfunc main() {}\n",
				"/test/project/lib/helper.go":    "package lib\nfunc helper() {}\n",
				"/test/project/build/output.txt": "build output\n",
				"/test/project/src/feature.go":   "package src\nfunc feature() {}\n",
			}).
			WithArgs("-e go --skip /test/project/lib /test/project")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("/test/project/main.go").
			AssertFileInOutput("/test/project/src/feature.go").
			AssertFileCount(2)
	})

	t.Run("skip multiple directories", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithWorkingDir("/test/project").
			WithFiles(map[string]string{
				"/test/project/main.go":          "package main\nfunc main() {}\n",
				"/test/project/lib/helper.go":    "package lib\nfunc helper() {}\n",
				"/test/project/build/output.txt": "build output\n",
				"/test/project/src/feature.go":   "package src\nfunc feature() {}\n",
				"/test/project/test/test.go":     "package test\nfunc test() {}\n",
			}).
			WithArgs("-e go --skip /test/project/lib -s /test/project/test /test/project")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("/test/project/main.go").
			AssertFileInOutput("/test/project/src/feature.go").
			AssertFileCount(2)
	})

	t.Run("skip nested directories", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithWorkingDir("/test/project").
			WithFiles(map[string]string{
				"/test/project/main.go":                 "package main\nfunc main() {}\n",
				"/test/project/lib/helper.go":           "package lib\nfunc helper() {}\n",
				"/test/project/lib/internal/util.go":    "package internal\nfunc util() {}\n",
				"/test/project/lib/external/wrapper.go": "package external\nfunc wrapper() {}\n",
			}).
			WithArgs("-e go --skip /test/project/lib /test/project")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("/test/project/main.go").
			AssertFileCount(1)
	})

	t.Run("skip files", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithWorkingDir("/test/project").
			WithFiles(map[string]string{
				"/test/project/main.go":                 "package main\nfunc main() {}\n",
				"/test/project/lib/helper.go":           "package lib\nfunc helper() {}\n",
				"/test/project/lib/internal/util.go":    "package internal\nfunc util() {}\n",
				"/test/project/lib/external/wrapper.go": "package external\nfunc wrapper() {}\n",
			}).
			WithArgs("-s /test/project/lib/helper.go /test/project")

		result := env.Run()

		fmt.Println(result.Output)

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("/test/project/main.go").
			AssertFileInOutput("/test/project/lib/internal/util.go").
			AssertFileInOutput("/test/project/lib/external/wrapper.go").
			AssertFileCount(3)
	})

	t.Run("skip directory with special characters", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithWorkingDir("/test/project").
			WithFiles(map[string]string{
				"/test/project/main.go":            "package main\nfunc main() {}\n",
				"/test/project/lib-v1.2/helper.go": "package lib\nfunc helper() {}\n",
				"/test/project/lib@temp/util.go":   "package temp\nfunc util() {}\n",
				"/test/project/src/feature.go":     "package src\nfunc feature() {}\n",
			}).
			WithArgs("-e go --skip /test/project/lib-v1.2 --skip /test/project/lib@temp /test/project")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("/test/project/main.go").
			AssertFileInOutput("/test/project/src/feature.go").
			AssertFileCount(2)
	})

	t.Run("skip nonexistent directory", func(t *testing.T) {
		env := e2e.NewEnvironment(t).
			WithWorkingDir("/test/project").
			WithFiles(map[string]string{
				"/test/project/main.go":        "package main\nfunc main() {}\n",
				"/test/project/lib/helper.go":  "package lib\nfunc helper() {}\n",
				"/test/project/src/feature.go": "package src\nfunc feature() {}\n",
			}).
			WithArgs("-e go --skip /test/project/nonexistent /test/project")

		result := env.Run()

		validator := e2e.NewOutputValidator(t, result)
		validator.
			AssertSuccessfulRun().
			AssertFileInOutput("/test/project/main.go").
			AssertFileInOutput("/test/project/lib/helper.go").
			AssertFileInOutput("/test/project/src/feature.go").
			AssertFileCount(3)
	})
}
