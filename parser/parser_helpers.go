package parser

import (
	"hybroid/ast"
	"hybroid/lexer"
)

type ParserError struct {
	Token   lexer.Token
	Message string
}

func (p *Parser) error(token lexer.Token, err string) {
	p.Errors = append(p.Errors, ParserError{token, err})
}

func (p *Parser) createBinExpr(left ast.Node, operator lexer.Token, tokenType lexer.TokenType, lexeme string, right ast.Node) ast.Node {
	valueType := p.determineValueType(left, right)
	return ast.BinaryExpr{
		Left:      left,
		Operator:  lexer.Token{Type: tokenType, Lexeme: lexeme, Literal: "", Location: operator.Location},
		Right:     right,
		ValueType: valueType,
	}
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == lexer.Eof
}

func (p *Parser) advance() lexer.Token {
	t := p.tokens[p.current]
	if p.current < len(p.tokens)-1 {
		p.current++
	}
	return t
}

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

func isFx(valueType ast.PrimitiveValueType) bool {
	return valueType == ast.FixedPoint || valueType == ast.Fixed || valueType == ast.Radian || valueType == ast.Degree
}
