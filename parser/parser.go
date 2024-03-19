package parser

import (
	"hybroid/lexer"
)

type NodeType int

const (
	VariableDeclarationStmt NodeType = iota
	FuncDeclarationStmt              // TODO: Implement this

	DirectiveStmt

	AddStmt
	RemoveStmt

	Prog

	AssignmentExpr
	LiteralExpr
	UnaryExpr
	BinaryExpr
	GroupingExpr
	ListExpr

	Identifier
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
	current int
	tokens  []lexer.Token
	Errors  []ParserError
}

func New(tokens []lexer.Token) *Parser {
	return &Parser{0, tokens, make([]ParserError, 0)}
}

func (p *Parser) statement() *Node {
	switch p.peek().Type {
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
	}

	return p.expression()
}

func (p *Parser) addToStmt() *Node {
	add := Node{
		NodeType: AddStmt,
		Token:    p.peek(-1),
	}

	add.Expression = p.expression()

	if _, ok := p.consume(lexer.To, "expected keyword 'to' after expression in an add statement"); !ok {
		return &add
	}

	if ident, ok := p.consume(lexer.Identifier, "expected identifier after keyword 'to'"); ok {
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

	if _, ok := p.consume(lexer.From, "expected keyword 'from' after expression in a remove statement"); !ok {
		return &remove
	}

	if ident, ok := p.consume(lexer.Identifier, "expected identifier after keyword 'from'"); ok {
		remove.Identifier = ident.Lexeme
	}

	return &remove
}

func (p *Parser) variableDeclaration() *Node {
	variable := Node{
		NodeType: VariableDeclarationStmt,
		Token:    p.peek(-1),
	}

	ident, identOk := p.consume(lexer.Identifier, "expected identifier in variable declaration")
	if !identOk {
		return &Node{Token: p.peek(-1)}
	}
	variable.Identifier = ident.Lexeme

	if _, ok := p.consume(lexer.Equal, "expected '=' after identifier in variable declaration"); !ok {
		return &Node{Token: p.peek(-1)}
	}
	variable.Expression = p.expression()

	return &variable
}

func (p *Parser) ParseTokens() Program {
	program := Program{}

	// Expect environment directive call as node
	statement := p.statement()
	if !p.verifyEnvironmentDirective(statement) {
		return program
	}

	for !p.isAtEnd() {
		statement := p.statement()
		if statement != nil {
			program.Body = append(program.Body, *statement)
		}
	}

	return program
}
