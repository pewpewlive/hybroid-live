from . import types, mappings, helpers


class Type:
    type: types.Type

    def __init__(self, type: str):
        api_type = types.Type.from_str(type)
        assert api_type is not None, "type cannot be None"
        self.type = api_type

    def __eq__(self, other: types.Type) -> bool:
        return self.type == other

    def __repr__(self):
        return f"Type({self.type})"

    def generate(self, param: bool, name: str) -> str:
        # Type.generate cannot generate a map, call Parameter.generate instead
        assert self.type is not types.Type.MAP, "Cannot generate APIType.MAP"

        # Early return if the type is of raw entity
        if self.type is types.Type.ENTITY:
            return "&RawEntityType{}"

        # The mapping of parameter types
        PARAM_TYPES = {
            types.Type.BOOL: "NewBasicType(ast.Bool)",
            types.Type.NUMBER: "NewBasicType(ast.Number)",
            types.Type.FIXED: "NewFixedPointType()",
            types.Type.LIST: "NewBasicType(ast.List)",
            types.Type.STRING: "NewBasicType(ast.Text)",
        }

        if self.type in (types.Type.LIST, types.Type.STRING):
            match name:
                case "GetEntitiesInRadius" | "GetAllEntities":
                    return "NewWrapperType(NewBasicType(ast.List), &RawEntityType{})"
                case "SetEntityFlippingMeshes" | "SetEntityMesh":
                    return "NewPathType(ast.MeshEnv)"
                case "PlaySound" | "PlayAmbientSound":
                    return "NewPathType(ast.SoundEnv)"

        if self.type is types.Type.NUMBER and name == "GetEntityType":
            return 'NewEnumType("Pewpew", "EntityType")'

        # The mapping of callback types, not including the `taken_callback` exception,
        # which is a map entry, and not a parameter, so it is dealt with later
        CALLBACK_TYPES = {
            "AddUpdateCallback": "NewFunctionType([]Type{},[]Type{})",
            "SetEntityUpdateCallback": "NewFunctionType([]Type{&RawEntityType{}},[]Type{})",
            "SetEntityWallCollision": "NewFunctionType([]Type{&RawEntityType{},NewFixedPointType(),NewFixedPointType()},[]Type{})",
            "SetEntityPlayerCollision": "NewFunctionType([]Type{&RawEntityType{},NewBasicType(ast.Number),&RawEntityType{}}, []Type{})",
            "SetEntityWeaponCollision": 'NewFunctionType([]Type{&RawEntityType{},NewBasicType(ast.Number),NewEnumType("Pewpew","WeaponType")},[]Type{NewBasicType(ast.Bool)})',
        }

        if param and self.type is types.Type.CALLBACK:
            return CALLBACK_TYPES[name]
        elif param:
            return PARAM_TYPES[self.type]

        match self.type:
            case types.Type.BOOL:
                return "&BoolVal{}"
            case types.Type.NUMBER:
                return "&NumberVal{}"
            case types.Type.FIXED:
                return "&FixedVal{}"
            case types.Type.STRING:
                return "&StringVal{}"
            case types.Type.CALLBACK:
                return "&FunctionVal{Params:[]Type{&RawEntityType{},NewBasicType(ast.Number),&RawEntityType{}}}"
            case _:
                raise Exception(
                    "invalid combination of `self`, `param`, or `name`"
                )  # Should not happen!


class Value:
    name: str
    type: Type
    map_entries: list["Value"]
    enum: str | None

    def __init__(self, raw: dict):
        type = Type(raw.get("type", "unknown"))

        map_entries = raw.get("map_entries", {})
        if type is types.Type.MAP and len(map_entries) == 0:
            raise Exception("map entries should not be empty when the type is MAP")

        self.name = raw.get("name", "")
        self.type = type
        self.map_entries = [Value(entry) for entry in map_entries]

        self.enum = raw.get("enum", None)
        if self.enum is not None:
            self.enum = mappings.get_enum(self.enum, None, helpers.pewpew_conversion)

    def __repr__(self):
        return f"Value({self.name}, {self.type}, {self.map_entries}, {self.enum})"

    def generate(self, lib_name: str, func_name: str) -> str:
        if self.enum is not None:
            return f'NewEnumType("{lib_name}", "{self.enum}")'

        if len(self.map_entries) != 0:
            MAP_TEMPLATE = "NewStructType([]StructField{%s})"

            fields = []
            for entry in self.map_entries:
                if entry.enum is not None:
                    entry_value = f'NewEnumVal("Pewpew", "{entry.enum}", true)'
                else:
                    entry_value = entry.type.generate(False, entry.name)

                fields.append((entry.name, entry_value))

            return MAP_TEMPLATE % (
                ",".join(
                    f'NewStructField("{name}", {value}, true)' for name, value in fields
                )
            )

        return self.type.generate(True, func_name)

    def generate_docs(self) -> str:
        if self.enum is not None:
            return f"{self.enum} {self.name}"
        elif len(self.map_entries) != 0:
            entries = ",\n  ".join(entry.generate_docs() for entry in self.map_entries)
            return f"struct {{\n  {entries}\n}}"
        else:
            return f"{self.type.type.to_str()} {self.name}"


class Function(types.Function):
    def __init__(self, lib: str, raw: dict):
        self.name = mappings.get_function(
            lib, raw["func_name"], helpers.pewpew_conversion
        )
        self.description = raw["comment"]
        self.parameters = [Value(param) for param in raw["parameters"]]
        self.returns = [Value(type) for type in raw["return_types"]]

    def generate(self, lib_name: str) -> str:
        VALUE_TEMPLATE = "NewFunction({params})"

        value_args = {
            "params": ",".join(
                param.generate(lib_name, self.name) for param in self.parameters
            ),
        }

        if len(self.returns) != 0:
            VALUE_TEMPLATE += ".WithReturns({returns})"
            value_args |= {
                "returns": ",".join(
                    ret.generate(lib_name, self.name) for ret in self.returns
                )
            }

        FUNCTION_TEMPLATE = '"{name}":{{\nName:"{name}",Value:{value},IsPub:true,\n}}'

        return FUNCTION_TEMPLATE.format_map(
            {"name": self.name, "value": VALUE_TEMPLATE.format_map(value_args)}
        )

    def generate_docs(self, lib_name: str) -> str:
        FUNCTION_TEMPLATE = (
            "### `{name}`\n\n```rs\n{name}({params}){returns}\n```\n{description}\n"
        )

        return FUNCTION_TEMPLATE.format_map(
            {
                "name": self.name,
                "params": ", ".join(param.generate_docs() for param in self.parameters),
                "returns": (
                    " -> " + ", ".join(ret.generate_docs() for ret in self.returns)
                    if len(self.returns) != 0
                    else ""
                ),
                "description": self.description,
            }
        )


class Enum(types.Enum):
    def __init__(self, raw: dict):
        self.name = mappings.get_enum(raw["name"], None, helpers.pewpew_conversion)
        self.variants = []
        for variant in raw["values"]:
            self.variants.append(
                mappings.get_enum(self.name, variant, helpers.pewpew_conversion)
            )

    def generate(self) -> str:
        ENUM_TEMPLATE = '"{name}":NewEnumVal("Pewpew","{name}",true,{variants})'

        return ENUM_TEMPLATE.format_map(
            {
                "name": self.name,
                "variants": ",".join(f'"{variant}"' for variant in self.variants),
            }
        )

    def generate_docs(self) -> str:
        ENUM_TEMPLATE = "### `{name}`\n\n{variants}"

        return ENUM_TEMPLATE.format_map(
            {
                "name": self.name,
                "variants": "\n".join(f"- `{variant}`" for variant in self.variants),
            }
        )
