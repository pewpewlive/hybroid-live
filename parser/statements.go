package parser

import (
	"fmt"
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/tokens"
	"os"
	"runtime/debug"
)

func (p *Parser) statement() (returnNode ast.Node) {
	returnNode = &ast.Improper{Token: p.peek()}

	defer func() {
		if errMsg := recover(); errMsg != nil {
			// If the error is a parseError, synchronize to
			// the next statement. If not, propagate the panic.
			if _, ok := errMsg.(ast.Error); ok {
				p.synchronize()
			} else if _, ok := errMsg.(ParserError); ok {
				p.synchronize()
			} else {
				fmt.Printf("panic: %s\nstacktrace:\n", errMsg)
				debug.PrintStack()
				os.Exit(1)
			}
		}
	}()

	varDecl := p.variableDeclarationStmt()
	if varDecl != nil {
		returnNode = varDecl
		return
	}

	token := p.peek().Type
	next := p.peek(1).Type

	if token == tokens.Pub {
		switch next {
		case tokens.Alias:
			p.advance(2)
			returnNode = p.aliasDeclarationStmt(false)
			return
		case tokens.Fn:
			p.advance(2)
			returnNode = p.functionDeclarationStmt(false)
			return
		case tokens.Class:
			p.advance(2)
			returnNode = p.classDeclarationStmt(false)
			return
		case tokens.Entity:
			p.advance(2)
			returnNode = p.entityDeclarationStmt(false)
			return
		case tokens.Enum:
			p.advance(2)
			returnNode = p.enumDeclarationStmt(false)
			return
			// case tokens.Type:
			// 	p.advance(2)
			// 	node = p.TypeDeclarationStmt(false)
		}
	}

	if token == tokens.Struct && next != tokens.Identifier {
		returnNode = p.expression()
		return
	}

	switch token {
	// case tokens.Type:
	// 	p.advance()
	// 	node = p.TypeDeclarationStmt(true)
	case tokens.Alias:
		p.advance()
		returnNode = p.aliasDeclarationStmt(true)
		return
	// case tokens.Macro:
	// 	p.advance()
	// 	returnNode = p.macroDeclarationStmt()
	// 	return
	case tokens.Env:
		p.advance()
		returnNode = p.envStmt()
		return
	case tokens.Fn:
		p.advance()
		returnNode = p.functionDeclarationStmt(true)
		return
	case tokens.Return:
		p.advance()
		returnNode = p.returnStmt()
		return
	case tokens.Yield:
		p.advance()
		returnNode = p.yieldStmt()
		return
	case tokens.Break:
		p.advance()
		returnNode = &ast.BreakStmt{Token: p.peek(-1)}
		return
	case tokens.Destroy:
		p.advance()
		returnNode = p.destroyStmt()
		return
	case tokens.Continue:
		p.advance()
		returnNode = &ast.ContinueStmt{Token: p.peek(-1)}
		return
	case tokens.If:
		p.advance()
		returnNode = p.ifStmt(false, false, false)
		return
	case tokens.Repeat:
		p.advance()
		returnNode = p.repeatStmt()
		return
	case tokens.For:
		p.advance()
		returnNode = p.forStmt()
		return
	case tokens.Tick:
		p.advance()
		returnNode = p.tickStmt()
		return
	case tokens.Use:
		p.advance()
		returnNode = p.useStmt()
		return
	case tokens.Enum:
		p.advance()
		returnNode = p.enumDeclarationStmt(true)
		return
	case tokens.Class:
		p.advance()
		returnNode = p.classDeclarationStmt(true)
		return
	case tokens.Entity:
		p.advance()
		returnNode = p.entityDeclarationStmt(true)
		return
	case tokens.While:
		p.advance()
		returnNode = p.whileStmt()
		return
	case tokens.Match:
		p.advance()
		returnNode = p.matchStmt(false)
		return
	}

	if p.peek().Type == tokens.SemiColon {
		returnNode = ast.NewImproper(p.advance())
		return
	}

	expr := p.expressionStatement()

	if expr.GetType() == ast.NA {
		p.Alert(&alerts.ExpectedStatement{}, alerts.NewSingle(expr.GetToken()))
	}

	returnNode = expr
	return
}

func (p *Parser) expressionStatement() ast.Node {
	expr := p.expression()

	if expr.GetType() == ast.Identifier || expr.GetType() == ast.EnvironmentAccessExpression {
		return p.assignmentStmt(expr)
	}

	typ := expr.GetType()
	if typ != ast.CallExpression && typ != ast.MethodCallExpression {
		return &ast.Improper{Token: expr.GetToken()} // the error is not shown correctly
	}

	return expr
}

func (p *Parser) aliasDeclarationStmt(isLocal bool) ast.Node {
	typeToken := p.peek(-1)
	name, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "in alias declaration"), tokens.Identifier)
	if !ok {
		return ast.NewImproper(name)
	}
	if token, ok := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.Equal, "after identifier in alias declaration"), tokens.Equal); !ok {
		return ast.NewImproper(token)
	}

	if !p.CheckType() {
		return ast.NewImproper(p.peek())
	}
	aliased := p.Type()

	return &ast.AliasDeclarationStmt{
		Alias:       name,
		AliasedType: aliased,
		Token:       typeToken,
		IsLocal:     isLocal,
	}
}

func (p *Parser) envStmt() ast.Node {
	stmt := ast.EnvironmentStmt{}

	if p.Context.EnvStatement != nil {
		p.Alert(&alerts.EnvironmentRedaclaration{}, alerts.NewSingle(p.peek()))
	}

	expr := p.EnvPathExpr()
	if expr.GetType() != ast.EnvironmentPathExpression {
		p.Alert(&alerts.ExpectedEnvironmentPathExpression{}, alerts.NewSingle(p.peek()))
		return &ast.Improper{Token: expr.GetToken()}
	}

	if _, ok := p.consume(p.NewAlert(&alerts.ExpectedKeyword{}, alerts.NewSingle(p.peek()), tokens.As), tokens.As); !ok {
		return &ast.Improper{Token: expr.GetToken()}
	}

	envTypeExpr := p.EnvType()

	if envTypeExpr.Type == ast.InvalidEnv {
		return &ast.Improper{Token: envTypeExpr.GetToken()}
	}

	envPathExpr, _ := expr.(*ast.EnvPathExpr)
	stmt.EnvType = envTypeExpr
	stmt.Env = envPathExpr
	p.Context.EnvStatement = &stmt

	return &stmt
}

func (p *Parser) enumDeclarationStmt(local bool) ast.Node {
	enumStmt := &ast.EnumDeclarationStmt{
		IsLocal: local,
	}

	ident := p.expression()

	if ident.GetType() != ast.Identifier {
		p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(ident.GetToken()), "in enum declaration")
		return &ast.Improper{Token: ident.GetToken()}
	}

	enumStmt.Name = ident.GetToken()

	start, _ := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.LeftBrace), tokens.LeftBrace)

	if p.match(tokens.RightBrace) {
		enumStmt.Fields = make([]tokens.Token, 0)
		return enumStmt
	}

	expr := p.expression()
	if expr.GetType() != ast.Identifier {
		p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(expr.GetToken()), "in enum declaration")
		return &ast.Improper{Token: expr.GetToken()}
	}
	fields := []tokens.Token{expr.GetToken()}
	for p.match(tokens.Comma) {
		if p.check(tokens.RightBrace) {
			break
		}
		expr = p.expression()
		if expr.GetType() != ast.Identifier {
			p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(expr.GetToken()), "in enum declaration")
			return &ast.Improper{Token: expr.GetToken()}
		}
		fields = append(fields, expr.GetToken())
	}

	enumStmt.Fields = fields

	p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewMulti(start, p.peek()), tokens.RightBrace), tokens.RightBrace)

	return enumStmt
}

func (p *Parser) classDeclarationStmt(isLocal bool) ast.Node {
	stmt := &ast.ClassDeclarationStmt{
		IsLocal: isLocal,
	}
	stmt.Token = p.peek(-1)

	name, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "as the name of the class"), tokens.Identifier)

	if ok {
		stmt.Name = name
	} else {
		return &ast.Improper{Token: stmt.Token}
	}

	_, ok = p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.LeftBrace), tokens.LeftBrace)
	if !ok {
		return &ast.Improper{Token: stmt.Token}
	}
	stmt.Methods = []ast.MethodDeclarationStmt{}
	for !p.match(tokens.RightBrace) {
		if p.match(tokens.Fn) {
			method, ok := p.methodDeclarationStmt(stmt.IsLocal).(*ast.MethodDeclarationStmt)
			if ok {
				stmt.Methods = append(stmt.Methods, *method)
			}
		} else if p.match(tokens.New) {
			construct, ok := p.constructorDeclarationStmt().(*ast.ConstructorStmt)
			if ok {
				stmt.Constructor = construct
			}
		} else {
			field := p.fieldDeclarationStmt()
			if field.GetType() != ast.NA {
				stmt.Fields = append(stmt.Fields, *field.(*ast.FieldDeclarationStmt))
			} else {
				p.Alert(&alerts.UnknownStatementInsideClass{}, alerts.NewMulti(field.GetToken(), p.peek()))
			}
		}
	}

	return stmt
}

func (p *Parser) entityDeclarationStmt(isLocal bool) ast.Node {
	stmt := &ast.EntityDeclarationStmt{
		IsLocal: isLocal,
		Token:   p.peek(-1),
	}

	name, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "as the name of the entity"), tokens.Identifier)

	if !ok {
		return &ast.Improper{Token: stmt.Token}
	}
	stmt.Name = name

	_, ok = p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.LeftBrace), tokens.LeftBrace)
	if !ok {
		return &ast.Improper{Token: stmt.Token}
	}

	for !p.match(tokens.RightBrace) {
		if p.match(tokens.Fn) {
			method := p.methodDeclarationStmt(stmt.IsLocal)
			if method.GetType() != ast.MethodDeclarationStatement {
				stmt.Methods = append(stmt.Methods, *method.(*ast.MethodDeclarationStmt))
			}
			continue
		}
		if p.match(tokens.Spawn) {
			spawner := p.entityFunctionDeclarationStmt(p.peek(-1), ast.Spawn)
			if spawner.GetType() != ast.NA {
				stmt.Spawner = spawner.(*ast.EntityFunctionDeclarationStmt)
			}
			continue
		}

		if p.match(tokens.Destroy) {
			destroyer := p.entityFunctionDeclarationStmt(p.peek(-1), ast.Destroy)
			if destroyer.GetType() != ast.NA {
				stmt.Destroyer = destroyer.(*ast.EntityFunctionDeclarationStmt)
			}
			continue
		}
		if p.check(tokens.Identifier) {
			switch p.peek().Lexeme {
			case "WeaponCollision":
				cb := p.entityFunctionDeclarationStmt(p.advance(), ast.WeaponCollision)
				if cb.GetType() != ast.NA {
					stmt.Callbacks = append(stmt.Callbacks, cb.(*ast.EntityFunctionDeclarationStmt))
				}
				continue
			case "WallCollision":
				cb := p.entityFunctionDeclarationStmt(p.advance(), ast.WallCollision)
				if cb.GetType() != ast.NA {
					stmt.Callbacks = append(stmt.Callbacks, cb.(*ast.EntityFunctionDeclarationStmt))
				}
				continue
			case "PlayerCollision":
				cb := p.entityFunctionDeclarationStmt(p.advance(), ast.PlayerCollision)
				if cb.GetType() != ast.NA {
					stmt.Callbacks = append(stmt.Callbacks, cb.(*ast.EntityFunctionDeclarationStmt))
				}
				continue
			case "Update":
				cb := p.entityFunctionDeclarationStmt(p.advance(), ast.Update)
				if cb.GetType() != ast.NA {
					stmt.Callbacks = append(stmt.Callbacks, cb.(*ast.EntityFunctionDeclarationStmt))
				}
				continue
			}
		}
		field := p.fieldDeclarationStmt()
		if field.GetType() != ast.NA {
			stmt.Fields = append(stmt.Fields, *field.(*ast.FieldDeclarationStmt))
		} else {
			p.Alert(&alerts.UnknownStatementInsideEntity{}, alerts.NewMulti(field.GetToken(), p.peek()))
		}
	}

	return stmt
}

func (p *Parser) entityFunctionDeclarationStmt(token tokens.Token, functionType ast.EntityFunctionType) ast.Node {
	stmt := &ast.EntityFunctionDeclarationStmt{
		Type:  functionType,
		Token: token,
	}

	stmt.Generics = p.genericParameters()
	stmt.Params = p.parameters(tokens.LeftParen, tokens.RightParen)
	stmt.Returns = p.returnings()
	p.Context.FunctionReturns.Push("entityFunctionDeclarationStmt", len(stmt.Returns))

	var success bool
	stmt.Body, success = p.getBody()
	if !success {
		return ast.NewImproper(stmt.Token)
	}

	p.Context.FunctionReturns.Pop("entityFunctionDeclarationStmt")

	return stmt
}

func (p *Parser) destroyStmt() ast.Node {
	stmt := ast.DestroyStmt{
		Token: p.peek(-1),
	}

	expr := p.self()
	exprType := expr.GetType()

	if exprType != ast.Identifier && exprType != ast.EnvironmentAccessExpression && exprType != ast.SelfExpression {
		p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(expr.GetToken()), "or environment access expression")
	}
	stmt.Identifier = expr
	stmt.Generics, _ = p.genericArguments()
	stmt.Args = p.arguments()

	return &stmt
}

func (p *Parser) constructorDeclarationStmt() ast.Node {
	stmt := &ast.ConstructorStmt{Token: p.peek(-1)}

	stmt.Generics = p.genericParameters()
	stmt.Params = p.parameters(tokens.LeftParen, tokens.RightParen)
	stmt.Return = p.returnings()
	var success bool
	stmt.Body, success = p.getBody()
	if !success {
		return &ast.Improper{Token: stmt.Token}
	}

	return stmt
}

func (p *Parser) fieldDeclarationStmt() ast.Node {
	stmt := ast.FieldDeclarationStmt{
		Token: p.peek(),
	}

	typ, ident := p.TypeAndIdentifier()
	if ident.GetType() != ast.Identifier {
		p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(ident.GetToken()), "in field declaration")
		return ast.NewImproper(ident.GetToken())
	}

	idents := []tokens.Token{ident.GetToken()}

	stmt.Type = typ
	for p.match(tokens.Comma) {
		ident := p.advance()
		if ident.Type != tokens.Identifier {
			p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(ident), "in field declaration")
		}

		idents = append(idents, ident)
	}

	stmt.Identifiers = idents

	if !p.match(tokens.Equal) {
		stmt.Values = []ast.Node{}
		return &stmt
	}

	expr := p.expression()
	if expr.GetType() == ast.NA {
		p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(p.peek()))
	}

	exprs := []ast.Node{expr}
	for p.match(tokens.Comma) {
		expr = p.expression()
		if expr.GetType() == ast.NA {
			p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(p.peek()))
		}
		exprs = append(exprs, expr)
	}
	stmt.Values = exprs

	return &stmt
}

func (p *Parser) methodDeclarationStmt(IsLocal bool) ast.Node {
	fnDec := p.functionDeclarationStmt(IsLocal)

	if fnDec.GetType() != ast.FunctionDeclarationStatement {
		return fnDec
	} else {
		FnDec := fnDec.(*ast.FunctionDeclarationStmt)
		return &ast.MethodDeclarationStmt{
			IsLocal:  FnDec.IsLocal,
			Name:     FnDec.Name,
			Return:   FnDec.Return,
			Params:   FnDec.Params,
			Generics: FnDec.GenericParams,
			Body:     FnDec.Body,
		}
	}
}

func (p *Parser) ifStmt(else_exists bool, is_else bool, is_elseif bool) *ast.IfStmt {
	ifStm := ast.IfStmt{
		Token: p.peek(-1),
	}

	var expr ast.Node
	if !is_else {
		expr = p.multiComparison()
	}
	ifStm.BoolExpr = expr
	ifStm.Body, _ = p.getBody()

	if is_else || is_elseif {
		return &ifStm
	}
	for p.match(tokens.Else) {
		if else_exists {
			p.Alert(&alerts.MoreThanOneElseStatement{}, alerts.NewSingle(p.peek(-1)))
		}
		var ifbody *ast.IfStmt
		if p.match(tokens.If) {
			ifbody = p.ifStmt(else_exists, false, true)
			ifStm.Elseifs = append(ifStm.Elseifs, ifbody)
		} else {
			else_exists = true
			ifbody = p.ifStmt(else_exists, true, false)
			ifStm.Else = ifbody
		}
	}

	return &ifStm
}

func (p *Parser) assignmentStmt(expr ast.Node) ast.Node {
	idents := []ast.Node{expr}

	for p.match(tokens.Comma) {
		expr := p.expression()
		idents = append(idents, expr)
	}
	values := []ast.Node{}
	if p.match(tokens.Equal) {
		expr2 := p.expression()
		values = append(values, expr2)
		for p.match(tokens.Comma) {
			expr2 := p.expression()
			values = append(values, expr2)
		}
		expr = &ast.AssignmentStmt{Identifiers: idents, Values: values, Token: p.peek(-1)}
	} else if p.match(tokens.PlusEqual, tokens.MinusEqual, tokens.SlashEqual, tokens.StarEqual, tokens.CaretEqual, tokens.ModuloEqual, tokens.BackSlashEqual) {
		assignOp := p.peek(-1)
		op := p.getOp(assignOp)

		expr2 := p.expression()
		binExpr := p.createBinExpr(idents[len(values)], op, op.Type, op.Lexeme, &ast.GroupExpr{Expr: expr2, ValueType: expr2.GetValueType(), Token: expr2.GetToken()})
		values = append(values, binExpr)
		for p.match(tokens.Comma) {
			expr2 := p.expression()
			binExpr := p.createBinExpr(idents[len(values)], op, op.Type, op.Lexeme, &ast.GroupExpr{Expr: expr2, ValueType: expr2.GetValueType(), Token: expr2.GetToken()})

			values = append(values, binExpr)
		}
		expr = &ast.AssignmentStmt{Identifiers: idents, Values: values, Token: assignOp}
	} else {
		p.Alert(&alerts.ExpectedAssignmentSymbol{}, alerts.NewSingle(p.peek()))
	}

	return expr
}

func (p *Parser) returnStmt() ast.Node {
	returnStmt := &ast.ReturnStmt{
		Token: p.peek(-1),
		Args:  []ast.Node{},
	}

	if p.Context.FunctionReturns.Count() == 0 {
		return returnStmt
	}

	if p.Context.FunctionReturns.Top().Item != 0 {
		args, _ := p.returnArgs()
		returnStmt.Args = args
	}

	return returnStmt
}

func (p *Parser) returnArgs() ([]ast.Node, bool) {
	args := []ast.Node{}
	expr := p.expression()
	if expr.GetType() == ast.NA {
		return args, false
	}
	args = append(args, expr)
	for p.match(tokens.Comma) {
		expr = p.expression()
		if expr.GetType() == ast.NA {
			p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(p.peek()))
		}
		args = append(args, expr)
	}
	return args, true
}

func (p *Parser) yieldStmt() ast.Node {
	yieldStmt := ast.YieldStmt{
		Token: p.peek(-1),
	}

	if p.peek().Type == tokens.RightBrace {
		return &yieldStmt
	}
	args := []ast.Node{}
	expr := p.expression()
	args = append(args, expr)
	for p.match(tokens.Comma) {
		expr = p.expression()
		if expr.GetType() == ast.NA {
			p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(p.peek()))
		}
		args = append(args, expr)
	}
	yieldStmt.Args = args

	return &yieldStmt
}

func (p *Parser) functionDeclarationStmt(IsLocal bool) ast.Node {
	fnDec := ast.FunctionDeclarationStmt{}

	fnDec.IsLocal = IsLocal

	ident, _ := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "for a function name"), tokens.Identifier)

	fnDec.Name = ident
	fnDec.GenericParams = p.genericParameters()
	fnDec.Params = p.parameters(tokens.LeftParen, tokens.RightParen)

	fnDec.Return = p.returnings()
	p.Context.FunctionReturns.Push("functionDeclarationStmt", len(fnDec.Return))

	var success bool
	fnDec.Body, success = p.getBody()
	if !success {
		return ast.NewImproper(fnDec.Name)
	}

	p.Context.FunctionReturns.Pop("functionDeclarationStmt")

	return &fnDec
}

func (p *Parser) repeatStmt() ast.Node {
	repeatStmt := ast.RepeatStmt{
		Token: p.peek(-1),
	}

	gotIterator := false
	if p.check(tokens.Number) ||
		p.check(tokens.Fixed) ||
		p.check(tokens.FixedPoint) ||
		p.check(tokens.Radian) ||
		p.check(tokens.Degree) ||
		p.check(tokens.Identifier) {

		repeatStmt.Iterator = p.expression()
		gotIterator = true
	}

	repeatStmt.Skip = &ast.Improper{Token: repeatStmt.Token}
	repeatStmt.Start = &ast.Improper{Token: repeatStmt.Token}

	variableAssigned := false
	iteratorAssgined := false
	skipAssigned := false
	startAssigned := false

	for i := 0; i < 4; i++ {
		if p.match(tokens.With) {
			identExpr := p.expression()
			if variableAssigned {
				p.Alert(&alerts.DuplicateKeywordInRepeatStatement{}, alerts.NewSingle(p.peek(-1)), "with")
			}
			variableAssigned = true
			if identExpr.GetType() != ast.Identifier {
				p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(identExpr.GetToken()), "after keyword 'with'")
			} else {
				repeatStmt.Variable = identExpr.(*ast.IdentifierExpr)
			}
		} else if p.match(tokens.To) {
			if iteratorAssgined {
				p.Alert(&alerts.DuplicateKeywordInRepeatStatement{}, alerts.NewSingle(p.peek(-1)), "to")
			}
			iteratorAssgined = true
			if gotIterator {
				p.Alert(&alerts.RedefinitionOfIteratorInRepeatStatement{}, alerts.NewSingle(p.peek(-1)))
			} else {
				repeatStmt.Iterator = p.expression()
				if repeatStmt.Iterator.GetType() == ast.NA {
					p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(repeatStmt.Iterator.GetToken()))
				}
			}
		} else if p.match(tokens.By) {
			if skipAssigned {
				p.Alert(&alerts.DuplicateKeywordInRepeatStatement{}, alerts.NewSingle(p.peek(-1)), "by")
			}
			skipAssigned = true
			repeatStmt.Skip = p.expression()
			if repeatStmt.Skip.GetType() == ast.NA {
				p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(repeatStmt.Skip.GetToken()))
			}
		} else if p.match(tokens.From) {
			if startAssigned {
				p.Alert(&alerts.DuplicateKeywordInRepeatStatement{}, alerts.NewSingle(p.peek(-1)), "from")
			}
			startAssigned = true
			repeatStmt.Start = p.expression()
			if repeatStmt.Start.GetType() == ast.NA {
				p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(repeatStmt.Start.GetToken()))
			}
		}
	}

	if repeatStmt.Iterator == nil {
		p.Alert(&alerts.MissingIteratorInRepeatStatement{}, alerts.NewSingle(repeatStmt.Token))
		repeatStmt.Iterator = &ast.LiteralExpr{Token: repeatStmt.Token, Value: "1", ValueType: ast.Number}
	}

	var success bool
	repeatStmt.Body, success = p.getBody()
	if !success {
		return ast.NewImproper(repeatStmt.Token)
	}

	return &repeatStmt
}

func (p *Parser) whileStmt() ast.Node {
	whileStmt := &ast.WhileStmt{}

	condtion := p.expression()

	if condtion.GetType() == ast.NA {
		p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(condtion.GetToken()))
		return ast.NewImproper(condtion.GetToken())
	}

	whileStmt.Condtion = condtion

	var success bool
	whileStmt.Body, success = p.getBody()
	if !success {
		return ast.NewImproper(whileStmt.Token)
	}

	return whileStmt
}

func (p *Parser) forStmt() ast.Node {
	forStmt := ast.ForStmt{
		Token: p.peek(-1),
	}

	if p.peek().Type == tokens.Identifier &&
		p.peek(1).Type == tokens.Comma {
		identExpr := p.expression()
		if identExpr.GetType() != ast.Identifier {
			p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(identExpr.GetToken()), "after keyword 'for' in for loop statement")
		} else {
			forStmt.First = identExpr.(*ast.IdentifierExpr)
		}
		p.match(tokens.Comma)
		identExpr = p.expression()
		if identExpr.GetType() != ast.Identifier {
			p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(identExpr.GetToken()))
		} else {
			forStmt.Second = identExpr.(*ast.IdentifierExpr)
		}
	} else {
		identExpr := p.expression()
		if identExpr.GetType() != ast.Identifier {
			p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(identExpr.GetToken()), "after keyword 'for' in for loop statement")
		} else {
			forStmt.First = identExpr.(*ast.IdentifierExpr)
		}
	}

	p.consume(p.NewAlert(&alerts.ExpectedKeyword{}, alerts.NewSingle(p.peek()), tokens.In), tokens.In)

	forStmt.Iterator = p.expression()

	if forStmt.Iterator == nil {
		p.Alert(&alerts.NoIteratorProvidedInForLoopStatement{}, alerts.NewSingle(forStmt.Token))
		forStmt.Iterator = &ast.LiteralExpr{Token: forStmt.Token, Value: "[1]", ValueType: ast.List}
	}

	var success bool
	forStmt.Body, success = p.getBody()
	if !success {
		return ast.NewImproper(forStmt.Token)
	}

	return &forStmt
}

func (p *Parser) tickStmt() ast.Node {
	tickStmt := ast.TickStmt{
		Token: p.peek(-1),
	}

	if p.match(tokens.With) {
		identExpr := p.expression()
		if identExpr.GetType() != ast.Identifier {
			p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(identExpr.GetToken()), "after keyword 'with'")
			return &tickStmt
		} else {
			tickStmt.Variable = identExpr.(*ast.IdentifierExpr)
		}
	}

	var success bool
	tickStmt.Body, success = p.getBody()
	if !success {
		return ast.NewImproper(tickStmt.Token)
	}

	return &tickStmt
}

func (p *Parser) variableDeclarationStmt() ast.Node {
	variable := ast.VariableDeclarationStmt{
		Token:   p.peek(),
		IsLocal: true,
		IsConst: false,
	}

	var typ *ast.TypeExpr
	var ide ast.Node
	nextToken := p.peek().Type

	if nextToken == tokens.Const {
		variable.Token = p.advance()
		variable.IsLocal = false
		variable.IsConst = true

		typ, ide = p.TypeAndIdentifier()

	} else if nextToken == tokens.Let {
		variable.Token = p.advance()
		variable.IsLocal = true
		variable.IsConst = false

		typ, ide = p.TypeAndIdentifier()

	} else if nextToken == tokens.Pub {
		currentStart := p.current

		variable.Token = p.advance()
		variable.IsLocal = false
		variable.IsConst = false

		typ, ide = p.TypeAndIdentifier()

		nextToken = p.peek().Type

		if nextToken != tokens.Equal {
			p.disadvance(p.current - currentStart)

			return nil
		}
	} else if !p.CheckType() {
		return nil
	} else {
		currentStart := p.current

		typ = p.Type()
		token := p.advance()
		if token.Type != tokens.Identifier {
			p.disadvance(p.current - currentStart)

			return nil
		}
		ide = &ast.IdentifierExpr{Name: token, ValueType: ast.Invalid}

		nextToken = p.peek().Type

		if nextToken != tokens.Equal {
			p.disadvance(p.current - currentStart)

			return nil
		}
	}

	if ide.GetType() != ast.Identifier {
		ideToken := ide.GetToken()
		p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(ideToken), "in variable declaration")
		return ast.NewImproper(ideToken)
	}

	idents := []tokens.Token{ide.GetToken()}
	variable.Type = typ
	for p.match(tokens.Comma) {
		ide := p.advance()
		if ide.Type != tokens.Identifier {
			p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(ide), "in variable declaration")
			return ast.NewImproper(ide)
		}

		idents = append(idents, ide)
	}

	variable.Identifiers = idents

	if !p.match(tokens.Equal) {
		variable.Values = []ast.Node{}
		return &variable
	}

	expr := p.expression()
	println(string(expr.GetType()))
	if expr.GetType() == ast.NA || expr.GetType() == ast.TypeExpression {
		p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(p.peek()))
	}

	exprs := []ast.Node{expr}
	for p.match(tokens.Comma) {
		expr = p.expression()
		if expr.GetType() == ast.NA || expr.GetType() == ast.TypeExpression {
			p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(p.peek()))
		}
		exprs = append(exprs, expr)
	}
	variable.Values = exprs

	return &variable
}

func (p *Parser) useStmt() ast.Node {
	useStmt := ast.UseStmt{}

	filepath := p.EnvPathExpr()
	if filepath.GetType() != ast.EnvironmentPathExpression {
		p.Alert(&alerts.ExpectedEnvironmentPathExpression{}, alerts.NewMulti(filepath.GetToken(), p.peek()))
		return ast.NewImproper(p.peek())
	}
	useStmt.Path = filepath.(*ast.EnvPathExpr)

	return &useStmt
}

func (p *Parser) matchStmt(isExpr bool) *ast.MatchStmt {
	matchStmt := ast.MatchStmt{}

	matchStmt.ExprToMatch = p.expression()

	start, _ := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.LeftBrace), tokens.LeftBrace)

	caseStmts, stop := p.caseStmt(isExpr)
	for !stop {
		matchStmt.Cases = append(matchStmt.Cases, caseStmts...)
		caseStmts, stop = p.caseStmt(isExpr)
		for i := range caseStmts {
			if caseStmts[i].Expression.GetToken().Lexeme == "else" {
				matchStmt.HasDefault = true
			}
		}
	}
	matchStmt.Cases = append(matchStmt.Cases, caseStmts...)

	p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewMulti(start, p.peek()), tokens.RightBrace), tokens.RightBrace)

	return &matchStmt
}

func (p *Parser) caseStmt(isExpr bool) ([]ast.CaseStmt, bool) {
	caseStmts := []ast.CaseStmt{}

	caseStmt := ast.CaseStmt{}
	if p.match(tokens.Else) {
		caseStmt.Expression = &ast.IdentifierExpr{
			Name:      p.peek(-1),
			ValueType: ast.Object,
		}
	} else {
		caseStmt.Expression = p.expression()
	}
	if caseStmt.Expression.GetType() == ast.NA {
		return caseStmts, true
	}
	caseStmts = append(caseStmts, caseStmt)
	for p.match(tokens.Comma) {
		caseStmt.Expression = p.expression()
		caseStmts = append(caseStmts, caseStmt)
		if caseStmt.Expression.GetType() == ast.NA {
			return caseStmts, true
		}
	}

	p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.FatArrow), tokens.FatArrow)

	if p.check(tokens.LeftBrace) {
		body, _ := p.getBody()
		for i := range caseStmts {
			caseStmts[i].Body = body
		}
	} else {
		expr := p.expression()
		if expr.GetType() == ast.NA {
			p.Alert(&alerts.ExpectedExpressionOrBody{}, alerts.NewSingle(p.peek()))
		}
		args := []ast.Node{expr}
		for p.match(tokens.Comma) {
			expr = p.expression()
			if expr.GetType() == ast.NA {
				p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(p.peek()))
			}
			args = append(args, expr)
		}

		var node ast.Node
		if isExpr {
			node = &ast.YieldStmt{
				Args:  args,
				Token: expr.GetToken(),
			}
		} else {
			node = &ast.ReturnStmt{
				Args:  args,
				Token: expr.GetToken(),
			}
		}

		body := []ast.Node{node}
		for i := range caseStmts {
			caseStmts[i].Body = body
		}
	}
	if p.check(tokens.RightBrace) {
		return caseStmts, true
	}

	return caseStmts, false
}
