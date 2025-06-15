package walker

import "hybroid/ast"

var NumberListVal = &ListVal{ValueType: NewBasicType(ast.Number)}
var VertexListVal = &ListVal{ValueType: NumberListVal.GetType()}
var MeshType = NewStructType([]StructField{
	NewStructField("vertexes", VertexListVal),
	NewStructField("segments", VertexListVal),
	NewStructField("colors", NumberListVal, true),
})
var MeshesType = (&ListVal{ValueType: MeshType}).GetType()
var SoundType = NewStructType([]StructField{
	NewStructField("attack", &NumberVal{}, true),
	NewStructField("decay", &NumberVal{}, true),
	NewStructField("sustain", &NumberVal{}, true),
	NewStructField("sustainPunch", &NumberVal{}, true),
	NewStructField("amplification", &NumberVal{}, true),
	NewStructField("harmonics", &NumberVal{}, true),
	NewStructField("harmonicsFalloff", &NumberVal{}, true),
	NewStructField("tremoloDepth", &NumberVal{}, true),
	NewStructField("tremoloFrequency", &NumberVal{}, true),
	NewStructField("frequency", &NumberVal{}, true),
	NewStructField("frequencyDeltaSweep", &NumberVal{}, true),
	NewStructField("frequencyJump1Onset", &NumberVal{}, true),
	NewStructField("frequencyJump2Onset", &NumberVal{}, true),
	NewStructField("frequencyJump1Amount", &NumberVal{}, true),
	NewStructField("frequencyJump2Amount", &NumberVal{}, true),
	NewStructField("frequencySweep", &NumberVal{}, true),
	NewStructField("vibratoFrequency", &NumberVal{}, true),
	NewStructField("vibratoDepth", &NumberVal{}, true),
	NewStructField("flangerOffset", &NumberVal{}, true),
	NewStructField("flangerOffsetSweep", &NumberVal{}, true),
	NewStructField("repeatFrequency", &NumberVal{}, true),
	NewStructField("lowPassCutoff", &NumberVal{}, true),
	NewStructField("lowPassCutoffSweep", &NumberVal{}, true),
	NewStructField("highPassCutoff", &NumberVal{}, true),
	NewStructField("highPassCutoffSweep", &NumberVal{}, true),
	NewStructField("bitCrush", &NumberVal{}, true),
	NewStructField("bitCrushSweep", &NumberVal{}, true),
	NewStructField("squareDuty", &NumberVal{}, true),
	NewStructField("squareDutySweep", &NumberVal{}, true),
	NewStructField("harmonicsFalloff", &NumberVal{}, true),
	NewStructField("normalization", &BoolVal{}, true),
	NewStructField("interpolateNoise", &BoolVal{}, true),
	NewStructField("compression", &NumberVal{}, true),
	NewStructField("harmonics", &NumberVal{}, true),
	NewStructField("harmonicsFalloff", &NumberVal{}, true),
	NewStructField("repeatFrequency", &NumberVal{}, true),
	NewStructField("sampleRate", &NumberVal{}, true),
	NewStructField("waveform", &StringVal{}),
})
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
			"Center": NewAliasType("Center", NewStructType([]StructField{
				NewStructField("x", &NumberVal{}),
				NewStructField("y", &NumberVal{}),
				NewStructField("z", &NumberVal{}),
			}), false),
			"Sound": NewAliasType("Sound", SoundType, false),
		},
		ConstValues: make(map[string]ast.Node),
	},
	importedWalkers: make([]*Walker, 0),
	UsedLibraries:   make([]ast.Library, 0),
	Classes:         make(map[string]*ClassVal),
	Entities:        make(map[string]*EntityVal),
	Enums:           make(map[string]*EnumVal),
}

var BuiltinLibraries = map[ast.Library]*Environment{
	ast.Pewpew: PewpewAPI,
	ast.Fmath:  FmathAPI,
	ast.Math:   MathAPI,
	ast.String: StringAPI,
	ast.Table:  TableAPI,
}
