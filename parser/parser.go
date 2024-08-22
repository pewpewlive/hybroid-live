package parser

import (
	"hybroid/ast"
	"hybroid/lexer"
)

type Parser struct {
	program []ast.Node
	current int
	tokens  []lexer.Token
	Errors  []ast.Error
	Context ParserContext
}

type ParserContext struct {
	FunctionReturns []int
}

func NewParser() Parser {
	return Parser{
		program: make([]ast.Node, 0),
		current: 0,
		tokens: make([]lexer.Token, 0),
		Errors: make([]ast.Error, 0),
		Context: ParserContext{
			FunctionReturns: make([]int, 0),
		},
	}
}

func (p *Parser) AssignTokens(tokens []lexer.Token) {
	p.tokens = tokens
}

// Appends an error to the ParserErrors
func (p *Parser) error(token lexer.Token, msg string) {
	errMsg := ast.Error{Token: token, Message: msg}
	p.Errors = append(p.Errors, errMsg)
	//panic(errMsg.Message)
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		switch p.peek().Type {
		case lexer.For, lexer.Fn, lexer.If, lexer.Repeat, lexer.Tick,
			lexer.Return, lexer.Let, lexer.While, lexer.Pub, lexer.Const,
			lexer.Break, lexer.Continue, lexer.Add, lexer.Remove,
			lexer.Class:
			return
		case lexer.Entity:
			if p.peek(1).Type == lexer.Identifier && p.peek(2).Type == lexer.LeftBrace {
				return
			}
		} // pub fn (entity thing) { }
		// entity yes{}

		p.advance()
	}
}

func (p *Parser) isMultiComparison() bool {
	return p.match(lexer.And, lexer.Or)
}

func (p *Parser) isComparison() bool {
	return p.match(lexer.Greater, lexer.GreaterEqual, lexer.Less, lexer.LessEqual, lexer.BangEqual, lexer.EqualEqual)
}

// Checks if the current position the parser is at is the End Of File
func (p *Parser) isAtEnd() bool {
	return p.peek().Type == lexer.Eof
}

// Advances by one into the next token and returns the previous token before advancing
func (p *Parser) advance() lexer.Token {
	t := p.tokens[p.current]
	if p.current < len(p.tokens)-1 {
		p.current++
	}
	return t
}

// Advances by one into the next token and returns the previous token before advancing
func (p *Parser) disadvance(amount int) lexer.Token {
	if p.current > 0 {
		p.current -= amount
	}
	return p.tokens[p.current]
}

func (p *Parser) getCurrent() int {
	return p.current
}

// Peeks into the current token or peeks at the token that is offset from the current position by the given offset
func (p *Parser) peek(offset ...int) lexer.Token {
	if offset == nil {
		return p.tokens[p.current]
	} else {
		if p.current+offset[0] >= len(p.tokens)-1 {
			return p.tokens[p.current]
		}
		return p.tokens[p.current+offset[0]]
	}
}

// Checks if the current type is the specified token type. Returns false if it's the End Of File
func (p *Parser) check(tokenType lexer.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == tokenType
}

// Matches the given list of tokens and advances if they match.
func (p *Parser) match(types ...lexer.TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

// Consumes a list of tokens, advancing if they match and returns true. Consume also advances if none of the tokens were able to match, and returns false
func (p *Parser) consume(message string, types ...lexer.TokenType) (lexer.Token, bool) {
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

func (p *Parser) ParseTokens() []ast.Node {
	statement := p.statement()
	p.program = append(p.program, statement)
	if statement.GetType() != ast.EnvironmentStatement {
		p.error(statement.GetToken(), "expected environment as the first statement in a file")
	}
	for !p.isAtEnd() {
		statement := p.statement()
		if statement.GetType() != ast.NA {
			p.program = append(p.program, statement)
		}
	}

	return p.program
}

func (p *Parser) getBody() ([]ast.Node, bool) {
	body := make([]ast.Node, 0)
	if p.match(lexer.FatArrow) {
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
	} else if !p.check(lexer.LeftBrace) {
		body = []ast.Node{p.statement()}
		return body, true
	}
	if _, success := p.consume("expected opening of the body", lexer.LeftBrace); !success {
		return body, false
	}

	for !p.match(lexer.RightBrace) { // passed that
		if p.peek().Type == lexer.Eof {
			p.error(p.peek(), "expected body closure")
			return body, false
		}

		statement := p.statement()
		if statement.GetType() != ast.NA {
			body = append(body, statement)
		}
	}

	return body, true
}
