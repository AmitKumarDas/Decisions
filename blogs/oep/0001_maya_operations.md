### Motivation
Ops approach defined in this proposal allows business logic i.e. code to be more cohesive. This approach also lets developers to express their logic that is readable and hence easy to understand.

### Issues without a proper ops approach -- some observations
- Rise of utils _(IMO this is an anti-pattern)_
- For example, below code snippets are used during integration tests
  - This ideally should be part of actual test logic
  - However, it is placed in file different from the original test logic
  - This leads to splitting of logic _(here test logic)_ across multiple packages
  - The resulting test logic is not **cohesive**
- Scattering of core test logic can be verified from
  - Methods names that do not reflect its true behavior
  - Explosion of arguments for a single method
- Once can also find test logic is tightly coupled with expect statements
  - In other words, this test logic is not usable in other scenarios
    - e.g. monitoring, self-heal, etc.


```go
// DeleteCSP ...
func (ops *Operations) DeleteCSP(spcName string, deleteCount int) {
  cspAPIList, err := ops.CSPClient.List(metav1.ListOptions{})
  Expect(err).To(BeNil())
  cspList := csp.
    ListBuilderForAPIObject(cspAPIList).
    List().
    Filter(csp.HasLabel(string(apis.StoragePoolClaimCPK), spcName), csp.IsStatus("Healthy"))
  cspCount := cspList.Len()
  Expect(deleteCount).Should(BeNumerically("<=", cspCount))

  for i := 0; i < deleteCount; i++ {
    _, err := ops.CSPClient.Delete(cspList.ObjectList.Items[i].Name, &metav1.DeleteOptions{})
    Expect(err).To(BeNil())
  }
}
```

```go
// IsHealthyCspCount ...
func (ops *Operations) IsHealthyCspCount(spcName string, expectedCspCount int) int {
  var cspCount int
  retries := maxRetry
  for i := 0; i < retries; i++ {
    cspCount = ops.GetHealthyCSPCount(spcName)
    if cspCount == expectedCspCount {
      return expectedCspCount
    }
    if retries == 0 {
      break
    }
    retries--
    time.Sleep(5 * time.Second)
  }
  return cspCount
}
```

```go
// ExecPod executes arbitrary command inside the pod
func (ops *Operations) ExecPod(podName, namespace, containerName string, command ...string) ([]byte, error) {
  var (
    execOut bytes.Buffer
    execErr bytes.Buffer
    err     error
  )
  config, err := ops.PodClient.GetConfig()
  Expect(err).To(BeNil(), "while getting config for exec'ing into pod")
  cset, err := ops.PodClient.GetClientSet()
  Expect(err).To(BeNil(), "while getting clientset for exec'ing into pod")
  req := cset.
    CoreV1().
    RESTClient().
    Post().
    Resource("pods").
    Name(podName).
    Namespace(namespace).
    SubResource("exec").
    Param("container", containerName).
    VersionedParams(&corev1.PodExecOptions{
      Container: containerName,
      Command:   command,
      Stdin:     false,
      Stdout:    true,
      Stderr:    true,
      TTY:       false,
      },
      scheme.ParameterCodec
    )

  exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
  Expect(err).To(BeNil(), "while exec'ing command in pod ", podName)

  err = exec.Stream(remotecommand.StreamOptions{
    Stdout: &execOut,
    Stderr: &execErr,
    Tty:    false,
  })
  Expect(err).To(BeNil(), "while streaming the command in pod ", podName, execOut.String(), execErr.String())
  Expect(execOut.Len()).Should(BeNumerically(">", 0), "while streaming the command in pod ", podName, execErr.String(), execOut.String())
  return execOut.Bytes(), nil
}
```

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
