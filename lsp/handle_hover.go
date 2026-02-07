package lsp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	w, ok := h.analyzedWalkers[params.TextDocument.URI]
	file, fileOk := h.files[params.TextDocument.URI]
	h.mu.Unlock()

	if !ok || !fileOk {
		log.Printf("Hover failed: walker_ok=%v, file_ok=%v for URI=%s", ok, fileOk, params.TextDocument.URI)
		return nil, nil
	}

	// 1. Get the word under the cursor
	word := getWordAt(file.Text, params.Position.Line, params.Position.Character)
	log.Printf("Hover word at line %d, char %d: %q", params.Position.Line, params.Position.Character, word)
	if word == "" {
		return nil, nil
	}

	// 2. Check for metadata (keywords, builtins, namespaces)
	detail, doc := getSymbolMetadata(word)
	if detail != "" {
		res := Hover{
			Contents: MarkupContent{
				Kind:  Markdown,
				Value: fmt.Sprintf("**%s** (%s)\n\n%s", word, detail, doc),
			},
		}
		return res, nil
	}

	// 3. Check for variables in current scope
	line := params.Position.Line + 1
	col := params.Position.Character + 1
	scope := w.GetScopeAt(line, col)
	if scope != nil {
		current := scope
		for current != nil {
			if variable, found := current.Variables[word]; found {
				typStr := variable.Value.GetType().String()
				res := Hover{
					Contents: MarkupContent{
						Kind:  Markdown,
						Value: fmt.Sprintf("**%s**: `%s`", word, typStr),
					},
				}
				return res, nil
			}
			current = current.Parent
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

	isWordChar := func(r rune) bool {
		return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
	}

	start := character
	for start > 0 && isWordChar(rune(l[start-1])) {
		start--
	}
	end := character
	for end < len(l) && isWordChar(rune(l[end])) {
		end++
	}

	if start == end {
		return ""
	}

	return l[start:end]
}
