from . import types, helpers

_API_MAPPING = {}

_PARAM_MAPPING = {}


def generate_api(fmath_lib: dict) -> str:
    return "package walker"


_FMATH_DOCS_TEMPLATE = """---
title: Fmath API
slug: appapi/fmath
sidebar:
  order: 2
---

<!-- This is an auto-generated file. To modify it, change utils/generate_api.py in Hybroid's repository. -->

## Functions

%s
"""


def _generate_function_docs(function: types.APIFunction) -> str:
    processed_name = _API_MAPPING.get(function.name, helpers.pascal_case(function.name))
    function_template = f"### `{processed_name}`\n"
    function_template += f"```rs\n{processed_name}({', '.join([_TYPE_MAPPING.get(param.type, 'unknown') + ' ' + _PARAM_MAPPING.get(param.name, helpers.camel_case(param.name)) for param in function.parameters])}) { ('-> ' + ', '.join([_TYPE_MAPPING.get(return_type.type, 'unknown') for return_type in function.return_types])) if len(function.return_types) > 0 else ''}\n```\n"
    function_template += f"{function.description}"

    return function_template


def generate_docs(fmath_lib: dict) -> str:
    functions = [types.APIFunction(function) for function in fmath_lib["functions"]]
    generated_functions = [_generate_function_docs(function) for function in functions]

    return _FMATH_DOCS_TEMPLATE % ("\n\n".join(generated_functions))
