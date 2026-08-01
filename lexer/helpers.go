package lexer

import (
	"io"
)

func (l *Lexer) advance() (byte, error) {
	b, err := l.source.ReadByte()
	if err != nil {
		return 0, err
	}
	l.buffer = append(l.buffer, b)
	if b == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	return b, nil
}

func (l *Lexer) consumeWhile(predicate func(byte) bool) error {
	var err error

	for b, err := l.peek(); err == nil && predicate(b); b, err = l.peek() {
		_, err := l.advance()
		if err != nil {
			return err
		}
	}

	return err
}

func (l *Lexer) bufferString() string {
	str := string(l.buffer)
	l.buffer = l.buffer[:0]
	return str
}

func (l *Lexer) peek(offset ...int) (byte, error) {
	peekOffset := 1
	if len(offset) == 1 && offset[0] > 0 {
		peekOffset = offset[0]
	}

	bytes, err := l.source.Peek(peekOffset)
	if err != nil && err != io.EOF {
		return 0, err
	}

	if len(bytes) < peekOffset {
		return 0, io.EOF
	}

	return bytes[peekOffset-1], nil
}

func (l *Lexer) isEOF() bool {
	_, err := l.peek()
	return err == io.EOF
}

func (l *Lexer) match(bytes ...byte) bool {
	if l.isEOF() {
		return false
	}

	for i, b := range bytes {
		b2, err := l.peek(i + 1)
		if err != nil || b != b2 {
			return false
		}
	}

	for range bytes {
		l.advance()
	}

	return true
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func isOctal(b byte) bool {
	return b >= '0' && b <= '7'
}

func isBinary(b byte) bool {
	return b == '0' || b == '1'
}

func isHex(b byte) bool {
	return isDigit(b) ||
		(b >= 'a' && b <= 'f') ||
		(b >= 'A' && b <= 'F')
}

func isAlphabetical(b byte) bool {
	return (b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') ||
		b == '_'
}

func isAlphanumeric(b byte) bool {
	return isAlphabetical(b) || isDigit(b)
}

func isWhitespace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\r' || b == '\t'
}
