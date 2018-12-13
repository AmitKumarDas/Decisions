## World of Operators


### OpenEBS - A Operator to install & manage OpenEBS components - WIP
```yaml
kind: OpenEBS
metadata:
  name: default
# specification of OpenEBS Operator
spec:
  # what version of openebs should be enabled/made available
  # version controls install, upgrade & rollback as well
  version: 0.8.0
```

- Below is the **generated** specification from above original spec
- This gets applied before operator logic kicks in
```yaml
kind: OpenEBS
metadata:
  name: default
# specification of OpenEBS Operator
spec:
  # what version of openebs should be enabled/made available
  # version controls install, upgrade & rollback as well
  version: 0.8.0
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
      castemplate: # optional; if mentioned will populate the spec
        name:
        apiVersion:
        config: # config to override castemplate default config
      spec: # will get populated via [1] castemplate && config or [2] can be inline; [1] overrides [2]
  # apiServer is the specification of maya api server
  apiServer:
    # A value of false will un-install api server
    enabled: true
    castemplate: # optional; if mentioned will populate the spec
      name:
      apiVersion:
      config: # config to override castemplate default config
    spec: # will get populated via [1] castemplate && config or [2] can be inline; [1] overrides [2]
```

### KubeNote Operator - A generic reconciler to spot kubernetes resource(s) & their state(s) - WIP
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


### KubeTest Operator - A generic reconciler to test kubernetes resource(s) - WIP
- Can be used to inject failures optionally
```yaml
kind: KubeTest
spec:
  # can monitor in a serial order or in parallel
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

### IscsiMonitor - Operator to monitor iscsi sessions - WIP
- IscsiMonitor operator filters specific nodes _(phase 1)_
- IscsiMonitor operator picks up iscsi targets and portals from PV objects _(phase 2)_
  - These details are used later to check iscsi sessions on the node
- Starts DaemonSet Pods on all the filtered nodes 
  - Check for the presence of iscsi binary _(phase 1)_
  - Run checks w.r.t iscsi session _(phase 2)_
- Each DaemonSet pod will create & update a `IscsiMonitorJob` custom resource during the process of its execution
- IscsiMonitor controller will watch its own spec as well as `IscsiMonitorJob` spec(s) & reconcile
- IscsiMonitor controller will delete DaemonSet once monitoring check is done for that reconcile operation
- Owner references will be set in such a way that all `IscsiMonitorJob` resources get deleted once DaemonSet gets deleted

```yaml
kind: IscsiMonitor
metadata:
  name: TestMyIscsi
spec:
  daemonset: # daemon set specifications
    castemplate: # optional; if mentioned will populate the spec
      name:
      config: # overrides default config present at cas template
    spec: # will get generated either via [1] above castemplate && config or [2] can be inline; [1] overrides [2]
status:
  phase:
  conditions:
  - check: Binary # phase 1
    type: Available
    status: "true"
    node:
    lastProbeTime:
    lastTransitionTime:
    lastUpdateTime:
    lastUpdatedBy:
    reason:
    message:
  - check: Session # phase 2
    type: Available
    status: "true"
    node:
    target:
    portal:
    lastProbeTime:
    lastTransitionTime:
    lastUpdateTime:
    lastUpdatedBy:
    reason:
    message:
```
```yaml
kind: IscsiMonitorJob
status:
  # Initializing, Completed, Failed, etc
  phase:
  conditions:
  - check: Binary # phase 1
    type: Available
    status: "true"
  - check: Binary # phase 1
    type: Permission
    status: "true"
  - check: Session # phase 2
    type: Available
    status:
    node:
    portal:
    target:
    lastTransitionTime:
    lastUpdateTime:
    lastUpdatedBy:
    reason:
    message:
```

### Case for Workflow Oriented Operator
- Read it here [CASTLang](https://github.com/AmitKumarDas/Decisions/blob/master/blogs/castlang_spinsoff_operator.md).

### Thinking in Code -- My Rough Notes
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

### Old & Rejected Drafts
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
