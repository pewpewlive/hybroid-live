package parser

import (
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/core"
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
	ignoreAlerts core.Stack[bool]
	syncedToken  tokens.Token
}

func NewParser(tokens []tokens.Token) Parser {
	parser := Parser{
		program: make([]ast.Node, 0),
		current: 0,
		tokens:  tokens,
		context: parserContext{
			ignoreAlerts: core.NewStack[bool]("IgnoreAlerts"),
		},
		Collector: alerts.NewCollector(),
	}

	if len(tokens) != 0 {
		parser.context.syncedToken = tokens[0]
	}

	parser.context.ignoreAlerts.Push("default", false)

	return parser
}

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

func (p *Parser) AlertSingle(alert alerts.Alert, token tokens.Token, args ...any) {
	args = append([]any{alerts.NewSingle(token)}, args...)
	p.Alert(alert, args...)
}

func (p *Parser) AlertMulti(alert alerts.Alert, tokenStart, tokenEnd tokens.Token, args ...any) {
	args = append([]any{alerts.NewMulti(tokenStart, tokenEnd)}, args...)
	p.Alert(alert, args...)
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

	if p.peek().Type == tokens.Entity && p.peek(1).Type == tokens.Identifier && (p.peek(2).Type == tokens.Less || p.peek(2).Type == tokens.LeftBrace) {
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

func (p *Parser) auxiliaryNode() ast.Node {
	current := p.current
	p.context.ignoreAlerts.Push("FieldDeclaration", true)
	node := p.fieldDeclaration(false)
	p.context.ignoreAlerts.Pop("FieldDeclaration")
	p.disadvance(p.current - current)
	if !ast.IsImproper(node, ast.NA) {
		return p.fieldDeclaration(false)
	}

	if p.match(tokens.Fn) {
		fnDec := p.functionDeclaration()

		if ast.IsImproper(fnDec, ast.FunctionDeclaration) {
			p.synchronizeDeclBody()
			return ast.NewImproper(fnDec.GetToken(), ast.MethodDeclaration)
		}

		fnDecl := fnDec.(*ast.FunctionDecl)
		return &ast.MethodDecl{
			IsPub:    fnDecl.IsPub,
			Name:     fnDecl.Name,
			Returns:  fnDecl.Returns,
			Params:   fnDecl.Params,
			Generics: fnDecl.Generics,
			Body:     fnDecl.Body,
		}
	} else if p.match(tokens.New) {
		constructor := p.constructorDeclaration()
		if ast.IsImproper(constructor, ast.ConstructorDeclaration) {
			p.synchronizeDeclBody()
			return ast.NewImproper(constructor.GetToken(), ast.ConstructorDeclaration)
		}
		if constructor.GetType() == ast.ConstructorDeclaration {
			return constructor
		}
	} else if p.match(tokens.Let) {
		field := p.fieldDeclaration(true)
		if ast.IsImproper(field, ast.VariableDeclaration) {
			p.synchronizeDeclBody()
			return ast.NewImproper(field.GetToken(), ast.VariableDeclaration)
		}
		if field.GetType() == ast.VariableDeclaration {
			return field
		}
	}
	var functionType ast.EntityFunctionType = ""
	if p.match(tokens.Spawn) {
		functionType = ast.Spawn
	} else if p.match(tokens.Destroy) {
		functionType = ast.Destroy
	} else if p.check(tokens.Identifier) {
		switch p.peek().Lexeme {
		case "WeaponCollision":
			functionType = ast.WeaponCollision
		case "WallCollision":
			functionType = ast.WallCollision
		case "PlayerCollision":
			functionType = ast.PlayerCollision
		case "Update":
			functionType = ast.Update
		}
		if functionType != "" {
			p.advance()
		}
	}
	if functionType != "" {
		entityFunction := p.entityFunctionDeclaration(p.peek(-1), functionType)
		if ast.IsImproper(entityFunction, ast.EntityFunctionDeclaration) {
			p.synchronizeDeclBody()
			return ast.NewImproper(entityFunction.GetToken(), ast.EntityFunctionDeclaration)
		}
		return entityFunction
	}

	// No auxiliary node found, try parsing a node instead (for error handling)
	return p.parseNode(p.synchronizeDeclBody)
}
