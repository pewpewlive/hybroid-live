package lsp

import (
	"hybroid/alerts"
	"hybroid/lexer"
	"hybroid/parser"
	"hybroid/tokens"
	"hybroid/walker"
	"strings"
)

type AnalysisResult struct {
	Diagnostics []Diagnostic
	Walker      *walker.Walker
	Tokens      []tokens.Token
}

func Analyze(uri DocumentURI, text string) AnalysisResult {
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
	walk := walker.NewWalker("in-memory", "in-memory")
	walk.SetProgram(program)

	walk.PreWalk(nil)
	walk.Walk()
	walk.PostWalk()
	diagnostics = append(diagnostics, alertsToDiagnostics(uri, walk.GetAlerts())...)

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
			Message: alert.Message(),
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
