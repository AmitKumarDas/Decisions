### Info
- Version: 1
- Last Updated On: 23 Feb 2019

### Motivation
Desire to apply, inject, merge configuration against components or services in a kubernetes cluster.

### High Level Design
- CASConfig is a kubernetes _custom resource_
- It is optionally controlled by a controller
  - Name of this controller is "cas-config controller"
- It can be embedded inside other resources, e.g.:
  - OpenebsCluster
  - ComponentSet
  - CstorVolumeSet
  - CstorPoolSet
  - JivaVolumeSet
- Its targets can be specified via various select rules
- It can patch, merge, add following configurations:
  - labels
  - annotations
  - envs
  - containers
  - sidecars
  - maincontainer
  - tolerations
  - etc

### Specifications of CASConfig
```yaml
kind: CASConfig
spec:
  scope:
  config:
  - policy:
    values:
      labels:
      annotations:
      envs:
      containers:
      sidecars:
      maincontainer:
    selector:
      byLabel:
      byAnnotation:
      byNamespace:
      byKind:
      byName:
```
