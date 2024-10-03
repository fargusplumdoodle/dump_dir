package tests

import (
	"strings"
	"testing"

	. "github.com/fargusplumdoodle/dump_dir/src/prompt"
)

func TestGenerateXMLPrompt(t *testing.T) {
	tests := []struct {
		name        string
		prompt      Prompt
		expectedXML string
	}{
		{
			name: "Empty Prompt",
			prompt: Prompt{
				ProcessedFiles: []FileInfo{},
				SkippedLarge:   []FileInfo{},
				SkippedBinary:  []FileInfo{},
			},
			expectedXML: `<?xml version="1.0" encoding="UTF-8"?>
<context>
  <fileList></fileList>
  <processedFiles></processedFiles>
  <skippedLarge></skippedLarge>
  <skippedBinary></skippedBinary>
</context>`,
		},
		{
			name: "Single Processed File",
			prompt: Prompt{
				ProcessedFiles: []FileInfo{
					{
						Path:     "src/main.go",
						Contents: "package main\n\nfunc main() {}\n",
						Status:   StatusParsed,
					},
				},
				SkippedLarge:  []FileInfo{},
				SkippedBinary: []FileInfo{},
			},
			expectedXML: `<?xml version="1.0" encoding="UTF-8"?>
<context>
  <fileList>
    <file path="src/main.go"></file>
  </fileList>
  <processedFiles>
    <file path="src/main.go"><![CDATA[package main

func main() {}
]]></file>
  </processedFiles>
  <skippedLarge></skippedLarge>
  <skippedBinary></skippedBinary>
</context>`,
		},
		{
			name: "Multiple Processed and Skipped Files",
			prompt: Prompt{
				ProcessedFiles: []FileInfo{
					{
						Path:     "src/main.go",
						Contents: "package main\n\nfunc main() {}\n",
						Status:   StatusParsed,
					},
					{
						Path:     "src/utils.go",
						Contents: "package utils\n\nfunc Helper() {}\n",
						Status:   StatusParsed,
					},
				},
				SkippedLarge: []FileInfo{
					{
						Path:     "data/large_file.txt",
						Contents: "<FILE TOO LARGE>",
						Status:   StatusSkippedTooLarge,
					},
				},
				SkippedBinary: []FileInfo{
					{
						Path:     "bin/executable",
						Contents: "<BINARY SKIPPED>",
						Status:   StatusSkippedBinary,
					},
				},
			},
			expectedXML: `<?xml version="1.0" encoding="UTF-8"?>
<context>
  <fileList>
    <file path="src/main.go"></file>
    <file path="src/utils.go"></file>
  </fileList>
  <processedFiles>
    <file path="src/main.go"><![CDATA[package main

func main() {}
]]></file>
    <file path="src/utils.go"><![CDATA[package utils

func Helper() {}
]]></file>
  </processedFiles>
  <skippedLarge>
    <file path="data/large_file.txt"></file>
  </skippedLarge>
  <skippedBinary>
    <file path="bin/executable"></file>
  </skippedBinary>
</context>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xmlPrompt := GenerateXMLPrompt(tt.prompt)
			outputXML, err := MarshalXMLPrompt(xmlPrompt)
			if err != nil {
				t.Fatalf("Error marshaling XML: %v", err)
			}
			// Trim spaces for comparison
			expected := strings.TrimSpace(tt.expectedXML)
			actual := strings.TrimSpace(outputXML)
			if expected != actual {
				t.Errorf("Generated XML does not match expected.\nExpected:\n%s\n\nGot:\n%s", expected, actual)
			}
		})
	}
}
