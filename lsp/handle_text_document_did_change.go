package lsp

import (
	"context"
	"encoding/json"
	"hybroid/core"
	"hybroid/evaluator"
	"os"
	"path/filepath"
	"strings"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *langHandler) handleTextDocumentDidOpen(ctx context.Context, conn notifier, req *jsonrpc2.Request) (result any, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params DidOpenTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	// Tracks whether the opened file is a "stray" single-file open (no
	// workspace root and no discoverable hybconfig.toml above it). In that
	// case we publish a one-shot Information diagnostic so the user knows
	// that any unresolved `use` references are because the rest of the
	// project is not in scope — the same hint tsserver shows for files
	// outside a tsconfig.
	singleFileMode := false

	h.mu.Lock()
	h.files[params.TextDocument.URI] = &File{
		LanguageID: params.TextDocument.LanguageID,
		Text:       params.TextDocument.Text,
		Version:    params.TextDocument.Version,
	}

	// Try to discover a Hybroid project root (hybconfig.toml) above this
	// file. This is the fallback for single-file opens: the client did not
	// give us a workspace root, so we look for one ourselves, matching the
	// behavior of tsserver (tsconfig.json), Pylance (extraPaths), and
	// clangd (compile_commands.json).
	if h.rootPath == "" {
		if path, perr := fromURI(params.TextDocument.URI); perr == nil {
			if absPath, aerr := filepath.Abs(path); aerr == nil {
				if root := findProjectRoot(absPath, h.rootMarkers); root != "" {
					h.rootPath = filepath.Clean(root)
					h.addFolder(h.rootPath)
				}
			}
		}
	}

	if h.eval == nil {
		if h.rootPath != "" {
			// We found a project root via the parent-directory walk. Run
			// the pre-analysis synchronously so the first didOpen gets
			// full-workspace diagnostics immediately. The same code path
			// is used by handleInitialize for folder opens, just without
			// the goroutine.
			if filesInfo, ferr := core.CollectFiles(h.rootPath); ferr == nil {
				ev := evaluator.NewEvaluator(filesInfo)
				ev.ParseAll(h.rootPath)
				ev.RunAnalysis()
				h.eval = ev
				for _, info := range filesInfo {
					p := info.Path()
					uri := toURI(filepath.Join(h.rootPath, p))
					if content, rerr := os.ReadFile(filepath.Join(h.rootPath, p)); rerr == nil {
						h.files[uri] = &File{
							LanguageID: "hybroid",
							Text:       string(content),
							Version:    0,
						}
					}
				}
			}
		} else if path, perr := fromURI(params.TextDocument.URI); perr == nil {
			// True single-file mode: no workspace, no project marker.
			// Build an ad-hoc evaluator that only knows about the opened
			// file. Unresolved `use` statements will surface as
			// hyb035W warnings (truthful), and we publish a one-shot
			// Information diagnostic to explain why.
			baseName := filepath.Base(path)
			h.eval = evaluator.NewEvaluator([]core.File{{
				FileName:      strings.TrimSuffix(baseName, filepath.Ext(baseName)),
				DirectoryPath: ".",
				FileExtension: filepath.Ext(baseName),
			}})
			singleFileMode = true
		}
	}
	h.mu.Unlock()

	h.markReady()
	h.analyzeAndPublish(ctx, conn, params.TextDocument.URI, params.TextDocument.Text)

	if singleFileMode {
		h.publishInfoOnce(ctx, conn, params.TextDocument.URI,
			"This file is open without its Hybroid project. Open the folder containing hybconfig.toml to resolve all `use` references.")
	}

	return nil, nil
}

func (h *langHandler) handleTextDocumentDidChange(ctx context.Context, conn notifier, req *jsonrpc2.Request) (result any, err error) {
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
		// Since we use TDSKFull in initialize, we assume the last change
		// contains the full text. Version is updated unconditionally —
		// a didChange with an empty ContentChanges list (which the LSP
		// allows) must still advance the version, otherwise downstream
		// publishDiagnostics carries the old version and editors treat
		// the diagnostic as stale.
		if len(params.ContentChanges) > 0 {
			file.Text = params.ContentChanges[len(params.ContentChanges)-1].Text
		}
		file.Version = params.TextDocument.Version
		fileText = file.Text
	}
	h.mu.Unlock()

	if ok {
		h.scheduleAnalysis(params.TextDocument.URI, fileText)
	}

	return nil, nil
}

func (h *langHandler) analyzeAndPublish(ctx context.Context, conn notifier, uri DocumentURI, text string) {
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
