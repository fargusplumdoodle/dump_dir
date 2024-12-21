package main

import (
	"fmt"
	. "github.com/fargusplumdoodle/dump_dir/src"
	"github.com/spf13/afero"
	"os"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	fs := afero.NewOsFs()
	clipboard := NewSystemClipboard()
	runCfg := RunConfig{
		Version:   version,
		Commit:    commit,
		Date:      date,
		Fs:        fs,
		Clipboard: clipboard,
	}

	if err := Run(os.Args[1:], runCfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
