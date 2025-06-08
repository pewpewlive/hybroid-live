package walker

import "hybroid/ast"

var TableEnv = &Environment{
	Name: "Table",
	Scope: Scope{
		Variables: map[string]*VariableVal{
			"Concat": {
				Name: "Concat",
				Value: NewFunction(NewWrapperType(NewBasicType(ast.List), NewBasicType(ast.Text)), NewBasicType(ast.Text), NewBasicType(ast.Number), NewBasicType(ast.Number)).
					WithReturns(NewBasicType(ast.Text)),
				IsPub:   true,
				IsConst: true,
			},
			"Insert": {
				Name:    "Insert",
				Value:   NewFunction(NewWrapperType(NewBasicType(ast.List), NewGeneric("T")), NewGeneric("T")).WithGenerics(NewGeneric("T")),
				IsPub:   true,
				IsConst: true,
			},
			"InsertAt": {
				Name:    "InsertAt",
				Value:   NewFunction(NewWrapperType(NewBasicType(ast.List), NewGeneric("T")), NewBasicType(ast.Number), NewGeneric("T")).WithGenerics(NewGeneric("T")),
				IsPub:   true,
				IsConst: true,
			},
			"Remove": {
				Name:    "Remove",
				Value:   NewFunction(NewWrapperType(NewBasicType(ast.List), NewGeneric("T")), NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Sort": {
				Name:    "Sort",
				Value:   NewFunction(NewWrapperType(NewBasicType(ast.List), NewGeneric("T"))),
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
