package walker

import (
	"fmt"
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/helpers"
	"hybroid/tokens"
)

func (w *Walker) CheckAccessibility(s *Scope, isLocal bool, token tokens.Token) {
	if s.Environment.Name != w.environment.Name && isLocal {
		w.AlertSingle(&alerts.ForeignLocalVariableAccess{}, token, token.Lexeme)
	}
}

// ONLY CALL THIS IF YOU ALREADY CALLED ResolveVariable
//
// Returns the variable of name token.Lexeme
func (w *Walker) GetVariable(s *Scope, token tokens.Token) *VariableVal {
	variable, ok := s.Variables[token.Lexeme]

	if !ok {
		return nil
	}
	w.CheckAccessibility(s, variable.IsLocal, token)

	return variable
}

func (w *Walker) TypeExists(name string) bool {
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

func (w *Walker) GetAliasType(s *Scope, token tokens.Token) *AliasType {
	alias, found := s.AliasTypes[token.Lexeme]
	if !found {
		// w.Error(w.Context.Node.GetToken(), fmt.Sprintf("no struct named %s exists", name))
		return nil
	}

	alias.IsUsed = true
	w.CheckAccessibility(s, alias.IsLocal, token)

	return alias
}

func (w *Walker) GetClass(s *Scope, token tokens.Token) *ClassVal {
	class, found := w.environment.Classes[token.Lexeme]
	if !found {
		// w.Error(w.Context.Node.GetToken(), fmt.Sprintf("no struct named %s exists", name))
		return nil
	}

	class.Type.IsUsed = true
	w.CheckAccessibility(s, class.IsLocal, token)

	return class
}

func (w *Walker) GetEntity(name string) *EntityVal {
	entityType, found := w.environment.Entities[name]
	if !found {
		// w.Error(w.Context.Node.GetToken(), fmt.Sprintf("no struct named %s exists", name))
		return nil
	}

	entityType.Type.IsUsed = true

	return entityType
}

func (s *Scope) AssignVariable(variable *VariableVal, value Value) Value {
	variable.Value = value

	return variable
}

func (w *Walker) DeclareVariable(s *Scope, value *VariableVal) (*VariableVal, bool) {
	if varFound, found := s.Variables[value.Name]; found {
		return varFound, false
	}

	s.Variables[value.Name] = value
	return value, true
}

func (w *Walker) DeclareClass(structVal *ClassVal) bool {
	if _, found := w.environment.Classes[structVal.Type.Name]; found {
		return false
	}

	w.environment.Classes[structVal.Type.Name] = structVal
	return true
}

func (w *Walker) DeclareEntity(entityVal *EntityVal) bool {
	if _, found := w.environment.Entities[entityVal.Type.Name]; found {
		return false
	}

	w.environment.Entities[entityVal.Type.Name] = entityVal
	return true
}

func (w *Walker) ResolveVariable(s *Scope, token tokens.Token) *Scope {
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
		for k, v := range s.Environment.UsedLibraries {
			if !v {
				continue
			}

			_, ok := LibraryEnvs[k].Scope.Variables[name]
			if ok {
				return &LibraryEnvs[k].Scope
			}
		}
		return nil
	}

	return w.ResolveVariable(s.Parent, token)
}

func ResolveTagScope[T ScopeTag](sc *Scope) (*Scope, *ScopeTag, *T) {
	if tag, ok := sc.Tag.(T); ok {
		return sc, &sc.Tag, &tag
	}

	if sc.Parent == nil {
		return nil, nil, nil
	}

	return ResolveTagScope[T](sc.Parent)
}

func (sc *Scope) ResolveReturnable() *ExitableTag {
	if sc.Parent == nil {
		return nil
	}

	if returnable := helpers.GetValOfInterface[ExitableTag](sc.Tag); returnable != nil {
		return returnable
	}

	if helpers.IsZero(sc.Tag) {
		return nil
	}

	return sc.Parent.ResolveReturnable()
}

func (w *Walker) ValidateArguments(generics map[string]Type, args []Type, params []Type, call ast.NodeCall) (int, bool) {

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

		if typFound, found := ResolveGenericType(&param); found {
			generic := (*typFound).(*GenericType)
			if typ, found := generics[generic.Name]; found {
				*typFound = typ
			} else {
				generics[generic.Name] = ResolveMatchingType(param, typeVal)
				param = typeVal
			}
		}

		if !TypeEquals(param, typeVal) {
			w.AlertSingle(&alerts.InvalidArgumentType{}, call.GetArgs()[i].GetToken(), typeVal.ToString(), param.ToString())
			return i, false
		}
	}
	return -1, true
}

func ResolveGenericType(typ *Type) (*Type, bool) {
	if (*typ).GetType() == Generic {
		return typ, true
	}

	if (*typ).GetType() == Wrapper {
		return ResolveGenericType(&(*typ).(*WrapperType).WrappedType)
	}

	return nil, false
}

func ResolveMatchingType(predefinedType Type, receivedType Type) Type {
	if predefinedType.GetType() == Wrapper && receivedType.GetType() == Wrapper {
		wrapper1 := predefinedType.(*WrapperType)
		wrapper2 := receivedType.(*WrapperType)

		if TypeEquals(wrapper1.Type, wrapper2.Type) {
			return ResolveMatchingType(wrapper1.WrappedType, wrapper2.WrappedType)
		}

		return wrapper2.Type
	}

	return receivedType
}

func (w *Walker) DetermineValueType(left Type, right Type) Type {
	if TypeEquals(left, right) {
		if left.GetType() == Fixed {
			return left
		}
		return right
	}
	if left.GetType() == Fixed {
		if right.GetType() == Fixed || right.PVT() == ast.Number {
			return left
		}
	}
	if right.GetType() == Fixed {
		if left.GetType() == Fixed || left.PVT() == ast.Number {
			return right
		}
	}

	return InvalidType
}

func (w *Walker) ValidateArithmeticOperands(left Type, right Type, expr *ast.BinaryExpr) bool {
	if left.PVT() == ast.Invalid {
		// w.Error(expr.Left.GetToken(), "cannot perform arithmetic on Invalid value")
		return false
	}

	if right.PVT() == ast.Invalid {
		// w.Error(expr.Right.GetToken(), "cannot perform arithmetic on Invalid value")
		return false
	}

	switch left.PVT() {
	case ast.List, ast.Map, ast.String, ast.Bool, ast.Entity, ast.Struct:
		// w.Error(expr.Left.GetToken(), "cannot perform arithmetic on a non-number value")
		return false
	}

	switch right.PVT() {
	case ast.List, ast.Map, ast.String, ast.Bool, ast.Entity, ast.Struct:
		// w.Error(expr.Right.GetToken(), "cannot perform arithmetic on a non-number value")
		return false
	}

	return true
}

func (w *Walker) ValidateReturnValues(returnArgs []ast.Node, _return []Type, expectReturn []Type, context string) {
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
		if _return[i] == InvalidType {
			continue
		}
		if expectReturn[i] == InvalidType {
			continue
		}
		if !TypeEquals(_return[i], expectReturn[i]) {
			w.AlertSingle(&alerts.TypeMismatch{}, returnArgs[i].GetToken(),
				_return[i].GetType(),
				expectReturn[i].GetType(),
				context,
			)
		}
	}
}

func (w *Walker) GetReturns(returns []*ast.TypeExpr, scope *Scope) []Type {
	returnType := EmptyReturn
	for i := range returns {
		returnType = append(returnType, w.TypeExpr(returns[i], scope))
	}

	return returnType
}

func (w *Walker) GetGenerics(genericArgs []*ast.TypeExpr, expectedGenerics []*GenericType, scope *Scope) map[string]Type {
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
			suppliedGenerics[expectedGenerics[i].Name] = w.TypeExpr(genericArgs[i], scope)
		}
	}

	return suppliedGenerics
}

func (w *Walker) TypeToValue(_type Type) Value {
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
		return &BoolVal{} // there is no func here
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
	case ast.Struct:
		val, _ := w.walkers[_type.(*NamedType).EnvName].environment.Classes[_type.ToString()]
		return val
	case ast.AnonStruct:
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
			variable, _ = walker.environment.Scope.Variables[enum.Name]
		}
		if variable == nil {
			panic(fmt.Sprintf("Enum variable could not be found when converting enum type to value (envName:%v, name:%v, walkerFound:%v)", enum.EnvName, enum.Name, found))
		}
		return variable
	case ast.Entity:
		val, _ := w.walkers[_type.(*NamedType).EnvName].environment.Entities[_type.ToString()]
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

func (w *Walker) GetContentsValueType(elems []ast.Node, scope *Scope) Type {
	valTypes := []Type{}
	if len(elems) == 0 {
		return ObjectTyp
	}
	val := w.GetNodeActualValue(&elems[0], scope)
	valTypes = append(valTypes, val.GetType())
	for i := range elems {
		if i == 0 {
			continue
		}
		val = w.GetNodeActualValue(&elems[i], scope)
		valTypes = append(valTypes, val.GetType())
		if !TypeEquals(valTypes[i-1], valTypes[i]) {
			w.AlertSingle(&alerts.MixedMapOrListContents{}, elems[i].GetToken(),
				"list",
				valTypes[i-1].ToString(),
				valTypes[i].ToString(),
			)
			return InvalidType
		}
	}

	return valTypes[0]
}

func (w *Walker) GetTypeFromString(str string) ast.PrimitiveValueType {
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
		return ast.AnonStruct
	default:
		return ast.Invalid
	}
}

func (w *Walker) DetermineCallTypeString(callType ProcedureType) string {
	if callType == Function {
		return "function"
	}

	return "method"
}

func (w *Walker) TypesToValues(types []Type) Values {
	vals := Values{}

	for _, v := range types {
		vals = append(vals, w.TypeToValue(v))
	}

	return vals
}

func (w *Walker) ReportExits(sender ExitableTag, scope *Scope) {
	receiver_ := scope.ResolveReturnable()

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

func DetermineCallTypeString(callType ProcedureType) string {
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
