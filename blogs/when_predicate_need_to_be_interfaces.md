### Motivation
Often we developers find ourselves trying to extract more mileage from our existing predicates. For those who want to 
understand a predicate; I am referring to the conditional logic. This is way better than typical if-else conditions.
So what should we do when complex logic _(read features)_ needs to be built & executed based on predicates _(there are
supposedly huge number of predicates)_. One way can be to make use of **map[predicateName]complexFuncType** or we can
have an elaborate design to handle all these cases.

### Code
```go
// pkg/zvol/v1alpha1/zvol.go
type zvol struct {
  // ...
}
```

```go
// pkg/zvol/v1alpha1/predicate.go

// Usual Approach
//
// typed predicate function
// and corresponding implementations
type predicateFunc func(*zvol) bool

func (z *zvol) IsInitialized() bool {
  return p.zvol.state != "Init" && p.zvol.state != ""
}

func (z *zvol) IsNotInitialized() bool {
  return !p.zvol.IsInitialized()
}

func IsNotInitialized() predicateFunc {
  return func(z *zvol) bool {
    return z.isNotInitialized()
  }
}

// Additional Implementations
//
// Need to make use of interface &
// corresponding implementation(s)
// if our requirement demands it
type predicate interface {
  check()       bool     // actual conditional logic
  name()        string   // name of this condition
  promCounter()          // some prometheus counter logic
  pass()        string   // success message when this condition passes
  fail()        string   // failure message when this condition fails
}

type predicateList []predicate

func (l predicateList) promCounter() {
  for _, p := range l {
    p.promCounter()
  }
}

type predicateListBuilder struct {
  zvol *zvol
  list predicateList
}

func PredicateListBuilderForObject(zvol *zvol) *predicateListBuilder {
  return &predicateListBuilder{zvol: zvol}
}

func (b *predicateListBuilder) WithIsNotInitialized() *predicateListBuilder {
  b.list = append(p.list, IsNotInitializedImpl(b.zvol))
}

func (b *predicateListBuilder) PromCounter() {
  b.list.promCounter()
}

type isNotInitialized struct {
  zvol *zvol
  cond  predicateFunc
}

func IsNotInitializedImpl(zvol *zvol) predicate {
  return &isNotInitialized{zvol: zvol}
}

func (p *isNotInitialized) check() bool {
  return p.cond()(p.zvol)
}

func (p *isNotInitialized) name() string {
  return "IsNotInitialized"
}

func (p *isNotInitialized) promCounter() {
  if !p.check() {
    return
  }
  // increment prometheus counts corresponding
  // to isNotInitialized
}

func (p *isNotInitialized) pass() string {
  return "isNotInitialized passed"
}

func (p *isNotInitialized) fail() string {
  return "isNotInitialized failed"
}
```
