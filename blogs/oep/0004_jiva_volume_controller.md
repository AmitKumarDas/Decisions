### JivaVolume as a Custom Resource
- A Jiva volume is represented as a custom resource
- A JivaVolume resource is watched by a controller named `jiva-volume-controller`

```yaml
kind: JivaVolume
metadata:
  name: pv-name
  annotations:
  labels:
spec:
  capacity:
  lun:
iscsi:
  iqn:
  targetPortal:
service:
  name: pv-name-svc
  namespace:
targetDeployment:
  name: pv-name-ctrl
  namespace:
replicaDeployment:
  name: pv-name-rep
  namespace:
status:
```

### JivaVolumeResize as a Custom Resource
- JivaVolumeResize represents an operation as well as a custom resource
- JivaVolumeResize resource is watched by a controller named `jiva-volume-resize-controller`
- Higher order resize workflow is as follows:
  - Verify if JivaVolume resource supports resize via `resize.jiva.openebs.io/enabled: true` annotation
  - Creates a JivaVolumeResize resource
- JivaVolumeResize controller workflow is as follows:
  - Tries to annotate JivaVolume resource with `resize.jiva.openebs.io/name: <name-of-resize-resource>`
  - Above might fail if an annotation exists already
  - It needs to wait to apply the annotation
  - Updates JivaVolumeResize's status once the resize is completed
  - Waits for JivaVolume's spec.capacity to get updated
  - Removes annotation once JivaVolume's spec.capacity gets updated
- JivaVolume controller has following workflow w.r.t resize
  - Verifies if JivaVolumeResize status is SUCCESS
  - Updates spec.capacity
- Higher order resize workflow then performs following actions:
  - Deletes JivaVolumeResize resource

```yaml
kind: JivaVolume
metadata:
  name: my-jiva
  annotations:
    resize.jiva.openebs.io/enabled: true
    resize.jiva.openebs.io/name: resize-6-to-12
spec:
  capacity: 6
status:
```

```yaml
kind: JivaVolumeResize
metadata:
  name: resize-6-to-12
  labels:
    jiva.openebs.io/name: my-jiva
spec:
  capacity: 12
status:
  phase: # INIT, IN-PROGRESS, FAILED, SUCCEEDED, COMPLETED
```

### Why capacity expansion requires a Custom Resource?
- One might imagine changing `spec.capacity` from existing value to a new value
- IMO a dedicated logic specializing capacity expansion is better
- This removes the need to reset to old value when capacity expansion fails / not authorized
- There may be cases when capacity expansion needs to wait for some other operation to complete
  - e.g. resize should not happen during a snapshot
  - e.g. resize should not happen if namespace does not have desired capacity
- The proposal here is to avoid modifying `spec.capacity` directly
  - Instead, annotate the resource & delegate the logic to a separate controller
  - In other words, we want the controller to update `spec.capacity` field via annotation(s)
