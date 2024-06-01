package walker

import (
	"hybroid/ast"
	"reflect"
)

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
		Tag:             UntaggedTag{},
		Variables:       map[string]VariableVal{},
		VariableIndexes: map[string]int{},
	}
	global := Global{
		Ctx: Context{
			Node:  ast.Improper{},
			Value: Unknown{},
			Ret:   ReturnType{},
		},
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
	Loop
)

type GetType int

const (
	YIELD GetType = iota
	RETURN
	CONTINUE
	BREAK
)

type ReturnableTag interface {
	SetReturn(state bool, types ...GetType) ScopeTag
}

func GetValOfInterface[T any, E any](val E) *T {
	value := reflect.ValueOf(val)
	ah := reflect.TypeFor[T]()
	if value.CanConvert(ah) {
		test := value.Convert(ah).Interface()
		tVal := test.(T)
		return &tVal
	}

	return nil
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

func (et FuncTag) SetReturn(state bool, types ...GetType) ScopeTag {
	et.Returns = append(et.Returns, state)
	return et
}

type MatchTag struct {
	ReturnAmount int
	ArmsYielded  int
	YieldValues  *ReturnType
}

func (et MatchTag) GetType() ScopeTagType {
	return MatchExpr
}

func (et MatchTag) SetReturn(state bool, types ...GetType) ScopeTag {
	if state {
		for _, v := range types {
			if v == YIELD {
				et.ArmsYielded++
			} else {
				et.ReturnAmount++
			}
		}
	}
	return et
}

type MultiPathTag struct {
	ReturnAmount   int
	YieldAmount    int
	ContinueAmount int
	BreakAmount    int
}

func (mp MultiPathTag) GetType() ScopeTagType {
	return MultiPath
}

func (et MultiPathTag) SetReturn(state bool, types ...GetType) ScopeTag {
	if state {
		for _, v := range types {
			if v == YIELD {
				et.YieldAmount++
			} else if v == RETURN {
				et.ReturnAmount++
			} else if v == CONTINUE {
				et.ContinueAmount++
			} else if v == BREAK {
				et.BreakAmount++
			}
		}
	}
	return et
}

type LoopTag struct {
	Continues []bool
	Breaks    []bool
	Returns   []bool
	Yields    []bool
}

func (lt LoopTag) GetType() ScopeTagType {
	return Loop
}

func (lt LoopTag) SetReturn(state bool, types ...GetType) ScopeTag {
	for _, v := range types {
		if v == YIELD {
			lt.Yields = append(lt.Yields, state)
		} else if v == RETURN {
			lt.Returns = append(lt.Returns, state)
		} else if v == CONTINUE {
			lt.Continues = append(lt.Continues, state)
		} else if v == BREAK {
			lt.Breaks = append(lt.Breaks, state)
		}
	}
	return lt
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
