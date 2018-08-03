### Anonymous To Rescue
This topic is about `anonymous functions` to rescue your programming logic. This is again a topic related to go based programming.

### First Class Functions versus Interfaces
I have written quite a good number of articles talking about advantages of using Interfaces in go. Is this one more that talks
about suitability of _**Interfaces over Functions**_? The answer is No! This is in fact the opposite i.e. we shall try to evaluate if we can do away with Interfaces & limit ourselves to first class functions i.e. `custom typed functions`. A reader
who has read all my go based blogs, will immediately judge my biasedness over defining typed function that implement a single method interface.

Let us take the same concept forward. To reiterate the concept we are talking about here is:
- Define an interface that exposes only a single method
- Define a typed function that implements above interface
- This leads to reduction in lines of code by avoiding specialized structs to implement above interface

### Changing the ingredients
Like I said earlier, we shall not define any interface this time. We shall start with typed functions instead. Let us take a
look at code snippet to get more clarity:

```go

```

### References
- https://golang.org/doc/codewalk/functions/
