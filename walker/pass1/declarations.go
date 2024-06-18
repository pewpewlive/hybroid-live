package pass1

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
	wkr "hybroid/walker"
)

func EnvStmt(w *wkr.Walker, node *ast.EnvironmentStmt, scope *wkr.Scope) {
	if scope.Environment.Type.Name != "UNKNOWN" {
		w.Error(node.GetToken(), "can't have more than one environment statement in a file")
		return
	}

	for i, v := range node.Env.Envs {
		if i < len(node.Env.Envs)-1 {
			scope.Environment.Type.Name += v.Lexeme + "::"
		}else {
			scope.Environment.Type.Name += v.Lexeme
		}
	}
	 
	if wlkr, found := (*w.Walkers)[w.Environment.Type.Name]; found {
		w.Error(node.GetToken(), fmt.Sprintf("cannot have two environments with the same name, path: %s",wlkr.Environment.Type.Path))
		return
	}

	(*w.Walkers)[w.Environment.Type.Name] = w
}

func FieldDeclaration(w *wkr.Walker, node *ast.FieldDeclarationStmt, container wkr.FieldContainer, scope *wkr.Scope) {
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
			explicitType := w.TypeExpr(node.Types[i])
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

	variables := VariableDeclaration(w, &varDecl, scope)
	node.Values = varDecl.Values
	for i := range variables {
		container.AddField(variables[i])
	}
}

func StructDeclaration(w *wkr.Walker, node *ast.StructDeclarationStmt, scope *wkr.Scope) {
	structVal := &wkr.StructVal{
		Type: *wkr.NewNamedType(node.Name.Lexeme),
		Fields: make([]*wkr.VariableVal, 0),
		Methods: map[string]*wkr.VariableVal{},
		Params: wkr.Types{},
	}

	structScope := wkr.NewScope(scope, &wkr.StructTag{StructVal: structVal}, wkr.SelfAllowing)

	params := make([]wkr.Type, 0)
	for _, param := range node.Constructor.Params {
		params = append(params, w.TypeExpr(param.Type))
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
		FieldDeclaration(w, &node.Fields[i], structVal, &structScope)
	}

	for i := range *node.Methods {
		params := make([]wkr.Type, 0)
		for _, param := range (*node.Methods)[i].Params {
			params = append(params, w.TypeExpr(param.Type))
		}

		ret := wkr.EmptyReturn
		for _, typee := range (*node.Methods)[i].Return {
			ret = append(ret, w.TypeExpr(typee))
			//fmt.Printf("%s\n", ret.values[len(ret.values)-1].Type.ToString())
		}
		variable := &wkr.VariableVal{
			Name:  (*node.Methods)[i].Name.Lexeme,
			Value: &wkr.FunctionVal{Params: params, Returns: ret},
			Token:  (*node.Methods)[i].GetToken(),
		}
		w.DeclareVariable(&structScope, variable, (*node.Methods)[i].Name)
		structVal.Methods[variable.Name] = variable
	}

	for i := range *node.Methods {
		MethodDeclaration(w, &(*node.Methods)[i], structVal, &structScope)
	}

	MethodDeclaration(w, &funcDeclaration, structVal, &structScope)
}

func MethodDeclaration(w *wkr.Walker, node *ast.MethodDeclarationStmt, container wkr.MethodContainer, scope *wkr.Scope) {
	funcExpr := ast.FunctionDeclarationStmt{
		Name:    node.Name,
		Return:  node.Return,
		Params:  node.Params,
		Body:    node.Body,
		IsLocal: true,
	}

	variable := FunctionDeclaration(w, &funcExpr, scope, wkr.Method)
	node.Body = funcExpr.Body
	container.AddMethod(variable)
}

func FunctionDeclaration(w *wkr.Walker, node *ast.FunctionDeclarationStmt, scope *wkr.Scope, procType wkr.ProcedureType) *wkr.VariableVal {
	ret := wkr.EmptyReturn
	for _, typee := range node.Return {
		ret = append(ret, w.TypeExpr(typee))
		//fmt.Printf("%s\n", ret.values[len(ret.values)-1].Type.ToString())
	}
	funcTag := &wkr.FuncTag{ReturnType: ret}
	fnScope := wkr.NewScope(scope, funcTag, wkr.ReturnAllowing)

	params := make([]wkr.Type, 0)
	for i, param := range node.Params {
		params = append(params, w.TypeExpr(param.Type))
		value := w.TypeToValue(params[i])
		w.DeclareVariable(&fnScope, &wkr.VariableVal{Name: param.Name.Lexeme, Value: value, Token: node.GetToken()}, node.Params[i].Name)
	}

	variable := &wkr.VariableVal{
		Name:  node.Name.Lexeme,
		Value: &wkr.FunctionVal{Params: params, Returns: ret},
		Token:  node.GetToken(),
	}
	if procType == wkr.Function {
		w.DeclareVariable(scope, variable, node.Name)
	}

	if scope.Parent != nil && !node.IsLocal {
		w.Error(node.GetToken(), "cannot declare a global function inside a local block")
	}

	WalkBody(w, &node.Body, funcTag, &fnScope)

	if funcTag, ok := fnScope.Tag.(*wkr.FuncTag); ok {
		if !funcTag.GetIfExits(wkr.Return) && !ret.Eq(&wkr.EmptyReturn) {
			w.Error(node.GetToken(), "not all code paths return a value")
		}
	}

	return variable
}

func EnumDeclaration(w *wkr.Walker, node *ast.EnumDeclarationStmt, scope *wkr.Scope) {
	enumVal := &wkr.EnumVal{
		Type:wkr.NewEnumType(node.Name.Lexeme),
	}

	if len(node.Fields) == 0 {
		w.Error(node.GetToken(), "can't declare an enum with no fields")
	}
	for _, v := range node.Fields {
		variable := &wkr.VariableVal{
			Name: v.Lexeme,
			Value: &wkr.EnumFieldVal{Type:enumVal.Type},
			IsConst: true,
		}
		enumVal.AddField(variable)
	}

	enumVar := &wkr.VariableVal{
		Name: enumVal.Type.Name,
		Value: enumVal,
		IsConst: true,
	}

	w.DeclareVariable(scope, enumVar, node.GetToken())
}

func VariableDeclaration(w *wkr.Walker, declaration *ast.VariableDeclarationStmt, scope *wkr.Scope) []*wkr.VariableVal {
	declaredVariables := []*wkr.VariableVal{}

	idents := len(declaration.Identifiers)
	values := make([]wkr.Value, idents)

	for i := range values {
		values[i] = &wkr.Invalid{}
	}

	valuesLength := len(declaration.Values)
	if valuesLength > idents {
		w.Error(declaration.Token, "too many values provided in declaration")
		return declaredVariables
	}

	for i := range declaration.Values {

		exprValue := GetNodeValue(w, &declaration.Values[i], scope)
		if declaration.Values[i].GetType() == ast.SelfExpression {
			w.Error(declaration.Values[i].GetToken(), "cannot assign self to a variable")
		}
		if types, ok := exprValue.(*wkr.Types); ok { 
			temp := values[i:]
			values = values[:i]
			w.AddTypesToValues(&values, types)
			values = append(values, temp...)
		} else {
			values[i] = exprValue
		}
	}

	if !declaration.IsLocal {
		w.Error(declaration.Token, "cannot declare a global variable inside a local block")
	}
	if declaration.Token.Type == lexer.Const && scope.Parent != nil {
		w.Error(declaration.Token, "cannot declare a global constant inside a local block")
	}

	for i, ident := range declaration.Identifiers {
		if ident.Lexeme == "_" {
			continue
		}

		valType := values[i].GetType()

		if declaration.Types[i] != nil {
			explicitType := w.TypeExpr(declaration.Types[i])
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
			Value: values[i],
			Name:  ident.Lexeme,
			Token:  ident,
		}
		declaredVariables = append(declaredVariables, variable)
		w.DeclareVariable(scope, variable, lexer.Token{Lexeme: ident.Lexeme, Location: declaration.Token.Location})
	}

	return declaredVariables
}