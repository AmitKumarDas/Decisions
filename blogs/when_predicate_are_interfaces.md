### Motivation
Often we developers find ourselves trying to extract more mileage from our existing predicates. For those who want to 
understand a predicate; I am referring to the conditional logic. This is way better than typical if-else conditions.
So what should we do when complex logic needs to be built & executed based on predicates. Read further to see the solution.

### Code
```go
// pkg/zvol/v1alpha1/zvol.go
type zvol struct {
  // ...
}
```
```go
// pkg/zvol/v1alpha1/predicate.go
type predicate interface {
  ok()      bool     // actual conditional logic
  name()    string   // name of this condition
  promCounter()      // some prometheus counter logic
  pass() string   // success message when this condition passes
  fail() string   // failure message when this condition fails
}

type isNotInitialized struct {
  zvol *zvol
}

func (p *isNotInitialized) ok() bool {
  return p.zvol.state == "Online"
}

func (p *isNotInitialized) name() string {
  return "IsNotInitialized"
}

func (p *isNotInitialized) promCounter() {
  // increment specific prometheus counts
}

func (p *isNotInitialized) pass() string {
  return "isNotInitialized passed"
}

func (p *isNotInitialized) fail() string {
  return "isNotInitialized failed"
}
```
