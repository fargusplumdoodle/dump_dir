package src

import (
	"fmt"
	"github.com/atotto/clipboard"
	"strings"
)

func GenerateSummary(matchingFiles []string, totalLines int) string {
	var summary strings.Builder
	summary.WriteString(boldMagenta("🔍 Matching files:\n"))
	for _, file := range matchingFiles {
		summary.WriteString(fmt.Sprintf("  - %s\n", file))
	}
	summary.WriteString(boldCyan(fmt.Sprintf("📚 Total files found: %d\n", len(matchingFiles))))
	summary.WriteString(boldCyan(fmt.Sprintf("📝 Total lines across all files: %d\n\n", totalLines)))
	return summary.String()
}

func CopyToClipboard(content string) bool {
	err := clipboard.WriteAll(content)
	if err != nil {
		fmt.Println(boldRed(fmt.Sprintf("❌ Error copying to clipboard: %v", err)))
		return false
	}
	return true
}

func PrintDetailedOutput(matchingFiles []string, detailedOutput string, totalLines int) {
	summary := GenerateSummary(matchingFiles, totalLines)

	if CopyToClipboard(detailedOutput) {
		summary += BoldGreen("✅ File contents have been copied to clipboard.\n")
	}

	fmt.Println(summary)
}

func FormatFileContent(path, contents string) string {
	return fmt.Sprintf("START FILE: %s\n%s\nEND FILE: %s\n\n", path, contents, path)
}

func PrintUsage() {
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
}
