package walker

import "hybroid/ast"

var TableAPI = &Environment{
	Name: "Table",
	Scope: Scope{
		Variables: map[string]*VariableVal{
			"Concat": {
				Name: "Concat",
				Value: NewFunction([]string{"list", "sep", "i", "j"}, NewWrapperType(NewBasicType(ast.List), NewBasicType(ast.Text)), NewBasicType(ast.Text), NewBasicType(ast.Number), NewBasicType(ast.Number)).
					WithReturns(NewBasicType(ast.Text)),
				IsPub:   true,
				IsConst: true,
			},
			"Insert": {
				Name:    "Insert",
				Value:   NewFunction([]string{"list", "value"}, NewWrapperType(NewBasicType(ast.List), NewGeneric("T")), NewGeneric("T")).WithGenerics(NewGeneric("T")),
				IsPub:   true,
				IsConst: true,
			},
			"InsertAt": {
				Name:    "InsertAt",
				Value:   NewFunction([]string{"list", "pos", "value"}, NewWrapperType(NewBasicType(ast.List), NewGeneric("T")), NewBasicType(ast.Number), NewGeneric("T")).WithGenerics(NewGeneric("T")),
				IsPub:   true,
				IsConst: true,
			},
			"Remove": {
				Name:    "Remove",
				Value:   NewFunction([]string{"list", "pos"}, NewWrapperType(NewBasicType(ast.List), NewGeneric("T")), NewBasicType(ast.Number)).WithReturns(NewGeneric("T")).WithGenerics(NewGeneric("T")),
				IsPub:   true,
				IsConst: true,
			},
			"Sort": {
				Name:    "Sort",
				Value:   NewFunction([]string{"list"}, NewWrapperType(NewBasicType(ast.List), NewGeneric("T"))).WithGenerics(NewGeneric("T")),
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
