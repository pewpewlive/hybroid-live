package pass1

import (
	"hybroid/ast"
	"hybroid/helpers"
	"hybroid/lexer"
	wkr "hybroid/walker"
)

func StructDeclarationStmt(w *wkr.Walker, node *ast.StructDeclarationStmt, scope *wkr.Scope) {
	structVal := &wkr.StructVal{
		Type:    *wkr.NewNamedType(node.Name.Lexeme),
		IsLocal: node.IsLocal,
		Fields:  make([]*wkr.VariableVal, 0),
		Methods: map[string]*wkr.VariableVal{},
		Params:  wkr.Types{},
	}

	structScope := wkr.NewScope(scope, &wkr.StructTag{StructVal: structVal}, wkr.SelfAllowing)

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

	if scope.Parent != nil && !node.IsLocal {
		w.Error(node.GetToken(), "cannot declare a global function inside a local block")
	}

	WalkBody(w, &node.Body, funcTag, fnScope)

	return variable
}

func EnumDeclarationStmt(w *wkr.Walker, node *ast.EnumDeclarationStmt, scope *wkr.Scope) {
	enumVal := &wkr.EnumVal{
		Type: wkr.NewEnumType(node.Name.Lexeme),
	}

	if len(node.Fields) == 0 {
		w.Error(node.GetToken(), "can't declare an enum with no fields")
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

	if !declaration.IsLocal && scope.Parent != nil {
		w.Error(declaration.Token, "cannot declare a global variable inside a local block")
	}
	if declaration.Token.Type == lexer.Const && scope.Parent != nil {
		w.Error(declaration.Token, "cannot declare a global constant inside a local block")
	}

	for i, ident := range declaration.Identifiers {
		if ident.Lexeme == "_" {
			continue
		}

		variable := &wkr.VariableVal{
			Value:   values[i],
			Name:    ident.Lexeme,
			IsLocal: declaration.IsLocal,
			Token:   ident,
		}
		declaredVariables = append(declaredVariables, variable)
		w.DeclareVariable(scope, variable, ident)
	}

	return declaredVariables
}
func IfStmt(w *wkr.Walker, node *ast.IfStmt, scope *wkr.Scope) {
	length := len(node.Elseifs) + 2
	mpt := wkr.NewMultiPathTag(length, scope.Attributes...)
	multiPathScope := wkr.NewScope(scope, mpt)
	ifScope := wkr.NewScope(multiPathScope, &wkr.UntaggedTag{})

	WalkBody(w, &node.Body, mpt, ifScope)

	for i := range node.Elseifs {
		ifScope := wkr.NewScope(multiPathScope, &wkr.UntaggedTag{})
		WalkBody(w, &node.Elseifs[i].Body, mpt, ifScope)
	}

	if node.Else != nil {
		elseScope := wkr.NewScope(multiPathScope, &wkr.UntaggedTag{})
		WalkBody(w, &node.Else.Body, mpt, elseScope)
	}

	w.ReportExits(mpt, scope)
}

func RepeatStmt(w *wkr.Walker, node *ast.RepeatStmt, scope *wkr.Scope) {
	repeatScope := wkr.NewScope(scope, &wkr.MultiPathTag{}, wkr.BreakAllowing, wkr.ContinueAllowing)
	lt := wkr.NewMultiPathTag(1, repeatScope.Attributes...)
	repeatScope.Tag = lt

	if node.Variable != nil {
		w.DeclareVariable(repeatScope, &wkr.VariableVal{Name: node.Variable.Name.Lexeme, Value: &wkr.Invalid{}, IsLocal: true}, node.Variable.Name)
	}

	WalkBody(w, &node.Body, lt, repeatScope)

	w.ReportExits(lt, scope)
}

func WhileStmt(w *wkr.Walker, node *ast.WhileStmt, scope *wkr.Scope) {
	whileScope := wkr.NewScope(scope, &wkr.MultiPathTag{}, wkr.BreakAllowing, wkr.ContinueAllowing)
	lt := wkr.NewMultiPathTag(1, whileScope.Attributes...)
	whileScope.Tag = lt

	_ = GetNodeValue(w, &node.Condtion, scope)

	WalkBody(w, &node.Body, lt, whileScope)
}

func ForloopStmt(w *wkr.Walker, node *ast.ForStmt, scope *wkr.Scope) {
	forScope := wkr.NewScope(scope, &wkr.MultiPathTag{}, wkr.BreakAllowing, wkr.ContinueAllowing)
	lt := wkr.NewMultiPathTag(1, forScope.Attributes...)
	forScope.Tag = lt

	if len(node.KeyValuePair) != 0 {
		w.DeclareVariable(forScope,
			&wkr.VariableVal{Name: node.KeyValuePair[0].Name.Lexeme, Value: &wkr.NumberVal{}},
			node.KeyValuePair[0].Name)
	}
	if len(node.KeyValuePair) == 2 {
		w.DeclareVariable(forScope,
			&wkr.VariableVal{Name: node.KeyValuePair[1].Name.Lexeme, Value: &wkr.Unknown{}},
			node.KeyValuePair[1].Name)
	}

	WalkBody(w, &node.Body, lt, forScope)

	w.ReportExits(lt, scope)
}

func TickStmt(w *wkr.Walker, node *ast.TickStmt, scope *wkr.Scope) {
	funcTag := &wkr.FuncTag{ReturnTypes: wkr.EmptyReturn}
	tickScope := wkr.NewScope(scope, funcTag, wkr.ReturnAllowing)
	
	if node.Variable.GetValueType() != ast.Unknown {
		w.DeclareVariable(tickScope, &wkr.VariableVal{Name: node.Variable.Name.Lexeme, Value: &wkr.NumberVal{}}, node.Token)
	}

	WalkBody(w, &node.Body, funcTag, tickScope)
}

func MatchStmt(w *wkr.Walker, node *ast.MatchStmt, scope *wkr.Scope) {
	casesLength := len(node.Cases) + 1
	if node.HasDefault {
		casesLength--
	}
	mpt := wkr.NewMultiPathTag(casesLength, scope.Attributes...)
	multiPathScope := wkr.NewScope(scope, mpt)

	// var has_default bool
	for i := range node.Cases {
		caseScope := wkr.NewScope(multiPathScope, &wkr.UntaggedTag{})

		WalkBody(w, &node.Cases[i].Body, mpt, caseScope)
	}

	w.ReportExits(mpt, scope)
}

func BreakStmt(w *wkr.Walker, node *ast.BreakStmt, scope *wkr.Scope) {
	if !scope.Is(wkr.BreakAllowing) {
		w.Error(node.GetToken(), "cannot use break outside of loops")
	}

	if returnable := scope.ResolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, wkr.Break)
		(*returnable).SetExit(true, wkr.All)
	}
}

func ContinueStmt(w *wkr.Walker, node *ast.ContinueStmt, scope *wkr.Scope) {
	if !scope.Is(wkr.ContinueAllowing) {
		w.Error(node.GetToken(), "cannot use break outside of loops")
	}

	if returnable := scope.ResolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, wkr.Continue)
		(*returnable).SetExit(true, wkr.All)
	}
}

func ReturnStmt(w *wkr.Walker, node *ast.ReturnStmt, scope *wkr.Scope) *wkr.Types {
	if !scope.Is(wkr.ReturnAllowing) {
		w.Error(node.GetToken(), "can't have a return statement outside of a function or method")
	}

	sc, _, _ := wkr.ResolveTagScope[*wkr.FuncTag](scope)
	for i := range node.Args {
		GetNodeValue(w, &node.Args[i], scope)
	}
	if sc == nil {
		return &wkr.EmptyReturn
	}

	if returnable := scope.ResolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, wkr.Return)
		(*returnable).SetExit(true, wkr.All)
	}

	return &wkr.EmptyReturn
}

func YieldStmt(w *wkr.Walker, node *ast.YieldStmt, scope *wkr.Scope) *wkr.Types {
	if !scope.Is(wkr.YieldAllowing) {
		w.Error(node.GetToken(), "cannot use yield outside of statement expressions") // wut
	}

	sc, _, _:= wkr.ResolveTagScope[*wkr.MatchExprTag](scope)
	yieldTypes := wkr.EmptyReturn
	for i := range node.Args {
		val := GetNodeValue(w, &node.Args[i], scope)
		valType := val.GetType()
		if types, ok := val.(*wkr.Types); ok {
			yieldTypes = append(yieldTypes, *types...)
		} else {
			yieldTypes = append(yieldTypes, valType)
		}
	}
	if sc, _, met :=  wkr.ResolveTagScope[*wkr.MatchExprTag](scope); sc != nil {
		if helpers.ListsAreSame((*met).YieldValues, wkr.EmptyReturn) {
			(*met).YieldValues = yieldTypes
		}else {
			(*met).YieldValues = append((*met).YieldValues, yieldTypes...)
		}
	}
	if sc == nil {
		return &wkr.EmptyReturn
	}

	if returnable := scope.ResolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, wkr.Yield)
		(*returnable).SetExit(true, wkr.All)
	}

	return &wkr.EmptyReturn
}