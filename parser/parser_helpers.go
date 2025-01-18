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
		Operator:  tokens.NewToken(tokenType, lexeme, "", operator.TokenLocation),
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
		return tokens.NewToken(tokens.Plus, "+", opEqual.Literal, opEqual.TokenLocation)
	case tokens.MinusEqual:
		return tokens.NewToken(tokens.Minus, "-", opEqual.Literal, opEqual.TokenLocation)
	case tokens.SlashEqual:
		return tokens.NewToken(tokens.Slash, "/", opEqual.Literal, opEqual.TokenLocation)
	case tokens.StarEqual:
		return tokens.NewToken(tokens.Star, "*", opEqual.Literal, opEqual.TokenLocation)
	case tokens.CaretEqual:
		return tokens.NewToken(tokens.Caret, "^", opEqual.Literal, opEqual.TokenLocation)
	case tokens.ModuloEqual:
		return tokens.NewToken(tokens.Modulo, "%", opEqual.Literal, opEqual.TokenLocation)
	default:
		return tokens.Token{}
	}
}

func (p *Parser) getParam(closing tokens.TokenType) ast.Param {
	typ := p.Type()
	peekType := p.peek().Type

	if peekType == tokens.Identifier {
		return ast.Param{Name: p.advance(), Type: typ}
	} else if peekType == tokens.Comma || peekType == closing {
		if typ.Name.GetType() == ast.Identifier && (typ.WrappedType != nil || typ.Fields != nil || typ.Params != nil || typ.Returns != nil) {
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
		p.Alert(&alerts.ExpectedOpeningMark{}, alerts.NewSingle(p.peek()), string(opening))
		return []ast.Param{}
	}

	open := p.peek(-1)

	var args []ast.Param
	if p.match(closing) {
		args = make([]ast.Param, 0)
	} else {
		var previous *ast.TypeExpr
		param := p.getParam(closing)
		if param.Type == nil {
			if len(args) == 0 {
				p.Alert(&alerts.ExpectedType{}, alerts.NewSingle(p.peek(-1))) //param.Name, "parameter need to be declared with a type before the name")
			} else {
				param.Type = previous
			}
		} else {
			previous = param.Type
		}
		args = append(args, param)
		for p.match(tokens.Comma) {
			param := p.getParam(closing)
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
		p.consume(p.NewAlert(&alerts.ExpectedEnclosingMark{}, alerts.NewMulti(open, p.peek()), string(closing)), closing)
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
			//p.error(token, "expected type identifier in generic parameters")
		} else {
			params = append(params, &ast.IdentifierExpr{Name: token, ValueType: ast.Invalid})
		}
	}

	p.consume(p.NewAlert(&alerts.ExpectedEnclosingMark{}, alerts.NewSingle(p.peek()), string(tokens.Greater)), tokens.Greater)
	//p.consumeOld("expected '>' in generic parameters", tokens.Greater)

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
	if _, ok := p.consume(p.NewAlert(&alerts.ExpectedOpeningMark{}, alerts.NewSingle(p.peek()), tokens.LeftParen), tokens.LeftParen); !ok {
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
		p.consume(p.NewAlert(&alerts.ExpectedEnclosingMark{}, alerts.NewSingle(p.peek()), string(tokens.RightParen)), tokens.RightParen)
		//p.consumeOld("expected closing paren after arguments", tokens.RightParen)
	}

	return args
}

func (p *Parser) returnings() []*ast.TypeExpr {
	ret := make([]*ast.TypeExpr, 0)

	if !p.match(tokens.ThinArrow) {
		p.Alert(&alerts.ExpectedReturnArrow{}, alerts.NewSingle(p.peek()))
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

		p.consume(p.NewAlert(&alerts.ExpectedEnclosingMark{}, alerts.NewSingle(p.peek()), string(tokens.RightParen)), tokens.RightParen)
	} else {
		ret = append(ret, p.Type())
	}

	return ret
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
