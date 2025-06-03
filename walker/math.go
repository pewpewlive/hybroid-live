package walker

import "hybroid/ast"

var MathEnv = &Environment{
	Name: "Math",
	Scope: Scope{
		Variables:  MathVariables,
		Tag:        &UntaggedTag{},
		AliasTypes: make(map[string]*AliasType),
	},
	importedWalkers: make([]*Walker, 0),
	UsedLibraries:   make([]Library, 0),
	Classes:         make(map[string]*ClassVal),
	Entities:        make(map[string]*EntityVal),
	Enums:           map[string]*EnumVal{},
}

var MathVariables = map[string]*VariableVal{
	"Pi": {
		Name:  "Pi",
		Value: &NumberVal{},
		IsPub: true,
	},
	"Huge": {
		Name:  "Huge",
		Value: &NumberVal{},
		IsPub: true,
	},
	"MaxInt": {
		Name:  "MaxInt",
		Value: &NumberVal{},
		IsPub: true,
	},
	"MinInt": {
		Name:  "MinInt",
		Value: &NumberVal{},
		IsPub: true,
	},

	"Abs": {
		Name:  "Abs",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsPub: true,
	},
	"Acos": {
		Name:  "Acos",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsPub: true,
	},
	"Atan": {
		Name:  "Atan",
		Value: NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsPub: true,
	},
	"Ceil": {
		Name:  "Ceil",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsPub: true,
	},
	"Floor": {
		Name:  "Floor",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsPub: true,
	},
	"Cos": {
		Name:  "Cos",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsPub: true,
	},
	"Sin": {
		Name:  "Sin",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsPub: true,
	},
	"Sincos": {
		Name:  "Sincos",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number), NewBasicType(ast.Number)),
		IsPub: true,
	},
	"Deg": {
		Name:  "Deg",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsPub: true,
	},
	"Rad": {
		Name:  "Rad",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsPub: true,
	},
	"Exp": {
		Name:  "Exp",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsPub: true,
	},
	"ToInt": {
		Name:  "ToInt",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsPub: true,
	},
	"Fmod": {
		Name:  "Fmod",
		Value: NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsPub: true,
	},
	"Ult": {
		Name:  "Ult",
		Value: NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Bool)),
		IsPub: true,
	},
	"Log": {
		Name:  "Log",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.String)),
		IsPub: true,
	},
	"Max": {
		Name:  "Max",
		Value: NewFunction(NewBasicType(ast.Number), NewVariadicType(NewBasicType(ast.Number))).WithReturns(NewBasicType(ast.Number)),
		IsPub: true,
	},
	"Min": {
		Name:  "Min",
		Value: NewFunction(NewBasicType(ast.Number), NewVariadicType(NewBasicType(ast.Number))).WithReturns(NewBasicType(ast.Number)),
		IsPub: true,
	},
	"Modf": {
		Name:  "Modf",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number), NewBasicType(ast.Number)),
		IsPub: true,
	},
	"Random": {
		Name:  "Random",
		Value: NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsPub: true,
	},
	"Sqrt": {
		Name:  "Sqrt",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsPub: true,
	},
	"Tan": {
		Name:  "Tan",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsPub: true,
	},
	"Type": {
		Name:  "Type",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.String)),
		IsPub: true,
	},
}
