package lsp

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strconv"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *langHandler) handleTextDocumentRename(_ context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params RenameParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	h.mu.Lock()
	eval := h.eval
	file, fileOk := h.files[params.TextDocument.URI]
	h.mu.Unlock()

	if eval == nil || !fileOk {
		return nil, nil
	}

	if isInCommentOrString(file.Text, params.Position.Line, params.Position.Character) {
		return nil, nil
	}

	path, _ := fromURI(params.TextDocument.URI)
	relPath := getRelPath(h.rootPath, path)
	h.evalMu.Lock()
	w := eval.AnalyzeFile(relPath)
	if w == nil {
		h.evalMu.Unlock()
		return nil, nil
	}
	defer h.evalMu.Unlock()

	word := getWordAt(file.Text, params.Position.Line, params.Position.Character)
	if word == "" {
		return nil, nil
	}

	newName := params.NewName
	if newName == "" || newName == word {
		return nil, nil
	}

	rootDir := h.rootPath
	if rootDir == "" {
		rootDir = filepath.Dir(path)
	}
	locations := h.findReferences(w, eval.Walkers(), eval.WalkerList(), word, params.Position.Line+1, params.Position.Character+1, true, rootDir)
	if len(locations) == 0 {
		return nil, nil
	}

	changes := make(map[DocumentURI][]TextEdit)
	seenEdits := make(map[string]bool)

	for _, loc := range locations {
		editKey := string(loc.URI) + ":" + strconv.Itoa(loc.Range.Start.Line) + ":" + strconv.Itoa(loc.Range.Start.Character)
		if seenEdits[editKey] {
			continue
		}
		seenEdits[editKey] = true

		changes[loc.URI] = append(changes[loc.URI], TextEdit{
			Range:   loc.Range,
			NewText: newName,
		})
	}

	return WorkspaceEdit{
		Changes: changes,
	}, nil
}
