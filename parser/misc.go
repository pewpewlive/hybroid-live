package parser

import (
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/tokens"
)

func (p *Parser) getFunctionParam() (ast.FunctionParam, bool) {
	functionParam := ast.FunctionParam{}
	typeExpr := p.typeExpr("in function parameters")
	peekType := p.peek().Type

	if peekType == tokens.Identifier {
		functionParam.Name = p.advance()
		functionParam.Type = typeExpr
		return functionParam, true
	}
	if typeExpr.Name.GetType() == ast.Identifier {
		functionParam.Name = typeExpr.Name.GetToken()
		functionParam.Type = nil
		return functionParam, true
	}

	p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(typeExpr.GetToken()), "in function parameters")

	return functionParam, false
}

func (p *Parser) functionParams(opening tokens.TokenType, closing tokens.TokenType) ([]ast.FunctionParam, bool) {
	if !p.match(opening) {
		p.Alert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), opening)
		return []ast.FunctionParam{}, false
	}

	var args []ast.FunctionParam
	if p.match(closing) {
		return args, true
	}
	var previous *ast.TypeExpr
	param, ok := p.getFunctionParam()

	success := ok
	if param.Type == nil {
		success = false
		p.Alert(&alerts.ExpectedType{}, alerts.NewSingle(param.Name))
	} else {
		previous = param.Type
	}
	args = append(args, param)
	for p.match(tokens.Comma) {
		param, ok := p.getFunctionParam()
		success = success && ok

		if param.Type == nil {
			if len(args) == 0 {
				success = false
				p.Alert(&alerts.ExpectedType{}, alerts.NewSingle(p.peek(-1)))
			} else {
				param.Type = previous
			}
		} else {
			previous = param.Type
		}

		args = append(args, param)
	}
	_, ok = p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), closing), closing)
	success = success && ok

	return args, success
}

func (p *Parser) genericParams() ([]*ast.IdentifierExpr, bool) {
	params := []*ast.IdentifierExpr{}
	if !p.match(tokens.Less) {
		return params, true
	}

	token, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "in generic parameters"), tokens.Identifier)
	success := ok
	if ok {
		params = append(params, &ast.IdentifierExpr{Name: token, ValueType: ast.Invalid})
	}

	for p.match(tokens.Comma) {
		token, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "in generic parameters"), tokens.Identifier)
		success = success && ok
		if ok {
			params = append(params, &ast.IdentifierExpr{Name: token, ValueType: ast.Invalid})
		}
	}

	_, ok = p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.Greater, "in generic parameters"), tokens.Greater)
	success = success && ok

	return params, success
}

// Prerequisite of calling this function is that you checked on peek and it was tokens.Less
func (p *Parser) tryGenericArgs() bool {
	p.context.ignoreAlerts.Push("TryGenericArgs", true)
	defer p.context.ignoreAlerts.Pop("TryGenericArgs")

	currentStart := p.current
	p.match(tokens.Less)

	p.typeExpr("in generic arguments")

	for p.match(tokens.Comma) {
		p.typeExpr("in generic arguments")
	}

	next := p.peek()

	p.disadvance(p.current - currentStart)

	if next.Type == tokens.Greater || next.Type == tokens.LeftParen {
		return true
	}
	return false
}

func (p *Parser) genericArgs() ([]*ast.TypeExpr, bool) {
	params := []*ast.TypeExpr{}
	if !p.match(tokens.Less) {
		return params, true // generic args are optional
	}

	expr := p.typeExpr("in generic arguments")
	params = append(params, expr)

	for p.match(tokens.Comma) {
		expr := p.typeExpr("in generic arguments")
		params = append(params, expr)
	}

	if p.match(tokens.Greater) {
		return params, true
	}

	p.Alert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.Greater)

	return params, false
}

func (p *Parser) functionArgs() ([]ast.Node, bool) {
	if _, ok := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.LeftParen), tokens.LeftParen); !ok {
		return nil, false
	}

	if p.match(tokens.RightParen) {
		return make([]ast.Node, 0), true
	}

	args, ok := p.expressions("in function arguments", false)
	if !ok {
		return args, false
	}
	p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.RightParen), tokens.RightParen)

	return args, true
}

func (p *Parser) functionReturns() *ast.TypeExpr {
	if !p.match(tokens.ThinArrow) {
		return nil
	}

	return p.typeExpr("in function returns")
}

func (p *Parser) identifier(typeContext string) *ast.IdentifierExpr {
	expr := p.expression()
	if expr.GetType() == ast.Identifier {
		return expr.(*ast.IdentifierExpr)
	}
	p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(expr.GetToken()), typeContext)

	return nil
}

// bool tells you if the parsing succeeded
func (p *Parser) identifiers(typeContext string, allowTrailing bool) ([]*ast.IdentifierExpr, bool) {
	idents := []*ast.IdentifierExpr{}
	ok := true

	ident := p.identifier(typeContext)
	if ident == nil {
		ok = false
	} else {
		idents = append(idents, ident)
	}

	for p.match(tokens.Comma) {
		if !p.check(tokens.Identifier) && allowTrailing {
			return idents, true
		}

		ident := p.identifier(typeContext)
		if ident == nil {
			ok = false
			continue
		}

		idents = append(idents, ident)
	}

	return idents, ok
}

// bool tells you if the parsing was successful or not
func (p *Parser) expressions(typeContext string, allowTrailing bool) ([]ast.Node, bool) {
	exprs := []ast.Node{}
	expr := p.expression()
	success := expr.GetType() != ast.NA
	if ast.IsImproper(expr, ast.NA) {
		p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(expr.GetToken()), typeContext)
	} else if expr.GetType() != ast.NA {
		exprs = append(exprs, expr)
	}
	for p.match(tokens.Comma) {
		exprStart := p.current
		expr := p.expression()
		success = success && expr.GetType() != ast.NA
		if expr.GetType() == ast.NA && allowTrailing {
			p.disadvance(p.current - exprStart)
			return exprs, success
		} else if ast.IsImproper(expr, ast.NA) {
			p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(expr.GetToken()), typeContext)
			continue
		}
		if expr.GetType() != ast.NA {
			exprs = append(exprs, expr)
		}
	}

	return exprs, success
}

func (p *Parser) identExprPairs(typeContext string, optional bool) ([]*ast.IdentifierExpr, []ast.Node, bool) {
	idents, ok := p.identifiers(typeContext, false)
	if !ok {
		return nil, nil, ok
	}

	if !p.match(tokens.Equal) {
		if optional {
			return idents, nil, ok
		}

		p.Alert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.Equal)
		return nil, nil, false
	}

	exprs, ok := p.expressions(typeContext, false)
	if !ok {
		return nil, nil, ok
	}

	return idents, exprs, true
}

func (p *Parser) keyValuePair(context string) (ast.Node, ast.Node, bool) {
	key := p.expression()
	if key.GetType() != ast.Identifier {
		p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(key.GetToken()), "as "+context)
		p.advance()
		return nil, nil, false
	}
	_, ok := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.Equal, "after "+context), tokens.Equal)
	if !ok {
		return nil, nil, false
	}

	expr := p.expression()
	if ast.IsImproper(expr, ast.NA) {
		p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(p.peek()), "as "+context+" value")
	}

	return key, expr, expr.GetType() != ast.NA
}

func (p *Parser) tryIdentifiers() bool {
	ident := p.advance()
	if ident.Type != tokens.Identifier {
		return false
	}

	for p.match(tokens.Comma) {
		ident := p.advance()
		if ident.Type != tokens.Identifier {
			return false
		}
	}

	return true
}

func (p *Parser) peekTypeVariableDecl() bool {
	currentStart := p.current
	p.context.ignoreAlerts.Push("PeekTypeVariableDecl", true)

	valid := false

	typeExpr := p.typeExpr("")
	ok := p.tryIdentifiers()
	valid = typeExpr.Name.GetType() != ast.NA && ok

	p.context.ignoreAlerts.Pop("PeekTypeVariableDecl")
	p.disadvance(p.current - currentStart)

	return valid
}

// Tells you whether the attempted parsed type is an actual TypeExpression (or Improper{Type:ast.TypeExpression}), in which case returning true, or Improper{Type:ast.NA}, in which case returning false
func (p *Parser) checkType(context string) (*ast.TypeExpr, bool) {
	currentStart := p.current
	p.context.ignoreAlerts.Push("CheckType", true)

	typeExpr := p.typeExpr("")

	p.context.ignoreAlerts.Pop("CheckType")
	valid := typeExpr != nil && !ast.IsImproper(typeExpr.Name, ast.NA)
	p.disadvance(p.current - currentStart)

	if valid {
		typeExpr = p.typeExpr(context) // so that potential alerts are not ignored
	}

	return typeExpr, valid
}

// Bool indicates whether it got a valid body or not; the success of the function
func (p *Parser) body(allowSingleSatement, allowArrow bool) ([]ast.Node, bool) {
	if !p.check(tokens.LeftBrace) && p.peek(1).Type == tokens.LeftBrace {
		p.advance()
	}
	body := make([]ast.Node, 0)
	if p.match(tokens.FatArrow) && allowArrow {
		args, ok := p.expressions("in fat arrow return arguments", false)
		if !ok {
			p.Alert(&alerts.ExpectedReturnArgs{}, alerts.NewSingle(p.peek()))
			return []ast.Node{}, false
		}
		body = []ast.Node{
			&ast.ReturnStmt{
				Token: args[0].GetToken(),
				Args:  args,
			},
		}
		return body, true
	} else if !p.check(tokens.LeftBrace) && allowSingleSatement {
		stmt := p.parseNode(p.synchronizeBody)
		if ast.IsImproperNotStatement(stmt) {
			p.Alert(&alerts.UnknownStatement{}, alerts.NewSingle(stmt.GetToken()))
			return body, false
		}
		body = []ast.Node{stmt}
		return body, true
	}
	if _, success := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.LeftBrace), tokens.LeftBrace); !success {
		return body, false
	}
	start := p.peek(-1)

	p.context.braceCounter.Increment()
	defer p.context.braceCounter.Decrement()

	for p.consumeTill("in body", start, tokens.RightBrace) {
		declaration := p.parseNode(p.synchronizeBody)
		if ast.IsImproperNotStatement(declaration) {
			p.Alert(&alerts.UnknownStatement{}, alerts.NewSingle(declaration.GetToken()))
			continue
		}
		body = append(body, declaration)
	}

	return body, true
}

func (p *Parser) limitedExpression(context string, types ...ast.NodeType) ast.Node {
	expr := p.expression()
	exprType := expr.GetType()
	if ast.IsImproper(expr, ast.NA) {
		p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(expr.GetToken()), context)
		return expr
	}
	ok := false
	for _, v := range types {
		if exprType == v {
			ok = true
		}
	}
	if !ok {
		p.Alert(&alerts.InvalidExpression{}, alerts.NewSingle(expr.GetToken()), string(exprType), context)
	}

	return expr
}

func (p *Parser) isCall(nodeType ast.NodeType) bool {
	return nodeType == ast.CallExpression ||
		nodeType == ast.MethodCallExpression ||
		nodeType == ast.NewExpession ||
		nodeType == ast.SpawnExpression
}

func (p *Parser) synchronizeBody() {
	braceCount := p.context.braceCounter.Value()
	expectedBlockCount := 0
	for !p.isAtEnd() {
		switch p.peek().Type {
		case tokens.Fn:
			p.advance()
			if p.peek().Type != tokens.LeftParen {
				p.disadvance()
				return
			}
		case tokens.LeftBrace:
			expectedBlockCount++
		case tokens.RightBrace:
			if expectedBlockCount == 0 {
				if braceCount == 0 {
					p.advance()
					continue
				}
				return
			}

			expectedBlockCount--
		case tokens.Entity:
			if p.peek(1).Type == tokens.Identifier && p.peek(2).Type == tokens.LeftBrace {
				return
			}
		case tokens.Let, tokens.Pub, tokens.Const, tokens.Class, tokens.Alias:
			return
		}

		p.advance()
	}
}

func (p *Parser) synchronizeDeclBody() {
	expectedBlockCount := 0
	for !p.isAtEnd() {
		switch p.peek().Type {
		case tokens.Fn:
			p.advance()
			if p.peek().Type != tokens.LeftParen {
				p.disadvance()
				return
			}
		case tokens.LeftBrace:
			expectedBlockCount++
		case tokens.RightBrace:
			if expectedBlockCount == 0 {
				return
			}

			expectedBlockCount--
		case tokens.Entity:
			if p.peek(1).Type == tokens.Identifier && p.peek(2).Type == tokens.LeftBrace {
				return
			}
		case tokens.Identifier:
			current := p.current
			p.advance()
			p.context.ignoreAlerts.Push("SynchronizeDeclBody", true)
			if _, ok := p.functionParams(tokens.LeftParen, tokens.RightParen); ok && p.check(tokens.LeftBrace) {
				p.disadvance(p.current - current)
				p.context.ignoreAlerts.Pop("SynchronizeDeclBody")
				return
			}
			p.disadvance(p.current - current)
			node := p.variableDeclaration(false)
			p.context.ignoreAlerts.Pop("SynchronizeDeclBody")
			if !ast.IsImproper(node, ast.NA) {
				p.disadvance(p.current - current)
				return
			}
		case tokens.Let, tokens.Pub, tokens.Const, tokens.Class, tokens.Alias, tokens.New, tokens.Spawn, tokens.Destroy:
			return
		}

		p.advance()
	}
}

func (p *Parser) synchronizeMatchBody() {
	expectedBlockCount := 0
	for !p.isAtEnd() {
		switch p.peek().Type {
		case tokens.LeftBrace:
			expectedBlockCount++
		case tokens.RightBrace:
			if expectedBlockCount == 0 {
				return
			}

			expectedBlockCount--
		default:
			current := p.current
			p.context.ignoreAlerts.Push("SynchronizeMatchBody", true)
			_, ok := p.expressions("", false)
			p.context.ignoreAlerts.Pop("SynchronizeMatchBody")
			if ok && p.check(tokens.FatArrow) {
				p.disadvance(p.current - current)
				return
			}
			p.disadvance(p.current - current)
		}

		p.advance()
	}
}

// Checks if there is a discrepancy between the line location of tokenStart and tokenEnd
//
// Returns true if there was a discrepancy. AllowNewLine is false by default.
func (p *Parser) coherencyCheck(tokenStart, tokenEnd tokens.Token, allowNewLine ...bool) bool {
	diffTolerance := 0
	if allowNewLine != nil && allowNewLine[0] {
		diffTolerance = 1
	}
	if tokenEnd.Line-tokenStart.Line > diffTolerance {
		p.Alert(&alerts.SyntaxIncoherency{}, alerts.NewMulti(tokenStart, tokenEnd), tokenEnd.Lexeme, tokenStart.Lexeme, diffTolerance == 1)
		return false
	}

	return true
}
