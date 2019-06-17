### Notes:
- Channels provide a way to **transfer messages (and ownership)** between threads 
  - Without fear of there being later (unsafe) access to the data by other threads
- The default channel provided by the standard library is a Multiple-Producer, Single-Consumer channel

```rust
use std::sync::mpsc::{channel, Sender, Receiver};
let (send, receive) = channel();
```

### !!!You don't Lock Code! You Lock the Data!!!
- **Locks can encapsulate data** such that access is only granted if the lock is held
- In Rust, you don't lock code, you lock data, and it is safer because of it
- Locks are typically represented by Mutexes
- Locks are shared between threads with an **Atomically Reference Counted structure** _(Arc)_
- Benefits:
  - This design of locking data prevents a lock from being acquired and never given up

```rust
use std::sync::{Arc, Mutex};
let data = Arc::new(Mutex::new(0));
```

### Further Notes:
- Traits like **Sync** and **Send** are implemented on types 
- Sync & Send symbolize if it can be sent or shared between threads safely
- These traits are not just documentation, they are intrinsic to the language

```rust
// Safe to share between threads.
use std::marker::Sync;
// Safe to transfer between threads.
use std::marker::Send;
```
