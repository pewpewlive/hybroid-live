import os

from . import types

_FUNCTION_BLACKLIST = [""]


def _generate_pewpew(pewpew_lib):
    output = "// AUTO-GENERATED, DO NOT MANUALLY MODIFY!\n\npackage api\n\n"

    enums = [types.APIEnum(enum) for enum in pewpew_lib["enums"]]
    functions = [types.APIFunction(function) for function in pewpew_lib["functions"]]

    for enum in enums:
        output += f"type {enum.name} int\n\n"
        output += "const (\n"

        for i, variant in enumerate(enum.values):
            if i == 0:
                output += f"  {variant} {enum.name} = iota\n"
            else:
                output += f"  {variant}\n"

        output += ")\n\n"

    return output


def generate_api_for_libs(pewpew_lib, fmath_lib):
    with open("pewpew.gen.go", mode="x") as f:
        pewpew_output = _generate_pewpew(pewpew_lib)
        f.write(pewpew_output)
