package src

import (
	"os"
)

func ValidateArgs() bool {
	if len(os.Args) < 3 {
		PrintUsage()
		return false
	}
	return true
}

func ParseArgs(args []string) (string, []string, []string, bool) {
	extension := args[0]
	var directories []string
	var skipDirs []string
	skipMode := false
	includeGitIgnored := false

	for _, arg := range args[1:] {
		if arg == "-s" {
			skipMode = true
			continue
		}
		if arg == "--include-gitignored-paths" {
			includeGitIgnored = true
			continue
		}
		if skipMode {
			skipDirs = append(skipDirs, arg)
			skipMode = false
		} else {
			directories = append(directories, arg)
		}
	}

	return extension, directories, skipDirs, includeGitIgnored
}
