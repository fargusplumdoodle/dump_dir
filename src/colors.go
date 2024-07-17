package src

import "github.com/fatih/color"

var (
	BoldGreen   = color.New(color.FgGreen, color.Bold).SprintfFunc()
	boldCyan    = color.New(color.FgCyan, color.Bold).SprintfFunc()
	boldMagenta = color.New(color.FgMagenta, color.Bold).SprintfFunc()
	boldRed     = color.New(color.FgRed, color.Bold).SprintfFunc()
)
