package walker

import (
	"fmt"
	"hybroid/ast"
	"strings"
)

func Action(w *Walker, node *ast.Node, wkrs *map[string]*Walker) {
	switch newNode := (*node).(type) {
	case *ast.EnvironmentStmt:
		if w.Environment.Name != "" {
			w.Error(newNode.GetToken(), "cannot have 2 environment declaration statements in one file")
			return
		}

		w.Environment.Name = strings.Join(newNode.Env.SubPaths, "::")

		for k, v := range *wkrs {
			if k == w.Environment.Name {
				w.Error(newNode.GetToken(), fmt.Sprintf("duplicate names found between %s and %s", w.Environment.Path, v.Environment.Path))
			}
		}
	default:
		w.Error(newNode.GetToken(), "first statement must be an environment declaration")
	}
}