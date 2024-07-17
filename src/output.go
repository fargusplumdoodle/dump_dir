package src

import (
	"fmt"
	"github.com/atotto/clipboard"
	"strings"
)

func GenerateSummary(processedFiles []FileInfo, totalLines int) string {
	var summary strings.Builder
	summary.WriteString(boldMagenta("\n🔍 Matching files:\n"))
	for _, file := range processedFiles {
		summary.WriteString(fmt.Sprintf("  - %s\n", file.Path))
	}
	summary.WriteString(boldCyan(fmt.Sprintf("\n📚 Total files found: %d\n", len(processedFiles))))
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

func FormatFileContent(path, contents string) string {
	return fmt.Sprintf("START FILE: %s\n%s\nEND FILE: %s\n\n", path, contents, path)
}

func GenerateDetailedOutput(processedFiles []FileInfo) (string, int) {
	var detailedOutput strings.Builder
	var totalLines int

	for _, fileInfo := range processedFiles {
		detailedOutput.WriteString(FormatFileContent(fileInfo.Path, fileInfo.Contents))
		totalLines += strings.Count(fileInfo.Contents, "\n")
	}

	return detailedOutput.String(), totalLines
}

func PrintDetailedOutput(processedFiles []FileInfo) {
	detailedOutput, totalLines := GenerateDetailedOutput(processedFiles)
	summary := GenerateSummary(processedFiles, totalLines)

	if CopyToClipboard(detailedOutput) {
		summary += BoldGreen("✅ File contents have been copied to clipboard.\n")
	}

	fmt.Println(summary)
}

func PrintUsage() {
	fmt.Println()
	fmt.Println(boldRed("❌ Error: Insufficient arguments"))
	fmt.Println()
	fmt.Println(boldCyan("Usage:"))
	fmt.Println("  dump_dir <file_extension1>[,<file_extension2>,...] <directory1> [directory2] ... [-s <skip_directory1>] [-s <skip_directory2>] ...")
	fmt.Println("  Use 'any' as file_extension to match all files")
	fmt.Println()
	fmt.Println(BoldGreen("Examples:"))
	fmt.Println("  dump_dir js ./project -s ./project/node_modules -s ./project/dist")
	fmt.Println("  dump_dir any ./project -s ./project/node_modules")
	fmt.Println("  dump_dir go,js,py ./project")
	fmt.Println("  dump_dir tsx ./README.md  # Will only print out the README.md file if it exists")
	fmt.Println()
	fmt.Println(boldMagenta("Description:"))
	fmt.Println("  This will search for files with the specified extensions (or all files if 'any' is used)")
	fmt.Println("  in the given directories, excluding any specified directories.")
	fmt.Println("  Multiple file extensions can be specified by separating them with commas.")
	fmt.Println()
}
