package src

import (
	"bufio"
	"github.com/gobwas/glob"
	"os"
	"path/filepath"
	"strings"
)

func GetIgnorePatterns(directories []string) []glob.Glob {
	var patterns []glob.Glob

	// Add patterns from global .gitignore
	homeDir, err := os.UserHomeDir()
	if err == nil {
		globalGitignore := filepath.Join(homeDir, ".gitignore_global")
		patterns = append(patterns, readIgnoreFile(globalGitignore)...)
	}

	// Add patterns from local .gitignore files
	for _, dir := range directories {
		localGitignore := filepath.Join(dir, ".gitignore")
		patterns = append(patterns, readIgnoreFile(localGitignore)...)
	}

	return patterns
}
func readIgnoreFile(path string) []glob.Glob {
	var patterns []glob.Glob
	file, err := os.Open(path)
	if err != nil {
		return patterns
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			if g, err := glob.Compile(line); err == nil {
				patterns = append(patterns, g)
			}
		}
	}

	return patterns
}
