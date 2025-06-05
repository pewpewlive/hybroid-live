package parser

import (
	"hybroid/alerts"
	"hybroid/tokens"
	"slices"
)

func (p *Parser) isMultiComparison() bool {
	return p.match(tokens.And, tokens.Or)
}

func (p *Parser) isComparison() (tokens.Token, bool) {
	op := tokens.Token{}

	isGreaterComp := p.check(tokens.Greater)
	isLessComp := p.check(tokens.Less)
	if isGreaterComp || isLessComp {
		compToken := p.peek()
		if p.peek(1).Type == tokens.Less || p.peek(1).Type == tokens.Greater {
			return op, false
		}
		if p.peek(1).Type == tokens.Equal {
			newToken, success := p.combineTokens(tokens.TokenType(int(compToken.Type)+1), 2)
			if !success {
				return op, false
			}
			return newToken, true
		}
		p.advance()
		return compToken, true
	}
	if !isGreaterComp && !isLessComp {
		return p.peek(), p.match(tokens.BangEqual, tokens.EqualEqual)
	}
	return op, false
}

// Checks if the current position the parser is at is the End Of File
func (p *Parser) isAtEnd() bool {
	return p.peek().Type == tokens.Eof
}

// Advances by the given offset into the next token and returns the previous token before advancing
func (p *Parser) advance(offset ...int) tokens.Token {
	currentOffset := 1
	if offset != nil {
		currentOffset = offset[0]
	}

	if currentOffset < 0 {
		panic("Attempt to advance with a negative offset. Use disadvance() instead!")
	}

	t := p.tokens[p.current]
	index := p.current + currentOffset

	if index < len(p.tokens) {
		p.current = index
	}

	return t
}

// Disadvances by the given offset into the previous tokens and returns the current token after disadvancing
func (p *Parser) disadvance(offset ...int) tokens.Token {
	currentOffset := 1
	if offset != nil {
		currentOffset = offset[0]
	}

	if currentOffset < 0 {
		panic("Attempt to disadvance with a negative offset (which moves forward). Use advance() instead!")
	}

	index := p.current - currentOffset

	if index >= 0 {
		p.current = index
	}

	return p.tokens[p.current]
}

// Peeks into the current token or peeks at the token that is offset from the current position by the given offset
func (p *Parser) peek(offset ...int) tokens.Token {
	index := p.current
	if offset != nil {
		index += offset[0]
	}

	if index >= 0 && index < len(p.tokens) {
		return p.tokens[index]
	} else {
		return p.tokens[p.current]
	}
}

// Checks if the current type is the specified token type. Returns false if it's the End Of File
func (p *Parser) check(tokens ...tokens.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return slices.Contains(tokens, p.peek().Type)
}

// Matches the given list of tokens and advances if they match.
func (p *Parser) match(tokens ...tokens.TokenType) bool {
	if p.check(tokens...) {
		p.advance()
		return true
	}

	return false
}

func (p *Parser) consumeTill(context string, start tokens.Token, types ...tokens.TokenType) bool {
	if p.isAtEnd() {
		p.AlertMulti(&alerts.ExpectedSymbol{}, start, p.peek(-1), types[0], context)
		return false
	}

	return !p.match(types...)
}

// Consumes one of the tokens in the given list and advances if it matches.
func (p *Parser) consume(alert alerts.Alert, typ tokens.TokenType) (tokens.Token, bool) {
	if p.isAtEnd() {
		p.AlertI(alert)
		return p.peek(), false // error
	}
	if p.check(typ) {
		return p.advance(), true
	}
	p.AlertI(alert)
	return p.peek(), false // error
}

// Helper function to run p.consume, with simpler alert creation
func (p *Parser) alertSingleConsume(alert alerts.Alert, token tokens.TokenType, args ...any) (tokens.Token, bool) {
	args = append([]any{alerts.NewSingle(p.peek()), token}, args...)
	return p.consume(p.NewAlert(alert, args...), token)
}

// Helper function to run p.consume, with simpler alert creation
func (p *Parser) alertMultiConsume(alert alerts.Alert, start, end tokens.Token, token tokens.TokenType, args ...any) (tokens.Token, bool) {
	args = append([]any{alerts.NewMulti(start, end), token}, args...)
	return p.consume(p.NewAlert(alert, args...), token)
}
