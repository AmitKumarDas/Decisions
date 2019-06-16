### References:
- https://github.com/nrc/r4cppp/blob/master/arrays.md


### Observations & / Notes
- Arrays can be borrowed by **Reference**
- Arrays are fixed size
- When arrays need to be _growable_. It is known as **Vector**

### Do you remember **SYNTAX**?
- Borrow & Fixed - &, datatype, size - &[i32, 4]
- Borrow & Not Fixed - &, datatype - &[i32]
- Define - value, size - [1_i32, 2, 3, 4]
- Borrow is represented as `&`
- Pointer is represented as `*`

### Now Please Shout
- _Arrays are **value types**_
- _Allocated on the **stack** like other values_
- An array object is a sequence of values

### Array object is NOT A POINTER
- An array object **is not a pointer to those values** (as in C).
- `let a = [1_i32, 2, 3, 4];`
  - will allocate 16 bytes on the **stack** 
  - and executing let b = a; will **copy** 16 bytes. 

### Coercion
```rust
let a: &[i32] = &[1, 2, 3, 4];
```
- How do you explain above coercion?
- RHS is a fixed length array of length four, allocated on the stack.
- We then take a reference to it (type &[i32; 4])
- That reference is **coerced** to type &[i32]

### Acess is bounds checked
- You can also check the length yourself by using len()
  - So clearly the length of the array is known somewhere. 
- **Bounds checking**, which is an integral part of **memory safety**. 
- The size is known dynamically (as opposed to statically in the case of fixed length arrays),
  - Hence slice types are dynamically sized types **(DSTs)**
  - _**FUN:** !!! DST ~ Dessert !!! You want to know Size of Your Dessert !!!_

### How's of DST
- Slice is just a sequence of values,
- Is the size stored as part of the slice?
  - **NO**
- Instead it is stored as part of the pointer
  - **Part of the Pointer?**
  - Remember that slices must always exist as pointer types
  - A pointer to a slice (like all pointers to DSTs) is a **fat pointer**
  - It is **two words** wide, rather than one, 
  - Contains the pointer to the **data** + **a payload**
  - In the case of slices, the payload is the length of the slice


### If you want a C-like array
- You need to explicitly make a pointer to the array
- This will **give you a Pointer to the FIRST element**

### Wow !!! Arrays and Traits !!!
- Rust arrays can implement traits
- And thus **Arrays can have Methods**

### This will error
```rust
fn foo(x: [i32])
```

- Since compiler does know the size & hence cannot allocate the value
- Solutions :
  - Smart Pointers
  - Borrowed Reference `fn foo(x : &[i32])`
  - Mutable Raw Pointer `fn foo(x : *mut [i32])`

### Samples
```rust
fn foo(a: &[i32; 4]) {
  println!("First: {}; last: {}", a[0], a[3]);
}

fn main() {
  foo(&[1, 2, 3, 4]);
}
```

