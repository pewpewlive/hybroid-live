package parser

import (
	"hybroid/lexer"
)

type Node struct {
	Node any
}

type NodeStatement struct {
	Statement any
}

type Program struct {
	Body []NodeStatement
}

type NodeVariableDeclarationStmt struct {
	Identifier string
	Expression any
}

type LiteralExpr struct {
	Value any
}

type UnaryExpr struct {
	Operator lexer.Token
	Right    any
}

type BinaryExpr struct {
	Left     any
	Operator lexer.Token
	Right    any
}

type GroupingExpr struct {
	Expression any
}

type IdentifierExpr struct {
	Symbol string
}

type Parser struct {
	current int
	tokens  []lexer.Token
}

func New() Parser {
	return Parser{}
}

func (p *Parser) isAtEnd(token lexer.Token) bool {
	return token.Type == lexer.Eof
}

func (p *Parser) advance() lexer.Token {
	t := p.tokens[p.current]
	p.current++
	return t
}

func (p *Parser) peek(offset ...int) lexer.Token {
	if len(offset) >= 1 {
		return p.tokens[p.current+offset[0]]
	} else {
		return p.tokens[p.current]
	}
}

func (p *Parser) previous() lexer.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) check(tokenType lexer.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == tokenType
}

func (p *Parser) match(types ...lexer.TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) consume(tokenType lexer.TokenType, message string) lexer.Token {
	if p.check(tokenType) {
		return p.advance()
	}

	return p.error(p.peek(), message) // error
}

func (p *Parser) primary() any {
	if p.check(lexer.False) {
		return LiteralExpr{"false"}
	}
	if p.check(lexer.True) {
		return LiteralExpr{"true"}
	}
	if p.check(lexer.Nil) {
		return LiteralExpr{"nil"}
	}

	if p.match(lexer.Number, lexer.FixedPoint, lexer.Degree, lexer.Radian, lexer.String) {
		return LiteralExpr{p.previous().Literal}
	}

	if p.check(lexer.LeftParen) {
		expr := p.expression()
		p.consume(lexer.RightParen, "Expect \")\" after expression.")
		return expr
	}

	return p.error(peek, "Expected expression.")
}

func (p *Parser) ParseTokens(tokens []lexer.Token) {
	p.tokens = tokens
	for i := 0; i < len(tokens); i++ {
		if p.isAtEnd(p.peek()) {
			return
		}
	}
}
