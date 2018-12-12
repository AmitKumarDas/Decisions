### Reference
https://medium.com/capital-one-tech/event-sourcing-with-aggregates-in-rust-4022af41cf67

### Tags
- event source, aggregate, domain driven design
- functional, fold, trait, enum, struct, mutable, immutable

### Introduction
These are my interpretations of above reference. I have read and interpreted some of the stuff for my own understanding.
Most of the phrases used by the author has been borrowed. Note that I am a newbie in Rust. Writing & interpreting makes
me build my mental models to understand Rust programming.

### Article
I am a huge advocate for:
- immutable streams of data,
- event sourcing,
- separation of command and query,
- and the benefits of treating state as a function of an event stream. 
- However, given all of that, I’ve only been able to implement a handful of _pure event sourcing systems_.

#### Event Sourcing - What?
- Event sourcing is the idea that your shared state is not mutable in place. 
- It is the result of sequentially applying immutable events that represent something that took place in the past.
- Once you start envisioning problems as event sourcing problems, it’s very difficult to see them any other way.
- The aggregate is a domain-driven-design (DDD) concept that fits well within event sourcing. 
- To put it as briefly as possible: 
  - You apply a command to an aggregate which then produces one or more events. 
- An aggregate can populate (re-hydrate) its state by sequential application of an event stream.
- Q - Is Aggregate a state ?

#### Low Level Design in Rust
- The first question is:  how do you represent an event? 
- Identify the domain entity(-ies) first
- Events that can happen to a object, e.g. 
  - reserved (an order of objects), 
  - released (an order is canceled), or
  - shipped (previously reserved actually leaves the warehouse)
- How about representing Events as Structs?
  - EntityShipped
  - EntityReserved
  - EntityReleased
- Ponder - Does above get you into if-else or switch case statements?
- Idea - **How about enums that are also structs?**

```rust
// nothing drastically different than golang
struct ProductEventData {
  id: string,
  quantity: usize,
  time: usize,
}

// inversion of control in play
// can do this in golang by using struct as go does not have enum
// in golang struct will have Released, Reserved & Shipped as fields of some type
// representing enum items as fields within a struct is preferred approach in golang when the case demands
enum ProductEvent {
  Released(ProductEventData),
  Reserved(ProductEventData),
  Shipped(ProductEventData),
}
```

- Q - Is it advantageous to treat all kinds of entity events as a single type?
- A - This is handy when it comes to build aggregate (i.e. calculated state) ?

- NOTE:
  - Aggregate is a calculate state
  - Aggregate is something on which commands/actions/imperatives are executed
  - If commands succeed then aggregate will return a list of events ready for emission
  - Do not let aggregate return events. Maintain de-coupling. Think easier unit testing.

- What are these commands/actions?
  - Ship action that will return Shipped(entity) as its response
  - Ship might also return multiple Events as its response. Some domains have this most often.

- How do we aggregate?
  - Aggregate should be able to apply all the events sequentially to their state
  - for example:- f (state 1 + event) = state 2

- Idea - Can we have a trait that describes an aggregate of events?
- Q - Is trait like an interface?

```rust
trait Aggregate {
    type Item;

    fn version(&self) -> u64;
    fn apply(&self, evt: &Self::Item) -> Self where Self:Sized;
}
```
- Above says:
  - all aggregates must have a version
  - all aggregates must have a apply method
  - apply method must take an event of type Item
  - apply must produce the type that implements the trait not the trait itself
  - Self:Sized means more predictable memory footprint & not something like a boxed trait object

- Q - Can you keep using Self all around? 
- Q - Can they refer to a type implicitly?
- Q - Can they be referred to explicitly?

- Enums are powerful in rust
  - Implement functions with enum as the receiver
  - Here nothing significant is done, other than building the event object itself
```rust
impl ProductEvent {
    // naming should be reserved & not reserve
    // its kind of fetching a Reserved instance
    // this is not a imperative or action or command
    // Q - why sku is typed to &str ?
    //
    // naming cane be better e.g.:
    // new_reserved vs reserved
    // new_shipped vs shipped
    // new_released vs released
    fn reserved(sku: &str, qty: usize) -> ProductEvent {
        ProductEvent::Reserved(ProductEventData {
            quantity: qty,
            sku: sku.to_owned(),
            timestamp: 1,
        })
    }
    fn shipped(sku: &str, qty: usize) -> ProductEvent {
        ProductEvent::Shipped(ProductEventData {
            quantity: qty,
            sku: sku.to_owned(),
            timestamp: 1,
        })
    }
    fn released(sku: &str, qty: usize) -> ProductEvent {
        ProductEvent::Released(ProductEventData {
            quantity: qty,
            sku: sku.to_owned(),
            timestamp: 1,
        })
    }
}
```

- Define the ProductAggregate i.e. calculated state
  - Aggregates are calculations for a single entity, not for the entire state of your application. 
  - Further, aggregates are short-lived. 
  - They live long enough to calculate state and validate a command, and that’s it. 
  - If you’re building a stateless service, it’s going to dispose of the aggregate at the end of the request.

```rust
#[derive(Debug)]
struct ProductAggregate {
    version: usize,
    qty_on_hand: usize,
    sku: String,
}

// fn are inside impl
// invert to golang ways to internalize
// code will look nested though
//
// impl against Enum
// impl against struct
//
// func vs. fn
// Q - public vs. private. How? When?
impl ProductAggregate {
    // seems like auto returns
    //
    // instantiation function
    fn new(sku: &str, initial_quantity: usize) -> ProductAggregate {
        ProductAggregate {
            version: 1,
            qty_on_hand: initial_quantity,
            sku: sku.to_owned(),
        }
    }
    
    // seems like auto returns
    // Err(...)
    // or Ok(...)
    //
    // Err is inferred to be of string type
    // Ok is interred to be of vector of events
    //
    // Hey who is calling this fn
    //
    // Have a look at format!
    // Have a look at vec!
    //
    // It returns the events after all
    // Do not get deviated by the Result & vec
    //
    // just a validation & no mutation
    fn reserve_quantity(&self, qty: usize) -> Result<Vec<ProductEvent>, String> {
        if qty > self.qty_on_hand {
            let msg = format!(
                "Cannot reserve more than on hand quantity ({})",
                self.qty_on_hand
            );
            Err(msg)
        } else if self.version == 0 {
            Err(
                "Cannot apply a command to an un-initialized aggregate. Did you forget something?"
                    .to_owned(),
            )
        } else {
            Ok(vec![ProductEvent::reserved(&self.sku, qty)])
        }
    }

    // this returns Released instance of the event
    // no mutation
    fn release_quantity(&self, qty: usize) -> Result<Vec<ProductEvent>, String> {
        Ok(vec![ProductEvent::released(&self.sku, qty)])
    }

    // this returns Shipped instance of the event
    // no mutation
    fn ship_quantity(&self, qty: usize) -> Result<Vec<ProductEvent>, String> {
        Ok(vec![ProductEvent::shipped(&self.sku, qty)])
    }

    // just an accessor
    // no mutation
    fn quantity(&self) -> usize {
        self.qty_on_hand
    }
}
```

- Functional Way via traits !!
  - Ability for the aggregate to validate incoming commands relies on the fact that it already has computed state.
  - To compute state, we pump an event stream through the apply method. 
  - A lot of people like mutable aggregates here, and they call apply repeatedly on the same aggregate.
  - However, mutable aggregates leads to difficult to diagnose problems in production 
  - Lets take a more functional approach and have apply method return a **brand new aggregate**, with the computed state:
    - here functional is also referring to the cause of immutability

```rust

trait Aggregate {
    type Item;

    fn version(&self) -> u64;
    fn apply(&self, evt: &Self::Item) -> Self where Self:Sized;
}

#[derive(Debug)]
struct ProductAggregate {
    version: usize,
    qty_on_hand: usize,
    sku: String,
}

// impl **trait** for **struct**
//
// ProductAggregate struct had a number of functions attached to it as impl
// struct functions were not mutating
//
// trait is similar to interface in golang
// Hey, does this not look like anonymous definitions
//
// Q - Should we really do this ?
// In other words define separate implementations for only for mutations. 
//
// Q - Does this result into more functional approach?
impl Aggregate for ProductAggregate {
    // concrete declaration
    type Item = ProductEvent;

    // self refers to ProductAggregate struct
    // look again - impl trait for struct
    fn version(&self) -> u64 {
        self.version
    }

    // self refers to ProductAggregate as impl trait for struct
    //
    // evt is bound to ProductEvent as we declared it sometimes back
    //
    // return type is bound to ProductAggregate as impl trait for struct
    //
    // &self, &Self & Self are confusing
    //
    // we still have the switch case even with the enum
    // so is there really any advantage ?
    //
    // Did you notice the de-construction of specific event the moment their is a match
    // Well, is this the reason, we built events as enums made out of struct? 
    // So that we can deconstruct & enjoy the rusty code
    // => reminds me of the lambda expression debates in Java 8 or so
    //
    // Did you notice the default match which is underscore _
    //
    // What the hell is clone() doing here? Why?
    // Note: self refers to ProductAggregate, self is passed as input
    // it is copied and re-allocated as a new ProductAggregate and this
    // new re-allocated thing is retured. This has lot of performance issues.
    fn apply(&self, evt: &ProductEvent) -> ProductAggregate {
      ProductAggregate {
        version: self.version + 1,
        qty_on_hand: match evt {
           &ProductEvent::Released(
              ProductEventData { quantity: qty, .. }) => {
                  self.qty_on_hand + qty
           },
           &ProductEvent::Reserved(
              ProductEventData { quantity: qty, .. }) => {
                    self.qty_on_hand - qty
           },
            _ => self.qty_on_hand,
          },
        sku: self.sku.clone(),
       }
   }
}
```

- Improvements
  - It was pointed out to me that there’s an even better way to do the apply. 
  - In my paranoid revolt against mutability, I overlooked Rust’s concept of the **move**. 
  - In Rust, when you assign something, the default behavior is to move the value from one place to the other (leaving the previous location unusable as a variable). 
  - **During this move, it’s safe (and often preferred) to do mutations because you’re guaranteed that no other code is referring to the thing you’re mutating at that moment.**
  - This is enforced at compile-time. Cool. Wow
  - Embracing the move and mutating the aggregate inside the apply avoids an extra allocation on the stack and avoids the call to clone() on the SKU string

```rust
// mut has come up for self
//
// so many functions with defaulting behaviour
// the panic inducing (i.e. overflow) math operation into an Option type
// defaults to 0
//
// self is the one that is returned
// there is no copy & allocate but a new immutable self
fn apply(mut self, evt: &ProductEvent) -> ProductAggregate {
  self.version += 1;
  self.qty_on_hand = match evt {
    &ProductEvent::Released(
       ProductEventData{quantity:qty,..}) => 
         self.qty_on_hand.checked_add(qty).unwrap_or(0),
    &ProductEvent::Reserved(
       ProductEventData{quantity:qty,..}) => 
         self.qty_on_hand.checked_sub(qty).unwrap_or(0),
    _ => self.qty_on_hand,
  };
  self
}
```

- Usage i.e. Caller code
```rust
// :: is used to represent type
// :: is used to invoke method calls
//
// There is no sign of Aggregate that deals with mutations
// Is that used by rust auto-magically?
//
// apply method takes a reference to the event
// it explicitly tells the consumer of the aggregate that it
// does not claim ownership of the event
let soap = ProductAggregate::new("SOAP", 100);
let u = soap.apply(&ProductEvent::reserved("SOAP", 10));

// if we use the second version i.e. efficient version of apply
// i.e. mut self
// and by mistake invoke soap.apply once again; it will
// result into rust's borrow-checker error
// i.e. value used here after move
let u = u.apply(&ProductEvent::reserved("SOAP", 30));
```

- More functional
  - The ulterior motive of using above designs was to make use of functional implementations
  - NOTE: repeated re-assignment of a thing to result into a new version of the thing
  - is the job of a **fold**

```rust
// convert a vector of events into an iterator
// & pipe them through a fold
// 
// these are otherwise known as streaming apis in some languages
// 
// Do you recongnize the anonymous function, its arguments & its body
// 
// Can we have a better name for example, apply_all
//
// Notice the instantiation right within the iterator
// it can perhaps also make use of variables declared earlier
//
// Did you notice even though the apply is a trait method it can be used
// in generic functional apis? Is this the difference between a trait and
// interface. In interface apply would have been a method receiver of some
// defined struct. 
// 
// However, IMO this too can be handled in golang's interface approach as well. 
// Just that we need to code everything to get this generic functional apis.
// In addition, since go does not support generics, the overall solution
// might become cumbersome.
let agg = v.iter()
    .fold(ProductAggregate::new("CEREAL", 100), |a, ref b| {
        a.apply(b)
    });
```

// same as above
//
// in rust, mutability is chosen by the consumer, not dictated by the struct
//
// init is the initial value given to fold i.e. first value
//
// b probably refers to the item that is being iterated i.e. ProductEvent
//
// Note that if we had used the earlier version of apply i.e. **copy & allocate**
// with lots of events then performance would have been worse
```rust
let agg = v.iter().fold(init, |a, ref b| a.apply(b));
```

### Conclusion
- I am not passing a mutable value throughout the fold. 
- Instead, I am returning a new, immutable value after each step. 
- At the end, the agg variable is immutable and will have all of my calculated state. 
- Because I’ve made all my command methods return events rather than mutate the aggregate itself, 
  - I can then issue a command directly on the agg variable without having to worry about changing mutability or
  - conflicting with Rust’s borrow checker.

### In My Opinion
- This is all about thinking in types/objects
- The more we can create distinct types and treat them as first class citizens, we can have more maintainable code
- Next time, you start coding think if you need to code it as a field of the struct or define it as a struct itself
- As this article suggests, you can isolate the mutating code as a separate implementation 
  - In other words, treat the mutation as a specific implementation of original struct i.e. entity
- The non-mutating logic can be functions of the direct struct i.e. entity itself
  - where each function returns an event or vector of events
- The original struct can also have a instantiation function & other accessor functions
- One can also treat the Events as ENUM wrapper over some data based structures
  - These wrapping over struct can be de-constructed later e.g. during switch case statements
  - It might be already clear that these events will be de-constructed typically inside the mutation implementation
- Caller code seems to more readable & simple due to the ability to use functional/streaming apis
