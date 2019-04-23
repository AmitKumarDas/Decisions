### UseCase 1
- I want to build a **Component Registration** feature
- Individual components should be able to register themselves
- Each component can register their capabilities
- Registrar is not tightly coupled with every component
- Each component implementation adheres to contracts exposed by Component interface

```go
// in github.com/openebs/maya/pkg/component/v1alpha1/interface.go
type Component interface {
  Name() string
  Create() error
  Delete() error
  Update() error
}
```

```go
// in github.com/openebs/maya/pkg/component/v1alpha1/component.go

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

func (b *RegistrarBuilder) WithComponent(obj Component) *RegistrarBuilder {
  if obj == nil {
    b.errs = append(b.errs, errors.New("failed to register: nil component instance"))
    return b
  }
  b.object = obj
  return b
}

func (b *RegistrarBuilder) WithCapabilities(c CapabilityList) *RegistrarBuilder {
  b.capabilities = c
  return b
}

func (b *RegistrarBuilder) Register() error {
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

```go
// in pkg/component/maya-apiserver/v1alpha1/component.go

import (
  component "github.com/openebs/maya/pkg/component/v1alpha1"
)

var _ = component.RegistrarBuilder().
  WithName("maya-apiserver").
  WithComponent(&MayaAPIServer{}).
  Register()

type MayaAPIServer struct {}
```
