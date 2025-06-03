package walker

import "hybroid/ast"

var TableEnv = &Environment{
	Name: "Table",
	Scope: Scope{
		Variables:  TableVariables,
		Tag:        &UntaggedTag{},
		AliasTypes: make(map[string]*AliasType),
	},
	importedWalkers: make([]*Walker, 0),
	UsedLibraries:   make([]Library, 0),
	Classes:         make(map[string]*ClassVal),
	Entities:        make(map[string]*EntityVal),
	Enums:           map[string]*EnumVal{},
}

var TableVariables = map[string]*VariableVal{
	"Concat": {
		Name: "Concat",
		Value: NewFunction(NewWrapperType(NewBasicType(ast.List), NewBasicType(ast.String)), NewBasicType(ast.String), NewBasicType(ast.Number), NewBasicType(ast.Number)).
			WithReturns(NewBasicType(ast.String)),
		IsPub: true,
	},
	"Insert": {
		Name:  "Insert",
		Value: NewFunction(NewWrapperType(NewBasicType(ast.List), NewGeneric("T")), NewGeneric("T")).WithGenerics(NewGeneric("T")),
		IsPub: true,
	},
	"InsertAt": {
		Name:  "InsertAt",
		Value: NewFunction(NewWrapperType(NewBasicType(ast.List), NewGeneric("T")), NewBasicType(ast.Number), NewGeneric("T")).WithGenerics(NewGeneric("T")),
		IsPub: true,
	},
	"Remove": {
		Name:  "Remove",
		Value: NewFunction(NewWrapperType(NewBasicType(ast.List), NewGeneric("T")), NewBasicType(ast.Number)),
		IsPub: true,
	},
	"Sort": {
		Name:  "Sort",
		Value: NewFunction(NewWrapperType(NewBasicType(ast.List), NewGeneric("T"))),
		IsPub: true,
	},
} // Table.Insert(list, 9)
