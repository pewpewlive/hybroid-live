import re

from . import imports


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


def format_string(
    string: str,
    string_format: list[Format],
    receiver: str,
    accessors: dict[str, str] | None = None,
) -> str:
    if accessors is None:
        accessors = {}

    if len(string_format) == 0:
        return f'"{string}"'

    specifiers = []
    for specifier in string_format:
        if type(specifier) is str:
            specifiers.append(f"{receiver}.{accessors.get(specifier, specifier)}")
        elif type(specifier) is dict:
            specifier, format = dict(specifier).popitem()
            imports.update_imports(format)
            specifiers.append(
                format.format(f"{receiver}.{accessors.get(specifier, specifier)}")
            )

    return f'fmt.Sprintf("{string}", {", ".join(specifiers)})'


def to_param_name(name: str) -> str:
    if not name:
        return name
    return name[0].lower() + name[1:]
