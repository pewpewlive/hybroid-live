package walker

import (
	"hybroid/ast"
	"hybroid/generator"
	"hybroid/helpers"
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
	AnonStruct
	Structure
	PewpewEntity
	Fixed
	Wrapper // List or Map
	Enum
	Unresolved
	NA
	NotKnown
)

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
	if !helpers.ListsAreSame(self.Params, ft.Params) {
		return false
	}
	if !helpers.ListsAreSame(self.Returns, ft.Returns) {
		return false
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
	length = len(self.Params)
	for i := range self.Returns {
		if i == length-1 {
			src.WriteString(self.Returns[i].ToString())
		} else {
			src.Append(self.Returns[i].ToString(), ", ")
		}
	}

	return src.String()
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

type AnonStructType struct {
	Fields map[string]*VariableVal
}

func NewAnonStructType(fields map[string]*VariableVal) *AnonStructType {
	return &AnonStructType{
		Fields: fields,
	}
}

func (self *AnonStructType) PVT() ast.PrimitiveValueType {
	return ast.AnonStruct
}

func (self *AnonStructType) GetType() ValueType {
	return AnonStruct
}

func (self *AnonStructType) _eq(other Type) bool {
	map1 := self.Fields
	map2 := other.(*AnonStructType).Fields

	for k, v := range map1 {
		containsK := false
		for k2, v2 := range map2 {
			if k == k2 && TypeEquals(v.GetType(), v2.GetType()) {
				containsK = true
			}
		}
		if !containsK {
			return false
		}
	}
	return true
}

func (self *AnonStructType) ToString() string {
	src := generator.StringBuilder{}

	src.WriteString("struct{")
	length := len(self.Fields) - 1
	index := 0
	for k, v := range self.Fields {
		if index == length {
			_type := v.Value.GetType()
			src.Append(_type.ToString(), " ", k)
		} else {
			_type := v.Value.GetType()
			src.Append(_type.ToString(), " ", k, ", ")
		}
		index++
	}
	src.WriteString("}")

	return src.String()
}

type NamedType struct {
	Name   string
	IsUsed bool
}

func NewNamedType(name string) *NamedType {
	return &NamedType{
		Name: name,
	}
}

// Type
func (self *NamedType) PVT() ast.PrimitiveValueType {
	return ast.Struct
}

func (self *NamedType) GetType() ValueType {
	return Structure
}

func (self *NamedType) _eq(othr Type) bool {
	other := othr.(*NamedType)
	if self.Name != other.Name {
		return false
	}

	return true
}

func (self *NamedType) ToString() string {
	return self.Name
}

type UnresolvedType struct {
	Expr ast.Node
}

func (self *UnresolvedType) PVT() ast.PrimitiveValueType {
	return ast.Unresolved
}

func (self *UnresolvedType) GetType() ValueType {
	return Unresolved
}

func (self *UnresolvedType) _eq(othr Type) bool {
	other := othr.(*UnresolvedType)
	if self.Expr != other.Expr {
		return false
	}

	return true
}

func (self *UnresolvedType) ToString() string {
	return "unresolved"
}

type EnumType struct {
	Name   string
	IsUsed bool
}

func NewEnumType(name string) *EnumType {
	return &EnumType{
		Name: name,
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
	if self.Type.GetType() != other.Type.GetType() {
		return false
	}

	if self.WrappedType.GetType() != other.WrappedType.GetType() {
		return false
	}

	if !self.Type._eq(other.Type) {
		return false
	}

	if !self.WrappedType._eq(other.WrappedType) {
		return false
	}

	return true
}

func (self *WrapperType) ToString() string {
	return self.Type.ToString() + "<" + self.WrappedType.ToString() + ">"
}

type NotAnyType struct{}

var NAType = &NotAnyType{}

// Type
func (self *NotAnyType) PVT() ast.PrimitiveValueType {
	return ast.Invalid
}

func (self *NotAnyType) GetType() ValueType {
	return NA
}

func (self *NotAnyType) _eq(_ Type) bool {
	return false
}

func (self *NotAnyType) ToString() string {
	return "NotAnyType"
}

func TypeEquals(t Type, other Type) bool {
	// Adds support for manipulation between Enums and numbers, there might be complications with this however
	// tIsEnum := t.PVT() == ast.Enum
	// otherTIsEnum := other.PVT() == ast.Enum

	// if !(tIsEnum && otherTIsEnum) && !(!tIsEnum && !otherTIsEnum) {
	// 	var underlyingType Type
	// 	var notEnum Type
	// 	if tIsEnum {
	// 		underlyingType = t.(*EnumType).UnderlyingType
	// 		notEnum = other
	// 	}else {
	// 		underlyingType = other.(*EnumType).UnderlyingType
	// 		notEnum = t
	// 	}

	// 	return TypeEquals(underlyingType,notEnum)
	// }
	tpvt := t.PVT()
	otherpvt := other.PVT()

	if tpvt == ast.Unknown {
		return true
	} else if otherpvt == ast.Unknown {
		return true
	}

	if tpvt == ast.Unresolved || tpvt == ast.Unknown {
		return true
	} else if otherpvt == ast.Unresolved {
		return true
	}

	if t.GetType() != other.GetType() {
		return false
	}

	return t._eq(other)
}

var InvalidType = NewBasicType(ast.Invalid)
