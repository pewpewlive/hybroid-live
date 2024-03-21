package parser

import (
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
	current int
	tokens  []lexer.Token
	Errors  []ParserError
}

func New(tokens []lexer.Token) *Parser {
	return &Parser{0, tokens, make([]ParserError, 0)}
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
	}
	return p.expression()
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
	exprs := []Node{}
	for i := 1; i < len(idents); i++ {
		expr := *p.expression()
		if expr.NodeType == CallExpr {
			exprs = append(exprs, expr)
			if p.peek(1) != lexer.Comma {
				break // x, y = fn(), fn()
			}
		} else {
			exprs = append(exprs, expr)
		}

		p.consume("need comatos, lexer.Comma")
	}
	variable.Value2 = &exprs

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
