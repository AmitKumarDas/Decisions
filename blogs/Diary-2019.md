### 23 Apr 2019
#### Maya Query Language - MayaQL
- A declarative workflow
- Understands kubernetes
- A k8s controller/watcher
- MayaQL & MayaQLResult
- Capabilities
  - can get one or more properties from a resource
  - can get one or more properties from a list of resources
  - can create a k8s resource
  - can delete a k8s resource
  - can build a k8s resource before creating it
  - can update a k8s resource
  - can pod exec a k8s pod
  - can grab log tails from a container
  - can display result in any desired format

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

### Detailed References:
- https://github.com/AmitKumarDas/Decisions/blob/master/blogs/build_as_you_go.md
