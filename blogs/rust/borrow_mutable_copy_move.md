### References:
- https://brson.github.io/rust-anthology/1/understanding-over-guesswork.html

### Notes:
- Borrowing types:
  - Immutably borrowing (&), 
  - Mutably borrowing (&mut), 
  - Copying (Copy trait),
  - Moving values
- At any given time there may be **any number of immutable borrows**
- Meanwhile there may only be **one mutable borrow**
- and a value _may not be used_ in the function _once it has been moved out_

### Shout With Me
- !!! Only one mutable borrow !!!

```rust
fn main() {
  // An owned, growable,
  // non-copyable string.
  let mut foo = String::from("foo");

  // Introduce a new scope
  {
    // Reference bar is created.
    let bar = &foo;
    // Error, bar is immutable.
    bar.push('c');
  } // bar is destroyed.

  // Error, bar does not exist.
  let baz = bar;

  // Works, reference and mutable; only one mutable borrow
  let rad = &mut foo;
  rad.push('c');
} // foo is destroyed.
```

### Further Details
- No worries about making sure each of their malloc() calls have a corresponding free()
- No need to rely on an outside tool (ref) to discover such errors
- The borrow checker determines when a value has been moved into a function call
  - and should not be further used in the caller
