### Motivation
Maya Ops represents a pipeline which is targeted to implement a set of ordered instructions. Ops approach defined in this proposal is primarily meant to be used while writing `test logic` or `tools` that cater to specific domain. In other words this pattern seems fit to be called as a domain specific language i.e. **DSL**. It allows business logic i.e. code to be more cohesive and specific to a targeted domain. This approach lets developers to express logic which is more readable and hence should be easy to understand.

### Graduation Criteria
- Users find MayaOps to be GitOps friendly to manage OpenEBS
- Users find MayaOps to be easier than writing shell scripts

#### UseCase - CStor Pool Health Check
```go
// pkg/tools/spc/health_check_csp/main.go

err := ops.
  New().
  Desc(`
    As a test developer, I want to verify the number
    of Healthy cstor pools of a given SPC against an
    expected value.
  `).
  Add(
    cspop.New().
      List(csp.ListOpts(string(apis.StoragePoolClaimCPK), spcName)).
      Filter(csp.IsStatus("Healthy")).
      VerifyLenEQ(count)
  ).
  Start()
```

#### UseCase - OpenEBS Health Check
```go
// pkg/tools/openebs/health_check/main.go

// NOTE: This example shows how one can 
// add pipe.Operation instances and execute
//
// NOTE: This can alternatively be implemented
// as pipe.Pipeline, which has the ability to add
// operations and execute at the end

// in pkg/ops/v1alpha1
//
// type OperationInstance func() Operation

var verifications = []OperationInstance{
  VerifyMayaAPIServer,
  VerifyNDMDaemonSet,
}

func Verify() error {
  return pipe.AddAll(verifications).Start()
}

func VerifyMayaAPIServer() ops.Operation {
  return podop.
    New(
      ops.Desc(`
        As an openebs admin, I want to test if maya api 
        server is installed and all its pods are in running
        state
      `)
    ).
    WithNamespace(openebs).
    List(pod.ListOpts(MayaAPIServerLabelSelector)).
    Filter(pod.IsRunning()).
    VerifyLenEQ(MayaAPIServerPodCount)
}

func VerifyNDMDaemonSet() ops.Operation{
  return podop.
    New(
      ops.Desc(`
        As an openebs admin, I want to test if NDM daemon set 
        is installed and all its pods are in running state    
      `)
    ).
    WithNamespace(openebs).
    List(pod.ListOpts(NDMDaemonSetLabelSelector)).
    Filter(pod.IsRunning()).
    VerifyLenEQ(NDMDaemonSetPodCount)
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
// pkg/ops/v1alpha1/doc.go

// What is Ops?
//
// Ops can be thought of as a pipeline of ordered operations
```

```go
// pkg/ops/v1alpha1/interface.go

// Pipeline defines the contracts
// supported by a pipeline
type Pipeline interface {

  // Start the pipeline of operations
  Start() error
}

// Operation defines the contracts
// supported by a pipeline operation
type Operation interface {

  // Run does the actual execution of 
  // a piped operation
  Run() error
}

// OperationInstance is a typed function that abstracts
// fetching an instance of Operation
type OperationInstance func() Operation
```

```go
// pkg/ops/v1alpha1/error.go

// FailedPipe represents an error that occurred
// while running a pipeline
type FailedPipe struct {
  Statement string
  Result    string
  Reason    string
}

// String representation of FailedPipe instance
type (f *FailedPipe) String() string {
  return "statement: %s,\nresult: %s,\nreason: %s\n"
}

// Error returns self which is an error
type (f *FailedPipe) Error() string {
  return f
}
```

```go
// pkg/ops/v1alpha1/pipe.go

// Option is a typed function that abstracts setting a 
// DefaultPipe instance
type Option func(*DefaultPipe)

// DefaultPipe is a default Pipeline implementation
type DefaultPipe struct {
  Description string
  Operations []Operation
}

// New returns a new instance of DefaultPipe
func New(opts ...Option) DefaultPipe {
  d := &DefaultPipe{Description: "default pipeline"}

  for _, option := range opts {
    option(d)
  }
  return d
}

func (d *DefaultPipe) handleError(err error) error {
  return &FailedPipe{
    Statement: d.Description,
    Result: "Failed",
    Reason: err.Error(),
  }
}

// Desc sets the description of this pipeline
func (d *DefaultPipe) Desc(msg string) *DefaultPipe{
  d.Description = msg
  return d
}

// Add adds the given operation to this pipeline
func (d *DefaultPipe) Add(op Operation) *DefaultPipe{
  d.Operations = append(d.Operations, op)
  return d
}

// AddAll adds all the given operations to this pipeline
func (d *DefaultPipe) AddAll(ops []OperationInstance) *DefaultPipe{
  for _, op := range ops {
    d.Operations = append(d.Operations, op())
  }
  return d
}

// Start starts execution of this pipeline by
// running its listed operations
func (d *DefaultPipe) Start() error {
  var err error
  
  if len(d.Operations) == 0 {
    return d.handleError(errors.New("failed to start pipeline: nil operations"))
  }
  
  for _, op := range d.Operations {
    err = op.Run()
    if err != nil {
      return d.handleError(err)
    }
  }

  return nil
}
```

```go
// pkg/ops/v1alpha1/operation.go

// Operation composes of all common fields
// required for any operation
type Operation struct {
  ID           string
  Namespace    string
  ShouldSkip   bool
  Errors       []error
  Store        map[string]string
}

// OperationOption is a typed function that
// abstracts building the operation instance
type OperationOption func(*Operation)

// WithStore provides an in-memory store to do the following:
//
//  [1] store results of each step of an operation &
//  [2] store config values if any
func WithStore(store map[string]string) OperationOption {
  return func(p *Operation) {
    if store == nil {
      p.Errors = append(
        p.Errors, 
        errors.New("failed to build operation: nil store provided"),
      )
      return
    }

    p.Store = store
  }
}

// NewOperation returns a new instance of operation
func NewOperation(opts ...OperationOption) *Operation {
  b := &Operation{}
  for _, o := range opts {
    o(b)
  }

  return b
}

// isContinue flags if this operation should continue 
// executing
func (p *Operation) isContinue() bool {
  return !(len(p.Errors) > 0 || p.ShouldSkip)
}

// errorOrNil returns error if any
func (p *Operation) errorOrNil() error {
  if len(p.Errors) > 0 {
    return errors.New("%v", p.Errors)
  }
  
  return nil
}
```

### Low Level Parts
- This is implementation of pkg/ops/v1alpha1/ interfaces

```go
// pkg/ops/kubernetes/pod/v1alpha1/operation.go

import (
  ops "github.com/openebs/maya/pkg/ops/v1alpha1"
  pod  "github.com/openebs/maya/pkg/kubernetes/pod/v1alpha1"
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

type Operation struct {
  ops.Operation
  Pod          *Pod
}

// Step is a typed function that abstracts 
// implementating a step of an operation
type Step func(*Operation)

type OperationRunner struct {
  Operation *Operation
  Steps     []Step
}

// New returns a new instance of Ops
func New(opts ...ops.OperationOption) *OperationRunner {
  o := &Operation{Operation: ops.NewOperation()}
  r := &OperationRunner{Operation: o}
  return r.withOptions(opts...)
}

func (o *OperationRunner) withOptions(opts ...ops.OperationOption) *OperationRunner{
  for _, option := range opts {
    option(o)
  }
  return o
}

// Steps builds this operation instance with steps that get
// executed as part of running this operation
func (p *OperationRunner) Steps(step ...OperationStep) *OperationRunner {
  p.Steps = append(p.Steps, step...)
  return p
}

// Run executes this operation by executing its steps in an ordered
// manner. The order that was used while creating this operation
// instance is used to execute as well.
//
// NOTE:
//  Run is an implementation of ops.Runner interface
func (p *OperationRunner) Run() error {
  for _, step := range p.Steps {
    if !p.isContinue() {
      break
    }

    step(p)
  }

  return p.errorOrNil()
}
```

```go
// pkg/ops/kubernetes/pod/v1alpha1/steps.go

// ShouldNotBeRunning adds should not be running check
// as a step in the operation
func (p *Operation) ShouldNotBeRunning() *Operation {
  p.Steps(shouldNotBeRunning())
  return p
}

// ShouldBeRunning adds should be running check
// as a step in the operation
func (p *Operation) ShouldBeRunning() *Operation {
  p.Steps(shouldBeRunning())
  return p
}

// shouldNotBeRunning errors if pod is running
func shouldNotBeRunning() OperationStep {
  return func(p *Operation) {
    isRunning := pod.From(p.Pod).IsRunning()
    if isRunning {
      err := errors.Errorf(
        "want - pod should not run: actual - pod {%s/%s} phase {%s}", 
        p.Pod.Namespace,
        p.Pod.Name,
        p.Pod.Status.Phase,
      )
      p.Errors = append(p.Errors, err)
    }
  }  
}

// shouldBeRunning errors if pod is not running
func shouldBeRunning() OperationStep {
  return func(p *Operation) {
    isRunning := pod.From(p.Pod).IsRunning()
    if !isRunning {
      err := errors.Errorf(
        "want - pod should run: actual - pod {%s/%s} phase {%s}", 
        p.Pod.Namespace,
        p.Pod.Name,
        p.Pod.Status.Phase,
      )
      p.Errors = append(p.Errors, err)
    }
  }
}
```

### Sample Usages
```go
cspop.New().
  GetFromKubernetes(PoolName).
  SetLabel("openebs.io/version", TargetVersion).
  UpdateToKubernetes().
  Run()
```

```go
unstructop.
  New(
    commonop.WithStore(store),
    commonop.WithGVK("openebs.io", "v1alpha1", "CStorPool"),
    commonop.WithNamespace(PoolNamespace),
  ).
  GetFromKubernetes(PoolName).
  SaveUIDToStoreWithKey("mypool.uid").
  SaveToStore(
    unstructop.GetValueFromPath(".spec.device"),
    unstructop.WithStoreKey("mypool.device"),
  ).
  Run()
```

```go
cspop.New(commonop.WithStore(store)).
  GetFromKubernetes(PoolName).
  SaveUIDToStoreWithKey("mycsp.uid").
  Run()
```

### Research
- Use of Gopkg.toml for every tool / folder
- Use of Dockerfile for every tool / folder
- Config Priority:
  - flags via github.com/spf13/pflag
  - environment variables
  - constants
