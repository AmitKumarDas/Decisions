## Rough
- install rust via rustup
```bash
$ curl https://sh.rustup.rs -sSf | sh
```
```bash
It will add the cargo, rustc, rustup and other commands to Cargo's bin 
directory, located at:

  /home/amit/.cargo/bin

This path will then be added to your PATH environment variable by modifying the
profile file located at:

  /home/amit/.profile

You can uninstall at any time with rustup self uninstall and these changes will
be reverted.
```
```bash

To get started you need Cargo's bin directory ($HOME/.cargo/bin) in your PATH 
environment variable. Next time you log in this will be done automatically.

To configure your current shell run source $HOME/.cargo/env
  ~ > source /home/amit/.cargo/env 
```
- update
```bash
If you’ve installed Rust via rustup, updating to the latest version is easy. 
From your shell, run the following update script:

$ rustup update
```
- linker
> Additionally, you’ll need a linker of some kind. It’s likely one is already installed, but when you try to compile a Rust program and get errors indicating that a linker could not execute, that means a linker isn’t installed on your system and you’ll need to install one manually. C compilers usually come with the correct linker. Check your platform’s documentation for how to install a C compiler.
- macro
```rust
println!("This is macro with !");
```
- Ahead-Of-Time compiled language
> Rust is an ahead-of-time compiled language, meaning you can compile a program and give the executable to someone else, and they can run it even without having Rust installed.
- dependency management ?
> cargo .. same as dep in go
- Does code compiles ?
```
cargo check
```
- Want to build?
```
cargo build
```


## Simple
- package
- file name
- var name
- error handling
- pointers
- pass by value
- pass by reference
- import
- data types

## Intermediate
- polymorphism
- functional

## Advanced
- heap
- stack
- garbage collection
