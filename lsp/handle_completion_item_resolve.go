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
	// var detail string
	// var documentation string
	detail := "default detail"
	documentation := "default documentation"
	if item.Data == 1 {
		detail = "PewPew API"
		documentation = "API for PewPew levels"
	} else if item.Data == 2 {
		detail = "Fmath API"
		documentation = "API for fixed-point math"
	}
	return CompletionItem{
		Label:         item.Label,
		Kind:          item.Kind,
		Data:          item.Data,
		Detail:        detail,
		Documentation: documentation,
	}, nil
}
