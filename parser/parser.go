package parser

import (
	"fmt"
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/tokens"
)

type Parser struct {
	alerts.AlertHandler

	program []ast.Node
	current int
	tokens  []tokens.Token
	Errors  []ast.Error
	Context ParserContext
}

type ParserContext struct {
	EnvStatement    *ast.EnvironmentStmt
	FunctionReturns []int
}

func NewParser() Parser {
	return Parser{
		program: make([]ast.Node, 0),
		current: 0,
		tokens:  make([]tokens.Token, 0),
		Context: ParserContext{
			EnvStatement:    nil,
			FunctionReturns: make([]int, 0),
		},
	}
}

func (p *Parser) AssignTokens(tokens []tokens.Token) {
	p.tokens = tokens
}

type ParserError struct{}

func (p *Parser) error(token tokens.Token, msg string) {
	errMsg := ast.Error{Token: token, Message: msg}
	fmt.Printf("%s\n", errMsg.Message)
	p.Errors = append(p.Errors, errMsg)
	panic(errMsg)
}

func (p *Parser) Alert(alertType alerts.Alert, args ...any) {
	p.Alert_(alertType, args...)

	if alertType.GetAlertType() == alerts.Error {
		panic(ParserError{})
	}
}

func (p *Parser) AlertI(alert alerts.Alert) {
	p.AlertI_(alert)

	if alert.GetAlertType() == alerts.Error {
		panic(ParserError{})
	}
}

func (p *Parser) synchronize() {
	//p.advance()
	for !p.isAtEnd() {
		switch p.peek().Type {
		case tokens.RightBrace:
			p.advance()
			return
		case tokens.For, tokens.Fn, tokens.If, tokens.Repeat, tokens.Tick,
			tokens.Return, tokens.Let, tokens.While, tokens.Pub, tokens.Const,
			tokens.Break, tokens.Continue, tokens.Add, tokens.Remove,
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

// Advances by one into the next token and returns the previous token before advancing
func (p *Parser) advance() tokens.Token {
	t := p.tokens[p.current]
	if p.current < len(p.tokens)-1 {
		p.current++
	}
	return t
}

// Advances by one into the next token and returns the previous token before advancing
func (p *Parser) disadvance(amount int) tokens.Token {
	if p.current > 0 {
		p.current -= amount
	}
	return p.tokens[p.current]
}

func (p *Parser) getCurrent() int {
	return p.current
}

// Peeks into the current token or peeks at the token that is offset from the current position by the given offset
func (p *Parser) peek(offset ...int) tokens.Token {
	if offset == nil {
		return p.tokens[p.current]
	} else {
		if p.current+offset[0] >= len(p.tokens)-1 || p.current+offset[0] < 1 {
			return p.tokens[p.current]
		}
		return p.tokens[p.current+offset[0]]
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

// Takes a list of tokens, advancing if the next token matches with any token from the list and returns true.
// Consume also advances if none of the tokens were able to match, and returns false
func (p *Parser) consumeOld(message string, types ...tokens.TokenType) (tokens.Token, bool) {
	if p.isAtEnd() {
		token := p.peek()
		p.error(token, message)
		return token, false // error
	}
	for _, tokenType := range types {
		if p.check(tokenType) {
			return p.advance(), true
		}
	}
	token := p.advance()
	p.error(token, message)
	return token, false // error
}

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
	failed := p.GetEnv()
	if failed {
		return []ast.Node{}
	}

	for !p.isAtEnd() {
		statement := p.statement()
		if statement == nil {
			continue
		}
		if statement.GetType() != ast.NA {
			p.program = append(p.program, statement)
		}
	}

	return p.program
}

func (p *Parser) GetEnv() bool {
	defer func() {
		if errMsg := recover(); errMsg != nil {
			if _, ok := errMsg.(ast.Error); ok {
			} else if _, ok := errMsg.(ParserError); ok {
			} else {
				panic(errMsg)
			}
		}
	}()

	if p.peek().Type != tokens.Env {
		p.Alert_(&alerts.ExpectedEnvironment{}, p.peek(), p.peek().Location)
		// unsynchronizable error.
		// if there is no env you cannot know which numbers are allowed
		return true
	}
	envStmt := p.statement()
	if envStmt.GetType() == ast.NA {
		return true
	}

	p.program = append(p.program, envStmt)

	return false
}

func (p *Parser) getBody() ([]ast.Node, bool) {
	body := make([]ast.Node, 0)
	if p.match(tokens.FatArrow) {
		args, ok := p.returnArgs()
		if !ok {
			p.error(p.peek(), "expected return arguments")
			return []ast.Node{}, false
		}
		body = []ast.Node{
			&ast.ReturnStmt{
				Token: args[0].GetToken(),
				Args:  args,
			},
		}
		return body, true
	} else if !p.check(tokens.LeftBrace) {
		body = []ast.Node{p.statement()}
		return body, true
	}
	if _, success := p.consumeOld("expected opening of the body", tokens.LeftBrace); !success {
		return body, false
	}
	start := p.peek(-1)

	for !p.match(tokens.RightBrace) { // passed that
		if p.peek().Type == tokens.Eof { // i say we debug and see the token content
			p.Alert(&alerts.ExpectedEnclosingMark{}, alerts.Multiline{StartToken: start, EndToken: p.peek(-1)}, string(tokens.RightBrace)) // no,
			//p.error(p.peek(), "expected body closure")// so we generate expected body closure error
			return body, false
		}

		statement := p.statement()
		if statement.GetType() != ast.NA {
			body = append(body, statement)
		}
	}

	return body, true
}
