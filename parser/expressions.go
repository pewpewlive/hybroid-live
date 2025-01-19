package parser

import (
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/tokens"
	"strings"
)

func (p *Parser) expression() ast.Node {
	return p.fn()
}

func (p *Parser) fn() ast.Node {
	if p.match(tokens.Fn) {
		fn := &ast.FunctionExpr{
			Token: p.peek(-1),
		}
		if p.check(tokens.LeftParen) {
			fn.Params = p.parameters(tokens.LeftParen, tokens.RightParen)
		} else {
			fn.Params = make([]ast.Param, 0)
			p.Alert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.LeftParen)
		}
		fn.Return = p.returnings()
		p.Context.FunctionReturns.Push("fn", len(fn.Return))

		var success bool
		fn.Body, success = p.getBody()
		if !success {
			return ast.NewImproper(fn.Token)
		}
		p.Context.FunctionReturns.Pop("fn")

		return fn
	} else {
		return p.multiComparison()
	}
}

func (p *Parser) multiComparison() ast.Node {
	expr := p.comparison()

	if p.isMultiComparison() {
		operator := p.peek(-1)
		right := p.multiComparison()
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

	if p.match(tokens.Plus, tokens.Minus) {
		operator := p.peek(-1)
		right := p.term()
		expr = &ast.BinaryExpr{Left: expr, Operator: operator, Right: right, ValueType: p.determineValueType(expr, right)}
	}

	return expr
}

func (p *Parser) factor() ast.Node {
	expr := p.concat()

	if p.match(tokens.Star, tokens.Slash, tokens.Caret, tokens.Modulo, tokens.BackSlash) {
		operator := p.peek(-1)
		right := p.factor()

		expr = &ast.BinaryExpr{Left: expr, Operator: operator, Right: right, ValueType: p.determineValueType(expr, right)}
	}

	return expr
}

func (p *Parser) concat() ast.Node {
	expr := p.unary()

	if p.match(tokens.Concat) {
		operator := p.peek(-1)
		right := p.concat()
		return &ast.BinaryExpr{Left: expr, Operator: operator, Right: right, ValueType: p.determineValueType(expr, right)}
	}

	return expr
}

func (p *Parser) unary() ast.Node {
	if p.match(tokens.Bang, tokens.Minus, tokens.Hash) {
		operator := p.peek(-1)
		right := p.unary()
		return &ast.UnaryExpr{Operator: operator, Value: right}
	}
	return p.entity()
}

func (p *Parser) entity() ast.Node {
	variable := p.accessorExprDepth2(nil)
	var expr ast.Node
	currentStart := p.current

	var conv *tokens.Token
	if p.match(tokens.Equal) {
		token := variable.GetToken()
		conv = &token

		expr = p.accessorExprDepth2(nil)
	} else {
		expr = variable
	}
	if p.match(tokens.Is, tokens.Isnt) {
		if conv != nil && variable.GetType() != ast.Identifier {
			p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(variable.GetToken()))
		}
		op := p.peek(-1)
		typ := p.Type()
		return &ast.EntityExpr{
			Expr:             expr,
			Type:             typ,
			ConvertedVarName: conv,
			Token:            expr.GetToken(),
			Operator:         op,
		}
	}

	p.disadvance(p.current - currentStart)

	return variable
}

func (p *Parser) call(caller ast.Node) ast.Node {
	hasGenerics := false
	args := []*ast.TypeExpr{}
	if p.check(tokens.Less) {
		var ok bool
		args, ok = p.genericArguments()
		if !ok {
			return caller
		}
		hasGenerics = false
	}
	if !p.check(tokens.LeftParen) {
		if hasGenerics {
			p.Alert(&alerts.ExpectedCallArgs{}, alerts.NewSingle(p.peek()))
		}
		return caller
	}

	callerType := caller.GetType()
	if callerType != ast.Identifier && callerType != ast.CallExpression && callerType != ast.EnvironmentAccessExpression && callerType != ast.MemberExpression {
		p.Alert(&alerts.InvalidCall{}, alerts.NewSingle(p.peek(-1)))
		return &ast.Improper{Token: p.peek(-1)}
	}

	call_expr := &ast.CallExpr{
		Caller:      caller,
		GenericArgs: args,
		Args:        p.arguments(),
	}

	return p.call(call_expr)
}

func (p *Parser) accessorExprDepth2(ident *ast.Node) ast.Node {
	expr, call := p.accessorExpr(ident)

	if call == nil {
		return p.call(expr)
	}

	args, _ := p.genericArguments()

	var methodCall ast.Node = &ast.MethodCallExpr{
		Identifier: expr,
		Call: &ast.CallExpr{
			Caller:      call,
			GenericArgs: args,
			Args:        p.arguments(),
		},
	}

	if p.check(tokens.Dot) || p.check(tokens.LeftBracket) {
		return p.accessorExprDepth2(&methodCall)
	}

	return p.call(methodCall)
}

func (p *Parser) accessorExpr(ident *ast.Node) (ast.Node, *ast.IdentifierExpr) {
	if ident == nil {
		expr := p.call(p.matchExpr())
		ident = &expr
	}

	isField, isMember := p.check(tokens.Dot), p.check(tokens.LeftBracket)

	if !isField && !isMember {
		return *ident, nil
	}

	start := p.advance()

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
		if propIdentifier.GetToken().Type != tokens.Identifier {
			p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(propIdentifier.GetToken()))
		}
		if p.check(tokens.LeftParen) || p.check(tokens.Less) {
			return *ident, propIdentifier.(*ast.IdentifierExpr)
		}
	} else if isMember {
		propIdentifier = p.expression()

		p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewMulti(start, p.peek()), tokens.RightBracket), tokens.RightBracket)
	}

	expr.SetPropertyIdentifier(propIdentifier)
	if p.check(tokens.Dot) || p.check(tokens.LeftBracket) {
		prop, call := p.accessorExpr(&propIdentifier)

		expr.SetProperty(prop)

		return expr, call
	}
	expr.SetProperty(propIdentifier)

	return expr, nil
}

func (p *Parser) matchExpr() ast.Node {
	if p.match(tokens.Match) {
		return &ast.MatchExpr{MatchStmt: *p.matchStmt(true)}
	}

	return p.macroCall()
}

func (p *Parser) macroCall() ast.Node {
	if p.match(tokens.At) {
		macroCall := &ast.MacroCallExpr{}
		caller := p.primary(true)
		callerType := caller.GetType()
		if callerType != ast.CallExpression {
			p.Alert(&alerts.ExpectedCallAfterMacroSymbol{}, alerts.NewSingle(caller.GetToken()))
			return &ast.Improper{}
		}
		macroCall.Caller = caller.(*ast.CallExpr)
		return macroCall
	}

	return p.new()
}

func (p *Parser) new() ast.Node {
	if p.match(tokens.New) {
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
	if p.match(tokens.Spawn) {
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
	if p.match(tokens.Self) {
		return &ast.SelfExpr{
			Token: p.peek(-1),
		}
	}

	return p.primary(true)
}

func (p *Parser) primary(allowStruct bool) ast.Node {
	if p.match(tokens.False) {
		return &ast.LiteralExpr{Value: "false", ValueType: ast.Bool, Token: p.peek(-1)}
	}
	if p.match(tokens.True) {
		return &ast.LiteralExpr{Value: "true", ValueType: ast.Bool, Token: p.peek(-1)}
	}

	if p.match(tokens.Number, tokens.Fixed, tokens.FixedPoint, tokens.Degree, tokens.Radian, tokens.String) {
		literal := p.peek(-1)
		var valueType ast.PrimitiveValueType
		env := p.Context.EnvStatement

		if env != nil {
			envType := env.EnvType.Type
			allowFX := envType == ast.Level
			switch literal.Type {
			case tokens.Number:
				if allowFX && strings.ContainsRune(literal.Lexeme, '.') {
					p.Alert(&alerts.ForbiddenTypeInEnvironment{}, alerts.NewSingle(literal), "float", []string{"level", "shared"})
				}
				valueType = ast.Number
			case tokens.Fixed:
				if !allowFX {
					p.Alert(&alerts.ForbiddenTypeInEnvironment{}, alerts.NewSingle(literal), "fixed", []string{"mesh", "sound"})
				}
				valueType = ast.Fixed
			case tokens.FixedPoint:
				if !allowFX {
					p.Alert(&alerts.ForbiddenTypeInEnvironment{}, alerts.NewSingle(literal), "fixedpoint", []string{"mesh", "sound"})
				}
				valueType = ast.FixedPoint
			case tokens.Degree:
				if !allowFX {
					p.Alert(&alerts.ForbiddenTypeInEnvironment{}, alerts.NewSingle(literal), "fixed", []string{"mesh", "sound"})
				}
				valueType = ast.Degree
			case tokens.Radian:
				if !allowFX {
					p.Alert(&alerts.ForbiddenTypeInEnvironment{}, alerts.NewSingle(literal), "radian", []string{"mesh", "sound"})
				}
				valueType = ast.Radian
			case tokens.String:
				valueType = ast.String
			}
		}

		return &ast.LiteralExpr{Value: literal.Literal, ValueType: valueType, Token: literal}
	}

	if p.match(tokens.LeftBrace) {
		return p.parseMap()
	}

	if p.match(tokens.LeftBracket) {
		return p.list()
	}

	if allowStruct && p.match(tokens.Struct) {
		return p.structExpr()
	}

	if p.match(tokens.Identifier) {
		token := p.peek(-1)
		if !p.match(tokens.Colon) {
			return &ast.IdentifierExpr{Name: token, ValueType: ast.Ident}
		}

		envPath := &ast.EnvPathExpr{
			Path: token,
		}

		next := p.advance()
		for p.match(tokens.Colon) {
			envPath.Combine(next)
			if next.Type != tokens.Identifier {
				p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(next))
				return &ast.Improper{Token: next}
			}
			next = p.advance()
		}
		envExpr := &ast.EnvAccessExpr{
			PathExpr: envPath,
		}
		envExpr.Accessed = &ast.IdentifierExpr{
			Name:      next,
			ValueType: ast.Invalid,
		}

		return envExpr
	}

	if p.match(tokens.LeftParen) {
		token := p.peek(-1)
		expr := p.expression()
		if expr.GetType() == ast.NA {
			p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(p.peek()))
		}
		p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewMulti(token, p.peek()), tokens.RightParen), tokens.RightParen)
		return &ast.GroupExpr{Expr: expr, Token: token, ValueType: expr.GetValueType()}
	}

	if p.match(tokens.Self) {
		return &ast.IdentifierExpr{Name: p.peek(-1)}
	}

	return &ast.Improper{Token: p.peek()}
}

func (p *Parser) list() ast.Node {
	token := p.peek(-1)
	list := make([]ast.Node, 0)
	if p.match(tokens.RightBracket) {
		return &ast.ListExpr{ValueType: ast.List, List: list, Token: token}
	}
	exprInList := p.expression()
	if exprInList.GetType() == ast.NA {
		p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(p.peek()))
	}
	list = append(list, exprInList)
	for p.match(tokens.Comma) {
		exprInList := p.expression()
		if exprInList.GetType() == ast.NA {
			p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(p.peek()))
			p.advance()
		}
		list = append(list, exprInList)
	}
	p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewMulti(token, p.peek()), tokens.RightBracket), tokens.RightBracket)

	return &ast.ListExpr{ValueType: ast.List, List: list, Token: token}
}

func (p *Parser) parseMap() ast.Node {
	token := p.peek(-1)
	parsedMap := make(map[tokens.Token]ast.Property, 0)
	for !p.check(tokens.RightBrace) {
		key := p.primary(true)

		var newKey tokens.Token
		switch key := key.(type) {
		case *ast.IdentifierExpr:
			newKey = key.GetToken()
		default:
			p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(key.GetToken()))
			p.advance()
			return &ast.Improper{Token: p.peek(-1)}
		}

		if _, ok := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.Equal, "after map key"), tokens.Equal); !ok {
			return &ast.Improper{Token: p.peek(-1)}
		}

		expr := p.expression()
		if expr.GetType() == ast.NA {
			p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(p.peek()))
		}

		if p.peek().Type == tokens.RightBrace {
			parsedMap[newKey] = ast.Property{Expr: expr, Type: expr.GetValueType()}
			break
		}

		if _, ok := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.Comma, "in map initialization"), tokens.Comma); !ok {
			return &ast.Improper{Token: p.peek(-1)}
		}

		parsedMap[newKey] = ast.Property{Expr: expr, Type: expr.GetValueType()}
	}
	p.advance()

	return &ast.MapExpr{Map: parsedMap, Token: token}
}

func (p *Parser) structExpr() ast.Node {
	structExpr := ast.StructExpr{
		Token:  p.peek(-1),
		Fields: make([]*ast.FieldDeclarationStmt, 0),
	}

	start, ok := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.LeftBrace), tokens.LeftBrace)
	if !ok {
		return &ast.Improper{Token: structExpr.Token}
	}
	if p.match(tokens.RightBrace) {
		return &structExpr
	}
	field := p.fieldDeclarationStmt()
	if field.GetType() != ast.NA {
		structExpr.Fields = append(structExpr.Fields, field.(*ast.FieldDeclarationStmt))
	} else {
		p.Alert(&alerts.ExpectedFieldDeclaration{}, alerts.NewSingle(field.GetToken()))
		return ast.NewImproper(field.GetToken())
	}

	for p.match(tokens.SemiColon) {
		if p.match(tokens.RightBrace) {
			return &structExpr
		}
		field := p.fieldDeclarationStmt()
		if field.GetType() != ast.NA {
			structExpr.Fields = append(structExpr.Fields, field.(*ast.FieldDeclarationStmt))
		} else {
			p.Alert(&alerts.ExpectedFieldDeclaration{}, alerts.NewSingle(field.GetToken()))
			return ast.NewImproper(field.GetToken())
		}
	}

	p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewMulti(start, p.peek()), tokens.RightBrace), tokens.RightBrace)

	return &structExpr
}

func (p *Parser) WrappedType() *ast.TypeExpr {
	typeExpr := ast.TypeExpr{}
	if p.check(tokens.Greater) {
		p.Alert(&alerts.EmptyWrappedType{}, alerts.NewSingle(p.peek()))
		return &typeExpr
	}

	return p.Type()
}

func (p *Parser) Type() *ast.TypeExpr {
	var expr ast.Node
	token := p.advance()
	if p.match(tokens.Colon) {
		if token.Type != tokens.Identifier {
			p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(token))
		}
		envAccess := &ast.EnvAccessExpr{
			PathExpr: &ast.EnvPathExpr{
				Path: token,
			},
		}
		next := p.advance()
		for p.match(tokens.Colon) {
			envAccess.PathExpr.Combine(next)
			if next.Type != tokens.Identifier {
				p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(next))
			}
			next = p.advance()
		}
		envAccess.Accessed = &ast.IdentifierExpr{
			Name: next,
		}
		expr = envAccess
	} else {
		expr = &ast.IdentifierExpr{
			Name:      token,
			ValueType: ast.Invalid,
		}
	}

	var typ *ast.TypeExpr = nil
	if expr.GetType() == ast.EnvironmentAccessExpression {
		typ = &ast.TypeExpr{Name: expr}
		typ.IsVariadic = p.match(tokens.Ellipsis)
		return typ
	}
	exprToken := expr.GetToken()

	switch exprToken.Type {
	case tokens.Identifier:
		typ = &ast.TypeExpr{}
		if p.match(tokens.Less) {
			typ.WrappedType = p.WrappedType()
			p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.Greater), tokens.Greater)
		}
		typ.Name = expr
	case tokens.Fn:
		typ = &ast.TypeExpr{}
		if !p.match(tokens.LeftParen) {
			typ.Name = expr
			break
		}
		if !p.match(tokens.RightParen) {
			_typ := p.Type()
			typ.Params = append(typ.Params, _typ)
			for p.match(tokens.Comma) {
				_typ := p.Type()
				typ.Params = append(typ.Params, _typ)
			}
			p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.RightParen), tokens.RightParen)
		}
		typ.Returns = p.returnings()
		typ.Name = expr
	case tokens.Struct:
		fields := p.parameters(tokens.LeftBrace, tokens.RightBrace)
		typ = &ast.TypeExpr{Name: expr, Fields: fields}
	case tokens.Entity:
		typ = &ast.TypeExpr{Name: &ast.IdentifierExpr{Name: exprToken}}
	default:
		p.Alert(&alerts.ExpectedType{}, alerts.NewSingle(expr.GetToken()))
		typ = &ast.TypeExpr{Name: &ast.Improper{Token: expr.GetToken()}}
	}
	typ.IsVariadic = p.match(tokens.Ellipsis)

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
	name, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "for a environment type expr"), tokens.Identifier)

	if !ok {
		return &ast.EnvTypeExpr{Type: ast.InvalidEnv, Token: name}
	}

	envType := StringToEnvType(name.Lexeme)

	if envType == ast.InvalidEnv {
		p.Alert(&alerts.InvalidEnvironmentType{}, alerts.NewSingle(name))
	}

	return &ast.EnvTypeExpr{Type: envType, Token: name}
}

func (p *Parser) EnvPathExpr() ast.Node {
	ident, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "for an environment path"), tokens.Identifier)

	if !ok {
		return &ast.Improper{Token: ident}
	}

	envPath := &ast.EnvPathExpr{
		Path: ident,
	}

	for p.match(tokens.Colon) {
		ident, ok = p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "in environment path"), tokens.Identifier)
		if !ok {
			return &ast.Improper{Token: envPath.GetToken()}
		}
		envPath.Combine(ident)
	}

	return envPath
}
