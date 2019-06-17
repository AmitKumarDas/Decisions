```rust
// Create a `Some(T)` and a None.
let maybe_foo = Some(0);
let not_foo = None;

// Unwrapping.
let foo = maybe_foo.unwrap();
let default = not_foo.unwrap_or(1);

// Match
let matched = match maybe_foo {
    Some(x) => x,
    None => -1,
};

// Map
let mapped = maybe_foo.map(|x| x as f64);
```
