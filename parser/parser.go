package parser

import (
	"hybroid/lexer"
)

type NodeType int

const (
	VariableDeclarationStmt NodeType = iota + 1
	FunctionDeclarationStmt

	DirectiveStmt

	AddStmt
	RemoveStmt
	ReturnStmt

	Prog

	AssignmentExpr
	LiteralExpr
	UnaryExpr
	BinaryExpr
	GroupingExpr
	ListExpr
	MapExpr
	CallExpr
	MemberExpr

	Identifier
)

type PrimitiveValueType int

const (
	Number PrimitiveValueType = iota + 1
	String
	Bool
	FixedPoint
	Fixed
	Radian
	Degree
	List
	Map
	Nil
	Func
	Entity
	Struct
	Ident

	Undefined
)

type Node struct {
	NodeType    NodeType
	Identifier  string
	Program     *Program
	Value       any
	Value2      any
	ValueType   PrimitiveValueType
	Left, Right *Node
	Expression  *Node
	IsLocal     bool
	Token       lexer.Token
}

type Program struct {
	Body []Node
}

type Parser struct {
	program Program
	current int
	tokens  []lexer.Token
	Errors  []ParserError
}

func New(tokens []lexer.Token) *Parser {
	return &Parser{Program{},0, tokens, make([]ParserError, 0)}
}

func (p *Parser) statement() *Node {
	token := p.peek().Type
	next := p.peek(1).Type

	if token == lexer.Pub && next == lexer.Fn {
		p.advance()
		token = p.peek().Type
	}

	switch token {
	case lexer.Let, lexer.Pub, lexer.Const:
		p.advance()
		return p.variableDeclaration()
	case lexer.At:
		p.advance()
		return p.directiveCall()
	case lexer.Add:
		p.advance()
		return p.addToStmt()
	case lexer.Remove:
		p.advance()
		return p.removeFromStmt()
	case lexer.Fn:
		p.advance()
		return p.functionDeclarationStmt()
	case lexer.Return:
		p.advance()
		return p.returnStmt()
	}
	expr := p.expression()
	if expr.NodeType == 0 {
		p.error(p.peek(), "expected expression")
	}
	return expr
}

func (p *Parser) returnStmt() *Node {
	returnStmt := Node{
		NodeType: ReturnStmt,
		Token: p.peek(-1),
	}
	args := []Node{}
	expr := p.assignment()
	if expr.NodeType == 0 {
		p.error(p.peek(), "expected expression")
	}
	args = append(args, *expr)
	for p.match(lexer.Comma) {
		expr = p.assignment()
		if expr.NodeType == 0 {
			p.error(p.peek(), "expected expression")
		}
		args = append(args, *expr)
	}

	returnStmt.Value = args

	return &returnStmt
}

func (p *Parser) functionDeclarationStmt() *Node {
	fnDec := Node{
		NodeType: FunctionDeclarationStmt,
		Token:    p.peek(-1),
	}

	fnDec.IsLocal = p.peek(-2).Type != lexer.Pub

	ident, ok := p.consume("expected a function name", lexer.Identifier)
	if !ok {
		return &fnDec
	}

	fnDec.Identifier = ident.Lexeme

	args := p.arguments()
	var params []lexer.Token

	for _, arg := range args {
		if arg.NodeType == Identifier {
			params = append(params, arg.Token)
			continue
		}
		p.error(arg.Token, "expected identifier in function declaration")
	}

	fnDec.Value = params

	prog := Program{}
	if _, success := p.consume("expected body of the function", lexer.LeftBrace); success {
		for !p.match(lexer.RightBrace) {
			statement := p.statement()
			if statement != nil {
				prog.Body = append(prog.Body, *statement)
			}
		}
	} // we might not be handling the case where there is no closing brace

	fnDec.Program = &prog

	return &fnDec
}

func (p *Parser) addToStmt() *Node {
	add := Node{
		NodeType: AddStmt,
		Token:    p.peek(-1),
	}

	add.Expression = p.expression()
	if add.NodeType == 0 {
		p.error(p.peek(), "expected expression")
	}

	if _, ok := p.consume("expected keyword 'to' after expression in an 'add' statement", lexer.To); !ok {
		return &add
	}

	if ident, ok := p.consume("expected identifier after keyword 'to'", lexer.Identifier); ok {
		add.Identifier = ident.Lexeme
	}

	return &add
}

func (p *Parser) removeFromStmt() *Node {
	remove := Node{
		NodeType: RemoveStmt,
		Token:    p.peek(-1),
	}

	remove.Expression = p.expression()
	if remove.NodeType == 0 {
		p.error(p.peek(), "expected expression")
	}

	if _, ok := p.consume("expected keyword 'from' after expression in a 'remove' statement", lexer.From); !ok {
		return &remove
	}

	if ident, ok := p.consume("expected identifier after keyword 'from'", lexer.Identifier); ok {
		remove.Identifier = ident.Lexeme
	}

	return &remove
}

func (p *Parser) variableDeclaration() *Node {
	variable := Node{
		NodeType: VariableDeclarationStmt,
		Token:    p.peek(-1), //let or pub, important
	}

	ident, _ := p.consume("expected identifier in variable declaration", lexer.Identifier)
	idents := []string{ident.Lexeme}
	for p.match(lexer.Comma) {
		ident, identOk := p.consume("expected identifier in variable declaration", lexer.Identifier)
		if !identOk {
			return &Node{Token: p.peek(-1)}
		}

		idents = append(idents, ident.Lexeme)
	}

	variable.Value = idents

	if _, ok := p.consume("expected '=' after identifier in variable declaration", lexer.Equal); !ok {
		return &Node{Token: p.peek(-1)}
	} // let a, b = name()
	
	expr := p.expression()
	if expr.NodeType == 0 {
		p.error(p.peek(), "expected expression")
	}
	
	exprs := []Node{*expr}
	for p.match(lexer.Comma) {
		expr = p.expression()
		if expr.NodeType == 0 {
			p.error(p.peek(), "expected expression")
		}
		exprs = append(exprs, *expr)
	}
	variable.Value2 = exprs

	return &variable
}

func (p *Parser) UpdateTokens(tokens []lexer.Token) {
	p.tokens = tokens
}

func (p *Parser) ParseTokens() Program {
	// Expect environment directive call as node
	statement := p.statement()
	if !p.verifyEnvironmentDirective(statement) {
		return p.program
	}
	p.program.Body = append(p.program.Body, *statement)

	for !p.isAtEnd() {
		statement := p.statement()
		if statement != nil {
			p.program.Body = append(p.program.Body, *statement)
		}
	}

	return p.program
}
