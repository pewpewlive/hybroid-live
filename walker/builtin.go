package walker

import "hybroid/ast"

var BuiltinEnv = &Environment{
	Name: "Builtin",
	Scope: Scope{
		Variables: BuiltinVariables,
		Tag:       &UntaggedTag{},
	},
	importedWalkers: make([]*Walker, 0),
	UsedLibraries:   make(map[Library]bool),
	Structs:         make(map[string]*ClassVal),
	Entities:        make(map[string]*EntityVal),
	CustomTypes:     make(map[string]*CustomType),
	AliasTypes: map[string]*AliasType{
		"Mesh":     NewAliasType("Mesh", MeshValueType),
		"Meshes":   NewAliasType("Meshes", MeshesValueType),
		"Vertex":   NewAliasType("Vertex", numberListVal.GetType()),
		"Vertexes": NewAliasType("Vertexes", vertexesVal.GetType()),
		"Segments": NewAliasType("Segments", vertexesVal.GetType()),
		"Segment":  NewAliasType("Segment", numberListVal.GetType()),
		"Colors":   NewAliasType("Segments", numberListVal.GetType()),
		"Center": NewAliasType("Center", NewStructType([]*VariableVal{
			{
				Name:  "x",
				Value: &NumberVal{},
			},
			{
				Name:  "y",
				Value: &NumberVal{},
			},
			{
				Name:  "z",
				Value: &NumberVal{},
			},
		}, false)),
		"Sound": NewAliasType("Sound", SoundValueType),
	},
}

var BuiltinVariables = map[string]*VariableVal{
	"ToString": {
		Name:    "ToString",
		Value:   NewFunction(NewBasicType(ast.Object)),
		IsUsed:  false,
		IsConst: true,
	},
	"ParseSound": {
		Name: "ParseSound",
		Value: NewFunction(NewBasicType(ast.String)).
			WithReturns(SoundValueType),
		IsConst: true,
	},
}
