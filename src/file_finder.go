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
	im, err := NewIgnoreManager(config.IncludeIgnored, config.SkipDirs)
	if err != nil {
		fmt.Printf(boldRed("❌ Error initializing IgnoreManager: %v\n"), err)
	}
	return &FileFinder{Config: config, Fs: fs, IgnoreManager: im}
}

func (ff *FileFinder) DiscoverFiles() []string {
	allDirs := ff.findAllDirectories()
	filesToProcess := ff.findMatchingFiles(allDirs)
	return append(filesToProcess, ff.Config.SpecificFiles...)
}

func (ff *FileFinder) findAllDirectories() []string {
	var allDirs []string
	for _, dir := range ff.Config.Directories {
		ff.walkDirectory(dir, &allDirs)
	}
	return allDirs
}

func (ff *FileFinder) walkDirectory(dir string, allDirs *[]string) {
	afero.Walk(ff.Fs, filepath.Clean(dir), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			PrintError("accessing path", path, err)
			return nil
		}
		if !info.IsDir() {
			return nil
		}
		if ff.shouldSkipDirectory(path) {
			return filepath.SkipDir
		}
		*allDirs = append(*allDirs, path)
		return nil
	})
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

func (ff *FileFinder) findMatchingFiles(dirs []string) []string {
	var matchingFiles []string
	for _, dir := range dirs {
		files, err := afero.ReadDir(ff.Fs, filepath.Clean(dir))
		if err != nil {
			PrintError("reading directory", dir, err)
			continue
		}
		for _, file := range files {
			if !file.IsDir() {
				ff.processFile(dir, file, &matchingFiles)
			}
		}
	}
	return matchingFiles
}

func (ff *FileFinder) processFile(dir string, file os.FileInfo, matchingFiles *[]string) {
	filePath := filepath.Join(dir, file.Name())
	if ff.shouldProcessFile(filePath) {
		*matchingFiles = append(*matchingFiles, filePath)
	}
}

func (ff *FileFinder) shouldProcessFile(filePath string) bool {
	return !ff.IgnoreManager.ShouldIgnore(filePath) &&
		ff.matchesExtensions(filepath.Base(filePath))
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
