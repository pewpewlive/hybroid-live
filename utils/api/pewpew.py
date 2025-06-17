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
	imports:       make([]Import, 0),
	UsedLibraries:   make([]ast.Library, 0),
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

{enums}

## Functions

{functions}
"""


def generate_docs(pewpew_lib: dict) -> str:
    enums = [api.Enum(enum).generate_docs() for enum in pewpew_lib["enums"]]
    functions = [
        api.Function("pewpew", function).generate_docs("Pewpew")
        for function in pewpew_lib["functions"]
    ]

    return _PEWPEW_DOCS_TEMPLATE.format_map(
        {
            "enums": "\n\n".join(enums),
            "functions": "\n\n".join(functions),
        }
    )
