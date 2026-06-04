package lsp

import (
	"os"
	"path/filepath"
	"testing"
)

// TestHandleDidOpen_SecondFileAfterRootFound_StaysInProjectMode documents
// the design choice that once a project root has been established (via
// the first didOpen discovering a hybconfig.toml ancestor), all
// subsequent didOpens stay in project mode — even if the new file lives
// outside the project tree.
//
// Rationale: a user who opens a project and then opens an unrelated
// .hyb file from elsewhere almost certainly wants the editor to keep
// the project context (so completion/hover/etc. all reference the
// project symbols). Dropping into single-file mode for the second
// file would be surprising.
//
// The test would catch a future refactor that accidentally enables
// per-file root resolution (i.e. moves the findProjectRoot call out
// of the `if h.rootPath == ""` guard).
func TestHandleDidOpen_SecondFileAfterRootFound_StaysInProjectMode(t *testing.T) {
	projectDir := writeProject(t, map[string]string{
		"hybconfig.toml": minimalHybConfig,
		"level.hyb":      minimalLevelSource,
	})
	uriA := toURI(filepath.Join(projectDir, "level.hyb"))

	h, conn := newTestHandler(t)

	// First open: in-project. Establishes h.rootPath.
	openForTest(t, h, conn, uriA, minimalLevelSource)
	if h.rootPath == "" {
		t.Fatalf("expected rootPath to be set after first didOpen")
	}
	firstRoot := h.rootPath
	firstInfoCount := countInfoDiags(conn, uriA)

	// Second open: a file in a completely separate tree (no
	// hybconfig.toml anywhere up). Even though the second file has
	// no project ancestor, the handler must NOT publish the
	// Information diagnostic — we are in project mode now.
	strayDir := t.TempDir()
	pathHasNoProjectMarker(t, strayDir)
	uriB := toURI(filepath.Join(strayDir, "stray.hyb"))

	openForTest(t, h, conn, uriB, "env Stray as Level\n")

	if h.rootPath != firstRoot {
		t.Errorf("rootPath changed from %q to %q after second didOpen",
			firstRoot, h.rootPath)
	}
	if countInfoDiags(conn, uriB) != 0 {
		t.Errorf("expected no Information diagnostic for %q (project mode already active), got %d",
			uriB, countInfoDiags(conn, uriB))
	}
	// Sanity: the first file's info count should also be unchanged.
	if got := countInfoDiags(conn, uriA); got != firstInfoCount {
		t.Errorf("first file's info diagnostic count changed: was %d, now %d",
			firstInfoCount, got)
	}
}

// TestFindProjectRoot_NonexistentFilePath verifies the cheap edge case:
// the function never touches filePath itself (only its directory), so
// passing a path whose file doesn't exist must not panic. The walk
// proceeds from the directory and either finds a marker or returns "".
//
// This is the kind of "obviously fine" code that gets a panic-only
// years later when someone tries to optimize by stat-ing the file.
func TestFindProjectRoot_NonexistentFilePath(t *testing.T) {
	dir := t.TempDir()
	pathHasNoProjectMarker(t, dir)

	// The file doesn't exist. findProjectRoot should walk up from
	// the directory and return "" since no ancestor has a marker.
	ghost := filepath.Join(dir, "this", "does", "not", "exist.hyb")

	got := findProjectRoot(ghost, []string{"hybconfig.toml"})
	if got != "" {
		t.Errorf("expected empty result for nonexistent path in marker-free tree, got %q", got)
	}

	// Sanity: with a marker present in the directory itself, the
	// function still returns the directory (it doesn't care that the
	// file is missing).
	withMarker := filepath.Join(dir, "also_missing.hyb")
	if err := os.WriteFile(filepath.Join(dir, "hybconfig.toml"), []byte(minimalHybConfig), 0o644); err != nil {
		t.Fatalf("write marker: %v", err)
	}
	got = findProjectRoot(withMarker, []string{"hybconfig.toml"})
	if filepath.Clean(got) != filepath.Clean(dir) {
		t.Errorf("expected dir %q for marker-in-dir, got %q", dir, got)
	}
}

// countInfoDiags returns the number of Information-severity diagnostics
// (severity 3) published for the given URI across all observed
// notifications. Used to assert the one-shot info-notice behavior
// without coupling tests to total notify counts.
func countInfoDiags(conn *fakeNotify, uri DocumentURI) int {
	n := 0
	for _, c := range conn.Notifies() {
		if c.Method != "textDocument/publishDiagnostics" {
			continue
		}
		p, ok := c.Params.(PublishDiagnosticsParams)
		if !ok || p.URI != uri {
			continue
		}
		for _, d := range p.Diagnostics {
			if d.Severity == 3 {
				n++
			}
		}
	}
	return n
}

// keep strings imported for future use.
