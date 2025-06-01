_POSSIBLE_PACKAGES = ["strings", "hybroid/ast"]
_imports = set()


def update_imports(string: str):
    for package in _POSSIBLE_PACKAGES:
        package = package[8:] if "hybroid/" in package else package
        if package in string:
            _imports.add(f'"{package}"')


def get_imports() -> str:
    return "\n  ".join(_imports)


def clear():
    _imports.clear()
