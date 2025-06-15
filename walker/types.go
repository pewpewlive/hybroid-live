package walker

import (
	"hybroid/ast"
	"hybroid/core"
	"hybroid/tokens"
)

type Type interface {
	PVT() ast.PrimitiveValueType
	GetType() ValueType
	//DO NOT USE ON ITS OWN, USE TypeEquals() INSTEAD
	_eq(other Type) bool
	String() string
}

type ValueType int

const (
	Basic ValueType = iota
	Fn
	Strct
	Named
	Fixed
	Wrapper // List or Map
	RawEntity
	Enum
	Alias
	Variadic
	Generic
	Path
	NA
)

type VariadicType struct {
	Type Type
}

func NewVariadicType(typ Type) *VariadicType {
	return &VariadicType{
		Type: typ,
	}
}

func (vt *VariadicType) PVT() ast.PrimitiveValueType {
	return vt.Type.PVT()
}

func (vt *VariadicType) GetType() ValueType {
	return Variadic
}

func (vt *VariadicType) _eq(other Type) bool {
	path := other.(*VariadicType)
	return vt.Type == path.Type
}

func (vt *VariadicType) String() string {
	return "..." + vt.Type.String()
}

type PathType struct {
	Env ast.Env
}

func NewPathType(envType ast.Env) *PathType {
	return &PathType{
		Env: envType,
	}
}

func (pt *PathType) PVT() ast.PrimitiveValueType {
	return ast.Path
}

func (pt *PathType) GetType() ValueType {
	return Path
}

func (pt *PathType) _eq(other Type) bool {
	path := other.(*PathType)
	return pt.Env == path.Env
}

func (pt *PathType) String() string {
	return string(pt.Env)
}

type AliasType struct {
	Token          tokens.Token
	Name           string
	UnderlyingType Type
	IsUsed         bool
	IsLocal        bool
}

func NewAliasType(name string, underlyingType Type, isLocal bool) *AliasType {
	return &AliasType{
		Name:           name,
		UnderlyingType: underlyingType,
		IsLocal:        isLocal,
	}
}

func (at *AliasType) PVT() ast.PrimitiveValueType {
	return at.UnderlyingType.PVT()
}

func (at *AliasType) GetType() ValueType {
	return at.UnderlyingType.GetType()
}

func (at *AliasType) _eq(other Type) bool {
	return TypeEquals(other, at.UnderlyingType)
}

func (at *AliasType) String() string {
	return at.Name + "(alias for " + at.UnderlyingType.String() + ")"
}

type FunctionType struct {
	Params   []Type
	Returns  []Type
	ProcType ProcedureType
}

func NewFunctionType(params []Type, returns []Type, procType ...ProcedureType) *FunctionType {
	pt := Function
	if procType != nil {
		pt = procType[0]
	}
	return &FunctionType{
		Params:   params,
		Returns:  returns,
		ProcType: pt,
	}
}

func (ft *FunctionType) PVT() ast.PrimitiveValueType {
	return ast.Func
}

func (ft *FunctionType) GetType() ValueType {
	return Fn
}

func (ft *FunctionType) _eq(other Type) bool {
	otherFT := other.(*FunctionType)
	if len(ft.Params) != len(otherFT.Params) {
		return false
	}
	for i := range ft.Params {
		if !TypeEquals(ft.Params[i], otherFT.Params[i]) {
			return false
		}
	}
	if len(ft.Returns) != len(otherFT.Returns) {
		return false
	}
	for i := range ft.Returns {
		if !TypeEquals(ft.Returns[i], otherFT.Returns[i]) {
			return false
		}
	}

	return true
}

func (ft *FunctionType) String() string {
	src := core.StringBuilder{}

	src.Write("fn(")

	length := len(ft.Params)
	for i := range ft.Params {
		if i == length-1 {
			src.Write(ft.Params[i].String())
		} else {
			src.Write(ft.Params[i].String(), ", ")
		}
	}
	src.Write(")")

	if len(ft.Returns) == 0 {
		return src.String()
	}

	src.Write(" ")
	length = len(ft.Returns)
	for i := range ft.Returns {
		if i == length-1 {
			src.Write(ft.Returns[i].String())
		} else {
			src.Write(ft.Returns[i].String(), ", ")
		}
	}

	return src.String()
}

type GenericType struct {
	Name string
}

func (gt *GenericType) PVT() ast.PrimitiveValueType {
	return ast.Generic
}

func (gt *GenericType) GetType() ValueType {
	return Generic
}

func (gt *GenericType) _eq(other Type) bool {
	g := other.(*GenericType)
	return g.Name == gt.Name
}

func (gt *GenericType) String() string {
	return gt.Name
}

func NewGeneric(name string) *GenericType {
	return &GenericType{
		Name: name,
	}
}

type RawEntityType struct{}

func (ret *RawEntityType) PVT() ast.PrimitiveValueType {
	return ast.Entity
}

func (ret *RawEntityType) GetType() ValueType {
	return RawEntity
}

func (ret *RawEntityType) _eq(other Type) bool {
	otherRET := other.(*RawEntityType)
	return otherRET.GetType() == RawEntity
}

func (ret *RawEntityType) String() string {
	return string(ast.Entity)
}

type BasicType struct {
	PrimitiveType ast.PrimitiveValueType
}

func NewBasicType(pvt ast.PrimitiveValueType) *BasicType {
	return &BasicType{
		PrimitiveType: pvt,
	}
}

func (bt *BasicType) PVT() ast.PrimitiveValueType {
	return bt.PrimitiveType
}

func (bt *BasicType) GetType() ValueType {
	return Basic
}

func (bt *BasicType) _eq(other Type) bool {
	basic := other.(*BasicType)
	return bt.PrimitiveType == basic.PrimitiveType
}

func (bt *BasicType) String() string {
	return string(bt.PrimitiveType)
}

type FixedPoint struct{}

func NewFixedPointType() *FixedPoint {
	return &FixedPoint{}
}

func (fp *FixedPoint) PVT() ast.PrimitiveValueType {
	return ast.Fixed
}

func (fp *FixedPoint) GetType() ValueType {
	return Fixed
}

func (fp *FixedPoint) _eq(other Type) bool {
	return true
}

func (fp *FixedPoint) String() string {
	return string(ast.Fixed)
}

type StructType struct {
	Fields map[string]StructField
}

func NewStructType(fields []StructField) *StructType {
	mapfields := map[string]StructField{}
	for i := range fields {
		mapfields[fields[i].Var.Name] = fields[i]
	}
	return &StructType{
		Fields: mapfields,
	}
}

func (st *StructType) PVT() ast.PrimitiveValueType {
	return ast.Struct
}

func (st *StructType) GetType() ValueType {
	return Strct
}

func (st *StructType) _eq(other Type) bool {
	map1 := st.Fields
	map2 := other.(*StructType).Fields

	for k, v := range map1 {
		v2, containsK := map2[k]
		if (containsK && !TypeEquals(v.Var.GetType(), v2.Var.GetType())) || (!v.Lenient && !containsK) {
			return false
		}
	}
	return true
}

func (st *StructType) String() string {
	src := core.StringBuilder{}

	src.Write("struct{")
	length := len(st.Fields) - 1
	index := 0
	for k, v := range st.Fields {
		if v.Lenient {
			src.Write("(optional)")
		}
		_type := v.Var.Value.GetType()
		src.Write(_type.String(), " ", k)
		if index != length {
			src.Write(", ")
		}
		index++
	}
	src.Write("}")

	return src.String()
}

type GenericWithType struct {
	GenericName string
	Type        Type
}

type NamedType struct {
	Pvt      ast.PrimitiveValueType
	EnvName  string
	Name     string
	IsUsed   bool
	Generics []GenericWithType
}

func NewNamedType(envName string, name string, primitive ast.PrimitiveValueType) *NamedType {
	return &NamedType{
		EnvName:  envName,
		Name:     name,
		Pvt:      primitive,
		Generics: []GenericWithType{},
	}
}

func (nt *NamedType) PVT() ast.PrimitiveValueType {
	return nt.Pvt
}

func (nt *NamedType) GetType() ValueType {
	return Named
}

func (nt *NamedType) _eq(othr Type) bool {
	other := othr.(*NamedType)
	if len(nt.Generics) != len(other.Generics) {
		return false
	}
	for i, v := range nt.Generics {
		if v.GenericName != other.Generics[i].GenericName || !TypeEquals(v.Type, other.Generics[i].Type) {
			return false
		}
	}
	return nt.Name == other.Name
}

func (nt *NamedType) String() string {
	if len(nt.Generics) == 0 {
		return nt.Name
	}

	src := core.StringBuilder{}
	src.Write("<")
	index := 0
	for i := range nt.Generics {
		src.Write(nt.Generics[i].Type.String())
		if index != len(nt.Generics)-1 {
			src.Write(", ")
		}
		index++
	}
	src.Write(">")
	return nt.Name + src.String()
}

type EnumType struct {
	Name    string
	EnvName string
	IsUsed  bool
}

func NewEnumType(envName, name string) *EnumType {
	return &EnumType{
		Name:    name,
		EnvName: envName,
	}
}

func (et *EnumType) PVT() ast.PrimitiveValueType {
	return ast.Enum
}

func (et *EnumType) GetType() ValueType {
	return Enum
}

func (et *EnumType) _eq(other Type) bool {
	return et.Name == other.(*EnumType).Name
}

func (et *EnumType) String() string {
	return et.Name
}

type WrapperType struct {
	WrappedType Type
	Type        Type
}

func NewWrapperType(_type Type, wrapped Type) *WrapperType {
	return &WrapperType{
		Type:        _type,
		WrappedType: wrapped,
	}
}

func (wt *WrapperType) PVT() ast.PrimitiveValueType {
	return wt.Type.PVT()
}

func (wt *WrapperType) GetType() ValueType {
	return Wrapper
}

func (wt *WrapperType) _eq(othr Type) bool {
	other := othr.(*WrapperType)
	if !TypeEquals(wt.Type, other.Type) {
		return false
	}

	if !TypeEquals(wt.WrappedType, other.WrappedType) {
		return false
	}

	return true
}

func (wt *WrapperType) String() string {
	return wt.Type.String() + "<" + wt.WrappedType.String() + ">"
}

type UnknownType struct{}

var UnknownTyp = &UnknownType{}

func (ot *UnknownType) PVT() ast.PrimitiveValueType {
	return ast.Object
}

func (ot *UnknownType) GetType() ValueType {
	return NA
}

func (ot *UnknownType) _eq(_ Type) bool {
	return false
}

func (ot *UnknownType) String() string {
	return "unknown"
}

type FuncSignature struct {
	Generics []*GenericType
	Params   []Type
	Returns  []Type
}

func NewFuncSignature(generics ...*GenericType) *FuncSignature {
	return &FuncSignature{
		Generics: generics,
		Params:   []Type{},
		Returns:  []Type{},
	}
}

func (fs *FuncSignature) WithParams(params ...Type) *FuncSignature {
	fs.Params = params
	return fs
}

func (fs *FuncSignature) WithReturns(returns ...Type) *FuncSignature {
	fs.Returns = returns
	return fs
}

func (fs *FuncSignature) String() string {
	src := core.StringBuilder{}

	src.Write("fn")

	if len(fs.Generics) != 0 {
		src.Write("<")
		for i := range fs.Generics {
			src.Write(fs.Generics[i].String())
		}
		src.Write(">")
	}
	if len(fs.Params) != 0 {
		src.Write("(")
		for i := range fs.Params {
			src.Write(fs.Params[i].String())
			if i != len(fs.Params)-1 {
				src.Write(", ")
			}
		}
		src.Write(")")
	}
	retLength := len(fs.Returns)
	if retLength != 0 {
		src.Write(" -> ")
		if retLength != 1 {
			src.Write("(")
		}
		for i := range fs.Returns {
			src.Write(fs.Returns[i].String())
		}
		if retLength != 1 {
			src.Write(")")
		}
	}

	return src.String()
}

func (fs *FuncSignature) Equals(other *FuncSignature) bool {
	if len(fs.Params) != len(other.Params) {
		return false
	}
	for i := range fs.Params {
		if !TypeEquals(fs.Params[i], other.Params[i]) {
			return false
		}
	}
	if len(fs.Returns) != len(other.Returns) {
		return false
	}
	for i := range fs.Returns {
		if !TypeEquals(fs.Returns[i], other.Returns[i]) {
			return false
		}
	}

	return true
}

func TypeEquals(t Type, other Type) bool {
	ttype := t.GetType()
	othertype := other.GetType()
	tpvt := t.PVT()
	otherpvt := other.PVT()

	if tpvt == ast.Object {
		return true
	} else if otherpvt == ast.Object {
		return true
	}

	if ttype == Named && tpvt == ast.Entity && othertype == RawEntity {
		return true
	}
	if othertype == Named && otherpvt == ast.Entity && ttype == RawEntity {
		return true
	}

	if t.GetType() != other.GetType() {
		return false
	}

	return t._eq(other)
}

var InvalidType = NewBasicType(ast.Invalid)
