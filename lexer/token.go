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

	// Literals
	Identifier
	String
	Number
	FixedPoint
	Degree
	Radian

	// Keywords
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

	Eof // EOF (End of File)
)

func (t TokenType) ToString() string {
	return [...]string{
		"LeftParen", "RightParen", "LeftBrace", "RightBrace", "LeftBracket", "RightBracket", "Comma", "Colon", "At", "Dot", "Concat", "Minus", "MinusEqual", "Plus", "PlusEqual ", "Slash", "SlashEqual", "Star", "StarEqual ", "Caret", "CaretEqual", "Bang", "BangEqual", "Equal", "EqualEqual", "FatArrow", "Greater", "GreaterEqual", "Less", "LessEqual", "Identifier", "String", "Number", "FixedPoint", "Degree", "Radian", "And", "Or", "True", "False", "Self", "Fn", "Tick", "Repeat", "For", "While", "If", "Else", "Nil", "Return", "Break", "Continue", "Let", "Pub", "In", "As", "To", "With", "Enum", "Use", "Spawn", "Trait", "Entity", "Find", "Remove", "Match", "Eof",
	}[t]
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal string
	Line    int
}

func (t Token) ToString() string {
	return fmt.Sprintf("Token {type: %v, lex: %v, lit: %v, line: %v}", t.Type.ToString(), t.Lexeme, t.Literal, t.Line)
}

func KeywordToToken(keyword string) (TokenType, bool) {
	token, ok := map[string]TokenType{
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
	}[keyword]

	return token, ok
}
