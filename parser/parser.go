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

func (p *Parser) getOp(opEqual lexer.Token) lexer.Token {
	switch opEqual.Type {
	case lexer.PlusEqual:
		return lexer.Token{Type: lexer.Plus, Location: opEqual.Location, Literal: opEqual.Literal, Lexeme: "+"}
	case lexer.MinusEqual:
		return lexer.Token{Type: lexer.Minus, Location: opEqual.Location, Literal: opEqual.Literal, Lexeme: "-"}
	case lexer.SlashEqual:
		return lexer.Token{Type: lexer.Slash, Location: opEqual.Location, Literal: opEqual.Literal, Lexeme: "/"}
	case lexer.StarEqual:
		return lexer.Token{Type: lexer.Star, Location: opEqual.Location, Literal: opEqual.Literal, Lexeme: "*"}
	case lexer.CaretEqual:
		return lexer.Token{Type: lexer.Caret, Location: opEqual.Location, Literal: opEqual.Literal, Lexeme: "^"}
	case lexer.ModuloEqual:
		return lexer.Token{Type: lexer.Modulo, Location: opEqual.Location, Literal: opEqual.Literal, Lexeme: "%"}
	default:
		return lexer.Token{}
	}
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
		if statement != nil {
			p.program = append(p.program, statement)
		}
	}

	return p.program
}
