package main

import (
	. "github.com/fargusplumdoodle/dump_dir/src"
	"github.com/spf13/afero"
	"os"
)

func main() {
	args := os.Args[1:]

	if !ValidateArgs(args) {
		PrintUsage()
		return
	}
	config := ParseArgs(args)

	fs := afero.NewOsFs()
	fileFinder := NewFileFinder(config, fs)
	fileProcessor := NewFileProcessor(fs)

	filePaths := fileFinder.DiscoverFiles()
	processedFiles := fileProcessor.ProcessFiles(filePaths)
	stats := CalculateStats(processedFiles)
	PrintDetailedOutput(stats)
}
