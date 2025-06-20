import os
import subprocess

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

    # Clean previous builds
    os.makedirs("build", exist_ok=True)
    for file in os.listdir("build"):
        os.remove("build/" + file)

    print(f"[+] Starting sequential build of Hybroid")

    for config in _CONFIGS:
        _build_platform(config, os.environ, os.getcwd())

    print("[+] Build job completed!")
