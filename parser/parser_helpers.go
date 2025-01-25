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
		Operator:  tokens.NewToken(tokenType, lexeme, "", operator.Position),
		Right:     right,
		ValueType: valueType,
	}
}

// Checks if the value type is expected to be a fixedpoint
func IsFx(valueType ast.PrimitiveValueType) bool {
	return valueType == ast.FixedPoint || valueType == ast.Fixed || valueType == ast.Radian || valueType == ast.Degree
}

func (p *Parser) getOp(opEqual tokens.Token) tokens.Token {
	switch opEqual.Type {
	case tokens.PlusEqual:
		return tokens.NewToken(tokens.Plus, "+", opEqual.Literal, opEqual.Position)
	case tokens.MinusEqual:
		return tokens.NewToken(tokens.Minus, "-", opEqual.Literal, opEqual.Position)
	case tokens.SlashEqual:
		return tokens.NewToken(tokens.Slash, "/", opEqual.Literal, opEqual.Position)
	case tokens.StarEqual:
		return tokens.NewToken(tokens.Star, "*", opEqual.Literal, opEqual.Position)
	case tokens.CaretEqual:
		return tokens.NewToken(tokens.Caret, "^", opEqual.Literal, opEqual.Position)
	case tokens.ModuloEqual:
		return tokens.NewToken(tokens.Modulo, "%", opEqual.Literal, opEqual.Position)
	default:
		return tokens.Token{}
	}
}

func (p *Parser) getParam(previous *ast.TypeExpr, closing tokens.TokenType) ast.Param {
	typ := p.Type()
	peekType := p.peek().Type

	if peekType == tokens.Identifier {
		return ast.Param{Name: p.advance(), Type: typ}
	} else if peekType == tokens.Comma || peekType == closing {
		if typ.Name.GetType() == ast.Identifier && previous != nil {
			return ast.Param{Name: typ.Name.GetToken()}
		} else {
			p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(typ.GetToken()))
			return ast.Param{Name: typ.GetToken()}
		}
	} else {
		p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()))
		return ast.Param{Name: p.advance()}
	}
}

func (p *Parser) parameters(opening tokens.TokenType, closing tokens.TokenType) []ast.Param {
	if !p.match(opening) {
		p.Alert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), opening)
		return []ast.Param{}
	}

	open := p.peek(-1)

	var args []ast.Param
	if p.match(closing) {
		args = make([]ast.Param, 0)
	} else {

		var previous *ast.TypeExpr
		param := p.getParam(nil, closing)
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
		for p.match(tokens.Comma) {
			param := p.getParam(previous, closing)
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
		p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewMulti(open, p.peek()), closing), closing)
	}

	return args
}

func (p *Parser) genericParameters() []*ast.IdentifierExpr {
	params := []*ast.IdentifierExpr{}
	if !p.match(tokens.Less) {
		return params
	}

	token := p.advance()
	if token.Type != tokens.Identifier {
		p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(token), "in generic parameters")
	} else {
		params = append(params, &ast.IdentifierExpr{Name: token, ValueType: ast.Invalid})
	}

	for p.match(tokens.Comma) {
		token := p.advance()
		if token.Type != tokens.Identifier {
			p.Alert(&alerts.ExpectedType{}, alerts.NewSingle(token))
		} else {
			params = append(params, &ast.IdentifierExpr{Name: token, ValueType: ast.Invalid})
		}
	}

	p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.Greater), tokens.Greater)

	return params
}

func (p *Parser) genericArguments() ([]*ast.TypeExpr, bool) {
	currentStart := p.current
	params := []*ast.TypeExpr{}
	if !p.match(tokens.Less) {
		return params, false
	}

	params = append(params, p.Type())

	for p.match(tokens.Comma) {
		params = append(params, p.Type())
	}

	if !p.match(tokens.Greater) {
		p.disadvance(p.current - currentStart)
		return params, false
	}

	return params, true
}

func (p *Parser) arguments() []ast.Node {
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

func (p *Parser) returnings() []*ast.TypeExpr {
	ret := make([]*ast.TypeExpr, 0)

	if !p.match(tokens.ThinArrow) {
		return ret
	}

	if p.match(tokens.LeftParen) {
		if p.match(tokens.RightParen) {
			return ret
		}
		ret = append(ret, p.Type())

		for p.match(tokens.Comma) {
			ret = append(ret, p.Type())
		}

		p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.RightParen), tokens.RightParen)
	} else {
		ret = append(ret, p.Type())
	}

	return ret
}

func (p *Parser) getIdentifiers() ([]tokens.Token, bool) {
	names := []tokens.Token{}
	ident, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek())), tokens.Identifier)
	if !ok {
		return names, false
	}
	names = append(names, ident)
	for p.match(tokens.Comma) {
		ident, ok = p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek())), tokens.Identifier)
		if !ok {
			return names, false
		}
		names = append(names, ident)
	}

	return names, true
}

func (p *Parser) getExpressions() ([]ast.Node, bool) {
	exprs := []ast.Node{}
	expr := p.expression()
	if expr.GetType() == ast.NA && ast.ImproperToNodeType(expr) == ast.NA {
		return exprs, false
	}
	exprs = append(exprs, expr)
	for p.match(tokens.Comma) {
		expr := p.expression()
		if expr.GetType() == ast.NA && ast.ImproperToNodeType(expr) == ast.NA {
			return exprs, false
		}
		exprs = append(exprs, expr)
	}

	return exprs, true
}

func (p *Parser) TypeAndIdentifier() (*ast.TypeExpr, ast.Node) {
	typ := p.Type()

	ident := p.advance()
	if ident.Type != tokens.Identifier {
		p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(ident))
	}

	return typ, &ast.IdentifierExpr{Name: ident, ValueType: ast.Invalid}
}

func (p *Parser) CheckType() bool {
	currentStart := p.current
	p.Context.IgnoreAlerts.Push("CheckType", true)

	typ := p.Type()

	p.Context.IgnoreAlerts.Pop("CheckType")
	p.disadvance(p.current - currentStart)

	return typ.Name.GetType() != ast.NA
}

func (p *Parser) getBody() ([]ast.Node, bool) {
	body := make([]ast.Node, 0)
	if p.match(tokens.FatArrow) {
		args, ok := p.returnArgs()
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
	} else if !p.check(tokens.LeftBrace) {
		body = []ast.Node{p.statement()}
		return body, true
	}
	if _, success := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.LeftBrace), tokens.LeftBrace); !success {
		return body, false
	}
	start := p.peek(-1)

	for !p.match(tokens.RightBrace) {
		if p.peek().Type == tokens.Eof {
			p.Alert(&alerts.ExpectedSymbol{}, alerts.NewMulti(start, p.peek(-1)), tokens.RightBrace)
			return body, false
		}

		statement := p.statement()
		if statement.GetType() != ast.NA {
			body = append(body, statement)
		}
	}

	return body, true
}
