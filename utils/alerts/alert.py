from . import helpers, imports, providers


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
        FUNCTION_TEMPLATE = "func ({} *{}) {}({}) {} {{{}}}"

        return FUNCTION_TEMPLATE.format(
            receiver, name, self.name, self.params, self.returns, self.code
        )


class Alert:
    name: str
    receiver: str
    type: str
    stage: str
    fields: dict[str, str]
    providers: dict[str, providers.Provider]
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

        alert_type = raw.get("type")
        if alert_type is None:
            raise ValueError(f"Type must not be None, Raw info: {raw}")
        self.type = alert_type

        self.stage = stage

        self.fields: dict[str, str] = {}
        self.providers: dict[str, providers.Provider] = {}
        for field, value in raw.get("fields", {}).items():
            if type(value) is dict:
                field_type = value["type"]
                provider_name = value.get("provider")
                if provider_name is not None:
                    provider = providers.PROVIDERS.get(provider_name)
                    if provider is None:
                        raise ValueError(
                            f"Unknown provider '{provider_name}' for field '{field}', Raw info: {raw}"
                        )
                    self.providers[field] = provider
                else:
                    self.fields[field] = field_type
                imports.update_imports(field_type)
            else:
                self.fields[field] = value
                imports.update_imports(value)

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

    def accessors(self) -> dict[str, str]:
        return {field: provider.accessor for field, provider in self.providers.items()}

    def distinct_providers(self) -> list[providers.Provider]:
        distinct: list[providers.Provider] = []
        for provider in self.providers.values():
            if provider not in distinct:
                distinct.append(provider)
        return distinct

    def generate(self) -> str:
        ALERT_TEMPLATE = """type {name} struct {{
            {body}
        }}
        {constructor}
        {functions}"""

        CONSTRUCTOR_TEMPLATE = """func New{name}({params}) *{name} {{
            return &{name}{{
                {inits}
            }}
        }}"""

        body = ["SnippetProvider"]
        body += [provider.embed for provider in self.distinct_providers()]
        body += [f"{field} {field_type}" for field, field_type in self.fields.items()]

        functions = [
            Function(
                name="Message",
                returns="string",
                code=f"return {helpers.format_string(self.message, self.message_format, self.receiver, self.accessors())}",
            ),
            Function(
                name="Note",
                returns="string",
                code=f"return {helpers.format_string(self.note, self.note_format, self.receiver, self.accessors())}",
            ),
            Function(name="Type", returns="Type", code=f"return {self.type}"),
            Function(
                name="ID",
                returns="string",
                code='return "hyb{:03d}{}"'.format(self.id, self.stage[0]),
            ),
        ]

        for provider in self.distinct_providers():
            functions.append(
                Function(
                    name=provider.setter,
                    returns=f"*{self.name}",
                    params=f"{provider.param} string",
                    code=(
                        f"{self.receiver}.{provider.accessor} = {provider.param}\n"
                        f"return {self.receiver}"
                    ),
                )
            )

        params = ", ".join(
            ["span core.Span"]
            + [
                f"{helpers.to_param_name(field)} {field_type}"
                for field, field_type in self.fields.items()
            ]
        )

        inits = ["SnippetProvider: SnippetProvider{span: span},"]
        for provider in self.distinct_providers():
            inits.append(f"{provider.embed}: {provider.init}, // Default value")
        for field in self.fields:
            inits.append(f"{field}: {field},")

        return ALERT_TEMPLATE.format_map(
            {
                "name": self.name,
                "body": "\n  ".join(body),
                "constructor": CONSTRUCTOR_TEMPLATE.format_map(
                    {
                        "name": self.name,
                        "params": params,
                        "inits": "\n  ".join(inits),
                    }
                ),
                "functions": "\n".join(
                    function.string(self.receiver, self.name) for function in functions
                ),
            }
        )
