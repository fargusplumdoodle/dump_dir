package src

import (
	"fmt"
	"os"
)

func ValidateArgs() bool {
	if len(os.Args) < 3 {
		fmt.Println()
		fmt.Println(boldRed("❌ Error: Insufficient arguments"))
		fmt.Println()
		fmt.Println(boldCyan("Usage:"))
		fmt.Println("  dump_dir <file_extension> <directory1> [directory2] ... [-s <skip_directory1>] [-s <skip_directory2>] ... [--include-gitignored-paths]")
		fmt.Println("  Use 'any' as file_extension to match all files")
		fmt.Println()
		fmt.Println(BoldGreen("Example:"))
		fmt.Println("  dump_dir js ./project -s ./project/node_modules -s ./project/dist")
		fmt.Println("  dump_dir any ./project -s ./project/node_modules --include-gitignored-paths")
		fmt.Println()
		fmt.Println(boldMagenta("Description:"))
		fmt.Println("  This will search for files with the specified extension (or all files if 'any' is used)")
		fmt.Println("  in the given directories, excluding any specified directories and .gitignore'd files by default.")
		fmt.Println("  Use --include-gitignored-paths to include files that would be ignored by .gitignore.")
		fmt.Println()
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
