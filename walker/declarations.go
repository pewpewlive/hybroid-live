package walker

import (
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/core"
	"hybroid/tokens"
)

// Rewrote
func (w *Walker) environmentDeclaration(node *ast.EnvironmentDecl) {
	if w.environment.Name != "" {
		w.Report(alerts.NewEnvironmentRedaclaration(node.GetToken().Span))
		return
	}
	switch node.EnvType.Token.Lexeme {
	case "Level":
		node.EnvType.Type = ast.LevelEnv
	case "Mesh":
		node.EnvType.Type = ast.MeshEnv
	case "Sound":
		node.EnvType.Type = ast.SoundEnv
	case "Shared":
		node.EnvType.Type = ast.SharedEnv
	default:
		w.Report(alerts.NewInvalidEnvironmentType(node.EnvType.Token.Span, node.EnvType.Token.Lexeme))
	}
	w.environment.Type = node.EnvType.Type
	w.environment.Name = node.Env.Path.Lexeme
	w.environment._envStmt = node
	if w2, ok := w.walkers[w.environment.Name]; ok && w2.environment.hybroidPath != w.environment.hybroidPath {
		w.Report(alerts.NewDuplicateEnvironmentNames(node.GetToken().Span, w.environment.hybroidPath, w2.environment.hybroidPath))
		return
	}

	w.walkers[w.environment.Name] = w
}

// Rewrote
func (w *Walker) aliasDeclaration(node *ast.AliasDecl, scope *Scope) {
	if scope.Parent != nil && node.IsPub {
		w.Report(alerts.NewPublicDeclarationInLocalScope(node.Token.Span))
	}
	if _, ok := scope.AliasTypes[node.Name.Lexeme]; ok {
		w.Report(alerts.NewRedeclaration(node.Token.Span, node.Name.Lexeme, "alias"))
		return
	}
	alias := NewAliasType(node.Name.Lexeme, w.typeExpression(node.Type, scope), node.IsPub)
	alias.Token = node.Token
	scope.AliasTypes[node.Name.Lexeme] = alias
}

func (w *Walker) classDeclaration(node *ast.ClassDecl, scope *Scope) {
	if scope.Parent != nil {
		w.Report(alerts.NewInvalidStmtInLocalBlock(node.Token.Span, "class declaration"))
		return
	}

	if node.Constructor == nil {
		w.Report(alerts.NewMissingConstructor(node.Token.Span, "new", "in class declaration"))
	}

	if w.typeExists(node.Name.Lexeme) {
		w.Report(alerts.NewTypeRedeclaration(node.Name.Span, node.Name.Lexeme))
	}

	classVal := &ClassVal{
		Token:   node.Name,
		Type:    *NewNamedType(w.environment.Name, node.Name.Lexeme, ast.Class),
		IsPub:   node.IsPub,
		Fields:  make(map[string]Field),
		Methods: map[string]*VariableVal{},
		New:     NewFunction(nil),
	}
	for _, param := range node.GenericParams {
		generic := NewGeneric(param.Name.Lexeme)
		classVal.Type.Generics = append(classVal.Type.Generics, GenericWithType{GenericName: generic.Name, Type: UnknownTyp})
	}

	// DECLARATIONS
	w.declareClass(classVal)
	classScope := w.NewScope(scope, &ClassTag{Val: classVal}, SelfAllowing)
	w.RegisterScope(classScope, node.Token, w.GetNodeEndToken(node))

	for i := range node.Fields {
		w.fieldDeclaration(&node.Fields[i], classVal, classScope, false)
	}

	for i := range node.Methods {
		w.methodDeclaration(&node.Methods[i], classVal, classScope, true)
	}

	if node.Constructor != nil {
		constructor := ast.MethodDecl{
			Name:     node.Constructor.Token,
			Params:   node.Constructor.Params,
			Generics: node.Constructor.Generics,
			IsPub:    true,
			Body:     node.Constructor.Body,
		}

		w.methodDeclaration(&constructor, classVal, classScope, true)  // declaration
		w.methodDeclaration(&constructor, classVal, classScope, false) // walking
		classVal.New = classVal.Methods["new"].Value.(*FunctionVal)
		delete(classVal.Methods, "new")
	}

	// WALKING
	for _, v := range classVal.Fields {
		if !v.Var.IsInit {
			w.Report(alerts.NewUninitializedFieldInConstructor(v.Var.Token.Span, v.Var.Name, "in class declaration"))
			break
		}
	}

	for i := range node.Methods {
		w.methodDeclaration(&node.Methods[i], classVal, classScope, false)
	}
}

func (w *Walker) entityDeclaration(node *ast.EntityDecl, scope *Scope) {
	et := &EntityTag{}
	entityScope := w.NewScope(scope, et, SelfAllowing)
	if scope.Parent != nil {
		w.Report(alerts.NewInvalidStmtInLocalBlock(node.Token.Span, "entity declaration"))
		return
	}
	if w.typeExists(node.Name.Lexeme) {
		w.Report(alerts.NewTypeRedeclaration(node.Name.Span, node.Name.Lexeme))
	}
	if node.Destroyer == nil {
		w.Report(alerts.NewMissingDestroy(node.Token.Span))
		return
	} else if node.Spawner == nil {
		w.Report(alerts.NewMissingConstructor(node.Token.Span, "spawn", "in entity declaration"))
		return
	}

	entityVal := NewEntityVal(w.environment.Name, node)
	for _, param := range node.GenericParams {
		generic := NewGeneric(param.Name.Lexeme)
		entityVal.Type.Generics = append(entityVal.Type.Generics, GenericWithType{GenericName: generic.Name, Type: UnknownTyp})
	}

	et.EntityVal = entityVal
	w.declareEntity(entityVal)

	w.RegisterScope(entityScope, node.Token, w.GetNodeEndToken(node))

	// DECLARATIONS
	for i := range node.Fields {
		w.fieldDeclaration(&node.Fields[i], entityVal, entityScope, false)
	}
	for i := range node.Methods {
		w.methodDeclaration(&node.Methods[i], entityVal, entityScope, true)
	}

	fn := w.entityFunctionDeclaration(node.Destroyer, entityScope)
	entityVal.Destroy = fn
	fn = w.entityFunctionDeclaration(node.Spawner, entityScope)
	entityVal.Spawn = fn

	//callbacks
	found := map[ast.EntityFunctionType][]tokens.Token{}
	for i := range node.Methods {
		w.methodDeclaration(&node.Methods[i], entityVal, entityScope, false)
	}
	for i := range node.Callbacks {
		found[node.Callbacks[i].Type] = append(found[node.Callbacks[i].Type], node.Callbacks[i].Token)
		w.entityFunctionDeclaration(node.Callbacks[i], entityScope)
	}
	for k := range found {
		if len(found[k]) > 1 {
			w.Report(alerts.NewRedeclaration(found[k][1].Span, string(k), "entity function"))
		}
	}

	for _, v := range entityVal.Fields {
		if !v.Var.IsInit {
			w.Report(alerts.NewUninitializedFieldInConstructor(v.Var.Token.Span, v.Var.Name, "in entity declaration"))
			break
		}
	}
}

func (w *Walker) entityFunctionDeclaration(node *ast.EntityFunctionDecl, scope *Scope) *FunctionVal {
	ft := &FuncTag{
		Return: false,
	}
	fnScope := w.NewScope(scope, ft, ReturnAllowing)
	w.RegisterScope(fnScope, node.Token, w.GetNodeEndToken(node))
	ft.Generics = w.getGenericParams(node.Generics, scope)

	ft.ReturnTypes = w.getReturns(node.Returns, fnScope)
	params := w.getParameters(node.Params, fnScope)

	funcSign := NewFuncSignature(ft.Generics...).
		WithParams(params...).
		WithReturns(ft.ReturnTypes...)

	switch node.Type {
	case ast.Spawn:
		if len(params) < 2 || !(params[0].GetType() == Fixed && params[1].GetType() == Fixed) {
			w.Report(alerts.NewInvalidSpawnerParameters(node.GetToken().Span))
			break
		}
		if node.Params[0].Name.Lexeme == "_" {
			w.Report(alerts.NewEmptyIdentifierOnSpawnParameters(node.Params[0].Name.Span))
		} else if variable, ok := fnScope.Variables["x"]; !ok {
			w.Report(alerts.NewInvalidSpawnerParameter(node.Params[0].Name.Span, "first", "x"))
		} else {
			variable.IsUsed = true
		}
		if node.Params[1].Name.Lexeme == "_" {
			w.Report(alerts.NewEmptyIdentifierOnSpawnParameters(node.Params[1].Name.Span))
		} else if variable, ok := fnScope.Variables["y"]; !ok {
			w.Report(alerts.NewInvalidSpawnerParameter(node.Params[0].Name.Span, "second", "y"))
		} else {
			variable.IsUsed = true
		}
	case ast.WallCollision:
		if !funcSign.Equals(WallCollisionSign) {
			w.Report(alerts.NewInvalidEntityFunctionSignature(node.GetToken().Span, funcSign.String(), WallCollisionSign.String(), string(node.Type)))
		}
	case ast.PlayerCollision:
		if !funcSign.Equals(PlayerCollisionSign) {
			w.Report(alerts.NewInvalidEntityFunctionSignature(node.GetToken().Span, funcSign.String(), PlayerCollisionSign.String(), string(node.Type)))
		}
	case ast.WeaponCollision:
		if !funcSign.Equals(WeaponCollisionSign) {
			w.Report(alerts.NewInvalidEntityFunctionSignature(node.GetToken().Span, funcSign.String(), WeaponCollisionSign.String(), string(node.Type)))
		}
	}

	w.walkFuncBody(node, &node.Body, ft, fnScope)

	if node.Type == ast.Destroy && !ft.GetIfExits(EntityDestruction) {
		w.Report(alerts.NewNotAllCodePathsExit(node.Token.Span, "destroy the entity"))
	}

	paramNames := make([]string, len(node.Params))
	for i, param := range node.Params {
		paramNames[i] = param.Name.Lexeme
	}

	return NewFunction(paramNames, params...).WithGenerics(ft.Generics...).WithReturns(ft.ReturnTypes...)
}

func (w *Walker) enumDeclaration(node *ast.EnumDecl, scope *Scope) {
	enumVal := &EnumVal{
		Type:   NewEnumType(scope.Environment.Name, node.Name.Lexeme),
		Fields: make(map[string]*VariableVal),
		IsPub:  node.IsPub,
	}

	for _, v := range node.Fields {
		if _, _, found := enumVal.ContainsField(v.Name.Lexeme); found {
			w.Report(alerts.NewDuplicateElement(v.GetToken().Span, "enum field", v.Name.Lexeme))
			continue
		}
		variable := NewVariable(v.Name, &EnumFieldVal{Type: enumVal.Type}, node.IsPub)
		enumVal.AddField(variable)
	}

	if w.typeExists(node.Name.Lexeme) {
		w.Report(alerts.NewTypeRedeclaration(node.Name.Span, node.Name.Lexeme))
		return
	}

	enumVal.Token = node.Name
	w.environment.Enums[node.Name.Lexeme] = enumVal
}

func (w *Walker) fieldDeclaration(node *ast.VariableDecl, container FieldContainer, scope *Scope, allowSelf bool) {
	if !allowSelf {
		scope.Attributes.Remove(SelfAllowing)
	}

	w.variableDeclaration(node, scope, true)
	for _, v := range node.Identifiers {
		variable, ok := scope.Variables[v.Name.Lexeme]
		if ok {
			scope.Variables[variable.Name] = variable
			container.AddField(variable)
		}
	}
	if !allowSelf {
		scope.Attributes.Add(SelfAllowing)
	}
}

func (w *Walker) methodDeclaration(node *ast.MethodDecl, container MethodContainer, scope *Scope, declare bool) {
	if !declare {
		variable, found := container.ContainsMethod(node.Name.Lexeme)
		if !found {
			panic("Method Declaration was called on declare = false, expecting the declaration to have already happened, but couldn't find the method.")
		}
		fn := variable.Value.(*FunctionVal)
		fnTag := &FuncTag{
			Return:      false,
			ReturnTypes: fn.Returns,
			Generics:    fn.Generics,
		}

		fnScope := w.NewScope(scope, fnTag, ReturnAllowing)
		w.RegisterScope(fnScope, node.Name, w.GetNodeEndToken(node))

		for i := range node.Params {
			param := &node.Params[i]
			variable := NewVariable(param.Name, w.typeToValue(fn.Params[i]))
			w.declareVariable(fnScope, variable)
		}
		w.walkFuncBody(node, &node.Body, fnTag, fnScope)
	} else {
		funcExpr := ast.FunctionDecl{
			Name:     node.Name,
			Returns:  node.Returns,
			Params:   node.Params,
			Generics: node.Generics,
			Body:     node.Body,
			IsPub:    false,
		}

		variable := w.functionDeclaration(&funcExpr, scope, Method)
		fn := variable.Value.(*FunctionVal)
		fn.ProcType = Method
		var methodType ast.MethodCallType = ast.EntityMethod
		if scope.Tag.GetType() == Class {
			methodType = ast.ClassMethod
		}
		namedType := container.GetType().(*NamedType)
		fn.MethodInfo = ast.NewMethodInfo(methodType, funcExpr.Name.Lexeme, namedType.Name, namedType.EnvName)
		container.AddMethod(variable)
	}
}

func (w *Walker) functionDeclaration(node *ast.FunctionDecl, scope *Scope, procType ProcedureType) *VariableVal {
	ft := &FuncTag{
		Return: false,
	}
	fnScope := w.NewScope(scope, ft, ReturnAllowing)
	w.RegisterScope(fnScope, node.Token, w.GetNodeEndToken(node))
	ft.Generics = w.getGenericParams(node.Generics, scope)

	ft.ReturnTypes = w.getReturns(node.Returns, fnScope)
	params := w.getParameters(node.Params, fnScope)

	paramNames := make([]string, len(node.Params))
	for i, param := range node.Params {
		paramNames[i] = param.Name.Lexeme
	}

	variable := &VariableVal{
		Name: node.Name.Lexeme,
		Value: NewFunction(paramNames, params...).
			WithGenerics(ft.Generics...).
			WithReturns(ft.ReturnTypes...),
		Token: node.Name,
		IsPub: node.IsPub,
	}

	if _, success := w.declareVariable(scope, variable); !success {
		w.Report(alerts.NewRedeclaration(node.Name.Span, node.Name.Lexeme, "variable"))
	}

	if procType == Function {
		w.walkFuncBody(node, &node.Body, ft, fnScope)
	}

	return variable
}

// Rewrote
func (w *Walker) variableDeclaration(declaration *ast.VariableDecl, scope *Scope, allowUnitialized bool) {
	//check if it's a public declaration in a local scope
	if declaration.IsPub && scope.Parent != nil {
		w.Report(alerts.NewPublicDeclarationInLocalScope(declaration.Token.Span))
		declaration.IsPub = false
	}

	var declType Type = UnknownTyp
	if declaration.Type != nil {
		declType = w.typeExpression(declaration.Type, scope)
	}
	variables := make([]*VariableVal, 0)
	values := make([]Value2, 0)
	exprCounter := 0
	for i := range declaration.Identifiers {
		ident := declaration.Identifiers[i]
		variable := NewVariable(ident.GetToken(), &Invalid{})

		if _, exists := scope.Variables[ident.Name.Lexeme]; exists {
			w.Report(alerts.NewRedeclaration(ident.Name.Span, ident.Name.Lexeme, "variable"))
		} else {
			variable.IsPub = declaration.IsPub
			variable.IsConst = declaration.IsConst
			variables = append(variables, variable)
		}

		if i <= len(values)-1 {
			variable.Value = values[i].Value
			variable.IsInit = true
		} else if exprCounter < len(declaration.Expressions) {
			val := w.GetActualNodeValue(&declaration.Expressions[exprCounter], scope)
			if vls, ok := val.(Values); ok {
				for _, v := range vls {
					values = append(values, Value2{v, i})
				}
			} else {
				values = append(values, Value2{val, i})
			}

			variable.Value = values[i].Value
			variable.IsInit = true
			exprCounter++
		} else if declaration.IsConst {
			w.Report(alerts.NewNoValueGivenForConstant(ident.Name.Span))
			continue
		} else if declaration.Type == nil {
			w.Report(alerts.NewExplicitTypeRequiredInDeclaration(ident.Name.Span, "to infer the value"))
			continue
		} else {
			val := w.typeToValue(declType)
			defaultVal := val.GetDefault()

			if defaultVal.Value == "nil" && !allowUnitialized {
				w.Report(alerts.NewExplicitTypeNotAllowed(declaration.Type.GetToken().Span, declType.String()))
				continue
			}

			variable.Value = val
			variable.IsInit = true
			declaration.Expressions = append(declaration.Expressions, defaultVal)
			exprCounter++
		}

		valType := variable.GetType()
		if declaration.IsConst {
			variable.Value = &ConstVal{
				Node: ident,
				Val:  variable.Value,
			}
			scope.ConstValues[variable.Name] = declaration.Expressions[values[i].Index]
			continue
		}
		if declType == nil {
			continue
		}
		if valType == UnknownTyp {
			w.Report(alerts.NewInvalidType(declaration.Expressions[values[i].Index].GetToken().Span, "unknown", "as a variable value"))
			continue
		}
		if declType.GetType() == RawEntity && valType.PVT() == ast.Number {
			variable.Value = &RawEntityVal{}
		} else if !TypeEquals(declType, valType) && declType != InvalidType && valType != InvalidType {
			w.Report(alerts.NewExplicitTypeMismatch(variable.Token.Span, declType.String(), valType.String()))
		}
	}

	for i := range variables {
		w.declareVariable(scope, variables[i])
	}
	exprsLen := len(declaration.Expressions)

	if exprCounter != exprsLen {
		for range exprsLen - exprCounter {
			val := w.GetActualNodeValue(&declaration.Expressions[exprCounter], scope)
			if vls, ok := val.(Values); ok {
				for _, v := range vls {
					values = append(values, Value2{v, exprCounter})
				}
			} else {
				values = append(values, Value2{val, exprCounter})
			}
			exprCounter++
		}
	}

	varsLen, valsLen := len(variables), len(values)
	if varsLen < valsLen {
		extraAmount := valsLen - varsLen
		if extraAmount == 1 {
			w.Report(alerts.NewTooManyElementsGiven(declaration.Expressions[exprsLen-1].GetToken().Span, extraAmount, "value", "in variable declaration"))
		} else {
			w.Report(alerts.NewTooManyElementsGiven(core.MergeSpans(declaration.Expressions[values[valsLen-extraAmount].Index].GetToken().Span, declaration.Expressions[values[valsLen-1].Index].GetToken().Span), extraAmount, "value", "in variable declaration"))
		}
	}
}
