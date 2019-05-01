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
```

```go
// pkg/kubernetes/pod/v1alpha1/podops.go

type PodOps struct {
  Errors       []error
  InitOptions  []PodInitOption
  BuildOptions []PodBuildOption
}

func NewPodOps() *PodOps {}

func (o *PodOps) Run() error {}
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
// cmd/upgrade/0.8.0-0.9.0/pool/registrar.go
// cmd/upgrade/0.8.0-0.9.0/pool/is_healthy.go
// cmd/upgrade/0.8.0-0.9.0/pool/is_replica_count.go
// cmd/upgrade/0.8.0-0.9.0/pool/upgrade.go
```
