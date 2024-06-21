package pass3

import (
	"fmt"
	"hybroid/ast"
	"hybroid/helpers"
	"hybroid/lexer"
	"hybroid/parser"
	wkr "hybroid/walker"
)

func AnonFnExpr(w *wkr.Walker, fn *ast.AnonFnExpr, scope *wkr.Scope) wkr.Value {
	fnScope := scope.AccessChild()
	funcTag := fnScope.Tag.(*wkr.FuncTag)
	WalkBody(w, &fn.Body, funcTag, fnScope)

	if !funcTag.GetIfExits(wkr.Return) && !funcTag.ReturnType.Eq(&wkr.EmptyReturn) {
		w.Error(fn.GetToken(), "not all code paths return a value")
	}

	return &funcTag.ReturnType
}

func AnonStructExpr(w *wkr.Walker, node *ast.AnonStructExpr, scope *wkr.Scope) *wkr.AnonStructVal {
	structTypeVal := &wkr.AnonStructVal{
		Fields: make(map[string]*wkr.VariableVal),
	}

	for i := range node.Fields {
		FieldDeclarationStmt(w, node.Fields[i], structTypeVal, scope)
	}

	return structTypeVal
}

func MatchExpr(w *wkr.Walker, node *ast.MatchExpr, scope *wkr.Scope) wkr.Value {
	matchScope := scope.AccessChild()
	mtt := matchScope.Tag.(*wkr.MatchExprTag)

	for i := range node.MatchStmt.Cases {
		caseScope := matchScope.AccessChild()
		WalkBody(w, &node.MatchStmt.Cases[i].Body, mtt, caseScope)
	}

	return mtt.YieldValues
}

func BinaryExpr(w *wkr.Walker, node *ast.BinaryExpr, scope *wkr.Scope) wkr.Value {
	left, right := GetNodeValue(w, &node.Left, scope), GetNodeValue(w, &node.Right, scope)
	leftType, rightType := left.GetType(), right.GetType()
	op := node.Operator
	switch op.Type {
	case lexer.Plus, lexer.Minus, lexer.Caret, lexer.Star, lexer.Slash, lexer.Modulo:
		w.ValidateArithmeticOperands(leftType, rightType, *node)
	default:
		if !wkr.TypeEquals(leftType, rightType) {
			w.Error(node.GetToken(), fmt.Sprintf("invalid comparison: types are not the same (left: %s, right: %s)", leftType.ToString(), rightType.ToString()))
		} else {
			return &wkr.BoolVal{}
		}
	}
	typ := DetermineValueType(w, leftType, rightType)

	if typ.PVT() == ast.Invalid {
		w.Error(node.GetToken(), fmt.Sprintf("invalid binary expression (left: %s, right: %s)", leftType.ToString(), rightType.ToString()))
		return &wkr.Invalid{}
	} else {
		return &wkr.BoolVal{}
	}
}

func EnvExpr(w *wkr.Walker, node *ast.EnvExpr, scope *wkr.Scope) wkr.Value {
	return &wkr.UnresolvedVal{
		Expr: node,
	}
}

func DetermineValueType(w *wkr.Walker, left wkr.Type, right wkr.Type) wkr.Type {
	if left.PVT() == ast.Unknown || right.PVT() == ast.Unknown {
		return wkr.NAType
	}
	if wkr.TypeEquals(left, right) {
		return right
	}
	if parser.IsFx(left.PVT()) && parser.IsFx(right.PVT()) {
		return left
	}

	return wkr.InvalidType
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

	sc := scope.ResolveVariable(ident.Name.Lexeme)
	if sc == nil {
		return &wkr.Invalid{}
	}

	variable := sc.GetVariable(ident.Name.Lexeme)

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
	}
	variable.IsUsed = true
	return variable.Value
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
	typeCall := wkr.DetermineCallTypeString(callType)

	callerToken := node.Caller.GetToken()
	val := GetNodeValue(w, &node.Caller, scope)

	valType := val.GetType().PVT()
	if valType != ast.Func {
		if valType != ast.Invalid {
			w.Error(callerToken, fmt.Sprintf("variable used as if it's a %s (type: %s)", typeCall, valType))
		} else {
			w.Error(callerToken, fmt.Sprintf("unkown %s", typeCall))
		}
		return &wkr.Invalid{}
	}

	variable, it_is := val.(*wkr.VariableVal)
	if it_is {
		val = variable.Value
	}
	fun, _ := val.(*wkr.FunctionVal)

	arguments := pass1.TypeifyNodeList(w, &node.Args, scope)
	index, failed := w.ValidateArguments(arguments, fun.Params, callerToken, typeCall)
	if !failed {
		argToken := node.Args[index].GetToken()
		w.Error(argToken, fmt.Sprintf("mismatched types: argument '%s' is not of expected type %s", argToken.Lexeme, fun.Params[index].ToString()))
	}

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
		val := GetNodeValue(w, &node.Identifier, scope)

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
		if valType != ast.String && valType != ast.Unknown {
			w.Error(node.Identifier.GetToken(), "variable is not a string")
			return &wkr.Invalid{}
		}
	} else if arrayType == ast.List {
		if valType != ast.Number && valType != ast.Unknown {
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

func DirectiveExpr(w *wkr.Walker, node *ast.DirectiveExpr, scope *wkr.Scope) *wkr.DirectiveVal {

	if node.Identifier.Lexeme != "Environment" {
		variable := GetNodeValue(w, &node.Expr, scope)
		variableToken := node.Expr.GetToken()

		variableType := variable.GetType().PVT()
		switch node.Identifier.Lexeme {
		case "Len":
			node.ValueType = ast.Number
			if variableType != ast.Map && variableType != ast.List && variableType != ast.String {
				w.Error(variableToken, "invalid expression in '@Len' directive")
			}
		case "MapToStr":
			node.ValueType = ast.String
			if variableType != ast.Map {
				w.Error(variableToken, "expected a map in '@MapToStr' directive")
			}
		case "ListToStr":
			node.ValueType = ast.List
			if variableType != ast.List {
				w.Error(variableToken, "expected a list in '@ListToStr' directive")
			}
		default:
			// TODO: Implement custom directives

			w.Error(node.Token, "unknown directive")
		}

	} else {

		ident, ok := node.Expr.(*ast.IdentifierExpr)
		if !ok {
			w.Error(node.Expr.GetToken(), "expected an identifier in '@Environment' directive")
		} else {
			name := ident.Name.Lexeme
			if name != "Level" && name != "Mesh" && name != "Sound" && name != "Shared" && name != "LuaGeneric" {
				w.Error(node.Expr.GetToken(), "invalid identifier in '@Environment' directive")
			}
		}
	}
	return &wkr.DirectiveVal{}
}

func SelfExpr(w *wkr.Walker, self *ast.SelfExpr, scope *wkr.Scope) wkr.Value {
	if !scope.Is(wkr.SelfAllowing) {
		w.Error(self.Token, "can't use self outside of struct/entity")
		return &wkr.Invalid{}
	}

	sc, _, structTag := wkr.ResolveTagScope[*wkr.StructTag](scope) // TODO: CHECK FOR ENTITY SCOPE

	if sc != nil {
		(*self).Type = ast.SelfStruct
		return (*structTag).StructVal
	} else {
		return &wkr.Invalid{}
	}
}

func NewExpr(w *wkr.Walker, new *ast.NewExpr, scope *wkr.Scope) wkr.Value {
	w.Context.Node = new
	val, found := w.GetStruct(new.Type.Lexeme)
	if !found {
		return val
	}
	structVal := val.(*wkr.StructVal)

	args := pass1.TypeifyNodeList(w, &new.Args, scope)
	index, failed := w.ValidateArguments(args, structVal.Params, new.Type, "new")
	if !failed {
		argToken := new.Args[index].GetToken()
		w.Error(argToken, fmt.Sprintf("mismatched types: argument '%s' is not of expected type %s", argToken.Lexeme, structVal.Params[index].ToString()))
	}

	return structVal
}

// lets keep it for now in the back of our mind
func TypeExpr(w *wkr.Walker, typee *ast.TypeExpr) wkr.Type {
	if typee == nil {
		return wkr.InvalidType
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
				Value: w.TypeToValue(TypeExpr(w, v.Type)),
				Token: v.Name,
			}
		}

		return &wkr.AnonStructType{
			Fields: fields,
		}
	case ast.Func:
		params := wkr.Types{}

		for _, v := range typee.Params {
			params = append(params, TypeExpr(w, v))
		}

		returns := wkr.Types{}
		for _, v := range typee.Returns {
			returns = append(returns, TypeExpr(w, v))
		}

		return &wkr.FunctionType{
			Params:  params,
			Returns: returns,
		}
	default:
		if structVal, found := w.Environment.Structs[typee.Name.GetToken().Lexeme]; found {
			return structVal.GetType()
		}
		if val := w.Environment.Scope.GetVariable(typee.Name.GetToken().Lexeme); val != nil {
			if val.GetType().PVT() == ast.Enum {
				return val.GetType()
			}
		}
		return wkr.InvalidType
	}
}
