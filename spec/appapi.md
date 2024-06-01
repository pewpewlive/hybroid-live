# Hybroid PewPew API

Augmented PewPew API: PewPew API, developed to use all Hybroid's features and to improve developer experience.

## APPAPI VS PPOL

APPAPI:

- For use in Hybroid
- Powerful and efficient abstractions, such as usage of structs and entities
- Converted into native PewPew and Fmath calls on transpile-time
- API split into namespaces for best organization
- Original PewPew and Fmath functions still available in `Origin` namespace
- Built right into Hybroid
- Functions adapt based on environment
- Includes improvements to the PewPew API
- PewPew enums treated as actual enums
- `tick` statement is an integral part of Hybroid language, being converted to `pewpew.add_update_callback()`

PPOL:

- For use in Lua (Hybroid may be compatible via Lua helper imports, however, no support will be provided)
- API exists in global scope
- Functions renamed into their shorthands
- Original PewPew and Fmath functions are removed
- Needs importing via `require()` function
- Includes improvements to the PewPew API
- PewPew enums are replaced with tables having indexes or having shorthands

## API

### `void PewPew.Level.SetSize(width fixed, height fixed)`

### `void PewPew.Level.SetSize(width fixed)`

Sets the level size to a desired one. If only one argument is given, the level will be square.
`void PewPew.Player.SetConfig(index int, config map)`
Configures the player.

### `entity PewPew.Entity.Ship`

`PewPew.Entity.Ship Spawn(x fixed, y fixed, index int)`
`PewPew.Entity.Ship Spawn(x fixed, y fixed, index int, PewPew.Types.WeaponConfig weapon_config)`
A spawnable player entity.
Example:

```rust
let playerShip = spawn Ship(100f, 100f, 0, {weaponType: Double, WeaponFreq: Hz10})
Print(@MapToStr(playerShip.GetConfig()))
```

`map GetConfig()`
Gets the configuration for the player.
`Damage(damage int)`
Damages the player.
