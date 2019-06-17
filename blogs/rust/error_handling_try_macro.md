### Notes:
- When working with functions which may return a Result<T, E> it is common to use the try!() macro.
- This macro expands to either unwrap the T value inside and assign it, or return the error up the call stack.
- This helps reduce visual 'noise' and assist in composition.

```rust
pub enum MyError {
  Io(io::Error),
  Utf8(Utf8Error)
}

impl From<io::Error> for MyError {
  fn from(err: io::Error) -> MyError {
    MyError::Io(err)
  }
}

fn open_and_read() -> Result<String, MyError> {
  let mut f = try!(File::open("foo.txt"));
  let mut s = String::new();
  let num_read = try!(f.read_to_string(&mut s));
  Ok(s)
}
```
