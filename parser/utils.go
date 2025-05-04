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

func (p *Parser) getFunctionParam(previous *ast.TypeExpr, closing tokens.TokenType) ast.FunctionParam {
	functionParam := ast.FunctionParam{}

	typeExpr := p.typeExpr("in function parameters")
	peekType := p.peek().Type

	if peekType == tokens.Identifier {
		functionParam.Name = p.advance()
		functionParam.Type = typeExpr
		return functionParam
	}
	if typeExpr.Name.GetType() == ast.Identifier {
		functionParam.Name = typeExpr.Name.GetToken()
		functionParam.Type = nil
	} else {
		p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(typeExpr.GetToken()), "in function parameters")
	}
	return functionParam
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
	param := p.getFunctionParam(nil, closing)
	if param.Type == nil {
		p.Alert(&alerts.ExpectedType{}, alerts.NewSingle(param.Name))
	} else {
		previous = param.Type
	}
	args = append(args, param)
	for p.match(tokens.Comma) {
		param := p.getFunctionParam(previous, closing)
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

	p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.Greater), tokens.Greater)

	return params
}

func (p *Parser) genericArgs() ([]*ast.TypeExpr, bool) {
	currentStart := p.current
	params := []*ast.TypeExpr{}
	if !p.match(tokens.Less) {
		return params, false
	}

	params = append(params, p.typeExpr("in generic arguments"))

	for p.match(tokens.Comma) {
		params = append(params, p.typeExpr("in generic arguments"))
	}

	if !p.match(tokens.Greater) {
		p.disadvance(p.current - currentStart)
		return params, false
	}

	return params, true
}

func (p *Parser) functionArgs() []ast.Node {
	if _, ok := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.LeftParen), tokens.LeftParen); !ok {
		return nil
	}

	var args []ast.Node
	if p.match(tokens.RightParen) {
		args = make([]ast.Node, 0)
	} else {
		arg := p.expression()
		args = append(args, arg)
		for p.match(tokens.Comma) {
			arg := p.expression()
			args = append(args, arg)
		}
		p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.RightParen), tokens.RightParen)
	}

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
	success := true
	if expr.GetType() == ast.NA {
		p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(p.peek()), typeContext)
		success = false
	} else {
		exprs = append(exprs, expr)
	}
	for p.match(tokens.Comma) {
		exprStart := p.current
		expr := p.expression()
		if expr.GetType() == ast.NA && allowTrailing {
			p.disadvance(p.current - exprStart)
			return exprs, success
		} else if expr.GetType() == ast.NA {
			p.Alert(&alerts.ExpectedExpression{}, alerts.NewSingle(p.peek()), typeContext)
			success = false
			continue
		}

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

	exprs, ok := p.expressions(typeContext, false)
	if !ok {
		return nil, nil, ok
	}

	return idents, exprs, true
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
		if p.context.FunctionReturns.Top().Item > 0 {
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
		} else {
			body = []ast.Node{p.statement()}
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

	for !p.match(tokens.RightBrace) {
		if p.isAtEnd() {
			p.Alert(&alerts.ExpectedSymbol{}, alerts.NewMulti(start, p.peek(-1)), tokens.RightBrace)
			return body, false
		}

		declaration := p.declaration()
		if declaration.GetType() != ast.NA {
			body = append(body, declaration)
		}
	}

	return body, true
}
