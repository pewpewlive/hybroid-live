package parser

import (
	"fmt"
	"hybroid/lexer"
)

type NodeType int

const (
	VariableDeclarationStmt NodeType = iota
	FunctionDeclarationStmt

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
	token := p.peek().Type;
	next := p.peek(1).Type;
	
	if token == lexer.Pub && next == lexer.Fn {
		p.advance()
	}

	fmt.Print("wha")

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
	}
	return p.expression()
}

func (p *Parser) functionDeclarationStmt() *Node {
	fnDec := Node{
		NodeType: FunctionDeclarationStmt,
		Token: p.peek(-1),
	}

	fnDec.IsLocal = p.peek(-2).Type != lexer.Pub
	
	ident, ok := p.consume(lexer.Identifier, "expected a function name")
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
	if _, success := p.consume(lexer.LeftBrace, "expected body of the function"); success {
		for !p.match(lexer.RightBrace) {
			statement := p.statement()
			if statement != nil {
				prog.Body = append(prog.Body, *statement)
			}
		}
	}// we might not be handling the case where there is no closing brace

	fnDec.Program = &prog

	return &fnDec
}

func (p *Parser) addToStmt() *Node {
	add := Node{
		NodeType: AddStmt,
		Token:    p.peek(-1),
	}

	add.Expression = p.expression()

	if _, ok := p.consume(lexer.To, "expected keyword 'to' after expression in an 'add' statement"); !ok {
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

	if _, ok := p.consume(lexer.From, "expected keyword 'from' after expression in a 'remove' statement"); !ok {
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
		Token:    p.peek(-1),//let or pub, important
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

func (p *Parser) UpdateTokens(tokens []lexer.Token) {
	p.tokens = tokens
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
