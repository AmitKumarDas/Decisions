## CASTLang

CASTLang is the _code name_ for next version of CASTemplate. This should provide all the benefits that CASTemplate provides 
as well help building **kubernetes workflows easier**. This version concentrates on RunTask. A RunTask can get executed via 
CASTemplate runner or via a new kubernetes controller. Care has been taken **NOT to modify CASTemplate** as it is user 
facing and hence is prone to breaking its users. This attempt to improvise CASTemplate code named **CASTLang** tries to make 
RunTask as independent and devops friendly than its earlier version.

#### Low Level Design
- The design boils down to:
  - Variable Declaration, Definition, Transformation, Actions & AutoSave
- `spec.let` & `spec.template` dictionaries will be stored at RunTask runner
- `spec.let` lets us declare variables with corresponding values
- `spec.template` lets us declare variables with corresponding values as go based templates
- `run` executes the actions

```yaml
kind: RunTask
spec:
  let:
  template:
  run:
output:
status:
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
spec:
  run:
    - id: 101
      action: list
      kind: PodList
      labelSelector: app=jiva
    - id: 102
      action: create
      kind: Pod
      content: ${@.spec.template.pod}
    - id: 103
      action: create
      kind: Pod
      content: ${@.spec.let.myPod}
```

```go
type RunTask struct {
  ObjectMeta
  Spec   RunTaskSpec `json:"spec"`
  Status RunTaskStatus `json:"status"`
  Output String `json:"output"` // template that is used to build the output of this task
}

type RunTaskSpec struct {
  Let      map[string]string `json:"let"` // dict of variable with its direct value
  Template map[string]string `json:"template"` // dict of variable with its templated value
  Run      []Runnable        `json:"run"` // list of runnables that get executed
}
```

### SpinOffs
- One can extend RunTask to meet their specific requirement
- I shall explain how `TestTask` extends from RunTask
- Most of stuff remains same baring `expect` which gets introduced as a new field
- Assumptions:
  - Depends on consuming Kubernetes APIs
  - Based on workflow or pipelining of tasks

```yaml
kind: TestTask
spec:
  let:
  template:
  run:
expect:
status:
```

```yaml
kind: TestTask
spec:
  let:
  template:
  run:
    - id: 101
      action: list
      kind: PodList
      labelSelector: app=jiva,org=openebs
expect:
  - pod: ${@.spec.run.101}
    match: 
      - status == Online
      - kind == PodList
      - namespace In default,openebs
      - labels == app=jiva,org=openebs
status:
```

```go
type TestTask struct {
  RunTask
  Expect []TestExpect `json:"expect"` // matchers used when running this task to build testing logic
}
```
