// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
package mapping

// AUTO-GENERATED ENUMS, DO NOT MANUALLY MODIFY!
var PewpewEnums = map[string]map[string]string{
	"EntityType": {
		"YellowBAF": "BAF", "BlueBAF": "BAF_BLUE", "RedBAF": "BAF_RED", "UFO": "UFO", "UFOBullet": "UFO_BULLET", "Asteroid": "ASTEROID", "Inertiac": "INERTIAC", "Mothership": "MOTHERSHIP", "MothershipBullet": "MOTHERSHIP_BULLET", "RollingCube": "ROLLING_CUBE", "RollingSphere": "ROLLING_SPHERE", "Wary": "WARY", "Crowder": "CROWDER", "CustomizableEntity": "CUSTOMIZABLE_ENTITY", "Ship": "SHIP", "Bomb": "BOMB", "WaryMissile": "WARY_MISSILE", "Spiny": "SPINY", "SuperMothership": "SUPER_MOTHERSHIP", "PlayerBullet": "PLAYER_BULLET", "BombExplosion": "BOMB_EXPLOSION", "PlayerExplosion": "PLAYER_EXPLOSION", "Bonus": "BONUS", "FloatingMessage": "FLOATING_MESSAGE", "Pointonium": "POINTONIUM", "Kamikaze": "KAMIKAZE", "BonusImplosion": "BONUS_IMPLOSION", "Mace": "MACE", "PlasmaField": "PLASMA_FIELD",
	}, "MothershipType": {
		"Triangle": "THREE_CORNERS", "Square": "FOUR_CORNERS", "Pentagon": "FIVE_CORNERS", "Hexagon": "SIX_CORNERS", "Heptagon": "SEVEN_CORNERS",
	}, "CannonFreq": {
		"Freq7_5": "FREQ_7_5", "Freq30": "FREQ_30", "Freq15": "FREQ_15", "Freq10": "FREQ_10", "Freq6": "FREQ_6", "Freq5": "FREQ_5", "Freq3": "FREQ_3", "Freq2": "FREQ_2", "Freq1": "FREQ_1",
	}, "CannonType": {
		"Single": "SINGLE", "TicToc": "TIC_TOC", "Double": "DOUBLE", "Triple": "TRIPLE", "FourDirections": "FOUR_DIRECTIONS", "DoubleSwipe": "DOUBLE_SWIPE", "Hemisphere": "HEMISPHERE", "Shotgun": "SHOTGUN", "Laser": "LASER",
	}, "BombType": {
		"Freeze": "FREEZE", "Repulsive": "REPULSIVE", "Atomize": "ATOMIZE", "SmallAtomize": "SMALL_ATOMIZE", "SmallFreeze": "SMALL_FREEZE",
	}, "MaceType": {
		"DamagePlayers": "DAMAGE_PLAYERS", "DamageEntities": "DAMAGE_ENTITIES",
	}, "BonusType": {
		"Reinstantiation": "REINSTANTIATION", "Shield": "SHIELD", "Speed": "SPEED", "Weapon": "WEAPON", "Mace": "MACE",
	}, "WeaponType": {
		"Bullet": "BULLET", "FreezeExplosion": "FREEZE_EXPLOSION", "RepulsiveExplosion": "REPULSIVE_EXPLOSION", "AtomizeExplosion": "ATOMIZE_EXPLOSION", "PlasmaField": "PLASMA_FIELD", "WallTrailLasso": "WALL_TRAIL_LASSO", "Mace": "MACE",
	}, "AsteroidSize": {
		"Small": "SMALL", "Medium": "MEDIUM", "Large": "LARGE", "VeryLarge": "VERY_LARGE",
	},
}

// AUTO-GENERATED VARIABLES, DO NOT MANUALLY MODIFY!
var PewpewVariables = map[string]string{
	"IncreasePlayerScore":                   "increase_score_of_player",
	"IncreasePlayerScoreStreak":             "increase_score_streak_of_player",
	"GetPlayerScore":                        "get_score_of_player",
	"GetPlayerScoreStreak":                  "get_score_streak_level",
	"GetPlayerConfig":                       "get_player_configuration",
	"DamageShip":                            "add_damage_to_player_ship",
	"IsEntityAlive":                         "entity_get_is_alive",
	"IsEntityBeingDestroyed":                "entity_get_is_started_to_be_destroyed",
	"NewEntity":                             "new_customizable_entity",
	"NewYellowBAF":                          "new_baf",
	"NewRedBAF":                             "new_baf_red",
	"NewBlueBAF":                            "new_baf_blue",
	"NewUFO":                                "new_ufo",
	"GetEntitiesInRadius":                   "get_entities_colliding_with_disk",
	"AddRotationToEntityMesh":               "customizable_entity_add_rotation_to_mesh",
	"SetEntityMusicResponse":                "customizable_entity_configure_music_response",
	"SetEntityMeshPosition":                 "customizable_entity_set_mesh_xyz",
	"SetEntityMeshXYZScale":                 "customizable_entity_set_mesh_xyz_scale",
	"SetEntityWallCollision":                "customizable_entity_configure_wall_collision",
	"SetUFOWallCollision":                   "ufo_set_enable_collisions_with_walls",
	"SetRollingCubeWallCollision":           "rolling_cube_set_enable_collisions_with_walls",
	"DestroyEntity":                         "entity_destroy",
	"SpawnEntity":                           "customizable_entity_start_spawning",
	"ExplodeEntity":                         "customizable_entity_start_exploding",
	"Print":                                 "print",
	"PrintDebugInfo":                        "print_debug_info",
	"SetLevelSize":                          "set_level_size",
	"AddWall":                               "add_wall",
	"RemoveWall":                            "remove_wall",
	"AddUpdateCallback":                     "add_update_callback",
	"GetNumberOfPlayers":                    "get_number_of_players",
	"StopGame":                              "stop_game",
	"GetPlayerInputs":                       "get_player_inputs",
	"ConfigurePlayer":                       "configure_player",
	"ConfigurePlayerHud":                    "configure_player_hud",
	"ConfigureShipWeapon":                   "configure_player_ship_weapon",
	"ConfigureShipWallTrail":                "configure_player_ship_wall_trail",
	"ConfigureShip":                         "configure_player_ship",
	"AddArrowToShip":                        "add_arrow_to_player_ship",
	"RemoveArrowFromShip":                   "remove_arrow_from_player_ship",
	"MakeShipTransparent":                   "make_player_ship_transparent",
	"SetShipSpeed":                          "set_player_ship_speed",
	"GetAllEntities":                        "get_all_entities",
	"GetEntityCount":                        "get_entity_count",
	"GetEntityType":                         "get_entity_type",
	"PlayAmbientSound":                      "play_ambient_sound",
	"PlaySound":                             "play_sound",
	"CreateExplosion":                       "create_explosion",
	"AddParticle":                           "add_particle",
	"NewAsteroid":                           "new_asteroid",
	"NewAsteroidWithSize":                   "new_asteroid_with_size",
	"NewBomb":                               "new_bomb",
	"NewBonus":                              "new_bonus",
	"NewCrowder":                            "new_crowder",
	"NewFloatingMessage":                    "new_floating_message",
	"NewInertiac":                           "new_inertiac",
	"NewKamikaze":                           "new_kamikaze",
	"NewMothership":                         "new_mothership",
	"NewMothershipBullet":                   "new_mothership_bullet",
	"NewPointonium":                         "new_pointonium",
	"NewPlasmaField":                        "new_plasma_field",
	"NewShip":                               "new_player_ship",
	"NewPlayerBullet":                       "new_player_bullet",
	"NewRollingCube":                        "new_rolling_cube",
	"NewRollingSphere":                      "new_rolling_sphere",
	"NewSpiny":                              "new_spiny",
	"NewSuperMothership":                    "new_super_mothership",
	"NewWary":                               "new_wary",
	"GetEntityPosition":                     "entity_get_position",
	"SetEntityPosition":                     "entity_set_position",
	"EntityMove":                            "entity_move",
	"SetEntityRadius":                       "entity_set_radius",
	"SetEntityUpdateCallback":               "entity_set_update_callback",
	"EntityReactToWeapon":                   "entity_react_to_weapon",
	"EntityAddMace":                         "entity_add_mace",
	"SetEntityPositionInterpolation":        "customizable_entity_set_position_interpolation",
	"SetEntityAngleInterpolation":           "customizable_entity_set_angle_interpolation",
	"SetEntityMesh":                         "customizable_entity_set_mesh",
	"SetEntityFlippingMeshes":               "customizable_entity_set_flipping_meshes",
	"SetEntityMeshColor":                    "customizable_entity_set_mesh_color",
	"SetEntityString":                       "customizable_entity_set_string",
	"SetEntityMeshZ":                        "customizable_entity_set_mesh_z",
	"SetEntityMeshScale":                    "customizable_entity_set_mesh_scale",
	"SetEntityMeshAngle":                    "customizable_entity_set_mesh_angle",
	"SkipEntityMeshAttributesInterpolation": "customizable_entity_skip_mesh_attributes_interpolation",
	"SetEntityVisibilityRadius":             "customizable_entity_set_visibility_radius",
	"SetEntityPlayerCollision":              "customizable_entity_set_player_collision_callback",
	"SetEntityWeaponCollision":              "customizable_entity_set_weapon_collision_callback",
	"SetEntityTag":                          "customizable_entity_set_tag",
	"GetEntityTag":                          "customizable_entity_get_tag",

	"EntityType":     "EntityType",
	"MothershipType": "MothershipType",
	"CannonFreq":     "CannonFrequency",
	"CannonType":     "CannonType",
	"BombType":       "BombType",
	"MaceType":       "MaceType",
	"BonusType":      "BonusType",
	"WeaponType":     "WeaponType",
	"AsteroidSize":   "AsteroidSize",
}
