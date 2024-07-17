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


class APIParameter(TypedDict):
    name: str | None
    type: APIType
    map_entries: list["APIParameter"]
    enum: str | None


def api_parameter(raw: dict) -> APIParameter:
    type = APIType.from_str(raw.get("type", "unknown"))
    assert type is not None, "type should not be None"

    map_entries: dict = raw.get("map_entries", {})
    if type is APIType.MAP and len(map_entries) == 0:
        raise Exception("map entries should not be empty when the type is MAP")

    return {
        "name": raw.get("name"),
        "type": type,
        "map_entries": list(map(api_parameter, map_entries)),
        "enum": raw.get("enum"),
    }


class APIFunction(TypedDict):
    return_types: list[APIParameter]
    func_name: str
    comment: str
    parameters: list[APIParameter]


def api_function(raw: dict) -> APIFunction:
    return {
        "return_types": list(map(api_parameter, raw["return_types"])),
        "func_name": raw["func_name"],
        "comment": raw["comment"],
        "parameters": list(map(api_parameter, raw["parameters"])),
    }


class APIEnum(TypedDict):
    name: str
    values: list[str]


def api_enum(raw: dict) -> APIEnum:
    return {"name": raw["name"], "values": raw["values"]}
