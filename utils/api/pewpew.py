from . import api, helpers, mappings, types


_PEWPEW_API_TEMPLATE = """// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
package walker

import "hybroid/ast"

// AUTO-GENERATED API, DO NOT MANUALLY MODIFY!
var PewpewAPI = &Environment{{
	Name: "Pewpew",
	Scope: Scope{{
		Variables: map[string]*VariableVal{{
            {},
        }},
		Tag: &UntaggedTag{{}},
		AliasTypes: make(map[string]*AliasType),
		ConstValues: make(map[string]ast.Node),
	}},
	importedWalkers: make([]*Walker, 0),
	UsedLibraries:   make([]Library, 0),
	Classes: make(map[string]*ClassVal),
	Entities: make(map[string]*EntityVal),
	Enums: {},
}}
"""

_PEWPEW_API_MAP_TEMPLATE = """// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
package mapping

// AUTO-GENERATED ENUMS, DO NOT MANUALLY MODIFY!
var PewpewEnums = map[string]map[string]string{{
    {enums}
}}

// AUTO-GENERATED VARIABLES, DO NOT MANUALLY MODIFY!
var PewpewVariables = map[string]string{{
    {functions},

    {enum_names},
}}
"""


def generate_api(pewpew_lib: dict) -> str:
    enums = [api.Enum(enum).generate() for enum in pewpew_lib["enums"]]
    functions = [
        api.Function("pewpew", function).generate("Pewpew")
        for function in pewpew_lib["functions"]
    ]

    return _PEWPEW_API_TEMPLATE.format(
        ",\n".join(functions), f"map[string]*EnumVal{{\n{",\n".join(enums)}}}"
    )


def generate_api_mapping() -> str:
    functions, enums = mappings.inverse_mappings("pewpew")
    functions = functions["pewpew"]

    ENUM_TEMPLATE = '"{name}": {{\n{values},\n}},'

    generated_enums = ""
    for hyb, enum in enums.items():
        _, enum = enum
        generated_enums += ENUM_TEMPLATE.format_map(
            {
                "name": hyb,
                "values": ",".join(f'"{hyb}":"{ppl}"' for hyb, ppl in enum.items()),
            }
        )

    return _PEWPEW_API_MAP_TEMPLATE.format_map(
        {
            "enums": generated_enums,
            "functions": ",\n".join(
                f'"{hyb}":"{ppl}"' for hyb, ppl in functions.items()
            ),
            "enum_names": ",\n".join(
                f'"{hyb}":"{enum[0]}"' for hyb, enum in enums.items()
            ),
        }
    )


_PEWPEW_DOCS_TEMPLATE = """---
title: PewPew API
slug: libraries/pewpew
sidebar:
  order: 1
---

<!-- This is an auto-generated file. To modify it, change https://github.com/pewpewlive/hybroid/blob/master/utils/generate_api.py -->

## Enums

%s

## Functions

%s
"""


def _generate_enum_docs(enum: api.Enum) -> str:
    enum_template = f"### `{enum.name}`\n"
    enum_template += "".join(
        [
            f"\n- `{mappings.get_function("pewpew", value, helpers.pascal_case)}`"
            for value in enum.variants
        ]
    )

    return enum_template


def _handle_params(parameters: list[api.Value]):
    params = []
    for param in parameters:
        if param.type == types.Type.MAP:
            params.append(
                "struct {\n  %s\n}" % "\n  ".join(_handle_params(param.map_entries))
            )
        else:
            pass
            # params.append(param.type.generate())

    return params


def _generate_function_docs(function: api.Function) -> str:
    processed_name = mappings.get_function("pewpew", function.name, helpers.pascal_case)
    # returns = (
    #     (
    #         "-> "
    #         + ", ".join(
    #             [return_type.generate() for return_type in function.returns]
    #         )
    #     )
    #     if len(function.returns) > 0
    #     else ""
    # )
    function_template = f"### `{processed_name}`\n"
    # function_template += f"```rs\n{processed_name}({', '.join(_handle_params(function.parameters))}) {returns}\n```\n"
    # function_template += f"{helpers.camel_case_all(function.description)}"

    return function_template


def generate_docs(pewpew_lib: dict) -> str:
    enums = [api.Enum(enum) for enum in pewpew_lib["enums"]]
    functions = [
        api.Function("pewpew", function) for function in pewpew_lib["functions"]
    ]

    generated_enums = [_generate_enum_docs(enum) for enum in enums]
    generated_functions = [_generate_function_docs(function) for function in functions]
    return _PEWPEW_DOCS_TEMPLATE % (
        "\n\n".join(generated_enums),
        "\n\n".join(generated_functions),
    )
