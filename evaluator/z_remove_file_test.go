package evaluator

import (
	"hybroid/core"
	"testing"
)

// TestRemoveFile_DropsAllPerFileState verifies that RemoveFile
// removes a file from every internal collection: walkers (by
// source path AND abs path), walkerList, files, programs,
// parseAlerts, and fileContents.
func TestRemoveFile_DropsAllPerFileState(t *testing.T) {
	ev := NewEvaluator([]core.File{
		{DirectoryPath: "src", FileName: "foo", FileExtension: ".hyb"},
		{DirectoryPath: "src", FileName: "bar", FileExtension: ".hyb"},
	})

	// Add content for both files so fileContents/programs/parseAlerts
	// have entries.
	if err := ev.UpdateFileContent("src/foo.hyb", "env Foo as Level\n\ntick {}\n"); err != nil {
		t.Fatalf("UpdateFileContent foo: %v", err)
	}
	if err := ev.UpdateFileContent("src/bar.hyb", "env Bar as Level\n\ntick {}\n"); err != nil {
		t.Fatalf("UpdateFileContent bar: %v", err)
	}

	// Snapshot pre-removal sizes for the bar file.
	ev.mu.Lock()
	preFiles := len(ev.files)
	preList := len(ev.walkerList)
	preWalkers := len(ev.walkers)
	prePrograms := len(ev.programs)
	preAlerts := len(ev.parseAlerts)
	preContents := len(ev.fileContents)
	ev.mu.Unlock()

	if preFiles != 2 || preList != 2 {
		t.Fatalf("setup: preFiles=%d preList=%d, want 2/2", preFiles, preList)
	}
	if preWalkers < 4 || prePrograms != 2 || preAlerts != 2 || preContents != 2 {
		t.Fatalf("setup: preWalkers=%d prePrograms=%d preAlerts=%d preContents=%d",
			preWalkers, prePrograms, preAlerts, preContents)
	}

	if !ev.RemoveFile("src/bar.hyb") {
		t.Fatal("RemoveFile returned false for known file")
	}

	ev.mu.Lock()
	postFiles := len(ev.files)
	postList := len(ev.walkerList)
	postWalkers := len(ev.walkers)
	postPrograms := len(ev.programs)
	postAlerts := len(ev.parseAlerts)
	postContents := len(ev.fileContents)
	ev.mu.Unlock()

	if postFiles != preFiles-1 {
		t.Errorf("files: got %d, want %d", postFiles, preFiles-1)
	}
	if postList != preList-1 {
		t.Errorf("walkerList: got %d, want %d", postList, preList-1)
	}
	if postWalkers >= preWalkers {
		t.Errorf("walkers: got %d, want <%d (one source path + one abs path should be dropped)", postWalkers, preWalkers)
	}
	if postPrograms != prePrograms-1 {
		t.Errorf("programs: got %d, want %d", postPrograms, prePrograms-1)
	}
	if postAlerts != preAlerts-1 {
		t.Errorf("parseAlerts: got %d, want %d", postAlerts, preAlerts-1)
	}
	if postContents != preContents-1 {
		t.Errorf("fileContents: got %d, want %d", postContents, preContents-1)
	}
}

// TestRemoveFile_UnknownPathReturnsFalse asserts that calling
// RemoveFile with a path that doesn't match any known file is a
// no-op (returns false) and doesn't mutate the evaluator.
func TestRemoveFile_UnknownPathReturnsFalse(t *testing.T) {
	ev := NewEvaluator([]core.File{
		{DirectoryPath: "src", FileName: "foo", FileExtension: ".hyb"},
	})
	if err := ev.UpdateFileContent("src/foo.hyb", "env Foo as Level\n\ntick {}\n"); err != nil {
		t.Fatalf("UpdateFileContent: %v", err)
	}

	ev.mu.Lock()
	preFiles := len(ev.files)
	preList := len(ev.walkerList)
	ev.mu.Unlock()

	if ev.RemoveFile("does/not/exist.hyb") {
		t.Error("RemoveFile returned true for unknown path")
	}
	if ev.RemoveFile("") {
		t.Error("RemoveFile returned true for empty path")
	}

	ev.mu.Lock()
	postFiles := len(ev.files)
	postList := len(ev.walkerList)
	ev.mu.Unlock()

	if postFiles != preFiles {
		t.Errorf("files mutated: got %d, want %d", postFiles, preFiles)
	}
	if postList != preList {
		t.Errorf("walkerList mutated: got %d, want %d", postList, preList)
	}
}

// TestRemoveFile_HandlesPathVariants verifies that RemoveFile
// resolves different path forms to the same canonical file. The
// LSP can call RemoveFile with abs paths, relative paths, paths
// with ".." segments, or paths with extra slashes — all should
// hit the same file.
func TestRemoveFile_HandlesPathVariants(t *testing.T) {
	ev := NewEvaluator([]core.File{
		{DirectoryPath: "src", FileName: "foo", FileExtension: ".hyb"},
	})
	if err := ev.UpdateFileContent("src/foo.hyb", "env Foo as Level\n\ntick {}\n"); err != nil {
		t.Fatalf("UpdateFileContent: %v", err)
	}

	cases := []string{
		"src/foo.hyb",
		"./src/foo.hyb",
		"src//foo.hyb",
	}
	for _, p := range cases {
		if !ev.RemoveFile(p) {
			t.Errorf("RemoveFile(%q) returned false, want true (should match canonical src/foo.hyb)", p)
		}
		// Re-add so the next variant has something to find.
		if err := ev.UpdateFileContent("src/foo.hyb", "env Foo as Level\n\ntick {}\n"); err != nil {
			t.Fatalf("re-add: %v", err)
		}
	}
}
