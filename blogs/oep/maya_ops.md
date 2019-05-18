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
- [ ] Ability to package operations to be used by Upgrade Executor
- [ ] Ability to package operations as a Docker image
- [ ] Ability to package operations as a Kubernetes Job
- [ ] Ability to package operations to be used by Litmus Executor

#### Assumptions of MayaOps
- Users of MayaOps need to learn MayaOps syntax
- MayaOps is Go code and will be written inside:
  - a .go file
  - MayaOps kubernetes custom resource

### High Level Design

#### Core
```go
// pkg/ops/v1alpha1/interface.go

// Ops exposes contracts of an operation
//
// NOTE:
//  An operation is typically identified with
// a set of ordered steps
type Ops interface {
  // Run does the actual execution of 
  // operation
  Run() error
}

// Manager exposes management related
// contract(s) w.r.t an operation
type Manager interface {
  Manage() error
}
```

```go
// pkg/ops/v1alpha1/manager.go

// SingleOpsManager is a concrete implementation
// of Manager interface. As the name suggests
// it manages/runs a single operation
type SingleOpsManager struct {
  Ops Ops
}

// NewSingleOpsManager returns a new instance of
// SingleOpsManager
func NewSingleOpsManager(o Ops) *SingleOpsManager {
  return &SingleOpsManager{Ops: o}
}

// Manage execute the underlying operation
func (s *SingleOpsManager) Manage() error {
  return s.Ops.Run()
}

// MultiOpsManager is a concrete implementation
// of Manager interface. As the name suggests
// it manages/runs a set of operations
type MultiOpsManager struct {
  Items []Ops
}

// NewMultiOpsManager returns a new instance of
// MultiOpsManager
func NewMultiOpsManager(o ...Ops) *MultiOpsManager {
  m := &GroupOpsManager{}
  m.Items = append(m.Items, o...)
  return m
}

// Manage executes all the underlying operations
// managed by this instance
func (g *GroupOpsManager) Manage() error {
  var err error
  for _, op := range g.Items {
    err = NewSingleOpsManager(op).Manage()
    if err != nil {
      return err
    }
  }
  
  return nil
}
```

```go
// pkg/ops/v1alpha1/ops.go


// GetStoreFunc abstracts fetching the in-memory
// storage required while executing various steps
// of an operation
type GetStoreFunc func() map[string]interface{}

type BaseOps struct {
  ID           string
  Namespace    string
  ShouldSkip   bool
  Errors       []error
  RunSteps     []OpsStep
  GetStore     GetStoreFunc
}

type BaseOpsOption func(*BaseOps)

func New(opts ...BaseOpsOption) *BaseOps {
  b := &BaseOps{}
  for _, o := range opts {
    o(b)
  }
  return b
}
```

```go
// pkg/ops/v1alpha1/registrar.go

import (
  ops "github.com/openebs/maya/pkg/ops/v1alpha1"
)

type Registrar struct {
  key string
  manager ops.Manager
  registry map[string]ops.Manager
}

func NewRegistrar() *Registrar {
  return &Registrar{}
}

func (r *Registrar) WithKey(key string) *Registrar {
  r.key = key
  return r
}

func (r *Registrar) WithManager(manager ops.Manager) *Registrar {
  r.manager = manager
  return r
} 

func (r *Registrar) Register (
  registry[r.key] = r.manager
)
```

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
}

// OpsOption abstracts building
// an Ops instance
type OpsOption func(*Ops)

// OpsStep abstracts implementation
// of an Ops step
type OpsStep func(*Ops)

// WithStore is a OpsOption function that
// provide in-memory map to store result 
// of an operation and at the same time
// a placeholder to store config values
func WithStore(store map[string]interface{}) OpsOption {
  return func(p *Ops) {
    p.GetStore = func() map[string]interface{} {
      return store
    }
  }
}

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
//  Run is an implementation of Ops interface
//
// NOTE:
//  This is a final operation
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
// exected operation steps
//
// NOTE:
//  This is a final operation
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

```go
// pkg/ops/kubernetes/pod/v1alpha1/list.go
```

#### UseCase -- UpgradeExecutor

```go
// cmd/upgrade/execute.go

upgrades map[string]ops.Manager

type Executor struct {
  UpgradePath string
}

func (e *Executor) Run() error {
  manager := upgrades[e.UpgradePath]
  if manager == nil {
    return errors.New("failed to run upgrade: un-supported upgrade path {%s}", e.upgradePath)
  }

  err := manager.Manage()
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
  multiOpsMgr := ops.NewMultiOpsManager(
    080-to-090.PodShouldBeRunning(store),
    080-to-090.PodImageShouldGetUpdated(store),
  )

  ops.Registrar(upgrades).
    WithKey(080-to-090.UpgradePath).
    WithManager(multiOpsMgr).
    Register()
}
```

```go
// cmd/upgrade/0.8.0-0.9.0/upgrade.go
import (
  pod "github.com/openebs/maya/pkg/ops/kubernetes/pod/v1alpha1"
)

const (
  UpgradePath string = "080-to-090"
  CStorPoolImage string = "openebs.io/cstor-pool:0.9.0"
)

func PodShouldBeRunning(store map[string]interface{}) Ops {
  return pod.New(
    pod.WithStore(store),
    pod.WithID("is-pod-running"),
    pod.WithDesc("pod should be running"),
  ).Steps(
    pod.GetFromStore(".pod.object"),
    pod.ShouldBeRunning(),
  )
}

func PodImageShouldGetUpdated(id, store map[string]interface{}) Ops {
  return pod.New(
    pod.WithStore(store),
    pod.WithID("update-pod-image"),
    pod.WithDesc("pod's image should get updated"),
  ).Steps(
    pod.GetFromStore(".pod.object"),
    pod.SetImage(CStorPoolImage),
    pod.Update(),
  )
}

// Alternatively
func SetPodImageIfRunning(id, store map[string]interface{}) Ops {
  return pod.New(
    pod.WithStore(store),
    pod.WithID("update-pod-image-if-running"),
    pod.WithDesc("pod's image should get updated if is running"),
  ).Steps(
    pod.GetFromStore(".pod.object"),
    pod.ShouldBeRunning(),
    pod.SetImage(CStorPoolImage),
    pod.Update(),
  )
}
```

### Further Thinking
- Add skip property that skips rest of the steps
```go
func SetPodImageIfRunning(id, store map[string]interface{}) Ops {
  return pod.New(
    pod.WithStore(store),
    pod.WithID("update-pod-image-if-running"),
    pod.WithDesc("pod's image should get updated if is running"),
  ).Steps(
    pod.GetFromStore(".pod.object"),
    pod.SkipIfVersionNotEqualsTo("0.8.2")
    pod.ShouldBeRunning(),
    pod.SetImage(CStorPoolImage),
    pod.Update(),
  )
}
```
- Add ability to run the steps directly against an Ops instance
```go
func SetCStorPoolPodImageIfRunning() {
  pod.New().
    GetFromKubernetes("my-cstor-pool-pod", "openebs").
    SkipIfVersionNotEqualsTo("0.8.2").
    ShouldBeRunning().
    SetImage(CStorPoolImage).
    Verify()
}
```

### MayaOps as a Kubernetes Custom Resource
```yaml
kind: MayaOps
metadata:
  name: upgrade-080-To-090
  namespace: openebs
spec:
  go:
    constants:
      SourceVersion: "0.8.2"
      TargetVersion: "0.9.0"
      PoolNamespace: openebs
      PoolName: my-cstor-pool

    funcs:
      getCSPUID: |
        cspops.New().
          UseStore(MayaStore).
          GetFromKubernetes(PoolName).
          SaveUIDToStore("csp.uid")

      updateCSPDeployment: |
        deployops.New().
          UseStore(MayaStore).
          GetFromKubernetes(PoolName, PoolNamespace).
          SkipIfVersionNotEqualsTo(SourceVersion).
          SetLabel("openebs.io/version", TargetVersion).
          SetLabel("openebs.io/cstor-pool", PoolName).
          UpdateContainer(
            "cstor-pool",
            container.WithImage("quay.io/openebs/cstor-pool:"+ TargetVersion),
            container.WithENV("OPENEBS_IO_CSTOR_ID", FromStore("csp.uid")),
            container.WithLivenessProbe(
              probe.NewBuilder().
                WithExec(
                  k8scmd.Exec{
                    "command": 
                      - "/bin/sh"
                      - "-c"
                      - "zfs set io.openebs:livenesstimestap='$(date)' cstor-$OPENEBS_IO_CSTOR_ID"
                  }
                ).
                WithFailureThreshold(3).
                WithInitialDelaySeconds(300).
                WithPeriodSeconds(10).
                WithSuccessThreshold(1).
                WithTimeoutSeconds(30).
                Build()
            ),
          ).
          UpdateContainer(
            "cstor-pool-mgmt",
            container.WithImage("quay.io/openebs/cstor-pool-mgmt:" + TargetVersion),
            container.WithPortsNil(),
          ).
          UpdateContainer(
            "maya-exporter",
            container.WithImage("quay.io/openebs/m-exporter:" + TargetVersion),
            container.WithCommand("maya-exporter"),
            container.WithArgs("-e=pool"),
            container.WithTCPPort(9500),
            container.WithPriviledged(),
            container.WithVolumeMount("device", "/dev"),
            container.WithVolumeMount("tmp", "/tmp"),
            container.WithVolumeMount("sparse", "/var/openebs/sparse"),
            container.WithVolumeMount("udev", "/run/udev"),
          ).
          UpdateToKubernetes()
```
