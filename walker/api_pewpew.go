package walker

import (
	"hybroid/ast"
)

var PewpewEnv = &Environment{
	Name: "Pewpew",
	Scope: Scope{
		Variables:   PewpewVariables,
		Tag:         &UntaggedTag{},
		AliasTypes:  make(map[string]*AliasType),
		ConstValues: make(map[string]ast.Node),
	},
	importedWalkers: make([]*Walker, 0),
	UsedLibraries:   make([]Library, 0),
	Classes:         make(map[string]*ClassVal),
	Entities:        make(map[string]*EntityVal),
	Enums: map[string]*EnumVal{
		"EntityType":     PewpewEntityType,
		"MothershipType": MothershipType,
		"CannonType":     CannonType,
		"CannonFreq":     CannonFrequency,
		"BombType":       BombType,
		"BonusType":      BonusType,
		"AsteroidSize":   AsteroidSize,
		"WeaponType":     WeaponType,
	},
}

var PewpewVariables = map[string]*VariableVal{
	//functions
	"Print": {
		Name:  "Print",
		Value: NewFunction(NewBasicType(ast.Object)),
		IsPub: true,
	},
	"PrintDebugInfo": {
		Name:  "PrintDebugInfo",
		Value: NewFunction(),
		IsPub: true,
	},
	"SetLevelSize": {
		Name:  "SetLevelSize",
		Value: NewFunction(NewFixedPointType(), NewFixedPointType()),
		IsPub: true,
	},
	"AddWall": {
		Name:  "AddWall",
		Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewFixedPointType()).WithReturns(NewBasicType(ast.Number)),
		IsPub: true,
	},
	"RemoveWall": {
		Name:  "SetLevelSize",
		Value: NewFunction(NewBasicType(ast.Number)),
		IsPub: true,
	},
	"AddUpdateCallback": {
		Name:  "AddUpdateCallback",
		Value: NewFunction(NewFunctionType([]Type{}, []Type{})),
		IsPub: true,
	},
	"GetNumberOfPlayers": {
		Name:  "GetNumberOfPlayers",
		Value: NewFunction().WithReturns(NewBasicType(ast.Number)),
		IsPub: true,
	},
	"IncreasePlayerScore": {
		Name:  "IncreasePlayerScore",
		Value: NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)),
		IsPub: true,
	},
	"IncreasePlayerScoreStreak": {
		Name:  "IncreasePlayerScoreStreak",
		Value: NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)),
		IsPub: true,
	},
	"GetPlayerScoreStreak": {
		Name:  "GetPlayerScoreStreak",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsPub: true,
	},
	"StopGame": {
		Name:  "StopGame",
		Value: NewFunction(),
		IsPub: true,
	},
	"GetPlayerInputs": {
		Name:  "GetPlayerInputs",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewFixedPointType()),
		IsPub: true,
	},
	"GetPlayerScore": {
		Name:  "GetPlayerScore",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)),
		IsPub: true,
	},
	"ConfigurePlayer": {
		Name: "ConfigurePlayer",
		Value: NewFunction(NewBasicType(ast.Number), NewStructType([]*VariableVal{
			{
				Name:  "has_lost",
				Value: NewBoolVal(),
			},
			{
				Name:  "shield",
				Value: NewNumberVal(),
			},
			{
				Name:  "camera_x_override",
				Value: &FixedVal{},
			},
			{
				Name:  "camera_y_override",
				Value: &FixedVal{},
			},
			{
				Name:  "camera_distance",
				Value: &FixedVal{},
			},
			{
				Name:  "camera_rotation_x_axis",
				Value: &FixedVal{},
			},
			{
				Name:  "move_joystick_color",
				Value: NewNumberVal(),
			},
			{
				Name:  "shoot_joystick_color",
				Value: NewNumberVal(),
			},
		}, true)),
		IsPub: true,
	},
	"ConfigurePlayerHud": {
		Name: "ConfigurePlayerHud",
		Value: NewFunction(NewBasicType(ast.Number), NewStructType([]*VariableVal{
			{
				Name:  "TopLeftLine",
				Value: &StringVal{},
			},
		}, true)),
		IsPub: true,
	},
	"GetPlayerConfig": {
		Name: "GetPlayerConfig",
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewStructType([]*VariableVal{
			{
				Name:  "shield",
				Value: NewNumberVal(),
			},
			{
				Name:  "has_lost",
				Value: NewBoolVal(),
			},
		}, false)),
		IsPub: true,
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
				Value: NewNumberVal(),
			},
		}, true)),
		IsPub: true,
	},
	"DamageShip": {
		Name:  "DamageShip",
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Number)),
		IsPub: true,
	},
	"AddArrowToShip": {
		Name:  "AddArrowToShip",
		Value: NewFunction(&RawEntityType{}, &RawEntityType{}, NewBasicType(ast.Number)).WithReturns(&RawEntityType{}),
		IsPub: true,
	},
	"RemoveArrowFromShip": {
		Name:  "RemoveArrowFromShip",
		Value: NewFunction(&RawEntityType{}, &RawEntityType{}),
		IsPub: true,
	},
	"SetShipSpeed": {
		Name:  "SetShipSpeed",
		Value: NewFunction(&RawEntityType{}, NewFixedPointType(), NewFixedPointType(), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}),
		IsPub: true,
	},
	"MakeShipTransparent": {
		Name:  "MakeShipTransparent",
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Number)),
		IsPub: true,
	},
	"GetAllEntities": {
		Name:  "GetAllEntities",
		Value: NewFunction().WithReturns(NewWrapperType(NewBasicType(ast.List), &RawEntityType{})),
		IsPub: true,
	},
	"GetEntitiesCollidingWithDisk": {
		Name:  "GetEntitiesCollidingWithDisk",
		Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewFixedPointType()).WithReturns(NewWrapperType(NewBasicType(ast.List), &RawEntityType{})),
		IsPub: true,
	},
	"GetEntitiesInRadius": {
		Name:  "GetEntitiesInRadius",
		Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewFixedPointType()).WithReturns(NewWrapperType(NewBasicType(ast.List), &RawEntityType{})),
		IsPub: true,
	},
	"GetEntityCount": {
		Name:  "GetEntityCount",
		Value: NewFunction(NewEnumType("Pewpew", "EntityType")).WithReturns(NewBasicType(ast.Number)),
		IsPub: true,
	},
	"GetEntityType": {
		Name:  "GetEntityType",
		Value: NewFunction(&RawEntityType{}).WithReturns(NewEnumType("Pewpew", "EntityType")),
		IsPub: true,
	},
	"PlayAmbientSound": {
		Name:  "PlayAmbientSound",
		Value: NewFunction(NewPathType(ast.SoundEnv), NewBasicType(ast.Number)),
		IsPub: true,
	},
	"PlaySound": {
		Name:  "PlaySound",
		Value: NewFunction(NewPathType(ast.SoundEnv), NewBasicType(ast.Number), NewFixedPointType(), NewFixedPointType()),
		IsPub: true,
	},
	"CreateExplosion": {
		Name:  "CreateExplosion",
		Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewBasicType(ast.Number), NewFixedPointType(), NewBasicType(ast.Number)),
		IsPub: true,
	},
	"NewAsteroid": {
		Name:  "NewAsteroid",
		Value: NewFunction(NewFixedPointType(), NewFixedPointType()).WithReturns(&RawEntityType{}),
		IsPub: true,
	},
	"NewAsteroidWithSize": {
		Name:  "NewAsteroidWithSize",
		Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}),
		IsPub: true,
	},
	"NewYellowBaf": {
		Name:  "NewYellowBaf",
		Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}),
		IsPub: true,
	},
	"NewRedBaf": {
		Name:  "NewRedBaf",
		Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}),
		IsPub: true,
	},
	"NewBlueBaf": {
		Name:  "NewBlueBaf",
		Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}),
		IsPub: true,
	},
	"NewBomb": {
		Name:  "NewBomb",
		Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewEnumType("Pewpew", "BonusType")).WithReturns(&RawEntityType{}),
		IsPub: true,
	},
	"NewBonus": {
		Name: "NewBonus",
		Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewEnumType("Pewpew", "BonusType"), NewStructType([]*VariableVal{
			{
				Name:  "box_duration",
				Value: NewNumberVal(),
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
				Value: NewNumberVal(),
			},
			{
				Name:  "speed_factor",
				Value: &FixedVal{},
			},
			{
				Name:  "speed_offset",
				Value: &FixedVal{},
			},
			{
				Name:  "speed_duration",
				Value: NewNumberVal(),
			},
			{
				Name: "taken_callback",
				Value: &FunctionVal{
					Params: []Type{
						&RawEntityType{},
						&RawEntityType{},
						&RawEntityType{},
					},
				},
			},
		}, true)).WithReturns(&RawEntityType{}),
		IsPub: true,
	},
	"NewCrowder": {
		Name:  "NewCrowder",
		Value: NewFunction(NewFixedPointType(), NewFixedPointType()).WithReturns(&RawEntityType{}),
		IsPub: true,
	},
	"NewFloatingMessage": {
		Name: "NewFloatingMessage",
		Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewBasicType(ast.String), NewStructType([]*VariableVal{
			{
				Name:  "scale",
				Value: &FixedVal{},
			},
			{
				Name:  "ticks_before_fade",
				Value: NewNumberVal(),
			},
			{
				Name:  "is_optional",
				Value: NewBoolVal(),
			},
		}, true)).WithReturns(&RawEntityType{}),
		IsPub: true,
	},

	"NewEntity": {
		Name:  "NewEntity",
		Value: NewFunction(NewFixedPointType(), NewFixedPointType()).WithReturns(&RawEntityType{}),
		IsPub: true,
	},
	"NewInertiac": {
		Name:  "NewInertiac",
		Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewFixedPointType()),
		IsPub: true,
	},
	"NewMothership": {
		Name:  "NewMothership",
		Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewEnumType("Pewpew", "MothershipType"), NewFixedPointType()).WithReturns(&RawEntityType{}),
		IsPub: true,
	},
	"NewPointonium": {
		Name:  "NewPointonium",
		Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}),
		IsPub: true,
	},
	"NewShip": {
		Name:  "NewShip",
		Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}),
		IsPub: true,
	},
	"NewBullet": {
		Name:  "NewBullet",
		Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}),
		IsPub: true,
	},
	"NewRollingCube": {
		Name:  "NewRollingCube",
		Value: NewFunction(NewFixedPointType(), NewFixedPointType()).WithReturns(&RawEntityType{}),
		IsPub: true,
	},
	"NewRollingSphere": {
		Name:  "NewRollingSphere",
		Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewFixedPointType()).WithReturns(&RawEntityType{}),
		IsPub: true,
	},
	"NewWary": {
		Name:  "NewWary",
		Value: NewFunction(NewFixedPointType(), NewFixedPointType()).WithReturns(&RawEntityType{}),
		IsPub: true,
	},
	"NewUfo": {
		Name:  "NewUfo",
		Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewFixedPointType()).WithReturns(&RawEntityType{}),
		IsPub: true,
	},
	"SetRollingCubeWallCollision": {
		Name:  "SetRollingCubeWallCollision",
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Bool)).WithReturns(&RawEntityType{}),
		IsPub: true,
	},
	"SetUFOWallCollision": {
		Name:  "SetUFOWallCollision",
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Bool)).WithReturns(&RawEntityType{}),
		IsPub: true,
	},
	"GetEntityPosition": {
		Name:  "GetEntityPosition",
		Value: NewFunction(&RawEntityType{}).WithReturns(NewFixedPointType(), NewFixedPointType()),
		IsPub: true,
	},
	"IsEntityAlive": {
		Name:  "IsEntityAlive",
		Value: NewFunction(&RawEntityType{}).WithReturns(NewBasicType(ast.Bool)),
		IsPub: true,
	},
	"IsEntityBeingDestroyed": {
		Name:  "IsEntityBeingDestroyed",
		Value: NewFunction(&RawEntityType{}).WithReturns(NewBasicType(ast.Bool)),
		IsPub: true,
	},
	"SetEntityPosition": {
		Name:  "SetEntityPosition",
		Value: NewFunction(&RawEntityType{}, NewFixedPointType(), NewFixedPointType()),
		IsPub: true,
	},
	"SetEntityRadius": {
		Name:  "SetEntityRadius",
		Value: NewFunction(&RawEntityType{}, NewFixedPointType()),
		IsPub: true,
	},
	"SetEntityCallback": {
		Name:  "SetEntityCallback",
		Value: NewFunction(&RawEntityType{}, NewFunctionType([]Type{&RawEntityType{}}, []Type{})),
		IsPub: true,
	},
	"DestroyEntity": {
		Name:  "DestroyEntity",
		Value: NewFunction(&RawEntityType{}),
		IsPub: true,
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
				Value: &FixedVal{},
			},
			{
				Name:  "y",
				Value: &FixedVal{},
			},
			{
				Name:  "player_index",
				Value: NewNumberVal(),
			},
		}, false)),
		IsPub: true,
	},
	"SetEntityInterpolation": {
		Name:  "SetEntityInterpolation",
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Bool)),
		IsPub: true,
	},
	"SetEntityMesh": {
		Name:  "SetEntityMesh",
		Value: NewFunction(&RawEntityType{}, NewPathType(ast.MeshEnv), NewBasicType(ast.Number)),
		IsPub: true,
	},
	"SetEntityFlippingMeshes": {
		Name:  "SetEntityFlippingMeshes",
		Value: NewFunction(&RawEntityType{}, NewPathType(ast.MeshEnv), NewBasicType(ast.Number), NewBasicType(ast.Number)),
		IsPub: true,
	},
	"SetEntityColor": {
		Name:  "SetEntityColor",
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Number)),
		IsPub: true,
	},
	"EntitySkipAttributesInterpolation": {
		Name:  "EntitySkipAttributesInterpolation",
		Value: NewFunction(&RawEntityType{}),
		IsPub: true,
	},
	"SetEntityMeshPosition": {
		Name:  "SetEntityMeshPosition",
		Value: NewFunction(&RawEntityType{}, NewFixedPointType(), NewFixedPointType(), NewFixedPointType()),
		IsPub: true,
	},
	"SetEntityMeshScale": {
		Name:  "SetEntityMeshScale",
		Value: NewFunction(&RawEntityType{}, NewFixedPointType()),
		IsPub: true,
	},
	"SetEntityMeshScaleXYZ": {
		Name:  "SetEntityMeshScaleXYZ",
		Value: NewFunction(&RawEntityType{}, NewFixedPointType(), NewFixedPointType(), NewFixedPointType()),
		IsPub: true,
	},
	"SetEntityMeshAngle": {
		Name:  "SetEntityMeshAngle",
		Value: NewFunction(&RawEntityType{}, NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewFixedPointType()),
		IsPub: true,
	},
	"ConfigureEntityMusicResponse": {
		Name: "ConfigureEntityMusicResponse",
		Value: NewFunction(&RawEntityType{}, NewStructType([]*VariableVal{
			{
				Name:  "color_start",
				Value: NewNumberVal(),
			},
			{
				Name:  "color_end",
				Value: NewNumberVal(),
			},
			{
				Name:  "scale_x_start",
				Value: &FixedVal{},
			},
			{
				Name:  "scale_x_end",
				Value: &FixedVal{},
			},
			{
				Name:  "scale_y_start",
				Value: &FixedVal{},
			},
			{
				Name:  "scale_y_end",
				Value: &FixedVal{},
			},
			{
				Name:  "scale_z_start",
				Value: &FixedVal{},
			},
			{
				Name:  "scale_z_end",
				Value: &FixedVal{},
			},
		}, true)),
		IsPub: true,
	},
	"RotateEntityMesh": {
		Name:  "RotateEntityMesh",
		Value: NewFunction(&RawEntityType{}, NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewFixedPointType()),
		IsPub: true,
	},
	"SetEntityVisibilityRadius": {
		Name:  "SetEntityVisibilityRadius",
		Value: NewFunction(&RawEntityType{}, NewFixedPointType()),
		IsPub: true,
	},
	"ConfigureEntityWallCollision": {
		Name:  "ConfigureEntityWallCollision",
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Bool), NewFunctionType([]Type{NewFixedPointType(), NewFixedPointType()}, []Type{})),
		IsPub: true,
	},
	"SetEntityWallCollision": {
		Name:  "SetEntityWallCollision",
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Bool)),
		IsPub: true,
	},
	"SetEntityPlayerCollision": {
		Name:  "SetEntityPlayerCollision",
		Value: NewFunction(&RawEntityType{}, NewFunctionType([]Type{&RawEntityType{}, NewBasicType(ast.Number), &RawEntityType{}}, []Type{})),
		IsPub: true,
	},
	"SetEntityWeaponCollision": {
		Name:  "SetEntityWeaponCollision",
		Value: NewFunction(&RawEntityType{}, NewFunctionType([]Type{&RawEntityType{}, NewBasicType(ast.Number), NewEnumType("Pewpew", "WeaponType")}, []Type{NewBasicType(ast.Bool)})),
		IsPub: true,
	},
	"SpawnEntity": {
		Name:  "SpawnEntity",
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Number)),
		IsPub: true,
	},
	"ExplodeEntity": {
		Name:  "ExplodeEntity",
		Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Number)),
		IsPub: true,
	},
}

var PewpewEntityType = NewEnumVal("Pewpew", "EntityType", true,
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
	"Mace",
	"PlasmaField",
)

var MaceType = NewEnumVal("Pewpew", "MaceType", true,
	"DamagePlayers",
	"DamageEntities",
)

var MothershipType = NewEnumVal("Pewpew", "MothershipType", true,
	"Triangle",
	"Square",
	"Pentagon",
	"Hexagon",
	"Heptagon",
)

var CannonType = NewEnumVal("Pewpew", "CannonType", true,
	"Single",
	"TicToc",
	"Double",
	"Triple",
	"FourDirections",
	"DoubleSwipe",
	"Hemisphere",
	"Shotgun",
	"Laser",
)

var CannonFrequency = NewEnumVal("Pewpew", "CannonFreq", true,
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

var BombType = NewEnumVal("Pewpew", "BombType", true,
	"Freeze",
	"Repulsive",
	"Atomize",
	"SmallAtomize",
	"SmallFreeze",
)

var BonusType = NewEnumVal("Pewpew", "BonusType", true,
	"Reinstantiation",
	"Shield",
	"Speed",
	"Weapon",
	"Mace",
)

var WeaponType = NewEnumVal("Pewpew", "WeaponType", true,
	"Bullet",
	"FreezeExplosion",
	"RepulsiveExplosion",
	"AtomizeExplosion",
	"PlasmaField",
	"WallTrailLasso",
	"Mace",
)

var AsteroidSize = NewEnumVal("Pewpew", "AsteroidSize", true,
	"Small",
	"Medium",
	"Large",
	"Enormous",
)
