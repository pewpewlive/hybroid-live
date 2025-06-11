package walker

import "hybroid/ast"

var NumberListVal = &ListVal{ValueType: NewBasicType(ast.Number)}
var VertexListVal = &ListVal{ValueType: NumberListVal.GetType()}
var MeshType = NewStructType([]*VariableVal{
	{Name: "vertexes", Value: VertexListVal},
	{Name: "segments", Value: VertexListVal},
	{Name: "colors", Value: NumberListVal},
}, false)
var MeshesType = (&ListVal{ValueType: MeshType}).GetType()
var SoundType = NewStructType([]*VariableVal{
	{Name: "attack", Value: &NumberVal{}},
	{Name: "decay", Value: &NumberVal{}},
	{Name: "sustain", Value: &NumberVal{}},
	{Name: "sustainPunch", Value: &NumberVal{}},
	{Name: "amplification", Value: &NumberVal{}},
	{Name: "harmonics", Value: &NumberVal{}},
	{Name: "harmonicsFalloff", Value: &NumberVal{}},
	{Name: "tremoloDepth", Value: &NumberVal{}},
	{Name: "tremoloFrequency", Value: &NumberVal{}},
	{Name: "frequency", Value: &NumberVal{}},
	{Name: "frequencyDeltaSweep", Value: &NumberVal{}},
	{Name: "frequencyJump1Onset", Value: &NumberVal{}},
	{Name: "frequencyJump2Onset", Value: &NumberVal{}},
	{Name: "frequencyJump1Amount", Value: &NumberVal{}},
	{Name: "frequencyJump2Amount", Value: &NumberVal{}},
	{Name: "frequencySweep", Value: &NumberVal{}},
	{Name: "vibratoFrequency", Value: &NumberVal{}},
	{Name: "vibratoDepth", Value: &NumberVal{}},
	{Name: "flangerOffset", Value: &NumberVal{}},
	{Name: "flangerOffsetSweep", Value: &NumberVal{}},
	{Name: "repeatFrequency", Value: &NumberVal{}},
	{Name: "lowPassCutoff", Value: &NumberVal{}},
	{Name: "lowPassCutoffSweep", Value: &NumberVal{}},
	{Name: "highPassCutoff", Value: &NumberVal{}},
	{Name: "highPassCutoffSweep", Value: &NumberVal{}},
	{Name: "bitCrush", Value: &NumberVal{}},
	{Name: "bitCrushSweep", Value: &NumberVal{}},
	{Name: "squareDuty", Value: &NumberVal{}},
	{Name: "squareDutySweep", Value: &NumberVal{}},
	{Name: "harmonicsFalloff", Value: &NumberVal{}},
	{Name: "normalization", Value: &BoolVal{}},
	{Name: "interpolateNoise", Value: &BoolVal{}},
	{Name: "compression", Value: &NumberVal{}},
	{Name: "harmonics", Value: &NumberVal{}},
	{Name: "harmonicsFalloff", Value: &NumberVal{}},
	{Name: "repeatFrequency", Value: &NumberVal{}},
	{Name: "sampleRate", Value: &NumberVal{}},
	{Name: "waveform", Value: &StringVal{}},
}, true)
var SoundsType = (&ListVal{ValueType: SoundType}).GetType()
var WeaponCollisionSign = NewFuncSignature().
	WithParams(NewBasicType(ast.Number), NewEnumType("Pewpew", "WeaponType")).
	WithReturns(NewBasicType(ast.Bool))
var PlayerCollisionSign = NewFuncSignature().
	WithParams(NewBasicType(ast.Number), &RawEntityType{})
var WallCollisionSign = NewFuncSignature().
	WithParams(NewFixedPointType(), NewFixedPointType())

var BuiltinEnv = &Environment{
	Name: "Builtin",
	Scope: Scope{
		Variables: map[string]*VariableVal{
			"ToString": {
				Name:    "ToString",
				Value:   NewFunction(NewBasicType(ast.Object)).WithReturns(NewBasicType(ast.Text)),
				IsUsed:  false,
				IsConst: true,
				IsPub:   true,
			},
			"ParseSound": {
				Name:    "ParseSound",
				Value:   NewFunction(NewBasicType(ast.Text)).WithReturns(SoundType),
				IsConst: true,
				IsPub:   true,
			},
		},
		Tag: &UntaggedTag{},
		AliasTypes: map[string]*AliasType{
			"Mesh":     NewAliasType("Mesh", MeshType, false),
			"Meshes":   NewAliasType("Meshes", MeshesType, false),
			"Vertex":   NewAliasType("Vertex", NumberListVal.GetType(), false),
			"Vertexes": NewAliasType("Vertexes", VertexListVal.GetType(), false),
			"Segments": NewAliasType("Segments", VertexListVal.GetType(), false),
			"Segment":  NewAliasType("Segment", NumberListVal.GetType(), false),
			"Colors":   NewAliasType("Segments", NumberListVal.GetType(), false),
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
			"Sound": NewAliasType("Sound", SoundType, false),
		},
		ConstValues: make(map[string]ast.Node),
	},
	importedWalkers: make([]*Walker, 0),
	UsedLibraries:   make([]Library, 0),
	Classes:         make(map[string]*ClassVal),
	Entities:        make(map[string]*EntityVal),
	Enums:           make(map[string]*EnumVal),
}

var BuiltinLibraries = map[Library]*Environment{
	Pewpew: PewpewAPI,
	Fmath:  FmathAPI,
	Math:   MathAPI,
	String: StringAPI,
	Table:  TableAPI,
}
