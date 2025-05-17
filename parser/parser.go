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
}

func NewParser(tokens []tokens.Token) Parser {
	parser := Parser{
		program: make([]ast.Node, 0),
		current: 0,
		tokens:  tokens,
		context: ParserContext{
			EnvDeclaration: nil,
			IgnoreAlerts:   helpers.NewStack[bool]("IgnoreAlerts"),
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

func IsCall(nodeType ast.NodeType) bool {
	return nodeType == ast.CallExpression ||
		nodeType == ast.MethodCallExpression ||
		nodeType == ast.NewExpession ||
		nodeType == ast.SpawnExpression
}

func (p *Parser) synchronize() {
	expectedBlockCount := 0
	for !p.isAtEnd() {
		switch p.peek().Type {
		case tokens.Fn:
			p.advance()
			if p.peek().Type != tokens.LeftParen {
				p.disadvance()
				return
			}
		case tokens.LeftBrace:
			expectedBlockCount++
		case tokens.RightBrace:
			if expectedBlockCount == 0 {
				return
			}
			expectedBlockCount--
		case tokens.Entity, tokens.Let, tokens.Pub, tokens.Const, tokens.Class, tokens.Alias:
			return
		default:
			current := p.current
			p.context.IgnoreAlerts.Push("Synchronize", true)

			expr := p.expression()
			exprType := expr.GetType()
			if exprType == ast.NA {
				exprType = expr.(*ast.Improper).Type
			}

			p.disadvance(p.current - current)
			p.context.IgnoreAlerts.Pop("Synchronize")

			if IsCall(exprType) {
				return
			}
		}

		p.advance()
	}
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
