package pass1

import (
	"hybroid/ast"
	"hybroid/walker"
)

func Action(w *walker.Walker, nodes *[]ast.Node, wlkrs *map[string]*walker.Walker) []ast.Node {
	w.Walkers = wlkrs
	w.Nodes = nodes

	newNodes := make([]ast.Node, 0)

	scope := &w.Environment.Scope
	for _, node := range *nodes {
		WalkNode(w, &node, scope)
	
		newNodes = append(newNodes, node)
	}

	return newNodes
}


func WalkNode(w *walker.Walker, node *ast.Node, scope *walker.Scope) {
	switch newNode := (*node).(type) {
	case *ast.EnvironmentStmt:
		EnvStmt(w, newNode, scope)
	case *ast.VariableDeclarationStmt:
		VariableDeclaration(w, newNode, scope)
	case *ast.FunctionDeclarationStmt:
		FunctionDeclaration(w, newNode, scope, walker.Function)
	case *ast.StructDeclarationStmt:
		StructDeclaration(w, newNode, scope)
	case *ast.EnumDeclarationStmt:
		EnumDeclaration(w, newNode, scope)
	case *ast.Improper:
		w.Error(newNode.GetToken(), "Improper statement: parser fault")
	default:
		w.Error(newNode.GetToken(), "Expected statement")
	}
}

func GetNodeValue(w *walker.Walker, node *ast.Node, scope *walker.Scope) walker.Value {
	var val walker.Value

	switch newNode := (*node).(type) {
	case *ast.LiteralExpr:
		val = w.LiteralExpr(newNode)
	case *ast.BinaryExpr:
		val = w.BinaryExpr(newNode, scope)
	case *ast.IdentifierExpr:
		val = w.IdentifierExpr(node, scope)
	case *ast.GroupExpr:
		val = w.GroupingExpr(newNode, scope)
	case *ast.ListExpr:
		val = w.ListExpr(newNode, scope)
	case *ast.UnaryExpr:
		val = w.UnaryExpr(newNode, scope)
	case *ast.CallExpr:
		val = w.CallExpr(newNode, scope, walker.Function)
	case *ast.MapExpr:
		val = w.MapExpr(newNode, scope)
	case *ast.DirectiveExpr:
		val = w.DirectiveExpr(newNode, scope)
	case *ast.AnonFnExpr:
		val = AnonFnExprPass1(w, newNode, scope)
	case *ast.AnonStructExpr:
		val = AnonStructExprPass1(w,newNode, scope)
	case *ast.MethodCallExpr:
		val = w.MethodCallExpr(node, scope)
	case *ast.MemberExpr:
		val = w.MemberExpr(newNode, scope)
	case *ast.FieldExpr:
		val = w.FieldExpr(newNode, scope)
	case *ast.NewExpr:
		val = w.NewExpr(newNode, scope)
	case *ast.SelfExpr:
		val = w.SelfExpr(newNode, scope)
	case *ast.MatchExpr:
		val = MatchExpr(w, newNode, scope)
	case *ast.EnvExpr:
		val = EnvExpr(w, newNode, scope)
	default:
		w.Error(newNode.GetToken(), "Expected expression")
		return &walker.Invalid{}
	}
	return val
}

func WalkBody(w *walker.Walker, body *[]ast.Node, tag walker.ExitableTag, scope *walker.Scope) {
	endIndex := -1
	for i := range *body {
		if tag.GetIfExits(walker.All) {
			w.Warn((*body)[i].GetToken(), "unreachable code detected")
			endIndex = i
			break
		}
		WalkNode(w, &(*body)[i], scope)
	}
	if endIndex != -1 {
		*body = (*body)[:endIndex]
	}
}