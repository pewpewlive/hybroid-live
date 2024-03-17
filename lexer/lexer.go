package lexer

import (
	"fmt"
	"strconv"
)

type LexerError struct {
	TokenType TokenType
	Line      int
	Column    int
	Message   string
}

func New(src []byte) Lexer {
	return Lexer{source: src}
}

func (l *Lexer) ChangeSrc(newSrc []byte) {
	l.source = newSrc
}

func (l *Lexer) lexerError(message string) {
	l.Errors = append(l.Errors, LexerError{Eof, l.line, l.column, message})
}

type Lexer struct {
	Tokens                       []Token
	start, current, line, column int
	source                       []byte
	Errors                       []LexerError
}

func (l *Lexer) advance() byte {
	t := l.source[l.current]
	l.current++
	l.column++
	return t
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

func (l *Lexer) isAtEndNext() bool {
	return l.current+1 >= len(l.source)
}

func (l *Lexer) addToken(token TokenType, literal string) {
	text := string(l.source)[l.start:l.current]
	l.Tokens = append(l.Tokens, Token{token, text, literal, l.line})
}

func (l *Lexer) matchChar(expected byte) bool {
	if l.isAtEnd() {
		return false
	}
	if l.source[l.current] != expected {
		return false
	}

	l.current++
	l.column++
	return true
}

func (l *Lexer) peek() byte {
	if l.isAtEnd() {
		return '\f'
	}

	return l.source[l.current]
}

func (l *Lexer) peekNext() byte {
	if l.isAtEndNext() {
		return '0'
	}

	return l.source[l.current+1]
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isHexDigit(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')
}

func isAlphabetical(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func isAlphanumeric(c byte) bool {
	return isAlphabetical(c) || isDigit(c)
}

func (l *Lexer) handleString() {
	for l.peek() != '"' && !l.isAtEnd() {
		if l.peek() == '\\' && l.peekNext() == '"' {
			l.advance()
		}
		if l.peek() == '\n' {
			l.line++
			l.column = 0
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
		l.addToken(Number, strNum)
	case "fx":
		l.addToken(FixedPoint, strNum)
	case "r":
		l.addToken(Radian, strNum)
	case "d":
		l.addToken(Degree, strNum)
	case "":
		l.addToken(Number, strNum)
	default:
		l.lexerError(fmt.Sprintf("invalid postfix `%s`", postfix))
	}
}

func tryParseNum(strNum string) bool { //bytes: num
	_, ok := strconv.ParseFloat(strNum, 64)

	return ok == nil
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
		l.addToken(LeftBrace, "") // the literal is empty because `{` is not a value
	case '}':
		l.addToken(RightBrace, "")
	case '(':
		l.addToken(LeftParen, "")
	case ')':
		l.addToken(RightParen, "")
	case '[':
		l.addToken(LeftBracket, "")
	case ']':
		l.addToken(LeftBracket, "")
	case ',':
		l.addToken(Comma, "")
	case ':':
		l.addToken(Colon, "")
	case '@':
		l.addToken(At, "")
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
					l.column = 0
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
	case ' ':
	case ';':
	case '\r':
	case '\t':
		break

	// Increment line count when hitting new line
	case '\n':
		l.line++
		l.column = 0

	case '"':
		l.handleString()

	default:
		if isDigit(c) {
			l.handleNumber()
		} else if isAlphabetical(c) {
			l.handleIdentifier()
		} else {
			l.lexerError(fmt.Sprintf("unexpected character `%c`", c))
		}
	}
}

func (l *Lexer) Tokenize() {
	l.line, l.start, l.current, l.column = 1, 0, 0, 0
	l.Tokens = make([]Token, 0)
	l.Errors = make([]LexerError, 0)

	for {
		if l.isAtEnd() {
			break
		}
		l.start = l.current
		l.scanToken()
	}

	l.Tokens = append(l.Tokens, Token{Eof, "", "", l.line}) // Append an EOF (End of File) token
	for _, lexerError := range l.Errors {
		fmt.Printf("Error: %v, at line: %v, column: %v", lexerError.Message, lexerError.Line, lexerError.Column)
	}
}
