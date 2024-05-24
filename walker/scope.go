package walker

import "go/ast"

type Context struct {
	Node  ast.Node
	Value Value
	Ret   ReturnType
}

type Global struct {
	Ctx          Context
	Scope        Scope
	foreignTypes map[string]Value
	StructTypes  map[string]*StructTypeVal
}

func NewGlobal() Global {
	scope := Scope{
		Variables:       map[string]VariableVal{},
		VariableIndexes: map[string]int{},
	}
	global := Global{
		Scope:        scope,
		foreignTypes: map[string]Value{},
		StructTypes:  map[string]*StructTypeVal{},
	}

	global.Scope.Global = &global
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

type EntityTag struct {
	//EntityType *StructTypeVal
}

func (et EntityTag) GetType() ScopeTagType {
	return Entity
}

type FuncTag struct {
	Returns    bool
	ReturnType ReturnType
}

func (et FuncTag) GetType() ScopeTagType {
	return Func
}

type MatchTag struct {
	ArmsYielded int
	YieldValues *ReturnType
}

func (et MatchTag) GetType() ScopeTagType {
	return MatchExpr
}

type MultiPathTag struct{}

func (mp MultiPathTag) GetType() ScopeTagType {
	return MultiPath
}

type ScopeAttribute int

const (
	ReturnAllowing ScopeAttribute = iota + 1
	YieldAllowing
	SelfAllowing
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
	Global *Global
	Parent *Scope

	Tag        ScopeTag
	Attributes ScopeAttributes

	Variables       map[string]VariableVal
	VariableIndexes map[string]int
}

func Contains(list []ScopeAttribute, thing ScopeAttribute) bool {
	for _, v := range list {
		if thing == v {
			return true
		}
	}
	return false
}

func (sc *Scope) Is(types ...ScopeAttribute) bool {
	if len(types) == 0 {
		return false
	}

	for _, v := range types {
		if !Contains(sc.Attributes, v) {
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
		Global: parent.Global,
		Parent: parent,

		Tag:        tag,
		Attributes: attrs,

		Variables:       map[string]VariableVal{},
		VariableIndexes: map[string]int{},
	}
}
