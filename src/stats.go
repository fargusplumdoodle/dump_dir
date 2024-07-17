package src

import (
	"fmt"
	"strings"
)

func CalculateStats(processedFiles []FileInfo) Stats {
	var totalLines, estimatedTokens int

	for _, fileInfo := range processedFiles {
		totalLines += strings.Count(fileInfo.Contents, "\n")
		estimatedTokens += estimateTokens(fileInfo.Contents)
	}

	return Stats{
		TotalFiles:      len(processedFiles),
		TotalLines:      totalLines,
		EstimatedTokens: estimatedTokens,
		ProcessedFiles:  processedFiles,
	}
}

func DisplayStats(stats Stats) string {
	var summary strings.Builder
	summary.WriteString(boldMagenta("\n🔍 Matching files:\n"))
	for _, file := range stats.ProcessedFiles {
		summary.WriteString(fmt.Sprintf("  - %s\n", file.Path))
	}
	summary.WriteString(boldCyan(fmt.Sprintf("\n📚 Total files found: %d\n", stats.TotalFiles)))
	summary.WriteString(boldCyan(fmt.Sprintf("📝 Total lines across all files: %d\n", stats.TotalLines)))
	summary.WriteString(boldCyan(fmt.Sprintf("🔢 Estimated tokens: %s\n\n", formatTokenCount(stats.EstimatedTokens))))

	return summary.String()
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
