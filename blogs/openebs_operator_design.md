
### OpenEBS Operator Specification -- WIP
- This is the typical default operator specification that gets generated and applied before the operator logic kicks in

```yaml
kind: OpenEBS
metadata:
  name: default
# specification of OpenEBS Operator
# specifications related to following managed via reconciliation
# - Install, Update, Upgrade, Probe / Monitor, Self Healing
# DOUBT: Should above be the job of openebs operator?
# THINK: Does it help the real human i.e. operator?
spec:
  # A value of false will uninstall all openebs components
  enabled: true
  # what version of openebs should be enabled/made available
  # version controls install, upgrade & rollback as well
  version: 0.7.0
  # volumeProvisioner is the specification of volume provisioner
  volumeProvisioner:
    # A value of false will un-install all instances of volume provisioner
    enabled: true
    # cluster provides a list of volume provisioners that interact with openebs
    # to provision volume
    cluster: 
    - type: remote
      name: remote-123
      id: 123
    - type: remote
      name: remote-231
      id: 231
    - type: local
      name: local-001
      nodeSelector:
      ha:
        support:
        - NodeFailure
        - NetworkFailure
  # apiServer is the specification of openebs api server
  apiServer:
    # A value of false will un-install all instances of api server
    enabled: true
    # nodeSelector translates to kubernetes pod's nodeSelector property
    nodeSelector:
    # desired to have ha support
    # absence of this property will imply no desire and hence no checks on the system
    ha:
      # What kind of High Availability support(s) are expected
      # All is a valid value
      # None is a valid value
      support:
      - NodeFailure
      - NetworkFailure
  # cstor is the specification related to cstor storage engine
  cstor:
    # A value of false will un-install all cstor options, templates, tunables
    enabled: true
    sparse:
      enabled: true
      # DOUBT: Can there be different flavors of sparse files ?
      # DOUBT: Should sparse be directly under spec ?
      # DOUBT: Should sparse be present at all ?
      support:
      - xyz
    snapshot:
      # A value of false will un-install all cstor snapshot related options, templates, tunables
      enabled: true
      # operator can provide the name of custom templates to invoke CRUD operation(s) w.r.t snapshot
      # default value implies operator will decide / set the name of the template
      template:
        create: default
        delete: default
        list: default
    volume:
      # A value of false will un-install all cstor volume related options, templates, tunables
      enabled: true
      # operator can provide the name of custom templates to invoke CRUD operation(s) w.r.t volume
      # default value implies operator will decide / set the name of the template
      template:
        create: default
        delete: default
        list: default
        read: default
  # jiva is the specification related jiva storage engine
  jiva:
    # A value of false will un-install all jiva related options, templates, tunables
    enabled: true
    volume:
      # A value of false will un-install all jiva volume related options, templates, tunables
      enabled: true
      template:
        create: default
        delete: default
        list: default
        read: default
  # localPV is the specification related to local PV storage engine
  localPV:
  # ndm is the specification related to ndm disk management solution
  ndm:
```

### Old Drafts
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
