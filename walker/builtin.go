package walker

import "hybroid/ast"

var BuiltinEnv = &Environment{
	Name: "Builtin",
	Scope: Scope{
		Variables: BuiltinVariables,
		Tag:       &UntaggedTag{},
		AliasTypes: map[string]*AliasType{
			"Mesh":     NewAliasType("Mesh", MeshValueType, false),
			"Meshes":   NewAliasType("Meshes", MeshesValueType, false),
			"Vertex":   NewAliasType("Vertex", numberListVal.GetType(), false),
			"Vertexes": NewAliasType("Vertexes", vertexesVal.GetType(), false),
			"Segments": NewAliasType("Segments", vertexesVal.GetType(), false),
			"Segment":  NewAliasType("Segment", numberListVal.GetType(), false),
			"Colors":   NewAliasType("Segments", numberListVal.GetType(), false),
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
			}, false), false),
			"Sound": NewAliasType("Sound", SoundValueType, false),
		},
	},
	importedWalkers: make([]*Walker, 0),
	UsedLibraries:   make([]Library, 0),
	Classes:         make(map[string]*ClassVal),
	Entities:        make(map[string]*EntityVal),
}

var BuiltinVariables = map[string]*VariableVal{
	"ToString": {
		Name:    "ToString",
		Value:   NewFunction(NewBasicType(ast.Object)),
		IsUsed:  false,
		IsConst: true,
		IsPub:   true,
	},
	"ParseSound": {
		Name: "ParseSound",
		Value: NewFunction(NewBasicType(ast.String)).
			WithReturns(SoundValueType),
		IsConst: true,
		IsPub:   true,
	},
}
