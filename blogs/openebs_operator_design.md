
### OpenEBS Operator Specification -- WIP
- This is the typical default operator specification that gets generated and applied before the operator logic kicks in

```yaml
kind: OpenEBS
metadata:
  name: default
# specification of OpenEBS Operator
# specifications related to following managed via reconciliation
# - Install, Update, Upgrade, Probe / Monitor, Self Healing
# DOUBT: Should above be the job of openebs operator?
# THINK: Does it help the real human i.e. operator?
spec:
  # what version of openebs should be enabled/made available
  # version controls install, upgrade & rollback as well
  version: 0.7.0
  # volumeProvisioner is the specification of volume provisioner
  volumeProvisioner:
    # A value of false will un-install volume provisioner
    enabled: true
    # remote volume provisioner(s) if any
    remote: 
    - name: remote-123
      id: 123
    - name: remote-231
      id: 231
    # local volume provisioner
    local:
      nodeSelector:
  # apiServer is the specification of openebs api server
  apiServer:
    # A value of false will un-install all instances of api server
    enabled: true
    # nodeSelector translates to kubernetes pod's nodeSelector property
    nodeSelector:
  # cstor is the specification related to cstor storage engine
  cstor:
    # A value of false will un-install all cstor options, templates, tunables
    enabled: true
    snapshot:
      # operator can provide the name of custom templates to invoke CRUD operation(s) w.r.t snapshot
      # default value implies operator will decide / set the name of the template
      casTemplate:
        create: default
        delete: default
        list: default
    volume:
      # operator can provide the name of custom templates to invoke CRUD operation(s) w.r.t volume
      # default value implies operator will decide / set the name of the template
      casTemplate:
        create: default
        delete: default
        list: default
        read: default
  # jiva is the specification related jiva storage engine
  jiva:
    # A value of false will un-install all jiva related options, templates, tunables
    enabled: true
    volume:
      # A value of false will un-install all jiva volume related options, templates, tunables
      enabled: true
      casTemplate:
        create: default
        delete: default
        list: default
        read: default
  # localPV is the specification related to local PV storage engine
  localPV:
  # ndm is the specification related to ndm disk management solution
  ndm:
```

### KubeMonitor Operator - A reconciler to test, monitor kubernetes resource(s) -- WIP
- Will be used to test OpenEBS Operator
- Can be used to inject failures optionally
- Can be used to test if all openebs PV have dependent resources
```yaml
kind: KubeMonitor
spec:
  resource:
  - kind: PersistentVolume
    namespace: default
    labelSelector: provisioner=openebs
    expect: Bound
    owns:
    - kind: Deployment
      namespace: default
      labelSelector: byOwnerName, app=target
      expect: Online
    - kind: Deployment
      namespace: default
      labelSelector: byOwnerName, app=replica
      expect: Online
status:
```

```yaml
kind: KubeMonitor
spec:
  type: serial
  resource:
  - kind: Pod
    labelSelector: app=jiva
    fault:
      state: Deleted
    expect: Online
  - kind: Deployment
    labelSelector: app=maya
    fault:
      state: Deleted
    expect: Online
  - kind: Deployment
    labelSelector: app=provisioner
    fault:
      state: Deleted
    expect: Online
  - kind: Deployment
    labelSelector: app=provisioner
    fault:
      image: openebs/bad-image:0.0.1
    expect: Online
status:
```

### Thinking in Code -- WIP
- pkg/state/v1alpha1/state.go
```go
type Reconciler interface {
  // Reconcile abstracts the reconciliation process
  Reconcile() *Response
}
type Expecter interface {
  // Expect abstracts expectation of the state of a resource
  Expect(s ExpectedStatus) *Response
}
type Status string
const (
  InProgress Status = "InProgress"
  Success    Status = "Success"
  Error      Status = "Error"
)
type ExpectedStatus string
const (
  Online   ExpectedStatus = "Online"
  NotFound ExpectedStatus = "NotFound"
  Bound    ExpectedStatus = "Bound"
)
type Result string
const (
  Updated  Result = "Updated"
  Created  Result = "Created"
  Noop     Result = "Noop"
  Deleted  Result = "Deleted"
  Patched  Result = "Patched"
  NoError  Result = "NoError"
  Verified Result = "Verified"
)
type Response struct {
  Alias     string
  Namespace string
  Name      string
  Kind      string
  Status    Status
  Result    Result
  Message   string
  LastUpdatedAt time
}
```
- pkg/pod/v1alpha1/pod.go
```go
type pod struct {}
type podList struct {}
func (p *pod) Reconcile() reconcile.Response {}
func (p *podList) Reconcile() reconcile.Response {}

type builder struct {}
func Builder() *builder {}
func (b *builder) WithImage(image string) *builder {}
func (b *builder) Build() *pod {}

type listBuilder struct {}
func ListBuilder() *listBuilder {}
func (b *listBuilder) WithCount(n int) *listBuilder {}
func (b *listBuilder) Build() *podList {}
```
- pkg/deployment/v1alpha1/deployment.go
```go
type deployment struct {}
func (p *deployment) Reconcile() reconcile.Response {}

type builder struct {}
func Builder() *builder {}
func (b *builder) WithReplicas(n int) *builder {}
func (b *builder) Build() *deployment {}

func GetMayaAPIServerSpecs(imageTag string) *appsv1.deployment {}
func GetOpenEBSProvisionerSpecs(imageTag string) *appsv1.deployment {}
```

### Old Drafts
#### Draft 1 -- dated 24/Oct/2018

- Review Comment - This is not a specification i.e. statement
- Review Comment - This is more like a program/logic
- Review Comment - This is not what enduser will really want as a specification

```yaml
# desired state
spec:
  # refers to openebs components that should be available & in the specified state
  component:

  # simple desired state specification
  - include: currentCASTemplate; currentRunTask; currentOpenEBSProvisioner; currentMAPIServerEnv
  - exclude: labelSelector app eq jiva version lt current

  # finely granular desired state specification
  - exclude: selector kind eq runtask
    if: selector name contains -0.6.0
  - include: selector kind eq castemplate
    replace: metadata.labels.abc
    value: |
    if: labelSelector version eq 0.6.0; pathSelector metadata.labels.abc ne default

  # refers to features supported by openebs and should work as specified here
  feaure:
  - support: CreateJivaVolume; NVMEIsAvailable; HostHasISCSIInstalled
  - support: CreateCstorVolume
    if: selector openebsVersion gt 0.7.0
  - support: VolumeOnSparseFiles
    if: labelSelector stage in test,qa
  - support: CreateCstorVolume
    check: latency lt 1msec
    with: replicas 3
    interval: 1day or 1 hr # check support every 1 day on success or every 1 hr on failure

status:
  state: failed
  reason:
  starttime:
  endtime:
  feature:
  - support: CreateJivaVolume; NVMEIsAvailable
    state: passed
  - support: CreateCstorVolume
    if: selector version gt 0.7.0
    state: passed
  - support: CreateCstorVolume
    check: latency lt 1msec
    with: replicas 3
    interval: 3days
    state: passed
    starttime:
    endtime:
    prevState: failed
  component:
  - include: currentCASTemplate, currentRunTask, currentOpenEBSProvisioner
    state: passed
  - exclude: labelSelector app eq jiva version lt current
    state: passed
  - exclude: selector kind eq runtask
    if: selector name contains -0.6.0
    state: passed
  - include: selector kind eq castemplate
    replace: metadata.labels.abc
    value: |
    if: labelSelector version eq 0.6.0, pathSelector metadata.labels.abc ne default
    state: failed
    reason: blah blah blah
    starttime:
    endtime:
    count:
```

##### Appendix
- currentCASTemplate - ensure CAS Templates for current version of openebs are available in the cluster
- currentRunTask - ensure Run Tasks for current version of openebs are available in the cluster
- currentOpenEBSProvisioner - ensure openebs provisioner for current version of openebs is available in the cluster
- currentMAPIServerEnv - ensure maya api server pod has environment variables for current version of openebs
