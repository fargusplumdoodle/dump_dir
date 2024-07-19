package src

import (
	"fmt"
	"github.com/atotto/clipboard"
	"strings"
)

func PrintUsage() {
	fmt.Println()
	fmt.Println(boldCyan("Usage:"))
	fmt.Println("  dump_dir [options] <file_extension1>[,<file_extension2>,...] <directory1> [directory2] ... [-s <skip_directory1>] [-s <skip_directory2>] ... [--include-ignored]")
	fmt.Println("  Use 'any' as file_extension to match all files")
	fmt.Println()
	fmt.Println(boldCyan("Options:"))
	fmt.Println("  -h, --help     Display this help information")
	fmt.Println("  -v, --version  Display the version of dump_dir")
	fmt.Println()
	fmt.Println(BoldGreen("Examples:"))
	fmt.Println("  dump_dir --help")
	fmt.Println("  dump_dir --version")
	fmt.Println("  dump_dir js ./project -s ./project/node_modules -s ./project/dist")
	fmt.Println("  dump_dir any ./project -s ./project/node_modules")
	fmt.Println("  dump_dir go,js,py ./project")
	fmt.Println("  dump_dir tsx ./README.md  # Will only print out the README.md file if it exists")
	fmt.Println("  dump_dir any ./project --include-ignored  # Include all files, even those normally ignored")
	fmt.Println()
	fmt.Println(boldMagenta("Description:"))
	fmt.Println("  This will search for files with the specified extensions (or all files if 'any' is used)")
	fmt.Println("  in the given directories, excluding any specified directories.")
	fmt.Println("  Multiple file extensions can be specified by separating them with commas.")
	fmt.Println("  Use --include-ignored to include files that would normally be ignored (e.g., those in .gitignore).")
	fmt.Println()
}

func PrintError(errorType string, filePath string, err error) {
	fmt.Printf(boldRed("❌ Error %s file %s: %v\n", errorType, filePath, err))
}

func CopyToClipboard(content string) bool {
	err := clipboard.WriteAll(content)
	if err != nil {
		fmt.Println(boldRed(fmt.Sprintf("❌ Error copying to clipboard: %v", err)))
		return false
	}
	return true
}

func FormatFileContent(path, contents string) string {
	return fmt.Sprintf("START FILE: %s\n%s\nEND FILE: %s\n\n", path, contents, path)
}

func GenerateDetailedOutput(stats Stats) string {
	var detailedOutput strings.Builder

	for _, fileInfo := range stats.ProcessedFiles {
		detailedOutput.WriteString(FormatFileContent(fileInfo.Path, fileInfo.Contents))
	}

	return detailedOutput.String()
}

func PrintDetailedOutput(stats Stats) {
	detailedOutput := GenerateDetailedOutput(stats)
	summary := DisplayStats(stats)

	if CopyToClipboard(detailedOutput) {
		summary += BoldGreen("✅ File contents have been copied to clipboard.\n")
	}

	fmt.Println(summary)
}
