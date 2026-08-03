package lsp

import (
	"context"
	"encoding/json"
	"os"
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

	path, pathErr := fromURI(params.TextDocument.URI)
	h.mu.Lock()
	projectMode := h.rootPath != ""
	rootPath := h.rootPath
	if !projectMode {
		delete(h.files, params.TextDocument.URI)
	}
	h.mu.Unlock()

	if projectMode && pathErr == nil {
		// Closing an editor tab does not remove a file from a Hybroid project.
		// Restore the disk-backed contents so unsaved buffer changes are
		// discarded while the file's environment remains available to every
		// other project file.
		if content, readErr := os.ReadFile(path); readErr == nil {
			h.mu.Lock()
			h.files[params.TextDocument.URI] = &File{
				LanguageID: "hybroid",
				Text:       string(content),
				Version:    0,
			}
			h.mu.Unlock()
			h.analyzeAndPublish(ctx, conn, params.TextDocument.URI, string(content))
			return nil, nil
		}

		// If the project file no longer exists on disk, this is a real removal
		// rather than a preview-tab close. Fall through to evaluator cleanup.
		h.mu.Lock()
		delete(h.files, params.TextDocument.URI)
		h.mu.Unlock()
	}

	// Drop the file's per-file state from the evaluator so its
	// walker, AST, and alerts are released. Without this, single-file
	// mode grows the evaluator's internal maps (walkers, walkerList,
	// files, programs, parseAlerts, fileContents) on every distinct
	// open in a long-running server.
	h.evalMu.Lock()
	if h.eval != nil {
		if pathErr == nil {
			relPath := getRelPath(rootPath, path)
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
