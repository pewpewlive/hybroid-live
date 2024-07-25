from . import types, helpers


_PEWPEW_API_TEMPLATE = """// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package walker

import (
	"hybroid/ast"
)

// AUTO-GENERATED ENUMS, DO NOT MANUALLY MODIFY!
%s

// AUTO-GENERATED API DEFINITION, DO NOT MANUALLY MODIFY!
var pewpewAPI = map[string]*VariableVal{
  // Enums
%s

  // Functions
%s
}

// AUTO-GENERATED API, DO NOT MANUALLY MODIFY!
var PewpewEnv = &Environment{
	Name: "Pewpew",
	Scope: Scope{
		Variables: pewpewAPI,
		Tag: &UntaggedTag{},
	},
	Structs: make(map[string]*StructVal),
	Entities: make(map[string]*EntityVal),
	CustomTypes: make(map[string]*CustomType),
}
"""

# The API mapping dictionary holds the initial mapping of lua enum variants/functions to Hybroid
# It also get populated with the converted case
_API_MAPPING = {
    # EntityType
    "BAF": "YellowBaf",
    # MothershipType
    "THREE_CORNERS": "Triangle",
    "FOUR_CORNERS": "Square",
    "FIVE_CORNERS": "Pentagon",
    "SIX_CORNERS": "Hexagon",
    "SEVEN_CORNERS": "Heptagon",
    # CannonFrequency
    "FREQ_7_5": "Freq7_5",
    # AsteroidSize
    "VERY_LARGE": "Huge",
}


def _generate_enum(enum: types.APIEnum) -> str:
    enum_template = f'var {enum.name} = NewEnumVal("{enum.name}", false,\n'
    enum_template += ",\n".join(
        [
            f'\t"{_API_MAPPING.get(value, helpers.convert_case(value))}"'
            for value in enum.values
        ]
    )
    enum_template += ",\n)"

    return enum_template


def _generate_enum_description(enum: types.APIEnum) -> str:
    enum_description = ",\n\t\t".join(
        [
            f'\n\t\tName: "{enum.name}"',
            f"Value: {enum.name}",
            "IsLocal: false",
            "IsConst: true",
        ]
    )

    return f'\t"{enum.name}": {{{enum_description},\n\t}},'


def _generate_function(enum: types.APIFunction) -> str:
    return ""


def generate_api(pewpew_lib: dict) -> str:
    enums = [types.APIEnum(enum) for enum in pewpew_lib["enums"]]
    functions = [types.APIFunction(function) for function in pewpew_lib["functions"]]

    generated_enums = [_generate_enum(enum) for enum in enums]
    generated_enum_descriptions = [_generate_enum_description(enum) for enum in enums]
    generated_functions = [_generate_function(function) for function in functions]

    return _PEWPEW_API_TEMPLATE % (
        "\n\n".join(generated_enums),
        "\n".join(generated_enum_descriptions),
        "\n".join(generated_functions),
    )


_PEWPEW_DOCS_TEMPLATE = """---
title: PewPew API
slug: appapi/pewpew
sidebar:
  order: 1
---

<!-- This is an auto-generated file. To modify it, change utils/generate_api.py in Hybroid's repository. -->

## Enums

%s

## Functions

%s
"""


def _generate_enum_docs(enum: types.APIEnum) -> str:
    enum_template = f"### `{enum.name}`\n"
    enum_template += "".join(
        [
            f"\n- `{_API_MAPPING.get(value, helpers.convert_case(value))}`"
            for value in enum.values
        ]
    )

    return enum_template


def generate_docs(pewpew_lib: dict) -> str:
    enums = [types.APIEnum(enum) for enum in pewpew_lib["enums"]]
    functions = [types.APIFunction(function) for function in pewpew_lib["functions"]]

    generated_enums = [_generate_enum_docs(enum) for enum in enums]
    return _PEWPEW_DOCS_TEMPLATE % ("\n\n".join(generated_enums), "To Be Added.")
