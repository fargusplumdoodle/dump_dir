package src

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var OsStat = os.Stat

// ErrInvalidMaxFileSize is a custom error type for invalid max filesize arguments
type ErrInvalidMaxFileSize struct {
	Value string
}

func (e ErrInvalidMaxFileSize) Error() string {
	return fmt.Sprintf("invalid max filesize: %s", e.Value)
}

func ValidateArgs(args []string) bool {
	return len(args) > 0
}

func ParseArgs(args []string) (Config, error) {
	config := Config{
		Action:        "dump_dir", // Default action
		SkipDirs:      []string{},
		SpecificFiles: []string{},
		Directories:   []string{},
		Extensions:    []string{},
		MaxFileSize:   500 * 1024, // Default to 500KB
	}

	if len(args) == 0 {
		return config, nil
	}

	// Check for version flag
	if args[0] == "--version" || args[0] == "-v" {
		config.Action = "version"
		return config, nil
	}

	// Check for help flag anywhere in the arguments
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			config.Action = "help"
			return config, nil
		}
	}

	// If we're here, it's the default dump_dir action
	config.Extensions = strings.Split(args[0], ",")

	skipMode := false
	for i := 1; i < len(args); i++ {
		arg := args[i]
		if arg == "--include-ignored" {
			config.IncludeIgnored = true
		} else if arg == "-s" || arg == "--skip" {
			skipMode = true
		} else if arg == "--max-filesize" || arg == "-m" {
			if i+1 < len(args) {
				size, err := parseFileSize(args[i+1])
				if err != nil {
					return config, ErrInvalidMaxFileSize{Value: args[i+1]}
				}
				config.MaxFileSize = size
				i++ // Skip the next argument as we've processed it
			} else {
				return config, ErrInvalidMaxFileSize{Value: ""}
			}
		} else if skipMode {
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

	return config, nil
}

func parseFileSize(sizeStr string) (int64, error) {
	sizeStr = strings.ToUpper(sizeStr)
	var multiplier int64 = 1

	if strings.HasSuffix(sizeStr, "KB") {
		multiplier = 1024
		sizeStr = strings.TrimSuffix(sizeStr, "KB")
	} else if strings.HasSuffix(sizeStr, "MB") {
		multiplier = 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "MB")
	} else if strings.HasSuffix(sizeStr, "B") {
		sizeStr = strings.TrimSuffix(sizeStr, "B")
	}

	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return size * multiplier, nil
}
