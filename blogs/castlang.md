## CASTLang

This is the next version of CASTemplate. This should provide all the benefits that CASTemplate provides as well make scripting
easier.

### CASTemplate - Good Parts - 1
- variable declaration & definition via go template functions before the task yaml
- task yaml makes use of these variables to be finally executed via go-templating
- Note - This is about variable definition, transformation of value, & templating
- NOTE - This is not about Kubernetes API invocation

#### Low Level Design - Good Parts - 1
- Variable declaration, definition, transformation, autosave
- `spec.let` & `spec.template` dictionaries will be stored at RunTask runner
- `run` executes the actions
- Saving the values need not always be nested inside map[string]interface{}
  - It can be keyed at dotted keys e.g. `nameParams := map[string]string{"pv.name": pvName}`
  - However, need to choose only one way to save them

```yaml
kind: RunTask
spec:
  let:
  template:
run:
output:
```

```yaml
kind: RunTask
spec:
  let:
  - myPod: |
      kind: Pod
      apiVersion: v1
      metadata:
        name: Hulla
  - isSvc: "true"
```

```yaml
kind: RunTask
spec:
  template:
  - pod: |
      {{- $name := .spec.let.name | default "cool" -}}
      {{- $ns := .spec.let.namespace | default "default" -}}
      kind: Pod
      apiVersion: v1
      metadata:
        name: $name
        namespace: $ns
  - val: |
      {{- $name := .spec.let.name | default "cool" -}}
      {{ $name | suffix "-dude" }}
```

```yaml
kind: RunTask
run:
  - id: 101
    action: list
    kind: Pod
    labelSelector: app=jiva
  - id: 102
    action: create
    podContent: ${spec.template.pod}
  - id: 103
    action: create
    podContent: ${spec.let.myPod}
```

```go
type RunTask struct {
  ObjectMeta
  Spec RunTaskSpec `json:"spec"`
}

type RunTaskSpec struct {
  Let      map[string]string `json:"let"` // dict of variable with its direct value
  Template map[string]string `json:"template"` // dict of variable with its templated value
}
```
