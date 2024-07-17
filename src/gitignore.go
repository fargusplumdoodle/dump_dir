package src

import (
	"bufio"
	"fmt"
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
		globalPatterns := readIgnoreFile(globalGitignore)
		patterns = append(patterns, globalPatterns...)
		fmt.Printf("Global .gitignore patterns (%s):\n", globalGitignore)
		for _, p := range globalPatterns {
			fmt.Printf("  %s\n", p)
		}
	}

	// Add patterns from local .gitignore files
	for _, dir := range directories {
		localGitignore := filepath.Join(dir, ".gitignore")
		localPatterns := readIgnoreFile(localGitignore)
		patterns = append(patterns, localPatterns...)
		fmt.Printf("Local .gitignore patterns (%s):\n", localGitignore)
		for _, p := range localPatterns {
			fmt.Printf("  %s\n", p)
		}
	}

	return patterns
}

func readIgnoreFile(path string) []glob.Glob {
	var patterns []glob.Glob
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Could not open ignore file %s: %v\n", path, err)
		return patterns
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			if g, err := glob.Compile(line); err == nil {
				patterns = append(patterns, g)
			} else {
				fmt.Printf("Error compiling ignore pattern '%s': %v\n", line, err)
			}
		}
	}

	return patterns
}

// Add this new function to log ignored files
func LogIgnoredFile(path string, pattern glob.Glob) {
	fmt.Printf("Ignoring file: %s (matched by pattern: %s)\n", path, pattern)
}
