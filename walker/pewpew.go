package walker

import (
	"hybroid/ast"
)

var PewpewEnv = &Environment{
	Name: "Pewpew",
	Scope: Scope{
		Variables: pewpewVariables,
		Tag: &UntaggedTag{},
	},
	Structs: make(map[string]*StructVal),
	Entities: make(map[string]*EntityVal),
	CustomTypes: make(map[string]*CustomType),
}

var pewpewVariables = map[string]*VariableVal{
	//enums
	"EntityType": {
		Name:    "EntityType",
		Value:   EntityType,
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
	"CannonFrequency": {
		Name:    "CannonFrequency",
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
		Value: NewFunction(NewBasicType(ast.String)),
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
		Value: NewFunction(NewBasicType(ast.Number), NewAnonStructType(map[string]*VariableVal{
			"has_lost": {
				Name: "has_lost",
				Value: &BoolVal{},
			},
			"shield": {
				Name: "shield",
				Value: &NumberVal{},
			},
			"camera_x_override": {
				Name: "camera_x_override",
				Value: &FixedVal{SpecificType: ast.Fixed},
			},
			"camera_y_override": {
				Name: "camera_y_override",
				Value: &FixedVal{SpecificType: ast.Fixed},
			},
			"camera_distance": {
				Name: "camera_distance",
				Value: &FixedVal{SpecificType: ast.Fixed},
			},
			"camera_rotation_x_axis": {
				Name: "camera_rotation_x_axis",
				Value: &FixedVal{SpecificType: ast.Fixed},
			},
			"move_joystick_color": {
				Name: "move_joystick_color",
				Value: &NumberVal{},
			},
			"shoot_joystick_color": {
				Name: "shoot_joystick_color",
				Value: &NumberVal{},
			},
		})),
		IsConst: true,
	},
	"ConfigurePlayerHud": {
		Name: "ConfigurePlayerHud",
		Value: NewFunction(NewBasicType(ast.Number), NewAnonStructType(map[string]*VariableVal{
			"TopLeftLine": {
				Name: "TopLeftLine",
				Value: &StringVal{},
			},
		})),
	},
	"GetPlayerConfig": {  
		Name:  "GetPlayerConfig", 
		Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewAnonStructType(map[string]*VariableVal{
			"shield": {
				Name: "shield",
				Value: &NumberVal{},
			},
			"has_lost": {
				Name: "has_lost",
				Value: &BoolVal{},
			},
		})),
		IsConst: true,
	}, 
	"ConfigureShipWeapon": {  
		Name:  "ConfigureShipWeapon",
		Value: NewFunction(&RawEntityType{}, NewAnonStructType(map[string]*VariableVal{
			"frequency": {
				Name: "frequency",
				Value: CannonFrequency,
			},
			"cannon": {
				Name: "cannon",
				Value: CannonType,
			},
			"duration": {
				Name: "duration",
				Value: &NumberVal{},
			},
		})),
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
		Value: NewFunction(NewEnumType("EntityType")).WithReturns(NewBasicType(ast.Number)),
		IsConst: true,
	},
	"GetEntityType": { 
		Name:  "GetEntityType", 
		Value: NewFunction(&RawEntityType{}).WithReturns(NewEnumType("EntityType")),
		IsConst: true,
	},

	// pewpew.play_ambient_sound(
	// 	sound_path: string,
	// 	index: int
	//   )
		// "PlayAmbientSound": { 
		// 	Name:  "PlayAmbientSound", 
		// 	Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)).WithReturns(&RawEntityType{}),
		// 	IsConst: true, // ok
		// },// let's like, ignore playAmbientSound and PlaySound for now
	// pewpew.play_sound(
	// 	sound_path: string,
	// 	index: int,
	// 	x: FixedPoint,
	// 	y: FixedPoint
	// )
		// "NewEntity": { 
		// 	Name:  "NewEntity", 
		// 	Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)).WithReturns(&RawEntityType{}),
		// 	IsConst: true,
	// },
		
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
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewEnumType("BonusType")).WithReturns(&RawEntityType{}),
		IsConst: true,
	},
	"NewBonus": { // let's intergrate this with the walker now
		Name:  "NewBonus", 
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewEnumType("BonusType"), NewAnonStructType(map[string]*VariableVal{
			"box_duration": {
				Name: "box_duration",
				Value: &NumberVal{}, 
			}, 
			"cannon": {
				Name: "cannon",
				Value: CannonType,
			},
			"frequency": {
				Name: "frequency",
				Value: CannonFrequency, 
			},
			"weapon_duration": {
				Name: "weapon_duration",
				Value: &NumberVal{},
			},
			"speed_factor": {
				Name: "speed_factor",
				Value: &FixedVal{ SpecificType: ast.Fixed },
			},
			"speed_offset": {
				Name: "speed_offset",
				Value: &FixedVal{ SpecificType: ast.Fixed },
			},
			"speed_duration": {
				Name: "speed_duration",
				Value: &NumberVal{},
			},
			"taken_callback": {
				Name: "taken_callback", 
				Value: &FunctionVal{
					Params: Types{
						&RawEntityType{},
					},
				},
			},
		})).WithReturns(&RawEntityType{}),
		IsConst: true,
	},
	"NewCrowder": { 
		Name:  "NewCrowder", 
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)).WithReturns(&RawEntityType{}),
		IsConst: true,
	},
	"NewFloatingMessage": { 
		Name:  "NewFloatingMessage", 
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewBasicType(ast.String), NewAnonStructType(map[string]*VariableVal{
			"scale": {
				Name: "scale",
				Value: &FixedVal{ SpecificType: ast.Fixed },
			},
			"ticks_before_fade": {
				Name: "ticks_before_fade",
				Value: &NumberVal{},
			},
			"is_optional": {
				Name: "is_optional",
				Value: &BoolVal{},
			},
		})).WithReturns(&RawEntityType{}),
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
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed), NewEnumType("MothershipType"), NewFixedPointType(ast.Degree)).WithReturns(&RawEntityType{}),
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
	"GetEntityAlive": {
		Name: "GetEntityAlive",
		Value: NewFunction(&RawEntityType{}).WithReturns(NewBasicType(ast.Bool)),
	},
	"IsEntityBeingDestroyed": {
		Name: "IsEntityBeingDestroyed",
		Value: NewFunction(&RawEntityType{}).WithReturns(NewBasicType(ast.Bool)),
	},
	"SetEntityPosition": {
		Name: "SetEntityPosition",
		Value: NewFunction(NewFixedPointType(ast.Fixed), NewFixedPointType(ast.Fixed)),
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
		Value: NewFunction(&RawEntityType{}, NewAnonStructType(map[string]*VariableVal{
			"type": {
				Name: "type",
				Value: WeaponType,
			},
			"x": {
				Name: "x",
				Value: &FixedVal{SpecificType: ast.Fixed},
			},
			"y": {
				Name: "y",
				Value: &FixedVal{SpecificType: ast.Fixed},
			},
			"player_index": {
				Name: "PlayerIndex",
				Value: &NumberVal{},
			},
		})),
	},
	// pewpew.customizable_entity_set_position_interpolation(
	// 	entity_id: EntityId,
	// 	enable: bool
	//   )

	// pewpew.customizable_entity_set_mesh(
	// 	entity_id: EntityId,
	// 	file_path: string,
	// 	index: int
	//   )

	// pewpew.customizable_entity_set_flipping_meshes(
	// 	entity_id: EntityId,
	// 	file_path: string,
	// 	index_0: int,
	// 	index_1: int
	//   )

	// pewpew.customizable_entity_set_mesh_color(
	// 	entity_id: EntityId,
	// 	color: int
	//   )
}

var EntityType = NewEnumVal("EntityType", false,
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

var MothershipType = NewEnumVal("MothershipType", false,
	"Triangle",
	"Square",
	"Pentagon",
	"Hexagon",
	"Heptagon",
)

var CannonType = NewEnumVal("CannonType", false,
	"Single",
	"TicToc",
	"Double",
	"Triple",
	"FourDirections",
	"DoubleSwipe",
	"Hemisphere",
)

var CannonFrequency = NewEnumVal("CannonFrequency", false,
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

var BombType = NewEnumVal("BombType", false,
	"Freeze",
	"Repulsive",
	"Atomize",
	"SmallAtomize",
	"SmallFreeze",
)

var BonusType = NewEnumVal("BonusType", false,
	"Reinstantiation",
	"Shield",
	"Speed",
	"Weapon",
)

var WeaponType = NewEnumVal("WeaponType", false,
	"Bullet",
	"FreezeExplosion",
	"RepulsiveExplosion",
	"AtomizeExplosion",
)

var AsteroidSize = NewEnumVal("AsteroidSize", false,
	"Small",
	"Medium",
	"Large",
	"Enormous",
)