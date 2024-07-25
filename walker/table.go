package walker

var TableEnv = &Environment{
	Name: "Table",
	Scope: Scope{
		Variables: tableVariables,
		Tag: &UntaggedTag{},
	},
	Structs: make(map[string]*StructVal),
	Entities: make(map[string]*EntityVal),
	CustomTypes: make(map[string]*CustomType),
}

var tableVariables = map[string]*VariableVal{
	
}