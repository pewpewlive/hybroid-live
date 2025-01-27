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

func (p *Parser) declaration() ast.Node {
	if p.match(tokens.Env) {
		p.Alert(&alerts.EnvironmentRedaclaration{}, alerts.NewSingle(p.peek()))
	}

	if varDecl := p.variableDeclaration(); varDecl != nil {
		return varDecl
	}

	if p.match(tokens.Fn) {
		return p.functionDeclaration()
	}

	return nil
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
		IsPub:   false,
		IsConst: false,
	}

	switch variableDecl.Token.Type {
	case tokens.Const:
		variableDecl.IsConst = true
	case tokens.Pub:
		variableDecl.IsPub = true
		if p.peek(1).Type == tokens.Const {
			variableDecl.IsConst = true
		}
	}

	currentStart := p.current
	p.match(tokens.Let, tokens.Const, tokens.Pub)
	if variableDecl.IsPub && variableDecl.IsConst {
		p.match(tokens.Const)
	}
	if variableDecl.Token.Type != tokens.Let {
		typeCheckStart := p.current
		if typeExpr, ok := p.checkType(); ok {
			variableDecl.Type = typeExpr
			if !p.check(tokens.Identifier) {
				p.disadvance(p.current - typeCheckStart)
				variableDecl.Type = nil
			}
		}
	}

	idents, exprs, ok := p.getIdentExprPairs("in variable declaration")
	if !ok {
		p.disadvance(p.current - currentStart)
		return nil
	}

	variableDecl.Identifiers = idents
	variableDecl.Expressions = exprs

	return &variableDecl
}

func (p *Parser) functionDeclaration() ast.Node {
	functionDecl := ast.FunctionDecl{}

	name, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "in function declaration"), tokens.Identifier)
	if !ok {
		return ast.NewImproper(p.peek(), ast.FunctionDeclaration)
	}
	functionDecl.Name = name
	functionDecl.Generics = p.genericParameters()
	functionDecl.Params = p.parameters(tokens.LeftParen, tokens.RightParen)
	functionDecl.ReturnTypes = p.returnings()

	p.Context.FunctionReturns.Push("functionDeclarationStmt", len(functionDecl.ReturnTypes))

	body, ok := p.getBody()
	if !ok {
		return ast.NewImproper(functionDecl.Name, ast.NA)
	}
	functionDecl.Body = body

	p.Context.FunctionReturns.Pop("functionDeclarationStmt")

	return &functionDecl
}

func (p *Parser) aliasDeclaration() ast.Node {
	return nil
}

func (p *Parser) enumDeclaration() ast.Node {
	return nil
}

func (p *Parser) classDeclaration() ast.Node {
	return nil
}

func (p *Parser) entityDeclaration() ast.Node {
	return nil
}

func (p *Parser) entityFunctionDeclaration(token tokens.Token, functionType ast.EntityFunctionType) ast.Node {
	return nil
}

func (p *Parser) constructorDeclaration() ast.Node {
	return nil
}

func (p *Parser) fieldDeclaration() ast.Node {
	return nil
}

func (p *Parser) methodDeclaration() ast.Node {
	return nil
}
