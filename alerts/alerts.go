package alerts

import (
	"fmt"
	"hybroid/tokens"
	"reflect"
	"strings"

	"github.com/mitchellh/colorstring"
)

type AlertStage int

const (
	// Syntax error
	Lexer AlertStage = iota
	Parser

	// Compile error
	Walker
	Eval
)

type AlertType int

const (
	Error AlertType = iota
	Warning
)

type Alert interface {
	GetMessage() string
	GetTokens() []tokens.Token
	GetLocations() []tokens.TokenLocation

	// Empty string means no note will be printed
	GetNote() string

	GetAlertType() AlertType

	// AddTokens(...tokens.Token)
}

type AlertHandler struct {
	Source     []byte
	SourcePath string

	Alerts    []Alert
	HasAlerts bool
}

func (ah *AlertHandler) Alert(alertType Alert, args ...any) {
	ah.HasAlerts = true

	alert := reflect.ValueOf(alertType).Elem()

	for i, arg := range args {
		alert.Field(i).Set(reflect.ValueOf(arg))
	}

	ah.Alerts = append(ah.Alerts, alert.Addr().Interface().(Alert))
}

func (ah *AlertHandler) PrintAlerts(alertStage AlertStage, source []byte, sourcePath string) {
	ah.Source = source
	ah.SourcePath = sourcePath

	switch alertStage {
	case Lexer, Parser:
		colorstring.Print("[red]Syntax error")
	case Walker, Eval:
		colorstring.Print("[red]Compilation error")
	}

	fmt.Printf(" in file: %s", sourcePath)

	for _, alert := range ah.Alerts {
		ah.PrintLocation(alert)
		ah.PrintMessage(alert)
		ah.PrintCodeSnippet(alert)
		ah.PrintNote(alert)

		ah.Finish()
	}
}

// func (ah *AlertHandler) PrintAlert(alert Alert) {
// 	ah.PrintLocation(alert)
// 	ah.PrintMessage(alert)
// 	ah.PrintCodeSnippet(alert)
// 	ah.PrintNote(alert)

// 	ah.Finish()
// }

func (ah *AlertHandler) SortTypes() {
	//TODO: error and warning sorting
}

func (ah *AlertHandler) PrintMessage(alert Alert) {
	fmt.Printf("%s.\n", alert.GetMessage())
}

func (ah *AlertHandler) PrintLocation(alert Alert) {
	location := CombineLocations(alert.GetLocations())
	fmt.Printf(" at line: %d\n", location.LineStart)
}

func (ah *AlertHandler) PrintCodeSnippet(alert Alert) {
	lineCount := 1
	columnCount := 0
	location := CombineLocations(alert.GetLocations())

	for i := 0; i < len(ah.Source); i++ {
		columnCount += 1

		if lineCount == location.LineStart && lineCount == location.LineEnd && ah.Source[i] == '\n' { // handles single line errors
			snippet := ah.Source[i-columnCount+1 : i-1]

			fmt.Printf("%d |   %s\n", lineCount, string(snippet))
			//fmt.Printf("%d", location.ColStart)
			fmt.Printf("  |   %s%s", strings.Repeat(" ", location.ColStart), strings.Repeat("^", location.ColEnd-location.ColStart))
		}

		// if lineCount == location.LineStart && lineCount != location.LineEnd && ah.Source[i] == '\n' { // handles multiple line errors
		// 	snippet := ah.Source[i-columnCount+1 : i-1]

		// 	fmt.Printf("%d |   %s\n", lineCount, string(snippet))
		// 	//fmt.Printf("%d", location.ColStart)
		// 	fmt.Printf("  |   %s%s", strings.Repeat(" ", location.ColStart), strings.Repeat("^", location.ColEnd-location.ColStart))
		// }

		if ah.Source[i] == '\n' {
			lineCount += 1
			columnCount = 0
		}
	}

	fmt.Print("\n\n")
}

func (ah *AlertHandler) PrintNote(alert Alert) {
	if alert.GetNote() != "" {
		fmt.Printf("%s\n", alert.GetNote())
	}
}

func CombineLocations(locations []tokens.TokenLocation) tokens.TokenLocation {
	if len(locations) == 0 {
		return tokens.TokenLocation{}
	}
	location := locations[0]

	for i, v := range locations {
		if i == 0 {
			continue
		}

		if v.ColStart < location.ColStart {
			location.ColStart = v.ColStart
		}
		if v.ColEnd > location.ColEnd {
			location.ColEnd = v.ColEnd
		}
		if v.LineStart < location.LineStart {
			location.LineStart = v.LineStart
		}
		if v.LineEnd > location.LineEnd {
			location.LineEnd = v.LineEnd
		}
	}

	return location
}

func (ah *AlertHandler) GetCodeSnippet(location tokens.TokenLocation) string {
	return ""
}

func (ah *AlertHandler) Finish() {
	ah.Source = []byte{}
}

// type Alert interface {
// 	GetToken() lexer.Token
// 	GetHeader() string
// 	GetMessage() string
// }

// type Error struct {
// 	Token   lexer.Token
// 	Message string
// }

// func (e Error) GetToken() lexer.Token {
// 	return e.Token
// }

// func (e Error) GetMessage() string {
// 	return e.Message
// }

// func (e Error) GetHeader() string {
// 	return "[red]Error"
// }

// type Warning struct {
// 	Token   lexer.Token
// 	Message string
// }

// func (w Warning) GetToken() lexer.Token {
// 	return w.Token
// }

// func (w Warning) GetMessage() string {
// 	return w.Message
// }

// func (e Warning) GetHeader() string {
// 	return "[yellow]Warning"
// }

/*

Original func:
fn function() {

	return


	...
	...
}

Error code snippet:
3 	return
...
6 	...			<
7 	...			<
...
8 }


walkBody() {

	for i := range node.Body {
		if unreachable_code {
			unreachable code  d
		}
	}
}
UnreachableCode{Node[]}
*/

// DescriptiveErrors = false
// error in file: etc. etc.

// GetLocation() -> TokenLocation
// thing.fn(fixed, number, fixed)
// Incorrect function signature
// *code snippet*
// Description of what happened and how to fix it
// Ex: function params are: number, number

// new IncorrectFuncSignature{Token[], Loc, FunctioCallNode}

// new MalformedValError{string, TokenLocation}
// GetErrorType() -> "Malformed number"
