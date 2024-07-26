package pass2

import (
	"fmt"
	"hybroid/ast"
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

	fnScope.Tag.(*wkr.FuncTag).ReturnTypes = returnTypes

	WalkBody(w, &fn.Body, fnScope)

	params := make([]wkr.Type, 0)
	for i, param := range fn.Params {
		params = append(params, TypeExpr(w, param.Type, w.Environment))
		w.GetVariable(fnScope, param.Name.Lexeme).Value = w.TypeToValue(params[i])
	}
	return &wkr.FunctionVal{ // returnTypes should contain a fn()
		Params:  params,
		Returns: returnTypes,
	}
}

func MatchExpr(w *wkr.Walker, node *ast.MatchExpr, scope *wkr.Scope) wkr.Value {
	matchScope := scope.AccessChild()
	matchScope.Tag = &wkr.MatchExprTag{YieldValues: make(wkr.Types, 0)}

	for i := range node.MatchStmt.Cases {
		caseScope := matchScope.AccessChild()
		GetNodeValue(w, &node.MatchStmt.Cases[i].Expression, matchScope)
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
// let a = Pewpew
// a::
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

	value := GetNodeValue(w, &node.Accessed, &walker.Environment.Scope)

	if value.GetType().PVT() == ast.Unresolved {
		value = GetNodeValue(walker, &value.(*wkr.UnresolvedVal).Expr, &walker.Environment.Scope)
	}

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

	args := []wkr.Type{}
	for i := range node.Args {
		args = append(args, GetNodeValue(w, &node.Args[i], scope).GetType())
	}
	w.ValidateArguments(args, fun.Params, node.Caller.GetToken())

	if len(fun.Returns) == 1 {
		return w.TypeToValue(fun.Returns[0]) 
	}
	return &fun.Returns
}

func PropertyExpr(w *wkr.Walker, node ast.Node, scope *wkr.Scope) wkr.Value {// READS CONTEXT
	var val wkr.Value
	val = &wkr.Invalid{}
	
	switch w.Context.Node.(type) {
	case *ast.FieldExpr:
		owner := w.Context.Value
		if !wkr.IsOfPrimitiveType(owner, ast.Struct, ast.Entity, ast.AnonStruct, ast.Enum) {
			token := w.Context.Node.GetToken()
			w.Error(token, fmt.Sprintf("variable '%s' is not a struct, entity, enum or anonymous struct", token.Lexeme))
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
			token := w.Context.Node.GetToken()
			w.Error(token, fmt.Sprintf("%s is not a map nor a list", token.Lexeme))
			break
		}
		if node.GetType() == ast.CallExpression {
			val := PropertyExpr(w, node.(*ast.CallExpr).Caller, scope)
			return CallExpr(w, val, node.(*ast.CallExpr), scope)
		}
		expr := GetNodeValue(w, &node, scope)
		if expr.GetType().PVT() != ast.String && expr.GetType().PVT() != ast.Number {
			w.Error(node.GetToken(), "value inside brackets must be a string or a number")
		}else {
			ownerType, _ := owner.GetType().(*wkr.WrapperType)
			val = w.TypeToValue(ownerType.WrappedType)
		}
	}

	if val.GetType().PVT() == ast.Invalid {
		w.Error(node.GetToken(), "invalid property")
	}

	return val
}

func FieldExpr(w *wkr.Walker, node *ast.FieldExpr, scope *wkr.Scope) wkr.Value {// WRITES CONTEXT
	var val wkr.Value
	if w.Context.Value.GetType().GetType() != wkr.NA {
		val = PropertyExpr(w, node.Identifier, scope)
	}else {
		val = GetNodeValue(w, &node.Identifier, scope)
	}

	if !wkr.IsOfPrimitiveType(val, ast.Struct, ast.Entity, ast.AnonStruct, ast.Enum) {
		w.Error(node.Identifier.GetToken(), fmt.Sprintf("variable '%s' is not a struct, entity, enum or anonymous struct", node.Identifier.GetToken().Lexeme))
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
	var val wkr.Value
	if w.Context.Value.GetType().GetType() != wkr.NA {
		val = PropertyExpr(w, node.Identifier, scope)
	}else {
		val = GetNodeValue(w, &node.Identifier, scope)
	}

	if val.GetType().GetType() != wkr.Wrapper {
		token := w.Context.Node.GetToken()
		w.Error(token, fmt.Sprintf("%s is not a map nor a list", token.Lexeme))
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
	_type := TypeExpr(w, new.Type, w.Environment)

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

	w.ValidateArguments(args, val.Params, new.Token)

	return val
}

func SpawnExpr(w *wkr.Walker, new *ast.SpawnExpr, scope *wkr.Scope) wkr.Value {
	w.Context.Node = new
	_type := TypeExpr(w, new.Type, w.Environment)

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

	w.ValidateArguments(args, val.SpawnParams, new.Token)

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
		// if w.Nodes[0].(*ast.EnvironmentStmt).EnvType.Type == ast.Level {
		// 	w.Error(expr.GetToken(), "cannot use the Math library in a Level environment")
		// 	return &wkr.Invalid{}
		// }
		return GetNodeValueFromExternalEnv(w, expr.Node, scope, wkr.MathEnv)
	}

	return &wkr.Invalid{}
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

func TypeExpr(w *wkr.Walker, typee *ast.TypeExpr, env *wkr.Environment) wkr.Type {
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
		typ = TypeExpr(w, &ast.TypeExpr{Name: expr.Accessed}, walker.Environment)
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
		fields := map[string]*wkr.VariableVal{}

		for _, v := range typee.Fields {
			fields[v.Name.Lexeme] = &wkr.VariableVal{
				Name:  v.Name.Lexeme,
				Value: w.TypeToValue(TypeExpr(w, v.Type, env)),
				Token: v.Name,
			}
		}

		typ = wkr.NewAnonStructType(fields, false)
	case ast.Func:
		params := wkr.Types{}

		for _, v := range typee.Params {
			params = append(params, TypeExpr(w, v, env))
		}

		returns := wkr.Types{}
		for _, v := range typee.Returns {
			returns = append(returns, TypeExpr(w, v, env))
		}

		typ = &wkr.FunctionType{
			Params:  params,
			Returns: returns,
		}
	case ast.Map, ast.List:
		wrapped := TypeExpr(w, typee.WrappedType, env)
		typ = wkr.NewWrapperType(wkr.NewBasicType(pvt), wrapped)
	case ast.Entity:
		typ = &wkr.RawEntityType{}
	default:
		typeName := typee.Name.GetToken().Lexeme
		if entityVal, found := w.Environment.Entities[typeName]; found {
			typ = entityVal.GetType()
			break
		}
		if structVal, found := env.Structs[typeName]; found {
			typ = structVal.GetType()
			break
		}
		if customType, found := w.Environment.CustomTypes[typeName]; found {
			typ = customType
			break
		}
		if val := w.GetVariable(&env.Scope, typeName); val != nil {
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