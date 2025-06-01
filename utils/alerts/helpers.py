import re

from . import imports

_imports = set()


def to_receiver(original: str) -> str:
    # Takes each capital letter in a name and
    # connects them into a single lowercase name.

    # For example: HelloThisISAnExample -> htisae

    receiver = "".join(re.findall(r"[A-Z]", original)).lower()

    # Exception to make sure it does not use Go's 'if' keyword
    if receiver == "if":
        receiver = "_if"

    return receiver


type Format = dict[str, str] | str


def format_string(string: str, string_format: list[Format], receiver: str) -> str:
    if len(string_format) == 0:
        return f'"{string}"'

    specifiers = []
    for specifier in string_format:
        if type(specifier) is str:
            specifiers.append(f"{receiver}.{specifier}")
        elif type(specifier) is dict:
            specifier, format = dict(specifier).popitem()
            imports.update_imports(format)
            specifiers.append(format.format(f"{receiver}.{specifier}"))

    return f'fmt.Sprintf("{string}", {", ".join(specifiers)})'
