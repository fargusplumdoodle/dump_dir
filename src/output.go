package src

import (
	"fmt"
	"strings"
)

func PrintUsage() {
	usage := `
` + boldCyan("Usage:") + `
  dump_dir [options] <path1> [path2] [options] ...

` + boldCyan("Options:") + `
  -h, --help                 Display this help information
  -v, --version              Display the version of dump_dir
  -s <directory>, --skip <directory>
                             Skip specified directory
  -e <extension[s]>, --extension <extension[s]>
                             Filter by specific file extensions
  --include-ignored          Include files that would normally be ignored
                             (e.g., those in .gitignore)
  -m <size>, --max-filesize <size>
                             Specify the maximum file size to process.
                             You can use units like B, KB, or MB
                             (e.g., 500KB, 2MB). If no unit is specified,
                             it defaults to 500KB.
  -g <pattern>, --glob <pattern>
                             Match file names with a glob pattern.
                             Does not support matching directory names
                             or ** patterns.
  -nc, --no-config           Ignore the .dump_dir.yml configuration file

` + BoldGreen("Common examples:") + `
  # Grab everything from ./project
  dump_dir ./project

  # Grab everything from ./src, but skip the csv_data directory
  dump_dir ./src -s ./src/csv_data  

  # Grab all of the .js files from the current directory 
  # skipping contents of the dist directory
  dump_dir . -e js --skip ./dist

  # Grab EVERYTHING including what is in your gitignore
  dump_dir . --include-ignored

  # Quickly grab one thing, ignoring your auto-included files
  dump_dir ./next.conf.ts -nc

  # Grab all of the test files
  dump_dir ./project --glob "*_test.go"

` + boldMagenta("Description:") + `
  dump_dir will find files based on your parameters
  and put their contents into your clipboard in a way
  that is easy for Large Language Models to understand.

  You can make a .dump_dir.yml config file to automatically
  include/exclude paths.

  More documentation at: https://github.com/fargusplumdoodle/dump_dir

`
	fmt.Print(usage)
}

func PrintError(errorType string, filePath string, err error) {
	fmt.Printf(boldRed("❌ Error %s file %s: %v\n", errorType, filePath, err))
}

func CopyToClipboard(clipboard ClipboardManager, content string) bool {
	err := clipboard.WriteAll(content)
	if err != nil {
		fmt.Println(boldRed(fmt.Sprintf("❌ Error copying to clipboard: %v", err)))
		return false
	}
	return true
}

func FormatFileContent(path, contents string) string {
	return fmt.Sprintf("START FILE: %s\n%s\nEND FILE: %s\n\n", path, contents, path)
}

func GenerateDetailedOutput(stats Stats) string {
	var detailedOutput strings.Builder

	for _, fileInfo := range stats.ProcessedFiles {
		detailedOutput.WriteString(FormatFileContent(fileInfo.Path, fileInfo.Contents))
	}

	return detailedOutput.String()
}

func PrintDetailedOutput(stats Stats, config RunConfig) {
	detailedOutput := GenerateDetailedOutput(stats)
	summary := DisplayStats(stats)

	if CopyToClipboard(config.Clipboard, detailedOutput) {
		summary += BoldGreen("✅ File contents have been copied to clipboard.\n")
	}

	fmt.Println(summary)
}
