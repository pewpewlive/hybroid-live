# The Hybroid syntax

## Comments

Comments in Hybroid are like in any other C-style language.

`//` indicates a single-line comment.

```rs
// I am a single-line comment!
Print("Hello, World!")
```

`/*` and `*/` indicate a multi-line comment.

```rs
/*
 I am a multi-line comment!
 Cool, right?
*/
Print("Hello, World!")
```

## Semicolons

Just like in Lua, semicolons are treated as a whitespace character.

## Environments

Environments are an important aspect of PPL and Hybroid. Not specifying the environment will result in a transpile-time error.

The environment definition must be the first statement in the file.

```rs
@Environment(Level)

// The rest of the code
```

The following environments are available:

- `Level` - for working with levels
  - When choosing the this environment, you get to use a subset of the Lua standard libraries: `table`, `string`, `fmath` (PPL-specific counterpart to `math`)
- `Mesh` - for working with meshes
  - When choosing the this environment, all of the standard libraries that are enabled globally in PPL are available (exceptions being `coroutine`, `io`, `os`, etc.)
- `Sound` - for working with sounds
  - Same as `Mesh`
- `Shared` - for creating constant files referenced in multiple environments
  - When choosing the this environment, `math` is disabled to work with `Level`, libraries open to `Level` are available
- `LuaGeneric` - or cothistandard Lua (for use in console applications, etc.)
  - When choosing the `LuaGeneric` environment, some features of the language would be disabled: `spawnable`s, `tick`, `spawn`, fixedpoint support, PPL libraries. All standard Lua libraries are available.

## Declaration of variables

```rs
// Local variables
let name = "Alpha"

// Global (public) variables
pub number_of_life = 42

// Reassignment
name = "blade"
```

## Declaration of constants

```rs
const PI = 3.14f
```

## Entities and spawning syntax

Entities are transpile-time classes. They are designed to provide OOP-like feel when working with entities. This feature is disallowed in `Generic` environments. Use `struct` keyword there instead.

### Defining an `entity`

```rs
entity Quadro {
  define {
    mesh_id2,
    speed,
    other
  }

  Spawn(x, y, speed) {
    self.speed = speed
    self.mesh_id2 = PewPew.NewEntity(x, y)
    PewPew.EntitySetMesh(self, "file_path", 0)
    PewPew.EntitySetMesh(self.mesh_id2, "file_path", 1)
    Fmath.Random_Int(0,6)
    return self
  }

  trait Update() {
    let x, y = Origin.pewpew.entity_get_position(self)
    x = x + 10fx * self.speed
    Origin.pewpew.entity_set_position(self, x, y)
  }

  trait WeaponCollision(index, wtype) {
  }

  trait PlayerCollision(index, ship_id) {
  }

  trait WallCollision(wall_x, wall_y) {
  }

  fn DamageOtherEntity(entity, x, y) {
    entity.damage(1)
  }
}
```

```lua
QuadroStates = {}

local function quadro_update(id)
  local x, y = pewpew.entity_get_position(id)
  x = x + 10fx * QuadroState[id].speed
  pewpew.entity_set_position(id, x, y)
end

local function quadro_weapon_collision(id, index, wtype)
  -- does stuff
end

local function quadro_player_collision(id, index, ship_id)
  -- does stuff
end

local function quadro_wall_collision(id, wall_x, wall_y)
  -- does stuff
end

function quadro_damage_other_entity(id, entity)
  (Type of entity)[id].damage(1)
end

function Quadro.new(x, y, speed)
  local id = pewpew.new_customizable_entity(x, y)
  QuadroState[id] = {}

  pewpew.entity_set_update_callback(id, quadro_update)
  pewpew.customizable_entity_set_weapon_collision_callback(id, quadro_weapon_collision)
  pewpew.customizable_entity_set_player_collision_callback(id, quadro_player_collision)
  pewpew.customizable_entity_configure_wall_collision(id, quadro_wall_collision)

  QuadroState[id].speed = speed
  QuadroState[id].mesh_id2 = pewpew.new_customizable_entity(x, y)
  pewpew.customizable_entity_set_mesh(id, "file_path", 0)
  pewpew.customizable_entity_set_mesh(QuadroState[id].mesh_id2, "file_path", 1)

  return id
end

```

### Creating an entity

```rs
let id = spawn Quadro with {x: 100fx, y: 100fx, speed: 10fx}

local id = Quadro.new(100fx, 100fx, 10fx)

// Invoking a Destroy trait
Destroy id
```

## Lua interop & importing

Original `pewpew`, `fmath`, `math`, `table` functions are available under `Origin` namespace.

Importing Lua libraries works as expected, just with omission of `/dynamic`.

```rs
use "mesh_helper.hyb" as mesh_helper_hybroid
```

You can write lua code with a special `@Lua` directive:

```rs
let number = 0

@Lua("number = number + 1")

Print(number) // -> 1
```

However, this is discouraged, as the transpiler can lose important context, such as variable declarations.

## Number Literals

In PPL, you use number literals with `fx` at the end of the number. But thankfully, Hybroid makes working with numbers easier, by giving several options.

### Fixedpoint Literal

Use `fx` to explicitly state you want to use fixedpoint numbers. This feature is disallowed in `Generic`, `Mesh` and `Sound` environments.

```rs
let speed = 100.2048fx
```

### Decimal Literal

If that's not what you want, Hybroid gives the option to use generic decimal literals by writing a float and adding `f` at the end

```rs
let a = 100.5f
let b = 3.14f
```

Behind the scenes, the transpiler will convert these numbers to their equivalent value based on the environment settings:

- On `Level` and `Shared` it will convert these numbers to their fixedpoint counterparts (`100.5f` will become `100.2048fx`)
- On `Mesh`, `Sound` and `Generic` it will stay as a decimal float, just without the 'f'

### Angle Literal

Hybroid also adds special literal support for angles.

```rs
let degrees = 180d
let pi = 3.14r
```

When using angle literals, the transpiler will automatically convert their values:

- The `d` literal allows you to write angles in degrees. They are automatically converted to radians and directly placed in the final Lua code.
- The `r` literal is functionally the same as a decimal `f` literal, keeping its value without the `r`. It is useful to denote when arguments are angles or just numbers.

### Other Literals

- `0x` is a hexadecimal literal. Example: `0xff`
- `0o` is an octal literal. Example: `0o07`
- `0b` is a binary literal. Example: `0b01`

## Loops

### Tick loops

In PPL, for updating every tick, `pewpew.add_update_callback` is used. Hybroid wraps it in a `tick` statement.

```rs
tick {
  Print("I am printed every tick!")
}
```

It is possible to create a `tick` statement with a time variable.

```rs
tick with time {
  Print(time .. " has elapsed")
}
```

### While loops

In Hybroid and PPL while loops are discouraged. However, you can still use them if you want or need to.

```rs
while true {
  Print("Running infinitely and as fast as possible!")
}
```

### Repeat loops

Repeat loops are simple for loops.

```rs
repeat 10 {
  Print("Hybroid is awesome!")
}
```

It is possible to create a `repeat` loop with an iteration variable.

```rs
repeat 10 with index {
  Print("This is " .. index .. "th iteration!") // -> This is 1th iteration!
}
```

### For loops

For loops are designed for advanced iterations.

```rs
let fruits = ["banana", "kiwi", "apple", "pear", "cherry"]

for fruit in fruits {
  Print(fruit)
}
```

It is possible to also get an index.

```rs
let fruits = ["banana", "kiwi", "apple", "pear", "cherry"]

for index, item in fruits {
  Print(index)
}
```

## Lists

In Lua, these structures are called "tables". These structures hold multiple data associated with a numeric index.

```rs
let fruits = ["banana", "kiwi", "apple", "pear", "cherry"]

Print(fruits[2]) // -> kiwi
```

To get the length of the list or , use `#` prefix.

```rs
let fruits = ["banana", "kiwi", "apple", "pear", "cherry"]

repeat #fruits with i {
  Print(fruits[i])
}
```

### Adding elements to the list

Using `add` keyword.

```rs
let fruits = ["banana", "kiwi", "apple", "pear", "cherry"]

add "watermelon" to fruits

Print(@ListToStr(fruits)) // -> ["banana", "kiwi", "apple", "pear", "cherry", "watermelon"]
```

### Finding the index of the item

Using `find` keyword. Only the first match is returned.

```rs
let fruits = ["banana", "kiwi", "apple", "pear", "cherry"]

Print(find "apple" in fruits) // -> 3
```

### Removing an element from the list

Using `remove` keyword.

```rs
let fruits = ["banana", "kiwi", "apple", "pear", "cherry"]

remove 4 from fruits

Print(@ListToStr(fruits)) // -> ["banana", "kiwi", "apple", "cherry"]
```

## Maps

In Lua, these structures are also called _tables_. These structures hold multiple data entries associated with a string index.

```rs
let inventory = {
  bananas: 2,
  apples: 5,
  kiwis: 10,
  pears: 0,
  cherries: 12
}

Print(fruits["apples"]) // -> 5

// or

Print(fruits.apples) // -> 5
```

### Adding elements to the map

Using `add` keyword.

```rs
let inventory = {
  bananas: 2,
  apples: 5,
  kiwis: 10,
  pears: 0,
  cherries: 12
}

add 10 as "watermelon" to inventory

Print(@MapToStr(fruits))

/*
-> {
  bananas: 2,
  apples: 5,
  kiwis: 10,
  pears: 0,
  cherries: 12,
  watermelons: 10
}
*/
```

### Finding the key of the item

Using `find` keyword. Only the first match is returned.

```rs
let inventory = {
  bananas: 2,
  apples: 5,
  kiwis: 10,
  pears: 0,
  cherries: 12
}

Print(find 10 in fruits) // -> "kiwis"
```

### Removing an element from the map

Using `remove` keyword.

```rs
let inventory = {
  bananas: 2,
  apples: 5,
  kiwis: 10,
  pears: 0,
  cherries: 12
}

remove "cherries" from fruits

Print(@MapToStr(fruits))

/*
-> {
  bananas: 2,
  apples: 5,
  kiwis: 10,
  pears: 0
}
*/
```

## Functions

Declaring a function works with the `fn` keyword. Functions are local by default.

```rs
fn Greet(name) {
  Print("Hello" .. name .. "!")
}

Greet("John") // -> Hello, John!
```

Functions can be annonymous, too! Useful for callbacks.

```rs
let Greet = fn (name) {
  Print("Hello" .. name .. "!")
}

Greet("John") // -> Hello, John!
```

## Directives

Directives are special functions that are evaluated in the transpiler. They work similarly to _macros_.

```rs
dir @Hello(name) { 
  "Hello ".. name .. "!" 
}

print(@Hello("John")) // -> Hello, John!
```

The generated code looks something like this:

```lua
print("Hello " .. "John" .. "!")
```

## Conditional statements

### If statement

```rs
let a = 10

if a == 10 {
  Print("It's 10!")
} else if a == 20 {
  Print("It's 20!")
} else {
  Print("It's a different number!")
}
```

If statements can also be used as expressions.

```rs
let a = 10

let check = if a == 10 {
  return "It's 10!"
} else if a == 20 {
  return "It's 20!"
} else {
  return "It's a different number!"
}

Print(check)
```

### Match statement

```rs
match a {
  1 => // if a is 1 or 10 then execute
  10 => {
    //execute
  }
  20 => {
    a = 24
    return
  }
  _ => { // else
    a = nil
  }
}

let a = 10
let check = match a {
  10 => "It's 10!"
  20 => "It's 20!"
  _ => "It's a different number!"
}

Print(check)
```

## Enums

Enums are converted to tables if compiling to Lua.

```rs
enum SandwichType {
  Blt,
  Panini,
  GrilledCheese,
  Ham
}
```

## Structures

Structures are classes that do not have inheritance. Using structures in any environment except `Generic` will result in a warning.

```rs
struct Rectangle {
  length: 0
  height: 0

  fn Init(length, height) {
    self.length = length
    self.height = height
  }

  fn Area() {
    return self.length * self.height
  }

  fn Perimeter() {
    return (self.length + self.height) * 2
  }
}

let rect = Rectangle.Init()
```
