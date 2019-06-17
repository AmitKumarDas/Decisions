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
