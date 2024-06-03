package parser

import (
	"hybroid/ast"
	"hybroid/lexer"
)

type Parser struct {
	program []ast.Node
	current int
	tokens  []lexer.Token
	Errors  []ast.Error
}

func New() *Parser {
	return &Parser{make([]ast.Node, 0), 0, make([]lexer.Token, 0), make([]ast.Error, 0)}
}

func (p *Parser) AssignTokens(tokens []lexer.Token) {
	p.tokens = tokens
}

func (p *Parser) ParseTokens() []ast.Node {
	// Expect environment directive call as the first node
	statement := p.statement()
	if !p.verifyEnvironmentDirective(statement) {
		return p.program
	}
	p.program = append(p.program, statement)

	for !p.isAtEnd() {
		statement := p.statement()
		if statement.GetType() != ast.NA {
			p.program = append(p.program, statement)
		}
	}

	return p.program
}
