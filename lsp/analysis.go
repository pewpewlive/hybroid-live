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
}

func Analyze(text string) AnalysisResult {
	diagnostics := make([]Diagnostic, 0)

	// Lexer
	lex := lexer.NewLexer(strings.NewReader(text))
	toks, err := lex.Tokenize()
	if err != nil {
		// If critical error, return early
		// But usually lexer accumulates alerts
	}
	diagnostics = append(diagnostics, alertsToDiagnostics(lex.GetAlerts())...)

	// Parser
	parse := parser.NewParser(toks)
	ast := parse.Parse()
	diagnostics = append(diagnostics, alertsToDiagnostics(parse.GetAlerts())...)

	// Walker
	// Assuming "Shared" env for generic LSP analysis for now, or infer from context?
	// For now, let's use a dummy path.
	walk := walker.NewWalker("in-memory", "in-memory")
	walk.SetProgram(ast)
	// We might need to mock walkers map for imports if we want to support multi-file analysis later.
	// For now, simple single-file.
	walk.PreWalk(nil)
	walk.Walk()
	walk.PostWalk()
	diagnostics = append(diagnostics, alertsToDiagnostics(walk.GetAlerts())...)

	return AnalysisResult{
		Diagnostics: diagnostics,
		Walker:      walk,
	}
}

func alertsToDiagnostics(alertsList []alerts.Alert) []Diagnostic {
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
		diags = append(diags, d)
	}
	return diags
}
