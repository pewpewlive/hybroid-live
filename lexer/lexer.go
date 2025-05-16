package lexer

import (
	"bufio"
	"hybroid/alerts"
	"hybroid/tokens"
	"io"
	"strconv"
	"strings"
)

type Lexer struct {
	alerts.Collector

	buffer []rune
	source *bufio.Reader

	line   int
	column int
}

func NewLexer(reader io.Reader) Lexer {
	return Lexer{
		Collector: alerts.NewCollector(),
		buffer:    make([]rune, 0),
		source:    bufio.NewReader(reader),
		line:      1,
		column:    1,
	}
}

func (l *Lexer) Alert(alertType alerts.Alert, args ...any) {
	l.Alert_(alertType, args...)
}

func (l *Lexer) Tokenize() ([]tokens.Token, error) {
	lexerTokens := make([]tokens.Token, 0)

	for {
		token, err := l.next()
		if err == io.EOF {
			newToken := tokens.NewToken(tokens.Eof, "", "", tokens.NewLocation(l.line, l.column, l.column))
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

	l.buffer = make([]rune, 0)

	token := tokens.Token{}
	token.Line = l.line

	c, err := l.advance()
	token.Column.Start = l.column - 1

	if err != nil {
		return nil, err
	}

	if isAlphabetical(c) {
		return l.handleIdentifier()
	}

	if isDigit(c) {
		return l.handleNumber()
	}

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
	case '@':
		token.Type = tokens.At
	case '#':
		token.Type = tokens.Hash
	case '|':
		token.Type = tokens.Pipe
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
		if l.match('=') {
			token.Type = tokens.LessEqual
		} else {
			token.Type = tokens.Less
		}
	case '>':
		if l.match('=') {
			token.Type = tokens.GreaterEqual
		} else {
			token.Type = tokens.Greater
		}
	case '%':
		if l.match('=') {
			token.Type = tokens.ModuloEqual
		} else {
			token.Type = tokens.Modulo
		}
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
	case '\\':
		if l.match('=') {
			token.Type = tokens.BackSlashEqual
		} else {
			token.Type = tokens.BackSlash
		}
	case ';':
		token.Type = tokens.SemiColon
	case '"':
		return l.handleString()
	default:
		token.Lexeme = string(token.Type)
		token.Line = l.line
		token.Column.End = l.column
		l.Alert(&alerts.UnsupportedCharacter{}, alerts.NewSingle(token), string(c))
		return nil, nil
	}

	token.Lexeme = string(token.Type)
	token.Line = l.line
	token.Column.End = l.column

	return &token, nil
}

func (l *Lexer) handleString() (*tokens.Token, error) {
	token := tokens.Token{
		Type:     tokens.String,
		Location: tokens.NewLocation(l.line, l.column-1, l.column),
	}

	for !l.match('"') && !l.isEOF() {
		if !l.match('\\', '"') {
			l.advance()
		}
	}
	token.Lexeme = l.bufferString()
	token.Literal = token.Lexeme[1 : len(token.Lexeme)-1]
	token.Line = l.line
	token.Column.End = l.column

	if token.Lexeme[len(token.Lexeme)-1] != '"' && l.isEOF() {
		l.Alert(&alerts.UnterminatedString{}, alerts.NewSingle(token))
	} else if strings.Contains(token.Literal, "\n") {
		l.Alert(&alerts.MultilineString{}, alerts.NewSingle(token))
	}

	return &token, nil
}

func (l *Lexer) handleNumber() (*tokens.Token, error) {
	token := tokens.Token{
		Type:     tokens.Number,
		Location: tokens.NewLocation(l.line, l.column-1, l.column),
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

		token.Line = l.line
		token.Column.End = l.column
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
		isValidDigit := func(r rune) bool { return isInRange(r) || r == '_' }

		for i, r := range token.Lexeme[2:] {
			if !isValidDigit(r) {
				location := token.Location
				location.Column.Start += i + 2
				location.Column.End = location.Column.Start + 1
				l.Alert(&alerts.InvalidDigitInLiteral{}, alerts.NewSingle(tokens.NewToken(tokens.Eof, "", "", location)), string(r), baseStr)
				return &token, nil
			}
		}

		literal, err := strconv.ParseInt(token.Lexeme, 0, 0)
		if err != nil {
			l.Alert(&alerts.MalformedNumber{}, alerts.NewSingle(token), token.Lexeme)
			return &token, nil
		}
		token.Literal = strconv.Itoa(int(literal))

		return &token, nil
	}

	isDigitOrUnderscore := func(r rune) bool { return isDigit(r) || r == '_' }
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

	token.Line = l.line
	token.Column.End = l.column
	token.Lexeme = l.bufferString()

	var literal float64
	if literal, err = strconv.ParseFloat(token.Lexeme, 64); err != nil {
		l.Alert(&alerts.MalformedNumber{}, alerts.NewSingle(token), token.Lexeme)
		return nil, err
	}
	token.Literal = strconv.FormatFloat(literal, 'f', -1, 64)

	postixLocation := tokens.NewLocation(l.line, l.column, l.column)
	err = l.consumeWhile(isAlphabetical)
	if err != nil {
		return nil, err
	}
	postixLocation.Line = l.line
	postixLocation.Column.End = l.column

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
		token.Location = postixLocation
		l.Alert(&alerts.InvalidNumberPostfix{}, alerts.NewSingle(token), postfix)
	}

	return &token, nil
}

func (l *Lexer) handleIdentifier() (*tokens.Token, error) {
	token := tokens.Token{
		Type:     tokens.Identifier,
		Location: tokens.NewLocation(l.line, l.column-1, l.column),
	}
	err := l.consumeWhile(isAlphanumeric)
	if err != nil {
		return nil, err
	}
	token.Line = l.line
	token.Column.End = l.column
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
