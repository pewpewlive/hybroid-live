package pass2

import (
	"hybroid/ast"
	"hybroid/walker"
	wkr "hybroid/walker"
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
	case *ast.VariableDeclarationStmt:
		VariableDeclarationStmt(w, newNode, scope)
	case *ast.FunctionDeclarationStmt:
		FunctionDeclarationStmt(w, newNode, scope, wkr.Function)
	case *ast.StructDeclarationStmt:
		StructDeclarationStmt(w, newNode, scope)
	case *ast.EnumDeclarationStmt:
		EnumDeclarationStmt(w, newNode, scope)
	case *ast.ReturnStmt:
		ReturnStmt(w, newNode, scope)
	case *ast.YieldStmt:
		YieldStmt(w, newNode, scope)
	case *ast.Improper:
		//w.Error(newNode.GetToken(), "Improper statement: parser fault")
	default:
		w.Error(newNode.GetToken(), "Expected statement")
	}
}

func GetNodeValue(w *walker.Walker, node *ast.Node, scope *walker.Scope) walker.Value {
	var val walker.Value

	switch newNode := (*node).(type) {
	case *ast.LiteralExpr:
		val = LiteralExpr(w, newNode)
	case *ast.BinaryExpr:
		val = BinaryExpr(w, newNode, scope)
	case *ast.IdentifierExpr:
		val = IdentifierExpr(w, node, scope)
	case *ast.GroupExpr:
		val = GroupingExpr(w, newNode, scope)
	case *ast.ListExpr:
		val = ListExpr(w, newNode, scope)
	case *ast.UnaryExpr:
		val = UnaryExpr(w, newNode, scope)
	case *ast.CallExpr:
		val = CallExpr(w, newNode, scope, wkr.Function)
	case *ast.MapExpr:
		val = MapExpr(w, newNode, scope)
	case *ast.AnonFnExpr:
		val = AnonFnExpr(w, newNode, scope)
	case *ast.AnonStructExpr:
		val = AnonStructExpr(w, newNode, scope)
	case *ast.MethodCallExpr:
		val = MethodCallExpr(w, node, scope)
	case *ast.MemberExpr:
		val = MemberExpr(w, newNode, scope)
	case *ast.FieldExpr:
		val = FieldExpr(w, newNode, scope)
	case *ast.NewExpr:
		val = NewExpr(w, newNode, scope)
	case *ast.SelfExpr:
		val = SelfExpr(w, newNode, scope)
	case *ast.MatchExpr:
		val = MatchExpr(w, newNode, scope)
	case *ast.EnvAccessExpr:
		val = EnvAccessExpr(w, newNode, scope)
	default:
		w.Error(newNode.GetToken(), "Expected expression")
		return &walker.Invalid{}
	}
	return val
}

func WalkBody(w *walker.Walker, body *[]ast.Node, scope *walker.Scope) {
	for i := range *body {
		WalkNode(w, &(*body)[i], scope)
	}
}

func TypeifyNodeList(w *wkr.Walker, nodes *[]ast.Node, scope *wkr.Scope) []wkr.Type {
	arguments := make([]wkr.Type, 0)
	for i := range *nodes {
		val := GetNodeValue(w, &(*nodes)[i], scope)
		if function, ok := val.(*wkr.FunctionVal); ok {
			arguments = append(arguments, function.Returns...)
		} else {
			arguments = append(arguments, val.GetType())
		}
	}
	return arguments
}
