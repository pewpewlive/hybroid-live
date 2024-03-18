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
	Body []any
}

type VariableDeclarationStmt struct {
	Identifier string
	Expression any
}

type AssignmentExpr struct {
	Asignee any
	Value   any
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

type ParserError struct {
	token lexer.Token
	err   string
}

type Parser struct {
	current int
	tokens  []lexer.Token
	Errors  []ParserError
}

func New() Parser {
	return Parser{}
}

func (p *Parser) error(token lexer.Token, err string) *ParserError {
	p.Errors = append(p.Errors, ParserError{
		token,
		err,
	})
	return &p.Errors[len(p.Errors)-1]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == lexer.Eof
}

func (p *Parser) advance() lexer.Token {
	t := p.tokens[p.current]
	p.current++
	return t
}

func (p *Parser) peek() lexer.Token {
	return p.tokens[p.current]
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

func (p *Parser) consume(tokenType lexer.TokenType, message string) (lexer.Token, *ParserError) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	return lexer.Token{}, p.error(p.tokens[p.current], message) // error
}

func (p *Parser) statement() any {
	switch p.peek().Type {
	case lexer.Let:
		p.advance()
		return p.variableDeclaration()
	}

	return p.expression()
}

func (p *Parser) variableDeclaration() any {
	ident, err1 := p.consume(lexer.Identifier, "Expected identifier in variable declaration.")
	if err1 != nil {
		return err1
	}

	_, err2 := p.consume(lexer.Equal, "Expected equal token following identifier in variable declaration.")
	if err2 != nil {
		return err2
	}

	return VariableDeclarationStmt{ident.Lexeme, p.expression()}
}

func (p *Parser) expression() any {
	return p.assignment()
}

func (p *Parser) assignment() any {
	expr := p.equality()

	if p.match(lexer.Equal) {
		value := p.assignment()
		expr = AssignmentExpr{expr, value}
	}

	return expr
}

func (p *Parser) equality() any {
	expr := p.comparison()

	if p.match(lexer.BangEqual, lexer.EqualEqual) {
		operator := p.previous()
		right := p.comparison()
		expr = BinaryExpr{expr, operator, right}
	}

	return expr
}

func (p *Parser) comparison() any {
	expr := p.term()

	if p.match(lexer.Greater, lexer.GreaterEqual, lexer.Less, lexer.LessEqual) {
		operator := p.previous()
		right := p.term()
		expr = BinaryExpr{expr, operator, right}
	}

	return expr
}

func (p *Parser) term() any { // 1 - 10
	expr := p.factor()

	if p.match(lexer.Plus, lexer.Minus) {
		operator := p.previous()
		right := p.term()
		expr = BinaryExpr{expr, operator, right}
	}

	return expr
}

func (p *Parser) factor() any {
	expr := p.unary()

	if p.match(lexer.Star, lexer.Slash) {
		operator := p.previous()
		right := p.factor()
		expr = BinaryExpr{expr, operator, right}
	}

	return expr
}

func (p *Parser) unary() any {
	if p.match(lexer.Bang, lexer.Minus) {
		operator := p.previous()
		right := p.unary()
		return UnaryExpr{operator, right}
	}

	return p.primary()
}

func (p *Parser) primary() any {
	if p.match(lexer.False) {
		return LiteralExpr{"false"}
	}
	if p.match(lexer.True) {
		return LiteralExpr{"true"}
	}
	if p.match(lexer.Nil) {
		return LiteralExpr{"nil"}
	}

	if p.match(lexer.Number, lexer.FixedPoint, lexer.Degree, lexer.Radian, lexer.String) {
		return LiteralExpr{p.previous().Literal}
	}

	if p.match(lexer.Identifier) {
		return IdentifierExpr{p.previous().Lexeme}
	}

	if p.match(lexer.LeftParen) {
		expr := p.expression()
		p.consume(lexer.RightParen, "Expect \")\" after expression.")
		return GroupingExpr{expr}
	}

	return p.error(p.peek(), "Expected expression.")
}

func (p *Parser) ParseTokens(tokens []lexer.Token) Program {
	p.tokens = tokens

	program := Program{}

	for !p.isAtEnd() {
		stmt := p.statement()
		program.Body = append(program.Body, stmt) // anyscript
	}

	return program
}
