### What objects to consider ?
Most of us (developers) have been introduced to the concept of proramming. Most of us again would have likely heard about 
Object Oriented Programming (OOPs). The essence of OOPs is about constructing objects by analyzing through the 
requirements in terms of real world objects.

### Is it all about _thinking in nouns_ ?
We have often heard of nouns becoming programming objects & the behaviours of these objects turn into methods or functions. 
If we conform to these OOPs essentials, does it ease our programming logic? Does this mean the logic will be understood 
by all its readers? Can we be assured of a good maintainance of this software?

### Loops & Conditions Rule The Roost
Take any programming logic from any real world software and more often than not, you will find yourselves navigating through 
its loops and if else branches. Does OOPs have anything to talk about loops & branches? How about trying to model these loops
and conditions into objects. Yes, I am referring to collections here. However, I am not talking about arrays or maps. I am 
talking about making a collection of these objects as first class citizens similar to the object itself. Let us look at a 
piece of code to understand further.

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
The difference in above thought process is what this article is all about. The very approach towards this thinking makes a 
lot of difference. It makes your software simple yet effective. In my opinion this is the foundation to build software 
that can be better maintained, is easy to unit test, and is easy to extend without injecting bugs.

How can we be so sure of above's advantages?
We talked about considering the collections of objects as first class objects themselves. Next lets try to manage the 
conditions. In general, when we hit a conditional requirement, we developers start to write if else clauses. This is one 
of the shortest keywords in any programming language that creates wonder. I bet a condition based keyword will have more 
occurences than that of loops.
