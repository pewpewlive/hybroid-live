package pass2 // THE ACTUAL WALKING

import (
	"hybroid/ast"
	"hybroid/lexer"
	"hybroid/walker"
	wkr "hybroid/walker"
)

func Action(w *walker.Walker, wlkrs map[string]*walker.Walker) {
	w.Walkers = wlkrs

	scope := &w.Environment.Scope
	for i := range w.Nodes {
		WalkNode(w, &w.Nodes[i], scope)
	}

	w.Walked = true;
}

func WalkNode(w *wkr.Walker, node *ast.Node, scope *wkr.Scope) {
	switch newNode := (*node).(type) {
	case *ast.EnvironmentStmt:
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
		val := GetNodeValue(w, &newNode.Caller, scope)
		CallExpr(w, val, newNode, scope)
	case *ast.EnvAccessExpr:
		_, newVersion := EnvAccessExpr(w, newNode)
		if newVersion != nil {
			*node = newVersion
		}
	case *ast.ClassDeclarationStmt:
		ClassDeclarationStmt(w, newNode, scope)
	case *ast.EnumDeclarationStmt:
		EnumDeclarationStmt(w, newNode, scope)
	case *ast.MatchStmt:
		MatchStmt(w, newNode, false, scope)
	case *ast.AssignmentStmt:
		AssignmentStmt(w, newNode, scope)
	case *ast.UseStmt:
		UseStmt(w, newNode, scope)
	case *ast.SpawnExpr:
		SpawnExpr(w, newNode, scope)
	case *ast.NewExpr:
		NewExpr(w, newNode, scope)
	// case *ast.TypeDeclarationStmt:
	// 	TypeDeclarationStmt(w, newNode, scope)
	case *ast.Improper:
		w.Error(newNode.GetToken(), "Improper statement: parser fault")
	case *ast.MacroDeclarationStmt:
	case *ast.EntityDeclarationStmt:
		EntityDeclarationStmt(w, newNode, scope)
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
		callVal := GetNodeValue(w, &newNode.Caller, scope)
		val = CallExpr(w, callVal, newNode, scope)
	case *ast.MapExpr:
		val = MapExpr(w, newNode, scope)
	case *ast.FunctionExpr:
		val = FunctionExpr(w, newNode, scope)
	case *ast.StructExpr:
		val = StructExpr(w, newNode, scope)
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
		var newVersion ast.Node
		val, newVersion = EnvAccessExpr(w, newNode)
		if newVersion != nil {
			*node = newVersion
		}
	case *ast.SpawnExpr:
		val = SpawnExpr(w, newNode, scope)
	// case *ast.CastExpr:
	// 	val = CastExpr(w, newNode, scope)
	case *ast.UseStmt:
	default:
		w.Error(newNode.GetToken(), "Expected expression")
		return &walker.Invalid{}
	}

	if scope.Node != nil {
		if w.Context.PewpewVarFound {
			scope.Node.Index = -1
			return val
		}
		// name := w.Context.Node.GetToken().Lexeme
		// if name == "meshes" && w.GetEnvStmt().EnvType.Type == ast.MeshEnv {
		// 	scope.Node.Index = -1
		// 	return val
		// }
		_, scope.Node.Index, _ = scope.Container.ContainsField((*node).GetToken().Lexeme)
		scope.Node.Index += 1
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

// func WalkBody(w *wkr.Walker, body *[]ast.Node, tag wkr.ExitableTag, scope *wkr.Scope) {
// 	endIndex := -1
// 	for i := range *body {
// 		if tag.GetIfExits(wkr.All) {
// 			w.Warn((*body)[i].GetToken(), "unreachable code detected")
// 			endIndex = i
// 			break
// 		}
// 		WalkNode(w, &(*body)[i], scope)
// 	}
// 	if endIndex != -1 {
// 		*body = (*body)[:endIndex]
// 	}
// }

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

func WalkParams(w *wkr.Walker, parameters []ast.Param, scope *wkr.Scope, declare func(name lexer.Token, value wkr.Value)) []wkr.Type {
	variadicParams := make(map[lexer.Token]int)
	params := make([]wkr.Type, 0)
	for i, param := range parameters {
		params = append(params, TypeExpr(w, param.Type, scope, false))
		if params[i].GetType() == wkr.Variadic {
			variadicParams[parameters[i].Name] = i
		}
		value := w.TypeToValue(params[i])
		declare(param.Name, value)
	}

	if len(variadicParams) > 1 {
		w.Error(parameters[0].Name, "can only have one vartiadic parameter")
	}else if len(variadicParams) != 0 {
		for k, v := range variadicParams {
			if v != len(parameters)-1 {
				w.Error(k, "variadic parameter should be last")
			}
		}
	}

	return params
}