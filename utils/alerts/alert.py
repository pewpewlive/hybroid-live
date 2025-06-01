from . import helpers, imports


class Function:
    name: str
    params: str
    returns: str
    code: str

    def __init__(
        self,
        name: str,
        returns: str,
        code: str,
        params: str = "",
    ):
        self.name = name
        self.params = params
        self.returns = returns
        self.code = code

    def string(self, receiver: str, name: str) -> str:
        FUNCTION_TEMPLATE = "func ({} *{}) {}({}) {} {{\n  {}\n}}"

        return FUNCTION_TEMPLATE.format(
            receiver, name, self.name, self.params, self.returns, self.code
        )


class Alert:
    name: str
    receiver: str
    type: str
    stage: str
    fields: dict[str, str]
    message: str
    message_format: list[helpers.Format]
    note: str
    note_format: list[helpers.Format]
    id: int

    def __init__(self, raw: dict, stage: str, id: int):
        name = raw.get("name")
        if name is None:
            raise ValueError(f"Name must not be None, Raw info: {raw}")
        self.name = name

        self.receiver = helpers.to_receiver(self.name)

        type = raw.get("type")
        if type is None:
            raise ValueError(f"Type must not be None, Raw info: {raw}")
        self.type = type

        self.stage = stage

        self.fields = {"Specifier": "Snippet"} | raw.get("fields", {})

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

    def generate(self) -> str:
        ALERT_TEMPLATE = "type {name} struct {{\n  {fields}\n}}\n\n{functions}"

        for type in self.fields.values():
            imports.update_imports(type)

        functions = [
            Function(
                name="Message",
                returns="string",
                code=f"return {helpers.format_string(self.message, self.message_format, self.receiver)}",
            ),
            Function(
                name="SnippetSpecifier",
                returns="Snippet",
                code=f"return {self.receiver}.Specifier",
            ),
            Function(
                name="Note",
                returns="string",
                code=f"return {helpers.format_string(self.note, self.note_format, self.receiver)}",
            ),
            Function(
                name="ID",
                returns="string",
                code='return "hyb{:03d}{}"'.format(self.id, self.stage[0]),
            ),
            Function(name="AlertType", returns="Type", code=f"return {self.type}"),
        ]

        return ALERT_TEMPLATE.format_map(
            {
                "name": self.name,
                "fields": "\n  ".join(
                    f"{field} {type}" for field, type in self.fields.items()
                ),
                "functions": "\n\n".join(
                    function.string(self.receiver, self.name) for function in functions
                ),
            }
        )
