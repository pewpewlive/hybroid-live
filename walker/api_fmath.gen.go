// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
package walker

import "hybroid/ast"

// AUTO-GENERATED API, DO NOT MANUALLY MODIFY!
var FmathAPI = &Environment{
	Name: "Fmath",
	Scope: Scope{
		Variables: map[string]*VariableVal{
			"MaxFixed": {
				Name: "MaxFixed", Value: NewFunction().WithReturns(NewFixedPointType()), IsPub: true,
			},
			"RandomFixed": {
				Name: "RandomFixed", Value: NewFunction(NewFixedPointType(), NewFixedPointType()).WithReturns(NewFixedPointType()), IsPub: true,
			},
			"RandomNumber": {
				Name: "RandomNumber", Value: NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)), IsPub: true,
			},
			"Sqrt": {
				Name: "Sqrt", Value: NewFunction(NewFixedPointType()).WithReturns(NewFixedPointType()), IsPub: true,
			},
			"FromFraction": {
				Name: "FromFraction", Value: NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewFixedPointType()), IsPub: true,
			},
			"ToNumber": {
				Name: "ToNumber", Value: NewFunction(NewFixedPointType()).WithReturns(NewBasicType(ast.Number)), IsPub: true,
			},
			"AbsFixed": {
				Name: "AbsFixed", Value: NewFunction(NewFixedPointType()).WithReturns(NewFixedPointType()), IsPub: true,
			},
			"ToFixed": {
				Name: "ToFixed", Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewFixedPointType()), IsPub: true,
			},
			"Sincos": {
				Name: "Sincos", Value: NewFunction(NewFixedPointType()).WithReturns(NewFixedPointType(), NewFixedPointType()), IsPub: true,
			},
			"Atan2": {
				Name: "Atan2", Value: NewFunction(NewFixedPointType(), NewFixedPointType()).WithReturns(NewFixedPointType()), IsPub: true,
			},
			"Tau": {
				Name: "Tau", Value: NewFunction().WithReturns(NewFixedPointType()), IsPub: true,
			},
		},
		Tag:         &UntaggedTag{},
		AliasTypes:  make(map[string]*AliasType),
		ConstValues: make(map[string]ast.Node),
	},
	imports:       make([]Import, 0),
	UsedLibraries: make([]ast.Library, 0),
	Classes:       make(map[string]*ClassVal),
	Entities:      make(map[string]*EntityVal),
	Enums:         make(map[string]*EnumVal),
}
