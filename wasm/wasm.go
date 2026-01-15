//go:build js

package wasm

import (
	"bufio"
	"fmt"
	"hybroid/alerts"
	"hybroid/generator"
	"hybroid/lexer"
	"hybroid/parser"
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

func compile(code string) (string, error) {
	l := lexer.NewLexer(strings.NewReader(code))
	tokensList, err := l.Tokenize()
	if err != nil {
		return "", err
	}

	if len(l.GetAlerts()) > 0 {
		hasError := false
		for _, a := range l.GetAlerts() {
			if a.AlertType() == alerts.Error {
				hasError = true
				break
			}
		}

		msg := formatAlerts(l.GetAlerts(), code)
		if hasError {
			return "", fmt.Errorf("%s", msg)
		}
		fmt.Println(msg) // Log warnings
	}

	p := parser.NewParser(tokensList)
	program := p.Parse()

	if len(p.GetAlerts()) > 0 {
		hasError := false
		for _, a := range p.GetAlerts() {
			if a.AlertType() == alerts.Error {
				hasError = true
				break
			}
		}

		msg := formatAlerts(p.GetAlerts(), code)
		if hasError {
			return "", fmt.Errorf("%s", msg)
		}
		fmt.Println(msg)
	}

	gen := generator.NewGenerator()
	generator.ResetGlobalGeneratorValues()
	gen.Generate(program, nil)

	return gen.GetSrc(), nil
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
	fmt.Println("Hybroid Live for WebAssembly, v0.1.0")
	js.Global().Set("hybroidCompile", compileWrapper())
}
