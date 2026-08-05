package lexer

import (
	"bufio"
	"hybroid/alerts"
	"hybroid/core"
	"hybroid/tokens"
	"io"
	"strconv"
	"strings"
)

type Lexer struct {
	alerts.Collector

	buffer []byte
	source *bufio.Reader

	line    int
	column  int
	srcByte int
}

func NewLexer(reader io.Reader) Lexer {
	return Lexer{
		Collector: alerts.NewCollector(),
		buffer:    make([]byte, 0),
		source:    bufio.NewReader(reader),
		line:      1,
		column:    1,
		srcByte:   0,
	}
}

func (l *Lexer) Tokenize() ([]tokens.Token, error) {
	lexerTokens := make([]tokens.Token, 0)

	for {
		token, err := l.next()
		if err == io.EOF {
			newToken := tokens.NewToken(tokens.Eof, "eof", "", l.span())
			lexerTokens = append(lexerTokens, newToken)
			break
		} else if err != nil && token == nil {
			return nil, err
		} else if token == nil {
			continue
		}

		lexerTokens = append(lexerTokens, *token)
	}

	return lexerTokens, nil
}

func (l *Lexer) next() (*tokens.Token, error) {
	if err := l.consumeWhile(isWhitespace); err != nil {
		return nil, err
	}

	l.buffer = l.buffer[:0]

	c, err := l.advance()
	if err != nil {
		return nil, err
	}

	if isAlphabetical(c) {
		return l.handleIdentifier()
	}
	if isDigit(c) {
		return l.handleNumber()
	}

	token := tokens.Token{}
	token.Span = l.span()
	switch c {
	case '{':
		token.Type = tokens.LeftBrace
	case '}':
		token.Type = tokens.RightBrace
	case '(':
		token.Type = tokens.LeftParen
	case ')':
		token.Type = tokens.RightParen
	case '[':
		token.Type = tokens.LeftBracket
	case ']':
		token.Type = tokens.RightBracket
	case ',':
		token.Type = tokens.Comma
	case ':':
		token.Type = tokens.Colon
	case '#':
		token.Type = tokens.Hash
	case '.':
		if l.match('.') {
			if l.match('.') {
				token.Type = tokens.Ellipsis
			} else {
				token.Type = tokens.Concat
			}
		} else {
			token.Type = tokens.Dot
		}
	case '+':
		if l.match('=') {
			token.Type = tokens.PlusEqual
		} else {
			token.Type = tokens.Plus
		}
	case '-':
		if l.match('=') {
			token.Type = tokens.MinusEqual
		} else if l.match('>') {
			token.Type = tokens.ThinArrow
		} else {
			token.Type = tokens.Minus
		}
	case '^':
		if l.match('=') {
			token.Type = tokens.CaretEqual
		} else {
			token.Type = tokens.Caret
		}
	case '*':
		if l.match('=') {
			token.Type = tokens.StarEqual
		} else {
			token.Type = tokens.Star
		}
	case '=':
		if l.match('=') {
			token.Type = tokens.EqualEqual
		} else if l.match('>') {
			token.Type = tokens.FatArrow
		} else {
			token.Type = tokens.Equal
		}
	case '!':
		if l.match('=') {
			token.Type = tokens.BangEqual
		} else {
			token.Type = tokens.Bang
		}
	case '<':
		token.Type = tokens.Less
	case '>':
		token.Type = tokens.Greater
	case '%':
		if l.match('=') {
			token.Type = tokens.ModuloEqual
		} else {
			token.Type = tokens.Modulo
		}
	case '|':
		token.Type = tokens.Pipe
	case '&':
		token.Type = tokens.Ampersand
	case '/':
		if l.match('/') {
			err := l.handleComment(false)
			return nil, err
		} else if l.match('*') {
			err := l.handleComment(true)
			return nil, err
		} else {
			if l.match('=') {
				token.Type = tokens.SlashEqual
			} else {
				token.Type = tokens.Slash
			}
		}
	case '~':
		if l.match('=') {
			token.Type = tokens.TildeEqual
		} else {
			token.Type = tokens.Tilde
		}
	case '\\':
		if l.match('=') {
			token.Type = tokens.BackSlashEqual
		} else {
			token.Type = tokens.BackSlash
		}
	case '"':
		return l.handleString()
	default:
		char := string(c)
		token.Lexeme = token.Type.String()
		if isUtf8NonAscii(c) {
			l.source.UnreadByte()
			r, _, _ := l.source.ReadRune()
			char = string(r)
		}
		token.Span.UpdateEnd(l.srcByte)
		l.Report(alerts.NewUnsupportedCharacter(token.Span, char))
		return nil, nil
	}

	token.Lexeme = token.Type.String()
	token.Span.UpdateEnd(l.srcByte)

	return &token, nil
}

func (l *Lexer) handleString() (*tokens.Token, error) {
	startLine := l.line
	startColumn := l.column - 1

	token := tokens.Token{
		Type: tokens.String,
		Span: core.NewSpan(l.srcByte, l.srcByte, startLine, startColumn),
	}

	for !l.match('"') && !l.isEOF() {
		if !l.match('\\', '"') {
			l.advance()
		}
	}
	token.Lexeme = l.bufferString()
	if len(token.Lexeme) == 0 {
		return nil, nil
	}
	if token.Lexeme[len(token.Lexeme)-1] != '"' {
		token.Span.UpdateEnd(l.srcByte)
		l.Report(alerts.NewUnterminatedString(token.Span))
	} else {
		token.Literal = token.Lexeme[1 : len(token.Lexeme)-1]
		if strings.Contains(token.Literal, "\n") {
			l.Report(alerts.NewMultilineString(token.Span))
		}
	}

	return &token, nil
}

func (l *Lexer) handleNumber() (*tokens.Token, error) {
	token := tokens.Token{
		Type: tokens.Number,
		Span: core.NewSpan(l.srcByte, l.srcByte, l.line, l.column-1),
	}

	base, err := l.peek()
	if err != nil && err != io.EOF {
		return nil, err
	}
	if l.buffer[0] == '0' && (base == 'x' || base == 'b' || base == 'o') {
		l.advance()

		err := l.consumeWhile(isAlphanumeric)
		if err != nil {
			return nil, err
		}

		token.Span.UpdateEnd(l.srcByte)
		token.Lexeme = l.bufferString()

		isInRange := isDigit
		var baseStr string
		switch base {
		case 'x':
			isInRange = isHex
			baseStr = "hex"
		case 'b':
			isInRange = isBinary
			baseStr = "binary"
		case 'o':
			isInRange = isOctal
			baseStr = "octal"
		}
		isValidDigit := func(b byte) bool { return isInRange(b) || b == '_' }

		for i, r := range token.Lexeme[2:] {
			if !isValidDigit(byte(r)) {
				span := token.Span
				span.StartByte += i + 2
				span.UpdateEnd(span.StartByte + 1)
				l.Report(alerts.NewInvalidDigitInLiteral(span, string(r), baseStr))
				return &token, nil
			}
		}

		literal, err := strconv.ParseInt(token.Lexeme, 0, 0)
		if err != nil {
			l.Report(alerts.NewMalformedNumber(token.Span, token.Lexeme))
			return &token, nil
		}
		token.Literal = strconv.Itoa(int(literal))

		return &token, nil
	}

	isDigitOrUnderscore := func(b byte) bool { return isDigit(b) || b == '_' }
	err = l.consumeWhile(isDigitOrUnderscore)
	if err != nil {
		return nil, err
	}

	next, err := l.peek(2)
	if err != nil && err != io.EOF {
		return nil, err
	}
	if isDigit(next) && l.match('.') {
		err = l.consumeWhile(isDigitOrUnderscore)
		if err != nil {
			return nil, err
		}
	}

	token.Span.UpdateEnd(l.srcByte)
	token.Lexeme = l.bufferString()

	var literal float64
	if literal, err = strconv.ParseFloat(token.Lexeme, 64); err != nil {
		l.Report(alerts.NewMalformedNumber(token.Span, token.Lexeme))
		return nil, err
	}
	token.Literal = strconv.FormatFloat(literal, 'f', -1, 64)

	postixLocation := l.span()
	err = l.consumeWhile(isAlphabetical)
	if err != nil {
		return nil, err
	}
	postixLocation.UpdateEnd(l.srcByte)

	postfix := l.bufferString()
	switch postfix {
	case "f":
		token.Type = tokens.Fixed
	case "fx":
		token.Type = tokens.FixedPoint
	case "r":
		token.Type = tokens.Radian
	case "d":
		token.Type = tokens.Degree
	case "":
		break
	default:
		tokenCopy := token
		tokenCopy.Span = postixLocation
		l.Report(alerts.NewInvalidNumberPostfix(tokenCopy.Span, postfix))
	}

	token.Span.UpdateEnd(l.srcByte)

	return &token, nil
}

func (l *Lexer) handleIdentifier() (*tokens.Token, error) {
	token := tokens.Token{
		Type: tokens.Identifier,
		Span: core.NewSpan(l.srcByte, l.srcByte, l.line, l.column-1),
	}
	err := l.consumeWhile(isAlphanumeric)
	if err != nil {
		return nil, err
	}
	token.Span.UpdateEnd(l.srcByte)
	token.Lexeme = l.bufferString()

	if keyword, found := tokens.KeywordToToken(token.Lexeme); found {
		token.Type = keyword
	}

	return &token, nil
}

func (l *Lexer) handleComment(multiline bool) error {
	if !multiline {
		_, err := l.source.ReadBytes('\n')
		if err != nil {
			return err
		}
		l.line++
		l.column = 1
		return nil
	} else {
		for !l.match('*', '/') && !l.isEOF() {
			if l.match('/', '*') {
				l.handleComment(true)
			} else {
				_, err := l.advance()
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
}
