package walker

import (
	"hybroid/ast"
)

var FmathEnv = &Environment{
	Name: "Fmath",
	Scope: Scope{
		Variables: FmathVariables,
		Tag:       &UntaggedTag{},
	},
	UsedWalkers:   make([]*Walker, 0),
	UsedLibraries: make(map[Library]bool),
	Structs:       make(map[string]*ClassVal),
	Entities:      make(map[string]*EntityVal),
	CustomTypes:   make(map[string]*CustomType),
	AliasTypes:    make(map[string]*AliasType),
}

var FmathVariables = map[string]*VariableVal{
	"MaxFixed": {
		Name:    "MaxFixed",
		Value:   NewFunction().WithReturns(NewFixedPointType(ast.Fixed)),
	},
	"RandomFixed": {
		Name:    "RandomFixed",
		Value:   NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)).WithReturns(NewFixedPointType(ast.Fixed)),
	},
	"RandomNum": {
		Name:    "RandomNum",
		Value:   NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
	},
	"Sqrt": {
		Name:    "Sqrt",
		Value:   NewFunction(NewFixedPointType(ast.Fixed)).WithReturns(NewFixedPointType(ast.Fixed)),
	},
	"FromFraction": {
		Name:    "FromFraction",
		Value:   NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewFixedPointType(ast.Fixed)),
	},
	"ToNum": {
		Name:    "ToNum",
		Value:   NewFunction(NewFixedPointType(ast.Fixed)).WithReturns(NewBasicType(ast.Number)),
	},
	"AbsFixed": {
		Name:    "AbsFixed",
		Value:   NewFunction(NewFixedPointType(ast.Fixed)).WithReturns(NewFixedPointType(ast.Fixed)),
	},
	"ToFixed": {
		Name:    "ToFixed",
		Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewFixedPointType(ast.Fixed)),
	},
	"Sincos": {
		Name:    "Sincos",
		Value:   NewFunction(NewFixedPointType(ast.Fixed)).WithReturns(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)),
	},
	"Atan2": {
		Name:    "Atan2",
		Value:   NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)).WithReturns(NewFixedPointType(ast.Fixed)),
	},
	"Tau": {
		Name:    "Tau",
		Value:   NewFunction().WithReturns(NewFixedPointType(ast.Fixed)),
	},
}
