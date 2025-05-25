package walker

import (
	"hybroid/ast"
	"hybroid/generator"
)

type Type interface {
	PVT() ast.PrimitiveValueType
	GetType() ValueType
	//DO NOT USE ON ITS OWN, USE TypeEquals() INSTEAD
	_eq(other Type) bool
	ToString() string
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
	NotKnown
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

func (vt *VariadicType) ToString() string {
	return "..." + vt.Type.ToString()
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

func (pt *PathType) ToString() string {
	return string(pt.Env)
}

type AliasType struct {
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

func (at *AliasType) ToString() string {
	return at.Name + "(alias for " + at.UnderlyingType.ToString() + ")"
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

func (ft *FunctionType) ToString() string {
	src := generator.StringBuilder{}

	src.Write("fn(")

	length := len(ft.Params)
	for i := range ft.Params {
		if i == length-1 {
			src.Write(ft.Params[i].ToString())
		} else {
			src.Write(ft.Params[i].ToString(), ", ")
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
			src.Write(ft.Returns[i].ToString())
		} else {
			src.Write(ft.Returns[i].ToString(), ", ")
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

func (gt *GenericType) ToString() string {
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

func (ret *RawEntityType) ToString() string {
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

// Type
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

func (bt *BasicType) ToString() string {
	return string(bt.PrimitiveType)
}

type FixedPoint struct {
	Specific ast.PrimitiveValueType
}

func NewFixedPointType(specific ast.PrimitiveValueType) *FixedPoint {
	return &FixedPoint{
		Specific: specific,
	}
}

// Type
func (fp *FixedPoint) PVT() ast.PrimitiveValueType {
	return fp.Specific
}

func (fp *FixedPoint) GetType() ValueType {
	return Fixed
}

func (fp *FixedPoint) _eq(other Type) bool {
	return true
}

func (fp *FixedPoint) ToString() string {
	return string(fp.Specific)
}

type StructType struct {
	Fields  map[string]Field
	Lenient bool
}

func NewStructType(fields []*VariableVal, lenient bool) *StructType {
	mapfields := map[string]Field{}
	for i := range fields {
		mapfields[fields[i].Name] = Field{Var: fields[i], Index: i}
	}
	return &StructType{
		Fields:  mapfields,
		Lenient: lenient,
	}
}

func NewStructTypeWithFields(fields map[string]Field, lenient bool) *StructType {
	return &StructType{
		Fields:  fields,
		Lenient: lenient,
	}
}

func (st *StructType) PVT() ast.PrimitiveValueType {
	return ast.AnonStruct
}

func (st *StructType) GetType() ValueType {
	return Strct
}

func (st *StructType) _eq(other Type) bool {
	map1 := st.Fields
	map2 := other.(*StructType).Fields
	if st.Lenient {
		return other._eq(st)
	}

	for k, v := range map1 {
		containsK := false
		for k2, v2 := range map2 {
			if k == k2 && TypeEquals(v.Var.GetType(), v2.Var.GetType()) {
				containsK = true
				break
			}
		}
		if !containsK {
			return false
		}
	}
	return true
}

func (st *StructType) ToString() string {
	src := generator.StringBuilder{}

	src.Write("struct{")
	length := len(st.Fields) - 1
	index := 0
	for k, v := range st.Fields {
		if index == length {
			_type := v.Var.Value.GetType()
			src.Write(_type.ToString(), " ", k)
		} else {
			_type := v.Var.Value.GetType()
			src.Write(_type.ToString(), " ", k, ", ")
		}
		index++
	}
	src.Write("}")

	return src.String()
}

type NamedType struct {
	Pvt     ast.PrimitiveValueType
	EnvName string
	Name    string
	IsUsed  bool
}

func NewNamedType(envName string, name string, primitive ast.PrimitiveValueType) *NamedType {
	return &NamedType{
		EnvName: envName,
		Name:    name,
		Pvt:     primitive,
	}
}

// Type
func (nt *NamedType) PVT() ast.PrimitiveValueType {
	return nt.Pvt
}

func (nt *NamedType) GetType() ValueType {
	return Named
}

func (nt *NamedType) _eq(othr Type) bool {
	other := othr.(*NamedType)
	return nt.Name == other.Name
}

func (nt *NamedType) ToString() string {
	return nt.Name
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

func (et *EnumType) ToString() string {
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

// Type
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

func (wt *WrapperType) ToString() string {
	return wt.Type.ToString() + "<" + wt.WrappedType.ToString() + ">"
}

type ObjectType struct{}

var ObjectTyp = &ObjectType{}

// Type
func (ot *ObjectType) PVT() ast.PrimitiveValueType {
	return ast.Object
}

func (ot *ObjectType) GetType() ValueType {
	return NA
}

func (ot *ObjectType) _eq(_ Type) bool {
	return false
}

func (ot *ObjectType) ToString() string {
	return "NotAnyType"
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

func (fs *FuncSignature) ToString() string {
	src := generator.StringBuilder{}

	src.Write("fn")

	if len(fs.Generics) != 0 {
		src.Write("<")
		for i := range fs.Generics {
			src.Write(fs.Generics[i].ToString())
		}
		src.Write(">")
	}
	if len(fs.Params) != 0 {
		src.Write("(")
		for i := range fs.Params {
			src.Write(fs.Params[i].ToString())
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
			src.Write(fs.Returns[i].ToString())
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

	/*
		if (ttype == Fixed && otherpvt == ast.Number) || (othertype == Fixed && tpvt == ast.Number) {
			return true
		}
	*/

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

var numberListVal = &ListVal{ValueType: NewBasicType(ast.Number)}
var vertexesVal = &ListVal{ValueType: numberListVal.GetType()}
var MeshValueType = NewStructType([]*VariableVal{
	{
		Name:  "vertexes",
		Value: vertexesVal,
	},
	{
		Name:  "segments",
		Value: vertexesVal,
	},
	{
		Name:  "colors",
		Value: numberListVal,
	},
}, false)
var MeshesValueType = (&ListVal{ValueType: MeshValueType}).GetType()

var SoundValueType = NewStructType([]*VariableVal{
	{
		Name:  "attack",
		Value: &NumberVal{},
	},
	{
		Name:  "decay",
		Value: &NumberVal{},
	},
	{
		Name:  "sustain",
		Value: &NumberVal{},
	},
	{
		Name:  "sustainPunch",
		Value: &NumberVal{},
	},
	{
		Name:  "amplification",
		Value: &NumberVal{},
	},
	{
		Name:  "harmonics",
		Value: &NumberVal{},
	},
	{
		Name:  "harmonicsFalloff",
		Value: &NumberVal{},
	},
	{
		Name:  "tremoloDepth",
		Value: &NumberVal{},
	},
	{
		Name:  "tremoloFrequency",
		Value: &NumberVal{},
	},
	{
		Name:  "frequency",
		Value: &NumberVal{},
	},
	{
		Name:  "frequencyDeltaSweep",
		Value: &NumberVal{},
	},
	{
		Name:  "frequencyJump1Onset",
		Value: &NumberVal{},
	},
	{
		Name:  "frequencyJump2Onset",
		Value: &NumberVal{},
	},
	{
		Name:  "frequencyJump1Amount",
		Value: &NumberVal{},
	},
	{
		Name:  "frequencyJump2Amount",
		Value: &NumberVal{},
	},
	{
		Name:  "frequencySweep",
		Value: &NumberVal{},
	},
	{
		Name:  "vibratoFrequency",
		Value: &NumberVal{},
	},
	{
		Name:  "vibratoDepth",
		Value: &NumberVal{},
	},
	{
		Name:  "flangerOffset",
		Value: &NumberVal{},
	},
	{
		Name:  "flangerOffsetSweep",
		Value: &NumberVal{},
	},
	{
		Name:  "repeatFrequency",
		Value: &NumberVal{},
	},
	{
		Name:  "lowPassCutoff",
		Value: &NumberVal{},
	},
	{
		Name:  "lowPassCutoffSweep",
		Value: &NumberVal{},
	},
	{
		Name:  "highPassCutoff",
		Value: &NumberVal{},
	},
	{
		Name:  "highPassCutoffSweep",
		Value: &NumberVal{},
	},
	{
		Name:  "bitCrush",
		Value: &NumberVal{},
	},
	{
		Name:  "bitCrushSweep",
		Value: &NumberVal{},
	},
	{
		Name:  "squareDuty",
		Value: &NumberVal{},
	},
	{
		Name:  "squareDutySweep",
		Value: &NumberVal{},
	},
	{
		Name:  "harmonicsFalloff",
		Value: &NumberVal{},
	},
	{
		Name:  "normalization",
		Value: &BoolVal{},
	},
	{
		Name:  "interpolateNoise",
		Value: &BoolVal{},
	},
	{
		Name:  "compression",
		Value: &NumberVal{},
	},
	{
		Name:  "harmonics",
		Value: &NumberVal{},
	},
	{
		Name:  "harmonicsFalloff",
		Value: &NumberVal{},
	},
	{
		Name:  "repeatFrequency",
		Value: &NumberVal{},
	},
	{
		Name:  "sampleRate",
		Value: &NumberVal{},
	},
	{
		Name:  "waveform",
		Value: &StringVal{},
	},
}, true)
var SoundsValueType = (&ListVal{ValueType: SoundValueType}).GetType()

var WeaponCollisionSign = NewFuncSignature().
	WithParams(NewBasicType(ast.Number), NewEnumType("Pewpew", "WeaponType")).
	WithReturns(NewBasicType(ast.Bool))

var PlayerCollisionSign = NewFuncSignature().
	WithParams(NewBasicType(ast.Number), &RawEntityType{})

var WallCollisionSign = NewFuncSignature().
	WithParams(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed))
