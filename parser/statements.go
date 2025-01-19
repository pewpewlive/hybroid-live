package parser

import (
	"hybroid/ast"
)

func (p *Parser) statement() (returnNode ast.Node) {
	return nil
}

func (p *Parser) expressionStatement() ast.Node {
	return nil
}

func (p *Parser) destroyStmt() ast.Node {
	return nil
}

func (p *Parser) ifStmt(else_exists bool, is_else bool, is_elseif bool) *ast.IfStmt {
	return nil
}

func (p *Parser) assignmentStmt(expr ast.Node) ast.Node {
	return nil
}

func (p *Parser) returnStmt() ast.Node {
	return nil
}

func (p *Parser) returnArgs() ([]ast.Node, bool) {
	return nil, false
}

func (p *Parser) yieldStmt() ast.Node {
	return nil
}

func (p *Parser) repeatStmt() ast.Node {
	return nil
}

func (p *Parser) whileStmt() ast.Node {
	return nil
}

func (p *Parser) forStmt() ast.Node {
	return nil
}

func (p *Parser) tickStmt() ast.Node {
	return nil
}

func (p *Parser) useStmt() ast.Node {
	return nil
}

func (p *Parser) matchStmt(isExpr bool) *ast.MatchStmt {
	return nil
}

func (p *Parser) caseStmt(isExpr bool) ([]ast.CaseStmt, bool) {
	return nil, false
}
