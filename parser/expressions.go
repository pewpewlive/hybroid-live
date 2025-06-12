package parser

import (
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/tokens"
)

func (p *Parser) expression() ast.Node {
	return p.mapExpr()
}

func (p *Parser) mapExpr() ast.Node {
	if p.match(tokens.LeftBrace) {
		return p.parseMap()
	}
	return p.list()
}

func (p *Parser) list() ast.Node {
	if p.match(tokens.LeftBracket) {
		return p.parseList()
	}
	return p.anonStruct()
}

func (p *Parser) anonStruct() ast.Node {
	if p.match(tokens.Struct) {
		return p.structExpr()
	}

	return p.fn()
}

func (p *Parser) fn() ast.Node {
	if p.match(tokens.Fn) {
		fn := &ast.FunctionExpr{
			Token: p.peek(-1),
		}
		gens, ok := p.genericParams()
		if !ok {
			return ast.NewImproper(fn.Token, ast.FunctionExpression)
		}
		fn.Generics = gens
		params, ok := p.functionParams(tokens.LeftParen, tokens.RightParen)
		if !ok {
			return ast.NewImproper(fn.Token, ast.FunctionExpression)
		}
		fn.Params = params
		returns, ok := p.functionReturns()
		if !ok {
			return ast.NewImproper(fn.Token, ast.FunctionExpression)
		}
		fn.Returns = returns

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
			p.AlertSingle(&alerts.ExpectedExpression{}, right.GetToken(), "as right value in binary expression")
		}
		expr = &ast.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) comparison() ast.Node {
	expr := p.term()
	if ast.IsImproper(expr, ast.NA) {
		return expr
	}

	if op, ok := p.isComparison(); ok {
		right := p.term()
		if ast.IsImproper(right, ast.NA) {
			p.AlertSingle(&alerts.ExpectedExpression{}, right.GetToken(), "as right value in binary expression")
		}
		expr = &ast.BinaryExpr{Left: expr, Operator: op, Right: right}
	}

	return expr
}

func (p *Parser) term() ast.Node {
	expr := p.factor()
	if ast.IsImproper(expr, ast.NA) {
		return expr
	}

	isLeftShift := p.peek().Type == tokens.Less && p.peek(1).Type == tokens.Less && p.peek(2).Type != tokens.Equal
	isRightShift := p.peek().Type == tokens.Greater && p.peek(1).Type == tokens.Greater && p.peek(2).Type != tokens.Equal
	isNormalOp := p.check(tokens.Plus, tokens.Minus, tokens.Ampersand, tokens.Pipe, tokens.Tilde)

	var op tokens.Token
	if isNormalOp {
		op = p.advance()
	} else if isLeftShift {
		newToken, success := p.combineTokens(tokens.LeftShift, 2)
		if !success {
			return expr
		}
		op = newToken
	} else if isRightShift {
		newToken, success := p.combineTokens(tokens.RightShift, 2)
		if !success {
			return expr
		}
		op = newToken
	}

	if isNormalOp || isLeftShift || isRightShift {
		right := p.term()
		if ast.IsImproper(right, ast.NA) {
			p.AlertSingle(&alerts.ExpectedExpression{}, right.GetToken(), "as right value in binary expression")
		}
		expr = &ast.BinaryExpr{Left: expr, Operator: op, Right: right}
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
			p.AlertSingle(&alerts.ExpectedExpression{}, right.GetToken(), "as right value in binary expression")
		}
		expr = &ast.BinaryExpr{Left: expr, Operator: operator, Right: right}
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
			p.AlertSingle(&alerts.ExpectedExpression{}, right.GetToken(), "as right value in concat expression")
		}
		return &ast.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) unary() ast.Node {
	if p.match(tokens.Bang, tokens.Minus, tokens.Hash) {
		operator := p.peek(-1)
		right := p.unary()
		if ast.IsImproper(right, ast.NA) {
			p.AlertSingle(&alerts.ExpectedExpression{}, right.GetToken(), "in unary expression")
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

	variable := p.AccessorExpr()
	var expr ast.Node
	var conv *tokens.Token
	if letMatched {
		if p.match(tokens.Equal) {
			tkn := variable.GetToken()
			conv = &tkn

			expr = p.AccessorExpr()
		} else {
			p.AlertSingle(&alerts.ExpectedSymbol{}, p.peek(), tokens.Equal, "in entity expression")
		}
	} else {
		expr = variable
		token = variable.GetToken()
	}
	if p.match(tokens.Is, tokens.Isnt) {
		if conv != nil {
			if variable.GetType() != ast.Identifier {
				p.AlertSingle(&alerts.ExpectedIdentifier{}, variable.GetToken())
			}
		}
		op := p.peek(-1)
		typ := p.typeExpr("in entity expression")

		return &ast.EntityEvaluationExpr{
			Expr:             expr,
			Type:             typ,
			ConvertedVarName: conv,
			Token:            token,
			Operator:         op,
		}
	}

	if letMatched {
		p.AlertSingle(&alerts.ExpectedKeyword{}, p.peek(), "is/isnt")
		return ast.NewImproper(token, ast.EntityEvaluationExpression)
	}

	return variable
}

func (p *Parser) AccessorExpr() ast.Node {
	expr := p.call(p.matchExpr())

accessCheck:
	var access *ast.AccessExpr
	if p.check(tokens.Dot, tokens.LeftBracket) {
		access = &ast.AccessExpr{
			Start:    expr,
			Accessed: []ast.Node{},
		}
	}
	for p.match(tokens.Dot, tokens.LeftBracket) {
		tokenType := p.peek(-1).Type

		if tokenType == tokens.Dot { // self.frames[]
			expr2 := p.primary()
			access.Accessed = append(access.Accessed, &ast.FieldExpr{
				Field: expr2,
			})
		} else {
			expr2 := p.expression()
			p.alertSingleConsume(&alerts.ExpectedSymbol{}, tokens.RightBracket, "in member expression")
			access.Accessed = append(access.Accessed, &ast.MemberExpr{
				Member: expr2,
			})
		}
		expr = p.call(access)
		if expr.GetType() == ast.CallExpression {
			goto accessCheck
		}
	}

	return expr
}

func (p *Parser) call(caller ast.Node) ast.Node {
	genericArgs := []*ast.TypeExpr{}
	hasGenerics := false
	if p.check(tokens.Less) {
		ok := p.tryGenericArgs()
		if !ok {
			return caller
		}
		hasGenerics = true
		genericArgs, ok = p.genericArgs()
		if !ok {
			return ast.NewImproper(caller.GetToken(), ast.CallExpression)
		}
	}

	if !hasGenerics && !p.check(tokens.LeftParen) {
		return caller
	}

	args, ok := p.functionArgs()
	if !ok {
		return ast.NewImproper(caller.GetToken(), ast.CallExpression)
	}

	callExpr := &ast.CallExpr{
		Caller:      caller,
		GenericArgs: genericArgs,
		Args:        args,
	}

	return p.call(callExpr)
}

func (p *Parser) matchExpr() ast.Node {
	if p.match(tokens.Match) {
		start := p.peek(-1)
		node := p.matchStatement(true)
		if ast.IsImproper(node, ast.NA) {
			return node
		}
		if node.GetType() == ast.NA {
			return ast.NewImproper(start, ast.MatchExpression)
		}
		return &ast.MatchExpr{MatchStmt: *node.(*ast.MatchStmt)}
	}

	return p.macroCall() // We do not parse macros right now
}

func (p *Parser) macroCall() ast.Node {
	// if p.match(tokens.At) {
	// 	macroCall := &ast.MacroCallExpr{}
	// 	caller := p.primary()
	// 	callerType := caller.GetType()
	// 	if callerType != ast.CallExpression {
	// 		p.AlertSingle(&alerts.ExpectedCallAfterMacroSymbol{}, caller.GetToken()))
	// 		return ast.NewImproper(p.peek(), ast.MacroCallExpression)
	// 	}
	// 	macroCall.Caller = caller.(*ast.CallExpr)
	// 	return macroCall
	// }

	return p.new()
}

func (p *Parser) new() ast.Node {
	if p.match(tokens.New) {
		expr := ast.NewExpr{
			Token: p.peek(-1),
		}
		// new<T, E>
		classGenericArgs, ok := p.genericArgs()
		if !ok {
			return ast.NewImproper(expr.Token, ast.NewExpession)
		}
		expr.ClassGenericArgs = classGenericArgs

		// new<T, E> Type<F, G>
		expr.Type = p.typeExpr("in new expression")
		if ast.IsImproper(expr.Type.Name, ast.NA) {
			return ast.NewImproper(expr.Token, ast.NewExpession)
		}

		// new<T, E> Type<F, G>(...)
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

		// spawn<T, E>
		classGenericArgs, ok := p.genericArgs()
		if !ok {
			return ast.NewImproper(expr.Token, ast.NewExpession)
		}
		expr.GenericArgs = classGenericArgs

		// spawn<T, E> Type<F, G>
		expr.Type = p.typeExpr("in spawn expression")

		// spawn<T, E> Type<F, G>(...)
		args, ok := p.functionArgs()
		if !ok {
			return ast.NewImproper(expr.Token, ast.NewExpession)
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

	return p.primary()
}

func (p *Parser) primary() ast.Node {
	if p.match(tokens.False) {
		return &ast.LiteralExpr{Value: "false", Token: p.peek(-1)}
	}
	if p.match(tokens.True) {
		return &ast.LiteralExpr{Value: "true", Token: p.peek(-1)}
	}

	if p.match(tokens.Number, tokens.Fixed, tokens.FixedPoint, tokens.Degree, tokens.Radian, tokens.String) {
		literal := p.peek(-1)
		return &ast.LiteralExpr{Value: literal.Literal, Token: literal}
	}

	if p.match(tokens.Identifier) {
		token := p.peek(-1)
		if !p.match(tokens.Colon) {
			return &ast.IdentifierExpr{Name: token}
		}

		envPath := &ast.EnvPathExpr{
			Path: token,
		}

		next := p.advance()
		for p.match(tokens.Colon) {
			envPath.Combine(next)
			if next.Type != tokens.Identifier {
				p.AlertSingle(&alerts.ExpectedIdentifier{}, next)
				return ast.NewImproper(next, ast.EnvironmentPathExpression)
			}
			next = p.advance()
		}
		envExpr := &ast.EnvAccessExpr{
			PathExpr: envPath,
		}
		envExpr.Accessed = &ast.IdentifierExpr{
			Name: next,
		}

		return envExpr
	}

	if p.match(tokens.LeftParen) {
		token := p.peek(-1)
		expr := p.expression()
		if ast.IsImproper(expr, ast.NA) {
			p.AlertSingle(&alerts.ExpectedExpression{}, p.peek(), "in group expression")
			return ast.NewImproper(token, ast.GroupExpression)
		}
		_, ok := p.alertSingleConsume(&alerts.ExpectedSymbol{}, tokens.RightParen, "in group expression")
		if !ok {
			return ast.NewImproper(token, ast.GroupExpression)
		}
		return &ast.GroupExpr{Expr: expr, Token: token}
	}

	if p.match(tokens.Self) {
		return &ast.IdentifierExpr{Name: p.peek(-1)}
	}

	return ast.NewImproper(p.peek(), ast.NA)
}

func (p *Parser) parseList() ast.Node {
	listExpr := &ast.ListExpr{
		Token: p.peek(-1),
	}

	if p.match(tokens.RightBracket) {
		return listExpr
	}

	list, ok := p.expressions("in list expression", true)
	if !ok && !p.sync(tokens.RightBracket) {
		return ast.NewImproper(listExpr.Token, ast.ListExpression)
	}
	listExpr.List = list

	_, ok = p.alertMultiConsume(&alerts.ExpectedSymbol{}, listExpr.Token, p.peek(), tokens.RightBracket, "in list expression")
	if !ok && p.sync(tokens.RightBracket) {
		p.advance()
	}

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

	key, value, ok := p.keyValuePair(true, "map key")
	if !ok {
		p.sync(tokens.Comma, tokens.RightBrace)
	} else {
		mapExpr.KeyValueList = append(mapExpr.KeyValueList, ast.Property{Key: key, Expr: value})
	}

	for p.match(tokens.Comma) {
		key, value, ok = p.keyValuePair(true, "map key")
		if !ok {
			p.sync(tokens.Comma, tokens.RightBrace)
		} else {
			mapExpr.KeyValueList = append(mapExpr.KeyValueList, ast.Property{Key: key, Expr: value})
		}
	}

	_, ok = p.alertMultiConsume(&alerts.ExpectedSymbol{}, mapExpr.Token, p.peek(), tokens.RightBrace, "in map expression")
	if !ok && p.sync(tokens.RightBrace) {
		p.advance()
	}

	return mapExpr
}

func (p *Parser) structExpr() ast.Node {
	structExpr := ast.StructExpr{
		Token:       p.peek(-1),
		Fields:      make([]*ast.IdentifierExpr, 0),
		Expressions: make([]ast.Node, 0),
	}

	start, ok := p.alertSingleConsume(&alerts.ExpectedSymbol{}, tokens.LeftBrace)
	if !ok {
		return ast.NewImproper(structExpr.Token, ast.StructExpression)
	}
	if p.match(tokens.RightBrace) {
		return &structExpr
	}

	field, value, ok := p.keyValuePair(false, "struct field")
	if !ok {
		p.sync(tokens.Comma, tokens.RightBrace)
	} else {
		structExpr.Fields = append(structExpr.Fields, field.(*ast.IdentifierExpr))
		structExpr.Expressions = append(structExpr.Expressions, value)
	}

	for p.match(tokens.Comma) {
		if p.check(tokens.RightBrace) {
			break
		}
		field, value, ok := p.keyValuePair(false, "struct field")
		if !ok {
			p.sync(tokens.Comma, tokens.RightBrace)
		} else {
			structExpr.Fields = append(structExpr.Fields, field.(*ast.IdentifierExpr))
			structExpr.Expressions = append(structExpr.Expressions, value)
		}
	}

	_, ok = p.alertMultiConsume(&alerts.ExpectedSymbol{}, start, p.peek(), tokens.RightBrace)
	if !ok && p.sync(tokens.RightBrace) {
		p.advance()
	}

	return &structExpr
}

func (p *Parser) wrappedTypeExpr(typeContext string) *ast.TypeExpr {
	typeExpr := ast.TypeExpr{}
	if p.check(tokens.Greater) {
		p.AlertSingle(&alerts.EmptyWrappedType{}, p.peek())
		return &typeExpr
	}

	return p.typeExpr(typeContext)
}

func (p *Parser) typeExpr(typeContext string) *ast.TypeExpr {
	var expr ast.Node
	token := p.advance()
	improperType := &ast.TypeExpr{Name: ast.NewImproper(token, ast.TypeExpression)}

	if p.match(tokens.Colon) {
		if token.Type != tokens.Identifier {
			p.AlertSingle(&alerts.ExpectedIdentifier{}, token)
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
				p.AlertSingle(&alerts.ExpectedIdentifier{}, next)
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
			Name: token,
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
			wrapped := p.wrappedTypeExpr(typeContext)
			if wrapped.Name.GetType() == ast.NA {
				return improperType
			}
			typeExpr.WrappedTypes = append(typeExpr.WrappedTypes, wrapped)
			for p.match(tokens.Comma) {
				wrapped := p.wrappedTypeExpr(typeContext)
				if wrapped.Name.GetType() == ast.NA {
					return improperType
				}
				typeExpr.WrappedTypes = append(typeExpr.WrappedTypes, wrapped)
			}
			_, ok := p.alertSingleConsume(&alerts.ExpectedSymbol{}, tokens.Greater)
			if !ok {
				return improperType
			}
		}
		typeExpr.Name = expr
	case tokens.Fn: // fn
		_, ok := p.alertSingleConsume(&alerts.ExpectedSymbol{}, tokens.LeftParen, "after 'fn' in type expression")
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
			_, ok := p.alertSingleConsume(&alerts.ExpectedSymbol{}, tokens.RightParen, "in fn type expression")
			if !ok {
				return improperType
			}
		}
		returns, ok := p.functionReturns()
		if returns != nil {
			if !ok {
				return improperType
			}
			typeExpr.Returns = returns
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
		p.disadvance()
		p.AlertSingle(&alerts.ExpectedType{}, expr.GetToken(), typeContext)
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
	envTypeExpr.Token = token

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
