import json
import os
import re

_FILE_TEMPLATE = """// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
  "fmt"
  {}
)

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
"""

_POSSIBLE_PACKAGES = ["strings"]
_imports = set()


def _to_receiver(original: str) -> str:
    # Takes each capital letter in a name and
    # connects them into a single lowercase name.

    # For example: HelloThisISAnExample -> htisae

    return "".join(re.findall(r"[A-Z]", original)).lower()


def _extract_imports(string: str):
    for package in _POSSIBLE_PACKAGES:
        if package in string:
            _imports.add(f'"{package}"')


type Format = dict[str, str] | str


def _format_string(string: str, string_format: list[Format], receiver: str) -> str:
    if len(string_format) == 0:
        return f'"{string}"'

    specifiers = []
    for specifier in string_format:
        if type(specifier) is str:
            specifiers.append(f"{receiver}.{specifier}")
        elif type(specifier) is dict:
            specifier, format = dict(specifier).popitem()
            _extract_imports(format)
            specifiers.append(format.format(f"{receiver}.{specifier}"))

    return f'fmt.Sprintf("{string}", {", ".join(specifiers)})'


class Alert:
    name: str
    receiver: str
    type: str
    stage: str
    params: dict[str, str]
    message: str
    message_format: list[Format]
    note: str
    note_format: list[Format]
    id: int

    def __init__(self, raw: dict, stage: str, id: int):
        name = raw.get("name")
        if name is None:
            raise ValueError(f"Name must not be None, Raw info: {raw}")
        self.name = name

        self.receiver = _to_receiver(self.name)

        type = raw.get("type")
        if type is None:
            raise ValueError(f"Type must not be None, Raw info: {raw}")
        self.type = type

        self.stage = stage

        self.params = {"Specifier": "Snippet"}
        self.params = self.params | raw.get("params", {})

        message = raw.get("message")
        if message is None:
            raise ValueError(f"Message must not be None, Raw info: {raw}")
        self.message = message

        self.message_format = raw.get("message_format", [])

        self.note = raw.get(
            "note", ""
        )  # Empty means that the alert will not print out a note

        self.note_format = raw.get("note_format", [])

        self.id = id

    def generate_str(self) -> str:
        type_template = "type {} struct {{\n  {}\n}}"
        function_template = "func ({} *{}) {}({}) {} {{\n  {}\n}}"

        alert = (
            type_template.format(
                self.name,
                "\n  ".join(f"{field} {type}" for field, type in self.params.items()),
            )
            + "\n\n"
        )

        alert_functions = [
            [
                "GetMessage",
                "",
                "string",
                f"return {_format_string(self.message, self.message_format, self.receiver)}",
            ],
            [
                "GetSpecifier",
                "",
                "Snippet",
                f"return {self.receiver}.Specifier",
            ],
            [
                "GetNote",
                "",
                "string",
                f"return {_format_string(self.note, self.note_format, self.receiver)}",
            ],
            [
                "GetID",
                "",
                "string",
                'return "hyb{:03d}{}"'.format(self.id, self.stage[0]),
            ],
            ["GetAlertType", "", "Type", f"return {self.type}"],
        ]

        for function in alert_functions:
            alert += (
                function_template.format(self.receiver, self.name, *function) + "\n\n"
            )

        return alert


def _generate_alerts(raw: dict, stage: str) -> str:
    alerts = []
    id = 1
    for alert in raw:
        alerts.append(Alert(alert, stage, id).generate_str())
        id += 1

    return _FILE_TEMPLATE.format(
        "\n  ".join(_imports)
    ) + "// AUTO-GENERATED, DO NOT MANUALLY MODIFY!\n".join(alerts)


def _generate_file(filename: str):
    # Read the .json file
    with open(filename, "r", encoding="utf-8") as f:
        alerts = json.load(f)

    new_filename = filename.replace(".json", ".gen.go")

    # Clear extracted imports
    _imports.clear()

    # Generate the .gen.go file
    with open(f"../../alerts/{new_filename}", "x", encoding="utf-8") as f:
        f.write(_generate_alerts(alerts, filename.split(".")[0].title()))


if __name__ == "__main__":
    # Change the directory to the hybroid/alerts folder
    os.chdir(os.path.dirname(__file__) + "/../alerts")

    # Delete already existing generated .gen.go files
    for file in os.listdir(os.getcwd()):
        if file.endswith(".gen.go"):
            os.remove(file)

    # Change the directory to the hybroid/utils/alerts folder
    os.chdir(os.path.dirname(__file__) + "/alerts")

    # Generate the alerts!
    for file in os.listdir(os.getcwd()):
        if file.endswith(".json"):
            print(f"[-] Generating alerts for {file}")
            _generate_file(file)

    print("[+] Alerts generated!")
