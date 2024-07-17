package src

import (
	"fmt"
	"github.com/atotto/clipboard"
	"strings"
)

func GenerateSummary(processedFiles []FileInfo, totalLines int, estimatedTokens int) string {
	var summary strings.Builder
	summary.WriteString(boldMagenta("\n🔍 Matching files:\n"))
	for _, file := range processedFiles {
		summary.WriteString(fmt.Sprintf("  - %s\n", file.Path))
	}
	summary.WriteString(boldCyan(fmt.Sprintf("\n📚 Total files found: %d\n", len(processedFiles))))
	summary.WriteString(boldCyan(fmt.Sprintf("📝 Total lines across all files: %d\n", totalLines)))
	summary.WriteString(boldCyan(fmt.Sprintf("🔢 Estimated tokens: %s\n\n", formatTokenCount(estimatedTokens))))

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

func GenerateDetailedOutput(processedFiles []FileInfo) (string, int, int) {
	var detailedOutput strings.Builder
	var totalLines int
	var estimatedTokens int

	for _, fileInfo := range processedFiles {
		detailedOutput.WriteString(FormatFileContent(fileInfo.Path, fileInfo.Contents))
		totalLines += strings.Count(fileInfo.Contents, "\n")
		estimatedTokens += estimateTokens(fileInfo.Contents)
	}

	return detailedOutput.String(), totalLines, estimatedTokens
}

func PrintDetailedOutput(processedFiles []FileInfo) {
	detailedOutput, totalLines, estimatedTokens := GenerateDetailedOutput(processedFiles)
	summary := GenerateSummary(processedFiles, totalLines, estimatedTokens)

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
	fmt.Println("  dump_dir <file_extension1>[,<file_extension2>,...] <directory1> [directory2] ... [-s <skip_directory1>] [-s <skip_directory2>] ... [--include-ignored]")
	fmt.Println("  Use 'any' as file_extension to match all files")
	fmt.Println()
	fmt.Println(BoldGreen("Examples:"))
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

func estimateTokens(content string) int {
	return len(strings.Fields(content))
}

func formatTokenCount(tokens int) string {
	if tokens < 100 {
		return fmt.Sprintf("%d", tokens)
	} else if tokens < 1000 {
		return fmt.Sprintf("%.1fk", float64(tokens)/1000)
	} else {
		return fmt.Sprintf("%dk", tokens/1000)
	}
}
