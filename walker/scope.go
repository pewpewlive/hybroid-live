package walker

import (
	"hybroid/ast"
	"hybroid/helpers"
)

type Context struct {
	Node  ast.Node
	Value Value
	Ret   ReturnType
}

type Namespace struct {
	Ctx          Context
	Scope        Scope
	foreignTypes map[string]Value
	StructTypes  map[string]*StructTypeVal
}

func NewNamespace() Namespace {
	scope := Scope{
		Tag:             UntaggedTag{},
		Variables:       map[string]VariableVal{},
		VariableIndexes: map[string]int{},
	}
	global := Namespace{
		Ctx: Context{
			Node:  ast.Improper{},
			Value: Unknown{},
			Ret:   ReturnType{},
		},
		Scope:        scope,
		foreignTypes: map[string]Value{},
		StructTypes:  map[string]*StructTypeVal{},
	}

	global.Scope.Namespace = &global
	return global
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
)

type ReturnableTag interface {
	SetReturn(state bool, types ...ExitType) ScopeTag
}

type ScopeTag interface {
	GetType() ScopeTagType
}

type UntaggedTag struct{}

func (ut UntaggedTag) GetType() ScopeTagType {
	return Untagged
}

type StructTag struct {
	StructType *StructTypeVal
}

func (st StructTag) GetType() ScopeTagType {
	return Struct
}

// to be used
type EntityTag struct {
	//EntityType *StructTypeVal
}

func (et EntityTag) GetType() ScopeTagType {
	return Entity
}

type FuncTag struct {
	Returns    []bool
	ReturnType ReturnType
}

func (et FuncTag) GetType() ScopeTagType {
	return Func
}

func (et FuncTag) SetReturn(state bool, types ...ExitType) ScopeTag {
	et.Returns = append(et.Returns, state)
	return et
}

type MatchExprTag struct {
	mpt         MultiPathTag
	ArmsYielded int
	YieldValues *ReturnType
}

func (met MatchExprTag) GetType() ScopeTagType {
	return MatchExpr
}

func (met MatchExprTag) SetReturn(state bool, types ...ExitType) ScopeTag {
	if state {
		for _, v := range types {
			if v == Yield {
				met.ArmsYielded++
			} else {
				met.mpt.SetReturn(state, types...)
			}
		}
	}
	return met
}

type MultiPathTag struct {
	Returns   []bool
	Yields    []bool
	Continues []bool
	Breaks    []bool
}

func (mpt MultiPathTag) GetType() ScopeTagType {
	return MultiPath
}

func (mpt MultiPathTag) SetReturn(state bool, types ...ExitType) ScopeTag {
	if state {
		for _, v := range types {
			if v == Yield {
				mpt.Yields = append(mpt.Yields, state)
			} else if v == Return {
				mpt.Returns = append(mpt.Returns, state)
			} else if v == Continue {
				mpt.Continues = append(mpt.Continues, state)
			} else if v == Break {
				mpt.Breaks = append(mpt.Breaks, state)
			}
		}
	}
	return mpt
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

func NewScopeAttributes(types ...ScopeAttribute) ScopeAttributes {
	return types
}

func (sa *ScopeAttributes) Add(_type ScopeAttribute) {
	for i := range *sa {
		if (*sa)[i] == _type {
			return
		}
	}
	*sa = append(*sa, _type)
}

var EmptyAttributes = ScopeAttributes{}

type Scope struct {
	Namespace *Namespace
	Parent    *Scope

	Tag        ScopeTag
	Attributes ScopeAttributes

	Variables       map[string]VariableVal
	VariableIndexes map[string]int
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

func NewScope(parent *Scope, tag ScopeTag) Scope {
	var attrs ScopeAttributes
	if parent == nil {
		attrs = EmptyAttributes
	} else {
		attrs = parent.Attributes
	}
	return Scope{
		Namespace: parent.Namespace,
		Parent:    parent,

		Tag:        tag,
		Attributes: attrs,

		Variables:       map[string]VariableVal{},
		VariableIndexes: map[string]int{},
	}
}
