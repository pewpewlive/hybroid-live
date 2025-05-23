package walker

import (
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/generator"
	"hybroid/tokens"
)

func (w *Walker) StructExpr(node *ast.StructExpr, scope *Scope) *AnonStructVal {
	anonStructScope := NewScope(scope, &UntaggedTag{})
	structTypeVal := NewAnonStructVal(make(map[string]Field), false)

	for i := range node.Fields {
		w.fieldDeclaration(node.Fields[i], structTypeVal, anonStructScope)
	}

	return structTypeVal
}

func (w *Walker) FunctionExpr(fn *ast.FunctionExpr, scope *Scope) Value {
	returnTypes := w.GetReturns(fn.Returns, scope)

	funcTag := &FuncTag{ReturnTypes: returnTypes}
	fnScope := NewScope(scope, funcTag, ReturnAllowing)

	w.WalkBody(&fn.Body, funcTag, fnScope)

	params := make([]Type, 0)
	for i, param := range fn.Params {
		params = append(params, w.TypeExpr(param.Type, scope))
		variable := NewVariable(param.Name, w.TypeToValue(params[i]))
		w.DeclareVariable(fnScope, variable)
	}
	return &FunctionVal{
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
	matchScope.Tag = &MatchExprTag{YieldTypes: make([]Type, 0)}
	mpt := NewMultiPathTag(casesLength)

	w.GetNodeValue(&node.MatchStmt.ExprToMatch, scope)

	for i := range node.MatchStmt.Cases {
		caseScope := NewScope(matchScope, mpt)
		w.GetNodeValue(&node.MatchStmt.Cases[i].Expression, matchScope)
		w.WalkBody(&node.MatchStmt.Cases[i].Body, mpt, caseScope)
	}

	yieldTypes := matchScope.Tag.(*MatchExprTag).YieldTypes
	node.ReturnAmount = len(yieldTypes)

	if node.ReturnAmount == 0 {
		return &Invalid{}
	} else if node.ReturnAmount == 1 {
		return w.TypeToValue(yieldTypes[0])
	}

	return w.TypesToValues(yieldTypes)
}

func (w *Walker) EntityExpr(node *ast.EntityExpr, scope *Scope) Value {
	val := w.GetNodeValue(&node.Expr, scope)
	typ := w.TypeExpr(node.Type, scope)

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
	default:
		return &Invalid{}
	}
}

func ConvertNodeToAccessFieldExpr(ident ast.Node, index int, exprType ast.SelfExprType, envName string, entityName string) *ast.AccessExpr {
	fieldExpr := &ast.FieldExpr{
		Index:      index,
		Field:      ident,
		ExprType:   exprType,
		EnvName:    envName,
		EntityName: entityName,
	}

	return &ast.AccessExpr{
		Start: &ast.SelfExpr{
			Token: ident.GetToken(),
			Type:  exprType,
		},
		Accessed: []ast.Node{
			fieldExpr,
		},
	}
}

func ConvertCallToMethodCall(call *ast.CallExpr, exprType ast.SelfExprType, envName string, name string) *ast.MethodCallExpr {
	copy := *call
	return &ast.MethodCallExpr{
		EnvName:     envName,
		TypeName:    name,
		ExprType:    exprType,
		Caller:      copy.Caller,
		GenericArgs: copy.GenericArgs,
		Args:        copy.Args,
		MethodName:  call.Caller.GetToken().Lexeme,
	}
}

func (w *Walker) IdentifierExpr(node *ast.Node, scope *Scope) Value {
	valueNode := *node
	ident := valueNode.(*ast.IdentifierExpr)
	identToken := ident.GetToken()

	sc := w.ResolveVariable(scope, ident.Name)
	if sc == nil {
		walker, found := w.walkers[ident.Name.Lexeme]
		if found {
			*node = &ast.LiteralExpr{
				Value: "\"" + walker.environment.luaPath + "\"",
			}
			return NewPathVal(walker.environment.luaPath, walker.environment.Type)
		}

		w.AlertSingle(&alerts.UndeclaredVariableAccess{}, identToken, identToken.Lexeme)
		return &Invalid{}
	}

	variable := w.GetVariable(sc, ident.Name)
	if sc.Tag.GetType() == Struct {
		class := sc.Tag.(*ClassTag).Val
		field, index, found := class.ContainsField(variable.Name)

		*node = ConvertNodeToAccessFieldExpr(ident, index, ast.SelfStruct, class.Type.EnvName, "")

		if found {
			return field
		}
		method, found := class.Methods[variable.Name]
		if found {
			return method
		}
		w.AlertSingle(&alerts.MethodOrFieldNotFound{}, identToken, ident.Name.Lexeme)
	} else if sc.Tag.GetType() == Entity {
		entity := sc.Tag.(*EntityTag).EntityType
		field, index, found := entity.ContainsField(variable.Name)

		*node = ConvertNodeToAccessFieldExpr(ident, index, ast.SelfEntity, entity.Type.EnvName, entity.Type.Name)

		if found {
			return field
		}
		method, found := entity.Methods[variable.Name]
		if found {
			return method
		}
		w.AlertSingle(&alerts.MethodOrFieldNotFound{}, identToken, variable.Name)
	} else if sc.Environment.Name == "Builtin" {
		scope.Environment.AddBuiltinVar(ident.Name.Lexeme)
		*node = &ast.BuiltinExpr{
			Name: ident.Name,
		}
	} else if sc.Environment.Name != w.environment.Name && scope != sc {
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
	}

	for i := range walker.environment.importedWalkers {
		if walker.environment.importedWalkers[i].environment.Name == w.environment.Name {
			w.AlertSingle(&alerts.ImportCycle{}, node.GetToken(), w.environment.hybroidPath, walker.environment.hybroidPath)
			return &Invalid{}, nil
		}
	}

	if walker.environment.luaPath == "/dynamic/level.lua" {
		if !walker.Walked {
			walker.Pass2()
		}
		value := w.GetNodeValue(&node.Accessed, &walker.environment.Scope)
		return value, nil
	}

	walker.environment.AddRequirement(walker.environment.luaPath)

	if !walker.Walked {
		walker.Pass2()
	}

	value := w.GetNodeValue(&node.Accessed, &walker.environment.Scope)
	return value, nil
}

func (w *Walker) GroupingExpr(node *ast.GroupExpr, scope *Scope) Value {
	return w.GetNodeValue(&node.Expr, scope)
}

func (w *Walker) ListExpr(node *ast.ListExpr, scope *Scope) Value {
	value := &ListVal{}
	value.ValueType = w.GetContentsValueType(node.List, scope)
	return value
}

// Before calling it is assumed that the value of the caller is already gotten
func (w *Walker) CallExpr(val Value, node *ast.Node, scope *Scope) Value {
	call := (*node).(*ast.CallExpr)

	valType := val.GetType().PVT()
	if valType == ast.Invalid {
		return &Invalid{}
	}
	if valType != ast.Func {
		w.AlertSingle(&alerts.InvalidCallerType{}, call.GetToken(), valType)
		return &Invalid{}
	}

	if variable, ok := val.(*VariableVal); ok {
		val = variable.Value
	}

	fn := val.(*FunctionVal)

	nodeGenerics := call.GenericArgs
	nodeArgs := call.Args
	if fn.ProcType == Method {
		caller := call.Caller.(*ast.AccessExpr)
		field := caller.Accessed[len(caller.Accessed)].(*ast.FieldExpr)

		*node = ConvertCallToMethodCall(call, field.ExprType, field.EnvName, field.EntityName)
		mcall := (*node).(*ast.MethodCallExpr)

		nodeGenerics = mcall.GenericArgs
		nodeArgs = mcall.Args
	}

	genericArgs := w.GetGenerics(nodeGenerics, fn.Generics, scope)
	actualParams := make([]Type, 0)
	for i := range fn.Params {
		if fn.Params[i].GetType() == Generic {
			actualParams = append(actualParams, genericArgs[fn.Params[i].(*GenericType).Name])
		} else {
			actualParams = append(actualParams, fn.Params[i])
		}
	}
	args := []Type{}
	for i := range call.Args {
		args = append(args, w.GetNodeValue(&nodeArgs[i], scope).GetType())
	}
	w.ValidateArguments(genericArgs, args, actualParams, call)
	actualReturns := make([]Type, 0)
	for i := range fn.Returns {
		if fn.Returns[i].GetType() == Generic {
			actualReturns = append(actualReturns, genericArgs[fn.Returns[i].(*GenericType).Name])
		} else {
			actualReturns = append(actualReturns, fn.Returns[i])
		}
	}

	returnLen := len(actualReturns)
	if returnLen == 0 {
		return &Invalid{}
	} else if returnLen == 1 {
		return w.TypeToValue(actualReturns[returnLen-1])
	}

	return w.TypesToValues(actualReturns)
}

// Rewrote
func (w *Walker) AccessExpr(node *ast.AccessExpr, scope *Scope) Value {
	val := w.GetNodeActualValue(&node.Start, scope)

	prevNode := node.Start
	for i := range node.Accessed {
		valPVT := val.GetType().PVT()
		if valPVT == ast.Invalid {
			return &Invalid{}
		}

		scopedVal, scopeable := val.(ScopeableValue)

		if valPVT != ast.List && valPVT != ast.Map && !scopeable {
			w.AlertSingle(&alerts.InvalidAccessValue{}, node.GetToken(), valPVT)
			return &Invalid{}
		}
		token := node.Accessed[i].GetToken()
		exprType := node.Accessed[i].GetType()

		// list and map error handling
		if valPVT == ast.List || valPVT == ast.Map {
			if exprType == ast.FieldExpression {
				w.AlertSingle(&alerts.FieldAccessOnListOrMap{}, token,
					prevNode.GetToken().Lexeme,
					valPVT,
				)
				return &Invalid{}
			}

			member := node.Accessed[i].(*ast.MemberExpr).Member
			memberVal := w.GetNodeActualValue(&member, scope)
			if (memberVal.GetType().PVT() != ast.Number && valPVT == ast.List) ||
				(memberVal.GetType().PVT() != ast.String && valPVT == ast.Map) {

				w.AlertSingle(&alerts.InvalidMemberIndex{}, token,
					valPVT,
					member.GetToken().Lexeme,
				)
			}

			val = w.TypeToValue(val.GetType().(*WrapperType).WrappedType)
			prevNode = node.Accessed[i]
			continue
		}

		//struct, class, entity, enum error handling
		if scopeable {
			if exprType == ast.MemberExpression {
				w.AlertSingle(&alerts.MemberAccessOnNonListOrMap{}, token,
					prevNode.GetToken().Lexeme,
					valPVT,
				)
				return &Invalid{}
			}

			field := node.Accessed[i].(*ast.FieldExpr).Field
			fieldVal := w.GetNodeValue(&field, scopedVal.Scopify(scope))

			if fieldVal.GetType().PVT() == ast.Invalid {
				w.AlertSingle(&alerts.InvalidField{}, token,
					valPVT,
					token.Lexeme,
				)
			}

			val = fieldVal
			prevNode = node.Accessed[i]
			continue
		}
	}

	return val
}

// Rewrote
func (w *Walker) MapExpr(node *ast.MapExpr, scope *Scope) Value {
	mapVal := MapVal{}

	var currentType Type = InvalidType
	keymap := make(map[string]bool)
	for i := range node.KeyValueList {
		prop := node.KeyValueList[i]
		key := prop.Key.GetToken()

		if _, alreadyExists := keymap[key.Lexeme]; alreadyExists {
			w.AlertSingle(&alerts.DuplicateKeyInMap{}, key)
		} else {
			keymap[key.Lexeme] = true
		}

		memberVal := w.GetNodeValue(&prop.Expr, scope)
		memberType := memberVal.GetType()

		if i != 0 && !TypeEquals(currentType, memberType) {
			w.AlertSingle(&alerts.MixedMapOrListContents{}, prop.Expr.GetToken(),
				currentType.ToString(),
				memberType.ToString(),
			)
			currentType = InvalidType
			break
		}

		currentType = memberType
	}

	mapVal.MemberType = currentType
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

// Rewrote
func (w *Walker) SelfExpr(self *ast.SelfExpr, scope *Scope) Value {
	sc, _, classTag := ResolveTagScope[*ClassTag](scope)

	if sc == nil {
		entitySc, _, entityTag := ResolveTagScope[*EntityTag](scope)
		if entitySc != nil {
			self.Type = ast.SelfEntity
			self.EntityName = (*entityTag).EntityType.Type.Name
			return (*entityTag).EntityType
		}
		w.AlertSingle(&alerts.InvalidUseOfSelf{}, self.Token)
		return &Invalid{}
	}

	(*self).Type = ast.SelfStruct
	return (*classTag).Val
}

func (w *Walker) NewExpr(new *ast.NewExpr, scope *Scope) Value {
	_type := w.TypeExpr(new.Type, scope)

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

	suppliedGenerics := w.GetGenerics(new.GenericArgs, val.Generics, scope)

	w.ValidateArguments(suppliedGenerics, args, val.Params, new)

	return val
}

func (w *Walker) SpawnExpr(new *ast.SpawnExpr, scope *Scope) Value {
	typeExpr := w.TypeExpr(new.Type, scope)

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

	suppliedGenerics := w.GetGenerics(new.GenericArgs, val.SpawnGenerics, scope)

	w.ValidateArguments(suppliedGenerics, args, val.SpawnParams, new)

	return val
}

func (w *Walker) GetNodeValueFromExternalEnv(expr ast.Node, env *Environment) Value {
	val := w.GetNodeValue(&expr, &env.Scope)
	_, isValues := val.(*Values)
	if !isValues && val.GetType().PVT() == ast.Invalid {
		// w.Error(expr.GetToken(), fmt.Sprintf("variable named '%s' doesn't exist", expr.GetToken().Lexeme))
	}
	return val
}

func (w *Walker) TypeExpr(typee *ast.TypeExpr, scope *Scope) Type {
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
				walker.Pass2()
			}
			env = walker.environment
		} else {

			for _, v := range walker.environment.Requirements() {
				if v == w.environment.luaPath {
					w.AlertSingle(&alerts.ImportCycle{}, typee.GetToken(), w.environment.hybroidPath, walker.environment.hybroidPath)
					return InvalidType
				}
			}

			w.environment.AddRequirement(walker.environment.luaPath)

			if !walker.Walked {
				walker.Pass2()
			}
			env = walker.environment
		}

		ident := &ast.IdentifierExpr{Name: expr.Accessed.GetToken(), ValueType: ast.Invalid}
		typ = w.TypeExpr(&ast.TypeExpr{Name: ident}, &env.Scope)
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
				Value: w.TypeToValue(w.TypeExpr(v.Type, scope)),
				Token: v.Name,
			})
		}

		typ = NewStructType(fields, false)
	case ast.Func:
		params := []Type{}

		for _, v := range typee.Params {
			params = append(params, w.TypeExpr(v, scope))
		}

		returns := w.GetReturns(typee.Returns, scope)

		typ = &FunctionType{
			Params:  params,
			Returns: returns,
		}
	case ast.Map, ast.List:
		wrapped := w.TypeExpr(typee.WrappedType, scope)
		typ = NewWrapperType(NewBasicType(pvt), wrapped)
	case ast.Entity:
		typ = &RawEntityType{}
	default:
		typeName := typee.Name.GetToken().Lexeme
		if entityVal, found := scope.Environment.Entities[typeName]; found {
			typ = entityVal.GetType()
			break
		}
		if structVal, found := scope.Environment.Classes[typeName]; found {
			typ = structVal.GetType()
			break
		}
		if aliasType, found := scope.AliasTypes[typeName]; found {
			typ = aliasType.UnderlyingType

			break
		}
		if val, ok := scope.Environment.Scope.Variables[typeName]; ok {
			if val.GetType().PVT() == ast.Enum {
				typ = val.GetType()
				w.CheckAccessibility(scope, val.IsLocal, typee.Name.GetToken())
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
			if alias, found := BuiltinEnv.Scope.AliasTypes[typeName]; found {
				typ = alias.UnderlyingType
				break
			}
		}

		if scope.Environment.Name != w.environment.Name {
			typ = InvalidType
			break
		}

		types := map[string]Type{}
		for i := range scope.Environment.importedWalkers {
			if !scope.Environment.importedWalkers[i].Walked {
				scope.Environment.importedWalkers[i].Pass2()
			}
			typ := w.TypeExpr(typee, &scope.Environment.importedWalkers[i].environment.Scope)
			if typ.PVT() != ast.Invalid {
				types[scope.Environment.importedWalkers[i].environment.Name] = typ
			}
		}

		for k, v := range scope.Environment.UsedLibraries {
			if !v {
				continue
			}

			typ := w.TypeExpr(typee, &LibraryEnvs[k].Scope)
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
			if alias, found := BuiltinEnv.Scope.AliasTypes[typeName]; found {
				typ = alias.UnderlyingType
				break
			}
		}

		typ = InvalidType
	}

	if typee.IsVariadic {
		return NewVariadicType(typ)
	}
	return typ
}
