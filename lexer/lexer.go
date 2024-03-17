package lexer

import (
	"fmt"
	"strconv"
)

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

	Eof // END OF FILE
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal string
	Line    int
}

func (t TokenType) ToString() string {
	return [...]string{
		"LeftParen", "RightParen", "LeftBrace", "RightBrace", "LeftBracket", "RightBracket", "Comma", "Colon", "At", "Dot", "Concat", "Minus", "MinusEqual", "Plus", "PlusEqual ", "Slash", "SlashEqual", "Star", "StarEqual ", "Caret", "CaretEqual", "Bang", "BangEqual", "Equal", "EqualEqual", "FatArrow", "Greater", "GreaterEqual", "Less", "LessEqual", "Identifier", "String", "Number", "FixedPoint", "Degree", "Radian", "And", "Or", "True", "False", "Self", "Fn", "Tick", "Repeat", "For", "While", "If", "Else", "Nil", "Return", "Break", "Continue", "Let", "Pub", "In", "As", "To", "With", "Enum", "Use", "Spawn", "Trait", "Entity", "Find", "Remove", "Match", "Eof",
	}[t]
}

func (t Token) ToString() string {
	return fmt.Sprintf("Token {type: %v, lex: %v, lit: %v, line: %v}", t.Type.ToString(), t.Lexeme, t.Literal, t.Line)
}

var tokens []Token
var start, current, line int
var source []byte

var passed bool

// patrons := map[int]string{
// 	0: "Terrence",
// 	1: "Evelyn",
// }

var keywords = map[string]TokenType{
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
}

func Advance() byte {
	t := source[current]
	current++
	return t
}

func IsAtEnd() bool {
	return current >= len(source)
}

func IsAtEndNext() bool {
	return current+1 >= len(source)
}

func AddToken(token TokenType, literal string) {
	text := string(source)[start:current]
	tokens = append(tokens, Token{token, text, literal, line})
}

func MatchChar(expected byte) bool {
	if IsAtEnd() {
		return false
	}
	if source[current] != expected {
		return false
	}

	current++
	return true
}

func Peek() byte {
	if IsAtEnd() {
		return '\f'
	}

	return source[current]
}

func PeekNext() byte {
	if IsAtEndNext() {
		return '0'
	}

	return source[current+1]
}

func IsDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func IsHexDigit(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')
}

func IsAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func IsAlphaNumeric(c byte) bool {
	return IsAlpha(c) || IsDigit(c)
}

func HandleString() {
	for Peek() != '"' && !IsAtEnd() {
		if Peek() == '\n' {
			line++
			LexerError("Multiline strings are not allowed.")
		}

		Advance()
	}

	if IsAtEnd() {
		LexerError("Unterminated string.")
		return
	}

	Advance()

	value := string(source)[start+1 : current-1]
	AddToken(String, value)
}

func HandleNumber() {
	if Peek() == 'x' {
		Advance()
		Advance()

		for IsHexDigit(Peek()) {
			Advance()
		}

		AddToken(Number, string(source[start:current]))

		return
	}

	for IsDigit(Peek()) {
		Advance()
	}

	if Peek() == '.' && IsDigit(PeekNext()) {
		Advance()

		for IsDigit(Peek()) {
			Advance()
		}
	}

	// parse a number to see if its a valid number

	strNum := string(source[start:current])
	if !TryParseNum(strNum) {
		LexerError(fmt.Sprintf("Invalid number: `%s`", strNum))
		return
	}
	// evaluate if its postfix: fx, r, d

	var postfix string
	postfixStart := current

	for IsAlpha(Peek()) {
		Advance()
	}

	postfix = string(source[postfixStart:current])
	switch postfix {

	case "f":
		AddToken(Number, strNum)
	case "fx":
		AddToken(FixedPoint, strNum)
	case "r":
		AddToken(Radian, strNum)
	case "d":
		AddToken(Degree, strNum)
	case "":
		AddToken(Number, strNum)
	default:
		LexerError(fmt.Sprintf("Invalid postfix: `%s`", postfix))
	}
}

func TryParseNum(strNum string) bool { //bytes: num
	_, ok := strconv.ParseFloat(strNum, 64)

	return ok == nil
}

func HandleIdentifier() {
	for IsAlphaNumeric(Peek()) {
		Advance()
	}

	text := string(source)[start:current]

	val, ok := keywords[text]
	if ok {
		AddToken(val, "")
		return
	}

	AddToken(Identifier, "")
}

func ScanToken() {
	c := Advance()

	switch c {

	case '{':
		AddToken(LeftBrace, "") // the literal is emplty because "{" is not a value
	case '}':
		AddToken(RightBrace, "")
	case '(':
		AddToken(LeftParen, "")
	case ')':
		AddToken(RightParen, "")
	case '[':
		AddToken(LeftBracket, "")
	case ']':
		AddToken(LeftBracket, "")
	case ',':
		AddToken(Comma, "")
	case ':':
		AddToken(Colon, "")
	case '@':
		AddToken(At, "")
	case '.':
		if MatchChar('.') {
			AddToken(Concat, "")
		} else {
			AddToken(Dot, "")
		}
	case '+':
		if MatchChar('=') {
			AddToken(PlusEqual, "")
		} else {
			AddToken(Plus, "")
		}
	case '-':
		if MatchChar('=') {
			AddToken(MinusEqual, "")
		} else {
			AddToken(Minus, "")
		}
	case '^':
		if MatchChar('=') {
			AddToken(CaretEqual, "")
		} else {
			AddToken(Caret, "")
		}
	case '*':
		if MatchChar('=') {
			AddToken(StarEqual, "")
		} else {
			AddToken(Star, "")
		}
	case '=':
		if MatchChar('=') {
			AddToken(EqualEqual, "")
		} else if MatchChar('>') {
			AddToken(FatArrow, "")
		} else {
			AddToken(Equal, "")
		}
	case '!':
		if MatchChar('=') {
			AddToken(BangEqual, "")
		} else {
			AddToken(Bang, "")
		}
	case '<':
		if MatchChar('=') {
			AddToken(LessEqual, "")
		} else {
			AddToken(Less, "")
		}
	case '>':
		if MatchChar('=') {
			AddToken(GreaterEqual, "")
		} else {
			AddToken(Greater, "")
		}

	case '/':
		if MatchChar('/') {
			for Peek() != '\n' && !IsAtEnd() {
				Advance()
			}
		} else if MatchChar('*') {
			// Handle multiLINE comment
			for (Peek() != '*' || PeekNext() != '/') && !IsAtEnd() {
				if Peek() == '\n' {
					line++
				}

				Advance()
			}

			Advance()
			Advance()
		} else {
			if MatchChar('=') {
				AddToken(SlashEqual, "")
			} else {
				AddToken(Slash, "")
			}
		}

	case ' ':
	case '\r':
	case '\t':
		break

	case '\n':
		line++

	case '"':
		HandleString()

	default:
		if IsDigit(c) {
			HandleNumber()
		} else if IsAlpha(c) {
			HandleIdentifier()
		} else {
			LexerError(fmt.Sprintf("Unexpected character: `%c`", c))
		}
	}
}

func Tokenize(src []byte) []Token {
	source = src
	line, start, current = 1, 0, 0
	passed = true
	tokens = make([]Token, 0)

	for {
		if IsAtEnd() {
			break
		}
		start = current
		ScanToken()
	}

	tokens = append(tokens, Token{Eof, "", "", line}) // END OF FILE
	return tokens
}

func LexerError(message string) {
	passed = false
	fmt.Printf("[Lexer Error] In file: %s  Line: %v  [ %s ]\n", "file.lc", line, message)
}
