### MayaLang
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

### Sample UseCase - CStor Pool Upgrade

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
