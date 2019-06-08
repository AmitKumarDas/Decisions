### Motivation
Maya Ops represents a pipeline which is targeted to implement a set of ordered instructions. Ops approach defined in this proposal is primarily meant to be used while writing `test logic` or `tools` that cater to specific domain. In other words this pattern seems fit to be called as a domain specific language i.e. **DSL**. It allows business logic i.e. code to be more cohesive and specific to a targeted domain. This approach lets developers to express logic which is more readable and hence easy to understand.

### Graduation Criteria
- Users find MayaOps to be easier than executing kubectl commands
- Users find MayaOps to be easier than writing shell scripts

#### UseCase - CStor Pool Health Check
```go
// pkg/tools/spc/health_check_csp/main.go

err := ops.New().
  Desc(`
    As a test developer, I want to verify the number
    of Healthy cstor pools of a given SPC against an
    expected value. I would also like to retry above
    expectation for a given number of attempts before
    giving up.
  `).
  WithRunner(
    cspops.New().Steps(
      cspops.List(csp.ListOpts(string(apis.StoragePoolClaimCPK), spcName)),
      cspops.Filter(csp.IsStatus("Healthy")),
      cspops.VerifyLenEQ(count),
      cspops.DeleteList(),    
    ),
    ops.RetryOnError(20, "3s"),  
  ).
  Run()
```

#### UseCase - OpenEBS Health Check
```go
// pkg/tools/openebs/health_check/main.go

type VerifyOpenEBSFn func() ops.Runner

var VerifyOpenEBSFns = []VerifyOpenEBSFn{
  VerifyMayaAPIServer,
  VerifyNDMDaemonSet,
}

func VerifyOpenEBS() error {
  for _, fn := range VerifyOpenEBSFns {
    err := fn().Run()
    if err != nil {
      return err
    }
  }

  return nil
}

func VerifyMayaAPIServer() ops.Runner {
  return ops.New().
    Desc(`
      As an openebs admin, I want to test if maya api 
      server is installed and all its pods are in running
      state
    `).
    WithRunner(
      podops.New().Steps(
        podops.WithNamespace(openebs),
        podops.List(pod.ListOpts(MayaAPIServerLabelSelector)),
        podops.Filter(pod.IsRunning()),
        podops.VerifyLenEQ(MayaAPIServerPodCount),
      ),
      ops.RetryOnError(10, "3s"),
    )
}

func VerifyNDMDaemonSet() ops.Runner{
  return ops.New().
    Desc(`
      As an openebs admin, I want to test if NDM daemon set 
      is installed and all its pods are in running state
    `).
    WithRunner(
      podops.New().Steps(
        podops.WithNamespace(openebs),
        podops.List(pod.ListOpts(NDMDaemonSetLabelSelector)),
        podops.Filter(pod.IsRunning()),
        podops.VerifyLenEQ(NDMDaemonSetPodCount),      
      ),
      ops.RetryOnError(10, "3s"),
    )
}
```

### High Level Design

#### Core - Drop 0
- Introducing **ops** pattern
- Ops is the shorthard notation for **MayaOperations**
- This concept is built around **builder pattern**
  - In this case each build method is executed immediately
- _Background_:
  - Ops is built on top of Maya's core builder & predicate functions
  - Maya's core builder pattern is however based on lazy executions
  - In other words, core builder execute in their final build methods
  - To repeat, Ops builder differs from core builder due to its immediate execution style

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

type FailedOperation struct {
  Statement string
  Result    string
  Reason    string
}

type (f *FailedOperation) String() string {
  return "statement: %s,\nresult: %s,\nreason: %s\n"
}

type (f *FailedOperation) Error() string {
  return f
}

type Option func(*DefaultRunner)

type DefaultRunner struct {
  Description string
  Retry       Retry
}

type Retry struct {
  Attempts  int
  Interval  string
}

func New() DefaultRunner {
  return &DefaultRunner{}
}

func (d *DefaultRunner) handleError(err error) error {
  return &FailedOperation{
    Statement: d.Description,
    Result: "Failed",
    Reason: err.Error(),
  }
}

func (d *DefaultRunner) setOptions(opts ...Option) {
  for _, option := range opts {
    option(d)
  }
}

func (d *DefaultRunner) Desc(msg string) *Default{
  d.Description = msg
  return d
}

func (d *DefaultRunner) Run(runner ops.Runner, opts ...ops.Option) error {
  d.setOptions(opts...)
  var err error
  for _ := range d.Retry.Attempts {
    err = runner.Run()
    if err == nil {
      return nil
    }
    sleep(d.Retry.Interval)
  }

  return d.handleError(err)
}


// TODO Refactor Below !!!


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
  SetLabel("openebs.io/version", TargetVersion).
  UpdateToKubernetes()
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
