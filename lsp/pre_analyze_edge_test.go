package lsp

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/sourcegraph/jsonrpc2"
)

// TestPreAnalyzeWorkspace_RootPathMissing covers the case where
// h.rootPath was set (e.g. via handleInitialize) but the directory is
// missing at pre-analyze time (deleted between init and pre-analyze,
// or the user pointed at a typo'd path). The contract the LSP must
// uphold: the server must not hang waiting for ready. Even when the
// workspace is unreachable, subsequent requests should be handled —
// typically by falling back to single-file mode on didOpen.
//
// Note: in the current implementation, core.CollectFiles tolerates
// a missing root and returns an empty file list with no error, so
// preAnalyzeWorkspace ends up running the full pipeline on an
// empty workspace. The test pins that behavior — if a future refactor
// makes CollectFiles strict (error on missing root) and preAnalyzeWorkspace
// early-returns without markReady, this test catches the hang.
func TestPreAnalyzeWorkspace_RootPathMissing(t *testing.T) {
	h, _ := newTestHandler(t)

	// Point rootPath at a non-existent subdirectory of t.TempDir().
	ghost := filepath.Join(t.TempDir(), "deleted", "now")
	h.rootPath = ghost
	h.addFolder(ghost)

	if h.conn == nil {
		t.Fatal("expected h.conn to be set by newTestHandler")
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("preAnalyzeWorkspace panicked: %v", r)
		}
	}()

	// Reset the ready state so we can observe whether
	// preAnalyzeWorkspace itself calls markReady. (newTestHandler
	// already called it once; we replace the channel.)
	h.mu.Lock()
	h.ready = make(chan struct{})
	h.readySet = false
	h.mu.Unlock()

	h.preAnalyzeWorkspace()

	// Give the function a moment to finish (it shouldn't spawn
	// anything async, but be defensive).
	time.Sleep(50 * time.Millisecond)

	if !h.waitReady(context.Background()) {
		t.Errorf("waitReady returned false after preAnalyzeWorkspace on missing root — server would hang")
	}
}

// TestHandleDidChange_NilParams covers the trivial input validation:
// a request with no params. The handler must respond with
// CodeInvalidParams (per LSP spec) rather than nil-erroring and
// panicking later in json.Unmarshal.
func TestHandleDidChange_NilParams(t *testing.T) {
	h, _ := newTestHandler(t)
	req := &jsonrpc2.Request{
		Method: "textDocument/didChange",
		Params: nil,
	}
	_, err := h.handleTextDocumentDidChange(context.Background(), h.conn, req)
	if err == nil {
		t.Fatal("expected non-nil error for nil params")
	}
	rpcErr, ok := err.(*jsonrpc2.Error)
	if !ok {
		t.Fatalf("expected *jsonrpc2.Error, got %T: %v", err, err)
	}
	if rpcErr.Code != jsonrpc2.CodeInvalidParams {
		t.Errorf("expected CodeInvalidParams (%d), got %d", jsonrpc2.CodeInvalidParams, rpcErr.Code)
	}
}
