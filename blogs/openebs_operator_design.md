### Install, Update, Upgrade, Probe, Checks, Healing via reconciliation
Is this and more openebs operator?

### Usecase #1 If any desired resource is deleted operator should re-install

- How will it know if a particular resource is present or absent ?
- Should every resource be tracked via name & version ? 
- Should every resource be tracked via counts ? 
- Should a resource be tracked based on tags ? 
- Should a resource be tracked based on any of above ?

### Usecase #2 If any resource is excluded operator should delete if resource is already installed ?

To Be Filled

### More Usecases 

### Specification
```yaml
# desired state
spec:
  component:

  # simple desired state specification
  - include: currentCASTemplate, currentRunTask, currentOpenEBSProvisioner
  - exclude: labelSelector app eq jiva version lt current

  # finely granular desired state specification
  - exclude: selector kind eq runtask
    if: selector name contains -0.6.0
  - include: selector kind eq castemplate
    replace: metadata.labels.abc
    value: |
    if: labelSelector version eq 0.6.0, pathSelector metadata.labels.abc ne default

  supports:
  - feature: CreateJivaVolume, NVMEIsAvailable

status:
  state: failed
  reason:
  starttime:
  endtime:
  supports:
  - feature: CreateJivaVolume, NVMEIsAvailable
    state: passed
  component:
  - include: currentCASTemplate, currentRunTask, currentOpenEBSProvisioner
    state: passed
  - exclude: labelSelector app eq jiva version lt current
    state: passed
  - exclude: selector kind eq runtask
    if: selector name contains -0.6.0
    state: passed
  - include: selector kind eq castemplate
    replace: metadata.labels.abc
    value: |
    if: labelSelector version eq 0.6.0, pathSelector metadata.labels.abc ne default
    state: failed
    reason: blah blah blah
    starttime:
    endtime:
    count:
```
