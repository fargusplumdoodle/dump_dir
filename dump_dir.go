package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
)

type FileInfo struct {
	Path     string
	Contents string
}

func main() {
	if !validateArgs() {
		return
	}

	extension, directories, skipDirs := parseArgs(os.Args[1:])

	matchingFiles, detailedOutput, totalLines := processDirectories(extension, directories, skipDirs)
	summary := generateSummary(matchingFiles, totalLines)

	if copyToClipboard(detailedOutput) {
		summary += "Detailed output has been copied to clipboard.\n"
	}

	fmt.Println(summary)
}

func validateArgs() bool {
	if len(os.Args) < 3 {
		fmt.Println("Usage: dump_dir <file_extension> <directory1> [directory2] ... [-s <skip_directory1>] [-s <skip_directory2>] ...")
		fmt.Println("Example: dump_dir js ./project -s ./project/node_modules -s ./project/dist")
		fmt.Println("This will search for all .js files in ./project, excluding the node_modules and dist directories.")
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
			detailedOutput.WriteString(fmt.Sprintf("File: %s\n%s\n\n", fileInfo.Path, fileInfo.Contents))
			totalLines += strings.Count(fileInfo.Contents, "\n")
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory %s: %v\n", dir, err)
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
	summary.WriteString("\nMatching files:\n")
	for _, file := range matchingFiles {
		summary.WriteString(fmt.Sprintf("- %s\n", file))
	}
	summary.WriteString(fmt.Sprintf("Total files found: %d\n", len(matchingFiles)))
	summary.WriteString(fmt.Sprintf("Total lines across all files: %d\n", totalLines))
	return summary.String()
}

func copyToClipboard(content string) bool {
	err := clipboard.WriteAll(content)
	if err != nil {
		fmt.Println("Error copying to clipboard:", err)
		return false
	}
	return true
}
