package pass3

import "hybroid/ast"

func Action(w *wkr.Walker, nodes *[]ast.Node, wlkrs *map[string]*walker.Walker) []ast.Node {
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