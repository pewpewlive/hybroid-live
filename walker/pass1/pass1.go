package pass1 // PREPARATION FOR THE ACTUAL WALKING

import (
	"hybroid/ast"
	wkr "hybroid/walker"
)

func Action(w *wkr.Walker, nodes []ast.Node, wkrs map[string]*wkr.Walker) {
	w.Walkers = wkrs
	w.Nodes = nodes

	scope := &w.Environment.Scope
	for i := range nodes {
		WalkNode(w, &nodes[i], scope)
	}
}

func WalkNode(w *wkr.Walker, node *ast.Node, scope *wkr.Scope) {
	// switch newNode := (*node).(type) {
	// 	// case *ast.EnvironmentDecl:
	// 	// 		switch name.Lexeme {
	// 	// case "Mesh":
	// 	// 	envTypeExpr.Type = ast.MeshEnv
	// 	// case "Level":
	// 	// 	envTypeExpr.Type = ast.LevelEnv
	// 	// case "Sound":
	// 	// 	envTypeExpr.Type = ast.SoundEnv
	// 	// default:
	// 	// 	p.Alert(&alerts.InvalidEnvironmentType{}, alerts.NewSingle(name))
	// 	// }

	// 	if w.Environment.Name != "" {
	// 		// w.Error(newNode.GetToken(), "cannot have 2 environment declaration statements in one file")
	// 	}

	// 	w.Environment.Name = newNode.Env.Path.Lexeme
	// 	// for k, v := range w.Walkers {
	// 	// 	if k == w.Environment.Name {
	// 	// 		 w.Error(newNode.GetToken(), fmt.Sprintf("duplicate names found between %s and %s", w.Environment.Path, v.Environment.Path))
	// 	// 	}
	// 	// }

	// 	w.Walkers[w.Environment.Name] = w
	// default:
	// }
}

func WalkBody(w *wkr.Walker, body *[]ast.Node, tag wkr.ExitableTag, scope *wkr.Scope) {
	endIndex := -1
	for i := range *body {
		if tag.GetIfExits(wkr.All) {
			// w.Warn((*body)[i].GetToken(), "unreachable code detected")
			endIndex = i
			break
		}
		WalkNode(w, &(*body)[i], scope)
	}
	if endIndex != -1 {
		*body = (*body)[:endIndex]
	}
}
