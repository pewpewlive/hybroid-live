package lsp

import (
	"context"
	"encoding/json"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *langHandler) HandleCompletionItemResolve(_ context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params CompletionItem
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	return h.completionResolve(&params)
}

func (h *langHandler) completionResolve(item *CompletionItem) (CompletionItem, error) {
	detail, documentation := getSymbolMetadata(item.Label)

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
