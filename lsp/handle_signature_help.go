package lsp

import (
	"context"
	"encoding/json"
	"fmt"
	"hybroid/walker"
	"strings"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *langHandler) handleTextDocumentSignatureHelp(_ context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params SignatureHelpParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	h.mu.Lock()
	w, ok := h.analyzedWalkers[params.TextDocument.URI]
	file, fileOk := h.files[params.TextDocument.URI]
	h.mu.Unlock()

	if !ok || !fileOk {
		return nil, nil
	}

	if isInCommentOrString(file.Text, params.Position.Line, params.Position.Character) {
		return nil, nil
	}

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
				switch ns {
				case "Pewpew":
					env = walker.PewpewAPI
				case "Fmath":
					env = walker.FmathAPI
				}
			}
		}

		if env != nil {
			if v, ok := env.Scope.Variables[lookupName]; ok {
				fnVal, _ = v.Value.(*walker.FunctionVal)
			}
		} else {
			if v, ok := walker.BuiltinEnv.Scope.Variables[lookupName]; ok {
				fnVal, _ = v.Value.(*walker.FunctionVal)
			} else if v, ok := walker.PewpewAPI.Scope.Variables[lookupName]; ok {
				fnVal, _ = v.Value.(*walker.FunctionVal)
			} else if v, ok := walker.FmathAPI.Scope.Variables[lookupName]; ok {
				fnVal, _ = v.Value.(*walker.FunctionVal)
			}
		}
	}

	if fnVal == nil {
		return nil, nil
	}

	labels := make([]string, len(fnVal.Params))
	paramsInfo := make([]ParameterInformation, len(fnVal.Params))
	for i, p := range fnVal.Params {
		labels[i] = fmt.Sprintf("param%d: %s", i+1, p.String())
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

	l := lines[line]
	if character > len(l) {
		character = len(l)
	}

	contentBefore := l[:character]

	openParenIdx := strings.LastIndex(contentBefore, "(")
	if openParenIdx == -1 {
		return "", 0
	}

	commas := strings.Count(contentBefore[openParenIdx:], ",")
	activeParam := commas

	nameEnd := openParenIdx
	for nameEnd > 0 && (l[nameEnd-1] == ' ' || l[nameEnd-1] == '\t') {
		nameEnd--
	}

	nameStart := nameEnd
	for nameStart > 0 && IsWordChar(rune(l[nameStart-1])) {
		nameStart--
	}

	if nameStart == nameEnd {
		return "", 0
	}

	return l[nameStart:nameEnd], activeParam

}
