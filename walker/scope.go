package walker

import (
	"hybroid/ast"
	"hybroid/helpers"
)

type Context struct {
	Node  ast.Node
	Value Value
	PewpewVarFound bool
	PewpewVarName string
} 

func (c *Context) Clear() {
	c.Node = &ast.Improper{}
	c.Value = &Unknown{}
	c.PewpewVarFound = false
}

type ScopeTagType int

const (
	Untagged ScopeTagType = iota
	Struct
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

type StructTag struct {
	StructVal *StructVal
}

func (st *StructTag) GetType() ScopeTagType {
	return Struct
}

type EntityTag struct {
	EntityType *EntityVal
}

func (et *EntityTag) GetType() ScopeTagType {
	return Entity
}

type FuncTag struct {
	Returns     []bool
	ReturnTypes Types

	Generics    []*GenericType
}

func (et *FuncTag) GetType() ScopeTagType {
	return Func
}

func (et *FuncTag) SetExit(state bool, etype ExitType) {
	if etype != Return && etype != All {
		return
	}
	et.Returns = append(et.Returns, state)
}

func (self *FuncTag) GetIfExits(et ExitType) bool {
	if et != Return && et != All {
		return false
	}
	if len(self.Returns) == 0 {
		return false
	}

	for _, v := range self.Returns {
		if v {
			return true
		}
	}

	return false
}

type MatchExprTag struct {
	Mpt         *MultiPathTag
	YieldValues Types
}

func (met *MatchExprTag) GetType() ScopeTagType {
	return MatchExpr
}

func (met *MatchExprTag) SetExit(state bool, typ ExitType) {
	met.Mpt.SetExit(state, typ)
}

func (self *MatchExprTag) GetIfExits(et ExitType) bool {
	return self.Mpt.GetIfExits(et)
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

var EmptyAttributes = ScopeAttributes{}

type Scope struct {
	Environment *Environment
	Parent      *Scope

	Tag        ScopeTag
	Attributes ScopeAttributes

	Variables         map[string]*VariableVal

	Body *[]*ast.Node
}

func (sc *Scope) Is(types ...ScopeAttribute) bool {
	if len(types) == 0 {
		return false
	}

	for _, v := range types {
		if !helpers.Contains(sc.Attributes, v) {
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
		attrs = parent.Attributes
	}
	for _, v := range extraAttrs {
		attrs.Add(v)
	}
	scope := Scope{
		Environment: parent.Environment,
		Parent:      parent,

		Tag:        tag,
		Attributes: attrs,

		Variables: map[string]*VariableVal{},
	}
	return &scope
}
