from . import api, mappings


_FMATH_API_TEMPLATE = """// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
package walker

import "hybroid/ast"

// AUTO-GENERATED API, DO NOT MANUALLY MODIFY!
var FmathAPI = &Environment{{
	Name: "Fmath",
	Scope: Scope{{
		Variables: map[string]*VariableVal{{
            {},
        }},
		Tag: &UntaggedTag{{}},
		AliasTypes: make(map[string]*AliasType),
		ConstValues: make(map[string]ast.Node),
	}},
	importedWalkers: make([]*Walker, 0),
	UsedLibraries: make([]ast.Library, 0),
	Classes: make(map[string]*ClassVal),
	Entities: make(map[string]*EntityVal),
	Enums: make(map[string]*EnumVal),
}}
"""

_FMATH_API_MAP_TEMPLATE = """// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
package mapping

// AUTO-GENERATED VARIABLES, DO NOT MANUALLY MODIFY!
var FmathVariables = map[string]string{{
    {functions},
}}
"""


def generate_api(fmath_lib: dict) -> str:
    functions = [
        api.Function("fmath", function).generate("Fmath")
        for function in fmath_lib["functions"]
    ]

    return _FMATH_API_TEMPLATE.format(",\n".join(functions))


def generate_api_mapping() -> str:
    functions, _ = mappings.inverse_mappings("fmath")
    functions = functions["fmath"]

    return _FMATH_API_MAP_TEMPLATE.format_map(
        {"functions": ",\n".join(f'"{hyb}":"{ppl}"' for hyb, ppl in functions.items())}
    )


_FMATH_DOCS_TEMPLATE = """---
title: Fmath API
slug: libraries/fmath
sidebar:
  order: 2
---

<!-- This is an auto-generated file. To modify it, change utils/generate_api.py in Hybroid's repository. -->

## Functions

{functions}
"""


def generate_docs(fmath_lib: dict) -> str:
    functions = [
        api.Function("fmath", function).generate_docs("Fmath")
        for function in fmath_lib["functions"]
    ]

    return _FMATH_DOCS_TEMPLATE.format_map({"functions": "\n\n".join(functions)})
