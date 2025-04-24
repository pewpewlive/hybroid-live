package parser

import (
	"fmt"
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/tokens"
	"os"
	"runtime/debug"
)

func (p *Parser) OLDstatement() (returnNode ast.Node) {
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

	varDecl := p.variableDeclaration()
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
			returnNode = p.aliasDeclaration()
			return
		case tokens.Fn:
			p.advance(2)
			returnNode = p.functionDeclaration()
			return
		case tokens.Class:
			p.advance(2)
			returnNode = p.classDeclaration()
			return
		case tokens.Entity:
			p.advance(2)
			// returnNode = p.entityDeclaration()
			return
		case tokens.Enum:
			p.advance(2)
			returnNode = p.enumDeclaration()
			return
			// case tokens.Type:
			// 	p.advance(2)
			// 	node = p.TypeDeclaration()
		}
	}

	if token == tokens.Struct && next != tokens.Identifier {
		returnNode = p.expression()
		return
	}

	switch token {
	// case tokens.Type:
	// 	p.advance()
	// 	node = p.TypeDeclaration(true)
	case tokens.Alias:
		p.advance()
		returnNode = p.aliasDeclaration()
		return
	// case tokens.Macro:
	// 	p.advance()
	// 	returnNode = p.macroDeclaration()
	// 	return
	case tokens.Env:
		p.advance()
		returnNode = p.environmentDeclaration()
		return
	case tokens.Fn:
		p.advance()
		returnNode = p.functionDeclaration()
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
		returnNode = p.enumDeclaration()
		return
	case tokens.Class:
		p.advance()
		returnNode = p.classDeclaration()
		return
	case tokens.Entity:
		p.advance()
		// returnNode = p.entityDeclaration()
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
		returnNode = ast.NewImproper(p.advance(), ast.NA)
		return
	}

	expr := p.expressionStatement()

	if expr.GetType() == ast.NA {
		p.Alert(&alerts.ExpectedStatement{}, alerts.NewSingle(expr.GetToken()))
	}

	returnNode = expr
	return
}

func (p *Parser) OLDexpressionStatement() ast.Node {
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

func (p *Parser) OLDaliasDeclarationStmt(isLocal bool) ast.Node {
	typeToken := p.peek(-1)
	name, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "in alias declaration"), tokens.Identifier)
	if !ok {
		return ast.NewImproper(name, ast.NA)
	}
	if token, ok := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.Equal, "after identifier in alias declaration"), tokens.Equal); !ok {
		return ast.NewImproper(token, ast.NA)
	}

	aliased, ok := p.checkType()
	if !ok {
		return ast.NewImproper(p.peek(), ast.NA)
	}

	return &ast.AliasDecl{
		Name:  name,
		Type:  aliased,
		Token: typeToken,
		IsPub: isLocal,
	}
}

func (p *Parser) OLDenvStmt() ast.Node {
	stmt := ast.EnvironmentDecl{}

	if p.Context.EnvDeclaration != nil {
		p.Alert(&alerts.EnvironmentRedaclaration{}, alerts.NewSingle(p.peek()))
	}

	expr := p.envPathExpr()
	if expr.GetType() != ast.EnvironmentPathExpression {
		p.Alert(&alerts.ExpectedEnvironmentPathExpression{}, alerts.NewSingle(p.peek()))
		return &ast.Improper{Token: expr.GetToken()}
	}

	if _, ok := p.consume(p.NewAlert(&alerts.ExpectedKeyword{}, alerts.NewSingle(p.peek()), tokens.As), tokens.As); !ok {
		return &ast.Improper{Token: expr.GetToken()}
	}

	envTypeExpr := p.envTypeExpr()

	if envTypeExpr.Type == ast.InvalidEnv {
		return &ast.Improper{Token: envTypeExpr.GetToken()}
	}

	envPathExpr, _ := expr.(*ast.EnvPathExpr)
	stmt.EnvType = envTypeExpr
	stmt.Env = envPathExpr
	p.Context.EnvDeclaration = &stmt

	return &stmt
}

func (p *Parser) OLDenumDeclarationStmt(local bool) ast.Node {
	enumStmt := &ast.EnumDecl{
		IsPub: local,
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

func (p *Parser) OLDclassDeclarationStmt(isLocal bool) ast.Node {
	stmt := &ast.ClassDecl{
		IsPub: isLocal,
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
	stmt.Methods = []ast.MethodDecl{}
	for !p.match(tokens.RightBrace) {
		// if p.match(tokens.Fn) {
		// 	// method, ok := p.methodDeclaration().(*ast.MethodDecl)
		// 	// if ok {
		// 	// 	stmt.Methods = append(stmt.Methods, *method)
		// 	// }
		// } else if p.match(tokens.New) {
		// 	construct, ok := p.constructorDeclaration().(*ast.ConstructorDecl)
		// 	if ok {
		// 		stmt.Constructor = construct
		// 	}
		// } else {
		// 	field := p.fieldDeclaration()
		// 	if field.GetType() != ast.NA {
		// 		stmt.Fields = append(stmt.Fields, *field.(*ast.FieldDecl))
		// 	} else {
		// 		p.Alert(&alerts.UnknownStatement{}, alerts.NewMulti(field.GetToken(), p.peek()), "in class declaration")
		// 	}
		// }
	}

	return stmt
}

func (p *Parser) OLDentityDeclarationStmt(isLocal bool) ast.Node {
	stmt := &ast.EntityDecl{
		IsPub: isLocal,
		Token: p.peek(-1),
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
			// method := p.methodDeclaration()
			// if method.GetType() != ast.MethodDeclaration {
			// 	stmt.Methods = append(stmt.Methods, *method.(*ast.MethodDecl))
			// }
			continue
		}
		if p.match(tokens.Spawn) {
			// spawner := p.entityFunctionDeclaration(p.peek(-1), ast.Spawn)
			// if spawner.GetType() != ast.NA {
			// 	stmt.Spawner = spawner.(*ast.EntityFunctionDecl)
			// }
			continue
		}

		if p.match(tokens.Destroy) {
			// destroyer := p.entityFunctionDeclaration(p.peek(-1), ast.Destroy)
			// if destroyer.GetType() != ast.NA {
			// 	stmt.Destroyer = destroyer.(*ast.EntityFunctionDecl)
			// }
			continue
		}
		if p.check(tokens.Identifier) {
			switch p.peek().Lexeme {
			case "WeaponCollision":
				// cb := p.entityFunctionDeclaration(p.advance(), ast.WeaponCollision)
				// if cb.GetType() != ast.NA {
				// 	stmt.Callbacks = append(stmt.Callbacks, cb.(*ast.EntityFunctionDecl))
				// }
				continue
			case "WallCollision":
				// cb := p.entityFunctionDeclaration(p.advance(), ast.WallCollision)
				// if cb.GetType() != ast.NA {
				// 	stmt.Callbacks = append(stmt.Callbacks, cb.(*ast.EntityFunctionDecl))
				// }
				continue
			case "PlayerCollision":
				// cb := p.entityFunctionDeclaration(p.advance(), ast.PlayerCollision)
				// if cb.GetType() != ast.NA {
				// 	stmt.Callbacks = append(stmt.Callbacks, cb.(*ast.EntityFunctionDecl))
				// }
				continue
			case "Update":
				// cb := p.entityFunctionDeclaration(p.advance(), ast.Update)
				// if cb.GetType() != ast.NA {
				// 	stmt.Callbacks = append(stmt.Callbacks, cb.(*ast.EntityFunctionDecl))
				// }
				continue
			}
		}
		// field := p.fieldDeclaration()
		// if field.GetType() != ast.NA {
		// 	stmt.Fields = append(stmt.Fields, *field.(*ast.FieldDecl))
		// } else {
		// 	p.Alert(&alerts.UnknownStatement{}, alerts.NewMulti(field.GetToken(), p.peek()), "in entity declaration")
		// }
	}

	return stmt
}

func (p *Parser) OLDentityFunctionDeclarationStmt(token tokens.Token, functionType ast.EntityFunctionType) ast.Node {
	stmt := &ast.EntityFunctionDecl{
		Type:  functionType,
		Token: token,
	}

	stmt.Generics = p.genericParams()
	stmt.Params = p.functionParams(tokens.LeftParen, tokens.RightParen)
	stmt.Return = p.functionReturns()
	//p.Context.FunctionReturns.Push("entityFunctionDeclarationStmt", len(stmt.Return))

	var success bool
	stmt.Body, success = p.body(true, true)
	if !success {
		return ast.NewImproper(stmt.Token, ast.NA)
	}

	p.Context.FunctionReturns.Pop("entityFunctionDeclarationStmt")

	return stmt
}

func (p *Parser) OLDdestroyStmt() ast.Node {
	stmt := ast.DestroyStmt{
		Token: p.peek(-1),
	}

	expr := p.self()
	exprType := expr.GetType()

	if exprType != ast.Identifier && exprType != ast.EnvironmentAccessExpression && exprType != ast.SelfExpression {
		p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(expr.GetToken()), "or environment access expression")
	}
	stmt.Identifier = expr
	stmt.Generics, _ = p.genericArgs()
	stmt.Args = p.functionArgs()

	return &stmt
}

func (p *Parser) OLDconstructorDeclarationStmt() ast.Node {
	stmt := &ast.ConstructorDecl{Token: p.peek(-1)}

	stmt.Generics = p.genericParams()
	stmt.Params = p.functionParams(tokens.LeftParen, tokens.RightParen)
	//stmt.Return = p.functionReturns()
	var success bool
	stmt.Body, success = p.body(true, true)
	if !success {
		return &ast.Improper{Token: stmt.Token}
	}

	return stmt
}

func (p *Parser) OLDfieldDeclarationStmt() ast.Node {
	stmt := ast.FieldDecl{
		Token: p.peek(),
	}

	typ, ident := p.typeAndIdentifier()
	if ident.GetType() != ast.Identifier {
		p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(ident.GetToken()), "in field declaration")
		return ast.NewImproper(ident.GetToken(), ast.NA)
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

func (p *Parser) OLDmethodDeclarationStmt(IsLocal bool) ast.Node {
	fnDec := p.functionDeclaration()

	if fnDec.GetType() != ast.FunctionDeclaration {
		return fnDec
	} else {
		FnDec := fnDec.(*ast.FunctionDecl)
		return &ast.MethodDecl{
			IsPub:    FnDec.IsPub,
			Name:     FnDec.Name,
			Return:   FnDec.Return,
			Params:   FnDec.Params,
			Generics: FnDec.Generics,
			Body:     FnDec.Body,
		}
	}
}

func (p *Parser) OLDifStmt(else_exists bool, is_else bool, is_elseif bool) *ast.IfStmt {
	ifStm := ast.IfStmt{
		Token: p.peek(-1),
	}

	var expr ast.Node
	if !is_else {
		expr = p.multiComparison()
	}
	ifStm.BoolExpr = expr
	ifStm.Body, _ = p.body(true, true)

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

func (p *Parser) OLDassignmentStmt(expr ast.Node) ast.Node {
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
		op := tokens.Token{Literal: assignOp.Literal, Position: assignOp.Position}
		switch assignOp.Type {
		case tokens.PlusEqual:
			op.Type = tokens.Plus
		case tokens.MinusEqual:
			op.Type = tokens.Minus
		case tokens.SlashEqual:
			op.Type = tokens.Slash
		case tokens.StarEqual:
			op.Type = tokens.Star
		case tokens.CaretEqual:
			op.Type = tokens.Caret
		case tokens.ModuloEqual:
			op.Type = tokens.Modulo
		}
		op.Lexeme = string(op.Type)

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

func (p *Parser) OLDreturnStmt() ast.Node {
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

func (p *Parser) OLDreturnArgs() ([]ast.Node, bool) {
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

func (p *Parser) OLDyieldStmt() ast.Node {
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

func (p *Parser) OLDfunctionDeclarationStmt(IsLocal bool) ast.Node {
	fnDec := ast.FunctionDecl{}

	fnDec.IsPub = IsLocal

	ident, _ := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "for a function name"), tokens.Identifier)

	fnDec.Name = ident
	fnDec.Generics = p.genericParams()
	fnDec.Params = p.functionParams(tokens.LeftParen, tokens.RightParen)

	fnDec.Return = p.functionReturns()
	//p.Context.FunctionReturns.Push("functionDeclarationStmt", len(fnDec.Return))

	var success bool
	fnDec.Body, success = p.body(true, true)
	if !success {
		return ast.NewImproper(fnDec.Name, ast.NA)
	}

	p.Context.FunctionReturns.Pop("functionDeclarationStmt")

	return &fnDec
}

func (p *Parser) OLDrepeatStmt() ast.Node {
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
				p.Alert(&alerts.DuplicateKeyword{}, alerts.NewSingle(p.peek(-1)), tokens.With)
			}
			variableAssigned = true
			if identExpr.GetType() != ast.Identifier {
				p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(identExpr.GetToken()), "after keyword 'with'")
			} else {
				repeatStmt.Variable = identExpr.(*ast.IdentifierExpr)
			}
		} else if p.match(tokens.To) {
			if iteratorAssgined {
				p.Alert(&alerts.DuplicateKeyword{}, alerts.NewSingle(p.peek(-1)), tokens.To)
			}
			iteratorAssgined = true
			if gotIterator {
				p.Alert(&alerts.IteratorRedefinition{}, alerts.NewSingle(p.peek(-1)), "in repeat statement")
			} else {
				repeatStmt.Iterator = p.expression()
				if repeatStmt.Iterator.GetType() == ast.NA {
					p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(repeatStmt.Iterator.GetToken()))
				}
			}
		} else if p.match(tokens.By) {
			if skipAssigned {
				p.Alert(&alerts.DuplicateKeyword{}, alerts.NewSingle(p.peek(-1)), tokens.By)
			}
			skipAssigned = true
			repeatStmt.Skip = p.expression()
			if repeatStmt.Skip.GetType() == ast.NA {
				p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(repeatStmt.Skip.GetToken()))
			}
		} else if p.match(tokens.From) {
			if startAssigned {
				p.Alert(&alerts.DuplicateKeyword{}, alerts.NewSingle(p.peek(-1)), tokens.From)
			}
			startAssigned = true
			repeatStmt.Start = p.expression()
			if repeatStmt.Start.GetType() == ast.NA {
				p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(repeatStmt.Start.GetToken()))
			}
		}
	}

	if repeatStmt.Iterator == nil {
		p.Alert(&alerts.MissingIterator{}, alerts.NewSingle(repeatStmt.Token), "in repeat statement")
		repeatStmt.Iterator = &ast.LiteralExpr{Token: repeatStmt.Token, Value: "1", ValueType: ast.Number}
	}

	var success bool
	repeatStmt.Body, success = p.body(true, true)
	if !success {
		return ast.NewImproper(repeatStmt.Token, ast.NA)
	}

	return &repeatStmt
}

func (p *Parser) OLDwhileStmt() ast.Node {
	whileStmt := &ast.WhileStmt{}

	condtion := p.expression()

	if condtion.GetType() == ast.NA {
		p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(condtion.GetToken()))
		return ast.NewImproper(condtion.GetToken(), ast.NA)
	}

	whileStmt.Condtion = condtion

	var success bool
	whileStmt.Body, success = p.body(true, true)
	if !success {
		return ast.NewImproper(whileStmt.Token, ast.NA)
	}

	return whileStmt
}

func (p *Parser) OLDforStmt() ast.Node {
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
		p.Alert(&alerts.MissingIterator{}, alerts.NewSingle(forStmt.Token), "in for statement")
		forStmt.Iterator = &ast.LiteralExpr{Token: forStmt.Token, Value: "[1]", ValueType: ast.List}
	}

	var success bool
	forStmt.Body, success = p.body(true, true)
	if !success {
		return ast.NewImproper(forStmt.Token, ast.NA)
	}

	return &forStmt
}

func (p *Parser) OLDtickStmt() ast.Node {
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
	tickStmt.Body, success = p.body(true, true)
	if !success {
		return ast.NewImproper(tickStmt.Token, ast.NA)
	}

	return &tickStmt
}

func (p *Parser) OLDvariableDeclarationStmt() ast.Node {
	variable := ast.VariableDecl{
		Token:   p.peek(),
		IsPub:   true,
		IsConst: false,
	}

	var typ *ast.TypeExpr
	var ide ast.Node
	nextToken := p.peek().Type

	if nextToken == tokens.Const {
		variable.Token = p.advance()
		variable.IsPub = false
		variable.IsConst = true

		typ, ide = p.typeAndIdentifier()

	} else if nextToken == tokens.Let {
		variable.Token = p.advance()
		variable.IsPub = true
		variable.IsConst = false

		// let a, b, c = 10, 20, 30

		currentStart := p.current
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
	} else if nextToken == tokens.Pub {
		currentStart := p.current

		variable.Token = p.advance()
		variable.IsPub = false
		variable.IsConst = false

		typ, ide = p.typeAndIdentifier()

		nextToken = p.peek().Type

		if nextToken != tokens.Equal {
			p.disadvance(p.current - currentStart)

			return nil
		}
	} else {
		currentStart := p.current

		typ = p.typeExpr()
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
		return ast.NewImproper(ideToken, ast.NA)
	}

	idents := []tokens.Token{ide.GetToken()}
	variable.Type = typ
	for p.match(tokens.Comma) {
		ide := p.advance()
		if ide.Type != tokens.Identifier {
			p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(ide), "in variable declaration")
			return ast.NewImproper(ide, ast.NA)
		}

		idents = append(idents, ide)
	}

	variable.Identifiers = idents

	if !p.match(tokens.Equal) {
		variable.Expressions = []ast.Node{}
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
	variable.Expressions = exprs

	return &variable
}

func (p *Parser) OLDuseStmt() ast.Node {
	useStmt := ast.UseStmt{}

	filepath := p.envPathExpr()
	if filepath.GetType() != ast.EnvironmentPathExpression {
		p.Alert(&alerts.ExpectedEnvironmentPathExpression{}, alerts.NewMulti(filepath.GetToken(), p.peek()))
		return ast.NewImproper(p.peek(), ast.NA)
	}
	useStmt.Path = filepath.(*ast.EnvPathExpr)

	return &useStmt
}

func (p *Parser) OLDmatchStmt(isExpr bool) *ast.MatchStmt {
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

func (p *Parser) OLDcaseStmt(isExpr bool) ([]ast.CaseStmt, bool) {
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
		body, _ := p.body(true, true)
		for i := range caseStmts { // "hello" =>
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
