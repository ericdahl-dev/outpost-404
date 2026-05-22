package game

import (
	"os"
	"path/filepath"
	"sync"
)

var (
	repoRootOnce sync.Once
	repoRoot     string
)

// RepoRoot returns the module root (directory containing go.mod).
func RepoRoot() string {
	repoRootOnce.Do(func() {
		dir, err := os.Getwd()
		if err != nil {
			repoRoot = "."
			return
		}
		for {
			if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
				repoRoot = dir
				return
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				repoRoot = "."
				return
			}
			dir = parent
		}
	})
	return repoRoot
}

// ScriptsDir returns the repo scripts/ directory.
func ScriptsDir() string {
	return filepath.Join(RepoRoot(), "scripts")
}

// DataDir returns the repo data/ directory.
func DataDir() string {
	return filepath.Join(RepoRoot(), "data")
}

// TestFixturePath joins path segments under test/internal/game/testdata/.
func TestFixturePath(parts ...string) string {
	base := []string{RepoRoot(), "test", "internal", "game", "testdata"}
	return filepath.Join(append(base, parts...)...)
}
