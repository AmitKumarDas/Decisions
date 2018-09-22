### Logic that takes care of boilerplate
- This is yet another series on code review
- The review is against a well known lib, i.e. **github.com/pkg/errors**

```go
// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}
```

Author(s) of above lib understand that type checking the error & then eextracting the root cause will be spread everywhere a 
library tries to make use of this `errors` package. So they come up with this brilliant piece i.e. a function that handles
both interface definition & getting out the root cause. The caller code will do something below now:

```go
switch err := errors.Cause(err).(type) {
case *MyError:
  // handle specifically
default:
  // unknown error
}
```
