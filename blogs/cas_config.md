### Info
- Version: 4
- Last Updated On: 01 Mar 2019

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
- CAS Config controller **will not** execute below actions:
  - create a resource
  - delete a resource
- CAS Config controller will execute below actions only:
  - patch resources
  - rollback the patch in-case of any errors
    - e.g. patch with old config info
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
  // Scope indicates if the config is meant
  // to be used within the scope of
  // namespace or cluster
  Scope     Scope      `json:"scope"`
  
  // Policies consist of a list of various 
  // config options that should be patched
  // against matching resources
  Policies  []Policy   `json:"policies"`
}

// ScopeType defines the type of a scope
type ScopeType

const (
  // ClusterScope is used to grant the entire 
  // config to be applied against matching 
  // resources found in a cluster
  ClusterScope    ScopeType = "cluster"
  
  // NamespaceScope is used to grant the entire
  // config to be applied against matching
  // resources found in particular namespace or
  // namespaces
  NamespaceScope  ScopeType = "namespace"
)

// Scope provides permission options
// that if satisfied will enable this
// config to be applied against matching
// resources
type Scope struct {
  // Level flags the permission to execute
  // this config within a cluster or one 
  // or more namespaces
  Level   ScopeType   `json:"level"`
  
  // Values provide specific namespaces
  // that are permitted for this config
  // to operate against
  Values  []string    `json:"values"`
}

// Policy provides the actual configuration
// values along with matching resources.
// Matching resources here refer to those
// resources which are eligible to have
// this config operated against
type Policy struct {
  // Name of this policy
  //
  // NOTE:
  //  Future: It may be used as a reference in 
  // matching resources to grant or deny this 
  // particular policy or a set of policies 
  // versus allowing itself to be applied 
  // against the entire config
  Name      string      `json:"name"`
  
  // Values represent the actual configuration
  // that gets applied against the matching
  // resources
  Values    Values      `json:"values"`
  
  // Include provides a way to select resources
  // that match this include rules
  Include   Include     `json:"include"`
}

type Values struct {
  Labels        map[string]string `json:"labels"`
  Annotations   map[string]string `json:"annotations"`
  Environments  map[string]string `json:"envs"`
  Containers    []Container       `json:"containers"`
  Sidecars      []Container       `json:"sidecars"`
  MainContainer Container         `json:"mainContainer"`
  From          []FromOperation   `json:"from"`
}

type FromOperation struct {
  Kind      string    `json:"kind"`
  Name      string    `json:"name"`
  Path      string    `json:"path"`
  AsType    string    `json:"asType"`
  AsKey     string    `json:"asKey"`
}

type Include struct {
  Labels          map[string]string   `json:"labels"`
  Annotations     map[string]string   `json:"annotations"`
  Namespaces      []string            `json:"namespaces"`
  Kinds           []string            `json:"kinds"`
  Names           []string            `json:"names"`
  ServiceAccounts []string            `json:"serviceAccounts"`
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
      labels:
      annotations:
      namespaces:
      kinds:
      names:
      serviceAccounts:
```
