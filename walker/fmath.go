package walker

import (
	"hybroid/ast"
)

var FmathEnv = &Environment{
	Name: "Fmath",
	Scope: Scope{
		Variables:  FmathVariables,
		Tag:        &UntaggedTag{},
		AliasTypes: make(map[string]*AliasType),
	},
	importedWalkers: make([]*Walker, 0),
	UsedLibraries:   make([]Library, 0),
	Classes:         make(map[string]*ClassVal),
	Entities:        make(map[string]*EntityVal),
}

var FmathVariables = map[string]*VariableVal{
	"MaxFixed": {
		Name:  "MaxFixed",
		Value: NewFunction().WithReturns(NewFixedPointType()),
		IsPub: true,
	},
	"RandomFixed": {
		Name:  "RandomFixed",
		Value: NewFunction(NewFixedPointType(), NewFixedPointType()).WithReturns(NewFixedPointType()),
		IsPub: true,
	},
	"RandomNum": {
		Name:  "RandomNum",
		Value: NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsPub: true,
	},
	"Sqrt": {
		Name:  "Sqrt",
		Value: NewFunction(NewFixedPointType()).WithReturns(NewFixedPointType()),
		IsPub: true,
	},
	"FromFraction": {
		Name:  "FromFraction",
		Value: NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewFixedPointType()),
		IsPub: true,
	},
	"ToNum": {
		Name:  "ToNum",
		Value: NewFunction(NewFixedPointType()).WithReturns(NewBasicType(ast.Number)),
		IsPub: true,
	},
	"AbsFixed": {
		Name:  "AbsFixed",
		Value: NewFunction(NewFixedPointType()).WithReturns(NewFixedPointType()),
		IsPub: true,
	},
	"ToFixed": {
		Name:  "ToFixed",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewFixedPointType()),
		IsPub: true,
	},
	"Sincos": {
		Name:  "Sincos",
		Value: NewFunction(NewFixedPointType()).WithReturns(NewFixedPointType(), NewFixedPointType()),
		IsPub: true,
	},
	"Atan2": {
		Name:  "Atan2",
		Value: NewFunction(NewFixedPointType(), NewFixedPointType()).WithReturns(NewFixedPointType()),
		IsPub: true,
	},
	"Tau": {
		Name:  "Tau",
		Value: NewFunction().WithReturns(NewFixedPointType()),
		IsPub: true,
	},
}
