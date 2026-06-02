package lsp

import (
	"context"
	"encoding/json"
	"hybroid/core"
	"hybroid/evaluator"
	"path/filepath"
	"strings"

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

	if h.eval == nil {
		path, err := fromURI(params.TextDocument.URI)
		if err == nil {
			baseName := filepath.Base(path)
			h.eval = evaluator.NewEvaluator([]core.FileInformation{{
				FileName:      strings.TrimSuffix(baseName, filepath.Ext(baseName)),
				DirectoryPath: ".",
				FileExtension: filepath.Ext(baseName),
			}})
		}
	}
	h.mu.Unlock()

	h.markReady()
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

	var fileText string
	h.mu.Lock()
	file, ok := h.files[params.TextDocument.URI]
	if ok {
		// Since we use TDSKFull in initialize, we assume the last change contains the full text
		if len(params.ContentChanges) > 0 {
			file.Text = params.ContentChanges[len(params.ContentChanges)-1].Text
			file.Version = params.TextDocument.Version
		}
		fileText = file.Text
	}
	h.mu.Unlock()

	if ok {
		h.scheduleAnalysis(params.TextDocument.URI, fileText)
	}

	return nil, nil
}

func (h *langHandler) analyzeAndPublish(ctx context.Context, conn *jsonrpc2.Conn, uri DocumentURI, text string) {
	path, err := fromURI(uri)
	if err != nil {
		return
	}

	h.mu.Lock()
	eval := h.eval
	var openFiles []struct {
		URI     DocumentURI
		Version int
	}
	for u, f := range h.files {
		openFiles = append(openFiles, struct {
			URI     DocumentURI
			Version int
		}{u, f.Version})
	}
	h.mu.Unlock()

	if eval == nil {
		return
	}

	relPath := getRelPath(h.rootPath, path)
	relPath = filepath.ToSlash(filepath.Clean(relPath))

	h.evalMu.Lock()
	eval.UpdateFileContent(relPath, text)
	eval.RunAnalysis()

	type diagInfo struct {
		uri     DocumentURI
		version int
		diags   []Diagnostic
	}
	diagBatch := make([]diagInfo, 0, len(openFiles))
	for _, info := range openFiles {
		p, ferr := fromURI(info.URI)
		if ferr != nil {
			continue
		}
		rPath := getRelPath(h.rootPath, p)
		rPath = filepath.ToSlash(filepath.Clean(rPath))
		diagBatch = append(diagBatch, diagInfo{
			uri:     info.URI,
			version: info.Version,
			diags:   alertsToDiagnostics(info.URI, eval.GetAlerts(rPath)),
		})
	}
	h.evalMu.Unlock()

	for _, info := range diagBatch {
		params := PublishDiagnosticsParams{
			URI:         info.uri,
			Diagnostics: info.diags,
		}
		version := info.version
		params.Version = &version
		conn.Notify(ctx, "textDocument/publishDiagnostics", params)
	}
}
