package src

import (
	"bufio"
	"fmt"
	"github.com/gobwas/glob"
	"os"
	"path/filepath"
	"strings"
)

func ProcessDirectories(extension string, directories, skipDirs []string, ignorePatterns []glob.Glob, includeGitIgnored bool) ([]string, string, int) {
	var matchingFiles []string
	var detailedOutput strings.Builder
	var totalLines int

	for _, dir := range directories {
		files, output, lines := processDirectory(dir, extension, skipDirs, ignorePatterns, includeGitIgnored)
		matchingFiles = append(matchingFiles, files...)
		detailedOutput.WriteString(output)
		totalLines += lines
	}

	return matchingFiles, detailedOutput.String(), totalLines
}

func processDirectory(dir, extension string, skipDirs []string, ignorePatterns []glob.Glob, includeGitIgnored bool) ([]string, string, int) {
	var matchingFiles []string
	var detailedOutput strings.Builder
	var totalLines int

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignore .git directory unless --include-gitignored-paths is used
		if !includeGitIgnored && info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		if info.IsDir() {
			for _, skipDir := range skipDirs {
				if strings.HasPrefix(path, skipDir) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		relPath, _ := filepath.Rel(dir, path)
		if !includeGitIgnored {
			for _, pattern := range ignorePatterns {
				if pattern.Match(relPath) {
					return nil
				}
			}
		}

		if extension == "any" || strings.HasSuffix(info.Name(), "."+extension) {
			fileInfo, err := processFile(path)
			if err != nil {
				return err
			}
			matchingFiles = append(matchingFiles, fileInfo.Path)
			detailedOutput.WriteString(FormatFileContent(fileInfo.Path, fileInfo.Contents))
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

	// Read the first 512 bytes to check if it's a binary file
	buffer := make([]byte, 512)
	bytesRead, err := file.Read(buffer)
	if err != nil {
		return FileInfo{}, err
	}

	// Reset the file pointer to the beginning
	_, err = file.Seek(0, 0)
	if err != nil {
		return FileInfo{}, err
	}

	// Check if the file is binary
	if isBinary(buffer[:bytesRead]) {
		return FileInfo{Path: path, Contents: "<BINARY SKIPPED>"}, nil
	}

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

func isBinary(buffer []byte) bool {
	const maxCheck = 1024 // Maximum number of bytes to check
	if len(buffer) > maxCheck {
		buffer = buffer[:maxCheck]
	}

	controlChars := 0
	for _, b := range buffer {
		if b == 0 {
			return true // Null byte, definitely binary
		}
		if b < 7 || (b > 14 && b < 32) {
			controlChars++
		}
	}

	// If more than 30% non-UTF8 control characters, assume binary
	return float64(controlChars)/float64(len(buffer)) > 0.3
}
