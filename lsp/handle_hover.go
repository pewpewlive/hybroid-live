package lsp

import (
	"context"
	"encoding/json"
	"fmt"
	"hybroid/core"
	"hybroid/walker"
	"math"
	"strconv"
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
	relPath := getRelPath(h.rootPath, path)
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

	// 1.5. Check for numeric literal hover (e.g. 90d, 10.5f -> show computed fixed-point value)
	numLit := getNumericLiteralAt(file.Text, params.Position.Line, params.Position.Character)
	if len(numLit) > 1 {
		suffix := numLit[len(numLit)-1]
		numStr := numLit[:len(numLit)-1]
		if val, err := strconv.ParseFloat(numStr, 64); err == nil {
			switch suffix {
			case 'd':
				rad := val * math.Pi / 180
				fxVal := fixedToFxStr(rad)
				return &Hover{
					Contents: MarkupContent{
						Kind:  Markdown,
						Value: fmt.Sprintf("`%s` = `%sfx`", numLit, fxVal),
					},
				}, nil
			case 'f', 'r':
				fxVal := fixedToFxStr(val)
				return &Hover{
					Contents: MarkupContent{
						Kind:  Markdown,
						Value: fmt.Sprintf("`%s` = `%sfx`", numLit, fxVal),
					},
				}, nil
			}
		}
	}

	// 2. Check for metadata (keywords, builtins, namespaces, entities)
	detail, doc := getSymbolMetadata(w, eval.Walkers(), word)
	if detail != "" {
		display := fmt.Sprintf("**%s**", word)
		if doc != "" {
			// If doc is a namespace (single word, no spaces, starts with uppercase)
			// it's likely a namespace returned for non-prefixed symbols.
			// This is a bit of a hack since we are reusing the doc field.
			if !strings.Contains(doc, " ") && len(doc) > 0 && doc[0] >= 'A' && doc[0] <= 'Z' {
				display = fmt.Sprintf("**%s** (%s)", word, doc)
				doc = "" // Clear it so it doesn't show as description
			}
		}

		value := display + " (" + detail + ")"
		if doc != "" {
			value += "\n\n" + doc
		}

		res := Hover{
			Contents: MarkupContent{
				Kind:  Markdown,
				Value: value,
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

// getNumericLiteralAt extracts a numeric literal token at the given position,
// including decimal points (e.g. "10.5f", "90d", "3.14r").
func getNumericLiteralAt(text string, line, character int) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(text, "\n")
	if line < 0 || line >= len(lines) {
		return ""
	}
	l := lines[line]
	if character < 0 || character >= len(l) {
		return ""
	}

	isNumLitChar := func(r rune) bool {
		return (r >= '0' && r <= '9') || r == '.' || r == '_'
	}

	// Scan forward to find the suffix character (d, f, r)
	end := character
	for end < len(l) && (isNumLitChar(rune(l[end])) || (l[end] >= 'a' && l[end] <= 'z')) {
		end++
	}
	// Scan backward over digits and dots
	start := character
	for start > 0 && isNumLitChar(rune(l[start-1])) {
		start--
	}
	// Include leading minus sign for negative literals
	if start > 0 && l[start-1] == '-' {
		start--
	}

	if start == end {
		return ""
	}

	return l[start:end]
}

// fixedToFxStr converts a float64 to its fixed-point string representation.
func fixedToFxStr(f float64) string {
	absF := math.Abs(f)
	integer := math.Min(math.Floor(absF), float64(int64(2)<<51))
	var sign string
	if f < 0 {
		sign = "-"
	}

	frac := math.Floor((absF - integer) * 4096)
	fracStr := ""
	if frac != 0 {
		fracStr = "." + strconv.FormatFloat(frac, 'f', -1, 64)
	}

	return fmt.Sprintf("%s%v%s", sign, integer, fracStr)
}
