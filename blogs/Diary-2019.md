### 23 Apr 2019
#### component
- catalog can be renamed to component
- a component represents a set of resources
- component will be built programatically
- No YAML
- pkg/component/v1alpha1/
- pkg/component/maya-apiserver/v1alpha1/
- pkg/component/openebs-provisioner/v1alpha1/
- BDD can make use of above
- openebs operator can make use of above

#### cas config
- casconfig represent the configurations that can be injected
- a component can have various capabilities over a cas config
- a component e.g. maya-apiserver has the capability to set image
- similarly cstorvolume does not understand image as a capability
- each component will support their own capabilities
- a component can be modified based on the provided:
  - config & 
  - supported capabilities
- A config may or may not belong to a component

```go
// in pkg/component/v1alpha1/interface.go
type Component interface {
  Name() string
  Create() error
  Delete() error
  Update() error
}
```

```go
// in pkg/component/v1alpha1/component.go
type Capability func(*Config) error
type CapabilityList []Capability
type IsConfig func(*Config) bool
var ComponentRegistrar map[string]Component
var CapabilityListRegistrar map[string]CapabilityList
```
