package src

import (
	"fmt"
	"sort"
	"strings"
)

func CalculateStats(processedFiles []FileInfo) Stats {
	var totalLines, estimatedTokens int
	var skippedLarge, skippedBinary, parsedFiles []FileInfo

	sortedFiles := SortFileList(processedFiles)

	for _, fileInfo := range sortedFiles {
		switch fileInfo.Status {
		case StatusParsed:
			totalLines += strings.Count(fileInfo.Contents, "\n")
			estimatedTokens += estimateTokens(fileInfo.Contents)
			parsedFiles = append(parsedFiles, fileInfo)
		case StatusSkippedTooLarge:
			skippedLarge = append(skippedLarge, fileInfo)
		case StatusSkippedBinary:
			skippedBinary = append(skippedBinary, fileInfo)
		}
	}

	return Stats{
		TotalFiles:      len(processedFiles),
		TotalLines:      totalLines,
		EstimatedTokens: estimatedTokens,
		ProcessedFiles:  sortedFiles,
		ParsedFiles:     parsedFiles,
		SkippedLarge:    skippedLarge,
		SkippedBinary:   skippedBinary,
	}
}

func DisplayStats(stats Stats) string {
	var summary strings.Builder
	printFileList(&summary, "🔍 Parsed files:", stats.ParsedFiles)
	printFileList(&summary, "🪨 Skipped large files:", stats.SkippedLarge)
	printFileList(&summary, "💽 Skipped binary files:", stats.SkippedBinary)

	summary.WriteString(boldCyan(fmt.Sprintf("\n📚 Total files found: %d\n", stats.TotalFiles)))
	summary.WriteString(boldCyan(fmt.Sprintf("📝 Total lines across all parsed files: %d\n", stats.TotalLines)))
	summary.WriteString(boldCyan(fmt.Sprintf("🔢 Estimated tokens: %s\n\n", formatTokenCount(stats.EstimatedTokens))))

	return summary.String()
}

func printFileList(summary *strings.Builder, heading string, files []FileInfo) {
	if len(files) == 0 {
		return
	}
	summary.WriteString(boldMagenta(fmt.Sprintf("\n%s\n", heading)))

	for _, file := range files {
		summary.WriteString(fmt.Sprintf("- %s\n", file.Path))
	}
	summary.WriteString("\n")
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

func SortFileList(files []FileInfo) []FileInfo {
	sort.Slice(files, func(i, j int) bool {
		// Split the paths into components
		pathI := strings.Split(files[i].Path, "/")
		pathJ := strings.Split(files[j].Path, "/")

		// Compare each component
		for k := 0; k < len(pathI) && k < len(pathJ); k++ {
			if pathI[k] != pathJ[k] {
				return pathI[k] < pathJ[k]
			}
		}

		// If all components are the same up to this point, shorter path comes first
		return len(pathI) < len(pathJ)
	})

	return files
}
