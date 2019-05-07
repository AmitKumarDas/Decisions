```rust
// Questions:
//
// How do we inherit?
// How is polymorphism implemented?

#[derive(Debug)]
pub enum Error {
    RefOverflow,
    RefUnderflow,
}

impl Display for Error {
  fn fmt(&self, f: &mut Formatter) -> fmt::Result {
    match *self {
      Error::RefOverflow => write!(f, "Refcnt overflow"),
      Error::RefUnderflow => write!(f, "Refcnt underflow"),
    }
  }
}

impl StdError for Error {
  fn description(&self) -> &str {
    match *self {
      Error::RefOverflow => "Refcnt overflow",
      Error::RefUnderflow => "Refcnt underflow",
    }
  }
}
```
```rust
//
// Hawk Eyes
// =>
// ->
// :
// ::
// !
// &self
// &mut
// &str
```
