//go:build js && wasm

package wasm

import (
	"bufio"
	"errors"
	"fmt"
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/generator"
	"hybroid/lexer"
	"hybroid/parser"
	"hybroid/walker"
	"strings"
	"syscall/js"
)

// Helper to reconstruct lines from source code
func sourceToLines(source string) map[int][]byte {
	lines := make(map[int][]byte)
	scanner := bufio.NewScanner(strings.NewReader(source))
	lineNum := 1
	for scanner.Scan() {
		// Copy the bytes because scanner reuses the buffer
		txt := scanner.Text()
		lines[lineNum] = []byte(txt)
		lineNum++
	}
	return lines
}

// formatAlerts converts a list of alerts into a formatted string with ANSI color codes, including error locations and code snippets.
func formatAlerts(alertsList []alerts.Alert, source string) string {
	lines := sourceToLines(source)
	var sb strings.Builder

	for _, alert := range alertsList {
		msg := ""
		switch alert.AlertType() {
		case alerts.Error:
			msg = fmt.Sprintf("[light_red][bold]error[%s]: [reset]", alert.ID())
		case alerts.Warning:
			msg = fmt.Sprintf("[light_yellow][bold]warning[%s]: [default]", alert.ID())
		}
		sb.WriteString(msg)
		sb.WriteString(fmt.Sprintf("[bold]%s[reset]\n", alert.Message()))

		// Location
		tokensList := alert.SnippetSpecifier().GetTokens()
		if len(tokensList) > 0 {
			sb.WriteString(fmt.Sprintf("  at line %d:%d\n", tokensList[0].Line, tokensList[0].Column.Start))
		}

		// Snippet
		snippet := alert.SnippetSpecifier().GetSnippet(lines, alert)
		sb.WriteString(snippet)

		// Note
		if alert.Note() != "" {
			sb.WriteString(fmt.Sprintf("note: %s\n", alert.Note()))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// processAlerts processes a list of alerts, returning an error if any are errors,
// or accumulating warnings in the provided strings.Builder.
func processAlerts(alertsList []alerts.Alert, source string, warnings *strings.Builder) error {
	if len(alertsList) == 0 {
		return nil
	}

	hasError := false
	for _, a := range alertsList {
		if a.AlertType() == alerts.Error {
			hasError = true
			break
		}
	}

	msg := formatAlerts(alertsList, source)
	if hasError {
		return errors.New(msg)
	}
	warnings.WriteString(msg)
	return nil
}

func compile(code string) (string, error) {
	var warnings strings.Builder

	l := lexer.NewLexer(strings.NewReader(code))
	tokensList, err := l.Tokenize()
	if err != nil {
		return "", err
	}

	if err := processAlerts(l.GetAlerts(), code, &warnings); err != nil {
		return "", err
	}

	p := parser.NewParser(tokensList)
	program := p.Parse()

	if err := processAlerts(p.GetAlerts(), code, &warnings); err != nil {
		return "", err
	}

	walker.SetupLibraryEnvironments()
	w := walker.NewWalker("main.hyb", "main.lua")
	w.SetProgram(program)

	// Single file compilation, so no other walkers to share context with
	walkers := make(map[string]*walker.Walker)
	w.PreWalk(walkers)
	w.Walk()
	w.PostWalk()

	if err := processAlerts(w.GetAlerts(), code, &warnings); err != nil {
		return "", err
	}

	gen := generator.NewGenerator()
	generator.ResetGlobalGeneratorValues()

	gen.SetEnv(w.Env().Name, w.Env().Type)
	gen.GenerateUsedLibraries(w.Env().UsedLibraries)

	if w.Env().Type != ast.LevelEnv {
		gen.Generate(w.Program(), w.Env().UsedBuiltinVars)
	} else {
		gen.Generate(w.Program(), []string{})
	}

	res := gen.GetSrc()
	if warnings.Len() > 0 {
		res = warnings.String() + "[default]============\n\n" + res
	}

	return res, nil
}

func compileWrapper() js.Func {
	compileFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return "expected 1 argument"
		}
		if args[0].Type() != js.TypeString {
			return "expected string"
		}
		code := args[0].String()
		output, err := compile(code)
		if err != nil {
			// Errors are returned instead of printed
			// fmt.Printf("unable to compile code: %s\n", err)
			return err.Error()
		}
		return output
	})
	return compileFunc
}

func init() {
	fmt.Println("Hybroid Live for WebAssembly v0.1.0 has been initialized.")
	js.Global().Set("hybroidCompile", compileWrapper())
}
