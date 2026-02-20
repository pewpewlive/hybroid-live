package lsp

import (
	"context"
	"encoding/json"
	"hybroid/walker"
	"path/filepath"
	"strings"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *langHandler) handleTextDocumentDefinition(_ context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params DocumentDefinitionParams
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
	if word == "" {
		return nil, nil
	}

	// 2. Resolve the definition
	rootDir := h.rootPath
	if rootDir == "" {
		rootDir = filepath.Dir(path)
	}
	loc := h.resolveDefinition(w, eval.Walkers(), word, params.Position.Line+1, params.Position.Character+1, rootDir)
	if loc != (Location{}) {
		return loc, nil
	}

	return nil, nil
}

func (h *langHandler) resolveDefinition(w *walker.Walker, walkers map[string]*walker.Walker, label string, line, col int, rootPath string) Location {
	// absHybPath resolves a relative HybroidPath to an absolute path for URI generation
	absHybPath := func(hybPath string) string {
		if filepath.IsAbs(hybPath) {
			return hybPath
		}
		return filepath.Join(rootPath, hybPath)
	}

	// Check if the label is an environment name (for `use MyHelper`, or namespace prefix in `Pewpew:X`)
	if walkers != nil {
		if w2, ok := walkers[label]; ok {
			// Navigate to the env declaration (first token of file, line 1)
			envToken := w2.Env().GetEnvToken()
			if envToken.Lexeme != "" {
				return toLSPLocation(absHybPath(w2.Env().HybroidPath()), envToken)
			}
		}
	}

	// Handle Namespace:Symbol or Namespace.Symbol
	if strings.Contains(label, ":") || strings.Contains(label, ".") {
		parts := strings.FieldsFunc(label, func(r rune) bool { return r == ':' || r == '.' })
		if len(parts) == 2 {
			ns := parts[0]
			sym := parts[1]

			var env *walker.Environment
			switch ns {
			case "Pewpew":
				env = walker.PewpewAPI
			case "Fmath":
				env = walker.FmathAPI
			case "Math":
				env = walker.MathAPI
			case "String":
				env = walker.StringAPI
			case "Table":
				env = walker.TableAPI
			}

			if env == nil && walkers != nil {
				if w2, ok := walkers[ns]; ok {
					env = w2.Env()
				}
			}

			// If not a namespace, check if it's an entity/enum/class in the current walker
			if env == nil && w != nil {
				if ev, ok := w.Env().Enums[ns]; ok {
					if field, _, found := ev.ContainsField(sym); found {
						return toLSPLocation(absHybPath(w.Env().HybroidPath()), field.Token)
					}
				}
				if ev, ok := w.Env().Entities[ns]; ok {
					if v, _, found := ev.ContainsField(sym); found {
						return toLSPLocation(absHybPath(w.Env().HybroidPath()), v.Token)
					}
					if v, found := ev.ContainsMethod(sym); found {
						return toLSPLocation(absHybPath(w.Env().HybroidPath()), v.Token)
					}
				}
				if cv, ok := w.Env().Classes[ns]; ok {
					if v, _, found := cv.ContainsField(sym); found {
						return toLSPLocation(absHybPath(w.Env().HybroidPath()), v.Token)
					}
					if v, found := cv.ContainsMethod(sym); found {
						return toLSPLocation(absHybPath(w.Env().HybroidPath()), v.Token)
					}
				}
			}

			if env != nil {
				if v, ok := env.Scope.Variables[sym]; ok && (env.Name == "Pewpew" || v.IsPub) {
					return toLSPLocation(absHybPath(env.HybroidPath()), v.Token)
				}
				if ev, ok := env.Enums[sym]; ok && (env.Name == "Pewpew" || ev.IsPub) {
					return toLSPLocation(absHybPath(env.HybroidPath()), ev.Token)
				}
			}
		}
	}

	// Check current scope variables
	scope := w.GetScopeAt(line, col)
	if scope != nil {
		current := scope
		for current != nil {
			if v, ok := current.Variables[label]; ok {
				// If it's a builtin, we might not have a meaningful hybroidPath
				if current.Environment.Name == "Builtin" || current.Environment.Name == "Pewpew" {
					return Location{}
				}
				return toLSPLocation(absHybPath(current.Environment.HybroidPath()), v.Token)
			}
			current = current.Parent
		}
	}

	// Check current walker's top-level enums, entities, classes
	if w != nil {
		env := w.Env()
		if ev, ok := env.Enums[label]; ok {
			return toLSPLocation(absHybPath(env.HybroidPath()), ev.Token)
		}
		if ev, ok := env.Entities[label]; ok {
			return toLSPLocation(absHybPath(env.HybroidPath()), ev.Token)
		}
		if cv, ok := env.Classes[label]; ok {
			return toLSPLocation(absHybPath(env.HybroidPath()), cv.Token)
		}

		// Check imported namespaces via 'use'
		for _, imp := range env.Imports() {
			if imp.ThroughUse {
				impEnv := imp.Env()
				if v, ok := impEnv.Scope.Variables[label]; ok && v.IsPub {
					return toLSPLocation(absHybPath(impEnv.HybroidPath()), v.Token)
				}
				if ev, ok := impEnv.Enums[label]; ok && ev.IsPub {
					return toLSPLocation(absHybPath(impEnv.HybroidPath()), ev.Token)
				}
				if cv, ok := impEnv.Classes[label]; ok && cv.IsPub {
					return toLSPLocation(absHybPath(impEnv.HybroidPath()), cv.Token)
				}
				if ev, ok := impEnv.Entities[label]; ok && ev.IsPub {
					return toLSPLocation(absHybPath(impEnv.HybroidPath()), ev.Token)
				}
			}
		}

		// Check used libraries
		for _, lib := range env.ImportedLibraries {
			libEnv := walker.BuiltinLibraries[lib]
			if libEnv != nil {
				if v, ok := libEnv.Scope.Variables[label]; ok {
					return toLSPLocation(absHybPath(libEnv.HybroidPath()), v.Token)
				}
				if ev, ok := libEnv.Enums[label]; ok {
					return toLSPLocation(absHybPath(libEnv.HybroidPath()), ev.Token)
				}
			}
		}
	}

	return Location{}
}
