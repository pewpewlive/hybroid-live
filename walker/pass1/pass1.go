package pass1

import (
	"fmt"
	"hybroid/ast"
	wkr "hybroid/walker"
	"strings"
)

func Action(w *wkr.Walker, nodes []ast.Node, wlkrs map[string]*wkr.Walker) {
	w.Walkers = wlkrs
	w.Nodes = nodes

	scope := &w.Environment.Scope
	for i := range nodes {
		WalkNode(w, &nodes[i], scope)
	}
}

func WalkNode(w *wkr.Walker, node *ast.Node, scope *wkr.Scope) {
	switch newNode := (*node).(type) {
	case *ast.EnvironmentStmt:
		if w.Environment.Name != "" {
			w.Error(newNode.GetToken(), "cannot have 2 environment declaration statements in one file")
		}

		w.Environment.Name = strings.Join(newNode.Env.SubPaths, "::")
		for k, v := range w.Walkers {
			if k == w.Environment.Name {
				w.Error(newNode.GetToken(), fmt.Sprintf("duplicate names found between %s and %s", w.Environment.Path, v.Environment.Path))
			}
		}

		w.Walkers[w.Environment.Name] = w
	case *ast.VariableDeclarationStmt:
		VariableDeclarationStmt(w, newNode, scope)
	case *ast.IfStmt:
		IfStmt(w, newNode, scope)
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
		callerVal := GetNodeValue(w, &newNode.Caller, scope)
		CallExpr(w, callerVal, newNode, scope)
	case *ast.StructDeclarationStmt:
		StructDeclarationStmt(w, newNode, scope)
	case *ast.MatchStmt:
		MatchStmt(w, newNode, scope)
	case *ast.EnumDeclarationStmt:
		EnumDeclarationStmt(w, newNode, scope)
	case *ast.MacroDeclarationStmt:
		MacroDeclarationStmt(w, newNode, scope)
	case *ast.PewpewExpr:
		PewpewExpr(w, newNode, scope)
	// case *ast.TypeDeclarationStmt:
	// 	TypeDeclarationStmt(w, newNode, scope)
	case *ast.UseStmt:
	case *ast.AssignmentStmt:
	case *ast.EnvAccessExpr:
	case *ast.EntityDeclarationStmt:
		EntityDeclarationStmt(w, newNode, scope)
	default:
		w.Error(newNode.GetToken(), "Expected statement")
	}
}

func GetNodeValue(w *wkr.Walker, node *ast.Node, scope *wkr.Scope) wkr.Value {
	var val wkr.Value

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
		callerVal := GetNodeValue(w, &newNode.Caller, scope)
		val = CallExpr(w, callerVal, newNode, scope)
	case *ast.MapExpr:
		val = MapExpr(w, newNode, scope)
	case *ast.AnonFnExpr:
		val = AnonFnExpr(w, newNode, scope)
	case *ast.AnonStructExpr:
		val = AnonStructExpr(w, newNode, scope)
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
	case *ast.SpawnExpr:
		val = SpawnExpr(w, newNode, scope)
	// case *ast.CastExpr:
	// 	val = CastExpr(w, newNode, scope)
	case *ast.PewpewExpr:
		val = PewpewExpr(w, newNode, scope)
	case *ast.FmathExpr:
		val = FmathExpr(w, newNode, scope)
	case *ast.StandardExpr:
		val = StandardExpr(w, newNode, scope)
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
