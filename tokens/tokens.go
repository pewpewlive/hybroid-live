package tokens

import "hybroid/core"

// Make sure to run `go install golang.org/x/tools/cmd/stringer@latest` before trying to generate

//go:generate stringer -type=TokenType -linecomment

type TokenType int

const (
	// Unused tokens
	// At @
	// SemiColon ;
	// Pipe |
	// Ampersand &
	// LeftShift <<
	// LeftShiftEqual <<=
	// RightShift >>
	// RightShiftEqual >>=
	// Find
	// Macro

	// Tokens

	Hash            TokenType = iota // #
	LeftParen                        // (
	RightParen                       // )
	LeftBrace                        // {
	RightBrace                       // }
	LeftBracket                      // [
	RightBracket                     // ]
	Comma                            // ,
	Colon                            // :
	Dot                              // .
	Concat                           // ..
	Ellipsis                         // ...
	Minus                            // -
	MinusEqual                       // -=
	Plus                             // +
	PlusEqual                        // +=
	Slash                            // /
	SlashEqual                       // /=
	BackSlash                        // \
	BackSlashEqual                   // \=
	Star                             // *
	StarEqual                        // *=
	Caret                            // ^
	CaretEqual                       // ^=
	Bang                             // !
	BangEqual                        // !=
	Equal                            // =
	EqualEqual                       // ==
	FatArrow                         // =>
	ThinArrow                        // ->
	Greater                          // >
	GreaterEqual                     // >=
	Less                             // <
	LessEqual                        // <=
	Modulo                           // %
	ModuloEqual                      // %=
	LeftShift                        // <<
	LeftShiftEqual                   // <<=
	RightShift                       // >>
	RightShiftEqual                  // >>=
	Pipe                             // |
	PipeEqual                        // |=
	Ampersand                        // &
	AmpersandEqual                   // &=

	// Literals

	Degree     // degree
	Fixed      // fixed
	FixedPoint // fixedPoint
	Identifier // identifier
	Number     // number
	Radian     // radian
	String     // string

	// Keywords

	Is       // is
	Isnt     // isnt
	Alias    // alias
	And      // and
	As       // as
	Break    // break
	By       // by
	Const    // const
	Continue // continue
	Else     // else
	Entity   // entity
	Enum     // enum
	Env      // env
	False    // false
	Fn       // fn
	For      // for
	If       // if
	In       // in
	From     // from
	To       // to
	Let      // let
	Match    // match
	New      // new
	Or       // or
	Pub      // pub
	Repeat   // repeat
	Return   // return
	Self     // self
	Spawn    // spawn
	Struct   // struct
	Class    // class
	Tick     // tick
	True     // true
	Use      // use
	While    // while
	With     // with
	Yield    // yield
	Destroy  // destroy

	Eof // EOF (End of File)
)

var keywords = map[string]TokenType{
	"is":       Is,
	"isnt":     Isnt,
	"alias":    Alias,
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
	"fn":       Fn,
	"to":       To,
	"for":      For,
	"if":       If,
	"in":       In,
	"let":      Let,
	"match":    Match,
	"new":      New,
	"or":       Or,
	"pub":      Pub,
	"repeat":   Repeat,
	"return":   Return,
	"self":     Self,
	"spawn":    Spawn,
	"struct":   Struct,
	"class":    Class,
	"tick":     Tick,
	"true":     True,
	"use":      Use,
	"from":     From,
	"while":    While,
	"with":     With,
	"yield":    Yield,
	"destroy":  Destroy,
}

func KeywordToToken(keyword string) (TokenType, bool) {
	v, found := keywords[keyword]
	return v, found
}

type Location struct {
	Line   int
	Column core.Span[int]
}

func NewLocation(line, columnStart, columnEnd int) Location {
	return Location{
		Line:   line,
		Column: core.NewSpan(columnStart, columnEnd),
	}
}

type Token struct {
	Location

	Type    TokenType
	Lexeme  string
	Literal string
}

func NewToken(tokenType TokenType, lexeme, literal string, location Location) Token {
	return Token{
		Type:     tokenType,
		Lexeme:   lexeme,
		Literal:  literal,
		Location: location,
	}
}
