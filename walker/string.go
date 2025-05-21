package walker

import "hybroid/ast"

var StringEnv = &Environment{
	Name: "String",
	Scope: Scope{
		Variables: StringVariables,
		Tag:       &UntaggedTag{},
	},
	importedWalkers: make([]*Walker, 0),
	UsedLibraries:   make(map[Library]bool),
	Structs:         make(map[string]*ClassVal),
	Entities:        make(map[string]*EntityVal),
	CustomTypes:     make(map[string]*CustomType),
	AliasTypes:      make(map[string]*AliasType),
}

var StringVariables = map[string]*VariableVal{
	"Byte": {
		Name:  "Byte",
		Value: NewFunction(NewBasicType(ast.String), NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
	},
	"Char": {
		Name:  "Char",
		Value: NewFunction(NewWrapperType(NewBasicType(ast.List), NewBasicType(ast.Number))).WithReturns(NewBasicType(ast.String)),
	},
	"Find": {
		Name: "Find",
		Value: NewFunction(NewBasicType(ast.String), NewBasicType(ast.String), NewBasicType(ast.Number)).
			WithReturns(NewBasicType(ast.Number), NewBasicType(ast.Number), NewWrapperType(NewBasicType(ast.List), NewBasicType(ast.Number))),
	},
	"Format": {
		Name:  "Format",
		Value: NewFunction(NewBasicType(ast.String), NewVariadicType(NewBasicType(ast.Object))).WithReturns(NewBasicType(ast.String)),
	},
	"Gsub": {
		Name:  "Gsub",
		Value: NewFunction(NewBasicType(ast.String), NewBasicType(ast.String), NewBasicType(ast.String)).WithReturns(NewBasicType(ast.String)),
	},
	"Gmatch": {
		Name:  "Gmatch",
		Value: NewFunction(NewBasicType(ast.String), NewBasicType(ast.String)).WithReturns(NewBasicType(ast.Number), NewBasicType(ast.Number)),
	},
	"Dump": {
		Name:  "Dump",
		Value: NewFunction(NewFunctionType(Types{}, Types{}), NewBasicType(ast.Bool)).WithReturns(NewBasicType(ast.String)),
	},
	"Len": {
		Name:  "Len",
		Value: NewFunction(NewBasicType(ast.String)).WithReturns(NewBasicType(ast.Number)),
	},
	"Lower": {
		Name:  "Lower",
		Value: NewFunction(NewBasicType(ast.String)).WithReturns(NewBasicType(ast.String)),
	},
	"Upper": {
		Name:  "Upper",
		Value: NewFunction(NewBasicType(ast.String)).WithReturns(NewBasicType(ast.String)),
	},
	"Match": {
		Name: "Match",
		Value: NewFunction(NewBasicType(ast.String), NewBasicType(ast.String), NewBasicType(ast.Number)).
			WithReturns(NewVariadicType(NewBasicType(ast.String))),
	},
	"Rep": {
		Name:  "Rep",
		Value: NewFunction(NewBasicType(ast.String), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.String)),
	},
	"Reverse": {
		Name:  "Reverse",
		Value: NewFunction(NewBasicType(ast.String)).WithReturns(NewBasicType(ast.String)),
	},
	"Sub": {
		Name:  "Sub",
		Value: NewFunction(NewBasicType(ast.String), NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.String)),
	},
	"Pack": {
		Name:  "Pack",
		Value: NewFunction(NewBasicType(ast.String), NewBasicType(ast.String), NewBasicType(ast.String), NewVariadicType(NewBasicType(ast.String))).WithReturns(NewBasicType(ast.String)),
	},
	"PackSize": {
		Name:  "PackSize",
		Value: NewFunction(NewBasicType(ast.String)).WithReturns(NewBasicType(ast.Number)),
	},
	"Unpack": {
		Name: "Unpack",
		Value: NewFunction(NewBasicType(ast.String), NewBasicType(ast.String), NewBasicType(ast.Number)).
			WithReturns(NewVariadicType(NewBasicType(ast.String)), NewBasicType(ast.Number)),
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
