package pass2

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
	wkr "hybroid/walker"
	"hybroid/walker/pass1"
)

func Action(w *wkr.Walker) {
	InspectScopes(w, &w.Environment.Scope)
}

func InspectScopes(w *wkr.Walker, s *wkr.Scope) {
	for i := range s.Children {
		InspectScopes(w, s.Children[i])
	}
}

func InspectVariables(w *wkr.Walker, s *wkr.Scope) {
	for i := range s.Variables {
		InspectVar(w, s.Variables[i], []*ast.EnvExpr{})
	}
}

func InspectVar(w *wkr.Walker, variable *wkr.VariableVal, recursionPotentials []*ast.EnvExpr) {
	if variable.Value.GetType().GetType() == wkr.Unresolved {
		unresolved := variable.Value.GetType().(*wkr.UnresolvedType)
		newValue = GetEnvFieldType(w, unresolved.Expr, nil, 0)
		if variable.Value.GetType().PVT() != ast.Un {
		}
		return
	}

	if variable.Value.GetType().PVT() == ast.Unresolved {
		if unresolvedVal, ok := variable.Value.(*wkr.UnresolvedVal); ok {
			for _, v := range recursionPotentials {
				if v == unresolvedVal.Expr {
					w.Error(v.GetToken(), "Failed resolve: found recursion (initial)")
					w.Error(unresolvedVal.Expr.GetToken(), "Failed resolve: found recursion (cause of recursion)")
					variable.Value = &wkr.Invalid{}
					return
				}
			}
			if !variable.IsLocal {
				recursionPotentials = append(recursionPotentials, &ast.EnvExpr{
					Envs: []ast.Node{
						&ast.IdentifierExpr{
							Name: lexer.Token{Lexeme: w.Environment.Type.Name, Location: variable.Token.Location},
						},
						&ast.IdentifierExpr{
							Name: variable.Token,
						},
					},
				})
			}
			variable.Value = ResolveEnvExpr(w, unresolvedVal.Expr, recursionPotentials)
		}
		return
	}

	if mapVal, is := variable.Value.(*wkr.MapVal); is {
		for i := range mapVal.Members {
			InspectVal(w, &mapVal.Members[i])
		}
	} else if listVal, is := variable.Value.(*wkr.ListVal); is {
		for i := range listVal.Values {
			InspectVal(w, &listVal.Values[i])
		}
	}
}

func InspectVal(w *wkr.Walker, val *wkr.Value) {
	if unresolvedVal, ok := (*val).(*wkr.UnresolvedVal); ok {
		*val = ResolveEnvExpr(w, unresolvedVal.Expr, []*ast.EnvExpr{})
	}
}

func ResolveEnvExpr(w *wkr.Walker, expr *ast.EnvExpr, recursionPotentials []*ast.EnvExpr) wkr.Value {
	return GetEnvFieldValue(w, expr, nil, recursionPotentials, 0)
}

func GetEnvFieldValue(w *wkr.Walker, envExpr *ast.EnvExpr, owner wkr.Value, recursionPotentials []*ast.EnvExpr, depth int) wkr.Value {
	if depth > len(envExpr.Envs)-1 {
		return owner
	}
	if owner.GetType().PVT() != ast.Environment {
		w.Error(envExpr.Envs[depth-1].GetToken(), fmt.Sprintf("Resolve failed: expected type Environment, got %s", owner.GetType().ToString()))
		return &wkr.Invalid{}
	}
	env := owner.(*wkr.EnvironmentVal)
	previousDepth := depth
	depth += 1
	switch node := envExpr.Envs[previousDepth].(type) {
	case *ast.IdentifierExpr:
		if owner == nil {
			value := pass1.IdentifierExpr(w, &envExpr.Envs[previousDepth], &w.Environment.Scope)
			variable := value.(*wkr.VariableVal)
			InspectVar(w, variable, recursionPotentials)
			return GetEnvFieldValue(w, envExpr, variable.Value, recursionPotentials, depth)
		} else {
			if child, found := env.Childern[node.Name.Lexeme]; found {
				return GetEnvFieldValue(child, envExpr, child.Environment, recursionPotentials, depth)
			}
			if variable, found := env.Variables[node.Name.Lexeme]; found {
				InspectVar(w, variable, recursionPotentials)
				return GetEnvFieldValue(w, envExpr, variable, recursionPotentials, depth)
			}
		}
	case *ast.MapExpr, *ast.ListExpr, *ast.SelfExpr, *ast.DirectiveExpr, *ast.GroupExpr, *ast.BinaryExpr, *ast.AnonFnExpr, *ast.AnonStructExpr, *ast.MatchExpr:
		w.Error(node.GetToken(), fmt.Sprintf("Resolve failed: Cannot have a %s in %s", node.GetType(), envExpr.GetType()))
	default:
		if owner == nil {
			value := pass1.GetNodeValue(w, &node, &w.Environment.Scope)
			InspectVal(w, &value)
			return GetEnvFieldValue(w, envExpr, value, recursionPotentials, depth)
		} else {
			value := pass1.GetNodeValue(w, &node, &env.Scope)
			InspectVal(w, &value)
			return GetEnvFieldValue(w, envExpr, value, recursionPotentials, depth)
		}
	}

	return &wkr.Invalid{}
}

func GetEnvFieldType(w *wkr.Walker, envExpr *ast.EnvExpr, owner wkr.Value, depth int) wkr.Value {
	if depth > len(envExpr.Envs)-1 {
		return owner
	}
	if owner.GetType().PVT() != ast.Environment {
		w.Error(envExpr.Envs[depth-1].GetToken(), fmt.Sprintf("Resolve failed: expected type Environment, got %s", owner.GetType().ToString()))
		return &wkr.Invalid{}
	}
	env := owner.(*wkr.EnvironmentVal)
	previousDepth := depth
	depth += 1
	switch node := envExpr.Envs[previousDepth].(type) {
	case *ast.IdentifierExpr:
		if owner == nil {
			for i, v := range *w.Walkers {
				if v.Environment.Type.Name == node.Name.Lexeme {
					return GetEnvFieldType(w, envExpr, (*w.Walkers)[i].Environment, depth)
				}
			}
			if variable, found := w.Environment.Scope.Variables[node.Name.Lexeme]; found {
				return GetEnvFieldType(w, envExpr, variable.Value, depth)
			} else {
				w.Error(node.Name, "Resolve failed: couldn't find environment or variable named so")
				return &wkr.Invalid{}
			}

		}
		if depth == len(envExpr.Envs)-1 {
			if typee, found := env.Structs[node.Name.Lexeme]; found {
				return GetEnvFieldType(w, envExpr, typee, depth)
			}
			if typee, found := env.Variables[node.Name.Lexeme]; found && typee.GetType().PVT() == ast.Enum {
				return GetEnvFieldType(w, envExpr, typee, depth)
			}
		}
		if child, found := env.Childern[node.Name.Lexeme]; found {
			return GetEnvFieldType(child, envExpr, child.Environment, depth)
		}
		if variable, found := env.Variables[node.Name.Lexeme]; found {
			return GetEnvFieldType(w, envExpr, variable, depth)
		}
	case *ast.MapExpr, *ast.ListExpr, *ast.SelfExpr, *ast.DirectiveExpr, *ast.GroupExpr, *ast.BinaryExpr, *ast.AnonFnExpr, *ast.AnonStructExpr, *ast.MatchExpr:
		w.Error(node.GetToken(), fmt.Sprintf("Resolve failed: Cannot have a %s in %s for a type", node.GetType(), envExpr.GetType()))
	}

	return &wkr.Invalid{}
}
