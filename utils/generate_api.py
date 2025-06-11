import os
import subprocess
import requests

from api import pewpew, fmath


def _generate(api: str | None, lib: str, extension: str, output: str):
    with open(f"{api or ""}{lib}.gen.{extension}", mode="x", encoding="utf-8") as f:
        f.write(output)


def _clean_gen_files(extension: str):
    # Delete already existing generated .gen.go files
    for file in os.listdir(os.getcwd()):
        if file.endswith(f".gen.{extension}"):
            os.remove(file)


if __name__ == "__main__":
    # Get the latest raw docs from the ppl-docs repo
    raw_json = requests.get(
        "https://raw.githubusercontent.com/pewpewlive/ppl-docs/master/raw_documentation.json"
    )
    assert raw_json.status_code == 200, "failed to get raw doc json"

    [pewpew_lib, fmath_lib] = raw_json.json()

    # Generation for API
    # Go to the walker directory where the following steps will be executed
    os.chdir(os.path.dirname(__file__) + "/../walker")
    _clean_gen_files("go")

    # Generate API!
    _generate("api_", "pewpew", "go", pewpew.generate_api(pewpew_lib))
    _generate("api_", "fmath", "go", fmath.generate_api(fmath_lib))

    # Format generated go file
    subprocess.run(
        ["gofmt", "-s", "-w", f"walker/"],
        cwd=os.path.dirname(__file__) + "/..",
    )

    # Mapping generation for API
    # Go to the generator directory where the following steps will be executed
    os.chdir(os.path.dirname(__file__) + "/../generator/mapping")
    _clean_gen_files("go")

    _generate(None, "pewpew", "go", pewpew.generate_api_mapping())
    _generate(None, "fmath", "go", fmath.generate_api_mapping())

    # Format generated go file
    subprocess.run(
        ["gofmt", "-s", "-w", f"generator/mapping/"],
        cwd=os.path.dirname(__file__) + "/..",
    )

    print("[+] API generated!")

    # Generation for docs
    # Go to the docs directory where the following steps will be executed
    os.chdir(os.path.dirname(__file__) + "/../docs/src/content/docs/libraries")
    _clean_gen_files("md")

    # Generate docs!
    _generate("api_", "pewpew", "md", pewpew.generate_docs(pewpew_lib))
    _generate("api_", "fmath", "md", fmath.generate_docs(fmath_lib))
    print("[+] Docs generated!")
