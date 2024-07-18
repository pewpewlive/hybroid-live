import os
import requests

from generator.libraries import *


if __name__ == "__main__":
    # Go to the root of this script file to perform further steps
    os.chdir(os.path.dirname(__file__))

    # Delete already existing generated .go files
    for file in os.listdir(os.getcwd()):
        if file.endswith(".go"):
            os.remove(file)

    # Get the latest raw docs from the ppl-docs repo
    raw_json = requests.get(
        "https://raw.githubusercontent.com/pewpewlive/ppl-docs/master/raw_documentation.json"
    )
    assert raw_json.status_code == 200, "failed to get raw doc json"

    [pewpew_lib, fmath_lib] = raw_json.json()

    # Generate!
    generate_api_for_libs(pewpew_lib, fmath_lib)
