## Being Functional With Interfaces
I would be talking about go language in particular here. We all have heard of or might even have experienced the benefits of
interface oriented programming. To those who are not aware of the benefits, an interface can be used to switch implementations at runtime thereby providing maximum flexibility while implementing programs.

### Single Method Interface
While interfaces are used most often for above said benefits, there is one good practice with respect to interface that should 
be followed religiously. This good habit that I am talking about is exposing only a single method from an interface. If a 
developer is careful enough to follow this practice for the entire module or the entire project, then it can be assumed that 
code implementation abides by most of the best coding practices. 

However, ensuring use of single method only interfaces is not simple. In some languages it leads to a lot of code, files and 
the entire business logic is lost with these explosion of interfaces. Go helps us a lot in this regard. Any structure in go 
can adhere to an interface contract just by defining a method that has same signature to that of the targeted interface. On 
similar lines, in go language one can define a typed function that models an interface & its implementation. This is easy
if followed after our first rule i.e. _have interfaces that expose only a single method_.

Let us dive into a code snippet to provide more clarity:

```go
type SingleMethod interface {
  Wow(x, y int) int
}

type SingleMethodFunc func(x, y int) int

func (s SingleMethodFunc) Wow(x,y int) int {
  // no concrete implementation here, implementation can be built late
  // can write below syntax as s is a function
  return s(s, y)
}
```

With above snippet we are now clear that a developer needs to code a few lines to achieve above mentioned good practices.
However, some of you readers might be wondering _if this functional way versus typical struct based implementation has any
benefits_? After all, one can just implement a struct that implements the interface and be idiomatic.

The answer to above confusion, lies in being functional. Now, this is not about _being funcctional_ versus _being idiomatic_
debate. If one is able to make use of functions, then there is no need to code a struct that should implement every other 
interface which is obviosuly going to grow based on the core belief of this article i.e. build single method interfaces.

### When functional, then write less code
Let me clarify once more on writing functions that implement single method interfaces. This does not mean one should stop
writing structures that implement interface. This just implies, one has to write little less code. How is this possible?

Let us get inside a code snippet once again:

```go
func main() {
  // anonymous function helps us in writing less code 
  // i.e. we do not need to write struct & properties
  // i.e. less code to bother / maintain
  // i.e. a just in time implementation
  // just cast the anonymous function to functional interface type
  lesscode := SingleMethodFunc(
    func(x, y int) int{
      // can do anything here -- probably return a hardcoded int
      // this logic is delegated to client -- runtime polymorphism
      // will be very helpful writing UnitTest code -- imagine !!
      return x+y
    }
  )
  
  fmt.Println(useinterface(lesscode))
}

func useinterface(i SingleMethod) int{
  return i.Wow(1, 2)
}
```

### Interfaces are great
Interfaces are great and when used along with functions they can be amazing. I have tried to write down my thoughts from 
this original [post](https://www.bartfokker.nl/posts/decorators/) to ingrain this as a habit of mine.

