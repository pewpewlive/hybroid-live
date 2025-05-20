package walker

import (
	"hybroid/ast"
	"hybroid/helpers"
	"hybroid/tokens"
	"strings"
)

// (w *Walker) func alker, node *ast.TypeDeclarationStmt, scope *Scope) {
// 	w.Environment.CustomTypes[node.Alias.Lexeme] = NewCustomType(node.Alias.Lexeme, TypeExpr(node.AliasedType, w.Environment))
// }

func (w *Walker) AliasDeclarationStmt(node *ast.AliasDecl, scope *Scope) {
	w.environment.AliasTypes[node.Name.Lexeme] = NewAliasType(node.Name.Lexeme, w.TypeExpr(node.Type, &w.environment.Scope, true))
}

func (w *Walker) ClassDeclarationStmt(node *ast.ClassDecl, scope *Scope) {
	if node.Constructor == nil {
		// w.Error(node.Name, "structs must be declared with a constructor")
		return
	}

	if w.TypeExists(node.Name.Lexeme) {
		// w.Error(node.Name, "a type with this name already exists")
	}

	generics := make([]*GenericType, 0)

	for _, param := range node.Constructor.Generics {
		generics = append(generics, NewGeneric(param.Name.Lexeme))
	}

	classVal := &ClassVal{
		Type:     *NewNamedType(w.environment.Name, node.Name.Lexeme, ast.Struct),
		IsLocal:  node.IsPub,
		Fields:   make(map[string]Field),
		Methods:  map[string]*VariableVal{},
		Generics: generics,
		Params:   Types{},
	}

	// DECLARATIONS
	w.DeclareClass(classVal)

	classScope := NewScope(scope, &ClassTag{Val: classVal}, SelfAllowing)

	params := make([]Type, 0)
	for _, param := range node.Constructor.Params {
		params = append(params, w.TypeExpr(param.Type, scope, true))
	}
	classVal.Params = params

	funcDeclaration := ast.MethodDecl{
		Name:     node.Constructor.Token,
		Params:   node.Constructor.Params,
		Generics: node.Constructor.Generics,
		IsPub:    true,
		Body:     node.Constructor.Body,
	}

	for i := range node.Fields {
		w.FieldDeclarationStmt(&node.Fields[i], classVal, classScope)
	}

	for i := range node.Methods {
		w.MethodDeclarationStmt(&node.Methods[i], classVal, classScope)
	}

	w.MethodDeclarationStmt(&funcDeclaration, classVal, classScope)

	// WALKING
	w.MethodDeclarationStmt(&funcDeclaration, classVal, classScope)

	for _, v := range classVal.Fields {
		if !v.Var.IsInit {
			// w.Error(node.GetToken(), "all fields need to be initialized in constructor (found '%s')", v.Var.Name)
			break
		}
	}

	for i := range node.Methods {
		w.MethodDeclarationStmt(&node.Methods[i], classVal, classScope)
	}
}

func (w *Walker) EntityDeclarationStmt(node *ast.EntityDecl, scope *Scope) {
	et := &EntityTag{}
	entityScope := NewScope(scope, et, SelfAllowing)

	if scope.Parent != nil {
		// w.Error(node.Token, "can't declare an entity inside a local block")
	}

	if w.TypeExists(node.Name.Lexeme) {
		// w.Error(node.Name, "a type with this name already exists")
	}

	entityVal := NewEntityVal(w.environment.Name, node.Name.Lexeme, node.IsPub)

	// DECLARATIONS
	for i := range node.Fields {
		w.FieldDeclarationStmt(&node.Fields[i], entityVal, entityScope)
	}

	et.EntityType = entityVal

	w.DeclareEntity(entityVal)

	for i := range node.Methods {
		w.MethodDeclarationStmt(&node.Methods[i], entityVal, entityScope)
	}

	//callbacks
	found := map[ast.EntityFunctionType][]tokens.Token{}

	if node.Destroyer == nil {
		// w.Error(node.Token, "entities must be declared with a destroyer")
	} else {
		w.EntityFunctionDeclarationStmt(node.Destroyer, entityVal, entityScope)
	}

	if node.Spawner == nil {
		// w.Error(node.Token, "entities must be declared with a spawner")
	} else {
		w.EntityFunctionDeclarationStmt(node.Spawner, entityVal, entityScope)
	}

	// WALKING
	if node.Destroyer != nil {
		w.EntityFunctionDeclarationStmt(node.Destroyer, entityVal, entityScope)
	}

	for i := range node.Methods {
		w.MethodDeclarationStmt(&node.Methods[i], entityVal, entityScope)
	}

	for i := range node.Callbacks {
		found[node.Callbacks[i].Type] = append(found[node.Callbacks[i].Type], node.Callbacks[i].Token)
		w.EntityFunctionDeclarationStmt(node.Callbacks[i], entityVal, entityScope)
	}

	for k := range found {
		if len(found[k]) > 1 {
			// for i := range found[k] {
			// 	w.Error(found[k][i], fmt.Sprintf("multiple instances of the same entity function is not allowed (%s)", k))
			// }
		}
	}
}

func (w *Walker) EntityFunctionDeclarationStmt(node *ast.EntityFunctionDecl, entityVal *EntityVal, scope *Scope) {
	generics := make([]*GenericType, 0)

	for _, param := range node.Generics {
		generics = append(generics, NewGeneric(param.Name.Lexeme))
	}

	ret := w.GetReturns(node.Return, scope)

	ft := &FuncTag{
		Generics:    generics,
		ReturnTypes: ret,
		Returns:     make([]bool, len(ret)),
	}
	fnScope := NewScope(scope, ft, ReturnAllowing)
	params := w.WalkParams(node.Params, scope, func(name tokens.Token, value Value) {
		w.DeclareVariable(fnScope, &VariableVal{
			Name:    name.Lexeme,
			Value:   value,
			IsLocal: true,
			Token:   node.GetToken(),
		}, name)
	})

	funcSign := NewFuncSignature().
		WithParams(params...).
		WithReturns(ret...)

	w.context.Clear()

	if node.Type != ast.Destroy || entityVal.DestroyParams != nil {
		w.WalkBody(&node.Body, ft, fnScope)

		if !ft.GetIfExits(Return) && len(ft.ReturnTypes) != 0 {
			// w.Error(node.GetToken(), "not all code paths return")
		}
	}

	switch node.Type {
	case ast.Spawn:
		for _, v := range entityVal.Fields {
			if !v.Var.IsInit {
				// w.Error(node.GetToken(), "all fields need to be initialized in spawner")
				break
			}
		}
		if len(params) < 2 || !(params[0].GetType() == Fixed && params[1].GetType() == Fixed) {
			// w.Error(node.Token, "first two parameters of %s must be of fixed type", node.Type)
		}
		if len(ret) != 0 {
			// w.Error(node.Token, "spawner must have no return types")
		}
		entityVal.SpawnParams = params
	case ast.Destroy:
		entityVal.DestroyParams = params
		entityVal.DestroyGenerics = generics
	case ast.WallCollision:
		if !funcSign.Equals(WallCollisionSign) {
			// w.Error(node.Token, "wrong function signature: expected %s", WallCollisionSign.ToString())
		}
	case ast.PlayerCollision:
		if !funcSign.Equals(PlayerCollisionSign) {
			// w.Error(node.Token, "wrong function signature: expected %s", PlayerCollisionSign.ToString())
		}
	case ast.WeaponCollision:
		if !funcSign.Equals(WeaponCollisionSign) {
			// w.Error(node.Token, "wrong function signature: expected %s", WeaponCollisionSign.ToString())
		}
	}
}

func (w *Walker) EnumDeclarationStmt(node *ast.EnumDecl, scope *Scope) {
	enumVal := &EnumVal{
		Type:   NewEnumType(scope.Environment.Name, node.Name.Lexeme),
		Fields: make(map[string]*VariableVal),
	}

	if len(node.Fields) == 0 {
		// w.Error(node.GetToken(), "can't declare an enum with no fields")
	}
	for _, v := range node.Fields {
		variable := &VariableVal{
			Name:    v.Name.Lexeme,
			Value:   &EnumFieldVal{Type: enumVal.Type},
			IsLocal: node.IsPub,
			IsConst: true,
		}
		enumVal.AddField(variable)
	}

	enumVar := &VariableVal{
		Name:    enumVal.Type.Name,
		Value:   enumVal,
		IsLocal: node.IsPub,
		IsConst: true,
	}

	if w.TypeExists(enumVar.Name) {
		// w.Error(node.Name, "a type with this name already exists")
	}

	w.DeclareVariable(scope, enumVar, node.GetToken())
}

func (w *Walker) FieldDeclarationStmt(node *ast.FieldDecl, container FieldContainer, scope *Scope) {
	varDecl := ast.VariableDecl{
		Identifiers: node.Identifiers,
		Type:        node.Type,
		Expressions: node.Values,
		IsPub:       false,
		Token:       node.Token,
	}

	variables := w.VariableDeclarationStmt(&varDecl, scope)
	node.Values = varDecl.Expressions
	for i := range variables {
		container.AddField(variables[i])
	}
}

func (w *Walker) MethodDeclarationStmt(node *ast.MethodDecl, container MethodContainer, scope *Scope) {
	if variable, found := container.ContainsMethod(node.Name.Lexeme); found {
		fn := variable.Value.(*FunctionVal)
		fnTag := &FuncTag{
			Returns:     make([]bool, 0),
			ReturnTypes: fn.Returns,

			Generics: fn.Generics,
		}

		fnScope := NewScope(scope, fnTag, ReturnAllowing)

		for i, param := range node.Params {
			w.DeclareVariable(fnScope, &VariableVal{Name: param.Name.Lexeme, Value: w.TypeToValue(fn.Params[i]), IsLocal: true, IsInit: true}, param.Name)
		}

		w.WalkBody(&node.Body, fnTag, fnScope)
	} else {
		funcExpr := ast.FunctionDecl{
			Name:     node.Name,
			Return:   node.Return,
			Params:   node.Params,
			Generics: node.Generics,
			Body:     node.Body,
			IsPub:    false,
		}

		variable := w.FunctionDeclarationStmt(&funcExpr, scope, Method)
		container.AddMethod(variable)
	}
}

func (w *Walker) FunctionDeclarationStmt(node *ast.FunctionDecl, scope *Scope, procType ProcedureType) *VariableVal {
	if node.Name.Lexeme == "Bounce" {
		print("breakpoint")
	}
	generics := make([]*GenericType, 0)

	for _, param := range node.Generics {
		generics = append(generics, NewGeneric(param.Name.Lexeme))
	}

	funcTag := &FuncTag{Generics: generics}
	fnScope := NewScope(scope, funcTag, ReturnAllowing)

	ret := w.GetReturns(node.Return, fnScope)

	funcTag.ReturnTypes = ret

	params := make([]Type, 0)
	for i, param := range node.Params {
		params = append(params, w.TypeExpr(param.Type, fnScope, true))
		w.DeclareVariable(fnScope, &VariableVal{Name: param.Name.Lexeme, Value: w.TypeToValue(params[i]), IsLocal: true, IsInit: true}, param.Name)
	}

	variable := &VariableVal{
		Name:    node.Name.Lexeme,
		Value:   &FunctionVal{Params: params, Returns: ret, Generics: generics},
		Token:   node.GetToken(),
		IsLocal: node.IsPub,
	}
	w.DeclareVariable(scope, variable, variable.Token)

	if procType == Function {
		w.WalkBody(&node.Body, funcTag, fnScope)

		if !funcTag.GetIfExits(Return) && len(ret) != 0 {
			// w.Error(node.GetToken(), "not all code paths return")
		}
	}

	return variable
}

func (w *Walker) VariableDeclarationStmt(declaration *ast.VariableDecl, scope *Scope) []*VariableVal {
	declaredVariables := []*VariableVal{}

	types := make([]Type, 0)

	index := 0
	for i := range declaration.Expressions {
		index++
		exprValue := w.GetNodeValue(&declaration.Expressions[i], scope)
		// if declaration.Values[i].GetType() == ast.SelfExpression {
		// 	// w.Error(declaration.Values[i].GetToken(), "cannot assign self to a variable")
		// }
		if _typs, ok := exprValue.(*Types); ok {
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
		// w.Error(declaration.Token, "too many values given in variable declaration")
	} else if identsLength > trueValuesLength {
		filledAll := true
		for i := index; i < identsLength; i++ {
			if declaration.Type != nil {
				typ := w.TypeExpr(declaration.Type, scope, true)
				val := w.TypeToValue(typ)
				_default := val.GetDefault()
				if _default.Value == "nil" {
					types = append(types, nil)
				} else {
					types = append(types, typ)
				}

				declaration.Expressions = append(declaration.Expressions, _default)
			} else {
				// w.Error(declaration.Identifiers[i], "variable is uninitialized and no explicit type was given")
				filledAll = false
			}
		}
		if !filledAll {
			// w.Error(declaration.Token, "too few values given in variable declaration")
			return []*VariableVal{}
		}
	}

	if declaration.IsPub && scope.Parent != nil {
		// w.Error(declaration.Token, "cannot declare a global variable inside a local block")
	}
	if declaration.Token.Type == tokens.Const && scope.Parent != nil {
		// w.Error(declaration.Token, "cannot declare a global constant inside a local block")
	}

	for i, ident := range declaration.Identifiers {
		if ident.Name.Lexeme == "_" {
			continue
		}

		valType := types[i]
		if w.environment.Type == ast.MeshEnv && ident.Name.Lexeme == "meshes" {
			if len(declaration.Identifiers) > 1 {
				// w.Error(ident, "'meshes' variable cannot be declared with other variables")
			} else if !declaration.IsPub {
				// w.Error(ident, "'meshes' has to be global")
			}

			if !TypeEquals(valType, MeshesValueType) {
				// w.Error(ident, "'meshes' needs to be of type %s", MeshesValueType.ToString())
			}
		}

		if w.environment.Type == ast.SoundEnv && ident.Name.Lexeme == "sounds" {
			if len(declaration.Identifiers) > 1 {
				// w.Error(ident, "'sounds' variable cannot be declared with other variables")
			} else if !declaration.IsPub {
				// w.Error(ident, "'sounds' has to be global")
			}

			if !TypeEquals(valType, SoundsValueType) {
				// w.Error(ident, "'sounds' needs to be of type %s", SoundsValueType.ToString())
			}
		}

		if declaration.Type == nil && types[i] == nil {
			// w.Error(declaration.Token, "Must provide an explicit type for an uninitialized variable")
		}
		if declaration.Type != nil && types[i] != nil {
			explicitType := w.TypeExpr(declaration.Type, scope, false)
			if !TypeEquals(valType, explicitType) {
				// w.Error(declaration.Identifiers[i], "Given value is %s, but explicit type is %s", valType.ToString(), explicitType.ToString())
			}
		} else if types[i] != nil && valType.PVT() == ast.Invalid {
			// w.Error(declaration.Expressions[i].GetToken(), "value is invalid")
		}

		var val Value
		if types[i] == nil {
			if declaration.Type == nil {
				val = &Invalid{}
			} else {
				val = w.TypeToValue(w.TypeExpr(declaration.Type, scope, false))
			}
		} else {
			if types[i].GetType() == Wrapper && types[i].(*WrapperType).WrappedType.PVT() == ast.Object {
				if declaration.Type == nil {
					// w.Error(declaration.Identifiers[i], "cannot infer the wrapped type of the map/list")
				} else {
					val = w.TypeToValue(w.TypeExpr(declaration.Type, scope, false))
				}
			} else {
				val = w.TypeToValue(types[i])
			}
		}

		variable := &VariableVal{
			Value:   val,
			Name:    ident.Name.Lexeme,
			IsLocal: declaration.IsPub,
			IsConst: declaration.IsConst,
			IsInit:  types[i] != nil,
			Token:   ident.Name,
		}
		declaredVariables = append(declaredVariables, variable)
		w.DeclareVariable(scope, variable, ident.Name)
	}

	return declaredVariables
}

func (w *Walker) DeclareConversion(scope *Scope) {
	if len(w.context.Conversions) == 1 {
		conv := w.context.Conversions[0]
		w.DeclareVariable(scope, &VariableVal{
			Name:   conv.Name.Lexeme,
			Value:  conv.Entity,
			IsInit: true,
		}, conv.Name)
	}
	w.context.Conversions = make([]EntityConversion, 0)
}

func (w *Walker) IfStmt(node *ast.IfStmt, scope *Scope) {
	length := len(node.Elseifs) + 2
	mpt := NewMultiPathTag(length, scope.Attributes...)
	multiPathScope := NewScope(scope, mpt)
	ifScope := NewScope(multiPathScope, &UntaggedTag{})

	w.context.Conversions = make([]EntityConversion, 0)

	boolExpr := w.GetNodeValue(&node.BoolExpr, scope)
	if boolExpr.GetType().PVT() != ast.Bool {
		// w.Error(node.BoolExpr.GetToken(), "if condition is not a comparison")
	}
	w.DeclareConversion(ifScope)
	w.WalkBody(&node.Body, mpt, ifScope)

	for i := range node.Elseifs {
		boolExpr := w.GetNodeValue(&node.Elseifs[i].BoolExpr, scope)
		if boolExpr.GetType().PVT() != ast.Bool {
			// w.Error(node.Elseifs[i].BoolExpr.GetToken(), "if condition is not a comparison")
		}
		ifScope := NewScope(multiPathScope, &UntaggedTag{})
		w.DeclareConversion(ifScope)
		w.WalkBody(&node.Elseifs[i].Body, mpt, ifScope)
	}

	if node.Else != nil {
		elseScope := NewScope(multiPathScope, &UntaggedTag{})
		w.WalkBody(&node.Else.Body, mpt, elseScope)
	}
}

func (w *Walker) AssignmentStmt(assignStmt *ast.AssignmentStmt, scope *Scope) {
	type Value2 struct {
		Value
		Index int
	}

	values := []Value2{}
	index := 1
	for i := range assignStmt.Values { // function()
		index++
		exprValue := w.GetNodeValue(&assignStmt.Values[i], scope)
		// if assignStmt.Values[i].GetType() == ast.SelfExpression {
		// 	// w.Error(assignStmt.Values[i].GetToken(), "cannot assign self to a variable")
		// }
		if types, ok := exprValue.(*Types); ok {
			for j := range *types {
				values = append(values, Value2{w.TypeToValue((*types)[j]), i})
			}
		} else {
			values = append(values, Value2{exprValue, i})
		}
	}

	variablesLength := len(assignStmt.Identifiers)
	valuesLength := len(values)
	if variablesLength < valuesLength {
		// w.Error(assignStmt.Token, "too many values given in variable declaration")
	} else if variablesLength > valuesLength {
		// w.Error(assignStmt.Token, "too few values given in variable declaration")
	}

	for i := index; i < variablesLength; i++ {
		values = append(values, Value2{&Invalid{}, i})
	}

	for i := range assignStmt.Identifiers {
		value := w.GetNodeValue(&assignStmt.Identifiers[i], scope)
		variable, ok := value.(*VariableVal)
		if !ok {
			variable = &VariableVal{
				Name:   "",
				Value:  value,
				IsInit: true,
			}
		}
		if variable.IsConst {
			// variableToken := assignStmt.Identifiers[i].GetToken()
			// w.Error(variableToken, "cannot modify '%s' because it is const", variableToken.Lexeme)
			continue
		}

		variableType := variable.GetType()
		if !variable.IsInit {
			variable.IsInit = true
		}

		valType := values[i].GetType()

		if !TypeEquals(variableType, valType) {
			// variableName := assignStmt.Identifiers[i].GetToken().Lexeme
			// w.Error(assignStmt.Values[values[i].Index].GetToken(), "mismatched types: '%s' is of type %s but a value of %s was given to it", variableName, variableType.ToString(), valType.ToString())
		}

		if vr, ok := values[i].Value.(*VariableVal); ok {
			values[i] = Value2{vr.Value, values[i].Index}
		}

		//variable.Value = values[i]
	}
}

func (w *Walker) RepeatStmt(node *ast.RepeatStmt, scope *Scope) {
	repeatScope := NewScope(scope, &MultiPathTag{}, BreakAllowing, ContinueAllowing)
	lt := NewMultiPathTag(1, repeatScope.Attributes...)
	repeatScope.Tag = lt

	end := w.GetNodeValue(&node.Iterator, scope)
	endType := end.GetType()
	//if !parser.IsFx(endType.PVT()) && endType.PVT() != ast.Number {
	// w.Error(node.Iterator.GetToken(), "invalid value type of iterator")
	//} else if variable, ok := end.(*VariableVal); ok {
	//if fixedpoint, ok := variable.Value.(*FixedVal); ok {
	//	endType = NewBasicType(fixedpoint.SpecificType)
	//}
	//} else {
	if fixedpoint, ok := end.(*FixedVal); ok {
		endType = NewBasicType(fixedpoint.SpecificType)
	}
	//}
	if node.Start.GetType() == ast.NA {
		node.Start = &ast.LiteralExpr{Token: node.Start.GetToken(), ValueType: endType.PVT(), Value: "1"}
	}
	start := w.GetNodeValue(&node.Start, scope)
	if node.Skip.GetType() == ast.NA {
		node.Skip = &ast.LiteralExpr{Token: node.Skip.GetToken(), ValueType: endType.PVT(), Value: "1"}
	}
	skip := w.GetNodeValue(&node.Skip, scope)

	repeatType := end.GetType()
	startType := start.GetType()
	skipType := skip.GetType()

	if !(TypeEquals(repeatType, startType) && TypeEquals(startType, skipType)) {
		// w.Error(node.Start.GetToken(), fmt.Sprintf("all value types must be the same (iter:%s, start:%s, by:%s)", repeatType.ToString(), startType.ToString(), skipType.ToString()))
	}

	if node.Variable != nil {
		w.DeclareVariable(repeatScope, &VariableVal{Name: node.Variable.Name.Lexeme, Value: w.TypeToValue(repeatType), IsLocal: true}, node.Variable.Name)
	}

	w.WalkBody(&node.Body, lt, repeatScope)
}

func (w *Walker) WhileStmt(node *ast.WhileStmt, scope *Scope) {
	whileScope := NewScope(scope, &MultiPathTag{}, BreakAllowing, ContinueAllowing)
	lt := NewMultiPathTag(1, whileScope.Attributes...)
	whileScope.Tag = lt

	_ = w.GetNodeValue(&node.Condition, scope)

	w.WalkBody(&node.Body, lt, whileScope)
}

func (w *Walker) ForloopStmt(node *ast.ForStmt, scope *Scope) {
	forScope := NewScope(scope, &MultiPathTag{}, BreakAllowing, ContinueAllowing)
	lt := NewMultiPathTag(1, forScope.Attributes...)
	forScope.Tag = lt

	w.DeclareVariable(forScope,
		&VariableVal{Name: node.First.Name.Lexeme, Value: &NumberVal{}},
		node.First.Name)

	valType := w.GetNodeValue(&node.Iterator, scope).GetType()
	wrapper, ok := valType.(*WrapperType)
	if !ok {
		// w.Error(node.Iterator.GetToken(), "iterator must be of type map or list")
	} else if node.Second != nil {
		node.OrderedIteration = wrapper.PVT() == ast.List
		w.DeclareVariable(forScope,
			&VariableVal{Name: node.Second.Name.Lexeme, Value: w.TypeToValue(wrapper.WrappedType)},
			node.Second.Name)
	}

	w.WalkBody(&node.Body, lt, forScope)
}

func (w *Walker) TickStmt(node *ast.TickStmt, scope *Scope) {
	funcTag := &FuncTag{ReturnTypes: EmptyReturn}
	tickScope := NewScope(scope, funcTag, ReturnAllowing)

	if node.Variable != nil {
		w.DeclareVariable(tickScope, &VariableVal{Name: node.Variable.Name.Lexeme, Value: &NumberVal{}}, node.Token)
	}

	w.WalkBody(&node.Body, funcTag, tickScope)
}

func (w *Walker) MatchStmt(node *ast.MatchStmt, isExpr bool, scope *Scope) {
	val := w.GetNodeValue(&node.ExprToMatch, scope)
	if val.GetType().PVT() == ast.Invalid {
		// w.Error(node.ExprToMatch.GetToken(), "variable is of type invalid")
	}
	valType := val.GetType()
	casesLength := len(node.Cases) + 1
	if node.HasDefault {
		casesLength--
	}
	mpt := NewMultiPathTag(casesLength, scope.Attributes...)
	multiPathScope := NewScope(scope, mpt)

	for i := range node.Cases {
		caseScope := NewScope(multiPathScope, &UntaggedTag{})

		if !isExpr {
			w.WalkBody(&node.Cases[i].Body, mpt, caseScope)
		}

		if node.Cases[i].Expression.GetToken().Lexeme == "else" {
			continue
		}

		caseValType := w.GetNodeValue(&node.Cases[i].Expression, scope).GetType()
		if !TypeEquals(valType, caseValType) {
			// w.Error(
			// node.Cases[i].Expression.GetToken(),
			// fmt.Sprintf("mismatched types: arm expression (%s) and match expression (%s)",
			// 	caseValType.ToString(),
			// 	valType.ToString()))
		}
	}
}

func (w *Walker) BreakStmt(node *ast.BreakStmt, scope *Scope) {
	if !scope.Is(BreakAllowing) {
		// w.Error(node.GetToken(), "cannot use break outside of loops")
	}

	if returnable := scope.ResolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, Break)
		(*returnable).SetExit(true, All)
	}
}

func (w *Walker) ContinueStmt(node *ast.ContinueStmt, scope *Scope) {
	if !scope.Is(ContinueAllowing) {
		// w.Error(node.GetToken(), "cannot use break outside of loops")
	}

	if returnable := scope.ResolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, Continue)
		(*returnable).SetExit(true, All)
	}
}

func (w *Walker) ReturnStmt(node *ast.ReturnStmt, scope *Scope) *Types {
	if !scope.Is(ReturnAllowing) {
		// w.Error(node.GetToken(), "can't have a return statement outside of a function or method")
	}

	ret := EmptyReturn
	for i := range node.Args {
		val := w.GetNodeValue(&node.Args[i], scope) // we need to check waht happens here
		valType := val.GetType()
		if types, ok := val.(*Types); ok {
			ret = append(ret, *types...)
		} else {
			ret = append(ret, valType)
		}
	}
	sc, _, funcTag := ResolveTagScope[*FuncTag](scope)
	if sc == nil {
		return &ret
	}

	errorMsg := w.ValidateReturnValues(ret, (*funcTag).ReturnTypes) // wait
	if errorMsg != "" {
		// w.Error(node.GetToken(), errorMsg)
	}

	if returnable := scope.ResolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, Return)
		(*returnable).SetExit(true, All)
	}

	return &ret
}

func (w *Walker) YieldStmt(node *ast.YieldStmt, scope *Scope) *Types {
	if !scope.Is(YieldAllowing) {
		// w.Error(node.GetToken(), "cannot use yield outside of statement expressions") // wut
	}

	ret := EmptyReturn
	for i := range node.Args {
		val := w.GetNodeValue(&node.Args[i], scope)
		valType := val.GetType()
		if types, ok := val.(*Types); ok {
			ret = append(ret, *types...)
		} else {
			ret = append(ret, valType)
		}
	}

	sc, _, matchExprT := ResolveTagScope[*MatchExprTag](scope)

	if sc == nil {
		return &ret
	}

	matchExprTag := *matchExprT

	if helpers.ListsAreSame(matchExprTag.YieldValues, EmptyReturn) {
		matchExprTag.YieldValues = ret
	} else {
		errorMsg := w.ValidateReturnValues(ret, matchExprTag.YieldValues)
		if errorMsg != "" {
			errorMsg = strings.Replace(errorMsg, "return", "yield", -1)
			// w.Error(node.GetToken(), errorMsg)
		}
	}

	if returnable := scope.ResolveReturnable(); returnable != nil {
		(*returnable).SetExit(true, Yield)
		(*returnable).SetExit(true, All)
	}

	return &ret
}

func (w *Walker) UseStmt(node *ast.UseStmt, scope *Scope) {
	if scope.Parent != nil {
		// w.Error(node.GetToken(), "cannot have a use statement inside a local block")
		return
	}

	path := node.Path.Path.Lexeme

	switch path {
	case "Pewpew":
		w.environment.UsedLibraries[Pewpew] = true
		return
	}

	switch path {
	case "Pewpew":
		if w.environment.Type != ast.LevelEnv {
			// w.Error(node.GetToken(), "cannot use the pewpew library in a non-level environment")
		}
		w.environment.UsedLibraries[Pewpew] = true
		return
	case "Fmath":
		if w.environment.Type != ast.LevelEnv {
			// w.Error(node.GetToken(), "cannot use the fmath library in a non-level environment")
		}
		w.environment.UsedLibraries[Fmath] = true
		return
	case "Math":
		w.environment.UsedLibraries[Math] = true
		return
	case "String":
		w.environment.UsedLibraries[String] = true
		return
	case "Table":
		w.environment.UsedLibraries[Table] = true
		return
	}

	envStmt := w.environment.EnvStmt
	envName := node.Path.Path.Lexeme
	walker, found := w.walkers[envName]

	if !found {
		// w.Error(node.GetToken(), fmt.Sprintf("Environment named '%s' doesn't exist", envName))
		return
	}

	if walker.environment.luaPath == "/dynamic/level.lua" {
		w.environment.UsedWalkers = append(w.environment.UsedWalkers, walker)
		return
	}

	for _, v := range walker.environment.EnvStmt.Requirements {
		if v == w.environment.luaPath {
			// w.Error(node.GetToken(), fmt.Sprintf("import cycle detected: this environment and '%s' are using each other", walker.Environment.Name))
			return
		}
	}

	success := envStmt.AddRequirement(walker.environment.luaPath)

	if !success {
		// w.Error(node.GetToken(), fmt.Sprintf("Environment '%s' is already used", envName))
		return
	}

	w.environment.UsedWalkers = append(w.environment.UsedWalkers, walker)
}

func (w *Walker) DestroyStmt(node *ast.DestroyStmt, scope *Scope) {
	val := w.GetNodeValue(&node.Identifier, scope)
	valType := val.GetType()

	if valType.PVT() == ast.Invalid {
		// w.Error(node.Identifier.GetToken(), "invalid variable given in destroy expression")
		return
	} else if valType.PVT() != ast.Entity {
		// w.Error(node.Identifier.GetToken(), "variable given in destroy statement is not an entity")
		return
	}

	if variable, ok := val.(*VariableVal); ok {
		val = variable.Value
	}

	entityVal := val.(*EntityVal)

	node.EnvName = entityVal.Type.EnvName
	node.EntityName = entityVal.Type.Name

	args := make([]Type, 0)
	for i := range node.Args {
		args = append(args, w.GetNodeValue(&node.Args[i], scope).GetType())
	}

	suppliedGenerics := w.GetGenerics(node, node.Generics, entityVal.DestroyGenerics, scope)

	w.ValidateArguments(suppliedGenerics, args, entityVal.DestroyParams, node.Token)
}
