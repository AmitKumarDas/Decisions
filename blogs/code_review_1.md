## Code Review
This is first of the series of articles tagged under _Code Review_. The purpose of writing a code review of any open source 
code is to learn coding from real world projects. Learning here involves writing down the same source code with comments of 
my own, which in turn will negate any first conclusion bias-es I might chance upon during my first code reading.

As we shall see, there are a lot of things that can be learnt from reading, reviewing & subsequent writing down the review
analysis than by just reading the best practices and writing software logic for an project. Have a careful look at the 
code comments titled under _Notes:_ section. These `// Notes:` are actually my review comments.

### Show me the code !!
I shall do a code review of https://github.com/pkg/errors/ project.

```go
package errors

// New returns an error with the supplied message.
// New also records the stack trace at the point it was called.
//
// Notes:
// - Naming seems to be perfect. Only one New across the package. Intuitive.
// - Returns an interface than a struct. This follows good practice.
// - Instantiation does not mandate the caller to learn a lot of things i.e. one simple argument.
func New(message string) error {
	return &fundamental{
		msg:   message,
		stack: callers(),
	}
}

// fundamental is an error that has a message and a stack, but no caller.
//
// Notes:
// - There is a reason to make this private to the package
// - Composes stack struct's pointer
// - Indicates a specific usage of stack
type fundamental struct {
	msg string
	*stack
}

// Notes:
// - Implementation of error interface. Good.
func (f *fundamental) Error() string { return f.msg }

// Notes:
// - Implementation of Formatter interface. Good.
func (f *fundamental) Format(s fmt.State, verb rune) {
  // ... code goes here
}

// WithStack annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
//
// Notes:
// - Naming is apt w.r.t instantiation of withStack struct
// - Caller need not bother about multiple input arguments
// - Interesting things are all handled in this instantiation function
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	return &withStack{
		err,
		callers(),
	}
}

// Notes:
// - Composes stack struct
// - Is a specific implementation of stack struct pointer
// - Is a specific implementation to multiple interfaces
type withStack struct {
	error
	*stack
}

// Notes:
// - Implementation to Causer interface
func (w *withStack) Cause() error { return w.error }

// Notes:
// - Implementation to Formatter interface
func (w *withStack) Format(s fmt.State, verb rune) {
	// code goes here
}

// WithMessage annotates err with a new message.
// If err is nil, WithMessage returns nil.
//
// Notes:
// - Naming is apt w.r.t corresponding struct that is instantiated
// - Callers are abstracted away from internals & only need to provide arguments that they understand
// - Returns an interface. This is idiomatic go code
func WithMessage(err error, message string) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   message,
	}
}

type withMessage struct {
	cause error
	msg   string
}

// Notes:
// - This is an implementation of error interface
// - It is a good practice to code to interface
func (w *withMessage) Error() string { return w.msg + ": " + w.cause.Error() }

// Notes:
// - This is an implementation of Causer interface
// - It is a good practice to code to interface
func (w *withMessage) Cause() error  { return w.cause }

// Notes:
// - This is an implementation of Formatter interface
// - It is a good practice to code to interface
func (w *withMessage) Format(s fmt.State, verb rune) {
	// code goes here
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
//
// Notes:
// - This is the major piece that is marketed to the developers
// - No learning curve involved when things become simple by just talking about a single function
// - Callers need not bother other than 2 arguments that they understand & have full control
// - Interesting logic e.g. withMessage & withStack are all taken care of
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	err = &withMessage{
		cause: err,
		msg:   message,
	}
	return &withStack{
		err,
		callers(),
	}
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is call, and the format specifier.
// If err is nil, Wrapf returns nil.
//
// Notes:
// - I wanted this more often than Wrap, but never knew this one existed
// - Naming is go idiomatic. No need to spend hours on a name & yet be precise in naming
// - Callers need not bother about providing arguments that is specific to this package
// - All interesting things are taken care of in this function
// - In other words, callers are abstracted from the internals in this package
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	err = &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
	return &withStack{
		err,
		callers(),
	}
}
```
