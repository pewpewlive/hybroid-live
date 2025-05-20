package parser

import (
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/tokens"
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
			fn.Params, _ = p.functionParams(tokens.LeftParen, tokens.RightParen)
		} else {
			fn.Params = make([]ast.FunctionParam, 0)
			p.Alert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.LeftParen, "in function parameters")
		}
		fn.Return = p.functionReturns()

		var success bool
		fn.Body, success = p.body(false, true)
		if !success {
			return ast.NewImproper(fn.Token, ast.FunctionExpression)
		}

		return fn
	} else {
		return p.multiComparison()
	}
}

func (p *Parser) multiComparison() ast.Node {
	expr := p.comparison()
	if ast.IsImproper(expr, ast.NA) {
		return expr
	}

	if p.isMultiComparison() {
		operator := p.peek(-1)
		right := p.multiComparison()
		if ast.IsImproper(right, ast.NA) {
			p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(right.GetToken()), "as right value in binary expression")
		}
		expr = &ast.BinaryExpr{Left: expr, Operator: operator, Right: right, ValueType: ast.Bool}
	}

	return expr
}

func (p *Parser) comparison() ast.Node {
	expr := p.term()
	if ast.IsImproper(expr, ast.NA) {
		return expr
	}

	if p.isComparison() {
		operator := p.peek(-1)
		right := p.term()
		if ast.IsImproper(right, ast.NA) {
			p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(right.GetToken()), "as right value in binary expression")
		}
		expr = &ast.BinaryExpr{Left: expr, Operator: operator, Right: right, ValueType: ast.Bool}
	}

	return expr
}

func (p *Parser) term() ast.Node {
	expr := p.factor()
	if ast.IsImproper(expr, ast.NA) {
		return expr
	}
	if p.match(tokens.Plus, tokens.Minus) {
		operator := p.peek(-1)
		right := p.term()
		if ast.IsImproper(right, ast.NA) {
			p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(right.GetToken()), "as right value in binary expression")
		}
		expr = &ast.BinaryExpr{Left: expr, Operator: operator, Right: right, ValueType: ast.Uninitialized}
	}

	return expr
}

func (p *Parser) factor() ast.Node {
	expr := p.concat()
	if ast.IsImproper(expr, ast.NA) {
		return expr
	}

	if p.match(tokens.Star, tokens.Slash, tokens.Caret, tokens.Modulo, tokens.BackSlash) {
		operator := p.peek(-1)
		right := p.factor()
		if ast.IsImproper(right, ast.NA) {
			p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(right.GetToken()), "as right value in binary expression")
		}
		expr = &ast.BinaryExpr{Left: expr, Operator: operator, Right: right, ValueType: ast.Uninitialized}
	}

	return expr
}

func (p *Parser) concat() ast.Node {
	expr := p.unary()
	if ast.IsImproper(expr, ast.NA) {
		return expr
	}

	if p.match(tokens.Concat) {
		operator := p.peek(-1)
		right := p.concat()
		if ast.IsImproper(right, ast.NA) {
			p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(right.GetToken()), "as right value in concat expression")
		}
		return &ast.BinaryExpr{Left: expr, Operator: operator, Right: right, ValueType: ast.Uninitialized}
	}

	return expr
}

func (p *Parser) unary() ast.Node {
	if p.match(tokens.Bang, tokens.Minus, tokens.Hash) {
		operator := p.peek(-1)
		right := p.unary()
		if ast.IsImproper(right, ast.NA) {
			p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(right.GetToken()), "in unary expression")
		}
		return &ast.UnaryExpr{Operator: operator, Value: right}
	}
	return p.entity()
}

func (p *Parser) entity() ast.Node {
	letMatched := false
	var token tokens.Token
	if p.match(tokens.Let) {
		letMatched = true
		token = p.peek(-1)
	}

	variable := p.accessorExprDepth2(nil)
	var expr ast.Node
	var conv *tokens.Token
	if letMatched {
		if p.match(tokens.Equal) {
			tkn := variable.GetToken()
			conv = &tkn

			expr = p.accessorExprDepth2(nil)
		} else {
			p.Alert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.Equal, "in entity expression")
		}
	} else {
		expr = variable
		token = variable.GetToken()
	}
	if p.match(tokens.Is, tokens.Isnt) {
		if conv != nil {
			if variable.GetType() != ast.Identifier {
				p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(variable.GetToken()))
			}
		}
		op := p.peek(-1)
		typ := p.typeExpr("in entity expression")

		return &ast.EntityExpr{
			Expr:             expr,
			Type:             typ,
			ConvertedVarName: conv,
			Token:            token,
			Operator:         op,
		}
	}

	if letMatched {
		p.Alert(&alerts.ExpectedKeyword{}, alerts.NewSingle(p.peek()), "is/isnt")
		return ast.NewImproper(token, ast.EntityExpression)
	}

	return variable
}

func (p *Parser) call(caller ast.Node) ast.Node {
	hasGenerics := false
	genericArgs := []*ast.TypeExpr{}
	if p.check(tokens.Less) {
		ok := p.tryGenericArgs()
		if !ok {
			return caller
		}
		genericArgs, _ = p.genericArgs()
		hasGenerics = true
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
		p.functionArgs()
		return ast.NewImproper(caller.GetToken(), ast.CallExpression)
	}

	args, ok := p.functionArgs()
	if !ok {
		return ast.NewImproper(caller.GetToken(), ast.CallExpression)
	}

	call_expr := &ast.CallExpr{
		Caller:      caller,
		GenericArgs: genericArgs,
		Args:        args,
	}

	return p.call(call_expr)
}

func (p *Parser) accessorExprDepth2(ident *ast.Node) ast.Node {
	expr, call := p.accessorExpr(ident)

	if call == nil {
		return expr
	}

	genericsArgs, _ := p.genericArgs()
	args, ok := p.functionArgs()
	if !ok {
		return ast.NewImproper(call.GetToken(), ast.MethodCallExpression)
	}

	var methodCall ast.Node = &ast.MethodCallExpr{
		Identifier: expr,
		Call: &ast.CallExpr{
			Caller:      call,
			GenericArgs: genericsArgs,
			Args:        args,
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
			p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(propIdentifier.GetToken()), "in field expression")
		}
		if p.check(tokens.LeftParen) || p.check(tokens.Less) {
			return *ident, propIdentifier.(*ast.IdentifierExpr)
		}
	} else if isMember {
		propIdentifier = p.expression()
		if ast.IsImproper(propIdentifier, ast.NA) {
			p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(propIdentifier.GetToken()), "in member expression")
		}

		p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewMulti(start, p.peek()), tokens.RightBracket, "in member expression"), tokens.RightBracket)
	}

	expr.SetPropertyIdentifier(propIdentifier)
	if p.check(tokens.Dot) || p.check(tokens.LeftBracket) {
		prop, call := p.accessorExpr(&propIdentifier)

		expr.SetProperty(prop)

		return expr, call
	}
	expr.SetProperty(propIdentifier)

	return p.call(expr), nil
}

func (p *Parser) matchExpr() ast.Node {
	if p.match(tokens.Match) {
		start := p.peek(-1)
		node := p.matchStmt(true)
		if ast.IsImproper(node, ast.NA) {
			return node
		}
		if node.GetType() == ast.NA {
			return ast.NewImproper(start, ast.MatchExpression)
		}
		return &ast.MatchExpr{MatchStmt: *node.(*ast.MatchStmt)}
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
			return ast.NewImproper(p.peek(), ast.MacroCallExpression)
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

		expr.Type = p.typeExpr("in new expression")
		args, ok := p.functionArgs()
		if !ok {
			return ast.NewImproper(expr.Token, ast.NewExpession)
		}
		expr.Args = args

		return &expr
	}

	return p.spawn()
}

func (p *Parser) spawn() ast.Node {
	if p.match(tokens.Spawn) {
		expr := ast.SpawnExpr{
			Token: p.peek(-1),
		}

		expr.Type = p.typeExpr("in spawn expression")
		args, ok := p.functionArgs()
		if !ok {
			return ast.NewImproper(expr.Token, ast.SpawnExpression)
		}
		expr.Args = args

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

		switch literal.Type {
		case tokens.Number:
			valueType = ast.Number
		case tokens.Fixed:
			valueType = ast.Fixed
		case tokens.FixedPoint:
			valueType = ast.FixedPoint
		case tokens.Degree:
			valueType = ast.Degree
		case tokens.Radian:
			valueType = ast.Radian
		case tokens.String:
			valueType = ast.String
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
				return ast.NewImproper(next, ast.EnvironmentPathExpression)
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
		if ast.IsImproper(expr, ast.NA) {
			p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(p.peek()), "in group expression")
			return ast.NewImproper(token, ast.GroupExpression)
		}
		_, ok := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.RightParen, "in group expression"), tokens.RightParen)
		if !ok {
			return ast.NewImproper(token, ast.GroupExpression)
		}
		return &ast.GroupExpr{Expr: expr, Token: token, ValueType: expr.GetValueType()}
	}

	if p.match(tokens.Self) {
		return &ast.IdentifierExpr{Name: p.peek(-1)}
	}

	return ast.NewImproper(p.peek(), ast.NA)
}

func (p *Parser) list() ast.Node {
	listExpr := &ast.ListExpr{
		ValueType: ast.List,
		Token:     p.peek(-1),
	}

	if p.match(tokens.RightBracket) {
		return listExpr
	}

	list, ok := p.expressions("in list expression", true)
	if !ok {
		return ast.NewImproper(listExpr.Token, ast.ListExpression)
	}
	listExpr.List = list

	p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewMulti(listExpr.Token, p.peek()), tokens.RightBracket, "in list expression"), tokens.RightBracket)

	return listExpr
}

func (p *Parser) parseMap() ast.Node {
	mapExpr := &ast.MapExpr{
		Token:        p.peek(-1),
		KeyValueList: make([]ast.Property, 0),
	}

	if p.match(tokens.RightBrace) {
		return mapExpr
	}

	p.context.braceCounter.Increment()
	defer p.context.braceCounter.Decrement()

	key, value, ok := p.keyValuePair("map key")
	if !ok {
		return ast.NewImproper(mapExpr.Token, ast.MapExpression)
	}
	mapExpr.KeyValueList = append(mapExpr.KeyValueList, ast.Property{Key: key, Expr: value, Type: value.GetValueType()})

	for p.match(tokens.Comma) {
		key, value, ok = p.keyValuePair("map key")
		if !ok {
			return ast.NewImproper(mapExpr.Token, ast.MapExpression)
		}
		mapExpr.KeyValueList = append(mapExpr.KeyValueList, ast.Property{Key: key, Expr: value, Type: value.GetValueType()})
	}
	p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewMulti(mapExpr.Token, p.peek()), tokens.RightBrace, "in map expression"), tokens.RightBrace)

	return mapExpr
}

func (p *Parser) structExpr() ast.Node {
	structExpr := ast.StructExpr{
		Token:  p.peek(-1),
		Fields: make([]*ast.FieldDecl, 0),
	}

	start, ok := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.LeftBrace), tokens.LeftBrace)
	if !ok {
		return ast.NewImproper(structExpr.Token, ast.StructExpression)
	}
	if p.match(tokens.RightBrace) {
		return &structExpr
	}
	p.context.braceCounter.Increment()
	defer p.context.braceCounter.Decrement()

	field, value, ok := p.keyValuePair("struct field")
	if !ok {
		return ast.NewImproper(structExpr.Token, ast.StructExpression)
	}
	structExpr.Fields = append(structExpr.Fields, &ast.FieldDecl{
		Identifiers: []*ast.IdentifierExpr{field.(*ast.IdentifierExpr)},
		Values:      []ast.Node{value},
		Token:       field.GetToken(),
	})

	for p.match(tokens.SemiColon) {
		if p.check(tokens.RightBrace) {
			break
		}
		field, value, ok := p.keyValuePair("struct field")
		if !ok {
			return ast.NewImproper(structExpr.Token, ast.StructExpression)
		}
		structExpr.Fields = append(structExpr.Fields, &ast.FieldDecl{
			Identifiers: []*ast.IdentifierExpr{field.(*ast.IdentifierExpr)},
			Values:      []ast.Node{value},
			Token:       field.GetToken(),
		})
	}

	p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewMulti(start, p.peek()), tokens.RightBrace), tokens.RightBrace)

	return &structExpr
}

func (p *Parser) wrappedTypeExpr(typeContext string) *ast.TypeExpr {
	typeExpr := ast.TypeExpr{}
	if p.check(tokens.Greater) {
		p.Alert(&alerts.EmptyWrappedType{}, alerts.NewSingle(p.peek()))
		return &typeExpr
	}

	return p.typeExpr(typeContext)
}

func (p *Parser) typeExpr(typeContext string) *ast.TypeExpr {
	var expr ast.Node
	token := p.advance()
	improperType := &ast.TypeExpr{Name: ast.NewImproper(token, ast.TypeExpression)}

	if token.Type == tokens.LeftParen {
		tuple := &ast.TupleExpr{LeftParen: token}

		types := []*ast.TypeExpr{}

		typ := p.typeExpr(typeContext)
		if typ.Name.GetType() == ast.NA {
			return improperType
		}
		types = append(types, typ)

		for p.match(tokens.Comma) {
			typ := p.typeExpr(typeContext)
			if typ.Name.GetType() == ast.NA {
				return improperType
			}
			types = append(types, typ)
		}
		_, ok := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.RightParen, "in tuple expression"), tokens.RightParen)

		tuple.Types = types
		if !ok {
			return improperType
		}

		return &ast.TypeExpr{Name: tuple}
	}

	if p.match(tokens.Colon) {
		if token.Type != tokens.Identifier {
			p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(token))
			return improperType
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
				return improperType
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

	typeExpr := ast.TypeExpr{}
	if expr.GetType() == ast.EnvironmentAccessExpression {
		typeExpr = ast.TypeExpr{Name: expr}
		typeExpr.IsVariadic = p.match(tokens.Ellipsis)
		return &typeExpr
	}
	exprToken := expr.GetToken()

	switch exprToken.Type {
	case tokens.Identifier:
		if p.match(tokens.Less) {
			typeExpr.WrappedType = p.wrappedTypeExpr(typeContext)
			_, ok := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.Greater), tokens.Greater)
			if !ok {
				return improperType
			}
		}
		typeExpr.Name = expr
	case tokens.Fn: // fn
		_, ok := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.LeftParen, "after 'fn' in type expression"), tokens.LeftParen)
		if !ok {
			return improperType
		}
		if !p.match(tokens.RightParen) {
			typ := p.typeExpr(typeContext)
			if typ.Name.GetType() == ast.NA {
				return improperType
			}
			typeExpr.Params = append(typeExpr.Params, typ)
			for p.match(tokens.Comma) {
				typ := p.typeExpr(typeContext)
				if typ.Name.GetType() == ast.NA {
					return improperType
				}
				typeExpr.Params = append(typeExpr.Params, typ)
			}
			_, ok := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.RightParen, "in fn type expression"), tokens.RightParen)
			if !ok {
				return improperType
			}
		}
		returns := p.functionReturns()
		if returns != nil {
			if returns.Name.GetType() == ast.NA {
				return improperType
			}
			typeExpr.Return = returns
		}
		typeExpr.Name = expr
	case tokens.Struct:
		fields, ok := p.functionParams(tokens.LeftBrace, tokens.RightBrace)
		if !ok {
			return improperType
		}
		typeExpr.Name = expr
		typeExpr.Fields = fields
	case tokens.Entity:
		typeExpr.Name = &ast.IdentifierExpr{Name: exprToken}
	default:
		p.Alert(&alerts.ExpectedType{}, alerts.NewSingle(expr.GetToken()), typeContext)
		typeExpr.Name = ast.NewImproper(expr.GetToken(), ast.NA)
	}
	typeExpr.IsVariadic = p.match(tokens.Ellipsis)

	return &typeExpr
}

func (p *Parser) envTypeExpr() *ast.EnvTypeExpr {
	envTypeExpr := &ast.EnvTypeExpr{
		Type: ast.InvalidEnv,
	}

	token, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "for an environment type"), tokens.Identifier)
	if ok {
		p.coherencyCheck(p.peek(-2), token)
	}

	return envTypeExpr
}

func (p *Parser) envPathExpr() ast.Node {
	envPath := &ast.EnvPathExpr{}

	ident, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "for an environment path"), tokens.Identifier)
	if !ok {
		return ast.NewImproper(ident, ast.EnvironmentPathExpression)
	}
	envPath.Path = ident

	for p.match(tokens.Colon) {
		ident, ok = p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "in environment path"), tokens.Identifier)
		if !ok {
			return ast.NewImproper(ident, ast.EnvironmentPathExpression)
		}
		envPath.Combine(ident)
	}

	return envPath
}
