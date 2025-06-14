package walker

import "hybroid/ast"

var StringAPI = &Environment{
	Name: "String",
	Scope: Scope{
		Variables: map[string]*VariableVal{
			"Byte": {
				Name:    "Byte",
				Value:   NewFunction(NewBasicType(ast.Text), NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Char": {
				Name:    "Char",
				Value:   NewFunction(NewWrapperType(NewBasicType(ast.List), NewBasicType(ast.Number))).WithReturns(NewBasicType(ast.Text)),
				IsPub:   true,
				IsConst: true,
			},
			"Find": {
				Name: "Find",
				Value: NewFunction(NewBasicType(ast.Text), NewBasicType(ast.Text), NewBasicType(ast.Number)).
					WithReturns(NewBasicType(ast.Number), NewBasicType(ast.Number), NewWrapperType(NewBasicType(ast.List), NewBasicType(ast.Number))),
				IsPub:   true,
				IsConst: true,
			},
			"Format": {
				Name:    "Format",
				Value:   NewFunction(NewBasicType(ast.Text), NewVariadicType(NewBasicType(ast.Object))).WithReturns(NewBasicType(ast.Text)),
				IsPub:   true,
				IsConst: true,
			},
			"Gsub": {
				Name:    "Gsub",
				Value:   NewFunction(NewBasicType(ast.Text), NewBasicType(ast.Text), NewBasicType(ast.Text)).WithReturns(NewBasicType(ast.Text)),
				IsPub:   true,
				IsConst: true,
			},
			"Gmatch": {
				Name:    "Gmatch",
				Value:   NewFunction(NewBasicType(ast.Text), NewBasicType(ast.Text)).WithReturns(NewBasicType(ast.Number), NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Dump": {
				Name:    "Dump",
				Value:   NewFunction(NewFunctionType([]Type{}, []Type{}), NewBasicType(ast.Bool)).WithReturns(NewBasicType(ast.Text)),
				IsPub:   true,
				IsConst: true,
			},
			"Len": {
				Name:    "Len",
				Value:   NewFunction(NewBasicType(ast.Text)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Lower": {
				Name:    "Lower",
				Value:   NewFunction(NewBasicType(ast.Text)).WithReturns(NewBasicType(ast.Text)),
				IsPub:   true,
				IsConst: true,
			},
			"Upper": {
				Name:    "Upper",
				Value:   NewFunction(NewBasicType(ast.Text)).WithReturns(NewBasicType(ast.Text)),
				IsPub:   true,
				IsConst: true,
			},
			"Match": {
				Name: "Match",
				Value: NewFunction(NewBasicType(ast.Text), NewBasicType(ast.Text), NewBasicType(ast.Number)).
					WithReturns(NewVariadicType(NewBasicType(ast.Text))),
				IsPub:   true,
				IsConst: true,
			},
			"Rep": {
				Name:    "Rep",
				Value:   NewFunction(NewBasicType(ast.Text), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Text)),
				IsPub:   true,
				IsConst: true,
			},
			"Reverse": {
				Name:    "Reverse",
				Value:   NewFunction(NewBasicType(ast.Text)).WithReturns(NewBasicType(ast.Text)),
				IsPub:   true,
				IsConst: true,
			},
			"Sub": {
				Name:    "Sub",
				Value:   NewFunction(NewBasicType(ast.Text), NewBasicType(ast.Number), NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Text)),
				IsPub:   true,
				IsConst: true,
			},
			"Pack": {
				Name:    "Pack",
				Value:   NewFunction(NewBasicType(ast.Text), NewBasicType(ast.Text), NewBasicType(ast.Text), NewVariadicType(NewBasicType(ast.Text))).WithReturns(NewBasicType(ast.Text)),
				IsPub:   true,
				IsConst: true,
			},
			"PackSize": {
				Name:    "PackSize",
				Value:   NewFunction(NewBasicType(ast.Text)).WithReturns(NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
			"Unpack": {
				Name: "Unpack",
				Value: NewFunction(NewBasicType(ast.Text), NewBasicType(ast.Text), NewBasicType(ast.Number)).
					WithReturns(NewVariadicType(NewBasicType(ast.Text)), NewBasicType(ast.Number)),
				IsPub:   true,
				IsConst: true,
			},
		},
		Tag:         &UntaggedTag{},
		AliasTypes:  make(map[string]*AliasType),
		ConstValues: make(map[string]ast.Node),
	},
	importedWalkers: make([]*Walker, 0),
	UsedLibraries:   make([]ast.Library, 0),
	Classes:         make(map[string]*ClassVal),
	Entities:        make(map[string]*EntityVal),
	Enums:           make(map[string]*EnumVal),
}

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
