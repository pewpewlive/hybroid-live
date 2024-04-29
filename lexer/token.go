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
	Colon                         // :
	At                            // @
	Wall                          // |

	// One or two character tokens
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
	Identifier
	String
	Number
	FixedPoint
	Degree
	Radian
	Fixed

	// Keywords
	By
	Add
	And
	Or
	True
	False
	Self
	Fn
	Tick
	Repeat
	For
	While
	If
	Else
	Nil
	Return
	Break
	Continue
	Let
	Pub
	Const
	In
	As
	To
	With
	Enum
	Use
	Spawn
	Trait
	Entity
	Find
	Remove
	Match
	From
	Struct
	Neww

	Eof // EOF (End of File)
)

var keywords = map[string]TokenType{
	"by":       By,
	"add":      Add,
	"and":      And,
	"or":       Or,
	"true":     True,
	"false":    False,
	"self":     Self,
	"fn":       Fn,
	"tick":     Tick,
	"repeat":   Repeat,
	"for":      For,
	"while":    While,
	"if":       If,
	"else":     Else,
	"nil":      Nil,
	"return":   Return,
	"break":    Break,
	"continue": Continue,
	"let":      Let,
	"pub":      Pub,
	"const":    Const,
	"in":       In,
	"as":       As,
	"to":       To,
	"with":     With,
	"enum":     Enum,
	"use":      Use,
	"spawn":    Spawn,
	"trait":    Trait,
	"entity":   Entity,
	"find":     Find,
	"remove":   Remove,
	"match":    Match,
	"from":     From,
	"struct":   Struct,
	"new":      Neww,
}

func (t TokenType) ToString() string {
	return [...]string{
		"LeftParen", "RightParen", "LeftBrace", "RightBrace", "LeftBracket", "RightBracket", "Comma", "Colon", "At", "Wall", "Dot", "Concat", "Minus", "MinusEqual", "Plus", "PlusEqual", "Slash", "SlashEqual", "Star", "StarEqual", "Caret", "CaretEqual", "Bang", "BangEqual", "Equal", "EqualEqual", "FatArrow", "Greater", "GreaterEqual", "Less", "LessEqual", "Identifier", "String", "Number", "FixedPoint", "Degree", "Radian", "Fixed", "By", "Add", "And", "Or", "True", "False", "Self", "Fn", "Tick", "Repeat", "For", "While", "If", "Else", "Nil", "Return", "Break", "Continue", "Let", "Pub", "Const", "In", "As", "To", "With", "Enum", "Use", "Spawn", "Trait", "Entity", "Find", "Remove", "Match", "From", "Struct", "New", "Eof",
	}[t]
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
