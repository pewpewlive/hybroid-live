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
	CstmType
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

func (self *VariadicType) PVT() ast.PrimitiveValueType {
	return self.Type.PVT()
}

func (self *VariadicType) GetType() ValueType {
	return Variadic
}

func (self *VariadicType) _eq(other Type) bool {
	path := other.(*VariadicType)
	return self.Type == path.Type
}

func (self *VariadicType) ToString() string {
	return "..." + self.Type.ToString()
}

type PathType struct {
	EnvType ast.EnvType
}

func NewPathType(envType ast.EnvType) *PathType {
	return &PathType{
		EnvType: envType,
	}
}

func (self *PathType) PVT() ast.PrimitiveValueType {
	return ast.Path
}

func (self *PathType) GetType() ValueType {
	return Path
}

func (self *PathType) _eq(other Type) bool {
	path := other.(*PathType)
	return self.EnvType == path.EnvType
}

func (self *PathType) ToString() string {
	return string(self.EnvType)
}

type CustomType struct {
	Name           string
	UnderlyingType Type
}

func NewCustomType(name string, underlyingType Type) *CustomType {
	return &CustomType{
		Name:           name,
		UnderlyingType: underlyingType,
	}
}

func (self *CustomType) PVT() ast.PrimitiveValueType {
	return self.UnderlyingType.PVT()
}

func (self *CustomType) GetType() ValueType {
	return CstmType
}

func (self *CustomType) _eq(other Type) bool {
	ret := other.(*CustomType)
	return ret.Name == self.Name
}

func (self *CustomType) ToString() string {
	return self.Name + "(" + self.UnderlyingType.ToString() + ")"
}

type FunctionType struct {
	Params  Types
	Returns Types
}

func NewFunctionType(params Types, returns Types) *FunctionType {
	return &FunctionType{
		Params:  params,
		Returns: returns,
	}
}

func (self *FunctionType) PVT() ast.PrimitiveValueType {
	return ast.Func
}

func (self *FunctionType) GetType() ValueType {
	return Fn
}

func (self *FunctionType) _eq(other Type) bool {
	ft := other.(*FunctionType)
	if len(self.Params) != len(ft.Params) {
		return false
	}
	for i := range self.Params {
		if !TypeEquals(self.Params[i], ft.Params[i]) {
			return false
		}
	}
	if len(self.Returns) != len(ft.Returns) {
		return false
	}
	for i := range self.Returns {
		if !TypeEquals(self.Returns[i], ft.Returns[i]) {
			return false
		}
	}

	return true
}

func (self *FunctionType) ToString() string {
	src := generator.StringBuilder{}

	src.WriteString("fn(")

	length := len(self.Params)
	for i := range self.Params {
		if i == length-1 {
			src.WriteString(self.Params[i].ToString())
		} else {
			src.Append(self.Params[i].ToString(), ", ")
		}
	}
	src.WriteString(")")

	if len(self.Returns) == 0 {
		return src.String()
	}

	src.WriteString(" ")
	length = len(self.Returns)
	for i := range self.Returns {
		if i == length-1 {
			src.WriteString(self.Returns[i].ToString())
		} else {
			src.Append(self.Returns[i].ToString(), ", ")
		}
	}

	return src.String()
}

type GenericType struct {
	Name string
}

func (self *GenericType) PVT() ast.PrimitiveValueType {
	return ast.Generic
}

func (self *GenericType) GetType() ValueType {
	return Generic
}

func (self *GenericType) _eq(other Type) bool {
	g := other.(*GenericType)
	return g.Name == self.Name
}

func (self *GenericType) ToString() string {
	return self.Name
}

func NewGeneric(name string) *GenericType {
	return &GenericType{
		Name: name,
	}
}

type RawEntityType struct{}

func (self *RawEntityType) PVT() ast.PrimitiveValueType {
	return ast.Entity
}

func (self *RawEntityType) GetType() ValueType {
	return RawEntity
}

func (self *RawEntityType) _eq(other Type) bool {
	ret := other.(*RawEntityType)
	return ret.GetType() == RawEntity
}

func (self *RawEntityType) ToString() string {
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
func (self *BasicType) PVT() ast.PrimitiveValueType {
	return self.PrimitiveType
}

func (self *BasicType) GetType() ValueType {
	return Basic
}

func (self *BasicType) _eq(other Type) bool {
	basic := other.(*BasicType)
	return self.PrimitiveType == basic.PrimitiveType
}

func (self *BasicType) ToString() string {
	return string(self.PrimitiveType)
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
func (self *FixedPoint) PVT() ast.PrimitiveValueType {
	return self.Specific
}

func (self *FixedPoint) GetType() ValueType {
	return Fixed
}

func (self *FixedPoint) _eq(other Type) bool {
	return true
}

func (self *FixedPoint) ToString() string {
	return string(self.Specific)
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

func (self *StructType) PVT() ast.PrimitiveValueType {
	return ast.AnonStruct
}

func (self *StructType) GetType() ValueType {
	return Strct
}

func (self *StructType) _eq(other Type) bool {
	map1 := self.Fields
	map2 := other.(*StructType).Fields
	if self.Lenient {
		return other._eq(self)
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

func (self *StructType) ToString() string {
	src := generator.StringBuilder{}

	src.WriteString("struct{")
	length := len(self.Fields) - 1
	index := 0
	for k, v := range self.Fields {
		if index == length {
			_type := v.Var.Value.GetType()
			src.Append(_type.ToString(), " ", k)
		} else {
			_type := v.Var.Value.GetType()
			src.Append(_type.ToString(), " ", k, ", ")
		}
		index++
	}
	src.WriteString("}")

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
func (self *NamedType) PVT() ast.PrimitiveValueType {
	return self.Pvt
}

func (self *NamedType) GetType() ValueType {
	return Named
}

func (self *NamedType) _eq(othr Type) bool {
	other := othr.(*NamedType)
	return self.Name == other.Name
}

func (self *NamedType) ToString() string {
	return self.Name
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

func (self *EnumType) PVT() ast.PrimitiveValueType {
	return ast.Enum
}

func (self *EnumType) GetType() ValueType {
	return Enum
}

func (self *EnumType) _eq(other Type) bool {
	return self.Name == other.(*EnumType).Name
}

func (self *EnumType) ToString() string {
	return self.Name
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
func (self *WrapperType) PVT() ast.PrimitiveValueType {
	return self.Type.PVT()
}

func (self *WrapperType) GetType() ValueType {
	return Wrapper
}

func (self *WrapperType) _eq(othr Type) bool {
	other := othr.(*WrapperType)
	if !TypeEquals(self.Type, other.Type) {
		return false
	}

	if !TypeEquals(self.WrappedType, other.WrappedType) {
		return false
	}

	return true
}

func (self *WrapperType) ToString() string {
	return self.Type.ToString() + "<" + self.WrappedType.ToString() + ">"
}

type ObjectType struct{}

var ObjectTyp = &ObjectType{}

// Type
func (self *ObjectType) PVT() ast.PrimitiveValueType {
	return ast.Object
}

func (self *ObjectType) GetType() ValueType {
	return NA
}

func (self *ObjectType) _eq(_ Type) bool {
	return false
}

func (self *ObjectType) ToString() string {
	return "NotAnyType"
}

type FuncSignature struct {
	Generics []*GenericType
	Params []Type
	Returns []Type
}

func NewFuncSignature(generics ...*GenericType) *FuncSignature {
	return &FuncSignature{
		Generics: generics,
		Params: []Type{},
		Returns: []Type{},
	}
}

func (self *FuncSignature) WithParams(params ...Type) *FuncSignature {
	self.Params = params
	return self
}

func (self *FuncSignature) WithReturns(returns ...Type) *FuncSignature {
	self.Returns = returns
	return self
}

func (self *FuncSignature) ToString() string {
	src := generator.StringBuilder{}

	src.WriteString("fn")

	if len(self.Generics) != 0 {
		src.WriteString("<")
		for i := range self.Generics {
			src.WriteString(self.Generics[i].ToString())
		}
		src.WriteString(">")
	}
	if len(self.Params) != 0 {
		src.WriteString("(")
		for i := range self.Params {
			src.WriteString(self.Params[i].ToString())
			if i != len(self.Params)-1 {
				src.WriteString(", ")
			}
		}
		src.WriteString(")")
	}
	retLength := len(self.Returns)
	if retLength != 0 {
		src.WriteString(" -> ")
		if retLength != 1 {
			src.WriteString("(")
		}
		for i := range self.Returns {
			src.WriteString(self.Returns[i].ToString())
		}
		if retLength != 1 {
			src.WriteString(")")
		}
	}

	return src.String()
}

func (self *FuncSignature) Equals(other *FuncSignature) bool {
	if len(self.Params) != len(other.Params) {
		return false
	}
	for i := range self.Params {
		if !TypeEquals(self.Params[i], other.Params[i]) {
			return false
		}
	}
	if len(self.Returns) != len(other.Returns) {
		return false
	}
	for i := range self.Returns {
		if !TypeEquals(self.Returns[i], other.Returns[i]) {
			return false
		}
	}

	return true
}

func TypeEquals(t Type, other Type) bool {
	tpvt := t.PVT()
	otherpvt := other.PVT()

	if tpvt == ast.Object {
		return true
	} else if otherpvt == ast.Object {
		return true
	}

	if t.GetType() == Named && t.PVT() == ast.Entity && other.GetType() == RawEntity {
		return true
	}
	if other.GetType() == Named && other.PVT() == ast.Entity && t.GetType() == RawEntity {
		return true
	}

	if t.GetType() != other.GetType() {
		return false
	}

	return t._eq(other)
}

var InvalidType = NewBasicType(ast.Invalid)

var MeshValueType = NewStructType([]*VariableVal{
	{
		Name:  "vertexes",
		Value: &ListVal{ValueType: NewWrapperType(NewBasicType(ast.List), NewBasicType(ast.Number))},
	},
	{
		Name:  "segments",
		Value: &ListVal{ValueType: NewWrapperType(NewBasicType(ast.List), NewBasicType(ast.Number))},
	},
	{
		Name:  "colors",
		Value: &ListVal{ValueType: NewBasicType(ast.Number)},
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