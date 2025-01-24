package src

import (
	"bufio"
	"github.com/spf13/afero"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
)

var ExecCommand = exec.Command

type IgnoreManager struct {
	fs afero.Fs

	ignorePatterns []glob.Glob
	ignoreDirs     []string
	skipPaths      []string
	includeIgnored bool
}

func NewIgnoreManager(fs afero.Fs, includeIgnored bool, skipPaths []string) (*IgnoreManager, error) {
	im := &IgnoreManager{
		fs:             fs,
		includeIgnored: includeIgnored,
		skipPaths:      skipPaths,
	}
	err := im.loadIgnorePatterns()
	if err != nil {
		return nil, err
	}
	return im, nil
}

func (im *IgnoreManager) loadIgnorePatterns() error {
	// Add .git to ignore patterns
	im.ignoreDirs = append(im.ignoreDirs, ".git")

	// Load global gitignore
	globalGitignorePath, err := getGlobalGitignorePath()
	if err == nil && globalGitignorePath != "" {
		im.loadIgnoreFile(globalGitignorePath)
	}

	// Load project-specific gitignore
	im.loadIgnoreFile(".gitignore")

	return nil
}

func getGlobalGitignorePath() (string, error) {
	cmd := ExecCommand("git", "config", "--global", "--get", "core.excludesfile")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func (im *IgnoreManager) loadIgnoreFile(path string) {
	file, err := im.fs.Open(path)
	if err != nil {
		return // Ignore errors, as the file might not exist
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pattern := strings.TrimSpace(scanner.Text())
		if pattern != "" && !strings.HasPrefix(pattern, "#") {
			if strings.HasSuffix(pattern, "/") {
				// Directory pattern
				im.ignoreDirs = append(im.ignoreDirs, strings.TrimSuffix(pattern, "/"))
			} else {
				// File pattern
				if g, err := glob.Compile(pattern); err == nil {
					im.ignorePatterns = append(im.ignorePatterns, g)
				}
			}
		}

	}
}

func (im *IgnoreManager) ShouldIgnore(path string) bool {
	if im.includeIgnored {
		return false
	}
	// Check if the path or any of its parent directories should be ignored
	dir := path
	for dir != "." && dir != "/" {
		baseName := filepath.Base(dir)
		for _, ignoreDir := range im.ignoreDirs {
			if baseName == ignoreDir {
				return true
			}
		}
		dir = filepath.Dir(dir)
	}

	for _, pattern := range im.ignorePatterns {
		if pattern.Match(path) {
			return true
		}
	}

	for _, skipPath := range im.skipPaths {
		if skipPath == path {
			return true
		}
	}

	return false
}
