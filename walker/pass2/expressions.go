package pass2

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
	wkr "hybroid/walker"
)

func AnonStructExpr(w *wkr.Walker, node *ast.AnonStructExpr, scope *wkr.Scope) *wkr.AnonStructVal {
	anonStructScope := wkr.NewScope(scope, &wkr.UntaggedTag{})
	structTypeVal := wkr.NewAnonStructVal(make(map[string]wkr.Field), false)

	for i := range node.Fields {
		FieldDeclarationStmt(w, node.Fields[i], structTypeVal, anonStructScope)
	}

	return structTypeVal
}

func AnonFnExpr(w *wkr.Walker, fn *ast.AnonFnExpr, scope *wkr.Scope) wkr.Value {
	returnTypes := wkr.EmptyReturn
	for i := range fn.Return {
		returnTypes = append(returnTypes, TypeExpr(w, fn.Return[i], scope, true))
	}
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

	for i := range node.MatchStmt.Cases {
		caseScope := wkr.NewScope(matchScope, mpt)
		GetNodeValue(w, &node.MatchStmt.Cases[i].Expression, matchScope)
		WalkBody(w, &node.MatchStmt.Cases[i].Body, mpt, caseScope)
	}

	yieldValues := matchScope.Tag.(*wkr.MatchExprTag).YieldValues

	node.ReturnAmount = len(yieldValues)

	return yieldValues
}

func BinaryExpr(w *wkr.Walker, node *ast.BinaryExpr, scope *wkr.Scope) wkr.Value {
	left, right := GetNodeValue(w, &node.Left, scope), GetNodeValue(w, &node.Right, scope)
	leftType, rightType := left.GetType(), right.GetType()
	op := node.Operator
	switch op.Type {
	case lexer.Plus, lexer.Minus, lexer.Caret, lexer.Star, lexer.Slash, lexer.Modulo:
		w.ValidateArithmeticOperands(leftType, rightType, node)
		typ := w.DetermineValueType(leftType, rightType)

		if typ.PVT() == ast.Invalid {
			w.Error(node.GetToken(), fmt.Sprintf("invalid binary expression (left: %s, right: %s)", leftType.ToString(), rightType.ToString()))
			return &wkr.Invalid{}
		}
		return &wkr.NumberVal{}
	case lexer.Concat:
		if !wkr.TypeEquals(leftType, wkr.NewBasicType(ast.String)) && !wkr.TypeEquals(rightType, wkr.NewBasicType(ast.String)) {
			w.Error(node.GetToken(), fmt.Sprintf("invalid concatenation: left is %s and right is %s", leftType.ToString(), rightType.ToString()))
			return &wkr.Invalid{}
		}
		return &wkr.StringVal{}
	default:
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
	default:
		return &wkr.Invalid{}
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
				Value: "\""+walker.Environment.Path+"\"",
			}
			return wkr.NewPathVal(walker.Environment.Path, walker.Environment.Type)
		}

		return &wkr.Invalid{}
	}

	variable := w.GetVariable(sc, ident.Name.Lexeme)

	if sc.Tag.GetType() == wkr.Struct { 
		_struct := sc.Tag.(*wkr.StructTag).StructVal
		_, index, _ := _struct.ContainsField(variable.Name)
		selfExpr := &ast.FieldExpr{
			Index: index,
			Identifier: &ast.SelfExpr{
				Token: valueNode.GetToken(),
				Type:  ast.SelfStruct,
			},
		}

		identExpr := &ast.IdentifierExpr{
			Name: valueNode.GetToken(),
		}
		selfExpr.Property = identExpr
		*node = selfExpr
	} else if sc.Tag.GetType() == wkr.Entity {
		entity := sc.Tag.(*wkr.EntityTag).EntityType
		_, index, _ := entity.ContainsField(variable.Name)
		selfExpr := &ast.FieldExpr{
			Index: index,
			Identifier: &ast.SelfExpr{
				Token: valueNode.GetToken(),
				Type:  ast.SelfEntity,
			},
		}

		identExpr := &ast.IdentifierExpr{
			Name: valueNode.GetToken(),
		}
		selfExpr.Property = identExpr
		*node = selfExpr
	}

	variable.IsUsed = true
	return variable.Value
}

func EnvAccessExpr(w *wkr.Walker, node *ast.EnvAccessExpr) (wkr.Value, ast.Node) {
	envName := node.PathExpr.Nameify()

	if node.Accessed.GetType() == ast.Identifier {
		name := node.Accessed.(*ast.IdentifierExpr).Name.Lexeme
		path := envName+"::"+name
		walker, found := w.Walkers[path]
		if found {
			return wkr.NewPathVal(walker.Environment.Path, walker.Environment.Type), &ast.LiteralExpr{
				Value: "\""+walker.Environment.Path+"\"",
			}
		}
	}

	walker, found := w.Walkers[envName]
	if !found {
		w.Error(node.GetToken(), fmt.Sprintf("Environment named '%s' doesn't exist", envName))
		return &wkr.Invalid{}, nil
	}

	envStmt := w.GetEnvStmt()

	for _, v := range walker.GetEnvStmt().Requirements {
		if v == w.Environment.Path {
			w.Error(node.GetToken(), fmt.Sprintf("import cycle detected: this environment and '%s' are using each other", walker.Environment.Name))
		}
	}

	envStmt.AddRequirement(walker.Environment.Path)

	if (!walker.Walked) {
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

func CallExpr(w *wkr.Walker, val wkr.Value, node *ast.CallExpr, scope *wkr.Scope) wkr.Value {
	valType := val.GetType().PVT()
	if valType != ast.Func {
		w.Error(node.Caller.GetToken(), "caller is not a function") 
		return &wkr.Invalid{}                                       
	}

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

	if len(fun.Returns) == 1 {
		return w.TypeToValue(fun.Returns[0]) 
	}
	return &fun.Returns
}

func GetGenerics(w *wkr.Walker, node ast.Node, genericArgs []*ast.TypeExpr, expectedGenerics []*wkr.GenericType, scope *wkr.Scope) map[string]wkr.Type {
	receivedGenericsLength := len(genericArgs)
	expectedGenericsLength := len(expectedGenerics)

	suppliedGenerics := map[string]wkr.Type{}
	if receivedGenericsLength > expectedGenericsLength {
		w.Error(node.GetToken(), "too many generic arguments supplied")
	}else {
		for i := range genericArgs {
			suppliedGenerics[expectedGenerics[i].Name] = TypeExpr(w, genericArgs[i], scope, true)
		}
	}

	return suppliedGenerics
}

func FieldExpr(w *wkr.Walker, node *ast.FieldExpr, scope *wkr.Scope) wkr.Value {// WRITES CONTEXT
	var val wkr.Value
	if node.Identifier.GetToken().Lexeme == "EnumTest" {
		println("breakpoint")
	}
	if w.Context.Value.GetType().GetType() != wkr.NA {
		scopeable, ok :=  w.Context.Value.(wkr.ScopeableValue)
		if !ok {
			w.Error(w.Context.Node.GetToken(), fmt.Sprintf("variable '%s' is not a struct, entity, enum or anonymous struct", w.Context.Node.GetToken().Lexeme))
			return &wkr.Invalid{}
		}
		w.Context.Clear()
		return GetNodeValue(w, &node.Property, scopeable.Scopify(scope, node))
	}else {
		val = GetNodeValue(w, &node.Identifier, scope)
	}

	if !wkr.IsOfPrimitiveType(val, ast.Struct, ast.Entity, ast.AnonStruct, ast.Enum) {
		w.Error(node.Identifier.GetToken(), fmt.Sprintf("variable '%s' is not a struct, entity, enum or anonymous struct", node.Identifier.GetToken().Lexeme))
		return &wkr.Invalid{}
	}
	owner := val.(wkr.ScopeableValue)
	w.Context.Value = owner
	w.Context.Node = node.Property
	finalValue := GetNodeValue(w, &node.Property, owner.Scopify(scope, node))
	w.Context.Clear()
	return finalValue
}

func MemberExpr(w *wkr.Walker, node *ast.MemberExpr, scope *wkr.Scope) wkr.Value {// WRITES CONTEXT
	var val wkr.Value
	if w.Context.Value.GetType().GetType() != wkr.NA {
		val = w.Context.Value
	}else {
		val = GetNodeValue(w, &node.Identifier, scope)
	}
	valType := val.GetType()
	if valType.GetType() != wkr.Wrapper {
		token := w.Context.Node.GetToken()
		w.Error(token, fmt.Sprintf("%s is not a map nor a list (found %s)", token.Lexeme, valType.ToString()))
		return &wkr.Invalid{}
	}
	
	if node.Property.GetValueType() != ast.Ident {
	}else if valType.PVT() == ast.Map {
		if node.GetValueType() != ast.String {
			w.Error(node.Property.GetToken(), "expected string inside brackets for map accessing")
		}
	}else if valType.PVT() == ast.List {
		if node.GetValueType() != ast.Number {
			w.Error(node.Property.GetToken(), "expected number inside brackets for list accessing")
		}
	}
	property := w.TypeToValue(val.GetType().(*wkr.WrapperType).WrappedType)
	w.Context.Value = property
	w.Context.Node = node.Property
	nodePropertyType := node.Property.GetType()
	if nodePropertyType == ast.Identifier {
		w.Context.Clear()
		return property
	}else if nodePropertyType == ast.CallExpression {
		w.Context.Clear()
		return CallExpr(w, property, node.Property.(*ast.CallExpr), scope)
	}else if nodePropertyType == ast.LiteralExpression {
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
	return GetNodeValue(w, &node.Value, scope)
}

func SelfExpr(w *wkr.Walker, self *ast.SelfExpr, scope *wkr.Scope) wkr.Value {
	if !scope.Is(wkr.SelfAllowing) {
		w.Error(self.Token, "can't use self outside of struct/entity")
		return &wkr.Invalid{}
	}

	sc, _, structTag := wkr.ResolveTagScope[*wkr.StructTag](scope) // TODO: CHECK FOR ENTITY SCOPE

	if sc == nil {
		entitySc, _, entityTag := wkr.ResolveTagScope[*wkr.EntityTag](scope)
		if entitySc != nil {
			return (*entityTag).EntityType
		}

		return &wkr.Invalid{}
	}

	(*self).Type = ast.SelfStruct
	return (*structTag).StructVal
}

func NewExpr(w *wkr.Walker, new *ast.NewExpr, scope *wkr.Scope) wkr.Value {
	w.Context.Node = new
	_type := TypeExpr(w, new.Type, scope, false)

	if _type.PVT() == ast.Invalid {
		w.Error(new.Type.GetToken(), "invalid type given in new expression")
		return &wkr.Invalid{}
	} else if _type.PVT() != ast.Struct {
		w.Error(new.Type.GetToken(), "type given in new expression is not a struct")
		return &wkr.Invalid{}
	}

	val := w.TypeToValue(_type).(*wkr.StructVal)

	args := make([]wkr.Type, 0)
	for i := range new.Args {
		args = append(args, GetNodeValue(w, &new.Args[i], scope).GetType())
	}

	suppliedGenerics := GetGenerics(w, new, new.Generics, val.Generics, scope)

	w.ValidateArguments(suppliedGenerics, args, val.Params, new.Token)

	return val
}

func SpawnExpr(w *wkr.Walker, new *ast.SpawnExpr, scope *wkr.Scope) wkr.Value {
	w.Context.Node = new
	_type := TypeExpr(w, new.Type, scope, false)

	if _type.PVT() == ast.Invalid {
		w.Error(new.Type.GetToken(), "invalid type given in spawn expression")
		return &wkr.Invalid{}
	} else if _type.PVT() != ast.Entity {
		w.Error(new.Type.GetToken(), "type given in spawn expression is not an entity")
		return &wkr.Invalid{}
	}

	val := w.TypeToValue(_type).(*wkr.EntityVal)

	args := make([]wkr.Type, 0)
	for i := range new.Args {
		args = append(args, GetNodeValue(w, &new.Args[i], scope).GetType())
	}

	suppliedGenerics := GetGenerics(w, new, new.Generics, val.SpawnGenerics, scope)

	w.ValidateArguments(suppliedGenerics, args, val.SpawnParams, new.Token)

	return val
}

func GetNodeValueFromExternalEnv(w *wkr.Walker, expr ast.Node, scope *wkr.Scope, env *wkr.Environment) wkr.Value {
	env.Scope.Parent = scope
	env.Scope.Attributes = scope.Attributes
	val := GetNodeValue(w, &expr, &env.Scope)
	_, isTypes := val.(*wkr.Types)
	if !isTypes && val.GetType().PVT() == ast.Invalid {
		w.Error(expr.GetToken(), fmt.Sprintf("variable named '%s' doesn't exist", expr.GetToken().Lexeme))
	}
	return val
}

func PewpewExpr(w *wkr.Walker, expr *ast.PewpewExpr, scope *wkr.Scope) wkr.Value {
	return GetNodeValueFromExternalEnv(w, expr.Node, scope, wkr.PewpewEnv)
}

func FmathExpr(w *wkr.Walker, expr *ast.FmathExpr, scope *wkr.Scope) wkr.Value {
	return GetNodeValueFromExternalEnv(w, expr.Node, scope, wkr.FmathEnv)
}

func StandardExpr(w *wkr.Walker, expr *ast.StandardExpr, scope *wkr.Scope) wkr.Value {
	if expr.Library == ast.MathLib {
		if w.Nodes[0].(*ast.EnvironmentStmt).EnvType.Type == ast.Level {
			w.Error(expr.GetToken(), "cannot use the Math library in a Level environment")
			return &wkr.Invalid{}
		}

		return GetNodeValueFromExternalEnv(w, expr.Node, scope, wkr.MathEnv)
	}else if expr.Library == ast.StringLib {
		return GetNodeValueFromExternalEnv(w, expr.Node, scope, wkr.StringEnv)
	}else {
		return GetNodeValueFromExternalEnv(w, expr.Node, scope, wkr.TableEnv)
	}
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
		path := expr.PathExpr.Nameify()

		walker, found := w.Walkers[path]
		if !found {
			w.Error(expr.PathExpr.GetToken(), "Environment name so doesn't exist")
			return wkr.InvalidType
		}
		typ = TypeExpr(w, &ast.TypeExpr{Name: expr.Accessed}, &walker.Environment.Scope, throw)
		if typee.IsVariadic {
			return wkr.NewVariadicType(typ)
		}
	}

	if typee.Name.GetToken().Type == lexer.Entity {
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
		fields := map[string]wkr.Field{}

		for i, v := range typee.Fields {
			fields[v.Name.Lexeme] = wkr.NewField(i, &wkr.VariableVal{
				Name:  v.Name.Lexeme,
				Value: w.TypeToValue(TypeExpr(w, v.Type, scope, throw)),
				Token: v.Name,
			})
		}

		typ = wkr.NewAnonStructType(fields, false)
	case ast.Func:
		params := wkr.Types{}

		for _, v := range typee.Params {
			params = append(params, TypeExpr(w, v, scope, throw))
		}

		returns := wkr.Types{}
		for _, v := range typee.Returns {
			returns = append(returns, TypeExpr(w, v, scope, throw))
		}

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
		if entityVal, found := w.Environment.Entities[typeName]; found {
			typ = entityVal.GetType()
			break
		}
		if structVal, found := scope.Environment.Structs[typeName]; found {
			typ = structVal.GetType()
			break
		}
		if customType, found := w.Environment.CustomTypes[typeName]; found {
			typ = customType
			break
		}
		if val := w.GetVariable(scope, typeName); val != nil {
			if val.GetType().PVT() == ast.Enum {
				typ = val.GetType()
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
