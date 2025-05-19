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

func (p *Parser) declaration() (returnNode ast.Node) {
	returnNode = ast.NewImproper(p.peek(), ast.NA)
	p.context.IsPub = false

	if p.match(tokens.Env) {
		return p.environmentDeclaration()
	}

	if p.match(tokens.Pub) {
		p.context.IsPub = true
	}

	switch {
	case p.match(tokens.Fn):
		returnNode = p.functionDeclaration()
	case p.match(tokens.Enum):
		returnNode = p.enumDeclaration()
	case p.match(tokens.Alias):
		returnNode = p.aliasDeclaration()
	case p.match(tokens.Class):
		returnNode = p.classDeclaration()
	case p.match(tokens.Entity):
		returnNode = p.entityDeclaration()
	case p.check(tokens.Let) || p.check(tokens.Const) || p.context.IsPub:
		returnNode = p.variableDeclaration(false)
	default:
		returnNode = p.statement()
	}

	if returnNode.GetType() == ast.NA && p.peekTypeVariableDecl() && !p.context.IsPub {
		returnNode = p.variableDeclaration(true)
	} else if returnNode.GetType() == ast.NA {
		current := p.current
		p.context.IgnoreAlerts.Push("ExpressionStatement", true)
		node := p.expressionStatement()
		p.context.IgnoreAlerts.Push("ExpressionStatement", false)
		p.disadvance(p.current - current)
		if !ast.IsImproper(node, ast.NA) {
			returnNode = p.expressionStatement()
		}
	}

	p.context.IsPub = false
	if returnNode.GetType() == ast.NA {
		p.synchronize()
	}

	return
}

func (p *Parser) auxiliaryDeclaration() ast.Node {
	if p.match(tokens.Fn) {
		fnDec := p.functionDeclaration()

		if ast.IsImproper(fnDec, ast.FunctionDeclaration) {
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
		return p.entityFunctionDeclaration(p.peek(-1), functionType)
	}

	// No auxiliary declaration found, try normal declaration anyway
	return p.declaration()
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

func (p *Parser) variableDeclaration(bypassKeyword bool) ast.Node {
	variableDecl := ast.VariableDecl{
		Token:   p.peek(),
		IsPub:   p.context.IsPub,
		IsConst: false,
	}

	if variableDecl.IsPub {
		variableDecl.Token = p.peek(-1)
	}

	if !bypassKeyword && p.match(tokens.Const) {
		variableDecl.IsConst = true
	} else if !bypassKeyword && p.match(tokens.Let) && variableDecl.IsPub {
		p.Alert(&alerts.UnexpectedKeyword{}, alerts.NewSingle(p.peek(-1)), variableDecl.Token.Lexeme, "in variable declaration")
	}

	typeCheckStart := p.current
	if typeExpr, ok := p.checkType("in variable declaration"); ok {
		variableDecl.Type = typeExpr
		if !p.check(tokens.Identifier) {
			p.disadvance(p.current - typeCheckStart)
			variableDecl.Type = nil
		}
	}

	idents, exprs, ok := p.identExprPairs("in variable declaration", variableDecl.Type != nil)
	if !ok {
		return ast.NewImproper(variableDecl.Token, ast.VariableDeclaration)
	}

	variableDecl.Identifiers = idents
	variableDecl.Expressions = exprs

	return &variableDecl
}

func (p *Parser) functionDeclaration() ast.Node {
	functionDecl := ast.FunctionDecl{
		IsPub: p.context.IsPub,
	}

	name, nameOk := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "as the name of the function"), tokens.Identifier)
	if !nameOk && (!p.check(tokens.Less) || !p.check(tokens.LeftParen)) {
		p.advance()
	}

	functionDecl.Name = name
	functionDecl.Generics = p.genericParams()
	functionDecl.Params, _ = p.functionParams(tokens.LeftParen, tokens.RightParen)
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
		IsPub: p.context.IsPub,
		Token: p.peek(-1),
	}

	name, _ := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "in enum declaration"), tokens.Identifier)
	enumStmt.Name = name

	start, ok := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.LeftBrace), tokens.LeftBrace)
	if !ok {
		return ast.NewImproper(enumStmt.Token, ast.EnumDeclaration)
	}

	fields, _ := p.expressions("in enum declaration", true)
	for _, v := range fields {
		if v.GetType() == ast.Identifier {
			enumStmt.Fields = append(enumStmt.Fields, v.(*ast.IdentifierExpr))
		} else {
			p.Alert(&alerts.InvalidEnumVariantName{}, alerts.NewSingle(v.GetToken()))
			continue
		}
	}

	p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewMulti(start, p.peek()), tokens.RightBrace), tokens.RightBrace)

	return enumStmt
}

func (p *Parser) aliasDeclaration() ast.Node {
	aliasDecl := &ast.AliasDecl{
		IsPub: p.context.IsPub,
		Token: p.peek(-1),
	}

	name, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "as the name in alias declaration"), tokens.Identifier)
	if ok {
		aliasDecl.Name = name
	} else if !p.check(tokens.Equal) {
		p.advance()
	}

	p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.Equal, "after name in alias declaration"), tokens.Equal)

	aliasDecl.Type, ok = p.checkType("in alias declaration")
	if !ok {
		p.Alert(&alerts.ExpectedType{}, alerts.NewSingle(aliasDecl.Type.GetToken()), "to be aliased in alias declaration")
	}

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
	start := p.peek(-1)
	stmt.Methods = []ast.MethodDecl{}
	for p.doesntEndWith("in class declaration", start, tokens.RightBrace) {
		auxiliaryDeclaration := p.auxiliaryDeclaration()
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
	start := p.peek(-1)

	for p.doesntEndWith("in entity declaration", start, tokens.RightBrace) {
		auxiliaryDeclaration := p.auxiliaryDeclaration()
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

	stmt.Generics = p.genericParams()
	stmt.Params, _ = p.functionParams(tokens.LeftParen, tokens.RightParen)
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

	stmt.Generics = p.genericParams()
	stmt.Params, _ = p.functionParams(tokens.LeftParen, tokens.RightParen)
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
