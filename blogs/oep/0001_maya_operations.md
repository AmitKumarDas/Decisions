### Motivation
There have been various requirements that would like to have `Maya` i.e. control plane for openebs provide a workflow 
or pipeline to run an ordered set of tasks. Maya has so far strived to provide this in form of CASTemplates & RunTasks.
While this is convenient it has lot of drawbacks. CASTemplates & RunTasks are provided as a declarative YAML. Finally, 
when higher order logic is required it resorts to golang based templating. One needs to understand the fact that 
`templating` cannot replace a programming language and is never turing complete.

Hence, this proof of concept is yet another attempt to provide **MayaOps** _(similar in spirit to GitOps or Github actions 
or so on)_. As one might have already guessed, this attempt will avoid declarative approach. However, without getting too 
much into programmatic versus declarative approach, let us list down the features of `MayaOps`.

#### What MayaOps should provide
- [ ] Simple way to code operations supported by Maya
- [ ] Provide proper error messages for invalid operations
- [ ] Operation steps can be ordered

### High Level Design

#### Core - Drop 0
- Introduces **ops** pattern
- Ops is the shorthard notation for Operations
- The concept is built using **builder pattern**
  - In this case each build method is executed immediately
- _Background_:
  - Ops is built on top of Maya's core builder & predicate functions
  - Maya's core builder pattern is however based on lazy executions
  - In other words, core builder execute in their final build methods
  - To repeat, Ops builder differs from core builder due to its immediate execution style
- To avoid clash among both these builders, this is termed as **operations** or **ops**

```go
// pkg/ops/v1alpha1/interface.go

// Runner abstracts execution of an
// operation
type Runner interface {
  // Run does the actual execution of 
  // operation
  Run() error
}

// Verifier abstracts verfication of
// an operation
type Verifier interface {
  // Verify verifies if operation execution
  // was successful
  Verify() error
}
```

```go
// pkg/ops/v1alpha1/ops.go

// OperationListRunner is a concrete implementation
// of Runner interface. As the name suggests
// it runs a list of operations
type OperationListRunner struct {
  Items []Runner
}

// NewOperationListRunner returns a new instance of
// OperationListRunner
func NewOperationListRunner(o ...Runner) *OperationListRunner {
  m := &OperationListRunner{}
  m.Items = append(m.Items, o...)
  return m
}

// Run executes all the underlying operations
// managed by this instance
func (g *OperationListRunner) Run() error {
  var err error
  for _, op := range g.Items {
    err = op.Run()
    if err != nil {
      return err
    }
  }
  
  return nil
}

// OperationListVerifier is a concrete implementation
// of Verifier interface. As the name suggests
// it verifies a list of already executed operations
type OperationListVerifier struct {
  Items []Verifier
}

// NewOperationListVerifier returns a new instance of
// OperationListVerifier
func NewOperationListVerifier(o ...Verifier) *OperationListVerifier {
  m := &OperationListVerifier{}
  m.Items = append(m.Items, o...)
  return m
}

// Verify verifies if all its operations
// were executed successfully
func (g *OperationListVerifier) Verify() error {
  var err error
  for _, op := range g.Items {
    err = op.Verify()
    if err != nil {
      return err
    }
  }

  return nil
}
```

```go
// pkg/ops/v1alpha1/ops.go

// BaseOps composes of all common fields
// required for any Ops structure
type BaseOps struct {
  ID           string
  Namespace    string
  ShouldSkip   bool
  Errors       []error
  Store        map[string]string
}

// BaseOpsOption is a custom function that
// abstracts building the BaseOps instance
type BaseOpsOption func(*BaseOps)

// New returns a new instance of BaseOps
func New(opts ...BaseOpsOption) *BaseOps {
  b := &BaseOps{}
  for _, o := range opts {
    o(b)
  }
  return b
}
```

### Low Level Parts
- This is implementation of pkg/ops/v1alpha1/ interfaces

```go
// pkg/ops/kubernetes/pod/v1alpha1/pod.go

import (
  ops "github.com/openebs/maya/pkg/ops/v1alpha1"
  pod "github.com/openebs/maya/pkg/kubernetes/pod/v1alpha1"
)

// NOTE:
//  This file (along with its package) deals
// with pod related operations only. 
//
// NOTE:
//  Each operation step is named in a way that
// reflects the actual behaviour.
// 
//  For example, A pod specific condition such as 
// IsRunning() should be defined as two different
// functions named as:
// - ShouldBeRunning() &
// - ShouldNotBeRunning()
//
//  These functions can be set in an order necessary 
// to build one operation.
//
//  One can take the reference from BDD specs
// to decide a name for any operation based step

type Ops struct {
  ops.Base
  Pod          *Pod
  RunSteps     []OpsStep
}

// OpsOption abstracts building
// an Ops instance
type OpsOption func(*Ops)

// OpsStep abstracts implementation
// of a step participating in the
// operation
type OpsStep func(*Ops)

func (o *Ops) withOptions(opts ...OpsOption) *Ops{
  for _, option := range opts {
    option(o)
  }
  return o
}

// New returns a new instance of Ops
func New(opts ...OpsOption) *Ops {
  o := &Ops{ops.Base: ops.New()}
  return o.withOptions(opts...)
}

// From returns a new instance of Ops
// by making use of the provided base
// Ops instance
func From(base *BaseOps, opts ...OpsOption) *Ops {
  b := base
  if b == nil {
    b = ops.New()
    b.Errors = append(b.Errors, errors.New("failed to init pod ops: nil base provided"))
  }

  o := &Ops{ops.Base: b}
  return o.withOptions(opts...)
}

// Steps sets Ops instance with the
// provided steps that will be run later
//
// NOTE:
//  These steps form the core of pod 
// operations
func (p *Ops) Steps(opts ...OpsStep) *Ops {
  p.RunSteps = append(p.RunSteps, opts...)
  return p
}

// Run executes the operations as steps in an ordered
// manner. The order that was used while creating
// this ops instance is used to execute as well.
//
// NOTE:
//  Run is an implementation of ops.Runner 
// interface
func (p *Ops) Run() error {
  if len(p.Errors) > 0 {
    return errors.New("%v", p.Errors)
  }

  for _, step := range p.RunSteps {
    if len(p.Errors) > 0 {
      return errors.New("%v", p.Errors)
    }

    if p.ShouldSkip {
      break
    }

    step(p)
  }

  return nil
}

// Verify provides the status of previously 
// executed operation steps
//
// NOTE:
//  Verify is an implementation of ops.Verifier
// interface
func (p *Ops) Verify() error {
  if len(p.Errors) > 0 {
    return errors.New("%v", p.Errors)
  }

  return nil
}

// isContinue flags if this operation
// should continue executing subsequent
// steps
func (p *Ops) isContinue() bool {
  return !(len(p.Errors) > 0 || p.ShouldSkip)
}

// WithStore provides an in-memory store
// to do the following:
//
//  [1] store results of operation steps &
//  [2] store config values
func WithStore(store map[string]string) OpsStep {
  return func(p *Ops) {
    p.WithStore(store)
  }
}

// WithStore provides an in-memory store
// to do the following:
//
//  [1] store results of operation steps &
//  [2] store config values
func (p *Ops) WithStore(store map[string]string) *Ops {
  if !p.isContinue() {
    return p
  }
  
  if store == nil {
    err := errors.New("nil store was provided")
    p.Errors = append(p.Errors, err)
    return p
  }

  p.Store = store
  return p
}

// ShouldBeRunning sets error if pod is not
// running
func ShouldBeRunning() OpsStep {
  return func(p *Ops) {
    p.ShouldBeRunning()
  }
}

// ShouldBeRunning sets error if pod is not
// running
func (p *Ops) ShouldBeRunning() *Ops {
  if !p.isContinue() {
    return p
  }

  isRunning := pod.From(p.Pod).IsRunning()
  if !isRunning {
    err := errors.Errorf(
      "pod {%s} is not running in namespace {%s}", 
      p.Pod.Name, 
      p.Pod.Namespace,
    )
    p.Errors = append(p.Errors, err)
  }

  return p
}
```

### Sample Usages
```go
cspops.New().
  GetFromKubernetes(PoolName).
  SetLabel("openebs.io/version", TargetVersion)
```

```go
unstructops.New().
  WithStore(store).
  WithGVK("openebs.io", "v1alpha1", "CStorPool").
  WithNamespace(PoolNamespace).
  GetFromKubernetes(PoolName).
  SaveUIDToStoreWithKey("pool.uid").
  SaveToStore(
    unstructops.GetValueFromPath(".spec.device"),
    unstructops.WithStoreKey("pool.device"),
  )   
```

```go
cspops.New().
  WithStore(store).
  GetFromKubernetes(PoolName).
  SaveUIDToStoreWithKey("csp.uid")
```
