package walker

import (
	"fmt"
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/tokens"
	"strconv"
)

func (w *Walker) checkAccessibility(s *Scope, isLocal bool, token tokens.Token) {
	if s.Environment.Name != w.environment.Name && isLocal {
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
	w.checkAccessibility(s, variable.IsLocal, token)

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

	return false
}

func (w *Walker) getAliasType(s *Scope, token tokens.Token) *AliasType {
	alias, found := s.AliasTypes[token.Lexeme]
	if !found {
		// w.Error(w.Context.Node.GetToken(), fmt.Sprintf("no struct named %s exists", name))
		return nil
	}

	alias.IsUsed = true
	w.checkAccessibility(s, alias.IsLocal, token)

	return alias
}

func (w *Walker) getClass(s *Scope, token tokens.Token) *ClassVal {
	class, found := w.environment.Classes[token.Lexeme]
	if !found {
		// w.Error(w.Context.Node.GetToken(), fmt.Sprintf("no struct named %s exists", name))
		return nil
	}

	class.Type.IsUsed = true
	w.checkAccessibility(s, class.IsLocal, token)

	return class
}

func (w *Walker) getEntity(name string) *EntityVal {
	entityType, found := w.environment.Entities[name]
	if !found {
		// w.Error(w.Context.Node.GetToken(), fmt.Sprintf("no struct named %s exists", name))
		return nil
	}

	entityType.Type.IsUsed = true

	return entityType
}

func (s *Scope) assignVariable(variable *VariableVal, value Value) Value {
	variable.Value = value

	return variable
}

func (w *Walker) declareVariable(s *Scope, value *VariableVal) (*VariableVal, bool) {
	if varFound, found := s.Variables[value.Name]; found {
		return varFound, false
	}

	s.Variables[value.Name] = value
	return value, true
}

func (w *Walker) declareClass(structVal *ClassVal) bool {
	if _, found := w.environment.Classes[structVal.Type.Name]; found {
		return false
	}

	w.environment.Classes[structVal.Type.Name] = structVal
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
			_, ok := LibraryEnvs[v].Scope.Variables[name]
			if ok {
				return &LibraryEnvs[v].Scope
			}
		}
		return nil
	}

	return w.resolveVariable(s.Parent, token)
}

func resolveTagScope[T ScopeTag](sc *Scope) (*Scope, *ScopeTag, *T) {
	if tag, ok := sc.Tag.(T); ok {
		return sc, &sc.Tag, &tag
	}

	if sc.Parent == nil {
		return nil, nil, nil
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

func convertNodeToAccessFieldExpr(ident ast.Node, index int, exprType ast.SelfExprType, envName string, entityName string) *ast.AccessExpr {
	fieldExpr := &ast.FieldExpr{
		Index:      index,
		Field:      ident,
		ExprType:   exprType,
		EnvName:    envName,
		EntityName: entityName,
	}

	return &ast.AccessExpr{
		Start: &ast.SelfExpr{
			Token: ident.GetToken(),
			Type:  exprType,
		},
		Accessed: []ast.Node{
			fieldExpr,
		},
	}
}

func convertCallToMethodCall(call *ast.CallExpr, exprType ast.SelfExprType, envName string, name string) *ast.MethodCallExpr {
	copy := *call
	return &ast.MethodCallExpr{
		EnvName:     envName,
		TypeName:    name,
		ExprType:    exprType,
		Caller:      copy.Caller,
		GenericArgs: copy.GenericArgs,
		Args:        copy.Args,
		MethodName:  call.Caller.GetToken().Lexeme,
	}
}

func (w *Walker) validateArguments(generics map[string]Type, args []Type, params []Type, call ast.NodeCall) (int, bool) {

	paramCount := len(params)
	if paramCount > len(args) {
		// w.Error(callToken, "too few arguments given in call")
		return -1, true
	}
	var param Type
	for i, typeVal := range args {
		if i >= paramCount-1 {
			if params[paramCount-1].GetType() == Variadic {
				param = params[paramCount-1].(*VariadicType).Type
			} else if i > paramCount-1 {
				// w.Error(callToken, "too many arguments given in call")
				return -1, true
			} else {
				param = params[i]
			}
		} else {
			param = params[i]
		}

		if typFound, found := resolveGenericType(&param); found {
			generic := (*typFound).(*GenericType)
			if typ, found := generics[generic.Name]; found {
				*typFound = typ
			} else {
				generics[generic.Name] = resolveMatchingType(param, typeVal)
				param = typeVal
			}
		}

		if !TypeEquals(param, typeVal) {
			w.AlertSingle(&alerts.InvalidArgumentType{}, call.GetArgs()[i].GetToken(), typeVal.String(), param.String())
			return i, false
		}
	}
	return -1, true
}

func resolveGenericType(typ *Type) (*Type, bool) {
	if (*typ).GetType() == Generic {
		return typ, true
	}

	if (*typ).GetType() == Wrapper {
		return resolveGenericType(&(*typ).(*WrapperType).WrappedType)
	}

	return nil, false
}

func resolveMatchingType(predefinedType Type, receivedType Type) Type {
	if predefinedType.GetType() == Wrapper && receivedType.GetType() == Wrapper {
		wrapper1 := predefinedType.(*WrapperType)
		wrapper2 := receivedType.(*WrapperType)

		if TypeEquals(wrapper1.Type, wrapper2.Type) {
			return resolveMatchingType(wrapper1.WrappedType, wrapper2.WrappedType)
		}

		return wrapper2.Type
	}

	return receivedType
}

func (w *Walker) validateArithmeticOperands(left Type, right Type, node *ast.BinaryExpr) Type {
	if !isNumerical(left.PVT()) {
		w.AlertSingle(&alerts.TypeMismatch{}, node.Left.GetToken(), "a numerical type", left, "in arithmetic expression")
		return InvalidType
	}
	if !isNumerical(right.PVT()) {
		w.AlertSingle(&alerts.TypeMismatch{}, node.Right.GetToken(), "a numerical type", right, "in arithmetic expression")
		return InvalidType
	}
	if !TypeEquals(left, right) {
		w.AlertSingle(&alerts.ArithmeticTypesMismatch{}, node.Left.GetToken(), left, right)
		return InvalidType
	}

	return left
}

func (w *Walker) validateConditionalOperands(leftVal Value, rightVal Value, node *ast.BinaryExpr) Value {
	left, right := leftVal.GetType(), rightVal.GetType()
	if left.PVT() != ast.Bool {
		w.AlertSingle(&alerts.TypeMismatch{}, node.Left.GetToken(), "a numerical type", left, "in arithmetic expression")
		return NewBoolVal()
	}
	if right.PVT() != ast.Bool {
		w.AlertSingle(&alerts.TypeMismatch{}, node.Right.GetToken(), "a numerical type", right, "in arithmetic expression")
		return NewBoolVal()
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
	return NewBoolVal()
}

func (w *Walker) validateReturnValues(returnArgs []ast.Node, _return []Value2, expectReturn []Type, context string) {
	retLen := len(_return)
	expRetLen := len(expectReturn)
	if retLen < expRetLen {
		requiredAmount := expRetLen - retLen
		w.AlertSingle(&alerts.TooFewValuesGiven{}, returnArgs[len(returnArgs)-1].GetToken(), requiredAmount, context)
	} else if retLen > expRetLen {
		extraAmount := retLen - expRetLen
		w.AlertMulti(&alerts.TooFewValuesGiven{},
			returnArgs[len(returnArgs)-extraAmount].GetToken(),
			returnArgs[len(returnArgs)-1].GetToken(),
			extraAmount,
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
				_return[i].GetType(),
				expectReturn[i],
				fmt.Sprintf(context+" (arg %d)", i+1),
			)
		}
	}
}

func (w *Walker) getReturns(returns []*ast.TypeExpr, scope *Scope) []Type {
	returnType := EmptyReturn
	for i := range returns {
		returnType = append(returnType, w.typeExpression(returns[i], scope))
	}

	return returnType
}

func (w *Walker) getGenericParams(genericParams []*ast.IdentifierExpr) []*GenericType {
	generics := make([]*GenericType, 0)

	for _, generic := range genericParams {
		for i := range generics {
			if generics[i].Name == generic.Name.Lexeme {
				w.AlertSingle(&alerts.DuplicateGenericParameter{}, generic.GetToken(), generic.Name.Lexeme)
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
		w.AlertMulti(&alerts.TooManyValuesGiven{},
			genericArgs[expectedGenericsLength-extraAmount].GetToken(),
			genericArgs[expectedGenericsLength-1].GetToken(),
			extraAmount,
			"in generic arguments",
		)
	} else {
		for i := range genericArgs {
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
	case ast.Radian, ast.Fixed, ast.FixedPoint, ast.Degree:
		return &FixedVal{SpecificType: _type.PVT()}
	case ast.Bool:
		return NewBoolVal()
	case ast.Func:
		ft := _type.(*FunctionType)
		return &FunctionVal{
			Params:  ft.Params,
			Returns: ft.Returns,
		}
	case ast.String:
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
		val := w.walkers[_type.(*NamedType).EnvName].environment.Classes[_type.String()]
		return val
	case ast.Struct:
		return &AnonStructVal{
			Fields: _type.(*StructType).Fields,
		}
	case ast.Enum:
		enum := _type.(*EnumType)
		walker, found := w.walkers[enum.EnvName]
		var variable *VariableVal
		switch enum.EnvName {
		case "Pewpew":
			variable = PewpewEnv.Scope.Variables[enum.Name]
		default:
			variable = walker.environment.Scope.Variables[enum.Name]
		}
		if variable == nil {
			panic(fmt.Sprintf("Enum variable could not be found when converting enum type to value (envName:%v, name:%v, walkerFound:%v)", enum.EnvName, enum.Name, found))
		}
		return variable
	case ast.Entity:
		val := w.walkers[_type.(*NamedType).EnvName].environment.Entities[_type.String()]
		return val
	case ast.Object:
		return &Unknown{}
	case ast.Generic:
		return &GenericVal{
			Type: _type.(*GenericType),
		}
	default:
		return &Invalid{}
	}
}

func (w *Walker) getContentsValueType(elems []ast.Node, scope *Scope) Type {
	valTypes := []Type{}
	if len(elems) == 0 {
		return ObjectTyp
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
				valTypes[i-1].String(),
				valTypes[i].String(),
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
		return ast.FixedPoint
	case "text":
		return ast.String
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

func (w *Walker) determineCallTypeString(callType ProcedureType) string {
	if callType == Function {
		return "function"
	}

	return "method"
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
	receiver.SetExit(sender.GetIfExits(Break), Break)
	receiver.SetExit(sender.GetIfExits(Continue), Continue)
	receiver.SetExit(sender.GetIfExits(All), All)
}

type ProcedureType int

const (
	Function ProcedureType = iota
	Method
)

func isNumerical(pvt ast.PrimitiveValueType) bool {
	return isOfPrimitiveType(pvt, ast.Number, ast.Fixed, ast.FixedPoint, ast.Degree, ast.Radian)
}

func isOfPrimitiveType(pvt ast.PrimitiveValueType, types ...ast.PrimitiveValueType) bool {
	if types == nil {
		return false
	}
	for _, prim := range types {
		if pvt == prim {
			return true
		}
	}

	return false
}

func determineCallTypeString(callType ProcedureType) string {
	if callType == Function {
		return "function"
	}

	return "method"
}

func SetupLibraryEnvironments() {
	PewpewEnv.Scope.Environment = PewpewEnv
	FmathEnv.Scope.Environment = FmathEnv
	MathEnv.Scope.Environment = MathEnv
	StringEnv.Scope.Environment = StringEnv
	TableEnv.Scope.Environment = TableEnv
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
