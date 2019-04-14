#### Prototype 1
```yaml
kind: JivaVolume
metadata:
  namespace:
  name:
  annotations:
  labels:
spec:
  capacity:
  casType:
  fsType:
  lun:
iscsi:
  iqn:
  targetPortal:
casConfig:
service:
target:
replica:
references:
  persistentVolume:
  persistenVolumeClaim:
  storageClass:
status:
```

#### Prototype 2
- casConfig is wrong
  - since cas config is not strictly typed
  - since cas config is tightly coupled with go templating
- config is replacement for casConfig
  - config will represent most of common kubernetes types
  - e.g. envs, labels, annotations, pod affinity, node affinity, etc
- references are bad
  - since spec should be the desired state on which reconcile should work
  - spec should lead to owned resources/objects
- spec should have simple annotations
  - complex annotations typically indicate code smell
  - annotations should be used to refer to external / in-direct resources
```yaml
kind: JivaVolume
metadata:
  namespace: # openebs system namespace
  name: pv-name # can be pv name or anything
  annotations:
  labels:
    openebs.io/match-name: # can be pv name if metadata.name != pv name
spec:
  capacity:
  casType:
  fsType:
  lun:
iscsi:
  iqn:
  targetPortal:
config:
  envs:
  labels:
    openebs.io/match-name: # cann be pv name if metadata.name != pv name
  annotations:
service:
  config:
  name: pv-name # can be pv name or anything 
  generateName: pv-name # can be pv name or anything
target:
  config:
  name: pv-name # can be pv name or anything;
  generateName: pv-name # can be pv name or anything;
replica:
  config:
  name: pv-name # can be pv name or anything;
  generateName: pv-name # can be pv name or anything
status:
```
