package src

import (
	"fmt"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"strings"
)

type FileFinder struct {
	Config        Config
	IgnoreManager *IgnoreManager
	Fs            afero.Fs
}

func NewFileFinder(config Config, fs afero.Fs) *FileFinder {
	im, err := NewIgnoreManager(fs, config.IncludeIgnored, config.SkipDirs)
	if err != nil {
		fmt.Printf(boldRed("❌ Error initializing IgnoreManager: %v\n"), err)
	}
	return &FileFinder{Config: config, Fs: fs, IgnoreManager: im}
}

func (ff *FileFinder) DiscoverFiles() []string {
	// Use a map to track unique files
	uniqueFiles := make(map[string]bool)

	// Process directories
	for _, dir := range ff.Config.Directories {
		files := ff.findMatchingFilesInDir(dir)
		for _, file := range files {
			uniqueFiles[file] = true
		}
	}

	// Add specific files if they match criteria
	for _, file := range ff.Config.SpecificFiles {
		if ff.shouldProcessFile(file) {
			uniqueFiles[file] = true
		}
	}

	// Convert map keys to slice
	result := make([]string, 0, len(uniqueFiles))
	for file := range uniqueFiles {
		result = append(result, file)
	}
	return result
}

func (ff *FileFinder) findMatchingFilesInDir(rootDir string) []string {
	var matchingFiles []string

	afero.Walk(ff.Fs, filepath.Clean(rootDir), func(path string, info os.FileInfo, err error) error {
		path = NormalizePath(path)
		if err != nil {
			PrintError("accessing path", path, err)
			return nil
		}

		if info.IsDir() {
			if ff.shouldSkipDirectory(path) {
				return filepath.SkipDir
			}
			return nil
		}

		// Process file if it matches criteria
		if ff.shouldProcessFile(path) {
			matchingFiles = append(matchingFiles, NormalizePath(path))
		}
		return nil
	})

	return matchingFiles
}

func (ff *FileFinder) shouldSkipDirectory(path string) bool {
	if ff.IgnoreManager.ShouldIgnore(path) {
		fmt.Printf("Skipping ignored directory: %s\n", path)
		return true
	}
	for _, skipDir := range ff.Config.SkipDirs {
		if ff.isSubdirectory(path, skipDir) {
			fmt.Printf("Skipping directory: %s\n", path)
			return true
		}
	}
	return false
}

func (ff *FileFinder) isSubdirectory(path, parentDir string) bool {
	normalizedParent := filepath.Clean(parentDir)
	return path == normalizedParent || strings.HasPrefix(path, normalizedParent+string(os.PathSeparator))
}

func (ff *FileFinder) shouldProcessFile(filePath string) bool {
	if ff.IgnoreManager.ShouldIgnore(filePath) {
		return false
	}

	// Check glob patterns first if they exist
	if len(ff.Config.GlobPatterns) > 0 {
		filename := filepath.Base(filePath)
		for _, pattern := range ff.Config.GlobPatterns {
			matched, err := filepath.Match(pattern, filename)
			if err != nil {
				fmt.Printf(boldRed("❌ Error matching glob pattern %s: %v\n"), pattern, err)
				continue
			}
			if matched {
				return true
			}
		}
		return false
	}

	// If no glob patterns, fall back to extension matching
	return ff.matchesExtensions(filepath.Base(filePath))
}

func (ff *FileFinder) matchesExtensions(filename string) bool {
	if len(ff.Config.Extensions) == 0 {
		return true
	}
	for _, ext := range ff.Config.Extensions {
		if strings.HasSuffix(filename, "."+ext) {
			return true
		}
	}
	return false
}
