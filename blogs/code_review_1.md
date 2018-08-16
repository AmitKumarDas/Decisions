## Code Review
This is first of the series of articles tagged under _Code Review_. The purpose of writing a code review of any open source 
code is to learn coding from real world projects. Learning here involves writing down the same source code with comments of 
my own, which in turn will negate any first conclusion bias-es I might chance upon during my first code reading.

As we shall see, there are a lot of things that can be learnt from reading, reviewing & subsequent writing down the review
analysis than by just reading the best practices and writing software logic for an project. Have a careful look at the 
code comments titled under Notes: section. These _Notes:_ are actually my review.

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
```
