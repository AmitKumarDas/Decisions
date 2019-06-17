### Crash Learning
```rust
fn example_simple()
fn example_params(x: u64, y: &u64, z: &mut u64)
fn example_returns(x: u64) -> u64
fn example_generic<U: Read>(reader: U) -> u64
fn example_generic_alt<U>(reader: U) -> u64
    where U: Read
```


```rust
// amitd: initial learnings in rust
// tag: fmt, format

// - we implement trait for struct
// - fn returns are indicated by arrow i.e. ->
// - :: is used between package::Struct
// - package is small cased
// - Struct is capital cased
// - macro end with !
// - line of code end with ;
// - Ok(()) does not end with ;
// - What is Ok(())
// - line of code can also end with ? and then with ;
// - impl Display is like implement you own custom fmt::Formatter
// - impl Display needs a mutable fmt::Formatter
// -- f: &mut fmt::Formatter
// - &self as the first arg determines this as a method
// - method chaining is cool as always
// -- self.children.iter().enumerate()
// -- Notice these methods are small cased
// - fmt is also a small cased Trait contract
// - during format {} is essential. Do not forget {}

impl Display for NexusBdev {
  fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
    writeln!(f, "uuid: {}\n", self.uuid)?;
    for (i, child) in self.children.iter().enumerate() {
        writeln!(f, "\t child[{index}]: {uri}\n", index = i, uri = child)?;
    }
    Ok(())
  }
}
```
```rust
// {} is essential
// = is confusing; but lets try to understand
// -- left uuid is used in fmt
// -- right uuid is the one that is passed as arg

pub fn get_storage_service_nqn(uuid: &str) -> String {
  format!("nqn.2019-05.io.openebs:{uuid}", uuid = uuid)
}
```
