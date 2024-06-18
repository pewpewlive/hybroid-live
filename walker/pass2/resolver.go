package pass2

import (
	"hybroid/ast"
	wkr "hybroid/walker"
)

func Action(w *wkr.Walker, )

func WalkNode(w *wkr.Walker, node *ast.Node, scope *wkr.Scope) {
	switch newNode := (*node).(type) {
	case *ast.VariableDeclarationStmt:
		VariableDeclaration(w, newNode, scope)
	case *ast.IfStmt:
		IfStmt(w, newNode, scope)
	case *ast.AssignmentStmt:
		Assignment(w, newNode, scope)
	case *ast.FunctionDeclarationStmt:
		FunctionDeclaration(w, newNode, scope, wkr.Function)
	case *ast.ReturnStmt:
		ReturnStmt(w, newNode, scope)
	case *ast.YieldStmt:
		YieldStmt(w, newNode, scope)
	case *ast.BreakStmt:
		BreakStmt(w, newNode, scope)
	case *ast.ContinueStmt:
		ContinueStmt(w, newNode, scope)
	case *ast.RepeatStmt:
		Repeat(w, newNode, scope)
	case *ast.WhileStmt:
		While(w, newNode, scope)
	case *ast.ForStmt:
		Forloop(w, newNode, scope)
	case *ast.TickStmt:
		Tick(w, newNode, scope)
	case *ast.CallExpr:
		w.CallExpr(newNode, scope, wkr.Function)
	case *ast.MethodCallExpr:
		w.MethodCallExpr(node, scope)
	case *ast.DirectiveExpr:
		w.DirectiveExpr(newNode, scope)
	case *ast.UseStmt:
		Use(w,newNode, scope)
	case *ast.StructDeclarationStmt:
		StructDeclaration(w, newNode, scope)
	case *ast.MatchStmt:
		Match(w, newNode, false, scope)
	case *ast.Improper:
		w.Error(newNode.GetToken(), "Improper statement: parser fault")
	default:
		w.Error(newNode.GetToken(), "Expected statement")
	}
}

func GetNodeValue(w *wkr.Walker, node *ast.Node, scope *wkr.Scope) wkr.Value {
	var val wkr.Value

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
		val = w.CallExpr(newNode, scope, wkr.Function)
	case *ast.MapExpr:
		val = w.MapExpr(newNode, scope)
	case *ast.DirectiveExpr:
		val = w.DirectiveExpr(newNode, scope)
	case *ast.AnonFnExpr:
		val = w.AnonFnExpr(newNode, scope)
	case *ast.AnonStructExpr:
		val = AnonStructExpr(w, newNode, scope)
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
	default:
		w.Error(newNode.GetToken(), "Expected expression")
		return &wkr.Invalid{}
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
