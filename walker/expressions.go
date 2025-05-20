package walker

import (
	"hybroid/ast"
	"hybroid/generator"
	"hybroid/tokens"
)

func (w *Walker) StructExpr(node *ast.StructExpr, scope *Scope) *AnonStructVal {
	anonStructScope := NewScope(scope, &UntaggedTag{})
	structTypeVal := NewAnonStructVal(make(map[string]Field), false)

	for i := range node.Fields {
		w.FieldDeclarationStmt(node.Fields[i], structTypeVal, anonStructScope)
	}

	return structTypeVal
}

func (w *Walker) GetReturns(returns *ast.TypeExpr, scope *Scope) Types {
	returnTypes := EmptyReturn
	if returns != nil {
		if returns.Name.GetType() != ast.TupleExpression {
			returnTypes = append(returnTypes, w.TypeExpr(returns, scope, true))
		} else {
			types := returns.Name.(*ast.TupleExpr).Types
			for _, typee := range types {
				returnTypes = append(returnTypes, w.TypeExpr(typee, scope, true))
			}
		}
	}

	return returnTypes
}

func (w *Walker) FunctionExpr(fn *ast.FunctionExpr, scope *Scope) Value {
	returnTypes := w.GetReturns(fn.Return, scope)

	funcTag := &FuncTag{ReturnTypes: returnTypes}
	fnScope := NewScope(scope, funcTag, ReturnAllowing)

	w.WalkBody(&fn.Body, funcTag, fnScope)

	params := make([]Type, 0)
	for i, param := range fn.Params {
		params = append(params, w.TypeExpr(param.Type, scope, true))
		w.DeclareVariable(fnScope, &VariableVal{Name: param.Name.Lexeme, Value: w.TypeToValue(params[i]), IsLocal: true}, param.Name)
	}
	return &FunctionVal{ // returnTypes should contain a fn()
		Params:  params,
		Returns: returnTypes,
	}
}

func (w *Walker) MatchExpr(node *ast.MatchExpr, scope *Scope) Value {
	mtt := &MatchExprTag{}

	matchScope := NewScope(scope, mtt, YieldAllowing)
	casesLength := len(node.MatchStmt.Cases) + 1
	if node.MatchStmt.HasDefault {
		casesLength--
	}
	matchScope.Tag = &MatchExprTag{YieldValues: make(Types, 0)}
	mpt := NewMultiPathTag(casesLength)

	w.GetNodeValue(&node.MatchStmt.ExprToMatch, scope)

	for i := range node.MatchStmt.Cases {
		caseScope := NewScope(matchScope, mpt)
		w.GetNodeValue(&node.MatchStmt.Cases[i].Expression, matchScope)
		w.WalkBody(&node.MatchStmt.Cases[i].Body, mpt, caseScope)
	}

	yieldValues := matchScope.Tag.(*MatchExprTag).YieldValues

	node.ReturnAmount = len(yieldValues)

	return yieldValues
}

func (w *Walker) EntityExpr(node *ast.EntityExpr, scope *Scope) Value {
	val := w.GetNodeValue(&node.Expr, scope)
	typ := w.TypeExpr(node.Type, scope, false)

	if ident, ok := node.Type.Name.(*ast.IdentifierExpr); ok {
		switch ident.Name.Lexeme {
		case "Asteroid", "YellowBaf", "Inertiac", "Mothership",
			"MothershipBullet", "RollingCube", "RollingSphere",
			"Ufo", "Wary", "Crowder", "Ship", "Bomb", "BlueBaf",
			"RedBaf", "WaryMissile", "UfoBullet", "PlayerBullet",
			"BombExplosion", "PlayerExplosion", "Bonus", "FloatingMessage",
			"Pointonium", "BonusImplosion":
			typ = &RawEntityType{}
			node.OfficialEntityType = true
		}
	}

	if typ.PVT() != ast.Entity {
		// w.Error(node.Token, "type given in entity expression is not an entity type")
	} else if !node.OfficialEntityType {
		varName := tokens.Token{}
		if node.ConvertedVarName != nil {
			varName = *node.ConvertedVarName
			w.context.Conversions = append(w.context.Conversions, NewEntityConversion(varName, w.TypeToValue(typ).(*EntityVal)))
		} else if len(w.context.Conversions) != 0 {
			w.context.Conversions = append(w.context.Conversions, NewEntityConversion(varName, w.TypeToValue(typ).(*EntityVal)))
		}
		entityVal := w.TypeToValue(typ).(*EntityVal)
		node.EntityName = entityVal.Type.Name
		node.EnvName = entityVal.Type.EnvName
	} else if node.ConvertedVarName != nil {
		// w.Error(*node.ConvertedVarName, "can't convert an entity to an official entity")
	}

	if val.GetType().GetType() != RawEntity {
		// w.Error(node.Token, "value given in entity expression is not an entity")
	}

	return &BoolVal{}
}

func (w *Walker) BinaryExpr(node *ast.BinaryExpr, scope *Scope) Value {
	left, right := w.GetNodeValue(&node.Left, scope), w.GetNodeValue(&node.Right, scope)
	leftType, rightType := left.GetType(), right.GetType()
	op := node.Operator
	switch op.Type {
	case tokens.Plus, tokens.Minus, tokens.Caret, tokens.Star, tokens.Slash, tokens.Modulo, tokens.BackSlash:
		w.ValidateArithmeticOperands(leftType, rightType, node)
		typ := w.DetermineValueType(leftType, rightType)

		if typ.PVT() == ast.Invalid {
			// w.Error(node.GetToken(), fmt.Sprintf("invalid binary expression (left: %s, right: %s)", leftType.ToString(), rightType.ToString()))
			return &Invalid{}
		}

		return w.TypeToValue(typ)
	case tokens.Concat:
		if !TypeEquals(leftType, NewBasicType(ast.String)) && !TypeEquals(rightType, NewBasicType(ast.String)) {
			// w.Error(node.GetToken(), fmt.Sprintf("invalid concatenation: left is %s and right is %s", leftType.ToString(), rightType.ToString()))
			return &Invalid{}
		}
		return &StringVal{}
	default:
		if op.Type == tokens.Or {
			if node.Left.GetType() == ast.EntityExpression && node.Left.(*ast.EntityExpr).ConvertedVarName != nil {
				// w.Error(node.Left.GetToken(), "conversion of entity is not possible in a binary expression with 'or' operator")
			} else if node.Right.GetType() == ast.EntityExpression && node.Right.(*ast.EntityExpr).ConvertedVarName != nil {
				// w.Error(node.Right.GetToken(), "conversion of entity is not possible in a binary expression with 'or' operator")
			}
		}

		if !TypeEquals(leftType, rightType) {
			// w.Error(node.GetToken(), fmt.Sprintf("invalid comparison: types are not the same (left: %s, right: %s)", leftType.ToString(), rightType.ToString()))
			return &Invalid{}
		}
		return &BoolVal{}
	}
}

func (w *Walker) LiteralExpr(node *ast.LiteralExpr) Value {
	switch node.ValueType {
	case ast.String:
		return &StringVal{}
	case ast.Fixed, ast.Radian, ast.FixedPoint, ast.Degree:
		return &FixedVal{SpecificType: node.ValueType}
	case ast.Bool:
		return &BoolVal{}
	case ast.Number:
		return &NumberVal{}
	case ast.Entity:
		return &RawEntityVal{}
	default:
		return &Invalid{}
	}
}

func ConvertNodeToFieldExpr(ident ast.Node, index int, exprType ast.SelfExprType, envName string, entityName string) *ast.FieldExpr {
	fieldExpr := &ast.FieldExpr{
		Index: index,
		Identifier: &ast.SelfExpr{
			Token: ident.GetToken(),
			Type:  exprType,
		},
		ExprType:   exprType,
		EnvName:    envName,
		EntityName: entityName,
	}

	fieldExpr.Property = ident
	if access, ok := ident.(ast.Accessor); ok {
		fieldExpr.PropertyIdentifier = access.GetIdentifier()
	} else {
		fieldExpr.PropertyIdentifier = &ast.IdentifierExpr{Name: ident.GetToken()}
	}

	return fieldExpr
}

func ConvertCallToMethodCall(call *ast.CallExpr, exprType ast.SelfExprType, envName string, name string) *ast.MethodCallExpr {
	copy := *call
	return &ast.MethodCallExpr{
		EnvName:  envName,
		TypeName: name,
		ExprType: exprType,
		Identifier: &ast.SelfExpr{
			EntityName: name,
			Type:       exprType,
		},
		Call:       &copy,
		MethodName: call.Caller.GetToken().Lexeme,
	}
}

func (w *Walker) IdentifierExpr(node *ast.Node, scope *Scope) Value {
	valueNode := *node
	ident := valueNode.(*ast.IdentifierExpr)

	sc := w.ResolveVariable(scope, ident.Name.Lexeme)
	if sc == nil {
		walker, found := w.walkers[ident.Name.Lexeme]
		if found {
			*node = &ast.LiteralExpr{
				Value: "\"" + walker.environment.luaPath + "\"",
			}
			return NewPathVal(walker.environment.luaPath, walker.environment.Type)
		}

		return &Invalid{}
	}

	variable, notAllowed := w.GetVariable(sc, ident.Name.Lexeme)
	if notAllowed {
		// w.Error(ident.GetToken(), "Not allowed to access a local variable from a different environment")
	}

	if sc.Tag.GetType() == Struct {
		class := sc.Tag.(*ClassTag).Val
		if variable.Value.GetType().GetType() == Fn {
			w.context.Value2 = class
			return variable
		}
		field, index, found := class.ContainsField(variable.Name)

		*node = ConvertNodeToFieldExpr(ident, index, ast.SelfStruct, class.Type.EnvName, "")

		if found {
			return field
		}
		method, found := class.Methods[variable.Name]
		if found {
			return method
		}
	} else if sc.Tag.GetType() == Entity {
		entity := sc.Tag.(*EntityTag).EntityType
		if variable.Value.GetType().GetType() == Fn {
			w.context.Value2 = entity
			return variable
		}
		field, index, found := entity.ContainsField(variable.Name)

		*node = ConvertNodeToFieldExpr(ident, index, ast.SelfEntity, entity.Type.EnvName, entity.Type.Name)

		if found {
			return field
		}
		method, found := entity.Methods[variable.Name]
		if found {
			return method
		}
	} else if sc.Environment.Name == "Builtin" {
		scope.Environment.AddBuiltinVar(ident.Name.Lexeme)
		*node = &ast.BuiltinExpr{
			Name: ident.Name,
		}
	} else if sc.Environment.Name != w.environment.Name {
		*node = &ast.EnvAccessExpr{
			PathExpr: &ast.EnvPathExpr{
				Path: tokens.Token{
					Lexeme:   sc.Environment.Name,
					Location: ident.GetToken().Location,
				},
			},
			Accessed: ident,
		}
	}

	if w.context.PewpewVarFound {
		name, found := generator.PewpewEnums[w.context.PewpewVarName][ident.Name.Lexeme]
		if found {
			ident.Name.Lexeme = name
		}
	}

	switch sc.Environment.Name {
	case "Pewpew":
		ident.Name.Lexeme = generator.PewpewVariables[ident.Name.Lexeme]
		if variable.GetType().GetType() != Fn {
			w.context.PewpewVarFound = true
			w.context.PewpewVarName = ident.Name.Lexeme
		}
	case "Fmath":
		ident.Name.Lexeme = generator.FmathFunctions[ident.Name.Lexeme]
	case "Math":
		ident.Name.Lexeme = generator.MathVariables[ident.Name.Lexeme]
	case "String":
		ident.Name.Lexeme = generator.StringVariables[ident.Name.Lexeme]
	case "Table":
		ident.Name.Lexeme = generator.TableVariables[ident.Name.Lexeme]
	}

	variable.IsUsed = true
	return variable
}

func (w *Walker) EnvAccessExpr(node *ast.EnvAccessExpr) (Value, ast.Node) {
	envName := node.PathExpr.Path.Lexeme

	if node.Accessed.GetType() == ast.Identifier {
		name := node.Accessed.(*ast.IdentifierExpr).Name.Lexeme
		path := envName + ":" + name
		walker, found := w.walkers[path]
		if found {
			return NewPathVal(walker.environment.luaPath, walker.environment.Type), &ast.LiteralExpr{
				Value: "\"" + walker.environment.luaPath + "\"",
			}
		}
	}

	switch envName {
	case "Pewpew":
		if w.environment.Type != ast.LevelEnv {
			// w.Error(node.GetToken(), "cannot use the pewpew library in a non-level environment")
		}
		return w.GetNodeValueFromExternalEnv(node.Accessed, PewpewEnv), nil
	case "Fmath":
		if w.environment.Type != ast.LevelEnv {
			// w.Error(node.GetToken(), "cannot use the fmath library in a non-level environment")
		}
		return w.GetNodeValueFromExternalEnv(node.Accessed, FmathEnv), nil
	case "Math":
		return w.GetNodeValueFromExternalEnv(node.Accessed, MathEnv), nil
	case "String":
		return w.GetNodeValueFromExternalEnv(node.Accessed, StringEnv), nil
	case "Table":
		return w.GetNodeValueFromExternalEnv(node.Accessed, TableEnv), nil
	}

	walker, found := w.walkers[envName]
	if !found {
		// w.Error(node.GetToken(), fmt.Sprintf("Environment named '%s' doesn't exist", envName))
		return &Invalid{}, nil
	}

	if walker.environment.Name == w.environment.Name {
		// w.Error(node.GetToken(), "cannot access self")
		return &Invalid{}, nil
	} else if walker.environment.luaPath == "/dynamic/level.lua" {
		if !walker.Walked {
			walker.Action(w.walkers)
		}
		value := w.GetNodeValue(&node.Accessed, &walker.environment.Scope)
		return value, nil
	}

	for _, v := range walker.environment.EnvStmt.Requirements {
		if v == w.environment.luaPath {
			// w.Error(node.GetToken(), fmt.Sprintf("import cycle detected: this environment and '%s' are using each other", walker.Environment.Name))
			return &Invalid{}, nil
		}
	}

	walker.environment.EnvStmt.AddRequirement(walker.environment.luaPath)

	if !walker.Walked {
		walker.Action(w.walkers)
	}

	value := w.GetNodeValue(&node.Accessed, &walker.environment.Scope)

	return value, nil
}

func (w *Walker) GroupingExpr(node *ast.GroupExpr, scope *Scope) Value {
	return w.GetNodeValue(&node.Expr, scope)
}

func (w *Walker) ListExpr(node *ast.ListExpr, scope *Scope) Value {
	var value ListVal
	for i := range node.List {
		val := w.GetNodeValue(&node.List[i], scope)
		if val.GetType().PVT() == ast.Invalid {
			// w.Error(node.List[i].GetToken(), fmt.Sprintf("variable '%s' inside list is invalid", node.List[i].GetToken().Lexeme))
		}
		value.Values = append(value.Values, val)
	}
	value.ValueType = GetContentsValueType(value.Values)
	return &value
}

func (w *Walker) CallExpr(val Value, node *ast.CallExpr, scope *Scope) (Value, ast.Node) {
	valType := val.GetType().PVT()
	if valType != ast.Func {
		// w.Error(node.Caller.GetToken(), "caller is not a function")
		return &Invalid{}, node
	}

	var finalNode ast.Node
	finalNode = node

	if entity, ok := w.context.Value2.(*EntityVal); ok {
		caller := node.Caller.GetToken().Lexeme
		_, contains := entity.ContainsMethod(caller)
		if !contains {
			_, index, _ := entity.ContainsField(caller)
			finalNode = ConvertNodeToFieldExpr(node, index, ast.SelfEntity, entity.Type.EnvName, entity.Type.Name)
			goto skip
		}
		finalNode = ConvertCallToMethodCall(node, ast.SelfEntity, entity.Type.EnvName, entity.Type.Name)
		w.context.Value2 = &Unknown{}
	} else if class, ok := w.context.Value2.(*ClassVal); ok {
		caller := node.Caller.GetToken().Lexeme
		_, contains := class.ContainsMethod(caller)
		if !contains {
			_, index, _ := class.ContainsField(caller)
			finalNode = ConvertNodeToFieldExpr(node, index, ast.SelfStruct, class.Type.EnvName, class.Type.Name)
			goto skip
		}
		finalNode = ConvertCallToMethodCall(node, ast.SelfStruct, class.Type.EnvName, class.Type.Name)
		w.context.Value2 = &Unknown{}
	}

skip:

	variable, it_is := val.(*VariableVal)
	if it_is {
		val = variable.Value
	}
	fun, _ := val.(*FunctionVal)

	suppliedGenerics := w.GetGenerics(node, node.GenericArgs, fun.Generics, scope)

	args := []Type{}
	for i := range node.Args {
		args = append(args, w.GetNodeValue(&node.Args[i], scope).GetType())
	}
	w.ValidateArguments(suppliedGenerics, args, fun.Params, node.Caller.GetToken())

	for i := range fun.Returns {
		if fun.Returns[i].GetType() == Generic {
			fun.Returns[i] = suppliedGenerics[fun.Returns[i].(*GenericType).Name]
		}
	}

	node.ReturnAmount = len(fun.Returns)

	if node.ReturnAmount == 1 {
		return w.TypeToValue(fun.Returns[0]), finalNode
	} else if node.ReturnAmount == 0 {
		return &Invalid{}, finalNode
	}
	return &fun.Returns, finalNode
}

func (w *Walker) GetGenerics(node ast.Node, genericArgs []*ast.TypeExpr, expectedGenerics []*GenericType, scope *Scope) map[string]Type {
	receivedGenericsLength := len(genericArgs)
	expectedGenericsLength := len(expectedGenerics)

	suppliedGenerics := map[string]Type{}
	if receivedGenericsLength > expectedGenericsLength {
		// w.Error(node.GetToken(), "too many generic arguments supplied")
	} else {
		for i := range genericArgs {
			suppliedGenerics[expectedGenerics[i].Name] = w.TypeExpr(genericArgs[i], scope, true)
		}
	}

	return suppliedGenerics
}

// Writes to context
func (w *Walker) FieldExpr(node *ast.FieldExpr, scope *Scope) Value {
	if node.Identifier.GetToken().Lexeme == "converted" {
		print("brekpoint")
	}
	var scopeable ScopeableValue
	var val Value
	if w.context.Node.GetType() != ast.NA {
		val = w.context.Value
	} else {
		val = w.GetNodeValue(&node.Identifier, scope)
	}
	if variable, ok := val.(*VariableVal); ok {
		val = variable.Value
	}
	if val.GetType().GetType() == Named && val.GetType().PVT() == ast.Entity {
		node.ExprType = ast.SelfEntity
		named := val.GetType().(*NamedType)
		node.EntityName = named.Name
		node.EnvName = named.EnvName
	}

	if scpbl, ok := val.(ScopeableValue); ok {
		scopeable = scpbl
	} else {
		// w.Error(node.Identifier.GetToken(), "variable is not of type class, struct, entity or enum")
		return &Invalid{}
	}

	newScope := scopeable.Scopify(scope, node)
	w.context.Value = val
	w.context.Node = node

	propVal := w.GetNodeValue(&node.PropertyIdentifier, newScope)
	if propVal.GetType().PVT() == ast.Invalid {
		// ident := node.PropertyIdentifier.GetToken()
		// w.Error(ident, fmt.Sprintf("'%s' doesn't exist", ident.Lexeme))
	}
	w.context.Value = propVal
	w.context.Node = node.Property

	defer w.context.Clear()
	if node.Property.GetType() != ast.Identifier {
		return w.GetNodeValue(&node.Property, newScope)
	} // var1[1]["test"].method()
	return propVal
}

// Writes to context
func (w *Walker) MemberExpr(node *ast.MemberExpr, scope *Scope) Value {
	var val Value
	if w.context.Value.GetType().GetType() != NA {
		val = w.context.Value
	} else {
		val = w.GetNodeValue(&node.Identifier, scope)
	}
	valType := val.GetType()
	if valType.GetType() != Wrapper {
		// token := w.Context.Node.GetToken()
		// w.Error(token, fmt.Sprintf("%s is not a map nor a list (found %s)", token.Lexeme, valType.ToString()))
		return &Invalid{}
	}

	propValPVT := w.GetNodeValue(&node.PropertyIdentifier, scope).GetType().PVT()
	if valType.PVT() == ast.Map {
		if propValPVT != ast.String {
			// w.Error(node.Property.GetToken(), "expected string inside brackets for map accessing")
		}
	} else if valType.PVT() == ast.List {
		if propValPVT != ast.Number {
			// w.Error(node.Property.GetToken(), "expected number inside brackets for list accessing")
		}
	}
	property := w.TypeToValue(valType.(*WrapperType).WrappedType)
	w.context.Value = property
	w.context.Node = node.Property
	nodePropertyType := node.Property.GetType()
	if nodePropertyType == ast.Identifier {
		w.context.Clear()
		return property
	} else if nodePropertyType == ast.CallExpression {
		w.context.Clear()
		val, newNode := w.CallExpr(property, node.Property.(*ast.CallExpr), scope)
		node.Property = newNode
		return val
	} else if nodePropertyType == ast.LiteralExpression {
		w.context.Clear()
		return property
	}
	final := w.GetNodeValue(&node.Property, scope)
	w.context.Clear()
	return final
}

func (w *Walker) MapExpr(node *ast.MapExpr, scope *Scope) Value {
	mapVal := MapVal{Members: []Value{}}
	// for _, v := range node.Map {
	// 	val := w.GetNodeValue(&v.Expr, scope)
	// 	mapVal.Members = append(mapVal.Members, val)
	// }
	mapVal.MemberType = GetContentsValueType(mapVal.Members)
	return &mapVal
}

func (w *Walker) UnaryExpr(node *ast.UnaryExpr, scope *Scope) Value {
	val := w.GetNodeValue(&node.Value, scope)
	valType := val.GetType()
	valPVT := valType.PVT()

	// token := node.GetToken()

	if valPVT == ast.Invalid {
		// w.Error(token, "value is invalid")
		return val
	}

	switch node.Operator.Type {
	case tokens.Bang:
		if valPVT != ast.Bool {
			// w.Error(token, "value must be a bool to be negated")
		}
	case tokens.Hash:
		if valType.GetType() == Wrapper && valType.(*WrapperType).Type.PVT() != ast.List {
			// w.Error(token, "value must be a list")
		} else if valType.GetType() != Wrapper {
			// w.Error(token, "value must be a list")
		}
		return &NumberVal{}
	case tokens.Minus:
		if valPVT != ast.Number && valType.GetType() != Fixed {
			// w.Error(token, "value must be a number or fixed")
		}
	}

	return val
}

func (w *Walker) SelfExpr(self *ast.SelfExpr, scope *Scope) Value {
	if !scope.Is(SelfAllowing) {
		// w.Error(self.Token, "can't use self outside of struct/entity")
		return &Invalid{}
	}

	sc, _, structTag := ResolveTagScope[*ClassTag](scope) // TODO: CHECK FOR ENTITY SCOPE

	if sc == nil {
		entitySc, _, entityTag := ResolveTagScope[*EntityTag](scope)
		if entitySc != nil {
			self.Type = ast.SelfEntity
			self.EntityName = (*entityTag).EntityType.Type.Name
			return (*entityTag).EntityType
		}

		return &Invalid{}
	}

	(*self).Type = ast.SelfStruct
	return (*structTag).Val
}

func (w *Walker) NewExpr(new *ast.NewExpr, scope *Scope) Value {
	_type := w.TypeExpr(new.Type, scope, false)

	if _type.PVT() == ast.Invalid {
		// w.Error(new.Type.GetToken(), "invalid type given in new expression")
		return &Invalid{}
	} else if _type.PVT() != ast.Struct {
		// w.Error(new.Type.GetToken(), "type given in new expression is not a struct")
		return &Invalid{}
	}

	val := w.TypeToValue(_type).(*ClassVal)

	args := make([]Type, 0)
	for i := range new.Args {
		args = append(args, w.GetNodeValue(&new.Args[i], scope).GetType())
	}

	suppliedGenerics := w.GetGenerics(new, new.Generics, val.Generics, scope)

	w.ValidateArguments(suppliedGenerics, args, val.Params, new.Token)

	return val
}

func (w *Walker) SpawnExpr(new *ast.SpawnExpr, scope *Scope) Value {
	typeExpr := w.TypeExpr(new.Type, scope, false)

	if typeExpr.PVT() == ast.Invalid {
		// w.Error(new.Type.GetToken(), "invalid type given in spawn expression")
		return &Invalid{}
	} else if typeExpr.PVT() != ast.Entity {
		// w.Error(new.Type.GetToken(), "type given in spawn expression is not an entity")
		return &Invalid{}
	}

	val := w.TypeToValue(typeExpr).(*EntityVal)

	args := make([]Type, 0)
	for i := range new.Args {
		args = append(args, w.GetNodeValue(&new.Args[i], scope).GetType())
	}

	suppliedGenerics := w.GetGenerics(new, new.Generics, val.SpawnGenerics, scope)

	w.ValidateArguments(suppliedGenerics, args, val.SpawnParams, new.Token)

	return val
}

func (w *Walker) MethodCallExpr(mcall *ast.MethodCallExpr, scope *Scope) (Value, ast.Node) {
	val := w.GetNodeValue(&mcall.Identifier, scope)
	if variable, ok := val.(*VariableVal); ok {
		val = variable.Value
	}
	valType := val.GetType()

	if valType.PVT() == ast.Invalid {
		// w.Error(mcall.Identifier.GetToken(), "value is invalid")
		return &Invalid{}, mcall
	}

	if structVal, ok := val.(*ClassVal); ok {
		mcall.EnvName = structVal.Type.EnvName
		mcall.TypeName = structVal.Type.Name
		mcall.ExprType = ast.SelfStruct
	} else if entityVal, ok := val.(*EntityVal); ok {
		mcall.EnvName = entityVal.Type.EnvName
		mcall.TypeName = entityVal.Type.Name
		mcall.ExprType = ast.SelfEntity
	}

	callToken := mcall.Call.Caller.GetToken()
	if methodContainer, ok := val.(MethodContainer); ok {
		method, found := methodContainer.ContainsMethod(callToken.Lexeme)

		mcall.MethodName = callToken.Lexeme

		if found {
			val, _ := w.CallExpr(method, mcall.Call, scope)
			return val, mcall
		}
	}
	if fieldContainer, ok := val.(FieldContainer); ok && valType.PVT() != ast.Enum {
		_, _, found := fieldContainer.ContainsField(callToken.Lexeme)

		if found {
			fieldExpr := &ast.FieldExpr{
				Property:           mcall.Call,
				PropertyIdentifier: mcall.Call.Caller,
				Identifier:         mcall.Identifier,
			}

			return w.FieldExpr(fieldExpr, scope), fieldExpr
		}
	} else {
		// w.Error(mcall.Identifier.GetToken(), "value is not of type class or entity")
		return &Invalid{}, mcall
	}

	// w.Error(mcall.GetToken(), "no method found")

	return &Invalid{}, mcall
}

func (w *Walker) GetNodeValueFromExternalEnv(expr ast.Node, env *Environment) Value {
	val := w.GetNodeValue(&expr, &env.Scope)
	_, isTypes := val.(*Types)
	if !isTypes && val.GetType().PVT() == ast.Invalid {
		// w.Error(expr.GetToken(), fmt.Sprintf("variable named '%s' doesn't exist", expr.GetToken().Lexeme))
	}
	return val
}

// func (w *Walker) CastExpr(cast *ast.CastExpr, scope *Scope) Value {
// 	val := w.GetNodeValue(&cast.Value, scope)
// 	typ := w.TypeExpr(cast.Type, w.Environment)

// 	if typ.GetType() != CstmType {
// 		return &Invalid{}
// 	}

// 	cstm := typ.(*CustomType)

// 	if !TypeEquals(val.GetType(), cstm.UnderlyingType) {
// 		// w.Error(cast.Value.GetToken(), fmt.Sprintf("expression type is %s, but underlying type is %s", val.GetType().ToString(), cstm.UnderlyingType.ToString()))
// 		return &Invalid{}
// 	}

// 	return NewCustomVal(cstm)
// }

func (w *Walker) TypeExpr(typee *ast.TypeExpr, scope *Scope, throw bool) Type {
	if typee == nil {
		return InvalidType
	}

	var typ Type
	if typee.Name.GetType() == ast.EnvironmentAccessExpression {
		expr, _ := typee.Name.(*ast.EnvAccessExpr)
		path := expr.PathExpr.Path.Lexeme

		walker, found := w.walkers[path]
		var env *Environment
		if !found {
			switch path {
			case "Pewpew":
				env = PewpewEnv
			case "Fmath":
				env = FmathEnv
			case "Math":
				env = MathEnv
			case "String":
				env = StringEnv
			case "Table":
				env = TableEnv
			default:
				// w.Error(expr.GetToken(), "Environment name so doesn't exist")
				return InvalidType
			}

		} else if walker.environment.luaPath == "/dynamic/level.lua" {
			if !walker.Walked {
				walker.Action(w.walkers)
			}
			env = walker.environment
		} else {

			for _, v := range walker.environment.EnvStmt.Requirements {
				if v == w.environment.luaPath {
					// w.Error(typee.GetToken(), fmt.Sprintf("import cycle detected: this environment and '%s' are using each other", walker.Environment.Name))
					return InvalidType
				}
			}

			w.environment.EnvStmt.AddRequirement(walker.environment.luaPath)

			if !walker.Walked {
				walker.Action(w.walkers)
			}

			env = walker.environment
		}

		ident := &ast.IdentifierExpr{Name: expr.Accessed.GetToken(), ValueType: ast.Invalid}
		typ = w.TypeExpr(&ast.TypeExpr{Name: ident}, &env.Scope, throw)
		if typee.IsVariadic {
			return NewVariadicType(typ)
		}
		return typ
	}

	if typee.Name.GetToken().Type == tokens.Entity {
		typ = &RawEntityType{}
		if typee.IsVariadic {
			return NewVariadicType(typ)
		}
		return typ
	}

	pvt := w.GetTypeFromString(typee.Name.GetToken().Lexeme)
	switch pvt {
	case ast.Bool, ast.String, ast.Number:
		typ = NewBasicType(pvt)
	case ast.Fixed, ast.FixedPoint, ast.Radian, ast.Degree:
		typ = NewFixedPointType(pvt)
	case ast.Enum:
		typ = NewBasicType(ast.Enum)
	case ast.AnonStruct:
		fields := []*VariableVal{}

		for _, v := range typee.Fields {
			fields = append(fields, &VariableVal{
				Name:  v.Name.Lexeme,
				Value: w.TypeToValue(w.TypeExpr(v.Type, scope, throw)),
				Token: v.Name,
			})
		}

		typ = NewStructType(fields, false)
	case ast.Func:
		params := Types{}

		for _, v := range typee.Params {
			params = append(params, w.TypeExpr(v, scope, throw))
		}

		returns := w.GetReturns(typee.Return, scope)

		typ = &FunctionType{
			Params:  params,
			Returns: returns,
		}
	case ast.Map, ast.List:
		wrapped := w.TypeExpr(typee.WrappedType, scope, throw)
		typ = NewWrapperType(NewBasicType(pvt), wrapped)
	case ast.Entity:
		typ = &RawEntityType{}
	default:
		typeName := typee.Name.GetToken().Lexeme
		if entityVal, found := scope.Environment.Entities[typeName]; found {
			typ = entityVal.GetType()
			w.CheckAccessibility(scope, entityVal.IsLocal, typee)
			break
		}
		if structVal, found := scope.Environment.Structs[typeName]; found {
			typ = structVal.GetType()
			w.CheckAccessibility(scope, structVal.IsLocal, typee)
			break
		}
		if customType, found := scope.Environment.CustomTypes[typeName]; found {
			typ = customType
			//w.CheckAccessibility(scope, customType.IsLocal, typee)
			break
		}
		if aliasType, found := scope.Environment.AliasTypes[typeName]; found {
			typ = aliasType.UnderlyingType

			break
		}
		if val, _ := w.GetVariable(&scope.Environment.Scope, typeName); val != nil {
			if val.GetType().PVT() == ast.Enum {
				typ = val.GetType()
				w.CheckAccessibility(scope, val.IsLocal, typee)
				break
			}
		}

		sc, _, fnTag := ResolveTagScope[*FuncTag](scope)

		if sc != nil {
			fnTag := *fnTag
			for _, v := range fnTag.Generics {
				if v.Name == typeName {
					return v
				}
			}
		}

		if len(scope.Environment.UsedLibraries) != 0 {
			if alias, found := BuiltinEnv.AliasTypes[typeName]; found {
				typ = alias.UnderlyingType
				break
			}
		}

		if scope.Environment.Name != w.environment.Name {
			typ = InvalidType
			break
		}

		types := map[string]Type{}
		for i := range scope.Environment.UsedWalkers {
			if !scope.Environment.UsedWalkers[i].Walked {
				scope.Environment.UsedWalkers[i].Action(w.walkers)
			}
			typ := w.TypeExpr(typee, &scope.Environment.UsedWalkers[i].environment.Scope, false)
			if typ.PVT() != ast.Invalid {
				types[scope.Environment.UsedWalkers[i].environment.Name] = typ
			}
		}

		for k, v := range scope.Environment.UsedLibraries {
			if !v {
				continue
			}

			typ := w.TypeExpr(typee, &LibraryEnvs[k].Scope, false)
			if typ.PVT() != ast.Invalid {
				types[LibraryEnvs[k].Name] = typ
			}
		}

		// if len(types) > 1 {
		// 	errorMsg := "conflicting types between: "
		// 	for k, v := range types {
		// 		errorMsg += k + ":" + v.ToString() + ", "
		// 	}
		// 	errorMsg = errorMsg[:len(errorMsg)-1]
		// 	// w.Error(typee.GetToken(), errorMsg)
		// } else if len(types) == 1 {
		// 	for k, v := range types {
		// 		typee.Name = &ast.EnvAccessExpr{
		// 			PathExpr: &ast.EnvPathExpr{
		// 				Path: lexer.Token{
		// 					Lexeme:   k,
		// 					Location: typee.Name.GetToken().Location,
		// 				},
		// 			},
		// 			Accessed: &ast.IdentifierExpr{
		// 				Name: typee.Name.GetToken(),
		// 			},
		// 		}
		// 		return v
		// 	}
		// }
		for k, v := range types {
			typee.Name = &ast.EnvAccessExpr{
				PathExpr: &ast.EnvPathExpr{
					Path: tokens.Token{
						Lexeme:   k,
						Location: typee.Name.GetToken().Location,
					},
				},
				Accessed: &ast.IdentifierExpr{
					Name: typee.Name.GetToken(),
				},
			}
			return v
		}

		if len(scope.Environment.UsedLibraries) != 0 {
			if alias, found := BuiltinEnv.AliasTypes[typeName]; found {
				typ = alias.UnderlyingType

				break
			}
		}

		typ = InvalidType
	}

	if typee.IsVariadic {
		return NewVariadicType(typ)
	}
	if throw && typ.PVT() == ast.Invalid {
		// w.Error(typee.GetToken(), "invalid type")
	}
	return typ
}
