### References:
- https://github.com/nrc/r4cppp/blob/master/arrays.md
  - It helps to refer to above tutorial.
  - I have compiled above content as per my personal liking

### Observations & / Notes
- Arrays are fixed size
- Arrays can be **borrowed by reference** - and they now become slices

### Do you remember **SYNTAX**?
- Borrow & Fixed == &, datatype, size - &[i32, 4]
- Borrow & **Not Fixed** == _slice_ == &, datatype - &[i32]
- Define == value, size == [1_i32, 2, 3, 4]
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
  - In the case of slices, the **payload is the length of the slice**
  - _**FUN**: !!!Pay load is lengthy!!! like Bank EMIs that like to get paid long!!!_
- **QUERY:** !!!Is length same as size in terms of rust jargons!!!

### Fat Pointer for `let a: &[i32] = &[1, 2, 3, 4];`
- The pointer a will be 128 bits wide (on a 64 bit system) ~ _**2 words**_
- First 64 bits will store the address of the 1 in the sequence [1, 2, 3, 4], &
- Second 64 bits will contain 4 == payload == length
- Usually, as a Rust programmer, these fat pointers can just be treated as regular pointers
  - But it is good to know about (it can affect the things you can do with casts, for example)

### Range Syntax with Arrays
```rust
let a: [i32; 4] = [1, 2, 3, 4];
let b: &[i32] = &a;   // Slice of the whole array.
let c = &a[0..4];     // Another slice of the whole array, also has type &[i32].
let c = &a[1..3];     // The middle two elements, &[i32].
let c = &a[1..];      // The last three elements.
let c = &a[..3];      // The first three elements.
let c = &a[..];       // The whole array, again.
let c = &b[1..3];     // We can also slice a slice.
```
- **Trick: Left & Right Numbers** - Subtract the numbers to get the number of elements
- **Trick: One Sided Number** - left is excluded if left is set
- **Trick: One Sided Number** - right is included if right is set

### Amazing Ops w.r.t Range Syntax
- a..b produces an iterator which runs from a to b-1
- This can be combined with other iterators in the usual way, or can be used in for loops:

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

### Here comes the Vector
- Heap Allocated

### Vector
- is an owning reference
- Like `Box<_>` it has move semantics

### Three Friends - Array, Slice & Vector
- Array i.e. Fixed Sized - Value
- Slice - Borrowed Reference
- Vector - Box Pointer - Think it as a Smart Pointer

### Vector & Pointer
- Vector length is stored in a Pointer
- Here Pointer is Vec value itself
- Vector has no literals
  - However literal effect can be assumed via `vec!` macro
- We can create an empty vector by `Vec::new()`

### Understand these syntax
```rust
let v = vec![1, 2, 3, 4];      // A Vec<i32> with length 4.
let v: Vec<i32> = Vec::new();  // An empty vector of i32s.
```
- In the second case above, the type annotation is necessary
- Why?
  - so that the compiler can know what the vector is a vector of.
- If we were to use the vector, the type annotation would probably not be necessary.

### Vectors & Arrays - Similarity
- can be indexed
- can be ranged

### Vectors - Specialty
- Their size can change
- They can get longer or shorter as needed
- For example, **v.push(5)** would add the element 5 to the **end** of the vector
  - This would require that v is mutable
- Note that growing a vector can cause reallocation
  - **Safeguard:** Which for large vectors can mean a lot of copying.
- To guard against this you can **pre-allocate space** in a vector using **with_capacity**

### Samples
```rust
fn foo(a: &[i32; 4]) {
  println!("First: {}; last: {}", a[0], a[3]);
}

fn main() {
  foo(&[1, 2, 3, 4]);
}
```

