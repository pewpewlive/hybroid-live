package lexer

import "strconv"

type LexerError struct {
	TokenType TokenType
	Location  TokenLocation
	Message   string
}

func (self LexerError) GetToken() Token {
	return Token{Type:self.TokenType, Location: self.Location}
}

func (self LexerError) GetHeader() string {
	return "[red]Error"
}

func (self LexerError) GetMessage() string {
	return self.Message
}


func (l *Lexer) lexerError(message string) {
	l.Errors = append(l.Errors, LexerError{Eof, TokenLocation{LineStart: l.line, LineEnd: l.line, ColStart: l.columnStart, ColEnd: l.columnCurrent}, message})
}

func (l *Lexer) advance() byte {
	t := l.source[l.current]
	l.current++
	l.columnCurrent++
	return t
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source.)
}

func (l *Lexer) isAtEndNext() bool {
	return l.current+1 >= len(l.source)
}

func (l *Lexer) addToken(token TokenType, literal string) {
	text := string(l.source)[l.start:l.current]
	l.Tokens = append(l.Tokens, Token{token, text, literal, TokenLocation{LineStart: l.line, LineEnd: l.line, ColStart: l.columnStart + 1, ColEnd: l.columnCurrent + 1}})
}

func (l *Lexer) matchChar(expected byte) bool {
	if l.isAtEnd() {
		return false
	}
	if l.source[l.current] != expected {
		return false
	}

	l.current++
	l.columnCurrent++
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

func tryParseNum(strNum string) bool {
	_, ok := strconv.ParseFloat(strNum, 64)

	return ok == nil
}
