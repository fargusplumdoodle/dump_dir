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
