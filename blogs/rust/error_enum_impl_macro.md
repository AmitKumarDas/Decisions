```rust
// code is in macros
// error
// enum
#[derive(Debug, Fail)]
pub enum NvmeError {
  #[fail(display = "IO error: {}", error)]
  IoError { error: io::Error },

  #[fail(display = "nqn: {} not found", _0)]
  NqnNotFound(String),

  #[fail(display = "controller with nqn: {} not found", _0)]
  CtlNotFound(String),

  #[fail(display = "no nvmf subsystems found")]
  NoSubsystems,
}

// implement an enum -- powerful
// does this implement the macro of NvmeError?
// can macro be modeled as an interface?
impl From<io::Error> for NvmeError {
  fn from(err: io::Error) -> NvmeError {
    NvmeError::IoError {
      error: err,
    }
  }
}
```
