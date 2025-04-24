package pass2

import (
	"fmt"
	"hybroid/ast"
	"hybroid/generator"
	"hybroid/tokens"
	wkr "hybroid/walker"
)

func StructExpr(w *wkr.Walker, node *ast.StructExpr, scope *wkr.Scope) *wkr.AnonStructVal {
	anonStructScope := wkr.NewScope(scope, &wkr.UntaggedTag{})
	structTypeVal := wkr.NewAnonStructVal(make(map[string]wkr.Field), false)

	for i := range node.Fields {
		FieldDeclarationStmt(w, node.Fields[i], structTypeVal, anonStructScope)
	}

	return structTypeVal
}

func GetReturns(w *wkr.Walker, returns *ast.TypeExpr, scope *wkr.Scope) wkr.Types {
	returnTypes := wkr.EmptyReturn
	if returns != nil {
		if returns.Name.GetType() != ast.TupleExpression {
			returnTypes = append(returnTypes, TypeExpr(w, returns, scope, true))
		} else {
			types := returns.Name.(*ast.TupleExpr).Types
			for _, typee := range types {
				returnTypes = append(returnTypes, TypeExpr(w, typee, scope, true))
			}
		}
	}

	return returnTypes
}

func FunctionExpr(w *wkr.Walker, fn *ast.FunctionExpr, scope *wkr.Scope) wkr.Value {
	returnTypes := GetReturns(w, fn.Return, scope)

	funcTag := &wkr.FuncTag{ReturnTypes: returnTypes}
	fnScope := wkr.NewScope(scope, funcTag, wkr.ReturnAllowing)

	WalkBody(w, &fn.Body, funcTag, fnScope)

	params := make([]wkr.Type, 0)
	for i, param := range fn.Params {
		params = append(params, TypeExpr(w, param.Type, scope, true))
		w.DeclareVariable(fnScope, &wkr.VariableVal{Name: param.Name.Lexeme, Value: w.TypeToValue(params[i]), IsLocal: true}, param.Name)
	}
	return &wkr.FunctionVal{ // returnTypes should contain a fn()
		Params:  params,
		Returns: returnTypes,
	}
}

func MatchExpr(w *wkr.Walker, node *ast.MatchExpr, scope *wkr.Scope) wkr.Value {
	mtt := &wkr.MatchExprTag{}

	matchScope := wkr.NewScope(scope, mtt, wkr.YieldAllowing)
	casesLength := len(node.MatchStmt.Cases) + 1
	if node.MatchStmt.HasDefault {
		casesLength--
	}
	matchScope.Tag = &wkr.MatchExprTag{YieldValues: make(wkr.Types, 0)}
	mpt := wkr.NewMultiPathTag(casesLength)

	GetNodeValue(w, &node.MatchStmt.ExprToMatch, scope)

	for i := range node.MatchStmt.Cases {
		caseScope := wkr.NewScope(matchScope, mpt)
		GetNodeValue(w, &node.MatchStmt.Cases[i].Expression, matchScope)
		WalkBody(w, &node.MatchStmt.Cases[i].Body, mpt, caseScope)
	}

	yieldValues := matchScope.Tag.(*wkr.MatchExprTag).YieldValues

	node.ReturnAmount = len(yieldValues)

	return yieldValues
}

func EntityExpr(w *wkr.Walker, node *ast.EntityExpr, scope *wkr.Scope) wkr.Value {
	val := GetNodeValue(w, &node.Expr, scope)
	typ := TypeExpr(w, node.Type, scope, false)

	if ident, ok := node.Type.Name.(*ast.IdentifierExpr); ok {
		switch ident.Name.Lexeme {
		case "Asteroid", "YellowBaf", "Inertiac", "Mothership",
			"MothershipBullet", "RollingCube", "RollingSphere",
			"Ufo", "Wary", "Crowder", "Ship", "Bomb", "BlueBaf",
			"RedBaf", "WaryMissile", "UfoBullet", "PlayerBullet",
			"BombExplosion", "PlayerExplosion", "Bonus", "FloatingMessage",
			"Pointonium", "BonusImplosion":
			typ = &wkr.RawEntityType{}
			node.OfficialEntityType = true
		}
	}

	if typ.PVT() != ast.Entity {
		w.Error(node.Token, "type given in entity expression is not an entity type")
	} else if !node.OfficialEntityType {
		varName := tokens.Token{}
		if node.ConvertedVarName != nil {
			varName = *node.ConvertedVarName
			w.Context.Conversions = append(w.Context.Conversions, wkr.NewEntityConversion(varName, w.TypeToValue(typ).(*wkr.EntityVal)))
		} else if len(w.Context.Conversions) != 0 {
			w.Context.Conversions = append(w.Context.Conversions, wkr.NewEntityConversion(varName, w.TypeToValue(typ).(*wkr.EntityVal)))
		}
		entityVal := w.TypeToValue(typ).(*wkr.EntityVal)
		node.EntityName = entityVal.Type.Name
		node.EnvName = entityVal.Type.EnvName
	} else if node.ConvertedVarName != nil {
		w.Error(*node.ConvertedVarName, "can't convert an entity to an official entity")
	}

	if val.GetType().GetType() != wkr.RawEntity {
		w.Error(node.Token, "value given in entity expression is not an entity")
	}

	return &wkr.BoolVal{}
}

func BinaryExpr(w *wkr.Walker, node *ast.BinaryExpr, scope *wkr.Scope) wkr.Value {
	left, right := GetNodeValue(w, &node.Left, scope), GetNodeValue(w, &node.Right, scope)
	leftType, rightType := left.GetType(), right.GetType()
	op := node.Operator
	switch op.Type {
	case tokens.Plus, tokens.Minus, tokens.Caret, tokens.Star, tokens.Slash, tokens.Modulo, tokens.BackSlash:
		w.ValidateArithmeticOperands(leftType, rightType, node)
		typ := w.DetermineValueType(leftType, rightType)

		if typ.PVT() == ast.Invalid {
			w.Error(node.GetToken(), fmt.Sprintf("invalid binary expression (left: %s, right: %s)", leftType.ToString(), rightType.ToString()))
			return &wkr.Invalid{}
		}

		return w.TypeToValue(typ)
	case tokens.Concat:
		if !wkr.TypeEquals(leftType, wkr.NewBasicType(ast.String)) && !wkr.TypeEquals(rightType, wkr.NewBasicType(ast.String)) {
			w.Error(node.GetToken(), fmt.Sprintf("invalid concatenation: left is %s and right is %s", leftType.ToString(), rightType.ToString()))
			return &wkr.Invalid{}
		}
		return &wkr.StringVal{}
	default:
		if op.Type == tokens.Or {
			if node.Left.GetType() == ast.EntityExpression && node.Left.(*ast.EntityExpr).ConvertedVarName != nil {
				w.Error(node.Left.GetToken(), "conversion of entity is not possible in a binary expression with 'or' operator")
			} else if node.Right.GetType() == ast.EntityExpression && node.Right.(*ast.EntityExpr).ConvertedVarName != nil {
				w.Error(node.Right.GetToken(), "conversion of entity is not possible in a binary expression with 'or' operator")
			}
		}

		if !wkr.TypeEquals(leftType, rightType) {
			w.Error(node.GetToken(), fmt.Sprintf("invalid comparison: types are not the same (left: %s, right: %s)", leftType.ToString(), rightType.ToString()))
			return &wkr.Invalid{}
		}
		return &wkr.BoolVal{}
	}
}

func LiteralExpr(w *wkr.Walker, node *ast.LiteralExpr) wkr.Value {
	switch node.ValueType {
	case ast.String:
		return &wkr.StringVal{}
	case ast.Fixed, ast.Radian, ast.FixedPoint, ast.Degree:
		return &wkr.FixedVal{SpecificType: node.ValueType}
	case ast.Bool:
		return &wkr.BoolVal{}
	case ast.Number:
		return &wkr.NumberVal{}
	case ast.Entity:
		return &wkr.RawEntityVal{}
	default:
		return &wkr.Invalid{}
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

func IdentifierExpr(w *wkr.Walker, node *ast.Node, scope *wkr.Scope) wkr.Value {
	valueNode := *node
	ident := valueNode.(*ast.IdentifierExpr)

	sc := w.ResolveVariable(scope, ident.Name.Lexeme)
	if sc == nil {
		walker, found := w.Walkers[ident.Name.Lexeme]
		if found {
			*node = &ast.LiteralExpr{
				Value: "\"" + walker.Environment.Path + "\"",
			}
			return wkr.NewPathVal(walker.Environment.Path, walker.Environment.Type)
		}

		return &wkr.Invalid{}
	}

	variable, notAllowed := w.GetVariable(sc, ident.Name.Lexeme)
	if notAllowed {
		w.Error(ident.GetToken(), "Not allowed to access a local variable from a different environment")
	}

	if sc.Tag.GetType() == wkr.Struct {
		class := sc.Tag.(*wkr.ClassTag).Val
		if variable.Value.GetType().GetType() == wkr.Fn {
			w.Context.Value2 = class
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
	} else if sc.Tag.GetType() == wkr.Entity {
		entity := sc.Tag.(*wkr.EntityTag).EntityType
		if variable.Value.GetType().GetType() == wkr.Fn {
			w.Context.Value2 = entity
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
	} else if sc.Environment.Name != w.Environment.Name {
		*node = &ast.EnvAccessExpr{
			PathExpr: &ast.EnvPathExpr{
				Path: tokens.Token{
					Lexeme:   sc.Environment.Name,
					Position: ident.GetToken().Position,
				},
			},
			Accessed: ident,
		}
	}

	if w.Context.PewpewVarFound {
		name, found := generator.PewpewEnums[w.Context.PewpewVarName][ident.Name.Lexeme]
		if found {
			ident.Name.Lexeme = name
		}
	}

	switch sc.Environment.Name {
	case "Pewpew":
		ident.Name.Lexeme = generator.PewpewVariables[ident.Name.Lexeme]
		if variable.GetType().GetType() != wkr.Fn {
			w.Context.PewpewVarFound = true
			w.Context.PewpewVarName = ident.Name.Lexeme
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

func EnvAccessExpr(w *wkr.Walker, node *ast.EnvAccessExpr) (wkr.Value, ast.Node) {
	envName := node.PathExpr.Path.Lexeme

	if node.Accessed.GetType() == ast.Identifier {
		name := node.Accessed.(*ast.IdentifierExpr).Name.Lexeme
		path := envName + ":" + name
		walker, found := w.Walkers[path]
		if found {
			return wkr.NewPathVal(walker.Environment.Path, walker.Environment.Type), &ast.LiteralExpr{
				Value: "\"" + walker.Environment.Path + "\"",
			}
		}
	}

	switch envName {
	case "Pewpew":
		if w.Environment.Type != ast.LevelEnv {
			w.Error(node.GetToken(), "cannot use the pewpew library in a non-level environment")
		}
		return GetNodeValueFromExternalEnv(w, node.Accessed, wkr.PewpewEnv), nil
	case "Fmath":
		if w.Environment.Type != ast.LevelEnv {
			w.Error(node.GetToken(), "cannot use the fmath library in a non-level environment")
		}
		return GetNodeValueFromExternalEnv(w, node.Accessed, wkr.FmathEnv), nil
	case "Math":
		return GetNodeValueFromExternalEnv(w, node.Accessed, wkr.MathEnv), nil
	case "String":
		return GetNodeValueFromExternalEnv(w, node.Accessed, wkr.StringEnv), nil
	case "Table":
		return GetNodeValueFromExternalEnv(w, node.Accessed, wkr.TableEnv), nil
	}

	walker, found := w.Walkers[envName]
	if !found {
		w.Error(node.GetToken(), fmt.Sprintf("Environment named '%s' doesn't exist", envName))
		return &wkr.Invalid{}, nil
	}

	if walker.Environment.Name == w.Environment.Name {
		w.Error(node.GetToken(), "cannot access self")
		return &wkr.Invalid{}, nil
	} else if walker.Environment.Path == "/dynamic/level.lua" {
		if !walker.Walked {
			Action(walker, w.Walkers)
		}
		value := GetNodeValue(w, &node.Accessed, &walker.Environment.Scope)
		return value, nil
	}

	envStmt := w.GetEnvStmt()

	for _, v := range walker.GetEnvStmt().Requirements {
		if v == w.Environment.Path {
			w.Error(node.GetToken(), fmt.Sprintf("import cycle detected: this environment and '%s' are using each other", walker.Environment.Name))
			return &wkr.Invalid{}, nil
		}
	}

	envStmt.AddRequirement(walker.Environment.Path)

	if !walker.Walked {
		Action(walker, w.Walkers)
	}

	value := GetNodeValue(w, &node.Accessed, &walker.Environment.Scope)

	return value, nil
}

func GroupingExpr(w *wkr.Walker, node *ast.GroupExpr, scope *wkr.Scope) wkr.Value {
	return GetNodeValue(w, &node.Expr, scope)
}

func ListExpr(w *wkr.Walker, node *ast.ListExpr, scope *wkr.Scope) wkr.Value {
	var value wkr.ListVal
	for i := range node.List {
		val := GetNodeValue(w, &node.List[i], scope)
		if val.GetType().PVT() == ast.Invalid {
			w.Error(node.List[i].GetToken(), fmt.Sprintf("variable '%s' inside list is invalid", node.List[i].GetToken().Lexeme))
		}
		value.Values = append(value.Values, val)
	}
	value.ValueType = wkr.GetContentsValueType(value.Values)
	return &value
}

func CallExpr(w *wkr.Walker, val wkr.Value, node *ast.CallExpr, scope *wkr.Scope) (wkr.Value, ast.Node) {
	valType := val.GetType().PVT()
	if valType != ast.Func {
		w.Error(node.Caller.GetToken(), "caller is not a function")
		return &wkr.Invalid{}, node
	}

	var finalNode ast.Node
	finalNode = node

	if entity, ok := w.Context.Value2.(*wkr.EntityVal); ok {
		caller := node.Caller.GetToken().Lexeme
		_, contains := entity.ContainsMethod(caller)
		if !contains {
			_, index, _ := entity.ContainsField(caller)
			finalNode = ConvertNodeToFieldExpr(node, index, ast.SelfEntity, entity.Type.EnvName, entity.Type.Name)
			goto skip
		}
		finalNode = ConvertCallToMethodCall(node, ast.SelfEntity, entity.Type.EnvName, entity.Type.Name)
		w.Context.Value2 = &wkr.Unknown{}
	} else if class, ok := w.Context.Value2.(*wkr.ClassVal); ok {
		caller := node.Caller.GetToken().Lexeme
		_, contains := class.ContainsMethod(caller)
		if !contains {
			_, index, _ := class.ContainsField(caller)
			finalNode = ConvertNodeToFieldExpr(node, index, ast.SelfStruct, class.Type.EnvName, class.Type.Name)
			goto skip
		}
		finalNode = ConvertCallToMethodCall(node, ast.SelfStruct, class.Type.EnvName, class.Type.Name)
		w.Context.Value2 = &wkr.Unknown{}
	}

skip:

	variable, it_is := val.(*wkr.VariableVal)
	if it_is {
		val = variable.Value
	}
	fun, _ := val.(*wkr.FunctionVal)

	suppliedGenerics := GetGenerics(w, node, node.GenericArgs, fun.Generics, scope)

	args := []wkr.Type{}
	for i := range node.Args {
		args = append(args, GetNodeValue(w, &node.Args[i], scope).GetType())
	}
	w.ValidateArguments(suppliedGenerics, args, fun.Params, node.Caller.GetToken())

	for i := range fun.Returns {
		if fun.Returns[i].GetType() == wkr.Generic {
			fun.Returns[i] = suppliedGenerics[fun.Returns[i].(*wkr.GenericType).Name]
		}
	}

	node.ReturnAmount = len(fun.Returns)

	if node.ReturnAmount == 1 {
		return w.TypeToValue(fun.Returns[0]), finalNode
	} else if node.ReturnAmount == 0 {
		return &wkr.Invalid{}, finalNode
	}
	return &fun.Returns, finalNode
}

func GetGenerics(w *wkr.Walker, node ast.Node, genericArgs []*ast.TypeExpr, expectedGenerics []*wkr.GenericType, scope *wkr.Scope) map[string]wkr.Type {
	receivedGenericsLength := len(genericArgs)
	expectedGenericsLength := len(expectedGenerics)

	suppliedGenerics := map[string]wkr.Type{}
	if receivedGenericsLength > expectedGenericsLength {
		w.Error(node.GetToken(), "too many generic arguments supplied")
	} else {
		for i := range genericArgs {
			suppliedGenerics[expectedGenerics[i].Name] = TypeExpr(w, genericArgs[i], scope, true)
		}
	}

	return suppliedGenerics
}

// Writes to context
func FieldExpr(w *wkr.Walker, node *ast.FieldExpr, scope *wkr.Scope) wkr.Value {
	if node.Identifier.GetToken().Lexeme == "converted" {
		print("brekpoint")
	}
	var scopeable wkr.ScopeableValue
	var val wkr.Value
	if w.Context.Node.GetType() != ast.NA {
		val = w.Context.Value
	} else {
		val = GetNodeValue(w, &node.Identifier, scope)
	}
	if variable, ok := val.(*wkr.VariableVal); ok {
		val = variable.Value
	}
	if val.GetType().GetType() == wkr.Named && val.GetType().PVT() == ast.Entity {
		node.ExprType = ast.SelfEntity
		named := val.GetType().(*wkr.NamedType)
		node.EntityName = named.Name
		node.EnvName = named.EnvName
	}

	if scpbl, ok := val.(wkr.ScopeableValue); ok {
		scopeable = scpbl
	} else {
		w.Error(node.Identifier.GetToken(), "variable is not of type class, struct, entity or enum")
		return &wkr.Invalid{}
	}

	newScope := scopeable.Scopify(scope, node)
	w.Context.Value = val
	w.Context.Node = node

	propVal := GetNodeValue(w, &node.PropertyIdentifier, newScope)
	if propVal.GetType().PVT() == ast.Invalid {
		ident := node.PropertyIdentifier.GetToken()
		w.Error(ident, fmt.Sprintf("'%s' doesn't exist", ident.Lexeme))
	}
	w.Context.Value = propVal
	w.Context.Node = node.Property

	defer w.Context.Clear()
	if node.Property.GetType() != ast.Identifier {
		return GetNodeValue(w, &node.Property, newScope)
	} // var1[1]["test"].method()
	return propVal
}

// Writes to context
func MemberExpr(w *wkr.Walker, node *ast.MemberExpr, scope *wkr.Scope) wkr.Value {
	var val wkr.Value
	if w.Context.Value.GetType().GetType() != wkr.NA {
		val = w.Context.Value
	} else {
		val = GetNodeValue(w, &node.Identifier, scope)
	}
	valType := val.GetType()
	if valType.GetType() != wkr.Wrapper {
		token := w.Context.Node.GetToken()
		w.Error(token, fmt.Sprintf("%s is not a map nor a list (found %s)", token.Lexeme, valType.ToString()))
		return &wkr.Invalid{}
	}

	propValPVT := GetNodeValue(w, &node.PropertyIdentifier, scope).GetType().PVT()
	if valType.PVT() == ast.Map {
		if propValPVT != ast.String {
			w.Error(node.Property.GetToken(), "expected string inside brackets for map accessing")
		}
	} else if valType.PVT() == ast.List {
		if propValPVT != ast.Number {
			w.Error(node.Property.GetToken(), "expected number inside brackets for list accessing")
		}
	}
	property := w.TypeToValue(valType.(*wkr.WrapperType).WrappedType)
	w.Context.Value = property
	w.Context.Node = node.Property
	nodePropertyType := node.Property.GetType()
	if nodePropertyType == ast.Identifier {
		w.Context.Clear()
		return property
	} else if nodePropertyType == ast.CallExpression {
		w.Context.Clear()
		val, newNode := CallExpr(w, property, node.Property.(*ast.CallExpr), scope)
		node.Property = newNode
		return val
	} else if nodePropertyType == ast.LiteralExpression {
		w.Context.Clear()
		return property
	}
	final := GetNodeValue(w, &node.Property, scope)
	w.Context.Clear()
	return final
}

func MapExpr(w *wkr.Walker, node *ast.MapExpr, scope *wkr.Scope) wkr.Value {
	mapVal := wkr.MapVal{Members: []wkr.Value{}}
	for _, v := range node.Map {
		val := GetNodeValue(w, &v.Expr, scope)
		mapVal.Members = append(mapVal.Members, val)
	}
	mapVal.MemberType = wkr.GetContentsValueType(mapVal.Members)
	return &mapVal
}

func UnaryExpr(w *wkr.Walker, node *ast.UnaryExpr, scope *wkr.Scope) wkr.Value {
	val := GetNodeValue(w, &node.Value, scope)
	valType := val.GetType()
	valPVT := valType.PVT()

	token := node.GetToken()

	if valPVT == ast.Invalid {
		w.Error(token, "value is invalid")
		return val
	}

	switch node.Operator.Type {
	case tokens.Bang:
		if valPVT != ast.Bool {
			w.Error(token, "value must be a bool to be negated")
		}
	case tokens.Hash:
		if valType.GetType() == wkr.Wrapper && valType.(*wkr.WrapperType).Type.PVT() != ast.List {
			w.Error(token, "value must be a list")
		} else if valType.GetType() != wkr.Wrapper {
			w.Error(token, "value must be a list")
		}
		return &wkr.NumberVal{}
	case tokens.Minus:
		if valPVT != ast.Number && valType.GetType() != wkr.Fixed {
			w.Error(token, "value must be a number or fixed")
		}
	}

	return val
}

func SelfExpr(w *wkr.Walker, self *ast.SelfExpr, scope *wkr.Scope) wkr.Value {
	if !scope.Is(wkr.SelfAllowing) {
		w.Error(self.Token, "can't use self outside of struct/entity")
		return &wkr.Invalid{}
	}

	sc, _, structTag := wkr.ResolveTagScope[*wkr.ClassTag](scope) // TODO: CHECK FOR ENTITY SCOPE

	if sc == nil {
		entitySc, _, entityTag := wkr.ResolveTagScope[*wkr.EntityTag](scope)
		if entitySc != nil {
			self.Type = ast.SelfEntity
			self.EntityName = (*entityTag).EntityType.Type.Name
			return (*entityTag).EntityType
		}

		return &wkr.Invalid{}
	}

	(*self).Type = ast.SelfStruct
	return (*structTag).Val
}

func NewExpr(w *wkr.Walker, new *ast.NewExpr, scope *wkr.Scope) wkr.Value {
	_type := TypeExpr(w, new.Type, scope, false)

	if _type.PVT() == ast.Invalid {
		w.Error(new.Type.GetToken(), "invalid type given in new expression")
		return &wkr.Invalid{}
	} else if _type.PVT() != ast.Struct {
		w.Error(new.Type.GetToken(), "type given in new expression is not a struct")
		return &wkr.Invalid{}
	}

	val := w.TypeToValue(_type).(*wkr.ClassVal)

	args := make([]wkr.Type, 0)
	for i := range new.Args {
		args = append(args, GetNodeValue(w, &new.Args[i], scope).GetType())
	}

	suppliedGenerics := GetGenerics(w, new, new.Generics, val.Generics, scope)

	w.ValidateArguments(suppliedGenerics, args, val.Params, new.Token)

	return val
}

func SpawnExpr(w *wkr.Walker, new *ast.SpawnExpr, scope *wkr.Scope) wkr.Value {
	typeExpr := TypeExpr(w, new.Type, scope, false)

	if typeExpr.PVT() == ast.Invalid {
		w.Error(new.Type.GetToken(), "invalid type given in spawn expression")
		return &wkr.Invalid{}
	} else if typeExpr.PVT() != ast.Entity {
		w.Error(new.Type.GetToken(), "type given in spawn expression is not an entity")
		return &wkr.Invalid{}
	}

	val := w.TypeToValue(typeExpr).(*wkr.EntityVal)

	args := make([]wkr.Type, 0)
	for i := range new.Args {
		args = append(args, GetNodeValue(w, &new.Args[i], scope).GetType())
	}

	suppliedGenerics := GetGenerics(w, new, new.Generics, val.SpawnGenerics, scope)

	w.ValidateArguments(suppliedGenerics, args, val.SpawnParams, new.Token)

	return val
}

func MethodCallExpr(w *wkr.Walker, mcall *ast.MethodCallExpr, scope *wkr.Scope) (wkr.Value, ast.Node) {
	val := GetNodeValue(w, &mcall.Identifier, scope)
	if variable, ok := val.(*wkr.VariableVal); ok {
		val = variable.Value
	}
	valType := val.GetType()

	if valType.PVT() == ast.Invalid {
		w.Error(mcall.Identifier.GetToken(), "value is invalid")
		return &wkr.Invalid{}, mcall
	}

	if structVal, ok := val.(*wkr.ClassVal); ok {
		mcall.EnvName = structVal.Type.EnvName
		mcall.TypeName = structVal.Type.Name
		mcall.ExprType = ast.SelfStruct
	} else if entityVal, ok := val.(*wkr.EntityVal); ok {
		mcall.EnvName = entityVal.Type.EnvName
		mcall.TypeName = entityVal.Type.Name
		mcall.ExprType = ast.SelfEntity
	}

	callToken := mcall.Call.Caller.GetToken()
	if methodContainer, ok := val.(wkr.MethodContainer); ok {
		method, found := methodContainer.ContainsMethod(callToken.Lexeme)

		mcall.MethodName = callToken.Lexeme

		if found {
			val, _ := CallExpr(w, method, mcall.Call, scope)
			return val, mcall
		}
	}
	if fieldContainer, ok := val.(wkr.FieldContainer); ok && valType.PVT() != ast.Enum {
		_, _, found := fieldContainer.ContainsField(callToken.Lexeme)

		if found {
			fieldExpr := &ast.FieldExpr{
				Property:           mcall.Call,
				PropertyIdentifier: mcall.Call.Caller,
				Identifier:         mcall.Identifier,
			}

			return FieldExpr(w, fieldExpr, scope), fieldExpr
		}
	} else {
		w.Error(mcall.Identifier.GetToken(), "value is not of type class or entity")
		return &wkr.Invalid{}, mcall
	}

	w.Error(mcall.GetToken(), "no method found")

	return &wkr.Invalid{}, mcall
}

func GetNodeValueFromExternalEnv(w *wkr.Walker, expr ast.Node, env *wkr.Environment) wkr.Value {
	val := GetNodeValue(w, &expr, &env.Scope)
	_, isTypes := val.(*wkr.Types)
	if !isTypes && val.GetType().PVT() == ast.Invalid {
		w.Error(expr.GetToken(), fmt.Sprintf("variable named '%s' doesn't exist", expr.GetToken().Lexeme))
	}
	return val
}

// func CastExpr(w *wkr.Walker, cast *ast.CastExpr, scope *wkr.Scope) wkr.Value {
// 	val := GetNodeValue(w, &cast.Value, scope)
// 	typ := TypeExpr(w, cast.Type, w.Environment)

// 	if typ.GetType() != wkr.CstmType {
// 		return &wkr.Invalid{}
// 	}

// 	cstm := typ.(*wkr.CustomType)

// 	if !wkr.TypeEquals(val.GetType(), cstm.UnderlyingType) {
// 		w.Error(cast.Value.GetToken(), fmt.Sprintf("expression type is %s, but underlying type is %s", val.GetType().ToString(), cstm.UnderlyingType.ToString()))
// 		return &wkr.Invalid{}
// 	}

// 	return wkr.NewCustomVal(cstm)
// }

func TypeExpr(w *wkr.Walker, typee *ast.TypeExpr, scope *wkr.Scope, throw bool) wkr.Type {
	if typee == nil {
		return wkr.InvalidType
	}

	var typ wkr.Type
	if typee.Name.GetType() == ast.EnvironmentAccessExpression {
		expr, _ := typee.Name.(*ast.EnvAccessExpr)
		path := expr.PathExpr.Path.Lexeme

		walker, found := w.Walkers[path]
		var env *wkr.Environment
		if !found {
			switch path {
			case "Pewpew":
				env = wkr.PewpewEnv
			case "Fmath":
				env = wkr.FmathEnv
			case "Math":
				env = wkr.MathEnv
			case "String":
				env = wkr.StringEnv
			case "Table":
				env = wkr.TableEnv
			default:
				w.Error(expr.GetToken(), "Environment name so doesn't exist")
				return wkr.InvalidType
			}

		} else if walker.Environment.Path == "/dynamic/level.lua" {
			if !walker.Walked {
				Action(walker, w.Walkers)
			}
			env = walker.Environment
		} else {

			for _, v := range walker.GetEnvStmt().Requirements {
				if v == w.Environment.Path {
					w.Error(typee.GetToken(), fmt.Sprintf("import cycle detected: this environment and '%s' are using each other", walker.Environment.Name))
					return wkr.InvalidType
				}
			}

			w.GetEnvStmt().AddRequirement(walker.Environment.Path)

			if !walker.Walked {
				Action(walker, w.Walkers)
			}

			env = walker.Environment
		}

		ident := &ast.IdentifierExpr{Name: expr.Accessed.GetToken(), ValueType: ast.Invalid}
		typ = TypeExpr(w, &ast.TypeExpr{Name: ident}, &env.Scope, throw)
		if typee.IsVariadic {
			return wkr.NewVariadicType(typ)
		}
		return typ
	}

	if typee.Name.GetToken().Type == tokens.Entity {
		typ = &wkr.RawEntityType{}
		if typee.IsVariadic {
			return wkr.NewVariadicType(typ)
		}
		return typ
	}

	pvt := w.GetTypeFromString(typee.Name.GetToken().Lexeme)
	switch pvt {
	case ast.Bool, ast.String, ast.Number:
		typ = wkr.NewBasicType(pvt)
	case ast.Fixed, ast.FixedPoint, ast.Radian, ast.Degree:
		typ = wkr.NewFixedPointType(pvt)
	case ast.Enum:
		typ = wkr.NewBasicType(ast.Enum)
	case ast.AnonStruct:
		fields := []*wkr.VariableVal{}

		for _, v := range typee.Fields {
			fields = append(fields, &wkr.VariableVal{
				Name:  v.Name.Lexeme,
				Value: w.TypeToValue(TypeExpr(w, v.Type, scope, throw)),
				Token: v.Name,
			})
		}

		typ = wkr.NewStructType(fields, false)
	case ast.Func:
		params := wkr.Types{}

		for _, v := range typee.Params {
			params = append(params, TypeExpr(w, v, scope, throw))
		}

		returns := GetReturns(w, typee.Return, scope)

		typ = &wkr.FunctionType{
			Params:  params,
			Returns: returns,
		}
	case ast.Map, ast.List:
		wrapped := TypeExpr(w, typee.WrappedType, scope, throw)
		typ = wkr.NewWrapperType(wkr.NewBasicType(pvt), wrapped)
	case ast.Entity:
		typ = &wkr.RawEntityType{}
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

		sc, _, fnTag := wkr.ResolveTagScope[*wkr.FuncTag](scope)

		if sc != nil {
			fnTag := *fnTag
			for _, v := range fnTag.Generics {
				if v.Name == typeName {
					return v
				}
			}
		}

		if len(scope.Environment.UsedLibraries) != 0 {
			if alias, found := wkr.BuiltinEnv.AliasTypes[typeName]; found {
				typ = alias.UnderlyingType
				break
			}
		}

		if scope.Environment.Name != w.Environment.Name {
			typ = wkr.InvalidType
			break
		}

		types := map[string]wkr.Type{}
		for i := range scope.Environment.UsedWalkers {
			if !scope.Environment.UsedWalkers[i].Walked {
				Action(scope.Environment.UsedWalkers[i], w.Walkers)
			}
			typ := TypeExpr(w, typee, &scope.Environment.UsedWalkers[i].Environment.Scope, false)
			if typ.PVT() != ast.Invalid {
				types[scope.Environment.UsedWalkers[i].Environment.Name] = typ
			}
		}

		for k, v := range scope.Environment.UsedLibraries {
			if !v {
				continue
			}

			typ := TypeExpr(w, typee, &wkr.LibraryEnvs[k].Scope, false)
			if typ.PVT() != ast.Invalid {
				types[wkr.LibraryEnvs[k].Name] = typ
			}
		}

		// if len(types) > 1 {
		// 	errorMsg := "conflicting types between: "
		// 	for k, v := range types {
		// 		errorMsg += k + ":" + v.ToString() + ", "
		// 	}
		// 	errorMsg = errorMsg[:len(errorMsg)-1]
		// 	w.Error(typee.GetToken(), errorMsg)
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
						Position: typee.Name.GetToken().Position,
					},
				},
				Accessed: &ast.IdentifierExpr{
					Name: typee.Name.GetToken(),
				},
			}
			return v
		}

		if len(scope.Environment.UsedLibraries) != 0 {
			if alias, found := wkr.BuiltinEnv.AliasTypes[typeName]; found {
				typ = alias.UnderlyingType

				break
			}
		}

		typ = wkr.InvalidType
	}

	if typee.IsVariadic {
		return wkr.NewVariadicType(typ)
	}
	if throw && typ.PVT() == ast.Invalid {
		w.Error(typee.GetToken(), "invalid type")
	}
	return typ
}
