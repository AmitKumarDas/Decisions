
#### UseCase -- Upgrade - v1alpha2 - Drop 1

```go
// cmd/upgrade/app/v1alpha2/upgrade.go

// upgradePaths represents supported upgrade paths
type upgradePaths map[string]bool

// upgrades represents supported upgrades
type upgrades map[string]ops.Verifier

func buildUpgradePath(source, target string) string {
  return source + "-" + target
}

// GetInstance returns the upgrade instance
// based on the provided source & target
// versions
func GetInstance(source, target string) ops.Verifier {
  return upgrades[buildUpgradePath(source, target)]
}

// Registrar is a utility type to register
// supported upgrades
type Registrar struct {
  path       string
  instance   ops.Verifier
}

// NewRegistrar returns a new instance of 
// Registrar
func NewRegistrar() *Registrar {
  return &Registrar{}
}

// WithUpgradePath sets the upgrade path by 
// accepting source & target versions
func (r *Registrar) WithUpgradePath(source, target string) *Registrar {
  r.path = buildUpgradePath(source, target)
  return r
}

// WithUpgradeInstance sets the upgrade instance
// that is applicable to the particular upgrade
// path
func (r *Registrar) WithUpgradeInstance(instance ops.Verifier) *Registrar {
  r.instance = instance
  return r
} 

// Register registers an upgrade instance
func (r *Registrar) Register() {
  upgradePaths[r.path] = true
  upgrades[r.path] = r.instance
}
```

```go
// cmd/upgrade/app/v1alpha2/add_upgrade_082_to_090.go

import (
  u082to090 "github.com/openebs/maya/cmd/upgrade/app/v1alpha2/082-to-090"
)

func init() {
  NewRegistrar().
    WithUpgradePath("0.8.2", "0.9.0").
    WithUpgradeInstance(&u082to090.Upgrade{}).
    Register()
}
```

```go
// cmd/upgrade/app/v1alpha2/082-to-090/upgrade.go

const (
  SourceVersion string = "0.8.2"
  TargetVersion string = "0.9.0"
  
  // Get these at runtime 
  // e.g. from ConfigMap
  PoolNamespace string = "openebs"
  PoolName string = "my-cstor-pool"
)

var store = map[string]string{}

type Upgrade struct {}

func (u *Upgrade) Verify() error {
  return ops.NewOperationListVerifier(
    getCStorPoolUID(),
    updateCStorPoolDeployment(),
    setStoragePoolVersion(),
    setCStorPoolVersion(),
  ).Verify()
}

func getCStorPoolUID() ops.Verifier {
  return cspops.New().
    WithStore(store).
    GetFromKubernetes(PoolName).
    SaveUIDToStoreWithKey("csp.uid")
}

func updateCStorPoolDeployment() ops.VerifierVerifier {
  return deployops.New().
    WithStore(store).
    GetFromKubernetes(PoolName, PoolNamespace).
    SkipIfVersionNotEqualsTo(SourceVersion).
    SetLabel("openebs.io/version", TargetVersion).
    SetLabel("openebs.io/cstor-pool", PoolName).
    UpdateContainer(
      "cstor-pool",
      container.WithImage("quay.io/openebs/cstor-pool:"+ TargetVersion),
      container.WithENV("OPENEBS_IO_CSTOR_ID", GetValueFromStoreKey("csp.uid")),
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

### UseCase -- Upgrade as a Kubernetes Custom Resource - Drop 2
- Note that this is completely programmatic with some yaml toppings
- One will get typical compile time errors if any
  - However, these compile & runtime failures will be known at runtime
  - Consider this as a dynamic language with full benefits of any statically compiled one
- This is implemented as a Kubernetes custom resource
  - Its specifications is understood by a Kubernetes based controller
- This custom resource is called as **MayaLang**
- This resource will be watched by a controller named **mayalang-controller**
  - It will watch for any MayaLang resources
  - It will create a .go file based on this MayaLang resource
  - It will build a Job that consists of `m-lang` docker image
  - It will load this .go file at an executable path of above Job's Pod
    - Job Pod has got only one thing to do; i.e. `go run` above .go file
  - It will apply this Job against Kubernetes cluster

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
  Name        string `json:"name"`
  Disabled    bool   `json:"disabled"`
  Description string `json:"desc"`
  Body        string `json:"body"`
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
    vars:
      SourceVersion: "0.8.2"
      TargetVersion: "0.9.0"
      PoolNamespace: openebs
      PoolName: my-cstor-pool
    
    funcs:
    - name: getCStorPoolUID
      body: |
        cspops.New().
          WithStore(store).
          GetFromKubernetes(PoolName).
          SaveUIDToStoreWithKey("csp.uid")

    - name: updateCStorPoolDeployment
      body: |
        deployops.New().
          WithStore(store).
          GetFromKubernetes(PoolName, PoolNamespace).
          SkipIfVersionNotEqualsTo(SourceVersion).
          SetLabel("openebs.io/version", TargetVersion).
          SetLabel("openebs.io/cstor-pool", PoolName).
          UpdateContainer(
            "cstor-pool",
            container.WithImage("quay.io/openebs/cstor-pool:"+ TargetVersion),
            container.WithENV("OPENEBS_IO_CSTOR_ID", GetValueFromStoreKey("csp.uid")),
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
            ops.WithRetryInterval("3s")
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

### UseCase -- Lots of Custom Resources - Drop 3
#### Why?
- Dedicated controllers helps in achieving specific logic per resource kind

#### How?
- Upgrade will take help of some custom resources to perform upgrade
- These custom resources will have individual controllers
  - e.g. CStorPoolUpgrade, CStorVolumeUpgrade, JivaVolumeUpgrade, etc
- These controllers will make use of appropriate **UpgradeRecipe** to perform upgrade
  - Recipes provides a declarative model to code the requirements

#### Why UpgradeRecipe(s)?
- It helps in having only k8s controller related logic inside controllers
- Upgrade/Domain specific logic will reside in Recipes

#### UpgradeRecipe -- Deep Dive
- This will be a Kubernetes CR that is used to hold **static content**
- UpgradeRecipe is a thin abstraction around MayaLang
- It delegates to MayaLang to perform actions required to be at the desired upgrade state
- A recipe in itself does not have a controller
- A recipe is only used to create MayaLang resource
  - MayaLang resource has its controller/watcher

#### Upgrade controller(s)
- Controllers will watch for its respective resources
  - In this case these resources are:
    - CStorPoolUpgrade, 
    - CStorVolumeUpgrade, 
    - JivaVolumeUpgrade
- Controller logic will refer to corresponding recipe to perform following tasks:
  - Build a MayaLang resource based on the recipe
  - Fill in appropriate constants against MayaLang resource
  - Apply MayaLang resource against the K8s cluster
  - Delete MayaLang resource when its status shows completed
  - Update its watched resource's status

```yaml
kind: CStorPoolUpgrade
metadata:
  name: my-csp-upgrade
spec:
  sourceVersion:
  targetVersion:

  # optional; controller can decide the
  # recipe based on source & target
  # versions
  recipe:
    name: 

  # upgrade all cstor pools
  # optional; can be filled up by controller
  all: true

  # or upgrade based on the names
  # is optional 
  # controller sets spec.all to true
  # if this list is nil
  pools:
  - name: csp-pool1
  - name: csp-pool2
  - name: csp-pool3
status:
```

```yaml
kind: UpgradeRecipe
metadata:
  name: cstorpool-upgrade-from-082-to-090
  namespace: openebs
  labels:
    cstorpool.upgrade.openebs.io/sourceversion: 
    cstorpool.upgrade.openebs.io/targetversion:
spec:
  mlang:
    spec:
      go:
        vars:
          SourceVersion: "0.8.2"
          TargetVersion: "0.9.0"

          # PoolNamespace: "<filled-in-at-runtime>"
          # PoolName: "<filled-in-at-runtime>"

        funcs:
        - name: saveCStorPoolDetails
          body: |
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

        - name: setCStorPoolVersion
          body: |
            cspops.New().
              GetFromKubernetes(PoolName).
              SetLabel("openebs.io/version", TargetVersion)
```
