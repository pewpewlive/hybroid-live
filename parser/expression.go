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
		fn := ast.AnonFnExpr{
			Token: p.peek(-1),
		}
		fn.Params = p.parameters()

		ret := make([]ast.TypeExpr, 0)
		for p.check(lexer.Identifier) {
			ret = append(ret, p.Type())
			if !p.check(lexer.Comma) {
				break
			} else {
				p.advance()
			}// yes
		}// did you pull??
		
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

	return p.methodCall(nil)
}

func (p *Parser) methodCall(owner ast.Node) ast.Node {
	expr, call := p.fieldExpr(owner)

	if call != nil {
		return call
	}

	return expr
}

func (p *Parser) fieldExpr(owner ast.Node) (ast.Node, ast.Node) { 
	if owner == nil {
		expression := p.memberCall(nil) 
		expr := ast.FieldExpr{
			Identifier: expression, 
		}

		if p.match(lexer.Dot) {
			expr2, call := p.fieldExpr(expr.Identifier)
			expr.Property = expr2
			return expr, call
		} else {
			return expression, nil
		}
	}

	ident := p.memberCall(nil)

	if memberExpr, ok := ident.(ast.MemberExpr); ok {
		memberExpr.Owner = owner;
		ident = memberExpr
	}
	
	expr := ast.FieldExpr{
		Owner:      owner,
		Identifier: ident,
	}

	if ident.GetType() == ast.CallExpression && ident.(ast.CallExpr).Caller.GetType() == ast.Identifier { 
		call := ident.(ast.CallExpr)
		fieldExpr := ast.FieldExpr{
			Owner:owner,
			Identifier: call.Caller,
		}
		return owner, ast.MethodCallExpr{
			Owner: owner,
			Call: fieldExpr,
			Args:call.Args,
			Token:call.Token,
		}
	}

	if p.match(lexer.Dot) {
		expr2 := p.memberCall(nil)
		if memberExpr, ok := expr2.(ast.MemberExpr); ok {
			memberExpr.Owner = expr;
			expr2 = memberExpr
		}
		expr.Property = expr2
	}

	return expr, nil
}

func (p *Parser) memberCall(owner ast.Node) ast.Node {
	expr := p.member(owner)

	if p.check(lexer.LeftParen) {
		return p.call(expr)
	}

	return expr
}

func (p *Parser) member(owner ast.Node) ast.Node { 
	if owner == nil {
		expression := p.new()
		expr := ast.MemberExpr{
			Identifier: expression,
		}

		if p.match(lexer.LeftBracket) {
			expr2 := p.member(expr.Identifier)
			expr.Property = expr2
			return expr
		} else {
			return expression
		}
	}
	var expr ast.MemberExpr 

	prop := p.expression()
	p.consume("expected closing bracket", lexer.RightBracket)

	expr = ast.MemberExpr{
		Owner:      owner,
		Identifier: prop,
	}

	if p.match(lexer.LeftBracket) {
		expr2 := p.member(expr)
		expr.Property = expr2
	}

	return expr
}

func (p *Parser) new() ast.Node {
	if p.match(lexer.Neww) {
		expr := ast.NewExpr{
			Token: p.peek(-1),
		}

		typee, ok := p.consume("expected type after new keyword", lexer.Identifier)

		if ok {
			expr.Type = typee
		}

		expr.Args = p.arguments()

		return expr
	}

	return p.self()
}

func (p *Parser) self() ast.Node {
	if p.match(lexer.Self) {
		return ast.SelfExpr{
			Token: p.peek(-1),
		}
	}

	return p.primary()
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
