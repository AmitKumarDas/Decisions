### Why RUST?
I am a developer at my core. I like to build mental models by learning, observing, implementing programs using different
languages. Currently, RUST claims to be the best with respect to perforamce and memory usage. It is also one of the modern
day programming languages.

## Rough
- install rust via rustup
```bash
$ curl https://sh.rustup.rs -sSf | sh
```
```bash
It will add the cargo, rustc, rustup and other commands to Cargo's bin 
directory, located at:

  /home/amit/.cargo/bin

This path will then be added to your PATH environment variable by modifying the
profile file located at:

  /home/amit/.profile

You can uninstall at any time with rustup self uninstall and these changes will
be reverted.
```
```bash

To get started you need Cargo's bin directory ($HOME/.cargo/bin) in your PATH 
environment variable. Next time you log in this will be done automatically.

To configure your current shell run source $HOME/.cargo/env
  ~ > source /home/amit/.cargo/env 
```
- update
```bash
If you’ve installed Rust via rustup, updating to the latest version is easy. 
From your shell, run the following update script:

$ rustup update
```
- linker
> Additionally, you’ll need a linker of some kind. It’s likely one is already installed, but when you try to compile a Rust program and get errors indicating that a linker could not execute, that means a linker isn’t installed on your system and you’ll need to install one manually. C compilers usually come with the correct linker. Check your platform’s documentation for how to install a C compiler.
- macro
```rust
println!("This is macro with !");
```
- Ahead-Of-Time compiled language
> Rust is an ahead-of-time compiled language, meaning you can compile a program and give the executable to someone else, and they can run it even without having Rust installed.
- dependency management ?
> cargo .. same as dep in go
- Does code compiles ?
```
cargo check
```
- Want to build?
```
cargo build
```
- Want to build & run?
```
cargo run
```
- Want to release?
```
cargo build --release
```
- Create new project
```
cargo new my-project
```

### Mental Model - Day 1
- variable name starts with `let name`
- if mutable tell it explicitly `let mut name`
- pass mutable explicitly `func(&mut name)`
- func names use underscore `my_func_name()`
- Result is an enum Type with Ok & Error
- **!Tip** **!Note**
```bash
- Rust has a number of types named Result in its standard library
  - a generic Result 
  - as well as specific versions for submodules, such as io::Result
```
- Custom Message & Panic with `.expect("something bad happened")`
- **!DOUBT** - Does Rust require a separate BDD library when it has its own `expect` & `Result` enum?
- **!Amit Says** - Rust tries hard to make code look like a series of function calls with `.` separator
- **!Amit Says** - In Go we have `%s` `%t` `%v` `%#v` & so on. In Rust use `{}`
- `::` is used to refers to functions or traits directly without instantiation. Remember static calls in Java ?
- `:` is used for type casting
- **!Amit Says** - One line says its all. Assignment and Error Handling in One Line
```
let guess: u32 = guess.trim().parse().expect("Please type a number!");
```
- `;` is back
- Rust at its best to push everything into a single line. Below is a switch case like code block:
```
match guess.cmp(&secret_number) {
  Ordering::Less => println!("Too small!"),
  Ordering::Greater => println!("Too big!"),
  Ordering::Equal => println!("You win!"),
}
```
- Block of code within `{}` separated by comma `,`
- **DOUBT** - Don't we need a semi-colon `;` within a block `{}`?
- Semi-Colons `;` Commas ',' Nested Blocks `{}` inside a Block `{}`
```rust
match guess.cmp(&secret_number) {
  Ordering::Less => println!("Too small!"),
  Ordering::Greater => println!("Too big!"),
  Ordering::Equal => {
    println!("You win!");
    break;
  }
}
```
- Map like operator via `=>`
- **!Amit says** - Everything under the sun has been used in Rust
- A sample error handling within a loop
```rust
let guess: u32 = match guess.trim().parse() {
  Ok(num) => num,
  Err(_) => continue,
};
```
- **!Amit says** - Rust tries its best to flatten which will otherwise be nested code blocks in other languages
- In Go underscore `_` is used to ignore a return value. In Rust `_` is used to catch anything i.e. represent anything.

### Mental Model - Day 2
- mutable versus. immutable
> In cases where you’re using large data structures, mutating an instance in place may be faster than copying and returning newly allocated instances. With smaller data structures, creating new instances and writing in a more functional programming style may be easier to think through, so lower performance might be a worthwhile penalty for gaining that clarity.
- constants: Look how readability is considered.
```rust
const MAX_POINTS: u32 = 100_000;
```
- Use of underscore `_` for readability. They name it as `visual separator`.
- Shadowing. When is it useful?
```rust
let spaces = "   ";
let spaces = spaces.len();
```
- Shadowing. No need to think of different variable names. Great !!!
- Data types in rust can be: 1/ Scalar or 2/ Compound
- Scalar - Integer, floating point, numbers, characters, booleans
- u32 - is a integer scalar data type - is unsigned - will take 32 bits of space
- Similarly u8, u16, u32, u64, u128, usize && i8, i16, i32, i64, i128, isize
- usize & isize will result into either 32 bit or 64 bit based on the architecture program is running
- Are you unsure of integer type?
> So how do you know which type of integer to use? If you’re unsure, Rust’s defaults are generally good choices, and integer types default to i32: this type is generally the fastest, even on 64-bit systems. The primary situation in which you’d use isize or usize is when indexing some sort of collection.
- 
- What is the default in floating point?
> The default type is f64 because on modern CPUs it’s roughly the same speed as f32 but is capable of more precision.
- Tuples & destructing them as well. I guess `destructuring` is common in functional languages.
```rust
   let tup = (500, 6.4, 1);
   let (x, y, z) = tup;
   println!("The value of y is: {}", y);
```
- Destructing is made possible via pattern matching. This is all internal. A developer need not worry about patterns.
- One can make use of dot i.e. period `.` followed by the index postion to access an element of tuple
- Tuple vs. Array - Array demands each element be of same data type
- Arrays in Rust are fixed length
- Arrays in Rust have the data collected in stack than heap
- You may be wanting a Vector that can grow & shrink in size
```rust
let a = [1, 2, 3, 4, 5];
```
```rust
let a: [i32; 5] = [1, 2, 3, 4, 5];
```

### Mental Model - Day 3


### Mental Model - Ownership
- https://medium.com/@thomascountz/ownership-in-rust-part-1-112036b1126b
  - Things get more interesting when we start passing around values
  - Switching from using a **string literal**, which is stored on the **stack**, 
  - to using a **String type**, which is stored on the **heap**
- _string literals are typically copied by value_
- _string type are moved by data_
  - accessing a moved data will error out
- Elaborate Explanation
  - string literal is stored somewhere in read only memory & not in stack or heap
  - pointer to string literal is stored in stack
  - Because it’s a string literal, it usually shows up as a reference, 
  - meaning that we use a pointer to a string stored in permanent memory, 
  - and it’s guaranteed to be valid for the duration of the entire program, (it has a static lifetime).
- https://medium.com/@thomascountz/ownership-in-rust-part-2-c3e1da89956e
  - Pass by Reference is referred to as borrowing
    - Why has borrowing come up as a concept?
    - So that data can be passed without passing over the ownershiphttps://medium.chttps://medium.com/capital-one-tech/building-an-event-sourcing-crate-for-rust-2c4294eea165om/capital-one-tech/building-an-event-sourcing-crate-for-rust-2c4294eea165
    - Ownership still lies with the originating source
    - Hence there is no concept of returning ownership since it was never passed
    - Ownership implies the responsibility of deallocating space from memory belongs to original source

### Mental Model - 20-Nov-2018
- https://medium.com/capital-one-tech/event-sourcing-with-aggregates-in-rust-4022af41cf67

#### Mental Model - Events & Event Source
- You can build up the state by aggregating the events! also known as _re-hydrate_
- Events seem to be PAST TENSE VERBS
- e.g. RESERVED, RELEASED, SHIPPED
- **!Amit Says** Actions seem to be present tense
- e.g. RESERVE, RELEASE, SHIP

#### Mental Model - Structs as Enum
- Ask. Do you want to pattern match on content?
- Ask. Do you want to pattern match on type?
- _Take I_
- Typically we do it by creating a custom defined type & creating various constants out of it
- In other words, we create ENUMs when presented with such scenarios
- Then this is used as a field in some struct
- Finally our entire code base is littered with if else or switch case conditions
- _Take II_
- How about defining a enum out of struct?
- In other words, ENUM around the entire content?
```rust
struct ProductEventData {
    quantity: usize,
    timestamp: usize,
    sku: String,
}

enum ProductEvent {
    Reserved(ProductEventData),
    Released(ProductEventData),
    Shipped(ProductEventData),
}
```
- With above you can treat different kinds of product events as a single type
- Is above a great thing to have?

#### Mental Model - Aggregate
- Aggregate represents a _calculated state_
- So aggregate is basically a state
- Aggregate is the state that a command executes
- Internals - Command is an imperative to the aggregate for validation
- Internals - The response of this imperative is a list of eligible events ready for emission
- Internals - Some aggregates can emit the events directly. However this is not unit testable
- Internals - Good Practice - Event emission should be a separate concern
- **!Amit Says** - Command is what I referred to as Action - Imperative - Authoritative - Bossy
- Here comes the twist - Define aggregate
- Aggregate should be able to APPLY events SEQUENTIALLY to their state
- It looks like `func(state1 + event) = state2`
- Aplying event to an aggregate produces an aggregate with a new state
- **!Amit Says** - Hey! Should an aggregate be defined as a struct?
- Why do we call it as an aggregate?
- Is it an aggregate of events?
- In Rust you can define a trait that describles an aggregate of events
```rust
trait Aggregate {
    type Item;

    fn version(&self) -> u64;
    fn apply(&self, evt: &Self::Item) -> Self where Self:Sized;
}
```
- Above means:
  - Aggregates must have a version
  - Aggregates must have an apply method
  - Self implies the real thing & not the trait type
  - Self:Sized implies more predictable memory footprint & not something like a boxed trait object
- NOTE: Aggregates are calculations for a single entity, not for the entire state of your application
- Aggregates are short lived
- Aggregates live long enough to calculate state & validate a command & that's it

#### Mental Model - Entity
- We still need to define the entity on which aggregate is bound
```rust
#[derive(Debug)]
struct Product {
    version: usize,
    qty_on_hand: usize,
    sku: String,
}
```

### Mental Model - ?date?
- https://medium.com/capital-one-tech/building-an-event-sourcing-crate-for-rust-2c4294eea165


## Simple
- package
- file name
- var name
- error handling
- pointers
- pass by value
- pass by reference
- import
- data types

## Intermediate
- polymorphism
- functional

## Advanced
- heap
- stack
- garbage collection
