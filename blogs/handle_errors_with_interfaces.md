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
  - Dealing with structs is not bad but interfaces if possible will save us a lot of churn i.e. technical debt
- Compose interface(s) in struct(s)
- Pass interface(s) as arguments

### Functional Outlook
Let us relook at above snippet with more functional outlook.

```go
// Processor is a contract to process buffers
type Processor interface {
	Process(buf []byte) (message string, err error)
}

// ProcessorFn is a function based implementation of Processor interface
type ProcessorFn func(buf []byte) (message string, err error)

// Process is an implementation of Processor interface
func (fn ProcessorFn) Process(buf []byte) (message string, err error) {
	// concrete logic is not implemented here
	// just invoke **this** i.e. the function
	return fn(buf)
}

// ErrorHandler is a contract to handle error
type ErrorHandler interface {
	Handle(err error)
}

// ErrorHandlerFn implements ErrorHandler interface
type ErrorHandlerFn func(err error)

// Handle is an implementation of ErrorHandler interface
func (fn ErrorHandlerFn) Handle(err error) {
	// concrete logic is not implemented here
	// just invoke **this** i.e. the function
	return fn(err)
}

// sendErrFn sends error to an external system
//
// NOTE: This is a concrete implementation
var sendErrFn = ErrorHandlerFn(
	func(err error) {
		external.Send(err)
	}
)

// abcProcessorFn is a specific implementation of Processor
//
// NOTE: assume 'abc' as a 3rd party lib
//
// NOTE: This is a concrete implementation
var abcProcessorFn = ProcessorFn(
	func(buf []byte) (message string, err error){
		return abc.Buffer(buf)
	}
)

// abcProcessorWithErrSender is a ABC based Processor with error 
// sent to an external system
var abcProcessorWithErrSender = ProcessorFn(
	func(buf []byte) (message string, err error){
		message, err = abcProcessorFn()(buf)
		if err != nil {
			sendErrFn()(err)
		}
		
		return
	}
)

// processorWithErrHandler is a generic processor with error handler
func processorWithErrHandler(buf []byte, processor Processor, handler ErrorHandler) (string, error) {
	msg, err := processor.Process(buf)
	if err != nil {
		handler.Handle(err)
	}
	return msg, err
}

main() {
	// seems over-engineering
	buf = []byte("i have some bytes")
	msg, err := abcProcessorWithErrSender()(buf)
	
	// vs.
	// seems simple
	errFn := sendErrFn()
	msg, err := abcProcessor()(buf)
	errFn(err)
	
	// vs.
	// looks to be simpler
	// same as the first snippet with functional implementation of interfaces
	msg, err = processorWithErrHandler(buf, abcProcessorFn, sendErrFn)
}
```

### Choose your design carefully
While interfaces are good, it seems that building functions that implement interfaces seem to reduce lines of code to some
extent _i.e. no need to create specific structure for specific implementation of the interface_. However, minute 
modularization of logic seems to be an act of over engineering that should be avoided.
