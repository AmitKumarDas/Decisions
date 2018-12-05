### Code That Lasts
The title seems to go against the very nature of change. Nothing lasts for forever. However, mankind will always be in a 
constant disagreement with nature and will try to build fortress around their idea(s).

This holds true for developers writing code and expect it to last for many many releases to come. Not all of them really 
think about this. However, this does not matter since all the developers must go with the flow of change by churning the
code based on the bugs discovered or based on the newly articulated requirements that seem to go against the very 
assumption that was made in the first place.

Is there any hope then to have some pieces of logic that can last?

### Code that is malleable
Remember the mighty trees that cannot bend usually break or get uprooted during natural calamities. Surprisingly, the young 
plants which seem to be more fragile come out unscathed. The question that we need to ask ourselves is _"is your code 
malleable to the forces of change"_ ?

### Patterns as a language to communicate
We know each one in the team writes code in a different way. It is good in one way that language does not become a barrier to
write code. A developer is free to choose a particular pattern that suits the requirements. OpenEBS developers have been 
doing this for some time now. Alas, this has a lot of associated costs. Lot of review cycles, longer time to implement 
business logic, even longer effort required to implement test logic, un-testable blocks of code, lack of code re-usability. 
This made us realize about specific patterns that stand out and can help in writing **better**, **faster** yet **testable** 
code. Well, idiomatic maya is all about _internalizing_ these patterns into making openebs _antifragile_.

### Devil is in the details
The stuff that I am going to write below will change based on my experiences and observations. These details should be 
a living document. These are the details that I believe will _make code more malleable to changes without cracking or being 
prone to bugs_.

#### High Level
- start date - 16 Nov 2018
- updated on:
  - 26 Nov 2018
  - 27 Nov 2018

- API:
  - Should have well defined payloads also known as API types
  - These API types should be in a well defined namespace/packaging
  - e.g. `pkg/apis/org.group/entity-name/version/`
- Entity
  - Should have resources/entities that closely match (or exactly match) these API types
  - These entities should be in their own defined namespace (different than API namespaces)
  - e.g. `pkg/entity-name/version/`
  - EntityList i.e. a collection of these entities should be a well defined type
  - refer - https://medium.com/@amit.das/thinking-in-types-3c234eb17680
- EntityAddons
  - Above entities may have some or all of these:
  - Builder, ListBuilder
  - Errors,
  - RestClient,
  - Template e.g. GoTemplate
  - Predicate, Map, Filter, Middleware, etc functions i.e. streaming APIs
  - All of these addons can reside in the same namespace that is meant for entity
  - i.e. `pkg/entity-name/version`
- Unit Testable Entity
  - The entity along with their helpers must be Unit Testable
- Caller Code / Service Code
  - These is the place which ties up the business logic
  - These are the feature implementors
  - These should reside in a namespace of their own
  - These might make use of DependencyInjection library e.g. `wire`
  - These might need to make use of profiling library
    - e.g. `https://github.com/pkg/profile`
  - e.g. `pkg/installer/v1alpha1` 
  - or `pkg/controller/openebs/v1alpha1`
  - or `cmd/maya-apiserver/server/resource_endpoint.go`
- Unit Testable Caller Code
  - The caller code should also be unit testable
  - However, unit testable caller code is optional
  - Some may want to test this via integration tests

#### Low Level
- start date - 26 Nov 2018
- updated on:
  - 27 Nov 2018, 
  - 30 Nov 2018,
  - 05 Dec 2018

```go
// pkg/entity/v1alpha1/entity.go

package v1alpha1

// Interface defines the public contracts 
// of this namespace
type Interface interface {
  // contracts i.e. methods that need to be exposed
}

// structure that composes 
// - fields
// - 3rd party interfaces -- can be mocked
// - 3rd party functions -- can be mocked
type entity struct {
  p1 string
  p2 string
}

// builder is used to build the entity
// 
// NOTE:
// A builder is typically used when the wrapped entity 
// has lots of fields. Some of the times, caller is not 
// interested to fill all these fields. While at other 
// scenarios, caller is interested to run validations 
// against one or more of these fields. Interestingly,
// there will be cases where both of the above is 
// applicable.
//
// NOTE:
// This lets the caller code make use of multiple dot 
// notation to invoke multiple methods. This makes the 
// caller code readable IMO.
type builder struct {
  e       *entity
  default bool
  checks  []Predicate
  errors  []error
}

// OptionFunc helps in building the entity instance
//
// NOTE:
// This is an alternative approach to build an entity
//
// NOTE:
// However, this exhibit more malleable characteristics 
// than that of builder. This function signature is 
// fundamental to building powerful codebase. Developers
// familiar with **functional programming**, **streaming 
// APIs**, etc can relate more with this approach. The 
// power is realized when this typed function is combined 
// with following:
//
// 1/ predicates 
// 2/ transformers also called mappers 
// 3/ predicates as well as transformers
//
// Above is essentially functional programming and its 
// worth noting to learn more on this from different 
// languages.
//
// TODO:
// !!Help Needed!! ??Contributors??
//
// Write the names & stories they offer to this pattern 
// in different languages
type OptionFunc func(*entity)

// New returns a new instance of entity based on the 
// provided options
//
// NOTE:
// New is built on top of passing typed functions that 
// can be thought of as fragments to build the entity 
// & then return the same. The arguments are variadic, 
// which means passing zero arguments will compile fine.
func New(opts ...OptionFunc) *entity {
  var e entity
  for _, o := range opts {
    o(&e)
  }
  return &e
}

// Default is a variant of New. It returns an instance 
// of entity by filling up the default fields that are 
// not set previously.
func Default(opts ...OptionFunc) *entity {
  return default(New(opts))
}

// default sets default fields against the provided 
// entity if not set previously
func default(e *entity) *entity {
  // TODO
  // check if nil check is required ?
  if e == nil {
    e = &entity{}
  }
  if len(e.p1) == 0 {
    e.p1 = "noop"
  }
  return e
}

// WithP1 sets p1 field
//
// NOTE:
// This returns a typed function that enables WithP1
// to participate in various combinations. For example 
// it can be passed as an argument to New, Default, Map,
// MapAll, etc functions and be a part of building a 
// higher order feature.
func WithP1(p1 string) OptionFunc {
  return func(e *entity) {
    e.p1 = p1
  }
}

// Builder returns a new instance of builder
func Builder() *builder {
  // TODO
  // check if New returns a non-nil instance ?
  return &builder{e: New()}
}

// P2 sets p2 field against entity instance
func (b *builder) P2(p2 string) *builder {
  b.e.p2 = p2
  return b
}

// UseDefaults will turn on default fields' setters
func (b *builder) UseDefaults() *builder {
  b.default = true
  return b
}

// AddCheck adds a check to be done against entity
//
// NOTE:
// AddCheck makes use of Predicate rather than exposing
// new methods for every new condition
func (b *builder) AddCheck(c Predicate) *builder {
  b.checks = append(b.checks, c)
  return b
}

// validate will run the checks against the entity
func (b *builder) validate() error {
  for _, c := range b.checks {
    if !c(b.e) {
      b.errors = append(b.errors, ValidationFailedError(c))
    }
  }
  if len(b.errors) == 0 {
    return nil
  }
  return GenericValidationFailedError
}

// Build will return the final instance of entity
func (b *builder) Build() (*entity, error) {
  if b.default {
    b.e = default(b.e)
  }
  err := b.validate()
  if err != nil {
    return nil, err
  }
  return b.e, nil
}

// EntityList represents a list of entities
//
// NOTE:
// Defining a type to this simple collection
// is essential to utilize the power of functional
// programming. This type is typically used along
// with Predicate(s), Filter(s), Transformer(s), 
// Comparator(s), etc to build high level logic.
type EntityList []*entity

// Predicate abstracts conditional logic w.r.t an entity
type Predicate func(* entity) bool

// PredicateList is a typed representation of list of
// predicates
//
// NOTE:
// Defining a type against this collection enables
// building custom operators on top of predicates.
// For example, **OR**, **AND**, etc operators can be 
// built on top of predicates.
type PredicateList []Predicate

// All returns true if all the predicate conditions 
// pass successfully against the provided entity
//
// NOTE:
// All represents AND-ing of predicates found in
// PredicateList
func (l PredicateList) All(e *entity) bool {
  for _, p := range l {
    if !p(e) {
      return false
    }
  }
  return true
}

// Any returns true if at-least one of the predicate 
// condition pass successfully against the provided 
// entity
//
// NOTE:
// Any represents OR-ing of predicates found in
// PredicateList
func (l PredicateList) Any(e *entity) bool {
  if len(l) == 0 {
    return true
  }
  for _, p := range l {
    if p(e) {
      return true
    }
  }
  return false
}

// Map runs the provided option against each entity 
// if all predicates succeed
func (l EntityList) Map(opt OptionFunc, pl ...Predicate) {
  for _, e := range l {
    if PredicateList(pl).All(e) {
      opt(e)
    }
  }
}

// MapIfAny runs the provided option against each 
// entity if at-least one predicate succeeds
func (l EntityList) MapIfAny(opt OptionFunc, pl ...Predicate) {
  for _, e := range l {
    if PredicateList(pl).Any(e) {
      opt(e)
    }
  }
}

// MapAll runs all the provided options against each
// entity if all predicates succeed
func (l EntityList) MapAll(opts []OptionFunc, pl ...Predicate) EntityList {}

// MapAllIfAny runs all the provided options against 
// each entity if at-least one predicate succeeds
func (l EntityList) MapAllIfAny(opts []OptionFunc, pl ...Predicate) EntityList {}

// Filter returns a list of entities that succeed 
// against all the provided predicates
func (l EntityList) Filter(pl ...Predicate) EntityList {}

// FilterIfAny returns a list of entities that succeed
// against at-least one predicate
func (l EntityList) FilterIfAny(pl ...Predicate) EntityList {}

```
#### Low Level - References
- https://github.com/openebs/maya/tree/master/pkg/msg/v1alpha1
- https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/util/errors/errors.go
