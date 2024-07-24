package walker

import "hybroid/ast"

var StringEnv = &Environment{
	Name: "String",
	Scope: Scope{
		Variables: mathVariables,
		Tag: &UntaggedTag{},
	},
	Structs: make(map[string]*StructVal),
	Entities: make(map[string]*EntityVal),
	CustomTypes: make(map[string]*CustomType),
}

var stringVariables = map[string]*VariableVal{
	"Byte": {
		Name:  "Byte",
		Value: NewFunction(NewBasicType(ast.String), NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"Char": {
		Name:  "Char",
		Value: NewFunction(NewWrapperType(NewBasicType(ast.List), NewBasicType(ast.Number))).WithReturns(NewBasicType(ast.String)),
		IsConst: true,
	},
	"Find": {
		Name:  "Find",
		Value: NewFunction(NewBasicType(ast.String), NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number), NewBasicType(ast.Number)),
		IsConst: true,
	},
	// "Format": {
	// 	Name:  "Byte",
	// 	Value: NewFunction(NewWrapperType(NewBasicType(ast.List), NewBasicType(ast.String))),
	// 	IsConst: true,
	// },
	"Gsub": {
		Name:  "Gsub",
		Value: NewFunction(NewBasicType(ast.String), NewBasicType(ast.String), NewBasicType(ast.String)).WithReturns(NewBasicType(ast.String)),
		IsConst: true,
	},
	"Gmatch": {
		Name:  "Gmatch",
		Value: NewFunction(NewBasicType(ast.String), NewBasicType(ast.String)).WithReturns(NewBasicType(ast.Number), NewBasicType(ast.Number)),
		IsConst: true,
	},
	// "Find": {
	// 	Name:  "Find",
	// 	Value: NewFunction(NewBasicType(ast.String), NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number), NewBasicType(ast.Number)),
	// 	IsConst: true,
	// },
	// "Find": {
	// 	Name:  "Find",
	// 	Value: NewFunction(NewBasicType(ast.String), NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number), NewBasicType(ast.Number)),
	// 	IsConst: true,
	// },
	// "Find": {
	// 	Name:  "Find",
	// 	Value: NewFunction(NewBasicType(ast.String), NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number), NewBasicType(ast.Number)),
	// 	IsConst: true,
	// },
	// byte
	// char
	// dump // to do
	// find
	// format // to do 
	// gmatch
	// gsub
	// len
	// lower
	// match
	// rep
	// reverse
	// sub
	// upper
	// pack
	// packsize
	// unpack

}