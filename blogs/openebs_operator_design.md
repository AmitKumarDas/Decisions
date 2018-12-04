
### OpenEBS Operator Specification -- WIP
- This is the typical default operator specification that gets generated and applied before the operator logic kicks in

```yaml
kind: OpenEBS
metadata:
  name: default
# specification of OpenEBS Operator
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
      # template corresponds to a subset of k8s Deployment spec applicable for local volume provisioner
      template:
  # apiServer is the specification of maya api server
  apiServer:
    # A value of false will un-install api server
    enabled: true
    casTemplate:
      cstor:
        snapshot:
          create: default
          delete: default
          list: default
        volume:
          create: default
          delete: default
          list: default
          read: default
      jiva:
        volume:
          create: default
          delete: default
          list: default
          read: default
    # template corresponds to a subset of k8s Deployment spec applicable for maya api server
    template:
```

### KubeNote Operator - A generic reconciler to spot kubernetes resource(s) & their state(s)
- Use-Cases:
  - Find presence of OpenEBS components
  - Verify if all openebs PVs have dependent resources with appropriate states
    - i.e. **detect stale resources**

#### UseCase #1 - Presence of OpenEBS components
```yaml
kind: KubeNote
spec:
  resource:
  - kind: Deployment
    namespace: openebs
    name: maya-apiserver
    expect: Progressing, Available
    owns:
    - kind: Pod
      expect: Running
    - kind: ReplicaSet
      expect: Available
    links:
    - kind: Service
      name: maya-apiserver-svc
      expect: Present
status:
```

Notes:
- `links` & `owns` are two approaches to query resources that are either _linked to_ or _owned by_ original resource
- There can be different set of expectations based on the resource being queried
- e.g. a Pod can have below expectations:
  - Running, Initialized, Ready, PodScheduled
- Similarly a Container can have below expectations:
  - Ready
- A Deployment can have following expectations:
  - Progressing, Available
- **Tip**: It might be good to have multiple expectations from the targeted resource

#### UseCase #2 - Verify if any OpenEBS volumes are orphaned
```yaml
kind: KubeNote
spec:
  resource:
  - kind: PersistentVolume
    namespace: default
    labelSelector: provisioner=openebs
    expect: Bound
    links:
    - kind: Deployment
      namespace: default
      labelSelector: pv=${PV_NAME}, app=target
      expect: Online
    - kind: Deployment
      namespace: default
      labelSelector: pv=${PV_NAME}, app=replica
      expect: Online
status:
```
- **Tip**: Use golang os.expand for variable expansion 
  - e.g. `labelSelector: pv=${PV_NAME}` will result into `labelSelector: pv=actual-name-of-pv`


### KubeTest - A generic reconciler to test kubernetes resource(s)
- Can be used to inject failures optionally
```yaml
kind: KubeTest
spec:
  # can monitor a serial in a serial manner or in parallel
  # it can take quite some time to get the monitoring result of each resource
  type: serial
  resource:
  - kind: Pod
    labelSelector: app=jiva
    fault:
      state: Deleted
    expect: Running
  - kind: Deployment
    labelSelector: app=maya
    fault:
      state: Deleted
    expect: Available
  - kind: Deployment
    labelSelector: app=provisioner
    fault:
      state: Deleted
    expect: Available
  - kind: Deployment
    labelSelector: app=provisioner
    fault:
      image: openebs/bad-image:0.0.1
    expect: Available
status:
```

### IscsiMonitor -- An operator to monitor iscsi sessions
- IscsiMonitor operator picks up iscsi targets and portals from PV objects
- Starts DaemonSet Pods on all the nodes initiating iscsi connection
- Each pod will calculate iscsi session status(-es) for the required iscsi target(s)
- Each pod will create a `IscsiMonitorItem` spec during the process of its execution
- IscsiMonitor will watch for `IscsiMonitorItem` spec(s) & update the `status.iscsiSessions` field
- IscsiMonitor will calculate the status.state field based on `status.iscsiSessions` field
- IscsiMonitor will delete the DaemonSet if all `IscsiMonitorItem` specs are in COMPLETED state
- IscsiMonitor will resync if one or more `IscsiMonitorItem` specs are in INIT state
```yaml
kind: IscsiMonitor
spec:
  volumeSelector:
  - label:
    annotation:
    namespace:
    name:
  job:
    nodeSelector:
    image:
    command:
    args:
status:
  # state value is derived based on individual status of status.iscsiSessions
  state:
  iscsiSessions:
  - node:
    target:
    portal:
    status:
    lastProbeTime:
    lastTransitionTime:
    lastUpdatedBy:
    reason:
    message:
```
```yaml
kind: IscsiMonitorItem
status:
  # either INIT or COMPLETED
  state:
  iscsiSessions:
  - node:
    portal:
    target:
    status:
    lastProbeTime:
    lastTransitionTime:
    lastUpdatedBy:
    reason:
    message:
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
