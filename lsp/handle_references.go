package lsp

import (
	"context"
	"encoding/json"
	"hybroid/walker"
	"path/filepath"

	"github.com/sourcegraph/jsonrpc2"
)

// ReferenceParams extends TextDocumentPositionParams with reference context.
type ReferenceParams struct {
	TextDocumentPositionParams
	Context ReferenceContext `json:"context"`
}

// ReferenceContext controls whether the declaration itself should be included.
type ReferenceContext struct {
	IncludeDeclaration bool `json:"includeDeclaration"`
}

func (h *langHandler) handleTextDocumentReferences(_ context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params ReferenceParams
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

	word := getWordAt(file.Text, params.Position.Line, params.Position.Character)
	if word == "" {
		return nil, nil
	}

	rootDir := h.rootPath
	if rootDir == "" {
		rootDir = filepath.Dir(path)
	}

	locations := h.findReferences(w, eval.Walkers(), eval.WalkerList(), word, params.Position.Line+1, params.Position.Character+1, params.Context.IncludeDeclaration, rootDir)
	if len(locations) == 0 {
		return nil, nil
	}

	return locations, nil
}

func (h *langHandler) findReferences(w *walker.Walker, walkers map[string]*walker.Walker, walkerList []*walker.Walker, label string, line, col int, includeDecl bool, rootPath string) []Location {
	var locations []Location

	// absHybPath resolves a relative HybroidPath to an absolute path for URI generation
	absHybPath := func(hybPath string) string {
		if filepath.IsAbs(hybPath) {
			return hybPath
		}
		return filepath.Join(rootPath, hybPath)
	}

	// Check if the label is an environment name first
	if _, ok := walkers[label]; ok {
		key := walker.RefKey("env", label)
		for _, wk := range walkerList {
			refs, ok := wk.ReferenceMap[key]
			if !ok {
				continue
			}
			for _, ref := range refs {
				locations = append(locations, toLSPLocation(absHybPath(wk.Env().HybroidPath()), ref.Token))
			}
		}
		if includeDecl {
			declLoc := h.resolveDefinition(w, walkers, label, line, col, rootPath)
			if declLoc != (Location{}) {
				locations = append([]Location{declLoc}, locations...)
			}
		}
		return locations
	}

	// Determine the definition's environment name and variable name
	defEnvName := ""
	varName := label

	// Handle Namespace:Symbol
	if idx := findNsSeparator(label); idx >= 0 {
		ns := label[:idx]
		varName = label[idx+1:]

		// Determine the env name for the namespace
		switch ns {
		case "Pewpew", "Fmath", "Math", "String", "Table":
			defEnvName = ns
		default:
			if w2, ok := walkers[ns]; ok {
				defEnvName = w2.Env().Name
			}
		}
	} else {
		// Unqualified symbol — find where it's defined
		scope := w.GetScopeAt(line, col)
		if scope != nil {
			current := scope
			for current != nil {
				if _, ok := current.Variables[label]; ok {
					defEnvName = current.Environment.Name
					break
				}
				current = current.Parent
			}
		}

		// Check current walker's top-level enums, entities, classes
		if defEnvName == "" && w != nil {
			env := w.Env()
			if _, ok := env.Enums[label]; ok {
				defEnvName = env.Name
			} else if _, ok := env.Entities[label]; ok {
				defEnvName = env.Name
			} else if _, ok := env.Classes[label]; ok {
				defEnvName = env.Name
			}
		}

		// Check ThroughUse imports
		if defEnvName == "" && w != nil {
			for _, imp := range w.Env().Imports() {
				if imp.ThroughUse {
					if v, ok := imp.Env().Scope.Variables[label]; ok && v.IsPub {
						defEnvName = imp.Env().Name
						break
					} else if ev, ok := imp.Env().Enums[label]; ok && ev.IsPub {
						defEnvName = imp.Env().Name
						break
					} else if ev, ok := imp.Env().Entities[label]; ok && ev.IsPub {
						defEnvName = imp.Env().Name
						break
					} else if cv, ok := imp.Env().Classes[label]; ok && cv.IsPub {
						defEnvName = imp.Env().Name
						break
					}
				}
			}
		}

		// Check imported libraries
		if defEnvName == "" && w != nil {
			for _, lib := range w.Env().ImportedLibraries {
				libEnv := walker.BuiltinLibraries[lib]
				if libEnv != nil {
					if _, ok := libEnv.Scope.Variables[label]; ok {
						defEnvName = libEnv.Name
						break
					} else if _, ok := libEnv.Enums[label]; ok {
						defEnvName = libEnv.Name
						break
					} else if _, ok := libEnv.Entities[label]; ok {
						defEnvName = libEnv.Name
						break
					} else if _, ok := libEnv.Classes[label]; ok {
						defEnvName = libEnv.Name
						break
					}
				}
			}
		}
	}

	if defEnvName == "" {
		return nil
	}

	// Build the reference key
	key := walker.RefKey(defEnvName, varName)

	// Collect references from ALL walkers
	for _, wk := range walkerList {
		refs, ok := wk.ReferenceMap[key]
		if !ok {
			continue
		}
		for _, ref := range refs {
			locations = append(locations, toLSPLocation(absHybPath(wk.Env().HybroidPath()), ref.Token))
		}
	}

	// Optionally include the declaration itself
	if includeDecl {
		declLoc := h.resolveDefinition(w, walkers, label, line, col, rootPath)
		if declLoc != (Location{}) {
			locations = append([]Location{declLoc}, locations...)
		}
	}

	return locations
}

// findNsSeparator returns the index of ':' or '.' namespace separator, or -1.
func findNsSeparator(s string) int {
	for i, c := range s {
		if c == ':' || c == '.' {
			return i
		}
	}
	return -1
}
