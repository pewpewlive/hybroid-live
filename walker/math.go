package walker

import "hybroid/ast"

var MathEnv = &Environment{
	Name: "Math",
	Scope: Scope{
		Variables: MathVariables,
		Tag: &UntaggedTag{},
	},
	UsedWalkers: make([]*Walker, 0),
	UsedLibraries: make(map[Library]bool),
	Structs: make(map[string]*StructVal),
	Entities: make(map[string]*EntityVal),
	CustomTypes: make(map[string]*CustomType),
}

var MathVariables = map[string]*VariableVal{
	"Pi": {
		Name:  "Pi",
		Value: &NumberVal{},
		IsConst: true,
	},
	"Huge": {
		Name:  "Huge",
		Value: &NumberVal{},
		IsConst: true,
	},
	"MaxInt": {
		Name:  "MaxInt",
		Value: &NumberVal{},
		IsConst: true,
	},
	"MinInt": {
		Name:  "MinInt",
		Value: &NumberVal{},
		IsConst: true,
	},

	"Abs": {
		Name:  "Abs",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"Acos": {
		Name:  "Acos",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"Atan": {
		Name:  "Atan",
		Value: NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"Ceil": {
		Name:  "Ceil",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"Floor": {
		Name:  "Floor",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"Cos": {
		Name:  "Cos",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"Sin": {
		Name:  "Sin",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"Sincos": {
		Name:  "Sincos",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number), NewBasicType(ast.Number)),
		IsConst: true,
	},
	"Deg": {
		Name:  "Deg",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"Rad": {
		Name:  "Rad",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"Exp": {
		Name:  "Exp",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"ToInt": {
		Name:  "ToInt",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"Fmod": {
		Name:  "Fmod",
		Value: NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"Ult": {
		Name:  "Ult",
		Value: NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Bool)),
		IsConst: true,
	},
	"Log": {
		Name:  "Log",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.String)),
		IsConst: true,
	},
	"Max": {
		Name:  "Max",
		Value: NewFunction(NewBasicType(ast.Number), NewVariadicType(NewBasicType(ast.Number))).WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"Min": {
		Name:  "Min",
		Value: NewFunction(NewBasicType(ast.Number), NewVariadicType(NewBasicType(ast.Number))).WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"Modf": {
		Name:  "Modf",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number), NewBasicType(ast.Number)),
		IsConst: true,
	},
	"Random": {
		Name:  "Random",
		Value: NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"Sqrt": {
		Name: "Sqrt",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"Tan": {
		Name: "Tan",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"Type": {
		Name: "Type",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.String)),
		IsConst: true,
	},
}