package tokens

import "hybroid/helpers"

type TokenType string

const (
	// Tokens

	Hash           TokenType = "#"   // #
	LeftParen      TokenType = "("   // (
	RightParen     TokenType = ")"   // )
	LeftBrace      TokenType = "{"   // {
	RightBrace     TokenType = "}"   // }
	LeftBracket    TokenType = "["   // [
	RightBracket   TokenType = "]"   // ]
	Comma          TokenType = ","   // ,
	At             TokenType = "@"   // @
	Pipe           TokenType = "|"   // |
	Colon          TokenType = ":"   // :
	SemiColon      TokenType = ";"   // ;
	Dot            TokenType = "."   // .
	Concat         TokenType = ".."  // ..
	Ellipsis       TokenType = "..." // ...
	Minus          TokenType = "-"   // -
	MinusEqual     TokenType = "-="  // -=
	Plus           TokenType = "+"   // +
	PlusEqual      TokenType = "+="  // +=
	Slash          TokenType = "/"   // /
	SlashEqual     TokenType = "/="  // /=
	BackSlash      TokenType = "\\"  // \
	BackSlashEqual TokenType = "\\=" // \=
	Star           TokenType = "*"   // *
	StarEqual      TokenType = "*="  // *=
	Caret          TokenType = "^"   // ^
	CaretEqual     TokenType = "^="  // ^=
	Bang           TokenType = "!"   // !
	BangEqual      TokenType = "!="  // !=
	Equal          TokenType = "="   // =
	EqualEqual     TokenType = "=="  // ==
	FatArrow       TokenType = "=>"  // =>
	ThinArrow      TokenType = "->"  // ->
	Greater        TokenType = ">"   // >
	GreaterEqual   TokenType = ">="  // >=
	Less           TokenType = "<"   // <
	LessEqual      TokenType = "<="  // <=
	Modulo         TokenType = "%"   // %
	ModuloEqual    TokenType = "%="  // %=

	// Literals

	Degree     TokenType = "degree"
	Fixed      TokenType = "fixed"
	FixedPoint TokenType = "fixedPoint"
	Identifier TokenType = "identifier"
	Number     TokenType = "number"
	Radian     TokenType = "radian"
	String     TokenType = "string"

	// Keywords

	Is       TokenType = "is"
	Isnt     TokenType = "isnt"
	Alias    TokenType = "alias"
	And      TokenType = "and"
	As       TokenType = "as"
	Break    TokenType = "break"
	By       TokenType = "by"
	Const    TokenType = "const"
	Continue TokenType = "continue"
	Else     TokenType = "else"
	Entity   TokenType = "entity"
	Enum     TokenType = "enum"
	Env      TokenType = "env"
	False    TokenType = "false"
	Find     TokenType = "find"
	Fn       TokenType = "fn"
	For      TokenType = "for"
	If       TokenType = "if"
	In       TokenType = "in"
	From     TokenType = "from"
	To       TokenType = "to"
	Let      TokenType = "let"
	Macro    TokenType = "macro"
	Match    TokenType = "match"
	New      TokenType = "new"
	Or       TokenType = "or"
	Pub      TokenType = "pub"
	Repeat   TokenType = "repeat"
	Return   TokenType = "return"
	Self     TokenType = "self"
	Spawn    TokenType = "spawn"
	Struct   TokenType = "struct"
	Class    TokenType = "class"
	Tick     TokenType = "tick"
	True     TokenType = "true"
	Use      TokenType = "use"
	While    TokenType = "while"
	With     TokenType = "with"
	Yield    TokenType = "yield"
	Destroy  TokenType = "destroy"
	Type     TokenType = "type"

	Eof TokenType = "eof" // EOF (End of File)
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
	"find":     Find,
	"fn":       Fn,
	"to":       To,
	"for":      For,
	"if":       If,
	"in":       In,
	"let":      Let,
	"macro":    Macro,
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
	"type":     Type,
}

func KeywordToToken(keyword string) (TokenType, bool) {
	v, found := keywords[keyword]
	return v, found
}

type Location struct {
	Line   helpers.Span[int]
	Column helpers.Span[int]
}

func NewLocation(lineStart, lineEnd, columnStart, columnEnd int) Location {
	return Location{
		Line:   helpers.NewSpan(lineStart, lineEnd),
		Column: helpers.NewSpan(columnStart, columnEnd),
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
