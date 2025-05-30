package alerts

import "fmt"

type Type int

const (
	Error Type = iota
	Warning
)

func (t Type) GetColor() string {
	switch t {
	case Error:
		return "light_red"
	case Warning:
		return "light_yellow"
	default:
		panic(fmt.Sprintf("unexpected alert type %v", t))
	}
}

/*
TODO: add fix snippet
"fix": {
      "insert": "number",
      "where": "before"
    }
*/

type Alert interface {
	Message() string
	SnippetSpecifier() Snippet
	Note() string
	ID() string
	AlertType() Type
}
