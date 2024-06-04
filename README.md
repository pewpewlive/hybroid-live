<img src="https://hybroid.pewpew.live/Logo.png" alt="Hybroid Logo" width="128" height="128">

# Hybroid

Programming language, handcrafted for PewPew.

## 🚧 Notice 🚧

As Hybroid is still in alpha, the language features may have breaking changes when updating. This notice will be removed once the language goes into a stable state.

## Why was Hybroid created?

Hybroid was created to overcome the limitations of Lua and projects like PewPewScript, and also for us to learn the making of programming languages.

## Why should I choose Hybroid over other solutions?

That's because Hybroid comes with many benefits, and only a few downsides.

Benefits of Hybroid:

- Contains many new features which are missing in Lua
- Optimized OOP via structs and entities
- Automatic dead-code elimination
- Strict typing
- Certain PewPew APIs are now an integral part of Hybroid (such as `tick` statement)
- Native support for fixedpoint numbers (including support for degree-to-radian conversion, transpile-time float-to-fixedpoint conversion)
- Native support for PewPew Marketplace

However, Hybroid does come with certain limitations:

- Not beginner-friendly
- Limited support for Lua helpers and libraries

## Syntax

The preliminary syntax for Hybroid can be found [here](spec/syntax.md).

## Syntax highlighting

An experimental VS Code extension is available [here](https://github.com/pewpewlive/hybroid-vscode).

## Goals

- [ ] Full spec

- [x] Working lexer

- [ ] Working parser

- [ ] Working AST walker

- [x] Basic Lua codegen

- [ ] Advanced codegen

- [ ] APPAPI implementation

- [ ] Language server and syntax highlighting

## License

Hybroid is licensed under Apache 2.0 license.
