package parser

import (
	"fmt"
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/tokens"
	"os"
	"runtime/debug"
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

func (p *Parser) declaration() ast.Node {
	defer func() {
		p.Context.IsPub = false

		if errMsg := recover(); errMsg != nil {
			// If the error is a parseError, synchronize to
			// the next statement. If not, propagate the panic.
			if _, ok := errMsg.(ParserError); ok {
				p.synchronize()
				return
			} else {
				fmt.Printf("panic: %s\nstacktrace:\n", errMsg)
				debug.PrintStack()
				os.Exit(1)
			}
		}
	}()
	p.Context.IsPub = false

	if p.match(tokens.Env) {
		if p.environmentDeclaration().GetType() == ast.EnvironmentDeclaration {
			p.AlertPanic(&alerts.EnvironmentRedaclaration{}, alerts.NewSingle(p.peek()))
		}
	}

	if p.match(tokens.Pub) {
		p.Context.IsPub = true
	}

	switch {
	case p.match(tokens.Fn):
		return p.functionDeclaration()
	case p.check(tokens.Let) || p.check(tokens.Const):
		return p.variableDeclaration()
	case p.match(tokens.Enum):
		return p.enumDeclaration()
	case p.match(tokens.Alias):
		return p.aliasDeclaration()
	case p.match(tokens.Class):
		return p.classDeclaration()
	case p.match(tokens.Entity):
		return p.entityDeclaration()
	}

	return p.statement()
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
	} else if p.match(tokens.Let) {
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

	expr := p.envPathExpr()
	if expr.GetType() == ast.NA {
		return expr
	}

	if _, ok := p.consume(p.NewAlert(&alerts.ExpectedKeyword{}, alerts.NewSingle(p.peek()), tokens.As), tokens.As); !ok {
		return ast.NewImproper(p.peek(), ast.EnvironmentDeclaration)
	}

	envTypeExpr := p.envTypeExpr()
	if envTypeExpr.Type == ast.InvalidEnv {
		return envDecl
	}

	envPathExpr, _ := expr.(*ast.EnvPathExpr)
	envDecl.EnvType = envTypeExpr
	envDecl.Env = envPathExpr
	p.Context.EnvDeclaration = envDecl

	return envDecl
}

func (p *Parser) variableDeclaration() ast.Node {
	variableDecl := ast.VariableDecl{
		Token:   p.peek(),
		IsPub:   p.Context.IsPub,
		IsConst: false,
	}

	if variableDecl.IsPub {
		variableDecl.Token = p.peek(-1)
	}

	if p.match(tokens.Const) {
		variableDecl.IsConst = true
	} else if p.match(tokens.Let) && variableDecl.IsPub {
		p.Alert(&alerts.UnexpectedKeyword{}, alerts.NewSingle(p.peek(-1)), variableDecl.Token.Lexeme, "in variable declaration")
	}

	typeCheckStart := p.current
	if typeExpr, ok := p.checkType(); ok {
		variableDecl.Type = typeExpr
		if !p.check(tokens.Identifier) {
			p.disadvance(p.current - typeCheckStart)
			variableDecl.Type = nil
		}
	}

	idents, exprs, ok := p.identExprPairs("in variable declaration", variableDecl.Type != nil)
	if !ok {
		p.panic()
	}

	variableDecl.Identifiers = idents
	variableDecl.Expressions = exprs

	return &variableDecl
}

func (p *Parser) functionDeclaration() ast.Node {
	functionDecl := ast.FunctionDecl{
		IsPub: p.Context.IsPub,
	}

	name, _ := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "as the name of the function"), tokens.Identifier)

	functionDecl.Name = name
	functionDecl.Generics = p.genericParams()
	functionDecl.Params = p.functionParams(tokens.LeftParen, tokens.RightParen)
	functionDecl.Return = p.functionReturns()

	if functionDecl.Return == nil || functionDecl.Return.Name.GetType() == ast.NA {
		p.Context.FunctionReturns.Push("functionDeclarationStmt", 0)
	} else if functionDecl.Return.Name.GetType() == ast.TupleExpression {
		p.Context.FunctionReturns.Push("functionDeclarationStmt", len(functionDecl.Return.Name.(*ast.TupleExpr).Types))
	} else {
		p.Context.FunctionReturns.Push("functionDeclarationStmt", 1)
	}

	body, ok := p.body(false, true)
	if !ok {
		return ast.NewImproper(functionDecl.Name, ast.FunctionDeclaration)
	}
	functionDecl.Body = body

	p.Context.FunctionReturns.Pop("functionDeclarationStmt")

	return &functionDecl
}

func (p *Parser) enumDeclaration() ast.Node {
	enumStmt := &ast.EnumDecl{
		IsPub: p.Context.IsPub,
		Token: p.peek(-1),
	}

	name, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "in enum declaration"), tokens.Identifier)
	if ok {
		enumStmt.Name = name
	}

	start, ok := p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.LeftBrace), tokens.LeftBrace)
	if !ok {
		p.panic()
		return ast.NewImproper(p.peek(), ast.EnumDeclaration)
	}

	fields, ok := p.identifiers("in enum declaration", true)
	if !ok {
		p.panic()
		return enumStmt
	}

	_, ok = p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewMulti(start, p.peek()), tokens.RightBrace), tokens.RightBrace)
	if !ok {
		p.panic()
	}
	enumStmt.Fields = fields

	return enumStmt
}

func (p *Parser) aliasDeclaration() ast.Node {
	aliasDecl := &ast.AliasDecl{
		IsPub: p.Context.IsPub,
		Token: p.peek(-1),
	}

	name, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "as the name in alias declaration"), tokens.Identifier)
	if ok {
		aliasDecl.Name = name
	}

	p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.Equal, "after identifier in alias declaration"), tokens.Equal)

	aliasDecl.Type = p.typeExpr()

	return aliasDecl
}

func (p *Parser) classDeclaration() ast.Node {
	stmt := &ast.ClassDecl{
		IsPub: p.Context.IsPub,
		Token: p.peek(-1),
	}

	name, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "as the name of the class"), tokens.Identifier)

	if ok {
		stmt.Name = name
	} else {
		return ast.NewImproper(stmt.Token, ast.ClassDeclaration)
	}

	_, ok = p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.LeftBrace), tokens.LeftBrace)
	if !ok {
		return ast.NewImproper(stmt.Token, ast.ClassDeclaration)
	}
	stmt.Methods = []ast.MethodDecl{}
	for !p.match(tokens.RightBrace) {
		auxiliaryDeclaration := p.auxiliaryDeclaration()
		switch declaration := auxiliaryDeclaration.(type) {
		case *ast.ConstructorDecl:
			if stmt.Constructor != nil {
				p.Alert(&alerts.MoreThanOneConstructor{}, alerts.NewSingle(p.peek()))
			} else {
				stmt.Constructor = declaration
			}
		case *ast.FieldDecl:
			stmt.Fields = append(stmt.Fields, *declaration)
		case *ast.MethodDecl:
			stmt.Methods = append(stmt.Methods, *declaration)
		default:
			p.Alert(&alerts.UnknownStatement{}, alerts.NewSingle(p.peek()), "in class declaration")
		}
	}

	return stmt
}

func (p *Parser) entityDeclaration() ast.Node {
	stmt := &ast.EntityDecl{
		IsPub: p.Context.IsPub,
		Token: p.peek(-1),
	}

	name, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "as the name of the entity"), tokens.Identifier)

	if !ok {
		return ast.NewImproper(stmt.Token, ast.EntityDeclaration)
	}
	stmt.Name = name

	_, ok = p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.LeftBrace), tokens.LeftBrace)
	if !ok {
		return ast.NewImproper(stmt.Token, ast.EntityDeclaration)
	}

	for !p.match(tokens.RightBrace) {
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
			p.Alert(&alerts.UnknownStatement{}, alerts.NewMulti(auxiliaryDeclaration.GetToken(), p.peek()), "in entity declaration")
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
	stmt.Params = p.functionParams(tokens.LeftParen, tokens.RightParen)
	stmt.Return = p.functionReturns()
	if stmt.Return == nil || stmt.Return.Name.GetType() == ast.NA {
		p.Context.FunctionReturns.Push("entityFunctionDeclarationStmt", 0)
	} else if stmt.Return.Name.GetType() == ast.TupleExpression {
		p.Context.FunctionReturns.Push("entityFunctionDeclarationStmt", len(stmt.Return.Name.(*ast.TupleExpr).Types))
	} else {
		p.Context.FunctionReturns.Push("entityFunctionDeclarationStmt", 1)
	}

	var success bool
	stmt.Body, success = p.body(true, true)
	if !success {
		return ast.NewImproper(stmt.Token, ast.EntityFunctionDeclaration)
	}

	p.Context.FunctionReturns.Pop("entityFunctionDeclarationStmt")

	return stmt
}

func (p *Parser) constructorDeclaration() ast.Node {
	stmt := &ast.ConstructorDecl{Token: p.peek(-1)}

	stmt.Generics = p.genericParams()
	stmt.Params = p.functionParams(tokens.LeftParen, tokens.RightParen)
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
	if typeExpr, ok := p.checkType(); ok {
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
