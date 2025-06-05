from . import types, helpers, mappings


_PEWPEW_API_TEMPLATE = """// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package walker

// AUTO-GENERATED API, DO NOT MANUALLY MODIFY!
var PewpewEanv = &Environment{{
	Name: "Pewpew",
	Scope: Scope{{
		Variables: map[string]*VariableVal{{
            {functions}
        }},
		Tag: &UntaggedTag{{}},
		AliasTypes: make(map[string]*AliasType),
		ConstValues: make(map[string]ast.Node),
	}},
	importedWalkers: make([]*Walker, 0),
	UsedLibraries:   make([]Library, 0),
	Classes: make(map[string]*ClassVal),
	Entities: make(map[string]*EntityVal),
	Enums: map[string]*EnumVal{{
        {enum_descriptions}
	}},
}}


// AUTO-GENERATED API DEFINITION, DO NOT MANUALLY MODIFY!
var PewpewVaariables = 

// AUTO-GENERATED ENUMS, DO NOT MANUALLY MODIFY!
{enums}
"""


def generate_api(pewpew_lib: dict) -> str:
    enums = [types.APIEnum(enum) for enum in pewpew_lib["enums"]]
    functions = [types.APIFunction(function) for function in pewpew_lib["functions"]]

    enums, descriptions = zip(*[enum.generate() for enum in enums])

    return _PEWPEW_API_TEMPLATE.format_map(
        {
            "enums": "\n".join(enums),
            "enum_descriptions": "\n\t".join(descriptions),
            "functions": "\n".join(function.generate() for function in functions),
        }
    )


_PEWPEW_DOCS_TEMPLATE = """---
title: PewPew API
slug: appapi/pewpew
sidebar:
  order: 1
---

<!-- This is an auto-generated file. To modify it, change https://github.com/pewpewlive/hybroid/blob/master/utils/generate_api.py -->

## Enums

%s

## Functions

%s
"""


def _generate_enum_docs(enum: types.APIEnum) -> str:
    enum_template = f"### `{enum.name}`\n"
    enum_template += "".join(
        [f"\n- `{mappings.get(value, helpers.pascal_case)}`" for value in enum.values]
    )

    return enum_template


def _handle_params(parameters: list[types.APIParameter]):
    params = []
    for param in parameters:
        if param.type == types.APIType.MAP:
            params.append(
                "struct {\n  %s\n}" % "\n  ".join(_handle_params(param.map_entries))
            )
        else:
            params.append(param.type.to_str())

    return params


def _generate_function_docs(function: types.APIFunction) -> str:
    processed_name = mappings.get(function.name, helpers.pascal_case)
    return_types = (
        (
            "-> "
            + ", ".join(
                [return_type.type.to_str() for return_type in function.return_types]
            )
        )
        if len(function.return_types) > 0
        else ""
    )
    function_template = f"### `{processed_name}`\n"
    function_template += f"```rs\n{processed_name}({', '.join(_handle_params(function.parameters))}) {return_types}\n```\n"
    function_template += f"{helpers.camel_case_all(function.description)}"

    return function_template


def generate_docs(pewpew_lib: dict) -> str:
    enums = [types.APIEnum(enum) for enum in pewpew_lib["enums"]]
    functions = [types.APIFunction(function) for function in pewpew_lib["functions"]]

    generated_enums = [_generate_enum_docs(enum) for enum in enums]
    generated_functions = [_generate_function_docs(function) for function in functions]
    return _PEWPEW_DOCS_TEMPLATE % (
        "\n\n".join(generated_enums),
        "\n\n".join(generated_functions),
    )
