package walker

import (
	"fmt"
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/generator/mapping"
	"hybroid/tokens"
	"reflect"
	"strconv"
)

func (w *Walker) structExpression(node *ast.StructExpr, scope *Scope) *StructVal {
	structTypeVal := NewStructVal(make(map[string]StructField))

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
		structTypeVal.AddField(NewVariable(fieldToken, val))
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
		if casesLength < 1 {
			w.AlertSingle(&alerts.InsufficientCases{}, matchStmt.Token)
		}
	} else if casesLength < 2 {
		w.AlertSingle(&alerts.InsufficientCases{}, matchStmt.Token)
	}
	matchScope := NewScope(scope, &MatchExprTag{YieldTypes: make([]Type, 0)}, YieldAllowing)
	valToMatch := w.GetActualNodeValue(&matchStmt.ExprToMatch, scope)
	valType := valToMatch.GetType()

	var prevPathTag PathTag
	for i := range cases {
		pt := NewPathTag()
		caseScope := NewScope(matchScope, pt)
		w.walkBody(&matchStmt.Cases[i].Body, pt, caseScope)
		if i != 0 {
			prevPathTag.SetAllExitAND(pt)
		} else {
			prevPathTag = *pt
		}

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
	if !matchStmt.HasDefault {
		prevPathTag.SetAllFalse()
	}

	if !prevPathTag.GetIfExits(Yield) {
		w.AlertSingle(&alerts.NotAllCodePathsExit{},
			matchStmt.Token,
			"yield",
		)
	}
	yieldTypes := matchScope.Tag.(*MatchExprTag).YieldTypes
	node.ReturnAmount = len(yieldTypes)

	switch node.ReturnAmount {
	case 0:
		return &Invalid{}
	case 1:
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
		var context string
		if scope.Environment.Name != w.environment.Name {
			context = "in the environment " + scope.Environment.Name
		}
		w.AlertSingle(&alerts.UndeclaredVariableAccess{}, identToken, identToken.Lexeme, context)
		return &Invalid{}
	}

	if (sc.Tag.GetType() == Class || sc.Tag.GetType() == Entity) && !scope.Is(SelfAllowing) {
		sc = w.resolveVariable(scope.Parent, ident.Name)
		goto check
	}

	variable := w.getVariable(sc, ident.Name)
	if val, ok := sc.ConstValues[variable.Name]; variable.IsConst && ok {
		*node = val
		ref := reflect.ValueOf(*node).Elem()
		field := ref.FieldByName("Token")
		if !field.IsValid() {
			field = ref.FieldByName("Name")
		}
		if !field.IsValid() {
			field = ref.FieldByName("Operator")
		}
		if field.IsValid() {
			location := field.FieldByName("Location")
			location.Set(reflect.ValueOf(identToken.Location))
		}
		if !w.context.DontSetToUsed {
			w.SetVarToUsed(variable)
			return variable
		}
		return variable
	}
	if sc.Tag.GetType() == Class {
		class := sc.Tag.(*ClassTag).Val
		field, index, found := class.ContainsField(variable.Name)

		if found {
			variable = field
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
		} else if method, found2 := class.Methods[variable.Name]; found2 {
			*node = &ast.MethodExpr{
				MethodInfo: method.Value.(*FunctionVal).MethodInfo,
				Token:      identToken,
				Access: &ast.SelfExpr{
					Token: ident.GetToken(),
					Type:  ast.ClassMethod,
				},
			}
			variable = method
		} else {
			panic(fmt.Sprintf("ResolveVariableScope stopped on an entity scope, but there was no field or method found. (identifier: %s, env: %s)", ident.Name.Lexeme, w.environment.Name))
		}
	} else if sc.Tag.GetType() == Entity && scope.Is(SelfAllowing) {
		entity := sc.Tag.(*EntityTag).EntityVal
		field, index, found := entity.ContainsField(variable.Name)

		if found {
			variable = field
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
		} else if method, found2 := entity.Methods[variable.Name]; found2 {
			*node = &ast.MethodExpr{
				MethodInfo: method.Value.(*FunctionVal).MethodInfo,
				Token:      identToken,
				Access: &ast.SelfExpr{
					Token: ident.GetToken(),
					Type:  ast.EntityMethod,
				},
			}
			variable = method
		} else {
			panic(fmt.Sprintf("ResolveVariableScope stopped on an entity scope, but there was no field or method found. (identifier: %s, env: %s)", ident.Name.Lexeme, w.environment.Name))
		}
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

	if !w.context.DontSetToUsed {
		w.SetVarToUsed(variable)
		return variable
	}

	return variable
}

func (w *Walker) environmentAccessExpression(expr *ast.Node) Value {
	node := (*expr).(*ast.EnvAccessExpr)
	envName := node.PathExpr.Path.Lexeme
	var accessed ast.Node = node.Accessed
	defer func() {
		if (*expr).GetType() != ast.EnvironmentAccessExpression {
			return
		}
		if accessed.GetType() != ast.Identifier {
			*expr = accessed
			return
		}
		node.Accessed = accessed.(*ast.IdentifierExpr)
	}()

	path := envName + ":" + node.Accessed.Name.Lexeme
	walker, found := w.walkers[path]
	if found {
		*expr = &ast.LiteralExpr{
			Value: "\"" + walker.environment.luaPath + "\"",
		}
		return NewPathVal(walker.environment.luaPath, walker.environment.Type, walker.environment.Name)
	}

	var val Value
	switch envName {
	case "Pewpew":
		w.AddLibrary(ast.Pewpew)
		val = w.GetNodeValue(&accessed, &PewpewAPI.Scope)
	case "Fmath":
		w.AddLibrary(ast.Fmath)
		val = w.GetNodeValue(&accessed, &FmathAPI.Scope)
	case "Math":
		w.AddLibrary(ast.Math)
		if w.environment.Type == ast.LevelEnv {
			w.AlertSingle(&alerts.UnallowedLibraryUse{}, node.PathExpr.Path, "Math", "Level")
		}
		val = w.GetNodeValue(&accessed, &MathAPI.Scope)
	case "String":
		w.AddLibrary(ast.String)
		val = w.GetNodeValue(&accessed, &StringAPI.Scope)
	case "Table":
		w.AddLibrary(ast.Table)
		val = w.GetNodeValue(&accessed, &TableAPI.Scope)
	default:
		walker, found := w.walkers[envName]
		if !found {
			w.AlertSingle(&alerts.InvalidEnvironmentAccess{}, node.PathExpr.GetToken(), envName)
			return &Invalid{}
		}

		if walker.environment.Name == w.environment.Name {
			return w.GetNodeValue(&accessed, &w.environment.Scope)
		}

		if walker.environment.Type != ast.SharedEnv && (w.environment.Type == ast.MeshEnv || w.environment.Type == ast.SoundEnv) {
			w.AlertSingle(&alerts.UnallowedEnvironmentAccess{}, node.PathExpr.GetToken(), "non Shared", "Mesh or Sound")
			return &Invalid{}
		} else if w.environment.Type == ast.LevelEnv && (walker.environment.Type == ast.MeshEnv || walker.environment.Type == ast.SoundEnv) {
			w.AlertSingle(&alerts.UnallowedEnvironmentAccess{}, node.PathExpr.GetToken(), "Mesh or Sound", "Level")
			return &Invalid{}
		}

		if paths, isCycle := w.ResolveImportCycle(walker); isCycle {
			paths = append([]string{w.environment.hybroidPath}, paths...)
			w.AlertSingle(&alerts.ImportCycle{}, node.PathExpr.Path, paths)
			return &Invalid{}
		}

		if walker.environment.luaPath == "/dynamic/level.lua" {
			if !walker.Walked {
				walker.Walk()
			}
			value := w.GetNodeValue(&accessed, &walker.environment.Scope)
			return value
		}

		w.environment.AddRequirement(walker.environment.luaPath)

		if !walker.Walked {
			walker.Walk()
		}

		val = w.GetNodeValue(&accessed, &walker.environment.Scope)
	}
	return val
}

func (w *Walker) groupExpression(node *ast.GroupExpr, scope *Scope) Value {
	return w.GetNodeValue(&node.Expr, scope)
}

func (w *Walker) listExpression(node *ast.ListExpr, scope *Scope) Value {
	value := &ListVal{}
	if len(node.List) == 0 {
		if node.Type == nil {
			w.AlertSingle(&alerts.UnknownListOrMapContents{}, node.Token)
		} else {
			value = w.typeToValue(w.typeExpression(node.Type, scope)).(*ListVal)
			return value
		}
	}
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

	if envAccess, ok := call.Caller.(*ast.EnvAccessExpr); ok {
		receiver_ := scope.resolveReturnable()
		if (envAccess.Accessed.Name.Lexeme == "ExplodeEntity" || envAccess.Accessed.Name.Lexeme == "DestroyEntity") && envAccess.PathExpr.Path.Lexeme == "Pewpew" && receiver_ != nil {
			(*receiver_).SetExit(true, EntityDestruction)
		}
	}

	fn := *val.(*FunctionVal)

	nodeGenerics := call.GenericArgs
	nodeArgs := call.Args
	genericArgs := w.getGenerics(nodeGenerics, fn.Generics, scope)
	args := []Value{}
	for i := range call.Args {
		args = append(args, w.GetActualNodeValue(&nodeArgs[i], scope))
	}
	w.validateArguments(genericArgs, args, &fn, call)

	actualReturns := fn.Returns
	returnLen := len(actualReturns)

	if fn.ProcType == Method {
		method := call.Caller.(*ast.MethodExpr)
		call.Caller = method.Access

		*node = convertCallToMethodCall(call, method)
		mcall := (*node).(*ast.MethodCallExpr)
		mcall.ReturnAmount = returnLen
	} else {
		call.ReturnAmount = returnLen
	}

	switch returnLen {
	case 0:
		return &Invalid{}
	case 1:
		return w.typeToValue(actualReturns[returnLen-1])
	}

	return w.typesToValues(actualReturns)
}

// Rewrote
func (w *Walker) accessExpression(_node *ast.Node, scope *Scope) Value {
	w.context.DontSetToUsed = false
	node := (*_node).(*ast.AccessExpr)
	var val Value

	typeExpr := &ast.TypeExpr{Name: node.Start}
	typ := w.typeExpression(typeExpr, scope)
	et, isEnum := typ.(*EnumType)
	if isEnum {
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
				if num.Value != "" {
					n, err := strconv.ParseFloat(num.Value, 64)
					if err == nil && n < float64(1) {
						w.AlertSingle(&alerts.ListIndexOutOfBounds{}, member.GetToken())
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
			w.ignoreAlerts = true
			fieldVal := w.GetNodeValue(&field.Field, scopedVal.Scopify(scope))
			w.ignoreAlerts = false

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
			if fn, ok := innerVal.(*FunctionVal); ok && fn.ProcType == Method {
				newAccess := *node
				methodExpr := &ast.MethodExpr{
					MethodInfo: fn.MethodInfo,
					Token:      node.Accessed[i].GetToken(),
					Access:     &newAccess,
				}
				newAccess.Accessed = newAccess.Accessed[:i]
				node.Start = methodExpr
				node.Accessed = node.Accessed[i+1:]
				*_node = node
				i = 0
				if len(node.Accessed) == 0 {
					*_node = methodExpr
					return fn
				}
			} else if entityVal, ok := val.(*EntityVal); ok && (*prevNode).GetType() != ast.SelfExpression {
				*prevNode = &ast.EntityAccessExpr{
					Expr:       *prevNode,
					EntityName: entityVal.Type.Name,
					EnvName:    entityVal.Type.EnvName,
				}
			}

			if i != len(node.Accessed) {
				val = innerVal
			} else {
				val = fieldVal
			}
			prevNode = &node.Accessed[i]
		}
	}

	if _, ok := val.(*VariableVal); !ok {
		val = NewVariable((*prevNode).GetToken(), val)
	}

	// check if we got an enum variant and convert that to its constant value
	if enumVal, ok := val.(*VariableVal).Value.(*EnumFieldVal); ok && isEnum {
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

	if len(node.KeyValueList) == 0 {
		if node.Type == nil {
			w.AlertSingle(&alerts.UnknownListOrMapContents{}, node.Token)
		} else {
			mapVal.MemberType = w.typeExpression(node.Type, scope)
		}
	}

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
	sc, classTag := resolveTagScope[*ClassTag](scope)

	if sc == nil {
		entitySc, entityTag := resolveTagScope[*EntityTag](scope)
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
	explicitGenericArgs := w.getGenerics(new.GenericArgs, val.New.Generics, scope)

	w.validateArguments(explicitGenericArgs, args, val.New, new)

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

	fn := val.Spawn
	explicitGenericArgs := w.getGenerics(new.GenericArgs, fn.Generics, scope)
	w.validateArguments(explicitGenericArgs, args, fn, new)

	new.EnvName = val.Type.EnvName
	return val
}

func (w *Walker) ResolveImportCycle(walker *Walker) ([]string, bool) {
	if walker == w {
		return []string{}, false
	}
	for _, v := range walker.environment.Requirements() {
		if v == w.environment.luaPath {
			return []string{walker.environment.hybroidPath}, true
		}
	}

	for _, v := range walker.environment.importedWalkers {
		if path, isCycle := w.ResolveImportCycle(v); isCycle {
			return append([]string{walker.environment.hybroidPath}, path...), true
		}
	}

	return []string{}, false
}

func (w *Walker) typeExpression(typee *ast.TypeExpr, scope *Scope) Type {
	var typ Type = UnknownTyp
	if typee == nil {
		return typ
	}

	defer func() {
		if typ == UnknownTyp || typ.GetType() == Wrapper {
			return
		}
		wrappedLen := len(typee.WrappedTypes)
		if ((typ.GetType() == Named && typ.PVT() == ast.Enum) || typ.GetType() != Named) && wrappedLen > 0 {
			w.AlertMulti(&alerts.TooManyElementsGiven{},
				typee.WrappedTypes[0].GetToken(),
				typee.WrappedTypes[wrappedLen-1].GetToken(),
				wrappedLen,
				"wrapped type",
				"in non-class/non-entity expression",
			)
		}
	}()

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
			if paths, isCycle := w.ResolveImportCycle(walker); isCycle {
				paths = append([]string{w.environment.hybroidPath}, paths...)
				w.AlertSingle(&alerts.ImportCycle{}, typee.GetToken(), paths)
				return InvalidType
			}

			w.environment.AddRequirement(walker.environment.luaPath)

			if !walker.Walked {
				walker.Walk()
			}
			env = walker.environment
		}

		typ = w.typeExpression(&ast.TypeExpr{Name: expr.Accessed}, &env.Scope)
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
		fields := []StructField{}

		for _, v := range typee.Fields {
			fields = append(fields, StructField{
				Var: &VariableVal{
					Name:  v.Name.Lexeme,
					Value: w.typeToValue(w.typeExpression(v.Type, scope)),
					Token: v.Name,
				},
			})
		}

		typ = NewStructType(fields)
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
		var wrapped Type = InvalidType
		wrappedLen := len(typee.WrappedTypes)
		if wrappedLen == 0 {
			w.AlertSingle(&alerts.InvalidListOrMapWrappedType{}, typee.GetToken())
		} else if wrappedLen > 1 {
			w.AlertSingle(&alerts.InvalidListOrMapWrappedType{}, typee.WrappedTypes[1].GetToken())
		} else {
			wrapped = w.typeExpression(typee.WrappedTypes[0], scope)
		}
		typ = NewWrapperType(NewBasicType(pvt), wrapped)
	case ast.Entity:
		typ = &RawEntityType{}
	default:
		typeName := typee.Name.GetToken().Lexeme

		// check for types of the environment
		if val, ok := scope.Environment.Enums[typeName]; ok {
			val.Type.IsUsed = true
			typ = val.Type
			w.checkAccessibility(scope, val.IsPub, typee.Name.GetToken())
			break
		}
		if entityVal, found := scope.Environment.Entities[typeName]; found {
			entityVal.Type.IsUsed = true
			val := CopyEntityVal(entityVal)
			typ = &val.Type
			w.FillGenericsInNamedType(&val.Type, typee, scope)
			w.checkAccessibility(scope, val.IsPub, typee.Name.GetToken())
			break
		}
		if classVal, found := scope.Environment.Classes[typeName]; found {
			classVal.Type.IsUsed = true
			val := CopyClassVal(classVal)
			typ = &val.Type
			w.FillGenericsInNamedType(&val.Type, typee, scope)
			w.checkAccessibility(scope, val.IsPub, typee.Name.GetToken())
			break
		}
		if aliasType, found := scope.resolveAlias(typeName); found {
			aliasType.IsUsed = true
			typ = aliasType.UnderlyingType
			w.checkAccessibility(scope, aliasType.IsPub, typee.Name.GetToken())
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

func (w *Walker) FillGenericsInNamedType(named *NamedType, typ *ast.TypeExpr, scope *Scope) {
	typesLen, genericsLen := len(typ.WrappedTypes), len(named.Generics)

	if typesLen < genericsLen {
		w.AlertSingle(&alerts.TooFewElementsGiven{}, typ.GetToken(), genericsLen-typesLen, "wrapped type", fmt.Sprintf("for the type '%s", named.String()))
	}

	for i := range typ.WrappedTypes {
		if i > genericsLen-1 {
			w.AlertMulti(&alerts.TooManyElementsGiven{},
				typ.WrappedTypes[i].GetToken(),
				typ.WrappedTypes[typesLen-1].GetToken(),
				typesLen-genericsLen,
				"wrapped type",
				fmt.Sprintf("for the type '%s", named.String()),
			)
			return
		}
		typ := w.typeExpression(typ.WrappedTypes[i], scope)
		named.Generics[i].Type = typ
	}
}
