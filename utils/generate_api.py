import os
import subprocess
import requests

from api import pewpew, fmath


def _generate(lib, extension, output):
    with open(f"api_{lib}.gen.{extension}", mode="x", encoding="utf-8") as f:
        f.write(output)


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
    _generate("pewpew", "go", pewpew.generate_api(pewpew_lib))
    _generate("fmath", "go", fmath.generate_api(fmath_lib))

    # Format generated go file
    subprocess.run(
        ["gofmt", "-s", "-w", f"alerts/"],
        cwd=os.path.dirname(__file__) + "/..",
    )

    print("[+] API generated!")

    # Generation for docs
    # Go to the docs directory where the following steps will be executed
    os.chdir(os.getcwd() + "/../docs/src/content/docs/appapi")

    # Delete already existing generated .gen.md files
    for file in os.listdir(os.getcwd()):
        if file.endswith(".gen.md"):
            os.remove(file)

    # Generate docs!
    _generate("pewpew", "md", pewpew.generate_docs(pewpew_lib))
    _generate("fmath", "md", fmath.generate_docs(fmath_lib))
    print("[+] Docs generated!")
