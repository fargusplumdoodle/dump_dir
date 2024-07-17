package main

import (
	. "github.com/fargusplumdoodle/dump_dir/src"
	"os"
)

func main() {
	if !ValidateArgs() {
		return
	}

	extension, directories, skipDirs, specificFiles := ParseArgs(os.Args[1:])

	matchingFiles, detailedOutput, totalLines := ProcessDirectories(extension, directories, skipDirs, specificFiles)
	PrintDetailedOutput(matchingFiles, detailedOutput, totalLines)
}
