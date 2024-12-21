package src

import (
	"path/filepath"
	"strings"
)

func isDirectory(path string) (bool, error) {
	fileInfo, err := OsStat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

func NormalizePath(path string) string {
	if path == "" {
		return ""
	}

	// Clean the path first to handle any ../.. or // cases
	cleaned := filepath.Clean(path)

	// Don't modify paths that:
	// 1. Start with ../
	// 2. Are absolute paths
	// 3. Are empty or "."
	if strings.HasPrefix(cleaned, "..") || filepath.IsAbs(cleaned) || cleaned == "." || cleaned == "" {
		return cleaned
	}

	// If path already starts with ./, return it as is
	if strings.HasPrefix(cleaned, "./") {
		return cleaned
	}

	// Add ./ prefix
	return "./" + cleaned
}
