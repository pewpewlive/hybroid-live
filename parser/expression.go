package parser

import (
	"hybroid/ast"
	"hybroid/lexer"
	"strings"
)

func (p *Parser) list() ast.Node {
	token := p.peek(-1)
	list := make([]ast.Node, 0)
	for !p.match(lexer.RightBracket) {
		exprInList := p.expression()
		if exprInList.GetType() == ast.NA {
			p.error(p.peek(), "expected expression")
			break
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
	token := p.peek(-1)
	parsedMap := make(map[lexer.Token]ast.Property, 0)
	for !p.check(lexer.RightBrace) {
		key := p.primary()

		var newKey lexer.Token
		switch key := key.(type) {
		case ast.IdentifierExpr:
			newKey = key.GetToken()
		case ast.LiteralExpr:
			if key.GetValueType() != ast.String {
				p.error(key.GetToken(), "expected a string in map initialization")
			}
			newKey = key.GetToken()
		default:
			p.error(key.GetToken(), "expected either string or an identifier in map initialization")
			p.advance()
			return ast.Improper{Token: p.peek(-1)}
		}

		if _, ok := p.consume("expected ':' after map key", lexer.Colon); !ok {
			return ast.Improper{Token: p.peek(-1)}
		}

		expr := p.expression()
		if expr.GetType() == 0 {
			p.error(p.peek(), "expected expression")
		}

		if p.peek().Type == lexer.RightBrace {
			parsedMap[newKey] = ast.Property{Expr: expr, Type: expr.GetValueType()}
			break
		}

		if _, ok := p.consume("expected ',' or '}' after expression", lexer.Comma, lexer.RightBrace); !ok {
			return ast.Improper{Token: p.peek(-1)}
		}

		parsedMap[newKey] = ast.Property{Expr: expr, Type: expr.GetValueType()}
	}
	p.advance()

	return ast.MapExpr{Map: parsedMap, Token: token}
}

func (p *Parser) expression() ast.Node {
	return p.fn()
}

func (p *Parser) fn() ast.Node {
	if p.match(lexer.Fn) {
		fn := ast.AnonFnExpr{}
		fn.Params = p.parameters()

		ret := make([]ast.TypeExpr, 0)
		for p.check(lexer.Identifier) {
			ret = append(ret, p.Type())
			if !p.check(lexer.Comma) {
				break
			} else {
				p.advance()
			}
		}
		fn.Return = ret
		fn.Body = *p.getBody()
		return fn
	} else {
		return p.directive()
	}
}

func (p *Parser) directive() ast.Node {
	if !p.match(lexer.At) {
		return p.multiComparison()
	}

	return p.directiveCall()
}

func (p *Parser) multiComparison() ast.Node {
	expr := p.comparison()

	if p.isMultiComparison() {
		operator := p.peek(-1)
		right := p.comparison()
		expr = ast.BinaryExpr{Left: expr, Operator: operator, Right: right, ValueType: ast.Bool}
	}

	return expr
}

func (p *Parser) comparison() ast.Node {
	expr := p.term()

	if p.isComparison() {
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

	return ast.Invalid
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

	return p.self()
}

func (p *Parser) self() ast.Node { // somestruct.x
	if p.check(lexer.Self) {
		expr := ast.SelfExpr{
			Token: p.peek(),
			Value: p.memberCall(nil),
		}
		return expr
	}

	return p.memberCall(nil)
}

func (p *Parser) memberCall(owner ast.Node) ast.Node {
	expr := p.member(owner)

	if p.check(lexer.LeftParen) {
		return p.call(expr)
	}

	return expr
}

func (p *Parser) call(caller ast.Node) ast.Node {
	callerType := caller.GetType()
	if callerType != ast.Identifier && callerType != ast.MemberExpression && callerType != ast.CallExpression {
		p.error(p.peek(-1), "cannot call unidentified value")
		return ast.Improper{Token: p.peek(-1)}
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
	if _, ok := p.consume("expected opening paren after an identifier", lexer.LeftParen); !ok {
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

func (p *Parser) member(owner ast.Node) ast.Node {
	if owner == nil {
		expression := p.primary()
		expr := ast.MemberExpr{
			Owner:      owner,
			Property:   expression,
			Identifier: expression,
			Bracketed:  false,
		}

		if p.check(lexer.Dot) || p.check(lexer.LeftBracket) {
			expr2 := p.memberCall(expr)
			expr.Property = expr2
			return expr
		} else {
			return expression
		}

	} else {
		var expr ast.MemberExpr
		bracketed := false
		operator, _ := p.consume("expected '.' or '[' after member expression", lexer.Dot, lexer.LeftBracket)

		var prop ast.Node
		if operator.Type == lexer.Dot {
			prop = p.primary()
			if prop.GetType() != ast.Identifier {
				p.error(p.peek(-1), "expected identifier after '.'")
			}
			bracketed = false
		} else if operator.Type == lexer.LeftBracket {
			prop = p.expression()
			p.consume("expected closing bracket", lexer.RightBracket)
			bracketed = true
		}

		expr = ast.MemberExpr{
			Owner:      owner,
			Property:   prop,
			Bracketed:  bracketed,
			Identifier: prop,
		}

		if p.check(lexer.Dot) || p.check(lexer.LeftBracket) {
			expr2 := p.memberCall(expr)
			expr.Property = expr2
		}

		return expr
	}
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
		allowFX := ident.Name.Lexeme == "Level" || ident.Name.Lexeme == "Shared"

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

	if p.match(lexer.LeftBrace) {
		return p.parseMap()
	}

	if p.match(lexer.LeftBracket) {
		return p.list()
	}

	if p.match(lexer.Identifier) {
		token := p.peek(-1)
		return ast.IdentifierExpr{Name: token, ValueType: ast.Ident}
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

	if p.match(lexer.Self) {
		return ast.IdentifierExpr{Name: p.peek(-1)}
	}

	return ast.Improper{Token: p.peek()}
}

func (p *Parser) WrappedType() *ast.TypeExpr {
	typee := ast.TypeExpr{}
	if p.check(lexer.Greater) {
		p.error(p.peek(), "empty wrapped type")
		return &typee
	}
	expr2 := p.Type()
	return &expr2
}

func (p *Parser) Type() ast.TypeExpr {
	expr := p.primary()

	if expr.GetType() == ast.Identifier {
		typee := ast.TypeExpr{}

		if p.match(lexer.Less) {
			typee.WrappedType = p.WrappedType()
			p.consume("expected '>'", lexer.Greater)
		}
		typee.Name = expr.GetToken()
		return typee
	} else if expr.GetToken().Type == lexer.Fn {
		typee := ast.TypeExpr{}

		p.advance()
		typee.Params = make([]ast.TypeExpr, 0)
		typee.Returns = make([]ast.TypeExpr, 0)
		if p.match(lexer.LeftParen) {
			typee.Params = append(typee.Params, p.Type())

			for p.match(lexer.Comma) {
				typee.Params = append(typee.Params, p.Type())
			}
			p.consume("expected closing parenthesis in 'fn(...'", lexer.RightParen)
		}

		if p.check(lexer.Identifier) {
			typee.Returns = append(typee.Returns, p.Type())

			for p.match(lexer.Comma) {
				typee.Returns = append(typee.Returns, p.Type())
			}
		}

		typee.Name = expr.GetToken()
		return typee
	} else {
		p.error(expr.GetToken(), "Expected an identifier for a type")
		p.advance()
		return ast.TypeExpr{Name: expr.GetToken()}
	}

}
