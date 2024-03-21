package parser

import (
	"hybroid/lexer"
	"strings"
)

func (p *Parser) expression() *Node {
	return p.assignment()
}

func (p *Parser) assignment() *Node {
	expr := p.list()

	if p.match(lexer.Equal) {
		value := p.assignment()
		expr = &Node{NodeType: AssignmentExpr, Expression: expr, Right: value, Token: p.peek(-1)}
	} else if p.match(lexer.PlusEqual) {
		value := p.term()
		binExpr := createBinExpr(expr, p.peek(-1), lexer.Plus, "+", &Node{NodeType: GroupingExpr, Expression: value})
		expr = &Node{NodeType: AssignmentExpr, Expression: expr, Right: binExpr, Token: p.peek(-1)}
	} else if p.match(lexer.MinusEqual) {
		value := p.term()
		binExpr := createBinExpr(expr, p.peek(-1), lexer.Minus, "-", &Node{NodeType: GroupingExpr, Expression: value})
		expr = &Node{NodeType: AssignmentExpr, Expression: expr, Right: binExpr, Token: p.peek(-1)}
	} else if p.match(lexer.SlashEqual) {
		value := p.term()
		binExpr := createBinExpr(expr, p.peek(-1), lexer.Slash, "/", &Node{NodeType: GroupingExpr, Expression: value})
		expr = &Node{NodeType: AssignmentExpr, Expression: expr, Right: binExpr, Token: p.peek(-1)}
	} else if p.match(lexer.StarEqual) {
		value := p.term()
		binExpr := createBinExpr(expr, p.peek(-1), lexer.Star, "*", &Node{NodeType: GroupingExpr, Expression: value})
		expr = &Node{NodeType: AssignmentExpr, Expression: expr, Right: binExpr, Token: p.peek(-1)}
	} else if p.match(lexer.CaretEqual) {
		value := p.term()
		binExpr := createBinExpr(expr, p.peek(-1), lexer.Caret, "^", &Node{NodeType: GroupingExpr, Expression: value})
		expr = &Node{NodeType: AssignmentExpr, Expression: expr, Right: binExpr, Token: p.peek(-1)}
	}

	return expr
}

func (p *Parser) list() *Node {
	if !p.match(lexer.LeftBracket) {
		return p.parseMap()
	}

	token := p.peek(-1)
	list := make([]Node, 0)
	for !p.check(lexer.RightBracket) {
		exprInList := p.expression()
		if exprInList.NodeType == 0 {
			p.error(p.peek(), "expected expression")
		}

		token, _ := p.consume("expected ',' or ']' after expression", lexer.Comma, lexer.RightBracket)

		list = append(list, *exprInList)
		if token.Type == lexer.RightBracket || token.Type == lexer.Eof {
			break
		}
	}

	return &Node{NodeType: ListExpr, ValueType: List, Value: list, Token: token}
}

func (p *Parser) parseMap() *Node {
	if !p.match(lexer.LeftBrace) {
		return p.equality()
	}

	token := p.peek(-1)
	parsedMap := make(map[string]Node, 0)
	for !p.check(lexer.RightBrace) {
		key := p.primary()

		var newKey string
		switch key.ValueType {
		case Ident:
			newKey = key.Identifier
		case String:
			newKey = key.Token.Literal
		default:
			p.error(key.Token, "expected either string or an identifier in map initialization")
			return &Node{Token: p.peek(-1)}
		}

		if _, ok := p.consume("expected ':' after map key", lexer.Colon); !ok {
			return &Node{Token: p.peek(-1)}
		}

		expr := p.expression()
		if expr.NodeType == 0 {
			p.error(p.peek(), "expected expression")
		}

		if p.peek().Type == lexer.RightBrace {
			parsedMap[newKey] = *expr
			break
		}

		if _, ok := p.consume("expected ',' or '}' after expression", lexer.Comma); !ok {
			return &Node{Token: p.peek(-1)}
		}

		parsedMap[newKey] = *expr
	}
	p.advance()

	return &Node{NodeType: MapExpr, ValueType: Map, Value: parsedMap, Token: token}
}

func (p *Parser) equality() *Node {
	expr := p.comparison()

	if p.match(lexer.BangEqual, lexer.EqualEqual) {
		operator := p.peek(-1)
		right := p.comparison()
		expr = &Node{NodeType: BinaryExpr, Left: expr, Token: operator, Right: right}
	}

	return expr
}

func (p *Parser) comparison() *Node {
	expr := p.term()

	if p.match(lexer.Greater, lexer.GreaterEqual, lexer.Less, lexer.LessEqual) {
		operator := p.peek(-1)
		right := p.term()
		expr = &Node{NodeType: BinaryExpr, Left: expr, Token: operator, Right: right}
	}

	return expr
}

func (p *Parser) term() *Node {
	expr := p.factor()

	if p.match(lexer.Plus, lexer.Minus) {
		operator := p.peek(-1)
		right := p.term()
		expr = &Node{NodeType: BinaryExpr, Left: expr, Token: operator, Right: right, ValueType: Undefined}
	}

	return expr
}

func (p *Parser) factor() *Node {
	expr := p.unary()

	if p.match(lexer.Star, lexer.Slash) {
		operator := p.peek(-1)
		right := p.factor()
		expr = &Node{NodeType: BinaryExpr, Left: expr, Token: operator, Right: right, ValueType: Undefined}
	}

	return expr
}

func (p *Parser) unary() *Node {
	if p.match(lexer.Bang, lexer.Minus) {
		operator := p.peek(-1)
		right := p.unary()
		return &Node{NodeType: UnaryExpr, Token: operator, Right: right}
	}

	return p.memberCall() //p.memberCall()
}

func (p *Parser) memberCall() *Node { // for maps
	expr := p.member()

	if p.check(lexer.LeftParen) {
		return p.call(expr)
	}

	return expr
}

func (p *Parser) call(caller *Node) *Node {
	call_expr := Node{
		NodeType:   CallExpr,
		Identifier: caller.Token.Lexeme,
		Expression: caller,
		Value:      p.arguments(),
		Token:      caller.Token,
	}

	if p.check(lexer.LeftParen) { // name()()
		call_expr = *p.call(&call_expr)
	}

	return &call_expr
}

func (p *Parser) arguments() []Node {
	if _, ok := p.consume("expected opening paren after an identifier", lexer.LeftParen); !ok {
		return nil
	}

	var args []Node
	if p.match(lexer.RightParen) {
		args = make([]Node, 0)
	} else {
		args = append(args, *p.assignment())
		for p.match(lexer.Comma) {
			args = append(args, *p.assignment())
		}
		p.consume("expected closing paren after arguments", lexer.RightParen)
	}

	return args
}

func (p *Parser) member() *Node {
	expr := p.primary()

	for p.match(lexer.Dot, lexer.LeftBracket) {
		operator := p.peek(-1) // . or [
		var prop Node

		prop = *p.primary()
		if operator.Type == lexer.Dot && prop.NodeType != Identifier {
			p.error(prop.Token, "expected identifier after '.'")
		}
		if operator.Type == lexer.LeftBracket {
			if prop.NodeType != NodeType(String) {
				p.error(prop.Token, "expected string after '['")
			}
			p.consume("expected closing bracket", lexer.RightBracket)
		}

		expr = &Node{
			NodeType:   MemberExpr,
			Expression: expr,
			Value:      prop,
		}
	}

	return expr
}

func (p *Parser) primary() *Node {
	if p.match(lexer.False) {
		return &Node{NodeType: LiteralExpr, Value: "false", ValueType: Bool, Token: p.peek(-1)}
	}
	if p.match(lexer.True) {
		return &Node{NodeType: LiteralExpr, Value: "true", ValueType: Bool, Token: p.peek(-1)}
	}
	if p.match(lexer.Nil) {
		return &Node{NodeType: LiteralExpr, Value: "nil", ValueType: Nil, Token: p.peek(-1)}
	}

	if p.match(lexer.Number, lexer.Fixed, lexer.FixedPoint, lexer.Degree, lexer.Radian, lexer.String) {
		literal := p.peek(-1)
		var valueType PrimitiveValueType
		allowFX := p.program.Body[0].Expression.Identifier == "Level" || p.program.Body[0].Expression.Identifier == "Shared"
		switch literal.Type {
		case lexer.Number:
			if allowFX && strings.ContainsRune(literal.Lexeme,'.') {
				p.error(literal, "cannot have a float in a level or shared environment")
			}
			valueType = Number
		case lexer.Fixed:
			if !allowFX {
				p.error(literal, "cannot have a fixed in a mesh, sound or luageneric environment")
			}
			valueType = Fixed
		case lexer.FixedPoint:
			if !allowFX {
				p.error(literal, "cannot have a fixedpoint in a mesh, sound or luageneric environment")
			}
			valueType = FixedPoint
		case lexer.Degree:
			if !allowFX {
				p.error(literal, "cannot have a degree, sound or luageneric environment")
			}
			valueType = Degree
		case lexer.Radian:
			if !allowFX {
				p.error(literal, "cannot have a radian in a mesh, sound or luageneric environment")
			}
			valueType = Radian
		case lexer.String:
			valueType = String
		}
		return &Node{NodeType: LiteralExpr, Value: literal.Literal, ValueType: valueType, Token: literal}
	}

	if p.match(lexer.Identifier) {
		token := p.peek(-1)
		return &Node{NodeType: Identifier, Identifier: token.Lexeme, Token: token, ValueType: Ident}
	}

	if p.match(lexer.LeftParen) {
		token := p.peek(-1)
		expr := p.expression()
		if expr.NodeType == 0 {
			p.error(p.peek(), "expected expression")
		}
		p.consume("expected ')' after expression", lexer.RightParen)
		return &Node{NodeType: GroupingExpr, Expression: expr, Token: token, ValueType: expr.ValueType}
	}
	p.advance()
	return &Node{Token: p.peek()}
}
