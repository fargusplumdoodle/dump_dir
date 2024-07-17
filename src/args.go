package src

import (
	"os"
	"path/filepath"
	"strings"
)

func ValidateArgs(args []string) bool {
	if len(args) < 2 {
		return false
	}
	return true
}

func ParseArgs(args []string) Config {
	if len(args) < 1 {
		PrintUsage()
		return Config{}
	}

	config := Config{
		Extensions: strings.Split(args[0], ","),
	}

	skipMode := false
	for _, arg := range args[1:] {
		if arg == "--include-ignored" {
			config.IncludeIgnored = true
			continue
		}
		if arg == "-s" {
			skipMode = true
			continue
		}
		if skipMode {
			config.SkipDirs = append(config.SkipDirs, arg)
			skipMode = false
		} else {
			if fileInfo, err := os.Stat(arg); err == nil {
				if fileInfo.IsDir() {
					config.Directories = append(config.Directories, arg)
				} else {
					config.SpecificFiles = append(config.SpecificFiles, filepath.Clean(arg))
				}
			} else {
				config.Directories = append(config.Directories, arg) // Assume it's a directory if we can't stat it
			}
		}
	}

	return config
}
