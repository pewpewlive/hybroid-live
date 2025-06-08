package walker

import "hybroid/ast"

var MathEnv = &Environment{
	Name: "Math",
	Scope: Scope{
		Variables: map[string]*VariableVal{
			"Pi": {
				Name:    "Pi",
				Value:   &NumberVal{},
				IsPub:   true,
				IsConst: true,
			},
			"Huge": {
				Name:    "Huge",
				Value:   &NumberVal{},
				IsPub:   true,
				IsConst: true,
			},
			"MaxInt": {
				Name:    "MaxInt",
				Value:   &NumberVal{},
				IsPub:   true,
				IsConst: true,
			},
			"MinInt": {
				Name:    "MinInt",
				Value:   &NumberVal{},
				IsPub:   true,
				IsConst: true,
			},

			"Abs": {
				Name:    "Abs",
				Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Acos": {
				Name:    "Acos",
				Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Atan": {
				Name:    "Atan",
				Value:   NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Ceil": {
				Name:    "Ceil",
				Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Floor": {
				Name:    "Floor",
				Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Cos": {
				Name:    "Cos",
				Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Sin": {
				Name:    "Sin",
				Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Sincos": {
				Name:    "Sincos",
				Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number), NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Deg": {
				Name:    "Deg",
				Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Rad": {
				Name:    "Rad",
				Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Exp": {
				Name:    "Exp",
				Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"ToInt": {
				Name:    "ToInt",
				Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Fmod": {
				Name:    "Fmod",
				Value:   NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Ult": {
				Name:    "Ult",
				Value:   NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Bool)),
				IsPub:   true,
				IsConst: true,
			},
			"Log": {
				Name:    "Log",
				Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Text)),
				IsPub:   true,
				IsConst: true,
			},
			"Max": {
				Name:    "Max",
				Value:   NewFunction(NewBasicType(ast.Number), NewVariadicType(NewBasicType(ast.Number))).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Min": {
				Name:    "Min",
				Value:   NewFunction(NewBasicType(ast.Number), NewVariadicType(NewBasicType(ast.Number))).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Modf": {
				Name:    "Modf",
				Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number), NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Random": {
				Name:    "Random",
				Value:   NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Sqrt": {
				Name:    "Sqrt",
				Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Tan": {
				Name:    "Tan",
				Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Type": {
				Name:    "Type",
				Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Text)),
				IsPub:   true,
				IsConst: true,
			},
		},
		Tag:         &UntaggedTag{},
		AliasTypes:  make(map[string]*AliasType),
		ConstValues: make(map[string]ast.Node),
	},
	importedWalkers: make([]*Walker, 0),
	UsedLibraries:   make([]Library, 0),
	Classes:         make(map[string]*ClassVal),
	Entities:        make(map[string]*EntityVal),
	Enums:           make(map[string]*EnumVal),
}
