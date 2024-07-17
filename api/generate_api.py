import os
import requests

from generator.libraries import generate_api_from_libs

if __name__ == "__main__":
    # Go to the root of this script file to perform further steps
    os.chdir(os.path.dirname(__file__))

    # Delete already existing generated .go files
    for _, _, files in os.walk(os.getcwd()):
        for go_file in [
            filepath for filepath in files if filepath.split(".")[-1] == "go"
        ]:
            os.remove(go_file)

    # Get the latest raw docs from the ppl-docs repo
    raw_json = requests.get(
        "https://raw.githubusercontent.com/pewpewlive/ppl-docs/master/raw_documentation.json"
    )
    assert raw_json.status_code == 200, "failed to get raw doc json"

    [pewpew_lib, fmath_lib] = raw_json.json()

    # Generate!
    generate_api_from_libs(pewpew_lib, fmath_lib)
