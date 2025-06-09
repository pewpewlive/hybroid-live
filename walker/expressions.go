package walker

import (
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/generator/mapping"
	"hybroid/tokens"
	"strconv"
)

func (w *Walker) structExpression(node *ast.StructExpr, scope *Scope) *StructVal {
	structTypeVal := NewStructVal(make(map[string]Field), false)

	for i := range node.Fields {
		fieldToken := node.Fields[i].Name
		val := w.GetActualNodeValue(&node.Expressions[i], scope)
		if field, found := structTypeVal.Fields[fieldToken.Lexeme]; found {
			w.AlertSingle(&alerts.Redeclaration{}, fieldToken, field.Var.Name, "struct field")
			continue
		}
		if _, ok := val.(Values); ok {
			w.AlertSingle(&alerts.InvalidType{}, node.Expressions[i].GetToken(), val.GetType(), "in struct field declaration")
		}
		structTypeVal.Fields[fieldToken.Lexeme] = NewField(0, NewVariable(fieldToken, val))
	}

	return structTypeVal
}

func (w *Walker) functionExpression(fn *ast.FunctionExpr, scope *Scope) Value {
	generics := w.getGenericParams(fn.Generics, scope)
	returnTypes := w.getReturns(fn.Returns, scope)
	funcTag := &FuncTag{Generics: generics, ReturnTypes: returnTypes}
	fnScope := NewScope(scope, funcTag, ReturnAllowing)
	params := w.getParameters(fn.Params, fnScope)

	w.walkFuncBody(fn, &fn.Body, funcTag, fnScope)

	return &FunctionVal{
		Params:  params,
		Returns: returnTypes,
	}
}

func (w *Walker) matchExpression(node *ast.MatchExpr, scope *Scope) Value {
	matchStmt := node.MatchStmt

	cases := matchStmt.Cases
	casesLength := len(cases)
	if !matchStmt.HasDefault {
		w.AlertSingle(&alerts.DefaultCaseMissing{}, matchStmt.Token)
		casesLength++
		if casesLength < 1 {
			w.AlertSingle(&alerts.InsufficientCases{}, matchStmt.Token)
		}
	} else if casesLength < 2 {
		w.AlertSingle(&alerts.InsufficientCases{}, matchStmt.Token)
	}
	matchScope := NewScope(scope, &MatchExprTag{YieldTypes: make([]Type, 0)}, YieldAllowing)
	mpt := NewMultiPathTag(casesLength, YieldAllowing)

	valToMatch := w.GetActualNodeValue(&matchStmt.ExprToMatch, scope)
	valType := valToMatch.GetType()

	for i := range cases {
		caseScope := NewScope(matchScope, mpt)

		w.walkBody(&matchStmt.Cases[i].Body, mpt, caseScope)

		if matchStmt.Cases[i].Expressions[0].GetToken().Lexeme == "else" {
			if i != len(matchStmt.Cases)-1 {
				w.AlertSingle(&alerts.InvalidDefaultCasePlacement{}, matchStmt.Cases[i].Expressions[0].GetToken(), "in match expression")
			}
			continue
		}

		for j := range matchStmt.Cases[i].Expressions {
			caseValType := w.GetNodeValue(&matchStmt.Cases[i].Expressions[j], scope).GetType()
			if valType == InvalidType || caseValType == InvalidType {
				continue
			}
			if !TypeEquals(valType, caseValType) {
				w.AlertSingle(&alerts.InvalidCaseType{}, matchStmt.Cases[i].Expressions[j].GetToken(), valType, caseValType)
			}
		}
	}

	if !mpt.GetIfExits(Yield) {
		w.AlertMulti(&alerts.NotAllCodePathsExit{},
			cases[0].Expressions[0].GetToken(),
			cases[len(cases)-1].Expressions[0].GetToken(),
			"yield",
		)
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
			"Pointonium", "BonusImplosion", "Mace", "PlasmaField":
			typ = &RawEntityType{}
			node.OfficialEntityType = true
		}
	}

	if valType.PVT() != ast.Entity {
		w.AlertSingle(&alerts.TypeMismatch{}, node.Expr.GetToken(), "entity", valType.String(), "in entity evaluation expression")
	}
	if !node.OfficialEntityType {
		if !(typ.GetType() == Named && typ.PVT() == ast.Entity) {
			w.AlertSingle(&alerts.TypeMismatch{}, node.Type.GetToken(), "entity", typ.String(), "in entity evaluation expression")
			return NewBoolVal("false")
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
	left, right := w.GetActualNodeValue(&node.Left, scope), w.GetActualNodeValue(&node.Right, scope)
	leftType, rightType := left.GetType(), right.GetType()
	op := node.Operator
	switch op.Type {
	case tokens.Plus, tokens.Minus, tokens.Caret, tokens.Star, tokens.Slash, tokens.Modulo, tokens.BackSlash:
		return w.validateArithmeticOperands(left, right, node, "in arithmetic expression")
	case tokens.Concat:
		if leftType.PVT() != ast.Text {
			w.AlertSingle(&alerts.TypeMismatch{}, node.Left.GetToken(), "string", leftType, "in concatenation")
		}
		if rightType.PVT() != ast.Text {
			w.AlertSingle(&alerts.TypeMismatch{}, node.Right.GetToken(), "string", rightType, "in concatenation")
		}
		return &StringVal{}
	case tokens.Greater, tokens.GreaterEqual, tokens.Less, tokens.LessEqual, tokens.BangEqual, tokens.EqualEqual:
		if leftType == InvalidType || rightType == InvalidType {
			return &BoolVal{}
		}
		if !TypeEquals(leftType, rightType) {
			w.AlertSingle(&alerts.TypesMismatch{}, node.Left.GetToken(), "left value", leftType, "right value", rightType)
		}
		return &BoolVal{}
	case tokens.Pipe, tokens.Ampersand, tokens.LeftShift, tokens.RightShift, tokens.Tilde:
		if leftType == InvalidType || rightType == InvalidType {
			return &Invalid{}
		}
		if leftType.PVT() != ast.Number {
			w.AlertSingle(&alerts.TypeMismatch{}, node.Left.GetToken(), "number", leftType, "in bitwise expression")
		}
		if rightType.PVT() != ast.Number {
			w.AlertSingle(&alerts.TypeMismatch{}, node.Right.GetToken(), "number", rightType, "in bitwise expression")
		}
		return &NumberVal{}
	default: // logical comparison
		if op.Type == tokens.Or {
			var operand ast.Node
			if node.Left.GetType() == ast.EntityEvaluationExpression {
				operand = node.Left
			} else if node.Right.GetType() == ast.EntityEvaluationExpression {
				operand = node.Right
			}
			if operand != nil && operand.(*ast.EntityEvaluationExpr).ConvertedVarName != nil {
				w.AlertSingle(&alerts.EntityConversionWithOrCondition{}, operand.GetToken())
				return &BoolVal{}
			}
		}

		return w.validateConditionalOperands(left, right, node)
	}
}

func (w *Walker) literalExpression(node *ast.LiteralExpr) Value {
	switch node.Token.Type {
	case tokens.String:
		return &StringVal{}
	case tokens.Fixed, tokens.Radian, tokens.FixedPoint, tokens.Degree:
		return &FixedVal{}
	case tokens.True, tokens.False:
		return NewBoolVal(node.Value)
	case tokens.Number:
		return NewNumberVal(node.Value)
	default:
		return &Invalid{}
	}
}

func (w *Walker) identifierExpression(node *ast.Node, scope *Scope) Value {
	valueNode := *node
	ident := valueNode.(*ast.IdentifierExpr)
	identToken := ident.GetToken()

	sc := w.resolveVariable(scope, ident.Name)
check:
	if sc == nil {

		walker, found := w.walkers[ident.Name.Lexeme]
		if found {
			*node = &ast.LiteralExpr{
				Value: "\"" + walker.environment.luaPath + "\"",
				Token: ident.Name,
			}
			return NewPathVal(walker.environment.luaPath, walker.environment.Type, walker.environment.Name)
		}

		w.AlertSingle(&alerts.UndeclaredVariableAccess{}, identToken, identToken.Lexeme)
		return &Invalid{}
	}

	if (sc.Tag.GetType() == Class || sc.Tag.GetType() == Entity) && !scope.Is(SelfAllowing) {
		sc = w.resolveVariable(scope.Parent, ident.Name)
		goto check
	}

	variable := w.getVariable(sc, ident.Name)
	if val, ok := sc.ConstValues[variable.Name]; variable.IsConst && ok {
		*node = val
		return variable
	}
	if sc.Tag.GetType() == Class {
		class := sc.Tag.(*ClassTag).Val
		field, index, found := class.ContainsField(variable.Name)

		*node = &ast.AccessExpr{
			Start: &ast.SelfExpr{
				Token: ident.GetToken(),
				Type:  ast.ClassMethod,
			},
			Accessed: []ast.Node{
				&ast.FieldExpr{
					Index: index,
					Field: ident,
				},
			},
		}

		if found {
			return field
		}
		method, found := class.Methods[variable.Name]
		if found {
			return method
		}
		w.AlertSingle(&alerts.MethodOrFieldNotFound{}, identToken, ident.Name.Lexeme)
	} else if sc.Tag.GetType() == Entity && scope.Is(SelfAllowing) {
		entity := sc.Tag.(*EntityTag).EntityVal
		field, index, found := entity.ContainsField(variable.Name)

		*node = &ast.AccessExpr{
			Start: &ast.SelfExpr{
				Token:      ident.GetToken(),
				EntityName: entity.Type.Name,
				Type:       ast.EntityMethod,
			},
			Accessed: []ast.Node{
				&ast.FieldExpr{
					Index: index,
					Field: ident,
				},
			},
		}

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
		ident.Type = ast.Raw
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
		ident := node.Accessed.(*ast.IdentifierExpr)
		path := envName + ":" + ident.Name.Lexeme
		walker, found := w.walkers[path]
		if found {
			return NewPathVal(walker.environment.luaPath, walker.environment.Type, walker.environment.Name), &ast.LiteralExpr{
				Value: "\"" + walker.environment.luaPath + "\"",
			}
		}
	}

	switch envName {
	case "Pewpew":
		return w.GetNodeValue(&node.Accessed, &PewpewAPI.Scope), nil
	case "Fmath":
		return w.GetNodeValue(&node.Accessed, &FmathAPI.Scope), nil
	case "Math":
		if w.environment.Type == ast.LevelEnv {
			w.AlertSingle(&alerts.UnallowedLibraryUse{}, node.PathExpr.Path, "Math", "Level")
		}
		return w.GetNodeValue(&node.Accessed, &MathAPI.Scope), nil
	case "String":
		return w.GetNodeValue(&node.Accessed, &StringAPI.Scope), nil
	case "Table":
		return w.GetNodeValue(&node.Accessed, &TableAPI.Scope), nil
	}

	walker, found := w.walkers[envName]
	if !found {
		w.AlertSingle(&alerts.InvalidEnvironmentAccess{}, node.PathExpr.GetToken(), envName)
		return &Invalid{}, nil
	}

	if walker.environment.Name == w.environment.Name {
		w.AlertSingle(&alerts.EnvironmentAccessToItself{}, node.PathExpr.GetToken())
		return &Invalid{}, nil
	}

	if walker.environment.Type != ast.SharedEnv && (w.environment.Type == ast.MeshEnv || w.environment.Type == ast.SoundEnv) {
		w.AlertSingle(&alerts.UnallowedEnvironmentAccess{}, node.PathExpr.GetToken(), "non Shared", "Mesh or Sound")
		return &Invalid{}, nil
	} else if w.environment.Type == ast.LevelEnv && (walker.environment.Type == ast.MeshEnv || walker.environment.Type == ast.SoundEnv) {
		w.AlertSingle(&alerts.UnallowedEnvironmentAccess{}, node.PathExpr.GetToken(), "Mesh or Sound", "Level")
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

	w.environment.AddRequirement(walker.environment.luaPath)

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

	valType := val.GetType()
	if valType == InvalidType {
		return &Invalid{}
	}
	if valType.PVT() != ast.Func {
		w.AlertSingle(&alerts.InvalidCallerType{}, call.GetToken(), valType)
		return &Invalid{}
	}

	if variable, ok := val.(*VariableVal); ok {
		val = variable.Value
	}

	fn := *val.(*FunctionVal)

	nodeGenerics := call.GenericArgs
	nodeArgs := call.Args
	if fn.ProcType == Method {
		caller := call.Caller.(*ast.AccessExpr)
		caller.Accessed = caller.Accessed[:len(caller.Accessed)-1]
		if len(caller.Accessed) == 0 {
			call.Caller = caller.Start
		}

		*node = convertCallToMethodCall(call, fn.MethodInfo)
		mcall := (*node).(*ast.MethodCallExpr)

		nodeGenerics = mcall.GenericArgs
		nodeArgs = mcall.Args
	}

	genericArgs := w.getGenerics(nodeGenerics, fn.Generics, scope)
	args := []Value{}
	for i := range call.Args {
		args = append(args, w.GetActualNodeValue(&nodeArgs[i], scope))
	}
	w.validateArguments(genericArgs, args, &fn, call)

	actualReturns := fn.Returns
	returnLen := len(actualReturns)
	call.ReturnAmount = returnLen
	if returnLen == 0 {
		return &Invalid{}
	} else if returnLen == 1 {
		return w.typeToValue(actualReturns[returnLen-1])
	}

	return w.typesToValues(actualReturns)
}

// Rewrote
func (w *Walker) accessExpression(_node *ast.Node, scope *Scope) Value {
	node := (*_node).(*ast.AccessExpr)
	var val Value

	typeExpr := &ast.TypeExpr{Name: node.Start}
	typ := w.typeExpression(typeExpr, scope)
	if et, ok := typ.(*EnumType); ok {
		val = w.typeToValue(et)
		node.Start = typeExpr.Name
	} else if typ == UnknownTyp {
		val = w.GetActualNodeValue(&node.Start, scope)
	} else {
		val = &Invalid{}
	}

	prevNode := &node.Start
	for i := range node.Accessed {
		valType := val.GetType()
		if valType == InvalidType {
			return &Invalid{}
		}

		scopedVal, scopeable := val.(ScopeableValue)

		if valType.GetType() != Wrapper && !scopeable {
			w.AlertSingle(&alerts.InvalidAccessValue{}, (*prevNode).GetToken(), valType)
			return &Invalid{}
		}
		token := node.Accessed[i].GetToken()
		exprType := node.Accessed[i].GetType()

		// list and map error handling
		if valType.GetType() == Wrapper {
			if exprType == ast.FieldExpression {
				w.AlertSingle(&alerts.FieldAccessOnListOrMap{}, token,
					node.Accessed[i].GetToken().Lexeme,
					valType,
				)
				return &Invalid{}
			}

			member := node.Accessed[i].(*ast.MemberExpr).Member
			memberVal := w.GetActualNodeValue(&member, scope)
			if (memberVal.GetType().PVT() != ast.Number && valType.PVT() == ast.List) ||
				(memberVal.GetType().PVT() != ast.Text && valType.PVT() == ast.Map) {

				w.AlertSingle(&alerts.InvalidMemberIndex{}, token,
					valType,
					member.GetToken().Lexeme,
				)
			}
			if memberVal.GetType().PVT() == ast.Number && valType.PVT() == ast.List {
				num := memberVal.(*NumberVal)
				if num.Value != "unknown" {
					n, err := strconv.ParseFloat(num.Value, 64)
					if err == nil && n < 1 {
						w.AlertSingle(&alerts.ListIndexOutOfBounds{}, member.GetToken(), num.Value)
					} else if err == nil && n != float64(int64(n)) {
						w.AlertSingle(&alerts.InvalidListIndex{}, member.GetToken())
					}
				}
			}

			val = w.typeToValue(val.GetType().(*WrapperType).WrappedType)
			prevNode = &node.Accessed[i]
			continue
		}

		//struct, class, entity, enum error handling
		if scopeable {
			if exprType == ast.MemberExpression {
				w.AlertSingle(&alerts.MemberAccessOnNonListOrMap{}, token,
					node.Accessed[i].GetToken().Lexeme,
					valType,
				)
				return &Invalid{}
			}

			field := node.Accessed[i].(*ast.FieldExpr)
			fieldVal := w.GetNodeValue(&field.Field, scopedVal.Scopify(scope))

			if _, found := fieldVal.(*VariableVal); !found {
				w.AlertSingle(&alerts.InvalidField{}, token,
					(*prevNode).GetToken().Lexeme,
					token.Lexeme,
				)
				return &Invalid{}
			}
			innerVal := fieldVal.(*VariableVal).Value

			fc := val.(FieldContainer)
			_, index, found := fc.ContainsField(field.GetToken().Lexeme)
			if found && valType.GetType() == Named {
				field.Index = index
			}
			ok2 := true
			if fn, ok := innerVal.(*FunctionVal); ok && fn.ProcType == Method {
				ok2 = false
			}
			if entityVal, ok := val.(*EntityVal); ok && ok2 && (*prevNode).GetType() != ast.SelfExpression {
				*prevNode = &ast.EntityAccessExpr{
					Expr:       *prevNode,
					EntityName: entityVal.Type.Name,
					EnvName:    entityVal.Type.EnvName,
				}
			}

			val = fieldVal
			prevNode = &node.Accessed[i]
		}
	}

	if _, ok := val.(*VariableVal); !ok {
		val = NewVariable((*prevNode).GetToken(), val)
	}

	// check if we got an enum variant and convert that to its constant value
	if enumVal, ok := val.(*VariableVal).Value.(*EnumFieldVal); ok {
		if enumVal.Type.EnvName == "Pewpew" {
			ident := node.Accessed[0].(*ast.FieldExpr).Field.(*ast.IdentifierExpr)
			ident.Name.Lexeme = mapping.PewpewEnums[enumVal.Type.Name][ident.Name.Lexeme]
			return val
		}
		*_node = &ast.LiteralExpr{
			Value: strconv.Itoa(enumVal.Index),
			Token: node.GetToken(),
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
			w.AlertSingle(&alerts.DuplicateElement{}, key, "map key", key.Lexeme)
		} else {
			keymap[key.Lexeme] = true
		}

		memberVal := w.GetNodeValue(&prop.Expr, scope)
		memberType := memberVal.GetType()

		if i != 0 && !TypeEquals(currentType, memberType) {
			w.AlertSingle(&alerts.MixedMapOrListContents{}, prop.Expr.GetToken(),
				currentType.String(),
				memberType.String(),
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

	if valPVT == ast.Invalid {
		return val
	}

	token := node.Value.GetToken()

	switch node.Operator.Type {
	case tokens.Bang:
		if valPVT != ast.Bool {
			w.AlertSingle(&alerts.TypeMismatch{}, token, "bool", valType.String(), "after '!' in unary expression")
		}
	case tokens.Hash:
		if !(valType.GetType() == Wrapper && valType.(*WrapperType).Type.PVT() == ast.List) {
			w.AlertSingle(&alerts.TypeMismatch{}, token, "list", valType.String(), "after '#' in unary expression")
		}
		return &NumberVal{}
	case tokens.Minus:
		if !isNumerical(valPVT) {
			w.AlertSingle(&alerts.TypeMismatch{}, token, "a numerical type", valType.String(), "after '-' in unary expression")
		}
	}

	return val
}

// Rewrote
func (w *Walker) selfExpression(self *ast.SelfExpr, scope *Scope) Value {
	if !scope.Is(SelfAllowing) {
		w.AlertSingle(&alerts.InvalidUseOfSelf{}, self.Token)
		return &Invalid{}
	}
	sc, _, classTag := resolveTagScope[*ClassTag](scope)

	if sc == nil {
		entitySc, _, entityTag := resolveTagScope[*EntityTag](scope)
		if entitySc != nil {
			self.Type = ast.EntityMethod
			self.EntityName = (*entityTag).EntityVal.Type.Name
			return (*entityTag).EntityVal
		}
		w.AlertSingle(&alerts.InvalidUseOfSelf{}, self.Token)
		return &Invalid{}
	}

	(*self).Type = ast.ClassMethod
	return (*classTag).Val
}

func (w *Walker) newExpression(new *ast.NewExpr, scope *Scope) Value {
	_type := w.typeExpression(new.Type, scope)

	if _type == InvalidType {
		return &Invalid{}
	}
	if _type.PVT() != ast.Class {
		w.AlertSingle(&alerts.TypeMismatch{}, new.Type.GetToken(), "class", _type.String(), "in new expression")
		return &Invalid{}
	}
	val := w.typeToValue(_type).(*ClassVal)

	args := make([]Value, 0)
	for i := range new.Args {
		args = append(args, w.GetActualNodeValue(&new.Args[i], scope))
	}

	explicitClassGenericArgs := w.getGenerics(new.ClassGenericArgs, val.Generics, scope)
	explicitGenericArgs := w.getGenerics(new.GenericArgs, val.New.Generics, scope)
	for k, v := range explicitClassGenericArgs {
		explicitGenericArgs[k] = v
	}

	w.validateArguments(explicitGenericArgs, args, val.New, new)
	for _, v := range val.Generics {
		generic := explicitGenericArgs[v.Name]
		if generic.Type == UnknownTyp {
			continue
		}
		for i, v := range val.Fields {
			if v.Var.Value.GetType().GetType() == Generic {
				val.Fields[i].Var.Value = w.typeToValue(generic.Type)
			}
		}
		for i := range val.Methods {
			fn := val.Methods[i].Value.(*FunctionVal)
			for j, v3 := range fn.Params {
				if gen, ok := v3.(*GenericType); ok && gen.Name == v.Name {
					fn.Params[j] = generic.Type
				}
			}
			for j, v3 := range fn.Returns {
				if gen, ok := v3.(*GenericType); ok && gen.Name == v.Name {
					fn.Returns[j] = generic.Type
				}
			}
		}
		val.Type.Generics = append(val.Type.Generics, generic)
	}

	new.EnvName = val.Type.EnvName
	return val
}

func (w *Walker) spawnExpression(new *ast.SpawnExpr, scope *Scope) Value {
	_type := w.typeExpression(new.Type, scope)

	if _type == InvalidType {
		return &Invalid{}
	}
	if _type.PVT() != ast.Entity {
		w.AlertSingle(&alerts.TypeMismatch{}, new.Type.GetToken(), "entity", _type.String(), "in spawn expression")
		return &Invalid{}
	}
	val := w.typeToValue(_type).(*EntityVal)

	args := make([]Value, 0)
	for i := range new.Args {
		args = append(args, w.GetActualNodeValue(&new.Args[i], scope))
	}

	explicitEntityGenericArgs := w.getGenerics(new.EntityGenericArgs, val.Generics, scope)
	explicitGenericArgs := w.getGenerics(new.GenericArgs, val.Spawn.Generics, scope)
	for k, v := range explicitEntityGenericArgs {
		explicitGenericArgs[k] = v
	}

	w.validateArguments(explicitGenericArgs, args, val.Spawn, new)
	for _, v := range val.Generics {
		generic := explicitGenericArgs[v.Name]
		if generic.Type == UnknownTyp {
			continue
		}
		for i, v := range val.Fields {
			if v.Var.Value.GetType().GetType() == Generic {
				val.Fields[i].Var.Value = w.typeToValue(generic.Type)
			}
		}
		for i := range val.Methods {
			fn := val.Methods[i].Value.(*FunctionVal)
			for j, v3 := range fn.Params {
				if gen, ok := v3.(*GenericType); ok && gen.Name == v.Name {
					fn.Params[j] = generic.Type
				}
			}
			for j, v3 := range fn.Returns {
				if gen, ok := v3.(*GenericType); ok && gen.Name == v.Name {
					fn.Returns[j] = generic.Type
				}
			}
		}
		fn := val.Destroy
		for j, v3 := range fn.Params {
			if gen, ok := v3.(*GenericType); ok && gen.Name == v.Name {
				fn.Params[j] = generic.Type
			}
		}
		for j, v3 := range fn.Returns {
			if gen, ok := v3.(*GenericType); ok && gen.Name == v.Name {
				fn.Returns[j] = generic.Type
			}
		}
		val.Type.Generics = append(val.Type.Generics, generic)
	}

	new.EnvName = val.Type.EnvName
	return val
}

func (w *Walker) typeExpression(typee *ast.TypeExpr, scope *Scope) Type {
	if typee == nil {
		return UnknownTyp
	}

	var typ Type
	if typee.Name.GetType() == ast.EnvironmentAccessExpression {
		expr, _ := typee.Name.(*ast.EnvAccessExpr)
		path := expr.PathExpr.Path

		walker, found := w.walkers[path.Lexeme]
		var env *Environment
		if !found {
			switch path.Lexeme {
			case "Pewpew":
				env = PewpewAPI
			case "Fmath":
				env = FmathAPI
			case "Math":
				env = MathAPI
			case "String":
				env = StringAPI
			case "Table":
				env = TableAPI
			default:
				w.AlertSingle(&alerts.InvalidEnvironment{}, path)
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

		if expr.Accessed.GetType() != ast.Identifier {
			return UnknownTyp
		}
		ident := expr.Accessed.(*ast.IdentifierExpr)
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
	case ast.Bool, ast.Text, ast.Number:
		typ = NewBasicType(pvt)
	case ast.Fixed:
		typ = NewFixedPointType()
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

		// check for types of the environment
		if val, ok := scope.Environment.Enums[typeName]; ok {
			typ = val.Type
			w.checkAccessibility(scope, val.IsPub, typee.Name.GetToken())
			break
		}
		if entityVal, found := scope.Environment.Entities[typeName]; found {
			typ = entityVal.GetType()
			break
		}
		if structVal, found := scope.Environment.Classes[typeName]; found {
			typ = structVal.GetType()
			break
		}
		if aliasType, found := scope.resolveAlias(typeName); found {
			typ = aliasType.UnderlyingType
			break
		}
		if aliasType, found := BuiltinEnv.Scope.AliasTypes[typeName]; found {
			typ = aliasType.UnderlyingType
			break
		}

		if scope.Environment.Name != w.environment.Name {
			typ = UnknownTyp
			break
		}

		// Check for function generics
		if gen, ok := w.resolveGenericParam(typeName, scope); ok {
			return gen
		}

		types := []Type{}
		envs := []string{}
		for _, v := range scope.Environment.UsedLibraries {
			typ := w.typeExpression(typee, &BuiltinLibraries[v].Scope)
			if typ != InvalidType && typ != UnknownTyp {
				types = append(types, typ)
				envs = append(envs, BuiltinLibraries[v].Name)
			}
		}

		for i := range scope.Environment.importedWalkers {
			if !scope.Environment.importedWalkers[i].Walked {
				scope.Environment.importedWalkers[i].Walk()
			}
			typ := w.typeExpression(typee, &scope.Environment.importedWalkers[i].environment.Scope)
			if typ != InvalidType && typ != UnknownTyp {
				types = append(types, typ)
				envs = append(envs, scope.Environment.importedWalkers[i].environment.Name)
			}
		}

		if len(types) > 1 {
			w.AlertSingle(&alerts.EnvironmentAccessAmbiguity{}, typee.GetToken(), envs, typeName)
			typ = InvalidType
			break
		}
		if len(types) == 0 {
			switch typeName {
			case "MeshEnv":
				return NewPathType(ast.MeshEnv)
			case "SoundEnv":
				return NewPathType(ast.SoundEnv)
			case "SharedEnv":
				return NewPathType(ast.SharedEnv)
			case "LevelEnv":
				return NewPathType(ast.LevelEnv)
			}
			break
		}
		typee.Name = &ast.EnvAccessExpr{
			PathExpr: &ast.EnvPathExpr{
				Path: tokens.Token{
					Lexeme:   envs[0],
					Location: typee.Name.GetToken().Location,
				},
			},
			Accessed: &ast.IdentifierExpr{
				Name: typee.Name.GetToken(),
			},
		}
		typ = types[0]
	}

	if typee.IsVariadic {
		return NewVariadicType(typ)
	}
	return typ
}
