package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// GetProjectRoot is a variable holding a function to discover the project's root directory by locating the go.mod file.
var GetProjectRoot = func() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}

	dir := filepath.Dir(filename)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			return "", fmt.Errorf("project root not found (could not locate go.mod)")
		}
		dir = parentDir
	}
}
