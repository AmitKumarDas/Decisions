## Delegate concerns to middleware
The topic may sound quite absurd. Well, it sounds absurd as I have mixed English with Go language's jargon. However, I felt
that this topic is fairly simple enough to be remembered & hence is easy enough to be referred to later.

_NOTE - This topic is all about programming in go language._

### Into the topic
Lets try to understand each word of the title. Middleware is a pattern/jargon referred to in go programming. This is also 
known as decorator pattern in other languages. Well if its decorator that we are talking about, then concerns are those orthogonal aspects of programming. In other words, concerns such as logging, metrics, caching, authentication & so on are orthogonal to business logic. Hopefully, by now we are clear about this topic, which is about _**how to delegate non-business logic related code to decorator functions.**_

### How do I?
Now that we have somewhat been introduced to the topic of middleware i.e. decorators in go, the next question that comes up is
_how to implement one?_ Let us go through some middlewares then:

```go
type DoIt interface {
  Do(string) string
}

// DoItFunc is a typed function that implements DoIt interface
type DoItFunc func(string) string

// Do is the DoIt's method implementation
// The concrete implementation is not here
// This concrete logic is coded lazily
func (fn DoItFunc) Do(input string) (output string) {
  return fn(input)
}

// DoItMiddleware is the middleware we are talking about here
// Note that we are making use of single method interface which
// simplifies 
type DoItMiddleware func(DoIt) DoIt

// Here goes the delegation of concerns to middleware(s)
// Auditor
func DoItAuditor(a *myAuditPackage) DoItMiddleware {
  // concrete logic for DoItMiddleware function
  return func(d DoIt) DoIt {
    // annonymous function has the concrete logic 
    // that matches the DoIt interface & in addition decorates 
    // the original DoIt instance
    fn := func(input string) string {
      // do the audit
      // doit execution
      return d(input)
    }
    // cast to a DoIt implementor
    return DoItFunc(fn)
  }
}
```

### References
- https://www.bartfokker.nl/posts/decorators/
