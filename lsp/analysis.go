package lsp

import (
	"fmt"
	"hybroid/alerts"
	"hybroid/lexer"
	"hybroid/parser"
	"hybroid/tokens"
	"hybroid/walker"
	"path/filepath"
	"strings"
)

type AnalysisResult struct {
	Diagnostics []Diagnostic
	Walker      *walker.Walker
	Tokens      []tokens.Token
}

func Analyze(uri DocumentURI, text string, walkerMap map[string]*walker.Walker, skipWalk bool) AnalysisResult {
	diagnostics := make([]Diagnostic, 0)

	if walkerMap == nil {
		walkerMap = make(map[string]*walker.Walker)
	}

	// Lexer
	lex := lexer.NewLexer(strings.NewReader(text))
	toks, err := lex.Tokenize()
	if err != nil {
		// Log lexer error if needed
	}
	diagnostics = append(diagnostics, alertsToDiagnostics(uri, lex.GetAlerts())...)

	// Parser
	parse := parser.NewParser(toks)
	program := parse.Parse()
	diagnostics = append(diagnostics, alertsToDiagnostics(uri, parse.GetAlerts())...)

	// Walker
	path, _ := fromURI(uri)
	absPath, err := filepath.Abs(path)
	if err == nil {
		path = absPath
	}

	// We need to determine the lua path relative to the project.
	// If we don't have a project root, preserve the directory structure to avoid collisions.
	luaPath := path
	if before, ok := strings.CutSuffix(luaPath, ".hyb"); ok {
		luaPath = before + ".lua"
	}

	walk := walker.NewWalker(path, luaPath)
	walk.SetProgram(program)

	// Ensure this walker is in the map by its path BEFORE PreWalk
	walkerMap[path] = walk

	walk.PreWalk(walkerMap)
	if !skipWalk {
		walk.Walk()
		walk.PostWalk()
		diagnostics = append(diagnostics, alertsToDiagnostics(uri, walk.GetAlerts())...)
	}

	// Ensure it's also indexed by its environment name if registered during PreWalk
	if walk.Env().Name != "" {
		walkerMap[walk.Env().Name] = walk
	}

	return AnalysisResult{
		Diagnostics: diagnostics,
		Walker:      walk,
		Tokens:      toks,
	}
}

func alertsToDiagnostics(uri DocumentURI, alertsList []alerts.Alert) []Diagnostic {
	diags := make([]Diagnostic, 0)
	for _, alert := range alertsList {
		snippet := alert.SnippetSpecifier()
		toks := snippet.GetTokens()

		var startTok, endTok tokens.Token
		if len(toks) > 0 {
			startTok = toks[0]
			endTok = toks[len(toks)-1]
		} else {
			// No tokens? Default to 0,0
			startTok = tokens.Token{Location: tokens.NewLocation(1, 1, 1)}
			endTok = startTok
		}

		// Clamp to non-negative LSP positions. Tokens are nominally
		// 1-based, so the -1 conversion should always yield >= 0, but
		// malformed tokens (e.g. from a hand-constructed test alert or
		// a future generator bug) can produce line=0/col=0. Editors
		// reject negative positions, so we floor them.
		startLine := max(startTok.Location.Line-1, 0)
		startCol := max(startTok.Location.Column.Start-1, 0)
		endLine := max(endTok.Location.Line-1, 0)
		endCol := max(endTok.Location.Column.End-1, 0)

		d := Diagnostic{
			Range: Range{
				Start: Position{Line: startLine, Character: startCol},
				End:   Position{Line: endLine, Character: endCol},
			},
			Message: fmt.Sprintf("[%s] %s", alert.ID(), alert.Message()),
			Severity: func() int {
				if alert.AlertType() == alerts.Error {
					return 1 // Error
				}
				return 2 // Warning
			}(),
			Source: func() *string { s := "hybroid"; return &s }(),
		}

		if note := alert.Note(); note != "" {
			d.RelatedInformation = []DiagnosticRelatedInformation{
				{
					Location: Location{
						URI:   uri,
						Range: d.Range,
					},
					Message: note,
				},
			}
		}

		diags = append(diags, d)
	}
	return diags
}
