import os

from .types import api_function

_function_blacklist = [""]


def _generate_pewpew(pewpew_lib):
    output = "// AUTO-GENERATED, DO NOT MANUALLY MODIFY!\n\npackage api\n\n/*\n"

    for function in pewpew_lib["functions"]:
        if function["func_name"] in _function_blacklist:
            continue

        fn = api_function(function)
        output = output + f"func {fn["func_name"]}()\n"

    return output + "*/\n"


def generate_api_from_libs(pewpew_lib, fmath_lib):
    with open("pewpew.gen.go", mode="x") as f:
        pewpew_output = _generate_pewpew(pewpew_lib)
        f.write(pewpew_output)
