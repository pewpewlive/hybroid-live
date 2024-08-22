package walker

import (
	"hybroid/ast"
)

var PewpewEnv = &Environment{
	Name: "Pewpew",
	Scope: Scope{
		Variables: PewpewVariables,
		Tag:       &UntaggedTag{},
	},
	UsedWalkers:   make([]*Walker, 0),
	UsedLibraries: make(map[Library]bool),
	Structs:       make(map[string]*ClassVal),
	Entities:      make(map[string]*EntityVal),
	CustomTypes:   make(map[string]*CustomType),
}

var PewpewVariables = map[string]*VariableVal{
	//enums
	"EntityType": {
		Name:    "EntityType",
		Value:   PewpewEntityType,
		IsLocal: false,
	},
	"MothershipType": {
		Name:    "MothershipType",
		Value:   MothershipType,
		IsLocal: false,
	},
	"CannonType": {
		Name:    "CannonType",
		Value:   CannonType,
		IsLocal: false,
	},
	"CannonFreq": {
		Name:    "CannonFreq",
		Value:   CannonFrequency,
		IsLocal: false,
	},
	"BombType": {
		Name:    "BombType",
		Value:   BombType,
		IsLocal: false,
	},
	"BonusType": {
		Name:    "BonusType",
		Value:   BonusType,
		IsLocal: false,
	},
	"WeaponType": {
		Name:    "WeaponType",
		Value:   WeaponType,
		IsLocal: false,
	},
	"AsteroidSize": {
		Name:    "AsteroidSize",
		Value:   AsteroidSize,
		IsLocal: false,
	},

	//functions
	"Print": {
		Name:    "Print",
		Value:   NewFunction(NewBasicType(ast.Object)),
	},
	"PrintDebugInfo": {
		Name:    "PrintDebugInfo",
		Value:   NewFunction(),
	},
	"SetLevelSize": {
		Name:    "SetLevelSize",
		Value:   NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)),
	},
	"AddWall": {
		Name:    "AddWall",
		Value:   NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)).WithReturns(NewBasicType(ast.Number)),
	},
	"RemoveWall": {
		Name:    "SetLevelSize",
		Value:   NewFunction(NewBasicType(ast.Number)),
	},
	"AddUpdateCallback": {
		Name:    "AddUpdateCallback",
		Value:   NewFunction(NewFunctionType(Types{}, Types{})),
	},
	"GetNumberOfPlayers": {
		Name:    "GetNumberOfPlayers",
		Value:   NewFunction().WithReturns(NewBasicType(ast.Number)),
	},
	"IncreasePlayerScore": {
		Name:    "IncreasePlayerScore",
		Value:   NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)),
	},
	"IncreasePlayerScoreStreak": {
		Name:    "IncreasePlayerScoreStreak",
		Value:   NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)),
	},
	"GetPlayerScoreStreak": {
		Name:    "GetPlayerScoreStreak",
		Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
	},
	"StopGame": {
		Name:    "StopGame",
		Value:   NewFunction(),
	},
	"GetPlayerInputs": {
		Name:    "GetPlayerInputs",
		Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)),
	},
	"GetPlayerScore": {
		Name:    "GetPlayerScore",
		Value:   NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
	},
	"ConfigurePlayer": {
		Name: "ConfigurePlayer",
		Value: NewFunction(NewBasicType(ast.Number), NewStructType([]*VariableVal{
			{
				Name:  "has_lost",
				Value: &BoolVal{},
			},
			{
				Name:  "shield",
				Value: &NumberVal{},
			},
			{
				Name:  "camera_x_override",
				Value: &FixedVal{SpecificType: ast.Fixed},
			},
			{
				Name:  "camera_y_override",
				Value: &FixedVal{SpecificType: ast.Fixed},
			},
			{
				Name:  "camera_distance",
				Value: &FixedVal{SpecificType: ast.Fixed},
			},
			{
				Name:  "camera_rotation_x_axis",
				Value: &FixedVal{SpecificType: ast.Fixed},
			},
			{
				Name:  "move_joystick_color",
				Value: &NumberVal{},
			},
			{
				Name:  "shoot_joystick_color",
				Value: &NumberVal{},
			},
		}, true)),
	},
	"ConfigurePlayerHud": {
		Name: "ConfigurePlayerHud",
		Value: NewFunction(NewBasicType(ast.Number), NewStructType([]*VariableVal{
			{
				Name:  "TopLeftLine",
				Value: &StringVal{},
			},
		}, true)),
	},
	"GetPlayerConfig": {
		Name: "GetPlayerConfig",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewStructType([]*VariableVal{
			{
				Name:  "shield",
				Value: &NumberVal{},
			},
			{
				Name:  "has_lost",
				Value: &BoolVal{},
			},
		}, false)),
	},
	"ConfigureShipWeapon": {
		Name: "ConfigureShipWeapon",
		Value: NewFunction(&RawEntityType{}, NewStructType([]*VariableVal{
			{
				Name:  "frequency",
				Value: CannonFrequency,
			},
			{
				Name:  "cannon",
				Value: CannonType,
			},
			{
				Name:  "duration",
				Value: &NumberVal{},
			},
		}, true)),
	},
	"DamageShip": {
		Name:    "DamageShip",
		Value:   NewFunction(&RawEntityType{}, NewBasicType(ast.Number)),
	},
	"AddArrowToShip": {
		Name:    "AddArrowToShip",
		Value:   NewFunction(&RawEntityType{}, &RawEntityType{}, NewBasicType(ast.Number)).WithReturns(&RawEntityType{}),
	},
	"RemoveArrowFromShip": {
		Name:    "RemoveArrowFromShip",
		Value:   NewFunction(&RawEntityType{}, &RawEntityType{}),
	},
	"SetShipSpeed": {
		Name:    "SetShipSpeed",
		Value:   NewFunction(&RawEntityType{}, NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}),
	},
	"MakeShipTransparent": {
		Name:    "MakeShipTransparent",
		Value:   NewFunction(&RawEntityType{}, NewBasicType(ast.Number)),
	},
	"GetAllEntities": {
		Name:    "GetAllEntities",
		Value:   NewFunction().WithReturns(NewWrapperType(NewBasicType(ast.List), &RawEntityType{})),
	},
	"GetEntitiesInRadius": {
		Name:    "GetEntitiesInRadius",
		Value:   NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)).WithReturns(NewWrapperType(NewBasicType(ast.List), &RawEntityType{})),
	},
	"GetEntityCount": {
		Name:    "GetEntityCount",
		Value:   NewFunction(NewEnumType("Pewpew", "EntityType")).WithReturns(NewBasicType(ast.Number)),
	},
	"GetEntityType": {
		Name:    "GetEntityType",
		Value:   NewFunction(&RawEntityType{}).WithReturns(NewEnumType("Pewpew", "EntityType")),
	},
	"PlayAmbientSound": {
		Name:    "PlayAmbientSound",
		Value:   NewFunction(NewPathType(ast.SoundEnv), NewBasicType(ast.Number)),
	},
	"PlaySound": {
		Name:    "PlaySound",
		Value:   NewFunction(NewPathType(ast.SoundEnv), NewBasicType(ast.Number), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)),
	},
	"CreateExplosion": {
		Name:    "CreateExplosion",
		Value:   NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewBasicType(ast.Number), NewFixedPointType(ast.Fixed), NewBasicType(ast.Number)),
	},
	"NewAsteroid": {
		Name:    "NewAsteroid",
		Value:   NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)).WithReturns(&RawEntityType{}),
	},
	"NewAsteroidWithSize": {
		Name:    "NewAsteroidWithSize",
		Value:   NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}),
	},
	"NewYellowBaf": {
		Name:    "NewYellowBaf",
		Value:   NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Degree), NewFixedPointType(ast.Fixed), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}),
	},
	"NewRedBaf": {
		Name:    "NewRedBaf",
		Value:   NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Degree), NewFixedPointType(ast.Fixed), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}),
	},
	"NewBlueBaf": {
		Name:    "NewBlueBaf",
		Value:   NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Degree), NewFixedPointType(ast.Fixed), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}),
	},
	"NewBomb": {
		Name:    "NewBomb",
		Value:   NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewEnumType("Pewpew", "BonusType")).WithReturns(&RawEntityType{}),
	},
	"NewBonus": { // let's intergrate this with the walker now
		Name: "NewBonus",
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewEnumType("Pewpew", "BonusType"), NewStructType([]*VariableVal{
			{
				Name:  "box_duration",
				Value: &NumberVal{},
			},
			{
				Name:  "cannon",
				Value: CannonType,
			},
			{
				Name:  "frequency",
				Value: CannonFrequency,
			},
			{
				Name:  "weapon_duration",
				Value: &NumberVal{},
			},
			{
				Name:  "speed_factor",
				Value: &FixedVal{SpecificType: ast.Fixed},
			},
			{
				Name:  "speed_offset",
				Value: &FixedVal{SpecificType: ast.Fixed},
			},
			{
				Name:  "speed_duration",
				Value: &NumberVal{},
			},
			{
				Name: "taken_callback",
				Value: &FunctionVal{
					Params: Types{
						&RawEntityType{},
						&RawEntityType{},
						&RawEntityType{},
					},
				},
			},
		}, true)).WithReturns(&RawEntityType{}),
	},
	"NewCrowder": {
		Name:    "NewCrowder",
		Value:   NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)).WithReturns(&RawEntityType{}),
	},
	"NewFloatingMessage": {
		Name: "NewFloatingMessage",
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewBasicType(ast.String), NewStructType([]*VariableVal{
			{
				Name:  "scale",
				Value: &FixedVal{SpecificType: ast.Fixed},
			},
			{
				Name:  "ticks_before_fade",
				Value: &NumberVal{},
			},
			{
				Name:  "is_optional",
				Value: &BoolVal{},
			},
		}, true)).WithReturns(&RawEntityType{}),
	},

	"NewEntity": {
		Name:    "NewEntity",
		Value:   NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)).WithReturns(&RawEntityType{}),
	},
	"NewInertiac": {
		Name:    "NewInertiac",
		Value:   NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Degree)),
	},
	"NewMothership": {
		Name:    "NewMothership",
		Value:   NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewEnumType("Pewpew", "MothershipType"), NewFixedPointType(ast.Degree)).WithReturns(&RawEntityType{}),
	},
	"NewPointonium": {
		Name:  "NewPointonium",
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}),
	},
	"NewShip": {
		Name:  "NewShip",
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}),
	},
	"NewBullet": {
		Name:  "NewBullet",
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Degree), NewFixedPointType(ast.Number)).WithReturns(&RawEntityType{}),
	},
	"NewRollingCube": {
		Name:  "NewRollingCube",
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)).WithReturns(&RawEntityType{}),
	},
	"NewRollingSphere": {
		Name:  "NewRollingSphere",
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Degree), NewFixedPointType(ast.Fixed)).WithReturns(&RawEntityType{}),
	},
	"NewWary": {
		Name:  "NewWary",
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)).WithReturns(&RawEntityType{}),
	},
	"NewUfo": {
		Name:  "NewUfo",
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)).WithReturns(&RawEntityType{}),
	},
	"SetRollingCubeWallCollision": {
		Name:  "SetRollingCubeWallCollision",
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Bool)).WithReturns(&RawEntityType{}),
	},
	"SetUFOWallCollision": {
		Name:  "SetUFOWallCollision",
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Bool)).WithReturns(&RawEntityType{}),
	},
	"GetEntityPosition": {
		Name:  "GetEntityPosition",
		Value: NewFunction(&RawEntityType{}).WithReturns(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)),
	},
	"IsEntityAlive": {
		Name:  "IsEntityAlive",
		Value: NewFunction(&RawEntityType{}).WithReturns(NewBasicType(ast.Bool)),
	},
	"IsEntityBeingDestroyed": {
		Name:  "IsEntityBeingDestroyed",
		Value: NewFunction(&RawEntityType{}).WithReturns(NewBasicType(ast.Bool)),
	},
	"SetEntityPosition": {
		Name:  "SetEntityPosition",
		Value: NewFunction(&RawEntityType{}, NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)),
	},
	"SetEntityRadius": {
		Name:  "SetEntityRadius",
		Value: NewFunction(&RawEntityType{}, NewFixedPointType(ast.Fixed)),
	},
	"SetEntityCallback": {
		Name:  "SetEntityCallback",
		Value: NewFunction(&RawEntityType{}, NewFunctionType(Types{&RawEntityType{}}, Types{})),
	},
	"DestroyEntity": {
		Name:  "DestroyEntity",
		Value: NewFunction(&RawEntityType{}),
	},
	"EntityReactToWeapon": {
		Name: "EntityReactToWeapon",
		Value: NewFunction(&RawEntityType{}, NewStructType([]*VariableVal{
			{
				Name:  "type",
				Value: WeaponType,
			},
			{
				Name:  "x",
				Value: &FixedVal{SpecificType: ast.Fixed},
			},
			{
				Name:  "y",
				Value: &FixedVal{SpecificType: ast.Fixed},
			},
			{
				Name:  "PlayerIndex",
				Value: &NumberVal{},
			},
		}, true)),
	},
	"SetEntityInterpolation": {
		Name:  "SetEntityInterpolation",
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Bool)),
	},
	"SetEntityMesh": {
		Name:  "SetEntityMesh",
		Value: NewFunction(&RawEntityType{}, NewPathType(ast.MeshEnv), NewBasicType(ast.Number)),
	},
	"SetEntityFlippingMeshes": {
		Name:  "SetEntityFlippingMeshes",
		Value: NewFunction(&RawEntityType{}, NewPathType(ast.MeshEnv), NewBasicType(ast.Number), NewBasicType(ast.Number)),
	},
	"SetEntityColor": {
		Name:  "SetEntityColor",
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Number)),
	},
	"EntitySkipAttributesInterpolation": {
		Name:  "EntitySkipAttributesInterpolation",
		Value: NewFunction(&RawEntityType{}),
	},
	"NewString": {
		Name:  "NewString",
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.String)),
	},
	"SetEntityMeshPosition": {
		Name:  "SetEntityMeshPosition",
		Value: NewFunction(&RawEntityType{}, NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)),
	},
	"SetEntityMeshScale": {
		Name:  "SetEntityMeshScale",
		Value: NewFunction(&RawEntityType{}, NewFixedPointType(ast.Fixed)),
	},
	"SetEntityMeshScaleXYZ": {
		Name:  "SetEntityMeshScaleXYZ",
		Value: NewFunction(&RawEntityType{}, NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)),
	},
	"SetEntityMeshAngle": {
		Name:  "SetEntityMeshAngle",
		Value: NewFunction(&RawEntityType{}, NewFixedPointType(ast.Fixed)),
	},
	"ConfigureEntityMusicResponse": {
		Name: "ConfigureEntityMusicResponse",
		Value: NewFunction(&RawEntityType{}, NewStructType([]*VariableVal{
			{
				Name:  "color_start",
				Value: &NumberVal{},
			},
			{
				Name:  "color_end",
				Value: &NumberVal{},
			},
			{
				Name:  "scale_x_start",
				Value: &FixedVal{SpecificType: ast.Fixed},
			},
			 {
				Name:  "scale_x_end",
				Value: &FixedVal{SpecificType: ast.Fixed},
			},
			{
				Name:  "scale_y_start",
				Value: &FixedVal{SpecificType: ast.Fixed},
			},
			{
				Name:  "scale_y_end",
				Value: &FixedVal{SpecificType: ast.Fixed},
			},
			{
				Name:  "scale_z_start",
				Value: &FixedVal{SpecificType: ast.Fixed},
			},
			{
				Name:  "scale_z_end",
				Value: &FixedVal{SpecificType: ast.Fixed},
			},
		}, true)),
	},
	"RotateEntityMesh": {
		Name:  "RotateEntityMesh",
		Value: NewFunction(&RawEntityType{}, NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)),
	},
	"SetEntityVisibilityRadius": {
		Name:  "SetEntityVisibilityRadius",
		Value: NewFunction(&RawEntityType{}, NewFixedPointType(ast.Fixed)),
	},
	"ConfigureEntityWallCollision": {
		Name:  "ConfigureEntityWallCollision",
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Bool), NewFunctionType(Types{&RawEntityType{}, NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)}, Types{})),
	},
	"SetEntityPlayerCollision": {
		Name:  "SetEntityPlayerCollision",
		Value: NewFunction(&RawEntityType{}, NewFunctionType(Types{&RawEntityType{}, NewBasicType(ast.Number), &RawEntityType{}}, Types{})),
	},
	"SetEntityWeaponCollision": {
		Name:  "SetEntityWeaponCollision",
		Value: NewFunction(&RawEntityType{}, NewFunctionType(Types{&RawEntityType{}, NewBasicType(ast.Number), NewEnumType("Pewpew", "WeaponType")}, Types{NewBasicType(ast.Bool)})),
	},
	"SpawnEntity": {
		Name:  "SpawnEntity",
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Number)),
	},
	"ExplodeEntity": {
		Name:  "ExplodeEntity",
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
