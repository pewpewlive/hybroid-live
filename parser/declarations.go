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
	}

	p.AlertPanic(&alerts.ExpectedStatement{}, alerts.NewSingle(p.peek()))
	return ast.NewImproper(p.peek(), ast.NA)
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

	name, _ := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "in function declaration"), tokens.Identifier)

	functionDecl.Name = name
	functionDecl.Generics = p.genericParams()
	functionDecl.Params = p.functionParams(tokens.LeftParen, tokens.RightParen)
	functionDecl.ReturnTypes = p.functionReturns()

	p.Context.FunctionReturns.Push("functionDeclarationStmt", len(functionDecl.ReturnTypes))

	body, ok := p.body(false, true, true)
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

	name, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "in alias declaration"), tokens.Identifier)
	if ok {
		aliasDecl.Name = name
	}

	p.consume(p.NewAlert(&alerts.ExpectedSymbol{}, alerts.NewSingle(p.peek()), tokens.Equal, "after identifier in alias declaration"), tokens.Equal)

	aliasDecl.Type = p.typeExpr()

	return aliasDecl
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
