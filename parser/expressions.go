package parser

import (
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/tokens"
)

func (p *Parser) expression() ast.Node {
	return &ast.LiteralExpr{Token: p.advance()}
}

func (p *Parser) fn() ast.Node {
	return nil
}

func (p *Parser) multiComparison() ast.Node {
	return nil
}

func (p *Parser) comparison() ast.Node {
	return nil
}

func (p *Parser) determineValueType(left ast.Node, right ast.Node) ast.PrimitiveValueType {
	return ast.String
}

func (p *Parser) term() ast.Node {
	return nil
}

func (p *Parser) factor() ast.Node {
	return nil
}

func (p *Parser) concat() ast.Node {
	return nil
}

func (p *Parser) unary() ast.Node {
	return nil
}

func (p *Parser) entity() ast.Node {
	return nil
}

func (p *Parser) call(caller ast.Node) ast.Node {
	return nil
}

func (p *Parser) accessorExprDepth2(ident *ast.Node) ast.Node {
	return nil
}

func (p *Parser) accessorExpr(ident *ast.Node) (ast.Node, *ast.IdentifierExpr) {
	return nil, nil
}

func (p *Parser) matchExpr() ast.Node {
	return nil
}

func (p *Parser) macroCall() ast.Node {
	return nil
}

func (p *Parser) new() ast.Node {
	return nil
}

func (p *Parser) spawn() ast.Node {
	return nil
}

func (p *Parser) self() ast.Node {
	return nil
}

func (p *Parser) primary(allowStruct bool) ast.Node {
	return nil
}

func (p *Parser) list() ast.Node {
	return nil
}

func (p *Parser) parseMap() ast.Node {
	return nil
}

func (p *Parser) structExpr() ast.Node {
	return nil
}

func (p *Parser) wrappedTypeExpr() *ast.TypeExpr {
	return nil
}

func (p *Parser) typeExpr() *ast.TypeExpr {
	name, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "for a type"), tokens.Identifier)
	if !ok {
		return nil
	}

	return &ast.TypeExpr{Name: &ast.IdentifierExpr{Name: name}}
}

func (p *Parser) envTypeExpr() *ast.EnvTypeExpr {
	envTypeExpr := ast.EnvTypeExpr{
		Type: ast.InvalidEnv,
	}
	name, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "for an environment type"), tokens.Identifier)
	envTypeExpr.Token = name
	if !ok {
		return &envTypeExpr
	}

	switch name.Lexeme {
	case "Mesh":
		envTypeExpr.Type = ast.MeshEnv
	case "Level":
		envTypeExpr.Type = ast.LevelEnv
	case "Sound":
		envTypeExpr.Type = ast.SoundEnv
	default:
		p.Alert(&alerts.InvalidEnvironmentType{}, alerts.NewSingle(name))
	}

	return &envTypeExpr
}

func (p *Parser) envPathExpr() ast.Node {
	envPath := &ast.EnvPathExpr{}

	ident, ok := p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "for an environment path"), tokens.Identifier)
	if !ok {
		return ast.NewImproper(ident, ast.EnvironmentPathExpression)
	}
	envPath.Path = ident

	for p.match(tokens.Colon) {
		ident, ok = p.consume(p.NewAlert(&alerts.ExpectedIdentifier{}, alerts.NewSingle(p.peek()), "in environment path"), tokens.Identifier)
		if !ok {
			return ast.NewImproper(ident, ast.EnvironmentPathExpression)
		}
		envPath.Combine(ident)
	}

	return envPath
}
