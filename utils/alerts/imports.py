_POSSIBLE_PACKAGES = ["strings", "hybroid/ast"]
_DEFAULT = ['"fmt"', '"hybroid/core"']
_imports = set(_DEFAULT)


def update_imports(string: str):
    for package in _POSSIBLE_PACKAGES:
        package = package[8:] if "hybroid/" in package else package
        if package in string:
            _imports.add(f'"{package}"')


def get_imports() -> str:
    return "\n  ".join(_imports)


def clear():
    _imports = set(_DEFAULT)
