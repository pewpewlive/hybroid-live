package lsp

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestFindProjectRoot_EmptyInputs is a deterministic sanity check
// for the empty-input contract: filePath=="" or markers==nil/[]
// must return "" without touching the filesystem.
func TestFindProjectRoot_EmptyInputs(t *testing.T) {
	cases := []struct {
		name     string
		filePath string
		markers  []string
	}{
		{"empty filePath", "", []string{"hybconfig.toml"}},
		{"nil markers", "/some/path", nil},
		{"empty markers", "/some/path", []string{}},
		{"both empty", "", nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := findProjectRoot(c.filePath, c.markers)
			if got != "" {
				t.Errorf("findProjectRoot(%q, %v) = %q, want \"\"", c.filePath, c.markers, got)
			}
		})
	}
}

// TestFindProjectRoot_Determinism asserts that the function is pure:
// same inputs always produce the same output. (Today this is
// trivially true because the function only does os.Stat, but the
// test pins the contract.)
func TestFindProjectRoot_Determinism(t *testing.T) {
	dir := t.TempDir()
	// Create a marker in the temp dir.
	if err := os.WriteFile(filepath.Join(dir, "hybconfig.toml"), []byte(""), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	marker := []string{"hybconfig.toml"}
	child := filepath.Join(dir, "sub", "child.hyb")

	first := findProjectRoot(child, marker)
	second := findProjectRoot(child, marker)
	if first != second {
		t.Errorf("non-deterministic: %q vs %q", first, second)
	}
	if first != dir {
		t.Errorf("got %q, want %q (the dir containing the marker)", first, dir)
	}
}

// TestFindProjectRoot_PicksClosestMarker covers the closest-marker
// contract: when nested directories each contain a marker, the
// function returns the closest one (not the outermost).
func TestFindProjectRoot_PicksClosestMarker(t *testing.T) {
	dir := t.TempDir()
	// Marker in the root.
	if err := os.WriteFile(filepath.Join(dir, "hybconfig.toml"), []byte(""), 0o644); err != nil {
		t.Fatalf("setup outer: %v", err)
	}
	// Marker in a subdir — should win.
	sub := filepath.Join(dir, "sub")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatalf("setup subdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sub, "hybconfig.toml"), []byte(""), 0o644); err != nil {
		t.Fatalf("setup inner: %v", err)
	}

	grandchild := filepath.Join(sub, "leaf", "file.hyb")
	got := findProjectRoot(grandchild, []string{"hybconfig.toml"})
	if got != sub {
		t.Errorf("got %q, want %q (closest marker wins)", got, sub)
	}
}

// TestFindProjectRoot_StopsAtFilesystemRoot asserts that the function
// terminates when it reaches the filesystem root (the walk has
// `parent == dir` as its stop condition). Without that, this test
// would hang and the test framework would kill it with a timeout.
func TestFindProjectRoot_StopsAtFilesystemRoot(t *testing.T) {
	// A path that doesn't exist, in a parent that doesn't contain
	// a marker. The walk should climb all the way to "/" (or the
	// platform equivalent) and return "".
	done := make(chan string, 1)
	go func() {
		done <- findProjectRoot("/nonexistent/path/that/does/not/exist/leaf.hyb", []string{"hybconfig.toml"})
	}()
	select {
	case got := <-done:
		if got != "" {
			t.Errorf("got %q, want \"\" (no marker in any parent)", got)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("findProjectRoot did not terminate — likely walking past filesystem root")
	}
}

// FuzzFindProjectRoot exercises findProjectRoot with arbitrary
// inputs to catch panics, infinite loops, and contract violations.
//
// The `markersCSV` parameter is a comma-separated list of marker
// names (Go's fuzzer only supports string, []byte, and a handful of
// numeric types for fuzz arguments — not []string, so we encode the
// list as a string). Empty entries are allowed; the split + filter
// step preserves the empty-marker edge case the production code
// handles (an empty marker name never matches a real file because
// os.Stat("") errors).
//
// Invariants checked:
//
//  1. No panic on any input (the harness wraps the call in a
//     defer/recover to convert a panic into a test failure with
//     the failing input attached — fuzzer can then minimize the
//     corpus entry).
//
//  2. Empty filePath or empty markers list ⇒ "" (this is also
//     covered by TestFindProjectRoot_EmptyInputs, but the fuzzer
//     ensures it holds for any path-shaped input).
//
//  3. Determinism: same inputs always return the same string.
//     (Catches accidental introduction of global state, time
//     dependencies, or random number generators.)
//
//  4. If the result is non-empty, it must be an ancestor of
//     filepath.Dir(filePath). The function walks strictly
//     upward, so the result can never be a sibling or child.
//
//  5. The result is bounded by the filesystem root — the walk
//     stops when filepath.Dir(dir) == dir. A test fixture with
//     a hang-prone input (covered in
//     TestFindProjectRoot_StopsAtFilesystemRoot) is the
//     deterministic counterpart; the fuzzer covers everything else.
func FuzzFindProjectRoot(f *testing.F) {
	// Seed corpus: a mix of realistic and adversarial inputs.
	// CSV encoding: "hybconfig.toml" | "" | "hybconfig.toml,another.toml"
	f.Add("/some/file.hyb", "hybconfig.toml")
	f.Add("", "hybconfig.toml")
	f.Add("/some/file.hyb", "")
	f.Add("relative/path.hyb", "hybconfig.toml")
	f.Add("/path/with spaces/and unicode 漢字.hyb", "hybconfig.toml")
	f.Add("/a/b/c/../../d/e.hyb", "hybconfig.toml,another.toml")
	// NUL byte — historically a panic source in C-string-based
	// code, harmless in Go but worth fuzzing.
	f.Add("/path/\x00/nul.hyb", "hybconfig.toml")
	// Very long path.
	f.Add("/"+string(make([]byte, 4096))+"file.hyb", "hybconfig.toml")
	// Markers with weird names.
	f.Add("/some/file.hyb", ",..,/,hybconfig.toml")
	f.Add("/some/file.hyb", string(make([]byte, 1024)))

	f.Fuzz(func(t *testing.T, filePath string, markersCSV string) {
		markers := splitMarkersCSV(markersCSV)

		// Run twice and assert determinism in a single iteration
		// (cheaper than splitting into two fuzz cases).
		var first, second string
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("findProjectRoot panicked on input %q, %v: %v", filePath, markers, r)
				}
			}()
			first = findProjectRoot(filePath, markers)
		}()
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("findProjectRoot panicked on second call with input %q, %v: %v", filePath, markers, r)
				}
			}()
			second = findProjectRoot(filePath, markers)
		}()

		if first != second {
			t.Fatalf("non-deterministic: %q vs %q for input %q, %v", first, second, filePath, markers)
		}

		// Empty input contract.
		if filePath == "" || len(markers) == 0 {
			if first != "" {
				t.Errorf("empty input %q, %v: got %q, want \"\"", filePath, markers, first)
			}
			return
		}

		// Result contract: if non-empty, it's an ancestor of filePath's dir.
		if first != "" {
			dir := filepath.Clean(filepath.Dir(filePath))
			rel, err := filepath.Rel(first, dir)
			if err != nil {
				t.Errorf("filepath.Rel(%q, %q) errored: %v", first, dir, err)
				return
			}
			// rel is ".." if first is a strict ancestor of dir,
			// "." if first == dir, or a relative path if first
			// is a descendant. The function only walks up, so
			// ".." and "." are the only valid outcomes.
			if rel != ".." && rel != "." {
				t.Errorf("result %q is not an ancestor or equal to %q (rel=%q)", first, dir, rel)
			}
		}
	})
}

// splitMarkersCSV splits a comma-separated marker list, filtering
// out empty entries (which never match a real file via os.Stat).
// It's defined here (rather than reused from elsewhere) to keep
// the fuzz test self-contained.
func splitMarkersCSV(s string) []string {
	if s == "" {
		return nil
	}
	out := make([]string, 0, 4)
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	out = append(out, s[start:])
	// Filter empties.
	filtered := out[:0]
	for _, m := range out {
		if m != "" {
			filtered = append(filtered, m)
		}
	}
	return filtered
}

// TestFindProjectRoot_NonEmptyResultContainsMarker is a non-fuzz
// regression test: when the function returns a non-empty result for
// a real filesystem, the result directory must contain at least one
// of the input markers as a direct child. Catches a class of bugs
// where the function could return a parent that doesn't have the
// marker (e.g. due to a stale os.Stat cache, wrong join, etc.).
func TestFindProjectRoot_NonEmptyResultContainsMarker(t *testing.T) {
	dir := t.TempDir()
	markerName := "hybconfig.toml"
	if err := os.WriteFile(filepath.Join(dir, markerName), []byte(""), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	child := filepath.Join(dir, "deeply", "nested", "file.hyb")
	got := findProjectRoot(child, []string{markerName})
	if got == "" {
		t.Fatalf("expected non-empty result, got \"\"")
	}

	entries, err := os.ReadDir(got)
	if err != nil {
		t.Fatalf("ReadDir(%q): %v", got, err)
	}
	found := false
	for _, e := range entries {
		if e.Name() == markerName && !e.IsDir() {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("result dir %q does not contain marker %q (entries: %v)", got, markerName, entryNames(entries))
	}
}

// entryNames is a small helper for the assertion error message.
func entryNames(entries []os.DirEntry) []string {
	out := make([]string, len(entries))
	for i, e := range entries {
		out[i] = e.Name()
	}
	return out
}

// TestFindProjectRoot_MultipleMarkersOrderIndependent verifies that
// the choice of marker in the result is determined by the filesystem
// layout, not by the order in the markers slice. (The current
// implementation iterates markers in slice order and returns on the
// first match, but that order doesn't affect which dir is chosen —
// it only matters when two markers coexist in the same dir, which
// is degenerate.)
func TestFindProjectRoot_MultipleMarkersOrderIndependent(t *testing.T) {
	dir := t.TempDir()
	// Same dir has two markers.
	if err := os.WriteFile(filepath.Join(dir, "a.toml"), []byte(""), 0o644); err != nil {
		t.Fatalf("setup a: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "b.toml"), []byte(""), 0o644); err != nil {
		t.Fatalf("setup b: %v", err)
	}
	child := filepath.Join(dir, "sub", "file.hyb")

	got1 := findProjectRoot(child, []string{"a.toml", "b.toml"})
	got2 := findProjectRoot(child, []string{"b.toml", "a.toml"})

	// Both calls should pick the same dir, even if they picked
	// different markers within it.
	if got1 != got2 {
		t.Errorf("marker order changed result: %q vs %q", got1, got2)
	}
	if got1 != dir {
		t.Errorf("got %q, want %q", got1, dir)
	}
}
