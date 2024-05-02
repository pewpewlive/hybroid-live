# The Hybroid syntax

## Comments

- [x] Completed

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

- [x] Completed

Environments are an important aspect of PPL and Hybroid. Not specifying the environment will result in a transpile-time error.

The environment definition must be the first statement in the file.

```rs
@Environment(Level)

// The rest of the code
```

The following environments are available:

- `Level` - for working with levels
  - When choosing this environment, you get to use a subset of the Lua standard libraries: `table`, `string`, `fmath` (PPL-specific counterpart to `math`)
- `Mesh` - for working with meshes
  - When choosing this environment, all of the standard libraries that are enabled globally in PPL are available (exceptions being `coroutine`, `io`, `os`, etc.)
- `Sound` - for working with sounds
  - Same as `Mesh`
- `Shared` - for creating constant files referenced in multiple environments
  - When choosing this environment, `math` is disabled to work with `Level`, libraries open to `Level` are available
- `LuaGeneric` - for using standard Lua (e.g. console applications, etc.)
  - When choosing this environment, some features of the language would be disabled: `spawnable`s, `tick`, `spawn`, fixedpoint support, PPL libraries. All standard Lua libraries are available.

## Declaration of variables

- [x] Completed

```rs
// Local variables
let name = "Alpha"

// Global (public) variables
pub meaning_of_life = 42

// Reassignment
name = "blade"
```

## Typed declarations

- [ ] Completed

Types allow you to explicitly describe a variable's type. In Hybroid, types are not always necessary. Types might be necessary when you want to describe a complex type variable, or if the variable is left undefined. Types are what allows Hybroid to make sure you can write valid code without much headache and without the need to debug a lot.

```rs
let a: number // variable uninitialized, type required
let num = 1 // variable initialized, type inferred
let numbers: list<number> = [] // list is empty, list value type required
let callback: fn(text, bool) // function uninitialized, type required
```

## Declaration of constants

- [ ] Completed

```rs
const PI = 3.14f
```

## Entities and spawning syntax

- [ ] Completed

Entities are transpile-time classes. They are designed to provide OOP-like feel when working with entities. This feature is disallowed in `Generic` environments. Use `struct` keyword there instead.

### Defining an `entity`

```rs
entity Quadro {
  id: number

  mesh_id: number
  mesh_id2: number
  mesh_id3: number
  
  x = 1f y: fixed z: fixed = 0f
  
  /* you can also do this
  x = 1f; y: fixed; z: fixed = 0f

  x = 1f
  y: fixed
  z: fixed = 0f
  */
  speed = 10f
  damage = 2

  Spawn(x fixed, y fixed, speed fixed) {
    self.speed = speed 
    self.mesh_id2 = PewPew.NewEntity(x, y)
    PewPew.SetMesh(self, "file_path", 0)
    PewPew.SetMesh(self.mesh_id2, "file_path", 1)
    
    return self
  }

  Destroy() {
    PewPew.start_exploding(self, 30)
  }

  Update() {
    let x, y = PewPew.GetPosition(self)
    x = x + 10fx * self.speed
    PewPew.GetPosition(self, x, y)
  }

  WeaponCollision(index, wtype) {
  }

  PlayerCollision(index, ship_id) {
  }

  WallCollision(wall_x, wall_y) {
  }

  fn DamageOtherEntity(entity OtherEntity) {
    entity.Damage(self.damage)
  }
}
```

The Hybroid code shown gets generated into Lua like so:

```lua
QuadroStates = {}

local function quadro_update(id)
  local x, y = pewpew.entity_get_position(id)
  x = x + 10fx * QuadroState[id][8] -- this is speed because all of the entity fields get mapped to their indexes
  pewpew.entity_set_position(id, x, y)
end

local function quadro_weapon_collision(id, index, wtype)
end

local function quadro_player_collision(id, index, ship_id)
end

local function quadro_wall_collision(id, wall_x, wall_y)
end

function quadro_damage_other_entity(id, entity)
  other_entity_damage(entity, QuadroStates[id][9])
end

function Quadro.new(x, y, speed)
  local id = pewpew.new_customizable_entity(x, y)
  QuadroState[id] = {id, 0, 1, 2, 1fx, 0fx, 0fx, 10fx, 2} -- set default values specified in the entity fields

  pewpew.entity_set_update_callback(id, quadro_update)
  pewpew.customizable_entity_set_weapon_collision_callback(id, quadro_weapon_collision)
  pewpew.customizable_entity_set_player_collision_callback(id, quadro_player_collision)
  pewpew.customizable_entity_configure_wall_collision(id, quadro_wall_collision)

  QuadroState[id][8] = speed
  QuadroState[id][3] = pewpew.new_customizable_entity(x, y)
  pewpew.customizable_entity_set_mesh(id, "file_path", 0)
  pewpew.customizable_entity_set_mesh(QuadroState[id][3], "file_path", 1)

  return id
end
```

### Creating an entity

```rs
let quadro = spawn Quadro(100fx, 100fx, 10fx)

destroy quadro
```

## Lua interop & importing

- [ ] Completed

Original `pewpew`, `fmath`, `math`, `table` functions are available under `Origin` namespace.

Importing Lua libraries works as expected, just with omission of `/dynamic`.

```rs
use "mesh_helper.hyb" as mesh_helper_hybroid
use "shared.hyb"
```

You can write lua code with a special `@Lua` directive:

```rs
let number = 0

@Lua("number = number + 1")

Print(number) // -> 1
```

However, this is discouraged, as the transpiler can lose important context, such as variable declarations.

## Number Literals

- [x] Completed

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

- [ ] Completed

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

Repeat loops are simple `for` loops.

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

`for` loops are designed for advanced iterations.

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

- [x] Completed

In Lua, these structures are called "tables". These structures hold multiple data associated with a numeric index.

```rs
let fruits = ["banana", "kiwi", "apple", "pear", "cherry"]

Print(fruits[2]) // -> kiwi
```

To get the length of the list or , use `#` prefix.

```rs
let fruits = ["banana", "kiwi", "apple", "pear", "cherry"]

repeat @Len(fruits) with i {
  Print(fruits[i])
}
```

### Adding elements to the list

- [ ] Completed

Using `add` keyword.

```rs
let fruits = ["banana", "kiwi", "apple", "pear", "cherry"]

add "watermelon" to fruits

Print(@ListToStr(fruits)) // -> ["banana", "kiwi", "apple", "pear", "cherry", "watermelon"]
```

### Finding the index of the item

- [ ] Completed

Using `find` keyword. Only the first match is returned.

```rs
let fruits = ["banana", "kiwi", "apple", "pear", "cherry"]

Print(find "apple" in fruits) // -> 3
```

### Removing an element from the list

- [ ] Completed

Using `remove` keyword.

```rs
let fruits = ["banana", "kiwi", "apple", "pear", "cherry"]

remove 4 from fruits

Print(@ListToStr(fruits)) // -> ["banana", "kiwi", "apple", "cherry"]
```

## Maps

- [x] Completed

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

- [ ] Completed

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

- [ ] Completed

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

- [ ] Completed

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

- [x] Completed

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

- [ ] Completed

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

- [ ] Completed

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

- [ ] Completed

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

- [x] Completed

Structures are classes that do not have inheritance.

```rs

struct Rectangle {
  mesh_id1: number
  mesh_id1, mesh_id1, mesh_id1, mesh_id1 = 0f, 0f, 0f, 0f, 0f

  x, y = 0,0

  New(length number, height number) {
    self.length = length
    self.height = height
    return self
  }

  fn Area() {
    return self.length * self.height
  }

  fn Perimeter() {
    return (self.length + self.height) * 2
  }

  fn Move() {
    x += 5
  }
}

let rect = new Rectangle(100, 100)


Print(rect.Area())
```

```lua
function Rectangle_New(length, height)
  local new = {0, 0, nil}
  new[1] = length
  new[2] = height
  return new
end

function Rectangle_Area(self)
  return self[1] * self[2]
end

function Rectangle_Perimeter(self)
  return (self[1] + self[2]) * 2
end

function Rectangle_Move(self)
  self.x = self.x + 5
end

local rect = Rectangle_New(100, 100)

print(Area(rect))
```

