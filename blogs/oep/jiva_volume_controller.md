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
- A Resize workflow has following events:
  - JivaVolume resource should enable resize via `resize.jiva.openebs.io/enabled: true` annotation
  - JivaVolume resource should be annotated with `resize.jiva.openebs.io/name: <name-of-resize-resource>`
  - A JivaVolumeResize resource is created
- Once resize operation is completed, resize controller updates JivaVolumeResize's status
- JivaVolume controller keeps a track of JivaVolumeResize if JivaVolume has resize annotation
- If a resize operation is completed successfully, JivaVolume controller performs the following:
  - Updates the capacity field
- Resize workflow then performs following actions:
  - Removes `resize.jiva.openebs.io/name` annotation
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
```
