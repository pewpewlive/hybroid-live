package pass1

import (
	"hybroid/ast"
	"hybroid/helpers"
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

	WalkBody(w, &fn.Body, funcTag, fnScope)

	return &funcTag.ReturnTypes
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
			Owner:      selfExpr,
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
			Owner:      selfExpr,
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

func CallExpr(w *wkr.Walker, node *ast.CallExpr, scope *wkr.Scope, callType wkr.ProcedureType) wkr.Value {
	val := GetNodeValue(w, &node.Caller, scope)

	valType := val.GetType().PVT()
	if valType != ast.Func {
		return &wkr.Unknown{}
	}

	for i := range node.Args {
		GetNodeValue(w, &node.Args[i], scope)
	}

	variable, it_is := val.(*wkr.VariableVal)
	if it_is {
		val = variable.Value
	}
	fun, _ := val.(*wkr.FunctionVal)

	if len(fun.Returns) == 1 {
		return w.TypeToValue(fun.Returns[0])
	}
	return &fun.Returns
}

func MethodCallExpr(w *wkr.Walker, node *ast.Node, scope *wkr.Scope) wkr.Value {
	method := (*node).(*ast.MethodCallExpr)

	ownerVal := GetNodeValue(w, &method.Owner, scope)

	if container := helpers.GetValOfInterface[wkr.FieldContainer](ownerVal); container != nil {
		container := *container
		if _, _, contains := container.ContainsField(method.MethodName); contains {
			expr := ast.CallExpr{
				Caller: method.Call,
				Args:   method.Args,
			}
			val := CallExpr(w, &expr, scope, wkr.Function)
			*node = &expr
			return val
		}
	}

	method.TypeName = ownerVal.GetType().(*wkr.NamedType).Name
	*node = method

	callExpr := ast.CallExpr{
		Caller: method.Call,
		Args:   method.Args,
	}

	return CallExpr(w, &callExpr, scope, wkr.Method)
}

func FieldExpr(w *wkr.Walker, node *ast.FieldExpr, scope *wkr.Scope) wkr.Value {
	if node.Owner == nil {
		val := GetNodeValue(w, &node.Identifier, scope)

		if val.GetType().PVT() == ast.Unresolved {
			return &wkr.Unknown{}
		}

		if !wkr.IsOfPrimitiveType(val, ast.Struct, ast.Entity, ast.AnonStruct, ast.Enum) {
			return &wkr.Unknown{}
		}

		var fieldVal wkr.Value
		if node.Property == nil {
			return val
		} else {
			w.Context.Value = val
			w.Context.Node = node.Identifier
			fieldVal = GetNodeValue(w, &node.Property, scope)
		}
		return fieldVal
	}
	owner := w.Context.Value
	variable := &wkr.VariableVal{Value: &wkr.Invalid{}}
	ident := node.Identifier.GetToken()
	var isField, isMethod bool
	if container, is := owner.(wkr.FieldContainer); is {
		field, index, containsField := container.ContainsField(ident.Lexeme)
		isField = containsField
		if containsField {
			node.Index = index
			variable = field
		}
	}
	if container, is := owner.(wkr.MethodContainer); is && !isField {
		method, containsMethod := container.ContainsMethod(ident.Lexeme)
		isMethod = containsMethod
		if isMethod {
			node.Index = -1
			variable = method
		}
	}

	if node.Property != nil {
		w.Context.Value = variable.Value
		val := GetNodeValue(w, &node.Property, scope)
		return val
	}

	return variable.Value
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

func MemberExpr(w *wkr.Walker, node *ast.MemberExpr, scope *wkr.Scope) wkr.Value {
	if node.Owner == nil {
		val := GetNodeValue(w, &node.Identifier, scope)

		if val.GetType().PVT() == ast.Unresolved {
			return &wkr.Unknown{}
		}

		var memberVal wkr.Value
		if node.Property == nil {
			return val
		} else {
			w.Context.Value = val
			memberVal = GetNodeValue(w, &node.Property, scope)
		}
		return memberVal
	}

	//val := GetNodeValue(w, &node.Identifier, scope)
	//valType := val.GetType().PVT()
	//array := w.Context.Value
	//arrayType := array.GetType().PVT()

	// if arrayType == ast.Map {
	// 	if valType != ast.String && valType != ast.Unknown {
	// 		w.Error(node.Identifier.GetToken(), "variable is not a string")
	// 		return &wkr.Invalid{}
	// 	}
	// } else if arrayType == ast.List {
	// 	if valType != ast.Number && valType != ast.Unknown {
	// 		w.Error(node.Identifier.GetToken(), "variable is not a number")
	// 		return &wkr.Invalid{}
	// 	}
	// }

	// if arrayType != ast.List && arrayType != ast.Map {
	// 	w.Error(node.Identifier.GetToken(), fmt.Sprintf("variable '%s' is not a list, or map", node.Identifier.GetToken().Lexeme))
	// 	return &wkr.Invalid{}
	// }

	// if variable, ok := array.(*wkr.VariableVal); ok {
	// 	array = variable.Value
	// }

	//wrappedValType := array.GetType().(*wkr.WrapperType).WrappedType
	//wrappedVal := w.TypeToValue(wrappedValType)

	// if node.Property != nil {
	// 	w.Context.Value = wrappedVal
	// 	return GetNodeValue(w, &node.Property, scope)
	// }

	//return wrappedVal
	return &wkr.Unknown{}
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
	if new.Type.GetType() == ast.Identifier {
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
