import os
import subprocess
import sys

_CONFIGS = [
    {"env": {"GOOS": "windows", "GOARCH": "arm64"}, "name": "windows-arm64"},
    {"env": {"GOOS": "darwin", "GOARCH": "arm64"}, "name": "macos-arm64"},
    {"env": {"GOOS": "linux", "GOARCH": "arm64"}, "name": "linux-arm64"},
    {"env": {"GOOS": "windows", "GOARCH": "amd64"}, "name": "windows-x86_64"},
    {"env": {"GOOS": "darwin", "GOARCH": "amd64"}, "name": "macos-x86_64"},
    {"env": {"GOOS": "linux", "GOARCH": "amd64"}, "name": "linux-x86_64"},
    {"env": {"GOOS": "windows", "GOARCH": "386"}, "name": "windows-x86"},
    {"env": {"GOOS": "linux", "GOARCH": "386"}, "name": "linux-x86"},
    {"env": {"GOOS": "linux", "GOARCH": "arm"}, "name": "linux-arm"},
    {"env": {"GOOS": "js", "GOARCH": "wasm"}, "name": "wasm"},
]


def _build_platform(platform, env, cwd):
    print(f"[-] Building {platform['name']}")
    filename = "hybroid-" + platform["name"]

    if platform["env"]["GOOS"] == "js":
        filename += ".wasm"
    elif platform["env"]["GOOS"] == "windows":
        filename += ".exe"

    # Build!
    subprocess.run(
        ["go", "build", "-o", f"./build/{filename}", "hybroid"],
        check=True,
        env=env | platform["env"],
        cwd=cwd,
    )


if __name__ == "__main__":
    # Go to the hybroid root directory to build
    os.chdir(os.path.dirname(__file__) + "/..")
    os.makedirs("build", exist_ok=True)

    target = None
    if len(sys.argv) > 1:
        target = sys.argv[1]
        valid_targets = [c["name"] for c in _CONFIGS]
        if target not in valid_targets:
            print(f"Error: Target '{target}' not found.")
            print(f"Available targets: {', '.join(valid_targets)}")
            sys.exit(1)
    
    # Clean previous builds only if building everything
    if not target:
        for file in os.listdir("build"):
            os.remove("build/" + file)

    print(f"[+] Starting build of Hybroid")

    for config in _CONFIGS:
        if target and config["name"] != target:
            continue
        _build_platform(config, os.environ, os.getcwd())

    print("[+] Build job completed!")
