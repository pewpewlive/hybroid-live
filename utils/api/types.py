import enum


class Type(enum.Enum):
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
                return Type.BOOL
            case "int32" | "int64":
                return Type.NUMBER
            case "fixedpoint":
                return Type.FIXED
            case "entityid":
                return Type.ENTITY
            case "list":
                return Type.LIST
            case "map":
                return Type.MAP
            case "string":
                return Type.STRING
            case "callback":
                return Type.CALLBACK

    def to_str(self) -> str:
        match self:
            case Type.BOOL:
                return "bool"
            case Type.NUMBER:
                return "number"
            case Type.FIXED:
                return "fixed"
            case Type.STRING:
                return "text"
            case Type.LIST:
                return "list"
            case Type.MAP:
                return "struct"
            case Type.CALLBACK:
                return "fn"
            case Type.ENTITY:
                return "entity"


class Function:
    name: str
    description: str
    parameters: list
    returns: list


class Enum:
    name: str
    variants: list[str]
