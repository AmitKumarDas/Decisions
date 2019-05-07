```rust
impl Debug for File {
  fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
    f.debug_struct("File")
      .field("pos", &self.pos)
      .field("rdr", &self.rdr)
      .field("wtr", &self.wtr)
      .field("can_read", &self.can_read)
      .field("can_write", &self.can_write)
      .finish()
  }
}
```
