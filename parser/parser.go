package parser

import (
	"hybroid/ast"
	"hybroid/lexer"
)

type Parser struct {
	program  []ast.Node
	current  int
	tokens   []lexer.Token
	Errors   []ast.Error
	Warnings []ast.Warning
}

func New() *Parser {
	return &Parser{make([]ast.Node, 0), 0, make([]lexer.Token, 0), make([]ast.Error, 0), make([]ast.Warning, 0)}
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

func (p *Parser) getParam() ast.Param {
	paramName := p.expression()
	paramType := p.Type()
	if paramName.GetType() != ast.Identifier {
		p.error(paramName.GetToken(), "expected an identifier in parameter")
	}
	return ast.Param{Type: paramType, Name: paramName.GetToken()}
}

func (p *Parser) parameters() []ast.Param {
	if _, ok := p.consume("expected opening paren after an identifier", lexer.LeftParen); !ok {
		return nil
	}

	var args []ast.Param
	if p.match(lexer.RightParen) {
		args = make([]ast.Param, 0)
	} else {
		args = append(args, p.getParam())
		for p.match(lexer.Comma) {
			args = append(args, p.getParam())
		}
		p.consume("expected closing paren after parameters", lexer.RightParen)
	}

	return args
}

func (p *Parser) arguments() []ast.Node {
	if _, ok := p.consume("expected opening paren", lexer.LeftParen); !ok {
		return nil
	}

	var args []ast.Node
	if p.match(lexer.RightParen) {
		args = make([]ast.Node, 0)
	} else {
		arg := p.expression()
		args = append(args, arg)
		for p.match(lexer.Comma) {
			arg := p.expression()
			args = append(args, arg)
		}
		p.consume("expected closing paren after arguments", lexer.RightParen)
	}

	return args
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
