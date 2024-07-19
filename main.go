package main

import (
	. "github.com/fargusplumdoodle/dump_dir/src"
	"github.com/spf13/afero"
	"os"
)

const version = "1.0.0" // You can update this as needed

func main() {
	args := os.Args[1:]

	if !ValidateArgs(args) {
		PrintUsage()
		return
	}
	config := ParseArgs(args)

	switch config.Action {
	case "help":
		PrintUsage()
		return
	case "version":
		PrintVersion()
		return
	case "dump_dir":
		performDumpDir(config)
	}
}

func performDumpDir(config Config) {
	fs := afero.NewOsFs()
	fileFinder := NewFileFinder(config, fs)
	fileProcessor := NewFileProcessor(fs)

	filePaths := fileFinder.DiscoverFiles()
	processedFiles := fileProcessor.ProcessFiles(filePaths)
	stats := CalculateStats(processedFiles)
	PrintDetailedOutput(stats)
}

func PrintVersion() {
	println("dump_dir version:", version)
}
