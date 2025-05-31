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
	scope.AliasTypes[node.Name.Lexeme] = NewAliasType(node.Name.Lexeme, w.typeExpression(node.Type, scope), !node.IsPub)
}

func (w *Walker) classDeclaration(node *ast.ClassDecl, scope *Scope) {
	if node.Constructor == nil {
		w.AlertSingle(&alerts.MissingConstructor{}, node.Token, "new", "in class declaration")
	}

	if scope.Parent != nil {
		w.AlertSingle(&alerts.InvalidStmtInLocalBlock{}, node.Token, "class declaration")
		return
	}

	if w.typeExists(node.Name.Lexeme) {
		w.AlertSingle(&alerts.TypeRedeclaration{}, node.Name, node.Name.Lexeme)
	}

	generics := make([]*GenericType, 0)
	if node.Constructor != nil {
		for _, param := range node.Constructor.Generics {
			generics = append(generics, NewGeneric(param.Name.Lexeme))
		}
	}

	classVal := &ClassVal{
		Type:     *NewNamedType(w.environment.Name, node.Name.Lexeme, ast.Class),
		IsLocal:  node.IsPub,
		Fields:   make(map[string]Field),
		Methods:  map[string]*VariableVal{},
		Generics: generics,
		Params:   []Type{},
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
		params := make([]Type, 0)
		for _, param := range node.Constructor.Params {
			params = append(params, w.typeExpression(param.Type, scope))
		}
		classVal.Params = params

		constructor := ast.MethodDecl{
			Name:     node.Constructor.Token,
			Params:   node.Constructor.Params,
			Generics: node.Constructor.Generics,
			IsPub:    true,
			Body:     node.Constructor.Body,
		}

		w.methodDeclaration(&constructor, classVal, classScope, true)  // declaration
		w.methodDeclaration(&constructor, classVal, classScope, false) // walking
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

	entityVal := NewEntityVal(w.environment.Name, node.Name.Lexeme, node.IsPub)

	// DECLARATIONS
	for i := range node.Fields {
		w.fieldDeclaration(&node.Fields[i], entityVal, entityScope, false)
	}

	et.EntityType = entityVal
	w.declareEntity(entityVal)

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

	if node.Destroyer == nil {
		w.AlertSingle(&alerts.MissingDestroy{}, node.Token)
	} else {
		fn := w.entityFunctionDeclaration(node.Destroyer, entityScope)
		entityVal.DestroyGenerics = fn.Generics
		entityVal.DestroyParams = fn.Params
	}
	if node.Spawner == nil {
		w.AlertSingle(&alerts.MissingConstructor{}, node.Token, "spawn", "in entity declaration")
	} else {
		fn := w.entityFunctionDeclaration(node.Spawner, entityScope)
		entityVal.SpawnGenerics = fn.Generics
		entityVal.SpawnParams = fn.Params
	}
	for _, v := range entityVal.Fields {
		if !v.Var.IsInit {
			w.AlertSingle(&alerts.UninitializedFieldInConstructor{}, v.Var.Token, v.Var.Name, "in entity declaration")
			break
		}
	}
}

func (w *Walker) entityFunctionDeclaration(node *ast.EntityFunctionDecl, scope *Scope) *FunctionVal {
	generics := make([]*GenericType, 0)
	for _, param := range node.Generics {
		generics = append(generics, NewGeneric(param.Name.Lexeme))
	}
	ret := w.getReturns(node.Returns, scope)

	ft := &FuncTag{
		Generics:    generics,
		ReturnTypes: ret,
		Returns:     make([]bool, 0),
	}
	fnScope := NewScope(scope, ft, ReturnAllowing)
	params := w.getParameters(node.Params, fnScope)

	funcSign := NewFuncSignature(generics...).
		WithParams(params...).
		WithReturns(ret...)

	w.walkFuncBody(node, &node.Body, ft, fnScope)

	switch node.Type {
	case ast.Spawn:
		if len(params) < 2 || !(params[0].GetType() == Fixed && params[1].GetType() == Fixed) {
			w.AlertSingle(&alerts.InvalidSpawnerParameters{}, node.GetToken())
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

	return NewFunction(params...).WithGenerics(generics...).WithReturns(ret...)
}

func (w *Walker) enumDeclaration(node *ast.EnumDecl, scope *Scope) {
	enumVal := &EnumVal{
		Type:   NewEnumType(scope.Environment.Name, node.Name.Lexeme),
		Fields: make(map[string]*VariableVal),
	}

	for _, v := range node.Fields {
		variable := NewVariable(v.Name, &EnumFieldVal{Type: enumVal.Type}, node.IsPub).Const()
		success := enumVal.AddField(variable)
		if !success {
			w.AlertSingle(&alerts.DuplicateElement{}, v.GetToken(), "enum field", v.Name.Lexeme)
		}
	}

	enumVar := NewVariable(node.Name, enumVal, node.IsPub).Const()

	if w.typeExists(enumVar.Name) {
		w.AlertSingle(&alerts.TypeRedeclaration{}, node.Name, node.Name.Lexeme)
		return
	}

	w.declareVariable(scope, enumVar)
}

func (w *Walker) fieldDeclaration(node *ast.FieldDecl, container FieldContainer, scope *Scope, allowSelf bool) {
	varDecl := ast.VariableDecl{
		Identifiers: node.Identifiers,
		Type:        node.Type,
		Expressions: node.Values,
		IsPub:       false,
		Token:       node.Token,
	}

	if !allowSelf {
		scope.Attributes.Remove(SelfAllowing)
	}

	w.variableDeclaration(&varDecl, scope, true)
	node.Values = varDecl.Expressions
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
			Returns:     make([]bool, 0),
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
		container.AddMethod(variable)
	}
}

func (w *Walker) functionDeclaration(node *ast.FunctionDecl, scope *Scope, procType ProcedureType) *VariableVal {
	generics := w.getGenericParams(node.Generics)

	funcTag := &FuncTag{Generics: generics}
	fnScope := NewScope(scope, funcTag, ReturnAllowing)

	ret := w.getReturns(node.Returns, fnScope)
	funcTag.ReturnTypes = ret

	params := w.getParameters(node.Params, fnScope)

	variable := &VariableVal{
		Name: node.Name.Lexeme,
		Value: NewFunction2(procType, params...).
			WithGenerics(generics...).
			WithReturns(ret...),
		Token: node.GetToken(),
		IsPub: node.IsPub,
	}
	w.declareVariable(scope, variable)

	if procType == Function {
		w.walkFuncBody(node, &node.Body, funcTag, fnScope)
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

	exprs := declaration.Expressions
	values := make([]Value2, 0)

	// get all values from the right side of the declaration
	for i := range declaration.Expressions {
		exprValue := w.GetNodeValue(&declaration.Expressions[i], scope)
		if vls, ok := exprValue.(Values); ok {
			for _, v := range vls {
				values = append(values, Value2{v, i})
			}
			continue
		}
		values = append(values, Value2{exprValue, i})
	}

	//compare values with the identifiers on the left side
	valuesLen := len(values)
	for i := range declaration.Identifiers {
		ident := declaration.Identifiers[i]
		if _, alreadyExists := scope.Variables[ident.Name.Lexeme]; alreadyExists {
			w.AlertSingle(&alerts.Redeclaration{}, ident.Name, ident.Name.Lexeme, "variable")
			continue
		}
		if i+1 > valuesLen && declaration.Type == nil {
			requiredAmount := len(declaration.Identifiers) - valuesLen
			lastToken := declaration.Token
			if len(exprs) != 0 {
				lastToken = exprs[len(exprs)-1].GetToken()
			}
			w.AlertSingle(&alerts.TooFewValuesGiven{},
				lastToken,
				requiredAmount,
				"variable declaration",
			)
			break
		}

		if ident.Name.Lexeme == "_" {
			continue
		}

		var value Value = nil
		if i < valuesLen {
			value = values[i].Value
		}
		variable := NewVariable(ident.Name, nil, declaration.IsPub)
		variable.IsInit = true

		if declaration.IsConst {
			if declaration.Type != nil {
				w.AlertSingle(&alerts.UnnecessaryTypeInConstDeclaration{},
					declaration.Type.GetToken(),
				)
			}

			if value == nil {
				w.AlertSingle(&alerts.NoValueGivenForConstant{}, ident.Name)
				value = &Unknown{}
			}

			variable.Value = &ConstVal{
				Node: ident,
				Val:  value,
			}
			variable.IsConst = true
			w.declareVariable(scope, variable)
			continue
		}

		var explicitType Type = nil
		if declaration.Type != nil {
			explicitType = w.typeExpression(declaration.Type, scope)
		}
		if explicitType == nil && value == nil {
			if allowUnitialized {
				declaration.Expressions = append(declaration.Expressions, &ast.LiteralExpr{Value: "nil"})
			} else {
				w.AlertSingle(&alerts.ExplicitTypeRequiredInDeclaration{}, ident.Name)
			}
			variable.IsInit = false
			value = &Unknown{}
		} else if explicitType != nil && value == nil {
			if allowUnitialized {
				print()
			}
			value = w.typeToValue(explicitType)
			defaultVal := value.GetDefault()

			if defaultVal.Value == "nil" {
				variable.IsInit = false
				if !allowUnitialized {
					w.AlertSingle(&alerts.ExplicitTypeNotAllowed{}, declaration.Type.GetToken(), explicitType.String())
				}
			}

			declaration.Expressions = append(declaration.Expressions, defaultVal)
		} else if explicitType != nil && value != nil {
			if explicitType.GetType() == RawEntity && value.GetType().PVT() == ast.Number {
				value = &RawEntityVal{}
			} else if !TypeEquals(explicitType, value.GetType()) && explicitType != InvalidType && value.GetType() != InvalidType {
				w.AlertSingle(&alerts.ExplicitTypeMismatch{},
					variable.Token,
					explicitType.String(),
					value.GetType().String(),
				)
			}
		}

		variable.Value = value
		w.declareVariable(scope, variable)
	}

	identsLen := len(declaration.Identifiers)
	if identsLen < valuesLen {
		extraAmount := valuesLen - identsLen
		if extraAmount == 1 {
			w.AlertSingle(&alerts.TooManyValuesGiven{},
				exprs[valuesLen-1].GetToken(),
				extraAmount,
				"in variable declaration",
			)
		} else {
			w.AlertMulti(&alerts.TooManyValuesGiven{},
				exprs[values[valuesLen-extraAmount].Index].GetToken(),
				exprs[values[valuesLen-1].Index].GetToken(),
				extraAmount,
				"in variable declaration",
			)
		}
	}
}
