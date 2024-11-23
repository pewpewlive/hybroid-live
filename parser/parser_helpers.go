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
		Operator:  tokens.Token{Type: tokenType, Lexeme: lexeme, Literal: "", Location: operator.Location},
		Right:     right,
		ValueType: valueType,
	}
}

// Checks if the value type is expected to be a fixedpoint
func IsFx(valueType ast.PrimitiveValueType) bool {
	return valueType == ast.FixedPoint || valueType == ast.Fixed || valueType == ast.Radian || valueType == ast.Degree
}

func (p *Parser) PeekIsType() bool {
	tokenType := p.peek().Type

	if tokenType == tokens.Fn {
		return p.peek(1).Type == tokens.LeftParen
	}

	return !(tokenType != tokens.Identifier && tokenType != tokens.Fn && tokenType != tokens.Struct && tokenType != tokens.Entity /* && tokens.Type != tokens.DotDotDot*/)
}

func (p *Parser) getOp(opEqual tokens.Token) tokens.Token {
	switch opEqual.Type {
	case tokens.PlusEqual:
		return tokens.Token{Type: tokens.Plus, Location: opEqual.Location, Literal: opEqual.Literal, Lexeme: "+"}
	case tokens.MinusEqual:
		return tokens.Token{Type: tokens.Minus, Location: opEqual.Location, Literal: opEqual.Literal, Lexeme: "-"}
	case tokens.SlashEqual:
		return tokens.Token{Type: tokens.Slash, Location: opEqual.Location, Literal: opEqual.Literal, Lexeme: "/"}
	case tokens.StarEqual:
		return tokens.Token{Type: tokens.Star, Location: opEqual.Location, Literal: opEqual.Literal, Lexeme: "*"}
	case tokens.CaretEqual:
		return tokens.Token{Type: tokens.Caret, Location: opEqual.Location, Literal: opEqual.Literal, Lexeme: "^"}
	case tokens.ModuloEqual:
		return tokens.Token{Type: tokens.Modulo, Location: opEqual.Location, Literal: opEqual.Literal, Lexeme: "%"}
	default:
		return tokens.Token{}
	}
}

func (p *Parser) getParam() ast.Param {
	typ, ide := p.TypeWithVar()
	if ide.GetType() != ast.Identifier {
		p.error(ide.GetToken(), "expected identifier as parameter")
	}
	return ast.Param{Type: typ, Name: ide.GetToken()}
}

func (p *Parser) parameters(opening tokens.TokenType, closing tokens.TokenType) []ast.Param {
	if !p.match(opening) {
		p.Alert(&alerts.ExpectedParenthesis{}, p.peek(), p.peek().Location, "(")
		return []ast.Param{}
	}

	var args []ast.Param
	if p.match(closing) {
		args = make([]ast.Param, 0)
	} else {
		var previous *ast.TypeExpr
		param := p.getParam()
		if param.Type == nil {
			if len(args) == 0 {
				p.Alert(&alerts.ExpectedType{}, p.peek(-1), p.peek(-1).Location) //param.Name, "parameter need to be declared with a type before the name")
			} else {
				param.Type = previous
			}
		} else {
			previous = param.Type
		}
		args = append(args, param)
		for p.match(tokens.Comma) {
			param := p.getParam()
			if param.Type == nil {
				if len(args) == 0 {
					p.Alert(&alerts.ExpectedType{}, p.peek(-1), p.peek(-1).Location)
				} else {
					param.Type = previous
				}
			} else {
				previous = param.Type
			}
			args = append(args, param)
		}
		p.consume(p.NewAlert(&alerts.ExpectedParenthesis{}, p.peek(), p.peek().Location, ")"), closing)
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
		p.error(token, "expected type identifier in generic parameters")
	} else {
		params = append(params, &ast.IdentifierExpr{Name: token, ValueType: ast.Invalid})
	}

	for p.match(tokens.Comma) {
		token := p.advance()
		if token.Type != tokens.Identifier {
			p.error(token, "expected type identifier in generic parameters")
		} else {
			params = append(params, &ast.IdentifierExpr{Name: token, ValueType: ast.Invalid})
		}
	}

	p.consumeOld("expected '>' in generic parameters", tokens.Greater)

	return params
}

func (p *Parser) genericArguments() ([]*ast.TypeExpr, bool) {
	current := p.getCurrent()
	params := []*ast.TypeExpr{}
	if !p.match(tokens.Less) {
		return params, false
	}

	params = append(params, p.Type())

	for p.match(tokens.Comma) {
		params = append(params, p.Type())
	}

	if !p.match(tokens.Greater) {
		p.disadvance(p.getCurrent() - current)
		return params, false
	}

	return params, true
}

func (p *Parser) arguments() []ast.Node {
	if _, ok := p.consumeOld("expected opening paren", tokens.LeftParen); !ok {
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
		p.consumeOld("expected closing paren after arguments", tokens.RightParen)
	}

	return args
}

func (p *Parser) returnings() []*ast.TypeExpr {
	ret := make([]*ast.TypeExpr, 0)
	if !p.match(tokens.ThinArrow) {
		return ret
	}
	isList := false
	if p.match(tokens.LeftParen) {
		isList = true
	}
	if !p.PeekIsType() {
		return ret
	}
	ret = append(ret, p.Type())
	for isList && p.match(tokens.Comma) {
		if !p.PeekIsType() {
			return ret
		}
		ret = append(ret, p.Type())
	}
	if isList {
		p.consumeOld("expected closing parenthesis", tokens.RightParen)
	}
	return ret
}

func (p *Parser) TypeWasVar(typ *ast.TypeExpr) *ast.IdentifierExpr {
	if typ.WrappedType != nil {
		return nil
	}
	if typ.Params != nil {
		return nil
	}
	if typ.Returns != nil {
		return nil
	}
	return &ast.IdentifierExpr{Name: typ.Name.GetToken(), ValueType: ast.Object}
}

func (p *Parser) TypeWithVar() (*ast.TypeExpr, ast.Node) {
	typ := p.Type()

	if typ.Name != nil && typ.Name.GetToken().Lexeme == "type" {
		print("")
	}

	node := p.primary(false)

	if node.GetType() != ast.Identifier {
		if ident := p.TypeWasVar(typ); ident != nil {
			return nil, ident
		} else {
			return typ, node
		}
	}
	ident := node.(*ast.IdentifierExpr)

	return typ, ident
}
