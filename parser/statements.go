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
		returnNode = p.returnStatement()
	case tokens.Yield:
		returnNode = p.yieldStatement()
	case tokens.Break:
		returnNode = &ast.BreakStmt{Token: p.peek(-1)}
	case tokens.Destroy:
		returnNode = p.destroyStatement()
	case tokens.Continue:
		returnNode = &ast.ContinueStmt{Token: p.peek(-1)}
	case tokens.If:
		returnNode = p.ifStatement(false, false, false)
	case tokens.Repeat:
		returnNode = p.repeatStatement()
	case tokens.For:
		returnNode = p.forStatement()
	case tokens.Tick:
		returnNode = p.tickStatement()
	case tokens.Use:
		returnNode = p.useStatement()
	case tokens.While:
		returnNode = p.whileStatement()
	case tokens.Match:
		returnNode = p.matchStatement(false)
	}

	if ast.IsImproper(returnNode, ast.NA) {
		p.disadvance()
	}

	return
}

func (p *Parser) expressionStatement() ast.Node {
	expr := p.expression()
	exprType := expr.GetType()

	if exprType == ast.Identifier || exprType == ast.EnvironmentAccessExpression ||
		exprType == ast.MemberExpression || exprType == ast.FieldExpression {
		return p.assignmentStatement(expr)
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

func (p *Parser) destroyStatement() ast.Node {
	destroyStmt := ast.DestroyStmt{
		Token: p.peek(-1),
	}

	entityArgs, ok := p.genericArgs()
	if !ok {
		return ast.NewImproper(destroyStmt.Token, ast.DestroyStatement)
	}
	destroyStmt.EntityGenericArgs = entityArgs

	expr := p.AccessorExpr()
	if ast.IsImproper(expr, ast.NA) {
		p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(expr.GetToken()), "in destroy statement")
		return ast.NewImproper(destroyStmt.Token, ast.DestroyStatement)
	} else if expr.GetType() != ast.CallExpression {
		if expr.GetType() != ast.NA {
			p.Alert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), "in destroy statement")
		}
		return ast.NewImproper(destroyStmt.Token, ast.DestroyStatement)
	}
	call := expr.(*ast.CallExpr)

	destroyStmt.Identifier = call.Caller
	destroyStmt.GenericArgs = call.GenericArgs
	destroyStmt.Args = call.Args

	return &destroyStmt
}

func (p *Parser) ifStatement(else_exists bool, is_else bool, is_elseif bool) ast.Node {
	ifStmt := ast.IfStmt{
		Token: p.peek(-1),
	}

	var expr ast.Node = nil
	if !is_else {
		expr = p.multiComparison()
		if ast.IsImproper(expr, ast.NA) {
			return ast.NewImproper(ifStmt.Token, ast.IfStatement)
		}
	}
	ifStmt.BoolExpr = expr
	body, ok := p.body(true, false)
	if !ok {
		return ast.NewImproper(ifStmt.Token, ast.IfStatement)
	}
	ifStmt.Body = body

	if is_else || is_elseif {
		return &ifStmt
	}
	for p.match(tokens.Else) {
		var ifStmt2 ast.Node
		if p.match(tokens.If) {
			if else_exists {
				p.Alert(&alerts.ElseIfBlockAfterElseBlock{}, alerts.NewSingle(p.peek(-1)))
			}
			ifStmt2 = p.ifStatement(else_exists, false, true)
			if ifStmt2.GetType() == ast.NA {
				return ast.NewImproper(ifStmt.Token, ast.IfStatement)
			}
			ifStmt.Elseifs = append(ifStmt.Elseifs, ifStmt2.(*ast.IfStmt))
		} else {
			if else_exists {
				p.Alert(&alerts.MoreThanOneElseBlock{}, alerts.NewSingle(p.peek(-1)))
			}
			else_exists = true
			ifStmt2 = p.ifStatement(else_exists, true, false)
			if ifStmt2.GetType() == ast.NA {
				return ast.NewImproper(ifStmt.Token, ast.IfStatement)
			}
			ifStmt.Else = ifStmt2.(*ast.IfStmt)
		}
	}

	return &ifStmt
}

func (p *Parser) assignmentStatement(expr ast.Node) ast.Node {
	idents := []ast.Node{expr}

	for p.match(tokens.Comma) {
		exp := p.expression()
		if ast.IsImproper(exp, ast.NA) {
			return ast.NewImproper(expr.GetToken(), ast.NA)
		}
		idents = append(idents, exp)
	}
	if p.match(tokens.Equal) {
		equal := p.peek(-1)
		exprs, ok := p.expressions("in assignment statement", false)
		if !ok {
			return ast.NewImproper(expr.GetToken(), ast.AssignmentStatement)
		}
		return &ast.AssignmentStmt{Identifiers: idents, Values: exprs, AssignOp: equal, Token: idents[0].GetToken()}
	}

	isLeftShiftEqual := p.peek().Type == tokens.Less && p.peek(1).Type == tokens.Less && p.peek(2).Type == tokens.Equal
	isRightShiftEqual := p.peek().Type == tokens.Greater && p.peek(1).Type == tokens.Greater && p.peek(2).Type == tokens.Equal
	isNormalCompoundOp := p.check(tokens.PlusEqual, tokens.MinusEqual, tokens.SlashEqual, tokens.StarEqual, tokens.CaretEqual, tokens.ModuloEqual, tokens.BackSlashEqual, tokens.AmpersandEqual, tokens.PipeEqual, tokens.TildeEqual)
	var op tokens.Token
	if isNormalCompoundOp {
		op = p.advance()
	} else if isLeftShiftEqual {
		newToken, success := p.combineTokens(tokens.LeftShiftEqual, 3)
		if !success {
			return ast.NewImproper(expr.GetToken(), ast.NA)
		}
		op = newToken
	} else if isRightShiftEqual {
		newToken, success := p.combineTokens(tokens.RightShiftEqual, 3)
		if !success {
			return ast.NewImproper(expr.GetToken(), ast.NA)
		}
		op = newToken
	}

	if isNormalCompoundOp || isLeftShiftEqual || isRightShiftEqual {
		exprs, ok := p.expressions("in assignment statement", false)
		if !ok {
			return ast.NewImproper(expr.GetToken(), ast.AssignmentStatement)
		}
		return &ast.AssignmentStmt{Identifiers: idents, Values: exprs, AssignOp: op, Token: idents[0].GetToken()}
	}
	return ast.NewImproper(expr.GetToken(), ast.NA)
}

func (p *Parser) returnStatement() ast.Node {
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

func (p *Parser) yieldStatement() ast.Node {
	yieldStmt := &ast.YieldStmt{
		Token: p.peek(-1),
	}

	if p.peek().Line != yieldStmt.Token.Line {
		return yieldStmt
	}
	yieldStmt.Args, _ = p.expressions("in yield statement", false)

	return yieldStmt
}

func (p *Parser) repeatStatement() ast.Node {
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
		repeatStmt.Iterator = &ast.LiteralExpr{Token: repeatStmt.Token, Value: "1"}
	}

	var success bool
	repeatStmt.Body, success = p.body(false, false)
	if !success {
		return ast.NewImproper(repeatStmt.Token, ast.RepeatStatement)
	}

	return &repeatStmt
}

func (p *Parser) whileStatement() ast.Node {
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

func (p *Parser) forStatement() ast.Node {
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

func (p *Parser) tickStatement() ast.Node {
	tickStmt := ast.TickStmt{
		Token: p.peek(-1),
	}

	if p.match(tokens.With) {
		identExpr := p.expression()
		if identExpr.GetType() != ast.Identifier {
			p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(identExpr.GetToken()), "after keyword 'with'")
			return ast.NewImproper(tickStmt.Token, ast.TickStatement)
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

func (p *Parser) useStatement() ast.Node {
	useStmt := &ast.UseStmt{
		Token: p.peek(-1),
	}

	filepath := p.envPathExpr()
	if filepath.GetType() != ast.EnvironmentPathExpression {
		p.AlertMulti(&alerts.ExpectedEnvironmentPathExpression{}, filepath.GetToken(), p.peek())
		return ast.NewImproper(p.peek(), ast.UseStatement)
	}
	useStmt.PathExpr = filepath.(*ast.EnvPathExpr)

	return useStmt
}

func (p *Parser) matchStatement(isExpr bool) ast.Node {
	var matchType ast.NodeType
	if isExpr {
		matchType = ast.MatchExpression
	} else {
		matchType = ast.MatchStatement
	}
	matchStmt := ast.MatchStmt{
		Token: p.peek(-1),
	}

	matchStmt.ExprToMatch = p.expression()

	start, ok := p.alertSingleConsume(&alerts.ExpectedSymbol{}, tokens.LeftBrace)
	if !ok {
		return ast.NewImproper(matchStmt.Token, matchType)
	}

	for p.consumeTill("in match statement", start, tokens.RightBrace) {
		node, ok := p.caseStatement(isExpr)
		if !ok {
			p.synchronizeMatchBody()
			continue
		}
		caseStmt := node.(*ast.CaseStmt)
		if caseStmt.Expressions[0].GetToken().Lexeme == "else" {
			if matchStmt.HasDefault {
				p.Alert(&alerts.MoreThanOneDefaultCase{}, alerts.NewSingle(caseStmt.Expressions[0].GetToken()))
				continue
			}
			matchStmt.HasDefault = true
		}
		matchStmt.Cases = append(matchStmt.Cases, caseStmt)
	}

	return &matchStmt
}

func (p *Parser) caseStatement(isExpr bool) (ast.Node, bool) {
	token := p.peek()
	caseStmt := &ast.CaseStmt{}

	exprs := []ast.Node{}
	if p.match(tokens.Else) {
		exprs = append(exprs, &ast.IdentifierExpr{Name: p.peek(-1)})
	} else {
		exprs2, ok := p.expressions("in match case", false)
		if !ok {
			return ast.NewImproper(token, ast.CaseStatement), false
		}
		exprs = exprs2
	}

	caseStmt.Expressions = exprs
	_, ok := p.alertSingleConsume(&alerts.ExpectedSymbol{}, tokens.FatArrow, "in match case")
	if !ok {
		return ast.NewImproper(token, ast.CaseStatement), false
	}

	body := ast.Body{}
	if !p.check(tokens.LeftBrace) {
		args, ok := p.expressions("after '=>' in match case", false)
		if !ok {
			return ast.NewImproper(token, ast.CaseStatement), false
		}
		var argsStmt ast.Node
		if isExpr {
			argsStmt = &ast.YieldStmt{
				Args:  args,
				Token: args[0].GetToken(),
			}
		} else {
			argsStmt = &ast.ReturnStmt{
				Args:  args,
				Token: args[0].GetToken(),
			}
		}
		body.Append(argsStmt)
	} else {
		body2, ok2 := p.body(false, false)
		if !ok2 {
			return caseStmt, false
		}
		body = body2
	}
	caseStmt.Body = body

	return caseStmt, true
}
