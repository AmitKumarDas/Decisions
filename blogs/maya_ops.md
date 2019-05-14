### Motivation
It has been my strongest desire to let `Maya` i.e. control plane for openebs provide a workflow or pipeline to run an ordered
set of tasks. Maya has so far strived to provide this in form of CASTemplates & RunTasks. While this is convenient it has 
lot of drawbacks. CASTemplates & RunTasks are provided as a declarative YAML. Finally, when higher order logic is required 
it resorts to golang based templating. One needs to understand the fact that `templating` cannot replace a programming 
language and is never turing complete.

Hence, this proof of concept is yet another attempt to provide **MayaOps** _(similar in spirit to GitOps or Github actions 
or so on)_. As one might have already guessed, this attempt will avoid declarative approach. However, without getting too 
much into programmatic versus declarative approach, let us list down the features of `MayaOps`.

#### What MayaOps should provide
- [ ] Simple way to code operations supported by Maya
- [ ] Provide proper error messages for invalid operations
- [ ] Operation steps can be ordered
- [ ] Ability to package operations to be used by Upgrade Executor
- [ ] Ability to package operations to be used as a replacement for RunTask(s)
- [ ] Ability to package operations to be used inside Ginkgo & Gomega code
- [ ] Ability to package operations as a Docker image
- [ ] Ability to package operations as a Kubernetes Job
- [ ] Ability to package operations to be used by Litmus Executor

#### Assumptions of MayaOps
- Users of MayaOps need to learn MayaOps syntax
- MayaOps is Go code and will be written inside a .go file
- MayaOps syntax will expose one or more hooks that needs to be filled in
- MayaOps SDK can packaged Maya Ops library as well as hooks as a single Docker image

### High Level Design

#### Core
```go
// pkg/ops/v1alpha1/interface.go

// Ops exposes contracts of an operation
//
// NOTE:
//  An operation is typically represented as
// a set of ordered steps
type Ops interface {
  Initializer
  Runner
}

type Initializer interface {
  // Init sets initialization options if any
  // before running the operation
  Init() error
}

type Runner interface {
  // Run does the actual execution of 
  // operation
  Run() error
}

type Factory interface {
  Instance() Runner
}
```

```go
// pkg/ops/v1alpha1/runner.go

type SingleRunner struct {
  Ops Ops
}

func NewSingleRunner(o Ops) *SingleRunner {
  return &SingleRunner{Ops: o}
}

func (s *SingleRunner) Run() error {
  err := s.Ops.Init()
  if err != nil {
    return err
  }

  return s.Ops.Run()
}

type GroupRunner struct {
  Items []Ops
}

func NewGroupRunner(o ...Ops) *GroupRunner {
  m := &GroupRunner{}
  m.Items = append(m.Items, o...)
  return m
}

func (g *GroupRunner) Run() error {
  var err error
  for _, op := range g.Items {
    err = NewSingleRunner(op).Run()
    if err != nil {
      return err
    }
  }
  
  return nil
}
```

```go
// pkg/kubernetes/pod/v1alpha1/podops.go

// NOTE:
//  This file (along with its package) deals
// with pod related operations only.

// GetStoreFunc abstracts fetching the in-memory
// storage required while executing various steps
// of an operation
type GetStoreFunc func() map[string]interface{}

type Ops struct {
  ID           string
  Namespace    string
  Pod          *Pod
  Errors       []error
  InitList     []OpsInit
  RunSteps     []OpsStep
  GetStore     GetStoreFunc
}

// OpsInit abstracts implementation
// of initializing an instance of Ops
type OpsInit func(*Ops)

// OpsStep abstracts implementation
// of an Ops step
type OpsStep func(*Ops)

// WithOpsStore is a Ops init option to
// provide in-memory storage
func WithOpsStore(store map[string]interface{}) OpsInit {
  return func(p *Ops) {
    p.GetStore = func() map[string]interface{} {
      return store
    }
  }
}

// NewOps returns a new instance of Ops
func NewOps(inits ...OpsInit) *Ops {
  p := &Ops{}
  p.InitList = append(p.InitList, inits...)
  return p
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

// Init runs the initialization options
// for this ops instance
//
// NOTE:
//  Init is an implementation of Ops interface
func (p *Ops) Init() error {
  for _, init := range p.InitList {
    if len(p.Errors) > 0 {
      return errors.New("%v", p.Errors)
    }
    init(p)
  }
  return nil
}

// Run executes the operations as steps in an ordered
// manner. The order that was used while creating
// this ops instance is used to execute as well.
//
// NOTE:
//  Run is an implementation of Ops interface
func (p *Ops) Run() error {
  for _, step := range p.RunSteps {
    if len(p.Errors) > 0 {
      return errors.New("%v", p.Errors)
    }
    step(p)
  }
  return nil
}
```

```go
// pkg/kubernetes/pod/v1alpha1/podlistops.go
```

#### UseCase -- UpgradeExecutor

```go
// cmd/upgrade/registrar.go

import (
  ops "github.com/openebs/maya/pkg/ops/v1alpha1"
)

type registrar struct {
  map[string]ops.Runner
}

type RegistrarBuilder struct {
  path string
  factory ops.Factory
}

func NewRegistrarBuilder() *RegistrarBuilder {
  return &RegistrarBuilder{}
}

func (r *RegistrarBuilder) WithPath(path UpgradePath) *RegistrarBuilder {
  r.path = path
  return r
}

func (r *RegistrarBuilder) WithFactory(factory ops.Factory) *RegistrarBuilder {
  r.factory = factory
  return r
} 

func (r *RegistrarBuilder) Register (
  registrar[r.path] = r.factory
)
```

```go
// cmd/upgrade/execute.go

type Executor struct {
  UpgradePath string
}

func (e *Executor) Run() error {
  factory := registrar[e.UpgradePath]
  if factory == nil {
    return errors.New("failed to run upgrade: invalid upgrade path {%s}", e.upgradePath)
  }

  err := factory.Instance().Run()
  if err != nil {
    return errors.Wrapf(err, "failed to execute upgrade for path {%s}", e.upgradePath)
  }

  return nil
}
```

```go
// cmd/upgrade/register_upgradepath_080_090.go

import (
  ops "github.com/openebs/maya/pkg/ops/v1alpah1"
  080-to-090 "github.com/openebs/maya/cmd/upgrade/0.8.0-0.9.0"
)

init() {
  store := map[string]interface{}{}
  grpRunner := ops.NewGroupRunner(
    080-to-090.NewPodShouldBeRunning("pod101", "pod should be running", store),
    080-to-090.NewPodUpdateImage("pod201", "pod's image should get updated", store),
  )

  RegistrarBuilder().
    WithPath(080-to-090.UpgradePath).
    WithRunner(grpRunner).
    Register()
}
```

```go
// cmd/upgrade/0.8.0-0.9.0/common.go

type Base struct {
  ID string
  Desc string
  Store map[string]interface{}
}

const (
  UpgradePath string = "080-to-090"

  PathToPodObject string = "taskResult.podInfo.object"
)
```

```go
// cmd/upgrade/0.8.0-0.9.0/pod_should_be_running.go

type PodShouldBeRunning struct {
  Base
}

func NewPodShouldBeRunning(id, desc string, store map[string]interface{}) *PodShouldBeRunning {
  return &PodShouldBeRunning{ID: id, Desc: desc, Store: store}
}

// Runner implements ops.Factory interface
func (i *PodShouldBeRunning) Instance() Ops {
  return pod.Ops(
    pod.WithOpsStore(i.Store),
    pod.WithOpsID(i.ID),
    pod.WithOpsDesc(i.Desc),
  ).Steps(
    pod.WithObjectFromStore(PathToPodObject),
    pod.ShouldBeRunning(),
    pod.SaveTuple("name", "namespace"),
  )
}
```

```go
// cmd/upgrade/0.8.0-0.9.0/pod_update_image.go

type PodUpdateImage struct {
  Base
}

func NewPodUpdateImage(id, desc string, store map[string]interface{}) *PodUpdateImage {
  return &PodUpdateImage{ID: id, Desc: desc, Store: store}
}

// Runner implements ops.Factory interface
func (i *PodUpdateImage) Instance() Ops {
  return pod.Ops(
    pod.WithOpsStore(i.Store),
    pod.WithOpsID(i.ID),
    pod.WithOpsDesc(i.Desc),
  ).Steps(
    pod.WithObjectFromStore(PathToPodObject),
    pod.SetImage("openebs.io/cstor-pool:0.12.1"),
    pod.Update(),
  )
}
```
