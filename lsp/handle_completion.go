package lsp

import (
	"context"
	"encoding/json"
	"hybroid/ast"
	"hybroid/evaluator"
	"hybroid/walker"
	"path/filepath"
	"strings"

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

	h.mu.Lock()
	eval := h.eval
	file, fileOk := h.files[params.TextDocument.URI]
	h.mu.Unlock()

	if eval == nil || !fileOk {
		return nil, nil
	}

	path, _ := fromURI(params.TextDocument.URI)
	relPath := getRelPath(h.rootPath, path)
	h.evalMu.Lock()
	w := eval.AnalyzeFile(relPath)
	result, err = h.completion(file, w, &params)
	h.evalMu.Unlock()

	return result, err
}

func (h *langHandler) completion(file *File, w *walker.Walker, params *CompletionParams) ([]CompletionItem, error) {
	items := make([]CompletionItem, 0)
	seen := make(map[string]bool)

	h.mu.Lock()
	eval := h.eval
	h.mu.Unlock()

	// 1. Get Context (Namespace and partial word)
	namespace, _, operator := h.getNamespaceContext(file.Text, params.Position)

	if namespace != "" {
		// Namespace-specific completion
		return h.namespaceCompletion(namespace, operator, w, eval, params.TextDocument.URI)
	}

	// 1.5. Check for 'env <name> as ' context -> suggest environment types
	if isEnvAsContext(file.Text, params.Position) {
		envTypes := []string{"Level", "Shared", "Mesh", "Sound"}
		for _, et := range envTypes {
			items = append(items, CompletionItem{
				Label:  et,
				Kind:   KeywordCompletion,
				Detail: "environment type",
			})
		}
		return items, nil
	}

	// 2. Local Scope and Current File Symbols
	if w != nil {
		line := params.Position.Line + 1
		col := params.Position.Character + 1
		scope := w.GetScopeAt(line, col)
		if scope != nil {
			current := scope
			for current != nil {
				for name, variable := range current.Variables {
					if !seen[name] {
						kind := VariableCompletion
						if _, ok := variable.Value.(*walker.FunctionVal); ok {
							kind = FunctionCompletion
						}
						items = append(items, CompletionItem{
							Label:  name,
							Kind:   kind,
							Detail: variable.Value.GetType().String(),
							Data:   params.TextDocument.URI,
						})
						seen[name] = true
					}
				}
				current = current.Parent
			}
		}

		// Types from current environment (Enums, Entities, Classes)
		env := w.Env()
		for name, ev := range env.Enums {
			if !seen[name] {
				items = append(items, CompletionItem{
					Label:  name,
					Kind:   EnumCompletion,
					Detail: "enum " + ev.Type.Name,
					Data:   params.TextDocument.URI,
				})
				seen[name] = true
			}
		}
		for name := range env.Entities {
			if !seen[name] {
				items = append(items, CompletionItem{
					Label:  name,
					Kind:   ClassCompletion,
					Detail: "entity " + name,
					Data:   params.TextDocument.URI,
				})
				seen[name] = true
			}
		}
		for name := range env.Classes {
			if !seen[name] {
				items = append(items, CompletionItem{
					Label:  name,
					Kind:   ClassCompletion,
					Detail: "class " + name,
					Data:   params.TextDocument.URI,
				})
				seen[name] = true
			}
		}

		// 3. Symbols from 'use' imported namespaces (WITHOUT prefix)
		for _, imp := range w.Env().Imports() {
			if imp.ThroughUse {
				env := imp.Env()
				// Add variables
				for name, variable := range env.Scope.Variables {
					if variable.IsPub && !seen[name] {
						kind := VariableCompletion
						if _, ok := variable.Value.(*walker.FunctionVal); ok {
							kind = FunctionCompletion
						}
						items = append(items, CompletionItem{
							Label:  name,
							Kind:   kind,
							Detail: variable.Value.GetType().String(),
							Data:   params.TextDocument.URI,
						})
						seen[name] = true
					}
				}
				// Add enums
				for name, ev := range env.Enums {
					if ev.IsPub && !seen[name] {
						items = append(items, CompletionItem{
							Label:  name,
							Kind:   EnumCompletion,
							Detail: "enum " + ev.Type.Name,
							Data:   params.TextDocument.URI,
						})
						seen[name] = true
					}
				}
			}
		}

		// 4. Builtin Libraries ONLY if explicitly imported via 'use'
		for _, lib := range w.Env().ImportedLibraries {
			var libEnv *walker.Environment
			switch lib {
			case ast.Pewpew:
				libEnv = walker.PewpewAPI
			case ast.Fmath:
				libEnv = walker.FmathAPI
			case ast.Math:
				libEnv = walker.MathAPI
			case ast.String:
				libEnv = walker.StringAPI
			case ast.Table:
				libEnv = walker.TableAPI
			}

			if libEnv != nil {
				for name, variable := range libEnv.Scope.Variables {
					if !seen[name] {
						kind := VariableCompletion
						if _, ok := variable.Value.(*walker.FunctionVal); ok {
							kind = FunctionCompletion
						}
						items = append(items, CompletionItem{
							Label:  name,
							Kind:   kind,
							Detail: variable.Value.GetType().String(),
							Data:   params.TextDocument.URI,
						})
						seen[name] = true
					}
				}
				for name, ev := range libEnv.Enums {
					if !seen[name] {
						items = append(items, CompletionItem{
							Label:  name,
							Kind:   EnumCompletion,
							Detail: "enum " + ev.Type.Name,
							Data:   params.TextDocument.URI,
						})
						seen[name] = true
					}
				}
				for name, cv := range libEnv.Classes {
					if !seen[name] {
						items = append(items, CompletionItem{
							Label:  name,
							Kind:   ClassCompletion,
							Detail: "class " + cv.Type.Name,
							Data:   params.TextDocument.URI,
						})
						seen[name] = true
					}
				}
				for name, εν := range libEnv.Entities {
					if !seen[name] {
						items = append(items, CompletionItem{
							Label:  name,
							Kind:   ClassCompletion,
							Detail: "entity " + εν.Type.Name,
							Data:   params.TextDocument.URI,
						})
						seen[name] = true
					}
				}
			}
		}
	}

	// 5. Standard Keywords
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

	// 6. Native Types
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

	// 7. Namespaces (Always available as prefixes)
	builtinNamespaces := []string{"Pewpew", "Fmath", "Math", "String", "Table"}
	for _, ns := range builtinNamespaces {
		if !seen[ns] {
			items = append(items, CompletionItem{
				Label: ns,
				Kind:  ModuleCompletion,
			})
			seen[ns] = true
		}
	}

	// 8. Custom Environments/Namespaces from eval
	if eval != nil {
		for name := range eval.Walkers() {
			// Skip absolute paths, only use environment names
			if filepath.IsAbs(name) || strings.ContainsAny(name, "/\\") {
				continue
			}

			// Add the namespace itself
			if !seen[name] {
				items = append(items, CompletionItem{
					Label:  name,
					Kind:   ModuleCompletion,
					Detail: "Environment",
				})
				seen[name] = true
			}
		}
	}

	// 9. Builtin Functions (Always available)
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

	return items, nil
}

func (h *langHandler) getNamespaceContext(text string, pos Position) (string, string, string) {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(text, "\n")
	if pos.Line < 0 || pos.Line >= len(lines) {
		return "", "", ""
	}
	line := lines[pos.Line]

	// If at the very start of a line, no namespace
	if pos.Character <= 0 {
		return "", "", ""
	}

	// Search backwards for : or . starting from the character BEFORE the cursor
	// because at Namespace:| the character at pos.Character might be space or newline
	curr := pos.Character
	if curr > len(line) {
		curr = len(line)
	}

	// We might be at Namespace:Part|
	// Scan back for the start of the current "word"
	wordStart := curr
	for wordStart > 0 && isWordChar(rune(line[wordStart-1])) {
		wordStart--
	}

	// Now check if the character before wordStart is ':' or '.'
	if wordStart > 0 && (line[wordStart-1] == ':' || line[wordStart-1] == '.') {
		operator := string(line[wordStart-1])
		nsEnd := wordStart - 1
		nsStart := nsEnd
		for nsStart > 0 && (isWordChar(rune(line[nsStart-1])) || line[nsStart-1] == ':') {
			nsStart--
		}
		if nsStart < nsEnd {
			namespace := line[nsStart:nsEnd]
			partial := line[wordStart:curr]
			return namespace, partial, operator
		}
	}

	return "", "", ""
}

func isWordChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
}

func (h *langHandler) namespaceCompletion(namespace string, operator string, w *walker.Walker, eval *evaluator.Evaluator, uri DocumentURI) ([]CompletionItem, error) {
	items := make([]CompletionItem, 0)

	var targetEnv *walker.Environment
	envName := namespace
	enumName := namespace

	if strings.Contains(namespace, ":") {
		parts := strings.Split(namespace, ":")
		envName = parts[0]
		enumName = parts[len(parts)-1]
	}

	if operator == ":" || strings.Contains(namespace, ":") {
		// Check builtins
		switch envName {
		case "Pewpew":
			targetEnv = walker.PewpewAPI
		case "Fmath":
			targetEnv = walker.FmathAPI
		case "Math":
			targetEnv = walker.MathAPI
		case "String":
			targetEnv = walker.StringAPI
		case "Table":
			targetEnv = walker.TableAPI
		default:
			// Check custom environments
			if eval != nil {
				if w2, ok := eval.Walkers()[envName]; ok {
					targetEnv = w2.Env()
				}
			}
		}
	}

	if targetEnv == nil && operator == "." {
		var foundEnum *walker.EnumVal
		// Search w.Env().Enums for it
		if w != nil {
			if ev, ok := w.Env().Enums[enumName]; ok {
				foundEnum = ev
			} else {
				// Check imports via 'use'
				for _, imp := range w.Env().Imports() {
					if imp.ThroughUse {
						if ev, ok := imp.Env().Enums[enumName]; ok && ev.IsPub {
							foundEnum = ev
							break
						}
					}
				}
				// Check explicitly imported builtin libraries if not found
				if foundEnum == nil {
					for _, lib := range w.Env().ImportedLibraries {
						var libEnv *walker.Environment
						switch lib {
						case ast.Pewpew:
							libEnv = walker.PewpewAPI
						case ast.Fmath:
							libEnv = walker.FmathAPI
						case ast.Math:
							libEnv = walker.MathAPI
						case ast.String:
							libEnv = walker.StringAPI
						case ast.Table:
							libEnv = walker.TableAPI
						}
						if libEnv != nil {
							if ev, ok := libEnv.Enums[enumName]; ok { // Builtin enums act as pub
								foundEnum = ev
								break
							}
						}
					}
				}
			}
		}

		if foundEnum != nil {
			for name := range foundEnum.Fields {
				items = append(items, CompletionItem{
					Label:  name,
					Kind:   FieldCompletion,
					Detail: "enum field",
					Data:   uri,
				})
			}
			return items, nil
		}
		return items, nil
	}

	if targetEnv != nil && operator == "." {
		if ev, ok := targetEnv.Enums[enumName]; ok {
			for name := range ev.Fields {
				items = append(items, CompletionItem{
					Label:  name,
					Kind:   FieldCompletion,
					Detail: "enum field",
					Data:   uri,
				})
			}
			return items, nil
		}
	}

	if targetEnv == nil {
		return items, nil
	}

	// Add symbols from target environment
	isBuiltin := targetEnv.Name == "Pewpew" || targetEnv.Name == "Fmath" || targetEnv.Name == "Math" || targetEnv.Name == "String" || targetEnv.Name == "Table"

	for name, variable := range targetEnv.Scope.Variables {
		if isBuiltin || variable.IsPub {
			kind := VariableCompletion
			if _, ok := variable.Value.(*walker.FunctionVal); ok {
				kind = FunctionCompletion
			}
			items = append(items, CompletionItem{
				Label:  name,
				Kind:   kind,
				Detail: variable.Value.GetType().String(),
				Data:   uri,
			})
		}
	}

	for name, ev := range targetEnv.Enums {
		if isBuiltin || ev.IsPub {
			items = append(items, CompletionItem{
				Label:  name,
				Kind:   EnumCompletion,
				Detail: "enum " + ev.Type.Name,
				Data:   uri,
			})
		}
	}

	for name, cv := range targetEnv.Classes {
		if isBuiltin || cv.IsPub {
			items = append(items, CompletionItem{
				Label:  name,
				Kind:   ClassCompletion,
				Detail: "class " + name,
				Data:   uri,
			})
		}
	}

	for name, ev := range targetEnv.Entities {
		if isBuiltin || ev.IsPub {
			items = append(items, CompletionItem{
				Label:  name,
				Kind:   ClassCompletion,
				Detail: "entity " + name,
				Data:   uri,
			})
		}
	}

	return items, nil
}

func getWordBefore(text string, line, character int) string {
	// This is now redundant or can be refactored, but I'll keep it if needed elsewhere
	// Actually, it was used in the previous version. I'll leave it for now.
	text = strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(text, "\n")
	if line < 0 || line >= len(lines) {
		return ""
	}
	l := lines[line]
	if character <= 0 || character > len(l) {
		return ""
	}

	end := character - 1
	for end > 0 && (l[end-1] == ' ' || l[end-1] == '\t') {
		end--
	}

	start := end
	for start > 0 && isWordChar(rune(l[start-1])) {
		start--
	}

	if start == end {
		return ""
	}

	return l[start:end]
}

// isEnvAsContext checks if the cursor is positioned after 'env <name> as '
// to provide environment type completions.
func isEnvAsContext(text string, pos Position) bool {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(text, "\n")
	if pos.Line < 0 || pos.Line >= len(lines) {
		return false
	}
	line := lines[pos.Line]
	before := strings.TrimSpace(line[:min(pos.Character, len(line))])
	// Match pattern: "env <identifier> as" (with optional partial word after)
	parts := strings.Fields(before)
	if len(parts) >= 3 && parts[0] == "env" && parts[2] == "as" {
		return true
	}
	return false
}
