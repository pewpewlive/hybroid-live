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
	luaPath := filepath.Base(path)
	if strings.HasSuffix(luaPath, ".hyb") {
		luaPath = strings.TrimSuffix(luaPath, ".lua") + ".lua" // wait, suffix is .hyb
		luaPath = strings.TrimSuffix(filepath.Base(path), ".hyb") + ".lua"
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

		d := Diagnostic{
			Range: Range{
				Start: Position{
					Line:      startTok.Location.Line - 1,
					Character: startTok.Location.Column.Start - 1,
				},
				End: Position{
					Line:      endTok.Location.Line - 1,
					Character: endTok.Location.Column.End - 1,
				},
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
