# mox

Mox is an interpreted programming language built with Golang.

## Core Features

### Data Types

Variables are defined using the `let` keyword, and constants with `const` keyword. Mox supports the following data types:

- `int`
- `float`
- `string`
- `bool`

```mox
let x = 5;          // Integer
let y = 3.14;      // Float
let name = "Mox";  // String
let isActive = true; // Boolean

// Explicit type declaration
let z: int = 10;
let pi: float = 3.1415;
let greeting: string = "Hello, Mox!";
```


### Functions

Functions are defined using `fn` keyword.

```mox
fn add(x, y) {
    return x + y;
}


// With type declaration

fn multiply(x: int, y: int): int {
    return x * y;
}

// Without return statement (implicit returns)
fn greet(name: string) {
    f"Hello %{name}!"
}
```

### Control Flow

Mox supports standard control flow constructs like `if`, `else`, and `while`.

```mox
let num = 10;

while (num > 0) {
    num -= 1;
}

// If-Else

if (num == 0) {
    f"%{num} is zero";
} else if (num == 10) {
    f"%{num} is ten";
} else {
    f"%{num} is neither zero nor ten";
}
```

### Classes

```mox

class Person {
    name: string;
    age: int;
    private ssn: string;
}

impl Person {
    fn init(name: string, age: int, ssn: string) {
        self.name = name;
        self.age = age;
        self.dsn = dsn;
    }
    fn greet() {
        f"Hello, my name is %{self.name} and I am %{self.age} years old.";
    }

    private fn set_ssn(ssn: string) {
        self.ssn = ssn;
    }

    fn get_ssn(): string {
        return self.ssn;
    }
}
```

### Expressions

Mox supports the following prefix operator Expressions

```typescript
-5             // (minus)
!<expression>  // (bang)
```

And the following infix ops

```ts
a + - * / < > b
==
!=
>=
<=
```


### Imports (TBD)

Every mox fils can be imported as a module.
Imports are also an expression, so if-expressions and conditional imports are allowed in imports.


```ts
import { foo };
import { foo.bar } bar;
import { if (cond) foo else baz } jeff;
```
