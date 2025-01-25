package alerts

import (
	"fmt"
	"hybroid/helpers"
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
	GetNote() string
	GetID() string
	GetAlertType() AlertType
}

type AlertHandler struct {
	Source []byte

	Alerts []Alert

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
	ah.Alerts = append(ah.Alerts, ah.NewAlert(alertType, args...))
}

func (ah *AlertHandler) AlertI_(alertType Alert) {
	ah.Alerts = append(ah.Alerts, alertType)
}

func (ah *AlertHandler) PrintAlerts(alertStage AlertStage, sourcePath string) {
	//FIXME: ah.Source = source

	var errMsg string
	switch alertStage {
	case Lexer, Parser:
		errMsg = "[light_red]Syntax error[%s]: "
	case Walker, Eval:
		errMsg = "[light_red]Compilation error[%s]: "
	}

	for _, alert := range ah.Alerts {
		color.Printf(errMsg, alert.GetID())
		ah.printLocation(alert, sourcePath)
		ah.printCodeSnippet(alert)
		ah.printMessage(alert)
		ah.printNote(alert)
		fmt.Println()
	}

	ah.Finish()
}

func (ah *AlertHandler) sortTypes() {
	//TODO: error and warning sorting
}

func (ah *AlertHandler) printMessage(alert Alert) {
	color.Println(alert.GetMessage())
}

func (ah *AlertHandler) printLocation(alert Alert, file string) {
	location := combineLocations(alert.GetSpecifier().GetTokens())
	color.Printf("%s:%d:%d\n", file, location.Start, location.Start)
}

func (ah *AlertHandler) printCodeSnippet(alert Alert) {
	lineCount := 1
	columnCount := 0
	specifier := alert.GetSpecifier()
	location := combineLocations(specifier.GetTokens())

	for i := 0; i < len(ah.Source); i++ {
		columnCount += 1

		if lineCount == location.Start && (ah.Source[i] == '\n' || i == len(ah.Source)-1) {
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

func (ah *AlertHandler) printNote(alert Alert) {
	if alert.GetNote() != "" {
		color.Printf("[cyan] %s= note:[default] %s\n", strings.Repeat(" ", len(strconv.Itoa(ah.currentLine))), alert.GetNote())
		return
	}

	fmt.Print("\n")
}

func combineLocations(tks []tokens.Token) helpers.Span[int] {
	if len(tks) == 0 {
		return helpers.Span[int]{}
	}
	location := tks[0].Position

	for i, v := range tks {
		loc := v.Position
		if i == 0 {
			continue
		}

		if loc.Start < location.Start {
			location.Start = loc.Start
		}
		if loc.End > location.End {
			location.End = loc.End
		}
	}

	return location
}

func (ah *AlertHandler) getCodeSnippet(location helpers.Span[int]) string {
	return ""
}

func (ah *AlertHandler) Finish() {
	ah.Source = []byte{}
}
