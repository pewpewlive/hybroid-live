package pass1

import (
	"hybroid/ast"
	"hybroid/lexer"
	wkr "hybroid/walker"
)

func AnonStructExpr(w *wkr.Walker, node *ast.AnonStructExpr, scope *wkr.Scope) *wkr.AnonStructVal {
	anonStructScope := wkr.NewScope(scope, &wkr.UntaggedTag{})
	structTypeVal := wkr.NewAnonStructVal(make(map[string]*wkr.VariableVal), false)

	for i := range node.Fields {
		FieldDeclarationStmt(w, node.Fields[i], structTypeVal, anonStructScope)
	}

	return structTypeVal
}

func AnonFnExpr(w *wkr.Walker, fn *ast.AnonFnExpr, scope *wkr.Scope) wkr.Value {
	returnTypes := wkr.EmptyReturn
	for i := range fn.Return {
		returnTypes = append(returnTypes, TypeExpr(w, fn.Return[i]))
	}
	funcTag := &wkr.FuncTag{ReturnTypes: returnTypes}
	fnScope := wkr.NewScope(scope, funcTag, wkr.ReturnAllowing)

	params := WalkParams(w, fn.Params, scope, func(name lexer.Token, value wkr.Value) {
		w.DeclareVariable(fnScope, &wkr.VariableVal{
			Name: name.Lexeme,
			Value: value,
			IsLocal: true,
			Token: name,
		}, name)
	})

	WalkBody(w, &fn.Body, funcTag, fnScope)

	return &wkr.FunctionVal{
		Returns: returnTypes,
		Params: params,
	}
}

func MatchExpr(w *wkr.Walker, node *ast.MatchExpr, scope *wkr.Scope) wkr.Value {
	mtt := &wkr.MatchExprTag{}

	matchScope := wkr.NewScope(scope, mtt, wkr.YieldAllowing)
	casesLength := len(node.MatchStmt.Cases) + 1
	if node.MatchStmt.HasDefault {
		casesLength--
	}
	mpt := wkr.NewMultiPathTag(casesLength)

	for i := range node.MatchStmt.Cases {
		caseScope := wkr.NewScope(matchScope, mpt)
		WalkBody(w, &node.MatchStmt.Cases[i].Body, mpt, caseScope)
	}

	for _, v := range mtt.YieldValues {
		if v.PVT() == ast.Unresolved {
			return wkr.NewUnresolvedVal(node)
		}
	}

	return mtt.YieldValues
}

func BinaryExpr(w *wkr.Walker, node *ast.BinaryExpr, scope *wkr.Scope) wkr.Value {
	op := node.Operator
	switch op.Type {
	case lexer.Plus, lexer.Minus, lexer.Caret, lexer.Star, lexer.Slash, lexer.Modulo:
		return &wkr.NumberVal{}
	default:
		return &wkr.BoolVal{}
	}
}

func LiteralExpr(w *wkr.Walker, node *ast.LiteralExpr) wkr.Value {
	switch node.ValueType {
	case ast.String:
		return &wkr.StringVal{}
	case ast.Fixed:
		return &wkr.FixedVal{SpecificType: ast.Fixed}
	case ast.Radian:
		return &wkr.FixedVal{SpecificType: ast.Radian}
	case ast.FixedPoint:
		return &wkr.FixedVal{SpecificType: ast.FixedPoint}
	case ast.Degree:
		return &wkr.FixedVal{SpecificType: ast.Degree}
	case ast.Bool:
		return &wkr.BoolVal{}
	case ast.Number:
		return &wkr.NumberVal{}
	default:
		return &wkr.Unknown{}
	}
}

func IdentifierExpr(w *wkr.Walker, node *ast.Node, scope *wkr.Scope) wkr.Value {
	valueNode := *node
	ident := valueNode.(*ast.IdentifierExpr)

	sc := w.ResolveVariable(scope, ident.Name.Lexeme)
	if sc == nil {
		return &wkr.Unknown{}
	}

	variable := w.GetVariable(sc, ident.Name.Lexeme)

	if sc.Tag.GetType() == wkr.Struct {
		_struct := sc.Tag.(*wkr.StructTag).StructVal
		_, index, _ := _struct.ContainsField(variable.Name)
		selfExpr := &ast.FieldExpr{
			Identifier: &ast.SelfExpr{
				Token: valueNode.GetToken(),
				Type:  ast.SelfStruct,
			},
		}

		fieldExpr := &ast.FieldExpr{
			Identifier: valueNode,
			Index:      index,
		}
		selfExpr.Property = fieldExpr
		*node = selfExpr
	} else if sc.Tag.GetType() == wkr.Entity {
		entity := sc.Tag.(*wkr.EntityTag).EntityType
		_, index, _ := entity.ContainsField(variable.Name)
		selfExpr := &ast.FieldExpr{
			Identifier: &ast.SelfExpr{
				Token: valueNode.GetToken(),
				Type:  ast.SelfEntity,
			},
		}

		fieldExpr := &ast.FieldExpr{
			Identifier: valueNode,
			Index:      index,
		}
		selfExpr.Property = fieldExpr
		*node = selfExpr
	}
	variable.IsUsed = true
	return variable.Value
}

func EnvAccessExpr(w *wkr.Walker, node *ast.EnvAccessExpr, scope *wkr.Scope) wkr.Value {
	return wkr.NewUnresolvedVal(node)
}

func GroupingExpr(w *wkr.Walker, node *ast.GroupExpr, scope *wkr.Scope) wkr.Value {
	return GetNodeValue(w, &node.Expr, scope)
}

func ListExpr(w *wkr.Walker, node *ast.ListExpr, scope *wkr.Scope) wkr.Value {
	var value wkr.ListVal
	for i := range node.List {
		val := GetNodeValue(w, &node.List[i], scope)
		// if val.GetType().PVT() == ast.Invalid {
		// 	w.Error(node.List[i].GetToken(), fmt.Sprintf("variable '%s' inside list is invalid", node.List[i].GetToken().Lexeme))
		// }
		value.Values = append(value.Values, val)
	}
	value.ValueType = wkr.GetContentsValueType(value.Values)
	return &value
}

func CallExpr(w *wkr.Walker, val wkr.Value, node *ast.CallExpr, scope *wkr.Scope) wkr.Value { 
	valType := val.GetType().PVT() 
	if valType != ast.Func {
		return &wkr.Invalid{}                                       
	}

	variable, it_is := val.(*wkr.VariableVal) 
	if it_is {
		val = variable.Value
	}
	fun, _ := val.(*wkr.FunctionVal)

	for i := range node.Args {
		GetNodeValue(w, &node.Args[i], scope).GetType()
	}

	if len(fun.Returns) == 1 {
		return w.TypeToValue(fun.Returns[0]) 
	}
	return &fun.Returns
}

func PropertyExpr(w *wkr.Walker, node ast.Node, scope *wkr.Scope) wkr.Value {
	var val wkr.Value
	val = &wkr.Invalid{}
	
	switch w.Context.Node.(type) {
	case *ast.FieldExpr:
		owner := w.Context.Value
		if !wkr.IsOfPrimitiveType(owner, ast.Struct, ast.Entity, ast.AnonStruct, ast.Enum) {
			break
		}
		if node.GetType() == ast.CallExpression {
			val := PropertyExpr(w, node.(*ast.CallExpr).Caller, scope)
			return CallExpr(w, val, node.(*ast.CallExpr), scope)
		}
		fieldContainer, _ := owner.(wkr.FieldContainer)
		if value, index, ok := fieldContainer.ContainsField(node.GetToken().Lexeme); ok {
			w.Context.Node.(*ast.FieldExpr).Index = index

			val = value
		}else {
			methodContainer, ok := owner.(wkr.MethodContainer)
			if !ok {
			}else if value, ok := methodContainer.ContainsMethod(node.GetToken().Lexeme); ok {
				val = value
			}
		}
	case *ast.MemberExpr:
		owner := w.Context.Value
		if owner.GetType().GetType() != wkr.Wrapper {
			break
		}
		if node.GetType() == ast.CallExpression {
			val := PropertyExpr(w, node.(*ast.CallExpr).Caller, scope)
			return CallExpr(w, val, node.(*ast.CallExpr), scope)
		}
		expr := GetNodeValue(w, &node, scope)
		if expr.GetType().PVT() != ast.String && expr.GetType().PVT() != ast.Number {
		}else {
			ownerType, _ := owner.GetType().(*wkr.WrapperType)
			val = w.TypeToValue(ownerType.WrappedType)
		}
	}

	return val
}

func FieldExpr(w *wkr.Walker, node *ast.FieldExpr, scope *wkr.Scope) wkr.Value {// WRITES CONTEXT
	val := GetNodeValue(w, &node.Identifier, scope)

	if !wkr.IsOfPrimitiveType(val, ast.Struct, ast.Entity, ast.AnonStruct, ast.Enum) {
		return &wkr.Invalid{}
	}

	w.Context.Value = val
	w.Context.Node = node
	finalValue := PropertyExpr(w, node.Property, scope)
	if member, ok := node.Property.(*ast.MemberExpr); ok {
		return MemberExpr(w, member, scope)
	}else if field, ok := node.Property.(*ast.FieldExpr); ok {
		return FieldExpr(w, field, scope)
	}
	w.Context.Clear()
	return finalValue
}

func MemberExpr(w *wkr.Walker, node *ast.MemberExpr, scope *wkr.Scope) wkr.Value {// WRITES CONTEXT
	val := GetNodeValue(w, &node.Identifier, scope)

	if val.GetType().GetType() != wkr.Wrapper {
		return &wkr.Invalid{}
	}
	w.Context.Value = val
	w.Context.Node = node
	finalValue := PropertyExpr(w, node.Property, scope)
	if member, ok := node.Property.(*ast.MemberExpr); ok {
		return MemberExpr(w, member, scope)
	}else if field, ok := node.Property.(*ast.FieldExpr); ok {
		return FieldExpr(w, field, scope)
	}
	w.Context.Clear()
	return finalValue
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
		return &wkr.Unknown{}
	}

	sc, _, structTag := wkr.ResolveTagScope[*wkr.StructTag](scope) // TODO: CHECK FOR ENTITY SCOPE

	if sc != nil {
		(*self).Type = ast.SelfStruct
		return (*structTag).StructVal
	} else {
		return &wkr.Unknown{}
	}
}

func NewExpr(w *wkr.Walker, new *ast.NewExpr, scope *wkr.Scope) wkr.Value {
	w.Context.Node = new
	var val *wkr.StructVal
	var found bool
	if new.Type.Name.GetType() == ast.Identifier {
		val, found = w.GetStruct(new.Type.GetToken().Lexeme)
	} else {
		return &wkr.Unknown{}
	}
	if !found {
		return &wkr.Unknown{}
	}

	return val
}

func SpawnExpr(w *wkr.Walker, new *ast.SpawnExpr, scope *wkr.Scope) wkr.Value {
	w.Context.Node = new
	var val *wkr.EntityVal
	var found bool
	if new.Type.GetType() == ast.Identifier {
		val, found = w.GetEntity(new.Type.GetToken().Lexeme)
	} else {
		return &wkr.Unknown{}
	}
	if !found {
		return &wkr.Unknown{}
	}

	return val
}

func GetNodeValueFromExternalEnv(w *wkr.Walker, expr ast.Node, scope *wkr.Scope, env *wkr.Environment) wkr.Value {
	env.Scope.Parent = scope
	env.Scope.Attributes = scope.Attributes
	return GetNodeValue(w, &expr, &env.Scope)
}

func PewpewExpr(w *wkr.Walker, expr *ast.PewpewExpr, scope *wkr.Scope) wkr.Value {
	return GetNodeValueFromExternalEnv(w, expr.Node, scope, wkr.PewpewEnv)
}

func FmathExpr(w *wkr.Walker, expr *ast.FmathExpr, scope *wkr.Scope) wkr.Value {
	return GetNodeValueFromExternalEnv(w, expr.Node, scope, wkr.FmathEnv)
}

func StandardExpr(w *wkr.Walker, expr *ast.StandardExpr, scope *wkr.Scope) wkr.Value {
	val := GetNodeValueFromExternalEnv(w, expr.Node, scope, wkr.MathEnv)
	return val
}

// func CastExpr(w *wkr.Walker, cast *ast.CastExpr, scope *wkr.Scope) wkr.Value {
// 	val := GetNodeValue(w, &cast.Value, scope)

// 	if (val.GetType().GetType() == wkr.Unresolved) {
// 		return val
// 	}

// 	typ := TypeExpr(w, cast.Type)

// 	if typ.GetType() == wkr.Unresolved {
// 		return &wkr.UnresolvedVal{
// 			Expr: typ.(*wkr.UnresolvedType).Expr,
// 		}
// 	}

// 	if typ.GetType() != wkr.CstmType {
// 		w.Error(cast.Type.GetToken(), "can only accept custom types in cast")
// 	}

// 	return &wkr.Unknown{}
// } 

func TypeExpr(w *wkr.Walker, typee *ast.TypeExpr) wkr.Type {
	if typee == nil {
		return wkr.InvalidType
	}
	var typ wkr.Type

	if typee.Name.GetType() == ast.EnvironmentAccessExpression {
		expr, _ := typee.Name.(*ast.EnvAccessExpr)
		typ = &wkr.UnresolvedType{
			Expr: expr,
		}
		if typee.IsVariadic {
			return wkr.NewVariadicType(typ)
		}
		return typ
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
		fields := map[string]*wkr.VariableVal{}

		for _, v := range typee.Fields {
			fields[v.Name.Lexeme] = &wkr.VariableVal{
				Name:  v.Name.Lexeme,
				Value: w.TypeToValue(TypeExpr(w, v.Type)),
				Token: v.Name,
			}
		}

		typ = wkr.NewAnonStructType(fields, false)
	case ast.Func:
		params := wkr.Types{}

		for _, v := range typee.Params {
			params = append(params, TypeExpr(w, v))
		}

		returns := wkr.Types{}
		for _, v := range typee.Returns {
			returns = append(returns, TypeExpr(w, v))
		}

		typ = &wkr.FunctionType{
			Params:  params,
			Returns: returns,
		}
	case ast.Map, ast.List:
		wrapped := TypeExpr(w, typee.WrappedType)
		typ = wkr.NewWrapperType(wkr.NewBasicType(pvt), wrapped)
	case ast.Entity:
		typ = &wkr.RawEntityType{}
	default:
		typeeName := typee.Name.GetToken().Lexeme
		if entityVal, found := w.CurrentEnvironment.Entities[typeeName]; found {
			typ = entityVal.GetType()
			break
		}
		if structVal, found := w.CurrentEnvironment.Structs[typeeName]; found {
			typ = structVal.GetType()
			break
		}
		if customType, found := w.CurrentEnvironment.CustomTypes[typeeName]; found {
			typ = customType
			break
		}
		if val := w.GetVariable(&w.CurrentEnvironment.Scope, typeeName); val != nil {
			if val.GetType().PVT() == ast.Enum {
				typ = val.GetType()
				break
			}
		}
		typ = wkr.InvalidType
	}

	if typee.IsVariadic {
		return wkr.NewVariadicType(typ)
	}
	return typ
}
