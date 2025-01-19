package parser

import (
	"fmt"
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/helpers"
	"hybroid/tokens"
)

type Parser struct {
	alerts.AlertHandler

	program []ast.Node
	current int
	tokens  []tokens.Token
	Context ParserContext
}

type ParserContext struct {
	EnvDeclaration  *ast.EnvironmentDecl
	FunctionReturns helpers.Stack[int]

	// ONLY USE WHENEVER YOU ARE CHECKING NODES AND MAKE SURE YOU DIDNT FORGET TO DISABLE IT
	IgnoreAlerts helpers.Stack[bool]
}

func NewParser() Parser {
	parser := Parser{
		program: make([]ast.Node, 0),
		current: 0,
		tokens:  make([]tokens.Token, 0),
		Context: ParserContext{
			EnvDeclaration:  nil,
			IgnoreAlerts:    helpers.NewStack[bool]("IgnoreAlerts"),
			FunctionReturns: helpers.NewStack[int]("FunctionReturns"),
		},
	}

	parser.Context.IgnoreAlerts.Push("default", false)
	parser.Context.FunctionReturns.Push("default", 0)

	return parser
}

func (p *Parser) AssignTokens(tokens []tokens.Token) {
	p.tokens = tokens
}

type ParserError struct{}

func (p *Parser) Alert(alertType alerts.Alert, args ...any) {
	if p.Context.IgnoreAlerts.Top().Item {
		return
	}

	p.Alert_(alertType, args...)
}

func (p *Parser) AlertPanic(alertType alerts.Alert, args ...any) {
	if p.Context.IgnoreAlerts.Top().Item {
		return
	}

	p.Alert_(alertType, args...)

	if alertType.GetAlertType() == alerts.Error {
		panic(ParserError{})
	}
}

func (p *Parser) AlertI(alert alerts.Alert) {
	if p.Context.IgnoreAlerts.Top().Item {
		return
	}

	p.AlertI_(alert)

	if alert.GetAlertType() == alerts.Error {
		//panic(ParserError{})
	}
}

func (p *Parser) synchronize() {
	//p.advance()
	for !p.isAtEnd() {
		switch p.peek().Type {
		case tokens.For, tokens.Fn, tokens.If, tokens.Repeat, tokens.Tick,
			tokens.Return, tokens.Let, tokens.While, tokens.Pub, tokens.Const,
			tokens.Break, tokens.Continue,
			tokens.Class:
			return
		case tokens.Entity:
			if p.peek(1).Type == tokens.Identifier && p.peek(2).Type == tokens.LeftBrace {
				return
			}
		} // pub fn (entity thing) { }
		// entity yes{}

		p.advance()
	}
}

func (p *Parser) isMultiComparison() bool {
	return p.match(tokens.And, tokens.Or)
}

func (p *Parser) isComparison() bool {
	return p.match(tokens.Greater, tokens.GreaterEqual, tokens.Less, tokens.LessEqual, tokens.BangEqual, tokens.EqualEqual)
}

// Checks if the current position the parser is at is the End Of File
func (p *Parser) isAtEnd() bool {
	return p.peek().Type == tokens.Eof
}

// Advances by the given offset into the next token and returns the previous token before advancing
func (p *Parser) advance(offset ...int) tokens.Token {
	currentOffset := 1
	if offset != nil {
		currentOffset = offset[0]
	}

	if currentOffset < 0 {
		panic("Attempt to advance with a negative offset. Use disadvance() instead!")
	}

	t := p.tokens[p.current]
	index := p.current + currentOffset

	if index < len(p.tokens) {
		p.current = index
	}

	return t
}

// Disadvances by the given offset into the previous tokens and returns the current token after disadvancing
func (p *Parser) disadvance(offset ...int) tokens.Token {
	currentOffset := 1
	if offset != nil {
		currentOffset = offset[0]
	}

	if currentOffset < 0 {
		panic("Attempt to disadvance with a negative offset (which moves forward). Use advance() instead!")
	}

	index := p.current - currentOffset

	if index >= 0 {
		p.current = index
	}

	return p.tokens[p.current]
}

// Peeks into the current token or peeks at the token that is offset from the current position by the given offset
func (p *Parser) peek(offset ...int) tokens.Token {
	index := p.current
	if offset != nil {
		index += offset[0]
	}

	if index >= 0 && index < len(p.tokens) {
		return p.tokens[index]
	} else {
		return p.tokens[p.current]
	}
}

// Checks if the current type is the specified token type. Returns false if it's the End Of File
func (p *Parser) check(tokenType tokens.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == tokenType
}

// Matches the given list of tokens and advances if they match.
func (p *Parser) match(types ...tokens.TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

// Consumes one of the tokens in the given list and advances if it matches.
func (p *Parser) consume(alert alerts.Alert, types ...tokens.TokenType) (tokens.Token, bool) {
	if p.isAtEnd() {
		token := p.peek()
		p.AlertI(alert)
		return token, false // error
	}
	for _, tokenType := range types {
		if p.check(tokenType) {
			return p.advance(), true
		}
	}
	token := p.advance()
	p.AlertI(alert)
	return token, false // error
}

func (p *Parser) ParseTokens() []ast.Node {
	firstToken := p.peek()
	if firstToken.Type != tokens.Env {
		p.Alert(&alerts.ExpectedEnvironment{}, alerts.NewSingle(firstToken))
		return []ast.Node{}
	}
	envDecl := p.environmentDeclaration()
	if envDecl.GetType() == ast.NA {
		return []ast.Node{}
	}

	p.program = append(p.program, envDecl)

	for !p.isAtEnd() {
		statement := p.declaration()
		if statement == nil {
			continue
		}
		if statement.GetType() != ast.NA {
			p.program = append(p.program, statement)
		}
	}

	fmt.Printf("%+#v", p.program)

	return p.program
}
