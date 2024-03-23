package parser

import (
	"hybroid/ast"
	"hybroid/lexer"
	"strings"
)

func (p *Parser) expression() ast.Node {
	return p.list()
}

func (p *Parser) list() ast.Node {
	if !p.match(lexer.LeftBracket) {
		return p.parseMap()
	}

	token := p.peek(-1)
	list := make([]ast.Node, 0)
	for !p.check(lexer.RightBracket) {
		exprInList := p.expression()
		if exprInList.GetType() == 0 {
			p.error(p.peek(), "expected expression")
		}

		token, _ := p.consume("expected ',' or ']' after expression", lexer.Comma, lexer.RightBracket)

		list = append(list, exprInList)
		if token.Type == lexer.RightBracket || token.Type == lexer.Eof {
			break
		}
	}

	return ast.ListExpr{ValueType: ast.List, List: list, Token: token}
}

func (p *Parser) parseMap() ast.Node {
	if !p.match(lexer.LeftBrace) {
		return p.directive()
	}

	token := p.peek(-1)
	parsedMap := make(map[string]ast.Property, 0)
	for !p.check(lexer.RightBrace) {
		key := p.primary()

		var newKey string
		switch key := key.(type) {
		case ast.IdentifierExpr:
			newKey = key.Name
		case ast.LiteralExpr:
			if key.GetValueType() != ast.String {
				p.error(key.GetToken(), "expected a string in map initialization")
			}
			newKey = key.GetToken().Literal
		default:
			p.error(key.GetToken(), "expected either string or an identifier in map initialization")
			return ast.Unknown{Token: p.peek(-1)}
		}

		if _, ok := p.consume("expected ':' after map key", lexer.Colon); !ok {
			return ast.Unknown{Token: p.peek(-1)}
		}

		expr := p.expression()
		if expr.GetType() == 0 {
			p.error(p.peek(), "expected expression")
		}

		if p.peek().Type == lexer.RightBrace {
			parsedMap[newKey] = ast.Property{Expr: expr, Type: expr.GetValueType()}
			break
		}

		if _, ok := p.consume("expected ',' or '}' after expression", lexer.Comma); !ok {
			return ast.Unknown{Token: p.peek(-1)}
		}

		parsedMap[newKey] = ast.Property{Expr: expr, Type: expr.GetValueType()}
	}
	p.advance()

	return ast.MapExpr{ValueType: ast.Map, Map: parsedMap, Token: token}
}

func (p *Parser) directive() ast.Node {
	if !p.match(lexer.At) {
		return p.multiComparison()
	}

	return p.directiveCall()
}

func (p *Parser) multiComparison() ast.Node {
	expr := p.equality()

	if p.match(lexer.And, lexer.Or) {
		operator := p.peek(-1)
		right := p.equality()
		expr = ast.BinaryExpr{Left: expr, Operator: operator, Right: right, ValueType: ast.Bool}
	}

	return expr
}

func (p *Parser) equality() ast.Node {
	expr := p.comparison()

	if p.match(lexer.BangEqual, lexer.EqualEqual) {
		operator := p.peek(-1)
		right := p.comparison()
		expr = ast.BinaryExpr{Left: expr, Operator: operator, Right: right, ValueType: ast.Bool}
	}

	return expr
}

func (p *Parser) comparison() ast.Node {
	expr := p.term()

	if p.match(lexer.Greater, lexer.GreaterEqual, lexer.Less, lexer.LessEqual) {
		operator := p.peek(-1)
		right := p.term()
		expr = ast.BinaryExpr{Left: expr, Operator: operator, Right: right, ValueType: ast.Bool}
	}

	return expr
}

func (p *Parser) determineValueType(left ast.Node, right ast.Node) ast.PrimitiveValueType {
	if left.GetValueType() == right.GetValueType() {
		return left.GetValueType()
	}
	if IsFx(left.GetValueType()) && IsFx(right.GetValueType()) {
		return ast.FixedPoint
	}

	return ast.Undefined
}

func (p *Parser) term() ast.Node {
	expr := p.factor()

	if p.match(lexer.Plus, lexer.Minus) {
		operator := p.peek(-1)
		right := p.term()
		expr = ast.BinaryExpr{Left: expr, Operator: operator, Right: right, ValueType: p.determineValueType(expr, right)}
	}

	return expr
}

func (p *Parser) factor() ast.Node {
	expr := p.unary()

	if p.match(lexer.Star, lexer.Slash, lexer.Caret, lexer.Modulo) {
		operator := p.peek(-1)
		right := p.factor()

		expr = ast.BinaryExpr{Left: expr, Operator: operator, Right: right, ValueType: p.determineValueType(expr, right)}
	}

	return expr
}

func (p *Parser) unary() ast.Node {
	if p.match(lexer.Bang, lexer.Minus) {
		operator := p.peek(-1)
		right := p.unary()
		return ast.UnaryExpr{Operator: operator, Value: right}
	}

	return p.memberCall()
}

func (p *Parser) memberCall() ast.Node {
	expr := p.member()

	if p.check(lexer.LeftParen) {
		return p.call(expr)
	}

	return expr
}

func (p *Parser) call(caller ast.Node) ast.Node {
	callerType := caller.GetType()
	if callerType != ast.Identifier && callerType != ast.MemberExpression && callerType != ast.CallExpression {
		p.error(p.peek(-1), "cannot call unidentified value")
		return ast.Unknown{Token: p.peek(-1)}
	}

	call_expr := ast.CallExpr{
		Identifier: caller.GetToken().Lexeme,
		Caller:     caller,
		Args:       p.arguments(),
		Token:      caller.GetToken(),
	}

	if p.check(lexer.LeftParen) {
		expr := p.call(call_expr)
		if expr.GetType() == ast.CallExpression {
			call_expr = expr.(ast.CallExpr)
		}
	}

	return call_expr
}

func (p *Parser) arguments() []ast.Node {
	if _, ok := p.consume("expected opening paren after an identifier", lexer.LeftParen); !ok {
		return nil
	}

	var args []ast.Node
	if p.match(lexer.RightParen) {
		args = make([]ast.Node, 0)
	} else {
		args = append(args, p.expression())
		for p.match(lexer.Comma) {
			args = append(args, p.expression())
		}
		p.consume("expected closing paren after arguments", lexer.RightParen)
	}

	return args
}

func (p *Parser) member() ast.Node {
	expr := p.primary()

	for p.match(lexer.Dot, lexer.LeftBracket) {
		operator := p.peek(-1)

		prop := p.primary()
		if operator.Type == lexer.Dot && prop.GetType() != ast.Identifier {
			p.error(p.peek(-1), "expected identifier after '.'")
		}
		if operator.Type == lexer.LeftBracket {
			if prop.GetType() == ast.LiteralExpression && prop.(ast.LiteralExpr).ValueType != ast.String {
				p.error(prop.(ast.LiteralExpr).Token, "expected string after '['")
			} else if prop.GetType() != ast.LiteralExpression {
				p.error(prop.(ast.LiteralExpr).Token, "expected string after '['")
			}
			p.consume("expected closing bracket", lexer.RightBracket)
		}

		expr = ast.MemberExpr{
			Identifier: expr,
			Property:   prop,
			Token:      expr.GetToken(),
		}
	}

	return expr
}

func (p *Parser) primary() ast.Node {
	if p.match(lexer.False) {
		return ast.LiteralExpr{Value: "false", ValueType: ast.Bool, Token: p.peek(-1)}
	}
	if p.match(lexer.True) {
		return ast.LiteralExpr{Value: "true", ValueType: ast.Bool, Token: p.peek(-1)}
	}
	if p.match(lexer.Nil) {
		return ast.LiteralExpr{Value: "nil", ValueType: ast.Nil, Token: p.peek(-1)}
	}

	if p.match(lexer.Number, lexer.Fixed, lexer.FixedPoint, lexer.Degree, lexer.Radian, lexer.String) {
		literal := p.peek(-1)
		var valueType ast.PrimitiveValueType
		ident := p.program[0].(ast.DirectiveExpr).Expr.(ast.IdentifierExpr)
		allowFX := ident.Name == "Level" || ident.Name == "Shared"
		switch literal.Type {
		case lexer.Number:
			if allowFX && strings.ContainsRune(literal.Lexeme, '.') {
				p.error(literal, "cannot have a float in a level or shared environment")
			}
			valueType = ast.Number
		case lexer.Fixed:
			if !allowFX {
				p.error(literal, "cannot have a fixed in a mesh, sound or luageneric environment")
			}
			valueType = ast.Fixed
		case lexer.FixedPoint:
			if !allowFX {
				p.error(literal, "cannot have a fixedpoint in a mesh, sound or luageneric environment")
			}
			valueType = ast.FixedPoint
		case lexer.Degree:
			if !allowFX {
				p.error(literal, "cannot have a degree, sound or luageneric environment")
			}
			valueType = ast.Degree
		case lexer.Radian:
			if !allowFX {
				p.error(literal, "cannot have a radian in a mesh, sound or luageneric environment")
			}
			valueType = ast.Radian
		case lexer.String:
			valueType = ast.String
		}
		return ast.LiteralExpr{Value: literal.Literal, ValueType: valueType, Token: literal}
	}

	if p.match(lexer.Identifier) {
		token := p.peek(-1)
		return ast.IdentifierExpr{Name: token.Lexeme, Token: token, ValueType: ast.Ident}
	}

	if p.match(lexer.LeftParen) {
		token := p.peek(-1)
		expr := p.expression()
		if expr.GetType() == 0 {
			p.error(p.peek(), "expected expression")
		}
		p.consume("expected ')' after expression", lexer.RightParen)
		return ast.GroupExpr{Expr: expr, Token: token, ValueType: expr.GetValueType()}
	}
	p.advance()
	return ast.Unknown{}
}
