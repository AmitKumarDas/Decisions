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
// pkg/ops/v1alpha1/ops.go

// Register exposes contract to enable any
// structure to get itself registered as
// an Ops interface
type Register interface {
  Instance() Ops
}

// Ops exposes all contracts necessary to
// participate in operation based workflow(s)
type Ops interface {
  // Init sets initialization options if any
  // before running the operations
  Init() error
  
  // Run does the actual execution of 
  // operations
  Run() error
}
```

```go
// pkg/kubernetes/pod/v1alpha1/podops.go

// GetStoreFunc abstracts fetching the in-memory
// storage required while executing various steps
// of an operation
type GetStoreFunc func() map[string]interface{}

type PodOps struct {
  ID           string
  Namespace    string
  Pod          *Pod
  Errors       []error
  InitList     []PodInitOption
  Steps        []PodBuildOption
  GetStore     GetStoreFunc
}

// PodInitOption abstracts the implementation
// of initializing a PodOps instance
type PodInitOption func(*PodOps)

// PodBuildOption abstract the implementation
// of building a PodOps instance
type PodBuildOption func(*PodOps)

// WithOpsStore is a PodOps init option to
// provide in-memory storage
func WithOpsStore(store map[string]interface{}) PodInitOption {
  return func(p *PodOps) {
    p.GetStore = func() map[string]interface{} {
      return store
    }
  }
}

// Ops returns a new instance of PodOps
func Ops(inits ...PodInitOption) *PodOps {
  p := &PodOps{}
  p.InitList = append(p.InitList, inits...)
  return p
}

// WithSteps sets PodOps instance with the
// provided steps that will be run later
//
// NOTE:
//  These steps form the core of pod 
// operations
func (p *PodOps) WithSteps(opts ...PodBuildOption) *PodOps {
  p.Steps = append(p.Steps, opts...)
  return p
}

// Init runs the initialization options
// for this pod ops instance
//
// NOTE:
//  Init is an implementation of Ops interface
func (p *PodOps) Init() error {
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
// this pod ops instance is used to execute as well.
//
// NOTE:
//  Run is an implementation of Ops interface
func (p *PodOps) Run() error {
  for _, step := range p.Steps {
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

#### Usage
```go
// cmd/upgrade/execute.go
```

```go
// cmd/upgrade/registrar.go
```

```go
// cmd/upgrade/0.8.0-0.9.0/pod/registrar.go

type Base struct {
  ID string
  Desc string
  Store map[string]interface{}
}
```

```go
// cmd/upgrade/0.8.0-0.9.0/pod/should_be_running.go

type ShouldBeRunning struct {
  Base
}

func NewShouldBeRunning(id, desc string, store map[string]interface{}) *ShouldBeRunning {
  return &ShouldBeRunning{ID: id, Desc: desc, Store: store}
}

// Instance implements ops.Register interface
func (i *ShouldBeRunning) Instance() Ops {
  return pod.Ops(
    pod.WithOpsStore(i.Store),
    pod.WithOpsID(i.ID),
    pod.WithOpsDesc(i.Desc),
  ).WithSteps(
    pod.WithOpsStoreObject("taskResult.podInfo.object"),
    pod.ShouldBeRunning(),
    pod.SaveTuple("name", "namespace"),
  )
}
```

```go
// cmd/upgrade/0.8.0-0.9.0/pod/update_image.go

type UpdateImage struct {
  Base
}

func NewUpdateImage(id, desc string, store map[string]interface{}) *UpdateImage {
  return &UpdateImage{ID: id, Desc: desc, Store: store}
}

// Instance implements ops.Register interface
func (i *UpdateImage) Instance() Ops {
  return pod.Ops(
    pod.WithOpsStore(i.Store),
    pod.WithOpsID(i.ID),
    pod.WithOpsDesc(i.Desc),
  ).Steps(
    pod.WithOpsStoreObject("taskResult.podInfo.object"),
    pod.SetImage("openebs.io/cstor-pool:0.12.1"),
    pod.Update(),
  )
}
```
