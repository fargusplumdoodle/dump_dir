package src

type FileStatus string

const (
	StatusParsed          FileStatus = "PARSED"
	StatusSkippedBinary   FileStatus = "SKIPPED_BINARY"
	StatusSkippedTooLarge FileStatus = "SKIPPED_TOO_LARGE"
)

type FileInfo struct {
	Path     string
	Contents string
	Status   FileStatus
}

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
	ProcessedFiles  []FileInfo
	ParsedFiles     []FileInfo
	SkippedLarge    []FileInfo
	SkippedBinary   []FileInfo
}
