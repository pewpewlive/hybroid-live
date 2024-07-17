package pass2

import (
	"fmt"
	"hybroid/ast"
	"hybroid/helpers"
	"hybroid/parser"
	wkr "hybroid/walker"
	"strings"
)

func StructDeclarationStmt(w *wkr.Walker, node *ast.StructDeclarationStmt, scope *wkr.Scope) {
	if node.Constructor == nil {
		return;
	}
	structScope := scope.AccessChild()

	structVal := structScope.Tag.(*wkr.StructTag).StructVal

	params := make([]wkr.Type, 0)
	for _, param := range node.Constructor.Params {
		params = append(params, TypeExpr(w, param.Type, w.Environment))
	}
	structVal.Params = params

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

	for i := range node.Methods {
		params := make([]wkr.Type, 0)
		for _, param := range node.Methods[i].Params {
			params = append(params, TypeExpr(w, param.Type, w.Environment))
		}

		ret := wkr.EmptyReturn
		for _, typee := range node.Methods[i].Return {
			ret = append(ret, TypeExpr(w, typee, w.Environment))
			//fmt.Printf("%s\n", ret.values[len(ret.values)-1].Type.ToString())
		}
		variable := &wkr.VariableVal{
			Name:    node.Methods[i].Name.Lexeme,
			Value:   &wkr.FunctionVal{Params: params, Returns: ret},
			IsLocal: node.IsLocal,
			Token:   node.Methods[i].GetToken(),
		}
		*w.GetVariable(structScope, variable.Name) = *variable
	}

	for i := range node.Methods {
		MethodDeclarationStmt(w, &node.Methods[i], structVal, structScope)
	}

	MethodDeclarationStmt(w, &funcDeclaration, structVal, structScope)
}

func EntityDeclarationStmt(w *wkr.Walker, node *ast.EntityDeclarationStmt, scope *wkr.Scope) {
	entityScope := scope.AccessChild()

	entityVal, _ := w.GetEntity(node.Name.Lexeme)

	for i := range node.Fields { // debug?
		FieldDeclarationStmt(w, &node.Fields[i], entityVal, entityScope)
	}

	//spawn
	if node.Spawner != nil {
		EntityFunctionDeclarationStmt(w, node.Spawner, entityVal, entityScope)
	}
	//destroy
	if node.Destroyer != nil {
		EntityFunctionDeclarationStmt(w, node.Destroyer, entityVal, entityScope)
	}

	for i := range node.Callbacks {
		EntityFunctionDeclarationStmt(w, node.Callbacks[i], entityVal, entityScope)
	}
}

func EntityFunctionDeclarationStmt(w *wkr.Walker, node *ast.EntityFunctionDeclarationStmt, entityVal *wkr.EntityVal, scope *wkr.Scope) {

	fnScope := scope.AccessChild()
	params := make([]wkr.Type, 0)
	for _, param := range node.Params {
		params = append(params, TypeExpr(w, param.Type, w.Environment))
		w.GetVariable(fnScope, param.Name.Lexeme).Value = w.TypeToValue(params[len(params)-1])
	}
	WalkBody(w, &node.Body, fnScope)
} // that should do it

func FieldDeclarationStmt(w *wkr.Walker, node *ast.FieldDeclarationStmt, container wkr.FieldContainer, scope *wkr.Scope) {
	varDecl := ast.VariableDeclarationStmt{
		Identifiers: node.Identifiers,
		Types:       node.Types,
		Values:      node.Values,
		IsLocal:     true,
		Token:       node.Token,
	}
	// structType := container.GetType()
	// if len(node.Types) != 0 {
	// 	for i := range node.Types {
	// 		explicitType := TypeExpr(w, node.Types[i], w.Environment)
	// 		if wkr.TypeEquals(explicitType, structType) {
	// 			w.Error(node.Types[i].GetToken(), "cannot have a field with a value type of its struct")
	// 			return
	// 		}
	// 	}
	// } else if len(node.Types) != 0 {
	// 	for i := range node.Values {
	// 		valType := GetNodeValue(w, &node.Values[i], scope).GetType()
	// 		if wkr.TypeEquals(valType, structType) {
	// 			w.Error(node.Types[i].GetToken(), "cannot have a field with a value type of its struct")
	// 			return
	// 		}
	// 	}
	// }

	variables := VariableDeclarationStmt(w, &varDecl, scope)
	node.Values = varDecl.Values
	for i := range variables {
		variable, _, found := container.ContainsField(variables[i].Name)
		if found {
			variable.Value = variables[i].Value
		} else {
			container.AddField(variables[i])
		}
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
	method, _ := container.ContainsMethod(variable.Name)
	*method = *variable
}

func FunctionDeclarationStmt(w *wkr.Walker, node *ast.FunctionDeclarationStmt, scope *wkr.Scope, procType wkr.ProcedureType) *wkr.VariableVal {
	ret := wkr.EmptyReturn
	for _, typee := range node.Return {
		ret = append(ret, TypeExpr(w, typee, w.Environment))
	}
	fnScope := scope.AccessChild()
	funcTag := fnScope.Tag.(*wkr.FuncTag)
	funcTag.ReturnTypes = ret

	params := make([]wkr.Type, 0)
	for i, param := range node.Params {
		params = append(params, TypeExpr(w, param.Type, w.Environment))
		w.GetVariable(fnScope, param.Name.Lexeme).Value = w.TypeToValue(params[i])
	}

	variable := &wkr.VariableVal{
		Name:  node.Name.Lexeme,
		Value: &wkr.FunctionVal{Params: params, Returns: ret},
		Token: node.GetToken(),
	}
	if procType == wkr.Function {
		w.GetVariable(fnScope, variable.Name).Value = variable.Value
	}

	WalkBody(w, &node.Body, fnScope)

	return variable
}

func EnumDeclarationStmt(w *wkr.Walker, node *ast.EnumDeclarationStmt, scope *wkr.Scope) {
	enumVal := w.GetVariable(scope, node.Name.Lexeme).Value.(*wkr.EnumVal)

	for _, v := range node.Fields {
		variable := &wkr.VariableVal{
			Name:    v.Lexeme,
			Value:   &wkr.EnumFieldVal{Type: enumVal.Type},
			IsLocal: node.IsLocal,
			IsConst: true,
		}
		field, _, _ := enumVal.ContainsField(variable.Name)
		*field = *variable
	}

	enumVar := &wkr.VariableVal{
		Name:    enumVal.Type.Name,
		Value:   enumVal,
		IsLocal: node.IsLocal,
		IsConst: true,
	}

	*w.GetVariable(scope, enumVar.Name) = *enumVar
}

func VariableDeclarationStmt(w *wkr.Walker, declaration *ast.VariableDeclarationStmt, scope *wkr.Scope) []*wkr.VariableVal {
	variables := []*wkr.VariableVal{}

	idents := len(declaration.Identifiers)
	values := make([]wkr.Value, idents)

	for i := range values {
		values[i] = &wkr.Unknown{}
	}

	valuesLength := len(declaration.Values)
	if valuesLength > idents {
		w.Error(declaration.Token, "too many values provided in declaration")
		return variables
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

	for i, ident := range declaration.Identifiers {
		if ident.Lexeme == "_" {
			continue
		}
		if declaration.Types[i] != nil {
			valueType := values[i].GetType()
			explicitType := TypeExpr(w, declaration.Types[i], w.Environment)
			if valueType.PVT() == ast.Unknown {
				values[i] = w.TypeToValue(explicitType)
				declaration.Values = append(declaration.Values, values[i].GetDefault()) // only here does it set to default if no value is given
			} else if !wkr.TypeEquals(valueType, explicitType) {
				w.Error(declaration.Values[i].GetToken(), fmt.Sprintf("Given value is %s, but explicit type is %s", valueType.ToString(), explicitType.ToString()))
			}
		}
		_var := w.GetVariable(scope, ident.Lexeme)
		_var.Value = values[i]
		variables = append(variables, _var)
	}

	return variables
}
func IfStmt(w *wkr.Walker, node *ast.IfStmt, scope *wkr.Scope) {
	multiPathScope := scope.AccessChild()
	ifScope := multiPathScope.AccessChild()
	boolExpr := GetNodeValue(w, &node.BoolExpr, scope)
	if boolExpr.GetType().PVT() != ast.Bool {
		w.Error(node.BoolExpr.GetToken(), "if condition is not a comparison")
	}
	WalkBody(w, &node.Body, ifScope)

	for i := range node.Elseifs {
		boolExpr := GetNodeValue(w, &node.Elseifs[i].BoolExpr, scope)
		if boolExpr.GetType().PVT() != ast.Bool {
			w.Error(node.Elseifs[i].BoolExpr.GetToken(), "if condition is not a comparison")
		}
		ifScope := multiPathScope.AccessChild()
		WalkBody(w, &node.Elseifs[i].Body, ifScope)
	}

	if node.Else != nil {
		elseScope := multiPathScope.AccessChild()
		WalkBody(w, &node.Else.Body, elseScope)
	}
}

func AssignmentStmt(w *wkr.Walker, assignStmt *ast.AssignmentStmt, scope *wkr.Scope) {

	wIdents := []wkr.Value{}
	for i := range assignStmt.Identifiers {
		wIdents = append(wIdents, GetNodeValue(w, &assignStmt.Identifiers[i], scope))
	}

	for i := range assignStmt.Values {
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

	if len(assignStmt.Values) < len(assignStmt.Identifiers) {
		w.Error(assignStmt.Values[len(assignStmt.Values)-1].GetToken(), "not enough values provided in assignment")
	} else if len(assignStmt.Values) > len(assignStmt.Identifiers) {
		w.Error(assignStmt.Values[len(assignStmt.Values)-1].GetToken(), "too many values provided in assignment")
	}
}

func RepeatStmt(w *wkr.Walker, node *ast.RepeatStmt, scope *wkr.Scope) {
	repeatScope := scope.AccessChild()

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

	if !(wkr.TypeEquals(repeatType, startType) && wkr.TypeEquals(startType, skipType)) {
		w.Error(node.Start.GetToken(), fmt.Sprintf("all value types must be the same (iter:%s, start:%s, by:%s)", repeatType.ToString(), startType.ToString(), skipType.ToString()))
	} else {
		w.GetVariable(repeatScope, node.Variable.Name.Lexeme).Value = w.TypeToValue(repeatType)
	}

	WalkBody(w, &node.Body, repeatScope)
}

func WhileStmt(w *wkr.Walker, node *ast.WhileStmt, scope *wkr.Scope) {
	whileScope := scope.AccessChild()

	_ = GetNodeValue(w, &node.Condtion, scope)

	WalkBody(w, &node.Body, whileScope)
}

func ForloopStmt(w *wkr.Walker, node *ast.ForStmt, scope *wkr.Scope) {
	forScope := scope.AccessChild()

	if len(node.KeyValuePair) != 0 {
		w.GetVariable(forScope, node.KeyValuePair[0].Name.Lexeme).Value = &wkr.NumberVal{}
	}
	valType := GetNodeValue(w, &node.Iterator, scope).GetType()
	wrapper, ok := valType.(*wkr.WrapperType)
	if !ok {
		w.Error(node.Iterator.GetToken(), "iterator must be of type map or list")
	} else if len(node.KeyValuePair) == 2 {
		node.OrderedIteration = wrapper.PVT() == ast.List
		w.GetVariable(forScope, node.KeyValuePair[1].Name.Lexeme).Value = w.TypeToValue(wrapper.WrappedType)
	}

	WalkBody(w, &node.Body, forScope)
}

func TickStmt(w *wkr.Walker, node *ast.TickStmt, scope *wkr.Scope) {
	tickScope := scope.AccessChild()

	WalkBody(w, &node.Body, tickScope)
}

func MatchStmt(w *wkr.Walker, node *ast.MatchStmt, isExpr bool, scope *wkr.Scope) {
	val := GetNodeValue(w, &node.ExprToMatch, scope)
	if val.GetType().PVT() == ast.Invalid {
		w.Error(node.ExprToMatch.GetToken(), "variable is of type invalid")
	}
	valType := val.GetType()
	multiPathScope := scope.AccessChild()

	for i := range node.Cases {
		caseScope := multiPathScope.AccessChild()

		if !isExpr {
			WalkBody(w, &node.Cases[i].Body, caseScope)
		}

		if node.Cases[i].Expression.GetToken().Lexeme == "_" {
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
}

func BreakStmt(w *wkr.Walker, node *ast.BreakStmt, scope *wkr.Scope) {
	if !scope.Is(wkr.BreakAllowing) {
		w.Error(node.GetToken(), "cannot use break outside of loops")
	}
}

func ContinueStmt(w *wkr.Walker, node *ast.ContinueStmt, scope *wkr.Scope) {
	if !scope.Is(wkr.ContinueAllowing) {
		w.Error(node.GetToken(), "cannot use break outside of loops")
	}
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

	if helpers.ListsAreSame(matchExprTag.YieldValues, wkr.EmptyReturn) {
		matchExprTag.YieldValues = ret
	} else {
		errorMsg := w.ValidateReturnValues(ret, matchExprTag.YieldValues)
		if errorMsg != "" {
			errorMsg = strings.Replace(errorMsg, "return", "yield", -1)
			w.Error(node.GetToken(), errorMsg)
		}
	}

	return &ret
}

func UseStmt(w *wkr.Walker, node *ast.UseStmt, scope *wkr.Scope) {
	if scope.Parent != nil {
		w.Error(node.GetToken(), "cannot have a use statement inside a local block")
		return
	}
	envStmt := w.GetEnvStmt()
	envName := node.Path.Nameify()
	walker, found := w.Walkers[envName]

	if !found {
		w.Error(node.GetToken(), fmt.Sprintf("Environment named '%s' doesn't exist", envName))
		return
	}

	for _, v := range walker.GetEnvStmt().Requirements {
		if v == w.Environment.Path {
			w.Error(node.GetToken(), fmt.Sprintf("import cycle detected: this environment and '%s' are using each other", walker.Environment.Name))
			return
		}
	}

	success := envStmt.AddRequirement(walker.Environment.Path)

	if !success {
		w.Error(node.GetToken(), fmt.Sprintf("Environment '%s' is already used", envName))
		return
	}

	w.UsedWalkers = append(w.UsedWalkers, walker)
}
