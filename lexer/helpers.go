package lexer

import (
	"io"
	"unicode/utf8"
)

func (l *Lexer) advance() (rune, error) {
	r, _, err := l.source.ReadRune()
	if err != nil {
		return r, err
	}
	l.buffer = append(l.buffer, r)
	if r == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	return r, nil
}

func (l *Lexer) consumeWhile(predicate func(rune) bool) error {
	var err error

	for r, err := l.peek(); err == nil && predicate(r); r, err = l.peek() {
		_, err := l.advance()
		if err != nil {
			return err
		}
	}

	return err
}

func (l *Lexer) bufferString() string {
	str := string(l.buffer)
	l.buffer = nil
	return str
}

func (l *Lexer) peek(offset ...int) (rune, error) {
	peekOffset := 1
	if len(offset) == 1 {
		peekOffset = offset[0]
	}

	bytes, err := l.source.Peek(peekOffset * utf8.UTFMax)
	if err != nil && err != io.EOF {
		return utf8.RuneError, err
	}
	if len(bytes) < peekOffset {
		return utf8.RuneError, io.EOF
	}

	runes := make([]rune, 0)
	for {
		if len(bytes) == 0 || len(runes) == peekOffset {
			break
		}

		r, size := utf8.DecodeRune(bytes)
		runes = append(runes, r)
		if len(bytes)-size >= 0 {
			bytes = bytes[size:]
		}
	}

	return runes[peekOffset-1], nil
}

func (l *Lexer) isEOF() bool {
	_, err := l.peek()
	return err == io.EOF
}

func (l *Lexer) match(runes ...rune) bool {
	if l.isEOF() {
		return false
	}

	for i, r := range runes {
		r2, err := l.peek(i + 1)
		if err != nil || r != r2 {
			return false
		}
	}

	for range runes {
		l.advance()
	}

	return true
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isOctal(r rune) bool {
	return r >= '0' && r <= '7'
}

func isBinary(r rune) bool {
	return r == '0' || r == '1'
}

func isHex(r rune) bool {
	return isDigit(r) ||
		(r >= 'a' && r <= 'f') ||
		(r >= 'A' && r <= 'F')
}

func isAlphabetical(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		r == '_'
}

func isAlphanumeric(r rune) bool {
	return isAlphabetical(r) || isDigit(r)
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\r' || r == '\t'
}
