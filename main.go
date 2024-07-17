package main

import (
	. "github.com/fargusplumdoodle/dump_dir/src"
	"os"
)

func main() {
	args := os.Args[1:]

	if !ValidateArgs(args) {
		PrintUsage()
		return
	}
	config := ParseArgs(args)
	fileProcessor := NewFileProcessor(config)

	matchingFiles, detailedOutput, totalLines := fileProcessor.ProcessDirectories()
	PrintDetailedOutput(matchingFiles, detailedOutput, totalLines)
}
