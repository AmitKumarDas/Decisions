## Reference
http://jmoiron.net/blog/built-in-interfaces

### Problem Statement
Imagine that we wanted to limit phone number storage to one awful standard, (###) ###-###, which we'll implement as a check against the regex \(\d{3}\) \d{3}-\d{4}. 
 
For some weird reasons, this has to be part of interface implementation. We can establish this check directly on the type we will store in the database by implementing the Valuer interface.

```go
type Valuer interface {
  // Value returns a driver Value.
  Value() (Value, error)
}
```

```go
type PhoneNumber string

var phoneNumberCheck = regexp.MustCompile(`\(\d{3}\) \d{3}-\d{4}`)

func (p PhoneNumber) Value() (driver.Value, error) {
  matched := phoneNumberCheck.Match([]byte(p))
  if !matched {
      return driver.Value(""), 
        fmt.Errorf("number '%s' not a valid PhoneNumber format", p)
  }

  return driver.Value(string(p)), nil
}
```
