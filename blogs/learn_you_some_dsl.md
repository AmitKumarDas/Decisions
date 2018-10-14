## Learn You Some DSL !!!
### Chapters
- Where It All Began
- Shameless Cloner
- Illusion of Control
- OODA Loop

## Rough
The first time I chanced to understand a go template file was when I was reviewing a helm chart pull request (PR). This PR was
sent by an operator. On the other side of the PR (read code review) was done by me (i.e. a programmer). I was pleasantly
surprised by the thought process put in by this operator to send this PR. I wondered, how a person with an operator role, was
able to understand the semantics of go language & go templating & then able to build the logic with little knowledge of my library. 

Above speaks a lot on the usefulness of the subset of language constructs exposed by go templating. It is simple enough to be
used by a broader community. A community consisting of programmers, testers, analysts, first timers and so on. 

Below are some of the learnings & ideas I got while coding against go templating:

### 1/ Passing values to a template?

### 2/ Can these values be mutated while template is being executed?

### 3/ Can we build a Domain Specific Language (DSL) based out of go template?
This is where my heart lies. Can we reuse the pipes understood in go templating and design a usable DSL. A DSL that is easy to
understand by non-programmers and also the one which is easy to develop by the programmers. In my opinion, this is key to build
communities where non-programmers can design and construct layers of logic without getting too deep into programming syntaxes.
This results into just having the core functionalities (read functions, schemes, strategies) in the project while letting users
of this project tune these features with a handy tool (read DSL).

### 4/ Go Templating is not YAML templating!!!
Useful to keep reminding ourselves especially in these heady days of Kubernetes. Remember, one can draft email body using go
template as well.

### 5/ It is great to use pipes; even better if you create something of your own
The thing that stood out were the pipes & the ease of adding a custom function. Pipes were suddenly no longer the pergative of
unix shells. It was readily available in a high level language. It helped in assembling multiple functions and we just got our
own functional language which is type safe (courtesy golang).
