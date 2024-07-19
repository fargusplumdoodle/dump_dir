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
	fp := &FileFinder{Config: config, Fs: fs}
	err := UpdateFileProcessor(fp)
	if err != nil {
		fmt.Printf(boldRed("❌ Error initializing IgnoreManager: %v\n"), err)
	}
	return fp
}

func (ff *FileFinder) DiscoverFiles() []string {
	allDirs := ff.findAllDirectories()
	filesToProcess := ff.findMatchingFiles(allDirs)
	filesToProcess = append(filesToProcess, ff.Config.SpecificFiles...)
	return filesToProcess
}

func (ff *FileFinder) findAllDirectories() []string {
	var allDirs []string

	for _, dir := range ff.Config.Directories {
		afero.Walk(ff.Fs, dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				PrintError("accessing path", path, err)
				return nil
			}

			if info.IsDir() {
				if !ff.Config.IncludeIgnored && ff.IgnoreManager.ShouldIgnore(path) {
					fmt.Printf("Skipping ignored directory: %s\n", path)
					return filepath.SkipDir
				}
				for _, skipDir := range ff.Config.SkipDirs {
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

func (ff *FileFinder) findMatchingFiles(dirs []string) []string {
	var matchingFiles []string

	for _, dir := range dirs {
		files, err := afero.ReadDir(ff.Fs, dir)
		if err != nil {
			PrintError("reading directory", dir, err)
			continue
		}

		for _, file := range files {
			if !file.IsDir() {
				filePath := filepath.Join(dir, file.Name())
				if ff.Config.IncludeIgnored || !ff.IgnoreManager.ShouldIgnore(filePath) {
					if ff.matchesExtensions(file.Name()) {
						matchingFiles = append(matchingFiles, filePath)
					}
				}
			}
		}
	}

	return matchingFiles
}

func (ff *FileFinder) matchesExtensions(filename string) bool {
	if len(ff.Config.Extensions) == 1 && ff.Config.Extensions[0] == "any" {
		return true
	}
	for _, ext := range ff.Config.Extensions {
		if strings.HasSuffix(filename, "."+ext) {
			return true
		}
	}
	return false
}
