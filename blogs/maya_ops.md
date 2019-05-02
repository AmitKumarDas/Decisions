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

### High Level Design

#### Core
```go
// pkg/ops/v1alpha1/ops.go

type Register interface {
  Instance() Ops
}

type Ops interface {
  Init() error
  Run() error
}
```

```go
// pkg/kubernetes/pod/v1alpha1/podops.go

type GetStoreFunc func() map[string]interface{}

type PodOps struct {
  ID           string
  Namespace    string
  Pod          *Pod
  Errors       []error
  InitOptions  []PodInitOption
  BuildOptions []PodBuildOption
  GetStore     GetStoreFunc
}

type PodInitOption func(*PodOps)
type PodBuildOption func(*PodOps)

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
  p.InitOptions = append(p.InitOptions, inits...)
  return p
}

// Steps sets PodOps instance with the
// provided steps that will be run later
//
// NOTE:
//  These steps form the core of pod 
// operations
func (p *PodOps) Steps(opts ...PodBuildOption) *PodOps {
  p.BuildOptions = append(p.BuildOptions, opts...)
  return p
}

// Init runs the initialization options
// for this pod ops instance
//
// NOTE:
//  Init is an implementation of Ops interface
func (p *PodOps) Init() error {
  for _, i := range o.InitOptions {
    if len(p.Errors) > 0 {
      return errors.New("%v", p.Errors)
    }
    i(p)
  }
  return nil
}

// Run executes the build options in a ordered
// manner. The order that was used while creating
// this pod ops instance is used to execute as well.
//
// NOTE:
//  Run is an implementation of Ops interface
func (p *PodOps) Run() error {
  for _, b := range o.BuildOptions {
    if len(p.Errors) > 0 {
      return errors.New("%v", p.Errors)
    }
    b(p)
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
```

```go
// cmd/upgrade/0.8.0-0.9.0/pod/should_be_healthy.go

type ShouldBeHealthy struct {
  Store map[string]interface{}
}

func NewShouldBeHealthy(store map[string]interface{}) *ShouldBeHealthy {
  return &ShouldBeHealthy{Store: store}
}

// Instance implements ops.Register interface
func (i *ShouldBeHealthy) Instance() Ops {
  return pod.Ops(
    pod.WithOpsStore(i.Store),
    pod.WithOpsID("healthyPod"),
    pod.WithOpsDesc("op-pod-should-be-healthy"),
  ).Steps(
    pod.WithOpsStoreObject("taskResult.pod101.object"),
    pod.ShouldBeRunning(),
    pod.SaveTuple("name", "namespace"),
  )
}
```

```go
// cmd/upgrade/0.8.0-0.9.0/pod/is_replica_count.go
```

```go
// cmd/upgrade/0.8.0-0.9.0/pod/upgrade.go
```
