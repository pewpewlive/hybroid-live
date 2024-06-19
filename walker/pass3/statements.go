package pass3

import (
	"fmt"
	"hybroid/ast"
	"hybroid/parser"
	wkr "hybroid/walker"
)

func IfStmt(w *wkr.Walker, node *ast.IfStmt, scope *wkr.Scope) {
	length := len(node.Elseifs) + 2
	mpt := wkr.NewMultiPathTag(length, scope.Attributes...)
	multiPathScope := wkr.NewScope(scope, mpt)
	ifScope := wkr.NewScope(&multiPathScope, &wkr.UntaggedTag{})
	boolExpr := GetNodeValue(w, &node.BoolExpr, scope)
	if boolExpr.GetType().PVT() != ast.Bool {
		w.Error(node.BoolExpr.GetToken(), "if condition is not a comparison")
	}
	WalkBody(w, &node.Body, mpt, &ifScope)

	for i := range node.Elseifs {
		boolExpr := GetNodeValue(w, &node.Elseifs[i].BoolExpr, scope)
		if boolExpr.GetType().PVT() != ast.Bool {
			w.Error(node.Elseifs[i].BoolExpr.GetToken(), "if condition is not a comparison")
		}
		ifScope := wkr.NewScope(&multiPathScope, &wkr.UntaggedTag{})
		WalkBody(w, &node.Elseifs[i].Body, mpt, &ifScope)
	}

	if node.Else != nil {
		elseScope := wkr.NewScope(&multiPathScope, &wkr.UntaggedTag{})
		WalkBody(w, &node.Else.Body, mpt, &elseScope)
	}

	returnabl := scope.ResolveReturnable()

	if returnabl == nil {
		return
	}
	returnable := *returnabl

	returnable.SetExit(mpt.GetIfExits(wkr.Return), wkr.Return)
	returnable.SetExit(mpt.GetIfExits(wkr.Yield), wkr.Yield)
	returnable.SetExit(mpt.GetIfExits(wkr.Break), wkr.Break)
	returnable.SetExit(mpt.GetIfExits(wkr.Continue), wkr.Continue)
	returnable.SetExit(mpt.GetIfExits(wkr.All), wkr.All)
}

func AssignmentStmt(w *wkr.Walker, assignStmt *ast.AssignmentStmt, scope *wkr.Scope) {
	hasFuncs := false

	wIdents := []wkr.Value{}
	for i := range assignStmt.Identifiers {
		wIdents = append(wIdents, GetNodeValue(w, &assignStmt.Identifiers[i], scope))
	}

	for i := range assignStmt.Values {
		if assignStmt.Values[i].GetType() == ast.CallExpression {
			hasFuncs = true
		}
		value := GetNodeValue(w, &assignStmt.Values[i], scope)
		if i > len(wIdents)-1 {
			break
		}
		variableType := wIdents[i].GetType()
		valueType := value.GetType()
		if variableType.PVT() == ast.Invalid {
			w.Error(assignStmt.Identifiers[i].GetToken(), "cannot assign a value to an undeclared variable")
			continue
		}

		if !wkr.TypeEquals(variableType, valueType) {
			w.Error(assignStmt.Values[i].GetToken(), fmt.Sprintf("mismatched types: variable has a type of %s, but a value of %s was given to it.", variableType.ToString(), valueType.ToString()))
		}

		variable, ok := wIdents[i].(*wkr.VariableVal)

		if ok {
			if _, err := scope.AssignVariable(variable, value); err != nil {
				err.Token = variable.Token
				w.AddError(*err)
			}
		}
	}

	if hasFuncs {
		w.Error(assignStmt.GetToken(), "cannot have a function call in assignment")
	} else if len(assignStmt.Values) < len(assignStmt.Identifiers) {
		w.Error(assignStmt.Values[len(assignStmt.Values)-1].GetToken(), "not enough values provided in assignment")
	} else if len(assignStmt.Values) > len(assignStmt.Identifiers) {
		w.Error(assignStmt.Values[len(assignStmt.Values)-1].GetToken(), "too many values provided in assignment")
	}
}

func RepeatStmt(w *wkr.Walker, node *ast.RepeatStmt, scope *wkr.Scope) {
	repeatScope := scope.AccessChild()
	lt := repeatScope.Tag.(*wkr.LoopTag)

	end := GetNodeValue(w, &node.Iterator, scope)
	endType := end.GetType()
	if !parser.IsFx(endType.PVT()) && endType.PVT() != ast.Number {
		w.Error(node.Iterator.GetToken(), "invalid value type of iterator")
	} else if variable, ok := end.(*wkr.VariableVal); ok {
		if fixedpoint, ok := variable.Value.(*wkr.FixedVal); ok {
			endType = wkr.NewBasicType(fixedpoint.SpecificType)
		}
	} else {
		if fixedpoint, ok := end.(*wkr.FixedVal); ok {
			endType = wkr.NewBasicType(fixedpoint.SpecificType)
		}
	}
	if node.Start.GetType() == ast.NA {
		node.Start = &ast.LiteralExpr{Token: node.Start.GetToken(), ValueType: endType.PVT(), Value: "1"}
	}
	start := GetNodeValue(w, &node.Start, scope)
	if node.Skip.GetType() == ast.NA {
		node.Skip = &ast.LiteralExpr{Token: node.Skip.GetToken(), ValueType: endType.PVT(), Value: "1"}
	}
	skip := GetNodeValue(w, &node.Skip, scope)

	repeatType := end.GetType()
	startType := start.GetType()
	skipType := skip.GetType()

	if wkr.TypeEquals(repeatType, startType) && wkr.TypeEquals(startType, skipType) {
		w.Error(node.Start.GetToken(), fmt.Sprintf("all value types must be the same (iter:%s, start:%s, by:%s)", repeatType, startType, skipType))
	}

	WalkBody(w, &node.Body, lt, repeatScope)

	w.ReportExits(lt, scope)
}

func WhileStmt(w *wkr.Walker, node *ast.WhileStmt, scope *wkr.Scope) {
	whileScope := wkr.NewScope(scope, &wkr.LoopTag{}, wkr.BreakAllowing, wkr.ContinueAllowing)
	lt := wkr.NewLoopTag(whileScope.Attributes...)
	whileScope.Tag = lt

	_ = GetNodeValue(w, &node.Condtion, scope)

	WalkBody(w, &node.Body, lt, &whileScope)
}

func ForloopStmt(w *wkr.Walker, node *ast.ForStmt, scope *wkr.Scope) {
	forScope := scope.AccessChild()
	lt := forScope.Tag.(*wkr.LoopTag)

	if len(node.KeyValuePair) != 0 {
		w.DeclareVariable(forScope,
			&wkr.VariableVal{Name: node.KeyValuePair[0].Name.Lexeme, Value: &wkr.NumberVal{}},
			node.KeyValuePair[0].Name)
	}
	valType := GetNodeValue(w, &node.Iterator, scope).GetType()
	wrapper, ok := valType.(*wkr.WrapperType)
	if !ok {
		w.Error(node.Iterator.GetToken(), "iterator must be of type map or list")
	} else if len(node.KeyValuePair) == 2 {
		node.OrderedIteration = wrapper.PVT() == ast.List
		w.DeclareVariable(forScope,
			&wkr.VariableVal{Name: node.KeyValuePair[1].Name.Lexeme, Value: w.TypeToValue(wrapper.WrappedType)},
			node.KeyValuePair[1].Name)
	}

	WalkBody(w, &node.Body, lt, forScope)

	w.ReportExits(lt, scope)
}

func TickStmt(w *wkr.Walker, node *ast.TickStmt, scope *wkr.Scope) {
	tickScope := scope.AccessChild()

	if node.Variable.GetValueType() != ast.Unknown {
		w.DeclareVariable(tickScope, &wkr.VariableVal{Name: node.Variable.Name.Lexeme, Value: &wkr.NumberVal{}}, node.Token)
	}

	for i := range node.Body {
		WalkNode(w, &node.Body[i], tickScope)
	}
}

func VariableDeclarationStmt(w *wkr.Walker, declaration *ast.VariableDeclarationStmt, scope *wkr.Scope) []*wkr.VariableVal {
	declaredVariables := []*wkr.VariableVal{}

	idents := len(declaration.Identifiers)
	if len(declaration.Values) > idents {
		return declaredVariables
	}

	values := make([]wkr.Value, idents)

	for i := range values {
		values[i] = &wkr.Invalid{}
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

	return declaredVariables
}

func UseStmt(w *wkr.Walker, node *ast.UseStmt, scope *wkr.Scope) {
	variable := &wkr.VariableVal{
		Name: node.Variable.Name.Lexeme,
		Value: &wkr.EnvironmentVal{
			Type: &wkr.EnvironmentType{
				Name: node.Variable.Name.Lexeme,
			},
		},
		Token: node.GetToken(),
	}

	w.DeclareVariable(scope, variable, node.Variable.Name)
}

func MatchStmt(w *wkr.Walker, node *ast.MatchStmt, isExpr bool, scope *wkr.Scope) {
	val := GetNodeValue(w, &node.ExprToMatch, scope)
	if val.GetType().PVT() == ast.Invalid {
		w.Error(node.ExprToMatch.GetToken(), "variable is of type invalid")
	}
	valType := val.GetType()
	casesLength := len(node.Cases) + 1
	if node.HasDefault {
		casesLength--
	}
	mpt := wkr.NewMultiPathTag(casesLength, scope.Attributes...)
	multiPathScope := wkr.NewScope(scope, mpt)

	var has_default bool
	for i := range node.Cases {
		caseScope := wkr.NewScope(&multiPathScope, &wkr.UntaggedTag{})

		if !isExpr {
			WalkBody(w, &node.Cases[i].Body, mpt, &caseScope)
		}

		if node.Cases[i].Expression.GetToken().Lexeme == "_" {
			has_default = true
			continue
		}

		caseValType := GetNodeValue(w, &node.Cases[i].Expression, scope).GetType()
		if !wkr.TypeEquals(valType, caseValType) {
			w.Error(
				node.Cases[i].Expression.GetToken(),
				fmt.Sprintf("mismatched types: arm expression (%s) and match expression (%s)",
					caseValType.ToString(),
					valType.ToString()))
		}
	}

	if has_default && len(node.Cases) == 1 {
		w.Error(node.Cases[0].Expression.GetToken(), "cannot have a match statement/expression with one arm that is default")
	}

	if !has_default && isExpr {
		w.Error(node.GetToken(), "match expression has no default arm")
	}

	if isExpr {
		return
	}

	w.ReportExits(mpt, scope)
}

func FunctionDeclarationStmt(w *wkr.Walker, node *ast.FunctionDeclarationStmt, scope *wkr.Scope, procType wkr.ProcedureType) *wkr.VariableVal {
	ret := wkr.EmptyReturn
	for _, typee := range node.Return {
		ret = append(ret, TypeExpr(w, typee))
		//fmt.Printf("%s\n", ret.values[len(ret.values)-1].Type.ToString())
	}
	funcTag := &wkr.FuncTag{ReturnType: ret}
	fnScope := wkr.NewScope(scope, funcTag, wkr.ReturnAllowing)

	params := make([]wkr.Type, 0)
	for i, param := range node.Params {
		params = append(params, TypeExpr(w, param.Type))
		value := w.TypeToValue(params[i])
		w.DeclareVariable(&fnScope, &wkr.VariableVal{Name: param.Name.Lexeme, Value: value, Token: node.GetToken()}, node.Params[i].Name)
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

	WalkBody(w, &node.Body, funcTag, &fnScope)

	if funcTag, ok := fnScope.Tag.(*wkr.FuncTag); ok {
		if !funcTag.GetIfExits(wkr.Return) && !ret.Eq(&wkr.EmptyReturn) {
			w.Error(node.GetToken(), "not all code paths return a value")
		}
	}

	return variable
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

func StructDeclarationStmt(w *wkr.Walker, node *ast.StructDeclarationStmt, scope *wkr.Scope) {
	structVal := &wkr.StructVal{
		Type:    *wkr.NewNamedType(node.Name.Lexeme),
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
		FieldDeclarationStmt(w, &node.Fields[i], structVal, &structScope)
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
		w.DeclareVariable(&structScope, variable, (*node.Methods)[i].Name)
		structVal.Methods[variable.Name] = variable
	}

	for i := range *node.Methods {
		MethodDeclarationStmt(w, &(*node.Methods)[i], structVal, &structScope)
	}

	MethodDeclarationStmt(w, &funcDeclaration, structVal, &structScope)
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
