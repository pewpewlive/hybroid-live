package pass2

import (
	"fmt"
	"hybroid/ast"
	"hybroid/helpers"
	"hybroid/lexer"
	wkr "hybroid/walker"
)

func AnonStructExpr(w *wkr.Walker, node *ast.AnonStructExpr, scope *wkr.Scope) *wkr.AnonStructVal {
	anonStructScope := scope.AccessChild()
	structTypeVal := &wkr.AnonStructVal{
		Fields: make(map[string]*wkr.VariableVal),
	}

	for i := range node.Fields {
		FieldDeclarationStmt(w, node.Fields[i], structTypeVal, anonStructScope)
	}

	return structTypeVal
}

func AnonFnExpr(w *wkr.Walker, fn *ast.AnonFnExpr, scope *wkr.Scope) wkr.Value {
	returnTypes := wkr.EmptyReturn
	for i := range fn.Return {
		returnTypes = append(returnTypes, TypeExpr(w, fn.Return[i], w.Environment))
	}
	fnScope := scope.AccessChild()

	WalkBody(w, &fn.Body, fnScope)

	return &fnScope.Tag.(*wkr.FuncTag).ReturnTypes
}

func MatchExpr(w *wkr.Walker, node *ast.MatchExpr, scope *wkr.Scope) wkr.Value {
	matchScope := scope.AccessChild()
	matchScope.Tag = &wkr.MatchExprTag{YieldValues: make(wkr.Types, 0)}

	for i := range node.MatchStmt.Cases {
		caseScope := matchScope.AccessChild()
		WalkBody(w, &node.MatchStmt.Cases[i].Body, caseScope)
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

func IdentifierExpr(w *wkr.Walker, node *ast.Node, scope *wkr.Scope) wkr.Value { // we finna
	valueNode := *node
	ident := valueNode.(*ast.IdentifierExpr) // bishwhat

	sc := w.ResolveVariable(scope, ident.Name.Lexeme)
	if sc == nil {
		return &wkr.Invalid{}
	}

	variable := w.GetVariable(sc, ident.Name.Lexeme)
	variable.IsUsed = true
	return variable.Value
}

func EnvAccessExpr(w *wkr.Walker, node *ast.EnvAccessExpr) wkr.Value {
	envName := node.PathExpr.Nameify()

	walker, found := w.Walkers[envName]
	if !found {
		w.Error(node.PathExpr.GetToken(), fmt.Sprintf("Environment named '%s' doesn't exist", envName))
		return &wkr.Invalid{}
	}

	envStmt := w.GetEnvStmt()

	for _, v := range walker.GetEnvStmt().Requirements {
		if v == w.Environment.Path {
			w.Error(node.GetToken(), fmt.Sprintf("import cycle detected: this environment and '%s' are using each other", walker.Environment.Name))
		}
	}

	envStmt.AddRequirement(walker.Environment.Path)

	value := GetNodeValue(w, &node.Accessed, &walker.Environment.Scope)

	if value.GetType().PVT() == ast.Unresolved {
		value = GetNodeValue(walker, &value.(*wkr.UnresolvedVal).Expr, &walker.Environment.Scope)
	}

	return value
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

func CallExpr(w *wkr.Walker, node *ast.CallExpr, scope *wkr.Scope, callType wkr.ProcedureType) wkr.Value {
	val := GetNodeValue(w, &node.Caller, scope)

	valType := val.GetType().PVT()
	if valType != ast.Func {
		return &wkr.Invalid{}
	}

	variable, it_is := val.(*wkr.VariableVal)
	if it_is {
		val = variable.Value
	}
	fun, _ := val.(*wkr.FunctionVal)

	args := []wkr.Type{}
	for i := range node.Args {
		args = append(args, GetNodeValue(w, &node.Args[i], scope).GetType())
	}
	w.ValidateArguments(args, fun.Params, node.Caller.GetToken(), w.DetermineCallTypeString(callType))

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
				Identifier: method.MethodName,
				Caller:     method.Call,
				Args:       method.Args,
				Token:      method.Token,
			}
			val := CallExpr(w, &expr, scope, wkr.Function)
			*node = &expr
			return val
		}
	}

	method.TypeName = ownerVal.GetType().(*wkr.NamedType).Name
	*node = method

	callExpr := ast.CallExpr{
		Identifier: method.TypeName,
		Caller:     method.Call,
		Args:       method.Args,
		Token:      method.Token,
	}

	return CallExpr(w, &callExpr, scope, wkr.Method)
}

func FieldExpr(w *wkr.Walker, node *ast.FieldExpr, scope *wkr.Scope) wkr.Value {
	if node.Owner == nil {
		val := GetNodeValue(w, &node.Identifier, scope)

		if val.GetType().PVT() == ast.Unresolved { // i fixed it
			return &wkr.Unknown{} // OHHHHHH
		}

		if !wkr.IsOfPrimitiveType(val, ast.Struct, ast.Entity, ast.AnonStruct, ast.Enum) {
			w.Error(node.Identifier.GetToken(), fmt.Sprintf("variable '%s' is not a struct, entity, enum or anonymous struct", node.Identifier.GetToken().Lexeme))
			return &wkr.Invalid{}
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
	if !isField && !isMethod {
		w.Error(ident, fmt.Sprintf("variable '%s' does not contain '%s'", w.Context.Node.GetToken().Lexeme, ident.Lexeme))
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
		val := GetNodeValue(w, &node.Identifier, scope) // invalid, can we start again cuz i wanna go into this

		var memberVal wkr.Value
		if node.Property == nil {
			return val
		} else {
			w.Context.Value = val
			memberVal = GetNodeValue(w, &node.Property, scope)
		}
		return memberVal
	}

	val := GetNodeValue(w, &node.Identifier, scope)
	valType := val.GetType().PVT()
	array := w.Context.Value
	arrayType := array.GetType().PVT()

	if arrayType == ast.Map {
		if valType != ast.String {
			w.Error(node.Identifier.GetToken(), "variable is not a string")
			return &wkr.Invalid{}
		}
	} else if arrayType == ast.List {
		if valType != ast.Number {
			w.Error(node.Identifier.GetToken(), "variable is not a number")
			return &wkr.Invalid{}
		}
	}

	if arrayType != ast.List && arrayType != ast.Map {
		w.Error(node.Identifier.GetToken(), fmt.Sprintf("variable '%s' is not a list, or map", node.Identifier.GetToken().Lexeme))
		return &wkr.Invalid{}
	}

	if variable, ok := array.(*wkr.VariableVal); ok {
		array = variable.Value
	}

	wrappedValType := array.GetType().(*wkr.WrapperType).WrappedType
	wrappedVal := w.TypeToValue(wrappedValType)

	if node.Property != nil {
		w.Context.Value = wrappedVal
		return GetNodeValue(w, &node.Property, scope)
	}

	return wrappedVal
}

func SelfExpr(w *wkr.Walker, self *ast.SelfExpr, scope *wkr.Scope) wkr.Value {
	if !scope.Is(wkr.SelfAllowing) {
		w.Error(self.Token, "can't use self outside of struct/entity")
		return &wkr.Invalid{}
	}
	// AHHHHHHH IT HAS ALREADY TURNED INTO A SELF EXPRESSION FROM PASS 1
	sc, _, structTag := wkr.ResolveTagScope[*wkr.StructTag](scope) // TODO: CHECK FOR ENTITY SCOPE

	if sc == nil {
		entitySc, _, entityTag := wkr.ResolveTagScope[*wkr.EntityTag](scope) // try now
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
	var val *wkr.StructVal
	var found bool
	if new.Type.GetType() == ast.Identifier {
		val, found = w.GetStruct(new.Type.GetToken().Lexeme)
	} else {
		return &wkr.Unknown{}
	}
	if !found {
		return &wkr.Invalid{}
	}

	return val
}

func TypeExpr(w *wkr.Walker, typee *ast.TypeExpr, env *wkr.Environment) wkr.Type {
	if typee == nil {
		return wkr.InvalidType
	}
	if typee.Name.GetType() == ast.EnvironmentAccessExpression {
		expr, _ := typee.Name.(*ast.EnvAccessExpr)
		path := expr.PathExpr.Nameify()

		walker, found := w.Walkers[path]
		if !found {
			w.Error(expr.PathExpr.GetToken(), "Environment name so doesn't exist")
			return wkr.InvalidType
		}
		return TypeExpr(w, &ast.TypeExpr{Name: expr.Accessed}, walker.Environment)
	}

	pvt := w.GetTypeFromString(typee.Name.GetToken().Lexeme)
	switch pvt {
	case ast.Bool, ast.String, ast.Number, ast.Fixed, ast.FixedPoint, ast.Radian, ast.Degree:
		return wkr.NewBasicType(pvt)
	case ast.Enum:
		return wkr.NewBasicType(ast.Enum)
	case ast.AnonStruct:
		fields := map[string]*wkr.VariableVal{}

		for _, v := range typee.Fields {
			fields[v.Name.Lexeme] = &wkr.VariableVal{
				Name:  v.Name.Lexeme,
				Value: w.TypeToValue(TypeExpr(w, v.Type, env)),
				Token: v.Name,
			}
		}

		return &wkr.AnonStructType{
			Fields: fields,
		}
	case ast.Func:
		params := wkr.Types{}

		for _, v := range typee.Params {
			params = append(params, TypeExpr(w, v, env))
		}

		returns := wkr.Types{}
		for _, v := range typee.Returns {
			returns = append(returns, TypeExpr(w, v, env))
		}

		return &wkr.FunctionType{
			Params:  params,
			Returns: returns,
		}
	case ast.Map, ast.List:
		wrapped := TypeExpr(w, typee.WrappedType, env)
		return wkr.NewWrapperType(wkr.NewBasicType(pvt), wrapped) // lets goo
	default: // oh wait we're not even checking here for maps and lists, bruh
		if structVal, found := env.Structs[typee.Name.GetToken().Lexeme]; found {
			return structVal.GetType()
		}
		if val := w.GetVariable(&env.Scope, typee.Name.GetToken().Lexeme); val != nil {
			if val.GetType().PVT() == ast.Enum {
				return val.GetType()
			}
		}
		return wkr.InvalidType
	}
}
