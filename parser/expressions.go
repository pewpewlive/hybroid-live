package parser

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
	"strings"
)

func (p *Parser) expression() ast.Node {
	return /*p.cast(*/ p.fn() /*)*/
}

// func (p *Parser) cast(node ast.Node) ast.Node {
// 	if p.match(lexer.As) {
// 		if !p.PeekIsType() {
// 			p.error(p.peek(), "expected type after 'as'")
// 		}
// 		return &ast.CastExpr{
// 			Value: node,
// 			Type: p.Type(),
// 		}
// 	}

// 	return node
// }

func (p *Parser) fn() ast.Node {
	if p.match(lexer.Fn) {
		fn := &ast.FunctionExpr{
			Token: p.peek(-1),
		}
		if p.check(lexer.LeftParen) {
			fn.Params = p.parameters(lexer.LeftParen, lexer.RightParen)
		} else {
			fn.Params = make([]ast.Param, 0)
			p.error(p.peek(), "expected opening parenthesis for parameters")
		}
		fn.Return = p.returnings()

		var success bool
		fn.Body, success = p.getBody()
		if !success {
			return ast.NewImproper(fn.Token)
		}
		return fn
	} else {
		return p.multiComparison()
	}
}

func (p *Parser) multiComparison() ast.Node {
	expr := p.comparison()

	if p.isMultiComparison() {
		operator := p.peek(-1)
		right := p.comparison()
		expr = &ast.BinaryExpr{Left: expr, Operator: operator, Right: right, ValueType: ast.Bool}
	}

	return expr
}

func (p *Parser) comparison() ast.Node {
	expr := p.term()

	if p.isComparison() {
		operator := p.peek(-1)
		right := p.term()
		expr = &ast.BinaryExpr{Left: expr, Operator: operator, Right: right, ValueType: ast.Bool}
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
		expr = &ast.BinaryExpr{Left: expr, Operator: operator, Right: right, ValueType: p.determineValueType(expr, right)}
	}

	return expr
}

func (p *Parser) factor() ast.Node {
	expr := p.concat()

	if p.match(lexer.Star, lexer.Slash, lexer.Caret, lexer.Modulo) {
		operator := p.peek(-1)
		right := p.factor()

		expr = &ast.BinaryExpr{Left: expr, Operator: operator, Right: right, ValueType: p.determineValueType(expr, right)}
	}

	return expr
}

func (p *Parser) concat() ast.Node {
	expr := p.unary()

	if p.match(lexer.Concat) {
		operator := p.peek(-1)
		right := p.concat()
		return &ast.BinaryExpr{Left: expr, Operator: operator, Right: right, ValueType: p.determineValueType(expr, right)}
	}

	return expr
}

func (p *Parser) unary() ast.Node {
	if p.match(lexer.Bang, lexer.Minus) {
		operator := p.peek(-1)
		right := p.unary()
		return &ast.UnaryExpr{Operator: operator, Value: right}
	}
	return p.accessorExprDepth2(nil)
}

func (p *Parser) call(caller ast.Node) ast.Node {
	hasGenerics := false
	args := []*ast.TypeExpr{}
	if p.check(lexer.Less) {
		args = p.genericArguments()
		hasGenerics = false
	}
	if !p.check(lexer.LeftParen) {
		if hasGenerics {
			p.error(p.peek(), "expected call arguments after generic arguments")
		}
		return caller
	}

	callerType := caller.GetType()
	if callerType != ast.Identifier && callerType != ast.CallExpression && callerType != ast.EnvironmentAccessExpression {
		p.error(p.peek(-1), fmt.Sprintf("cannot call unidentified value (caller: %v)", callerType))
		return &ast.Improper{Token: p.peek(-1)}
	}

	call_expr := &ast.CallExpr{
		Caller: caller,
		GenericArgs: args,
		Args:   p.arguments(),
	}

	return p.call(call_expr)
}

func (p *Parser) accessorExprDepth2(ident *ast.Node) ast.Node {
	expr, call := p.accessorExpr(ident)

	if call == nil {
		return p.call(expr)
	}

	var methodCall ast.Node
	methodCall = &ast.MethodCallExpr{
		Identifier: expr,
		Call: &ast.CallExpr{
			Caller: call,
			GenericArgs: p.genericArguments(),
			Args: p.arguments(),
		},
	}
	
	if p.check(lexer.Dot) || p.check(lexer.LeftBracket) {
		return p.accessorExprDepth2(&methodCall)
	}
	
	return p.call(methodCall)
}

func (p *Parser) accessorExpr(ident *ast.Node) (ast.Node, *ast.IdentifierExpr) {
	if ident == nil {
		expr := p.call(p.matchExpr())
		ident = &expr
	}

	isField, isMember := p.check(lexer.Dot), p.check(lexer.LeftBracket)

	if !isField && !isMember {
		return *ident, nil
	}

	p.advance()

	var expr ast.Accessor
	if isField {
		expr = &ast.FieldExpr{
			Identifier: *ident,
		}
	} else {
		expr = &ast.MemberExpr{
			Identifier: *ident,
		}
	}

	var propIdentifier ast.Node
	if isField {
		propIdentifier = &ast.IdentifierExpr{Name: p.advance()}
		if propIdentifier.GetToken().Type != lexer.Identifier {
			p.error(propIdentifier.GetToken(), "expected identifier in field expression")
		}
		if p.check(lexer.LeftParen) || p.check(lexer.Less) {
			return *ident, propIdentifier.(*ast.IdentifierExpr)
		}
	}else if isMember {
		propIdentifier = p.expression()
		p.consume("expected closing bracket in member expression", lexer.RightBracket)
	}
	
	expr.SetPropertyIdentifier(propIdentifier)
	if p.check(lexer.Dot) || p.check(lexer.LeftBracket) {
		prop, call := p.accessorExpr(&propIdentifier)

		expr.SetProperty(prop)
	
		return expr, call
	}
	expr.SetProperty(propIdentifier)
	
	return expr, nil
}

func (p *Parser) matchExpr() ast.Node {
	if p.match(lexer.Match) {
		return &ast.MatchExpr{MatchStmt: *p.matchStmt(true)}
	}

	return p.macroCall()
}

func (p *Parser) macroCall() ast.Node {
	if p.match(lexer.At) {
		macroCall := &ast.MacroCallExpr{}
		caller := p.primary(true)
		callerType := caller.GetType()
		if callerType != ast.CallExpression {
			p.error(caller.GetToken(), "expected call after '@'")
			return &ast.Improper{}
		}
		macroCall.Caller = caller.(*ast.CallExpr)
		return macroCall
	}

	return p.new()
}

func (p *Parser) new() ast.Node {
	if p.match(lexer.New) {
		expr := ast.NewExpr{
			Token: p.peek(-1),
		}

		expr.Type = p.Type()
		expr.Args = p.arguments()

		return &expr
	}

	return p.spawn()
}

func (p *Parser) spawn() ast.Node {
	if p.match(lexer.Spawn) {
		expr := ast.SpawnExpr{
			Token: p.peek(-1),
		}

		expr.Type = p.Type()
		expr.Args = p.arguments()

		return &expr
	}

	return p.self()
}

func (p *Parser) self() ast.Node {
	if p.match(lexer.Self) {
		return &ast.SelfExpr{
			Token: p.peek(-1),
		}
	}

	return p.primary(true)
}

func (p *Parser) primary(allowStruct bool) ast.Node {
	if p.match(lexer.False) {
		return &ast.LiteralExpr{Value: "false", ValueType: ast.Bool, Token: p.peek(-1)}
	}
	if p.match(lexer.True) {
		return &ast.LiteralExpr{Value: "true", ValueType: ast.Bool, Token: p.peek(-1)}
	}

	if p.match(lexer.Number, lexer.Fixed, lexer.FixedPoint, lexer.Degree, lexer.Radian, lexer.String) {
		literal := p.peek(-1)
		var valueType ast.PrimitiveValueType
		env, ok := p.program[0].(*ast.EnvironmentStmt)
		if ok {
			envType := env.EnvType.Type
			allowFX := envType == ast.Level
			switch literal.Type {
			case lexer.Number:
				if allowFX && strings.ContainsRune(literal.Lexeme, '.') {
					p.error(literal, "cannot have a float in a level or shared environment")
				}
				valueType = ast.Number
			case lexer.Fixed:
				if !allowFX {
					p.error(literal, "cannot have a fixed in a mesh or sound environment")
				}
				valueType = ast.Fixed
			case lexer.FixedPoint:
				if !allowFX {
					p.error(literal, "cannot have a fixedpoint in a mesh, sound environment")
				}
				valueType = ast.FixedPoint
			case lexer.Degree:
				if !allowFX {
					p.error(literal, "cannot have a degree, sound environment")
				}
				valueType = ast.Degree
			case lexer.Radian:
				if !allowFX {
					p.error(literal, "cannot have a radian in a mesh or sound environment")
				}
				valueType = ast.Radian
			case lexer.String:
				valueType = ast.String
			}
		}

		return &ast.LiteralExpr{Value: literal.Literal, ValueType: valueType, Token: literal}
	}

	if p.match(lexer.LeftBrace) {
		return p.parseMap()
	}

	if p.match(lexer.LeftBracket) {
		return p.list()
	}

	if allowStruct && p.match(lexer.Struct) {
		return p.structExpr()
	}

	if p.match(lexer.Identifier) {
		token := p.peek(-1)
		if !p.match(lexer.Colon) {
			return &ast.IdentifierExpr{Name: token, ValueType: ast.Ident}
		}

		envPath := &ast.EnvPathExpr{
			Path: token,
		}

		next := p.advance()
		for p.match(lexer.Colon) {
			envPath.Combine(next)
			if next.Type != lexer.Identifier {
				p.error(next, "expected identifier in environment expression")
				return &ast.Improper{Token: next}
			}
			next = p.advance()
		}
		envExpr := &ast.EnvAccessExpr{
			PathExpr: envPath,
		}
		envExpr.Accessed = &ast.IdentifierExpr{
			Name: next,
			ValueType: ast.Invalid,
		}

		return envExpr
	}

	if p.match(lexer.LeftParen) {
		token := p.peek(-1)
		expr := p.expression()
		if expr.GetType() == ast.NA {
			p.error(p.peek(), "expected expression")
		}
		p.consume("expected ')' after expression", lexer.RightParen)
		return &ast.GroupExpr{Expr: expr, Token: token, ValueType: expr.GetValueType()}
	}

	if p.match(lexer.Self) {
		return &ast.IdentifierExpr{Name: p.peek(-1)}
	}

	return &ast.Improper{Token: p.peek()}
}

func (p *Parser) list() ast.Node {
	token := p.peek(-1)
	list := make([]ast.Node, 0)
	if p.match(lexer.RightBracket) {
		return &ast.ListExpr{ValueType: ast.List, List: list, Token: token}
	}
	exprInList := p.expression()
	if exprInList.GetType() == ast.NA {
		p.error(p.peek(), "expected expression")
	}
	list = append(list, exprInList)
	for !p.match(lexer.RightBracket) {
		p.consume("expected ',' after expression", lexer.Comma)

		exprInList := p.expression()
		if exprInList.GetType() == ast.NA {
			p.error(p.peek(), "expected expression")
			p.advance()
		}
		list = append(list, exprInList)
	}

	return &ast.ListExpr{ValueType: ast.List, List: list, Token: token}
}

func (p *Parser) parseMap() ast.Node {
	token := p.peek(-1)
	parsedMap := make(map[lexer.Token]ast.Property, 0)
	for !p.check(lexer.RightBrace) {
		key := p.primary(true)

		var newKey lexer.Token
		switch key := key.(type) {
		case *ast.IdentifierExpr:
			newKey = key.GetToken()
		// case *ast.LiteralExpr:
		// 	if key.GetValueType() != ast.String {
		// 		p.error(key.GetToken(), "expected a string in map initialization")
		// 	}
		// 	newKey = key.GetToken()
		default:
			p.error(key.GetToken(), "expected either string or an identifier in map initialization")
			p.advance()
			return &ast.Improper{Token: p.peek(-1)}
		}

		if _, ok := p.consume("expected '=' after map key", lexer.Equal); !ok {
			return &ast.Improper{Token: p.peek(-1)}
		}

		expr := p.expression()
		if expr.GetType() == ast.NA {
			p.error(p.peek(), "expected expression")
		}

		if p.peek().Type == lexer.RightBrace {
			parsedMap[newKey] = ast.Property{Expr: expr, Type: expr.GetValueType()}
			break
		}

		if _, ok := p.consume("expected ',' after expression", lexer.Comma); !ok {
			return &ast.Improper{Token: p.peek(-1)}
		}

		parsedMap[newKey] = ast.Property{Expr: expr, Type: expr.GetValueType()}
	}
	p.advance()

	return &ast.MapExpr{Map: parsedMap, Token: token}
}

func (p *Parser) structExpr() ast.Node {
	_struct := ast.StructExpr{
		Token:  p.peek(-1),
		Fields: make([]*ast.FieldDeclarationStmt, 0),
	}

	_, ok := p.consume("expected opening brace in anonymous struct expression", lexer.LeftBrace)
	if !ok {
		return &ast.Improper{Token: _struct.Token}
	}

	for !p.match(lexer.RightBrace) {
		field := p.fieldDeclarationStmt()
		if field.GetType() != ast.NA {
			_struct.Fields = append(_struct.Fields, field.(*ast.FieldDeclarationStmt))
		} else {
			p.error(field.GetToken(), "expected field declaration inside anonymous struct")
		}
		for p.match(lexer.Comma) {
			field := p.fieldDeclarationStmt()
			if field.GetType() != ast.NA {
				_struct.Fields = append(_struct.Fields, field.(*ast.FieldDeclarationStmt))
			} else {
				p.error(field.GetToken(), "expected field declaration inside anonymous struct")
			}
		}
	}

	return &_struct
}

func (p *Parser) WrappedType() *ast.TypeExpr {
	typee := ast.TypeExpr{}
	if p.check(lexer.Greater) {
		p.error(p.peek(), "empty wrapped type")
		return &typee
	}
	expr2 := p.Type()
	return expr2
}

func (p *Parser) Type() *ast.TypeExpr {
	var expr ast.Node
	token := p.advance()
	if p.match(lexer.Colon) {
		if token.Type != lexer.Identifier {
			p.error(token, "expected identifier")
		}
		envAccess := &ast.EnvAccessExpr{
			PathExpr: &ast.EnvPathExpr{
				Path: token,
			},
		}
		next := p.advance()
		for p.match(lexer.Colon) {
			envAccess.PathExpr.Combine(next)
			if next.Type != lexer.Identifier {
				p.error(next, "expected identifier in environment expression")
			}
			next = p.advance()
		}
		envAccess.Accessed = &ast.IdentifierExpr{
			Name: next,
		}
		expr = envAccess
	}else {
		expr = &ast.IdentifierExpr{
			Name: token,
			ValueType: ast.Invalid,
		}
	}

	var typ *ast.TypeExpr
	if expr.GetType() == ast.EnvironmentAccessExpression {
		typ = &ast.TypeExpr{Name: expr}
		typ.IsVariadic = p.match(lexer.DotDotDot)
		return typ
	}
	exprToken := expr.GetToken()

	switch exprToken.Type {
	// case lexer.DotDotDot:
	// 	p.advance()
	// 	typ := p.Type()
	// 	typ.IsVariadic = true
	// 	return typ
	case lexer.Identifier:
		typ = &ast.TypeExpr{}
		if p.match(lexer.Less) { // map<number>
			typ.WrappedType = p.WrappedType()
			p.consume("expected '>'", lexer.Greater)
		}
		typ.Name = expr
	case lexer.Fn:
		typ = &ast.TypeExpr{}
		p.advance()
		params := make([]*ast.TypeExpr, 0) // yes
		typ.Returns = make([]*ast.TypeExpr, 0)
		if p.match(lexer.LeftParen) { // because this fucks up
			if !p.match(lexer.RightParen) {
				params = append(params, p.Type())

				for p.match(lexer.Comma) {
					params = append(params, p.Type())
				}
				p.consume("expected closing parenthesis in 'fn(...'", lexer.RightParen)
			}
		}

		typ.Params = params
		typ.Returns = p.returnings()
		typ.Name = expr
	case lexer.Struct:
		p.advance()
		fields := p.parameters(lexer.LeftBrace, lexer.RightBrace)
		typ = &ast.TypeExpr{Name: expr, Fields: fields}
	case lexer.Entity:
		typ = &ast.TypeExpr{Name: &ast.IdentifierExpr{Name: exprToken}}
	default:
		//p.error(exprToken, "Improper type")
		p.advance()
		typ = &ast.TypeExpr{Name: expr}
	}
	typ.IsVariadic = p.match(lexer.DotDotDot)

	return typ
}

func StringToEnvType(name string) ast.EnvType {
	switch name {
	case "Mesh":
		return ast.MeshEnv
	case "Level":
		return ast.Level
	case "Sound":
		return ast.SoundEnv
	default:
		return ast.InvalidEnv
	}
}

func (p *Parser) EnvType() *ast.EnvTypeExpr {
	name, ok := p.consume("expected identifier for a environment type expr", lexer.Identifier)

	if !ok {
		return &ast.EnvTypeExpr{Type: ast.InvalidEnv, Token: name}
	}

	envType := StringToEnvType(name.Lexeme)

	if envType == ast.InvalidEnv {
		p.error(name, "expected 'Level', 'Mesh' or 'Sound' as environment type")
	}

	return &ast.EnvTypeExpr{Type: envType, Token: name}
}

func (p *Parser) EnvPathExpr() ast.Node {
	ident, ok := p.consume("expected identifier for an environment path", lexer.Identifier)

	if !ok {
		return &ast.Improper{Token: ident}
	}

	envPath := &ast.EnvPathExpr{
		Path: ident,
	}

	for p.match(lexer.Colon) {
		ident, ok = p.consume("expected identifier in environment path", lexer.Identifier)
		if !ok {
			return &ast.Improper{Token: envPath.GetToken()}
		}
		envPath.Combine(ident)
	}

	return envPath
}
