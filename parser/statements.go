package parser

import (
	"hybroid/ast"
	"hybroid/lexer"
)

func (p *Parser) statement() ast.Node {
	defer func() {
		if errMsg := recover(); errMsg != nil {
			// If the error is a parseError, synchronize to
			// the next statement. If not, propagate the panic.
			if _, ok := errMsg.(ast.Error); ok {
				//p. = true
				p.synchronize()
			} else {
				panic(errMsg)
			}
		}
	}()
	token := p.peek().Type
	next := p.peek(1).Type

	if token == lexer.Pub {
		switch next {
		case lexer.Fn:
			p.advance()
			token = p.peek().Type
		case lexer.Struct:
			p.advance()
			p.advance()
			return p.structDeclarationStmt()
		case lexer.Enum:
			p.advance()
			p.advance()
			return p.enumDeclarationStmt(false)
		}
	}

	if token == lexer.Struct && next != lexer.Identifier {
		return p.expression()
	}

	switch token {
	case lexer.Env:
		p.advance()
		return p.envStmt()
	case lexer.Let, lexer.Pub, lexer.Const:
		p.advance()
		return p.variableDeclarationStmt()
	case lexer.Add:
		p.advance()
		return p.addToStmt()
	case lexer.Remove:
		p.advance()
		return p.removeFromStmt()
	case lexer.Fn:
		p.advance()
		return p.functionDeclarationStmt()
	case lexer.Return:
		p.advance()
		return p.returnStmt()
	case lexer.Yield:
		p.advance()
		return p.yieldStmt()
	case lexer.Break:
		p.advance()
		return &ast.BreakStmt{Token: p.peek(-1)}
	case lexer.Continue:
		p.advance()
		return &ast.ContinueStmt{Token: p.peek(-1)}
	case lexer.Identifier, lexer.Self:
		return p.assignmentStmt()
	case lexer.If:
		p.advance()
		return p.ifStmt(false, false, false)
	case lexer.Repeat:
		p.advance()
		return p.repeatStmt()
	case lexer.For:
		p.advance()
		return p.forStmt()
	case lexer.Tick:
		p.advance()
		return p.tickStmt()
	case lexer.Use:
		p.advance()
		return p.useStmt()
	case lexer.Enum:
		p.advance()
		return p.enumDeclarationStmt(true)
	case lexer.Struct:
		p.advance()
		return p.structDeclarationStmt()
	case lexer.While:
		p.advance()
		return p.whileStmt()
	case lexer.Match:
		p.advance()
		return p.matchStmt(false)
	}

	expr := p.expression()
	if expr.GetType() == ast.NA {
		p.error(p.peek(), "expected statement")
		p.advance()
	}
	return expr
}

func (p *Parser) getBody() []ast.Node {
	body := make([]ast.Node, 0)
	if _, success := p.consume("expected opening of the body", lexer.LeftBrace); !success {
		return body
	}

	for !p.match(lexer.RightBrace) {
		if p.peek().Type == lexer.Eof {
			p.error(p.peek(), "expected body closure")
			break
		}

		statement := p.statement()
		if statement.GetType() != ast.NA {
			body = append(body, statement)
		}
	}

	return body
}

func (p *Parser) envStmt() ast.Node {
	stmt := ast.EnvironmentStmt{}

	expr := p.EnvPathExpr()
	if expr.GetType() != ast.EnvironmentPathExpression {
		p.error(expr.GetToken(), "expected environment path expression")
		return &ast.Improper{Token: expr.GetToken()}
	}

	if _, ok := p.consume("expected keyword 'as' after envrionment expression", lexer.As); !ok {
		return &ast.Improper{Token: expr.GetToken()}
	}

	envTypeExpr := p.EnvType()

	if envTypeExpr.Type == ast.InvalidEnv {
		return &ast.Improper{Token: envTypeExpr.GetToken()}
	}

	envPathExpr, _ := expr.(*ast.EnvPathExpr)
	stmt.EnvType = envTypeExpr
	stmt.Env = envPathExpr

	return &stmt
}

func (p *Parser) enumDeclarationStmt(local bool) ast.Node {
	enumStmt := &ast.EnumDeclarationStmt{
		IsLocal: local,
	}

	ident := p.expression()

	if ident.GetType() != ast.Identifier {
		p.error(ident.GetToken(), "expected identifier after 'enum' in enum declaration")
		return &ast.Improper{Token: ident.GetToken()}
	}

	enumStmt.Name = ident.GetToken()

	p.consume("expected opening of a body", lexer.LeftBrace)

	if p.match(lexer.RightBrace) {
		enumStmt.Fields = make([]lexer.Token, 0)
		return enumStmt
	}

	expr := p.expression()
	if expr.GetType() != ast.Identifier {
		p.error(expr.GetToken(), "expected identifier in enum declaration")
		return &ast.Improper{Token: expr.GetToken()}
	}
	fields := []lexer.Token{expr.GetToken()}
	for p.match(lexer.Comma) {
		if p.check(lexer.RightBrace) {
			break
		}
		expr = p.expression()
		if expr.GetType() != ast.Identifier {
			p.error(expr.GetToken(), "expected identifier in enum declaration")
			return &ast.Improper{Token: expr.GetToken()}
		}
		fields = append(fields, expr.GetToken())
	}

	enumStmt.Fields = fields

	p.consume("expected body closure", lexer.RightBrace)

	return enumStmt
}

func (p *Parser) structDeclarationStmt() ast.Node {
	stmt := ast.StructDeclarationStmt{
		IsLocal: p.peek(-1).Type != lexer.Pub,
	}
	stmt.Token = p.peek(-1)

	name, ok := p.consume("expected the name of the structure", lexer.Identifier)

	if ok {
		stmt.Name = name
	} else {
		return &ast.Improper{Token: stmt.Token}
	}

	_, ok = p.consume("expected opening of the struct body", lexer.LeftBrace)
	if !ok {
		return &ast.Improper{Token: stmt.Token}
	}
	stmt.Methods = &[]ast.MethodDeclarationStmt{}
	for !p.match(lexer.RightBrace) {
		if p.match(lexer.Fn) {
			method, ok := p.methodDeclarationStmt(stmt.IsLocal).(*ast.MethodDeclarationStmt)
			if ok {
				*stmt.Methods = append(*stmt.Methods, *method)
			}
		} else if p.match(lexer.New) {
			construct, ok := p.constructorDeclarationStmt().(*ast.ConstructorStmt)
			if ok {
				stmt.Constructor = construct
			}
		} else {
			field := p.fieldDeclarationStmt()
			if field.GetType() != ast.NA {
				stmt.Fields = append(stmt.Fields, *field.(*ast.FieldDeclarationStmt))
			} else {
				p.error(p.peek(), "unknown statement inside struct")
			}
		}
	}

	if stmt.Constructor == nil {
		stmt.Constructor = &ast.ConstructorStmt{
			Token: stmt.Token,
		}
	}

	return &stmt
}

func (p *Parser) constructorDeclarationStmt() ast.Node {
	stmt := ast.ConstructorStmt{Token: p.peek(-1)}

	stmt.Params = p.parameters(lexer.LeftParen, lexer.RightParen)

	stmt.Body = p.getBody()

	if stmt.Body == nil {
		return &ast.Improper{Token: stmt.Token}
	}

	return &stmt
}

func (p *Parser) fieldDeclarationStmt() ast.Node {
	stmt := ast.FieldDeclarationStmt{}

	typ, ident := p.TypeWithVar()
	if ident.GetType() != ast.Identifier {
		return &ast.Improper{}
	}

	idents := []lexer.Token{ident.GetToken()}
	types := []*ast.TypeExpr{typ}
	for p.match(lexer.Comma) {
		typ, ident := p.TypeWithVar()
		if ident.GetType() != ast.Identifier {
			return &ast.Improper{}
		}

		idents = append(idents, ident.GetToken())
		types = append(types, typ)
	}

	stmt.Identifiers = idents
	stmt.Types = types

	if !p.match(lexer.Equal) {
		stmt.Values = []ast.Node{}
		return &stmt
	}

	expr := p.expression()
	if expr.GetType() == ast.NA {
		p.error(p.peek(), "expected expression")
	}

	exprs := []ast.Node{expr}
	for p.match(lexer.Comma) {
		expr = p.expression()
		if expr.GetType() == ast.NA {
			p.error(p.peek(), "expected expression")
		}
		exprs = append(exprs, expr)
	}
	stmt.Values = exprs

	return &stmt
}

func (p *Parser) methodDeclarationStmt(IsLocal bool) ast.Node {
	fnDec := ast.MethodDeclarationStmt{
		IsLocal: IsLocal,
	}

	ident, ok := p.consume("expected a function name", lexer.Identifier)
	if !ok {
		return &fnDec
	}

	fnDec.Name = ident
	fnDec.Params = p.parameters(lexer.LeftParen, lexer.RightParen)

	ret := make([]*ast.TypeExpr, 0)
	for p.check(lexer.Identifier) {
		ret = append(ret, p.Type())
		if !p.check(lexer.Comma) {
			break
		} else {
			p.advance()
		}
	}
	fnDec.Return = ret
	fnDec.Body = p.getBody()

	return &fnDec
}

func (p *Parser) ifStmt(else_exists bool, is_else bool, is_elseif bool) *ast.IfStmt {
	ifStm := ast.IfStmt{
		Token: p.peek(-1),
	}

	var expr ast.Node
	if !is_else {
		expr = p.multiComparison()
		// if exprType == ast.Identifier && !(p.isMultiComparison() || p.check(lexer.LeftBrace)) {
		// 	for !p.check(lexer.LeftBrace) {
		// 		p.advance()
		// 	}
		// }
		// if exprType != ast.BinaryExpression && exprType != ast.Identifier && exprType != ast.UnaryExpression {
		// 	p.error(expr.GetToken(), "if condition is not a valid expression")
		// 	for !p.check(lexer.LeftBrace) {
		// 		p.advance()
		// 	}
		// }
	}
	ifStm.BoolExpr = expr
	ifStm.Body = p.getBody()
	if is_else || is_elseif {
		return &ifStm
	}
	for p.match(lexer.Else) {
		if else_exists {
			p.error(p.peek(-1), "cannot have two else statements in an if statement")
		}
		var ifbody *ast.IfStmt
		if p.match(lexer.If) {
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

func (p *Parser) assignmentStmt() ast.Node {
	expr := p.expression()

	idents := []ast.Node{expr}

	for p.match(lexer.Comma) { // memberExpr or IdentifierExpr
		identExpr := p.expression()
		idents = append(idents, identExpr)
	}

	if p.match(lexer.Equal) {
		values := []ast.Node{p.expression()}
		for p.match(lexer.Comma) {
			expr2 := p.expression()

			values = append(values, expr2)
		}
		expr = &ast.AssignmentStmt{Identifiers: idents, Values: values, Token: p.peek(-1)}
	} else if p.match(lexer.PlusEqual, lexer.MinusEqual, lexer.SlashEqual, lexer.StarEqual, lexer.CaretEqual, lexer.ModuloEqual) {
		assignOp := p.peek(-1)
		op := p.getOp(assignOp)
		if len(idents) > 1 {
			p.error(assignOp, "cannot assign to multiple variables with this operator")
		}
		expr2 := p.term()
		binExpr := p.createBinExpr(expr, op, op.Type, op.Lexeme, &ast.GroupExpr{Expr: expr2, ValueType: expr2.GetValueType(), Token: expr2.GetToken()})
		expr = &ast.AssignmentStmt{Identifiers: idents, Values: []ast.Node{binExpr}, Token: assignOp}
	}

	return expr
}

func (p *Parser) returnStmt() ast.Node {
	returnStmt := ast.ReturnStmt{
		Token: p.peek(-1),
	}

	args := []ast.Node{}
	expr := p.expression()
	if expr.GetType() == ast.NA {
		return &returnStmt
	}
	args = append(args, expr)
	for p.match(lexer.Comma) {
		expr = p.expression()
		if expr.GetType() == ast.NA {
			p.error(p.peek(), "expected expression")
		}
		args = append(args, expr)
	}
	returnStmt.Args = args

	return &returnStmt
}

func (p *Parser) yieldStmt() ast.Node {
	yieldStmt := ast.YieldStmt{
		Token: p.peek(-1),
	}

	if p.peek().Type == lexer.RightBrace {
		return &yieldStmt
	}
	args := []ast.Node{}
	expr := p.expression()
	args = append(args, expr)
	for p.match(lexer.Comma) {
		expr = p.expression()
		if expr.GetType() == ast.NA {
			p.error(p.peek(), "expected expression")
		}
		args = append(args, expr)
	}
	yieldStmt.Args = args

	return &yieldStmt
}

func (p *Parser) functionDeclarationStmt() ast.Node {
	fnDec := ast.FunctionDeclarationStmt{}

	fnDec.IsLocal = p.peek(-2).Type != lexer.Pub

	ident, ok := p.consume("expected a function name", lexer.Identifier)
	if !ok {
		return &fnDec
	}

	fnDec.Name = ident
	fnDec.Params = p.parameters(lexer.LeftParen, lexer.RightParen)

	fnDec.Return = p.returnings()
	fnDec.Body = p.getBody()

	return &fnDec
}

func (p *Parser) addToStmt() ast.Node {
	add := ast.AddStmt{
		Token: p.peek(-1),
	}

	add.Value = p.expression()
	if add.GetType() == ast.NA {
		p.error(p.peek(), "expected expression")
	}

	if _, ok := p.consume("expected keyword 'to' after expression in an 'add' statement", lexer.To); !ok {
		return &add
	}

	if ident, ok := p.consume("expected identifier after keyword 'to'", lexer.Identifier); ok {
		add.Identifier = ident.Lexeme
	}

	return &add
}

func (p *Parser) removeFromStmt() ast.Node {
	remove := ast.RemoveStmt{
		Token: p.peek(-1),
	}

	remove.Value = p.expression()
	if remove.GetType() == ast.NA {
		p.error(p.peek(), "expected expression")
	}

	if _, ok := p.consume("expected keyword 'from' after expression in a 'remove' statement", lexer.From); !ok {
		return &remove
	}

	if ident, ok := p.consume("expected identifier after keyword 'from'", lexer.Identifier); ok {
		remove.Identifier = ident.Lexeme
	}

	return &remove
}

func (p *Parser) repeatStmt() ast.Node {
	repeatStmt := ast.RepeatStmt{
		Token: p.peek(-1),
	}

	gotIterator := false
	if p.check(lexer.Number) ||
		p.check(lexer.Fixed) ||
		p.check(lexer.FixedPoint) ||
		p.check(lexer.Radian) ||
		p.check(lexer.Degree) ||
		p.check(lexer.Identifier) {

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
		if p.match(lexer.With) {
			identExpr := p.expression()
			if variableAssigned {
				p.error(p.peek(-1), "duplicate keyword 'with' in repeat statement")
			}
			variableAssigned = true
			if identExpr.GetType() != ast.Identifier {
				p.error(identExpr.GetToken(), "expected identifier expression after keyword 'with'")
			} else {
				repeatStmt.Variable = identExpr.(*ast.IdentifierExpr)
			}
		} else if p.match(lexer.To) {
			if iteratorAssgined {
				p.error(p.peek(-1), "duplicate keyword 'to' in repeat statement")
			}
			iteratorAssgined = true
			if gotIterator {
				p.error(p.peek(-1), "unnecessary redefinition of iterator")
			} else {
				repeatStmt.Iterator = p.expression()
				if repeatStmt.Iterator.GetType() == ast.NA {
					p.error(repeatStmt.Iterator.GetToken(), "unknown expression after keyword 'to'")
				}
			}
		} else if p.match(lexer.By) {
			if skipAssigned {
				p.error(p.peek(-1), "duplicate keyword 'by' in repeat statement")
			}
			skipAssigned = true
			repeatStmt.Skip = p.expression()
			if repeatStmt.Skip.GetType() == ast.NA {
				p.error(repeatStmt.Skip.GetToken(), "unknown expression after keyword 'by'")
			}
		} else if p.match(lexer.From) {
			if startAssigned {
				p.error(p.peek(-1), "duplicate keyword 'from' in repeat statement")
			}
			startAssigned = true
			repeatStmt.Start = p.expression()
			if repeatStmt.Start.GetType() == ast.NA {
				p.error(repeatStmt.Start.GetToken(), "unknown expression after keyword 'from'")
			}
		}
	}

	if repeatStmt.Iterator == nil {
		p.error(repeatStmt.Token, "no iterator provided in repeat statement")
		repeatStmt.Iterator = &ast.LiteralExpr{Token: repeatStmt.Token, Value: "1", ValueType: ast.Number}
	}

	repeatStmt.Body = p.getBody()

	return &repeatStmt
}

func (p *Parser) whileStmt() ast.Node {
	whileStmt := &ast.WhileStmt{}

	condtion := p.expression()

	if condtion.GetType() == ast.NA {
		p.error(condtion.GetToken(), "Expected an expressions after 'while'")
		return &ast.Improper{}
	}

	whileStmt.Condtion = condtion

	whileStmt.Body = p.getBody()

	return whileStmt
}

func (p *Parser) forStmt() ast.Node {
	forStmt := ast.ForStmt{
		Token: p.peek(-1),
	}

	if p.peek().Type == lexer.Identifier &&
		p.peek(1).Type == lexer.Comma {
		identExpr := p.expression()
		if identExpr.GetType() != ast.Identifier {
			p.error(identExpr.GetToken(), "expected identifier expression after keyword 'for'")
		} else {
			forStmt.KeyValuePair[0] = identExpr.(*ast.IdentifierExpr)
		}
		p.match(lexer.Comma)
		identExpr = p.expression()
		if identExpr.GetType() != ast.Identifier {
			p.error(identExpr.GetToken(), "expected identifier expression after a comma")
		} else {
			forStmt.KeyValuePair[1] = identExpr.(*ast.IdentifierExpr)
		}
	} else {
		identExpr := p.expression()
		if identExpr.GetType() != ast.Identifier {
			p.error(identExpr.GetToken(), "expected identifier expression after keyword 'for'")
		} else {
			forStmt.KeyValuePair[0] = identExpr.(*ast.IdentifierExpr)
		}
	}

	p.consume("expected keyword 'in' after for loop variables", lexer.In)

	forStmt.Iterator = p.expression()

	if forStmt.Iterator == nil {
		p.error(forStmt.Token, "no iterator provided in for loop statement")
		forStmt.Iterator = &ast.LiteralExpr{Token: forStmt.Token, Value: "[1]", ValueType: ast.List}
	}

	forStmt.Body = p.getBody()

	return &forStmt
}

func (p *Parser) tickStmt() ast.Node {
	tickStmt := ast.TickStmt{
		Token: p.peek(-1),
	}

	if p.match(lexer.With) {
		identExpr := p.expression()
		if identExpr.GetType() != ast.Identifier {
			p.error(identExpr.GetToken(), "expected identifier expression after keyword 'with'")
			return &tickStmt
		}
		tickStmt.Variable = *identExpr.(*ast.IdentifierExpr)
	}

	tickStmt.Body = p.getBody()

	return &tickStmt
}

func (p *Parser) variableDeclarationStmt() ast.Node {
	variable := ast.VariableDeclarationStmt{
		Token:   p.peek(-1),
		IsLocal: p.peek(-1).Type == lexer.Let,
	}

	typ, ide := p.TypeWithVar()
	if ide.GetType() != ast.Identifier {
		p.error(ide.GetToken(), "expected identifier in variable declaration")
		return &ast.Improper{}
	}

	idents := []lexer.Token{ide.GetToken()}
	types := []*ast.TypeExpr{typ}
	for p.match(lexer.Comma) {
		typ, ide = p.TypeWithVar()
		if ide.GetType() != ast.Identifier {
			p.error(ide.GetToken(), "expected identifier in variable declaration")
			return &ast.Improper{}
		}

		idents = append(idents, ide.GetToken())
		types = append(types, typ)
	}

	variable.Identifiers = idents
	variable.Types = types

	if !p.match(lexer.Equal) {
		variable.Values = []ast.Node{}
		return &variable
	}

	expr := p.expression()
	if expr.GetType() == ast.NA {
		p.error(p.peek(), "expected expression")
	}

	exprs := []ast.Node{expr}
	for p.match(lexer.Comma) {
		expr = p.expression()
		if expr.GetType() == ast.NA {
			p.error(p.peek(), "expected expression")
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
		p.error(p.peek(), "expected filepath")
		return &ast.Improper{}
	}
	useStmt.Path = filepath.(*ast.EnvPathExpr)

	return &useStmt
}

func (p *Parser) matchStmt(isExpr bool) *ast.MatchStmt {
	matchStmt := ast.MatchStmt{}

	matchStmt.ExprToMatch = p.expression()

	p.consume("expected opening of the match body", lexer.LeftBrace)

	caseStmts, stop := p.caseStmt(isExpr)
	for !stop {
		matchStmt.Cases = append(matchStmt.Cases, caseStmts...)
		caseStmts, stop = p.caseStmt(isExpr)
		for i := range caseStmts {
			if caseStmts[i].Expression.GetToken().Lexeme == "_" {
				matchStmt.HasDefault = true
			}
		}
	}
	matchStmt.Cases = append(matchStmt.Cases, caseStmts...)

	p.consume("expected closing of the match body", lexer.RightBrace)

	return &matchStmt
}

func (p *Parser) caseStmt(isExpr bool) ([]ast.CaseStmt, bool) {
	caseStmts := []ast.CaseStmt{}

	caseStmt := ast.CaseStmt{}
	caseStmt.Expression = p.expression()
	if caseStmt.Expression.GetType() == ast.NA {
		return caseStmts, true
	}
	caseStmts = append(caseStmts, caseStmt)
	for p.match(lexer.Comma) {
		caseStmt.Expression = p.expression()
		caseStmts = append(caseStmts, caseStmt)
		if caseStmt.Expression.GetType() == ast.NA {
			return caseStmts, true
		}
	}

	p.consume("expected fat arrow after expression in case", lexer.FatArrow)

	if p.check(lexer.LeftBrace) {
		body := p.getBody()
		for i := range caseStmts {
			caseStmts[i].Body = body
		}
	} else {
		expr := p.expression()
		if expr.GetType() == ast.NA {
			p.error(expr.GetToken(), "expected expression or '{' after fat arrow")
		}
		args := []ast.Node{expr}
		for p.match(lexer.Comma) {
			expr = p.expression()
			if expr.GetType() == ast.NA {
				p.error(expr.GetToken(), "expected expression")
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
	if p.check(lexer.RightBrace) {
		return caseStmts, true
	}

	return caseStmts, false
}
