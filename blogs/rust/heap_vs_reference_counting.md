### Notes
- Your program runs on a stack
- This stack refers to memory allocated in heap

### Box vs Rc
- Box<T> -- Reference pointer
- Rc<T>  -- Reference counting pointer
- Above are known as **container types**
- Also known as **trait objects**

### When to use what
- Rc<T> is used when we want multiple methods using a **read only reference**
- Hence, **shared ownership** over some data
- It counts the uses of reference pointing to the same data on the heap
- When last reference is dropped, the data itself will be dropped & the memory is freed
- However, Rc<T> is not thread safe
- Rc<T> can not be passed between threads

### Their homes
- Box<T> & Rc<T> themselves live in the stack
- However, the data that Box & Rc contain live on heap
- So both are zero cost abstractions

### Summary
- Trait & TraitObject are different
- Trait are similar to interface
- Trait can define a contract
- Trait may not have any contract definitions
- TraitObject is a pointer to a concrete type implementing Trait
- TraitObjects are dynamically Sized Traits
- However, Rust needs to know everything at compile time about size of types
- TraitObject's size are determined via **dynamic dispatch** that maps it to the right implementation
  - How?
  - It is achieved using pointer types such as:
  - &, Box<T>, Rc<T>, Arc<T>

### ProTip
- `Box` & `Rc` are just references (pointers) living in stack that point to objects living on heap
- `Box` does not provide a shared ownership
- `Rc` provides a shared ownership
- `Rc` is not threadsafe
- Trait is similar to interface in other languages
