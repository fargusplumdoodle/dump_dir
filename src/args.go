package src

import (
	"os"
	"path/filepath"
)

func ValidateArgs() bool {
	if len(os.Args) < 2 {
		PrintUsage()
		return false
	}
	return true
}

func ParseArgs(args []string) (string, []string, []string, []string) {
	extension := args[0]
	var directories []string
	var skipDirs []string
	var specificFiles []string
	skipMode := false

	for _, arg := range args[1:] {
		if arg == "-s" {
			skipMode = true
			continue
		}
		if skipMode {
			skipDirs = append(skipDirs, arg)
			skipMode = false
		} else {
			if fileInfo, err := os.Stat(arg); err == nil {
				if fileInfo.IsDir() {
					directories = append(directories, arg)
				} else {
					specificFiles = append(specificFiles, filepath.Clean(arg))
				}
			} else {
				directories = append(directories, arg) // Assume it's a directory if we can't stat it
			}
		}
	}

	return extension, directories, skipDirs, specificFiles
}
