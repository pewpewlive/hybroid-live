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
	items := make([]CompletionItem, 0)
	seen := make(map[string]bool)

	// 1. Keywords
	keywords := []string{
		"is", "isnt", "alias", "and", "as", "break", "by", "const", "continue",
		"else", "entity", "enum", "env", "false", "fn", "to", "for", "if", "in",
		"let", "match", "new", "or", "pub", "repeat", "return", "self", "spawn",
		"struct", "class", "tick", "true", "use", "from", "while", "with",
		"yield", "destroy", "every",
	}
	for _, kw := range keywords {
		if !seen[kw] {
			items = append(items, CompletionItem{
				Label: kw,
				Kind:  KeywordCompletion,
			})
			seen[kw] = true
		}
	}

	// 2. Native Types
	nativeTypes := []string{
		"number", "fixed", "text", "map", "list", "bool", "struct", "entity",
	}
	for _, nt := range nativeTypes {
		if !seen[nt] {
			items = append(items, CompletionItem{
				Label: nt,
				Kind:  TypeParameterCompletion,
			})
			seen[nt] = true
		}
	}

	// 3. Namespaces
	namespaces := []string{"Pewpew", "Fmath", "Math", "String", "Table"}
	for _, ns := range namespaces {
		if !seen[ns] {
			items = append(items, CompletionItem{
				Label: ns,
				Kind:  ModuleCompletion,
			})
			seen[ns] = true
		}
	}

	// 4. Environments
	environments := []string{"Level", "Mesh", "Sound", "Shared"}
	for _, env := range environments {
		if !seen[env] {
			items = append(items, CompletionItem{
				Label: env,
				Kind:  ConstantCompletion,
			})
			seen[env] = true
		}
	}

	// 5. Builtin Functions (from BuiltinEnv)
	for name, variable := range walker.BuiltinEnv.Scope.Variables {
		if !seen[name] {
			kind := VariableCompletion
			if _, ok := variable.Value.(*walker.FunctionVal); ok {
				kind = FunctionCompletion
			}
			items = append(items, CompletionItem{
				Label:  name,
				Kind:   kind,
				Detail: variable.Value.GetType().String(),
			})
			seen[name] = true
		}
	}

	// 5. Pewpew Symbols
	for name, variable := range walker.PewpewAPI.Scope.Variables {
		if !seen[name] {
			kind := VariableCompletion
			if _, ok := variable.Value.(*walker.FunctionVal); ok {
				kind = FunctionCompletion
			}
			items = append(items, CompletionItem{
				Label:  name,
				Kind:   kind,
				Detail: "Pewpew:" + variable.Value.GetType().String(),
			})
			seen[name] = true
		}
	}

	// 6. Fmath Symbols
	for name, variable := range walker.FmathAPI.Scope.Variables {
		if !seen[name] {
			kind := VariableCompletion
			if _, ok := variable.Value.(*walker.FunctionVal); ok {
				kind = FunctionCompletion
			}
			items = append(items, CompletionItem{
				Label:  name,
				Kind:   kind,
				Detail: "Fmath:" + variable.Value.GetType().String(),
			})
			seen[name] = true
		}
	}

	h.mu.Lock()
	w, ok := h.analyzedWalkers[uri]
	h.mu.Unlock()

	if !ok {
		return items, nil
	}

	// 5. Scope-specific completions
	line := params.Position.Line + 1
	col := params.Position.Character + 1

	scope := w.GetScopeAt(line, col)
	if scope == nil {
		return items, nil
	}

	current := scope
	for current != nil {
		for name, variable := range current.Variables {
			if !seen[name] {
				seen[name] = true
				kind := VariableCompletion
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
