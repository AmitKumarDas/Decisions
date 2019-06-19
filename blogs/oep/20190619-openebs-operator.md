### Motivation
OpenEBS has various components that need to be in running state to let its volumes be consumed by higher applications. As the
number of components grow, have inter-dependency, need upgrades it becomes un-manageable to let this be managed by a human.
OpenEBS operator should ideally take over most of the manual operations needed to let OpenEBS run smoothly.

### Points to Consider
- Number of openebs components can vary
- Some of the openebs components can be optional
- Each component might need a special handling

### Possible Custom Resources
- Operator
- ReleaseConfig
- ConfigInjector

### Sample Schemas

#### Operator - Sample 1
```yaml
kind: Operator
metadata:
  name: TheOne
  namespace: openebs
spec:
  # installs or upgrades to OpenEBS 1.1.0
  version: 1.1.0
```

#### Operator - Sample 2
```yaml
kind: Operator
metadata:
  name: TheOne
  namespace: openebs
spec:
  # installs or upgrades to OpenEBS 1.1.0
  version: 1.1.0
  components:
  # installs local provisioner only
  - name: local-provisioner
```

#### Operator - Sample 3
```yaml
kind: Operator
metadata:
  name: TheOne
  namespace: openebs
spec:
  # installs or upgrades to OpenEBS 1.1.0
  version: 1.1.0
  components:
  # installs specified components
  - name: local-provisioner
  - name: maya-api-server
  - name: external-provisioner
  - name: ndm
```

#### ReleaseConfig - Sample 1
```yaml
```

#### ConfigInjector - Sample 1
```yaml
```
