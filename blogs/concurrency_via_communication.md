WIP - These are my thoughts on the topic "Share Memory By Communicating".

Typical languages make use of threads & lock to operate on memory which is now referred to as shared memory. Go tries to
invert this problem of sharing. In other words, threads need not vie for ownership of a memory rather they are provided 
with the shared memory at appropriate time; removing the need for housekeeping stuff i.e. lock & un-lock. 

Well, there are other approaches as well, e.g. use of thread safe data structures. However, we shall concentrate on Go's
technique to handle concurrency in this article.

Imagine a pipe connected across all the threads. A memory object can pass via this pipe and reach to all the threads 
one at a time. This sounds like threads need not do the typical housekeeping stuff. In other words, threads need not worry
about putting up locks & then remembering to un-lock when job is done. This is otherwise treated as "**_explicitly dealings with locks_**". In Go terms, the pipe is known as channel & threads are goroutines. 

Now that we understand how only one go routine has access to shared data at any given time, let us get into other details.

References:
- https://blog.golang.org/share-memory-by-communicating
- https://golang.org/doc/codewalk/sharemem/
