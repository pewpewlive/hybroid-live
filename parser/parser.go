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
	context parserContext
}

type parserContext struct {
	isPub        bool
	ignoreAlerts helpers.Stack[bool]
	braceCounter helpers.Counter
}

func NewParser(tokens []tokens.Token) Parser {
	parser := Parser{
		program: make([]ast.Node, 0),
		current: 0,
		tokens:  tokens,
		context: parserContext{
			ignoreAlerts: helpers.NewStack[bool]("IgnoreAlerts"),
			braceCounter: helpers.NewCounter("BraceCounter"),
		},
		Collector: alerts.NewCollector(),
	}

	parser.context.ignoreAlerts.Push("default", false)

	return parser
}

type ParserError struct{}

func (p *Parser) Alert(alertType alerts.Alert, args ...any) {
	if p.context.ignoreAlerts.Top().Item {
		return
	}

	p.Alert_(alertType, args...)
}

func (p *Parser) AlertI(alert alerts.Alert) {
	if p.context.ignoreAlerts.Top().Item {
		return
	}

	p.AlertI_(alert)
}

func (p *Parser) Parse() []ast.Node {
	for !p.isAtEnd() {
		declaration := p.bodyNode(p.synchronizeBody)
		if declaration == nil {
			continue
		}
		if ast.IsImproperNotStatement(declaration) {
			p.Alert(&alerts.UnknownStatement{}, alerts.NewSingle(declaration.GetToken()))
			continue
		}
		if declaration.GetType() != ast.NA {
			p.program = append(p.program, declaration)
			continue
		}
	}

	return p.program
}
