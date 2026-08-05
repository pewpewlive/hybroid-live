class Provider:
    embed: str
    init: str
    accessor: str
    setter: str
    param: str

    def __init__(self, embed: str, init: str, accessor: str, setter: str, param: str):
        self.embed = embed
        self.init = init
        self.accessor = accessor
        self.setter = setter
        self.param = param


PROVIDERS = {
    "context": Provider(
        embed="ContextProvider",
        init='ContextProvider{context: ""}',
        accessor="context",
        setter="WithContext",
        param="ctx",
    ),
}
