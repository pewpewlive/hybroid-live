package walker

import (
	"fmt"
	"hybroid/ast"
	"hybroid/helpers"
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

func (w *Walker) identifierExpr(node *ast.Node, scope *Scope) Value {
	valueNode := *node
	ident := valueNode.(ast.IdentifierExpr)
	sc := scope.ResolveVariable(ident.Name.Lexeme)

	if sc != nil {
		newValue := sc.GetVariable(sc, ident.Name.Lexeme)

		if sc.Tag.GetType() == Struct {
			varIndex := sc.GetVariableIndex(sc, ident.Name.Lexeme)
			selfExpr := ast.FieldExpr{
				Identifier: ast.SelfExpr{
					Token: valueNode.GetToken(),
					Type:  ast.SelfStruct,
				},
			}

			fieldExpr := ast.FieldExpr{
				Owner:      selfExpr,
				Identifier: valueNode,
				Index:      varIndex,
			}
			selfExpr.Property = fieldExpr
			*node = selfExpr
		}
		return newValue.Value
	}
	return Invalid{}
}

func (w *Walker) groupingExpr(node *ast.GroupExpr, scope *Scope) Value {
	return w.GetNodeValue(&node.Expr, scope)
}

func (w *Walker) listExpr(node *ast.ListExpr, scope *Scope) Value {
	var value ListVal
	for i := range node.List {
		val := w.GetNodeValue(&node.List[i], scope)
		if val.GetType().Type == ast.Invalid {
			w.error(node.List[i].GetToken(), fmt.Sprintf("variable '%s' inside list is invalid", node.List[i].GetToken().Lexeme))
		}
		value.Values = append(value.Values, val)
	}
	value.ValueType = value.GetContentsValueType()
	return value
}

type ProcedureType int

const (
	Function ProcedureType = iota
	Method
)

func (w *Walker) determineCallTypeString(callType ProcedureType) string {
	if callType == Function {
		return "function"
	}

	return "method"
}

func (w *Walker) validateArguments(args []TypeVal, params []TypeVal, callToken lexer.Token, typeCall string) bool {
	if len(params) < len(args) {
		w.error(callToken, fmt.Sprintf("too many arguments given in %s call", typeCall))
		return false
	}
	if len(params) > len(args) {
		w.error(callToken, fmt.Sprintf("too few arguments given in %s call", typeCall))
		return false
	}
	for i, typeVal := range args {
		if !typeVal.Eq(params[i]) {
			return false
		}
	}
	return true
}

func (w *Walker) typeifyNodeList(nodes *[]ast.Node, scope *Scope) []TypeVal {
	arguments := make([]TypeVal, 0)
	for i := range *nodes {
		val := w.GetNodeValue(&(*nodes)[i], scope)
		if function, ok := val.(FunctionVal); ok {
			arguments = append(arguments, function.returnVal...)
		} else {
			arguments = append(arguments, val.GetType())
		}
	}
	return arguments
}

func (w *Walker) callExpr(node *ast.CallExpr, scope *Scope, callType ProcedureType) Value {
	typeCall := w.determineCallTypeString(callType)

	callerToken := node.Caller.GetToken()
	val := w.GetNodeValue(&node.Caller, scope)

	valType := val.GetType().Type
	if valType != ast.Func {
		if valType != ast.Invalid {
			w.error(callerToken, fmt.Sprintf("variable used as if it's a %s (type: %s)", typeCall, val.GetType().Type.ToString()))
		} else {
			w.error(callerToken, fmt.Sprintf("unkown %s", typeCall))
		}
		return Invalid{}
	}

	variable, it_is := val.(VariableVal)
	if it_is {
		val = variable.Value
	}
	fun, _ := val.(FunctionVal)

	arguments := w.typeifyNodeList(&node.Args, scope)
	w.validateArguments(arguments, fun.params, callerToken, typeCall)

	if len(fun.returnVal) == 1 {
		return w.GetValueFromType(fun.returnVal[0])
	}
	return CallVal{types: fun.returnVal}
}

func (w *Walker) methodCallExpr(node *ast.Node, scope *Scope) Value { // ty
	method := (*node).(ast.MethodCallExpr) // ok

	ownerVal := w.GetNodeValue(&method.Owner, scope)

	if container := helpers.GetValOfInterface[Container](ownerVal); container != nil {
		container := *container
		fields := container.GetFields()
		for _, value := range fields {
			if value.Name == method.MethodName {
				expr := ast.CallExpr{
					Identifier: method.MethodName,
					Caller:     method.Call,
					Args:       method.Args,
					Token:      method.Token,
				}
				val := w.callExpr(&expr, scope, Function)
				*node = expr
				return val
			}
		}
	}

	method.TypeName = ownerVal.GetType().Name
	*node = method

	callExpr := ast.CallExpr{
		Identifier: method.TypeName,
		Caller:     method.Call,
		Args:       method.Args,
		Token:      method.Token,
	}

	return w.callExpr(&callExpr, scope, Method)
}

func IsOfPrimitiveType(value Value, types ...ast.PrimitiveValueType) bool {
	if types == nil {
		return false
	}
	valType := value.GetType().Type
	for _, prim := range types {
		if valType == prim {
			return true
		}
	}

	return false
}

func (w *Walker) fieldExpr(node *ast.FieldExpr, scope *Scope) Value {
	if node.Owner == nil {
		val := w.GetNodeValue(&node.Identifier, scope)

		if !IsOfPrimitiveType(val, ast.Struct, ast.Entity, ast.Namespace) {
			w.error(node.Identifier.GetToken(), fmt.Sprintf("variable '%s' is not a struct, entity or namespace", node.Identifier.GetToken().Lexeme))
			return Invalid{}
		}

		var fieldVal Value
		if node.Property == nil {
			return val
		} else {
			scope.Namespace.Ctx.Value = val
			fieldVal = w.GetNodeValue(&node.Property, scope)
		}
		return fieldVal
	}
	ownerr := scope.Namespace.Ctx.Value
	variable := VariableVal{Value: Invalid{}}
	if IsOfPrimitiveType(ownerr, ast.Struct, ast.Entity, ast.Namespace) {
		if container := helpers.GetValOfInterface[Container](ownerr); container != nil {
			container := *container
			ident := node.Identifier.GetToken()
			val, index, contains := container.Contains(ident.Lexeme)

			if !contains {
				w.error(ident, fmt.Sprintf("no field or method named '%s' in '%s'", ident.Lexeme, node.Owner.GetToken().Lexeme))
				return Invalid{}
			} else {
				variable = val.(VariableVal)
				node.Index = index
			}
		}
	}

	if node.Property != nil {
		scope.Namespace.Ctx.Value = variable.Value
		val := w.GetNodeValue(&node.Property, scope)
		return val
	}

	return variable.Value
}

func (w *Walker) mapExpr(node *ast.MapExpr, scope *Scope) Value {
	mapVal := MapVal{Members: map[string]VariableVal{}}
	for k, v := range node.Map {
		//fmt.Printf("%s, ",v.Type.ToString())
		val := w.GetNodeValue(&v.Expr, scope)

		mapVal.Members[k.Lexeme] = VariableVal{
			Name:  k.Lexeme,
			Value: val,
			Node:  v.Expr,
		}
	}
	mapVal.MemberType = mapVal.GetContentsValueType()
	return mapVal
}

func (w *Walker) unaryExpr(node *ast.UnaryExpr, scope *Scope) Value {
	return w.GetNodeValue(&node.Value, scope)
}

func (w *Walker) memberExpr(node *ast.MemberExpr, scope *Scope) Value {
	if node.Owner == nil {
		val := w.GetNodeValue(&node.Identifier, scope)
		valType := val.GetType().Type

		if valType != ast.List && valType != ast.Map {
			w.error(node.Identifier.GetToken(), fmt.Sprintf("variable '%s' is not a list, map", node.Identifier.GetToken().Lexeme))
			return Invalid{}
		}

		var memberVal Value
		if node.Property == nil {
			return val
		} else {
			scope.Namespace.Ctx.Value = val
			memberVal = w.GetNodeValue(&node.Property, scope)
		}
		return memberVal
	}

	val := w.GetNodeValue(&node.Identifier, scope)
	valType := val.GetType()
	array := scope.Namespace.Ctx.Value
	arrayType := array.GetType()

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

	wrappedValType := TypeVal{Type: ast.Invalid}

	if variable, ok := array.(VariableVal); ok {
		array = variable.Value
	}

	if list, ok := array.(ListVal); ok {
		wrappedValType = list.ValueType
	} else if mapp, ok := array.(MapVal); ok {
		wrappedValType = mapp.MemberType
	}

	wrappedVal := w.GetValueFromType(wrappedValType)

	if node.Property != nil {
		scope.Namespace.Ctx.Value = wrappedVal
		return w.GetNodeValue(&node.Property, scope)
	}

	return wrappedVal
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

func (w *Walker) selfExpr(self *ast.SelfExpr, scope *Scope) Value {
	if !scope.Is(SelfAllowing) {
		w.error(self.Token, "can't use self outside of struct/entity")
		return Invalid{}
	}

	sc, _, structTag := ResolveTagScope[StructTag](scope) // TODO: CHECK FOR ENTITY SCOPE

	if sc != nil {
		(*self).Type = ast.SelfStruct
		return StructVal{Type: structTag.StructType}
	} else {
		return Invalid{}
	}
}

func (w *Walker) newExpr(new *ast.NewExpr, scope *Scope) StructVal {
	resolved := scope.ResolveStructType(new.Type.Lexeme)

	structTypeVal := resolved.GetStructType(resolved, new.Type.Lexeme)

	args := w.typeifyNodeList(&new.Args, scope)
	w.validateArguments(args, structTypeVal.Params, new.Type, "new")

	return StructVal{
		Type: structTypeVal,
	}
}

func (w *Walker) anonFnExpr(fn *ast.AnonFnExpr, scope *Scope) FunctionVal {
	ret := EmptyReturn
	for _, typee := range fn.Return {
		ret = append(ret, w.typeExpr(&typee))
	}

	fnScope := NewScope(scope, FuncTag{ReturnType: ret})
	fnScope.Attributes.Add(ReturnAllowing)

	params := make([]TypeVal, 0)
	for i, param := range fn.Params {
		params = append(params, w.typeExpr(&param.Type))
		value := w.GetValueFromType(params[i])
		fnScope.DeclareVariable(VariableVal{Name: param.Name.Lexeme, Value: value, Node: fn})
	}

	/*
		if len(ret.values) == 0 {
			ret.values = append(ret.values, TypeVal{Type: ast.Nil})
		}*/

	for _, node := range fn.Body {
		w.WalkNode(&node, &fnScope)
	}

	return FunctionVal{
		params:    params,
		returnVal: ret,
	}
}

func (w *Walker) matchExpr(node *ast.MatchExpr, scope *Scope) Returns {
	matchScope := NewScope(scope, MatchExprTag{})
	matchScope.Attributes.Add(YieldAllowing)

	w.matchStmt(&node.MatchStmt, true, scope)
	has_default := false
	for i := range node.MatchStmt.Cases {
		if node.MatchStmt.Cases[i].Expression.GetToken().Lexeme == "_" {
			has_default = true
		}
		caseScope := NewScope(&matchScope, UntaggedTag{})

		endIndex := -1
		for j := range node.MatchStmt.Cases[i].Body {
			matchTag, _ := matchScope.Tag.(MatchExprTag)
			if matchTag.ArmsYielded == i+1 {
				w.warn(node.MatchStmt.Cases[i].Body[j].GetToken(), "unreachable code detected")
				endIndex = j
				break
			}
			w.WalkNode(&node.MatchStmt.Cases[i].Body[j], &caseScope)
		}
		if endIndex != -1 {
			node.MatchStmt.Cases[i].Body = node.MatchStmt.Cases[i].Body[:endIndex]
		}
	}
	returnable := helpers.GetValOfInterface[ReturnableTag](scope.Tag)

	matchTag, _ := matchScope.Tag.(MatchExprTag)
	if matchTag.YieldValues == nil {
		matchTag.YieldValues = &EmptyReturn
	}
	node.ReturnAmount = len(*matchTag.YieldValues)

	if returnable == nil {
		return *matchTag.YieldValues
	}

	if !has_default {
		w.error(node.MatchStmt.GetToken(), "not all arms yield a value")
		scope.Tag = (*returnable).SetReturn(false, Return)
	} else {
		caseLength := len(node.MatchStmt.Cases)
		yields := caseLength == matchTag.ArmsYielded
		scope.Tag = (*returnable).SetReturn(yields, Yield)
		if !yields {
			w.error(node.MatchStmt.GetToken(), "not all arms yield a value")
		}
		returns := caseLength == len(matchTag.mpt.Returns)
		yieldss := caseLength == len(matchTag.mpt.Yields)
		breaks := caseLength == len(matchTag.mpt.Breaks)
		continues := caseLength == len(matchTag.mpt.Continues)

		scope.Tag = (*returnable).SetReturn(returns, Return)
		scope.Tag = (*returnable).SetReturn(yieldss, Yield)
		scope.Tag = (*returnable).SetReturn(breaks, Break)
		scope.Tag = (*returnable).SetReturn(continues, Continue)
	}

	return *matchTag.YieldValues
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
	params := []TypeVal{}
	if typee.Params != nil {
		for _, v := range *typee.Params {
			params = append(params, w.typeExpr(&v))
		}
	}

	returns := []TypeVal{}
	for _, v := range typee.Returns {
		returns = append(returns, w.typeExpr(&v))
	}

	typ := w.GetTypeFromString(typee.Name.Lexeme)
	if typ == ast.Invalid {
		if foreignType, ok := w.Namespace.foreignTypes[typee.Name.Lexeme]; ok {
			return foreignType.GetType()
		}
	}
	if typ == ast.Invalid {
		w.error(typee.GetToken(), "invalid type")
	}

	return TypeVal{
		Name:        typ.ToString(),
		Type:        typ,
		WrappedType: wrapped,
		Params:      params,
		Returns:     returns,
	}
}
