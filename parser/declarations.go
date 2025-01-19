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
	/*


	 */

	// switch expression {
	// case condition:

	// }asdsad

	return nil
}

func (p *Parser) environmentDeclaration() ast.Node {
	stmt := &ast.EnvironmentDecl{}

	expr := p.EnvPathExpr()
	if expr.GetType() == ast.NA {
		return expr
	}

	if _, ok := p.consume(p.NewAlert(&alerts.ExpectedKeyword{}, alerts.NewSingle(p.peek()), tokens.As), tokens.As); !ok {
		return ast.NewImproper(p.peek(), ast.EnvironmentDeclaration)
	}

	envTypeExpr := p.EnvType()
	if envTypeExpr.Type == ast.InvalidEnv {
		return stmt
	}

	envPathExpr, _ := expr.(*ast.EnvPathExpr)
	stmt.EnvType = envTypeExpr
	stmt.Env = envPathExpr
	p.Context.EnvDeclaration = stmt

	return stmt
}

func (p *Parser) variableDeclaration() ast.Node {
	variable := ast.VariableDecl{
		Token:   p.peek(),
		IsLocal: true,
		IsConst: false,
	}

	currentStart := p.current

	switch p.advance().Type {
	case tokens.Let:
		if _, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek())), tokens.Identifier); !ok {
			return ast.NewImproper(variable.Token, ast.VariableDeclaration)
		}

		names, ok := p.getIdentifiers()
		if !ok {
			return ast.NewImproper(variable.Token, ast.VariableDeclaration)
		}

		values, ok := p.getExpressions()
		if !ok {
			return ast.NewImproper(variable.Token, ast.VariableDeclaration)
		}

		variable.Identifiers = names
		variable.Values = values

		break
	case tokens.Const:
		variable.IsConst = true

		var typ *ast.TypeExpr
		if p.CheckType() {
			typ = p.Type()
		}

		variable.Type = typ

		names, ok := p.getIdentifiers()
		if !ok {
			return ast.NewImproper(variable.Token, ast.VariableDeclaration)
		}

		values, ok := p.getExpressions()
		if !ok {
			return ast.NewImproper(variable.Token, ast.VariableDeclaration)
		}

		variable.Identifiers = names
		variable.Values = values

		break
	case tokens.Pub:
		variable.IsLocal = false
		if p.match(tokens.Const) {
			variable.IsConst = true
		}

		var typ *ast.TypeExpr
		if p.CheckType() {
			typ = p.Type()
		}

		variable.Type = typ // TODO: finish

		names, ok := p.getIdentifiers()
		if !ok {
			p.disadvance(p.current - currentStart)
			return ast.NewImproper(variable.Token, ast.VariableDeclaration)
		}

		values, ok := p.getExpressions()
		if !ok {
			p.disadvance(p.current - currentStart)
			return ast.NewImproper(variable.Token, ast.VariableDeclaration)
		}

		variable.Identifiers = names
		variable.Values = values
		break
	}

	return &variable
}

func (p *Parser) functionDeclaration(IsLocal bool) ast.Node {
	return nil
}

func (p *Parser) aliasDeclaration(isLocal bool) ast.Node {
	return nil
}

func (p *Parser) enumDeclaration(local bool) ast.Node {
	return nil
}

func (p *Parser) classDeclaration(isLocal bool) ast.Node {
	return nil
}

func (p *Parser) entityDeclaration(isLocal bool) ast.Node {
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

func (p *Parser) methodDeclaration(IsLocal bool) ast.Node {
	return nil
}
