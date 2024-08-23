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

func AliasDeclarationStmt(w *wkr.Walker, node *ast.AliasDeclarationStmt, scope *wkr.Scope) {
	w.Environment.AliasTypes[node.Alias.Lexeme] = wkr.NewAliasType(node.Alias.Lexeme, TypeExpr(w, node.AliasedType, &w.Environment.Scope, true))
}

func ClassDeclarationStmt(w *wkr.Walker, node *ast.ClassDeclarationStmt, scope *wkr.Scope) {
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

	classVal := &wkr.ClassVal{
		Type:     *wkr.NewNamedType(w.Environment.Name, node.Name.Lexeme, ast.Struct),
		IsLocal:  node.IsLocal,
		Fields:   make(map[string]wkr.Field),
		Methods:  map[string]*wkr.VariableVal{},
		Generics: generics,
		Params:   wkr.Types{},
	}

	// DECLARATIONS
	w.DeclareClass(classVal)

	classScope := wkr.NewScope(scope, &wkr.ClassTag{Val: classVal}, wkr.SelfAllowing)

	params := make([]wkr.Type, 0)
	for _, param := range node.Constructor.Params {
		params = append(params, TypeExpr(w, param.Type, scope, true))
	}
	classVal.Params = params

	funcDeclaration := ast.MethodDeclarationStmt{
		Name:     node.Constructor.Token,
		Params:   node.Constructor.Params,
		Return:   node.Constructor.Return,
		Generics: node.Constructor.Generics,
		IsLocal:  true,
		Body:     node.Constructor.Body,
	}

	for i := range node.Fields {
		FieldDeclarationStmt(w, &node.Fields[i], classVal, classScope)
	}

	for i := range node.Methods {
		MethodDeclarationStmt(w, &node.Methods[i], classVal, classScope)
	}

	MethodDeclarationStmt(w, &funcDeclaration, classVal, classScope)

	// WALKING
	MethodDeclarationStmt(w, &funcDeclaration, classVal, classScope)

	for _, v := range classVal.Fields {
		if !v.Var.IsInit {
			w.Error(node.GetToken(), "all fields need to be initialized in constructor (found '%s')", v.Var.Name)
			break
		}
	}

	for i := range node.Methods {
		MethodDeclarationStmt(w, &node.Methods[i], classVal, classScope)
	}
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

	// DECLARATIONS
	for i := range node.Fields {
		FieldDeclarationStmt(w, &node.Fields[i], entityVal, entityScope)
	}

	et.EntityType = entityVal

	w.DeclareEntity(entityVal)

	for i := range node.Methods {
		MethodDeclarationStmt(w, &node.Methods[i], entityVal, entityScope)
	}

	//callbacks
	found := map[ast.EntityFunctionType][]lexer.Token{}

	if node.Destroyer == nil {
		w.Error(node.Token, "entities must be declared with a destroyer")
	} else {
		EntityFunctionDeclarationStmt(w, node.Destroyer, entityVal, entityScope)
	}

	if node.Spawner == nil {
		w.Error(node.Token, "entities must be declared with a spawner")
	} else {
		EntityFunctionDeclarationStmt(w, node.Spawner, entityVal, entityScope)
	}

	// WALKING
	if node.Destroyer != nil {
		EntityFunctionDeclarationStmt(w, node.Destroyer, entityVal, entityScope)
	}

	for i := range node.Methods {
		MethodDeclarationStmt(w, &node.Methods[i], entityVal, entityScope)
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
}

func EntityFunctionDeclarationStmt(w *wkr.Walker, node *ast.EntityFunctionDeclarationStmt, entityVal *wkr.EntityVal, scope *wkr.Scope) {
	generics := make([]*wkr.GenericType, 0)

	for _, param := range node.Generics {
		generics = append(generics, wkr.NewGeneric(param.Name.Lexeme))
	}

	ret := []wkr.Type{}

	for i := range node.Returns {
		ret = append(ret, TypeExpr(w, node.Returns[i], scope, true))
	}

	ft := &wkr.FuncTag{
		Generics:    generics,
		ReturnTypes: ret,
		Returns:     make([]bool, len(ret)),
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

	funcSign := wkr.NewFuncSignature().
		WithParams(params...).
		WithReturns(ret...)

	w.Context.Clear()

	if node.Type != ast.Destroy || entityVal.DestroyParams != nil {
		WalkBody(w, &node.Body, ft, fnScope)

		if !ft.GetIfExits(wkr.Return) && len(ft.ReturnTypes) != 0 {
			w.Error(node.GetToken(), "not all code paths return")
		}
	}

	switch node.Type {
	case ast.Spawn:
		for _, v := range entityVal.Fields {
			if !v.Var.IsInit {
				w.Error(node.GetToken(), "all fields need to be initialized in spawner")
				break
			}
		}
		if len(params) < 2 || !(params[0].GetType() == wkr.Fixed && params[1].GetType() == wkr.Fixed) {
			w.Error(node.Token, "first two parameters of %s must be of fixed type", node.Type)
		}
		if len(ret) != 0 {
			w.Error(node.Token, "spawner must have no return types")
		}
		entityVal.SpawnParams = params
	case ast.Destroy:
		entityVal.DestroyParams = params
		entityVal.DestroyGenerics = generics
	case ast.WallCollision:
		if !funcSign.Equals(wkr.WallCollisionSign) {
			w.Error(node.Token, "wrong function signature: expected %s", wkr.WallCollisionSign.ToString())
		}
	case ast.PlayerCollision:
		if !funcSign.Equals(wkr.PlayerCollisionSign) {
			w.Error(node.Token, "wrong function signature: expected %s", wkr.PlayerCollisionSign.ToString())
		}
	case ast.WeaponCollision:
		if !funcSign.Equals(wkr.WeaponCollisionSign) {
			w.Error(node.Token, "wrong function signature: expected %s", wkr.WeaponCollisionSign.ToString())
		}
	}
}

func EnumDeclarationStmt(w *wkr.Walker, node *ast.EnumDeclarationStmt, scope *wkr.Scope) {
	enumVal := &wkr.EnumVal{
		Type:   wkr.NewEnumType(scope.Environment.Name, node.Name.Lexeme),
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
		Type:        node.Type,
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
	if variable, found := container.ContainsMethod(node.Name.Lexeme); found {
		fn := variable.Value.(*wkr.FunctionVal)
		fnTag := &wkr.FuncTag{
			Returns:     make([]bool, 0),
			ReturnTypes: fn.Returns,

			Generics: fn.Generics,
		}

		fnScope := wkr.NewScope(scope, fnTag, wkr.ReturnAllowing)

		WalkBody(w, &node.Body, fnTag, fnScope)
	} else {
		funcExpr := ast.FunctionDeclarationStmt{
			Name:          node.Name,
			Return:        node.Return,
			Params:        node.Params,
			GenericParams: node.Generics,
			Body:          node.Body,
			IsLocal:       true,
		}

		variable := FunctionDeclarationStmt(w, &funcExpr, scope, wkr.Method)
		container.AddMethod(variable)
	}
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
		if procType == wkr.Function {
			w.DeclareVariable(fnScope, &wkr.VariableVal{Name: param.Name.Lexeme, Value: w.TypeToValue(params[i]), IsLocal: true, IsInit: true}, param.Name)
		}
	}

	variable := &wkr.VariableVal{
		Name:  node.Name.Lexeme,
		Value: &wkr.FunctionVal{Params: params, Returns: ret, Generics: generics},
		Token: node.GetToken(),
		IsLocal: node.IsLocal,
	}
	w.DeclareVariable(scope, variable, variable.Token)

	if procType == wkr.Function {
		WalkBody(w, &node.Body, funcTag, fnScope)

		if !funcTag.GetIfExits(wkr.Return) && len(ret) != 0 {
			w.Error(node.GetToken(), "not all code paths return")
		}
	}

	return variable
}

func VariableDeclarationStmt(w *wkr.Walker, declaration *ast.VariableDeclarationStmt, scope *wkr.Scope) []*wkr.VariableVal {
	declaredVariables := []*wkr.VariableVal{}

	types := make([]wkr.Type, 0)

	index := 0
	for i := range declaration.Values {
		index++
		exprValue := GetNodeValue(w, &declaration.Values[i], scope)
		// if declaration.Values[i].GetType() == ast.SelfExpression {
		// 	w.Error(declaration.Values[i].GetToken(), "cannot assign self to a variable")
		// }
		if _typs, ok := exprValue.(*wkr.Types); ok {
			for _, v := range *_typs {
				types = append(types, v)
			}
		} else {
			types = append(types, exprValue.GetType())
		}
	}

	identsLength := len(declaration.Identifiers)
	trueValuesLength := len(types)
	if identsLength < trueValuesLength {
		w.Error(declaration.Token, "too many values given in variable declaration")
	} else if identsLength > trueValuesLength {
		filledAll := true
		for i := index; i < identsLength; i++ {
			if declaration.Type != nil {
				typ := TypeExpr(w, declaration.Type, scope, true)
				val := w.TypeToValue(typ)
				_default := val.GetDefault()
				if _default.Value == "nil" {
					types = append(types, nil)
				} else {
					types = append(types, typ)
				}

				declaration.Values = append(declaration.Values, _default)
			} else {
				w.Error(declaration.Identifiers[i], "variable is uninitialized and no explicit type was given")
				filledAll = false
			}
		}
		if !filledAll {
			w.Error(declaration.Token, "too few values given in variable declaration")
			return []*wkr.VariableVal{}
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

		valType := types[i]
		if w.Environment.Type == ast.MeshEnv && ident.Lexeme == "meshes" {
			if len(declaration.Identifiers) > 1 {
				w.Error(ident, "'meshes' variable cannot be declared with other variables")
			} else if declaration.IsLocal {
				w.Error(ident, "'meshes' has to be global")
			}

			if !wkr.TypeEquals(valType, wkr.MeshesValueType) {
				w.Error(ident, "'meshes' needs to be of type %s", wkr.MeshesValueType.ToString())
			}
		}

		if w.Environment.Type == ast.SoundEnv && ident.Lexeme == "sounds" {
			if len(declaration.Identifiers) > 1 {
				w.Error(ident, "'sounds' variable cannot be declared with other variables")
			} else if declaration.IsLocal {
				w.Error(ident, "'sounds' has to be global")
			}

			if !wkr.TypeEquals(valType, wkr.SoundsValueType) {
				w.Error(ident, "'sounds' needs to be of type %s", wkr.SoundsValueType.ToString())
			}
		}

		if declaration.Type == nil && types[i] == nil {
			w.Error(declaration.Token, "Must provide an explicit type for an uninitialized variable")
		}
		if declaration.Type != nil && types[i] != nil {
			explicitType := TypeExpr(w, declaration.Type, scope, false)
			if !wkr.TypeEquals(valType, explicitType) {
				w.Error(declaration.Identifiers[i], "Given value is %s, but explicit type is %s", valType.ToString(), explicitType.ToString())
			}
		} else if types[i] != nil && valType.PVT() == ast.Invalid {
			w.Error(declaration.Values[i].GetToken(), "value is invalid")
		}

		var val wkr.Value
		if types[i] == nil {
			if declaration.Type == nil {
				val = &wkr.Invalid{}
			} else {
				val = w.TypeToValue(TypeExpr(w, declaration.Type, scope, false))
			}
		} else {
			if types[i].GetType() == wkr.Wrapper && types[i].(*wkr.WrapperType).WrappedType.PVT() == ast.Object {
				if declaration.Type == nil  {
					w.Error(declaration.Identifiers[i], "cannot infer the wrapped type of the map/list")
				}else {
					val = w.TypeToValue(TypeExpr(w, declaration.Type, scope, false))
				}
			}else {
				val = w.TypeToValue(types[i])
			}
		}

		variable := &wkr.VariableVal{
			Value:   val,
			Name:    ident.Lexeme,
			IsLocal: declaration.IsLocal,
			IsConst: declaration.IsConst,
			IsInit:  types[i] != nil,
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
	type Value struct {
		wkr.Value
		Index int
	}

	values := []Value{}
	index := 1
	for i := range assignStmt.Values { // function()
		index++
		exprValue := GetNodeValue(w, &assignStmt.Values[i], scope)
		// if assignStmt.Values[i].GetType() == ast.SelfExpression {
		// 	w.Error(assignStmt.Values[i].GetToken(), "cannot assign self to a variable")
		// }
		if types, ok := exprValue.(*wkr.Types); ok {
			for j := range *types {
				values = append(values, Value{w.TypeToValue((*types)[j]), i})
			}
		} else {
			values = append(values, Value{exprValue, i})
		}
	}

	variablesLength := len(assignStmt.Identifiers)
	valuesLength := len(values)
	if variablesLength < valuesLength {
		w.Error(assignStmt.Token, "too many values given in variable declaration")
	} else if variablesLength > valuesLength {
		w.Error(assignStmt.Token, "too few values given in variable declaration")
	}

	for i := index; i < variablesLength; i++ {
		values = append(values, Value{&wkr.Invalid{}, i})
	}

	for i := range assignStmt.Identifiers {
		value := GetNodeValue(w, &assignStmt.Identifiers[i], scope)
		variable, ok := value.(*wkr.VariableVal)
		if !ok {
			variable = &wkr.VariableVal{
				Name: "",
				Value: value,
				IsInit: true,
			}
		}
		if variable.IsConst {
			variableToken := assignStmt.Identifiers[i].GetToken()
			w.Error(variableToken, "cannot modify '%s' because it is const", variableToken.Lexeme)
			continue
		}

		variableType := variable.GetType()
		if !variable.IsInit {
			variable.IsInit = true
		}

		valType := values[i].GetType()

		if !wkr.TypeEquals(variableType, valType) {
			variableName := assignStmt.Identifiers[i].GetToken().Lexeme
			w.Error(assignStmt.Values[values[i].Index].GetToken(), "mismatched types: '%s' is of type %s but a value of %s was given to it", variableName, variableType.ToString(), valType.ToString())
		}

		if vr, ok := values[i].Value.(*wkr.VariableVal); ok {
			values[i] = Value{vr.Value, values[i].Index}
		}

		//variable.Value = values[i]
	}
}

func RepeatStmt(w *wkr.Walker, node *ast.RepeatStmt, scope *wkr.Scope) {
	repeatScope := wkr.NewScope(scope, &wkr.MultiPathTag{}, wkr.BreakAllowing, wkr.ContinueAllowing)
	lt := wkr.NewMultiPathTag(1, repeatScope.Attributes...)
	repeatScope.Tag = lt

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
	} else if node.KeyValuePair[1] != nil {
		node.OrderedIteration = wrapper.PVT() == ast.List
		w.DeclareVariable(forScope,
			&wkr.VariableVal{Name: node.KeyValuePair[1].Name.Lexeme, Value: w.TypeToValue(wrapper.WrappedType)},
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

	path := node.Path.Path.Lexeme

	switch path {
	case "Pewpew":
		w.Environment.UsedLibraries[wkr.Pewpew] = true
		return
	}

	switch path {
	case "Pewpew":
		w.Environment.UsedLibraries[wkr.Pewpew] = true
		return
	case "Fmath":
		w.Environment.UsedLibraries[wkr.Fmath] = true
		return
	case "Math":
		w.Environment.UsedLibraries[wkr.Math] = true
		return
	case "String":
		w.Environment.UsedLibraries[wkr.String] = true
		return
	case "Table":
		w.Environment.UsedLibraries[wkr.Table] = true
		return
	}

	envStmt := w.GetEnvStmt()
	envName := node.Path.Path.Lexeme
	walker, found := w.Walkers[envName]

	if !found {
		w.Error(node.GetToken(), fmt.Sprintf("Environment named '%s' doesn't exist", envName))
		return
	}

	if walker.Environment.Path == "/dynamic/level.lua" {
		w.Environment.UsedWalkers = append(w.Environment.UsedWalkers, walker)
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

	w.Environment.UsedWalkers = append(w.Environment.UsedWalkers, walker)
}

func DestroyStmt(w *wkr.Walker, node *ast.DestroyStmt, scope *wkr.Scope) {
	val := GetNodeValue(w, &node.Identifier, scope)
	valType := val.GetType()

	if valType.PVT() == ast.Invalid {
		w.Error(node.Identifier.GetToken(), "invalid variable given in destroy expression")
		return
	} else if valType.PVT() != ast.Entity {
		w.Error(node.Identifier.GetToken(), "variable given in destroy statement is not an entity")
		return
	}

	if variable, ok := val.(*wkr.VariableVal); ok {
		val = variable.Value
	}

	entityVal := val.(*wkr.EntityVal)

	node.EnvName = entityVal.Type.EnvName
	node.EntityName = entityVal.Type.Name

	args := make([]wkr.Type, 0)
	for i := range node.Args {
		args = append(args, GetNodeValue(w, &node.Args[i], scope).GetType())
	}

	suppliedGenerics := GetGenerics(w, node, node.Generics, entityVal.DestroyGenerics, scope)

	w.ValidateArguments(suppliedGenerics, args, entityVal.DestroyParams, node.Token)
}
