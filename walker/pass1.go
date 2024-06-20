package walker

import "hybroid/ast"

func Action(w *Walker, node *ast.Node, wkrs *map[string]*Walker) {
	switch newNode := (*node).(type) {
	case *ast.EnvironmentStmt:
		envsLength := len(newNode.Env.SubEnvs)

		var parent *ast.Environment
		for i := newNode.Env.SubEnvs {
			subEnvName := newNode.Env.SubEnvs[i].GetToken().Lexeme
			if i == envsLength-1 {

				continue
			}

			// env Env2::Env3::Env4 as Shared

			// env Env2::Env3

			// env BossMeshes::Outside as Mesh

			// env BossMeshes as Mesh

			// let a = VEnv2_variable.name

			// let Env3 = Env4

			// let b = Env2::Env3::Env4::a

			if parent == nil {
				wkrs[subEnvName] = NewWalker(w.Environment.Type.Path)
				parent = wkrs[subEnvName].Environment
			}else {

			}
		}

	default:
		w.Error(newNode.GetToken(), "first statement must be an environment declaration")
	}
}
