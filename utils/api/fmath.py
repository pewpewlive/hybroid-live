from . import api


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
	UsedLibraries: make([]Library, 0),
	Classes: make(map[string]*ClassVal),
	Entities: make(map[string]*EntityVal),
	Enums: make(map[string]*EnumVal),
}}
"""


def generate_api(fmath_lib: dict) -> str:
    functions = [
        api.Function(function).generate("Fmath") for function in fmath_lib["functions"]
    ]

    return _FMATH_API_TEMPLATE.format(",\n".join(functions))


_FMATH_DOCS_TEMPLATE = """---
title: Fmath API
slug: libraries/fmath
sidebar:
  order: 2
---

<!-- This is an auto-generated file. To modify it, change utils/generate_api.py in Hybroid's repository. -->

## Functions

%s
"""


# def _generate_function_docs(function: types.APIFunction) -> str:
#     processed_name = mappings.get(function.name, helpers.pascal_case)
#     function_template = f"### `{processed_name}`\n"
#     function_template += f"```rs\n{processed_name}({', '.join([_TYPE_MAPPING.get(param.type, 'unknown') + ' ' +  mappings.get(param.name, helpers.camel_case) for param in function.parameters])}) { ('-> ' + ', '.join([_TYPE_MAPPING.get(return_type.type, 'unknown') for return_type in function.return_types])) if len(function.return_types) > 0 else ''}\n```\n"
#     function_template += f"{function.description}"

#     return function_template


def generate_docs(fmath_lib: dict) -> str:
    return ""


#     functions = [types.APIFunction(function) for function in fmath_lib["functions"]]
#     generated_functions = [_generate_function_docs(function) for function in functions]

#     return _FMATH_DOCS_TEMPLATE % ("\n\n".join(generated_functions))
