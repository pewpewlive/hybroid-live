package parser

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
)

// Creates a BinaryExpr
func (p *Parser) createBinExpr(left ast.Node, operator lexer.Token, tokenType lexer.TokenType, lexeme string, right ast.Node) ast.Node {
	valueType := p.determineValueType(left, right)
	return &ast.BinaryExpr{
		Left:      left,
		Operator:  lexer.Token{Type: tokenType, Lexeme: lexeme, Literal: "", Location: operator.Location},
		Right:     right,
		ValueType: valueType,
	}
}

// Checks if the value type is expected to be a fixedpoint
func IsFx(valueType ast.PrimitiveValueType) bool {
	return valueType == ast.FixedPoint || valueType == ast.Fixed || valueType == ast.Radian || valueType == ast.Degree
}

func (p *Parser) getOp(opEqual lexer.Token) lexer.Token {
	switch opEqual.Type {
	case lexer.PlusEqual:
		return lexer.Token{Type: lexer.Plus, Location: opEqual.Location, Literal: opEqual.Literal, Lexeme: "+"}
	case lexer.MinusEqual:
		return lexer.Token{Type: lexer.Minus, Location: opEqual.Location, Literal: opEqual.Literal, Lexeme: "-"}
	case lexer.SlashEqual:
		return lexer.Token{Type: lexer.Slash, Location: opEqual.Location, Literal: opEqual.Literal, Lexeme: "/"}
	case lexer.StarEqual:
		return lexer.Token{Type: lexer.Star, Location: opEqual.Location, Literal: opEqual.Literal, Lexeme: "*"}
	case lexer.CaretEqual:
		return lexer.Token{Type: lexer.Caret, Location: opEqual.Location, Literal: opEqual.Literal, Lexeme: "^"}
	case lexer.ModuloEqual:
		return lexer.Token{Type: lexer.Modulo, Location: opEqual.Location, Literal: opEqual.Literal, Lexeme: "%"}
	default:
		return lexer.Token{}
	}
}

func (p *Parser) getParam() ast.Param {
	typ, ide := p.TypeWithVar()
	if ide.GetType() != ast.Identifier {
		p.error(ide.GetToken(), "expected identifier as parameter")
	}
	if typ == nil {
		p.error(ide.GetToken(), "parameters need to be declared with a type before the name")
	}
	return ast.Param{Type: typ, Name: ide.GetToken()}
}

func (p *Parser) parameters(opening lexer.TokenType, closing lexer.TokenType) []ast.Param {
	if _, ok := p.consume(fmt.Sprintf("expected %s", string(opening)), opening); !ok {
		return nil
	}

	var args []ast.Param
	if p.match(lexer.RightParen) {
		args = make([]ast.Param, 0)
	} else {
		args = append(args, p.getParam())
		for p.match(lexer.Comma) {
			args = append(args, p.getParam())
		}
		p.consume(fmt.Sprintf("expected %s after an identifier", string(closing)), closing)
	}

	return args
}

func (p *Parser) arguments() []ast.Node {
	if _, ok := p.consume("expected opening paren", lexer.LeftParen); !ok {
		return nil
	}

	var args []ast.Node
	if p.match(lexer.RightParen) {
		args = make([]ast.Node, 0)
	} else {
		arg := p.expression()
		args = append(args, arg)
		for p.match(lexer.Comma) {
			arg := p.expression()
			args = append(args, arg)
		}
		p.consume("expected closing paren after arguments", lexer.RightParen)
	}

	return args
}

func (p *Parser) returnings() []*ast.TypeExpr {
	ret := make([]*ast.TypeExpr, 0)
	if !p.match(lexer.ThinArrow) {
		return ret
	}
	isList := false
	if p.match(lexer.LeftParen) {
		isList = true
	}
	if !p.PeekIsType() {
		return ret
	}
	ret = append(ret, p.Type())
	for isList && p.match(lexer.Comma) {
		if !p.PeekIsType() {
			return ret
		}
		ret = append(ret, p.Type())
	}
	if isList {
		p.consume("expected closing parenthesis", lexer.RightParen)
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
	return &ast.IdentifierExpr{Name: typ.Name.GetToken(), ValueType: ast.Unknown}
}

func (p *Parser) TypeWithVar() (*ast.TypeExpr, ast.Node) {
	typ := p.Type()

	node := p.primary()

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
