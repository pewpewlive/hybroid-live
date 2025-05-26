package walker

import (
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/tokens"
)

func (w *Walker) structExpression(node *ast.StructExpr, scope *Scope) *AnonStructVal {
	anonStructScope := NewScope(scope, &UntaggedTag{})
	structTypeVal := NewAnonStructVal(make(map[string]Field), false)

	for i := range node.Fields {
		w.fieldDeclaration(node.Fields[i], structTypeVal, anonStructScope)
	}

	return structTypeVal
}

func (w *Walker) functionExpression(fn *ast.FunctionExpr, scope *Scope) Value {
	generics := w.getGenericParams(fn.Generics)
	returnTypes := w.getReturns(fn.Returns, scope)
	funcTag := &FuncTag{Generics: generics, ReturnTypes: returnTypes}
	fnScope := NewScope(scope, funcTag, ReturnAllowing)

	w.walkBody(&fn.Body, funcTag, fnScope)

	params := make([]Type, 0)
	for i, param := range fn.Params {
		params = append(params, w.typeExpression(param.Type, scope))
		variable := NewVariable(param.Name, w.typeToValue(params[i]))
		w.declareVariable(fnScope, variable)
	}
	return &FunctionVal{
		Params:  params,
		Returns: returnTypes,
	}
}

func (w *Walker) matchExpression(node *ast.MatchExpr, scope *Scope) Value {
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
		w.walkBody(&node.MatchStmt.Cases[i].Body, mpt, caseScope)
	}

	yieldTypes := matchScope.Tag.(*MatchExprTag).YieldTypes
	node.ReturnAmount = len(yieldTypes)

	if node.ReturnAmount == 0 {
		return &Invalid{}
	} else if node.ReturnAmount == 1 {
		return w.typeToValue(yieldTypes[0])
	}

	return w.typesToValues(yieldTypes)
}

func (w *Walker) entityEvaluationExpression(node *ast.EntityEvaluationExpr, scope *Scope) Value {
	val := w.GetNodeValue(&node.Expr, scope)
	valType := val.GetType()
	typ := w.typeExpression(node.Type, scope)

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

	if valType.PVT() != ast.Entity {
		w.AlertSingle(&alerts.InvalidType{}, node.Expr.GetToken(), "entity", valType.ToString(), "in entity evaluation expression")
	}
	if !node.OfficialEntityType {
		if typ.GetType() != Named && typ.PVT() != ast.Entity {
			w.AlertSingle(&alerts.InvalidType{}, node.Type.GetToken(), "entity", typ.ToString(), "in entity evaluation expression")
			return &BoolVal{}
		}
		entityVal := w.typeToValue(typ).(*EntityVal)
		if node.ConvertedVarName != nil {
			w.context.EntityCasts.Push(NewEntityCast(*node.ConvertedVarName, entityVal))
		}
		node.EntityName = entityVal.Type.Name
		node.EnvName = entityVal.Type.EnvName
	} else if node.ConvertedVarName != nil {
		w.AlertSingle(&alerts.OfficialEntityConversion{}, *node.ConvertedVarName)
	}

	return &BoolVal{}
}

func (w *Walker) binaryExpression(node *ast.BinaryExpr, scope *Scope) Value {
	left, right := w.GetNodeValue(&node.Left, scope), w.GetNodeValue(&node.Right, scope)
	leftType, rightType := left.GetType(), right.GetType()
	op := node.Operator
	switch op.Type {
	case tokens.Plus, tokens.Minus, tokens.Caret, tokens.Star, tokens.Slash, tokens.Modulo, tokens.BackSlash:
		w.validateArithmeticOperands(leftType, rightType, node)
		typ := w.determineValueType(leftType, rightType)

		if typ.PVT() == ast.Invalid {
			return &Invalid{}
		}
		return w.typeToValue(typ)
	case tokens.Concat:
		if !TypeEquals(leftType, NewBasicType(ast.String)) && !TypeEquals(rightType, NewBasicType(ast.String)) {
			return &Invalid{}
		}
		return &StringVal{}
	default: // comparison
		if op.Type == tokens.Or {
			var operand ast.Node
			if node.Left.GetType() == ast.EntityExpression {
				operand = node.Left
			} else if node.Right.GetType() == ast.EntityExpression {
				operand = node.Right
			}
			if operand != nil && operand.(*ast.EntityEvaluationExpr).ConvertedVarName != nil {
				w.AlertSingle(&alerts.EntityConversionWithOrCondition{}, operand.GetToken())
				return &Invalid{}
			}
		}

		if !TypeEquals(leftType, rightType) {
			w.AlertSingle(&alerts.TypeMismatch{}, node.Left.GetToken(),
				leftType.ToString(),
				rightType.ToString(),
				"in binary expression",
			)
			return &Invalid{}
		}
		return &BoolVal{}
	}
}

func (w *Walker) literalExpression(node *ast.LiteralExpr) Value {
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

func (w *Walker) identifierExpression(node *ast.Node, scope *Scope) Value {
	valueNode := *node
	ident := valueNode.(*ast.IdentifierExpr)
	identToken := ident.GetToken()

	sc := w.resolveVariable(scope, ident.Name)
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

	variable := w.getVariable(sc, ident.Name)
	if sc.Tag.GetType() == Class {
		class := sc.Tag.(*ClassTag).Val
		field, index, found := class.ContainsField(variable.Name)

		*node = convertNodeToAccessFieldExpr(ident, index, ast.SelfClass, class.Type.EnvName, "")

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

		*node = convertNodeToAccessFieldExpr(ident, index, ast.SelfEntity, entity.Type.EnvName, entity.Type.Name)

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

	variable.IsUsed = true
	return variable
}

func (w *Walker) environmentAccessExpression(node *ast.EnvAccessExpr) (Value, ast.Node) {
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
		return w.GetNodeValue(&node.Accessed, &PewpewEnv.Scope), nil
	case "Fmath":
		return w.GetNodeValue(&node.Accessed, &FmathEnv.Scope), nil
	case "Math":
		if w.environment.Type == ast.LevelEnv {
			w.AlertSingle(&alerts.UnallowedLibraryUse{}, node.PathExpr.Path, "Math", "Level")
		}
		return w.GetNodeValue(&node.Accessed, &MathEnv.Scope), nil
	case "String":
		return w.GetNodeValue(&node.Accessed, &StringEnv.Scope), nil
	case "Table":
		return w.GetNodeValue(&node.Accessed, &TableEnv.Scope), nil
	}

	walker, found := w.walkers[envName]
	if !found {
		w.AlertSingle(&alerts.InvalidEnvironmentAccess{}, node.PathExpr.GetToken())
		return &Invalid{}, nil
	}

	if walker.environment.Name == w.environment.Name {
		w.AlertSingle(&alerts.EnvironmentAccessToItself{}, node.PathExpr.GetToken())
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
			walker.Walk()
		}
		value := w.GetNodeValue(&node.Accessed, &walker.environment.Scope)
		return value, nil
	}

	walker.environment.AddRequirement(walker.environment.luaPath)

	if !walker.Walked {
		walker.Walk()
	}

	value := w.GetNodeValue(&node.Accessed, &walker.environment.Scope)
	return value, nil
}

func (w *Walker) groupExpression(node *ast.GroupExpr, scope *Scope) Value {
	return w.GetNodeValue(&node.Expr, scope)
}

func (w *Walker) listExpression(node *ast.ListExpr, scope *Scope) Value {
	value := &ListVal{}
	value.ValueType = w.getContentsValueType(node.List, scope)
	return value
}

// Before calling it is assumed that the value of the caller is already gotten
func (w *Walker) callExpression(val Value, node *ast.Node, scope *Scope) Value {
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

		*node = convertCallToMethodCall(call, field.ExprType, field.EnvName, field.EntityName)
		mcall := (*node).(*ast.MethodCallExpr)

		nodeGenerics = mcall.GenericArgs
		nodeArgs = mcall.Args
	}

	genericArgs := w.getGenerics(nodeGenerics, fn.Generics, scope)
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
	w.validateArguments(genericArgs, args, actualParams, call)
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
		return w.typeToValue(actualReturns[returnLen-1])
	}

	return w.typesToValues(actualReturns)
}

// Rewrote
func (w *Walker) accessExpression(node *ast.AccessExpr, scope *Scope) Value {
	val := w.GetActualNodeValue(&node.Start, scope)

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
			memberVal := w.GetActualNodeValue(&member, scope)
			if (memberVal.GetType().PVT() != ast.Number && valPVT == ast.List) ||
				(memberVal.GetType().PVT() != ast.String && valPVT == ast.Map) {

				w.AlertSingle(&alerts.InvalidMemberIndex{}, token,
					valPVT,
					member.GetToken().Lexeme,
				)
			}

			val = w.typeToValue(val.GetType().(*WrapperType).WrappedType)
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
func (w *Walker) mapExpression(node *ast.MapExpr, scope *Scope) Value {
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

func (w *Walker) unaryExpression(node *ast.UnaryExpr, scope *Scope) Value {
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
func (w *Walker) selfExpression(self *ast.SelfExpr, scope *Scope) Value {
	sc, _, classTag := resolveTagScope[*ClassTag](scope)

	if sc == nil {
		entitySc, _, entityTag := resolveTagScope[*EntityTag](scope)
		if entitySc != nil {
			self.Type = ast.SelfEntity
			self.EntityName = (*entityTag).EntityType.Type.Name
			return (*entityTag).EntityType
		}
		w.AlertSingle(&alerts.InvalidUseOfSelf{}, self.Token)
		return &Invalid{}
	}

	(*self).Type = ast.SelfClass
	return (*classTag).Val
}

func (w *Walker) newExpression(new *ast.NewExpr, scope *Scope) Value {
	_type := w.typeExpression(new.Type, scope)

	if _type.PVT() == ast.Invalid {
		return &Invalid{}
	}
	if _type.PVT() != ast.Class {
		w.AlertSingle(&alerts.InvalidType{}, new.Type.GetToken(), "class", _type.ToString(), "in new expression")
		return &Invalid{}
	}

	val := w.typeToValue(_type).(*ClassVal)

	args := make([]Type, 0)
	for i := range new.Args {
		args = append(args, w.GetNodeValue(&new.Args[i], scope).GetType())
	}
	suppliedGenerics := w.getGenerics(new.GenericArgs, val.Generics, scope)
	w.validateArguments(suppliedGenerics, args, val.Params, new)

	return val
}

func (w *Walker) spawnExpression(new *ast.SpawnExpr, scope *Scope) Value {
	_type := w.typeExpression(new.Type, scope)

	if _type.PVT() == ast.Invalid {
		return &Invalid{}
	}
	if _type.PVT() != ast.Entity {
		w.AlertSingle(&alerts.InvalidType{}, new.Type.GetToken(), "entity", _type.ToString(), "in spawn expression")
		return &Invalid{}
	}

	val := w.typeToValue(_type).(*EntityVal)

	args := make([]Type, 0)
	for i := range new.Args {
		args = append(args, w.GetNodeValue(&new.Args[i], scope).GetType())
	}
	suppliedGenerics := w.getGenerics(new.GenericArgs, val.SpawnGenerics, scope)
	w.validateArguments(suppliedGenerics, args, val.SpawnParams, new)

	return val
}

func (w *Walker) typeExpression(typee *ast.TypeExpr, scope *Scope) Type {
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
				walker.Walk()
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
				walker.Walk()
			}
			env = walker.environment
		}

		ident := &ast.IdentifierExpr{Name: expr.Accessed.GetToken(), ValueType: ast.Invalid}
		typ = w.typeExpression(&ast.TypeExpr{Name: ident}, &env.Scope)
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

	pvt := w.getTypeFromString(typee.Name.GetToken().Lexeme)
	switch pvt {
	case ast.Bool, ast.String, ast.Number:
		typ = NewBasicType(pvt)
	case ast.Fixed, ast.FixedPoint, ast.Radian, ast.Degree:
		typ = NewFixedPointType(pvt)
	case ast.Enum:
		typ = NewBasicType(ast.Enum)
	case ast.Struct:
		fields := []*VariableVal{}

		for _, v := range typee.Fields {
			fields = append(fields, &VariableVal{
				Name:  v.Name.Lexeme,
				Value: w.typeToValue(w.typeExpression(v.Type, scope)),
				Token: v.Name,
			})
		}

		typ = NewStructType(fields, false)
	case ast.Func:
		params := []Type{}

		for _, v := range typee.Params {
			params = append(params, w.typeExpression(v, scope))
		}

		returns := w.getReturns(typee.Returns, scope)

		typ = &FunctionType{
			Params:  params,
			Returns: returns,
		}
	case ast.Map, ast.List:
		wrapped := w.typeExpression(typee.WrappedType, scope)
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
				w.checkAccessibility(scope, val.IsLocal, typee.Name.GetToken())
				break
			}
		}

		sc, _, fnTag := resolveTagScope[*FuncTag](scope)

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
				scope.Environment.importedWalkers[i].Walk()
			}
			typ := w.typeExpression(typee, &scope.Environment.importedWalkers[i].environment.Scope)
			if typ.PVT() != ast.Invalid {
				types[scope.Environment.importedWalkers[i].environment.Name] = typ
			}
		}

		for k, v := range scope.Environment.UsedLibraries {
			if !v {
				continue
			}

			typ := w.typeExpression(typee, &LibraryEnvs[k].Scope)
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
