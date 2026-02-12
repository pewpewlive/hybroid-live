package lsp

import (
	"context"
	"encoding/json"
	"fmt"
	"hybroid/core"
	"hybroid/walker"
	"path/filepath"
	"strings"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *langHandler) handleTextDocumentHover(_ context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params HoverParams
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
	relPath, _ := filepath.Rel(h.rootPath, path)
	w := eval.AnalyzeFile(relPath)
	if w == nil {
		return nil, nil
	}

	// 1. Get the word under the cursor
	word := getWordAt(file.Text, params.Position.Line, params.Position.Character)
	core.DebugLog("Hover word at line %d, char %d: %q", params.Position.Line, params.Position.Character, word)
	if word == "" {
		return nil, nil
	}

	// 2. Check for metadata (keywords, builtins, namespaces, entities)
	detail, doc := getSymbolMetadata(w, eval.Walkers(), word)
	if detail != "" {
		res := Hover{
			Contents: MarkupContent{
				Kind:  Markdown,
				Value: fmt.Sprintf("**%s** (%s)\n\n%s", word, detail, doc),
			},
		}
		return res, nil
	}

	// 3. Check for variables or members in current scope
	line := params.Position.Line + 1
	col := params.Position.Character + 1
	scope := w.GetScopeAt(line, col)
	if scope != nil {
		// Handle member access hover (e.g. ship.x)
		if strings.Contains(word, ".") || strings.Contains(word, ":") {
			parts := strings.FieldsFunc(word, func(r rune) bool { return r == '.' || r == ':' })
			if len(parts) >= 2 {
				base := parts[0]
				if variable, found := scope.GetVariable(base); found {
					currentVal := variable.Value
					for i := 1; i < len(parts); i++ {
						member := parts[i]
						if container, ok := currentVal.(walker.FieldContainer); ok {
							if v, _, found := container.ContainsField(member); found {
								currentVal = v.Value
								if i == len(parts)-1 {
									return &Hover{
										Contents: MarkupContent{
											Kind:  Markdown,
											Value: fmt.Sprintf("**%s**: `%s`", word, currentVal.GetType().String()),
										},
									}, nil
								}
								continue
							}
						}
						if container, ok := currentVal.(walker.MethodContainer); ok {
							if v, found := container.ContainsMethod(member); found {
								currentVal = v.Value
								if i == len(parts)-1 {
									return &Hover{
										Contents: MarkupContent{
											Kind:  Markdown,
											Value: fmt.Sprintf("**%s**: `%s` (method)", word, currentVal.GetType().String()),
										},
									}, nil
								}
								continue
							}
						}
						break
					}
				}
			}
		}

		if variable, found := scope.GetVariable(word); found {
			typStr := variable.Value.GetType().String()
			res := Hover{
				Contents: MarkupContent{
					Kind:  Markdown,
					Value: fmt.Sprintf("**%s**: `%s`", word, typStr),
				},
			}
			return res, nil
		}
	}

	return nil, nil
}

func getWordAt(text string, line, character int) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(text, "\n")
	if line < 0 || line >= len(lines) {
		return ""
	}
	l := lines[line]
	if character < 0 || character >= len(l) {
		return ""
	}

	start := character
	for start > 0 && IsWordChar(rune(l[start-1])) {
		start--
	}
	end := character
	for end < len(l) && IsWordChar(rune(l[end])) {
		end++
	}

	if start == end {
		return ""
	}

	return l[start:end]
}
