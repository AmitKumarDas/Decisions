## WIP - Delegate concerns to middleware
The topic may sound quite absurd. Well, it sounds absurd as I have mixed English with Go language's jargon. However, I felt
that this topic is fairly simple enough to be remembered & hence is easy enough to be referred to later.

_NOTE - This topic is all about programming in go language._

### Into the topic
Lets try to understand each word of the title. Middleware is a pattern/jargon referred in go programming. This is also known as
decorator pattern in other languages. Well if its decorator that we are talking about, then concerns are those orthogonal
aspects of programming. In other words, concerns such as logging, metrics, caching, authentication & so on are orthogonal to
business logic. Hopefully, by now we are clear about this topic, which is about _**how to delegate non-business logic related
code to decorator functions.**_

### How do I?
