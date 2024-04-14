package parser

import (
	"hybroid/ast"
	"hybroid/lexer"
)

// Appends an error to the ParserErrors
func (p *Parser) error(token lexer.Token, msg string) {
	errMsg := ast.Error{Token: token, Message: msg}
	p.Errors = append(p.Errors, errMsg)
	//panic(errMsg.Message)
}

func (p *Parser) warn(token lexer.Token, msg string) {
	warn := ast.Warning{Token: token, Message: msg}
	p.Warnings = append(p.Warnings, warn)
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		switch p.peek().Type {
		case lexer.For, lexer.Fn, lexer.If, lexer.Repeat, lexer.Tick,
			lexer.Return, lexer.Let, lexer.While, lexer.Pub, lexer.Const,
			lexer.Break, lexer.Continue, lexer.Add, lexer.Remove:
			return
		}

		p.advance()
	}
}

func (p *Parser) isMultiComparison() bool {
	return p.match(lexer.And, lexer.Or)
}

func (p *Parser) isComparison() bool {
	return p.match(lexer.Greater, lexer.GreaterEqual, lexer.Less, lexer.LessEqual, lexer.BangEqual, lexer.EqualEqual)
}

// Creates a BinaryExpr
func (p *Parser) createBinExpr(left ast.Node, operator lexer.Token, tokenType lexer.TokenType, lexeme string, right ast.Node) ast.Node {
	valueType := p.determineValueType(left, right)
	return ast.BinaryExpr{
		Left:      left,
		Operator:  lexer.Token{Type: tokenType, Lexeme: lexeme, Literal: "", Location: operator.Location},
		Right:     right,
		ValueType: valueType,
	}
}

// Checks if the current position the parser is at is the End Of File
func (p *Parser) isAtEnd() bool {
	return p.peek().Type == lexer.Eof
}

// Advances by one into the next token and returns the previous tokne before advancing
func (p *Parser) advance() lexer.Token {
	t := p.tokens[p.current]
	if p.current < len(p.tokens)-1 {
		p.current++
	}
	return t
}

// Peeks into the current token or peeks at the token that is offet from the current position by the given offset
func (p *Parser) peek(offset ...int) lexer.Token {
	if offset == nil {
		return p.tokens[p.current]
	} else {
		if p.current+offset[0] >= len(p.tokens)-1 {
			return p.tokens[p.current]
		}
		return p.tokens[p.current+offset[0]]
	}
}

// Checks if the current type is the specified token type. Returns false if it's the End Of File
func (p *Parser) check(tokenType lexer.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == tokenType
}

// Matches the given list of tokens and advances if they match.
func (p *Parser) match(types ...lexer.TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}

	return false
}

// Consumes a list of tokens, advancing if they match and returns true. Consume also advances if none of the tokens were able to match, and returns false
func (p *Parser) consume(message string, types ...lexer.TokenType) (lexer.Token, bool) {
	if p.isAtEnd() {
		token := p.peek()
		p.error(token, message)
		return token, false // error
	}
	for _, tokenType := range types {
		if p.check(tokenType) {
			return p.advance(), true
		}
	}
	token := p.advance()
	p.error(token, message)
	return token, false // error
}

// Checks if the value type is expected to be a fixedpoint
func IsFx(valueType ast.PrimitiveValueType) bool {
	return valueType == ast.FixedPoint || valueType == ast.Fixed || valueType == ast.Radian || valueType == ast.Degree
}
