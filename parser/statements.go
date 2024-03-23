package parser

import (
	"hybroid/ast"
	"hybroid/lexer"
)

func (p *Parser) statement() ast.Node {
	token := p.peek().Type
	next := p.peek(1).Type

	if token == lexer.Pub && next == lexer.Fn {
		p.advance()
		token = p.peek().Type
	}

	switch token {
	case lexer.Let, lexer.Pub, lexer.Const:
		p.advance()
		return p.variableDeclarationStmt()
	case lexer.Add:
		p.advance()
		return p.addToStmt()
	case lexer.Remove:
		p.advance()
		return p.removeFromStmt()
	case lexer.Fn:
		p.advance()
		return p.functionDeclarationStmt()
	case lexer.Return:
		p.advance()
		return p.returnStmt()
	case lexer.Identifier:
		return p.assignmentStmt()
	case lexer.If:
		p.advance()
		return p.ifStmt()
	case lexer.Repeat:
		p.advance()
		return p.repeatStmt()
	case lexer.Tick:
		p.advance()
		return p.tickStmt()
	}
	expr := p.expression()
	if expr.GetType() == 0 {
		p.error(p.peek(), "expected expression")
	}
	return expr
}

func (p *Parser) ifStmt() ast.Node { //
	ifStm := ast.IfStmt{}
	expr := p.expression()

	body := make([]ast.Node, 0)
	if _, success := p.consume("expected body of the function", lexer.LeftBrace); success {
		for !p.match(lexer.RightBrace) {
			if p.peek().Type == lexer.Eof {
				p.error(p.peek(), "expected body closure")
				break
			}
			statement := p.statement()
			if statement != nil {
				body = append(body, statement)
			}
		}
	}

	ifStm.Body = body
	ifStm.BoolExpr = expr

	return ifStm
}

func (p *Parser) assignmentStmt() ast.Node {
	expr := p.expression()

	idents := []ast.Node{expr}

	for p.match(lexer.Comma) { // memberExpr or IdentifierExpr
		identExpr := p.expression()
		idents = append(idents, identExpr)
	}

	if p.match(lexer.Equal) {
		values := []ast.Node{p.expression()}
		for p.match(lexer.Comma) {
			expr2 := p.expression()

			values = append(values, expr2)
		}
		expr = ast.AssignmentStmt{Identifiers: idents, Values: values, Token: p.peek(-1)}
	} else if p.match(lexer.PlusEqual, lexer.MinusEqual, lexer.SlashEqual, lexer.StarEqual, lexer.CaretEqual, lexer.ModuloEqual) {
		assignOp := p.peek(-1)
		op := p.getOp(assignOp)
		if len(idents) > 1 {
			p.error(assignOp, "cannot assign to multiple variables with this operator")
		}
		expr2 := p.term()
		binExpr := p.createBinExpr(expr, op, op.Type, op.Lexeme, ast.GroupExpr{Expr: expr2, ValueType: expr2.GetValueType(), Token: expr2.GetToken()})
		expr = ast.AssignmentStmt{Identifiers: idents, Values: []ast.Node{binExpr}, Token: assignOp}
	}

	return expr
}

func (p *Parser) returnStmt() ast.Node {
	returnStmt := ast.ReturnStmt{
		Token: p.peek(-1),
	}

	if p.peek().Type == lexer.RightBrace {
		return returnStmt
	}
	args := []ast.Node{}
	expr := p.expression()
	args = append(args, expr)
	for p.match(lexer.Comma) {
		expr = p.expression()
		if expr.GetType() == 0 {
			p.error(p.peek(), "expected expression")
		}
		args = append(args, expr)
	}

	returnStmt.Args = args

	return returnStmt
}

func (p *Parser) functionDeclarationStmt() ast.Node {
	fnDec := ast.FunctionDeclarationStmt{}

	fnDec.IsLocal = p.peek(-2).Type != lexer.Pub

	ident, ok := p.consume("expected a function name", lexer.Identifier)
	if !ok {
		return fnDec
	}

	fnDec.Name = ident

	args := p.arguments()
	var params []lexer.Token

	for _, arg := range args {
		if arg.GetType() == ast.Identifier {
			params = append(params, arg.GetToken())
			continue
		}
		p.error(arg.GetToken(), "expected identifier in function declaration")
	}

	fnDec.Params = params

	body := make([]ast.Node, 0)
	if token, success := p.consume("expected body of the function", lexer.LeftBrace); success { // hjere
		for !p.match(lexer.RightBrace) {
			if p.peek().Type == lexer.Eof {
				p.error(token, "expected body closure")
				break
			}
			statement := p.statement()
			if statement != nil {
				body = append(body, statement)
			}
		}
	}

	fnDec.Body = body

	return fnDec
}

func (p *Parser) addToStmt() ast.Node {
	add := ast.AddStmt{
		Token: p.peek(-1),
	}

	add.Value = p.expression()
	if add.GetType() == 0 {
		p.error(p.peek(), "expected expression")
	}

	if _, ok := p.consume("expected keyword 'to' after expression in an 'add' statement", lexer.To); !ok {
		return add
	}

	if ident, ok := p.consume("expected identifier after keyword 'to'", lexer.Identifier); ok {
		add.Identifier = ident.Lexeme
	}

	return add
}

func (p *Parser) removeFromStmt() ast.Node {
	remove := ast.RemoveStmt{
		Token: p.peek(-1),
	}

	remove.Value = p.expression()
	if remove.GetType() == 0 {
		p.error(p.peek(), "expected expression")
	}

	if _, ok := p.consume("expected keyword 'from' after expression in a 'remove' statement", lexer.From); !ok {
		return remove
	}

	if ident, ok := p.consume("expected identifier after keyword 'from'", lexer.Identifier); ok {
		remove.Identifier = ident.Lexeme
	}

	return remove
}

func (p *Parser) repeatStmt() ast.Node {
	repeatStmt := ast.RepeatStmt{
		Token: p.peek(-1),
	}

	repeatStmt.Skip = ast.LiteralExpr{Value: "1", ValueType: ast.Number, Token: repeatStmt.Token}
	repeatStmt.Start = ast.LiteralExpr{Value: "1", ValueType: ast.Number, Token: repeatStmt.Token}

	for i := 0; i < 4; i++ {
		if p.match(lexer.With) {
			identExpr := p.expression()
			if identExpr.GetType() != ast.Identifier {
				p.error(identExpr.GetToken(), "expected identifier expression after keyword 'with'")
				return repeatStmt
			}
			repeatStmt.Variable = identExpr.(ast.IdentifierExpr)
		} else if p.match(lexer.To) {
			repeatStmt.Iterator = p.expression()
		} else if p.match(lexer.By) {
			repeatStmt.Skip = p.expression()
		} else if p.match(lexer.From) {
			repeatStmt.Start = p.expression()
		}
	}

	body := make([]ast.Node, 0)
	if _, success := p.consume("expected body of the repeat statement", lexer.LeftBrace); success {
		for !p.match(lexer.RightBrace) {
			if p.peek().Type == lexer.Eof {
				p.error(p.peek(), "expected body closure")
				break
			}
			statement := p.statement()
			if statement != nil {
				body = append(body, statement)
			}
		}
	}

	repeatStmt.Body = body

	return repeatStmt
}

func (p *Parser) tickStmt() ast.Node {
	tickStmt := ast.TickStmt{
		Token: p.peek(-1),
	}

	if p.match(lexer.With) {
		identExpr := p.expression()
		if identExpr.GetType() != ast.Identifier {
			p.error(identExpr.GetToken(), "expected identifier expression after keyword 'with'")
			return tickStmt
		}
		tickStmt.Variable = identExpr.(ast.IdentifierExpr)
	}

	body := make([]ast.Node, 0)
	if _, success := p.consume("expected body of the tick statement", lexer.LeftBrace); success {
		for !p.match(lexer.RightBrace) {
			if p.peek().Type == lexer.Eof {
				p.error(p.peek(), "expected body closure")
				break
			}
			statement := p.statement()
			if statement != nil {
				body = append(body, statement)
			}
		}
	}

	tickStmt.Body = body

	return tickStmt
}

func (p *Parser) variableDeclarationStmt() ast.Node {
	variable := ast.VariableDeclarationStmt{
		Token: p.peek(-1), //let or pub, important
	}

	ident, _ := p.consume("expected identifier in variable declaration", lexer.Identifier)
	idents := []string{ident.Lexeme}
	for p.match(lexer.Comma) {
		ident, identOk := p.consume("expected identifier in variable declaration", lexer.Identifier)
		if !identOk {
			return ast.Unknown{Token: p.peek(-1)}
		}

		idents = append(idents, ident.Lexeme)
	}

	variable.Identifiers = idents

	if !p.match(lexer.Equal) {
		variable.Values = []ast.Node{}
		return variable
	}

	expr := p.expression()
	if expr.GetType() == 0 {
		p.error(p.peek(), "expected expression")
	}

	exprs := []ast.Node{expr}
	for p.match(lexer.Comma) {
		expr = p.expression()
		if expr.GetType() == 0 {
			p.error(p.peek(), "expected expression")
		}
		exprs = append(exprs, expr)
	}
	variable.Values = exprs

	return variable
}
