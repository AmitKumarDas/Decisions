### Introduction
Most of us (_developers_) have been introduced to the concept of programming via books or our own experience. It is more than
likely some of us would have heard about Object Oriented Programming (OOPs). The essence of OOPs is all about constructing
objects by analyzing software requirements in terms of real world objects.

### Need to move beyond _thinking in nouns to construct objects_ ?
We have often heard of nouns becoming programmable objects & the behaviours of these objects turn into methods or functions. 
Theory aside, if we conform to these OOPs essentials, does it ease our programming logic? Does this mean this logic will be
understood by all its readers (read reviewers, fellow maintainers)? Can we be assured of a good maintainance of this 
software?

### Loops & Conditions Rule The Roost
Take any programming logic from any real world software and more often than not, you will find yourselves navigating through 
its loops and if else branches. _Does OOPs have anything to talk about loops & branches?_ _How about modeling these loops
and conditions into defined types?_ Yes, I am referring to collections here. However, I am not talking about arrays or maps.
_**I am indicating to make a collection of these objects as first class citizens similar to the object itself.**_ Let us look
at a piece of code to understand further.

```go
type storage struct {}

type storageList struct {
  items []storage
}

type storageMap struct {
  items map[string]storage
}
```

### Does above make any difference?
The difference in above thought process is what this article is all about. The very approach of modelling collections as 
custom types makes a lot of difference. It makes your software simple yet effective. In my opinion this is the foundation to
build software that can be better maintained, is easy to unit test, and is easy to extend without injecting bugs.

_How can we be so sure of above's advantages?_
We talked about considering the collections of objects as typed objects. To understand the full benefits, let us park collections for a moment and try to think about conditions. 

In general, when we hit a conditional requirement, we developers start to write _if else_ clauses. This is one of the
shortest keywords in any programming language that creates wonder. I bet a condition based keyword will have more occurences
than that of loops. Getting back to our mode of thinking in types, can these conditions be treated as first class citizens?
In other words should we define custom type to represent a condition?

In my opinion, solution lies in composing these collection types along with conditional types to arrive at writing logic that
is very effective and yet is simple to maintain.

### Creating your own conditionals i.e. custom type to represent a condition
Why? This ensures decoupling conditions from loop based logic. _(Remember that we have already moved the loops into custom defined types i.e. collection types.)_ Let us resort to some code samples to understand better.

```go
// Predicate abstracts evaluating storage based condition
type Predicate(s storage) bool

// isOnline evaluates if storage is online or not
//
// NOTE:
// Function adheres to Predicate functional type
//
// NOTE:
// One can find domain specific logic only
func isOnline(s storage) bool {...}

// Filter returns a list of storages that satisfy the provided
// condition
//
// !!! IMPORTANT !!!
//
// This is the gist of this article. Loop meets the condition here.
// Domain logic is not seen here. What we have done here is - decouple the
// business logic from programming constructs i.e. loops & conditions.
func (l storageList) Filter(p Predicate) (u storageList) {
  for _, s := range l {
    if p(s) {
      u.items = append(u.items, s)
    }
  }
  return
}

// Onlines returns a list of storages that are online
//
// NOTE:
// Programming is fun now.
func (l storageList) Onlines() (u storageList) {
  return l.Filter(isOnline)
}
```

Above tries to model a simple storage related requirement by considering a collection based type & a conditional type.
With similar core pieces in place, we can venture to stitch together (read permutations & combinations) of complex
requirements that will be easy to code and maintain and even unit test.

### Conclusion
It is all about **thinking in types**. Some may view above as functional approach as well. However, the core idea of
thinking in terms of specific structures, functions or objects in addition to language provided keywords like `for loop`, 
`switch case` conditions, etc can help us in getting some amazing stuff done with a maintainable codebase. 

### We are all familiar with `sort` Interface
- How about implementing sort.Interface to sort the storage collections !!!
```go
func (l storageList) Len() int {
	return len(l)
}

func (l storageList) Less(i, j int) bool {
	return l[i].LessThan(l[j]) // implement LessThan method
}

func (l storageList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
```
- NOTE: Most of the projects seem to use a custom list to implement `sort` interface. 
- NOTE: It could be great if these projects look beyond supporting sort & start using it to build domain logic as well.

### Ever thought of building custom operators as well ?
How about requirements driving us to implement various operators on top of collections & conditions? e.g. `or`, `eq`, `ne`,
`gt`, `lt`, and so on.

### Fun Side
- I tend to call above as _first-class domain types_. A jargon that will obviously be mis-interpreted. Hence this article.
