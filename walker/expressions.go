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

// func (w *Walker) selfExpr(node *ast.SelfExpr, scope *Scope) Value {

// }

func (w *Walker) callExpr(node *ast.CallExpr, scope *Scope) Value {
	callerToken := node.Caller.GetToken()
	sc := scope.ResolveVariable(callerToken.Lexeme)

	if sc == nil { //make sure in the future member calls are also taken into account
		w.error(node.Token, "undeclared function")
		return Invalid{}
	} else {
		fn := sc.GetVariable(sc, callerToken.Lexeme)
		fun, ok := fn.Value.(FunctionVal)

		arguments := make([]TypeVal, 0)
		for _, arg := range node.Args {
			val := w.GetNodeValue(&arg, scope)
			if function, ok := val.(FunctionVal); ok {
				arguments = append(arguments, function.returnVal.values...)
			} else {
				arguments = append(arguments, val.GetType())
			}
		}
		if !ok {
			w.error(callerToken, "variable used as if it's a function")
		} else if len(fun.params) < len(arguments) {
			w.error(callerToken, "too many arguments given in function call")
		} else if len(fun.params) > len(arguments) {
			w.error(callerToken, "too few arguments given in function call")
		}

		return CallVal{types: fun.returnVal}
	}
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
		sc := scope.ResolveVariable(node.Identifier.GetToken().Lexeme)

		var array Value
		if sc == nil {
			w.error(node.Identifier.GetToken(), fmt.Sprintf("undeclared variable \"%s\"", node.Identifier.GetToken().Lexeme))
			return Invalid{}
		} else {
			array = sc.GetVariable(sc, node.Identifier.GetToken().Lexeme).Value
		}

		next, ok := node.Property.(ast.MemberExpr)

		if ok {
			arrayType := array.GetType()
			if arrayType.Type == ast.Namespace {
				return Unknown{}
			}
			if arrayType.Type != ast.List && arrayType.Type != ast.Map {
				w.error(node.Identifier.GetToken(), "variable is not a list, map or a namespace")
				return Invalid{}
			}
			return w.memberExpr(array, &next, scope)
		} else {
			w.error(node.GetToken(), "expected member expression")
			return Invalid{}
		}
	}

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
		next, ok := node.Property.(ast.MemberExpr)
		if ok {
			return w.memberExpr(w.GetValueFromType(wrappedValType), &next, scope)
		} else {
			w.error(node.GetToken(), "expected member expression")
			return Invalid{}
		}
	}

	return w.GetValueFromType(wrappedValType)
}

func (w *Walker) directiveExpr(node *ast.DirectiveExpr, scope *Scope) DirectiveVal {

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
	if len(ret.values) == 0 {
		ret.values = append(ret.values, TypeVal{Type: ast.Nil})
	}

	for _, node := range fn.Body {
		w.WalkNode(&node, &fnScope)
	}

	if w.bodyReturns(&fn.Body, &ret, &fnScope) == nil && ret.values[0].Type != ast.Nil {
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
