package parser

import (
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
	IsPub           bool
	IgnoreAlerts    helpers.Stack[bool]
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

func (p *Parser) panic() {
	panic(ParserError{})
}

func (p *Parser) synchronize() {
	for !p.isAtEnd() {
		switch p.peek().Type {
		// case tokens.For, tokens.If, tokens.Repeat, tokens.Tick,
		// 	tokens.Return, , tokens.While, , ,
		// 	tokens.Break, tokens.Continue,
		// 	:
		// 	return
		case tokens.Entity, tokens.Fn, tokens.Let, tokens.Pub, tokens.Const, tokens.Class:
			return
		}

		p.advance()
	}
}

func (p *Parser) Parse() []ast.Node {
	if p.match(tokens.Env) {
		envDecl := p.environmentDeclaration()
		if envDecl.GetType() != ast.EnvironmentDeclaration {
			return []ast.Node{}
		}
		p.program = append(p.program, envDecl)
	} else {
		p.Alert(&alerts.ExpectedEnvironment{}, alerts.NewSingle(p.peek()))
		return []ast.Node{}
	}

	for !p.isAtEnd() {
		statement := p.declaration()
		if statement == nil {
			continue
		}
		if statement.GetType() != ast.NA {
			p.program = append(p.program, statement)
		}
	}

	return p.program
}
