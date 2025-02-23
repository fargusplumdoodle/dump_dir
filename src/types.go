package src

import (
	"fmt"
	"github.com/spf13/afero"
)

type FileStatus string

const (
	StatusParsed          FileStatus = "PARSED"
	StatusSkippedBinary   FileStatus = "SKIPPED_BINARY"
	StatusSkippedTooLarge FileStatus = "SKIPPED_TOO_LARGE"
)

type FileInfo struct {
	Path     string
	Contents string
	Status   FileStatus
}

type Config struct {
	Action         string
	Extensions     []string
	Directories    []string
	SkipDirs       []string
	SpecificFiles  []string
	IncludeIgnored bool
	MaxFileSize    int64
	GlobPatterns   []string
	NoConfig       bool
}

func (c *Config) AddSkipDir(path string) {
	if path == "" {
		return
	}

	normalizedPath := NormalizePath(path)

	for _, existingPath := range c.SkipDirs {
		if existingPath == normalizedPath {
			return
		}
	}

	c.SkipDirs = append(c.SkipDirs, normalizedPath)
	return
}

func (c *Config) AddIncludePath(path string) error {
	if path == "" {
		return nil
	}
	normalizedPath := NormalizePath(path)

	// Check if path is in skip directories
	for _, skipDir := range c.SkipDirs {
		if normalizedPath == skipDir {
			return nil
		}
	}

	// Check if path already exists in either list
	for _, existingPath := range c.SpecificFiles {
		if existingPath == normalizedPath {
			return nil
		}
	}
	for _, existingPath := range c.Directories {
		if existingPath == normalizedPath {
			return nil
		}
	}

	// Determine if it's a directory or file
	isDir, err := isDirectory(normalizedPath)
	if err != nil {
		return fmt.Errorf("error checking path type for %s: %w", normalizedPath, err)
	}

	// Add to appropriate list
	if isDir {
		c.Directories = append(c.Directories, normalizedPath)
	} else {
		c.SpecificFiles = append(c.SpecificFiles, normalizedPath)
	}

	return nil
}

type RunConfig struct {
	Fs        afero.Fs
	Clipboard ClipboardManager
	Version   string
	Commit    string
	Date      string
}

type Stats struct {
	TotalFiles      int
	TotalLines      int
	EstimatedTokens int
	ProcessedFiles  []FileInfo
	ParsedFiles     []FileInfo
	SkippedLarge    []FileInfo
	SkippedBinary   []FileInfo
}
