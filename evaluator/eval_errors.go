package evaluator

type ErrorType int

const (
	Lexer ErrorType = iota
	Parser
	Walker
	Eval
)

type AlertType int

const (
	Error	AlertType = iota
	Warning
)

type Alert interface {
	GetMessage() string
	GetCodeSnippet() string
	GetNote() string
	GetAlertType() AlertType
	GetErrorType() ErrorType
}

type SomeError struct {
	
}

func (e SomeError) GetMessage() string {
	return "Something does not exist."
}

func (e SomeError) GetCodeSnippet() string {
	return ""
}

func (e SomeError) GetNote() string {
	return "Dont do that, do that."
}

func (e SomeError) GetAlertType() AlertType {
	return Error
}

func (e SomeError) GetErrorType() ErrorType {
	return Parser
}


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