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

	if _, ok := p.consume(p.NewAlert(&alerts.ExpectedKeyword{}, alerts.NewSingle(p.peek()), tokens.As), tokens.As); !ok {
		return ast.NewImproper(p.peek(), ast.EnvironmentDeclaration)
	}

	envPathExpr, _ := pathExpr.(*ast.EnvPathExpr)

	envDecl.EnvType = p.envTypeExpr()
	envDecl.Env = envPathExpr

	return envDecl
}

func (p *Parser) variableDeclaration(matchedLetOrConst bool) ast.Node {
	variableDecl := ast.VariableDecl{
		IsPub:   p.context.isPub,
		IsConst: false,
	}

	// check if variableDecl started with let or const
	if matchedLetOrConst {
		variableDecl.Token = p.peek(-1)
		variableDecl.IsConst = variableDecl.Token.Type == tokens.Const
		if variableDecl.IsPub {
			variableDecl.Token = p.peek(-2)
		}
		if variableDecl.IsPub && variableDecl.Token.Type == tokens.Let {
			p.Alert(&alerts.UnexpectedKeyword{}, alerts.NewSingle(variableDecl.Token), variableDecl.Token.Lexeme, "in variable declaration")
		}
	} else {
		if variableDecl.IsPub {
			variableDecl.Token = p.peek(-1)
		} else {
			variableDecl.Token = p.peek()
		}
	}

	typeCheckStart := p.current
	typeExpr, ok := p.checkType("in variable declaration")
	if ok {
		variableDecl.Type = typeExpr
		if !p.check(tokens.Identifier) {
			p.disadvance(p.current - typeCheckStart)
			variableDecl.Type = nil
		}
	} else if !matchedLetOrConst {
		return ast.NewImproper(variableDecl.Token, ast.NA)
	}

	idents, exprs, ok := p.identExprPairs("in variable declaration", variableDecl.Type != nil)
	if !ok {
		if matchedLetOrConst {
			return ast.NewImproper(variableDecl.Token, ast.VariableDeclaration)
		}
		if p.peek(-1).Type == tokens.Equal {
			return ast.NewImproper(variableDecl.Token, ast.VariableDeclaration)
		}
		return ast.NewImproper(variableDecl.Token, ast.NA)
	}

	variableDecl.Identifiers = idents
	variableDecl.Expressions = exprs

	return &variableDecl
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
	if !nameOk && !p.check(tokens.Less) && !p.check(tokens.LeftParen) {
		p.advance()
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
	functionDecl.Return = p.functionReturns()

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

	start, ok := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.LeftBrace), tokens.LeftBrace)
	if !ok {
		return ast.NewImproper(enumStmt.Token, ast.EnumDeclaration)
	}

	p.context.braceCounter.Increment()
	defer p.context.braceCounter.Decrement()

	fields, _ := p.expressions("in enum declaration", true)
	for _, v := range fields {
		if v.GetType() == ast.Identifier {
			enumStmt.Fields = append(enumStmt.Fields, v.(*ast.IdentifierExpr))
		} else {
			p.Alert(&alerts.InvalidEnumVariantName{}, alerts.NewSingle(v.GetToken()))
			return ast.NewImproper(enumStmt.Token, ast.EnumDeclaration)
		}
	}

	_, ok = p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewMulti(start, p.peek()), tokens.RightBrace), tokens.RightBrace)
	if !ok {
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

	_, ok = p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.Equal, "after name in alias declaration"), tokens.Equal)
	if !ok {
		return ast.NewImproper(aliasDecl.Token, ast.AliasDeclaration)
	}

	typ, ok := p.checkType("in alias declaration")
	if !ok {
		p.Alert(&alerts.ExpectedType{}, alerts.NewSingle(aliasDecl.Type.GetToken()), "in alias declaration")
	}
	if typ.GetType() == ast.NA {
		return ast.NewImproper(aliasDecl.Token, ast.AliasDeclaration)
	}
	aliasDecl.Type = typ

	return aliasDecl
}

func (p *Parser) classDeclaration() ast.Node {
	stmt := &ast.ClassDecl{
		IsPub: p.context.isPub,
		Token: p.peek(-1),
	}
	if p.context.isPub {
		stmt.Token = p.peek(-2)
	}

	name, _ := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "as the name of the class"), tokens.Identifier)
	stmt.Name = name

	_, ok := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.LeftBrace), tokens.LeftBrace)
	if !ok {
		return ast.NewImproper(stmt.Token, ast.ClassDeclaration)
	}
	p.context.braceCounter.Increment()
	defer p.context.braceCounter.Decrement()

	start := p.peek(-1)
	stmt.Methods = []ast.MethodDecl{}
	for p.consumeTill("in class declaration", start, tokens.RightBrace) {
		auxiliaryDeclaration := p.auxiliaryNode()
		switch declaration := auxiliaryDeclaration.(type) {
		case *ast.ConstructorDecl:
			if stmt.Constructor != nil {
				p.Alert(&alerts.MoreThanOneConstructor{}, alerts.NewMulti(declaration.GetToken(), p.peek(-1)))
			} else {
				stmt.Constructor = declaration
			}
		case *ast.FieldDecl:
			stmt.Fields = append(stmt.Fields, *declaration)
		case *ast.MethodDecl:
			stmt.Methods = append(stmt.Methods, *declaration)
		default:
			p.Alert(&alerts.UnknownStatement{}, alerts.NewMulti(auxiliaryDeclaration.GetToken(), p.peek(-1)), "in class declaration")
		}
	}

	return stmt
}

func (p *Parser) entityDeclaration() ast.Node {
	stmt := &ast.EntityDecl{
		IsPub: p.context.isPub,
		Token: p.peek(-1),
	}
	if p.context.isPub {
		stmt.Token = p.peek(-2)
	}

	name := p.advance()
	stmt.Name = name

	if !p.match(tokens.LeftBrace) {
		p.disadvance(2)
		return ast.NewImproper(stmt.Token, ast.NA)
	}
	if name.Type != tokens.Identifier {
		p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "as the name of the entity")
	}
	p.context.braceCounter.Increment()
	defer p.context.braceCounter.Decrement()

	start := p.peek(-1)
	for p.consumeTill("in entity declaration", start, tokens.RightBrace) {
		auxiliaryDeclaration := p.auxiliaryNode()
		if auxiliaryDeclaration.GetType() == ast.FieldDeclaration {
			stmt.Fields = append(stmt.Fields, *auxiliaryDeclaration.(*ast.FieldDecl))
			continue
		}
		if auxiliaryDeclaration.GetType() == ast.MethodDeclaration {
			stmt.Methods = append(stmt.Methods, *auxiliaryDeclaration.(*ast.MethodDecl))
			continue
		}
		if auxiliaryDeclaration.GetType() != ast.EntityFunctionDeclaration {
			p.Alert(&alerts.UnknownStatement{}, alerts.NewMulti(auxiliaryDeclaration.GetToken(), p.peek(-1)), "in entity declaration")
			continue
		}

		funcDecl := auxiliaryDeclaration.(*ast.EntityFunctionDecl)
		funcType := funcDecl.Type

		switch funcDecl.Type {
		case ast.Spawn:
			if stmt.Spawner != nil {
				p.Alert(&alerts.MoreThanOneEntityFunction{}, alerts.NewMulti(funcDecl.GetToken(), p.peek(-1)), string(funcType))
			} else {
				stmt.Spawner = funcDecl
			}
		case ast.Destroy:
			if stmt.Destroyer != nil {
				p.Alert(&alerts.MoreThanOneEntityFunction{}, alerts.NewMulti(funcDecl.GetToken(), p.peek(-1)), string(funcType))
			} else {
				stmt.Destroyer = funcDecl
			}
		default:
			var wasFound bool
			for i := range stmt.Callbacks {
				if stmt.Callbacks[i].Type == funcType {
					p.Alert(&alerts.MoreThanOneEntityFunction{}, alerts.NewMulti(funcDecl.GetToken(), p.peek(-1)), string(funcType))
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
	stmt.Return = p.functionReturns()

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
	returns := p.functionReturns()
	if returns != nil {
		p.Alert(&alerts.ReturnsInConstructor{}, alerts.NewSingle(returns.GetToken()))
	}
	var success bool
	stmt.Body, success = p.body(false, true)
	if !success {
		return ast.NewImproper(stmt.Token, ast.ConstructorDeclaration)
	}

	return stmt
}

func (p *Parser) fieldDeclaration() ast.Node {
	fieldDecl := ast.FieldDecl{
		Token: p.peek(),
	}

	typeCheckStart := p.current
	if typeExpr, ok := p.checkType("in field declaration"); ok {
		fieldDecl.Type = typeExpr
		if !p.check(tokens.Identifier) {
			p.disadvance(p.current - typeCheckStart)
			fieldDecl.Type = nil
		}
	}

	idents, values, ok := p.identExprPairs("in field declaration", fieldDecl.Type != nil)
	if !ok {
		return ast.NewImproper(fieldDecl.Token, ast.FieldDeclaration)
	}

	fieldDecl.Identifiers = idents
	fieldDecl.Values = values

	return &fieldDecl
}
