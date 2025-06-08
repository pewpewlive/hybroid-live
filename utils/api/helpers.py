import re


def pascal_case(original: str) -> str:
    """Converts a string from `snake_case` to `PascalCase`."""

    return original.title().replace("_", "")


def camel_case(original: str) -> str:
    """Converts a string from `snake_case` to `camelCase`."""

    string = pascal_case(original)
    return string[0].lower() + string[1:]


def camel_case_all(original: str) -> str:
    """Converts all occurences of `snake_case` to `camelCase`."""

    camel_case_matches = re.findall(r"`([\w_]+)`", original)

    for match in camel_case_matches:
        original = original.replace(match, camel_case(match))

    return original


def pewpew_conversion(original: str) -> str:
    conversions = {
        "customizable_entity_set": "set_entity",
        "customizable_entity_get": "get_entity",
        "customizable_entity_skip": "skip_entity",
        "collision_callback": "collision",
        "player_ship": "ship",
        "entity_set": "set_entity",
        "entity_get": "get_entity",
        "fixedpoint": "fixed",
    }

    for ppl, hyb in conversions.items():
        original = original.replace(ppl, hyb)

    return pascal_case(original)
