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

	if p.match(tokens.Env) {
		if p.environmentDeclaration().GetType() == ast.EnvironmentDeclaration {
			p.Alert(&alerts.EnvironmentRedaclaration{}, alerts.NewSingle(p.peek()))
		}
	}

	if p.match(tokens.Pub) {
		p.Context.IsPub = true
	}

	if p.match(tokens.Fn) {
		return p.functionDeclaration()
	}

	if varDecl := p.variableDeclaration(); varDecl != nil {
		return varDecl
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
		IsPub:   p.Context.IsPub,
		IsConst: false,
	}
	p.Context.IsPub = false

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

	idents, exprs, ok := p.getIdentExprPairs("in variable declaration", variableDecl.Type != nil)
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
	p.Context.IsPub = false

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
