package main

import (
	"bufio"
	"fmt"
	"github.com/gobwas/glob"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/fatih/color"
)

type FileInfo struct {
	Path     string
	Contents string
}

var (
	boldGreen   = color.New(color.FgGreen, color.Bold).SprintfFunc()
	boldCyan    = color.New(color.FgCyan, color.Bold).SprintfFunc()
	boldMagenta = color.New(color.FgMagenta, color.Bold).SprintfFunc()
	boldRed     = color.New(color.FgRed, color.Bold).SprintfFunc()
)

func main() {
	if !validateArgs() {
		return
	}

	extension, directories, skipDirs, includeGitIgnored := parseArgs(os.Args[1:])

	ignorePatterns := []glob.Glob{}
	if !includeGitIgnored {
		ignorePatterns = getIgnorePatterns(directories)
	}

	matchingFiles, detailedOutput, totalLines := processDirectories(extension, directories, skipDirs, ignorePatterns, includeGitIgnored)
	summary := generateSummary(matchingFiles, totalLines)

	if copyToClipboard(detailedOutput) {
		summary += boldGreen("✅ File contents have been copied to clipboard.\n")
	}

	fmt.Println(summary)
}

func validateArgs() bool {
	if len(os.Args) < 3 {
		fmt.Println()
		fmt.Println(boldRed("❌ Error: Insufficient arguments"))
		fmt.Println()
		fmt.Println(boldCyan("Usage:"))
		fmt.Println("  dump_dir <file_extension> <directory1> [directory2] ... [-s <skip_directory1>] [-s <skip_directory2>] ... [--include-gitignored-paths]")
		fmt.Println("  Use 'any' as file_extension to match all files")
		fmt.Println()
		fmt.Println(boldGreen("Example:"))
		fmt.Println("  dump_dir js ./project -s ./project/node_modules -s ./project/dist")
		fmt.Println("  dump_dir any ./project -s ./project/node_modules --include-gitignored-paths")
		fmt.Println()
		fmt.Println(boldMagenta("Description:"))
		fmt.Println("  This will search for files with the specified extension (or all files if 'any' is used)")
		fmt.Println("  in the given directories, excluding any specified directories and .gitignore'd files by default.")
		fmt.Println("  Use --include-gitignored-paths to include files that would be ignored by .gitignore.")
		fmt.Println()
		return false
	}
	return true
}

func parseArgs(args []string) (string, []string, []string, bool) {
	extension := args[0]
	var directories []string
	var skipDirs []string
	skipMode := false
	includeGitIgnored := false

	for _, arg := range args[1:] {
		if arg == "-s" {
			skipMode = true
			continue
		}
		if arg == "--include-gitignored-paths" {
			includeGitIgnored = true
			continue
		}
		if skipMode {
			skipDirs = append(skipDirs, arg)
			skipMode = false
		} else {
			directories = append(directories, arg)
		}
	}

	return extension, directories, skipDirs, includeGitIgnored
}

func getIgnorePatterns(directories []string) []glob.Glob {
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

func processDirectories(extension string, directories, skipDirs []string, ignorePatterns []glob.Glob, includeGitIgnored bool) ([]string, string, int) {
	var matchingFiles []string
	var detailedOutput strings.Builder
	var totalLines int

	for _, dir := range directories {
		files, output, lines := processDirectory(dir, extension, skipDirs, ignorePatterns, includeGitIgnored)
		matchingFiles = append(matchingFiles, files...)
		detailedOutput.WriteString(output)
		totalLines += lines
	}

	return matchingFiles, detailedOutput.String(), totalLines
}

func processDirectory(dir, extension string, skipDirs []string, ignorePatterns []glob.Glob, includeGitIgnored bool) ([]string, string, int) {
	var matchingFiles []string
	var detailedOutput strings.Builder
	var totalLines int

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignore .git directory unless --include-gitignored-paths is used
		if !includeGitIgnored && info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		if info.IsDir() {
			for _, skipDir := range skipDirs {
				if strings.HasPrefix(path, skipDir) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		relPath, _ := filepath.Rel(dir, path)
		if !includeGitIgnored {
			for _, pattern := range ignorePatterns {
				if pattern.Match(relPath) {
					return nil
				}
			}
		}

		if extension == "any" || strings.HasSuffix(info.Name(), "."+extension) {
			fileInfo, err := processFile(path)
			if err != nil {
				return err
			}
			matchingFiles = append(matchingFiles, fileInfo.Path)
			detailedOutput.WriteString(fmt.Sprintf("📄 File: %s\n%s\n\n", boldCyan(fileInfo.Path), fileInfo.Contents))
			totalLines += strings.Count(fileInfo.Contents, "\n")
		}
		return nil
	})

	if err != nil {
		fmt.Printf(boldRed("❌ Error walking directory %s: %v\n"), dir, err)
	}

	return matchingFiles, detailedOutput.String(), totalLines
}

func processFile(path string) (FileInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return FileInfo{}, err
	}
	defer file.Close()

	var contents strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		contents.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		return FileInfo{}, err
	}

	return FileInfo{Path: path, Contents: contents.String()}, nil
}

func generateSummary(matchingFiles []string, totalLines int) string {
	var summary strings.Builder
	summary.WriteString(boldMagenta("🔍 Matching files:\n"))
	for _, file := range matchingFiles {
		summary.WriteString(fmt.Sprintf("  - %s\n", file))
	}
	summary.WriteString(boldCyan(fmt.Sprintf("📚 Total files found: %d\n", len(matchingFiles))))
	summary.WriteString(boldCyan(fmt.Sprintf("📝 Total lines across all files: %d\n\n", totalLines)))
	return summary.String()
}

func copyToClipboard(content string) bool {
	err := clipboard.WriteAll(content)
	if err != nil {
		fmt.Println(boldRed(fmt.Sprintf("❌ Error copying to clipboard: %v", err)))
		return false
	}
	return true
}
