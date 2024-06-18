package pass2

import (
	"hybroid/ast"
	wkr "hybroid/walker"
)

func AnonFnExpr(w *wkr.Walker, fn *ast.AnonFnExpr, scope *wkr.Scope) *wkr.FunctionVal {
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
		Params:    params,
		Returns: ret,
	}
}