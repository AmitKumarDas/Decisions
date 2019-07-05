### Motivation
OpenEBS has various components that need to be in running state to let its volumes be consumed by higher applications. As the
number of components grow, have inter-dependency, need upgrades it becomes un-manageable to let this be managed by a human.
OpenEBS operator should ideally take over most of the manual operations needed to let OpenEBS run smoothly.

### Points to Consider
- Number of openebs components can vary
- Some of the openebs components can be optional
- Each component might need a special handling

### Possible Custom Resources
- Operator
- ReleaseConfig
- ConfigInjector

### Sample Schemas
#### Operator Custom Resource
This is the only custom resource that is watched and is reconciled by the operator.

#### Operator - Sample 1
```yaml
kind: Operator
metadata:
  name: the-one-and-only
  namespace: openebs
spec:
  # installs or upgrades OpenEBS to
  # this specified version
  version:

  # config is optional; has detailed
  # specifications of supported components
  #
  # operator searches for a config with
  # same name as that of operator
  config:

  # disable is optional;
  #
  # if set to true it disables reconciling the
  # operator
  disable:

  # components is optional; operator
  # installs or upgrades all components
  # if this is not set
  #
  # can install specified components as well
  #
  # will un-install the ones that are not
  # specified and are running due to 
  # earlier installations
  components:

  # crds is optional; operator installs or updates
  # the crds
  #
  # can install or update specified crds as well
  #
  # will un-install the ones that are not specified
  # and were installed due to earlier installations
  crds:
  
  # storageclasses is optional; operator installs or
  # updates the storageclasses
  #
  # can install or update specified storageclasses
  # as well
  #
  # will un-install the ones that are not specified
  # and were installed due to earlier installations
  storageclasses:
```

#### Operator - Sample 2
```yaml
kind: Operator
metadata:
  name: the-one-and-only
  namespace: openebs
spec:
  version: 1.0.0  
```


#### Operator - Sample 3
```yaml
kind: Operator
metadata:
  name: the-one-and-only
  namespace: openebs
spec:
  version: 1.1.0
  components:
  - name: LocalProvisioner
```

#### Operator - Sample 4
```yaml
kind: Operator
metadata:
  name: the-one-and-only
  namespace: openebs
spec:
  version: 1.1.0  
  components:
  - name: LocalProvisioner
  - name: MayaAPIServer
  - name: ExternalCSIProvisioner
  - name: NDM
```

#### Operator - Sample 5
```yaml
kind: Operator
metadata:
  name: the-one-and-only
  namespace: openebs
spec:
  version: 1.1.0
  config:
    name: the-one-and-only
    namespace: openebs
```

#### Operator - Sample 6
```yaml
kind: Operator
metadata:
  name: the-one-and-only
  namespace: openebs
spec:
  version: 1.1.0  
  crds:
  - name: CRDAll
```

#### Operator - Sample 7
```yaml
kind: Operator
metadata:
  name: the-one-and-only
  namespace: openebs
spec:
  version: 1.1.0
  storageclasses:
  - name: StorageClassJiva
  - name: StorageClassCStor
  - name: StorageClassLocalPV
```

#### Operator - Sample 8
```yaml
kind: Operator
metadata:
  name: the-one-and-only
  namespace: openebs
spec:
  version: 1.1.0
  disable: true
```

#### OperatorConfig Custom Resource
OperatorConfig custom resource provides necessary configuration to reconcile openebs operator. This does not have a 
dedicated controller/watcher of its own.

#### OperatorConfig - Sample 1
```yaml
kind: OperatorConfig
metadata:
  name: the-one-and-only
  namespace: openebs
spec:
  # optional k8s jobs that are run before reconciling
  # the operator
  #
  # out-of-band pre tasks can be handled here
  #
  # pre tasks can be validation, patch, update, create, etc
  preTasks:

  # optional k8s jobs that are run after reconciling
  # the operator
  #
  # out-of-band post tasks can be handled here
  #
  # post tasks should be limited to validation or monitoring
  # related activities
  postTasks:

  # labelSelector is optional
  #
  # label selector selects various openebs components via labels
  # or label expressions
  labelSelector:

  # components are understood internally
  #
  # operator understands if a component is a deployment, or STS, 
  # or DaemonSet, or Job, etc.
  components:
```

#### OperatorConfig - Sample 2
```yaml
kind: OperatorConfig
metadata:
  name: the-one-and-only
  namespace: openebs
spec:
  components:
  - name: MayaAPIServer
    image: quay.io/openebs/m-apiserver
    imageTag: 1.0.0
    namespace: openebs
    labelSelector:
    sidecars:
    - name: maya-exporter
      image: quay.io/openebs/m-exporter
      imageTag: 1.0.0
  - name: NDM
    image: quay.io/openebs/ndm
    imageTag: 0.5.0
    namespace: openebs
    labelSelector:
```

#### OperatorConfig - Sample 3
```yaml
kind: OperatorConfig
metadata:
  name: the-one-and-only
  namespace: openebs
spec:
  preTasks:
  - name: DeleteDeprecatedCRDs
    image: quay.io/openebs/m-RemoveDeprecated
    imageTag: 0.0.1
    namespace: openebs
  postTasks:
  - name: IsOpenEBSRunning
    image: quay.io/openebs/m-IsRunning
    imageTag: 0.0.1
    namespace: openebs
  components:
  - name: MayaAPIServer
    image: quay.io/openebs/m-apiserver
    imageTag: 1.0.0
    namespace: openebs
    labelSelector:
    sidecars:
    - name: maya-exporter
      image: quay.io/openebs/m-exporter
      imageTag: 1.0.0
```

#### OperatorConfig - Sample 4
```yaml
kind: OperatorConfig
metadata:
  name: the-one-and-only
  namespace: openebs
spec:
  labelSelector:
    matchLabels: 
      openebs.io/is-operator-managed: true
    matchExpressions:
    - key: openebs.io/version
      operator: In
      values: ["0.8.0", "0.7.0"]
  components:
  - name: MayaAPIServer
    image: quay.io/openebs/m-apiserver
    imageTag: 1.0.0
    namespace: openebs
    labelSelector:
    sidecars:
    - name: maya-exporter
      image: quay.io/openebs/m-exporter
      imageTag: 1.0.0
```

#### ConfigInjector - Definition
```yaml
kind: ConfigInjector
metadata:
  name:
  namespace:
spec:
  policies:
  - name:   # name given to this injection
    select: # select kubernetes resource(s)
    apply:  # apply values against above selected kubernetes resource(s) as target(s)
```

#### ConfigInjector - Sample 1
- Update a container
  - Where the container name is `cstor-istgt`
  - Add or Update following env variables:
    - QueueDepth
    - Luworkers
```yaml
kind: ConfigInjector
metadata:
  name: 
  namespace: 
spec:
  policies:
  - name: inject-perf-tunables
    select: 
      kind: Deployment
      name: my-cstor-deploy
      namespace: openebs
    apply:
      containers:
      - name: cstor-istgt
        env:
        - name: QueueDepth
          value: 6
        - name: Luworkers
          value: 3
```

#### ConfigInjector - Sample 2
- Remove an annotation from custom resource if exists
- Where custom resource kind is `CStorVolumeReplica`
```yaml
kind: ConfigInjector
metadata:
  name: 
  namespace: 
spec:
  policies:
  - name: remove-ann
    select: 
      kind: CStorVolumeReplica
      labelSelector:
        matchLabels: 
          openebs.io/pv: pvc-abc-abc
      namespace: openebs
    remove:
      annotations:
      - openebs.io/vsm
```

#### ConfigInjector - Sample 3
- Remove an annotation from custom resource if exists
- Add / Update an annotation from custom resource
- Where custom resource kind is `CStorVolumeReplica`
```yaml
kind: ConfigInjector
metadata:
  name: 
  namespace: 
spec:
  policies:
  - name: update-annotations
    select: 
      kind: CStorVolumeReplica
      labelSelector:
        matchLabels: 
          openebs.io/pv: pvc-abc-abc
      namespace: openebs
    remove:
      annotations:
      - openebs.io/vsm
    apply:
      annotations:
        openebs.io/version: 1.0.0
```

#### References
- [OEP on CSI](https://github.com/openebs/openebs/pull/2617/)
