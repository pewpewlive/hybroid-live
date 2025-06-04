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

type Context struct {
	Node        ast.Node
	Value       Value
	EntityCasts core.Queue[EntityCast]
}

func (c *Context) Clear() {
	c.Node = &ast.Improper{}
	c.Value = &Unknown{}
	c.EntityCasts.Clear()
}

type ScopeTagType int

const (
	Untagged ScopeTagType = iota
	Class
	Entity
	Func
	MultiPath
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
	Returns     []bool
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
	ft.Returns = append(ft.Returns, state)
}

func (ft *FuncTag) GetIfExits(et ExitType) bool {
	if et != Return && et != All {
		return false
	}
	if len(ft.Returns) == 0 {
		return false
	}

	for _, v := range ft.Returns {
		if v {
			return true
		}
	}

	return false
}

type MatchExprTag struct {
	Mpt        *MultiPathTag
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

type MultiPathTag struct {
	Exits map[ExitType][]bool
}

func (mpt *MultiPathTag) GetType() ScopeTagType {
	return MultiPath
}

func (mpt *MultiPathTag) SetExit(state bool, typ ExitType) {
	if _, found := mpt.Exits[typ]; !found {
		return
	}
	for i := range mpt.Exits[typ] {
		if !mpt.Exits[typ][i] {
			mpt.Exits[typ][i] = state
			break
		}
	}
}

func (mpt *MultiPathTag) GetIfExits(et ExitType) bool {
	exits := mpt.Exits[et]

	if len(exits) == 0 {
		return false
	}

	for _, v := range exits {
		if !v {
			return false
		}
	}

	return true
}

func NewMultiPathTag(requirement int, attrs ...ScopeAttribute) *MultiPathTag {
	exits := map[ExitType][]bool{
		All: make([]bool, requirement),
	}
	for _, v := range attrs {
		switch v {
		case YieldAllowing:
			exits[Yield] = make([]bool, requirement)
		case BreakAllowing:
			exits[Break] = make([]bool, requirement)
		case ContinueAllowing:
			exits[Continue] = make([]bool, requirement)
		case ReturnAllowing:
			exits[Return] = make([]bool, requirement)
		}
	}

	return &MultiPathTag{Exits: exits}
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
