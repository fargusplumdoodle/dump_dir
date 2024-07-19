package src

import (
	"bufio"
	"fmt"
	"github.com/spf13/afero"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const FilesPerGoroutine = 1
const MaxFileSize = 10 * 1024 * 10 // 10 MB

type FileFinder struct {
	Config        Config
	IgnoreManager *IgnoreManager
	Fs            afero.Fs
}

func NewFileFinder(config Config, fs afero.Fs) *FileFinder {
	fp := &FileFinder{Config: config, Fs: fs}
	err := UpdateFileProcessor(fp)
	if err != nil {
		fmt.Printf(boldRed("❌ Error initializing IgnoreManager: %v\n"), err)
	}
	return fp
}

func (fp *FileFinder) DiscoverFiles() []FileInfo {
	allDirs := fp.findAllDirectories()
	filesToProcess := fp.findMatchingFiles(allDirs)
	filesToProcess = append(filesToProcess, fp.Config.SpecificFiles...)
	return fp.processFilesParallel(filesToProcess)
}

func (fp *FileFinder) processFilesParallel(files []string) []FileInfo {
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
				fileInfo, err := fp.processFile(file)
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

func (fp *FileFinder) findAllDirectories() []string {
	var allDirs []string

	for _, dir := range fp.Config.Directories {
		afero.Walk(fp.Fs, dir, func(path string, info os.FileInfo, err error) error {
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
	}

	return allDirs
}

func (fp *FileFinder) findMatchingFiles(dirs []string) []string {
	var matchingFiles []string

	for _, dir := range dirs {
		files, err := afero.ReadDir(fp.Fs, dir)
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

func (fp *FileFinder) matchesExtensions(filename string) bool {
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

func (fp *FileFinder) processFile(path string) (FileInfo, error) {
	file, err := fp.Fs.Open(path)
	if err != nil {
		return FileInfo{}, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return FileInfo{}, fmt.Errorf("getting file info: %w", err)
	}
	if info.Size() == 0 {
		return FileInfo{Path: path, Contents: "<EMPTY FILE>"}, nil
	}

	if info.Size() > MaxFileSize {
		return FileInfo{Path: path, Contents: fmt.Sprintf("<FILE TOO LARGE: %d bytes>", info.Size())}, nil
	}

	isBinary, err := fp.fileIsBinary(file)
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

func (fp *FileFinder) fileIsBinary(file afero.File) (bool, error) {
	buffer := make([]byte, 512)
	bytesRead, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false, fmt.Errorf("reading file: %w", err)
	}
	buffer = buffer[:bytesRead]

	_, err = file.Seek(0, 0)
	if err != nil {
		return false, fmt.Errorf("seeking file: %w", err)
	}
	const maxCheck = 1024
	if len(buffer) > maxCheck {
		buffer = buffer[:maxCheck]
	}

	controlChars := 0
	for _, b := range buffer {
		if b == 0 {
			return true, nil
		}
		if b < 7 || (b > 14 && b < 32) {
			controlChars++
		}
	}

	return float64(controlChars)/float64(len(buffer)) > 0.3, nil
}
