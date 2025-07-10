<img src="https://hybroid.pewpew.live/Logo.png" alt="Hybroid Logo" width="128" height="128">

# Hybroid Live [![Go](https://github.com/pewpewlive/hybroid/actions/workflows/go.yml/badge.svg)](https://github.com/pewpewlive/hybroid/actions/workflows/go.yml)

Programming language, handcrafted for PewPew.

## 🚧 Notice 🚧

As Hybroid Live is still in alpha, the language features may have breaking changes when updating. This notice will be removed once the language goes into a stable state.

## Why was Hybroid Live created?

Hybroid Live was created to overcome the limitations and shortcomings of Lua, as well as to provide a better developer experience.

## Pros and cons of Hybroid Live

Benefits of Hybroid Live:

- Contains many new features which are missing in Lua (such as enums, classes, etc.)
- PewPew Live specific features (such as the tick loop or entities)
- State of the art error messages, inspired by Rust and Scala
- Optimized OOP via classes and entities
- Familiar syntax reminiscent of Rust and other popular languages
- Automatic dead-code elimination
- Strict typing
- Native support for fixedpoint numbers (including support for degree-to-radian and float-to-fixedpoint conversion)

However, Hybroid Live does come with certain limitations:

- Not beginner-friendly
- No support for Lua interoperability

## Syntax

The syntax for Hybroid Live can be found [in the spec](spec.md). The syntax specification may not be up to date.

## Syntax highlighting

An experimental VS Code extension is available [on github](https://github.com/pewpewlive/hybroid-vscode).

## Getting Started

Check out the "Getting Started" section of the [Hybroid documentation site](https://pewpewlive.github.io/hybroid-docs/get-started/quick-start/).

## Building for release

Run `utils/build_hybroid.py`.

## License

Hybroid Live is licensed under Apache 2.0 license.
