package walker

import (
	"fmt"
	"hybroid/ast"
	"hybroid/helpers"
	"hybroid/lexer"
	"hybroid/parser"
)

func (w *Walker) determineValueType(left Type, right Type) Type {
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

func (w *Walker) binaryExpr(node *ast.BinaryExpr, scope *Scope) Value {
	left, right := w.GetNodeValue(&node.Left, scope), w.GetNodeValue(&node.Right, scope)
	leftType, rightType := left.GetType(), right.GetType()
	op := node.Operator
	switch op.Type {
	case lexer.Plus, lexer.Minus, lexer.Caret, lexer.Star, lexer.Slash, lexer.Modulo:
		w.validateArithmeticOperands(leftType, rightType, *node)
	default:
		if !TypeEquals(leftType, rightType) {
			w.error(node.GetToken(), fmt.Sprintf("invalid comparison: types are not the same (left: %s, right: %s)",leftType.ToString(), rightType.ToString()))
		} else {
			return &BoolVal{}
		}
	}
	typ := w.determineValueType(leftType, rightType)

	if typ.PVT() == ast.Invalid {
		w.error(node.GetToken(), fmt.Sprintf("invalid binary expression (left: %s, right: %s)",leftType.ToString(), rightType.ToString()))
		return &Invalid{}
	} else {
		return &BoolVal{}
	}
}

func (w *Walker) literalExpr(node *ast.LiteralExpr) Value {
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

func (w *Walker) identifierExpr(node *ast.Node, scope *Scope) Value {
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

func (w *Walker) groupingExpr(node *ast.GroupExpr, scope *Scope) Value {
	return w.GetNodeValue(&node.Expr, scope)
}

func (w *Walker) listExpr(node *ast.ListExpr, scope *Scope) Value {
	var value ListVal
	for i := range node.List {
		val := w.GetNodeValue(&node.List[i], scope)
		if val.GetType().PVT() == ast.Invalid {
			w.error(node.List[i].GetToken(), fmt.Sprintf("variable '%s' inside list is invalid", node.List[i].GetToken().Lexeme))
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

func (w *Walker) determineCallTypeString(callType ProcedureType) string {
	if callType == Function {
		return "function"
	}

	return "method"
}

func (w *Walker) validateArguments(args []Type, params []Type, callToken lexer.Token, typeCall string) (int, bool) {
	if len(params) < len(args) {
		w.error(callToken, fmt.Sprintf("too many arguments given in %s call", typeCall))
		return -1, true
	}
	if len(params) > len(args) {
		w.error(callToken, fmt.Sprintf("too few arguments given in %s call", typeCall))
		return -1, true
	}
	for i, typeVal := range args {
		if !TypeEquals(typeVal, params[i]) {
			return i, false
		}
	}
	return -1, true
}

func (w *Walker) typeifyNodeList(nodes *[]ast.Node, scope *Scope) []Type {
	arguments := make([]Type, 0)
	for i := range *nodes {
		val := w.GetNodeValue(&(*nodes)[i], scope)
		if function, ok := val.(*FunctionVal); ok {
			arguments = append(arguments, function.returns...)
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

	valType := val.GetType().PVT()
	if valType != ast.Func {
		if valType != ast.Invalid {
			w.error(callerToken, fmt.Sprintf("variable used as if it's a %s (type: %s)", typeCall, valType.ToString()))
		} else {
			w.error(callerToken, fmt.Sprintf("unkown %s", typeCall))
		}
		return &Invalid{}
	}

	variable, it_is := val.(*VariableVal)
	if it_is {
		val = variable.Value
	}
	fun, _ := val.(*FunctionVal)

	arguments := w.typeifyNodeList(&node.Args, scope)
	index, failed := w.validateArguments(arguments, fun.params, callerToken, typeCall)
	if !failed {
		argToken := node.Args[index].GetToken()
		w.error(argToken, fmt.Sprintf("mismatched types: argument '%s' is not of expected type %s", argToken.Lexeme, fun.params[index].ToString()))
	}

	if len(fun.returns) == 1 {
		return w.TypeToValue(fun.returns[0])
	}
	return &fun.returns
}

func (w *Walker) methodCallExpr(node *ast.Node, scope *Scope) Value {
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
			val := w.callExpr(&expr, scope, Function)
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

	return w.callExpr(&callExpr, scope, Method)
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

func (w *Walker) fieldExpr(node *ast.FieldExpr, scope *Scope) Value {
	if node.Owner == nil {
		val := w.GetNodeValue(&node.Identifier, scope)

		if !IsOfPrimitiveType(val, ast.Struct, ast.Entity, ast.AnonStruct, ast.Enum) {
			w.error(node.Identifier.GetToken(), fmt.Sprintf("variable '%s' is not a struct, entity, enum or anonymous struct", node.Identifier.GetToken().Lexeme))
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
		w.error(ident, fmt.Sprintf("variable '%s' does not contain '%s'", w.Context.Node.GetToken().Lexeme, ident.Lexeme))
	}

	if node.Property != nil {
		w.Context.Value = variable.Value
		val := w.GetNodeValue(&node.Property, scope)
		return val
	}

	return variable.Value
}

func (w *Walker) mapExpr(node *ast.MapExpr, scope *Scope) Value {
	mapVal := MapVal{Members: []Value{}}
	for _, v := range node.Map {
		val := w.GetNodeValue(&v.Expr, scope)
		mapVal.Members = append(mapVal.Members, val)
	}
	mapVal.MemberType = GetContentsValueType(mapVal.Members)
	return &mapVal
}

func (w *Walker) unaryExpr(node *ast.UnaryExpr, scope *Scope) Value {
	return w.GetNodeValue(&node.Value, scope)
}

func (w *Walker) memberExpr(node *ast.MemberExpr, scope *Scope) Value {
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
			w.error(node.Identifier.GetToken(), "variable is not a string")
			return &Invalid{}
		}
	} else if arrayType == ast.List {
		if valType != ast.Number && valType != 0 {
			w.error(node.Identifier.GetToken(), "variable is not a number")
			return &Invalid{}
		}
	}

	if arrayType != ast.List && arrayType != ast.Map {
		w.error(node.Identifier.GetToken(), fmt.Sprintf("variable '%s' is not a list, or map", node.Identifier.GetToken().Lexeme))
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

func (w *Walker) directiveExpr(node *ast.DirectiveExpr, scope *Scope) *DirectiveVal {

	if node.Identifier.Lexeme != "Environment" {
		variable := w.GetNodeValue(&node.Expr, scope)
		variableToken := node.Expr.GetToken()

		variableType := variable.GetType().PVT()
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

		ident, ok := node.Expr.(*ast.IdentifierExpr)
		if !ok {
			w.error(node.Expr.GetToken(), "expected an identifier in '@Environment' directive")
		} else {
			name := ident.Name.Lexeme
			if name != "Level" && name != "Mesh" && name != "Sound" && name != "Shared" && name != "LuaGeneric" {
				w.error(node.Expr.GetToken(), "invalid identifier in '@Environment' directive")
			}
		}
	}
	return &DirectiveVal{}
}

func (w *Walker) selfExpr(self *ast.SelfExpr, scope *Scope) Value {
	if !scope.Is(SelfAllowing) {
		w.error(self.Token, "can't use self outside of struct/entity")
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

func (w *Walker) newExpr(new *ast.NewExpr, scope *Scope) Value {
	w.Context.Node = new
	val, found := w.GetStruct(new.Type.Lexeme)
	if !found {
		return val
	}
	structVal := val.(*StructVal)

	args := w.typeifyNodeList(&new.Args, scope)
	index, failed := w.validateArguments(args, structVal.Params, new.Type, "new")
	if !failed {
		argToken := new.Args[index].GetToken()
		w.error(argToken, fmt.Sprintf("mismatched types: argument '%s' is not of expected type %s", argToken.Lexeme, structVal.Params[index].ToString()))
	}

	return structVal
}

func (w *Walker) anonFnExpr(fn *ast.AnonFnExpr, scope *Scope) *FunctionVal {
	ret := EmptyReturn
	for _, typee := range fn.Return {
		ret = append(ret, w.typeExpr(typee))
	}

	funcTag := &FuncTag{ReturnType: ret}
	fnScope := NewScope(scope, funcTag)
	fnScope.Attributes.Add(ReturnAllowing)

	params := make([]Type, 0)
	for i, param := range fn.Params {
		params = append(params, w.typeExpr(param.Type))
		value := w.TypeToValue(params[i])
		fnScope.DeclareVariable(&VariableVal{Name: param.Name.Lexeme, Value: value, Token: param.Name})
	}

	w.WalkBody(&fn.Body, funcTag, &fnScope)

	if !funcTag.GetIfExits(Return) && !ret.Eq(&EmptyReturn) {
		w.error(fn.GetToken(), "not all code paths return a value")
	}

	return &FunctionVal{
		params:    params,
		returns: ret,
	}
}

func (w *Walker) anonStructExpr(node *ast.AnonStructExpr, scope *Scope) *AnonStructVal {
	structTypeVal := &AnonStructVal{// t
		Fields:       map[string]*VariableVal{},
	}

	for i := range node.Fields {
		w.fieldDeclaration(node.Fields[i], structTypeVal, scope)
	}

	return structTypeVal
}

func (w *Walker) matchExpr(node *ast.MatchExpr, scope *Scope) Value {
	casesLength := len(node.MatchStmt.Cases)+1
	if node.MatchStmt.HasDefault {
		casesLength--
	}
	matchScope := NewScope(scope, &MatchExprTag{})
	matchScope.Attributes.Add(YieldAllowing)
	mtt := &MatchExprTag{mpt:NewMultiPathTag(casesLength, matchScope.Attributes...)}
	matchScope.Tag = mtt

	w.match(&node.MatchStmt, true, scope)
	for i := range node.MatchStmt.Cases {
		caseScope := NewScope(&matchScope, &UntaggedTag{})
		w.WalkBody(&node.MatchStmt.Cases[i].Body, mtt, &caseScope)
	}
	returnable := scope.ResolveReturnable()

	matchTag, _ := matchScope.Tag.(*MatchExprTag)
	if matchTag.YieldValues == nil {
		matchTag.YieldValues = &EmptyReturn
	}
	node.ReturnAmount = len(*matchTag.YieldValues)

	if returnable == nil {
		return matchTag.YieldValues
	}

	if !matchTag.GetIfExits(Yield) {
		w.error(node.MatchStmt.GetToken(), "not all cases yield")
		(*returnable).SetExit(false, Yield)
	}else {
		(*returnable).SetExit(true, Yield)
	}
	(*returnable).SetExit(matchTag.GetIfExits(Return), Return)
	(*returnable).SetExit(matchTag.GetIfExits(Break), Break)
	(*returnable).SetExit(matchTag.GetIfExits(Continue), Continue)

	return matchTag.YieldValues
}

func (w *Walker) typeExpr(typee *ast.TypeExpr) Type {
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
				Value: w.TypeToValue(w.typeExpr(v.Type)),
				Token: v.Name,
			}
		}

		return &AnonStructType{
			Fields: fields,
		}
	case ast.Func:
		params := Types{}

		for _, v := range typee.Params {
			params = append(params, w.typeExpr(v))
		}

		returns := Types{}
		for _, v := range typee.Returns {
			returns = append(returns, w.typeExpr(v))
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
