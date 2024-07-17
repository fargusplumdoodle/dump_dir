package src

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FileProcessor struct {
	Config Config
}

func NewFileProcessor(config Config) *FileProcessor {
	return &FileProcessor{Config: config}
}

func (fp *FileProcessor) ProcessDirectories() ([]string, string, int) {
	var matchingFiles []string
	var detailedOutput strings.Builder
	var totalLines int

	// Process specific files first
	for _, file := range fp.Config.SpecificFiles {
		fileInfo, err := processFile(file)
		if err != nil {
			fmt.Printf(boldRed("❌ Error processing file %s: %v\n"), file, err)
			continue
		}
		matchingFiles = append(matchingFiles, fileInfo.Path)
		detailedOutput.WriteString(FormatFileContent(fileInfo.Path, fileInfo.Contents))
		totalLines += strings.Count(fileInfo.Contents, "\n")
	}

	// Process directories
	for _, dir := range fp.Config.Directories {
		files, output, lines := fp.processDirectory(dir)
		matchingFiles = append(matchingFiles, files...)
		detailedOutput.WriteString(output)
		totalLines += lines
	}

	return matchingFiles, detailedOutput.String(), totalLines
}

func (fp *FileProcessor) processDirectory(dir string) ([]string, string, int) {
	var matchingFiles []string
	var detailedOutput strings.Builder
	var totalLines int

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf(boldRed("❌ Error accessing path %s: %v\n"), path, err)
			return nil // Continue walking despite the error
		}

		if info.IsDir() {
			for _, skipDir := range fp.Config.SkipDirs {
				if strings.HasPrefix(path, skipDir) {
					fmt.Printf("Skipping directory: %s\n", path)
					return filepath.SkipDir
				}
			}
			return nil
		}

		// Check if the file is in the specificFiles list
		for _, specificFile := range fp.Config.SpecificFiles {
			if path == specificFile {
				return nil // Skip processing here as it's already been processed
			}
		}

		if fp.matchesExtensions(info.Name()) {
			fileInfo, err := processFile(path)
			if err != nil {
				fmt.Printf(boldRed("❌ Error processing file %s: %v\n"), path, err)
				return nil // Continue walking despite the error
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

func (fp *FileProcessor) matchesExtensions(filename string) bool {
	if len(fp.Config.Extensions) == 1 && fp.Config.Extensions[0] == "any" {
		return true
	}
	for _, ext := range fp.Config.Extensions {
		if strings.HasSuffix(filename, "."+ext) {
			return true
		}
	}
	return false
}

func processFile(path string) (FileInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return FileInfo{}, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Check if the file is empty
	if info, err := file.Stat(); err != nil {
		return FileInfo{}, fmt.Errorf("error getting file info: %w", err)
	} else if info.Size() == 0 {
		return FileInfo{Path: path, Contents: "<EMPTY FILE>"}, nil
	}

	// Read the first 512 bytes to check if it's a binary file

	// Check if the file is binary
	isBinary, err := fileIsBinary(file)
	if err != nil {
		return FileInfo{}, fmt.Errorf("error checking if file is binary: %w", err)

	}
	if isBinary {
		return FileInfo{Path: path, Contents: "<BINARY SKIPPED>"}, nil
	}

	var contents strings.Builder
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024*1024), 10*1024*1024) // Increase buffer size

	for scanner.Scan() {
		contents.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		return FileInfo{}, fmt.Errorf("error scanning file: %w", err)
	}

	return FileInfo{Path: path, Contents: contents.String()}, nil
}

func fileIsBinary(file *os.File) (bool, error) {
	buffer := make([]byte, 512)
	bytesRead, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false, fmt.Errorf("error reading file: %w", err)
	}
	buffer = buffer[:bytesRead]

	// Reset the file pointer to the beginning
	_, err = file.Seek(0, 0)
	if err != nil {
		return false, fmt.Errorf("error seeking file: %w", err)
	}
	const maxCheck = 1024 // Maximum number of bytes to check
	if len(buffer) > maxCheck {
		buffer = buffer[:maxCheck]
	}

	controlChars := 0
	for _, b := range buffer {
		if b == 0 {
			return true, nil // Null byte, definitely binary
		}
		if b < 7 || (b > 14 && b < 32) {
			controlChars++
		}
	}

	// If more than 30% non-UTF8 control characters, assume binary
	return float64(controlChars)/float64(len(buffer)) > 0.3, nil
}
