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
	returnNode = ast.NewImproper(p.peek(), ast.NA)

	defer func() {
		if errMsg := recover(); errMsg != nil {
			// If the error is a parseError, synchronize to
			// the next statement. If not, propagate the panic.
			if _, ok := errMsg.(ParserError); ok {
				p.synchronize()
			} else {
				fmt.Printf("panic: %s\nstacktrace:\n", errMsg)
				debug.PrintStack()
				os.Exit(1)
			}
		}
	}()

	token := p.peek().Type

	switch token {
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
	case tokens.While:
		p.advance()
		returnNode = p.whileStmt()
		return
	case tokens.Match:
		p.advance()
		returnNode = p.matchStmt(false)
		return
	}

	returnNode = p.expressionStatement()
	return
}

func (p *Parser) expressionStatement() ast.Node {
	expr := p.expression()

	if expr.GetType() == ast.Identifier || expr.GetType() == ast.EnvironmentAccessExpression {
		return p.assignmentStmt(expr)
	}

	typ := expr.GetType()
	if typ != ast.CallExpression && typ != ast.MethodCallExpression {
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
	destroyStmt.Args = p.functionArgs()

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
			ifbody = p.ifStmt(else_exists, false, true)
			ifStmt.Elseifs = append(ifStmt.Elseifs, ifbody)
		} else {
			if else_exists {
				p.Alert(&alerts.MoreThanOneElseStatement{}, alerts.NewSingle(p.peek(-1)))
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
		exprs, _ := p.expressions("in assignment statement")
		return &ast.AssignmentStmt{Identifiers: idents, Values: exprs, Token: p.peek(-1)}
	}
	if p.match(tokens.PlusEqual, tokens.MinusEqual, tokens.SlashEqual, tokens.StarEqual, tokens.CaretEqual, tokens.ModuloEqual, tokens.BackSlashEqual) {
		if len(idents) > 1 {
			p.Alert(&alerts.MultipleIdentifiersInCompoundAssignment{}, alerts.NewMulti(expr.GetToken(), idents[len(idents)].GetToken()))
		}
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
		op.Lexeme = string(op.Type)

		expr2 := p.expression()
		if expr2.GetType() == ast.NA {
			p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(p.peek()), "in assignment statement")
		} else {
			binExpr := p.createBinExpr(idents[1], op, op.Type, op.Lexeme, &ast.GroupExpr{Expr: expr2, ValueType: expr2.GetValueType(), Token: expr2.GetToken()})
			values = append(values, binExpr)
		}
		for p.match(tokens.Comma) {
			expr2 := p.expression()

			if expr2.GetType() == ast.NA {
				p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(p.peek()), "in assignment statement")
			} else {
				binExpr := p.createBinExpr(idents[1], op, op.Type, op.Lexeme, &ast.GroupExpr{Expr: expr2, ValueType: expr2.GetValueType(), Token: expr2.GetToken()})
				values = append(values, binExpr)
			}
		}
		return &ast.AssignmentStmt{Identifiers: idents, Values: values, Token: assignOp}
	}
	p.Alert(&alerts.ExpectedAssignmentSymbol{}, alerts.NewSingle(p.peek()))

	return ast.NewImproper(expr.GetToken(), ast.AssignmentStatement)
}

func (p *Parser) returnStmt() ast.Node {
	returnStmt := &ast.ReturnStmt{
		Token: p.peek(-1),
		Args:  []ast.Node{},
	}

	if p.context.FunctionReturns.Count() == 0 {
		return returnStmt
	}

	if p.context.FunctionReturns.Top().Item != 0 {
		returnStmt.Args, _ = p.expressions("in return arguments")
	}

	return returnStmt
}

func (p *Parser) yieldStmt() ast.Node {
	yieldStmt := ast.YieldStmt{
		Token: p.peek(-1),
	}

	if p.peek().Type == tokens.RightBrace {
		return &yieldStmt
	}

	yieldStmt.Args, _ = p.expressions("in yield statement")

	return &yieldStmt
}

func (p *Parser) repeatStmt() ast.Node {
	repeatStmt := ast.RepeatStmt{
		Token: p.peek(-1),
	}

outer:
	for range 4 {
		token := p.peek()
		switch token.Type {
		case tokens.With:
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
			it := p.expression()
			if repeatStmt.Iterator != nil {
				p.Alert(&alerts.IteratorRedefinition{}, alerts.NewSingle(p.peek(-1)), "in repeat statement")
				continue
			}
			repeatStmt.Iterator = it
		case tokens.By:
			skip := p.expression()
			if repeatStmt.Skip != nil {
				p.Alert(&alerts.DuplicateKeyword{}, alerts.NewSingle(p.peek(-1)), tokens.By)
				continue
			}
			repeatStmt.Skip = skip
		case tokens.From:
			start := p.expression()
			if repeatStmt.Start != nil {
				p.Alert(&alerts.DuplicateKeyword{}, alerts.NewSingle(p.peek(-1)), tokens.From)
			}
			repeatStmt.Start = start
		case tokens.LeftBrace:
			break outer
		default:
			repeatStmt.Iterator = p.expression()
		}
	}

	if repeatStmt.Iterator == nil {
		p.Alert(&alerts.MissingIterator{}, alerts.NewSingle(repeatStmt.Token), "in repeat statement")
		repeatStmt.Iterator = &ast.LiteralExpr{Token: repeatStmt.Token, Value: "1", ValueType: ast.Number}
	}

	var success bool
	repeatStmt.Body, success = p.body(false, false)
	if !success {
		return ast.NewImproper(repeatStmt.Token, ast.NA)
	}

	return &repeatStmt
}

func (p *Parser) whileStmt() ast.Node {
	whileStmt := &ast.WhileStmt{}

	condition := p.multiComparison()

	if condition.GetType() == ast.NA {
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
	if forStmt.Iterator.GetType() == ast.NA {
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
	tickStmt.Body, success = p.body(true, true)
	if !success {
		return ast.NewImproper(tickStmt.Token, ast.NA)
	}

	return &tickStmt
}

func (p *Parser) useStmt() ast.Node {
	return nil
}

func (p *Parser) matchStmt(isExpr bool) *ast.MatchStmt {
	return nil
}

func (p *Parser) caseStmt(isExpr bool) ([]ast.CaseStmt, bool) {
	return nil, false
}
