package parser

import (
	"hybroid/ast"
	"hybroid/lexer"
)

func (p *Parser) matchDirectiveStmt(expr ast.Node) bool {
	directive, ok := expr.(ast.DirectiveStmt)
	if !ok {
		return false
	}

	ident, valueType := directive.Expr.(ast.IdentifierExpr).Name, directive.Expr.GetValueType()

	switch directive.Identifier {
	case "Environment":
		if ident != "Level" && ident != "Mesh" && ident != "Sound" && ident != "Shared" && ident != "LuaGeneric" {
			p.error(directive.Expr.GetToken(), "invalid expression in '@Environment' directive")
			return false
		}
	case "Len":
		if valueType != ast.String && valueType != ast.Map && valueType != ast.List {
			p.error(directive.Expr.GetToken(), "invalid expression in '@Len' directive")
			return false
		}
	case "MapToStr":
		if valueType != ast.Map {
			p.error(directive.Expr.GetToken(), "invalid expression in '@MapToStr' directive")
			return false
		}
	case "ListToStr":
		if valueType != ast.List {
			p.error(directive.Expr.GetToken(), "invalid expression in '@ListToStr' directive")
			return false
		}
	default:
		// TODO: Add support for custom directives

		p.error(directive.Expr.GetToken(), "invalid directive call")
		return false
	}
	return true
}

func (p *Parser) directiveCall() ast.Node {
	directiveNode := ast.DirectiveStmt{}

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
	if _, ok := statement.(ast.DirectiveStmt); !ok {
		p.error(lexer.Token{Type: lexer.Eof, Lexeme: "", Literal: "", Location: lexer.TokenLocation{}}, "the first statement in code has to be an '@Environment' directive")
		return false
	} else {
		if statement.(ast.DirectiveStmt).Identifier != "Environment" {
			p.error(statement.GetToken(), "the first statement in code has to be an '@Environment' directive")
			return false
		}
	}

	return true
}
