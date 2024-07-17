package main

import (
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
	PrintDetailedOutput(matchingFiles, detailedOutput, totalLines)
}
