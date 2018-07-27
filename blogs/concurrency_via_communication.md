These are my thoughts on the topic "Share Memory By Communicating".

Typical languages make use of threads & lock to operate on memory which is now referred to as shared memory. Go tries to
invert this problem of sharing. In other words, threads need not vie for ownership of a memory rather they are provided 
with the shared memory removing the need for housekeeping stuff i.e. lock & un-lock. 

Well, there are other approaches as well, e.g. use of thread safe data structures. However, we shall concentrate on Go's
technique to handle concurrency.

Imagine a pipe connected across all the threads. A memory object can pass via this pipe and reach to all the threads 
one at a time. This sounds like threads need not manage the house

