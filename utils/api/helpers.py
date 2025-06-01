import re


def to_pascal_case(original: str) -> str:
    """Converts a string from `snake_case` to `PascalCase`."""

    return original.title().replace("_", "")


def to_camel_case(original: str) -> str:
    """Converts a string from `snake_case` to `camelCase`."""

    string = to_pascal_case(original)
    return string[0].lower() + string[1:]


def to_camel_case_all(original: str) -> str:
    """Converts all occurences of `snake_case` to `camelCase`."""

    camel_case_matches = re.findall(r"`([\w_]+)`", original)

    for match in camel_case_matches:
        original = original.replace(match, to_camel_case(match))

    return original
