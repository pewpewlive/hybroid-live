package parser

import (
	"hybroid/ast"
	"hybroid/lexer"
)

func (p *Parser) directiveCall() ast.Node {
	directiveNode := ast.DirectiveExpr{}

	ident, identOk := p.consume("expected identifier in directive call", lexer.Identifier)
	if !identOk {
		return directiveNode
	}
	directiveNode.Identifier = ident.Lexeme
	directiveNode.Token = ident

	if _, ok := p.consume("expected '(' after directive call", lexer.LeftParen); !ok {
		return directiveNode
	}

	directiveNode.Expr = p.expression()

	if _, ok := p.consume("expected ')' after directive call", lexer.RightParen); !ok {
		return directiveNode
	}

	return directiveNode
}

func (p *Parser) verifyEnvironmentDirective(statement ast.Node) bool {
	if _, ok := statement.(ast.DirectiveExpr); !ok {
		p.error(lexer.Token{Type: lexer.Eof, Lexeme: "", Literal: "", Location: lexer.TokenLocation{}}, "the first statement in code has to be an '@Environment' directive")
		return false
	} else {
		if statement.(ast.DirectiveExpr).Identifier != "Environment" {
			p.error(statement.GetToken(), "the first statement in code has to be an '@Environment' directive")
			return false
		}
	}

	return true
}
