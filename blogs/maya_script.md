### Motivation
What if there is an easy way to code adhoc stuff or automations with respect to 
Kubernetes whithout messing with shell scripts.

### High Level Design - _MayaScript_
- It is a kubernetes custom resource
- It can be optionally managed by a controller
  - Name of its optional controller is "maya-script controller"
- It understands invoking kubernetes API calls
- It understands invoking openebs/maya functions
- It executes itself as a go template
- It can save the success result of each line of script code
- It can save the error result (if any) of the line of script code
- It can save the time taken to execute each line of script code
- It exits execution if any line of script code results in an error
- It can accept one or more yaml document(s)
- It should execute based on the order provided in the specifications

### _Specification_ of Maya Script
```go
type MayaScript struct {
  metav1.TyeMeta
  metav1.ObjectMeta
  Spec   MayaScriptSpec
  Status MayaScriptStatus
}

type MayaScriptSpec struct {
  Script map[string]string `json:"script"`
}

type MayaScriptStatus struct {}
```
```yaml
kind: MayaScript
spec:
  script:
    desc: I would like to manage my kubernetes scripts in a better way
    given: a kubernetes cluster with 3 nodes
    withSCYaml: |
    withClusterRoleYaml: |
    withAPIServerYaml: |
    withNDMYaml: |
    withProvisionerYaml: |
    it: should install openebs without problems
    installSC: $(k8sApply .yaml.withSC)
    installClusterRole: $(k8sApply .yaml.withClusterRole)
    installMayaAPIServer: $(k8sApply .yaml.withAPIServer)
    installNDM: $(k8sApply .yaml.withNDM)
    installProvisioner: $(k8sApply .yaml.withProvisioner)
    mayaOption: $($opts := k8sKind "deploy" | list k8sName "maya-apiserver" | list k8sPath ".spec.replicas")
    verifyMayaWith3Replicas: $(k8sGetNumValue $opts | eq 3)
```

```yaml
kind: MayaScript
spec:
  status:
    output:
      desc: I would like to manage my kubernetes scripts in a better way
      given: a kubernetes cluster with 3 nodes
      withSCYaml: ...
      withClusterRoleYaml: ...
      withAPIServerYaml: ...
      withNDMYaml: ...
      withProvisionerYaml: ...
      it: should install openebs without problems
      installSC: ok, 7 secs
      installClusterRole: ok, 5 secs
      installMayaAPIServer: ok, 2 secs
      installNDM: ok, 10 secs
      installProvisioner: ok, 5 secs
      mayaOption: ok, 1 secs
      verifyMayaWith3Replicas: ok, 1 secs
```
