### References:
- https://brson.github.io/rust-anthology/1/understanding-over-guesswork.html

### Notes:
- An io::Error and a Utf8Error are different types and cannot be returned in the same Result<T,E>
  - Since the E value would differ and violate Rust's strong typing
- This is typically solved by creating a new Error
  - Which is an enumeration over the possible errors
- Then there are the Into<T> and From<T> traits which can be implemented to provide seamless interaction.

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
```
