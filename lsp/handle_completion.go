package lsp

import (
	"context"
	"encoding/json"
	"hybroid/walker"

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
	items := []CompletionItem{
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
	}

	h.mu.Lock()
	w, ok := h.analyzedWalkers[uri]
	h.mu.Unlock()

	if !ok {
		return items, nil
	}

	// LSP lines are 0-based, Hybroid tokens are 1-based.
	line := params.Position.Line + 1
	col := params.Position.Character + 1

	scope := w.GetScopeAt(line, col)
	if scope == nil {
		return items, nil
	}

	// Traverse scope chain
	seen := make(map[string]bool)
	current := scope
	for current != nil {
		for name, variable := range current.Variables {
			if !seen[name] {
				seen[name] = true
				kind := VariableCompletion
				
				// Attempt to deduce kind
				if _, ok := variable.Value.(*walker.FunctionVal); ok {
					kind = FunctionCompletion
				}

				items = append(items, CompletionItem{
					Label:  name,
					Kind:   kind,
					Detail: variable.Value.GetType().String(),
				})
			}
		}
		current = current.Parent
	}

	return items, nil
}
