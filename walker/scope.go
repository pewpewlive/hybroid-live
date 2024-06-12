package walker

import (
	"hybroid/ast"
	"hybroid/helpers"
)

type Context struct {
	Node  ast.Node
	Value Value
	Ret   Types
}

func NewEnvironment(path string) EnvironmentVal {
	scope := Scope{
		Tag:             &UntaggedTag{},
		Variables:       map[string]VariableVal{},
		VariableIndexes: map[string]int{},
	}
	global := EnvironmentVal{
		Scope:        scope,
		StructTypes:  map[string]*StructVal{},
	}

	global.Scope.Environment = &global
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
	Loop
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

// to be used
type EntityTag struct {
	//EntityType *StructTypeVal
}

func (et *EntityTag) GetType() ScopeTagType {
	return Entity
}

type LoopTag struct {
	Exits map[ExitType][]bool
}

func NewLoopTag(attrs ...ScopeAttribute) *LoopTag {
	exits := map[ExitType][]bool{
		All: make([]bool, 0),
	}
	
	for _, v := range attrs {
		switch v {
		case YieldAllowing:
			exits[Yield] = make([]bool, 0)
		case BreakAllowing:
			exits[Break] = make([]bool, 0)
		case ContinueAllowing:
			exits[Continue] = make([]bool, 0)
		case ReturnAllowing:
			exits[Return] = make([]bool, 0)
		}
	}
	return &LoopTag{
		Exits:exits,
	}
}

func (lt *LoopTag) GetType() ScopeTagType {
	return Loop
}

func (lt *LoopTag) SetExit(state bool, typ ExitType) {
	if _, found := lt.Exits[typ]; !found {
		return
	}
	lt.Exits[typ] = append(lt.Exits[typ], state)
}

func (lt *LoopTag) GetIfExits(et ExitType) bool {
	if _, found := lt.Exits[et]; !found {
		return false
	}
	exits := lt.Exits[et]

	for _, v := range exits {
		if v {
			return true
		}
	}

	return false
}

type FuncTag struct {
	Returns    []bool
	ReturnType Types
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
	if et != Return {
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
	mpt         *MultiPathTag
	YieldValues *Types
}

func (met *MatchExprTag) GetType() ScopeTagType {
	return MatchExpr
}

func (met *MatchExprTag) SetExit(state bool, typ ExitType) {
	met.mpt.SetExit(state, typ)
}

func (self *MatchExprTag) GetIfExits(et ExitType) bool {
	return self.mpt.GetIfExits(et)
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
	if et == All {
		exitTimes := 0
		for k := range mpt.Exits {
			if k == All {
				continue
			}

			for _, v := range mpt.Exits[k] {
				if v {
					exitTimes++
				}
			} 
		}
		return exitTimes == len(mpt.Exits[All])
	}
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
	Environment *EnvironmentVal
	Parent      *Scope

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
		Environment: parent.Environment,
		Parent:      parent,

		Tag:        tag,
		Attributes: attrs,

		Variables:       map[string]VariableVal{},
		VariableIndexes: map[string]int{},
	}
}
