package parser

import (
	"hybroid/ast"
	"hybroid/lexer"
)

func (p *Parser) directiveCall() ast.Node {
	directiveNode := ast.DirectiveExpr{}

	ident, identOk := p.consume("expected identifier in directive call", lexer.Identifier)
	if !identOk {
		return &directiveNode
	}
	directiveNode.Identifier = ident
	directiveNode.Token = ident

	if _, ok := p.consume("expected '(' after directive call", lexer.LeftParen); !ok {
		return &directiveNode
	}

	directiveNode.Expr = p.expression()

	if _, ok := p.consume("expected ')' after directive call", lexer.RightParen); !ok {
		return &directiveNode
	}

	return &directiveNode
}