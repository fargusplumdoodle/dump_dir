package src

type FileInfo struct {
	Path     string
	Contents string
}

type Config struct {
	Action         string
	Extensions     []string
	Directories    []string
	SkipDirs       []string
	SpecificFiles  []string
	IncludeIgnored bool
}

type Stats struct {
	TotalFiles      int
	TotalLines      int
	EstimatedTokens int
	ProcessedFiles  []FileInfo
}
