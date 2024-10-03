package prompt

import "encoding/xml"

type Prompt struct {
	ProcessedFiles []FileInfo
	SkippedLarge   []FileInfo
	SkippedBinary  []FileInfo
}

type XMLFile struct {
	Path     string `xml:"path,attr"`
	Contents string `xml:",cdata"`
}

type XMLPrompt struct {
	XMLName       xml.Name  `xml:"context"`
	FileList      []XMLFile `xml:"fileList>file"`
	Processed     []XMLFile `xml:"processedFiles>file"`
	SkippedLarge  []XMLFile `xml:"skippedLarge>file"`
	SkippedBinary []XMLFile `xml:"skippedBinary>file"`
}

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
