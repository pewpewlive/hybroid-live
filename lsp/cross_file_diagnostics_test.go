package lsp

import (
	"context"
	"testing"
	"time"

	"github.com/sourcegraph/jsonrpc2"
)

// openForTest is a tiny helper that wraps the boilerplate of issuing a
// didOpen for a single file and waiting for analyzeAndPublish to complete.
// It returns the URI and the count of notifies after open (so callers can
// assert on the delta after a follow-up action).
func openForTest(t *testing.T, h *langHandler, conn *fakeNotify, uri DocumentURI, body string) int {
	t.Helper()
	req := newTestRequest("textDocument/didOpen", DidOpenTextDocumentParams{
		TextDocument: TextDocumentItem{
			URI:        uri,
			LanguageID: "hybroid",
			Version:    0,
			Text:       body,
		},
	})
	if _, err := h.handleTextDocumentDidOpen(context.Background(), h.conn, req); err != nil {
		t.Fatalf("didOpen: %v", err)
	}
	return conn.Count()
}

// TestAnalyzeAndPublish_RePublishesAllOpenFiles locks in the production
// behavior that a didChange to one file re-publishes diagnostics for
// every open file. This is what keeps editors from showing stale
// diagnostics when a change in one file affects the type resolution of
// another (e.g. a `use` reference, a class redefinition).
//
// If a future refactor narrows the publish loop to "only the changed
// URI", this test will fail and force a conscious decision.
func TestAnalyzeAndPublish_RePublishesAllOpenFiles(t *testing.T) {
	h, conn := newTestHandler(t)

	uriA := DocumentURI("file:///a.hyb")
	uriB := DocumentURI("file:///b.hyb")
	uriC := DocumentURI("file:///c.hyb")

	openForTest(t, h, conn, uriA, "env A as Level\n")
	openForTest(t, h, conn, uriB, "env B as Level\n")
	openForTest(t, h, conn, uriC, "env C as Level\n")
	baseline := conn.Count()

	// Now edit A. After the debounce, every open file should get a fresh
	// publishDiagnostics with the new version.
	changeReq := newTestRequest("textDocument/didChange", DidChangeTextDocumentParams{
		TextDocument: VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: TextDocumentIdentifier{URI: uriA},
			Version:                1,
		},
		ContentChanges: []TextDocumentContentChangeEvent{
			{Text: "env A as Level\n// edited\n"},
		},
	})
	if _, err := h.handleTextDocumentDidChange(context.Background(), h.conn, changeReq); err != nil {
		t.Fatalf("didChange: %v", err)
	}

	// Wait for the debounced analysis to complete and publish.
	time.Sleep(500 * time.Millisecond)
	got := conn.Count()
	if got < baseline+3 {
		t.Fatalf("expected at least 3 new publishes (one per open file), got %d new", got-baseline)
	}

	// The 3 most recent publishes must cover all three URIs, in some order.
	published := map[DocumentURI]bool{}
	for _, c := range conn.Notifies()[baseline:] {
		if c.Method != "textDocument/publishDiagnostics" {
			continue
		}
		p, ok := c.Params.(PublishDiagnosticsParams)
		if !ok {
			continue
		}
		published[p.URI] = true
	}
	for _, want := range []DocumentURI{uriA, uriB, uriC} {
		if !published[want] {
			t.Errorf("expected a publishDiagnostics for %q after editing A, got: %v", want, published)
		}
	}
}

// TestDidClose_ClearsDiagnosticsWithNullVersion verifies the LSP-spec
// shape of the close notification: `diagnostics: []` AND no `version`
// field. Editors (VS Code) only clear stale diagnostics when the version
// is omitted — a stale version would leave the squiggle visible.
func TestDidClose_ClearsDiagnosticsWithNullVersion(t *testing.T) {
	h, conn := newTestHandler(t)
	uri := DocumentURI("file:///x.hyb")

	openForTest(t, h, conn, uri, "env X as Level\n")
	baseline := conn.Count()

	closeReq := newTestRequest("textDocument/didClose", DidCloseTextDocumentParams{
		TextDocument: TextDocumentIdentifier{URI: uri},
	})
	if _, err := h.handleTextDocumentDidClose(context.Background(), h.conn, closeReq); err != nil {
		t.Fatalf("didClose: %v", err)
	}

	// The close notification must be the last one, with empty diagnostics
	// and no version field. The LSP spec is explicit: when a file closes,
	// the server SHOULD publish an empty list to clear the editor state.
	nots := conn.Notifies()
	if len(nots) == 0 {
		t.Fatal("expected at least one notification (the close clear)")
	}
	last := nots[len(nots)-1]
	if last.Method != "textDocument/publishDiagnostics" {
		t.Fatalf("last notification was %q, want publishDiagnostics", last.Method)
	}
	p, ok := last.Params.(PublishDiagnosticsParams)
	if !ok {
		t.Fatalf("last notification params type = %T", last.Params)
	}
	if p.URI != uri {
		t.Errorf("cleared URI = %q, want %q", p.URI, uri)
	}
	if len(p.Diagnostics) != 0 {
		t.Errorf("cleared diagnostics = %v, want []", p.Diagnostics)
	}
	if p.Version != nil {
		t.Errorf("cleared version = %v, want nil (omitted)", p.Version)
	}
	_ = baseline
}

// TestDidChange_UnknownURI_DoesNotPanic verifies the defensive path:
// didChange for a URI that was never opened must not crash the server.
// In production this can happen if a client sends stale notifications
// after a file is closed, or if the editor's state diverges from the
// server's.
func TestDidChange_UnknownURI_DoesNotPanic(t *testing.T) {
	h, conn := newTestHandler(t)
	unknown := DocumentURI("file:///never-opened.hyb")

	req := newTestRequest("textDocument/didChange", DidChangeTextDocumentParams{
		TextDocument: VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: TextDocumentIdentifier{URI: unknown},
			Version:                1,
		},
		ContentChanges: []TextDocumentContentChangeEvent{
			{Text: "anything"},
		},
	})
	// Must not panic, must not return an error.
	res, err := h.handleTextDocumentDidChange(context.Background(), h.conn, req)
	if err != nil {
		t.Errorf("didChange for unknown URI returned error: %v", err)
	}
	if res != nil {
		t.Errorf("didChange for unknown URI returned non-nil result: %v", res)
	}
	// Crucially, no publishDiagnostics should be sent for an unknown URI.
	for _, c := range conn.Notifies() {
		if c.Method != "textDocument/publishDiagnostics" {
			continue
		}
		p, ok := c.Params.(PublishDiagnosticsParams)
		if !ok {
			continue
		}
		if p.URI == unknown {
			t.Errorf("unexpected publishDiagnostics for never-opened URI %q", unknown)
		}
	}
}

// TestHandleInitialize_NullRootUri verifies that the initialize handler
// does not panic when the client sends rootUri=null (some clients do
// this when no folder is open in the editor). The current production
// code is null-safe because params.RootURI is a DocumentURI (string)
// and the "" branch skips the URI parse.
func TestHandleInitialize_NullRootUri(t *testing.T) {
	h, conn := newTestHandler(t)

	req := newTestRequest("initialize", InitializeParams{
		ProcessID: 1234,
		// RootURI is intentionally the zero value (empty string),
		// which is what json.Unmarshal produces for a missing/null
		// rootUri field.
	})
	res, err := h.handleInitialize(context.Background(), h.conn, req)
	if err != nil {
		t.Fatalf("handleInitialize: %v", err)
	}
	if h.rootPath != "" {
		t.Errorf("expected rootPath empty when RootURI missing, got %q", h.rootPath)
	}
	// The capabilities response should still come back; clients rely on
	// the result being non-nil to proceed.
	if res == nil {
		t.Errorf("expected non-nil InitializeResult")
	}
	// No notifies should have been sent during initialize.
	if conn.Count() != 0 {
		t.Errorf("expected 0 notifies during initialize, got %d", conn.Count())
	}
}

// TestHandleDidChange_EmptyContentChanges verifies the degenerate input:
// a didChange with no actual content changes. The handler must not
// panic, must not crash the timer, and the produced publish (if any)
// must be for the requested URI and carry a non-nil version.
//
// Note: as of writing, the handler still re-publishes diagnostics even
// when ContentChanges is empty (the schedule path is unconditional once
// the URI is in h.files). We document that behavior here so a future
// "optimization" that breaks the publish is caught.
func TestHandleDidChange_EmptyContentChanges(t *testing.T) {
	h, conn := newTestHandler(t)
	uri := DocumentURI("file:///x.hyb")
	openForTest(t, h, conn, uri, "env X as Level\n")
	baseline := conn.Count()

	req := newTestRequest("textDocument/didChange", DidChangeTextDocumentParams{
		TextDocument: VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: TextDocumentIdentifier{URI: uri},
			Version:                1,
		},
		ContentChanges: []TextDocumentContentChangeEvent{},
	})
	if _, err := h.handleTextDocumentDidChange(context.Background(), h.conn, req); err != nil {
		t.Fatalf("didChange: %v", err)
	}
	// Wait long enough for the debounced analysis to fire.
	time.Sleep(500 * time.Millisecond)
	if conn.Count() <= baseline {
		t.Fatalf("expected at least one publishDiagnostics for didChange even with empty changes, got 0")
	}
	// The publish must be for our URI with the new version.
	var found bool
	for _, c := range conn.Notifies()[baseline:] {
		if c.Method != "textDocument/publishDiagnostics" {
			continue
		}
		p, ok := c.Params.(PublishDiagnosticsParams)
		if !ok || p.URI != uri {
			continue
		}
		if p.Version == nil || *p.Version != 1 {
			t.Errorf("expected publish version=1, got %v", p.Version)
		}
		found = true
	}
	if !found {
		t.Errorf("expected a publishDiagnostics for %q after didChange, got none", uri)
	}
}

// keep jsonrpc2 import alive for tests that need its types.
var _ = jsonrpc2.CodeInvalidParams
