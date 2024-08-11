package pass2

import (
	"fmt"
	"hybroid/ast"
	"hybroid/helpers"
	"hybroid/lexer"
	"hybroid/parser"
	wkr "hybroid/walker"
	"strings"
)

// func TypeDeclarationStmt(w *wkr.Walker, node *ast.TypeDeclarationStmt, scope *wkr.Scope) {
// 	w.Environment.CustomTypes[node.Alias.Lexeme] = wkr.NewCustomType(node.Alias.Lexeme, TypeExpr(w, node.AliasedType, w.Environment))
// }

func StructDeclarationStmt(w *wkr.Walker, node *ast.StructDeclarationStmt, scope *wkr.Scope) {
	if node.Constructor == nil {
		w.Error(node.Name, "structs must be declared with a constructor")
		return
	}

	if w.TypeExists(node.Name.Lexeme) {
		w.Error(node.Name, "a type with this name already exists")
	}

	generics := make([]*wkr.GenericType, 0)
	
	for _, param := range node.Constructor.Generics {
		generics = append(generics, wkr.NewGeneric(param.Name.Lexeme))
	}

	structVal := &wkr.StructVal{
		Type:    *wkr.NewNamedType(w.Environment.Name, node.Name.Lexeme, ast.Struct),
		IsLocal: node.IsLocal,
		Fields:  make(map[string]wkr.Field),
		Methods: map[string]*wkr.VariableVal{},
		Generics: generics,
		Params:  wkr.Types{},
	}

	structScope := wkr.NewScope(scope, &wkr.StructTag{StructVal: structVal}, wkr.SelfAllowing)

	params := make([]wkr.Type, 0)
	for _, param := range node.Constructor.Params {
		params = append(params, TypeExpr(w, param.Type, scope, true))
	}
	structVal.Params = params

	w.DeclareStruct(structVal)

	funcDeclaration := ast.MethodDeclarationStmt{
		Name:    node.Constructor.Token,
		Params:  node.Constructor.Params,
		Return:  node.Constructor.Return,
		Generics: node.Constructor.Generics,
		IsLocal: true,
		Body:    node.Constructor.Body,
	}

	for i := range node.Fields {
		FieldDeclarationStmt(w, &node.Fields[i], structVal, structScope)
	}

	for i := range node.Methods {
		params := make([]wkr.Type, 0)
		for _, param := range node.Methods[i].Params {
			params = append(params, TypeExpr(w, param.Type, scope, true))
		}

		ret := wkr.EmptyReturn
		for _, typee := range node.Methods[i].Return {
			ret = append(ret, TypeExpr(w, typee, scope, true))
			//fmt.Printf("%s\n", ret.values[len(ret.values)-1].Type.ToString())
		}
		variable := &wkr.VariableVal{
			Name:    node.Methods[i].Name.Lexeme,
			Value:   &wkr.FunctionVal{Params: params, Returns: ret},
			IsLocal: node.IsLocal,
			Token:   node.Methods[i].GetToken(),
		}
		w.DeclareVariable(structScope, variable, node.Methods[i].Name)
		structVal.Methods[variable.Name] = variable
	}

	for i := range node.Methods {
		MethodDeclarationStmt(w, &node.Methods[i], structVal, structScope)
	}

	MethodDeclarationStmt(w, &funcDeclaration, structVal, structScope)
}

func EntityDeclarationStmt(w *wkr.Walker, node *ast.EntityDeclarationStmt, scope *wkr.Scope) {
	et := &wkr.EntityTag{}
	entityScope := wkr.NewScope(scope, et, wkr.SelfAllowing)

	if scope.Parent != nil {
		w.Error(node.Token, "can't declare an entity inside a local block")
	}

	if w.TypeExists(node.Name.Lexeme) {
		w.Error(node.Name, "a type with this name already exists")
	}

	entityVal := wkr.NewEntityVal(w.Environment.Name, node.Name.Lexeme, node.IsLocal)

	//fields
	for i := range node.Fields { 
		FieldDeclarationStmt(w, &node.Fields[i], entityVal, entityScope)
	}

	et.EntityType = entityVal

	w.DeclareEntity(entityVal)

	//callbacks
	found := map[ast.EntityFunctionType][]lexer.Token{}
	
	//spawn
	if node.Spawner == nil {
		w.Error(node.Token, "entities must be declared with a spawner")
	} else {
		EntityFunctionDeclarationStmt(w, node.Spawner, entityVal, entityScope)
	}
	//destroy
	if node.Destroyer == nil {
		w.Error(node.Token, "entities must be declared with a destroyer")
	} else {
		EntityFunctionDeclarationStmt(w, node.Destroyer, entityVal, entityScope)
	}

	for i := range node.Callbacks {
		found[node.Callbacks[i].Type] = append(found[node.Callbacks[i].Type], node.Callbacks[i].Token)
		EntityFunctionDeclarationStmt(w, node.Callbacks[i], entityVal, entityScope)
	}

	for k := range found {
		if len(found[k]) > 1 {
			for i := range found[k] {
				w.Error(found[k][i], fmt.Sprintf("multiple instances of the same entity function is not allowed (%s)", k))
			}
		}
	}

	//methods
	for i := range node.Methods {
		MethodDeclarationStmt(w, &node.Methods[i], entityVal, entityScope)
	}
}

func EntityFunctionDeclarationStmt(w *wkr.Walker, node *ast.EntityFunctionDeclarationStmt, entityVal *wkr.EntityVal, scope *wkr.Scope) {
	generics := make([]*wkr.GenericType, 0)
	
	for _, param := range node.Generics {
		generics = append(generics, wkr.NewGeneric(param.Name.Lexeme))
	}
	
	ft := &wkr.FuncTag{
		Generics: generics,
		ReturnTypes: wkr.EmptyReturn,
		Returns:     make([]bool, 0),
	}
	fnScope := wkr.NewScope(scope, ft, wkr.ReturnAllowing)
	params := WalkParams(w, node.Params, scope, func(name lexer.Token, value wkr.Value) {
		w.DeclareVariable(fnScope, &wkr.VariableVal{
			Name:    name.Lexeme,
			Value:   value,
			IsLocal: true,
			Token:   node.GetToken(),
		}, name)
	})

	w.Context.Clear()

	WalkBody(w, &node.Body, ft, fnScope)
	switch node.Type {
	case ast.Spawn:
		if len(params) < 2 || !(params[0].GetType() == wkr.Fixed && params[1].GetType() == wkr.Fixed) {
			w.Error(node.Token, fmt.Sprintf("first two parameters of %s must be of fixed type", node.Type))
		}
		entityVal.SpawnParams = params
	case ast.Destroy:
		entityVal.DestroyParams = params

	case ast.WallCollision:
		if len(params) < 2 || len(params) > 2 || !(params[0].GetType() == wkr.Fixed && params[1].GetType() == wkr.Fixed) {
			w.Error(node.Token, fmt.Sprintf("first two parameters of %s must be of fixed type", node.Type))
		}
	case ast.PlayerCollision:
		if len(params) < 2 || len(params) > 2 || !(params[0].PVT() == ast.Number && params[1].GetType() == wkr.RawEntity) {
			w.Error(node.Token, "first parameter must be a number (player index) and second must be an entity_id (ship id)")
		}		
	case ast.WeaponCollision:
		// TODO:
		//first need WeaponType and by extension the implementation of pewpew library
	}
} 

func EnumDeclarationStmt(w *wkr.Walker, node *ast.EnumDeclarationStmt, scope *wkr.Scope) {
	enumVal := &wkr.EnumVal{
		Type: wkr.NewEnumType(node.Name.Lexeme),
		Fields: make(map[string]*wkr.VariableVal),
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

	if w.TypeExists(enumVar.Name) {
		w.Error(node.Name, "a type with this name already exists")
	}

	w.DeclareVariable(scope, enumVar, node.GetToken())
}

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
		GenericParams: node.Generics,
		Body:    node.Body,
		IsLocal: true,
	}

	variable := FunctionDeclarationStmt(w, &funcExpr, scope, wkr.Method)
	node.Body = funcExpr.Body
	container.AddMethod(variable)
}

func FunctionDeclarationStmt(w *wkr.Walker, node *ast.FunctionDeclarationStmt, scope *wkr.Scope, procType wkr.ProcedureType) *wkr.VariableVal {
	generics := make([]*wkr.GenericType, 0)
	
	for _, param := range node.GenericParams {
		generics = append(generics, wkr.NewGeneric(param.Name.Lexeme))
	}

	funcTag := &wkr.FuncTag{Generics: generics}
	fnScope := wkr.NewScope(scope, funcTag, wkr.ReturnAllowing)

	ret := wkr.EmptyReturn 
	for _, typee := range node.Return {
		ret = append(ret, TypeExpr(w, typee, fnScope, true))
	}

	funcTag.ReturnTypes = ret

	params := make([]wkr.Type, 0)
	for i, param := range node.Params {
		params = append(params, TypeExpr(w, param.Type, fnScope, true))
		w.DeclareVariable(fnScope, &wkr.VariableVal{Name: param.Name.Lexeme, Value: w.TypeToValue(params[i]), IsLocal: true}, param.Name)
	}

	w.Context.Clear()

	variable := &wkr.VariableVal{
		Name:  node.Name.Lexeme,
		Value: &wkr.FunctionVal{Params: params, Returns: ret, Generics: generics},
		Token: node.GetToken(),
	}
	if procType == wkr.Function {
		w.DeclareVariable(scope, variable, variable.Token)
	}

	WalkBody(w, &node.Body, funcTag, fnScope)

	return variable
}


func VariableDeclarationStmt(w *wkr.Walker, declaration *ast.VariableDeclarationStmt, scope *wkr.Scope) []*wkr.VariableVal {
	declaredVariables := []*wkr.VariableVal{}

	identsLength := len(declaration.Identifiers)
	valuesLength := len(declaration.Values)
	values := make([]wkr.Value, valuesLength)

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
	trueValuesLength := len(values)

	if !declaration.IsLocal && scope.Parent != nil {
		w.Error(declaration.Token, "cannot declare a global variable inside a local block")
	}
	if declaration.Token.Type == lexer.Const && scope.Parent != nil {
		w.Error(declaration.Token, "cannot declare a global constant inside a local block")
	}

	if identsLength < trueValuesLength {
		w.Error(declaration.Token, "too many values given in variable declaration")
	}else if trueValuesLength > identsLength {
		w.Error(declaration.Token, "too few values given in variable declaration")
		return []*wkr.VariableVal{}
	}

	for i, ident := range declaration.Identifiers {
		if ident.Lexeme == "_" {
			continue
		}
		if w.Environment.Type == ast.MeshEnv && ident.Lexeme == "meshes" {
			if len(declaration.Identifiers) > 1 {
				w.Error(ident, "'meshes' variable cannot be declared with another variables")
			}else if declaration.IsLocal {
				w.Error(ident, "'meshes' has to be global")
			}
			expectedValueType := &wkr.ListVal{ValueType: wkr.NewAnonStructType(map[string]wkr.Field{
					"vertexes": wkr.NewField(0, &wkr.VariableVal{
						Name: "vertexes",
						Value: &wkr.ListVal{ValueType: wkr.NewWrapperType(wkr.NewBasicType(ast.List), wkr.NewBasicType(ast.Number))},
					}),
					"segments": wkr.NewField(1, &wkr.VariableVal{
						Name: "segments",
						Value:  &wkr.ListVal{ValueType: wkr.NewWrapperType(wkr.NewBasicType(ast.List), wkr.NewBasicType(ast.Number))},
					}),
					"colors": wkr.NewField(2, &wkr.VariableVal{
						Name: "colors",
						Value: &wkr.ListVal{ValueType: wkr.NewBasicType(ast.Number)},
					}),
				}, true),
			}
			if !wkr.TypeEquals(values[i].GetType(), expectedValueType.GetType()) {
				w.Error(ident, fmt.Sprintf("'meshes' needs to be of type %s", expectedValueType.GetType().ToString()))
			}
		}

		if declaration.Types[i] != nil {
			valueType := values[i].GetType()
			explicitType := TypeExpr(w, declaration.Types[i], scope, false)
			if valueType.PVT() == ast.Object {
				values[i] = w.TypeToValue(explicitType)
				declaration.Values = append(declaration.Values, values[i].GetDefault()) 
			} else if !wkr.TypeEquals(valueType, explicitType) {
				w.Error(declaration.Values[i].GetToken(), fmt.Sprintf("Given value is %s, but explicit type is %s", valueType.ToString(), explicitType.ToString()))
			}
		}

		if values[i].GetType().PVT() == ast.Invalid {
			w.Error(declaration.Values[i].GetToken(), "Given value is invalid")
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

	boolExpr := GetNodeValue(w, &node.BoolExpr, scope)
	if boolExpr.GetType().PVT() != ast.Bool {
		w.Error(node.BoolExpr.GetToken(), "if condition is not a comparison")
	}
	WalkBody(w, &node.Body, mpt, ifScope)

	for i := range node.Elseifs {
		boolExpr := GetNodeValue(w, &node.Elseifs[i].BoolExpr, scope)
		if boolExpr.GetType().PVT() != ast.Bool {
			w.Error(node.Elseifs[i].BoolExpr.GetToken(), "if condition is not a comparison")
		}
		ifScope := wkr.NewScope(multiPathScope, &wkr.UntaggedTag{})
		WalkBody(w, &node.Elseifs[i].Body, mpt, ifScope)
	}

	if node.Else != nil {
		elseScope := wkr.NewScope(multiPathScope, &wkr.UntaggedTag{})
		WalkBody(w, &node.Else.Body, mpt, elseScope)
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
	repeatScope := wkr.NewScope(scope, &wkr.MultiPathTag{}, wkr.BreakAllowing, wkr.ContinueAllowing)
	lt := wkr.NewMultiPathTag(1, repeatScope.Attributes...)
	repeatScope.Tag = lt

	if node.Variable != nil {
		w.DeclareVariable(repeatScope, &wkr.VariableVal{Name: node.Variable.Name.Lexeme, Value: &wkr.Invalid{}, IsLocal: true}, node.Variable.Name)
	}

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
	}

	if node.Variable != nil {
		w.DeclareVariable(repeatScope, &wkr.VariableVal{Name: node.Variable.Name.Lexeme, Value: w.TypeToValue(repeatType), IsLocal: true}, node.Variable.Name)
	}

	WalkBody(w, &node.Body, lt, repeatScope)
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
	valType := GetNodeValue(w, &node.Iterator, scope).GetType()
	wrapper, ok := valType.(*wkr.WrapperType)
	if !ok {
		w.Error(node.Iterator.GetToken(), "iterator must be of type map or list")
	} else if len(node.KeyValuePair) == 2 {
		node.OrderedIteration = wrapper.PVT() == ast.List
		w.DeclareVariable(forScope,
			&wkr.VariableVal{Name: node.KeyValuePair[1].Name.Lexeme, Value:w.TypeToValue(wrapper.WrappedType)},
			node.KeyValuePair[1].Name)
	}

	WalkBody(w, &node.Body, lt, forScope)
}

func TickStmt(w *wkr.Walker, node *ast.TickStmt, scope *wkr.Scope) {
	funcTag := &wkr.FuncTag{ReturnTypes: wkr.EmptyReturn}
	tickScope := wkr.NewScope(scope, funcTag, wkr.ReturnAllowing)

	if node.Variable != nil {
		w.DeclareVariable(tickScope, &wkr.VariableVal{Name: node.Variable.Name.Lexeme, Value: &wkr.NumberVal{}}, node.Token)
	}

	WalkBody(w, &node.Body, funcTag, tickScope)
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

	for i := range node.Cases {
		caseScope := wkr.NewScope(multiPathScope, &wkr.UntaggedTag{})

		if !isExpr {
			WalkBody(w, &node.Cases[i].Body, mpt, caseScope)
		}

		if node.Cases[i].Expression.GetToken().Lexeme == "else" {
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

	ret := wkr.EmptyReturn
	for i := range node.Args {
		val := GetNodeValue(w, &node.Args[i], scope) // we need to check waht happens here
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

	errorMsg := w.ValidateReturnValues(ret, (*funcTag).ReturnTypes) // wait
	if errorMsg != "" {
		w.Error(node.GetToken(), errorMsg)
	}

	if returnable := scope.ResolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, wkr.Return)
		(*returnable).SetExit(true, wkr.All)
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

	if returnable := scope.ResolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, wkr.Yield)
		(*returnable).SetExit(true, wkr.All)
	}

	return &ret
}

func UseStmt(w *wkr.Walker, node *ast.UseStmt, scope *wkr.Scope) {
	if scope.Parent != nil {
		w.Error(node.GetToken(), "cannot have a use statement inside a local block")
		return
	}

	if strings.ToLower(node.Path.Path.Lexeme) == "pewpew" {
		w.UsedLibraries[wkr.Pewpew] = true
		return;
	}

	envStmt := w.GetEnvStmt()
	envName := node.Path.Path.Lexeme
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
