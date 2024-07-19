package src

import (
	"os"
	"strings"
)

var OsStat = os.Stat

func ValidateArgs(args []string) bool {
	return len(args) > 0
}

func ParseArgs(args []string) Config {
	config := Config{
		Action:        "dump_dir", // Default action
		SkipDirs:      []string{},
		SpecificFiles: []string{},
		Directories:   []string{},
		Extensions:    []string{},
	}

	if len(args) == 0 {
		return config
	}

	// Check for version flag
	if args[0] == "--version" || args[0] == "-v" {
		config.Action = "version"
		return config
	}

	// Check for help flag anywhere in the arguments
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			config.Action = "help"
			return config
		}
	}

	// If we're here, it's the default dump_dir action
	config.Extensions = strings.Split(args[0], ",")

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
			if fileInfo, err := OsStat(arg); err == nil {
				if fileInfo.IsDir() {
					config.Directories = append(config.Directories, arg)
				} else {
					config.SpecificFiles = append(config.SpecificFiles, arg)
				}
			} else {
				// If we can't stat it, assume it's a directory
				config.Directories = append(config.Directories, arg)
			}
		}
	}

	return config
}
