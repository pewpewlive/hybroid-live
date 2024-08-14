package walker

import (
	"hybroid/ast"
)

var PewpewEnv = &Environment{
	Name: "Pewpew",
	Scope: Scope{
		Variables: PewpewVariables,
		Tag: &UntaggedTag{},
	},
	UsedWalkers: make([]*Walker, 0),
	UsedLibraries: make(map[Library]bool),
	Structs: make(map[string]*StructVal),
	Entities: make(map[string]*EntityVal),
	CustomTypes: make(map[string]*CustomType),
}

var PewpewVariables = map[string]*VariableVal{
	//enums
	"EntityType": {
		Name:    "EntityType",
		Value:   PewpewEntityType,
		IsLocal: false,
		IsConst: true,
	},
	"MothershipType": {
		Name:    "MothershipType",
		Value:   MothershipType,
		IsLocal: false,
		IsConst: true,
	},
	"CannonType": {
		Name:    "CannonType",
		Value:   CannonType,
		IsLocal: false,
		IsConst: true,
	},
	"CannonFreq": {
		Name:    "CannonFreq",
		Value:   CannonFrequency,
		IsLocal: false,
		IsConst: true,
	},
	"BombType": {
		Name:    "BombType",
		Value:   BombType,
		IsLocal: false,
		IsConst: true,
	},
	"BonusType": {
		Name:    "BonusType",
		Value:   BonusType,
		IsLocal: false,
		IsConst: true,
	},
	"WeaponType": {
		Name:    "WeaponType",
		Value:   WeaponType,
		IsLocal: false,
		IsConst: true,
	},
	"AsteroidSize": {
		Name:    "AsteroidSize",
		Value:   AsteroidSize,
		IsLocal: false,
		IsConst: true,
	},

	//functions
	"Print": {
		Name:  "Print",
		Value: NewFunction(NewBasicType(ast.Object)),
		IsConst: true,
	},
	"PrintDebugInfo": {
		Name:  "PrintDebugInfo",
		Value: NewFunction(),
		IsConst: true,
	},
	"SetLevelSize": {
		Name:  "SetLevelSize",
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)),
		IsConst: true,
	},
	"AddWall": {
		Name:  "AddWall",
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)).WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"RemoveWall": {
		Name:  "SetLevelSize",
		Value: NewFunction(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"AddUpdateCallback": {
		Name:  "AddUpdateCallback",
		Value: NewFunction(NewFunctionType(Types{}, Types{})),
		IsConst: true,
	},
	"GetNumberOfPlayers": {
		Name:  "GetNumberOfPlayers",
		Value: NewFunction().WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"IncreasePlayerScore": {
		Name:  "IncreasePlayerScore",
		Value: NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)),
		IsConst: true,
	},
	"IncreasePlayerScoreStreak": {
		Name:  "IncreasePlayerScoreStreak",
		Value: NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)),
		IsConst: true,
	},
	"GetPlayerScoreStreak": {
		Name:  "GetPlayerScoreStreak",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"StopGame": {
		Name:  "StopGame",
		Value: NewFunction(),
		IsConst: true,
	},
	"GetPlayerInputs": {
		Name:  "GetPlayerInputs",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)),
		IsConst: true,
	},
	"GetPlayerScore": {
		Name:  "GetPlayerScore",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"ConfigurePlayer": { 
		Name:  "ConfigurePlayer",
		Value: NewFunction(NewBasicType(ast.Number), NewStructType(map[string]Field{
			"has_lost": NewField(0, &VariableVal{
				Name: "has_lost",
				Value: &BoolVal{},
			}),
			"shield": NewField(1, &VariableVal{
				Name: "shield",
				Value: &NumberVal{},
			}),
			"camera_x_override": NewField(2, &VariableVal{
				Name: "camera_x_override",
				Value: &FixedVal{SpecificType: ast.Fixed},
			}),
			"camera_y_override": NewField(3, &VariableVal{
				Name: "camera_y_override",
				Value: &FixedVal{SpecificType: ast.Fixed},
			}),
			"camera_distance": NewField(4, &VariableVal{
				Name: "camera_distance",
				Value: &FixedVal{SpecificType: ast.Fixed},
			}),
			"camera_rotation_x_axis": NewField(5, &VariableVal{
				Name: "camera_rotation_x_axis",
				Value: &FixedVal{SpecificType: ast.Fixed},
			}),
			"move_joystick_color": NewField(6, &VariableVal{
				Name: "move_joystick_color",
				Value: &NumberVal{},
			}),
			"shoot_joystick_color": NewField(7, &VariableVal{
				Name: "shoot_joystick_color",
				Value: &NumberVal{},
			}),
		}, true)),
		IsConst: true,
	},
	"ConfigurePlayerHud": {
		Name: "ConfigurePlayerHud",
		Value: NewFunction(NewBasicType(ast.Number), NewStructType(map[string]Field{
			"TopLeftLine": NewField(0, &VariableVal{
				Name: "TopLeftLine",
				Value: &StringVal{},
			}),
		}, true)),
	},
	"GetPlayerConfig": {  
		Name:  "GetPlayerConfig", 
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewStructType(map[string]Field{
			"shield": NewField(0, &VariableVal{
				Name: "shield",
				Value: &NumberVal{},
			}),
			"has_lost": NewField(1, &VariableVal{
				Name: "has_lost",
				Value: &BoolVal{},
			}),
		}, false)),
		IsConst: true,
	}, 
	"ConfigureShipWeapon": {  
		Name:  "ConfigureShipWeapon",
		Value: NewFunction(&RawEntityType{}, NewStructType(map[string]Field{
			"frequency": NewField(0, &VariableVal{
				Name: "frequency",
				Value: CannonFrequency,
			}),
			"cannon": NewField(1, &VariableVal{
				Name: "cannon",
				Value: CannonType,
			}),
			"duration": NewField(2, &VariableVal{
				Name: "duration",
				Value: &NumberVal{},
			}),
		}, true)),
		IsConst: true,
	},
	"DamageShip": { 
		Name:  "DamageShip", 
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Number)),
		IsConst: true,
	},
	"AddArrowToShip": { 
		Name:  "AddArrowToShip", 
		Value: NewFunction(&RawEntityType{}, &RawEntityType{}, NewBasicType(ast.Number)).WithReturns(&RawEntityType{}), 
		IsConst: true,
	},
	"RemoveArrowFromShip": { 
		Name:  "RemoveArrowFromShip", 
		Value: NewFunction(&RawEntityType{}, &RawEntityType{}),
		IsConst: true,
	}, 
	"SetShipSpeed": { 
		Name:  "SetShipSpeed", 
		Value: NewFunction(&RawEntityType{}, NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}),
		IsConst: true,
	},
	"MakeShipTransparent": { 
		Name:  "MakeShipTransparent", 
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Number)),
		IsConst: true,
	},  
	"GetAllEntities": { 
		Name:  "GetAllEntities", 
		Value: NewFunction().WithReturns(NewWrapperType(NewBasicType(ast.List), &RawEntityType{})),
		IsConst: true,
	},
	"GetEntitiesInRadius": { 
		Name:  "GetEntitiesInRadius", 
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)).WithReturns(NewWrapperType(NewBasicType(ast.List), &RawEntityType{})),
		IsConst: true,
	}, 
	"GetEntityCount": { 
		Name:  "GetEntityCount", 
		Value: NewFunction(NewEnumType("Pewpew", "EntityType")).WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"GetEntityType": { 
		Name:  "GetEntityType", 
		Value: NewFunction(&RawEntityType{}).WithReturns(NewEnumType("Pewpew", "EntityType")),
		IsConst: true,
	},
	"PlayAmbientSound": { 
		Name:  "PlayAmbientSound", 
		Value: NewFunction(NewPathType(ast.SoundEnv), NewBasicType(ast.Number)),
		IsConst: true,
	},
	"PlaySound": { 
		Name:  "PlaySound", 
		Value: NewFunction(NewPathType(ast.SoundEnv), NewBasicType(ast.Number), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)),
		IsConst: true,
	},
	"CreateExplosion": { 
		Name:  "CreateExplosion", 
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewBasicType(ast.Number), NewFixedPointType(ast.Fixed), NewBasicType(ast.Number)),
		IsConst: true,
	},
	"NewAsteroid": { 
		Name:  "NewAsteroid", 
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)).WithReturns(&RawEntityType{}),
		IsConst: true,
	},
	"NewAsteroidWithSize": { 
		Name:  "NewAsteroidWithSize", 
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}),
		IsConst: true,
	}, 
	"NewYellowBaf": { 
		Name:  "NewYellowBaf", 
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Degree), NewFixedPointType(ast.Fixed), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}),
		IsConst: true,
	},
	"NewRedBaf": { 
		Name:  "NewRedBaf", 
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Degree), NewFixedPointType(ast.Fixed), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}),
		IsConst: true,
	},
	"NewBlueBaf": { 
		Name:  "NewBlueBaf", 
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Degree), NewFixedPointType(ast.Fixed), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}),
		IsConst: true,
	},
	"NewBomb": { 
		Name:  "NewBomb", 
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewEnumType("Pewpew", "BonusType")).WithReturns(&RawEntityType{}),
		IsConst: true,
	},
	"NewBonus": { // let's intergrate this with the walker now
		Name:  "NewBonus", 
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewEnumType("Pewpew", "BonusType"), NewStructType(map[string]Field{
			"box_duration": NewField(0, &VariableVal{
				Name: "box_duration",
				Value: &NumberVal{}, 
			}), 
			"cannon": NewField(1, &VariableVal{
				Name: "cannon",
				Value: CannonType,
			}),
			"frequency": NewField(2, &VariableVal{
				Name: "frequency",
				Value: CannonFrequency, 
			}),
			"weapon_duration": NewField(3, &VariableVal{
				Name: "weapon_duration",
				Value: &NumberVal{},
			}),
			"speed_factor": NewField(4, &VariableVal{
				Name: "speed_factor",
				Value: &FixedVal{ SpecificType: ast.Fixed },
			}),
			"speed_offset": NewField(5, &VariableVal{
				Name: "speed_offset",
				Value: &FixedVal{ SpecificType: ast.Fixed },
			}),
			"speed_duration": NewField(6, &VariableVal{
				Name: "speed_duration",
				Value: &NumberVal{},
			}),
			"taken_callback": NewField(7, &VariableVal{
				Name: "taken_callback", 
				Value: &FunctionVal{
					Params: Types{
						&RawEntityType{},
					},
				},
			}),
		}, true)).WithReturns(&RawEntityType{}),
		IsConst: true,
	},
	"NewCrowder": { 
		Name:  "NewCrowder", 
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)).WithReturns(&RawEntityType{}),
		IsConst: true,
	},
	"NewFloatingMessage": { 
		Name:  "NewFloatingMessage", 
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewBasicType(ast.String), NewStructType(map[string]Field{
			"scale": NewField(0, &VariableVal{
				Name: "scale",
				Value: &FixedVal{ SpecificType: ast.Fixed },
			}),
			"ticks_before_fade": NewField(1, &VariableVal{
				Name: "ticks_before_fade",
				Value: &NumberVal{},
			}),
			"is_optional": NewField(2, &VariableVal{
				Name: "is_optional",
				Value: &BoolVal{},
			}),
		}, true)).WithReturns(&RawEntityType{}),
		IsConst: true,
	},

	"NewEntity": {
		Name:  "NewEntity", 
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)).WithReturns(&RawEntityType{}),
		IsConst: true,
	},
	"NewInertiac": {
		Name: "NewInertiac",
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Degree)),
		IsConst: true,
	},
	"NewMothership": {
		Name: "NewMothership",
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewEnumType("Pewpew", "MothershipType"), NewFixedPointType(ast.Degree)).WithReturns(&RawEntityType{}),
		IsConst: true,
	},
	"NewPointonium": {
		Name: "NewPointonium",
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}),
	},
	"NewShip": {
		Name: "NewShip",
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}),
	},
	"NewBullet": {
		Name: "NewBullet",
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Degree), NewFixedPointType(ast.Number)).WithReturns(&RawEntityType{}),
	},
	"NewRollingCube": {
		Name: "NewRollingCube",
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)).WithReturns(&RawEntityType{}),
	},
	"NewRollingSphere": {
		Name: "NewRollingSphere",
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Degree), NewFixedPointType(ast.Fixed)).WithReturns(&RawEntityType{}),
	},
	"NewWary": {
		Name: "NewWary",
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)).WithReturns(&RawEntityType{}),
	},
	"NewUfo": {
		Name: "NewUfo",
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)).WithReturns(&RawEntityType{}),
	},
	"SetRollingCubeWallCollision": {
		Name: "SetRollingCubeWallCollision",
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Bool)).WithReturns(&RawEntityType{}),
	},
	"SetUFOWallCollision": {
		Name: "SetUFOWallCollision",
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Bool)).WithReturns(&RawEntityType{}),
	},
	"GetEntityPosition": {
		Name: "GetEntityPosition",
		Value: NewFunction(&RawEntityType{}).WithReturns(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)),
	},
	"IsEntityAlive": {
		Name: "IsEntityAlive",
		Value: NewFunction(&RawEntityType{}).WithReturns(NewBasicType(ast.Bool)),
	},
	"IsEntityBeingDestroyed": {
		Name: "IsEntityBeingDestroyed",
		Value: NewFunction(&RawEntityType{}).WithReturns(NewBasicType(ast.Bool)),
	},
	"SetEntityPosition": {
		Name: "SetEntityPosition",
		Value: NewFunction(&RawEntityType{}, NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)),
	},
	"SetEntityRadius": {
		Name: "SetEntityRadius",
		Value: NewFunction(&RawEntityType{}, NewFixedPointType(ast.Fixed)),
	},
	"SetEntityCallback": {
		Name: "SetEntityCallback",
		Value: NewFunction(&RawEntityType{}, NewFunctionType(Types{ &RawEntityType{} }, Types{})),
	},
	"DestroyEntity": {
		Name: "DestroyEntity",
		Value: NewFunction(&RawEntityType{}),
	},
	"EntityReactToWeapon": {
		Name: "EntityReactToWeapon",
		Value: NewFunction(&RawEntityType{}, NewStructType(map[string]Field{
			"type": NewField(0, &VariableVal{
				Name: "type",
				Value: WeaponType,
			}),
			"x": NewField(1, &VariableVal{
				Name: "x",
				Value: &FixedVal{SpecificType: ast.Fixed},
			}),
			"y": NewField(2, &VariableVal{
				Name: "y",
				Value: &FixedVal{SpecificType: ast.Fixed},
			}),
			"player_index": NewField(3, &VariableVal{
				Name: "PlayerIndex",
				Value: &NumberVal{},
			}),
		}, true)),
	},
	"SetEntityInterpolation": {
		Name: "SetEntityInterpolation",
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Bool)),
	},
	"SetEntityMesh": {
		Name: "SetEntityMesh",
		Value: NewFunction(&RawEntityType{}, NewPathType(ast.MeshEnv), NewBasicType(ast.Number)),
	},
	"SetEntityFlippingMeshes": {
		Name: "SetEntityFlippingMeshes",
		Value: NewFunction(&RawEntityType{}, NewPathType(ast.MeshEnv), NewBasicType(ast.Number), NewBasicType(ast.Number)),
	},
	"SetEntityColor": {
		Name: "SetEntityColor",
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Number)),
	},
	"EntitySkipAttributesInterpolation": {
		Name: "EntitySkipAttributesInterpolation",
		Value: NewFunction(&RawEntityType{}),
	},
	"NewString": {
		Name: "NewString",
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.String)),
	},
	"SetEntityMeshPosition": {
		Name: "SetEntityMeshPosition",
		Value: NewFunction(&RawEntityType{}, NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)),
	},
	"SetEntityMeshScale": {
		Name: "SetEntityMeshScale",
		Value: NewFunction(&RawEntityType{}, NewFixedPointType(ast.Fixed)),
	},
	"SetEntityMeshScaleXYZ": {
		Name: "SetEntityMeshScaleXYZ",
		Value: NewFunction(&RawEntityType{}, NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)),
	},
	"SetEntityMeshAngle": {
		Name: "SetEntityMeshAngle",
		Value: NewFunction(&RawEntityType{}, NewFixedPointType(ast.Fixed)),
	},
	"ConfigureEntityMusicResponse": {
		Name: "ConfigureEntityMusicResponse",
		Value: NewFunction(&RawEntityType{}, NewStructType(map[string]Field{
			"color_start": NewField(0, &VariableVal{
				Name: "color_start",
				Value: &NumberVal{},
			}),
			"color_end": NewField(1, &VariableVal{
				Name: "color_end",
				Value: &NumberVal{},
			}),
			"scale_x_start": NewField(2, &VariableVal{
				Name: "scale_x_start",
				Value: &FixedVal{ SpecificType: ast.Fixed },
			}),
			"scale_x_end": NewField(3, &VariableVal{
				Name: "scale_x_end",
				Value: &FixedVal{ SpecificType: ast.Fixed },
			}),
			"scale_y_start": NewField(4, &VariableVal{
				Name: "scale_y_start",
				Value: &FixedVal{ SpecificType: ast.Fixed },
			}),
			"scale_y_end": NewField(5, &VariableVal{
				Name: "scale_y_end",
				Value: &FixedVal{ SpecificType: ast.Fixed },
			}),
			"scale_z_start": NewField(6, &VariableVal{
				Name: "scale_z_start",
				Value: &FixedVal{ SpecificType: ast.Fixed },
			}),
			"scale_z_end": NewField(7, &VariableVal{
				Name: "scale_z_end",
				Value: &FixedVal{ SpecificType: ast.Fixed },
			}),
		}, true)),
	},
	"RotateEntityMesh": {
		Name: "RotateEntityMesh",
		Value: NewFunction(&RawEntityType{}, NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)),
	},
	"SetEntityVisibilityRadius": {
		Name: "SetEntityVisibilityRadius",
		Value: NewFunction(&RawEntityType{}, NewFixedPointType(ast.Fixed)),
	},
	"ConfigureEntityWallCollision": {
		Name: "ConfigureEntityWallCollision",
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Bool), NewFunctionType(Types{&RawEntityType{}, NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)}, Types{})),
	},
	"SetEntityPlayerCollision": {
		Name: "SetEntityPlayerCollision",
		Value: NewFunction(&RawEntityType{}, NewFunctionType(Types{&RawEntityType{}, NewBasicType(ast.Number), &RawEntityType{}}, Types{})),
	},
	"SetEntityWeaponCollision": {
		Name: "SetEntityWeaponCollision",
		Value: NewFunction(&RawEntityType{}, NewFunctionType(Types{&RawEntityType{}, NewBasicType(ast.Number), NewEnumType("Pewpew", "WeaponType")}, Types{NewBasicType(ast.Bool)})),
	},
	"SpawnEntity": {
		Name: "SpawnEntity",
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Number)),
	},
	"ExplodeEntity": {
		Name: "ExplodeEntity",
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Number)),
	},
}

var PewpewEntityType = NewEnumVal("Pewpew", "EntityType", false,
	"Asteroid",
	"YellowBaf",
	"Inertiac",
	"Mothership",
	"MothershipBullet",
	"RollingCube",
	"RollingSphere",
	"Ufo",
	"Wary",
	"Crowder",
	"CustomizableEntity",
	"Ship",
	"Bomb",
	"BlueBaf",
	"RedBaf",
	"WaryMissile",
	"UfoBullet",
	"PlayerBullet",
	"BombExplosion",
	"PlayerExplosion",
	"Bonus",
	"FloatingMessage",
	"Pointonium",
	"BonusImplosion",
)

var MothershipType = NewEnumVal("Pewpew", "MothershipType", false,
	"Triangle",
	"Square",
	"Pentagon",
	"Hexagon",
	"Heptagon",
)

var CannonType = NewEnumVal("Pewpew", "CannonType", false,
	"Single",
	"TicToc",
	"Double",
	"Triple",
	"FourDirections",
	"DoubleSwipe",
	"Hemisphere",
)

var CannonFrequency = NewEnumVal("Pewpew", "CannonFreq", false,
	"Freq30",
	"Freq15",
	"Freq10",
	"Freq7_5",
	"Freq6",
	"Freq5",
	"Freq3",
	"Freq2",
	"Freq1",
)

var BombType = NewEnumVal("Pewpew", "BombType", false,
	"Freeze",
	"Repulsive",
	"Atomize",
	"SmallAtomize",
	"SmallFreeze",
)

var BonusType = NewEnumVal("Pewpew", "BonusType", false,
	"Reinstantiation",
	"Shield",
	"Speed",
	"Weapon",
)

var WeaponType = NewEnumVal("Pewpew", "WeaponType", false,
	"Bullet",
	"FreezeExplosion",
	"RepulsiveExplosion",
	"AtomizeExplosion",
)

var AsteroidSize = NewEnumVal("Pewpew", "AsteroidSize", false,
	"Small",
	"Medium",
	"Large",
	"Enormous",
)