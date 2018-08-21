### Code Review - Episode 2
This article is part of a series of blogs related to code review. I keep writing down the source code of my own project or
one from some open source community and start providing my own review comments as code level commments. This helps me in 
avoiding any first conclusion bias(-es) I might form when I look at the code for the first time. This also helps me to think
by putting on a user's hat and review the code from its readability point of view.

I shall divide this blog into two parts:
- Part 1 talks about my recents thoughts about writing `practical code` vs. `smart code` vs. `generic code`
- Part 2 talks about one of my latest source code and its code review which validates the discussions of Part 1

## Part 1
### Smart code vs. Generic code vs. Practical code
We need not elaborate deeply on what is generic code. I guess most of us understand this. This typically means having a lot
of interfaces, higher order functions and so on to implement a given requirement. This will enable us to achieve 
extensibility. The only side effect will be more lines of code. Now this is all fine when the requirement demands such a 
need.

So what is smart code then? Smart code will get 10 on 10 from the review hubs. Smart code is perhaps the best way to have 
the logic implemented for some piece of requirement that is un-thinkable by most of the fellow developers. 

However, I have a different take with respect to code that is smart. I feel a code is given a smart tag, when a fellow 
maintainer or reviewer looks at the code and does not feel at home. This other fella somehow feels it must be good code, 
since it looks terse and makes use of some of the language features that is not so common. Only the one who lives in the 
language day in & day out can write such a code. There are no linting issues and it has very good code coverage. This pushes 
this code reviewer to approve the code with a Looks Good To Me `LGTM`. Needless to say, this reviewer feels guilty about 
approving this code which works fine but misses something which is difficult to describe. What is this something that is so
difficult to explain?

### Being practical is user friendly
There are so many factors that a developer needs to keep in mind, while implementing a logic. There are chances this 
developer goes too far keeping best practices in mind. At the end of the day, this logic will get `10/10` from various 
linting and other tools. This sounds great. However, the thing that I have come to realize rather late, is the ability to 
write code that is easy to understand by the reviewers and maintainers is _**the most important**_ that should matter.

I shall delve into code that tries to uncover if the same is more practical and easy to understand versus if it is more 
smart versus if it is more generic.

## Part 2 
### Here are the code reviews
**Check for `Notes:` section in the code level comments**

```go
package v1alpha1

import (
  "fmt"
)

// Artifact has the artifact that will be installed or uninstalled
type Artifact struct {
	// Doc represents the original artifact itself
	Doc string
}

// ArtifactList is the list of artifacts that will be installed or uninstalled
type ArtifactList struct {
	Items []*Artifact
}

// ArtifactListOptions is used to list artifacts based on these options
//
// Notes:
// - This was introduced to provide a generic list option argument. 
// - contd... In other words, there might be more options in future.
// - contd... Using this struct helps avoiding adding paramters to the function signature.
// - contd... This adheres to open-closed principle
//
// - Having an argument in the interface method enables decoupling the creation & invocation parts
// - e.g. creation of an instance does not need this as an argument
// - contd... this argument will be provided lazily during invocation of the interface method
// - Many a times, logic does not have the argument handy during creation of an instance
// - So this way of exposing via interface & interface method seems great
type ArtifactListOptions struct {
  Version string
}

// ArtifactLister abstracts listing of Artifacts
type ArtifactLister interface {
	List(options ArtifactListOptions) (ArtifactList, error)
}

// ArtifactListerFunc is a functional implementation of ArtifactLister
//
// Notes:
// - This helps in avoiding creating a dedicated structure. Go vs. Java. Go wins here.
// - Less code is good
// - Be careful of terse code that may seem smart but is a painful w.r.t maintainance
// - This particular functional implementation might be a learning curve to fellow developers.
// - contd... be mindful of above
type ArtifactListerFunc func(options ArtifactListOptions) (ArtifactList, error)

// List is an implementation of ArtifactLister
func (fn ArtifactListerFunc) List(options ArtifactListOptions) (ArtifactList, error) {
	return fn(options)
}

// VersionedArtifactLister returns a new instance of ArtifactLister that is
// capable of returning artifacts based on the provided version
//
// Notes:
// - Just a function suffices. No need of a dedicated structure as noted earlier
// - One can similarly create functions that cater to specific requirements without the need for structs
// - One is able to pass the options late during invocation time. Good
func VersionedArtifactLister() ArtifactLister {
	return ArtifactListerFunc(func(options ArtifactListOptions) (ArtifactList, error) {
		switch options.Version {
		case "0.7.0":
			return RegisteredArtifactsFor070(), nil
		default:
			return ArtifactList{}, fmt.Errorf("invalid version '%s': failed to list artifacts by version", version)
		}
	})
}

// Take 2
//
// One the whole above code seems good. 
// Above code is not too verbose and is not too terse as well. 
//
// However, is there any other approach that can be more pragmatic?
//
// Let us check below code snippets:

// VersionArtifactLister abstracts fetching a list of artifacts based on 
// provided version
//
// Notes:
// - The interface is no more
// - The functional implementation of interface is gone as well
// - A higher order function is just what is needed
// - The argument of this functionn type makes use of version which seems to be very natural
// - We still achieve decoupling of defining a function from its invocation as happened earlier
// - contd... decoupling is achieved by passing function first
// - contd... & later by invoking this function when its argument(s) are available
//
// - The main debate of practical code comes up here
// - If the requirements will never make use of any other filtering abilities other than `version`
// - contd... then this code seems to be more practical than the one above
type VersionArtifactLister func(version string) (ArtifactList, error)

// ListArtifactsByVersion returns artifacts based on the provided version
//
// Notes:
// - This is an implicit implementation of VersionArtifactLister function type
func ListArtifactsByVersion(version string) (ArtifactList, error) {
  switch version {
  case "0.7.0":
	  return RegisteredArtifactsFor070(), nil
  default:
	  return ArtifactList{}, fmt.Errorf("invalid version '%s': failed to list artifacts by version", version)
  }
}

// Take 3
//
// In this take, we shall see that being functional solves most of our needs i.e.practical & readable code
// Let us check below code snippets.

type ArtifactListOptions struct {
	Version string
}

// ArtifactLister abstracts fetching a list of artifacts
//
// Notes:
// - Sticking to higher order function
// - Name of this function is changed to make it more apt
// - The argument of this function is different here
// - We still achieve decoupling of defining a function from its invocation as happened earlier
// - contd... decoupling is achieved by passing function first
// - contd... & later by invoking this function when its argument(s) are available
type ArtifactLister func(options ArtifactListOptions) (ArtifactList, error)

// ListArtifactsByVersion returns artifacts based on the provided version
//
// Notes:
// - This is an implicit implementation of ArtifactLister function type
func ListArtifactsByVersion(options ArtifactListOptions) (ArtifactList, error) {
  switch options.Version {
  case "0.7.0":
	  return RegisteredArtifactsFor070(), nil
  default:
	  return ArtifactList{}, fmt.Errorf("invalid version '%s': failed to list artifacts by version", options.Version)
  }
}
```
