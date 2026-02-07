package lsp

import (
	"context"
	"encoding/json"

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
	result := Analyze(text)

	h.mu.Lock()
	h.analyzedWalkers[uri] = result.Walker
	h.mu.Unlock()

	// Publish Diagnostics
	params := PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: result.Diagnostics,
	}

	conn.Notify(ctx, "textDocument/publishDiagnostics", params)
}
