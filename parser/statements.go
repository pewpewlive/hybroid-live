package parser

import (
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/tokens"
)

func (p *Parser) statement() (returnNode ast.Node) {
	returnNode = ast.NewImproper(p.peek(), ast.NA)

	switch p.advance().Type {
	case tokens.Return:
		returnNode = p.returnStmt()
	case tokens.Yield:
		returnNode = p.yieldStmt()
	case tokens.Break:
		returnNode = &ast.BreakStmt{Token: p.peek(-1)}
	case tokens.Destroy:
		returnNode = p.destroyStmt()
	case tokens.Continue:
		returnNode = &ast.ContinueStmt{Token: p.peek(-1)}
	case tokens.If:
		returnNode = p.ifStmt(false, false, false)
	case tokens.Repeat:
		returnNode = p.repeatStmt()
	case tokens.For:
		returnNode = p.forStmt()
	case tokens.Tick:
		returnNode = p.tickStmt()
	case tokens.Use:
		returnNode = p.useStmt()
	case tokens.While:
		returnNode = p.whileStmt()
	case tokens.Match:
		returnNode = p.matchStmt(false)
	}

	if returnNode.GetType() == ast.NA {
		p.disadvance()
	}

	return
}

func (p *Parser) expressionStatement() ast.Node {
	expr := p.expression()
	exprType := expr.GetType()

	if exprType == ast.Identifier || exprType == ast.EnvironmentAccessExpression ||
		exprType == ast.MemberExpression || exprType == ast.FieldExpression {
		return p.assignmentStmt(expr)
	}

	if exprType == ast.NA {
		improperType := expr.(*ast.Improper).Type
		if p.isCall(improperType) {
			return expr
		}
	}

	if !p.isCall(exprType) {
		return ast.NewImproper(expr.GetToken(), ast.NA)
	}

	return expr
}

func (p *Parser) destroyStmt() ast.Node {
	destroyStmt := ast.DestroyStmt{
		Token: p.peek(-1),
	}

	expr := p.self()
	exprType := expr.GetType()

	if exprType != ast.Identifier && exprType != ast.EnvironmentAccessExpression && exprType != ast.SelfExpression {
		p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(expr.GetToken()), "or environment access expression or a self expression in destroy statement")
	}
	destroyStmt.Identifier = expr
	destroyStmt.Generics, _ = p.genericArgs()
	args, ok := p.functionArgs()
	if !ok {
		return ast.NewImproper(destroyStmt.Token, ast.DestroyStatement)
	}
	destroyStmt.Args = args

	return &destroyStmt
}

func (p *Parser) ifStmt(else_exists bool, is_else bool, is_elseif bool) *ast.IfStmt {
	ifStmt := ast.IfStmt{
		Token: p.peek(-1),
	}

	var expr ast.Node = nil
	if !is_else {
		expr = p.multiComparison()
	}
	ifStmt.BoolExpr = expr
	ifStmt.Body, _ = p.body(true, false)

	if is_else || is_elseif {
		return &ifStmt
	}
	for p.match(tokens.Else) {
		var ifbody *ast.IfStmt
		if p.match(tokens.If) {
			if else_exists {
				p.Alert(&alerts.ElseIfBlockAfterElseBlock{}, alerts.NewSingle(p.peek(-1)))
			}
			ifbody = p.ifStmt(else_exists, false, true)
			ifStmt.Elseifs = append(ifStmt.Elseifs, ifbody)
		} else {
			if else_exists {
				p.Alert(&alerts.MoreThanOneElseBlock{}, alerts.NewSingle(p.peek(-1)))
			}
			else_exists = true
			ifbody = p.ifStmt(else_exists, true, false)
			ifStmt.Else = ifbody
		}
	}

	return &ifStmt
}

func (p *Parser) assignmentStmt(expr ast.Node) ast.Node {
	idents := []ast.Node{expr}

	for p.match(tokens.Comma) {
		expr := p.expression()
		idents = append(idents, expr)
	}
	values := []ast.Node{}
	if p.match(tokens.Equal) {
		exprs, _ := p.expressions("in assignment statement", false)
		return &ast.AssignmentStmt{Identifiers: idents, Values: exprs, Token: p.peek(-1)}
	}
	if p.match(tokens.PlusEqual, tokens.MinusEqual, tokens.SlashEqual, tokens.StarEqual, tokens.CaretEqual, tokens.ModuloEqual, tokens.BackSlashEqual) {
		assignOp := p.peek(-1)
		op := tokens.Token{Literal: assignOp.Literal, Location: assignOp.Location}
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
		case tokens.BackSlashEqual:
			op.Type = tokens.BackSlash
		}

		expr2 := p.expression()
		if ast.IsImproper(expr2, ast.NA) {
			p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(p.peek()), "in assignment statement")
		} else {
			values = append(values, expr2)
		}
		for p.match(tokens.Comma) {
			expr2 := p.expression()

			if ast.IsImproper(expr2, ast.NA) {
				p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(p.peek()), "in assignment statement")
				continue
			}
			values = append(values, expr2)
		}
		return &ast.AssignmentStmt{Identifiers: idents, Values: values, AssignOp: op, Token: idents[0].GetToken()}
	}
	//p.Alert(&alerts.ExpectedAssignmentSymbol{}, alerts.NewSingle(p.peek()))
	return ast.NewImproper(expr.GetToken(), ast.NA)
}

func (p *Parser) returnStmt() ast.Node {
	returnStmt := &ast.ReturnStmt{
		Token: p.peek(-1),
		Args:  []ast.Node{},
	}

	if p.peek().Line != returnStmt.Token.Line {
		return returnStmt
	}
	returnStmt.Args, _ = p.expressions("in return arguments", false)

	return returnStmt
}

func (p *Parser) yieldStmt() ast.Node {
	yieldStmt := &ast.YieldStmt{
		Token: p.peek(-1),
	}

	if p.peek().Line != yieldStmt.Token.Line {
		return yieldStmt
	}
	yieldStmt.Args, _ = p.expressions("in yield statement", false)

	return yieldStmt
}

func (p *Parser) repeatStmt() ast.Node {
	repeatStmt := ast.RepeatStmt{
		Token: p.peek(-1),
	}

	allowedExprTypes := []ast.NodeType{ast.Identifier, ast.FieldExpression, ast.MemberExpression, ast.CallExpression, ast.MethodCallExpression, ast.LiteralExpression, ast.BinaryExpression, ast.UnaryExpression}

	i := 0
outer:
	for {
		token := p.peek()
		switch token.Type {
		case tokens.With:
			p.advance()
			identExpr := p.expression()
			if identExpr.GetType() != ast.Identifier {
				p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(identExpr.GetToken()), "after keyword 'with'")
				continue
			}
			if repeatStmt.Variable != nil {
				p.Alert(&alerts.DuplicateKeyword{}, alerts.NewSingle(p.peek(-1)), tokens.With)
				continue
			}
			repeatStmt.Variable = identExpr.(*ast.IdentifierExpr)
		case tokens.To:
			p.advance()
			it := p.limitedExpression("as iterator in repeat statement", allowedExprTypes...)
			if repeatStmt.Iterator != nil {
				p.Alert(&alerts.IteratorRedefinition{}, alerts.NewSingle(p.peek(-1)), "in repeat statement")
				continue
			}
			repeatStmt.Iterator = it
		case tokens.By:
			p.advance()
			skip := p.limitedExpression("as skip expression in repeat statement", allowedExprTypes...)
			if repeatStmt.Skip != nil {
				p.Alert(&alerts.DuplicateKeyword{}, alerts.NewSingle(p.peek(-1)), tokens.By)
				continue
			}
			repeatStmt.Skip = skip
		case tokens.From:
			p.advance()
			start := p.limitedExpression("as from expression in repeat statement", allowedExprTypes...)
			if repeatStmt.Start != nil {
				p.Alert(&alerts.DuplicateKeyword{}, alerts.NewSingle(p.peek(-1)), tokens.From)
			}
			repeatStmt.Start = start
		case tokens.LeftBrace:
			break outer
		default:
			if i != 0 {
				break outer
			}
			repeatStmt.Iterator = p.expression()
		}
		i += 1
	}

	if repeatStmt.Iterator == nil {
		p.Alert(&alerts.MissingIterator{}, alerts.NewSingle(repeatStmt.Token), "in repeat statement")
		repeatStmt.Iterator = &ast.LiteralExpr{Token: repeatStmt.Token, Value: "1", ValueType: ast.Number}
	}

	var success bool
	repeatStmt.Body, success = p.body(false, false)
	if !success {
		return ast.NewImproper(repeatStmt.Token, ast.RepeatStatement)
	}

	return &repeatStmt
}

func (p *Parser) whileStmt() ast.Node {
	whileStmt := &ast.WhileStmt{
		Token: p.peek(-1),
	}

	condition := p.multiComparison()

	if ast.IsImproper(condition, ast.NA) {
		p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(condition.GetToken()))
		return ast.NewImproper(condition.GetToken(), ast.WhileStatement)
	}

	whileStmt.Condition = condition

	var success bool
	whileStmt.Body, success = p.body(false, false)
	if !success {
		return ast.NewImproper(whileStmt.Token, ast.WhileStatement)
	}

	return whileStmt
}

func (p *Parser) forStmt() ast.Node {
	forStmt := ast.ForStmt{
		Token: p.peek(-1),
	}

	identExpr := p.expression()
	if identExpr.GetType() != ast.Identifier {
		p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(identExpr.GetToken()), "after keyword 'for' in for loop statement")
	} else {
		forStmt.First = identExpr.(*ast.IdentifierExpr)
	}

	if p.match(tokens.Comma) {
		identExpr = p.expression()
		if identExpr.GetType() != ast.Identifier {
			p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(identExpr.GetToken()))
		} else {
			forStmt.Second = identExpr.(*ast.IdentifierExpr)
		}
	}

	p.consume(p.NewAlert(&alerts.ExpectedKeyword{}, alerts.NewSingle(p.peek()), tokens.In), tokens.In)

	forStmt.Iterator = p.expression()
	if ast.IsImproper(forStmt.Iterator, ast.NA) {
		p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(forStmt.Iterator.GetToken()))
	}

	// if forStmt.Iterator == nil {
	// 	p.Alert(&alerts.MissingIterator{}, alerts.NewSingle(forStmt.Token), "in for statement")
	// 	forStmt.Iterator = &ast.LiteralExpr{Token: forStmt.Token, Value: "[1]", ValueType: ast.List}
	// }

	var success bool
	forStmt.Body, success = p.body(true, true)
	if !success {
		return ast.NewImproper(forStmt.Token, ast.ForStatement)
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
	tickStmt.Body, success = p.body(false, false)
	if !success {
		return ast.NewImproper(tickStmt.Token, ast.TickStatement)
	}

	return &tickStmt
}

func (p *Parser) useStmt() ast.Node {
	useStmt := ast.UseStmt{}

	filepath := p.envPathExpr()
	if filepath.GetType() != ast.EnvironmentPathExpression {
		p.Alert(&alerts.ExpectedEnvironmentPathExpression{}, alerts.NewMulti(filepath.GetToken(), p.peek()))
		return ast.NewImproper(p.peek(), ast.UseStatement)
	}
	useStmt.Path = filepath.(*ast.EnvPathExpr)

	return &useStmt
}

func (p *Parser) matchStmt(isExpr bool) *ast.MatchStmt {
	matchStmt := ast.MatchStmt{
		Token: p.peek(-1),
	}

	matchStmt.ExprToMatch = p.expression()

	start, _ := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.LeftBrace), tokens.LeftBrace)

	caseStmts, stop := p.caseStmt(isExpr)
	for !stop {
		matchStmt.Cases = append(matchStmt.Cases, caseStmts...)
		caseStmts, stop = p.caseStmt(isExpr)
		for i := range caseStmts {
			if caseStmts[i].Expression.GetToken().Lexeme == "else" {
				if matchStmt.HasDefault {
					p.Alert(&alerts.MoreThanOneDefaultCase{}, alerts.NewSingle(caseStmts[i].Expression.GetToken()))
					continue
				}
				matchStmt.HasDefault = true
			}
			matchStmt.Cases = append(matchStmt.Cases, caseStmts[i])
		}
	}

	p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewMulti(start, p.peek()), tokens.RightBrace), tokens.RightBrace)

	if len(matchStmt.Cases) < 1 {
		p.Alert(&alerts.InsufficientCases{}, alerts.NewMulti(matchStmt.Token, p.peek(-1)))
	}
	if !matchStmt.HasDefault && isExpr {
		p.Alert(&alerts.DefaultCaseMissing{}, alerts.NewMulti(matchStmt.Token, p.peek(-1)))
	}

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
		body, _ := p.body(true, false)
		for i := range caseStmts { // "hello" =>
			caseStmts[i].Body = body
		}
		if p.check(tokens.RightBrace) {
			return caseStmts, true
		}

		return caseStmts, false
	}
	expr := p.expression()
	if expr.GetType() == ast.NA {
		p.Alert(&alerts.ExpectedExpressionOrBody{}, alerts.NewSingle(p.peek()))
	}
	args := []ast.Node{expr}
	for p.match(tokens.Comma) {
		expr = p.expression()
		if ast.IsImproper(expr, ast.NA) {
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
	if p.check(tokens.RightBrace) {
		return caseStmts, true
	}

	return caseStmts, false
}
