package walker

import "hybroid/ast"

var TableEnv = &Environment{
	Name: "Table",
	Scope: Scope{
		Variables: tableVariables,
		Tag:       &UntaggedTag{},
	},
	Structs:     make(map[string]*StructVal),
	Entities:    make(map[string]*EntityVal),
	CustomTypes: make(map[string]*CustomType),
}

var tableVariables = map[string]*VariableVal{
	"Concat": {
		Name: "Concat",
		Value: NewFunction(NewWrapperType(NewBasicType(ast.List), NewBasicType(ast.String)), NewBasicType(ast.String), NewBasicType(ast.Number), NewBasicType(ast.Number)).
			WithReturns(NewBasicType(ast.String)),
		IsLocal: false,
		IsConst: true,
	},
}