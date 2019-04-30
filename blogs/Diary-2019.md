### 30 Apr 2019
#### Group & Version
- openebs.io:v1alpha1, openebs.io:v1alpha2
- pool.openebs.io:v1alpha1, pool.openebs.io:v1beta1
- task.openebs.io:v1alpha1

### 29 Apr 2019
#### MayaQL to RunTask
- Keep the legacy of RunTask
- Also thought of MayaActions
  - However, how to manage the plural
- Probably `RunTask` makes use of `Action`
```yaml
kind: RunTask
spec:
  actions:
  - run: GetPodName
    desc: It gets the name of a Pod
```

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
- Catalog can be renamed to Component
- A component represents a set of resources
- Component will be built programatically
- No YAML
- Code as be structured as follows:
  - pkg/component/v1alpha1/
  - pkg/component/maya-apiserver/v1alpha1/
  - pkg/component/openebs-provisioner/v1alpha1/
- BDD can make use of above if & when required to install
- Openebs operator can make use of above as well

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
