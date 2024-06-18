package walker

import (
	"fmt"
	"hybroid/ast"
	"hybroid/helpers"
	"hybroid/lexer"
	"hybroid/parser"
)

func (w *Walker) DetermineValueType(left Type, right Type) Type {
	if left.PVT() == 0 || right.PVT() == 0 {
		return NAType
	}
	if TypeEquals(left, right) {
		return right
	}
	if parser.IsFx(left.PVT()) && parser.IsFx(right.PVT()) {
		return left
	}

	return InvalidType
}

func (w *Walker) BinaryExpr(node *ast.BinaryExpr, scope *Scope) Value {
	left, right := w.GetNodeValue(&node.Left, scope), w.GetNodeValue(&node.Right, scope)
	leftType, rightType := left.GetType(), right.GetType()
	op := node.Operator
	switch op.Type {
	case lexer.Plus, lexer.Minus, lexer.Caret, lexer.Star, lexer.Slash, lexer.Modulo:
		w.validateArithmeticOperands(leftType, rightType, *node)
	default:
		if !TypeEquals(leftType, rightType) {
			w.Error(node.GetToken(), fmt.Sprintf("invalid comparison: types are not the same (left: %s, right: %s)",leftType.ToString(), rightType.ToString()))
		} else {
			return &BoolVal{}
		}
	}
	typ := w.DetermineValueType(leftType, rightType)

	if typ.PVT() == ast.Invalid {
		w.Error(node.GetToken(), fmt.Sprintf("invalid binary expression (left: %s, right: %s)",leftType.ToString(), rightType.ToString()))
		return &Invalid{}
	} else {
		return &BoolVal{}
	}
}

func (w *Walker) LiteralExpr(node *ast.LiteralExpr) Value {
	switch node.ValueType {
	case ast.String:
		return &StringVal{}
	case ast.Fixed:
		return &FixedVal{ast.Fixed}
	case ast.Radian:
		return &FixedVal{ast.Radian}
	case ast.FixedPoint:
		return &FixedVal{ast.FixedPoint}
	case ast.Degree:
		return &FixedVal{ast.Degree}
	case ast.Bool:
		return &BoolVal{}
	case ast.Number:
		return &NumberVal{}
	default:
		return &Invalid{}
	}
}

func (w *Walker) IdentifierExpr(node *ast.Node, scope *Scope) Value {
	valueNode := *node
	ident := valueNode.(*ast.IdentifierExpr)

	sc := scope.ResolveVariable(ident.Name.Lexeme)
	if sc == nil {
		return &Invalid{}
	}

	variable := sc.GetVariable(ident.Name.Lexeme)

	if sc.Tag.GetType() == Struct {
		_struct := sc.Tag.(*StructTag).StructVal
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

func (w *Walker) GroupingExpr(node *ast.GroupExpr, scope *Scope) Value {
	return w.GetNodeValue(&node.Expr, scope)
}

func (w *Walker) ListExpr(node *ast.ListExpr, scope *Scope) Value {
	var value ListVal
	for i := range node.List {
		val := w.GetNodeValue(&node.List[i], scope)
		if val.GetType().PVT() == ast.Invalid {
			w.Error(node.List[i].GetToken(), fmt.Sprintf("variable '%s' inside list is invalid", node.List[i].GetToken().Lexeme))
		}
		value.Values = append(value.Values, val)
	}
	value.ValueType = GetContentsValueType(value.Values)
	return &value
}

type ProcedureType int

const (
	Function ProcedureType = iota
	Method
)

func (w *Walker) DetermineCallTypeString(callType ProcedureType) string {
	if callType == Function {
		return "function"
	}

	return "method"
}

func (w *Walker) ValidateArguments(args []Type, params []Type, callToken lexer.Token, typeCall string) (int, bool) {
	if len(params) < len(args) {
		w.Error(callToken, fmt.Sprintf("too many arguments given in %s call", typeCall))
		return -1, true
	}
	if len(params) > len(args) {
		w.Error(callToken, fmt.Sprintf("too few arguments given in %s call", typeCall))
		return -1, true
	}
	for i, typeVal := range args {
		if !TypeEquals(typeVal, params[i]) {
			return i, false
		}
	}
	return -1, true
}

func (w *Walker) TypeifyNodeList(nodes *[]ast.Node, scope *Scope) []Type {
	arguments := make([]Type, 0)
	for i := range *nodes {
		val := w.GetNodeValue(&(*nodes)[i], scope)
		if function, ok := val.(*FunctionVal); ok {
			arguments = append(arguments, function.Returns...)
		} else {
			arguments = append(arguments, val.GetType())
		}
	}
	return arguments
}

func (w *Walker) CallExpr(node *ast.CallExpr, scope *Scope, callType ProcedureType) Value {
	typeCall := w.DetermineCallTypeString(callType)

	callerToken := node.Caller.GetToken()
	val := w.GetNodeValue(&node.Caller, scope)

	valType := val.GetType().PVT()
	if valType != ast.Func {
		if valType != ast.Invalid {
			w.Error(callerToken, fmt.Sprintf("variable used as if it's a %s (type: %s)", typeCall, valType.ToString()))
		} else {
			w.Error(callerToken, fmt.Sprintf("unkown %s", typeCall))
		}
		return &Invalid{}
	}

	variable, it_is := val.(*VariableVal)
	if it_is {
		val = variable.Value
	}
	fun, _ := val.(*FunctionVal)

	arguments := w.TypeifyNodeList(&node.Args, scope)
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

func (w *Walker) MethodCallExpr(node *ast.Node, scope *Scope) Value {
	method := (*node).(*ast.MethodCallExpr)

	ownerVal := w.GetNodeValue(&method.Owner, scope)

	if container := helpers.GetValOfInterface[FieldContainer](ownerVal); container != nil {
		container := *container
		if _, _, contains := container.ContainsField(method.MethodName); contains {
			expr := ast.CallExpr{
				Identifier: method.MethodName,
				Caller:     method.Call,
				Args:       method.Args,
				Token:      method.Token,
			}
			val := w.CallExpr(&expr, scope, Function)
			*node = &expr
			return val
		}
	}

	method.TypeName = ownerVal.GetType().(*NamedType).Name
	*node = method

	callExpr := ast.CallExpr{
		Identifier: method.TypeName,
		Caller:     method.Call,
		Args:       method.Args,
		Token:      method.Token,
	}

	return w.CallExpr(&callExpr, scope, Method)
}

func IsOfPrimitiveType(value Value, types ...ast.PrimitiveValueType) bool {
	if types == nil {
		return false
	}
	valType := value.GetType().PVT()
	for _, prim := range types {
		if valType == prim {
			return true
		}
	}

	return false
}

func (w *Walker) FieldExpr(node *ast.FieldExpr, scope *Scope) Value {
	if node.Owner == nil {
		val := w.GetNodeValue(&node.Identifier, scope)

		if !IsOfPrimitiveType(val, ast.Struct, ast.Entity, ast.AnonStruct, ast.Enum) {
			w.Error(node.Identifier.GetToken(), fmt.Sprintf("variable '%s' is not a struct, entity, enum or anonymous struct", node.Identifier.GetToken().Lexeme))
			return &Invalid{}
		}

		var fieldVal Value
		if node.Property == nil {
			return val
		} else {
			w.Context.Value = val
			w.Context.Node = node.Identifier
			fieldVal = w.GetNodeValue(&node.Property, scope)
		}
		return fieldVal
	}
	owner := w.Context.Value
	variable := &VariableVal{Value: &Invalid{}}
	ident := node.Identifier.GetToken()
	var isField, isMethod bool
	if container, is := owner.(FieldContainer); is {
		field, index, containsField := container.ContainsField(ident.Lexeme)
		isField = containsField
		if containsField {
			node.Index = index
			variable = field
		}
	}
	if container, is := owner.(MethodContainer); is && !isField {
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
		val := w.GetNodeValue(&node.Property, scope)
		return val
	}

	return variable.Value
}

func (w *Walker) MapExpr(node *ast.MapExpr, scope *Scope) Value {
	mapVal := MapVal{Members: []Value{}}
	for _, v := range node.Map {
		val := w.GetNodeValue(&v.Expr, scope)
		mapVal.Members = append(mapVal.Members, val)
	}
	mapVal.MemberType = GetContentsValueType(mapVal.Members)
	return &mapVal
}

func (w *Walker) UnaryExpr(node *ast.UnaryExpr, scope *Scope) Value {
	return w.GetNodeValue(&node.Value, scope)
}

func (w *Walker) MemberExpr(node *ast.MemberExpr, scope *Scope) Value {
	if node.Owner == nil {
		val := w.GetNodeValue(&node.Identifier, scope)

		var memberVal Value
		if node.Property == nil {
			return val
		} else {
			w.Context.Value = val
			memberVal = w.GetNodeValue(&node.Property, scope)
		}
		return memberVal
	}

	val := w.GetNodeValue(&node.Identifier, scope)
	valType := val.GetType().PVT()
	array := w.Context.Value
	arrayType := array.GetType().PVT()

	if arrayType == ast.Map {
		if valType != ast.String && valType != 0 {
			w.Error(node.Identifier.GetToken(), "variable is not a string")
			return &Invalid{}
		}
	} else if arrayType == ast.List {
		if valType != ast.Number && valType != 0 {
			w.Error(node.Identifier.GetToken(), "variable is not a number")
			return &Invalid{}
		}
	}

	if arrayType != ast.List && arrayType != ast.Map {
		w.Error(node.Identifier.GetToken(), fmt.Sprintf("variable '%s' is not a list, or map", node.Identifier.GetToken().Lexeme))
		return &Invalid{}
	}

	if variable, ok := array.(*VariableVal); ok {
		array = variable.Value
	}

	wrappedValType := array.GetType().(*WrapperType).WrappedType
	wrappedVal := w.TypeToValue(wrappedValType)

	if node.Property != nil {
		w.Context.Value = wrappedVal
		return w.GetNodeValue(&node.Property, scope)
	}

	return wrappedVal
}

func (w *Walker) DirectiveExpr(node *ast.DirectiveExpr, scope *Scope) *DirectiveVal {

	if node.Identifier.Lexeme != "Environment" {
		variable := w.GetNodeValue(&node.Expr, scope)
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
	return &DirectiveVal{}
}

func (w *Walker) SelfExpr(self *ast.SelfExpr, scope *Scope) Value {
	if !scope.Is(SelfAllowing) {
		w.Error(self.Token, "can't use self outside of struct/entity")
		return &Invalid{}
	}

	sc, _, structTag := ResolveTagScope[*StructTag](scope) // TODO: CHECK FOR ENTITY SCOPE

	if sc != nil {
		(*self).Type = ast.SelfStruct
		return (*structTag).StructVal
	} else {
		return &Invalid{}
	}
}

func (w *Walker) NewExpr(new *ast.NewExpr, scope *Scope) Value {
	w.Context.Node = new
	val, found := w.GetStruct(new.Type.Lexeme)
	if !found {
		return val
	}
	structVal := val.(*StructVal)

	args := w.TypeifyNodeList(&new.Args, scope)
	index, failed := w.ValidateArguments(args, structVal.Params, new.Type, "new")
	if !failed {
		argToken := new.Args[index].GetToken()
		w.Error(argToken, fmt.Sprintf("mismatched types: argument '%s' is not of expected type %s", argToken.Lexeme, structVal.Params[index].ToString()))
	}

	return structVal
}

func (w *Walker) TypeExpr(typee *ast.TypeExpr) Type {
	if typee == nil {
		return InvalidType
	}
	pvt := w.GetTypeFromString(typee.Name.Lexeme)
	switch pvt {
	case ast.Bool, ast.String, ast.Number, ast.Fixed, ast.FixedPoint, ast.Radian, ast.Degree:
		return NewBasicType(pvt)
	case ast.Enum:
		return NewBasicType(ast.Enum)
	case ast.AnonStruct:
		fields := map[string]*VariableVal{}

		for _, v := range typee.Fields {
			fields[v.Name.Lexeme] = &VariableVal{
				Name: v.Name.Lexeme,
				Value: w.TypeToValue(w.TypeExpr(v.Type)),
				Token: v.Name,
			}
		}

		return &AnonStructType{
			Fields: fields,
		}
	case ast.Func:
		params := Types{}

		for _, v := range typee.Params {
			params = append(params, w.TypeExpr(v))
		}

		returns := Types{}
		for _, v := range typee.Returns {
			returns = append(returns, w.TypeExpr(v))
		}

		return &FunctionType{
			Params: params,
			Returns: returns,
		}
	default:
		if structVal, found := w.Environment.Structs[typee.Name.Lexeme]; found {
			return structVal.GetType()
		}
		if val := w.Environment.Scope.GetVariable(typee.Name.Lexeme); val != nil {
			if val.GetType().PVT() == ast.Enum {
				return val.GetType()
			}
		}
		return InvalidType
	}
}
