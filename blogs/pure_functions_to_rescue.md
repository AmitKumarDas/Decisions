### Pure functions to rescue
This topic is about `first class functions` to rescue your programming logic. This is again a topic related to go based programming.

### First Class Functions versus Interfaces
I have written quite a good number of articles talking about advantages of using Interfaces in go. Is this one more that talks
about suitability of _**Interfaces over Functions**_? The answer is No! This is in fact the opposite i.e. we shall try to evaluate if we can do away with Interfaces & limit ourselves to first class functions i.e. `typed functions`. A reader
who has read all my go based blogs, will immediately judge my biasedness over defining typed function that implement a single method interface.

Let us take the same concept forward. To reiterate the concept we are talking about here is outlined in following lines:
- Define an interface that exposes only a single method
- Define a typed function that implements above interface
- This leads to reduction in lines of code by avoiding specialized structs to implement above interface

### Changing the ingredients
Like I said earlier, we shall not define any interface this time. We shall start with typed functions instead. Let us take a
look at code snippet to get more clarity:

```go
// a data structure for illustration purposes
type data struct {
  daytime bool
}

// action defines a typed function
//
// NOTE:
//  This is similar to a contract signature of an interface
type action func(given *data) (response *data, ok bool)

// actOne is a specific action 
//
// NOTE:
//  This has the concrete implementation itself unlike a 
// typed function that implements a interface
func actOne(given *data) (response *data, ok bool) {
  response = given
  ok = true
  return
}

// actNoop is a specific action
//
// NOTE:
//  This has the concrete implementation itself unlike a 
// typed function that implements a interface
func actNoop(given *data) (response *data, ok bool) {
  return
}

// strategy is a typed function that can be used to choose an action
type strategy func(given *data) action

// actBasedOnTime is a specific strategy
//
// NOTE:
//  This has the concrete implementation itself unlike a 
// typed function that implements a interface
//
// NOTE:
//  Anonymous function can be returned 
type actBasedOnTime(time time.Date) strategy {
  return func(given *data) action {
    if time.IsDay() && given.daytime {
      return actOne
    }
    return actNoop
  }
}

func main() {
  data = &data{}
  actor := actBasedOnTime(time.New())(data)
  resp, ok := actor(data)
}
```

### What did we see?
In my opinion, I find this to be a lot more appealing than writing a interface & necessary typed function adhering those
interfaces. In short a lot less code.

### References
- https://golang.org/doc/codewalk/functions/
