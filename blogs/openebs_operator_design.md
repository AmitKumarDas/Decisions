
### WIP Draft -- OpenEBS Operator Specification
```yaml
kind: OpenEBS
metadata:
  name: default
# specification of OpenEBS Operator
# specifications related to following managed via reconciliation
# - Install, Update, Upgrade, Probe / Monitor, Self Healing
# query: Should above be the job of openebs operator?
# think: Does it help the real human i.e. operator?
spec:
  # maya is the specification of control plane
  maya:
    # what version of maya should be enabled/made available
    version: 0.7.0
    ha:
      # what kind of High Availability support(s)
      support:
      - NodeFailure
      - NetworkFailure
    # what storage engine(s) should be enabled
    # query: should this be casEngine? casType?
    cas:
    - cstor
    - jiva
    # what disk management solutions should be enabled
    disk:
    - ndm
  # cstor is the specification related to cstor storage engine
  cstor:
    onSparse: true
  # jiva is the specification related jiva storage engine
  jiva:
  # localPV is the specification related to local PV storage engine
  localPV:
  # ndm is the specification related to ndm disk management solution
  ndm:
```

### Rough Drafts
#### Draft 1 -- dated 24/Oct/2018

- Review Comment - This is not a specification i.e. statement
- Review Comment - This is more like a program/logic
- Review Comment - This is not what enduser will really want as a specification

```yaml
# desired state
spec:
  # refers to openebs components that should be available & in the specified state
  component:

  # simple desired state specification
  - include: currentCASTemplate; currentRunTask; currentOpenEBSProvisioner; currentMAPIServerEnv
  - exclude: labelSelector app eq jiva version lt current

  # finely granular desired state specification
  - exclude: selector kind eq runtask
    if: selector name contains -0.6.0
  - include: selector kind eq castemplate
    replace: metadata.labels.abc
    value: |
    if: labelSelector version eq 0.6.0; pathSelector metadata.labels.abc ne default

  # refers to features supported by openebs and should work as specified here
  feaure:
  - support: CreateJivaVolume; NVMEIsAvailable; HostHasISCSIInstalled
  - support: CreateCstorVolume
    if: selector openebsVersion gt 0.7.0
  - support: VolumeOnSparseFiles
    if: labelSelector stage in test,qa
  - support: CreateCstorVolume
    check: latency lt 1msec
    with: replicas 3
    interval: 1day or 1 hr # check support every 1 day on success or every 1 hr on failure

status:
  state: failed
  reason:
  starttime:
  endtime:
  feature:
  - support: CreateJivaVolume; NVMEIsAvailable
    state: passed
  - support: CreateCstorVolume
    if: selector version gt 0.7.0
    state: passed
  - support: CreateCstorVolume
    check: latency lt 1msec
    with: replicas 3
    interval: 3days
    state: passed
    starttime:
    endtime:
    prevState: failed
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

##### Appendix
- currentCASTemplate - ensure CAS Templates for current version of openebs are available in the cluster
- currentRunTask - ensure Run Tasks for current version of openebs are available in the cluster
- currentOpenEBSProvisioner - ensure openebs provisioner for current version of openebs is available in the cluster
- currentMAPIServerEnv - ensure maya api server pod has environment variables for current version of openebs
