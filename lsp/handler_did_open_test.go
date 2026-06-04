package lsp

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sourcegraph/jsonrpc2"
)

// pathHasNoProjectMarker walks up from dir and returns true if no
// hybconfig.toml is found anywhere up to the filesystem root. Tests that
// exercise the "no project" branch of handleTextDocumentDidOpen call this
// first to assert hermeticity — otherwise a stray marker from a parent
// test run or an unrelated project above /tmp would silently flip the
// branch being tested.
func pathHasNoProjectMarker(t *testing.T, dir string) {
	t.Helper()
	root, err := filepath.Abs(dir)
	if err != nil {
		t.Fatalf("Abs: %v", err)
	}
	if got := findProjectRoot(root, []string{"hybconfig.toml"}); got != "" {
		t.Skipf("skipping: ancestor of %s contains hybconfig.toml at %s — cannot exercise true single-file mode here", root, got)
	}
}

// TestHandleDidOpen_EmptyProjectNoFiles asserts that didOpen into a
// directory with a hybconfig.toml marker but no .hyb files doesn't panic
// and leaves the handler in a usable state. This was a previously-uncovered
// edge: preAnalyzeWorkspace iterates an empty file list and must still
// mark the handler ready so subsequent requests don't hang.
func TestHandleDidOpen_EmptyProjectNoFiles(t *testing.T) {
	dir := writeProject(t, map[string]string{
		"hybconfig.toml": minimalHybConfig,
	})
	uri := toURI(filepath.Join(dir, "level.hyb"))

	h, _ := newTestHandler(t)
	req := newTestRequest("textDocument/didOpen", DidOpenTextDocumentParams{
		TextDocument: TextDocumentItem{
			URI:        uri,
			LanguageID: "hybroid",
			Version:    0,
			Text:       minimalLevelSource,
		},
	})

	_, err := h.handleTextDocumentDidOpen(context.Background(), h.conn, req)
	if err != nil {
		t.Fatalf("handleTextDocumentDidOpen: %v", err)
	}

	if h.rootPath == "" {
		t.Errorf("expected h.rootPath to be set to %q", dir)
	}
	if h.eval == nil {
		t.Errorf("expected h.eval to be set even with 0 .hyb files")
	}
	// markReady should have been called; a subsequent waitReady must not block.
	if !h.waitReady(context.Background()) {
		t.Errorf("expected waitReady to return true immediately")
	}
}

// TestHandleDidOpen_EmptyTextSingleFile verifies that a didOpen with empty
// text in single-file mode doesn't crash and produces a publishDiagnostics
// notification (possibly with 0 diagnostics, but the notification must be
// sent so the editor clears any stale state).
func TestHandleDidOpen_EmptyTextSingleFile(t *testing.T) {
	dir := t.TempDir()
	pathHasNoProjectMarker(t, dir)
	uri := toURI(filepath.Join(dir, "level.hyb"))

	h, conn := newTestHandler(t)
	req := newTestRequest("textDocument/didOpen", DidOpenTextDocumentParams{
		TextDocument: TextDocumentItem{
			URI:        uri,
			LanguageID: "hybroid",
			Version:    0,
			Text:       "",
		},
	})

	_, err := h.handleTextDocumentDidOpen(context.Background(), h.conn, req)
	if err != nil {
		t.Fatalf("handleTextDocumentDidOpen: %v", err)
	}

	if conn.CountByMethod("textDocument/publishDiagnostics") == 0 {
		t.Errorf("expected at least one publishDiagnostics for empty text")
	}
}

// TestHandleDidOpen_EditAfterSingleFileMode verifies the single-file
// evaluator is re-used across didChange — i.e. opening in single-file mode
// establishes h.eval, and a subsequent didChange publishes fresh
// diagnostics for that URI only (not for any other open files).
func TestHandleDidOpen_EditAfterSingleFileMode(t *testing.T) {
	dir := t.TempDir()
	pathHasNoProjectMarker(t, dir)
	uri := toURI(filepath.Join(dir, "level.hyb"))

	h, conn := newTestHandler(t)
	openReq := newTestRequest("textDocument/didOpen", DidOpenTextDocumentParams{
		TextDocument: TextDocumentItem{
			URI:        uri,
			LanguageID: "hybroid",
			Version:    0,
			Text:       minimalLevelSource,
		},
	})
	if _, err := h.handleTextDocumentDidOpen(context.Background(), h.conn, openReq); err != nil {
		t.Fatalf("didOpen: %v", err)
	}
	// Reset notification counter so the next assertion is unambiguous.
	notifiesAfterOpen := conn.Count()

	changeReq := newTestRequest("textDocument/didChange", DidChangeTextDocumentParams{
		TextDocument: VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: TextDocumentIdentifier{URI: uri},
			Version:                1,
		},
		ContentChanges: []TextDocumentContentChangeEvent{
			{Text: "env TestLevel as Level\n\ntick {\n  let x = 1\n}\n"},
		},
	})
	if _, err := h.handleTextDocumentDidChange(context.Background(), h.conn, changeReq); err != nil {
		t.Fatalf("didChange: %v", err)
	}
	// Wait for debounce to fire. The timer is set in didChange; with
	// lintDebounce=1ms, the callback should run almost immediately. We
	// sleep 500ms for headroom on slow CI; nothing about the test
	// depends on the exact duration.
	time.Sleep(500 * time.Millisecond)
	got := conn.Count()
	if got <= notifiesAfterOpen {
		t.Errorf("expected at least one new publishDiagnostics after didChange, got %d total (was %d)", got, notifiesAfterOpen)
	}

	// All post-change notifications must be for the same URI.
	for _, c := range conn.Notifies()[notifiesAfterOpen:] {
		if c.Method != "textDocument/publishDiagnostics" {
			continue
		}
		p, ok := c.Params.(PublishDiagnosticsParams)
		if !ok {
			t.Fatalf("unexpected params type %T", c.Params)
		}
		if p.URI != uri {
			t.Errorf("didChange published diagnostics for unexpected URI %q (want %q)", p.URI, uri)
		}
	}
}

// TestHandleDidOpen_SingleFileNoProject_PublishesInfo is the regression
// test for the single-file fallback: a didOpen with no discoverable
// project root must publish a one-shot Information diagnostic AND
// normal error/warning diagnostics from the single-file evaluator.
func TestHandleDidOpen_SingleFileNoProject_PublishesInfo(t *testing.T) {
	dir := t.TempDir()
	pathHasNoProjectMarker(t, dir)
	uri := toURI(filepath.Join(dir, "level.hyb"))

	h, conn := newTestHandler(t)
	req := newTestRequest("textDocument/didOpen", DidOpenTextDocumentParams{
		TextDocument: TextDocumentItem{
			URI:        uri,
			LanguageID: "hybroid",
			Version:    0,
			// Use `use` of an unknown module so the walker produces
			// hyb035W — this proves the evaluator actually ran.
			Text: "env TestLevel as Level\n\nuse NoSuchModule\n\ntick {\n}\n",
		},
	})

	_, err := h.handleTextDocumentDidOpen(context.Background(), h.conn, req)
	if err != nil {
		t.Fatalf("handleTextDocumentDidOpen: %v", err)
	}

	if h.rootPath != "" {
		t.Errorf("expected h.rootPath empty in single-file mode, got %q", h.rootPath)
	}
	if h.eval == nil {
		t.Errorf("expected h.eval to be set even in single-file mode")
	}

	// Find the Information diagnostic for our URI.
	var infoDiag *Diagnostic
	for _, c := range conn.Notifies() {
		if c.Method != "textDocument/publishDiagnostics" {
			continue
		}
		p, ok := c.Params.(PublishDiagnosticsParams)
		if !ok || p.URI != uri {
			continue
		}
		for i := range p.Diagnostics {
			if p.Diagnostics[i].Severity == 3 {
				infoDiag = &p.Diagnostics[i]
				break
			}
		}
	}
	if infoDiag == nil {
		t.Fatalf("expected one Information-severity diagnostic for %q", uri)
	}
	if infoDiag.Source == nil || *infoDiag.Source != "hybroid" {
		t.Errorf("info diag source = %v, want pointer to \"hybroid\"", infoDiag.Source)
	}
	if !strings.Contains(infoDiag.Message, "hybconfig.toml") {
		t.Errorf("info diag message %q does not mention hybconfig.toml", infoDiag.Message)
	}
}

// TestHandleDidOpen_SingleFileInProject_DiscoversRoot verifies the parent
// walk: a file inside a project tree (hybconfig.toml in an ancestor) must
// set h.rootPath, run full pre-analysis, and NOT publish the Information
// notice.
func TestHandleDidOpen_SingleFileInProject_DiscoversRoot(t *testing.T) {
	root := writeProject(t, map[string]string{
		"hybconfig.toml":  minimalHybConfig,
		"level.hyb":      minimalLevelSource,
		"helpers/util.hyb": "env Helpers as Level\n",
	})
	uri := toURI(filepath.Join(root, "level.hyb"))

	h, conn := newTestHandler(t)
	req := newTestRequest("textDocument/didOpen", DidOpenTextDocumentParams{
		TextDocument: TextDocumentItem{
			URI:        uri,
			LanguageID: "hybroid",
			Version:    0,
			Text:       minimalLevelSource,
		},
	})

	_, err := h.handleTextDocumentDidOpen(context.Background(), h.conn, req)
	if err != nil {
		t.Fatalf("handleTextDocumentDidOpen: %v", err)
	}

	if h.rootPath == "" {
		t.Fatalf("expected h.rootPath to be set to %q", root)
	}
	if filepath.Clean(h.rootPath) != filepath.Clean(root) {
		t.Errorf("rootPath = %q, want %q", h.rootPath, root)
	}
	if h.eval == nil {
		t.Errorf("expected h.eval to be set after project pre-analysis")
	}

	// No Information diagnostic should have been published for our URI.
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
				t.Errorf("did not expect Information diagnostic in project mode, got %+v", d)
			}
		}
	}
}

// TestHandleDidOpen_RepeatedOpen_NoDuplicateInfo verifies the one-shot
// behavior of publishInfoOnce: a second didOpen of the same URI in
// single-file mode must NOT republish the Information notice.
func TestHandleDidOpen_RepeatedOpen_NoDuplicateInfo(t *testing.T) {
	dir := t.TempDir()
	pathHasNoProjectMarker(t, dir)
	uri := toURI(filepath.Join(dir, "level.hyb"))

	h, conn := newTestHandler(t)
	params := DidOpenTextDocumentParams{
		TextDocument: TextDocumentItem{
			URI:        uri,
			LanguageID: "hybroid",
			Version:    0,
			Text:       minimalLevelSource,
		},
	}

	if _, err := h.handleTextDocumentDidOpen(context.Background(), h.conn, newTestRequest("textDocument/didOpen", params)); err != nil {
		t.Fatalf("first didOpen: %v", err)
	}
	// Second open with same URI.
	if _, err := h.handleTextDocumentDidOpen(context.Background(), h.conn, newTestRequest("textDocument/didOpen", params)); err != nil {
		t.Fatalf("second didOpen: %v", err)
	}

	infoCount := 0
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
				infoCount++
			}
		}
	}
	if infoCount != 1 {
		t.Errorf("expected exactly 1 Information diagnostic across 2 didOpens, got %d", infoCount)
	}
}

// TestHandleDidOpen_FileContentsStored verifies that the file map records
// the open's text and version exactly as received, so later didChange
// handlers (which read h.files) see the right baseline.
func TestHandleDidOpen_FileContentsStored(t *testing.T) {
	dir := t.TempDir()
	pathHasNoProjectMarker(t, dir)
	uri := toURI(filepath.Join(dir, "level.hyb"))

	h, _ := newTestHandler(t)
	body := "env TestLevel as Level\n"
	req := newTestRequest("textDocument/didOpen", DidOpenTextDocumentParams{
		TextDocument: TextDocumentItem{
			URI:        uri,
			LanguageID: "hybroid",
			Version:    7,
			Text:       body,
		},
	})

	if _, err := h.handleTextDocumentDidOpen(context.Background(), h.conn, req); err != nil {
		t.Fatalf("handleTextDocumentDidOpen: %v", err)
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	f, ok := h.files[uri]
	if !ok {
		t.Fatalf("h.files[%q] not set", uri)
	}
	if f.Text != body {
		t.Errorf("stored text = %q, want %q", f.Text, body)
	}
	if f.Version != 7 {
		t.Errorf("stored version = %d, want 7", f.Version)
	}
	if f.LanguageID != "hybroid" {
		t.Errorf("stored languageId = %q, want %q", f.LanguageID, "hybroid")
	}
}

// keep jsonrpc2 import alive for tests that need its types.
var _ = jsonrpc2.CodeInvalidParams
