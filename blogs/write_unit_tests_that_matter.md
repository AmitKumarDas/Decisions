Ever wondered if Unit Tests should be written in a particular way? Should they follow a typical pattern? Go community has been
advocating the use of `Table Driven Tests` to derive more benefits. Though `Table Driven Tests` are nothing new, it has become
more common these days. I see a clear value in writing tests that are readable and can be extensible. In this article, I shall
present some examples that will bring home the point of _**writing tests that matter**_.

So how does a test matter? To a layman it is feeding a system with scenarios and compare the resulting output with expected 
(or calculated) results. If we developers, stick to this basic rule that helps a layman to understand the target under test, 
then I believe we have done a great job. In other words, let the reader/reviewer (or anybody trying to understand the system
by browsing the test code) not bother about intricate details of the system, or peculiar language semantics while reading the
test code.
