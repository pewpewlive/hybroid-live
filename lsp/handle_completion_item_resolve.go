package lsp

import (
	"context"
	"encoding/json"
	"hybroid/walker"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *langHandler) HandleCompletionItemResolve(ctx context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params CompletionItem
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	if !h.waitReady(ctx) {
		return params, nil
	}

	return h.completionResolve(&params)
}

func (h *langHandler) completionResolve(item *CompletionItem) (CompletionItem, error) {
	h.mu.Lock()
	eval := h.eval
	h.mu.Unlock()

	var walkers map[string]*walker.Walker
	var w *walker.Walker
	if eval != nil {
		h.evalMu.Lock()
		defer h.evalMu.Unlock()
		walkers = eval.Walkers()
		if item.Data != nil {
			if uri, ok := item.Data.(string); ok {
				path, ferr := fromURI(DocumentURI(uri))
				if ferr == nil {
					relPath := getRelPath(h.rootPath, path)
					w = eval.AnalyzeFile(relPath)
				}
			}
		}
	}

	detail, documentation := getSymbolMetadata(w, walkers, item.Label)

	if detail == "" {
		detail = item.Detail
	}
	if documentation == "" {
		documentation = item.Documentation
	}

	return CompletionItem{
		Label:         item.Label,
		Kind:          item.Kind,
		Tags:          item.Tags,
		Detail:        detail,
		Documentation: documentation,
		Deprecated:    item.Deprecated,
		Preselect:     item.Preselect,
		SortText:      item.SortText,
		FilterText:    item.FilterText,
		InsertText:    item.InsertText,
		TextEdit:      item.TextEdit,
		Data:          item.Data,
	}, nil
}
