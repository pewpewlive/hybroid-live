package lexer

import "fmt"

type TokenType int

const (
	// Single character tokens
	LeftParen    TokenType = iota // (
	RightParen                    // )
	LeftBrace                     // {
	RightBrace                    // }
	LeftBracket                   // [
	RightBracket                  // ]
	Comma                         // ,
	At                            // @
	Pipe                          // |

	// One or two character tokens
	Colon        // :
	DoubleColon  // ::
	Dot          // .
	Concat       // ..
	Minus        // -
	MinusEqual   // -=
	Plus         // +
	PlusEqual    // +=
	Slash        // /
	SlashEqual   // /=
	Star         // *
	StarEqual    // *=
	Caret        // ^
	CaretEqual   // ^=
	Bang         // !
	BangEqual    // !=
	Equal        // =
	EqualEqual   // ==
	FatArrow     // =>
	Greater      // >
	GreaterEqual // >=
	Less         // <
	LessEqual    // <=
	Modulo       // %
	ModuloEqual  // %=

	// Literals
	Degree
	Fixed
	FixedPoint
	Identifier
	Number
	Radian
	String

	// Keywords
	Add
	And
	As
	Break
	By
	Const
	Continue
	Else
	Entity
	Enum
	Env
	False
	Find
	Fn
	For
	From
	If
	In
	Let
	Match
	New
	Or
	Pub
	Remove
	Repeat
	Return
	Self
	Spawn
	Struct
	Tick
	To
	True
	Use
	While
	With
	Yield

	Eof // EOF (End of File)
)

var keywords = map[string]TokenType{
	"add":      Add,
	"and":      And,
	"as":       As,
	"break":    Break,
	"by":       By,
	"const":    Const,
	"continue": Continue,
	"else":     Else,
	"entity":   Entity,
	"enum":     Enum,
	"env":      Env,
	"false":    False,
	"find":     Find,
	"fn":       Fn,
	"for":      For,
	"from":     From,
	"if":       If,
	"in":       In,
	"let":      Let,
	"match":    Match,
	"new":      New,
	"or":       Or,
	"pub":      Pub,
	"remove":   Remove,
	"repeat":   Repeat,
	"return":   Return,
	"self":     Self,
	"spawn":    Spawn,
	"struct":   Struct,
	"tick":     Tick,
	"to":       To,
	"true":     True,
	"use":      Use,
	"while":    While,
	"with":     With,
	"yield":    Yield,
}

var tokens = [...]string{"LeftParen", "RightParen", "LeftBrace", "RightBrace", "LeftBracket", "RightBracket", "Comma", "At", "Pipe", "Colon", "DoubleColon", "Dot", "Concat", "Minus", "MinusEqual", "Plus", "PlusEqual", "Slash", "SlashEqual", "Star", "StarEqual", "Caret", "CaretEqual", "Bang", "BangEqual", "Equal", "EqualEqual", "FatArrow", "Greater", "GreaterEqual", "Less", "LessEqual", "Modulo", "ModuloEqual", "Degree", "Fixed", "FixedPoint", "Identifier", "Number", "Radian", "String", "Add", "And", "As", "Break", "By", "Const", "Continue", "Else", "Entity", "Enum", "Env", "False", "Find", "Fn", "For", "From", "If", "In", "Let", "Match", "New", "Or", "Pub", "Remove", "Repeat", "Return", "Self", "Spawn", "Struct", "Tick", "To", "True", "Use", "While", "With", "Yield", "Eof"}

func (t TokenType) ToString() string {
	return tokens[t]
}

type TokenLocation struct {
	LineStart int
	ColStart  int
	LineEnd   int
	ColEnd    int
}

type Token struct {
	Type     TokenType
	Lexeme   string
	Literal  string
	Location TokenLocation
}

func (t Token) ToString() string {
	return fmt.Sprintf("Token (%v), Lex: '%v', Lit: '%v', Ln: %v, ColStart: %v, ColEnd: %v", t.Type.ToString(), t.Lexeme, t.Literal, t.Location.LineStart, t.Location.ColStart, t.Location.ColEnd)
}

func KeywordToToken(keyword string) (TokenType, bool) {
	token, ok := keywords[keyword]

	return token, ok
}
