package src

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const FilesPerGoroutine = 1
const MaxFileSize = 10 * 1024 * 1024 // 10 MB

type FileProcessor struct {
	Config        Config
	IgnoreManager *IgnoreManager
}

func NewFileProcessor(config Config) *FileProcessor {
	fp := &FileProcessor{Config: config}
	err := UpdateFileProcessor(fp)
	if err != nil {
		fmt.Printf(boldRed("❌ Error initializing IgnoreManager: %v\n"), err)
	}
	return fp
}

func (fp *FileProcessor) ProcessDirectories() []FileInfo {
	// Step 1: Find all directories and subdirectories
	allDirs := fp.findAllDirectories()

	// Step 2: Find all matching files in subdirectories
	filesToProcess := fp.findMatchingFiles(allDirs)

	// Add specifically mentioned files
	filesToProcess = append(filesToProcess, fp.Config.SpecificFiles...)

	// Step 3: Process all found files in parallel
	return fp.processFilesParallel(filesToProcess)
}

func (fp *FileProcessor) processFilesParallel(files []string) []FileInfo {
	var wg sync.WaitGroup
	fileInfoChan := make(chan FileInfo, len(files))

	// Process files in chunks
	for i := 0; i < len(files); i += FilesPerGoroutine {
		end := i + FilesPerGoroutine
		if end > len(files) {
			end = len(files)
		}

		wg.Add(1)
		go func(chunk []string) {
			defer wg.Done()
			for _, file := range chunk {
				fileInfo, err := processFile(file)
				if err != nil {
					PrintError("processing", file, err)
					continue
				}
				fileInfoChan <- fileInfo
			}
		}(files[i:end])
	}

	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(fileInfoChan)
	}()

	// Collect results
	var processedFiles []FileInfo
	for fileInfo := range fileInfoChan {
		processedFiles = append(processedFiles, fileInfo)
	}

	return processedFiles
}

func (fp *FileProcessor) findAllDirectories() []string {
	var allDirs []string

	for _, dir := range fp.Config.Directories {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				PrintError("accessing path", path, err)
				return nil
			}

			if info.IsDir() {
				if !fp.Config.IncludeIgnored && fp.IgnoreManager.ShouldIgnore(path) {
					fmt.Printf("Skipping ignored directory: %s\n", path)
					return filepath.SkipDir
				}
				for _, skipDir := range fp.Config.SkipDirs {
					if strings.HasPrefix(path, skipDir) {
						fmt.Printf("Skipping directory: %s\n", path)
						return filepath.SkipDir
					}
				}
				allDirs = append(allDirs, path)
			}
			return nil
		})

		if err != nil {
			fmt.Printf(boldRed("❌ Error walking directory %s: %v\n"), dir, err)
		}
	}

	return allDirs
}

func (fp *FileProcessor) findMatchingFiles(dirs []string) []string {
	var matchingFiles []string

	for _, dir := range dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			PrintError("reading directory", dir, err)
			continue
		}

		for _, file := range files {
			if !file.IsDir() {
				filePath := filepath.Join(dir, file.Name())
				if fp.Config.IncludeIgnored || !fp.IgnoreManager.ShouldIgnore(filePath) {
					if fp.matchesExtensions(file.Name()) {
						matchingFiles = append(matchingFiles, filePath)
					}
				}
			}
		}
	}

	return matchingFiles
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
		return FileInfo{}, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	// Check if the file is empty
	info, err := file.Stat()
	if err != nil {
		return FileInfo{}, fmt.Errorf("getting file info: %w", err)
	}
	if info.Size() == 0 {
		return FileInfo{Path: path, Contents: "<EMPTY FILE>"}, nil
	}

	// Check if the file size exceeds the maximum allowed size
	if info.Size() > MaxFileSize {
		return FileInfo{Path: path, Contents: fmt.Sprintf("<FILE TOO LARGE: %d bytes>", info.Size())}, nil
	}

	// Check if the file is binary
	isBinary, err := fileIsBinary(file)
	if err != nil {
		return FileInfo{}, fmt.Errorf("checking if file is binary: %w", err)
	}
	if isBinary {
		return FileInfo{Path: path, Contents: "<BINARY SKIPPED>"}, nil
	}

	var contents strings.Builder
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024*1024), MaxFileSize)

	for scanner.Scan() {
		contents.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		if err == bufio.ErrTooLong {
			return FileInfo{Path: path, Contents: "<FILE EXCEEDS BUFFER SIZE>"}, nil
		}
		return FileInfo{}, fmt.Errorf("scanning file: %w", err)
	}

	return FileInfo{Path: path, Contents: contents.String()}, nil
}

func fileIsBinary(file *os.File) (bool, error) {
	buffer := make([]byte, 512)
	bytesRead, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false, fmt.Errorf("reading file: %w", err)
	}
	buffer = buffer[:bytesRead]

	// Reset the file pointer to the beginning
	_, err = file.Seek(0, 0)
	if err != nil {
		return false, fmt.Errorf("seeking file: %w", err)
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
