```go
// pkg/reconcile/v1alpha1/reconcile.go

type Interface interface {
  Run() error
}
```

```go
// pkg/reconcile/entity/v1alpha1/entity.go
// e.g. pkg/reconcile/spc/v1alpha1/spc.go

type Reconcile struct {
  Name string
  Namespace string
}

type Predicate func(*Reconcile) bool

type PredicateList []Predicate

type (l PredicateList) All(r *Reconcile) bool {
  for _, cond := range l {
    if !cond(r) {return false}
  }

  return true
}

type (l PredicateList) Any() bool {
  if len(l) == 0 {
    return true
  }

  for _, cond := range l {
    if cond(r) {return true}
  }
  return false
}

type Action func(*Reconcile) error

func (r *Reconcile) Run() error {
  // NOTE:
  //   in some reconciliation scenarios
  // order of actions is important
  var actions = []Action {
    Delete(IsObjectNotExists),
    Create(IsChildrenNotExists),
    Update(IsDiff),
  }

  for _, act := range actions {
    err := act(r)
    if err != nil {
      return err
    }
  }
  
  return nil
}

// Delete handles the cases when reconciliation
// needs to execute delete related operations
//
// NOTE:
//  Most common example will be deleting the
// children of the watched object
func Delete(conds ...Predicate) Action {
  return func(r *Reconcile) error {
    checks := PredicateList{conds}
    if !checks.All(r) {
      // checks fail; so nothing to do
      return nil
    }
    
    // logic to delete
    // ...
  }
}

// Create handles the cases when reconciliation
// needs to execute create related operations
//
// NOTE:
//  Most common example will be creating the
// children of the watched object
func Create(conds ...Predicate) Action {
  return func(r *Reconcile) error {
    checks := PredicateList{conds}
    if !checks.All(r) {
      // checks fail; so nothing to do
      return nil
    }
    
    // logic to create
    // ...
  }
}

// Update handles the cases when reconciliation
// needs to execute update related operations
//
// NOTE:
//  Most common example will be updating the
// children of the watched object
func Update(conds ...Predicate) Action {
  return func(r *Reconcile) error {
    checks := PredicateList{conds}
    if !checks.All(r) {
      // checks fail; so nothing to do
      return nil
    }
    
    // logic to update
    // ...
  }
}
```
