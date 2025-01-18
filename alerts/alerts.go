package alerts

import (
	"fmt"
	"hybroid/tokens"
	"reflect"
	"strconv"
	"strings"

	color "github.com/mitchellh/colorstring"
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
	GetSpecifier() SnippetSpecifier

	// Empty string means no note will be printed
	GetNote() string

	GetID() string
	GetAlertType() AlertType

	// AddTokens(...tokens.Token)
}

type AlertHandler struct {
	Source []byte

	Alerts    []Alert
	HasAlerts bool

	currentLine int
}

func (ah *AlertHandler) NewAlert(alert Alert, args ...any) Alert {
	alertValue := reflect.ValueOf(alert).Elem()
	alertType := reflect.TypeOf(alert).Elem()

	fieldsSet := 0
	panicMessage := "Attempt to construct %s{} field `%s` of type `%s`%s, with `%s` at %d"

	for i, arg := range args {
		field := alertValue.Field(i)
		fieldType := field.Type()
		argValue := reflect.ValueOf(arg)
		argType := argValue.Type()

		if field.Kind() == reflect.Interface {
			if !argType.Implements(fieldType) {
				panic(fmt.Sprintf(panicMessage, alertValue.Type().Name(), alertType.Field(i).Name, fieldType, " (interface)", argType, i+1))
			}
			field.Set(argValue)
		} else {
			if argType == reflect.TypeFor[tokens.TokenType]() {
				argType = reflect.TypeFor[string]()
				argValue = argValue.Convert(argType)
			}

			if argType != fieldType {
				panic(fmt.Sprintf(panicMessage, alertValue.Type().Name(), alertType.Field(i).Name, fieldType, "", argType, i+1))
			}

			field.Set(argValue)
		}

		fieldsSet++
	}

	for i := 0; i < alertType.NumField(); i++ {
		field := alertValue.Field(i)
		if !field.IsZero() {
			continue
		}

		if defaultValue, ok := alertType.Field(i).Tag.Lookup("default"); ok {
			argValue := reflect.ValueOf(defaultValue)
			field.Set(argValue)
			fieldsSet++
		}
	}

	if alertValue.NumField() != fieldsSet {
		panicMessage := "Attempt to construct %s{} with invalid amount of arguments: expected %d, but got: %d"
		panic(fmt.Sprintf(panicMessage, alertValue.Type().Name(), alertValue.NumField(), fieldsSet))
	}

	return alertValue.Addr().Interface().(Alert)
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
		errMsg = "[light_red]Syntax error[%s]: "
	case Walker, Eval:
		errMsg = "[light_red]Compilation error[%s]: "
	}

	for _, alert := range ah.Alerts {
		color.Printf(errMsg, alert.GetID())
		ah.PrintLocation(alert, sourcePath)
		ah.PrintCodeSnippet(alert)
		ah.PrintMessage(alert)
		ah.PrintNote(alert)
		fmt.Println()
	}

	ah.Finish()
}

func (ah *AlertHandler) SortTypes() {
	//TODO: error and warning sorting
}

func (ah *AlertHandler) PrintMessage(alert Alert) {
	color.Println(alert.GetMessage())
}

func (ah *AlertHandler) PrintLocation(alert Alert, file string) {
	location := CombineLocations(alert.GetSpecifier().GetTokens())
	color.Printf("%s:%d:%d\n", file, location.Line.Start, location.Column.Start)
}

func (ah *AlertHandler) PrintCodeSnippet(alert Alert) {
	lineCount := 1
	columnCount := 0
	specifier := alert.GetSpecifier()
	location := CombineLocations(specifier.GetTokens()) // how will this work

	// Alert(&err, Multiple{tk, tk2, tk3}, Fix{"thing"}, Singleline{tk1})

	for i := 0; i < len(ah.Source); i++ {
		columnCount += 1

		if lineCount == location.Line.Start && (ah.Source[i] == '\n' || i == len(ah.Source)-1) {
			color.Println(specifier.GetSnippet(string(ah.Source), i, columnCount, lineCount))
			break
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
		color.Printf("[cyan] %s= note:[default] %s\n", strings.Repeat(" ", len(strconv.Itoa(ah.currentLine))), alert.GetNote())
		return
	}

	fmt.Print("\n")
}

func CombineLocations(tks []tokens.Token) tokens.TokenLocation {
	if len(tks) == 0 {
		return tokens.TokenLocation{}
	}
	location := tks[0].TokenLocation

	for i, v := range tks {
		loc := v.TokenLocation
		if i == 0 {
			continue
		}

		if loc.Column.Start < location.Column.Start {
			location.Column.Start = loc.Column.Start
		}
		if loc.Column.End > location.Column.End {
			location.Column.End = loc.Column.End
		}
		if loc.Line.Start < location.Line.Start {
			location.Line.Start = loc.Line.Start
		}
		if loc.Line.End > location.Line.End {
			location.Line.End = loc.Line.End
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
