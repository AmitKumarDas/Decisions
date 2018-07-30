### My learnings on Sharing Memory By Communicating

Most programming languages make use of threads & lock to operate on memory _referred to as shared memory_. Go tries 
to invert this problem of sharing by passing this memory object. In other words, threads need not vie for ownership 
of a memory rather they are provided with this shared memory at appropriate time; removing the need for extra checks
i.e. remembering to use lock & un-lock. 

Well, there are other approaches as well, e.g. use of thread safe data structures. However, we shall concentrate on Go's
technique to handle concurrency in this article.

Imagine a pipe connected across all the threads. A memory object can pass via this pipe and reach to all the threads 
one at a time. This sounds like threads need not be coded with any orthogonal concerns. In other words, threads need 
not worry about putting up locks & then remembering to un-lock when job is done. This worry is otherwise termed as 
"**_explicitly dealing with locks_**". 

In Go terminology, this pipe is called as channel & threads as goroutines. Now that we are told about how only one go
routine has access to shared memory at any given time, let us get into usage details.

For example, this is the older approach i.e. a structure with _lock_ as one of its propertues:

```go
type Resources struct {
  data []*Resource
  lock *sync.Mutex
}
```

Above is typical coding style that indicative of further code/functions to invoke lock & unlock. This is apparent when
above struct gets used in a function that is invoked by multiple threads. One would notice that this function has 
a lot to do in terms of remembering when to use lock & when to unlock the same. Unlock should further take into account failure cases as well.

```go
func Poller(res *Resources){
   // a never ending loop
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

While above is a very simple example to understand use of channel(s) instead of mutex, there is one important fact that is
still under wraps. And that is, each receive or send on channels are blocking. In fact these very blocking nature of channels, ensures proper working of _passing memory by communicating_. If one needs to use non-blocking sends, receives from one or more channels, then _select with a default case_ can achieve it. 

In addition, there may be understanding (_code readability_) issues w.r.t difference in implementation of channel versus
using the mutex approach. These complexities can be attributed to the way we implement passing the memory via channels.
However, this readability issue(s) can be addressed if one follows effective coding practices. 

Let me list down some coding practices, I have encountered with w.r.t use of channels:
- One channel may not be sufficient to replace a resource's mutex 
- Multiple channels may be needed to orchestrate the entire concurrency
- Make use a function that:
  - Defines a channel
  - Implements the channel handler code which is usually a switch case
  - Embeds this channel handler logic in an immediately invocable anonymous function i.e. iife
  - Run this IIFE i.e. core channel handler as a goroutine
  - Function can finally return this channel, which will be used outside to feed into this channel
- A resource is typically retrieved from a channel with a loop logic i.e. `for := range`
- A resource can also be retrieved from a channel within a select case which is in turn within a never ending for loop

On the whole, if one gets a good grasp of functional programming, understanding use of function closures, then implementing
concurrency via channels will get simpler.

Below is a code snippet that I have copied from one of the reference links. This snippet is the go implementation to 
what I just pointed out earlier.

```go
// StateMonitor maintains a map that stores the state of the URLs being
// polled, and prints the current state every updateInterval nanoseconds.
// It returns a chan State to which resource state should be sent.
func StateMonitor(updateInterval time.Duration) chan<- State {
	updates := make(chan State)
	urlStatus := make(map[string]string)
	ticker := time.NewTicker(updateInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				logState(urlStatus)
			case s := <-updates:
				urlStatus[s.url] = s.status
			}
		}
	}()
	return updates
}
```

References:
- https://blog.golang.org/share-memory-by-communicating
- https://golang.org/doc/codewalk/sharemem/
