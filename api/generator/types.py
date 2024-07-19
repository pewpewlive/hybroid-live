from enum import Enum
from typing import TypedDict


class APIType(Enum):
    BOOL = 1
    NUMBER = 2
    FIXED = 3
    RAW_ENTITY = 4
    LIST = 5
    MAP = 6
    STRING = 7
    CALLBACK = 8

    @staticmethod
    def from_str(string: str):
        match string.lower():
            case "boolean":
                return APIType.BOOL
            case "int32" | "int64":
                return APIType.NUMBER
            case "fixedpoint":
                return APIType.FIXED
            case "entityid":
                return APIType.RAW_ENTITY
            case "list":
                return APIType.LIST
            case "map":
                return APIType.MAP
            case "string":
                return APIType.STRING
            case "callback":
                return APIType.CALLBACK


class APIParameter:
    name: str | None
    type: APIType
    map_entries: list["APIParameter"]
    enum: str | None

    def __init__(self, raw: dict):
        type = APIType.from_str(raw.get("type", "unknown"))
        assert type is not None, "type should not be None"

        map_entries: dict = raw.get("map_entries", {})
        if type is APIType.MAP and len(map_entries) == 0:
            raise Exception("map entries should not be empty when the type is MAP")

        self.name = raw.get("name")
        self.type = type
        self.map_entries = [APIParameter(entry) for entry in map_entries]
        self.enum = raw.get("enum")


class APIFunction:
    name: str
    description: str
    parameters: list[APIParameter]
    return_types: list[APIParameter]

    def __init__(self, raw: dict):
        self.name = raw["func_name"]
        self.description = raw["comment"]
        self.parameters = [APIParameter(param) for param in raw["parameters"]]
        self.return_types = [APIParameter(type) for type in raw["return_types"]]


class APIEnum:
    name: str
    values: list[str]

    def __init__(self, raw: dict):
        self.name = raw["name"]
        self.values = raw["values"]
