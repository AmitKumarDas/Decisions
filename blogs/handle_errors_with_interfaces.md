### Handle errors with interfaces
Why not? Errors are crucial to proper functioning of business logic. Errors are integral part of business logic.

### Understand by a code snippet
This is a go snippet code extracted from [link](http://redhat-crypto.gitlab.io/defensive-coding-guide/#chap-Defensive_Coding-Go). Let us try to understand this snippet:

```go
type Processor interface {
	Process(buf []byte) (message string, err error)
}

type ErrorHandler interface {
	Handle(err error)
}

func RegularError(buf []byte, processor Processor,
	handler ErrorHandler) (message string, err error) {
	message, err = processor.Process(buf)
	if err != nil {
		handler.Handle(err)
		return "", err
	}
	return
}
```

### What Coding Practices Does Above Teach?
- _**Think through interfaces**_ versus _**world of structs only**_
- Interface ensures adhering to write logic in an idiomatic way
- Utility functions are better off with interfaces or primitives as arguments
  - Think twice if you are passing or returning structures
  - Dealing with structs is not bad but interfaces if possible will save us a lot of trouble
- Compose interface(s) in struct(s)
- Pass interface(s) as arguments

### Concluding with a snippet
