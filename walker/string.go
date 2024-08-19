package walker

import "hybroid/ast"

var StringEnv = &Environment{
	Name: "String",
	Scope: Scope{
		Variables: StringVariables,
		Tag:       &UntaggedTag{},
	},
	UsedWalkers:   make([]*Walker, 0),
	UsedLibraries: make(map[Library]bool),
	Structs:       make(map[string]*ClassVal),
	Entities:      make(map[string]*EntityVal),
	CustomTypes:   make(map[string]*CustomType),
}

var StringVariables = map[string]*VariableVal{
	"Byte": {
		Name:    "Byte",
		Value:   NewFunction(NewBasicType(ast.String), NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"Char": {
		Name:    "Char",
		Value:   NewFunction(NewWrapperType(NewBasicType(ast.List), NewBasicType(ast.Number))).WithReturns(NewBasicType(ast.String)),
		IsConst: true,
	},
	"Find": {
		Name: "Find",
		Value: NewFunction(NewBasicType(ast.String), NewBasicType(ast.String), NewBasicType(ast.Number)).
			WithReturns(NewBasicType(ast.Number), NewBasicType(ast.Number), NewWrapperType(NewBasicType(ast.List), NewBasicType(ast.Number))),
		IsConst: true,
	},
	"Format": {
		Name:    "Format",
		Value:   NewFunction(NewBasicType(ast.String), NewVariadicType(NewBasicType(ast.Object))).WithReturns(NewBasicType(ast.String)),
		IsConst: true,
	},
	"Gsub": {
		Name:    "Gsub",
		Value:   NewFunction(NewBasicType(ast.String), NewBasicType(ast.String), NewBasicType(ast.String)).WithReturns(NewBasicType(ast.String)),
		IsConst: true,
	},
	"Gmatch": {
		Name:    "Gmatch",
		Value:   NewFunction(NewBasicType(ast.String), NewBasicType(ast.String)).WithReturns(NewBasicType(ast.Number), NewBasicType(ast.Number)),
		IsConst: true,
	},
	"Dump": {
		Name:    "Dump",
		Value:   NewFunction(NewFunctionType(Types{}, Types{}), NewBasicType(ast.Bool)).WithReturns(NewBasicType(ast.String)),
		IsConst: true,
	},
	"Len": {
		Name:    "Len",
		Value:   NewFunction(NewBasicType(ast.String)).WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"Lower": {
		Name:    "Lower",
		Value:   NewFunction(NewBasicType(ast.String)).WithReturns(NewBasicType(ast.String)),
		IsConst: true,
	},
	"Upper": {
		Name:    "Upper",
		Value:   NewFunction(NewBasicType(ast.String)).WithReturns(NewBasicType(ast.String)),
		IsConst: true,
	},
	"Match": {
		Name: "Match",
		Value: NewFunction(NewBasicType(ast.String), NewBasicType(ast.String), NewBasicType(ast.Number)).
			WithReturns(NewVariadicType(NewBasicType(ast.String))),
		IsConst: true,
	},
	"Rep": {
		Name:    "Rep",
		Value:   NewFunction(NewBasicType(ast.String), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.String)),
		IsConst: true,
	},
	"Reverse": {
		Name:    "Reverse",
		Value:   NewFunction(NewBasicType(ast.String)).WithReturns(NewBasicType(ast.String)),
		IsConst: true,
	},
	"Sub": {
		Name:    "Sub",
		Value:   NewFunction(NewBasicType(ast.String), NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.String)),
		IsConst: true,
	},
	"Pack": {
		Name:    "Pack",
		Value:   NewFunction(NewBasicType(ast.String), NewBasicType(ast.String), NewBasicType(ast.String), NewVariadicType(NewBasicType(ast.String))).WithReturns(NewBasicType(ast.String)),
		IsConst: true,
	},
	"PackSize": {
		Name:    "PackSize",
		Value:   NewFunction(NewBasicType(ast.String)).WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"Unpack": {
		Name: "Unpack",
		Value: NewFunction(NewBasicType(ast.String), NewBasicType(ast.String), NewBasicType(ast.Number)).
			WithReturns(NewVariadicType(NewBasicType(ast.String)), NewBasicType(ast.Number)),
		IsConst: true,
	},
	// byte [done]
	// char [tobeverified]
	// dump [tobeverified]
	// find [tobeverified]
	// format [tobeverified]
	// gmatch [tobeverified]
	// gsub [tobeverified]
	// len [tobeverified]
	// lower [done]
	// match [tobeverified]
	// rep [tobeverified]
	// reverse [tobeverified]
	// sub [tobeverified]
	// upper [done]
	// pack [done]
	// packsize [done]
	// unpack [tobeverified]

}
