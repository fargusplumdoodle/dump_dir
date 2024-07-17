package src

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
)

type IgnoreManager struct {
	ignorePatterns []glob.Glob
}

func NewIgnoreManager() (*IgnoreManager, error) {
	im := &IgnoreManager{}
	err := im.loadIgnorePatterns()
	if err != nil {
		return nil, err
	}
	return im, nil
}

func (im *IgnoreManager) loadIgnorePatterns() error {
	if g, err := glob.Compile(".git"); err == nil {
		im.ignorePatterns = append(im.ignorePatterns, g)
	}
	// Load global gitignore
	homeDir, err := os.UserHomeDir()
	if err == nil {
		globalGitignorePath := filepath.Join(homeDir, ".gitignore_global")
		im.loadIgnoreFile(globalGitignorePath)
	}

	// Load project-specific gitignore
	im.loadIgnoreFile(".gitignore")

	return nil
}

func (im *IgnoreManager) loadIgnoreFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		return // Ignore errors, as the file might not exist
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pattern := strings.TrimSpace(scanner.Text())
		if pattern != "" && !strings.HasPrefix(pattern, "#") {
			if g, err := glob.Compile(pattern); err == nil {
				im.ignorePatterns = append(im.ignorePatterns, g)
			}
		}
	}
}

func (im *IgnoreManager) ShouldIgnore(path string) bool {
	// Always ignore .git directories
	if strings.Contains(path, string(os.PathSeparator)+".git"+string(os.PathSeparator)) {
		return true
	}

	for _, pattern := range im.ignorePatterns {
		if pattern.Match(path) {
			return true
		}
	}
	return false
}

func UpdateFileProcessor(fp *FileProcessor) error {
	im, err := NewIgnoreManager()
	if err != nil {
		return err
	}
	fp.IgnoreManager = im
	return nil
}
