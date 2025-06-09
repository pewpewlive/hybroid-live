# The Hybroid Live syntax

## Comments

- [x] Completed

Comments in Hybroid Live are like in any other C-style language.

`//` indicates a single-line comment.

```rs
// I am a single-line comment!
Pewpew:Print("Hello, World!")
```

`/*` and `*/` indicate a multi-line comment.

```rs
/*
 I am a multi-line comment!
 Cool, right?
*/
Pewpew:Print("Hello, World!")
```

Comments are ignored, and will not show up in the generated Lua source code

## Semicolons

Just like in Lua, semicolons are treated as a whitespace character.

## Environments

- [x] Completed

Environments are an important aspect of PPL and Hybroid Live. Not specifying the environment will result in a transpile-time error.

The environment definition must be the first statement in the file.

```rs
env HelloWorld as Level

// The rest of the code
```

The following environments are available:

- `Level` - for working with levels
  - When choosing this environment, you get to use a subset of the Lua standard libraries: `table`, `string`, `fmath` (PPL-specific counterpart to `math`)
- `Mesh` - for working with meshes
  - When choosing this environment, all of the standard libraries that are enabled globally in PPL are available (exceptions being `coroutine`, `io`, `os`, etc.)
- `Sound` - for working with sounds
  - Same as `Mesh`

## Declaration of variables

- [x] Completed

```rs
// Local variables
let name = "Alpha"

// Global (public) variables
pub meaning_of_life = 42

// Assignment
name = "blade"
```

## Typed declarations

- [x] Completed

Types allow you to explicitly describe a variable's type. In Hybroid Live, types are not always necessary. Types might be necessary when you want to describe a complex type variable, or if the variable is left undefined. Types are what allows Hybroid Live to make sure you can write valid code without much headache and without the need to debug a lot.

```rs
let number a  // variable uninitialized, type required
let num = 1 // variable initialized, type inferred
list<number> numbers = [] // list is empty, list value type required
pub fn(text, bool) callback // function uninitialized, type required
```

## Declaration of constants

- [x] Completed

```rs
const PI = 3.14f
```

## Entities and spawning syntax

- [x] Completed

Entities are transpile-time classes. They are designed to provide OOP-like feel when working with entities. This feature is disallowed in `Generic` environments. Use `struct` keyword there instead.

### Defining an `entity`

```rs
entity Quadro {
  fixed x, y
  fixed speed

  spawn(fixed x, y, speed) {
    self.x = x
    self.y = y
    self.speed = speed
  }

  destroy() {
    Pewpew:ExplodeEntity(self, 30)
  }

  Update() {
    let x, y = Pewpew:GetPosition(self)
    x = x + 10f * self.speed
    Pewpew:SetPosition(self, x, y)
  }

  fn DamageOtherEntity(entity OtherEntity) {
    entity.Damage(self.damage)
  }
}
```

### Creating an entity

```rs
let quadro = spawn Quadro(100fx, 100fx, 10fx)

destroy quadro()
```

## Number Literals

- [x] Completed

In PPL, you use number literals with `fx` at the end of the number. But thankfully, Hybroid Live makes working with numbers easier, by giving several options.

### Fixedpoint Literal

Use `fx` to explicitly state you want to use fixedpoint numbers. This feature is disallowed in `Generic`, `Mesh` and `Sound` environments.

```rs
let speed = 100.2048fx
```

### Decimal Literal

If that's not what you want, Hybroid Live gives the option to use generic decimal literals by writing a float and adding `f` at the end

```rs
let a = 100.5f
let b = 3.14f
```

Behind the scenes, the transpiler will convert these numbers to their equivalent value based on the environment settings:

- On `Level` and `Shared` it will convert these numbers to their fixedpoint counterparts (`100.5f` will become `100.2048fx`)
- On `Mesh` and `Sound` it will stay as a decimal float, just without the 'f'

### Angle Literal

Hybroid Live also adds special literal support for angles.

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

- [x] Completed

### Tick loops

In PPL, for updating every tick, `pewpew.add_update_callback` is used. Hybroid Live wraps it in a `tick` statement.

```rs
tick {
  Pewpew:Print("I am printed every tick!")
}
```

It is possible to create a `tick` statement with a time variable.

```rs
tick with time {
  Pewpew:Print(time .. " has elapsed")
}
```

### While loops

In Hybroid Live and PPL while loops are discouraged. However, you can still use them if you want or need to.

```rs
while true {
  Pewpew:Print("Running infinitely and as fast as possible!")
}
```

### Repeat loops

Repeat loops are simple `for` loops.

```rs
repeat 10 {
  Pewpew:Print("Hybroid Live is awesome!")
}
```

It is possible to create a `repeat` loop with an iteration variable.

```rs
repeat 10 with index {
  Pewpew:Print("This is " .. index .. "th iteration!") // -> This is 1th iteration!
}
```

You can also specify the skip amount and the range, just like in lua.

```rs
repeat by 2 from 4 to 10 with index {
  Pewpew:Print("This is " .. index .. "th iteration!") // -> This is 4th iteration!
}
```

### For loops

`for` loops are designed for advanced iterations.

```rs
let fruits = ["banana", "kiwi", "apple", "pear", "cherry"]

for fruit in fruits {
  Pewpew:Print(fruit)
}
```

It is possible to also get an index.

```rs
let fruits = ["banana", "kiwi", "apple", "pear", "cherry"]

for index, item in fruits {
  Pewpew:Print(index)
}
```

## Lists

- [x] Completed

In Lua, these structures are called "tables". These structures hold multiple data associated with a numeric index.

```rs
let fruits = ["banana", "kiwi", "apple", "pear", "cherry"]

Pewpew:Print(fruits[2]) // -> kiwi
```

To get the length of the list, or use `#` prefix.

```rs
let fruits = ["banana", "kiwi", "apple", "pear", "cherry"]

repeat #fruits with i {
  Pewpew:Print(fruits[i])
}
```

### Adding elements to the list

- [ ] Completed

Using `add` keyword.

```rs
let fruits = ["banana", "kiwi", "apple", "pear", "cherry"]

add "watermelon" to fruits

Pewpew:Print(@ListToStr(fruits)) // -> ["banana", "kiwi", "apple", "pear", "cherry", "watermelon"]
```

### Finding the index of the item

- [ ] Completed

Using `find` keyword. Only the first match is returned.

```rs
let fruits = ["banana", "kiwi", "apple", "pear", "cherry"]

Pewpew:Print(find "apple" in fruits) // -> 3
```

### Removing an element from the list

- [ ] Completed

Using `remove` keyword.

```rs
let fruits = ["banana", "kiwi", "apple", "pear", "cherry"]

remove 4 from fruits

Pewpew:Print(@ListToStr(fruits)) // -> ["banana", "kiwi", "apple", "cherry"]
```

## Maps

- [x] Completed

In Lua, these structures are also called _tables_. These structures hold multiple data entries associated with a string index.

```rs
let inventory = {
  bananas = 2,
  apples = 5,
  kiwis = 10,
  pears = 0,
  cherries = 12, // trailing comma is optional!
}

Pewpew:Print(fruits["apples"]) // -> 5
```

### Adding elements to the map

- [x] Completed

Using `add` keyword.

```rs
let inventory = {
  bananas = 2,
  apples = 5,
  kiwis = 10,
  pears = 0,
  cherries = 12,
}

add 10 as "watermelon" to inventory

Pewpew:Print(ToString(fruits))

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
  bananas = 2,
  apples = 5,
  kiwis = 10,
  pears = 0,
  cherries = 12,
}

Pewpew:Print(find 10 in fruits) // -> "kiwis"
```

### Removing an element from the map

- [ ] Completed

Using `remove` keyword.

```rs
let inventory = {
  bananas = 2,
  apples = 5,
  kiwis = 10,
  pears = 0,
  cherries = 12,
}

remove "cherries" from fruits

Pewpew:Print(@MapToStr(fruits))

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
fn Greet(text name) {
  Pewpew:Print("Hello" .. name .. "!")
}

Greet("John") // -> Hello, John!
```

Functions can be annonymous, too! Useful for callbacks.

```rs
let Greet = fn(name) {
  Pewpew:Print("Hello" .. name .. "!")
}

Greet("John") // -> Hello, John!
```

## Macros

- [ ] Completed

Macros are special functions that are evaluated in the transpiler.

```rs
macro CoolMacro(params) => "Hello " .. params 

macro HandleEntity($params) => {
  let id = new $params()
  id.AddCallback()
} 
```

When you use them:

```rs
Pewpew:Print(CoolMacro("John" .. "!"))
```

The expanded code looks something like this:

```rs
Pewpew:Print("Hello " .. "John" .. "!")
```

## Conditional statements

- [x] Completed

### If statement

```rs
let a = 10

if a == 10 {
  Pewpew:Print("It's 10!")
} else if a == 20 {
  Pewpew:Print("It's 20!")
} else {
  Pewpew:Print("It's a different number!")
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

Pewpew:Print(check)
```

### Match statement

```rs
match a {
  1, 10 => {
    // if a is 1 or 10 then execute
  }
  20 => {
    a = 24
    return
  }
  else => {
    a = nil
  }
}

let a = 10
let check = match a {
  10 => "It's 10!"
  20 => "It's 20!"
  else => "It's a different number!"
}

Pewpew:Print(check)
```

## Enums

- [x] Completed

Enums are converted to tables if compiling to Lua.

```rs
enum SandwichType {
  Blt,
  Panini,
  GrilledCheese,
  Ham,
}
```

## Classes

- [x] Completed

Classes are structs that allow methods.

```rs

class Rectangle {
  number width, height

  new(number width, height) {
    self.width = width
    self.height = height
  }

  fn Area() -> number {
    return width * height
  }
}

let rect = new Rectangle(100, 100)

Pewpew:Print(rect.Area())
```
