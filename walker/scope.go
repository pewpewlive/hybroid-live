package walker

import (
	"hybroid/ast"
	"hybroid/core"
	"hybroid/tokens"
	"slices"
)

type EntityCast struct {
	Name   tokens.Token
	Entity *EntityVal
}

func NewEntityCast(name tokens.Token, val *EntityVal) EntityCast {
	return EntityCast{
		Name:   name,
		Entity: val,
	}
}

type Binding struct {
	Scope   *Scope
	VarName string
}

type Context struct {
	EntityCasts   core.Queue[EntityCast]
	DontSetToUsed bool
}

func (c *Context) Clear() {
	c.DontSetToUsed = false
	c.EntityCasts.Clear()
}

type ScopeTagType int

const (
	Untagged ScopeTagType = iota
	Class
	Entity
	Func
	NormalPath
	MatchExpr
)

type ExitType int

const (
	Yield ExitType = iota
	Return
	Continue
	Break
	All
)

type ExitableTag interface {
	ScopeTag
	SetExit(state bool, _type ExitType)
	GetIfExits(typ ExitType) bool
}

type ScopeTag interface {
	GetType() ScopeTagType
}

type UntaggedTag struct{}

func (ut *UntaggedTag) GetType() ScopeTagType {
	return Untagged
}

type ClassTag struct {
	Val *ClassVal
}

func (st *ClassTag) GetType() ScopeTagType {
	return Class
}

type EntityTag struct {
	EntityVal *EntityVal
}

func (et *EntityTag) GetType() ScopeTagType {
	return Entity
}

type FuncTag struct {
	Return      bool
	ReturnTypes []Type

	Generics []*GenericType
}

func (ft *FuncTag) GetType() ScopeTagType {
	return Func
}

func (ft *FuncTag) SetExit(state bool, etype ExitType) {
	if etype != Return && etype != All {
		return
	}
	if state {
		ft.Return = true
	}
}

func (ft *FuncTag) GetIfExits(et ExitType) bool {
	return ft.Return
}

type MatchExprTag struct {
	Mpt        *PathTag
	YieldTypes []Type
}

func (met *MatchExprTag) GetType() ScopeTagType {
	return MatchExpr
}

func (met *MatchExprTag) SetExit(state bool, typ ExitType) {
	met.Mpt.SetExit(state, typ)
}

func (met *MatchExprTag) GetIfExits(et ExitType) bool {
	return met.Mpt.GetIfExits(et)
}

type PathTag struct {
	Exits map[ExitType]bool
}

func (mpt *PathTag) GetType() ScopeTagType {
	return NormalPath
}

func (mpt *PathTag) SetAllFalse() {
	for typ := range mpt.Exits {
		mpt.Exits[typ] = false
	}
}

func (mpt *PathTag) SetExit(state bool, typ ExitType) {
	if _, found := mpt.Exits[typ]; !found {
		return
	}
	if !mpt.Exits[typ] {
		mpt.Exits[typ] = state
	}
}

func (mpt *PathTag) SetAllExitAND(other *PathTag) {
	for typ := range mpt.Exits {
		if _, found := other.Exits[typ]; !found {
			continue
		}
		mpt.Exits[typ] = mpt.Exits[typ] && other.Exits[typ]
	}
}

func (mpt *PathTag) GetIfExits(et ExitType) bool {
	return mpt.Exits[et]
}

func NewPathTag(attrs ...ScopeAttribute) *PathTag {
	exits := map[ExitType]bool{
		All: false,
	}
	for _, v := range attrs {
		switch v {
		case YieldAllowing:
			exits[Yield] = false
		case BreakAllowing:
			exits[Break] = false
		case ContinueAllowing:
			exits[Continue] = false
		case ReturnAllowing:
			exits[Return] = false
		}
	}

	return &PathTag{Exits: exits}
}

type ScopeAttribute int

const (
	ReturnAllowing ScopeAttribute = iota + 1
	YieldAllowing
	SelfAllowing
	BreakAllowing
	ContinueAllowing
)

type ScopeAttributes []ScopeAttribute

func (sa *ScopeAttributes) Add(attribute ScopeAttribute) {
	for i := range *sa {
		if (*sa)[i] == attribute {
			return
		}
	}
	*sa = append(*sa, attribute)
}

func (sa *ScopeAttributes) Remove(attribute ScopeAttribute) {
	for i := range *sa {
		if (*sa)[i] == attribute {
			*sa = append((*sa)[:i], (*sa)[i+1:]...)
			return
		}
	}
}

var EmptyAttributes = ScopeAttributes{}

type Scope struct {
	Environment *Environment
	Parent      *Scope

	Tag        ScopeTag
	Attributes ScopeAttributes

	Variables   map[string]*VariableVal
	AliasTypes  map[string]*AliasType
	ConstValues map[string]ast.Node

	Body *[]*ast.Node
}

func (sc *Scope) resolveAlias(typeName string) (*AliasType, bool) {
	if alias, found := sc.AliasTypes[typeName]; found {
		return alias, true
	}

	if sc.Parent == nil {
		return nil, false
	}

	return sc.Parent.resolveAlias(typeName)
}

func (sc *Scope) Is(types ...ScopeAttribute) bool {
	if len(types) == 0 {
		return false
	}

	for _, v := range types {
		if !slices.Contains(sc.Attributes, v) {
			return false
		}
	}

	return true
}

func NewScope(parent *Scope, tag ScopeTag, extraAttrs ...ScopeAttribute) *Scope {
	var attrs ScopeAttributes
	if parent == nil {
		attrs = EmptyAttributes
	} else {
		attrs = append(attrs, parent.Attributes...)
	}
	for _, v := range extraAttrs {
		attrs.Add(v)
	}
	scope := Scope{
		Environment: parent.Environment,
		Parent:      parent,

		Tag:        tag,
		Attributes: attrs,

		Variables:   map[string]*VariableVal{},
		AliasTypes:  map[string]*AliasType{},
		ConstValues: make(map[string]ast.Node),
	}
	return &scope
}
