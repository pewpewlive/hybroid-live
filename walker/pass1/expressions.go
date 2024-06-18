package pass1

import (
	"hybroid/ast"
	"hybroid/lexer"
	wkr "hybroid/walker"
)

func AnonStructExprPass1(w *wkr.Walker, node *ast.AnonStructExpr, scope *wkr.Scope) *wkr.AnonStructVal {
	structTypeVal := &wkr.AnonStructVal{
		Fields:       map[string]*wkr.VariableVal{},
	}

	for i := range node.Fields {
		FieldDeclaration(w, node.Fields[i], structTypeVal, scope)
	}

	return structTypeVal
}

func AnonFnExprPass1(w *wkr.Walker, fn *ast.AnonFnExpr, scope *wkr.Scope) *wkr.FunctionVal {
	ret := wkr.EmptyReturn
	for _, typee := range fn.Return {
		ret = append(ret, w.TypeExpr(typee))
	}

	funcTag := &wkr.FuncTag{ReturnType: ret}
	fnScope := wkr.NewScope(scope, funcTag)
	fnScope.Attributes.Add(wkr.ReturnAllowing)

	params := make([]wkr.Type, 0)
	for i, param := range fn.Params {
		params = append(params, w.TypeExpr(param.Type))
		value := w.TypeToValue(params[i])
		w.DeclareVariable(&fnScope, &wkr.VariableVal{Name: param.Name.Lexeme, Value: value, Token: param.Name}, param.Name)
	}

	WalkBody(w, &fn.Body, funcTag, &fnScope)

	if !funcTag.GetIfExits(wkr.Return) && !ret.Eq(&wkr.EmptyReturn) {
		w.Error(fn.GetToken(), "not all code paths return a value")
	}

	return &wkr.FunctionVal{
		Params:  params,
		Returns: ret,
	}
}

func MatchExpr(w *wkr.Walker, node *ast.MatchExpr, scope *wkr.Scope) wkr.Value {
	casesLength := len(node.MatchStmt.Cases) + 1
	if node.MatchStmt.HasDefault {
		casesLength--
	}
	matchScope := wkr.NewScope(scope, &wkr.MatchExprTag{})
	matchScope.Attributes.Add(wkr.YieldAllowing)
	mtt := &wkr.MatchExprTag{Mpt: wkr.NewMultiPathTag(casesLength, matchScope.Attributes...)}
	matchScope.Tag = mtt

	for i := range node.MatchStmt.Cases {
		caseScope := wkr.NewScope(&matchScope, &wkr.UntaggedTag{})
		WalkBody(w, &node.MatchStmt.Cases[i].Body, mtt, &caseScope)
	}

	return mtt.YieldValues
}

func BinaryExpr(w *wkr.Walker, node *ast.BinaryExpr, scope *wkr.Scope) wkr.Value {
	left, right := GetNodeValue(&node.Left, scope), GetNodeValue(&node.Right, scope)
	leftType, rightType := left.GetType(), right.GetType()
	op := node.Operator
	switch op.Type {
	case lexer.Plus, lexer.Minus, lexer.Caret, lexer.Star, lexer.Slash, lexer.Modulo:
		w.validateArithmeticOperands(leftType, rightType, *node)
	default:
		if !TypeEquals(leftType, rightType) {
			w.Error(node.GetToken(), fmt.Sprintf("invalid comparison: types are not the same (left: %s, right: %s)",leftType.ToString(), rightType.ToString()))
		} else {
			return &BoolVal{}
		}
	}
	typ := w.DetermineValueType(leftType, rightType)

	if typ.PVT() == ast.Invalid {
		w.Error(node.GetToken(), fmt.Sprintf("invalid binary expression (left: %s, right: %s)",leftType.ToString(), rightType.ToString()))
		return &Invalid{}
	} else {
		return &BoolVal{}
	}
}

func EnvExpr(w *wkr.Walker, node *ast.EnvExpr, scope *wkr.Scope) wkr.Value {
	return &wkr.UnresolvedVal{
		EnvAccess: node,
	}
}