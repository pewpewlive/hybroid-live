package parser

import "hybroid/lexer"

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
		return p.equality()
	}

	token := p.peek(-1)
	list := make([]Node, 0)
	for !p.check(lexer.RightBracket) {
		exprInList := p.expression()

		if p.peek().Type == lexer.RightBracket {
			list = append(list, *exprInList)
			break
		}

		if _, ok := p.consume(lexer.Comma, "expected ',' or ']' after value"); !ok {
			return &Node{Token: p.peek(-1)}
		}

		list = append(list, *exprInList)
	}
	p.advance()

	return &Node{NodeType: ListExpr, ValueType: List, Value: list, Token: token}
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

	return p.primary()
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

	if p.match(lexer.Number, lexer.FixedPoint, lexer.Degree, lexer.Radian, lexer.String) {
		literal := p.peek(-1)
		var valueType PrimitiveValueType
		switch literal.Type {
		case lexer.Number:
			valueType = Number
		case lexer.FixedPoint:
			valueType = FixedPoint
		case lexer.Degree:
			valueType = Degree
		case lexer.Radian:
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
		p.consume(lexer.RightParen, "expected ')' after expression")
		return &Node{NodeType: GroupingExpr, Expression: expr, Token: token, ValueType: expr.ValueType}
	}

	p.advance()
	p.error(p.peek(-1), "expected expression")
	return &Node{Token: p.peek(-1)}
}
