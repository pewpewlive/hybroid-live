from enum import Enum

from . import mappings, helpers


class APIType(Enum):
    BOOL = 1
    NUMBER = 2
    FIXED = 3
    ENTITY = 4
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
                return APIType.ENTITY
            case "list":
                return APIType.LIST
            case "map":
                return APIType.MAP
            case "string":
                return APIType.STRING
            case "callback":
                return APIType.CALLBACK

    def to_str(self) -> str:
        match self:
            case APIType.BOOL:
                return "bool"
            case APIType.NUMBER:
                return "number"
            case APIType.FIXED:
                return "fixed"
            case APIType.STRING:
                return "text"
            case APIType.LIST:
                return "list"
            case APIType.MAP:
                return "struct"
            case APIType.CALLBACK:
                return "fn"
            case APIType.ENTITY:
                return "entity"

    def generate(self, param: bool, name: str) -> str:
        # APIType.generate_str cannot generate a map, call APIParameter.generate instead
        assert self is not APIType.MAP, "Cannot generate APIType.MAP"

        # Early return if the type is of raw entity
        if self is APIType.ENTITY:
            return "&RawEntityType{}"

        # The mapping of parameter types
        _PARAM_TYPE_MAPPING = {
            APIType.BOOL: "NewBasicType(ast.Bool)",
            APIType.NUMBER: "NewBasicType(ast.Number)",
            APIType.FIXED: "NewFixedPointType()",
            APIType.LIST: "NewBasicType(ast.List)",
            APIType.STRING: "NewBasicType(ast.String)",
        }

        # The mapping of callback types, not including the `taken_callback` exception,
        # which is a map entry, and not a parameter, so it is dealt with later
        _CALLBACK_TYPE_MAPPING = {
            "AddUpdateCallback": "NewFunctionType(Types{}, Types{})",
            "SetEntityCallback": "NewFunctionType(Types{&RawEntityType{}}, Types{})",
            "ConfigureEntityWallCollision": "NewFunctionType(Types{&RawEntityType{}, NewFixedPointType(), NewFixedPointType(ast.Fixed)}, Types{})",
            "SetEntityPlayerCollision": "NewFunctionType(Types{&RawEntityType{}, NewBasicType(ast.Number), &RawEntityType{}}, Types{})",
            "SetEntityWeaponCollision": 'NewFunctionType(Types{&RawEntityType{}, NewBasicType(ast.Number), NewEnumType("WeaponType")}, Types{NewBasicType(ast.Bool)})',
        }

        if param and self is APIType.CALLBACK:
            return _CALLBACK_TYPE_MAPPING[name]
        elif param:
            return _PARAM_TYPE_MAPPING[self]  # TODO: fghfgh

        match self:
            case APIType.BOOL:
                return "&BoolVal{}"
            case APIType.NUMBER:
                return "&NumberVal{}"
            case APIType.FIXED:
                return "&FixedVal{SpecificType: ast.Fixed}"
            case APIType.STRING:
                return "&StringVal{}"
            case APIType.CALLBACK:
                return "&FunctionVal{Params: Types{&RawEntityType{},&RawEntityType{},&RawEntityType{}}}"
            case _:
                raise Exception(
                    "invalid combination of `self`, `param`, or `name`"
                )  # Should no happen!


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

    def generate(self) -> str:
        return ""


class APIEnum:
    name: str
    values: list[str]

    def __init__(self, raw: dict):
        self.name = raw["name"]
        self.values = raw["values"]

    def generate(self) -> tuple[str, str]:
        ENUM_TEMPLATE = 'var {name} = NewEnumVal("Pewpew", "{name}", false, {values})'
        DESCRIPTION_TEMPLATE = (
            '"{name}": {{Name: "{name}", Value: {name}, IsPub: true, IsConst: true}},'
        )

        return ENUM_TEMPLATE.format_map(
            {
                "name": self.name + "a",
                "values": ", ".join(
                    f'"{mappings.get(value, helpers.pascal_case)}"'
                    for value in self.values
                ),
            }
        ), DESCRIPTION_TEMPLATE.format_map({"name": self.name + "a"})
