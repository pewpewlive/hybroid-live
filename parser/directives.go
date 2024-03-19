package parser

import "hybroid/lexer"

func (p *Parser) matchDirectiveStmt(expr *Node) bool {
	switch expr.Identifier {
	case "Environment":
		ident := expr.Expression.Identifier
		if ident != "Level" && ident != "Mesh" && ident != "Sound" && ident != "Shared" && ident != "LuaGeneric" {
			p.error(expr.Expression.Token, "invalid expression in '@Environment' directive")
			return false
		}
	case "Len":
		if expr.Expression.ValueType != String && expr.Expression.ValueType != Map && expr.Expression.ValueType != List {
			p.error(expr.Expression.Token, "invalid expression in '@Len' directive")
			return false
		}
	case "MapToStr":
		if expr.Expression.ValueType != Map {
			p.error(expr.Expression.Token, "invalid expression in '@MapToStr' directive")
			return false
		}
	case "ListToStr":
		if expr.Expression.ValueType != List {
			p.error(expr.Expression.Token, "invalid expression in '@ListToStr' directive")
			return false
		}
	default:
		// TODO: Add support for custom directives

		p.error(expr.Token, "invalid directive call")
		return false
	}
	return true
}

func (p *Parser) directiveCall() *Node {
	directiveNode := Node{
		NodeType: DirectiveStmt,
		// Identifier: ident.Lexeme,
		// Expression: expr,
		// Token:      ident,
	}

	ident, identOk := p.consume(lexer.Identifier, "expected identifier in directive call")
	if !identOk {
		return &directiveNode
	}
	directiveNode.Identifier = ident.Lexeme
	directiveNode.Token = ident

	if _, ok := p.consume(lexer.LeftParen, "expected '(' after directive call"); !ok {
		return &directiveNode
	}

	directiveNode.Expression = p.expression()

	if _, ok := p.consume(lexer.RightParen, "expected ')' after directive call"); !ok {
		return &directiveNode
	}

	p.matchDirectiveStmt(&directiveNode)

	return &directiveNode
}

func (p *Parser) verifyEnvironmentDirective(statement *Node) bool {
	if statement == nil {
		p.error(lexer.Token{Type: lexer.Eof, Lexeme: "", Literal: "", Location: lexer.TokenLocation{}}, "the first statement in code has to be an '@Environment' directive")
		return false
	} else {
		if statement.NodeType != DirectiveStmt || statement.Identifier != "Environment" {
			p.error(statement.Token, "the first statement in code has to be an '@Environment' directive")
			return false
		}
	}

	return true
}
