package walker

import "hybroid/ast"

var BuiltinEnv = &Environment{
	Name: "Builtin",
	Scope: Scope{
		Variables: BuiltinVariables,
		Tag:       &UntaggedTag{},
	},
	UsedWalkers:   make([]*Walker, 0),
	UsedLibraries: make(map[Library]bool),
	Structs:       make(map[string]*ClassVal),
	Entities:      make(map[string]*EntityVal),
	CustomTypes:   make(map[string]*CustomType),
}

var BuiltinVariables = map[string]*VariableVal{
	"ToString": {
		Name:  "ToString",
		Value: NewFunction(NewBasicType(ast.Object)),
		IsUsed: false,
		IsConst: true,
	},
	"ParseSound": {
		Name: "ParseSound",
		Value: NewFunction(NewBasicType(ast.String)).
			WithReturns(SoundValueType),
		IsConst: true,
	},
}