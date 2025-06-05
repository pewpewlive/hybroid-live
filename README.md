<img src="https://hybroid.pewpew.live/Logo.png" alt="Hybroid Logo" width="128" height="128">

# Hybroid Live [![Go](https://github.com/pewpewlive/hybroid/actions/workflows/go.yml/badge.svg)](https://github.com/pewpewlive/hybroid/actions/workflows/go.yml)

Programming language, handcrafted for PewPew.

## 🚧 Notice 🚧

As Hybroid Live is still in alpha, the language features may have breaking changes when updating. This notice will be removed once the language goes into a stable state.

## Why was Hybroid Live created?

Hybroid Live was created to overcome the limitations of Lua and projects like PewPewScript, and also for us to learn the making of programming languages.

## Why should I choose Hybroid Live over other solutions?

That's because Hybroid Live comes with many benefits, and only a few downsides.

Benefits of Hybroid Live:

- Contains many new features which are missing in Lua
- Optimized OOP via classes and entities
- Automatic dead-code elimination
- Strict typing
- Certain PewPew APIs are now an integral part of Hybroid Live (such as `tick` statement)
- Native support for fixedpoint numbers (including support for degree-to-radian conversion, transpile-time float-to-fixedpoint conversion)
- Native support for PewPew Marketplace

However, Hybroid Live does come with certain limitations:

- Not beginner-friendly
- No support for Lua helpers and libraries

## Syntax

The preliminary syntax for Hybroid Live can be found [here](spec.md).

## Syntax highlighting

An experimental VS Code extension is available [here](https://github.com/pewpewlive/hybroid-vscode).

## Building for release

Run `utils/build_hybroid.py`.

## License

Hybroid Live is licensed under Apache 2.0 license.
