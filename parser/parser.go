package parser

import (
	"fmt"
	"hybroid/lexer"
)

type NodeType int

const (
	VariableDeclarationStmt NodeType = iota
	DirectiveStmt
	AssignmentExpr
	LiteralExpr
	UnaryExpr
	BinaryExpr
	GroupingExpr
	IdentifierExpr
)

type PrimitiveValueType int

const (
	Number PrimitiveValueType = iota + 1
	String
	Bool
	FixedPoint
	Radian
	Degree
	List
	Map
	Nil

	Undefined
)

type Node struct {
	NodeType   NodeType
	Identifier string
	Value      any
	// EntityType
	// StructType
	ValueType   PrimitiveValueType
	Left, Right *Node
	Expression  *Node
	Token       lexer.Token
}

type Program struct {
	Body []Node
}

// type VariableDeclarationStmt struct {
// 	Identifier string
// 	Expression any
// }

// type AssignmentExpr struct {
// 	Asignee any
// 	Value   any
// }

// type LiteralExpr struct {
// 	Value any
// }

// type UnaryExpr struct {
// 	Operator lexer.Token
// 	Right    any
// }

// type BinaryExpr struct {
// 	Left     any
// 	Operator lexer.Token
// 	Right    any
// }

// type GroupingExpr struct {
// 	Expression any
// }

// type IdentifierExpr struct {
// 	Symbol string
// }

type ParserError struct {
	token   lexer.Token
	Message string
}

func (pe *ParserError) Msg() string {
	return fmt.Sprintf("Error: %v, at line: %v (%v)", pe.Message, pe.token.Line, pe.token.ToString())
}

type Parser struct {
	current int
	tokens  []lexer.Token
	Errors  []ParserError
}

func New() Parser {
	return Parser{}
}

func (p *Parser) error(token lexer.Token, err string) *ParserError {
	p.Errors = append(p.Errors, ParserError{
		token,
		err,
	})
	return &p.Errors[len(p.Errors)-1]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == lexer.Eof
}

func (p *Parser) advance() lexer.Token {
	t := p.tokens[p.current]
	p.current++
	return t
}

func (p *Parser) peek(offset ...int) lexer.Token {
	if offset == nil {
		return p.tokens[p.current]
	} else {
		return p.tokens[p.current+offset[0]]
	}
}

func (p *Parser) check(tokenType lexer.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == tokenType
}

func (p *Parser) match(types ...lexer.TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) consume(tokenType lexer.TokenType, message string) (lexer.Token, *ParserError) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	return lexer.Token{}, p.error(p.tokens[p.current], message) // error
}

func (p *Parser) statement() *Node {
	switch p.peek().Type {
	case lexer.Let:
		p.advance()
		return p.variableDeclaration()
		// case lexer.At:
		// 	p.advance()
		// 	p.directiveCall()
	}

	return p.expression()
}

// func (p *Parser) matchDirective(ident string) *Node {
// 	switch ident {
// 	case "Environment":

// 	case "Len":
// 	}
// }

// func (p *Parser) directiveCall() *Node {
// 	ident, err1 := p.consume(lexer.Identifier, "expected identifier in directive call")
// 	if err1 != nil {
// 		return nil
// 	}

// 	_, err2 := p.consume(lexer.LeftParen, "expected '(' after directive call")
// 	if err2 != nil {
// 		return nil
// 	}

// 	expr := p.expression()
// 	if expr == nil {
// 		return nil
// 	}

// 	_, err4 := p.consume(lexer.LeftParen, "expected ')' after directive call")
// 	if err4 != nil {
// 		return nil
// 	}

// 	directiveNode := Node {
// 		NodeType: DirectiveStmt,
// 	}

// }

func (p *Parser) validateOperands(left *Node, right *Node) bool {
	if left.ValueType == 0 {
		p.error(left.Token, "cannot perform arithmetic on nil value")
		return false
	} else if right.ValueType == 0 {
		p.error(right.Token, "cannot perform arithmetic on nil value")
		return false
	} else if left.ValueType == Undefined {
		p.error(left.Token, "cannot perform arithmetic on undefined value")
		return false
	} else if right.ValueType == Undefined {
		p.error(right.Token, "cannot perform arithmetic on undefined value")
		return false
	} else {
		if (left.ValueType == List || left.ValueType == Map || left.ValueType == String) ||
			(right.ValueType == List || right.ValueType == Map || right.ValueType == String) {

			p.error(left.Token, "cannot perform arithmetic on extraenous value")
			return false
		} else if left.ValueType != right.ValueType {
			p.error(left.Token, "left operand and right operand don't have the same type")
			return false
		}
	}

	return true
}

func (p *Parser) variableDeclaration() *Node {
	ident, err1 := p.consume(lexer.Identifier, "expected identifier in variable declaration")
	if err1 != nil {
		return &Node{Token: p.peek(-1)}
	}

	_, err2 := p.consume(lexer.Equal, "expected equal token following identifier in variable declaration")
	if err2 != nil {
		return &Node{Token: p.peek(-1)}
	}

	return &Node{
		NodeType:   VariableDeclarationStmt,
		Identifier: ident.Lexeme,
		Expression: p.expression(),
	}
}

func (p *Parser) expression() *Node {
	return p.assignment()
}

func (p *Parser) assignment() *Node {
	expr := p.equality()

	if p.match(lexer.Equal) {
		value := p.assignment()
		expr = &Node{NodeType: AssignmentExpr, Expression: expr, Value: *value}
	}

	return expr
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

func (p *Parser) term() *Node { // 1 - 10, 1 + 10
	expr := p.factor()

	if p.match(lexer.Plus, lexer.Minus) {
		operator := p.peek(-1)
		right := p.term()
		if p.validateOperands(expr, right) {
			expr = &Node{NodeType: BinaryExpr, Left: expr, Token: operator, Right: right, ValueType: expr.ValueType}
		} else {
			expr = &Node{NodeType: BinaryExpr, Left: expr, Token: operator, Right: right, ValueType: Undefined}
		}
	}

	return expr
}

func (p *Parser) factor() *Node {
	expr := p.unary()

	if p.match(lexer.Star, lexer.Slash) {
		operator := p.peek(-1)
		right := p.factor()
		if p.validateOperands(expr, right) {
			expr = &Node{NodeType: BinaryExpr, Left: expr, Token: operator, Right: right, ValueType: expr.ValueType}
		} else {
			expr = &Node{NodeType: BinaryExpr, Left: expr, Token: operator, Right: right, ValueType: Undefined}
		}
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
		return &Node{NodeType: LiteralExpr, Identifier: token.Lexeme, Token: token}
	}

	if p.match(lexer.LeftParen) {
		token := p.peek(-1)
		expr := p.expression()
		p.consume(lexer.RightParen, "expected ')' after expression")
		return &Node{NodeType: GroupingExpr, Expression: expr, Token: token, ValueType: expr.ValueType}
	}

	if p.match(lexer.LeftBracket) {
		token := p.peek(-1)
		list := make([]Node, 0)
		for !p.check(lexer.RightBracket) {
			exprInList := p.expression()

			if p.peek().Type == lexer.RightBracket {
				list = append(list, *exprInList)
				break
			}

			_, err := p.consume(lexer.Comma, "expected ',' or ']' after value")
			if err != nil {
				return &Node{Token: p.peek(-1)}
			}

			list = append(list, *exprInList)
		}
		p.advance()
		return &Node{NodeType: LiteralExpr, ValueType: List, Value: list, Token: token}
	}

	p.advance()
	p.error(p.peek(-1), "expected expression")
	return &Node{Token: p.peek(-1)}
}

func (p *Parser) ParseTokens(tokens []lexer.Token) Program {
	p.tokens = tokens

	program := Program{}

	for !p.isAtEnd() {
		stmt := p.statement()
		if stmt != nil {
			program.Body = append(program.Body, *stmt)
		}
	}

	return program
}
