package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run script.go <file_extension> <directory1> [directory2] ...")
		os.Exit(1)
	}

	extension := os.Args[1]
	directories := os.Args[2:]

	var matchingFiles []string
	var totalLines int
	var detailedOutput strings.Builder

	for _, dir := range directories {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), "."+extension) {
				matchingFiles = append(matchingFiles, path)
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()

				detailedOutput.WriteString(fmt.Sprintf("File: %s\n", path))

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					detailedOutput.WriteString(scanner.Text() + "\n")
					totalLines++
				}
				detailedOutput.WriteString("\n")

				if err := scanner.Err(); err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			fmt.Printf("Error walking directory %s: %v\n", dir, err)
		}
	}

	// Generate summary
	summary := fmt.Sprintf("\n")
	summary += fmt.Sprintf("Matching files:\n")
	for _, file := range matchingFiles {
		summary += fmt.Sprintf("- %s\n", file)
	}
	summary += fmt.Sprintf("Total files found: %d\n", len(matchingFiles))
	summary += fmt.Sprintf("Total lines across all files: %d\n", totalLines)

	// Copy detailed output to clipboard
	err := clipboard.WriteAll(detailedOutput.String())
	if err != nil {
		fmt.Println("Error copying to clipboard:", err)
	} else {
		summary += "Detailed output has been copied to clipboard.\n"
	}

	// Print summary to terminal
	fmt.Println(summary)
}
