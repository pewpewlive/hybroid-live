package pass3

import (
	"hybroid/ast"
	wkr "hybroid/walker"
)

func Action(w *wkr.Walker) {

}

func WalkNode(w *wkr.Walker, node *ast.Node, scope *wkr.Scope) {
	switch newNode := (*node).(type) {
	case *ast.VariableDeclarationStmt:
		VariableDeclarationStmt(w, newNode, scope)
	case *ast.IfStmt:
		IfStmt(w, newNode, scope)
	case *ast.AssignmentStmt:
		AssignmentStmt(w, newNode, scope)
	case *ast.FunctionDeclarationStmt:
		FunctionDeclarationStmt(w, newNode, scope, wkr.Function)
	case *ast.ReturnStmt:
		ReturnStmt(w, newNode, scope)
	case *ast.YieldStmt:
		YieldStmt(w, newNode, scope)
	case *ast.BreakStmt:
		BreakStmt(w, newNode, scope)
	case *ast.ContinueStmt:
		ContinueStmt(w, newNode, scope)
	case *ast.RepeatStmt:
		RepeatStmt(w, newNode, scope)
	case *ast.WhileStmt:
		WhileStmt(w, newNode, scope)
	case *ast.ForStmt:
		ForloopStmt(w, newNode, scope)
	case *ast.TickStmt:
		TickStmt(w, newNode, scope)
	case *ast.CallExpr:
		CallExpr(w, newNode, scope, wkr.Function)
	case *ast.MethodCallExpr:
		MethodCallExpr(w, node, scope)
	case *ast.DirectiveExpr:
		DirectiveExpr(w, newNode, scope)
	case *ast.UseStmt:
		UseStmt(w, newNode, scope)
	case *ast.StructDeclarationStmt:
		StructDeclarationStmt(w, newNode, scope)
	case *ast.MatchStmt:
		MatchStmt(w, newNode, false, scope)
	case *ast.Improper:
		w.Error(newNode.GetToken(), "Improper statement: parser fault")
	default:
		w.Error(newNode.GetToken(), "Expected statement")
	}
}

func GetNodeValue(w *wkr.Walker, node *ast.Node, scope *wkr.Scope) wkr.Value {
	var val wkr.Value
	val = &wkr.Invalid{}

	switch newNode := (*node).(type) {
	case *ast.CallExpr:
		val = CallExpr(w, newNode, scope, wkr.Function)
	case *ast.DirectiveExpr:
		val = DirectiveExpr(w, newNode, scope)
	case *ast.AnonFnExpr:
		val = AnonFnExpr(w, newNode, scope)
	case *ast.AnonStructExpr:
		val = AnonStructExpr(w, newNode, scope)
	case *ast.MethodCallExpr:
		val = MethodCallExpr(w, node, scope)
	case *ast.NewExpr:
		val = NewExpr(w, newNode, scope)
	case *ast.MatchExpr:
		val = MatchExpr(w, newNode, scope)
	default:
		w.Error(newNode.GetToken(), "Expected expression")
	}

	return val
}

func WalkBody(w *wkr.Walker, body *[]ast.Node, tag wkr.ExitableTag, scope *wkr.Scope) {
	endIndex := -1
	for i := range *body {
		if tag.GetIfExits(wkr.All) {
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
