package lsp

import (
	"os"
	"path/filepath"
)

// findProjectRoot walks up from filePath's directory looking for any of the
// given marker filenames. Returns the directory that contains the first
// marker found, or "" if none is found up to the filesystem root.
//
// The walk is bounded by the filesystem (it stops at the root directory) and
// terminates as soon as a marker is found, so it is cheap in practice.
//
// This is used as a fallback for single-file opens: when the client does not
// provide a workspace root, we still want to locate a Hybroid project if the
// file lives inside one. The convention matches TypeScript's tsconfig.json
// walk, Pylance's extraPaths, and clangd's compile_commands.json discovery.
func findProjectRoot(filePath string, markers []string) string {
	if filePath == "" || len(markers) == 0 {
		return ""
	}

	dir := filepath.Dir(filePath)
	dir = filepath.Clean(dir)

	for {
		for _, marker := range markers {
			candidate := filepath.Join(dir, marker)
			if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
				return dir
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}
