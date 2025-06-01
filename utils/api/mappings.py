# The API mapping dictionary holds the initial mapping of Lua enum variants/functions to Hybroid
# Gets populated as the generation goes
_API_MAPPING = {
    # EntityType
    "BAF": "YellowBaf",
    # MothershipType
    "THREE_CORNERS": "Triangle",
    "FOUR_CORNERS": "Square",
    "FIVE_CORNERS": "Pentagon",
    "SIX_CORNERS": "Hexagon",
    "SEVEN_CORNERS": "Heptagon",
    # CannonFrequency
    "FREQ_7_5": "Freq7_5",
    # AsteroidSize
    "VERY_LARGE": "VeryLarge",
}


def get(key: str, fn) -> str:
    """
    Attempts to find value in mapping by key, if not found, uses the specified function to generate one.
    """

    value = _API_MAPPING.get(key, None)
    if value is None:
        value = fn(key)
        _API_MAPPING[key] = value

    return value
