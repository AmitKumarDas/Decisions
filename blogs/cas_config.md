### Info
- Version: 3
- Last Updated On: 28 Feb 2019

### Motivation
Desire to apply, inject, merge configuration against targeted resources in a kubernetes cluster. We shall also attempt to run commands against the targeted resources if it fits into the design of cas config.

### High Level Design
- CASConfig is a kubernetes _custom resource_
- It can be **optionally** controlled by its own controller
  - Name of its controller is called as "cas-config controller"
- It can be embedded inside other resources, e.g.:
  - OpenebsCluster
  - Upgrade (_Near Future_)
  - CstorVolume (_Future_)
  - CstorPoolClaim (_Future_)
  - JivaVolume (_Future_)
- The embedding object should embed only a single config specification
- A config can have one or more targets
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
- CAS Config will not execute below actions:
  - create a resource
  - delete a resource
- CAS Config will execute below action only:
  - patch resources
  - get resources
  - list resources

### Specifications of CASConfig
```go
type CASConfig struct {
  metav1.TypeMeta
  metav1.ObjectMeta
  
  Spec    Spec    `json:"spec"`
  Status  Status  `json:"status"`
}

type Spec struct {
  Scope     Scope      `json:"scope"`
  Policies  []Policy   `json:"policies"`
}

type ScopeType

const (
  ClusterScope    ScopeType = "cluster"
  NamespaceScope  ScopeType = "namespace"
)

type Scope struct {
  Level   ScopeType   `json:"level"`
  Values  []string    `json:"values"`
}

type Policy struct {
  Name      string      `json:"name"`
  Values    Values      `json:"values"`
  Include   Include     `json:"include"`
}

type Values struct {
  Labels        map[string]string `json:"labels"`
  Annotations   map[string]string `json:"annotations"`
  Environments  map[string]string `json:"envs"`
  Containers    []Container       `json:"containers"`
  Sidecars      []Container       `json:"sidecars"`
}

type ValueOps struct {
  LabelOps    []LabelOperation  `json:"labels"`
}

type LabelOperation struct {
  Kind      string    `json:"kind"`
  Name      string    `json:"name"`
  Selector  Selector  `json:"select"`
  Path      string    `json:"path"`
}

type Selector struct {
  ByLabels        map[string]string `json:"byLabels"`
  ByLabelops      []LabelOperation  `json:"byLabelOps"`
  ByAnnotations   map[string]string `json:"byAnnotations"`
  ByNamespace     string            `json:"byNamespace"`
  ByKind          string            `json:"byKind"`
  ByName          string            `json:"byName"`
}
```

```yaml
kind: CASConfig
spec:
  scope:
    level: # namespace or cluster
    values: # names of namespaces
  policies: # policies will be applied in this order
  - name:
    values:
      labels:
      annotations:
      envs:
      containers:
      sidecars:
      maincontainer:
      from:
        kind:
        name:
        path:
        asType: # optional
        asKey: # optional
    include:
      labels: # map[string]string
      annotations: # map[string]string
      namespaces: # []string
      kinds: # []string
      names: # []string
      serviceAccounts: # []string
```
