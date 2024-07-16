package parser

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
	"strings"
)

func (p *Parser) expression() ast.Node {
	return p.fn()
}

func (p *Parser) fn() ast.Node {
	if p.match(lexer.Fn) {
		fn := &ast.AnonFnExpr{
			Token: p.peek(-1),
		}
		fn.Params = p.parameters(lexer.LeftParen, lexer.RightParen)

		ret := make([]*ast.TypeExpr, 0)
		for p.check(lexer.Identifier) {
			ret = append(ret, p.Type())
			if !p.check(lexer.Comma) {
				break
			} else {
				p.advance()
			}
		}

		fn.Return = ret
		var success bool
		fn.Body, success = p.getBody()
		if !success {
			return &ast.Improper{}
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
	return p.accessorExprDepth2(nil, nil, ast.NA)
}

func (p *Parser) call(caller ast.Node) ast.Node {
	if !p.check(lexer.LeftParen) {
		return caller
	}

	callerType := caller.GetType()
	if callerType != ast.Identifier && callerType != ast.CallExpression {
		p.error(p.peek(-1), fmt.Sprintf("cannot call unidentified value (caller: %v)", callerType))
		return &ast.Improper{Token: p.peek(-1)}
	}

	call_expr := &ast.CallExpr{
		Identifier: caller.GetToken().Lexeme,
		Caller:     caller,
		Args:       p.arguments(),
		Token:      caller.GetToken(),
	}

	if p.check(lexer.LeftParen) {
		expr := p.call(call_expr)
		if expr.GetType() == ast.CallExpression {
			call_expr = expr.(*ast.CallExpr)
		}
	}

	return call_expr
}

func (p *Parser) resolveProperty(node *ast.Node) ast.Accessor {
	var accessor ast.Accessor
	if (*node).GetType() == ast.FieldExpression {
		accessor = (*node).(*ast.FieldExpr)
	} else {
		accessor = (*node).(*ast.MemberExpr)
	}
	prop := accessor.GetProperty()
	propAccessor := (*prop).(ast.Accessor)

	if *propAccessor.GetProperty() == nil {
		accessor.SetProperty(nil)
		*node = accessor
		return propAccessor
	}

	last := p.resolveProperty(prop)

	if fieldExpr, ok := (*node).(*ast.FieldExpr); ok {
		fieldExpr.Property = *prop
		*node = fieldExpr
	} else if memberExpr, ok := (*node).(*ast.MemberExpr); ok {
		memberExpr.Property = *prop
		*node = memberExpr
	}

	return last
}

func (p *Parser) accessorExprDepth2(owner ast.Accessor, ident ast.Node, nodeType ast.NodeType) ast.Node {
	expr := p.accessorExprDepth1(owner, ident, nodeType)

	if expr.GetType() != ast.FieldExpression && expr.GetType() != ast.MemberExpression {
		return p.call(expr)
	}
	if !p.check(lexer.LeftParen) {
		return expr
	}

	acesss := expr.(ast.Accessor)
	args := p.arguments()
	beforeExpr := acesss.DeepCopy()
	last := p.resolveProperty(&expr)
	if last.GetType() == ast.FieldExpression {
		expr = &ast.MethodCallExpr{
			Owner:      expr,
			Call:       beforeExpr,
			MethodName: last.GetToken().Lexeme,
			Args:       args,
			Token:      last.GetToken(),
		}
	} else {
		expr = &ast.CallExpr{
			Caller:     beforeExpr,
			Identifier: last.GetToken().Lexeme,
			Args:       args,
			Token:      last.GetToken(),
		}
	}

	expr = p.call(expr)

	isField, isMember := p.check(lexer.Dot), p.check(lexer.LeftBracket)

	var propNodeType ast.NodeType
	if isField {
		propNodeType = ast.FieldExpression
	} else {
		propNodeType = ast.MemberExpression
	}

	if !isField && !isMember {
		return expr
	}

	accessorExpr := p.accessorExprDepth2(owner, expr, propNodeType)

	return accessorExpr
}

func (p *Parser) accessorExprDepth1(owner ast.Accessor, ident ast.Node, nodeType ast.NodeType) ast.Node {
	if ident == nil {
		ident = p.matchExpr()
	}
	if owner == nil {
		ident = p.call(ident)
	}

	isField, isMember := p.check(lexer.Dot), p.check(lexer.LeftBracket)

	if !isField && !isMember {
		if owner == nil {
			return ident
		} else {
			if nodeType == ast.FieldExpression {
				return &ast.FieldExpr{
					Owner:      owner,
					Identifier: ident,
				}
			} else {
				return &ast.MemberExpr{
					Owner:      owner,
					Identifier: ident,
				}
			}
		}
	}

	p.advance()

	var propNodeType ast.NodeType
	if isField {
		propNodeType = ast.FieldExpression
	} else {
		propNodeType = ast.MemberExpression
	}
	if nodeType == ast.NA {
		nodeType = propNodeType
	}

	var expr ast.Accessor
	var prop ast.Node
	var propIdent ast.Node
	if nodeType == ast.FieldExpression {
		expr = &ast.FieldExpr{
			Owner:      owner,
			Identifier: ident,
		}
	} else {
		expr = &ast.MemberExpr{
			Owner:      owner,
			Identifier: ident,
		}
	}
	if propNodeType == ast.FieldExpression {
		propIdent = p.new()
	} else {
		propIdent = p.expression()
		p.consume("expected closing bracket in member expression", lexer.RightBracket)
	}
	prop = p.accessorExprDepth1(expr, propIdent, propNodeType)

	expr.SetProperty(prop)

	return expr
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
		caller := p.primary()
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

	return p.primary()
}

func (p *Parser) primary() ast.Node {
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
		}

		return &ast.LiteralExpr{Value: literal.Literal, ValueType: valueType, Token: literal}
	}

	if p.match(lexer.LeftBrace) {
		return p.parseMap()
	}

	if p.match(lexer.LeftBracket) {
		return p.list()
	}

	if p.match(lexer.Struct) {
		return p.anonStruct()
	}

	if p.match(lexer.Identifier) {
		token := p.peek(-1)
		if p.match(lexer.DoubleColon) {
			envPath := &ast.EnvPathExpr{
				SubPaths: []string{
					token.Lexeme,
				},
			}
	
			next := p.accessorExprDepth2(nil, nil, ast.NA)
			for p.match(lexer.DoubleColon) {
				envPath.SubPaths = append(envPath.SubPaths, next.GetToken().Lexeme)
				if next.GetType() != ast.Identifier {
					p.error(next.GetToken(), "expected identifier in environment expression")
					return &ast.Improper{Token: next.GetToken()}
				}
				next = p.accessorExprDepth2(nil, nil, ast.NA)
			}
			envExpr := &ast.EnvAccessExpr{
				PathExpr: envPath,
			}
			envExpr.Accessed = next
	
			return envExpr
		}
		return &ast.IdentifierExpr{Name: token, ValueType: ast.Ident}
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

	return &ast.ListExpr{ValueType: ast.List, List: list, Token: token}
}

func (p *Parser) parseMap() ast.Node {
	token := p.peek(-1)
	parsedMap := make(map[lexer.Token]ast.Property, 0)
	for !p.check(lexer.RightBrace) {
		key := p.primary()

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

		if _, ok := p.consume("expected ',' or '}' after expression", lexer.Comma, lexer.RightBrace); !ok {
			return &ast.Improper{Token: p.peek(-1)}
		}

		parsedMap[newKey] = ast.Property{Expr: expr, Type: expr.GetValueType()}
	}
	p.advance()

	return &ast.MapExpr{Map: parsedMap, Token: token}
}

func (p *Parser) anonStruct() ast.Node {
	anonStruct := ast.AnonStructExpr{
		Token:  p.peek(-1),
		Fields: make([]*ast.FieldDeclarationStmt, 0),
	}

	_, ok := p.consume("expected opening brace in anonymous struct expression", lexer.LeftBrace)
	if !ok {
		return &ast.Improper{Token: anonStruct.Token}
	}

	for !p.match(lexer.RightBrace) {
		field := p.fieldDeclarationStmt()
		if field.GetType() != ast.NA {
			anonStruct.Fields = append(anonStruct.Fields, field.(*ast.FieldDeclarationStmt))
		} else {
			p.error(field.GetToken(), "expected field declaration inside anonymous struct")
		}
		for p.match(lexer.Comma) {
			field := p.fieldDeclarationStmt()
			if field.GetType() != ast.NA {
				anonStruct.Fields = append(anonStruct.Fields, field.(*ast.FieldDeclarationStmt))
			} else {
				p.error(field.GetToken(), "expected field declaration inside anonymous struct")
			}
		}
	}

	return &anonStruct
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
	expr := p.primary()
	exprToken := expr.GetToken()

	if expr.GetType() == ast.EnvironmentAccessExpression {
		envAccess := expr.(*ast.EnvAccessExpr)
		if envAccess.Accessed.GetType() != ast.Identifier {
			p.error(envAccess.Accessed.GetToken(), "accessed type must be an identifier")
		}
		return &ast.TypeExpr{Name: expr}
	}

	switch exprToken.Type {
	case lexer.Identifier:
		typee := &ast.TypeExpr{}

		if p.match(lexer.Less) { // map<number>
			typee.WrappedType = p.WrappedType()
			p.consume("expected '>'", lexer.Greater)
		}
		typee.Name = expr
		return typee
	case lexer.Fn:
		typee := &ast.TypeExpr{}

		p.advance()
		params := make([]*ast.TypeExpr, 0)
		typee.Returns = make([]*ast.TypeExpr, 0)
		if p.match(lexer.LeftParen) {
			params = append(params, p.Type())

			for p.match(lexer.Comma) {
				params = append(params, p.Type())
			}
			p.consume("expected closing parenthesis in 'fn(...'", lexer.RightParen)
		}

		typee.Params = params

		if p.check(lexer.Identifier) {
			typee.Returns = append(typee.Returns, p.Type())

			for p.match(lexer.Comma) {
				typee.Returns = append(typee.Returns, p.Type())
			}
		}

		typee.Name = expr
		return typee
	case lexer.Struct:
		fields := p.parameters(lexer.LeftBrace, lexer.RightBrace)
		return &ast.TypeExpr{Name: expr, Fields: fields}
	default:
		p.error(exprToken, "Expected an identifier for a type")
		p.advance()
		return &ast.TypeExpr{Name: expr}
	}
}

func StringToEnvType(name string) ast.EnvType {
	switch name {
	case "Mesh":
		return ast.Mesh
	case "Level":
		return ast.Level
	case "Sound":
		return ast.Sound
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
		SubPaths: []string{ident.Lexeme},
	}

	for p.match(lexer.DoubleColon) {
		ident, ok = p.consume("expected identifier in environment path", lexer.Identifier)
		if !ok {
			return &ast.Improper{Token: ident}
		}
		envPath.SubPaths = append(envPath.SubPaths, ident.Lexeme)
	}

	return envPath
}
