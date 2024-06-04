# Contribution Guide

This is meant to be a manual, teaching newcomers how to contribute to the project, be it the documentation or the codebase.

## Contributing to docs

Here lies the analysis of the composition of the docs, which will help you contribute to the documentation, as you will have a frame of reference as to how things are supposed to be structured.

The following documentation structure is valid to all the docs across the codebase. It must be ensured that their validity stays on par with this manual, so as to keep consistency.

**Here is a dummy documentation:**

# Component (folder) being documented

`# Component (folder) being documented`

_Overview_

## `file1.go`

`# file1.go`

Small description talking about what this file is useful for and what it contains.

### Types

_code snippet_
_description (optional)_

_..._

### Interfaces

`## Interfaces`

#### **_Foobar:_**

`### ***Foobar:***`

```go
type Foobar interface {
    foo(a int) bool
    bar(a bool) int
}
```

**Methods:**

1. `foo(a int) -> bool` - returns a boolean, explain more about what the function does
2. `bar(a bool) -> int` - returns an integer, same thing here

**Implementations:**

Here we don't need to enumerate through the implementations, we can just write the code snippet and explain more thoroughly wherever needed and so on.

For Example:

```go
type ExampleImplementation struct {
    num1 int
    num2 int
}
```

_description (optional)_

```go
type ExampleImplementation2 struct {
    nums []int
}
```

_description (optional)_

### **_Interface2:_**

`### ***Interface2:***`

_..._

If there are structures that don't implement any interface. You write them in the `Structures` section, the same way the interface section is written.

## Structures

`## Structures`

### **_StructName:_**

_code snippet_

**Methods:**\
_list_

**Associated Functions:**\
_list_

_..._

# `file2.go`

**-End of dummy documentation-**

Here are all the **sections**, ordered by their precedence:

1. Interfaces
2. Structures
3. Types
4. Enums
5. Global Variables

Here are all the **sub-sections**, ordered by their precedence:

1. _Name_ (interfaces only)
2. _Code Snippet_ (not written as a distinct subsection)
3. _Description_ (not a distinct subsection)
4. Methods
5. Extra Methods (interfaces only)
6. Associated Functions
7. Implementations (interfaces only)

It is noteworthy to say that functions, structures and any other value in the codebase may or may not have small helpful comments above them.

## Contributing to the codebase

It is recommended that you **check the documentation** of the codebase first before trying to contribute to it. If you have a proposal to make, you can **open an issue** and we can discuss it there.

The documentation for each component of the language can be found in its respective folder (the walker docs can be found in the walker folder).
