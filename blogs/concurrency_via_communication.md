These are my learnings on the topic "Share Memory By Communicating".

Most programming languages make use of threads & lock to operate on memory _referred to as shared memory_. Go tries 
to invert this problem of sharing by passing this memory object. In other words, threads need not vie for ownership 
of a memory rather they are provided with this shared memory at appropriate time; removing the need for housekeeping
stuff i.e. use of lock & un-lock. 

Well, there are other approaches as well, e.g. use of thread safe data structures. However, we shall concentrate on Go's
technique to handle concurrency in this article.

Imagine a pipe connected across all the threads. A memory object can pass via this pipe and reach to all the threads 
one at a time. This sounds like threads need not do the typical housekeeping stuff. In other words, threads need not 
worry about putting up locks & then remembering to un-lock when job is done. This worry is otherwise termed as 
"**_explicitly dealing with locks_**". 

In Go terms, this pipe is known as channel & threads are goroutines. Now that we understand how only one go routine has access to shared data at any given time, let us get into other details.

For example, there is a structure that has _lock_ i.e.

```go
type Resources struct {
  data []*Resource
  lock *sync.Mutex
}
```

Above is typical to other langauges and is an indication to future code that is not clean. This is apparent when
above struct gets used in a function that is invoked by multiple threads. One would notice that this function has 
a lot to do in terms of housekeeping i.e. use of lock & unlock.

```go
func Poller(res *Resources){
   // never ending loop
   for {
     res.lock.Lock()
     ...
     res.lock.Unlock()
     ...
     res.lock.Lock()
     ...
     res.lock.Unlock()
   }
}
```

Now imagine if above can be re-factored to follow _passing memory by communicating_ philosophy. We were told about pipes i.e. channels which can be used to pass these shared memories. In other words, go channel will pass `&Resources{...}`. Our Poller
function will look something like below:

```go
func Poller(in, out chan *Resource) {
  for i := range in {
    ...
    out <- i
  }
}
```

While above is a very simple example to understand use of channel(s) instead of mutex, there are a varitety of other complexities that can arise. These complexities can be attributed to the way we design a channel code. Like every
other complexity this too can be handled gracefully if one follows good coding practices. 

Let me list down some coding practices, I have encountered with w.r.t use of channels:
- One channel may not be sufficient to replace a resource's mutex 
- Multiple channels may be needed to orchestrate the entire concurrency
- Make use a function that:
  - Defines a channel
  - Implements the channel handler code which is usually a switch case
  - Embeds this channel handler logic in an immediately invocable anonymous function i.e. iife
  - Run this IIFE i.e. core channel handler as a goroutine
  - Function can finally return this channel, which will be used outside to feed into this channel
- Retrieving resource from a channel is blocking & hence this entire design works
- A resource is typically retrieved from a channel with a loop logic i.e. `for := range`
- A resource can also be retrieved from a channel within a select case which is in turn within a never ending for loop

References:
- https://blog.golang.org/share-memory-by-communicating
- https://golang.org/doc/codewalk/sharemem/
