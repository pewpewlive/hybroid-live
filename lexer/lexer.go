package lexer

import (
	"fmt"
)

type Lexer struct {
	Tokens []Token
	source []byte
	Errors []LexerError

	start, current, line, columnStart, columnCurrent int
}

func New() *Lexer {
	return &Lexer{make([]Token, 0), make([]byte, 0), make([]LexerError, 0), 0, 0, 1, 0, 0}
}

func (l *Lexer) AssignSource(src []byte) {
	l.source = src
}

func (l *Lexer) handleString() {
	for l.peek() != '"' && !l.isAtEnd() {
		if l.peek() == '\\' && l.peekNext() == '"' {
			l.advance()
		}
		if l.peek() == '\n' {
			l.line++
			l.columnStart = 0
			l.columnCurrent = 0
			l.lexerError("multiline strings are not allowed")
		}

		l.advance()
	}

	if l.isAtEnd() {
		l.lexerError("unterminated string")
		return
	}

	l.advance()

	value := string(l.source)[l.start+1 : l.current-1]
	l.addToken(String, value)
}

func (l *Lexer) handleNumber() {
	if l.peek() == 'x' {
		l.advance()
		l.advance()

		for isHexDigit(l.peek()) {
			l.advance()
		}

		l.addToken(Number, string(l.source[l.start:l.current]))

		return
	}

	for isDigit(l.peek()) {
		l.advance()
	}

	if l.peek() == '.' && isDigit(l.peekNext()) {
		l.advance()

		for isDigit(l.peek()) {
			l.advance()
		}
	}

	// Parse a number to see if its a valid number

	strNum := string(l.source[l.start:l.current])
	if !tryParseNum(strNum) {
		l.lexerError(fmt.Sprintf("invalid number `%s`", strNum))
		return
	}
	// Evaluate if it is a postfix: `fx`, `r`, `d`

	var postfix string
	postfixStart := l.current

	for isAlphabetical(l.peek()) {
		l.advance()
	}

	postfix = string(l.source[postfixStart:l.current])
	switch postfix {

	case "f":
		l.addToken(Fixed, strNum)
	case "fx":
		l.addToken(FixedPoint, strNum)
	case "r":
		l.addToken(Radian, strNum)
	case "d":
		l.addToken(Degree, strNum)
	case "":
		l.addToken(Number, strNum)
	default:
		l.lexerError(fmt.Sprintf("invalid postfix '%s'", postfix))
	}
}

func (l *Lexer) handleIdentifier() {
	for isAlphanumeric(l.peek()) {
		l.advance()
	}

	text := string(l.source)[l.start:l.current]

	val, ok := KeywordToToken(text)
	if ok {
		l.addToken(val, "")
		return
	}

	l.addToken(Identifier, "")
}

func (l *Lexer) scanToken() {
	c := l.advance()

	switch c {

	case '{':
		l.addToken(LeftBrace, "")
	case '}':
		l.addToken(RightBrace, "")
	case '(':
		l.addToken(LeftParen, "")
	case ')':
		l.addToken(RightParen, "")
	case '[':
		l.addToken(LeftBracket, "")
	case ']':
		l.addToken(RightBracket, "")
	case ',':
		l.addToken(Comma, "")
	case ':':
		l.addToken(Colon, "")
	case '@':
		l.addToken(At, "")
	case '|':
		l.addToken(Wall, "")
	case '.':
		if l.matchChar('.') {
			l.addToken(Concat, "")
		} else {
			l.addToken(Dot, "")
		}
	case '+':
		if l.matchChar('=') {
			l.addToken(PlusEqual, "")
		} else {
			l.addToken(Plus, "")
		}
	case '-':
		if l.matchChar('=') {
			l.addToken(MinusEqual, "")
		} else {
			l.addToken(Minus, "")
		}
	case '^':
		if l.matchChar('=') {
			l.addToken(CaretEqual, "")
		} else {
			l.addToken(Caret, "")
		}
	case '*':
		if l.matchChar('=') {
			l.addToken(StarEqual, "")
		} else {
			l.addToken(Star, "")
		}
	case '=':
		if l.matchChar('=') {
			l.addToken(EqualEqual, "")
		} else if l.matchChar('>') {
			l.addToken(FatArrow, "")
		} else {
			l.addToken(Equal, "")
		}
	case '!':
		if l.matchChar('=') {
			l.addToken(BangEqual, "")
		} else {
			l.addToken(Bang, "")
		}
	case '<':
		if l.matchChar('=') {
			l.addToken(LessEqual, "")
		} else {
			l.addToken(Less, "")
		}
	case '>':
		if l.matchChar('=') {
			l.addToken(GreaterEqual, "")
		} else {
			l.addToken(Greater, "")
		}
	case '%':
		if l.matchChar('=') {
			l.addToken(ModuloEqual, "")
		} else {
			l.addToken(Modulo, "")
		}

	case '/':
		if l.matchChar('/') {
			for l.peek() != '\n' && !l.isAtEnd() {
				l.advance()
			}
		} else if l.matchChar('*') {
			// Handle multiline comments
			for (l.peek() != '*' || l.peekNext() != '/') && !l.isAtEnd() {
				if l.peek() == '\n' {
					l.line++
					l.columnStart = 0
					l.columnCurrent = 0
				}

				l.advance()
			}

			l.advance()
			l.advance()
		} else {
			if l.matchChar('=') {
				l.addToken(SlashEqual, "")
			} else {
				l.addToken(Slash, "")
			}
		}

	// Whitespace characters
	case ' ', ';', '\r', '\t':
		break

	// Increment line count when hitting new line
	case '\n':
		l.line++
		l.columnStart = 0
		l.columnCurrent = 0

	case '"':
		l.handleString()

	default:
		if isDigit(c) {
			l.handleNumber()
		} else if isAlphabetical(c) {
			l.handleIdentifier()
		} else {
			l.lexerError(fmt.Sprintf("unexpected character '%c'", c))
		}
	}
}

func (l *Lexer) Tokenize() {
	for !l.isAtEnd() {
		l.start = l.current
		l.columnStart = l.columnCurrent
		l.scanToken()
	}

	l.Tokens = append(l.Tokens, Token{
		Eof, "", "", TokenLocation{
			LineStart: l.line,
			LineEnd:   l.line,
			ColStart:  l.columnCurrent + 1,
			ColEnd:    l.columnCurrent + 1,
		},
	}) // Append an EOF (End of File) token
}
