package parser

import (
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/tokens"
)

func (p *Parser) environmentDeclaration() ast.Node {
	envDecl := &ast.EnvironmentDecl{}

	pathExpr := p.envPathExpr()
	if pathExpr.GetType() == ast.NA {
		return pathExpr
	}

	if _, ok := p.alertSingleConsume(&alerts.ExpectedKeyword{}, tokens.As, "in environment declaration"); !ok {
		return ast.NewImproper(p.peek(), ast.EnvironmentDeclaration)
	}

	envPathExpr, _ := pathExpr.(*ast.EnvPathExpr)

	envDecl.EnvType = p.envTypeExpr()
	envDecl.Env = envPathExpr

	return envDecl
}

func (p *Parser) simpleVariableDeclaration() ast.Node {
	varDecl := &ast.VariableDecl{
		Token: p.peek(-1),
		IsPub: p.context.isPub,
	}
	varDecl.IsConst = varDecl.Token.Type == tokens.Const
	if varDecl.IsPub && varDecl.Token.Type == tokens.Pub {
		current := p.current
		p.context.ignoreAlerts.Push("VariableDeclaration", true)
		typeExpr := p.typeExpr("")
		p.context.ignoreAlerts.Pop("VariableDeclaration")
		next := p.peek()
		p.disadvance(p.current - current)

		if typeExpr != nil && !ast.IsImproper(typeExpr.Name, ast.NA) {
			if next.Type == tokens.Equal && typeExpr.Name.GetType() == ast.Identifier {
			} else {
				varDecl.Type = p.typeExpr("in variable declaration")
			}
		} else {
			varDecl.Type = p.typeExpr("in variable declaration")
		}
	}
	if varDecl.IsPub && varDecl.Token.Type != tokens.Pub {
		varDecl.Token = p.peek(-2)
	}
	if varDecl.IsPub && varDecl.Token.Type == tokens.Let { // pub let a
		p.Alert(&alerts.UnexpectedKeyword{}, alerts.NewSingle(p.peek(-2)), "pub", "in variable declaration")
	}

	idents, exprs, ok := p.identExprPairs("in variable declaration", false)
	if !ok {
		return ast.NewImproper(varDecl.Token, ast.VariableDeclaration)
	}

	varDecl.Identifiers = idents
	varDecl.Expressions = exprs

	return varDecl
}

func (p *Parser) typedVariableDeclaration() ast.Node {
	varDecl := &ast.VariableDecl{
		Token: p.peek(-1),
		IsPub: p.context.isPub,
	}

	currentStart := p.current
	p.context.ignoreAlerts.Push("CheckType", true)

	typeExpr := p.typeExpr("")

	p.context.ignoreAlerts.Pop("CheckType")
	valid := typeExpr != nil && typeExpr.Name.GetType() != ast.NA
	p.disadvance(p.current - currentStart)

	if !valid {
		return ast.NewImproper(varDecl.Token, ast.NA)
	}
	typeExpr = p.typeExpr("in variable declaration")
	varDecl.Type = typeExpr

	idents, exprs, ok := p.identExprPairs("in variable declaration", true)
	if !ok {
		return ast.NewImproper(varDecl.Token, ast.NA)
	}
	varDecl.Identifiers = idents
	varDecl.Expressions = exprs

	return varDecl
}

func (p *Parser) functionDeclaration() ast.Node {
	functionDecl := ast.FunctionDecl{
		IsPub: p.context.isPub,
	}
	if functionDecl.IsPub {
		functionDecl.Token = p.peek(-2)
	} else {
		functionDecl.Token = p.peek(-1)
	}

	name, nameOk := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "as the name of the function"), tokens.Identifier)
	if !nameOk {
		return ast.NewImproper(functionDecl.Token, ast.NA)
	}

	functionDecl.Name = name
	generics, ok := p.genericParams()
	if !ok {
		return ast.NewImproper(functionDecl.Token, ast.FunctionDeclaration)
	}
	functionDecl.Generics = generics
	params, ok := p.functionParams(tokens.LeftParen, tokens.RightParen)
	if !ok {
		return ast.NewImproper(functionDecl.Name, ast.FunctionDeclaration)
	}
	functionDecl.Params = params
	returns, ok := p.functionReturns()
	if !ok {
		return ast.NewImproper(functionDecl.Name, ast.FunctionDeclaration)
	}
	functionDecl.Returns = returns

	body, ok := p.body(false, true)
	if !ok {
		return ast.NewImproper(functionDecl.Name, ast.FunctionDeclaration)
	}
	functionDecl.Body = body

	return &functionDecl
}

func (p *Parser) enumDeclaration() ast.Node {
	enumStmt := &ast.EnumDecl{
		Token: p.peek(-1),
		IsPub: p.context.isPub,
	}

	name, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "in enum declaration"), tokens.Identifier)
	if !ok {
		return ast.NewImproper(enumStmt.Token, ast.EnumDeclaration)
	}
	enumStmt.Name = name

	start, ok := p.alertSingleConsume(&alerts.ExpectedSymbol{}, tokens.LeftBrace)
	if !ok {
		return ast.NewImproper(enumStmt.Token, ast.EnumDeclaration)
	}

	fields, _ := p.expressions("in enum declaration", true)
	for _, v := range fields {
		if v.GetType() == ast.Identifier {
			enumStmt.Fields = append(enumStmt.Fields, v.(*ast.IdentifierExpr))
		} else {
			p.AlertSingle(&alerts.InvalidEnumVariantName{}, v.GetToken())
			p.sync(tokens.RightBrace)
			break
		}
	}

	_, ok = p.alertMultiConsume(&alerts.ExpectedSymbol{}, start, p.peek(), tokens.RightBrace)
	if !ok {
		if p.sync(tokens.RightBrace) {
			p.advance()
		}
		return ast.NewImproper(enumStmt.Token, ast.EnumDeclaration)
	}

	return enumStmt
}

func (p *Parser) aliasDeclaration() ast.Node {
	aliasDecl := &ast.AliasDecl{
		Token: p.peek(-1),
		IsPub: p.context.isPub,
	}

	name, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "as the name in alias declaration"), tokens.Identifier)
	if ok {
		aliasDecl.Name = name
	} else if !p.check(tokens.Equal) && p.peek(1).Type == tokens.Equal {
		p.advance()
	} else {
		return ast.NewImproper(aliasDecl.Token, ast.AliasDeclaration)
	}

	_, ok = p.alertSingleConsume(&alerts.ExpectedSymbol{}, tokens.Equal, "after name in alias declaration")
	if !ok {
		return ast.NewImproper(aliasDecl.Token, ast.AliasDeclaration)
	}

	typeExpr, ok := p.checkType("in alias declaration")
	if !ok {
		p.AlertSingle(&alerts.ExpectedType{}, typeExpr.GetToken(), "in alias declaration")
	}
	if typeExpr.GetType() == ast.NA {
		return ast.NewImproper(aliasDecl.Token, ast.AliasDeclaration)
	}
	aliasDecl.Type = typeExpr

	return aliasDecl
}

func (p *Parser) classDeclaration() ast.Node {
	stmt := &ast.ClassDecl{
		IsPub:         p.context.isPub,
		Token:         p.peek(-1),
		GenericParams: make([]*ast.IdentifierExpr, 0),
	}
	if p.context.isPub {
		stmt.Token = p.peek(-2)
	}

	name, _ := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "as the name of the class"), tokens.Identifier)
	stmt.Name = name

	if p.tryGenericArgs() {
		generics, ok := p.genericParams()
		if !ok && !p.sync(tokens.LeftBrace) {
			return ast.NewImproper(stmt.Token, ast.ClassDeclaration)
		} else {
			stmt.GenericParams = generics
		}
	}

	_, ok := p.alertSingleConsume(&alerts.ExpectedSymbol{}, tokens.LeftBrace)
	if !ok {
		return ast.NewImproper(stmt.Token, ast.ClassDeclaration)
	}
	start := p.peek(-1)
	stmt.Methods = []ast.MethodDecl{}
	for p.consumeTill("in class declaration", start, tokens.RightBrace) {
		auxiliaryDeclaration := p.auxiliaryNode()
		switch declaration := auxiliaryDeclaration.(type) {
		case *ast.ConstructorDecl:
			if stmt.Constructor != nil {
				p.AlertMulti(&alerts.MoreThanOneConstructor{}, declaration.GetToken(), p.peek(-1))
			} else {
				stmt.Constructor = declaration
			}
		case *ast.VariableDecl:
			stmt.Fields = append(stmt.Fields, *declaration)
		case *ast.MethodDecl:
			stmt.Methods = append(stmt.Methods, *declaration)
		default:
			p.AlertMulti(&alerts.UnknownStatement{}, auxiliaryDeclaration.GetToken(), p.peek(-1), "in class declaration")
		}
	}

	return stmt
}

func (p *Parser) entityDeclaration() ast.Node {
	stmt := &ast.EntityDecl{
		IsPub:         p.context.isPub,
		Token:         p.peek(-1),
		GenericParams: make([]*ast.IdentifierExpr, 0),
	}
	if p.context.isPub {
		stmt.Token = p.peek(-2)
	}

	name := p.advance()
	stmt.Name = name

	if p.tryGenericParams() {
		generics, ok := p.genericParams()
		if !ok && !p.sync(tokens.LeftBrace) {
			return ast.NewImproper(stmt.Token, ast.EntityDeclaration)
		} else {
			stmt.GenericParams = generics
		}
	}

	if !p.match(tokens.LeftBrace) {
		p.disadvance(2)
		return ast.NewImproper(stmt.Token, ast.NA)
	}
	if name.Type != tokens.Identifier {
		p.AlertSingle(&alerts.ExpectedIdentifier{}, p.peek(), "as the name of the entity")
	}
	start := p.peek(-1)
	for p.consumeTill("in entity declaration", start, tokens.RightBrace) {
		auxiliaryDeclaration := p.auxiliaryNode()
		if auxiliaryDeclaration.GetType() == ast.VariableDeclaration {
			stmt.Fields = append(stmt.Fields, *auxiliaryDeclaration.(*ast.VariableDecl))
			continue
		}
		if auxiliaryDeclaration.GetType() == ast.MethodDeclaration {
			stmt.Methods = append(stmt.Methods, *auxiliaryDeclaration.(*ast.MethodDecl))
			continue
		}
		if auxiliaryDeclaration.GetType() != ast.EntityFunctionDeclaration {
			p.AlertMulti(&alerts.UnknownStatement{}, auxiliaryDeclaration.GetToken(), p.peek(-1), "in entity declaration")
			continue
		}

		funcDecl := auxiliaryDeclaration.(*ast.EntityFunctionDecl)
		funcType := funcDecl.Type

		switch funcDecl.Type {
		case ast.Spawn:
			if stmt.Spawner != nil {
				p.AlertMulti(&alerts.MoreThanOneEntityFunction{}, funcDecl.GetToken(), p.peek(-1), string(funcType))
			} else {
				stmt.Spawner = funcDecl
			}
		case ast.Destroy:
			if stmt.Destroyer != nil {
				p.AlertMulti(&alerts.MoreThanOneEntityFunction{}, funcDecl.GetToken(), p.peek(-1), string(funcType))
			} else {
				stmt.Destroyer = funcDecl
			}
		default:
			var wasFound bool
			for i := range stmt.Callbacks {
				if stmt.Callbacks[i].Type == funcType {
					p.AlertMulti(&alerts.MoreThanOneEntityFunction{}, funcDecl.GetToken(), p.peek(-1), string(funcType))
					wasFound = true
				}
			}
			if !wasFound {
				stmt.Callbacks = append(stmt.Callbacks, funcDecl)
			}
		}

	}

	return stmt
}

func (p *Parser) entityFunctionDeclaration(token tokens.Token, functionType ast.EntityFunctionType) ast.Node {
	stmt := &ast.EntityFunctionDecl{
		Type:  functionType,
		Token: token,
	}

	generics, ok := p.genericParams()
	if !ok {
		return ast.NewImproper(stmt.Token, ast.EntityFunctionDeclaration)
	}
	stmt.Generics = generics
	params, ok := p.functionParams(tokens.LeftParen, tokens.RightParen)
	if !ok {
		return ast.NewImproper(stmt.Token, ast.EntityFunctionDeclaration)
	}
	stmt.Params = params
	returns, ok := p.functionReturns()
	if !ok {
		return ast.NewImproper(stmt.Token, ast.EntityFunctionDeclaration)
	}
	stmt.Returns = returns

	var success bool
	stmt.Body, success = p.body(true, true)
	if !success {
		return ast.NewImproper(stmt.Token, ast.EntityFunctionDeclaration)
	}

	return stmt
}

func (p *Parser) constructorDeclaration() ast.Node {
	stmt := &ast.ConstructorDecl{Token: p.peek(-1)}

	generics, ok := p.genericParams()
	if !ok {
		return ast.NewImproper(stmt.Token, ast.ConstructorDeclaration)
	}
	stmt.Generics = generics
	params, ok := p.functionParams(tokens.LeftParen, tokens.RightParen)
	if !ok {
		return ast.NewImproper(stmt.Token, ast.ConstructorDeclaration)
	}
	stmt.Params = params
	returns, _ := p.functionReturns()
	if returns != nil {
		p.AlertSingle(&alerts.ReturnsInConstructor{}, returns[0].GetToken())
	}
	var success bool
	stmt.Body, success = p.body(false, true)
	if !success {
		return ast.NewImproper(stmt.Token, ast.ConstructorDeclaration)
	}

	return stmt
}

func (p *Parser) fieldDeclaration(matchedLet bool) ast.Node {
	fieldDecl := ast.VariableDecl{
		Token: p.peek(),
	}

	typeCheckStart := p.current
	typeExpr, typeOk := p.checkType("in field declaration")
	if typeOk {
		fieldDecl.Type = typeExpr
		if !p.check(tokens.Identifier) {
			p.disadvance(p.current - typeCheckStart)
			fieldDecl.Type = nil
		}
	} else if !matchedLet {
		return ast.NewImproper(fieldDecl.Token, ast.NA)
	}

	idents, values, ok := p.identExprPairs("in field declaration", fieldDecl.Type != nil)
	if !ok {
		if matchedLet {
			return ast.NewImproper(fieldDecl.Token, ast.VariableDeclaration)
		} else {
			return ast.NewImproper(fieldDecl.Token, ast.NA)
		}
	}

	fieldDecl.Identifiers = idents
	fieldDecl.Expressions = values

	return &fieldDecl
}
