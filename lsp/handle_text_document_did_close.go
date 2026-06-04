package lsp

import (
	"context"
	"encoding/json"
	"path/filepath"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *langHandler) handleTextDocumentDidClose(ctx context.Context, conn notifier, req *jsonrpc2.Request) (result any, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params DidCloseTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	h.mu.Lock()
	delete(h.files, params.TextDocument.URI)
	h.mu.Unlock()

	// Drop the file's per-file state from the evaluator so its
	// walker, AST, and alerts are released. Without this, single-file
	// mode grows the evaluator's internal maps (walkers, walkerList,
	// files, programs, parseAlerts, fileContents) on every distinct
	// open in a long-running server.
	h.evalMu.Lock()
	if h.eval != nil {
		if p, perr := fromURI(params.TextDocument.URI); perr == nil {
			relPath := getRelPath(h.rootPath, p)
			relPath = filepath.ToSlash(filepath.Clean(relPath))
			h.eval.RemoveFile(relPath)
		}
	}
	h.evalMu.Unlock()

	conn.Notify(ctx, "textDocument/publishDiagnostics", PublishDiagnosticsParams{
		URI:         params.TextDocument.URI,
		Diagnostics: []Diagnostic{},
	})

	return nil, nil
}
