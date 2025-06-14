package walker

import (
	"fmt"
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/tokens"
	"math"
	"strconv"
)

func (w *Walker) checkAccessibility(s *Scope, isPub bool, token tokens.Token) {
	if s.Environment.Name != w.environment.Name && !isPub {
		w.AlertSingle(&alerts.ForeignLocalVariableAccess{}, token, token.Lexeme)
	}
}

// ONLY CALL THIS IF YOU ALREADY CALLED ResolveVariable
//
// Returns the variable of name token.Lexeme
func (w *Walker) getVariable(s *Scope, token tokens.Token) *VariableVal {
	variable, ok := s.Variables[token.Lexeme]

	if !ok {
		return nil
	}
	w.checkAccessibility(s, variable.IsPub, token)

	return variable
}

func (w *Walker) typeExists(name string) bool {
	if _, found := w.environment.Entities[name]; found {
		return true
	}
	if _, found := w.environment.Classes[name]; found {
		return true
	}
	if _, found := w.environment.Scope.AliasTypes[name]; found {
		return true
	}
	if _, found := w.environment.Enums[name]; found {
		return true
	}
	if w.getTypeFromString(name) != ast.Invalid {
		return true
	}

	return false
}

func (w *Walker) declareVariable(s *Scope, value *VariableVal) (*VariableVal, bool) {
	if value.Name == "_" {
		return value, false
	}
	if w.typeExists(value.Name) {
		w.AlertSingle(&alerts.ConflictingVariableNameWithType{}, value.Token, value.Name)
		return value, false
	}
	if varFound, found := s.Variables[value.Name]; found {
		return varFound, false
	}

	s.Variables[value.Name] = value
	return value, true
}

func (w *Walker) declareClass(classVal *ClassVal) bool {
	if _, found := w.environment.Classes[classVal.Type.Name]; found {
		return false
	}

	w.environment.Classes[classVal.Type.Name] = classVal
	return true
}

func (w *Walker) declareEntity(entityVal *EntityVal) bool {
	if _, found := w.environment.Entities[entityVal.Type.Name]; found {
		return false
	}

	w.environment.Entities[entityVal.Type.Name] = entityVal
	return true
}

func (w *Walker) resolveVariable(s *Scope, token tokens.Token) *Scope {
	name := token.Lexeme
	if _, found := s.Variables[name]; found {
		return s
	}

	if s.Parent == nil {
		_, ok := BuiltinEnv.Scope.Variables[name]
		if ok {
			return &BuiltinEnv.Scope
		}
		for i := range s.Environment.importedWalkers {
			_, ok := s.Environment.importedWalkers[i].environment.Scope.Variables[name]
			if ok {
				return &s.Environment.importedWalkers[i].environment.Scope
			}
		}
		for _, v := range s.Environment.UsedLibraries {
			_, ok := BuiltinLibraries[v].Scope.Variables[name]
			if ok {
				return &BuiltinLibraries[v].Scope
			}
		}
		return nil
	}

	return w.resolveVariable(s.Parent, token)
}

func resolveTagScope[T ScopeTag](sc *Scope) (*Scope, *T) {
	if tag, ok := sc.Tag.(T); ok {
		return sc, &tag
	}

	if sc.Parent == nil {
		return nil, nil
	}

	return resolveTagScope[T](sc.Parent)
}

func (sc *Scope) resolveReturnable() *ExitableTag {
	if sc.Parent == nil {
		return nil
	}

	if returnable, ok := sc.Tag.(ExitableTag); ok {
		return &returnable
	} else if sc.Tag == nil {
		return nil
	}

	return sc.Parent.resolveReturnable()
}

func convertCallToMethodCall(call *ast.CallExpr, mi ast.MethodInfo) *ast.MethodCallExpr {
	copy := *call
	return &ast.MethodCallExpr{
		MethodInfo:  mi,
		Caller:      copy.Caller,
		GenericArgs: copy.GenericArgs,
		Args:        copy.Args,
	}
}

func (w *Walker) validateArguments(generics map[string]Type, args []Value, fn *FunctionVal, call ast.NodeCall) {
	nodeArgs := call.GetArgs()
	paramCount := len(fn.Params)

	defer func() {
		for k := range generics {
			generic := generics[k]
			if generic == UnknownTyp {
				w.AlertSingle(&alerts.MissingGenericArgument{}, call.GetToken(), k)
			}
		}
	}()

	if len(args) > paramCount && fn.Params[paramCount-1].GetType() != Variadic {
		extraAmount := len(args) - paramCount
		w.AlertMulti(&alerts.TooManyElementsGiven{},
			nodeArgs[len(nodeArgs)-extraAmount].GetToken(),
			nodeArgs[len(nodeArgs)-1].GetToken(),
			extraAmount,
			"value",
			"in call arguments",
		)
		return
	}

	var param Type
	for i, arg := range args {
		if i >= paramCount-1 {
			if fn.Params[paramCount-1].GetType() == Variadic {
				param = fn.Params[paramCount-1].(*VariadicType).Type
			} else {
				param = fn.Params[i]
			}
		} else {
			param = fn.Params[i]
		}

		if _, ok := arg.(Values); ok {
			w.AlertSingle(&alerts.InvalidCallAsArgument{}, nodeArgs[i].GetToken())
			continue
		}

		argType := arg.GetType()

		if typFound, found := resolveGenericType(param); found {
			genericArg, ok := generics[typFound.Name]
			if !ok || genericArg == UnknownTyp {
				genericArg = resolveGenericArgType(param, argType)
				generics[typFound.Name] = genericArg
				param = argType
			} else if !TypeEquals(genericArg, argType) {
				w.AlertSingle(&alerts.TypesMismatch{}, nodeArgs[i].GetToken(),
					"generic argument", genericArg,
					"function argument", argType,
				)
			}
			continue
		}

		if param == InvalidType || argType == InvalidType {
			continue
		}

		if !TypeEquals(param, argType) {
			w.AlertSingle(&alerts.InvalidArgumentType{}, nodeArgs[i].GetToken(), argType.String(), param.String())
			return
		}
	}

	for i, retArg := range fn.Returns {
		if retArg.GetType() == Generic {
			generic := retArg.(*GenericType)
			if generics[generic.Name] == UnknownTyp {
				w.AlertSingle(&alerts.MissingGenericArgument{}, call.GetToken(), generic.Name)
				continue
			}
			fn.Returns[i] = generics[generic.Name]
		}
	}

	if paramCount > len(args) {
		if len(nodeArgs) == 0 {
			w.AlertSingle(&alerts.TooFewElementsGiven{}, call.GetToken(), paramCount-len(args), "value", "in call arguments")
		} else {
			w.AlertSingle(&alerts.TooFewElementsGiven{}, nodeArgs[len(nodeArgs)-1].GetToken(), paramCount-len(args), "value", "in call arguments")
		}
	}

}

func resolveGenericType(typ Type) (*GenericType, bool) {
	if typ.GetType() == Generic {
		return typ.(*GenericType), true
	}

	if typ.GetType() == Wrapper {
		return resolveGenericType(typ.(*WrapperType).WrappedType)
	}

	return nil, false
}

func resolveGenericArgType(genericParam Type, genericArg Type) Type {
	if genericParam.GetType() == Wrapper && genericArg.GetType() == Wrapper {
		wrapper1 := genericParam.(*WrapperType)
		wrapper2 := genericArg.(*WrapperType)

		if wrapper1.Type.GetType() != Generic {
			return resolveGenericArgType(wrapper1.WrappedType, wrapper2.WrappedType)
		}

		return wrapper2.Type
	}

	return genericArg
}

var ops = map[string]func(float64, float64) float64{
	"+":  func(a, b float64) float64 { return a + b },
	"-":  func(a, b float64) float64 { return a - b },
	"*":  func(a, b float64) float64 { return a * b },
	"/":  func(a, b float64) float64 { return a / b },
	"^":  func(a, b float64) float64 { return math.Pow(a, a) },
	"%":  func(a, b float64) float64 { return math.Mod(a, b) },
	"\\": func(a, b float64) float64 { return a / b },
}

// Validates the arithmetic operands, so for example a condition "value1 + value2", "value1 - value2"
//
// If both values are booleans, and the boolean values are known at compile time, the condition will be calculated and the returning Value will have the calculation in the BoolVal.
func (w *Walker) validateArithmeticOperands(leftVal Value, rightVal Value, node *ast.BinaryExpr, context string) Value {
	left, right := leftVal.GetType(), rightVal.GetType()
	if left == InvalidType || right == InvalidType {
		return &Invalid{}
	}
	if !isNumerical(left.PVT()) {
		w.AlertSingle(&alerts.TypeMismatch{}, node.Left.GetToken(), "a numerical type", left, context)
		return &Invalid{}
	}
	if !isNumerical(right.PVT()) {
		w.AlertSingle(&alerts.TypeMismatch{}, node.Right.GetToken(), "a numerical type", right, context)
		return &Invalid{}
	}
	if !TypeEquals(left, right) {
		w.AlertMulti(&alerts.TypesMismatch{}, node.Left.GetToken(), node.Right.GetToken(), "left value", left, "right value", right)
		return &Invalid{}
	}
	if left.PVT() != ast.Number {
		return w.typeToValue(left)
	}
	num1, num2 := leftVal.(*NumberVal), rightVal.(*NumberVal)
	if num1.Value == "" || num2.Value == "" {
		return &NumberVal{}
	}

	n1, err := strconv.ParseFloat(num1.Value, 64)
	if err != nil {
		return &NumberVal{}
	}
	n2, err2 := strconv.ParseFloat(num2.Value, 64)
	if err2 != nil {
		return &NumberVal{}
	}

	n3 := ops[node.Operator.Lexeme](n1, n2)

	return NewNumberVal(fmt.Sprintf("%v", n3))
}

// Validates the conditional operands, so for example a condition "value1 and value2", "value1 or value2"
//
// If both values are booleans, and the boolean values are known at compile time, the condition will be calculated and the returning Value will have the calculation in the BoolVal.
func (w *Walker) validateConditionalOperands(leftVal Value, rightVal Value, node *ast.BinaryExpr) Value {
	left, right := leftVal.GetType(), rightVal.GetType()

	if left == InvalidType || right == InvalidType {
		return &Invalid{}
	}
	if left.PVT() != ast.Bool {
		w.AlertSingle(&alerts.TypeMismatch{}, node.Left.GetToken(), "a boolean", left, "in logical comparison expression")
		return &BoolVal{}
	}
	if right.PVT() != ast.Bool {
		w.AlertSingle(&alerts.TypeMismatch{}, node.Right.GetToken(), "a boolean", right, "in logical comparison expression")
		return &BoolVal{}
	}
	leftBool, rightBool := leftVal.(*BoolVal), rightVal.(*BoolVal)

	leftCondition, ok := strconv.ParseBool(leftBool.Value)
	rightCondition, ok2 := strconv.ParseBool(rightBool.Value)
	if node.Operator.Type == tokens.And {
		if (ok == nil && !leftCondition) || (ok2 == nil && !rightCondition) {
			return NewBoolVal("false")
		} else if ok == nil && ok2 == nil {
			return NewBoolVal(strconv.FormatBool(leftCondition && rightCondition))
		}
	} else {
		if (ok == nil && leftCondition) || (ok2 == nil && !rightCondition) {
			return NewBoolVal("true")
		} else if ok == nil && ok2 == nil {
			return NewBoolVal(strconv.FormatBool(leftCondition || rightCondition))
		}
	}
	return &BoolVal{}
}

func (w *Walker) validateReturnValues(returnArgs []ast.Node, _return []Value2, expectReturn []Type, context string) {
	retLen := len(_return)
	expRetLen := len(expectReturn)
	if retLen < expRetLen {
		requiredAmount := expRetLen - retLen
		w.AlertSingle(&alerts.TooFewElementsGiven{}, returnArgs[len(returnArgs)-1].GetToken(), requiredAmount, "return value", context)
	} else if retLen > expRetLen {
		extraAmount := retLen - expRetLen
		w.AlertMulti(&alerts.TooFewElementsGiven{},
			returnArgs[len(returnArgs)-extraAmount].GetToken(),
			returnArgs[len(returnArgs)-1].GetToken(),
			extraAmount,
			"return value",
			context,
		)
	}
	for i := range _return {
		if i >= expRetLen {
			break
		}
		if _return[i].GetType() == InvalidType || expectReturn[i] == InvalidType {
			continue
		}
		if !TypeEquals(_return[i].GetType(), expectReturn[i]) {
			w.AlertSingle(&alerts.TypeMismatch{}, returnArgs[_return[i].Index].GetToken(),
				expectReturn[i],
				_return[i].GetType(),
				fmt.Sprintf(context+" (arg %d)", i+1),
			)
		}
	}
}

func (w *Walker) ifCondition(node *ast.Node, scope *Scope) {
	condition := w.GetActualNodeValue(node, scope)
	if condition.GetType() == InvalidType {
		return
	}
	if condition.GetType().PVT() != ast.Bool {
		w.AlertSingle(&alerts.InvalidCondition{}, (*node).GetToken(), "in if statement")
	} else if conditionValue := condition.(*BoolVal).Value; conditionValue != "" {
		w.AlertSingle(&alerts.LiteralCondition{}, (*node).GetToken(), conditionValue)
	}
}

func (w *Walker) getParameters(parameters []ast.FunctionParam, scope *Scope) []Type {
	variadicParams := make(map[tokens.Token]int)
	params := make([]Type, 0)
	for i, param := range parameters {
		params = append(params, w.typeExpression(param.Type, scope))
		if params[i].GetType() == Variadic {
			variadicParams[parameters[i].Name] = i
		}
		value := w.typeToValue(params[i])
		variable := NewVariable(param.Name, value)
		w.declareVariable(scope, variable)
	}

	if len(variadicParams) > 1 {
		for k := range variadicParams {
			w.AlertSingle(&alerts.MoreThanOneVariadicParameter{}, k)
			break
		}
	} else if len(variadicParams) == 1 {
		for k, v := range variadicParams {
			if v != len(parameters)-1 {
				w.AlertSingle(&alerts.VariadicParameterNotAtEnd{}, k)
			}
		}
	}

	return params
}

func (w *Walker) getReturns(returns []*ast.TypeExpr, scope *Scope) []Type {
	returnType := EmptyReturn
	for i := range returns {
		returnType = append(returnType, w.typeExpression(returns[i], scope))
	}

	return returnType
}

func (w *Walker) resolveGenericParam(name string, scope *Scope) (*GenericType, bool) {
	if scope.Parent == nil {
		return nil, false
	}

	if fn, ok := scope.Tag.(*FuncTag); ok {
		for _, v := range fn.Generics {
			if name == v.Name {
				return v, true
			}
		}
	} else if ct, ok := scope.Tag.(*ClassTag); ok {
		for _, v := range ct.Val.Type.Generics {
			if name == v.GenericName {
				return NewGeneric(name), true
			}
		}
	} else if et, ok := scope.Tag.(*EntityTag); ok {
		for _, v := range et.EntityVal.Type.Generics {
			if name == v.GenericName {
				return NewGeneric(name), true
			}
		}
	}

	return w.resolveGenericParam(name, scope.Parent)
}

func (w *Walker) getGenericParams(genericParams []*ast.IdentifierExpr, scope *Scope) []*GenericType {
	generics := make([]*GenericType, 0)

	for _, generic := range genericParams {
		if _, found := w.resolveGenericParam(generic.Name.Lexeme, scope); found {
			w.AlertSingle(&alerts.DuplicateElement{}, generic.GetToken(), "generic parameter", generic.Name.Lexeme)
			break
		}
		for i := range generics {
			if generics[i].Name == generic.Name.Lexeme {
				w.AlertSingle(&alerts.DuplicateElement{}, generic.GetToken(), "generic parameter", generic.Name.Lexeme)
				break
			}

		}
		generics = append(generics, NewGeneric(generic.Name.Lexeme))
	}

	return generics
}

func (w *Walker) getGenerics(genericArgs []*ast.TypeExpr, expectedGenerics []*GenericType, scope *Scope) map[string]Type {
	receivedGenericsLength := len(genericArgs)
	expectedGenericsLength := len(expectedGenerics)

	suppliedGenerics := map[string]Type{}
	if receivedGenericsLength > expectedGenericsLength {
		extraAmount := receivedGenericsLength - expectedGenericsLength
		w.AlertMulti(&alerts.TooManyElementsGiven{},
			genericArgs[receivedGenericsLength-extraAmount].GetToken(),
			genericArgs[receivedGenericsLength-1].GetToken(),
			extraAmount,
			"value",
			"in generic arguments",
		)
	} else {
		for i := range expectedGenerics {
			if i > len(genericArgs)-1 {
				suppliedGenerics[expectedGenerics[i].Name] = UnknownTyp
				continue
			}
			suppliedGenerics[expectedGenerics[i].Name] = w.typeExpression(genericArgs[i], scope)
		}
	}

	return suppliedGenerics
}

func (w *Walker) typeToValue(_type Type) Value {
	if _type.GetType() == RawEntity {
		return &RawEntityVal{}
	}
	if _type.GetType() == Variadic {
		return &ListVal{
			ValueType: _type.(*VariadicType).Type,
		}
	}
	switch _type.PVT() {
	case ast.Fixed:
		return &FixedVal{}
	case ast.Bool:
		return &BoolVal{}
	case ast.Func:
		ft := _type.(*FunctionType)
		return &FunctionVal{
			Params:  ft.Params,
			Returns: ft.Returns,
		}
	case ast.Text:
		return &StringVal{}
	case ast.Number:
		return &NumberVal{}
	case ast.List:
		return &ListVal{
			ValueType: _type.(*WrapperType).WrappedType,
		}
	case ast.Map:
		return &MapVal{
			MemberType: _type.(*WrapperType).WrappedType,
		}
	case ast.Class:
		named := _type.(*NamedType)
		val := CopyClassVal(w.walkers[named.EnvName].environment.Classes[named.Name])
		for _, v := range named.Generics {
			generic := v.Type
			name := v.GenericName

			if generic == UnknownTyp {
				continue
			}
			for i, v := range val.Fields {
				if v.Var.Value.GetType().GetType() == Generic {
					val.Fields[i].Var.Value = w.typeToValue(generic)
				}
			}
			for i := range val.Methods {
				fn := val.Methods[i].Value.(*FunctionVal)
				for j, v3 := range fn.Params {
					if gen, ok := v3.(*GenericType); ok && gen.Name == name {
						fn.Params[j] = generic
					}
				}
				for j, v3 := range fn.Returns {
					if gen, ok := v3.(*GenericType); ok && gen.Name == name {
						fn.Returns[j] = generic
					}
				}
			}
			fn := val.New
			for j, v3 := range fn.Params {
				if gen, ok := v3.(*GenericType); ok && gen.Name == name {
					fn.Params[j] = generic
				}
			}
			for j, v3 := range fn.Returns {
				if gen, ok := v3.(*GenericType); ok && gen.Name == name {
					fn.Returns[j] = generic
				}
			}
		}
		val.Type.Generics = named.Generics
		return &val
	case ast.Struct:
		return &StructVal{
			Fields: _type.(*StructType).Fields,
		}
	case ast.Entity:
		named := _type.(*NamedType)
		val := CopyEntityVal(w.walkers[named.EnvName].environment.Entities[named.Name])
		for _, v := range named.Generics {
			generic := v.Type
			name := v.GenericName

			if generic == UnknownTyp {
				continue
			}
			for i, v := range val.Fields {
				if v.Var.Value.GetType().GetType() == Generic {
					val.Fields[i].Var.Value = w.typeToValue(generic)
				}
			}
			for i := range val.Methods {
				fn := val.Methods[i].Value.(*FunctionVal)
				for j, v3 := range fn.Params {
					if gen, ok := v3.(*GenericType); ok && gen.Name == name {
						fn.Params[j] = generic
					}
				}
				for j, v3 := range fn.Returns {
					if gen, ok := v3.(*GenericType); ok && gen.Name == name {
						fn.Returns[j] = generic
					}
				}
			}
			fn := val.Spawn
			for j, v3 := range fn.Params {
				if gen, ok := v3.(*GenericType); ok && gen.Name == name {
					fn.Params[j] = generic
				}
			}
			for j, v3 := range fn.Returns {
				if gen, ok := v3.(*GenericType); ok && gen.Name == name {
					fn.Returns[j] = generic
				}
			}
			fn = val.Destroy
			for j, v3 := range fn.Params {
				if gen, ok := v3.(*GenericType); ok && gen.Name == name {
					fn.Params[j] = generic
				}
			}
			for j, v3 := range fn.Returns {
				if gen, ok := v3.(*GenericType); ok && gen.Name == name {
					fn.Returns[j] = generic
				}
			}
		}
		val.Type.Generics = named.Generics
		return &val
	case ast.Object:
		return &Unknown{}
	case ast.Generic:
		return &GenericVal{
			Type: _type.(*GenericType),
		}
	case ast.Enum:
		enumType := _type.(*EnumType)
		switch enumType.EnvName {
		case "Pewpew":
			return PewpewAPI.Enums[enumType.Name]
		default:
			return w.walkers[enumType.EnvName].environment.Enums[enumType.Name]
		}
	case ast.Path:
		pathType := _type.(*PathType)
		return NewPathVal("", pathType.Env, "")
	default:
		return &Invalid{}
	}
}

func (w *Walker) getContentsValueType(elems []ast.Node, scope *Scope) Type {
	valTypes := []Type{}
	if len(elems) == 0 {
		return UnknownTyp
	}
	val := w.GetActualNodeValue(&elems[0], scope)
	valTypes = append(valTypes, val.GetType())
	for i := range elems {
		if i == 0 {
			continue
		}
		val = w.GetActualNodeValue(&elems[i], scope)
		valTypes = append(valTypes, val.GetType())
		if !TypeEquals(valTypes[i-1], valTypes[i]) {
			w.AlertSingle(&alerts.MixedMapOrListContents{}, elems[i].GetToken(),
				"list",
				valTypes[i].String(),
				valTypes[i-1].String(),
			)
			return InvalidType
		}
	}

	return valTypes[0]
}

func (w *Walker) getTypeFromString(str string) ast.PrimitiveValueType {
	switch str {
	case "number":
		return ast.Number
	case "fixed":
		return ast.Fixed
	case "text":
		return ast.Text
	case "map":
		return ast.Map
	case "list":
		return ast.List
	case "fn":
		return ast.Func
	case "bool":
		return ast.Bool
	case "struct":
		return ast.Struct
	default:
		return ast.Invalid
	}
}

func (w *Walker) typesToValues(types []Type) Values {
	vals := Values{}

	for _, v := range types {
		vals = append(vals, w.typeToValue(v))
	}

	return vals
}

func (w *Walker) reportExits(sender ExitableTag, scope *Scope) {
	receiver_ := scope.resolveReturnable()

	if receiver_ == nil {
		return
	}

	receiver := *receiver_

	receiver.SetExit(sender.GetIfExits(Yield), Yield)
	receiver.SetExit(sender.GetIfExits(Return), Return)
	receiver.SetExit(sender.GetIfExits(EntityDestruction), EntityDestruction)
	receiver.SetExit(sender.GetIfExits(ControlFlow), ControlFlow)
}

type ProcedureType int

const (
	Function ProcedureType = iota
	Method
)

func isNumerical(pvt ast.PrimitiveValueType) bool {
	return pvt == ast.Number || pvt == ast.Fixed
}

func SetupLibraryEnvironments() {
	PewpewAPI.Scope.Environment = PewpewAPI
	FmathAPI.Scope.Environment = FmathAPI
	MathAPI.Scope.Environment = MathAPI
	StringAPI.Scope.Environment = StringAPI
	TableAPI.Scope.Environment = TableAPI
	BuiltinEnv.Scope.Environment = BuiltinEnv
}

// used in return, yield, assignment statement and variable declaration
type Value2 struct {
	Value
	Index int
}

type Values2 []Value2

func (v2 Values2) Types() *[]Type {
	vals := make([]Type, 0)
	for _, v := range v2 {
		vals = append(vals, v.Value.GetType())
	}
	return &vals
}
