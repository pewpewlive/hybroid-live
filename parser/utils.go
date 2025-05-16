package parser

import (
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/tokens"
)

// Creates a BinaryExpr
func (p *Parser) createBinExpr(left ast.Node, operator tokens.Token, tokenType tokens.TokenType, lexeme string, right ast.Node) ast.Node {
	valueType := p.determineValueType(left, right)
	return &ast.BinaryExpr{
		Left:      left,
		Operator:  tokens.NewToken(tokenType, lexeme, "", operator.Location),
		Right:     right,
		ValueType: valueType,
	}
}

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
		p.coherencyFailed(string(functionParam.Name.Lexeme), typeExpr.GetToken(), functionParam.Name)
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
	if !p.match2(opening) {
		p.Alert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), opening)
		return []ast.FunctionParam{}
	}

	var args []ast.FunctionParam
	if p.match2(closing) {
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
	for p.match2(tokens.Comma) {
		comma := p.peek(-1)
		param, ok := p.getFunctionParam()
		if !ok {
			p.coherencyFailed("function parameter", comma, param.Name, true)
		}

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
	if !p.match2(tokens.Less) {
		return params
	}

	token, ok := p.consume2(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "in generic parameters"), tokens.Identifier)
	if ok {
		params = append(params, &ast.IdentifierExpr{Name: token, ValueType: ast.Invalid})
	}

	for p.match2(tokens.Comma) {
		token, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "in generic parameters"), tokens.Identifier)
		if ok {
			p.coherencyFailed("generic parameter", p.peek(-2), token, true)

			params = append(params, &ast.IdentifierExpr{Name: token, ValueType: ast.Invalid})
		}
	}

	p.consume2(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.Greater, "in generic parameters"), tokens.Greater)

	return params
}

func (p *Parser) genericArgs() ([]*ast.TypeExpr, bool) {
	currentStart := p.current
	params := []*ast.TypeExpr{}
	if !p.match2(tokens.Less) {
		return params, false
	}

	startToken := p.peek(-1)

	expr := p.typeExpr("in generic arguments")
	if expr.GetType() != ast.NA {
		p.coherencyFailed("generic argument", startToken, expr.GetToken())
	}
	params = append(params, expr)

	for p.match2(tokens.Comma) {
		expr := p.typeExpr("in generic arguments")
		if expr.GetType() != ast.NA {
			p.coherencyFailed("generic argument", startToken, expr.GetToken(), true)
		}
		params = append(params, expr)
	}

	if !p.match2(tokens.Greater) {
		p.disadvance(p.current - currentStart)
		return params, false
	}

	return params, true
}

func (p *Parser) functionArgs() []ast.Node {
	if _, ok := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.LeftParen), tokens.LeftParen); !ok {
		return nil
	}

	if p.match(tokens.RightParen) {
		return make([]ast.Node, 0)
	}

	args, _ := p.expressions("in function arguments", false, true)
	p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.RightParen), tokens.RightParen)

	return args
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

	for p.match2(tokens.Comma) {
		if !p.check(tokens.Identifier) && allowTrailing {
			return idents, true
		}

		ident := p.identifier(typeContext)
		if ident == nil {
			ok = false
			continue
		}

		p.coherencyFailed("identifier", p.peek(-2), ident.Name, true)
		idents = append(idents, ident)
	}

	return idents, ok
}

// bool tells you if the parsing was successful or not
func (p *Parser) expressions(typeContext string, allowTrailing bool, allowNewLine bool) ([]ast.Node, bool) {
	previousToken := p.peek(-1)
	exprs := []ast.Node{}
	expr := p.expression()
	success := true
	if expr.GetType() == ast.NA {
		p.Alert(&alerts.ExpectedExpression{}, alerts.NewMulti(p.peek(-1), p.peek()), typeContext)
		success = false
	} else {
		p.coherencyFailed("expression", previousToken, expr.GetToken(), allowNewLine)
		exprs = append(exprs, expr)
	}
	for p.match2(tokens.Comma) {
		comma := p.peek(-1)
		exprStart := p.current
		expr := p.expression()
		if expr.GetType() == ast.NA && allowTrailing {
			p.disadvance(p.current - exprStart)
			return exprs, success
		} else if expr.GetType() == ast.NA {
			p.Alert(&alerts.ExpectedExpression{}, alerts.NewMulti(p.peek(-1), p.peek()), typeContext)
			success = false
			continue
		}

		p.coherencyFailed("expression", comma, expr.GetToken(), true)
		exprs = append(exprs, expr)
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

	exprs, ok := p.expressions(typeContext, false, false)
	if !ok {
		return nil, nil, ok
	}

	return idents, exprs, true
}

func (p *Parser) peekTypeVariableDecl() bool {
	currentStart := p.current
	p.context.IgnoreAlerts.Push("CheckType", true)

	typeExpr := p.typeExpr("")

	p.context.IgnoreAlerts.Pop("CheckType")
	expr := p.expression()
	valid := typeExpr != nil && typeExpr.Name.GetType() != ast.NA && expr.GetType() == ast.Identifier
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
		if p.coherencyFailed("body", p.peek(-2), p.peek(-1), true) {
			return body, false
		}
		args, ok := p.expressions("in fat arrow return arguments", false, false)
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
	if p.coherencyFailed("body", p.peek(-2), start) {
		return body, false
	}

	for !p.match(tokens.RightBrace) {
		if p.isAtEnd() {
			p.Alert(&alerts.ExpectedSymbol{}, alerts.NewMulti(start, p.peek(-1)), tokens.RightBrace)
			return body, false
		}

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
	if exprType == ast.NA {
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

// Checks if there is a discrepancy between the line location of tokenStart and tokenEnd
//
// Returns true if there was a discrepancy. AllowNewLine is false by default.
func (p *Parser) coherencyFailed(parsedSection string, tokenStart, tokenEnd tokens.Token, allowNewLine ...bool) bool {
	diffTolerance := 0
	if allowNewLine != nil && allowNewLine[0] {
		diffTolerance = 1
	}
	if tokenEnd.Line-tokenStart.Line > diffTolerance {
		p.Alert(&alerts.SyntaxIncoherency{}, alerts.NewMulti(tokenStart, tokenEnd), parsedSection, diffTolerance == 1)
		return true
	}

	return false
}
