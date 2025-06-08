# The function mapping dictionary holds the initial mapping of Lua functions to Hybroid
# Gets populated as the generation goes
_FUNCTION_MAPPING = {
    # Fmath Functions
    "random_int": "RandomNumber",
    "to_int": "ToNumber",
    # Pewpew Functions
    "increase_score_of_player": "IncreasePlayerScore",
    "increase_score_streak_of_player": "IncreasePlayerScoreStreak",
    "get_score_of_player": "GetPlayerScore",
    "get_score_streak_level": "GetPlayerScoreStreak",
    "get_player_configuration": "GetPlayerConfig",
    "add_damage_to_player_ship": "DamageShip",
    "entity_get_is_alive": "IsEntityAlive",
    "entity_get_is_started_to_be_destroyed": "IsEntityBeingDestroyed",
    "new_customizable_entity": "NewEntity",
    "new_baf": "NewYellowBAF",
    "new_baf_red": "NewRedBAF",
    "new_baf_blue": "NewBlueBAF",
    "new_ufo": "NewUFO",
    "get_entities_colliding_with_disk": "GetEntitiesInRadius",
    "customizable_entity_add_rotation_to_mesh": "AddRotationToEntityMesh",
    "customizable_entity_configure_music_response": "SetEntityMusicResponse",
    "customizable_entity_set_mesh_xyz": "SetEntityMeshPosition",
    "customizable_entity_configure_wall_collision": "SetEntityWallCollision",
    "ufo_set_enable_collisions_with_walls": "SetUFOWallCollision",
    "entity_destroy": "DestroyEntity",
    "customizable_entity_start_spawning": "SpawnEntity",
    "customizable_entity_start_exploding": "ExplodeEntity",
}

_ENUM_MAPPING = {
    "EntityType": (
        "EntityType",
        {
            "BAF": "YellowBAF",
            "BAF_BLUE": "BlueBAF",
            "BAF_RED": "RedBAF",
            "UFO": "UFO",
            "UFO_BULLET": "UFOBullet",
        },
    ),
    "MothershipType": (
        "MothershipType",
        {
            "THREE_CORNERS": "Triangle",
            "FOUR_CORNERS": "Square",
            "FIVE_CORNERS": "Pentagon",
            "SIX_CORNERS": "Hexagon",
            "SEVEN_CORNERS": "Heptagon",
        },
    ),
    "CannonFrequency": ("CannonFreq", {"FREQ_7_5": "Freq7_5"}),
}


def get_function(key: str, fn) -> str:
    """
    Attempts to find value in mapping by key, if not found, uses the specified function to generate one.
    """

    value = _FUNCTION_MAPPING.get(key, None)
    if value is None:
        value = fn(key)
        _FUNCTION_MAPPING[key] = value

    return value


def get_enum(key: str, value: str | None, fn) -> str:
    """
    Attempts to find value in mapping by key first, and return new name if not found, otherwise uses the specified function to generate new value names and returns them.
    """

    enum = _ENUM_MAPPING.get(key, None)
    if enum is None and value is None:
        _ENUM_MAPPING[key] = (key, {})
        return key
    elif enum is not None and value is None:
        return enum[0]
    elif enum is None and value is not None:
        for name, mapping in _ENUM_MAPPING.items():
            if mapping[0] == key:
                enum = mapping
                key = name

    assert value is not None, "value cannot be None with a valid parent Enum"
    assert enum is not None, "failed to find converted Enum"

    variant = enum[1].get(value, None)
    if variant is None:
        variant = fn(value)
        _ENUM_MAPPING[key][1][value] = variant

    return variant


def inverse_mapping() -> dict[str, str]:
    return {hyb: ppl for ppl, hyb in _FUNCTION_MAPPING.items()}
