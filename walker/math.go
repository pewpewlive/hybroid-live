package walker

import "hybroid/ast"

var MathEnv = &Environment{
	Name: "Math",
	Scope: Scope{
		Variables: MathVariables,
		Tag:       &UntaggedTag{},
	},
	UsedWalkers:   make([]*Walker, 0),
	UsedLibraries: make(map[Library]bool),
	Structs:       make(map[string]*ClassVal),
	Entities:      make(map[string]*EntityVal),
	CustomTypes:   make(map[string]*CustomType),
	AliasTypes:    make(map[string]*AliasType),
}

var MathVariables = map[string]*VariableVal{
	"Pi": {
		Name:    "Pi",
		Value:   &NumberVal{},
	},
	"Huge": {
		Name:    "Huge",
		Value:   &NumberVal{},
	},
	"MaxInt": {
		Name:    "MaxInt",
		Value:   &NumberVal{},
	},
	"MinInt": {
		Name:    "MinInt",
		Value:   &NumberVal{},
	},

	"Abs": {
		Name:    "Abs",
		Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
	},
	"Acos": {
		Name:    "Acos",
		Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
	},
	"Atan": {
		Name:    "Atan",
		Value:   NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
	},
	"Ceil": {
		Name:    "Ceil",
		Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
	},
	"Floor": {
		Name:    "Floor",
		Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
	},
	"Cos": {
		Name:    "Cos",
		Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
	},
	"Sin": {
		Name:    "Sin",
		Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
	},
	"Sincos": {
		Name:    "Sincos",
		Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number), NewBasicType(ast.Number)),
	},
	"Deg": {
		Name:    "Deg",
		Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
	},
	"Rad": {
		Name:    "Rad",
		Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
	},
	"Exp": {
		Name:    "Exp",
		Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
	},
	"ToInt": {
		Name:    "ToInt",
		Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
	},
	"Fmod": {
		Name:    "Fmod",
		Value:   NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
	},
	"Ult": {
		Name:    "Ult",
		Value:   NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Bool)),
	},
	"Log": {
		Name:    "Log",
		Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.String)),
	},
	"Max": {
		Name:    "Max",
		Value:   NewFunction(NewBasicType(ast.Number), NewVariadicType(NewBasicType(ast.Number))).WithReturns(NewBasicType(ast.Number)),
	},
	"Min": {
		Name:    "Min",
		Value:   NewFunction(NewBasicType(ast.Number), NewVariadicType(NewBasicType(ast.Number))).WithReturns(NewBasicType(ast.Number)),
	},
	"Modf": {
		Name:    "Modf",
		Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number), NewBasicType(ast.Number)),
	},
	"Random": {
		Name:    "Random",
		Value:   NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
	},
	"Sqrt": {
		Name:    "Sqrt",
		Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
	},
	"Tan": {
		Name:    "Tan",
		Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
	},
	"Type": {
		Name:    "Type",
		Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.String)),
	},
}
