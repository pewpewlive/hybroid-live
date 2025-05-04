package walker

import (
	"fmt"
	"hybroid/alerts"
	"hybroid/ast"
	"hybroid/helpers"
	"hybroid/tokens"
	"slices"
)

var LibraryEnvs = map[Library]*Environment{
	Pewpew: PewpewEnv,
	Fmath:  FmathEnv,
	Math:   MathEnv,
	String: StringEnv,
	Table:  TableEnv,
}

type Environment struct {
	Name            string
	Path            string // dynamic lua path
	Type            ast.EnvType
	Scope           Scope
	UsedWalkers     []*Walker
	UsedLibraries   map[Library]bool
	UsedBuiltinVars []string
	Structs         map[string]*ClassVal
	Entities        map[string]*EntityVal
	CustomTypes     map[string]*CustomType
	AliasTypes      map[string]*AliasType
}

func (e *Environment) AddBuiltinVar(name string) {
	if slices.Contains(e.UsedBuiltinVars, name) {
		return
	}

	e.UsedBuiltinVars = append(e.UsedBuiltinVars, name)
}

func NewEnvironment(path string) *Environment {
	scope := Scope{
		Tag:       &UntaggedTag{},
		Variables: map[string]*VariableVal{},
	}
	global := &Environment{
		Path:  path,
		Type:  ast.InvalidEnv,
		Scope: scope,
		UsedLibraries: map[Library]bool{
			Pewpew: false,
			Table:  false,
			String: false,
			Math:   false,
			Fmath:  false,
		},
		Structs:     map[string]*ClassVal{},
		Entities:    map[string]*EntityVal{},
		CustomTypes: map[string]*CustomType{},
		AliasTypes:  make(map[string]*AliasType),
	}

	global.Scope.Environment = global
	return global
}

type Library int

const (
	Pewpew Library = iota
	Fmath
	Math
	String
	Table
)

type Walker struct {
	alerts.Collector

	CurrentEnvironment *Environment
	Environment        *Environment
	Walkers            map[string]*Walker
	Nodes              []ast.Node
	Context            Context
	Walked             bool
}

// var pewpewEnv = &Environment{
// 	Path: "pewpew_path",
// 	Variables: map[string]*VariableVal{
// 		"WeaponType":
// 	},
// }

func NewWalker(path string) *Walker {
	walker := &Walker{
		Environment: NewEnvironment(path),
		Nodes:       []ast.Node{},
		Context: Context{
			Node:   &ast.Improper{},
			Value:  &Unknown{},
			Value2: &Unknown{},
		},
		Collector: alerts.NewCollector(),
	}
	walker.CurrentEnvironment = walker.Environment
	return walker
}

func (w *Walker) GetEnvStmt() *ast.EnvironmentDecl {
	return w.Nodes[0].(*ast.EnvironmentDecl)
}

// ONLY CALL THIS IF YOU ALREADY CALLED ResolveVariable
func (w *Walker) GetVariable(s *Scope, name string) (*VariableVal, bool) {
	variable, ok := s.Variables[name]

	if !ok {
		return nil, false
	}
	return variable, s.Environment.Name != w.Environment.Name && variable.IsLocal
}

func (w *Walker) TypeExists(name string) bool {
	if _, found := w.GetEntity(name); found {
		return true
	}
	if _, found := w.GetStruct(name); found {
		return true
	}

	return false
}

func (w *Walker) GetStruct(name string) (*ClassVal, bool) {
	structType, found := w.Environment.Structs[name]
	if !found {
		// w.Error(w.Context.Node.GetToken(), fmt.Sprintf("no struct named %s exists", name))
		return nil, false
	}

	structType.Type.IsUsed = true

	w.Environment.Structs[name] = structType

	return structType, true
}

func (w *Walker) GetEntity(name string) (*EntityVal, bool) {
	entityType, found := w.Environment.Entities[name]
	if !found {
		// w.Error(w.Context.Node.GetToken(), fmt.Sprintf("no struct named %s exists", name))
		return nil, false
	}

	entityType.Type.IsUsed = true

	w.Environment.Entities[name] = entityType

	return entityType, true
}

func (w *Walker) AssignVariableByName(s *Scope, name string, value Value) Value {
	scope := w.ResolveVariable(s, name)

	if scope == nil {
		// return &Invalid{}, &ast.Error{Message: "cannot assign to an undeclared variable"}
	}

	variable := scope.Variables[name]
	if variable.IsConst {
		// return &Invalid{}, &ast.Error{Message: "cannot assign to a constant variable"}
	}

	variable.Value = value

	scope.Variables[name] = variable

	temp := scope.Variables[name]

	return temp
}

func (s *Scope) AssignVariable(variable *VariableVal, value Value) Value {
	if variable.IsConst {
		// return &Invalid{}, &ast.Error{Message: "cannot assign to a constant variable"}
	}

	//variable.Value = value

	return variable
}

func (w *Walker) DeclareVariable(s *Scope, value *VariableVal, token tokens.Token) (*VariableVal, bool) {
	if varFound, found := s.Variables[value.Name]; found {
		// w.Error(token, fmt.Sprintf("variable with name '%s' already exists", varFound.Name))
		return varFound, false
	}

	s.Variables[value.Name] = value
	return value, true
}

func (w *Walker) DeclareClass(structVal *ClassVal) bool {
	if _, found := w.Environment.Structs[structVal.Type.Name]; found {
		return false
	}

	w.Environment.Structs[structVal.Type.Name] = structVal
	return true
}

func (w *Walker) DeclareEntity(entityVal *EntityVal) bool {
	if _, found := w.Environment.Entities[entityVal.Type.Name]; found {
		return false
	}

	w.Environment.Entities[entityVal.Type.Name] = entityVal
	return true
}

func (w *Walker) ResolveVariable(s *Scope, name string) *Scope {
	if _, found := s.Variables[name]; found {
		return s
	}

	if s.Parent == nil {
		_, ok := BuiltinEnv.Scope.Variables[name]
		if ok {
			return &BuiltinEnv.Scope
		}
		for i := range s.Environment.UsedWalkers {
			variable, _ := s.Environment.UsedWalkers[i].GetVariable(&s.Environment.UsedWalkers[i].Environment.Scope, name)
			if variable != nil {
				return &s.Environment.UsedWalkers[i].Environment.Scope
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

	return w.ResolveVariable(s.Parent, name)
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

func (w *Walker) ValidateArguments(generics map[string]Type, args []Type, params []Type, callToken tokens.Token) (int, bool) {

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
			// w.Error(callToken, fmt.Sprintf("argument is of type %s, but should be %s", typeVal.ToString(), param.ToString()))
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

func returnsAreValid(list1 []Type, list2 []Type) bool {
	if len(list1) != len(list2) {
		return false
	}

	for i, v := range list1 {
		if !TypeEquals(v, list2[i]) {
			return false
		}
	}
	return true
}

func (w *Walker) ValidateReturnValues(_return Types, expectReturn Types) string {
	returnValues, expectedReturnValues := _return, expectReturn
	if len(returnValues) < len(expectedReturnValues) { // debug?
		return "not enough return values given"
	} else if len(returnValues) > len(expectedReturnValues) {
		return "too many return values given"
	}
	if !returnsAreValid(returnValues, expectedReturnValues) {
		return "invalid return type(s)"
	}
	return ""
}

func (w *Walker) CheckAccessibility(s *Scope, isLocal bool, node ast.Node) {
	if s.Environment.Name != w.Environment.Name && isLocal {
		// w.Error(node.GetToken(), "Not allowed to access a local variable/type from a different environment")
	}
}

func (w *Walker) TypeToValue(_type Type) Value {
	if _type.GetType() == RawEntity {
		return &RawEntityVal{}
	}
	if _type.GetType() == CstmType {
		return NewCustomVal(_type.(*CustomType))
	}
	if _type.GetType() == Variadic {
		return &ListVal{
			ValueType: _type.(*VariadicType).Type,
			Values:    make([]Value, 0),
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
		val, _ := w.Walkers[_type.(*NamedType).EnvName].GetStruct(_type.ToString())
		return val
	case ast.AnonStruct:
		return &AnonStructVal{
			Fields: _type.(*StructType).Fields,
		}
	case ast.Enum:
		enum := _type.(*EnumType)
		walker, found := w.Walkers[enum.EnvName]
		var variable *VariableVal
		switch enum.EnvName {
		case "Pewpew":
			variable = PewpewEnv.Scope.Variables[enum.Name]
		default:
			variable, _ = walker.GetVariable(&walker.Environment.Scope, enum.Name)
		}
		if variable == nil {
			panic(fmt.Sprintf("Enum variable could not be found when converting enum type to value (envName:%v, name:%v, walkerFound:%v)", enum.EnvName, enum.Name, found))
		}
		return variable
	case ast.Entity:
		val, _ := w.Walkers[_type.(*NamedType).EnvName].GetEntity(_type.ToString())
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

func (w *Walker) AddTypesToValues(list *[]Value, tys *Types) {
	for _, typ := range *tys {
		val := w.TypeToValue(typ)
		*list = append(*list, val)
	}
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
