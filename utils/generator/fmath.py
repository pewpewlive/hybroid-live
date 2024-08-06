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


def generate_docs(fmath_lib: dict) -> str:
    return _FMATH_DOCS_TEMPLATE % ("To Be Added.")
