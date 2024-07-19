package tests

import (
	. "github.com/fargusplumdoodle/dump_dir/src"
	"github.com/spf13/afero"
	"os/exec"
	"path/filepath"
)

func setupTestFileSystem(fileSystem map[string]string) afero.Fs {
	fs := afero.NewMemMapFs()
	for path, content := range fileSystem {
		dir := filepath.Dir(path)
		fs.MkdirAll(dir, 0755)
		afero.WriteFile(fs, path, []byte(content), 0644)
	}
	return fs
}

// contains checks if a string slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func ResetExecCommand() {
	ExecCommand = exec.Command
}
