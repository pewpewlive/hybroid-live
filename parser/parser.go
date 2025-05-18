package parser

import (
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/helpers"
	"hybroid/tokens"
)

type Parser struct {
	alerts.Collector

	program []ast.Node
	current int
	tokens  []tokens.Token
	context ParserContext
}

type ParserContext struct {
	EnvDeclaration *ast.EnvironmentDecl
	IsPub          bool
	IgnoreAlerts   helpers.Stack[bool]
	BraceEntries   helpers.Stack[bool]
}

func NewParser(tokens []tokens.Token) Parser {
	parser := Parser{
		program: make([]ast.Node, 0),
		current: 0,
		tokens:  tokens,
		context: ParserContext{
			EnvDeclaration: nil,
			IgnoreAlerts:   helpers.NewStack[bool]("IgnoreAlerts"),
			BraceEntries:   helpers.NewStack[bool]("BraceEntries"),
		},
		Collector: alerts.NewCollector(),
	}

	parser.context.IgnoreAlerts.Push("default", false)

	return parser
}

type ParserError struct{}

func (p *Parser) Alert(alertType alerts.Alert, args ...any) {
	if p.context.IgnoreAlerts.Top().Item {
		return
	}

	p.Alert_(alertType, args...)
}

func (p *Parser) AlertI(alert alerts.Alert) {
	if p.context.IgnoreAlerts.Top().Item {
		return
	}

	p.AlertI_(alert)
}

func (p *Parser) Parse() []ast.Node {
	for !p.isAtEnd() {
		declaration := p.declaration()
		if declaration == nil {
			continue
		}
		if declaration.GetType() != ast.NA {
			p.program = append(p.program, declaration)
			continue
		}
	}

	return p.program
}
