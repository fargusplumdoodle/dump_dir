package src

import (
	"fmt"
	"github.com/atotto/clipboard"
	"strings"
)

func PrintUsage() {
	fmt.Println()
	fmt.Println(boldCyan("Usage:"))
	fmt.Println("  dump_dir [options] <file_extension1> [,<file_extension2>,...] <directory1> [directory2] [...other options]")
	fmt.Println("  Use 'any' as file_extension to match all files")
	fmt.Println()
	fmt.Println(boldCyan("Options:"))
	fmt.Println("  -h, --help                 Display this help information")
	fmt.Println("  -v, --version              Display the version of dump_dir")
	fmt.Println("  -s <directory>             Skip specified directory")
	fmt.Println("  --include-ignored          Include files that would normally be ignored (e.g., those in .gitignore)")
	fmt.Println("  -m <size>, --max-filesize <size>  Specify the maximum file size to process. You can use units like B, KB, or MB (e.g., 500KB, 2MB).")
	fmt.Println("                             If no unit is specified, it defaults to bytes.")
	fmt.Println("  -g, --glob <pattern>       Only include files matching the glob pattern (e.g., '*.txt', 'test_*.go')")
	fmt.Println()
	fmt.Println(BoldGreen("Examples:"))
	fmt.Println("  dump_dir js ./project -s ./project/node_modules -s ./project/dist")
	fmt.Println("  dump_dir any ./project")
	fmt.Println("  dump_dir go,js,py ./project")
	fmt.Println("  dump_dir any ./README.md ./main.go")
	fmt.Println("  dump_dir any ./project --include-ignored")
	fmt.Println("  dump_dir go ./project --max-filesize 1MB")
	fmt.Println("  dump_dir any ./project -g '*.txt' -g 'test_*.go'")
	fmt.Println("  dump_dir any ./project --glob '*.md'")
	fmt.Println()
	fmt.Println(boldMagenta("Description:"))
	fmt.Println("  This will search for files with the specified extensions (or all files if 'any' is used)")
	fmt.Println("  in the given directories, excluding any specified directories.")
	fmt.Println("  Multiple file extensions can be specified by separating them with commas.")
	fmt.Println("  Use --include-ignored to include files that would normally be ignored (e.g., those in .gitignore).")
	fmt.Println("  The tool respects .gitignore rules by default and ignores common version control directories.")
	fmt.Println("  You can set a maximum file size to process using the --max-filesize option.")
	fmt.Println()
	fmt.Println("  More documentation at: https://github.com/fargusplumdoodle/dump_dir")
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
