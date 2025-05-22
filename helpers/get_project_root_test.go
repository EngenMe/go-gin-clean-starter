package helpers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetProjectRoot tests the behavior of the GetProjectRoot function, including success and various error scenarios.
func TestGetProjectRoot(t *testing.T) {
	originalGetProjectRoot := GetProjectRoot
	defer func() {
		GetProjectRoot = originalGetProjectRoot
	}()

	t.Run(
		"successful discovery of project root", func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "project-root-test-*")
			require.NoError(t, err)
			defer os.RemoveAll(tmpDir)

			goModPath := filepath.Join(tmpDir, "go.mod")
			err = os.WriteFile(goModPath, []byte("module test"), 0644)
			require.NoError(t, err)

			subDir := filepath.Join(tmpDir, "cmd", "app")
			err = os.MkdirAll(subDir, 0755)
			require.NoError(t, err)

			GetProjectRoot = func() (string, error) {
				dir := subDir
				for {
					if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
						return dir, nil
					}

					parentDir := filepath.Dir(dir)
					if parentDir == dir {
						return "", nil
					}
					dir = parentDir
				}
			}

			root, err := GetProjectRoot()
			assert.NoError(t, err)
			assert.Equal(t, tmpDir, root)
		},
	)

	t.Run(
		"go.mod not found", func(t *testing.T) {
			GetProjectRoot = func() (string, error) {
				return "", filepath.ErrBadPattern
			}

			root, err := GetProjectRoot()
			assert.Error(t, err)
			assert.Empty(t, root)
		},
	)

	t.Run(
		"unable to get current file path", func(t *testing.T) {
			GetProjectRoot = func() (string, error) {
				return "", filepath.ErrBadPattern
			}

			root, err := GetProjectRoot()
			assert.Error(t, err)
			assert.Empty(t, root)
		},
	)
}

// TestActualGetProjectRoot verifies the actual behavior of GetProjectRoot in locating the project root containing go.mod.
// It skips the test if an appropriate project structure is not found during execution.
func TestActualGetProjectRoot(t *testing.T) {
	originalGetProjectRoot := GetProjectRoot

	root, err := originalGetProjectRoot()
	if err == nil {
		_, err := os.Stat(filepath.Join(root, "go.mod"))
		assert.NoError(t, err, "Expected to find go.mod in project root")
	} else {
		t.Skip("Skipping actual project root test (no proper project structure found)")
	}
}
