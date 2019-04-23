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

var componentRegistrar map[string]Component
var capabilityListRegistrar map[string]CapabilityList

type RegistrarBuilder struct {
  name string
  object Component
  capabilities CapabilityList
  errs []error
}

func NewRegistrarBuilder() *RegistrarBuilder {
  return &RegistrarBuilder{}
}

func (b *RegistrarBuilder) WithName(name string) *RegistrarBuilder {
  if name == "" {
    b.errs = append(b.errs, errors.New("failed to register: missing component name"))
    return b
  }
  b.name = name
  return b
}

func (b *RegistrarBuilder) WithObject(obj Component) *RegistrarBuilder {
  if obj == nil {
    b.errs = append(b.errs, errors.New("failed to register: nil component instance"))
    return b
  }
  b.object = obj
  return b
}

func (b *RegistrarBuilder) WithCapabilities(c CapabilityList) *RegistrarBuilder {
  if c == nil {
    b.errs = append(b.errs, errors.New("failed to register: nil component capabilities"))
    return b
  }
  b.capabilities = c
  return b
}

func (b *RegistrarBuilder) Build() error {
  if len(b.errs) > 0 {
    return errors.Errorf("%v", b.errs)
  }

  if componentRegistrar[b.name] != nil {
    return errors.Errorf("failed to register: component '%s' already registered", b.name)
  }

  componentRegistrar[b.name] = b.object
  capabilityListRegistrar[b.name] = b.capabilities

  return nil
}
```
