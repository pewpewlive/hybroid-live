import os

from . import types, pewpew, fmath


def _generate_api(lib, output):
    with open(f"api_{lib}.gen.go", mode="x", encoding="utf-8") as f:
        f.write(output)


def _generate_docs(lib, output):
    with open(f"{lib}.gen.md", mode="x", encoding="utf-8") as f:
        f.write(output)


def generate_api_for_libs(pewpew_lib, fmath_lib):
    _generate_api("pewpew", pewpew.generate_api(pewpew_lib))
    _generate_api("fmath", fmath.generate_api(fmath_lib))


def generate_docs_for_libs(pewpew_lib, fmath_lib):
    _generate_docs("pewpew", pewpew.generate_docs(pewpew_lib))
    _generate_docs("fmath", fmath.generate_docs(fmath_lib))
