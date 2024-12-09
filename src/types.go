package src

import "github.com/fargusplumdoodle/dump_dir/src/prompt"

type Config struct {
	Action         string
	Extensions     []string
	Directories    []string
	SkipDirs       []string
	SpecificFiles  []string
	IncludeIgnored bool
	MaxFileSize    int64
}

type Stats struct {
	TotalFiles      int
	TotalLines      int
	EstimatedTokens int
	ProcessedFiles  []prompt.FileInfo
	ParsedFiles     []prompt.FileInfo
	SkippedLarge    []prompt.FileInfo
	SkippedBinary   []prompt.FileInfo
}
