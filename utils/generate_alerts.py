import json
import os
import re

_FILE_TEMPLATE = """// AUTO-GENERATED, DO NOT MANUALLY MODIFY!

package alerts

import (
  "fmt"
  "hybroid/tokens"
)

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
"""


def _to_receiver(original: str) -> str:
    # Takes each capital letter in a name and
    # connects them into a single lowercase name.

    # For example: HelloThisISAnExample -> htisae

    return "".join(re.findall(r"[A-Z]", original)).lower()


def _format_string(string: str, string_format: list[str], receiver: str) -> str:
    if len(string_format) == 0:
        return f'"{string}"'

    specifiers = ", ".join(f"{receiver}.{specifier}" for specifier in string_format)

    return f'fmt.Sprintf("{string}", {specifiers})'


class Alert:
    name: str
    receiver: str
    type: str
    stage: str
    params: dict[str, str]
    message: str
    message_format: list[str]
    note: str
    note_format: list[str]

    def __init__(self, raw: dict, stage: str):
        self.name = raw.get("name", None)
        assert self.name is not None, f"Name must not be None, Raw info: {raw}"

        self.receiver = _to_receiver(self.name)

        self.type = raw.get("type", None)
        assert self.type is not None, f"Type must not be None, Raw info: {raw}"

        self.stage = stage

        self.params = {"Token": "tokens.Token", "Location": "tokens.TokenLocation"}
        self.params = self.params | raw.get("params", {})

        self.message = raw.get("message", None)
        assert self.message is not None, f"Message must not be None, Raw info: {raw}"

        self.message_format = raw.get("message_format", [])

        self.note = raw.get(
            "note", ""
        )  # Empty means that the alert will not print out a note

        self.note_format = raw.get("note_format", [])

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
                "GetTokens",
                "",
                "[]tokens.Token",
                f"return []tokens.Token{{{self.receiver}.Token}}",
            ],
            [
                "GetLocations",
                "",
                "[]tokens.TokenLocation",
                f"return []tokens.TokenLocation{{{self.receiver}.Location}}",
            ],
            [
                "GetNote",
                "",
                "string",
                f"return {_format_string(self.note, self.note_format, self.receiver)}",
            ],
            ["GetAlertType", "", "AlertType", f"return {self.type}"],
            ["GetAlertStage", "", "AlertStage", f"return {self.stage}"],
        ]

        for function in alert_functions:
            alert += (
                function_template.format(self.receiver, self.name, *function) + "\n\n"
            )

        return alert


def _generate_alerts(raw: dict, stage: str) -> str:
    alerts = []
    for alert in raw:
        alerts.append(Alert(alert, stage).generate_str())

    return _FILE_TEMPLATE + "// AUTO-GENERATED, DO NOT MANUALLY MODIFY!\n".join(alerts)


def _generate_file(filename: str):
    # Read the .json file
    with open(filename, "r", encoding="utf-8") as f:
        alerts = json.load(f)

    new_filename = filename.replace(".json", ".gen.go")

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
