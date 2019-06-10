## Reference
http://jmoiron.net/blog/built-in-interfaces

### Problem Statement
Imagine that we wanted to limit phone number storage to one awful 
standard, (###) ###-###, which we'll implement as a check against the 
regex \(\d{3}\) \d{3}-\d{4}. 
 
For some weird reasons, this has to be part of interface implementation. We can
establish this check directly on the type we will store in the database by
implementing the Valuer interface.

```go
type Valuer interface {
  // Value returns a driver Value.
  Value() (Value, error)
}
```
