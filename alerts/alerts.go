package alerts

import (
	"fmt"
	"hybroid/tokens"
	"reflect"
)

type AlertStage int

const (
	Lexer AlertStage = iota
	Parser
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
	GetAlertStage() AlertStage

	// AddTokens(...tokens.Token)
}

type AlertHandler struct {
	Alerts []Alert
}

func (ah *AlertHandler) Alert(alertType Alert, args ...any) {
	alert := reflect.ValueOf(alertType).Elem()

	for i, arg := range args {
		alert.Field(i).Set(reflect.ValueOf(arg))
	}

	ah.Alerts = append(ah.Alerts, alert.Addr().Interface().(Alert))
}

func (ah *AlertHandler) PrintAlerts() {
	for _, alert := range ah.Alerts {
		ah.PrintAlert(alert)
	}
}

func (ah *AlertHandler) PrintAlert(alert Alert) {
	ah.PrintMessage(alert)
	ah.PrintCodeSnippet(alert)
	ah.PrintNote(alert)
}

func (ah *AlertHandler) PrintMessage(alert Alert) {
	fmt.Printf("%s\n", alert.GetMessage())
}

func (ah *AlertHandler) PrintCodeSnippet(alert Alert) {

}

func (ah *AlertHandler) PrintNote(alert Alert) {
	if alert.GetNote() != "" {
		fmt.Printf("%s\n", alert.GetNote())
	}
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
