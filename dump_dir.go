package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/fatih/color"
)

type FileInfo struct {
	Path     string
	Contents string
}

var (
	boldGreen   = color.New(color.FgGreen, color.Bold).SprintfFunc()
	boldCyan    = color.New(color.FgCyan, color.Bold).SprintfFunc()
	boldMagenta = color.New(color.FgMagenta, color.Bold).SprintfFunc()
	boldRed     = color.New(color.FgRed, color.Bold).SprintfFunc()
)

func main() {
	if !validateArgs() {
		return
	}

	extension, directories, skipDirs := parseArgs(os.Args[1:])

	matchingFiles, detailedOutput, totalLines := processDirectories(extension, directories, skipDirs)
	summary := generateSummary(matchingFiles, totalLines)

	if copyToClipboard(detailedOutput) {
		summary += boldGreen("✅ File contents have been copied to clipboard.\n")
	}

	fmt.Println(summary)
}

func validateArgs() bool {
	if len(os.Args) < 3 {
		fmt.Println()
		fmt.Println(boldRed("❌ Error: Insufficient arguments"))
		fmt.Println()
		fmt.Println(boldCyan("Usage:"))
		fmt.Println("  dump_dir <file_extension> <directory1> [directory2] ... [-s <skip_directory1>] [-s <skip_directory2>] ...")
		fmt.Println()
		fmt.Println(boldGreen("Example:"))
		fmt.Println("  dump_dir js ./project -s ./project/node_modules -s ./project/dist")
		fmt.Println()
		fmt.Println(boldMagenta("Description:"))
		fmt.Println("  This will search for all .js files in ./project, excluding the node_modules and dist directories.")
		fmt.Println()
		return false
	}
	return true
}

func parseArgs(args []string) (string, []string, []string) {
	extension := args[0]
	var directories []string
	var skipDirs []string
	skipMode := false

	for _, arg := range args[1:] {
		if arg == "-s" {
			skipMode = true
			continue
		}
		if skipMode {
			skipDirs = append(skipDirs, arg)
			skipMode = false
		} else {
			directories = append(directories, arg)
		}
	}

	return extension, directories, skipDirs
}

func processDirectories(extension string, directories, skipDirs []string) ([]string, string, int) {
	var matchingFiles []string
	var detailedOutput strings.Builder
	var totalLines int

	for _, dir := range directories {
		files, output, lines := processDirectory(dir, extension, skipDirs)
		matchingFiles = append(matchingFiles, files...)
		detailedOutput.WriteString(output)
		totalLines += lines
	}

	return matchingFiles, detailedOutput.String(), totalLines
}

func processDirectory(dir, extension string, skipDirs []string) ([]string, string, int) {
	var matchingFiles []string
	var detailedOutput strings.Builder
	var totalLines int

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			for _, skipDir := range skipDirs {
				if strings.HasPrefix(path, skipDir) {
					return filepath.SkipDir
				}
			}
			return nil
		}
		if strings.HasSuffix(info.Name(), "."+extension) {
			fileInfo, err := processFile(path)
			if err != nil {
				return err
			}
			matchingFiles = append(matchingFiles, fileInfo.Path)
			detailedOutput.WriteString(fmt.Sprintf("📄 File: %s\n%s\n\n", boldCyan(fileInfo.Path), fileInfo.Contents))
			totalLines += strings.Count(fileInfo.Contents, "\n")
		}
		return nil
	})

	if err != nil {
		fmt.Printf(boldRed("❌ Error walking directory %s: %v\n"), dir, err)
	}

	return matchingFiles, detailedOutput.String(), totalLines
}

func processFile(path string) (FileInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return FileInfo{}, err
	}
	defer file.Close()

	var contents strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		contents.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		return FileInfo{}, err
	}

	return FileInfo{Path: path, Contents: contents.String()}, nil
}

func generateSummary(matchingFiles []string, totalLines int) string {
	var summary strings.Builder
	summary.WriteString(boldMagenta("🔍 Matching files:\n"))
	for _, file := range matchingFiles {
		summary.WriteString(fmt.Sprintf("  - %s\n", file))
	}
	summary.WriteString(boldCyan(fmt.Sprintf("📚 Total files found: %d\n", len(matchingFiles))))
	summary.WriteString(boldCyan(fmt.Sprintf("📝 Total lines across all files: %d\n\n", totalLines)))
	return summary.String()
}

func copyToClipboard(content string) bool {
	err := clipboard.WriteAll(content)
	if err != nil {
		fmt.Println(boldRed(fmt.Sprintf("❌ Error copying to clipboard: %v", err)))
		return false
	}
	return true
}
