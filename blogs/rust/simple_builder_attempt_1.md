```rust
pub struct Pod {
  // fields
}

impl Default for Pod {
  fn default() -> Self 
  {
    Pod {
      // set the fields with defaults
    }
  }
}

impl Pod {
  pub fn new() -> Pod 
  {
    Default::default()  
  }
  
  pub fn with_name(&mut self, name: &str) -> &mut Pod
  {
    self.name = name;
    self
  }
  
  pub fn with_namespace(&mut self, namespace: &str) -> &mut Pod
  {
    self.namespace = namespace;
    self
  }
  
  pub fn build(&self) -> Pod
  {
    self.clone()
  }
}
```
```rust

let pod = Pod::new()
            .with_name("my-jiva-pod")
            .with_namespace("openebs")
            .build();
```
