import json
import os
import subprocess

from alerts.alert import Alert
import alerts.imports as imports


def _generate(raw: dict, stage: str) -> str:
    FILE_TEMPLATE = """// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
    
package alerts

import (
  "fmt"
  {imports}
)

// AUTO-GENERATED, DO NOT MANUALLY MODIFY!
{alerts}
"""

    alerts: list[Alert] = []
    id = 1
    for alert in raw:
        alerts.append(Alert(alert, stage, id))
        id += 1

    return FILE_TEMPLATE.format_map(
        {
            "alerts": "\n\n// AUTO-GENERATED, DO NOT MANUALLY MODIFY!\n".join(
                alert.generate() for alert in alerts
            ),
            "imports": imports.get_imports(),
        }
    )


if __name__ == "__main__":
    # Change the directory to the hybroid/alerts folder
    os.chdir(os.path.dirname(__file__) + "/../alerts")

    # Delete already existing generated .gen.go files
    for file in os.listdir(os.getcwd()):
        if file.endswith(".gen.go"):
            os.remove(file)

    # Change the directory to the hybroid/utils/alerts/json folder
    os.chdir(os.path.dirname(__file__) + "/alerts/json")

    # Generate the alerts!
    for file in os.listdir(os.getcwd()):
        if file.endswith(".json"):
            print(f"[-] Generating alerts for {file}")

            # Read the .json file
            with open(file, "r", encoding="utf-8") as f:
                alerts = json.load(f)

            new_filename = file.replace(".json", ".gen.go")

            # Clear extracted imports
            imports.clear()

            # Generate the .gen.go file
            with open(f"../../../alerts/{new_filename}", "x", encoding="utf-8") as f:
                f.write(_generate(alerts, file.split(".")[0].title()))

    # Format generated go file
    subprocess.run(
        ["gofmt", "-s", "-w", f"alerts/"],
        cwd=os.path.dirname(__file__) + "/..",
    )
    print("[+] Alerts generated!")
