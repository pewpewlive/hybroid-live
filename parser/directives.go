package parser

import (
	"hybroid/ast"
	"hybroid/lexer"
)

func (p *Parser) matchDirectiveStmt(expr ast.Node) bool {
	directive, ok := expr.(ast.DirectiveExpr)
	if !ok {
		return false
	}

	ident := directive.Expr.(ast.IdentifierExpr).Name

	switch directive.Identifier {
	case "Environment":
		directive.ValueType = ast.Undefined
		if ident != "Level" && ident != "Mesh" && ident != "Sound" && ident != "Shared" && ident != "LuaGeneric" {
			p.error(directive.Expr.GetToken(), "invalid expression in '@Environment' directive")
			return false
		}
	case "Len":
		directive.ValueType = ast.Number
	case "MapToStr":
		directive.ValueType = ast.String
	case "ListToStr":
		directive.ValueType = ast.String
	default:
		// TODO: Add support for custom directives

		directive.ValueType = ast.Undefined
		p.error(directive.Expr.GetToken(), "invalid directive call")
		return false
	}
	return true
}

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

	p.matchDirectiveStmt(directiveNode)

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
