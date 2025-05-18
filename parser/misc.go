package parser

import (
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/tokens"
)

// Checks if the value type is expected to be a fixedpoint
func IsFx(valueType ast.PrimitiveValueType) bool {
	return valueType == ast.FixedPoint || valueType == ast.Fixed || valueType == ast.Radian || valueType == ast.Degree
}

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

func (p *Parser) functionParams(opening tokens.TokenType, closing tokens.TokenType) []ast.FunctionParam {
	if !p.match(opening) {
		p.Alert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), opening)
		return []ast.FunctionParam{}
	}

	var args []ast.FunctionParam
	if p.match(closing) {
		return args
	}
	var previous *ast.TypeExpr
	param, _ := p.getFunctionParam()

	if param.Type == nil {
		p.Alert(&alerts.ExpectedType{}, alerts.NewSingle(param.Name))
	} else {
		previous = param.Type
	}
	args = append(args, param)
	for p.match(tokens.Comma) {
		param, _ := p.getFunctionParam()

		if param.Type == nil {
			if len(args) == 0 {
				p.Alert(&alerts.ExpectedType{}, alerts.NewSingle(p.peek(-1)))
			} else {
				param.Type = previous
			}
		} else {
			previous = param.Type
		}

		args = append(args, param)
	}
	p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), closing), closing)

	return args
}

func (p *Parser) genericParams() []*ast.IdentifierExpr {
	params := []*ast.IdentifierExpr{}
	if !p.match(tokens.Less) {
		return params
	}

	token, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "in generic parameters"), tokens.Identifier)
	if ok {
		params = append(params, &ast.IdentifierExpr{Name: token, ValueType: ast.Invalid})
	}

	for p.match(tokens.Comma) {
		token, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "in generic parameters"), tokens.Identifier)
		if ok {
			params = append(params, &ast.IdentifierExpr{Name: token, ValueType: ast.Invalid})
		}
	}

	p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.Greater, "in generic parameters"), tokens.Greater)

	return params
}

// Prerequisite of calling this function is that you checked on peek and it was tokens.Less
func (p *Parser) tryGenericArgs() bool {
	p.context.IgnoreAlerts.Push("TryGenericArgs", true)
	defer func() {
		p.context.IgnoreAlerts.Pop("TryGenericArgs")
	}()

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
		return nil, nil, ok
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
	p.context.IgnoreAlerts.Push("PeekTypeVariableDecl", true)

	valid := false

	typeExpr := p.typeExpr("")
	ok := p.tryIdentifiers()
	valid = typeExpr.Name.GetType() != ast.NA && ok

	p.context.IgnoreAlerts.Pop("PeekTypeVariableDecl")
	p.disadvance(p.current - currentStart)

	return valid
}

func (p *Parser) checkType() (*ast.TypeExpr, bool) {
	currentStart := p.current
	p.context.IgnoreAlerts.Push("CheckType", true)

	typeExpr := p.typeExpr("")

	p.context.IgnoreAlerts.Pop("CheckType")
	valid := typeExpr != nil && typeExpr.Name.GetType() != ast.NA
	if !valid {
		p.disadvance(p.current - currentStart)
	}

	return typeExpr, valid
}

// Bool indicates whether it got a valid body or not; the success of the function
func (p *Parser) body(allowSingleSatement, allowArrow bool) ([]ast.Node, bool) {
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
		body = []ast.Node{p.statement()}
		return body, true
	}
	if _, success := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.LeftBrace), tokens.LeftBrace); !success {
		return body, false
	}
	start := p.peek(-1)

	for p.doesntEndWith("in body", start, tokens.RightBrace) {

		declaration := p.declaration()
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

func (p *Parser) synchronize() {
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
				if p.context.BraceEntries.Count() != 0 {
					p.context.BraceEntries.Pop("Brace")
					p.advance()
					continue
				}
				return
			}

			expectedBlockCount--
		case tokens.Entity, tokens.Let, tokens.Pub, tokens.Const, tokens.Class, tokens.Alias:
			return
		}

		p.advance()
	}
}
