package walker

import (
	"fmt"
	"hybroid/ast"
	"hybroid/lexer"
	"hybroid/parser"
)

func (w *Walker) determineValueType(left TypeVal, right TypeVal) TypeVal {
	if left.Type == 0 || right.Type == 0 {
		return TypeVal{Type: 0}
	}
	if left.Eq(right) {
		return right
	}
	if parser.IsFx(left.Type) && parser.IsFx(right.Type) {
		return left
	}

	return TypeVal{Type: ast.Invalid}
}

func (w *Walker) binaryExpr(node *ast.BinaryExpr, scope *Scope) Value {
	left, right := w.GetNodeValue(&node.Left, scope), w.GetNodeValue(&node.Right, scope)
	leftType, rightType := left.GetType(), right.GetType()
	op := node.Operator
	switch op.Type {
	case lexer.Plus, lexer.Minus, lexer.Caret, lexer.Star, lexer.Slash, lexer.Modulo:
		w.validateArithmeticOperands(leftType, rightType, *node)
	default:
		if !leftType.Eq(rightType) {
			w.error(node.GetToken(), fmt.Sprintf("invalid comparison: types are not the same (left: %s, right: %s)", leftType.Type.ToString(), rightType.Type.ToString()))
		} else {
			return BoolVal{}
		}
	}
	val := w.GetValueFromType(w.determineValueType(leftType, rightType))

	if val.GetType().Type == ast.Invalid {
		w.error(node.GetToken(), fmt.Sprintf("invalid binary expression (left: %s, right: %s)", leftType.Type.ToString(), rightType.Type.ToString()))
		return val
	} else {
		return val
	}
}

func (w *Walker) literalExpr(node *ast.LiteralExpr) Value {

	switch node.ValueType {
	case ast.String:
		return StringVal{}
	case ast.Fixed:
		return FixedVal{
			ast.Fixed}
	case ast.Radian:
		return FixedVal{
			ast.Radian}
	case ast.FixedPoint:
		return FixedVal{
			ast.FixedPoint}
	case ast.Degree:
		return FixedVal{
			ast.Degree}
	case ast.Bool:
		return BoolVal{}
	case ast.Nil:
		return NilVal{}
	case ast.Number:
		return NumberVal{}
	default:
		return Invalid{}
	}
}

func (w *Walker) identifierExpr(node *ast.IdentifierExpr, scope *Scope) Value {
	sc := scope.ResolveVariable(node.Name.Lexeme)

	if sc != nil {
		newValue := sc.GetVariable(sc, node.Name.Lexeme)

		/*
			if sc.Type == Structure {
				varIndex := sc.GetVariableIndex(sc, node.Name.Lexeme)

				selfExpr := ast.SelfExpr{
					Token: node.GetToken(),
					Value: node,
					Type:  ast.SelfStruct,
					Index: varIndex,
				}
			} /* else if sc.Type == Entity {
				varIndex := sc.GetVariableIndex(sc, node.Name.Lexeme)

				selfExpr := ast.SelfExpr{
					Token: newValue.Node.GetToken(),
					Value: newValue.Node,
					Type:  ast.SelfEntity,
					Index: varIndex,
				}
			}*/

		//fmt.Printf("%v %s\n", sc.Type, newValue.Name)
		return newValue
	} else {
		//w.error(node.Name, "unknown identifier")
		return Invalid{}
	}
}

func (w *Walker) groupingExpr(node *ast.GroupExpr, scope *Scope) Value {
	return w.GetNodeValue(&node.Expr, scope)
}

func (w *Walker) listExpr(node *ast.ListExpr, scope *Scope) Value {
	var value ListVal
	for _, expr := range node.List {
		value.Values = append(value.Values, w.GetNodeValue(&expr, scope))
	}
	value.ValueType = value.GetContentsValueType()
	return value
}

type CallType int 

const (
	UnknownCall CallType = iota
	Function
	Method
)

func (w *Walker) determineCallTypeString(callType CallType, node *ast.Node, call ast.CallExpr, scope *Scope) string {
	if callType == Function {
		return"function"
	}
	
	if callType == Method {
		return "method"
	}
	
	member, ok := call.Caller.(ast.MemberExpr)

	if !ok {
		return "function"
	}

	val := w.GetNodeValue(&member.Identifier, scope)
	valType := val.GetType()
	if valType.Type != ast.Struct && valType.Type != ast.Entity {
		return "function"
	}

	methodCallExpr := ast.MethodCallExpr{
		TypeName: valType.Name,
		Args: call.Args,
		Caller: call.Caller,
		Token: call.Token,
	}
	*node = methodCallExpr
	return "method"
}

func (w *Walker) callExpr(node ast.Node, scope *Scope, callType CallType) Value {
	callExpr := node.(ast.CallExpr)
	typeCall := w.determineCallTypeString(callType, &node, callExpr, scope)

	callerToken := callExpr.Caller.GetToken()
	val := w.GetNodeValue(&callExpr.Caller, scope)

	if val.GetType().Type != ast.Func {
		w.error(callerToken, fmt.Sprintf("variable used as if it's a %s (type: %s)", typeCall, val.GetType().Type.ToString()))
		return Invalid{} 
	}

	variable, it_is := val.(VariableVal)
	if it_is {
		val = variable.Value
	}
	fun, _ := val.(FunctionVal)

	arguments := make([]TypeVal, 0)
	for _, arg := range callExpr.Args {
		val := w.GetNodeValue(&arg, scope)
		if function, ok := val.(FunctionVal); ok {
			arguments = append(arguments, function.returnVal.values...)
		} else {
			arguments = append(arguments, val.GetType())
		}
	}
	if len(fun.params) < len(arguments) {
		w.error(callerToken, fmt.Sprintf("too many arguments given in %s call", typeCall))
	} else if len(fun.params) > len(arguments) {
		w.error(callerToken, fmt.Sprintf("too few arguments given in %s call", typeCall))
	}

	return CallVal{types: fun.returnVal}
}

func (w *Walker) methodCallExpr(node *ast.MethodCallExpr, scope *Scope) Value {
	// so, i think that uhh, that's really all we need to do here
	if node.Caller.GetType() == ast.SelfExpression { // 
		sc := scope.ResolveStructScope()
		node.TypeName = sc.WrappedType.Name
	}
	callExpr := ast.CallExpr{ 
		Identifier: node.TypeName, 
		Caller:     node.Caller,
		Args:       node.Args, 
		Token:      node.Token,
	}
	
	return w.callExpr(callExpr, scope, Method)
}

func (w *Walker) mapExpr(node *ast.MapExpr, scope *Scope) Value {
	mapVal := MapVal{Members: map[string]MapMemberVal{}}
	for k, v := range node.Map {
		//fmt.Printf("%s, ",v.Type.ToString())
		val := w.GetNodeValue(&v.Expr, scope)

		mapVal.Members[k.Lexeme] = MapMemberVal{
			Var: VariableVal{
				Name:  k.Lexeme,
				Value: val,
				Node:  v.Expr,
			},
			Owner: mapVal,
		}
	}
	mapVal.MemberType = mapVal.GetContentsValueType()
	return mapVal
}

func (w *Walker) unaryExpr(node *ast.UnaryExpr, scope *Scope) Value {
	return w.GetNodeValue(&node.Value, scope)
}

func (w *Walker) memberExpr(array Value, node *ast.MemberExpr, scope *Scope) Value {
	if node.Owner == nil {
		val := w.GetNodeValue(&node.Identifier, scope)//node.Identifier is updated with the index, mhm
		valType := val.GetType().Type// yes but owner wont get updated

		if valType != ast.List && valType != ast.Map && valType != ast.Namespace {
			w.error(node.Identifier.GetToken(), fmt.Sprintf("variable '%s' is not a list, map, or namespace", node.Identifier.GetToken().Lexeme))
			return Invalid{}
		}

		next, ok := (*node.Property).(ast.MemberExpr) // the owner at last is not getting generated or what like it just doesnt 
		// we finna debug oh yeah
		if ok {
			if valType == ast.Namespace { // TODO: RESOLVE THIS CASE IN THE SECOND WALKER STAGE
				return Unknown{} 
			}

			return w.memberExpr(val, &next, scope)
		} else {
			w.error(node.GetToken(), "expected member expression")
			return Invalid{}
		}
	} // lets go
	w.GetNodeValue(node.Owner, scope)// we just walk it we dont really need it for anything else
	val := w.GetNodeValue(&node.Identifier, scope)
	valType := val.GetType()
	arrayType := array.GetType()
	if node.Bracketed {
		if arrayType.Type == ast.Map {
			if valType.Type != ast.String && valType.Type != 0 {
				w.error(node.Identifier.GetToken(), "variable is not a string")
				return Invalid{}
			}
		} else if arrayType.Type == ast.List {
			if valType.Type != ast.Number && valType.Type != 0 {
				w.error(node.Identifier.GetToken(), "variable is not a number")
				return Invalid{}
			}
		}
	} else {
		if arrayType.Type == ast.Map {
			w.error(node.Identifier.GetToken(), "map members are accessed with '[]'")
		}
	}

	wrappedValType := TypeVal{Type: ast.Invalid}
	if list, ok := array.(ListVal); ok {
		wrappedValType = list.ValueType
	} else if mapp, ok := array.(MapVal); ok {
		wrappedValType = mapp.MemberType
	}

	if wrappedValType.Type == ast.Map || wrappedValType.Type == ast.List || wrappedValType.Type == ast.Namespace {
		if node.Property == nil {
			return w.GetValueFromType(wrappedValType)
		} // run itshit
		next, ok := (*node.Property).(ast.MemberExpr)
		if ok {
			return w.memberExpr(w.GetValueFromType(wrappedValType), &next, scope)
		} else {
			w.error(node.GetToken(), "expected member expression")
			return Invalid{}
		}
	}

	return w.GetValueFromType(wrappedValType)
}

func (w *Walker) directiveExpr(node *ast.DirectiveExpr, scope *Scope) DirectiveVal { //hello

	if node.Identifier.Lexeme != "Environment" {
		variable := w.GetNodeValue(&node.Expr, scope)
		variableToken := node.Expr.GetToken()

		variableType := variable.GetType().Type
		switch node.Identifier.Lexeme {
		case "Len":
			node.ValueType = ast.Number
			if variableType != ast.Map && variableType != ast.List && variableType != ast.String {
				w.error(variableToken, "invalid expression in '@Len' directive")
			}
		case "MapToStr":
			node.ValueType = ast.String
			if variableType != ast.Map {
				w.error(variableToken, "expected a map in '@MapToStr' directive")
			}
		case "ListToStr":
			node.ValueType = ast.List
			if variableType != ast.List {
				w.error(variableToken, "expected a list in '@ListToStr' directive")
			}
		default:
			// TODO: Implement custom directives

			w.error(node.Token, "unknown directive")
		}

	} else {

		ident, ok := node.Expr.(ast.IdentifierExpr)
		if !ok {
			w.error(node.Expr.GetToken(), "expected an identifier in '@Environment' directive")
		} else {
			name := ident.Name.Lexeme
			if name != "Level" && name != "Mesh" && name != "Sound" && name != "Shared" && name != "LuaGeneric" {
				w.error(node.Expr.GetToken(), "invalid identifier in '@Environment' directive")
			}
		}
	}
	return DirectiveVal{}
}

func (w *Walker) selfExpr(self *ast.SelfExpr, scope *Scope) Value {
	sc := scope.ResolveStructScope() // TODO: CHECK FOR ENTITY SCOPE

	if sc == nil {
		w.error(self.Token, "can't use self outside of struct/entity")
		return Invalid{}
	}

	if sc.Type == Structure {
		(*self).Type = ast.SelfStruct
		(*self).Index = sc.GetVariableIndex(sc, self.Value.GetToken().Lexeme)
	}

	if _, ok := self.Value.(ast.SelfExpr); ok { // so just "self"
		structTypeVal := sc.GetStructType(sc, sc.WrappedType.Name)
		return StructVal{
			Type: &structTypeVal,
		}
	} else {
		return w.GetNodeValue(&self.Value, sc)
	}
}

func (w *Walker) newExpr(new *ast.NewExpr, scope *Scope) StructVal {
	resolved := scope.ResolveStructType(new.Type.Lexeme)

	structTypeVal := resolved.GetStructType(resolved, new.Type.Lexeme)

	return StructVal{
		Type: &structTypeVal,
	}
}

func (w *Walker) anonFnExpr(fn *ast.AnonFnExpr, scope *Scope) FunctionVal {
	fnScope := NewScope(scope.Global, scope, ReturnAllowing)

	params := make([]TypeVal, 0)
	for i, param := range fn.Params {
		params = append(params, w.typeExpr(&param.Type))
		value := w.GetValueFromType(params[i])
		fnScope.DeclareVariable(VariableVal{Name: param.Name.Lexeme, Value: value, Node: fn})
	}

	var ret ReturnType
	for _, typee := range fn.Return {
		ret.values = append(ret.values, w.typeExpr(&typee))
	}

	/*
		if len(ret.values) == 0 {
			ret.values = append(ret.values, TypeVal{Type: ast.Nil})
		}*/

	for _, node := range fn.Body {
		w.WalkNode(&node, &fnScope)
	}

	if w.bodyReturns(&fn.Body, &ret, &fnScope) == nil && len(ret.values) != 0 {
		w.error(fn.GetToken(), "not all function paths return a value")
	}

	return FunctionVal{
		params:    params,
		returnVal: ret,
	}
}

func (w *Walker) typeExpr(typee *ast.TypeExpr) TypeVal {
	if typee == nil {
		return TypeVal{Type: ast.Invalid}
	}
	var wrapped *TypeVal
	if typee.WrappedType != nil {
		temp := w.typeExpr(typee.WrappedType)
		wrapped = &temp
	}
	params := make([]TypeVal, 0)
	for _, v := range typee.Params {
		params = append(params, w.typeExpr(&v))
	}
	returns := make([]TypeVal, 0)
	for _, v := range typee.Returns {
		returns = append(returns, w.typeExpr(&v))
	}
	//fmt.Printf("%s\n",typee.Name.Lexeme)
	return TypeVal{
		Type:        w.GetTypeFromString(typee.Name.Lexeme),
		WrappedType: wrapped,
		Params:      params,
		Returns:     ReturnType{values: returns},
	}
}
