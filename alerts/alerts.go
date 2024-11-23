package alerts

import (
	"fmt"
	"hybroid/tokens"
	"reflect"
	"strconv"
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

/*
TODO: add fix snippet
"fix": {
      "insert": "number",
      "where": "before"
    }
*/

type Alert interface {
	GetMessage() string
	GetTokens() []tokens.Token
	GetLocations() []tokens.TokenLocation

	// Empty string means no note will be printed
	GetNote() string

	GetName() string
	GetAlertType() AlertType

	// AddTokens(...tokens.Token)
}

type AlertHandler struct {
	Source []byte

	Alerts    []Alert
	HasAlerts bool

	currentLine int
}

func (ah *AlertHandler) NewAlert(alertType Alert, args ...any) Alert {
	alert := reflect.ValueOf(alertType).Elem()

	for i, arg := range args {
		if reflect.TypeOf(arg) != alert.Field(i).Type() {
			panic(fmt.Sprintf("Attempt to construct %s{} field `%s` of type `%s`, with `%s` at %d\n", alert.Type().Name(), reflect.TypeOf(alertType).Elem().Field(i).Name, alert.Field(i).Type(), reflect.TypeOf(arg), i+1))
		}
		alert.Field(i).Set(reflect.ValueOf(arg))
	}

	return alert.Addr().Interface().(Alert)
}

func (ah *AlertHandler) Alert_(alertType Alert, args ...any) {
	ah.HasAlerts = true
	ah.Alerts = append(ah.Alerts, ah.NewAlert(alertType, args...))
}

func (ah *AlertHandler) AlertI_(alertType Alert) {
	ah.HasAlerts = true
	ah.Alerts = append(ah.Alerts, alertType)
}

func (ah *AlertHandler) PrintAlerts(alertStage AlertStage, source []byte, sourcePath string) {
	ah.Source = source

	var errMsg string
	switch alertStage {
	case Lexer, Parser:
		errMsg = "[red]Syntax error:"
	case Walker, Eval:
		errMsg = "[red]Compilation error:"
	}

	for _, alert := range ah.Alerts {
		colorstring.Print(errMsg)
		colorstring.Printf(" [light_gray]%s[default] in %s", alert.GetName(), sourcePath)
		ah.PrintLocation(alert)
		ah.PrintMessage(alert)
		ah.PrintCodeSnippet(alert)
		ah.PrintNote(alert)
		fmt.Println()
	}

	ah.Finish()
}

func (ah *AlertHandler) SortTypes() {
	//TODO: error and warning sorting
}

func (ah *AlertHandler) PrintMessage(alert Alert) {
	colorstring.Printf("[red]message:[default] %s.\n", alert.GetMessage())
}

func (ah *AlertHandler) PrintLocation(alert Alert) {
	location := CombineLocations(alert.GetLocations())
	fmt.Printf(": %d\n", location.LineStart)
}

func (ah *AlertHandler) PrintCodeSnippet(alert Alert) {
	lineCount := 1
	columnCount := 0
	location := CombineLocations(alert.GetLocations())

	for i := 0; i < len(ah.Source); i++ {
		columnCount += 1

		if lineCount == location.LineStart && lineCount == location.LineEnd && ah.Source[i] == '\n' {
			ah.snippetPrintSingleLine(i, columnCount, lineCount, location)
		}

		// handles multiple line errors
		if lineCount == location.LineStart && lineCount != location.LineEnd && ah.Source[i] == '\n' {
			snippetStart := ah.Source[i-columnCount+1 : i-1]

			colorstring.Printf("[cyan]%d%s |[default]   %s\n", lineCount, strings.Repeat(" ", len(strconv.Itoa(location.LineEnd))-1), string(snippetStart))
			colorstring.Printf("[cyan]%s |[light_red]  %s^\n", strings.Repeat(" ", len(strconv.Itoa(location.LineEnd))), strings.Repeat("_", location.ColStart))
		}

		if lineCount != location.LineStart && lineCount == location.LineEnd && ah.Source[i] == '\n' {
			snippetEnd := ah.Source[i-columnCount+1 : i-1]

			colorstring.Printf("[dark_gray]...[default]%s[light_red]|\n", strings.Repeat(" ", len(strconv.Itoa(lineCount))))
			colorstring.Printf("[cyan]%d |[light_red] | [default]%s\n", lineCount, string(snippetEnd))
			colorstring.Printf("[cyan]%s |[light_red] |%s^\n", strings.Repeat(" ", len(strconv.Itoa(lineCount))), strings.Repeat("_", location.ColEnd))
		}

		if ah.Source[i] == '\n' {
			lineCount += 1
			columnCount = 0
		}
	}

	ah.currentLine = lineCount
}

func (ah *AlertHandler) PrintNote(alert Alert) {
	if alert.GetNote() != "" {
		colorstring.Printf("[cyan]%s= note:[default] %s\n", strings.Repeat(" ", len(strconv.Itoa(ah.currentLine))), alert.GetNote())
		return
	}

	fmt.Print("\n")
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

func (ah *AlertHandler) snippetPrintSingleLine(index, columnCount, lineCount int, location tokens.TokenLocation) {
	line := ah.Source[index-columnCount+1 : index-1]
	if location.ColStart > 80 {
		start := string(line[0:30])
		var content string
		short := false
		if location.ColEnd+20 > columnCount && columnCount < 120 {
			content = string(line[location.ColStart-20 : columnCount-1])
		} else {
			content = string(line[location.ColStart-20 : location.ColEnd+20])
			short = true
		}
		var shortt string
		if short {
			shortt = "..."
		}
		// TODO: move prints into separate functions
		colorstring.Printf("[cyan]%d |[default]   %s[dark_gray]...[default]%s[dark_gray]%s[reset]\n", lineCount, start, content, shortt)
		colorstring.Printf("[cyan]%s |[light_red]   %s%s\n", strings.Repeat(" ", len(strconv.Itoa(lineCount))), strings.Repeat(" ", len(start)+22), strings.Repeat("^", location.ColEnd-location.ColStart+1))
	} else {
		colorstring.Printf("[cyan]%d |[default]   %s\n", lineCount, line)
		colorstring.Printf("[cyan]%s |[light_red]   %s%s\n", strings.Repeat(" ", len(strconv.Itoa(lineCount))), strings.Repeat(" ", location.ColStart-1), strings.Repeat("^", location.ColEnd-location.ColStart+1))
	}
}

func (ah *AlertHandler) snippetPrintMultiLine() {

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
