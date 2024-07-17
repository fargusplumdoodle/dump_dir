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
			detailedOutput.WriteString(fmt.Sprintf("START FILE: %s\n%s\n\nEND FILE: %s\n\n", fileInfo.Path, fileInfo.Contents, fileInfo.Path))
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
