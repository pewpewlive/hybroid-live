package lexer

import (
	"hybroid/tokens"
	"strconv"
)

func (l *Lexer) advance() byte {
	t := l.source[l.current]
	l.current++
	l.columnCurrent++
	return t
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

func (l *Lexer) isAtEndNext() bool {
	return l.current+1 >= len(l.source)
}

func (l *Lexer) addToken(token tokens.TokenType, literal string) {
	text := string(l.source)[l.start:l.current]
	l.Tokens = append(l.Tokens, tokens.Token{Type: token, Lexeme: text, Literal: literal, Location: tokens.TokenLocation{LineStart: l.lineStart, LineEnd: l.lineCurrent, ColStart: l.columnStart + 1, ColEnd: l.columnCurrent + 1}})
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

func newLocation(lineStart, columnStart, lineCurrent, columnCurrent int) tokens.TokenLocation {
	return tokens.TokenLocation{LineStart: lineStart, ColStart: columnStart, LineEnd: lineCurrent, ColEnd: columnCurrent}
}
