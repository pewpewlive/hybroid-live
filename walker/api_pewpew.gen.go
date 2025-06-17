// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
package walker

import "hybroid/ast"

// AUTO-GENERATED API, DO NOT MANUALLY MODIFY!
var PewpewAPI = &Environment{
	Name: "Pewpew",
	Scope: Scope{
		Variables: map[string]*VariableVal{
			"Print": {
				Name: "Print", Value: NewFunction(NewBasicType(ast.Text)), IsPub: true,
			},
			"PrintDebugInfo": {
				Name: "PrintDebugInfo", Value: NewFunction(), IsPub: true,
			},
			"SetLevelSize": {
				Name: "SetLevelSize", Value: NewFunction(NewFixedPointType(), NewFixedPointType()), IsPub: true,
			},
			"AddWall": {
				Name: "AddWall", Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewFixedPointType()).WithReturns(NewBasicType(ast.Number)), IsPub: true,
			},
			"RemoveWall": {
				Name: "RemoveWall", Value: NewFunction(NewBasicType(ast.Number)), IsPub: true,
			},
			"AddUpdateCallback": {
				Name: "AddUpdateCallback", Value: NewFunction(NewFunctionType([]Type{}, []Type{})), IsPub: true,
			},
			"GetNumberOfPlayers": {
				Name: "GetNumberOfPlayers", Value: NewFunction().WithReturns(NewBasicType(ast.Number)), IsPub: true,
			},
			"IncreasePlayerScore": {
				Name: "IncreasePlayerScore", Value: NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)), IsPub: true,
			},
			"IncreasePlayerScoreStreak": {
				Name: "IncreasePlayerScoreStreak", Value: NewFunction(NewBasicType(ast.Number), NewBasicType(ast.Number)), IsPub: true,
			},
			"GetPlayerScoreStreak": {
				Name: "GetPlayerScoreStreak", Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)), IsPub: true,
			},
			"StopGame": {
				Name: "StopGame", Value: NewFunction(), IsPub: true,
			},
			"GetPlayerInputs": {
				Name: "GetPlayerInputs", Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewFixedPointType()), IsPub: true,
			},
			"GetPlayerScore": {
				Name: "GetPlayerScore", Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)), IsPub: true,
			},
			"ConfigurePlayer": {
				Name: "ConfigurePlayer", Value: NewFunction(NewBasicType(ast.Number), NewStructType([]StructField{NewStructField("has_lost", &BoolVal{}, true), NewStructField("shield", &NumberVal{}, true), NewStructField("camera_x_override", &FixedVal{}, true), NewStructField("camera_y_override", &FixedVal{}, true), NewStructField("camera_distance", &FixedVal{}, true), NewStructField("camera_rotation_x_axis", &FixedVal{}, true), NewStructField("move_joystick_color", &NumberVal{}, true), NewStructField("shoot_joystick_color", &NumberVal{}, true)})), IsPub: true,
			},
			"ConfigurePlayerHud": {
				Name: "ConfigurePlayerHud", Value: NewFunction(NewBasicType(ast.Number), NewStructType([]StructField{NewStructField("top_left_line", &StringVal{}, true)})), IsPub: true,
			},
			"GetPlayerConfig": {
				Name: "GetPlayerConfig", Value: NewFunction(NewBasicType(ast.Number)).WithReturns(NewStructType([]StructField{NewStructField("shield", &NumberVal{}, true), NewStructField("has_lost", &BoolVal{}, true)})), IsPub: true,
			},
			"ConfigureShipWeapon": {
				Name: "ConfigureShipWeapon", Value: NewFunction(&RawEntityType{}, NewStructType([]StructField{NewStructField("frequency", NewEnumVal("Pewpew", "CannonFreq", true), true), NewStructField("cannon", NewEnumVal("Pewpew", "CannonType", true), true), NewStructField("duration", &NumberVal{}, true)})), IsPub: true,
			},
			"ConfigureShipWallTrail": {
				Name: "ConfigureShipWallTrail", Value: NewFunction(&RawEntityType{}, NewStructType([]StructField{NewStructField("wall_length", &NumberVal{}, true)})), IsPub: true,
			},
			"ConfigureShip": {
				Name: "ConfigureShip", Value: NewFunction(&RawEntityType{}, NewStructType([]StructField{NewStructField("swap_inputs", &BoolVal{}, true)})), IsPub: true,
			},
			"DamageShip": {
				Name: "DamageShip", Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Number)), IsPub: true,
			},
			"AddArrowToShip": {
				Name: "AddArrowToShip", Value: NewFunction(&RawEntityType{}, &RawEntityType{}, NewBasicType(ast.Number)).WithReturns(NewBasicType(ast.Number)), IsPub: true,
			},
			"RemoveArrowFromShip": {
				Name: "RemoveArrowFromShip", Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Number)), IsPub: true,
			},
			"MakeShipTransparent": {
				Name: "MakeShipTransparent", Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Number)), IsPub: true,
			},
			"SetShipSpeed": {
				Name: "SetShipSpeed", Value: NewFunction(&RawEntityType{}, NewFixedPointType(), NewFixedPointType(), NewBasicType(ast.Number)).WithReturns(NewFixedPointType()), IsPub: true,
			},
			"GetAllEntities": {
				Name: "GetAllEntities", Value: NewFunction().WithReturns(NewWrapperType(NewBasicType(ast.List), &RawEntityType{})), IsPub: true,
			},
			"GetEntitiesInRadius": {
				Name: "GetEntitiesInRadius", Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewFixedPointType()).WithReturns(NewWrapperType(NewBasicType(ast.List), &RawEntityType{})), IsPub: true,
			},
			"GetEntityCount": {
				Name: "GetEntityCount", Value: NewFunction(NewEnumType("Pewpew", "EntityType")).WithReturns(NewBasicType(ast.Number)), IsPub: true,
			},
			"GetEntityType": {
				Name: "GetEntityType", Value: NewFunction(&RawEntityType{}).WithReturns(NewEnumType("Pewpew", "EntityType")), IsPub: true,
			},
			"PlayAmbientSound": {
				Name: "PlayAmbientSound", Value: NewFunction(NewPathType(ast.SoundEnv), NewBasicType(ast.Number)), IsPub: true,
			},
			"PlaySound": {
				Name: "PlaySound", Value: NewFunction(NewPathType(ast.SoundEnv), NewBasicType(ast.Number), NewFixedPointType(), NewFixedPointType()), IsPub: true,
			},
			"CreateExplosion": {
				Name: "CreateExplosion", Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewBasicType(ast.Number), NewFixedPointType(), NewBasicType(ast.Number)), IsPub: true,
			},
			"AddParticle": {
				Name: "AddParticle", Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewBasicType(ast.Number), NewBasicType(ast.Number)), IsPub: true,
			},
			"NewAsteroid": {
				Name: "NewAsteroid", Value: NewFunction(NewFixedPointType(), NewFixedPointType()).WithReturns(&RawEntityType{}), IsPub: true,
			},
			"NewAsteroidWithSize": {
				Name: "NewAsteroidWithSize", Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewEnumType("Pewpew", "AsteroidSize")).WithReturns(&RawEntityType{}), IsPub: true,
			},
			"NewYellowBAF": {
				Name: "NewYellowBAF", Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}), IsPub: true,
			},
			"NewRedBAF": {
				Name: "NewRedBAF", Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}), IsPub: true,
			},
			"NewBlueBAF": {
				Name: "NewBlueBAF", Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}), IsPub: true,
			},
			"NewBomb": {
				Name: "NewBomb", Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewEnumType("Pewpew", "BombType")).WithReturns(&RawEntityType{}), IsPub: true,
			},
			"NewBonus": {
				Name: "NewBonus", Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewEnumType("Pewpew", "BonusType"), NewStructType([]StructField{NewStructField("box_duration", &NumberVal{}, true), NewStructField("cannon", NewEnumVal("Pewpew", "CannonType", true), true), NewStructField("frequency", NewEnumVal("Pewpew", "CannonFreq", true), true), NewStructField("weapon_duration", &NumberVal{}, true), NewStructField("number_of_shields", &NumberVal{}, true), NewStructField("speed_factor", &FixedVal{}, true), NewStructField("speed_offset", &FixedVal{}, true), NewStructField("speed_duration", &NumberVal{}, true), NewStructField("taken_callback", &FunctionVal{Params: []Type{&RawEntityType{}, NewBasicType(ast.Number), &RawEntityType{}}}, true)})).WithReturns(&RawEntityType{}), IsPub: true,
			},
			"NewCrowder": {
				Name: "NewCrowder", Value: NewFunction(NewFixedPointType(), NewFixedPointType()).WithReturns(&RawEntityType{}), IsPub: true,
			},
			"NewFloatingMessage": {
				Name: "NewFloatingMessage", Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewBasicType(ast.Text), NewStructType([]StructField{NewStructField("scale", &FixedVal{}, true), NewStructField("dz", &FixedVal{}, true), NewStructField("ticks_before_fade", &NumberVal{}, true), NewStructField("is_optional", &BoolVal{}, true)})).WithReturns(&RawEntityType{}), IsPub: true,
			},
			"NewEntity": {
				Name: "NewEntity", Value: NewFunction(NewFixedPointType(), NewFixedPointType()).WithReturns(&RawEntityType{}), IsPub: true,
			},
			"NewInertiac": {
				Name: "NewInertiac", Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewFixedPointType()).WithReturns(&RawEntityType{}), IsPub: true,
			},
			"NewKamikaze": {
				Name: "NewKamikaze", Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewFixedPointType()).WithReturns(&RawEntityType{}), IsPub: true,
			},
			"NewMothership": {
				Name: "NewMothership", Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewEnumType("Pewpew", "MothershipType"), NewFixedPointType()).WithReturns(&RawEntityType{}), IsPub: true,
			},
			"NewMothershipBullet": {
				Name: "NewMothershipBullet", Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewBasicType(ast.Number), NewBasicType(ast.Bool)).WithReturns(&RawEntityType{}), IsPub: true,
			},
			"NewPointonium": {
				Name: "NewPointonium", Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}), IsPub: true,
			},
			"NewPlasmaField": {
				Name: "NewPlasmaField", Value: NewFunction(&RawEntityType{}, &RawEntityType{}, NewStructType([]StructField{NewStructField("length", &FixedVal{}, true), NewStructField("stiffness", &FixedVal{}, true)})).WithReturns(&RawEntityType{}), IsPub: true,
			},
			"NewShip": {
				Name: "NewShip", Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}), IsPub: true,
			},
			"NewPlayerBullet": {
				Name: "NewPlayerBullet", Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewBasicType(ast.Number)).WithReturns(&RawEntityType{}), IsPub: true,
			},
			"NewRollingCube": {
				Name: "NewRollingCube", Value: NewFunction(NewFixedPointType(), NewFixedPointType()).WithReturns(&RawEntityType{}), IsPub: true,
			},
			"NewRollingSphere": {
				Name: "NewRollingSphere", Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewFixedPointType()).WithReturns(&RawEntityType{}), IsPub: true,
			},
			"NewSpiny": {
				Name: "NewSpiny", Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewFixedPointType()).WithReturns(&RawEntityType{}), IsPub: true,
			},
			"NewSuperMothership": {
				Name: "NewSuperMothership", Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewEnumType("Pewpew", "MothershipType"), NewFixedPointType()).WithReturns(&RawEntityType{}), IsPub: true,
			},
			"NewWary": {
				Name: "NewWary", Value: NewFunction(NewFixedPointType(), NewFixedPointType()).WithReturns(&RawEntityType{}), IsPub: true,
			},
			"NewUFO": {
				Name: "NewUFO", Value: NewFunction(NewFixedPointType(), NewFixedPointType(), NewFixedPointType()).WithReturns(&RawEntityType{}), IsPub: true,
			},
			"RollingCubeSetEnableCollisionsWithWalls": {
				Name: "RollingCubeSetEnableCollisionsWithWalls", Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Bool)), IsPub: true,
			},
			"SetUFOWallCollision": {
				Name: "SetUFOWallCollision", Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Bool)), IsPub: true,
			},
			"GetEntityPosition": {
				Name: "GetEntityPosition", Value: NewFunction(&RawEntityType{}).WithReturns(NewFixedPointType(), NewFixedPointType()), IsPub: true,
			},
			"IsEntityAlive": {
				Name: "IsEntityAlive", Value: NewFunction(&RawEntityType{}).WithReturns(NewBasicType(ast.Bool)), IsPub: true,
			},
			"IsEntityBeingDestroyed": {
				Name: "IsEntityBeingDestroyed", Value: NewFunction(&RawEntityType{}).WithReturns(NewBasicType(ast.Bool)), IsPub: true,
			},
			"SetEntityPosition": {
				Name: "SetEntityPosition", Value: NewFunction(&RawEntityType{}, NewFixedPointType(), NewFixedPointType()), IsPub: true,
			},
			"EntityMove": {
				Name: "EntityMove", Value: NewFunction(&RawEntityType{}, NewFixedPointType(), NewFixedPointType()), IsPub: true,
			},
			"SetEntityRadius": {
				Name: "SetEntityRadius", Value: NewFunction(&RawEntityType{}, NewFixedPointType()), IsPub: true,
			},
			"SetEntityUpdateCallback": {
				Name: "SetEntityUpdateCallback", Value: NewFunction(&RawEntityType{}, NewFunctionType([]Type{&RawEntityType{}}, []Type{})), IsPub: true,
			},
			"DestroyEntity": {
				Name: "DestroyEntity", Value: NewFunction(&RawEntityType{}), IsPub: true,
			},
			"EntityReactToWeapon": {
				Name: "EntityReactToWeapon", Value: NewFunction(&RawEntityType{}, NewStructType([]StructField{NewStructField("type", NewEnumVal("Pewpew", "WeaponType", true), true), NewStructField("x", &FixedVal{}, true), NewStructField("y", &FixedVal{}, true), NewStructField("player_index", &NumberVal{}, true)})).WithReturns(NewBasicType(ast.Bool)), IsPub: true,
			},
			"EntityAddMace": {
				Name: "EntityAddMace", Value: NewFunction(&RawEntityType{}, NewStructType([]StructField{NewStructField("distance", &FixedVal{}, true), NewStructField("angle", &FixedVal{}, true), NewStructField("rotation_speed", &FixedVal{}, true), NewStructField("type", NewEnumVal("Pewpew", "MaceType", true), true)})).WithReturns(&RawEntityType{}), IsPub: true,
			},
			"SetEntityPositionInterpolation": {
				Name: "SetEntityPositionInterpolation", Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Bool)), IsPub: true,
			},
			"SetEntityAngleInterpolation": {
				Name: "SetEntityAngleInterpolation", Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Bool)), IsPub: true,
			},
			"SetEntityMesh": {
				Name: "SetEntityMesh", Value: NewFunction(&RawEntityType{}, NewPathType(ast.MeshEnv), NewBasicType(ast.Number)), IsPub: true,
			},
			"SetEntityFlippingMeshes": {
				Name: "SetEntityFlippingMeshes", Value: NewFunction(&RawEntityType{}, NewPathType(ast.MeshEnv), NewBasicType(ast.Number), NewBasicType(ast.Number)), IsPub: true,
			},
			"SetEntityMeshColor": {
				Name: "SetEntityMeshColor", Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Number)), IsPub: true,
			},
			"SetEntityString": {
				Name: "SetEntityString", Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Text)), IsPub: true,
			},
			"SetEntityMeshPosition": {
				Name: "SetEntityMeshPosition", Value: NewFunction(&RawEntityType{}, NewFixedPointType(), NewFixedPointType(), NewFixedPointType()), IsPub: true,
			},
			"SetEntityMeshZ": {
				Name: "SetEntityMeshZ", Value: NewFunction(&RawEntityType{}, NewFixedPointType()), IsPub: true,
			},
			"SetEntityMeshScale": {
				Name: "SetEntityMeshScale", Value: NewFunction(&RawEntityType{}, NewFixedPointType()), IsPub: true,
			},
			"SetEntityMeshXYZScale": {
				Name: "SetEntityMeshXYZScale", Value: NewFunction(&RawEntityType{}, NewFixedPointType(), NewFixedPointType(), NewFixedPointType()), IsPub: true,
			},
			"SetEntityMeshAngle": {
				Name: "SetEntityMeshAngle", Value: NewFunction(&RawEntityType{}, NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewFixedPointType()), IsPub: true,
			},
			"SkipEntityMeshAttributesInterpolation": {
				Name: "SkipEntityMeshAttributesInterpolation", Value: NewFunction(&RawEntityType{}), IsPub: true,
			},
			"SetEntityMusicResponse": {
				Name: "SetEntityMusicResponse", Value: NewFunction(&RawEntityType{}, NewStructType([]StructField{NewStructField("color_start", &NumberVal{}, true), NewStructField("color_end", &NumberVal{}, true), NewStructField("scale_x_start", &FixedVal{}, true), NewStructField("scale_x_end", &FixedVal{}, true), NewStructField("scale_y_start", &FixedVal{}, true), NewStructField("scale_y_end", &FixedVal{}, true), NewStructField("scale_z_start", &FixedVal{}, true), NewStructField("scale_z_end", &FixedVal{}, true)})), IsPub: true,
			},
			"AddRotationToEntityMesh": {
				Name: "AddRotationToEntityMesh", Value: NewFunction(&RawEntityType{}, NewFixedPointType(), NewFixedPointType(), NewFixedPointType(), NewFixedPointType()), IsPub: true,
			},
			"SetEntityVisibilityRadius": {
				Name: "SetEntityVisibilityRadius", Value: NewFunction(&RawEntityType{}, NewFixedPointType()), IsPub: true,
			},
			"SetEntityWallCollision": {
				Name: "SetEntityWallCollision", Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Bool), NewFunctionType([]Type{&RawEntityType{}, NewFixedPointType(), NewFixedPointType()}, []Type{})), IsPub: true,
			},
			"SetEntityPlayerCollision": {
				Name: "SetEntityPlayerCollision", Value: NewFunction(&RawEntityType{}, NewFunctionType([]Type{&RawEntityType{}, NewBasicType(ast.Number), &RawEntityType{}}, []Type{})), IsPub: true,
			},
			"SetEntityWeaponCollision": {
				Name: "SetEntityWeaponCollision", Value: NewFunction(&RawEntityType{}, NewFunctionType([]Type{&RawEntityType{}, NewBasicType(ast.Number), NewEnumType("Pewpew", "WeaponType")}, []Type{NewBasicType(ast.Bool)})), IsPub: true,
			},
			"SpawnEntity": {
				Name: "SpawnEntity", Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Number)), IsPub: true,
			},
			"ExplodeEntity": {
				Name: "ExplodeEntity", Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Number)), IsPub: true,
			},
			"SetEntityTag": {
				Name: "SetEntityTag", Value: NewFunction(&RawEntityType{}, NewBasicType(ast.Number)), IsPub: true,
			},
			"GetEntityTag": {
				Name: "GetEntityTag", Value: NewFunction(&RawEntityType{}).WithReturns(NewBasicType(ast.Number)), IsPub: true,
			},
		},
		Tag:         &UntaggedTag{},
		AliasTypes:  make(map[string]*AliasType),
		ConstValues: make(map[string]ast.Node),
	},
	imports:       make([]Import, 0),
	UsedLibraries: make([]ast.Library, 0),
	Classes:       make(map[string]*ClassVal),
	Entities:      make(map[string]*EntityVal),
	Enums: map[string]*EnumVal{
		"EntityType":     NewEnumVal("Pewpew", "EntityType", true, "Asteroid", "YellowBAF", "Inertiac", "Mothership", "MothershipBullet", "RollingCube", "RollingSphere", "UFO", "Wary", "Crowder", "CustomizableEntity", "Ship", "Bomb", "BlueBAF", "RedBAF", "WaryMissile", "UFOBullet", "Spiny", "SuperMothership", "PlayerBullet", "BombExplosion", "PlayerExplosion", "Bonus", "FloatingMessage", "Pointonium", "Kamikaze", "BonusImplosion", "Mace", "PlasmaField"),
		"MothershipType": NewEnumVal("Pewpew", "MothershipType", true, "Triangle", "Square", "Pentagon", "Hexagon", "Heptagon"),
		"CannonType":     NewEnumVal("Pewpew", "CannonType", true, "Single", "TicToc", "Double", "Triple", "FourDirections", "DoubleSwipe", "Hemisphere", "Shotgun", "Laser"),
		"CannonFreq":     NewEnumVal("Pewpew", "CannonFreq", true, "Freq30", "Freq15", "Freq10", "Freq7_5", "Freq6", "Freq5", "Freq3", "Freq2", "Freq1"),
		"BombType":       NewEnumVal("Pewpew", "BombType", true, "Freeze", "Repulsive", "Atomize", "SmallAtomize", "SmallFreeze"),
		"MaceType":       NewEnumVal("Pewpew", "MaceType", true, "DamagePlayers", "DamageEntities"),
		"BonusType":      NewEnumVal("Pewpew", "BonusType", true, "Reinstantiation", "Shield", "Speed", "Weapon", "Mace"),
		"WeaponType":     NewEnumVal("Pewpew", "WeaponType", true, "Bullet", "FreezeExplosion", "RepulsiveExplosion", "AtomizeExplosion", "PlasmaField", "WallTrailLasso", "Mace"),
		"AsteroidSize":   NewEnumVal("Pewpew", "AsteroidSize", true, "Small", "Medium", "Large", "VeryLarge")},
}
