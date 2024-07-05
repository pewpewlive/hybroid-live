package walker

import (
	"fmt"
	"hybroid/ast"
	"hybroid/helpers"
	"hybroid/lexer"
)

type Path struct {
	Src string
	Dst string
}

type Environment struct {
	Name      string
	Path      Path
	Scope     Scope
	Variables map[string]*VariableVal
	Structs   map[string]*StructVal
}

func NewEnvironment(path Path) *Environment {
	scope := Scope{
		Children: make([]*Scope, 0),

		Tag:       &UntaggedTag{},
		Variables: map[string]*VariableVal{},
	}
	global := &Environment{
		Path:    path,
		Scope:   scope,
		Structs: map[string]*StructVal{},
	}

	global.Scope.Environment = global
	return global
}

type Walker struct {
	Environment *Environment
	Walkers     *map[string]*Walker
	UsedWalkers    []*Walker
	Nodes       []ast.Node
	Errors      []ast.Error
	Warnings    []ast.Warning
	Context     Context
}

func NewWalker(path Path) *Walker {
	environment := NewEnvironment(path)
	walker := Walker{
		Environment: environment,
		Nodes:       []ast.Node{},
		Errors:      []ast.Error{},
		Warnings:    []ast.Warning{},
		Context: Context{
			Node:  &ast.Improper{},
			Value: &Unknown{},
			Ret:   Types{},
		},
	}
	return &walker
}

func (w *Walker) Error(token lexer.Token, msg string) {
	w.Errors = append(w.Errors, ast.Error{Token: token, Message: msg})
}

func (w *Walker) Warn(token lexer.Token, msg string) {
	w.Warnings = append(w.Warnings, ast.Warning{Token: token, Message: msg})
}

func (w *Walker) AddError(err ast.Error) {
	w.Errors = append(w.Errors, err)
}

func (s *Scope) GetVariable(name string) *VariableVal {
	variable := s.Variables[name]

	if variable == nil {
		return nil
	}

	variable.IsUsed = true

	s.Variables[name] = variable

	return s.Variables[name]
}

func (w *Walker) GetStruct(name string) (*StructVal, bool) {
	structType, found := w.Environment.Structs[name]
	if !found {
		//w.Error(w.Context.Node.GetToken(), fmt.Sprintf("no struct named %s", name, " exists"))
		return nil, false
	}

	structType.Type.IsUsed = true

	w.Environment.Structs[name] = structType

	return structType, true
}

func (s *Scope) AssignVariableByName(name string, value Value) (Value, *ast.Error) {
	scope := s.ResolveVariable(name)

	if scope == nil {
		return &Invalid{}, &ast.Error{Message: "cannot assign to an undeclared variable"}
	}

	variable := scope.Variables[name]
	if variable.IsConst {
		return &Invalid{}, &ast.Error{Message: "cannot assign to a constant variable"}
	}

	variable.Value = value

	scope.Variables[name] = variable

	temp := scope.Variables[name]

	return temp, nil
}

func (s *Scope) AssignVariable(variable *VariableVal, value Value) (Value, *ast.Error) {
	if variable.IsConst {
		return &Invalid{}, &ast.Error{Message: "cannot assign to a constant variable"}
	}

	variable.Value = value

	return variable, nil
}

func (w *Walker) DeclareVariable(s *Scope, value *VariableVal, token lexer.Token) (*VariableVal, bool) {
	if varFound, found := s.Variables[value.Name]; found {
		w.Error(token, fmt.Sprintf("variable with name '%s' already exists", varFound.Name))
		return varFound, false
	}

	s.Variables[value.Name] = value
	return value, true
}

func (w *Walker) DeclareStruct(structVal *StructVal) bool {
	if _, found := w.Environment.Structs[structVal.Type.Name]; found {
		return false
	}

	w.Environment.Structs[structVal.Type.Name] = structVal
	return true
}

func (s *Scope) ResolveVariable(name string) *Scope {
	if _, found := s.Variables[name]; found {
		return s
	}

	if s.Parent == nil {
		return nil
	}

	return s.Parent.ResolveVariable(name)
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

func (w *Walker) ValidateArithmeticOperands(left Type, right Type, expr ast.BinaryExpr) bool {
	//fmt.Printf("Validating operands: %v (%v) and %v (%v)\n", left.Val, left.Type, right.Val, right.Type)
	if left.PVT() == ast.Invalid {
		w.Error(expr.Left.GetToken(), "cannot perform arithmetic on Invalid value")
		return false
	}

	if right.PVT() == ast.Invalid {
		w.Error(expr.Right.GetToken(), "cannot perform arithmetic on Invalid value")
		return false
	}

	switch left.PVT() {
	case ast.List, ast.Map, ast.String, ast.Bool, ast.Entity, ast.Struct:
		w.Error(expr.Left.GetToken(), "cannot perform arithmetic on a non-number value")
		return false
	}

	switch right.PVT() {
	case ast.List, ast.Map, ast.String, ast.Bool, ast.Entity, ast.Struct:
		w.Error(expr.Right.GetToken(), "cannot perform arithmetic on a non-number value")
		return false
	}

	return true
}

func returnsAreValid(list1 []Type, list2 []Type) bool {
	if len(list1) != len(list2) {
		return false
	}

	for i, v := range list1 {
		//fmt.Printf("%s compared to %s\n", list1[i].ToString(), list2[i].ToString())
		if !TypeEquals(v, list2[i]) {
			return false
		}
	}
	return true
}

func (w *Walker) ValidateReturnValues(_return Types, expectReturn Types) string {
	returnValues, expectedReturnValues := _return, expectReturn
	if len(returnValues) < len(expectedReturnValues) {
		return "not enough return values given"
	} else if len(returnValues) > len(expectedReturnValues) {
		return "too many return values given"
	}
	if !returnsAreValid(returnValues, expectedReturnValues) {
		return "invalid return type(s)"
	}
	return ""
}

func (w *Walker) TypeToValue(_type Type) Value {
	switch _type.PVT() {
	case ast.Radian, ast.Fixed, ast.FixedPoint, ast.Degree:
		return &FixedVal{SpecificType: _type.PVT()}
	case ast.Bool:
		return &BoolVal{}
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
		val, _ := w.GetStruct(_type.ToString())
		return val
	case ast.AnonStruct:
		return &AnonStructVal{
			Fields: _type.(*AnonStructType).Fields,
		}
	case ast.Enum:
		return w.Environment.Scope.GetVariable(_type.(*EnumType).Name)
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
	receiver.SetExit(sender.GetIfExits(Yield), All)
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
