package lexer

import "fmt"

type TokenType string

const (
	// Tokens

	LeftParen      TokenType = "leftParen"      // (
	RightParen     TokenType = "rightParen"     // )
	LeftBrace      TokenType = "leftBrace"      // {
	RightBrace     TokenType = "rightBrace"     // }
	LeftBracket    TokenType = "leftBracket"    // [
	RightBracket   TokenType = "rightBracket"   // ]
	Comma          TokenType = "comma"          // ,
	At             TokenType = "at"             // @
	Pipe           TokenType = "pipe"           // |
	Colon          TokenType = "colon"          // :
	DoubleColon    TokenType = "doubleColon"    // ::
	Dot            TokenType = "dot"            // .
	Concat         TokenType = "concat"         // ..
	Minus          TokenType = "minus"          // -
	MinusEqual     TokenType = "minusEqual"     // -=
	Plus           TokenType = "plus"           // +
	PlusEqual      TokenType = "plusEqual"      // +=
	Slash          TokenType = "slash"          // /
	SlashEqual     TokenType = "slashEqual"     // /=
	BackSlash 	   TokenType = "backSlash" 		  // \
	BackSlashEqual TokenType = "backSlashEqual"	// \=
	Star           TokenType = "star"           // *
	StarEqual      TokenType = "starEqual"      // *=
	Caret          TokenType = "caret"          // ^
	CaretEqual     TokenType = "caretEqual"     // ^=
	Bang           TokenType = "bang"           // !
	BangEqual      TokenType = "bangEqual"      // !=
	Equal          TokenType = "equal"          // =
	EqualEqual     TokenType = "equalEqual"     // ==
	FatArrow       TokenType = "fatArrow"       // =>
	Greater        TokenType = "greater"        // >
	GreaterEqual   TokenType = "greaterEqual"   // >=
	Less           TokenType = "less"           // <
	LessEqual      TokenType = "lessEqual"      // <=
	Modulo         TokenType = "modulo"         // %
	ModuloEqual    TokenType = "moduloEqual"    // %=

	// Literals

	Degree     TokenType = "degree"
	Fixed      TokenType = "fixed"
	FixedPoint TokenType = "fixedPoint"
	Identifier TokenType = "identifier"
	Number     TokenType = "number"
	Radian     TokenType = "radian"
	String     TokenType = "string"

	// Keywords

	Add      TokenType = "add"
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
	From     TokenType = "from"
	If       TokenType = "if"
	In       TokenType = "in"
	Let      TokenType = "let"
	Macro    TokenType = "macro"
	Match    TokenType = "match"
	New      TokenType = "new"
	Or       TokenType = "or"
	Pub      TokenType = "pub"
	Remove   TokenType = "remove"
	Repeat   TokenType = "repeat"
	Return   TokenType = "return"
	Self     TokenType = "self"
	Spawn    TokenType = "spawn"
	Struct   TokenType = "struct"
	Tick     TokenType = "tick"
	To       TokenType = "to"
	True     TokenType = "true"
	Use      TokenType = "use"
	While    TokenType = "while"
	With     TokenType = "with"
	Yield    TokenType = "yield"
	Destroy  TokenType = "destroy"

	Eof TokenType = "eof" // EOF (End of File)
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
	"macro":    Macro,
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
	"destroy":  Destroy,
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
	return fmt.Sprintf("Token (%v), Lex: '%v', Lit: '%v', Ln: %v, ColStart: %v, ColEnd: %v", string(t.Type), t.Lexeme, t.Literal, t.Location.LineStart, t.Location.ColStart, t.Location.ColEnd)
}

func KeywordToToken(keyword string) (TokenType, bool) {
	v, found := keywords[keyword]
	return v, found
}
