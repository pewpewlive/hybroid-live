package lsp

import (
	"context"
	"encoding/json"
	"path/filepath"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *langHandler) handleTextDocumentDidOpen(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params DidOpenTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	h.mu.Lock()
	h.files[params.TextDocument.URI] = &File{
		LanguageID: params.TextDocument.LanguageID,
		Text:       params.TextDocument.Text,
		Version:    params.TextDocument.Version,
	}
	h.mu.Unlock()

	h.analyzeAndPublish(ctx, conn, params.TextDocument.URI, params.TextDocument.Text)

	return nil, nil
}

func (h *langHandler) handleTextDocumentDidChange(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params DidChangeTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	h.mu.Lock()
	file, ok := h.files[params.TextDocument.URI]
	if ok {
		// Since we use TDSKFull in initialize, we assume the last change contains the full text
		if len(params.ContentChanges) > 0 {
			file.Text = params.ContentChanges[len(params.ContentChanges)-1].Text
			file.Version = params.TextDocument.Version
		}
	}
	h.mu.Unlock()

	if file != nil {
		h.analyzeAndPublish(ctx, conn, params.TextDocument.URI, file.Text)
	}

	return nil, nil
}

func (h *langHandler) analyzeAndPublish(ctx context.Context, conn *jsonrpc2.Conn, uri DocumentURI, text string) {
	path, _ := fromURI(uri)

	h.mu.Lock()
	eval := h.eval
	h.mu.Unlock()

	if eval == nil {
		return
	}

	// 1. Update content in memory
	relPath, err := filepath.Rel(h.rootPath, path)
	if err != nil {
		relPath = path
	}
	eval.UpdateFileContent(relPath, text)

	// 2. Run project-wide analysis
	eval.RunAnalysis()

	// 3. Publish Diagnostics
	params := PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: alertsToDiagnostics(uri, eval.GetAlerts(relPath)),
	}

	conn.Notify(ctx, "textDocument/publishDiagnostics", params)
}
