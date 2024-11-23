import os
import requests

from generator.api_generator import *


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

    # Delete already existing generated .gen.go files
    for file in os.listdir(os.getcwd()):
        if file.endswith(".gen.go"):
            os.remove(file)

    # Generate API!
    generate_api_for_libs(pewpew_lib, fmath_lib)
    print("[+] API generated!")

    # Generation for docs
    # Go to the docs directory where the following steps will be executed
    os.chdir(os.getcwd() + "/../docs/src/content/docs/appapi")

    # Delete already existing generated .gen.md files
    for file in os.listdir(os.getcwd()):
        if file.endswith(".gen.md"):
            os.remove(file)

    # Generate docs!
    generate_docs_for_libs(pewpew_lib, fmath_lib)
    print("[+] Docs generated!")
