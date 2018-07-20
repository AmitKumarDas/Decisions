This is an extract from go standard library. Let us try to understand this snippet:

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

What Coding Practices Does Above Teach?
- Think through interfaces versus jumping into structs
- Interface ensures adhering to write logic in an idiomatic way
- Utility functions are better off with interfaces or primitives as arguments
- Compose interface(s) in struct(s)
- Pass interface(s) as arguments
