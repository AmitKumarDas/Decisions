
#### Depends on
- https://github.com/AmitKumarDas/Decisions/blob/master/blogs/oep/0001_maya_operations.md

#### High Level Design

```go
// cmd/upgrade/app/v1alpha2/upgrade.go

// upgradePaths represents supported upgrade paths
type upgradePaths map[string]bool

// upgrades represents supported upgrades
type upgrades map[string]ops.Verifier

func buildUpgradePath(kind, source, target string) string {
  return kind + "-" + source + "-" + target
}

// GetInstance returns the upgrade instance
// based on the provided source & target
// versions
func GetInstance(kind, source, target string) ops.Verifier {
  return upgrades[buildUpgradePath(kind, source, target)]
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
func (r *Registrar) WithUpgradePath(kind, source, target string) *Registrar {
  r.path = buildUpgradePath(kind, source, target)
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
// cmd/upgrade/app/v1alpha2/add_upgrade_csp_082_to_090.go

import (
  csp "github.com/openebs/maya/cmd/upgrade/app/v1alpha2/082-to-090/cstorpool"
)

func init() {
  NewRegistrar().
    WithUpgradePath("cstorpool", "0.8.2", "0.9.0").
    WithUpgradeInstance(&csp.Upgrade{}).
    Register()
}
```

```go
// cmd/upgrade/app/v1alpha2/082-to-090/cstorpool/upgrade.go

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

func updateCStorPoolDeployment() ops.Verifier {
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
