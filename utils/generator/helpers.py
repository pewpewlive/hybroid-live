import re


def convert_to_pascal_case(original: str) -> str:
    return original.title().replace("_", "")


def convert_to_camel_case(original: str) -> str:
    titled = original.title().replace("_", "")
    return titled[0].lower() + titled[1:]


def find_snake_in_docs_and_convert(original: str) -> str:
    camel_case_matches = re.findall(r"`([\w_]+)`", original)

    for match in camel_case_matches:
        original = original.replace(match, convert_to_camel_case(match))

    return original
