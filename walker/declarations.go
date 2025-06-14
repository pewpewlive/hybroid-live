package walker

import (
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/tokens"
)

// Rewrote
func (w *Walker) environmentDeclaration(node *ast.EnvironmentDecl) {
	if w.environment.Name != "" {
		w.AlertSingle(&alerts.EnvironmentRedaclaration{}, node.GetToken())
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
		w.AlertSingle(&alerts.InvalidEnvironmentType{}, node.EnvType.Token, node.EnvType.Token.Lexeme)
	}
	w.environment.Type = node.EnvType.Type
	w.environment.Name = node.Env.Path.Lexeme
	w.environment._envStmt = node
	if w2, ok := w.walkers[w.environment.Name]; ok {
		w.AlertSingle(&alerts.DuplicateEnvironmentNames{}, node.GetToken(), w.environment.hybroidPath, w2.environment.hybroidPath)
		return
	}

	w.walkers[w.environment.Name] = w
}

// Rewrote
func (w *Walker) aliasDeclaration(node *ast.AliasDecl, scope *Scope) {
	if scope.Parent != nil && node.IsPub {
		w.AlertSingle(&alerts.PublicDeclarationInLocalScope{}, node.Token)
	}
	if _, ok := scope.AliasTypes[node.Name.Lexeme]; ok {
		w.AlertSingle(&alerts.Redeclaration{}, node.Token, node.Name, "alias")
		return
	}
	alias := NewAliasType(node.Name.Lexeme, w.typeExpression(node.Type, scope), !node.IsPub)
	alias.Token = node.Token
	scope.AliasTypes[node.Name.Lexeme] = alias
}

func (w *Walker) classDeclaration(node *ast.ClassDecl, scope *Scope) {
	if scope.Parent != nil {
		w.AlertSingle(&alerts.InvalidStmtInLocalBlock{}, node.Token, "class declaration")
		return
	}

	if node.Constructor == nil {
		w.AlertSingle(&alerts.MissingConstructor{}, node.Token, "new", "in class declaration")
	}

	if w.typeExists(node.Name.Lexeme) {
		w.AlertSingle(&alerts.TypeRedeclaration{}, node.Name, node.Name.Lexeme)
	}

	classVal := &ClassVal{
		Token:   node.Name,
		Type:    *NewNamedType(w.environment.Name, node.Name.Lexeme, ast.Class),
		IsLocal: node.IsPub,
		Fields:  make(map[string]Field),
		Methods: map[string]*VariableVal{},
		New:     NewFunction(),
	}
	for _, param := range node.GenericParams {
		generic := NewGeneric(param.Name.Lexeme)
		classVal.Type.Generics = append(classVal.Type.Generics, GenericWithType{GenericName: generic.Name, Type: UnknownTyp})
	}

	// DECLARATIONS
	w.declareClass(classVal)
	classScope := NewScope(scope, &ClassTag{Val: classVal}, SelfAllowing)

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
			w.AlertSingle(&alerts.UninitializedFieldInConstructor{}, v.Var.Token, v.Var.Name, "in class declaration")
			break
		}
	}

	for i := range node.Methods {
		w.methodDeclaration(&node.Methods[i], classVal, classScope, false)
	}
}

func (w *Walker) entityDeclaration(node *ast.EntityDecl, scope *Scope) {
	et := &EntityTag{}
	entityScope := NewScope(scope, et, SelfAllowing)
	if scope.Parent != nil {
		w.AlertSingle(&alerts.InvalidStmtInLocalBlock{}, node.Token, "entity declaration")
		return
	}
	if w.typeExists(node.Name.Lexeme) {
		w.AlertSingle(&alerts.TypeRedeclaration{}, node.Name, node.Name.Lexeme)
	}
	if node.Destroyer == nil {
		w.AlertSingle(&alerts.MissingDestroy{}, node.Token)
		return
	} else if node.Spawner == nil {
		w.AlertSingle(&alerts.MissingConstructor{}, node.Token, "spawn", "in entity declaration")
		return
	}

	entityVal := NewEntityVal(w.environment.Name, node)
	for _, param := range node.GenericParams {
		generic := NewGeneric(param.Name.Lexeme)
		entityVal.Type.Generics = append(entityVal.Type.Generics, GenericWithType{GenericName: generic.Name, Type: UnknownTyp})
	}

	et.EntityVal = entityVal
	w.declareEntity(entityVal)

	// DECLARATIONS
	for i := range node.Fields {
		w.fieldDeclaration(&node.Fields[i], entityVal, entityScope, false)
	}

	for i := range node.Methods {
		w.methodDeclaration(&node.Methods[i], entityVal, entityScope, true)
	}

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
			w.AlertSingle(&alerts.Redeclaration{}, found[k][1], k, "entity function")
		}
	}

	fn := w.entityFunctionDeclaration(node.Destroyer, entityScope)
	entityVal.Destroy = fn
	fn = w.entityFunctionDeclaration(node.Spawner, entityScope)
	entityVal.Spawn = fn

	for _, v := range entityVal.Fields {
		if !v.Var.IsInit {
			w.AlertSingle(&alerts.UninitializedFieldInConstructor{}, v.Var.Token, v.Var.Name, "in entity declaration")
			break
		}
	}
}

func (w *Walker) entityFunctionDeclaration(node *ast.EntityFunctionDecl, scope *Scope) *FunctionVal {
	ft := &FuncTag{
		Return: false,
	}
	fnScope := NewScope(scope, ft, ReturnAllowing)
	ft.Generics = w.getGenericParams(node.Generics, scope)

	ft.ReturnTypes = w.getReturns(node.Returns, fnScope)
	params := w.getParameters(node.Params, fnScope)

	funcSign := NewFuncSignature(ft.Generics...).
		WithParams(params...).
		WithReturns(ft.ReturnTypes...)

	switch node.Type {
	case ast.Spawn:
		if len(params) < 2 || !(params[0].GetType() == Fixed && params[1].GetType() == Fixed) {
			w.AlertSingle(&alerts.InvalidSpawnerParameters{}, node.GetToken())
			break
		}
		if node.Params[0].Name.Lexeme == "_" {
			w.AlertSingle(&alerts.EmptyIdentifierOnSpawnParameters{}, node.Params[0].Name)
		} else {
			fnScope.Variables["x"].IsUsed = true // its used regardless of user input in the generator (to create the customizable entity)
		}
		if node.Params[1].Name.Lexeme == "_" {
			w.AlertSingle(&alerts.EmptyIdentifierOnSpawnParameters{}, node.Params[1].Name)
		} else {
			fnScope.Variables["y"].IsUsed = true // its used regardless of user input in the generator (to create the customizable entity)
		}
	case ast.WallCollision:
		if !funcSign.Equals(WallCollisionSign) {
			w.AlertSingle(&alerts.InvalidEntityFunctionSignature{}, node.GetToken(), funcSign, WallCollisionSign, node.Type)
		}
	case ast.PlayerCollision:
		if !funcSign.Equals(PlayerCollisionSign) {
			w.AlertSingle(&alerts.InvalidEntityFunctionSignature{}, node.GetToken(), funcSign, PlayerCollisionSign, node.Type)
		}
	case ast.WeaponCollision:
		if !funcSign.Equals(WeaponCollisionSign) {
			w.AlertSingle(&alerts.InvalidEntityFunctionSignature{}, node.GetToken(), funcSign, WeaponCollisionSign, node.Type)
		}
	}

	w.walkFuncBody(node, &node.Body, ft, fnScope)

	if node.Type == ast.Destroy && !ft.GetIfExits(EntityDestruction) {
		w.AlertSingle(&alerts.NotAllCodePathsExit{}, node.Token, "destroy the entity")
	}

	return NewFunction(params...).WithGenerics(ft.Generics...).WithReturns(ft.ReturnTypes...)
}

func (w *Walker) enumDeclaration(node *ast.EnumDecl, scope *Scope) {
	enumVal := &EnumVal{
		Type:   NewEnumType(scope.Environment.Name, node.Name.Lexeme),
		Fields: make(map[string]*VariableVal),
	}

	for _, v := range node.Fields {
		if _, _, found := enumVal.ContainsField(v.Name.Lexeme); found {
			w.AlertSingle(&alerts.DuplicateElement{}, v.GetToken(), "enum field", v.Name.Lexeme)
			continue
		}
		variable := NewVariable(v.Name, &EnumFieldVal{Type: enumVal.Type}, node.IsPub)
		enumVal.AddField(variable)
	}

	if w.typeExists(node.Name.Lexeme) {
		w.AlertSingle(&alerts.TypeRedeclaration{}, node.Name, node.Name.Lexeme)
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

		fnScope := NewScope(scope, fnTag, ReturnAllowing)

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
	fnScope := NewScope(scope, ft, ReturnAllowing)
	ft.Generics = w.getGenericParams(node.Generics, scope)

	ft.ReturnTypes = w.getReturns(node.Returns, fnScope)
	params := w.getParameters(node.Params, fnScope)

	variable := &VariableVal{
		Name: node.Name.Lexeme,
		Value: NewFunction(params...).
			WithGenerics(ft.Generics...).
			WithReturns(ft.ReturnTypes...),
		Token: node.Name,
		IsPub: node.IsPub,
	}

	if _, success := w.declareVariable(scope, variable); !success {
		w.AlertSingle(&alerts.Redeclaration{}, node.Name, node.Name.Lexeme, "variable")
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
		w.AlertSingle(&alerts.PublicDeclarationInLocalScope{}, declaration.Token)
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

		if _, alreadyExists := scope.Variables[ident.Name.Lexeme]; alreadyExists {
			w.AlertSingle(&alerts.Redeclaration{}, ident.Name, ident.Name.Lexeme, "variable")
		} else {
			variable.IsPub = declaration.IsPub
			variable.IsConst = declaration.IsConst
			variables = append(variables, variable)
		}

		if i <= len(values)-1 {
			variable.Value = values[i].Value
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
			exprCounter++
		} else if declaration.IsConst {
			w.AlertSingle(&alerts.NoValueGivenForConstant{}, ident.Name)
			continue
		} else if declaration.Type == nil {
			w.AlertSingle(&alerts.ExplicitTypeRequiredInDeclaration{}, ident.Name, "to infer the value")
			continue
		} else {
			val := w.typeToValue(declType)
			defaultVal := val.GetDefault()

			if defaultVal.Value == "nil" && !allowUnitialized {
				w.AlertSingle(&alerts.ExplicitTypeNotAllowed{}, declaration.Type.GetToken(), declType.String())
				continue
			}

			variable.Value = val
			declaration.Expressions = append(declaration.Expressions, defaultVal)
			exprCounter++
		}
		variable.IsInit = true

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
		if declType.GetType() == RawEntity && valType.PVT() == ast.Number {
			variable.Value = &RawEntityVal{}
		} else if !TypeEquals(declType, valType) && declType != InvalidType && valType != InvalidType {
			w.AlertSingle(&alerts.ExplicitTypeMismatch{},
				variable.Token,
				declType.String(),
				valType.String(),
			)
		}
	}

	for i := range variables {
		w.declareVariable(scope, variables[i])
	}
	exprsLen := len(declaration.Expressions)
	varsLen, valsLen := len(variables), len(values)
	if varsLen < valsLen {
		extraAmount := valsLen - varsLen
		if extraAmount == 1 {
			w.AlertSingle(&alerts.TooManyElementsGiven{},
				declaration.Expressions[exprsLen-1].GetToken(),
				extraAmount,
				"value",
				"in variable declaration",
			)
		} else {
			w.AlertMulti(&alerts.TooManyElementsGiven{},
				declaration.Expressions[values[valsLen-extraAmount].Index].GetToken(),
				declaration.Expressions[values[valsLen-1].Index].GetToken(),
				extraAmount,
				"value",
				"in variable declaration",
			)
		}
	}
}
