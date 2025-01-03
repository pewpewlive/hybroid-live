package lexer

import (
	"hybroid/alerts"
	"hybroid/tokens"
)

type Lexer struct {
	alerts.AlertHandler

	Tokens []tokens.Token
	source []byte

	start, current, lineStart, lineCurrent, columnStart, columnCurrent int
}

func NewLexer() Lexer {
	return Lexer{
		Tokens:        make([]tokens.Token, 0),
		source:        make([]byte, 0),
		start:         0,
		current:       0,
		lineStart:     1,
		lineCurrent:   1,
		columnStart:   0,
		columnCurrent: 0,
	}
}

func (l *Lexer) Alert(alertType alerts.Alert, args ...any) {
	l.Alert_(alertType, args...)
}

func (l *Lexer) AssignSource(src []byte) {
	l.source = src
}

func (l *Lexer) handleString() {
	hasMultilineStr := false
	stringStartCol := l.columnCurrent

	for l.peek() != '"' && !l.isAtEnd() {
		if l.peek() == '\\' && l.peekNext() == '"' {
			l.advance()
		}
		if l.peek() == '\n' {
			l.lineCurrent++
			l.columnStart = 0
			l.columnCurrent = 0
			hasMultilineStr = true
		}

		l.advance()
	}

	if l.isAtEnd() {
		l.Alert(&alerts.UnterminatedString{}, alerts.Singleline{Token: tokens.Token{Location: newLocation(l.lineStart, stringStartCol, l.lineStart, stringStartCol)}})
		return
	} else if hasMultilineStr {
		l.Alert(&alerts.MultilineString{}, alerts.Singleline{Token: tokens.Token{Location: newLocation(l.lineStart, stringStartCol, l.lineCurrent, l.columnCurrent+1)}})
	}

	l.advance()

	value := string(l.source)[l.start+1 : l.current-1]
	l.addToken(tokens.String, value)
}

func (l *Lexer) handleNumber() {
	if l.peek() == 'x' {
		l.advance()
		l.advance()

		for isHexDigit(l.peek()) {
			l.advance()
		}

		l.addToken(tokens.Number, string(l.source[l.start:l.current]))

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
		l.Alert(&alerts.MalformedNumber{}, alerts.Singleline{Token: tokens.Token{Location: newLocation(l.lineStart, l.columnStart, l.lineCurrent, l.columnCurrent)}})
		return
	}
	// Evaluate if it is a postfix: `fx`, `r`, `d`

	var postfix string
	postfixStart := l.current
	postfixColumn := l.columnCurrent + 1

	for isAlphabetical(l.peek()) {
		l.advance()
	}

	postfix = string(l.source[postfixStart:l.current])
	switch postfix {

	case "f":
		l.addToken(tokens.Fixed, strNum)
	case "fx":
		l.addToken(tokens.FixedPoint, strNum)
	case "r":
		l.addToken(tokens.Radian, strNum)
	case "d":
		l.addToken(tokens.Degree, strNum)
	case "":
		l.addToken(tokens.Number, strNum)
	default:
		l.Alert(&alerts.InvalidNumberPostfix{}, alerts.Singleline{Token: tokens.Token{Location: newLocation(l.lineStart, postfixColumn, l.lineCurrent, l.columnCurrent)}}, postfix)
	}
}

func (l *Lexer) handleIdentifier() {
	for isAlphanumeric(l.peek()) {
		l.advance()
	}

	text := string(l.source)[l.start:l.current]

	val, ok := tokens.KeywordToToken(text)
	if ok {
		l.addToken(val, "")
		return
	}

	l.addToken(tokens.Identifier, "")
}

func (l *Lexer) scanToken() {
	c := l.advance()

	switch c {

	case '{':
		l.addToken(tokens.LeftBrace, "")
	case '}':
		l.addToken(tokens.RightBrace, "")
	case '(':
		l.addToken(tokens.LeftParen, "")
	case ')':
		l.addToken(tokens.RightParen, "")
	case '[':
		l.addToken(tokens.LeftBracket, "")
	case ']':
		l.addToken(tokens.RightBracket, "")
	case ',':
		l.addToken(tokens.Comma, "")
	case ':':
		if l.matchChar(':') {
			l.addToken(tokens.DoubleColon, "")
		} else {
			l.addToken(tokens.Colon, "")
		}
	case '@':
		l.addToken(tokens.At, "")
	case '#':
		l.addToken(tokens.Hash, "")
	case '|':
		l.addToken(tokens.Pipe, "")
	case '.':
		if l.matchChar('.') {
			if l.matchChar('.') {
				l.addToken(tokens.DotDotDot, "")
			} else {
				l.addToken(tokens.Concat, "")
			}
		} else {
			l.addToken(tokens.Dot, "")
		}
	case '+':
		if l.matchChar('=') {
			l.addToken(tokens.PlusEqual, "")
		} else {
			l.addToken(tokens.Plus, "")
		}
	case '-':
		if l.matchChar('=') {
			l.addToken(tokens.MinusEqual, "")
		} else if l.matchChar('>') {
			l.addToken(tokens.ThinArrow, "")
		} else {
			l.addToken(tokens.Minus, "")
		}
	case '^':
		if l.matchChar('=') {
			l.addToken(tokens.CaretEqual, "")
		} else {
			l.addToken(tokens.Caret, "")
		}
	case '*':
		if l.matchChar('=') {
			l.addToken(tokens.StarEqual, "")
		} else {
			l.addToken(tokens.Star, "")
		}
	case '=':
		if l.matchChar('=') {
			l.addToken(tokens.EqualEqual, "")
		} else if l.matchChar('>') {
			l.addToken(tokens.FatArrow, "")
		} else {
			l.addToken(tokens.Equal, "")
		}
	case '!':
		if l.matchChar('=') {
			l.addToken(tokens.BangEqual, "")
		} else {
			l.addToken(tokens.Bang, "")
		}
	case '<':
		if l.matchChar('=') {
			l.addToken(tokens.LessEqual, "")
		} else {
			l.addToken(tokens.Less, "")
		}
	case '>':
		if l.matchChar('=') {
			l.addToken(tokens.GreaterEqual, "")
		} else {
			l.addToken(tokens.Greater, "")
		}
	case '%':
		if l.matchChar('=') {
			l.addToken(tokens.ModuloEqual, "")
		} else {
			l.addToken(tokens.Modulo, "")
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
					l.lineCurrent++
					l.columnStart = 0
					l.columnCurrent = 0
				}

				l.advance()
			}

			l.advance()
			l.advance()
		} else {
			if l.matchChar('=') {
				l.addToken(tokens.SlashEqual, "")
			} else {
				l.addToken(tokens.Slash, "")
			}
		}

	case '\\':
		if l.matchChar('=') {
			l.addToken(tokens.BackSlashEqual, "")
		} else {
			l.addToken(tokens.BackSlash, "")
		}
	case ';':
		l.addToken(tokens.SemiColon, "")
	// Whitespace characters
	case ' ', '\r', '\t':
		break

	// Increment lineCurrent count when hitting new lineCurrent
	case '\n':
		l.lineCurrent++
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
			l.Alert(&alerts.UnsupportedCharacter{}, alerts.Singleline{Token: tokens.Token{Location: newLocation(l.lineStart, l.columnStart, l.lineCurrent, l.columnCurrent)}}, c)
		}
	}
}

func (l *Lexer) Tokenize() {
	for !l.isAtEnd() {
		l.start = l.current
		l.lineStart = l.lineCurrent
		l.columnStart = l.columnCurrent
		l.scanToken()
	}

	l.Tokens = append(l.Tokens, tokens.Token{
		Type:    tokens.Eof,
		Lexeme:  "",
		Literal: "",
		Location: tokens.TokenLocation{
			LineStart: l.lineCurrent,
			LineEnd:   l.lineCurrent,
			ColStart:  l.columnCurrent + 1,
			ColEnd:    l.columnCurrent + 1,
		},
	}) // Append an EOF (End of File) token
}
