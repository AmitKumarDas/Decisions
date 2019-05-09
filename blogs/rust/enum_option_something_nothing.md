Reference - https://doc.rust-lang.org/stable/book/ch06-00-enums.html

### TIL _(Today I Learned)_
- Kind as a suffix to name an enum

```rust
// Define a type by enumerating its possible values. 
// How an enum can encode meaning along with data. 
//
// A useful enum, called Option, which expresses that 
// a value can be either something or nothing. 
//
// How pattern matching in the match expression makes 
// it easy to run different code for different values
// of an enum. It is kind of map over each enum values.
//
// How the if let construct is another convenient and
// concise idiom available to you to **handle** enums in 
// your code.
//
// Rust’s enums are most similar to **algebraic data types**
// in functional languages, such as F#, OCaml, and Haskell.
```

```rust
// V4 & V6 are variants of enum
enum IpAddrKind {
    V4,
    V6,
}
```
