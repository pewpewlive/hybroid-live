package lsp

import (
	"context"
	"encoding/json"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *langHandler) handleTextDocumentCompletion(_ context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params CompletionParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	return h.completion(params.TextDocument.URI, &params)
}

func (h *langHandler) completion(uri DocumentURI, params *CompletionParams) ([]CompletionItem, error) {
	return []CompletionItem{
		{
			Label: "PewPew",
			Kind:  ClassCompletion,
			Data:  1,
		},
		{
			Label: "Fmath",
			Kind:  ClassCompletion,
			Data:  2,
		},
	}, nil
}
