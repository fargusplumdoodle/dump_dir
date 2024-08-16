package src

import (
	"bufio"
	"fmt"
	"github.com/spf13/afero"
	"io"
	"strings"
	"sync"
)

const FilesPerGoroutine = 1

type FileProcessor struct {
	Fs     afero.Fs
	Config Config
}

func NewFileProcessor(fs afero.Fs, config Config) *FileProcessor {
	return &FileProcessor{Fs: fs, Config: config}
}

func (fp *FileProcessor) ProcessFiles(files []string) []FileInfo {
	var wg sync.WaitGroup
	fileInfoChan := make(chan FileInfo, len(files))

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

	go func() {
		wg.Wait()
		close(fileInfoChan)
	}()

	var processedFiles []FileInfo
	for fileInfo := range fileInfoChan {
		processedFiles = append(processedFiles, fileInfo)
	}

	return processedFiles
}

func (fp *FileProcessor) processFile(path string) (FileInfo, error) {
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
		return FileInfo{Path: path, Contents: "<EMPTY FILE>", Status: StatusParsed}, nil
	}

	if info.Size() > fp.Config.MaxFileSize {
		return FileInfo{Status: StatusSkippedTooLarge, Path: path, Contents: fmt.Sprintf("<FILE TOO LARGE: %d bytes>", info.Size())}, nil
	}

	isBinary, err := fp.fileIsBinary(file)
	if err != nil {
		return FileInfo{}, fmt.Errorf("checking if file is binary: %w", err)
	}
	if isBinary {
		return FileInfo{Status: StatusSkippedBinary, Path: path, Contents: "<BINARY SKIPPED>"}, nil
	}

	var contents strings.Builder
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024*1024), int(fp.Config.MaxFileSize))

	for scanner.Scan() {
		contents.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		if err == bufio.ErrTooLong {
			return FileInfo{Path: path, Contents: "<FILE EXCEEDS BUFFER SIZE>"}, nil
		}
		return FileInfo{}, fmt.Errorf("scanning file: %w", err)
	}

	return FileInfo{Status: StatusParsed, Path: path, Contents: contents.String()}, nil
}

func (fp *FileProcessor) fileIsBinary(file afero.File) (bool, error) {
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
