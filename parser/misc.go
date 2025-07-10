package parser

import (
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/tokens"
	"strings"
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
		p.Alert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), opening, "in function parameters")
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
	_, ok = p.alertSingleConsume(&alerts.ExpectedSymbol{}, closing, "in function parameters")
	success = success && ok

	return args, success
}

// Prerequisite of calling this function is that you checked on peek and it was tokens.Less
func (p *Parser) tryGenericParams(offset ...int) bool {
	p.context.ignoreAlerts.Push("TryGenericParams", true)
	defer p.context.ignoreAlerts.Pop("TryGenericParams")

	off := 0
	if offset != nil {
		off = offset[0]
	}
	currentStart := p.current
	if off != 0 {
		p.advance(off)
	}
	if !p.match(tokens.Less) {
		return false
	}

	p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "in generic parameters"), tokens.Identifier)

	for p.match(tokens.Comma) {
		p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "in generic parameters"), tokens.Identifier)
	}

	next := p.peek()

	p.disadvance(p.current - currentStart)
	if next.Type == tokens.Greater || next.Type == tokens.LeftParen {
		return true
	}
	return false
}

func (p *Parser) genericParams() ([]*ast.IdentifierExpr, bool) {
	params := []*ast.IdentifierExpr{}
	if !p.match(tokens.Less) {
		return params, true
	}

	token, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "in generic parameters"), tokens.Identifier)
	success := ok
	if ok {
		params = append(params, &ast.IdentifierExpr{Name: token})
	}

	for p.match(tokens.Comma) {
		token, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "in generic parameters"), tokens.Identifier)
		success = success && ok
		if ok {
			params = append(params, &ast.IdentifierExpr{Name: token})
		}
	}

	_, ok = p.alertSingleConsume(&alerts.ExpectedSymbol{}, tokens.Greater, "in generic parameters")
	success = success && ok

	return params, success
}

// Prerequisite of calling this function is that you checked on peek and it was tokens.Less
func (p *Parser) tryGenericArgs() bool {
	p.context.ignoreAlerts.Push("TryGenericArgs", true)
	defer p.context.ignoreAlerts.Pop("TryGenericArgs")

	currentStart := p.current
	if !p.match(tokens.Less) {
		return false
	}

	p.typeExpr("in generic arguments")
	for p.match(tokens.Comma) {
		p.typeExpr("in generic arguments")
	}
	next := p.peek()

	p.disadvance(p.current - currentStart)
	return next.Type == tokens.Greater
}

func (p *Parser) genericArgs() ([]*ast.TypeExpr, bool) {
	params := []*ast.TypeExpr{}
	if !p.match(tokens.Less) {
		return params, true // generic args are optional
	}

	success := true

	expr := p.typeExpr("in generic arguments")
	if ast.IsImproper(expr.Name, ast.NA) {
		p.AlertSingle(&alerts.ExpectedCallArgs{}, expr.GetToken())
		success = false
	} else {
		params = append(params, expr)
	}

	for p.match(tokens.Comma) {
		expr := p.typeExpr("in generic arguments")
		if ast.IsImproper(expr.Name, ast.NA) {
			p.AlertSingle(&alerts.ExpectedCallArgs{}, expr.GetToken())
			success = false
		} else {
			params = append(params, expr)
		}
	}

	if p.match(tokens.Greater) {
		return params, true
	}

	p.Alert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.Greater)

	return params, success
}

func (p *Parser) functionArgs() ([]ast.Node, bool) {
	if _, ok := p.alertSingleConsume(&alerts.ExpectedSymbol{}, tokens.LeftParen, "in function arguments"); !ok {
		return nil, false
	}

	if p.match(tokens.RightParen) {
		return make([]ast.Node, 0), true
	}

	args, ok := p.expressions("in function arguments", false)
	if !ok {
		return args, false
	}
	p.alertSingleConsume(&alerts.ExpectedSymbol{}, tokens.RightParen)

	return args, true
}

func (p *Parser) functionReturns() ([]*ast.TypeExpr, bool) {
	var returns []*ast.TypeExpr
	if !p.match(tokens.ThinArrow) {
		return returns, true
	}

	if p.match(tokens.LeftParen) {
		success := true
		typ := p.typeExpr("in function returns")
		if typ.GetType() != ast.NA {
			returns = append(returns, typ)
		} else {
			success = false
		}

		for p.match(tokens.Comma) {
			typ := p.typeExpr("in function returns")
			if typ.GetType() != ast.NA {
				returns = append(returns, typ)
			} else {
				success = false
			}
		}

		p.alertSingleConsume(&alerts.ExpectedSymbol{}, tokens.RightParen, "in function return types")
		return returns, success
	}

	returns = append(returns, p.typeExpr("in function return types"))
	return returns, returns[0].GetType() != ast.NA
}

func (p *Parser) identifier(typeContext string) *ast.IdentifierExpr {
	if p.peek().Type != tokens.Identifier {
		expr := p.expression()
		p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(expr.GetToken()), typeContext)
		return nil
	}
	return &ast.IdentifierExpr{
		Name: p.advance(),
		Type: ast.Other,
	}
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
	if previous := p.peek(-1); previous.Line != p.peek().Line {
		p.AlertSingle(&alerts.ExpectedIdentifier{}, previous, typeContext)
		return nil, nil, false
	}
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
	equal := p.peek(-1)
	if p.peek().Line != equal.Line {
		p.AlertSingle(&alerts.ExpectedExpression{}, equal, typeContext)
		return idents, []ast.Node{}, false
	}

	exprs, ok := p.expressions(typeContext, false)
	if !ok {
		return nil, nil, ok
	}

	return idents, exprs, true
}

func (p *Parser) keyValuePair(isMap bool, context string) (ast.Node, ast.Node, bool) {
	key := p.expression()
	condition := key.GetType() != ast.Identifier
	if isMap {
		condition = key.GetType() != ast.LiteralExpression || key.GetToken().Type != tokens.String
	}
	if condition {
		if isMap {
			p.Alert(&alerts.InvalidMapKey{}, alerts.NewSingle(key.GetToken()))
		} else {
			p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(key.GetToken()), "as "+context)
		}
		return nil, nil, false
	}
	_, ok := p.alertSingleConsume(&alerts.ExpectedSymbol{}, tokens.Equal, "after "+context)
	if !ok {
		return nil, nil, false
	}

	expr := p.expression()
	if ast.IsImproper(expr, ast.NA) {
		p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(p.peek()), "as "+context+" value")
	}

	return key, expr, expr.GetType() != ast.NA
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
func (p *Parser) body(allowSingleSatement, allowArrow bool) (ast.Body, bool) {
	if !p.check(tokens.LeftBrace) && p.peek(1).Type == tokens.LeftBrace {
		p.advance()
	}
	body := ast.NewBody()
	if p.match(tokens.FatArrow) && allowArrow {
		args, ok := p.expressions("in fat arrow return arguments", false)
		if !ok {
			p.Alert(&alerts.ExpectedReturnArgs{}, alerts.NewSingle(p.peek()))
			return ast.Body{}, false
		}
		body.Append(&ast.ReturnStmt{
			Token: args[0].GetToken(),
			Args:  args,
		})
		return body, true
	} else if !p.check(tokens.LeftBrace) && allowSingleSatement {
		stmt := p.parseNode(p.synchronizeBody)
		if ast.IsImproperNotStatement(stmt) {
			p.Alert(&alerts.UnknownStatement{}, alerts.NewSingle(stmt.GetToken()))
			return body, false
		}
		body.Append(stmt)
		return body, true
	}
	if _, success := p.alertSingleConsume(&alerts.ExpectedSymbol{}, tokens.LeftBrace); !success {
		return body, false
	}
	start := p.peek(-1)

	for p.consumeTill("in body", start, tokens.RightBrace) {
		declaration := p.parseNode(p.synchronizeBody)
		if ast.IsImproperNotStatement(declaration) {
			p.Alert(&alerts.UnknownStatement{}, alerts.NewSingle(declaration.GetToken()))
			continue
		}
		body.Append(declaration)
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

// this is used only for maps, lists and structs
func (p *Parser) sync(syncPoints ...tokens.TokenType) bool {
	expectedBlockCount := 0
	for !p.isAtEnd() {
		peekType := p.peek().Type
		if syncPoints != nil && expectedBlockCount == 0 {
			for _, v := range syncPoints {
				if peekType == v {
					return true
				}
			}
		}

		switch p.peek().Type {
		case tokens.LeftBrace:
			expectedBlockCount++
		case tokens.RightBrace:
			expectedBlockCount--
		}

		p.advance()
	}
	return false
}

func (p *Parser) synchronizeBody() {
	defer func() {
		p.context.syncedToken = p.peek()
	}()
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
				if p.context.syncedToken == p.peek() {
					break
				}
				return
			}

			expectedBlockCount--
		case tokens.Entity:
			if p.peek(1).Type == tokens.Identifier && p.peek(2).Type == tokens.LeftBrace {
				if p.context.syncedToken == p.peek() {
					break
				}
				return
			}
		case tokens.Let, tokens.Pub, tokens.Const, tokens.Class, tokens.Alias, tokens.Repeat, tokens.For, tokens.Destroy, tokens.Spawn, tokens.New:
			return
		case tokens.If:
			if p.peek(-1).Type != tokens.Else {
				return
			}
		case tokens.Match:
			if p.peek(-1).Type != tokens.Comma && p.peek(-1).Type != tokens.Equal {
				return
			}
		}

		p.advance()
	}
}

func (p *Parser) synchronizeDeclBody() {
	defer func() {
		p.context.syncedToken = p.peek()
	}()
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
				if p.context.syncedToken == p.peek() {
					break
				}
				return
			}

			expectedBlockCount--
		case tokens.Let, tokens.Pub, tokens.Const, tokens.Class, tokens.Alias, tokens.Repeat, tokens.For, tokens.Destroy, tokens.Spawn, tokens.New:
			return
		case tokens.If:
			if p.peek(-1).Type != tokens.Else {
				return
			}
		case tokens.Match:
			if p.peek(-1).Type != tokens.Comma && p.peek(-1).Type != tokens.Equal {
				return
			}
		case tokens.Identifier:
			peek := p.peek().Lexeme
			switch peek {
			case "Update", "WeaponCollision", "PlayerCollision", "WallCollision":
				return
			}
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

// Combines the locations as well as the lexemes and literals of the n tokens. It is assumed that you use this function to create an already existing token type and also that n is more than 1.
//
// It can fail when one of the given tokens is not in the same line as the previous ones. In this case it will disadvance back to the start.
//
// E.g. combining tokens '>', '>' and '=' would result with a token of '>>=' with the appropriate location. The type is of course 'RightshiftEqual', which you have to give.
// If it fails at '=', it will still give whatever it combined. In this case '>>'.
func (p *Parser) combineTokens(tokenType tokens.TokenType, n int) (tokens.Token, bool) {
	newToken := p.advance()
	newToken.Type = tokenType
	i := 1
	for i < n {
		next := p.advance()
		if newToken.Line != next.Line {
			p.disadvance(i + 1)
			return newToken, false
		}
		newToken.Column.Start = min(newToken.Column.Start, next.Column.Start)
		newToken.Column.End = max(newToken.Column.End, next.Column.End)
		newToken.Lexeme = strings.Join([]string{newToken.Lexeme, next.Lexeme}, "")
		i++
	}

	return newToken, true
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
		p.AlertMulti(&alerts.SyntaxIncoherency{}, tokenStart, tokenEnd, tokenEnd.Lexeme, tokenStart.Lexeme, diffTolerance == 1)
		return false
	}

	return true
}
