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
		node := p.parseNode(p.synchronizeBody)
		if node == nil {
			continue
		}
		if ast.IsImproperNotStatement(node) {
			p.Alert(&alerts.UnknownStatement{}, alerts.NewSingle(node.GetToken()))
			continue
		}
		if node.GetType() != ast.NA {
			p.program = append(p.program, node)
			continue
		}
	}

	return p.program
}

func (p *Parser) parseNode(syncFunc func()) (returnNode ast.Node) {
	returnNode = ast.NewImproper(p.peek(), ast.NA)
	p.context.isPub = false

	defer func() {
		p.context.isPub = false
		if returnNode.GetType() == ast.NA {
			syncFunc()
		}
	}()

	if p.match(tokens.Env) {
		returnNode = p.environmentDeclaration()
		return
	}

	if p.match(tokens.Pub) {
		p.context.isPub = true
	}

	if p.peek().Type == tokens.Entity && p.peek(1).Type == tokens.Identifier && p.peek(2).Type == tokens.LeftBrace {
		p.advance()
		returnNode = p.entityDeclaration()
		return
	}

	current := p.current
	p.context.ignoreAlerts.Push("VariableDeclaration", true)
	node := p.variableDeclaration(false)
	p.context.ignoreAlerts.Pop("VariableDeclaration")
	p.disadvance(p.current - current)
	if !ast.IsImproper(node, ast.NA) {
		returnNode = p.variableDeclaration(false)
		return
	}

	switch {
	case p.match(tokens.Let) || p.match(tokens.Const):
		returnNode = p.variableDeclaration(true)
	case p.match(tokens.Fn):
		returnNode = p.functionDeclaration()
	case p.match(tokens.Enum):
		returnNode = p.enumDeclaration()
	case p.match(tokens.Class):
		returnNode = p.classDeclaration()
	case p.match(tokens.Alias):
		returnNode = p.aliasDeclaration()
	default:
		if p.context.isPub {
			p.Alert(&alerts.UnexpectedKeyword{}, alerts.NewSingle(p.peek(-1)), tokens.Pub, "before statement")
			p.context.isPub = false
		}

		returnNode = p.statement()
	}

	p.context.isPub = false

	if ast.IsImproper(returnNode, ast.NA) {
		current := p.current
		p.context.ignoreAlerts.Push("ExpressionStatement", true)
		node := p.expressionStatement()
		p.context.ignoreAlerts.Pop("ExpressionStatement")
		p.disadvance(p.current - current)
		if !ast.IsImproper(node, ast.NA) {
			returnNode = p.expressionStatement()
		}
	}

	return
}
