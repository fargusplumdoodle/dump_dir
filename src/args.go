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
		Extensions:    []string{}, // Empty slice means all extensions
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

	skipMode := false
	extensionMode := false

	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch arg {
		case "--help", "-h":
			config.Action = "help"
			return config, nil
		case "--include-ignored":
			config.IncludeIgnored = true
		case "-s", "--skip":
			skipMode = true
		case "-e", "--extension":
			extensionMode = true
		case "--max-filesize", "-m":
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
		default:
			if skipMode {
				config.AddSkipDir(arg)
				skipMode = false
			} else if extensionMode {
				config.Extensions = append(config.Extensions, strings.Split(arg, ",")...)
				extensionMode = false
			} else {
				if err := config.AddIncludePath(arg); err != nil {
					fmt.Printf("Warning: Could not process path %s: %v\n", arg, err)
				}
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
