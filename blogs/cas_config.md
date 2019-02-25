### Info
- Version: 2
- Last Updated On: 25 Feb 2019

### Motivation
Desire to apply, inject, merge configuration against components or services in a kubernetes cluster.

### High Level Design
- CASConfig is a kubernetes _custom resource_
- It is optionally controlled by a controller
  - Name of this controller is "cas-config controller"
- It can be embedded inside other resources, e.g.:
  - OpenebsCluster
  - CstorVolume
  - CstorPoolClaim
  - JivaVolume
- Embedding object should embed only a single config specifications
- A config can have targets
  - A config target is selected via `spec.select` option
  - A config target is one on which config is applied i.e. injected
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
```go
type CASConfig struct {
  metav1.TypeMeta
  metav1.ObjectMeta
  
  Spec    Spec    `json:"spec"`
  Status  Status  `json:"status"`
}

type Spec struct {
  Scope   Scope     `json:"scope"`
  Groups  []Group   `json:"groups"`
}

type ScopeType

const (
  ClusterScope    ScopeType = "cluster"
  NamespaceScope  ScopeType = "namespace"
)

type Scope struct {
  Level   ScopeType   `json:"level"`
}

type Group struct {
  Name      string      `json:"name"`
  Values    Values        `json:"data"`
  Selector  Selector    `json:"selector"`
}

type Data struct {
  Labels      map[string]string `json:"labels"`
  Annotations map[string]string `json:"annotations"`
}
```

```yaml
kind: CASConfig
spec:
  scope:
    level:
    values:
  group:
  - name:
    values:
      labels:
      annotations:
      envs:
      containers:
      sidecars:
      maincontainer:
    valueOps:
    - from:
        kind:
        name:
        select:
        path:
    select:
      byLabels:
      byLabelOps:
      byAnnotation:
      byNamespace:
      byKind:
      byName:
```
