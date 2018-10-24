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
  - feature: CreateJivaVolume, NVMEIsAvailable, HostHasISCSIInstalled

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
