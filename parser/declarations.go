package parser

import (
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/tokens"
)

/*
declarations:

declaration     → classDecl
				| entityDecl
				| funDecl
				| varDecl
				| statement ;

classDecl       → "class" IDENTIFIER ( "<" IDENTIFIER )?
				  "{" function* "}" ;

funDecl        → "fun" function ;
varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;
*/

/*
statements:

statement      → exprStmt

	| forStmt
	| ifStmt
	| printStmt
	| returnStmt
	| whileStmt
	| block ;

exprStmt       → expression ";" ;
forStmt        → "for" "(" ( varDecl | exprStmt | ";" )

	expression? ";"
	expression? ")" statement ;

ifStmt         → "if" "(" expression ")" statement

	( "else" statement )? ;

printStmt      → "print" expression ";" ;
returnStmt     → "return" expression? ";" ;
whileStmt      → "while" "(" expression ")" statement ;
assignmentStmt → ( call "." )? IDENTIFIER "="

block          → "{" declaration* "}" ;

expressions:

expression     → logic_or ;

logic_or       → logic_and ( "or" logic_and )* ;
logic_and      → equality ( "and" equality )* ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;

unary          → ( "!" | "-" ) unary | call ;
call           → primary ( "(" arguments? ")" | "." IDENTIFIER )* ;
primary        → "true" | "false" | "nil" | "this"

	| NUMBER | STRING | IDENTIFIER | "(" expression ")"
	| "super" "." IDENTIFIER ;

utility:

function       → IDENTIFIER "(" parameters? ")" block ;
parameters     → IDENTIFIER ( "," IDENTIFIER )* ;
arguments      → expression ( "," expression )* ;
*/

func (p *Parser) bodyNode(syncFunc func()) (returnNode ast.Node) {
	returnNode = ast.NewImproper(p.peek(), ast.NA)
	p.context.IsPub = false

	defer func() {
		p.context.IsPub = false
		if returnNode.GetType() == ast.NA {
			syncFunc()
		}
	}()

	if p.match(tokens.Env) {
		returnNode = p.environmentDeclaration()
		return
	}

	if p.match(tokens.Pub) {
		p.context.IsPub = true
	}

	if p.peek().Type == tokens.Entity && p.peek(1).Type == tokens.Identifier && p.peek(2).Type == tokens.LeftBrace {
		p.advance()
		returnNode = p.entityDeclaration()
		return
	}

	current := p.current
	p.context.IgnoreAlerts.Push("VariableDeclaration", true)
	node := p.variableDeclaration(false)
	p.context.IgnoreAlerts.Pop("VariableDeclaration")
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
		if p.context.IsPub {
			p.Alert(&alerts.UnexpectedKeyword{}, alerts.NewSingle(p.peek(-1)), tokens.Pub, "before statement")
			p.context.IsPub = false
		}

		returnNode = p.statement()
	}

	p.context.IsPub = false

	if ast.IsImproper(returnNode, ast.NA) {
		current := p.current
		p.context.IgnoreAlerts.Push("ExpressionStatement", true)
		node := p.expressionStatement()
		p.context.IgnoreAlerts.Pop("ExpressionStatement")
		p.disadvance(p.current - current)
		if !ast.IsImproper(node, ast.NA) {
			returnNode = p.expressionStatement()
		}
	}

	return
}

func (p *Parser) declBodyNode() ast.Node {
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
			Return:   fnDecl.Return,
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
	} else if p.match(tokens.Let) || p.peekTypeVariableDecl() {
		field := p.fieldDeclaration()
		if field.GetType() == ast.FieldDeclaration {
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

	// No auxiliary declaration found, try normal declaration anyway
	return p.bodyNode(p.synchronizeDeclBody)
}

func (p *Parser) environmentDeclaration() ast.Node {
	envDecl := &ast.EnvironmentDecl{}

	pathExpr := p.envPathExpr()
	if pathExpr.GetType() == ast.NA {
		return pathExpr
	}

	if _, ok := p.consume(p.NewAlert(&alerts.ExpectedKeyword{}, alerts.NewSingle(p.peek()), tokens.As), tokens.As); !ok {
		return ast.NewImproper(p.peek(), ast.EnvironmentDeclaration)
	}

	envPathExpr, _ := pathExpr.(*ast.EnvPathExpr)

	envDecl.EnvType = p.envTypeExpr()
	envDecl.Env = envPathExpr
	p.context.EnvDeclaration = envDecl

	return envDecl
}

func (p *Parser) variableDeclaration(matchedLetOrConst bool) ast.Node {
	variableDecl := ast.VariableDecl{
		IsPub:   p.context.IsPub,
		IsConst: false,
	}

	// check if variableDecl started with let or const
	if matchedLetOrConst {
		variableDecl.Token = p.peek(-1)
		variableDecl.IsConst = variableDecl.Token.Type == tokens.Const
		if variableDecl.IsPub {
			variableDecl.Token = p.peek(-2)
		}
		if variableDecl.IsPub && variableDecl.Token.Type == tokens.Let {
			p.Alert(&alerts.UnexpectedKeyword{}, alerts.NewSingle(variableDecl.Token), variableDecl.Token.Lexeme, "in variable declaration")
		}
	} else {
		if variableDecl.IsPub {
			variableDecl.Token = p.peek(-1)
		} else {
			variableDecl.Token = p.peek()
		}
	}

	typeCheckStart := p.current
	typeExpr, ok := p.checkType("in variable declaration")
	if ok {
		variableDecl.Type = typeExpr
		if !p.check(tokens.Identifier) {
			p.disadvance(p.current - typeCheckStart)
			variableDecl.Type = nil
		}
	} else if !matchedLetOrConst {
		return ast.NewImproper(variableDecl.Token, ast.NA)
	}

	idents, exprs, ok := p.identExprPairs("in variable declaration", variableDecl.Type != nil)
	if !ok {
		if matchedLetOrConst {
			return ast.NewImproper(variableDecl.Token, ast.VariableDeclaration)
		}
		if p.peek(-1).Type == tokens.Equal {
			return ast.NewImproper(variableDecl.Token, ast.VariableDeclaration)
		}
		return ast.NewImproper(variableDecl.Token, ast.NA)
	}

	variableDecl.Identifiers = idents
	variableDecl.Expressions = exprs

	return &variableDecl
}

func (p *Parser) functionDeclaration() ast.Node {
	functionDecl := ast.FunctionDecl{
		IsPub: p.context.IsPub,
	}
	if functionDecl.IsPub {
		functionDecl.Token = p.peek(-2)
	} else {
		functionDecl.Token = p.peek(-1)
	}

	name, nameOk := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "as the name of the function"), tokens.Identifier)
	if !nameOk && !p.check(tokens.Less) && !p.check(tokens.LeftParen) {
		p.advance()
	}

	functionDecl.Name = name
	generics, ok := p.genericParams()
	if !ok {
		return ast.NewImproper(functionDecl.Token, ast.FunctionDeclaration)
	}
	functionDecl.Generics = generics
	params, ok := p.functionParams(tokens.LeftParen, tokens.RightParen)
	if !ok {
		return ast.NewImproper(functionDecl.Name, ast.FunctionDeclaration)
	}
	functionDecl.Params = params
	functionDecl.Return = p.functionReturns()

	body, ok := p.body(false, true)
	if !ok {
		return ast.NewImproper(functionDecl.Name, ast.FunctionDeclaration)
	}
	functionDecl.Body = body

	return &functionDecl
}

func (p *Parser) enumDeclaration() ast.Node {
	enumStmt := &ast.EnumDecl{
		Token: p.peek(-1),
		IsPub: p.context.IsPub,
	}

	name, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "in enum declaration"), tokens.Identifier)
	if !ok {
		return ast.NewImproper(enumStmt.Token, ast.EnumDeclaration)
	}
	enumStmt.Name = name

	start, ok := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.LeftBrace), tokens.LeftBrace)
	if !ok {
		return ast.NewImproper(enumStmt.Token, ast.EnumDeclaration)
	}

	p.context.BraceEntries.Push("Enum", false)
	defer func() {
		p.context.BraceEntries.Pop("Enum")
	}()

	fields, _ := p.expressions("in enum declaration", true)
	for _, v := range fields {
		if v.GetType() == ast.Identifier {
			enumStmt.Fields = append(enumStmt.Fields, v.(*ast.IdentifierExpr))
		} else {
			p.Alert(&alerts.InvalidEnumVariantName{}, alerts.NewSingle(v.GetToken()))
			return ast.NewImproper(enumStmt.Token, ast.EnumDeclaration)
		}
	}

	_, ok = p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewMulti(start, p.peek()), tokens.RightBrace), tokens.RightBrace)
	if !ok {
		return ast.NewImproper(enumStmt.Token, ast.EnumDeclaration)
	}

	return enumStmt
}

func (p *Parser) aliasDeclaration() ast.Node {
	aliasDecl := &ast.AliasDecl{
		Token: p.peek(-1),
		IsPub: p.context.IsPub,
	}

	name, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "as the name in alias declaration"), tokens.Identifier)
	if ok {
		aliasDecl.Name = name
	} else if !p.check(tokens.Equal) && p.peek(1).Type == tokens.Equal {
		p.advance()
	} else {
		return ast.NewImproper(aliasDecl.Token, ast.AliasDeclaration)
	}

	_, ok = p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.Equal, "after name in alias declaration"), tokens.Equal)
	if !ok {
		return ast.NewImproper(aliasDecl.Token, ast.AliasDeclaration)
	}

	typ, ok := p.checkType("in alias declaration")
	if !ok {
		p.Alert(&alerts.ExpectedType{}, alerts.NewSingle(aliasDecl.Type.GetToken()), "in alias declaration")
	}
	if typ.GetType() == ast.NA {
		return ast.NewImproper(aliasDecl.Token, ast.AliasDeclaration)
	}
	aliasDecl.Type = typ

	return aliasDecl
}

func (p *Parser) classDeclaration() ast.Node {
	stmt := &ast.ClassDecl{
		IsPub: p.context.IsPub,
		Token: p.peek(-1),
	}

	name, _ := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "as the name of the class"), tokens.Identifier)
	stmt.Name = name

	_, ok := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.LeftBrace), tokens.LeftBrace)
	if !ok {
		return ast.NewImproper(stmt.Token, ast.ClassDeclaration)
	}
	p.context.BraceEntries.Push("Class", false)
	defer func() {
		p.context.BraceEntries.Pop("Class")
	}()

	start := p.peek(-1)
	stmt.Methods = []ast.MethodDecl{}
	for p.consumeTill("in class declaration", start, tokens.RightBrace) {
		auxiliaryDeclaration := p.declBodyNode()
		switch declaration := auxiliaryDeclaration.(type) {
		case *ast.ConstructorDecl:
			if stmt.Constructor != nil {
				p.Alert(&alerts.MoreThanOneConstructor{}, alerts.NewMulti(declaration.GetToken(), p.peek(-1)))
			} else {
				stmt.Constructor = declaration
			}
		case *ast.FieldDecl:
			stmt.Fields = append(stmt.Fields, *declaration)
		case *ast.MethodDecl:
			stmt.Methods = append(stmt.Methods, *declaration)
		default:
			p.Alert(&alerts.UnknownStatement{}, alerts.NewMulti(auxiliaryDeclaration.GetToken(), p.peek(-1)), "in class declaration")
		}
	}

	return stmt
}

func (p *Parser) entityDeclaration() ast.Node {
	stmt := &ast.EntityDecl{
		IsPub: p.context.IsPub,
		Token: p.peek(-1),
	}

	name := p.advance()
	stmt.Name = name

	if !p.match(tokens.LeftBrace) {
		p.disadvance(2)
		return ast.NewImproper(stmt.Token, ast.NA)
	}
	if name.Type != tokens.Identifier {
		p.Alert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "as the name of the entity")
	}
	p.context.BraceEntries.Push("Entity", false)
	defer func() {
		p.context.BraceEntries.Pop("Entity")
	}()

	start := p.peek(-1)
	for p.consumeTill("in entity declaration", start, tokens.RightBrace) {
		auxiliaryDeclaration := p.declBodyNode()
		if auxiliaryDeclaration.GetType() == ast.FieldDeclaration {
			stmt.Fields = append(stmt.Fields, *auxiliaryDeclaration.(*ast.FieldDecl))
			continue
		}
		if auxiliaryDeclaration.GetType() == ast.MethodDeclaration {
			stmt.Methods = append(stmt.Methods, *auxiliaryDeclaration.(*ast.MethodDecl))
			continue
		}
		if auxiliaryDeclaration.GetType() != ast.EntityFunctionDeclaration {
			p.Alert(&alerts.UnknownStatement{}, alerts.NewMulti(auxiliaryDeclaration.GetToken(), p.peek(-1)), "in entity declaration")
			continue
		}

		funcDecl := auxiliaryDeclaration.(*ast.EntityFunctionDecl)
		funcType := funcDecl.Type

		switch funcDecl.Type {
		case ast.Spawn:
			if stmt.Spawner != nil {
				p.Alert(&alerts.MoreThanOneEntityFunction{}, alerts.NewMulti(funcDecl.GetToken(), p.peek(-1)), string(funcType))
			} else {
				stmt.Spawner = funcDecl
			}
		case ast.Destroy:
			if stmt.Destroyer != nil {
				p.Alert(&alerts.MoreThanOneEntityFunction{}, alerts.NewMulti(funcDecl.GetToken(), p.peek(-1)), string(funcType))
			} else {
				stmt.Destroyer = funcDecl
			}
		default:
			var wasFound bool
			for i := range stmt.Callbacks {
				if stmt.Callbacks[i].Type == funcType {
					p.Alert(&alerts.MoreThanOneEntityFunction{}, alerts.NewMulti(funcDecl.GetToken(), p.peek(-1)), string(funcType))
					wasFound = true
				}
			}
			if !wasFound {
				stmt.Callbacks = append(stmt.Callbacks, funcDecl)
			}
		}

	}

	return stmt
}

func (p *Parser) entityFunctionDeclaration(token tokens.Token, functionType ast.EntityFunctionType) ast.Node {
	stmt := &ast.EntityFunctionDecl{
		Type:  functionType,
		Token: token,
	}

	generics, ok := p.genericParams()
	if !ok {
		return ast.NewImproper(stmt.Token, ast.EntityFunctionDeclaration)
	}
	stmt.Generics = generics
	params, ok := p.functionParams(tokens.LeftParen, tokens.RightParen)
	if !ok {
		return ast.NewImproper(stmt.Token, ast.EntityFunctionDeclaration)
	}
	stmt.Params = params
	stmt.Return = p.functionReturns()

	var success bool
	stmt.Body, success = p.body(true, true)
	if !success {
		return ast.NewImproper(stmt.Token, ast.EntityFunctionDeclaration)
	}

	return stmt
}

func (p *Parser) constructorDeclaration() ast.Node {
	stmt := &ast.ConstructorDecl{Token: p.peek(-1)}

	generics, ok := p.genericParams()
	if !ok {
		return ast.NewImproper(stmt.Token, ast.ConstructorDeclaration)
	}
	stmt.Generics = generics
	params, ok := p.functionParams(tokens.LeftParen, tokens.RightParen)
	if !ok {
		return ast.NewImproper(stmt.Token, ast.ConstructorDeclaration)
	}
	stmt.Params = params
	returns := p.functionReturns()
	if returns != nil {
		p.Alert(&alerts.ReturnsInConstructor{}, alerts.NewSingle(returns.GetToken()))
	}
	var success bool
	stmt.Body, success = p.body(false, true)
	if !success {
		return ast.NewImproper(stmt.Token, ast.ConstructorDeclaration)
	}

	return stmt
}

func (p *Parser) fieldDeclaration() ast.Node {
	fieldDecl := ast.FieldDecl{
		Token: p.peek(),
	}

	typeCheckStart := p.current
	if typeExpr, ok := p.checkType("in field declaration"); ok {
		fieldDecl.Type = typeExpr
		if !p.check(tokens.Identifier) {
			p.disadvance(p.current - typeCheckStart)
			fieldDecl.Type = nil
		}
	}

	idents, values, ok := p.identExprPairs("in field declaration", fieldDecl.Type != nil)
	if !ok {
		return ast.NewImproper(fieldDecl.Token, ast.FieldDeclaration)
	}

	fieldDecl.Identifiers = idents
	fieldDecl.Values = values

	return &fieldDecl
}
