## World of Operators


### OpenEBS - A Operator to install & manage OpenEBS components - WIP
```yaml
kind: OpenEBS
metadata:
  name: default
# specification of OpenEBS Operator
spec:
  # what version of openebs should be enabled/made available
  # version controls install, upgrade & rollback all components
  # of openebs
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
        # currently this is cluster scoped
        # eventually, this should be namespaced
        namespace:
        config: # config to override castemplate default config; existing helm values should fit in
      spec: # will get populated via [1] castemplate && config or [2] can be inline; [1] overrides [2]
  # apiServer is the specification of maya api server
  apiServer:
    # A value of false will un-install api server
    enabled: true
    castemplate: # optional; if mentioned will populate the spec
      name:
      apiVersion:
      # currently this is cluster scoped
      # eventually, this should be namespaced
      namespace:
      config: # config to override castemplate default config; existing helm values should fit in
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
      config: # overrides default config present at cas template; existing helm values should fit in
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
Till now we have been talking in terms of kubernetes controllers where above specs are the desired states & reconcile
logic will do the hard work to get the desired state into actual state.

At any point of time, if you feel the need for a workflow based _(i.e. pipeline based task execution)_ approach, have a
look at [CASTLang](https://github.com/AmitKumarDas/Decisions/blob/master/blogs/castlang_spinsoff_operator.md).
