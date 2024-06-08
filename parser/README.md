# Parser

The parser converts tokens to an Abstract Syntax Tree (AST) for it to be processed further.

## `parser.go`

The parser file has the necessary structure and functions to work with the `Parser`.

### Structures

#### **_Parser_**

The `Parser` holds the tokens ready to be converted, an index pointing to the current token being parsed, as well as a list of errors generated along the way when processing the tokens. The resulting program is stored as a list of `ast.Node`s.

```go
type Parser struct {
 program []ast.Node
 current int
 tokens  []lexer.Token
 Errors  []ast.Error
}
```

#### **Constructor:** 
`NewParser() -> *Parser` - Creates a new `Parser` with an empty list of tokens, errors, empty program and returns a pointer to it.

#### **Methods:**

1. `AssignTokens(tokens []lexer.Token)` - Assigns a list of tokens (returned by the tokenizer) to be parsed.
2. `error(token lexer.Token, msg string)` - Appends an error with the following `token` and `msg` to the parser list of errors.
3. `synchronize()` - A method that tries to synchronize the parser back with the next tokens that are valid.
4. `isMultiComparison() -> bool` - Checks if the current token is either an `and` or `or` keyword.
5. `isComparison() -> bool` - Checks if the current token is either of these: `>`, `>=`, `<`, `<=`, `!=`, `==`.
6. `isAtEnd() -> bool` - Checks if the current position the parser is at is the End Of File.
7. `advance() -> lexer.Token` - Advances by one into the next token and returns the previous token before advancing.
8. `peek(offset ...int) -> lexer.Token` - Peeks into the current token or peeks at the token that is offset from the current position by the given optional `offset`.
9. `check(tokenType lexer.TokenType) -> bool` - Checks if the current token is the specified `tokenType`. Note: returns false if it's the End Of File.
10.  `match(types ...lexer.TokenType) -> bool` - Matches the given list of tokens and advances if they match.
11.  `consume(message string, types ...lexer.TokenType) -> (lexer.Token, bool)` - Consumes a list of tokens, advancing if they match and returns true. It also advances if none of the tokens were able to match, and returns false. Note: it creates errors with the specified `message` if necessary.
12. `ParseTokens() -> []ast.Node` - Parses the assigned tokens, and returns a program (list of `ast.Node`s) for it to be later walked.

## `parser_helpers.go`

The file that holds all of the necessary helpers for the [Parser](https://github.com/pewpewlive/hybroid/blob/master/parser/README.md#parsergo). It makes the code less repetitive and easier to grasp.

### **Methods:**

(This file adds additional methods to [Parser](https://github.com/pewpewlive/hybroid/blob/master/parser/README.md#parsergo))

1. `createBinExpr(left ast.Node, operator lexer.Token, tokenType lexer.TokenType, lexeme string, right ast.Node) -> ast.Node` - Evaluates the value type, creates a `BinaryExpr` with the respective parameters, and returns it.
2.  `getOp(opEqual lexer.Token) -> lexer.Token` - Returns the respective operation when using assignment operators. For example: if given `lexer.MinusEqual` (`-=`), it returns `lexer.Minus`.
3.  `getParam() -> ast.Param` - Attempts to get the current token (that is an identifier) and its type. Returns an `ast.Param` type with the respective values.
4.  `parameters() -> []ast.Param` - Returns a list of `ast.Param`s. Uses `getParam()` under the hood to get all of the parameters. Note: throws errors if the expression is missing parentheses. <!-- FIXME: Think of a better description -->
5.  `arguments() -> []ast.Node` - Returns a list of `ast.Param`s. Uses `getParam()` under the hood to get all of the parameters. Note: throws errors if the expression is missing parentheses. <!-- FIXME: Think of a better description -->

### **Extra functions:**

1. `IsFx(valueType ast.PrimitiveValueType) -> bool` - Checks if the `valueType` is expected to be a fixedpoint.

## `statements.go`

