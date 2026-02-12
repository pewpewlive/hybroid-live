package lsp

import (
	"context"
	"encoding/json"
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
	relPath, _ := filepath.Rel(h.rootPath, path)
	w := eval.AnalyzeFile(relPath)

	return h.completion(file, w, &params)
}

func (h *langHandler) completion(file *File, w *walker.Walker, params *CompletionParams) ([]CompletionItem, error) {
	items := make([]CompletionItem, 0)
	seen := make(map[string]bool)

	h.mu.Lock()
	eval := h.eval
	h.mu.Unlock()

	// 1. Cross-file symbols from other environments
	if eval != nil {
		for name, w2 := range eval.Walkers() {
			// Skip absolute paths, only use environment names
			if filepath.IsAbs(name) || strings.ContainsAny(name, "/\\") {
				continue
			}

			// Add the namespace itself
			if !seen[name] {
				items = append(items, CompletionItem{
					Label:         name,
					Kind:          ModuleCompletion,
					Detail:        "Environment",
					Documentation: "Custom environment: " + name,
				})
				seen[name] = true
			}

			// Add public symbols from this environment
			env := w2.Env()
			for varName, variable := range env.Scope.Variables {
				if variable.IsPub {
					fullLabel := name + ":" + varName
					if !seen[fullLabel] {
						items = append(items, CompletionItem{
							Label:  fullLabel,
							Kind:   FunctionCompletion, // Could be more specific
							Detail: name,
						})
						seen[fullLabel] = true
					}
				}
			}
			for enumName, ev := range env.Enums {
				if ev.IsPub {
					fullLabel := name + ":" + enumName
					if !seen[fullLabel] {
						items = append(items, CompletionItem{
							Label:  fullLabel,
							Kind:   EnumCompletion,
							Detail: name,
						})
						seen[fullLabel] = true
					}
				}
			}
		}
	}

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

	// 7. Enums from Pewpew and Fmath
	for name, ev := range walker.PewpewAPI.Enums {
		if !seen[name] {
			items = append(items, CompletionItem{
				Label:  name,
				Kind:   EnumCompletion,
				Detail: "Pewpew:enum " + ev.Type.Name,
			})
			seen[name] = true
		}
	}
	for name, ev := range walker.FmathAPI.Enums {
		if !seen[name] {
			items = append(items, CompletionItem{
				Label:  name,
				Kind:   EnumCompletion,
				Detail: "Fmath:enum " + ev.Type.Name,
			})
			seen[name] = true
		}
	}

	if w == nil {
		return items, nil
	}

	// 8. Enums from current walker
	for name, ev := range w.Env().Enums {
		if !seen[name] {
			items = append(items, CompletionItem{
				Label:  name,
				Kind:   EnumCompletion,
				Detail: "enum " + ev.Type.Name,
			})
			seen[name] = true
		}
	}

	// 9. Entities from current walker
	for name := range w.Env().Entities {
		if !seen[name] {
			items = append(items, CompletionItem{
				Label:  name,
				Kind:   ClassCompletion,
				Detail: "entity " + name,
			})
			seen[name] = true
		}
	}

	// 10. Classes from current walker
	for name := range w.Env().Classes {
		if !seen[name] {
			items = append(items, CompletionItem{
				Label:  name,
				Kind:   ClassCompletion,
				Detail: "class " + name,
			})
			seen[name] = true
		}
	}

	// Check if we are completing after a : or .
	if params.Context.TriggerKind == 2 && params.Context.TriggerCharacter != nil {
		trigger := *params.Context.TriggerCharacter
		if trigger == ":" || trigger == "." {
			// Get the word before trigger
			wordBefore := getWordBefore(file.Text, params.Position.Line, params.Position.Character)
			if wordBefore != "" {
				line := params.Position.Line + 1
				col := params.Position.Character + 1
				scope := w.GetScopeAt(line, col)

				// Check if wordBefore is an enum
				var enumVal *walker.EnumVal
				if ev, ok := w.Env().Enums[wordBefore]; ok {
					enumVal = ev
				} else if ev, ok := walker.PewpewAPI.Enums[wordBefore]; ok {
					enumVal = ev
				} else if ev, ok := walker.FmathAPI.Enums[wordBefore]; ok {
					enumVal = ev
				}

				if enumVal != nil {
					enumItems := make([]CompletionItem, 0)
					for fieldName, fieldVar := range enumVal.Fields {
						enumItems = append(enumItems, CompletionItem{
							Label:  fieldName,
							Kind:   EnumMemberCompletion,
							Detail: fieldVar.Value.GetType().String(),
						})
					}
					return enumItems, nil
				}

				// Check if wordBefore is a variable (instance) or an entity/class type
				var container walker.FieldContainer
				var methodContainer walker.MethodContainer

				if scope != nil {
					if v, ok := scope.GetVariable(wordBefore); ok {
						if fc, ok := v.Value.(walker.FieldContainer); ok {
							container = fc
						}
						if mc, ok := v.Value.(walker.MethodContainer); ok {
							methodContainer = mc
						}
					}
				}

				// If not found as instance, check as type (for static-like access like spawn/new)
				if container == nil && methodContainer == nil {
					if ev, ok := w.Env().Entities[wordBefore]; ok {
						container = ev
						methodContainer = ev
					} else if cv, ok := w.Env().Classes[wordBefore]; ok {
						container = cv
						methodContainer = cv
					}
				}

				if container != nil || methodContainer != nil {
					memberItems := make([]CompletionItem, 0)
					seenMembers := make(map[string]bool)
					if container != nil {
						// For EntityVal/ClassVal, Fields is map[string]Field
						// We need to access the actual map because FieldContainer interface doesn't expose it.
						// But wait, let's use type switch to be safe.
						switch c := container.(type) {
						case *walker.EntityVal:
							for name, f := range c.Fields {
								memberItems = append(memberItems, CompletionItem{
									Label:  name,
									Kind:   FieldCompletion,
									Detail: f.Var.Value.GetType().String(),
								})
								seenMembers[name] = true
							}
						case *walker.ClassVal:
							for name, f := range c.Fields {
								memberItems = append(memberItems, CompletionItem{
									Label:  name,
									Kind:   FieldCompletion,
									Detail: f.Var.Value.GetType().String(),
								})
								seenMembers[name] = true
							}
						case *walker.EnumVal:
							// Already handled above, but for completeness
							for name, f := range c.Fields {
								memberItems = append(memberItems, CompletionItem{
									Label:  name,
									Kind:   EnumMemberCompletion,
									Detail: f.Value.GetType().String(),
								})
								seenMembers[name] = true
							}
						}
					}
					if methodContainer != nil {
						switch c := methodContainer.(type) {
						case *walker.EntityVal:
							for name, m := range c.Methods {
								if !seenMembers[name] {
									memberItems = append(memberItems, CompletionItem{
										Label:  name,
										Kind:   MethodCompletion,
										Detail: m.Value.GetType().String(),
									})
								}
							}
						case *walker.ClassVal:
							for name, m := range c.Methods {
								if !seenMembers[name] {
									memberItems = append(memberItems, CompletionItem{
										Label:  name,
										Kind:   MethodCompletion,
										Detail: m.Value.GetType().String(),
									})
								}
							}
						}
					}
					return memberItems, nil
				}
			}
		}
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

func getWordBefore(text string, line, character int) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(text, "\n")
	if line < 0 || line >= len(lines) {
		return ""
	}
	l := lines[line]
	if character <= 0 || character > len(l) {
		return ""
	}

	// We are at the position of the trigger character (e.g. :).
	// We want the word BEFORE it.
	end := character - 1
	for end > 0 && (l[end-1] == ' ' || l[end-1] == '\t') {
		end--
	}

	start := end
	// Use a local isWordChar that excludes : and .
	isWordCharLocal := func(r rune) bool {
		return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
	}

	for start > 0 && isWordCharLocal(rune(l[start-1])) {
		start--
	}

	if start == end {
		return ""
	}

	return l[start:end]
}
