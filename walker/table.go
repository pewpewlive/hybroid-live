package walker

import "hybroid/ast"

var TableEnv = &Environment{
	Name: "Table",
	Scope: Scope{
		Variables: TableVariables,
		Tag:       &UntaggedTag{},
	},
	UsedWalkers:   make([]*Walker, 0),
	UsedLibraries: make(map[Library]bool),
	Structs:       make(map[string]*ClassVal),
	Entities:      make(map[string]*EntityVal),
	CustomTypes:   make(map[string]*CustomType),
}

var TableVariables = map[string]*VariableVal{
	"Concat": {
		Name: "Concat",
		Value: NewFunction(NewWrapperType(NewBasicType(ast.List), NewBasicType(ast.String)), NewBasicType(ast.String), NewBasicType(ast.Number), NewBasicType(ast.Number)).
			WithReturns(NewBasicType(ast.String)),
		IsConst: true,
	},
	"Insert": {
		Name:    "Insert",
		Value:   NewFunction(NewWrapperType(NewBasicType(ast.List), NewGeneric("T")), NewGeneric("T")).WithGenerics(NewGeneric("T")),
		IsConst: true,
	},
	"InsertAt": {
		Name:    "InsertAt",
		Value:   NewFunction(NewWrapperType(NewBasicType(ast.List), NewGeneric("T")), NewBasicType(ast.Number), NewGeneric("T")).WithGenerics(NewGeneric("T")),
		IsConst: true,
	},
	"Remove": {
		Name:    "Remove",
		Value:   NewFunction(NewWrapperType(NewBasicType(ast.List), NewGeneric("T")), NewBasicType(ast.Number)),
		IsConst: true,
	},
	"Sort": {
		Name:    "Sort",
		Value:   NewFunction(NewWrapperType(NewBasicType(ast.List), NewGeneric("T"))),
		IsConst: true,
	},
} // Table.Insert(list, 9)
