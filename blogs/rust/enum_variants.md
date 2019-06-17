### References:
- https://brson.github.io/rust-anthology/1/understanding-over-guesswork.html

### Notes:
- Represent variants with _encapsulated values_, _generics_, and even _structs_!
- Query: Is this good? It seems to be similar to a structure. Should we use a struct instead?

```rust
// Structure with generic
struct One<T> {
    foo: usize,
    bar: T
}
// 2-tuple
struct Two(usize, usize);
// Enum
enum Three {
    // Plain.
    Foo,
    // Variant with Tuple.
    Bar(usize),
    // Variant with Struct.
    Baz { x: u64, y: u64, z: u64, },
}
```
