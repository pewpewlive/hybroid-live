package walker

import "hybroid/ast"

var MathAPI = &Environment{
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
				Value:   NewFunction([]string{"x"}, NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Acos": {
				Name:    "Acos",
				Value:   NewFunction([]string{"x"}, NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Atan": {
				Name:    "Atan",
				Value:   NewFunction([]string{"y", "x"}, NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Ceil": {
				Name:    "Ceil",
				Value:   NewFunction([]string{"x"}, NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Floor": {
				Name:    "Floor",
				Value:   NewFunction([]string{"x"}, NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Cos": {
				Name:    "Cos",
				Value:   NewFunction([]string{"x"}, NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Sin": {
				Name:    "Sin",
				Value:   NewFunction([]string{"x"}, NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Sincos": {
				Name:    "Sincos",
				Value:   NewFunction([]string{"x"}, NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number), NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Deg": {
				Name:    "Deg",
				Value:   NewFunction([]string{"x"}, NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Rad": {
				Name:    "Rad",
				Value:   NewFunction([]string{"x"}, NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Exp": {
				Name:    "Exp",
				Value:   NewFunction([]string{"x"}, NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"ToInt": {
				Name:    "ToInt",
				Value:   NewFunction([]string{"x"}, NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Fmod": {
				Name:    "Fmod",
				Value:   NewFunction([]string{"x", "y"}, NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Ult": {
				Name:    "Ult",
				Value:   NewFunction([]string{"x", "y"}, NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Bool)),
				IsPub:   true,
				IsConst: true,
			},
			"Log": {
				Name:    "Log",
				Value:   NewFunction([]string{"x"}, NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Text)),
				IsPub:   true,
				IsConst: true,
			},
			"Max": {
				Name:    "Max",
				Value:   NewFunction([]string{"x", "rest"}, NewBasicType(ast.Number), NewVariadicType(NewBasicType(ast.Number))).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Min": {
				Name:    "Min",
				Value:   NewFunction([]string{"x", "rest"}, NewBasicType(ast.Number), NewVariadicType(NewBasicType(ast.Number))).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Modf": {
				Name:    "Modf",
				Value:   NewFunction([]string{"x"}, NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number), NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Random": {
				Name:    "Random",
				Value:   NewFunction([]string{"m", "n"}, NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Sqrt": {
				Name:    "Sqrt",
				Value:   NewFunction([]string{"x"}, NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Tan": {
				Name:    "Tan",
				Value:   NewFunction([]string{"x"}, NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Type": {
				Name:    "Type",
				Value:   NewFunction([]string{"x"}, NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Text)),
				IsPub:   true,
				IsConst: true,
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
