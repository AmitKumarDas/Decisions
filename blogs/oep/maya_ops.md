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


// GetStoreFunc abstracts fetching the in-memory
// storage required while executing various steps
// of an operation
type GetStoreFunc func() map[string]string

type BaseOps struct {
  ID           string
  Namespace    string
  ShouldSkip   bool
  Errors       []error
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
  key      string
  runner   ops.Runner
  registry map[string]ops.Runner
}

func NewRegistrar() *Registrar {
  return &Registrar{}
}

func (r *Registrar) WithKey(key string) *Registrar {
  r.key = key
  return r
}

func (r *Registrar) WithRunner(runner ops.Runner) *Registrar {
  r.runner = runner
  return r
} 

func (r *Registrar) Register (
  registry[r.key] = r.runner
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
  RunSteps     []OpsStep
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


#### UseCase -- Upgrade


```go
// cmd/upgrade/app/v1alpha2/upgrade.go

import (
  
)

const (
  // Get these from ConfigMap
  SourceVersion: "0.8.2"
  TargetVersion: "0.9.0"
  PoolNamespace: openebs
  PoolName: my-cstor-pool
)

func Execute() error {
  v := ops.NewOperationListVerifier(
    getCStorPoolUID(),
    updateCStorPoolDeployment(),
    setStoragePoolVersion(),
    setCStorPoolVersion(),
  )
  
  return v.Verify()
}

func getCStorPoolUID() ops.Verifier {
  return cspops.New().
    UseStore(MayaStore).
    GetFromKubernetes(PoolName).
    SaveUIDToStore("csp.uid")
}

func updateCStorPoolDeployment() ops.VerifierVerifier { {
  return deployops.New().
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
    UpdateToKubernetes().
    ShouldRolloutEventually(
      ops.WithRetryAttempts(10), 
      ops.WithRetryInterval(3s)
    )
}

func setStoragePoolVersion() ops.Verifier {
  return spops.New().
    GetFromKubernetes(PoolName).
    SetLabel("openebs.io/version", TargetVersion)
}

func setCStorPoolVersion() ops.Verifier {
  return cspops.New().
    GetFromKubernetes(PoolName).
    SetLabel("openebs.io/version", TargetVersion)
}
```


### MayaOps as a Kubernetes Custom Resource
- This is again programmatic with some yaml sugar
- This custom resource is called MayaLang

```go

type MayaLang struct {
  Spec MayaLangSpec
  Status MayaLangStatus
}

type MayaLangSpec struct {
  GO        MayaLangGO `json:"go"`
}

type MayaLangGO struct {
  Constants map[string]string  `json:"constants"`
  Functions []MayaLangFunction `json:"funcs"`
}

type MayaLangFunction struct {
  Name      string `json:"name"`
  Disabled  bool   `json:"disabled"`
  Body      string `json:"body"`
}

type MayaLangStatus struct {}

```

```yaml
kind: MayaLang
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
    - name: getCStorPoolUID
      body: |
        cspops.New().
          UseStore(MayaStore).
          GetFromKubernetes(PoolName).
          SaveUIDToStore("csp.uid")

    - name: updateCStorPoolDeployment
      body: |
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
          UpdateToKubernetes().
          ShouldRolloutEventually(
            ops.WithRetryAttempts(10), 
            ops.WithRetryInterval(3s)
          )

    - name: setStoragePoolVersion
      body: |
        spops.New().
          GetFromKubernetes(PoolName).
          SetLabel("openebs.io/version", TargetVersion)
          
    - name: setCStorPoolVersion
      body: |
        cspops.New().
          GetFromKubernetes(PoolName).
          SetLabel("openebs.io/version", TargetVersion)
```
