package pass3

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
	wkr "hybroid/walker"
	"strings"
)

func StructDeclarationStmt(w *wkr.Walker, node *ast.StructDeclarationStmt, scope *wkr.Scope) {
	structVal := &wkr.StructVal{
		Type:    *wkr.NewNamedType(node.Name.Lexeme),
		Fields:  make([]*wkr.VariableVal, 0),
		Methods: map[string]*wkr.VariableVal{},
		Params:  wkr.Types{},
	}

	structTag := &wkr.StructTag{StructVal: structVal}
	structScope := scope.AccessChild()
	scope.Tag = structTag

	params := make([]wkr.Type, 0)
	for _, param := range node.Constructor.Params {
		params = append(params, TypeExpr(w, param.Type))
	}
	structVal.Params = params

	w.DeclareStruct(structVal)

	funcDeclaration := ast.MethodDeclarationStmt{
		Name:    node.Constructor.Token,
		Params:  node.Constructor.Params,
		Return:  node.Constructor.Return,
		IsLocal: true,
		Body:    node.Constructor.Body,
	}

	for i := range node.Fields {
		FieldDeclarationStmt(w, &node.Fields[i], structVal, structScope)
	}

	for i := range *node.Methods {
		params := make([]wkr.Type, 0)
		for _, param := range (*node.Methods)[i].Params {
			params = append(params, TypeExpr(w, param.Type))
		}

		ret := wkr.EmptyReturn
		for _, typee := range (*node.Methods)[i].Return {
			ret = append(ret, TypeExpr(w, typee))
			//fmt.Printf("%s\n", ret.values[len(ret.values)-1].Type.ToString())
		}
		variable := &wkr.VariableVal{
			Name:    (*node.Methods)[i].Name.Lexeme,
			Value:   &wkr.FunctionVal{Params: params, Returns: ret},
			IsLocal: node.IsLocal,
			Token:   (*node.Methods)[i].GetToken(),
		}
		w.DeclareVariable(structScope, variable, (*node.Methods)[i].Name)
		structVal.Methods[variable.Name] = variable
	}

	for i := range *node.Methods {
		MethodDeclarationStmt(w, &(*node.Methods)[i], structVal, structScope)
	}

	MethodDeclarationStmt(w, &funcDeclaration, structVal, structScope)
}

func FieldDeclarationStmt(w *wkr.Walker, node *ast.FieldDeclarationStmt, container wkr.FieldContainer, scope *wkr.Scope) {
	varDecl := ast.VariableDeclarationStmt{
		Identifiers: node.Identifiers,
		Types:       node.Types,
		Values:      node.Values,
		IsLocal:     true,
		Token:       node.Token,
	}
	structType := container.GetType()
	if len(node.Types) != 0 {
		for i := range node.Types {
			explicitType := TypeExpr(w, node.Types[i])
			if wkr.TypeEquals(explicitType, structType) {
				w.Error(node.Types[i].GetToken(), "cannot have a field with a value type of its struct")
				return
			}
		}
	} else if len(node.Types) != 0 {
		for i := range node.Values {
			valType := GetNodeValue(w, &node.Values[i], scope).GetType()
			if wkr.TypeEquals(valType, structType) {
				w.Error(node.Types[i].GetToken(), "cannot have a field with a value type of its struct")
				return
			}
		}
	}

	variables := VariableDeclarationStmt(w, &varDecl, scope)
	node.Values = varDecl.Values
	for i := range variables {
		container.AddField(variables[i])
	}
}

func MethodDeclarationStmt(w *wkr.Walker, node *ast.MethodDeclarationStmt, container wkr.MethodContainer, scope *wkr.Scope) {
	funcExpr := ast.FunctionDeclarationStmt{
		Name:    node.Name,
		Return:  node.Return,
		Params:  node.Params,
		Body:    node.Body,
		IsLocal: true,
	}

	variable := FunctionDeclarationStmt(w, &funcExpr, scope, wkr.Method)
	node.Body = funcExpr.Body
	container.AddMethod(variable)
}

func FunctionDeclarationStmt(w *wkr.Walker, node *ast.FunctionDeclarationStmt, scope *wkr.Scope, procType wkr.ProcedureType) *wkr.VariableVal {
	ret := wkr.EmptyReturn
	for _, typee := range node.Return {
		ret = append(ret, TypeExpr(w, typee))
	}
	funcTag := &wkr.FuncTag{ReturnTypes: ret}
	fnScope := wkr.NewScope(scope, funcTag, wkr.ReturnAllowing)

	params := make([]wkr.Type, 0)
	for i, param := range node.Params {
		params = append(params, TypeExpr(w, param.Type))
		value := w.TypeToValue(params[i])
		w.DeclareVariable(fnScope, &wkr.VariableVal{
			Name:    param.Name.Lexeme,
			Value:   value,
			IsLocal: node.IsLocal,
			Token:   node.GetToken(),
		}, node.Params[i].Name)
	}

	variable := &wkr.VariableVal{
		Name:  node.Name.Lexeme,
		Value: &wkr.FunctionVal{Params: params, Returns: ret},
		Token: node.GetToken(),
	}
	if procType == wkr.Function {
		w.DeclareVariable(scope, variable, node.Name)
	}

	return variable
}

func EnumDeclarationStmt(w *wkr.Walker, node *ast.EnumDeclarationStmt, scope *wkr.Scope) {
	enumVal := &wkr.EnumVal{
		Type: wkr.NewEnumType(node.Name.Lexeme),
	}

	for _, v := range node.Fields {
		variable := &wkr.VariableVal{
			Name:    v.Lexeme,
			Value:   &wkr.EnumFieldVal{Type: enumVal.Type},
			IsLocal: node.IsLocal,
			IsConst: true,
		}
		enumVal.AddField(variable)
	}

	enumVar := &wkr.VariableVal{
		Name:    enumVal.Type.Name,
		Value:   enumVal,
		IsLocal: node.IsLocal,
		IsConst: true,
	}

	w.DeclareVariable(scope, enumVar, node.GetToken())
}

func VariableDeclarationStmt(w *wkr.Walker, declaration *ast.VariableDeclarationStmt, scope *wkr.Scope) []*wkr.VariableVal {
	declaredVariables := []*wkr.VariableVal{}

	idents := len(declaration.Identifiers)
	values := make([]wkr.Value, idents)

	for i := range values {
		values[i] = &wkr.Invalid{}
	}

	valuesLength := len(declaration.Values)
	if valuesLength > idents {
		return declaredVariables
	}

	for i := range declaration.Values {

		exprValue := GetNodeValue(w, &declaration.Values[i], scope)
		if types, ok := exprValue.(*wkr.Types); ok {
			temp := values[i:]
			values = values[:i]
			w.AddTypesToValues(&values, types)
			values = append(values, temp...)
		} else {
			values[i] = exprValue
		}
	}

	for i, ident := range declaration.Identifiers {
		if ident.Lexeme == "_" {
			continue
		}

		valType := values[i].GetType()

		if declaration.Types[i] != nil {
			explicitType := TypeExpr(w, declaration.Types[i])
			if valType == wkr.InvalidType && explicitType != wkr.InvalidType {
				values[i] = w.TypeToValue(explicitType)
				declaration.Values = append(declaration.Values, values[i].GetDefault())
			} else if !wkr.TypeEquals(valType, explicitType) {
				w.Error(declaration.Token, fmt.Sprintf("mismatched types: value type (%s) not the same with explict type (%s)",
					valType.ToString(),
					explicitType.ToString()))
			}
		}

		variable := &wkr.VariableVal{
			Value:   values[i],
			Name:    ident.Lexeme,
			IsLocal: declaration.IsLocal,
			Token:   ident,
		}
		declaredVariables = append(declaredVariables, variable)
		w.DeclareVariable(scope, variable, lexer.Token{Lexeme: ident.Lexeme, Location: declaration.Token.Location})
	}

	return declaredVariables
}

func ReturnStmt(w *wkr.Walker, node *ast.ReturnStmt, scope *wkr.Scope) *wkr.Types {
	if !scope.Is(wkr.ReturnAllowing) {
		w.Error(node.GetToken(), "can't have a return statement outside of a function or method")
	}

	ret := wkr.EmptyReturn
	for i := range node.Args {
		val := GetNodeValue(w, &node.Args[i], scope)
		valType := val.GetType()
		if types, ok := val.(*wkr.Types); ok {
			ret = append(ret, *types...)
		} else {
			ret = append(ret, valType)
		}
	}
	sc, _, funcTag := wkr.ResolveTagScope[*wkr.FuncTag](scope)
	if sc == nil {
		return &ret
	}

	if returnable := scope.ResolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, wkr.Return)
		(*returnable).SetExit(true, wkr.All)
	}

	errorMsg := w.ValidateReturnValues(ret, (*funcTag).ReturnTypes)
	if errorMsg != "" {
		w.Error(node.GetToken(), errorMsg)
	}

	return &ret
}

func YieldStmt(w *wkr.Walker, node *ast.YieldStmt, scope *wkr.Scope) *wkr.Types {
	if !scope.Is(wkr.YieldAllowing) {
		w.Error(node.GetToken(), "cannot use yield outside of statement expressions") // wut
	}

	ret := wkr.EmptyReturn
	for i := range node.Args {
		val := GetNodeValue(w, &node.Args[i], scope)
		valType := val.GetType()
		if types, ok := val.(*wkr.Types); ok {
			ret = append(ret, *types...)
		} else {
			ret = append(ret, valType)
		}
	}

	sc, _, matchExprT := wkr.ResolveTagScope[*wkr.MatchExprTag](scope)

	if sc == nil {
		return &ret
	}

	matchExprTag := *matchExprT

	if matchExprTag.YieldValues == nil {
		matchExprTag.YieldValues = &ret
	} else {
		errorMsg := w.ValidateReturnValues(ret, *matchExprTag.YieldValues)
		if errorMsg != "" {
			errorMsg = strings.Replace(errorMsg, "return", "yield", -1)
			w.Error(node.GetToken(), errorMsg)
		}
	}

	if returnable := scope.ResolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, wkr.Yield)
		(*returnable).SetExit(true, wkr.All)
	}

	return &ret
}