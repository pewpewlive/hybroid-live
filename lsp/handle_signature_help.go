package lsp

import (
	"context"
	"encoding/json"
	"fmt"
	"hybroid/walker"
	"strings"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *langHandler) handleTextDocumentSignatureHelp(ctx context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params SignatureHelpParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	if !h.waitReady(ctx) {
		return nil, nil
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

	path, err := fromURI(params.TextDocument.URI)
	if err != nil {
		return nil, nil
	}
	relPath := getRelPath(h.rootPath, path)
	h.evalMu.Lock()
	w := eval.AnalyzeFile(relPath)
	if w == nil {
		h.evalMu.Unlock()
		return nil, nil
	}
	defer h.evalMu.Unlock()

	funcName, activeParam := findCallContext(file.Text, params.Position.Line, params.Position.Character)
	if funcName == "" {
		return nil, nil
	}

	var fnVal *walker.FunctionVal

	line := params.Position.Line + 1
	col := params.Position.Character + 1
	scope := w.GetScopeAt(line, col)
	if scope != nil {
		current := scope
		for current != nil {
			if v, ok := current.Variables[funcName]; ok {
				if f, ok := v.Value.(*walker.FunctionVal); ok {
					fnVal = f
					break
				}
			}
			current = current.Parent
		}
	}

	if fnVal == nil {
		lookupName := funcName
		var env *walker.Environment

		if strings.Contains(funcName, ":") || strings.Contains(funcName, ".") {
			parts := strings.FieldsFunc(funcName, func(r rune) bool { return r == ':' || r == '.' })
			if len(parts) == 2 {
				ns := parts[0]
				lookupName = parts[1]
				env = resolveBuiltinEnvByName(ns)

				if env == nil {
					if w2, ok := eval.Walkers()[ns]; ok {
						env = w2.Env()
					}
				}
			}
		}

		if env != nil {
			if v, ok := env.Scope.Variables[lookupName]; ok {
				fnVal, _ = v.Value.(*walker.FunctionVal)
			}
		} else {
			// Check builtins
			if v, ok := walker.BuiltinEnv.Scope.Variables[lookupName]; ok {
				fnVal, _ = v.Value.(*walker.FunctionVal)
			}

			// Check current walker's context if available
			if fnVal == nil {
				env := w.Env()

				// 1. Check imports (ThroughUse)
				for _, imp := range env.Imports() {
					if imp.ThroughUse {
						if v, ok := imp.Env().Scope.Variables[lookupName]; ok && v.IsPub {
							if f, ok := v.Value.(*walker.FunctionVal); ok {
								fnVal = f
								break
							}
						}
					}
				}

				// 2. Check libraries (only those explicitly imported via 'use')
				if fnVal == nil {
					for _, lib := range env.ImportedLibraries {
						libEnv := walker.BuiltinLibraries[lib]
						if v, ok := libEnv.Scope.Variables[lookupName]; ok {
							if f, ok := v.Value.(*walker.FunctionVal); ok {
								fnVal = f
								break
							}
						}
					}
				}
			}
		}
	}

	if fnVal == nil {
		return nil, nil
	}

	labels := make([]string, len(fnVal.Params))
	paramsInfo := make([]ParameterInformation, len(fnVal.Params))
	for i, p := range fnVal.Params {
		if i < len(fnVal.ParamNames) && fnVal.ParamNames[i] != "" {
			labels[i] = fmt.Sprintf("%s %s", p.String(), fnVal.ParamNames[i])
		} else {
			labels[i] = fmt.Sprintf("param%d: %s", i+1, p.String())
		}
		paramsInfo[i] = ParameterInformation{
			Label: labels[i],
		}
	}

	signatureLabel := fmt.Sprintf("%s(%s)", funcName, strings.Join(labels, ", "))

	res := SignatureHelp{
		Signatures: []SignatureInformation{
			{
				Label:      signatureLabel,
				Parameters: paramsInfo,
			},
		},
		ActiveSignature: 0,
		ActiveParameter: activeParam,
	}

	return res, nil
}

func findCallContext(text string, line, character int) (string, int) {
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	if line < 0 || line >= len(lines) {
		return "", 0
	}

	runes := []rune(lines[line])
	if character > len(runes) {
		character = len(runes)
	}

	// Scan backwards from the cursor, tracking paren depth, to find the
	// opening paren of the call we are currently inside.
	depth := 0
	openParenIdx := -1
	for i := character - 1; i >= 0; i-- {
		c := runes[i]
		if c == ')' {
			depth++
		} else if c == '(' {
			if depth == 0 {
				openParenIdx = i
				break
			}
			depth--
		}
	}
	if openParenIdx == -1 {
		return "", 0
	}

	// Count commas at depth 0 between the opening paren and the cursor to
	// determine the active parameter.
	activeParam := 0
	d := 0
	for i := openParenIdx; i < character; i++ {
		c := runes[i]
		if c == '(' {
			d++
		} else if c == ')' {
			d--
		} else if c == ',' && d == 1 {
			activeParam++
		}
	}

	nameEnd := openParenIdx
	for nameEnd > 0 && (runes[nameEnd-1] == ' ' || runes[nameEnd-1] == '\t') {
		nameEnd--
	}

	nameStart := nameEnd
	for nameStart > 0 && IsWordChar(runes[nameStart-1]) {
		nameStart--
	}

	if nameStart == nameEnd {
		return "", 0
	}

	return string(runes[nameStart:nameEnd]), activeParam
}
