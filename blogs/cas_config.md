### Info
- Version: 4
- Last Updated On: 01 Mar 2019

### Motivation
Desire to apply, inject, merge configuration against targeted resources in a kubernetes cluster. This started off as an attempt to manage config during the design of OpenEBS operator. However, its utility is beyond Openebs Operator and can be
used as-is with other custom resources as well. 

_NOTE: It may also attempt to run commands against the targeted resources (e.g. `kubectl exec`) to set configuration. However, it **remains to be discussed** if this fits into the design of cas config._

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
  //
  // Level defaults to namespace scope. In 
  // other words this config can be applied
  // against any resource found in the
  // namespace of this config
  Level   ScopeType   `json:"level"`
  
  // Namespaces provide specific namespaces
  // that are permitted for this config
  // to operate against
  Namespaces  []string    `json:"namespaces"`
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

// Values represent configuration that gets
// applied against matching resources
//
// NOTE:
//  Values can be specified directly or can
// be extracted from other resources
type Values struct {
  // Labels to be applied against the matching
  // resources
  Labels        map[string]string `json:"labels"`
  
  // Annotations to be applied against the
  // matching resources
  Annotations   map[string]string `json:"annotations"`
  
  // Environments to be applied against the
  // matching resources
  Environments  map[string]string `json:"envs"`
  
  // Containers to be set against the matching
  // resources
  //
  // NOTE:
  //  This will replace any existing containers
  // of the matching resources
  Containers    []Container       `json:"containers"`
  
  // Sidecar containers to be set against the
  // matching resources
  //
  // NOTE:
  //  This will add/merge to the existing containers
  // of the matching resources
  Sidecars      []Container       `json:"sidecars"`
  
  // MainContainer will replace the main container
  // of the matching resources
  MainContainer Container         `json:"mainContainer"`
  
  // From provides a list of operations that leads
  // to extraction of configuration from external
  // resources and subsequently including these
  // extracted values into this config
  From          []FromOperation   `json:"from"`
}

// FromOperation provides the necessary properties
// required to extract value from an external
// resource and saving this extracted value
// in this config
type FromOperation struct {
  Kind      string    `json:"kind"`
  Name      string    `json:"name"`
  Path      string    `json:"path"`
  ToType    string    `json:"toType"`
  ToKey     string    `json:"toKey"`
}

// Include provides the matchers used
// to select resources
type Include struct {
  // Labels if found in resource(s) make
  // these resources eligible to be
  // applied with this config
  Labels          map[string]string   `json:"labels"`
  
  // Annotations if found in resource(s)
  // make these resources eligible to be
  // applied with this config
  Annotations     map[string]string   `json:"annotations"`
  
  // Namespaces of resources that are
  // eligible to be applied with this
  // config
  Namespaces      []string            `json:"namespaces"`
  
  // Kinds of resources that are eligible
  // to be applied with this config
  Kinds           []string            `json:"kinds"`
  
  // Names of resources that are eligible
  // to be applied with this config
  Names           []string            `json:"names"`
  
  // ServiceAccounts of resources that are eligible
  // to be applied with this config
  ServiceAccounts []string            `json:"serviceAccounts"`
}
```

```yaml
kind: CASConfig
spec:
  scope:
    level:
    namespaces:
  policies:
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
        toType:
        toKey:
    include:
      labels:
      annotations:
      namespaces:
      kinds:
      names:
      serviceAccounts:
```
