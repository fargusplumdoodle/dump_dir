package main

import (
	"fmt"
	. "github.com/fargusplumdoodle/dump_dir/src"
	"github.com/gobwas/glob"
	"os"
)

func main() {
	if !ValidateArgs() {
		return
	}

	extension, directories, skipDirs, includeGitIgnored := ParseArgs(os.Args[1:])

	ignorePatterns := []glob.Glob{}
	if !includeGitIgnored {
		ignorePatterns = GetIgnorePatterns(directories)
	}

	matchingFiles, detailedOutput, totalLines := ProcessDirectories(extension, directories, skipDirs, ignorePatterns, includeGitIgnored)
	summary := GenerateSummary(matchingFiles, totalLines)

	if CopyToClipboard(detailedOutput) {
		summary += BoldGreen("✅ File contents have been copied to clipboard.\n")
	}

	fmt.Println(summary)
}
